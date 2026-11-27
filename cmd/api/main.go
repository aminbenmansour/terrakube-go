// Command api serves the Terrakube JSON:API, the Terraform TFE-compatible
// endpoints and the state backend.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aminbenmansour/terrakube-go/internal/api"
	internalauth "github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/clilogin"
	"github.com/aminbenmansour/terrakube-go/internal/config"
	"github.com/aminbenmansour/terrakube-go/internal/db"
	graphqlapi "github.com/aminbenmansour/terrakube-go/internal/graphql"
	"github.com/aminbenmansour/terrakube-go/internal/httpapi"
	"github.com/aminbenmansour/terrakube-go/internal/jobstatus"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/logging"
	"github.com/aminbenmansour/terrakube-go/internal/notification"
	"github.com/aminbenmansour/terrakube-go/internal/pat"
	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/runtrigger"
	"github.com/aminbenmansour/terrakube-go/internal/schedule"
	"github.com/aminbenmansour/terrakube-go/internal/secretstore"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/teamtoken"
	"github.com/aminbenmansour/terrakube-go/internal/tfe"
	"github.com/aminbenmansour/terrakube-go/internal/tfoutput"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
	"github.com/aminbenmansour/terrakube-go/internal/ui"
	"github.com/aminbenmansour/terrakube-go/internal/webhook"
	"github.com/aminbenmansour/terrakube-go/migrations"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "terrakube-api: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := logging.New(cfg.LogLevel, cfg.Env)
	slog.SetDefault(logger)

	cipher, err := secretstore.NewCipher(cfg.EncryptionKeys, cfg.EncryptionActiveKey)
	if err != nil {
		return fmt.Errorf("configure secret encryption: %w", err)
	}
	if cipher == nil {
		logger.Warn("secret encryption is disabled: set TERRAKUBE_ENCRYPTION_KEYS")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if cfg.AutoMigrate {
		if err := db.Migrate(ctx, pool, migrations.Files, "."); err != nil {
			return err
		}
		logger.Info("database migrations applied")
		if cipher != nil {
			changed, err := secretstore.Backfill(ctx, store.New(pool), cipher)
			if err != nil {
				return fmt.Errorf("backfill secret encryption: %w", err)
			}
			if changed > 0 {
				logger.Info("encrypted plaintext secrets", "rows", changed)
			}
		}
	}

	q, err := queue.New(cfg.RedisURL)
	if err != nil {
		return err
	}
	defer q.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	err = q.Ping(pingCtx)
	cancel()
	if err != nil {
		return err
	}

	events, err := jobstatus.NewBus(cfg.RedisURL)
	if err != nil {
		return err
	}
	defer events.Close()

	storeBackend, err := storage.New(ctx, cfg.Storage)
	if err != nil {
		return err
	}
	logger.Info("storage configured", "type", cfg.Storage.Type)

	var internalAuth *internalauth.Validator
	if cfg.InternalTokenSecret == "" {
		logger.Warn("internal token authentication is disabled: set TERRAKUBE_INTERNAL_TOKEN_SECRET to enable /tfstate/v1")
	} else {
		internalAuth, err = internalauth.NewValidator(cfg.InternalTokenSecret)
		if err != nil {
			return fmt.Errorf("configure internal token authentication: %w", err)
		}
	}

	// Every job status update is published so the UI's GraphQL subscriptions
	// see it, whether the change came from the management API or the executor.
	base := store.New(pool)
	querier := jobstatus.Wrap(secretstore.Wrap(base, cipher), events.Publish)

	var verifiers []internalauth.Verifier
	if cfg.DexEnabled {
		verifiers = append(verifiers, internalauth.NewDexValidator(cfg.DexIssuerURL,
			append([]string{cfg.DexClientID}, cfg.DexExtraAudiences...), cfg.DexAudienceValidation))
		logger.Info("dex token validation enabled", "issuer", cfg.DexIssuerURL)
	}
	var pats *pat.Service
	var teamTokens *teamtoken.Service
	if cfg.PATSecret != "" {
		pats, err = pat.New(querier, cfg.PATSecret, cfg.AdminGroup)
		if err != nil {
			return fmt.Errorf("configure personal access tokens: %w", err)
		}
		verifiers = append(verifiers, pats)
		// Team tokens share the PAT secret, as upstream does.
		teamTokens, err = teamtoken.New(querier, cfg.PATSecret, cfg.AdminGroup)
		if err != nil {
			return fmt.Errorf("configure team tokens: %w", err)
		}
		verifiers = append(verifiers, teamTokens)
		logger.Info("personal and team access tokens enabled")
	}

	metrics := httpapi.NewMetrics()
	handler := httpapi.New(httpapi.Deps{
		Logger:       logger,
		Service:      "api",
		Started:      time.Now(),
		DB:           pool,
		Redis:        q,
		Storage:      storeBackend,
		RateLimit:    cfg.RateLimit,
		RateBurst:    cfg.RateBurst,
		Metrics:      metrics,
		MetricsToken: cfg.MetricsToken,
	})
	handler.Mount(api.MountPrefix, api.New(api.Deps{
		Store:       querier,
		Lists:       base,
		Logger:      logger,
		PublicURL:   cfg.PublicURL,
		AuthToken:   cfg.AuthToken,
		Verifiers:   verifiers,
		PATs:        pats,
		AdminGroup:  cfg.AdminGroup,
		CORSOrigins: cfg.CORSOrigins,
	}))
	handler.Mount("/pat/v1", api.NewPAT(api.Deps{
		Logger:    logger,
		AuthToken: cfg.AuthToken,
		Verifiers: verifiers,
		PATs:      pats,
	}))
	handler.Mount("/access-token/v1", api.NewAccessToken(api.Deps{
		Logger:       logger,
		AuthToken:    cfg.AuthToken,
		Verifiers:    verifiers,
		CORSOrigins:  cfg.CORSOrigins,
		AdminGroup:   cfg.AdminGroup,
		AccessTokens: querier,
		TeamTokens:   teamTokens,
	}))
	handler.Mount("/notification/v1", api.NewNotification(api.Deps{
		Store:       querier,
		Logger:      logger,
		AuthToken:   cfg.AuthToken,
		Verifiers:   verifiers,
		CORSOrigins: cfg.CORSOrigins,
	}))
	handler.Mount("/webhook/v1", webhook.New(webhook.Deps{
		Store:  querier,
		Queue:  q,
		Logger: logger,
	}))
	handler.Mount("/vcs/v1", api.NewVCS(api.Deps{
		Store:       querier,
		Logger:      logger,
		AuthToken:   cfg.AuthToken,
		Verifiers:   verifiers,
		CORSOrigins: cfg.CORSOrigins,
	}))
	if cfg.TerraformLogin {
		if pats == nil {
			return fmt.Errorf("terraform login requires TERRAKUBE_PAT_SECRET")
		}
		if cfg.DexIssuerURL == "" || cfg.DexClientID == "" {
			return fmt.Errorf("terraform login requires TERRAKUBE_DEX_ISSUER_URL and TERRAKUBE_DEX_CLIENT_ID")
		}
		cookieSigner, err := clilogin.NewCookieSigner(cfg.PATSecret)
		if err != nil {
			return fmt.Errorf("configure terraform login: %w", err)
		}
		loginService := clilogin.NewService(querier,
			clilogin.NewDex(cfg.DexIssuerURL, cfg.DexClientID, cfg.PublicURL+"/oauth/callback"),
			pats, cfg.TerraformLoginDefaultDays, cfg.TerraformLoginMaxDays)
		handler.Mount("/oauth", clilogin.New(clilogin.Deps{
			Service:   loginService,
			Cookie:    cookieSigner,
			PublicURL: cfg.PublicURL,
			Logger:    logger,
		}))
		logger.Info("terraform login enabled", "issuer", cfg.DexIssuerURL)
	}
	tfeDeps := tfe.Deps{
		Store:              querier,
		Storage:            storeBackend,
		Queue:              q,
		Layout:             tfstate.NewLayout(cfg.Storage.Type),
		Logger:             logger,
		PublicURL:          cfg.PublicURL,
		AuthToken:          cfg.AuthToken,
		UploadTicketSecret: cfg.UploadTicketSecret,
		InternalAuth:       internalAuth,
		DexClientID:        cfg.DexClientID,
		DexIssuerURL:       cfg.DexIssuerURL,
		TerraformLogin:     cfg.TerraformLogin,
	}
	graphqlDeps := graphqlapi.Deps{
		DB:         pool,
		Events:     events,
		Logger:     logger,
		AuthToken:  cfg.AuthToken,
		Verifiers:  verifiers,
		AdminGroup: cfg.AdminGroup,
	}
	// The UI calls these cross-origin when it is served from another host, and
	// browsers enforce CORS on the HTTP query endpoint just like any other.
	cors := httpapi.CORS(cfg.CORSOrigins)
	handler.Mount("/graphql/api/v1", cors(httpapi.LimitBody(jsonapi.MaxBodyBytes)(graphqlapi.New(graphqlDeps))))
	handler.Mount("/subscriptions", cors(graphqlapi.NewSubscriptions(graphqlDeps)))
	handler.Get("/.well-known/terraform.json", tfe.Discovery(tfeDeps))
	handler.Mount("/remote/tfe/v2", tfe.New(tfeDeps))
	handler.Mount("/tfoutput/v1", tfoutput.New(tfoutput.Deps{
		Store:        querier,
		Queue:        q,
		Storage:      storeBackend,
		Layout:       tfstate.NewLayout(cfg.Storage.Type),
		Logger:       logger,
		AuthToken:    cfg.AuthToken,
		InternalAuth: internalAuth,
	}))
	handler.Mount("/tfstate/v1", tfstate.New(tfstate.Deps{
		Store:             querier,
		Storage:           storeBackend,
		Logger:            logger,
		Layout:            tfstate.NewLayout(cfg.Storage.Type),
		Auth:              internalAuth,
		AuthToken:         cfg.AuthToken,
		MaxStateSizeBytes: cfg.StateMaxSizeBytes,
		PublicURL:         cfg.PublicURL,
	}))

	// The UI is mounted last so every API path wins over the catch-all SPA
	// fallback. Serving it from the API keeps it on one origin, which is why no
	// CORS exception is needed.
	if uiHandler := ui.New(ui.Deps{
		Dir:          cfg.UIDir,
		PublicURL:    cfg.PublicURL,
		DexClientID:  cfg.DexClientID,
		DexAuthority: cfg.DexIssuerURL,
		Logger:       logger,
	}); uiHandler != nil {
		handler.Mount("/", uiHandler)
	}

	if cfg.MetricsToken != "" {
		go metrics.Sample(ctx, logger, metricsSource{store: querier, queue: q}, 30*time.Second)
	}
	// Every replica drains the notification outbox; the claim is a conditional
	// update with SKIP LOCKED, so no leader lock is needed.
	go notification.NewDispatcher(querier, notification.NewDeliverer(true), cfg.PublicURL, logger).Run(ctx)
	startScheduler(ctx, logger, pool, querier, q)

	return httpapi.Serve(ctx, logger, cfg.HTTPAddr, handler, cfg.ShutdownTimeout)
}

