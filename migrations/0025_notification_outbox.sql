-- Notification outbox: a triggering job status change writes one row per
-- matching configuration, and a poller delivers it. Follows upstream's
-- changelog-2.33.0-notification.xml and changelog-2.33.1-notification-outbox-job-status.xml.
--
-- payload stays nullable: this implementation renders the payload at dispatch
-- time so the two writers (the executor and the management API) only need the
-- same SQL statement that changes the job status. The column exists for schema
-- parity with upstream.
CREATE TABLE notification_outbox (
    id               varchar(36) PRIMARY KEY,
    job_id           bigint      NOT NULL REFERENCES job (id) ON DELETE CASCADE,
    configuration_id varchar(36) NOT NULL REFERENCES notification_configuration (id) ON DELETE CASCADE,
    payload          text,
    status           varchar(20) NOT NULL DEFAULT 'PENDING',
    -- The status the payload is rendered for. Re-checked against the job at
    -- claim time so a transient "failed" that the executor corrected never
    -- sends a false alarm.
    job_status       varchar(20),
    attempt_count    integer     NOT NULL DEFAULT 0,
    last_attempt_at  timestamptz,
    next_attempt_at  timestamptz,
    last_error       text,
    created_date     timestamptz NOT NULL DEFAULT now(),
    updated_date     timestamptz,
    created_by       varchar(128),
    updated_by       varchar(128)
);

-- Idempotency per (run, configuration, status). SKIPPED rows leave the index so
-- a later, genuine occurrence of the same status can still notify.
CREATE UNIQUE INDEX uk_notification_outbox_run_configuration_status
    ON notification_outbox (job_id, configuration_id, job_status)
    WHERE status <> 'SKIPPED';

-- The poller's due-row scan: status, then next attempt, then creation order.
CREATE INDEX idx_notification_outbox_due
    ON notification_outbox (status, next_attempt_at, created_date);

-- The stuck-row sweep.
CREATE INDEX idx_notification_outbox_status_last_attempt
    ON notification_outbox (status, last_attempt_at);

-- The delivery-history read: one workspace's runs.
CREATE INDEX idx_notification_outbox_job ON notification_outbox (job_id);
