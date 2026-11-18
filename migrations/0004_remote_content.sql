-- Remote-content job linkage and executor step bookkeeping.
--
-- Mirrors the upstream `job.override_source` / `job.override_branch` columns
-- from the remote-state changelog. CLI-created runs store the uploaded
-- configuration tarball URL in override_source and use the reserved
-- `remote-content` branch, exactly like upstream.

ALTER TABLE job
    ADD COLUMN override_source text,
    ADD COLUMN override_branch varchar(64);

-- The executor resumes a run by looking for the next pending step.
CREATE INDEX step_job_number_idx ON step (job_id, step_number);
