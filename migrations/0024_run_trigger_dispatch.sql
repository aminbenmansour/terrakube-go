-- Run trigger dispatch: a completed job fans out to the destination workspaces
-- of every enabled edge. Columns and the outbox follow
-- api/src/main/resources/db/changelog/local/changelog-2.34.0-run-trigger-event.xml.

-- Which job caused this one, and how many hops in. Both default so existing
-- rows need no backfill; the depth is the runtime net against a cycle that
-- slipped past the write-time check.
ALTER TABLE job ADD COLUMN triggered_by_job_id bigint;
ALTER TABLE job ADD COLUMN cascade_depth integer NOT NULL DEFAULT 0;

-- Backs the fan-out idempotency lookup, one query per dependent.
CREATE INDEX idx_job_workspace_triggered_by_job_id ON job (workspace_id, triggered_by_job_id);

-- The durable outbox. A row is written in the same statement that marks a job
-- completed, so a crash between the two cannot lose a dispatch. It is internal
-- plumbing, not a JSON:API resource.
CREATE TABLE run_trigger_event (
    id              varchar(36) PRIMARY KEY,
    job_id          bigint      NOT NULL REFERENCES job (id) ON DELETE CASCADE,
    status          varchar(20) NOT NULL DEFAULT 'PENDING',
    attempt_count   integer     NOT NULL DEFAULT 0,
    last_attempt_at timestamptz,
    next_attempt_at timestamptz,
    last_error      text,
    created_date    timestamptz NOT NULL DEFAULT now(),
    updated_date    timestamptz,
    CONSTRAINT uk_run_trigger_event_job_id UNIQUE (job_id)
);

-- The poller's due-row scan: status, then creation order.
CREATE INDEX idx_run_trigger_event_status_created_date
    ON run_trigger_event (status, created_date);

-- The stuck-row sweep: PROCESSING rows whose lease expired.
CREATE INDEX idx_run_trigger_event_status_last_attempt
    ON run_trigger_event (status, last_attempt_at);
