-- name: GetWebhook :one
SELECT * FROM webhook WHERE id = $1;

-- name: GetWebhookByWorkspace :one
SELECT * FROM webhook WHERE workspace_id = $1 ORDER BY created_date, id LIMIT 1;

-- name: CreateWebhook :one
INSERT INTO webhook (id, workspace_id, remote_hook_id, migrated_v2, created_date, updated_date)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: UpdateWebhook :one
UPDATE webhook
SET remote_hook_id = $2,
    migrated_v2 = $3,
    sync_error = $4,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteWebhook :exec
DELETE FROM webhook WHERE id = $1;

-- name: ListWebhookEvents :many
SELECT * FROM webhook_event WHERE webhook_id = $1 ORDER BY priority, created_date, id;

-- name: GetWebhookEvent :one
SELECT * FROM webhook_event WHERE id = $1;

-- name: CreateWebhookEvent :one
INSERT INTO webhook_event (id, webhook_id, branch, path, path_type, event, template_id, priority,
                           pr_workflow_enabled, pr_apply_enabled, created_date, updated_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now(), now())
RETURNING *;

-- name: UpdateWebhookEvent :one
UPDATE webhook_event
SET branch = $2,
    path = $3,
    path_type = $4,
    event = $5,
    template_id = $6,
    priority = $7,
    pr_workflow_enabled = $8,
    pr_apply_enabled = $9,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteWebhookEvent :exec
DELETE FROM webhook_event WHERE id = $1;

-- name: RecordWebhookDelivery :execrows
-- The composite key makes a provider retry a no-op: zero rows means this
-- delivery id was already accepted for this webhook.
INSERT INTO webhook_delivery (webhook_id, delivery_id, created_date)
VALUES ($1, $2, now())
ON CONFLICT (webhook_id, delivery_id) DO NOTHING;

-- name: DeleteWebhookDelivery :exec
-- Compensates a recorded delivery whose run could not be created, so the
-- provider's retry is processed instead of being mistaken for a replay.
DELETE FROM webhook_delivery WHERE webhook_id = $1 AND delivery_id = $2;

-- name: PruneWebhookDeliveries :exec
DELETE FROM webhook_delivery WHERE created_date < now() - interval '30 days';
