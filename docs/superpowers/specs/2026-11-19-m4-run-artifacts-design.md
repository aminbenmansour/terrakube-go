# M4 run logs and artifacts — design

**Issues:** #33 (live logs and durable final logs), #34 (state, plan and
history after a run).

## Intent

A run must leave behind logs and state that survive an API or executor
restart, and every uploaded artifact must be discoverable from the job,
step or history record that produced it. A failed upload must fail the run,
not disappear into a warning log.

## Goals

- Final step logs are archived to object storage under the upstream
  `tfoutput/{org}/{job}/{step}.tfoutput` key.
- `/tfoutput/v1` serves the archived object for finished steps and falls
  back to `step.output` for rows written before this change.
- `history.output` carries the upstream state JSON URL.
- A plan, log or history upload failure marks the run failed.
- Live streaming, ordering, step boundaries and the existing 10 MiB bound
  stay unchanged.

## Non-goals

- Switching `step.output` to the upstream `/tfoutput/v1/...` URL. Every
  reader (TFE run logs, UI) would have to resolve URLs; not required by the
  issues. It remains a documented divergence.
- Multi-part log uploads, log retention policies, or a new schema.
- State upload from the executor. Terraform still writes state through the
  HTTP backend; the executor archives the resulting state into history.

## Approaches considered

1. **Archive alongside the inline column (chosen).** Keep `step.output` as
   text, write the same bounded text to `tfoutput/...`, and teach the log
   handler to prefer the object. Small, backward compatible, satisfies
   durability.
2. **Full upstream parity.** Set `step.output` to a URL and serve only from
   object storage. Better parity but touches every log reader and risks
   breaking TFE `runs/{id}/logs` for no acceptance-criteria gain.
3. **Keep Postgres only.** Rejected: #33 explicitly asks for object storage
   and #34 asks for artifacts linked to records.

## Decisions

- `tfstate.Layout.Output(orgID, jobID, stepID)` returns
  `tfoutput/{org}/{job}/{step}.tfoutput` for both local and S3 backends.
  Upstream's local backend uses a home-directory path; the storage
  abstraction uses the S3-style key for both. Documented.
- `completeStep` archives the log before writing the status. On upload
  failure it marks the step `failed` and returns an error when the intended
  status was `completed`. A `failed` or `cancelled` step keeps its status
  (the log failure is logged).
- `runApply` archives history before completing the step. A history failure
  marks the step `failed` and returns the error, so `processJob` fails the
  run.
- `history.output` is `{PublicURL}/tfstate/v1/organization/{org}/workspace/{ws}/state/{historyID}.json`.
  When `PublicURL` is empty the field stays empty.
- No schema or HTTP API changes.

## What is left to the builder

- Exact test names and stubs; the plan names them.
- Whether to archive an empty log (the plan archives it; an empty object is
  harmless and keeps the key predictable).
- Commit split; the plan uses one commit per task.
