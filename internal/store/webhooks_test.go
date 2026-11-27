package store

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
)

// TestWebhookQueries covers the workspace's one webhook, its triggers in
// priority order, the provider hook id round-trip and the cascade on delete.
func TestWebhookQueries(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
		ID: uuid.NewString(), Name: "webhook-store", ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization: %v", err)
	}
	workspace, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
		ID: uuid.NewString(), Name: "webhook", IacType: "terraform", OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	template, err := q.CreateTemplate(ctx, CreateTemplateParams{
		ID: uuid.NewString(), OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create template: %v", err)
	}

	hook, err := q.CreateWebhook(ctx, CreateWebhookParams{
		ID: uuid.NewString(), WorkspaceID: workspace.ID,
	})
	if err != nil {
		t.Fatalf("create webhook: %v", err)
	}
	if hook.MigratedV2 || hook.RemoteHookID != nil {
		t.Fatalf("new webhook = %+v", hook)
	}

	remoteHookID := "123456"
	if hook, err = q.UpdateWebhook(ctx, UpdateWebhookParams{
		ID: hook.ID, RemoteHookID: &remoteHookID,
	}); err != nil {
		t.Fatalf("update webhook: %v", err)
	}
	if hook.RemoteHookID == nil || *hook.RemoteHookID != remoteHookID {
		t.Fatalf("remote hook id = %v", hook.RemoteHookID)
	}

	found, err := q.GetWebhookByWorkspace(ctx, workspace.ID)
	if err != nil || found.ID != hook.ID {
		t.Fatalf("get by workspace = %+v/%v", found, err)
	}

	second, err := q.CreateWebhookEvent(ctx, CreateWebhookEventParams{
		ID: uuid.NewString(), WebhookID: hook.ID, Branch: ".*", Path: "**",
		PathType: "PATTERN", Event: "PUSH", TemplateID: template.ID, Priority: 2,
	})
	if err != nil {
		t.Fatalf("create second trigger: %v", err)
	}
	first, err := q.CreateWebhookEvent(ctx, CreateWebhookEventParams{
		ID: uuid.NewString(), WebhookID: hook.ID, Branch: "main", Path: "modules/**",
		PathType: "PATTERN", Event: "PULL_REQUEST", TemplateID: template.ID, Priority: 1,
		PrWorkflowEnabled: true,
	})
	if err != nil {
		t.Fatalf("create first trigger: %v", err)
	}

	events, err := q.ListWebhookEvents(ctx, hook.ID)
	if err != nil {
		t.Fatalf("list events: %v", err)
	}
	if len(events) != 2 || events[0].ID != first.ID || events[1].ID != second.ID {
		t.Fatalf("events = %+v, want priority order", events)
	}
	if !events[0].PrWorkflowEnabled || events[0].PrApplyEnabled {
		t.Fatalf("first trigger = %+v", events[0])
	}

	if updated, err := q.UpdateWebhookEvent(ctx, UpdateWebhookEventParams{
		ID: first.ID, Branch: "release/*", Path: "**", PathType: "REGEX", Event: "PUSH",
		TemplateID: template.ID, Priority: 0,
	}); err != nil {
		t.Fatalf("update event: %v", err)
	} else if updated.Branch != "release/*" || updated.PathType != "REGEX" || updated.PrWorkflowEnabled {
		t.Fatalf("updated trigger = %+v", updated)
	}

	delivery := RecordWebhookDeliveryParams{WebhookID: hook.ID, DeliveryID: "delivery-1"}
	if rows, err := q.RecordWebhookDelivery(ctx, delivery); err != nil || rows != 1 {
		t.Fatalf("record delivery = %d/%v, want one row", rows, err)
	}
	if rows, err := q.RecordWebhookDelivery(ctx, delivery); err != nil || rows != 0 {
		t.Fatalf("replayed delivery = %d/%v, want the conflict ignored", rows, err)
	}
	if err := q.DeleteWebhookDelivery(ctx, DeleteWebhookDeliveryParams(delivery)); err != nil {
		t.Fatalf("delete delivery: %v", err)
	}
	if rows, err := q.RecordWebhookDelivery(ctx, delivery); err != nil || rows != 1 {
		t.Fatalf("record after compensation = %d/%v", rows, err)
	}
	if err := q.PruneWebhookDeliveries(ctx); err != nil {
		t.Fatalf("prune deliveries: %v", err)
	}

	if err := q.DeleteWebhook(ctx, hook.ID); err != nil {
		t.Fatalf("delete webhook: %v", err)
	}
	if _, err := q.GetWebhook(ctx, hook.ID); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("webhook after delete = %v", err)
	}
	if events, err := q.ListWebhookEvents(ctx, hook.ID); err != nil || len(events) != 0 {
		t.Fatalf("events after delete = %+v/%v", events, err)
	}
}