// metricsSource feeds the job and queue gauges from the database and Redis.
type metricsSource struct {
	store store.Querier
	queue *queue.Queue
}

func (s metricsSource) JobStatusCounts(ctx context.Context) (map[string]int64, error) {
	rows, err := s.store.CountJobsByStatus(ctx)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int64, len(rows))
	for _, row := range rows {
		counts[row.Status] = row.Total
	}
	return counts, nil
}

func (s metricsSource) QueueDepth(ctx context.Context) (int64, error) {
	return s.queue.Depth(ctx, queue.JobsStream)
}

// startScheduler runs the workspace dispatcher on the replica holding a
// session-level Postgres advisory lock, so two API instances never queue the
// same tick twice. A replica that cannot take the lock retries; when the
// holder exits, its connection closes and the next replica takes over.
func startScheduler(ctx context.Context, logger *slog.Logger, pool *pgxpool.Pool, st store.Querier, q *queue.Queue) {
	const lockKey int64 = 0x74657272616b7562 // "terrakub"
	go func() {
		for ctx.Err() == nil {
			conn, err := pool.Acquire(ctx)
			if err != nil {
				logger.Warn("schedule dispatcher cannot acquire a connection", "error", err)
				if !waitFor(ctx, 30*time.Second) {
					return
				}
				continue
			}
			var locked bool
			if err := conn.QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", lockKey).Scan(&locked); err != nil || !locked {
				conn.Release()
				if err != nil {
					logger.Warn("schedule dispatcher lock failed", "error", err)
				}
				if !waitFor(ctx, 30*time.Second) {
					return
				}
				continue
			}
			logger.Info("schedule dispatcher active")
			go runtrigger.NewDispatcher(st, q, logger).Run(ctx)
			schedule.NewDispatcher(st, q, logger).Run(ctx)
			conn.Release()
			return
		}
	}()
}

func waitFor(ctx context.Context, d time.Duration) bool {
	select {
	case <-ctx.Done():
		return false
	case <-time.After(d):
		return true
	}
}
