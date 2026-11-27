-- name: ListNotificationConfigurationsByOrganization :many
SELECT * FROM notification_configuration WHERE organization_id = $1 ORDER BY name, id;

-- name: ListNotificationConfigurationsByWorkspace :many
SELECT * FROM notification_configuration WHERE workspace_id = $1 ORDER BY name, id;

-- name: GetNotificationConfiguration :one
SELECT * FROM notification_configuration WHERE id = $1;

-- name: CreateNotificationConfiguration :one
INSERT INTO notification_configuration (id, name, description, organization_id, workspace_id,
    channel_type, destination_url, signing_secret, active, message_style, created_date, updated_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now(), now())
RETURNING *;

-- name: UpdateNotificationConfiguration :one
UPDATE notification_configuration
SET name = $2,
    description = $3,
    channel_type = $4,
    destination_url = $5,
    signing_secret = $6,
    active = $7,
    message_style = $8,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteNotificationConfiguration :exec
DELETE FROM notification_configuration WHERE id = $1;

-- name: ListNotificationTriggersByConfiguration :many
SELECT * FROM notification_trigger WHERE configuration_id = $1 ORDER BY created_date, id;

-- name: GetNotificationTrigger :one
SELECT * FROM notification_trigger WHERE id = $1;

-- name: CreateNotificationTrigger :one
INSERT INTO notification_trigger (id, job_status, configuration_id, created_date, updated_date)
VALUES ($1, $2, $3, now(), now())
RETURNING *;

-- name: DeleteNotificationTrigger :exec
DELETE FROM notification_trigger WHERE id = $1;

-- name: ListNotificationConfigurationTemplateIDs :many
SELECT template_id FROM notification_configuration_template WHERE configuration_id = $1 ORDER BY template_id;

-- name: DeleteNotificationConfigurationTemplates :exec
DELETE FROM notification_configuration_template WHERE configuration_id = $1;

-- name: AddNotificationConfigurationTemplate :exec
INSERT INTO notification_configuration_template (configuration_id, template_id) VALUES ($1, $2);

-- name: ListNotificationSecrets :many
SELECT id, signing_secret FROM notification_configuration;

-- name: ListNotificationSecretsAfter :many
SELECT id, signing_secret FROM notification_configuration
WHERE id > $1 ORDER BY id LIMIT $2;

-- name: SetNotificationSecret :exec
UPDATE notification_configuration SET signing_secret = $2 WHERE id = $1 AND signing_secret = $3;

-- name: ClaimNotificationOutbox :one
-- One due row per call, oldest first. SKIP LOCKED lets every API replica drain
-- the outbox without a leader lock; the update itself is the claim.
UPDATE notification_outbox
SET status = 'SENDING',
    attempt_count = attempt_count + 1,
    last_attempt_at = now(),
    updated_date = now()
WHERE id = (
    SELECT id FROM notification_outbox
    WHERE status = 'PENDING'
      AND (next_attempt_at IS NULL OR next_attempt_at <= now())
    ORDER BY created_date, id
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;

-- name: RecordNotificationOutbox :execrows
-- Keyed on the attempt's own last_attempt_at so a result from an attempt that
-- outlived the stuck-row sweep cannot overwrite the newer attempt.
UPDATE notification_outbox
SET status = $3,
    last_error = $4,
    next_attempt_at = $5,
    updated_date = now()
WHERE id = $1
  AND status = 'SENDING'
  AND last_attempt_at = $2;

-- name: ReclaimStuckNotificationOutbox :execrows
UPDATE notification_outbox
SET status = 'PENDING', updated_date = now()
WHERE status = 'SENDING' AND last_attempt_at < $1 AND attempt_count < $2;

-- name: FailStuckNotificationOutbox :execrows
UPDATE notification_outbox
SET status = 'FAILED', last_error = $3, updated_date = now()
WHERE status = 'SENDING' AND last_attempt_at < $1 AND attempt_count >= $2;

-- name: RearmNotificationOutbox :execrows
UPDATE notification_outbox
SET status = 'PENDING',
    attempt_count = 0,
    last_error = NULL,
    next_attempt_at = NULL,
    updated_date = now()
WHERE id = $1 AND status = 'FAILED';

-- name: NotificationOutboxExistsForWorkspace :one
SELECT EXISTS (
    SELECT 1 FROM notification_outbox o
    JOIN job j ON j.id = o.job_id
    WHERE o.id = $1 AND j.workspace_id = $2
);

-- name: ListNotificationDeliveriesByWorkspace :many
SELECT o.*, c.name AS configuration_name, c.channel_type
FROM notification_outbox o
JOIN notification_configuration c ON c.id = o.configuration_id
JOIN job j ON j.id = o.job_id
WHERE j.workspace_id = $1
ORDER BY o.created_date DESC, o.id
LIMIT $2;
