// Package webhook receives provider webhook deliveries for a workspace's
// registered hook and starts a run for the matching trigger.
//
// The endpoint is public: the caller is authenticated by the provider's
// signature, using base64(workspace id) as the secret upstream registers with
// GitHub and GitLab. An unknown webhook, a mismatched repository or an event no
// trigger matches is acknowledged without starting a run.
package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/httpapi"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	vcsprovider "github.com/aminbenmansour/terrakube-go/internal/vcs"
)

// JobQueue publishes a created run to the executor.
type JobQueue interface {
	EnqueueJob(ctx context.Context, jobID int64) (string, error)
}

// Deps are the collaborators the receiver needs.
type Deps struct {
	Store  store.Querier
	Queue  JobQueue
	Logger *slog.Logger
}

// Receiver turns a verified delivery into a run.
type Receiver struct {
	deps Deps
}

// New builds the public /webhook/v1 router. The mount prefix is supplied by
// cmd/api, so the route here is the webhook id.
func New(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	receiver := &Receiver{deps: deps}
	r := chi.NewRouter()
	r.Use(httpapi.LimitBody(jsonapi.MaxBodyBytes))
	r.Post("/{webhookID}", receiver.serveHTTP)
	return r
}

func (r *Receiver) serveHTTP(w http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		r.deps.Logger.Warn("read webhook body", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	status, err := r.accept(request.Context(), chi.URLParam(request, "webhookID"), request.Header, body)
	if err != nil {
		r.deps.Logger.Error("process webhook delivery", "status", status, "error", err)
	}
	w.WriteHeader(status)
}

// accept verifies and handles one delivery. The returned status is what the
// provider sees; an error is only for logging.
func (r *Receiver) accept(ctx context.Context, webhookID string, headers http.Header, body []byte) (int, error) {
	if _, err := uuid.Parse(webhookID); err != nil {
		return http.StatusOK, nil
	}
	hook, err := r.deps.Store.GetWebhook(ctx, webhookID)
	if errors.Is(err, pgx.ErrNoRows) {
		return http.StatusOK, nil
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	workspace, err := r.deps.Store.GetWorkspace(ctx, hook.WorkspaceID)
	if errors.Is(err, pgx.ErrNoRows) {
		return http.StatusOK, nil
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if workspace.VcsID == nil || *workspace.VcsID == "" {
		return http.StatusOK, nil
	}
	connection, err := r.deps.Store.GetVCS(ctx, *workspace.VcsID)
	if errors.Is(err, pgx.ErrNoRows) {
		return http.StatusOK, nil
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	secret := base64.StdEncoding.EncodeToString([]byte(workspace.ID))
	var (
		delivery Delivery
		handled  bool
	)
	switch connection.VcsType {
	case vcsprovider.TypeGitHub:
		if !verifyGitHubSignature(secret, body, headers.Get("X-Hub-Signature-256")) {
			return http.StatusUnauthorized, nil
		}
		delivery, handled = parseGitHub(strings.ToLower(headers.Get("X-GitHub-Event")), body)
	case vcsprovider.TypeGitLab:
		if subtle.ConstantTimeCompare([]byte(headers.Get("X-Gitlab-Token")), []byte(secret)) != 1 {
			return http.StatusUnauthorized, nil
		}
		delivery, handled = parseGitLab(body)
	default:
		// Bitbucket and Azure DevOps registration is not implemented, so no
		// delivery for them can be genuine.
		return http.StatusOK, nil
	}
	if !handled || delivery.Event == eventPing || delivery.Event == "" {
		return http.StatusOK, nil
	}
	if NormalizeRepository(delivery.Repository) == "" ||
		NormalizeRepository(delivery.Repository) != NormalizeRepository(strValue(workspace.Source)) {
		r.deps.Logger.Info("webhook repository does not match the workspace source",
			"webhook_id", hook.ID, "workspace_id", workspace.ID)
		return http.StatusOK, nil
	}
	delivery.Files = trimDedupe(delivery.Files)

	events, err := r.deps.Store.ListWebhookEvents(ctx, hook.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	trigger, ok := matchTrigger(events, delivery)
	if !ok {
		r.deps.Logger.Info("no webhook trigger matches the delivery",
			"webhook_id", hook.ID, "event", delivery.Event, "branch", delivery.Branch)
		return http.StatusOK, nil
	}

	deliveryID := deliveryID(headerValue(headers, "X-GitHub-Delivery", "X-Gitlab-Event-UUID", "Idempotency-Key"), body)
	rows, err := r.deps.Store.RecordWebhookDelivery(ctx, store.RecordWebhookDeliveryParams{
		WebhookID: hook.ID, DeliveryID: deliveryID,
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if rows == 0 {
		r.deps.Logger.Info("ignoring replayed webhook delivery", "webhook_id", hook.ID, "delivery_id", deliveryID)
		return http.StatusOK, nil
	}
	_ = r.deps.Store.PruneWebhookDeliveries(ctx)

	if err := r.startRun(ctx, workspace, connection.VcsType, trigger, delivery); err != nil {
		// The run is missing, so let the provider's retry be processed.
		_ = r.deps.Store.DeleteWebhookDelivery(ctx, store.DeleteWebhookDeliveryParams{
			WebhookID: hook.ID, DeliveryID: deliveryID,
		})
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// startRun creates and queues the normal plan, approval and apply for the
// matching trigger. A template that has left the workspace's organization is
// treated as no configuration, not an error, as upstream does.
func (r *Receiver) startRun(ctx context.Context, workspace store.Workspace, provider string, trigger store.WebhookEvent, delivery Delivery) error {
	template, err := r.deps.Store.GetTemplate(ctx, trigger.TemplateID)
	if errors.Is(err, pgx.ErrNoRows) {
		r.deps.Logger.Info("webhook trigger template is gone", "webhook_id", trigger.WebhookID, "template_id", trigger.TemplateID)
		return nil
	}
	if err != nil {
		return err
	}
	if template.OrganizationID != workspace.OrganizationID {
		r.deps.Logger.Warn("webhook trigger template belongs to another organization",
			"webhook_id", trigger.WebhookID, "template_id", trigger.TemplateID)
		return nil
	}

	branch := delivery.Branch
	if delivery.Event == eventRelease && delivery.ReleaseTag != "" {
		branch = "refs/tags/" + delivery.ReleaseTag
	}
	if branch == "" {
		return nil
	}
	job, err := r.deps.Store.CreateJob(ctx, store.CreateJobParams{
		OrganizationID:    workspace.OrganizationID,
		WorkspaceID:       workspace.ID,
		Status:            "pending",
		Via:               viaFor(provider),
		Refresh:           true,
		PlanChanges:       true,
		CreatedBy:         nonEmptyPtr(delivery.Actor),
		Tcl:               template.Tcl,
		TemplateReference: stringPtr(template.ID),
		CommitID:          nonEmptyPtr(delivery.Commit),
		OverrideBranch:    stringPtr(branch),
	})
	if err != nil {
		return err
	}
	for _, step := range []struct {
		number int32
		name   string
	}{
		{100, "Terraform Plan"},
		{150, "Approve Plan"},
		{200, "Terraform Apply"},
	} {
		if _, err := r.deps.Store.CreateStep(ctx, store.CreateStepParams{
			ID:         uuid.NewString(),
			JobID:      job.ID,
			StepNumber: step.number,
			Name:       stringPtr(step.name),
			Status:     stringPtr("pending"),
			LogStatus:  stringPtr("pending"),
		}); err != nil {
			_, _ = r.deps.Store.UpdateJobStatus(ctx, store.UpdateJobStatusParams{ID: job.ID, Status: "failed"})
			return err
		}
	}
	_ = r.deps.Store.SetWorkspaceLastJob(ctx, store.SetWorkspaceLastJobParams{
		ID: workspace.ID, LastJobStatus: "pending",
	})
	if r.deps.Queue == nil {
		_, _ = r.deps.Store.UpdateJobStatus(ctx, store.UpdateJobStatusParams{ID: job.ID, Status: "failed"})
		return errors.New("run queue is not configured")
	}
	if _, err := r.deps.Queue.EnqueueJob(ctx, job.ID); err != nil {
		_, _ = r.deps.Store.UpdateJobStatus(ctx, store.UpdateJobStatusParams{ID: job.ID, Status: "failed"})
		return err
	}
	r.deps.Logger.Info("webhook run queued", "workspace_id", workspace.ID, "job_id", job.ID,
		"event", delivery.Event, "branch", delivery.Branch)
	return nil
}

func verifyGitHubSignature(secret string, body []byte, signature string) bool {
	if signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

func headerValue(headers http.Header, names ...string) string {
	for _, name := range names {
		if value := headers.Get(name); value != "" {
			return value
		}
	}
	return ""
}

func strValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func stringPtr(value string) *string {
	return &value
}
