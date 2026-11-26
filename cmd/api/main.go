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
	"github.com/aminbenmansour/terrakube-go/internal/pat"
	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/secretstore"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfe"
	"github.com/aminbenmansour/terrakube-go/internal/tfoutput"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
	"github.com/aminbenmansour/terrakube-go/internal/ui"
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
	if cfg.PATSecret != "" {
		pats, err = pat.New(querier, cfg.PATSecret, cfg.AdminGroup)
		if err != nil {
			return fmt.Errorf("configure personal access tokens: %w", err)
		}
		verifiers = append(verifiers, pats)
		logger.Info("personal access tokens enabled")
	}

	handler := httpapi.New(httpapi.Deps{
		Logger:    logger,
		Service:   "api",
		Started:   time.Now(),
		DB:        pool,
		Redis:     q,
		Storage:   storeBackend,
		RateLimit: cfg.RateLimit,
		RateBurst: cfg.RateBurst,
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

	return httpapi.Serve(ctx, logger, cfg.HTTPAddr, handler, cfg.ShutdownTimeout)
}
