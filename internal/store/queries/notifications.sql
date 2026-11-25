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
