# Threat model

Scope: the Go services in this repository — `cmd/api` (management API, TFE cloud
API, GraphQL, state backend), `cmd/executor`, `cmd/registry` — and the local
storage and deployment assets they use. Upstream Terrakube's Java code and the
React UI image are out of scope; see [SECURITY.md](../SECURITY.md).

## Assets

| Asset | Why it matters |
| --- | --- |
| Terraform state | Often contains credentials and private addresses |
| Variable values, SSH private keys, VCS tokens | Credentials the platform holds on users' behalf |
| PATs, team tokens, internal JWTs | Impersonate a user or a job |
| Database | Every asset above, plus the permission graph |
| Object storage | State archives and run logs |
| The executor host | Runs untrusted Terraform/OpenTofu code with credentials in its environment |

## Trust boundaries

1. **Browser/CLI → API.** Unauthenticated until a bearer token proves an
   identity (dev token, Dex OIDC JWT, or PAT). CORS decides which origins a
   browser may call; the REST surface is also reachable without CORS by any
   non-browser client.
2. **API → database.** Parameterised SQL only; secrets are encrypted before
   they reach a column (`internal/secretstore`).
3. **API → executor.** Job creation writes a queue message; the executor claims
   it. Both services share the database and Redis and are assumed to be on the
   same trusted network.
4. **Executor → Terraform child process.** Untrusted code from a user's
   repository. It sees the decrypted variables of its own workspace only, in a
   per-job working directory.
5. **API/executor → object storage.** Keys are namespaced by
   organization/workspace/job; the local backend rejects keys that escape its
   root.
6. **API → outbound HTTP (future).** VCS webhooks and notification senders will
   make requests to user-supplied URLs. This boundary is not implemented yet
   and has no guard; see accepted risks.

## Threats and controls

| Threat | Control |
| --- | --- |
| Read or write another organization's resources | RBAC (`internal/rbac`) resolves the caller's teams and grants per organization, project and workspace; cross-org access answers `404`/`403` and is covered by tests |
| Forge an executor token to read another workspace's state or logs | Internal JWTs are signed with a shared secret, classified as workspace-executor or service, and checked against the workspace and the active step before state or log reads |
| Steal a PAT from the database | A PAT is an HMAC-signed JWT; the database holds only its metadata row and a soft-delete flag, and the verifier checks revocation on every request. The token value is shown once, at creation |
| Read secrets from a database dump | Variable values, SSH keys and VCS credentials are encrypted at rest with a rotating key (`make rotate-keys`); secret columns are never returned in API responses |
| Secret printed in a log or error | Handlers never interpolate secret values; decryption failures report the record id, not the value |
| Oversized or slow request exhausts memory | JSON:API and GraphQL bodies are capped at 8 MiB (`413`), state and configuration uploads at 100 MiB each, archive extraction at 128 MiB/file, 1 GiB/archive and 100 000 entries; an optional per-caller rate limit answers `429` with `Retry-After` |
| Decompression or path traversal through an archive | Extraction rejects escaping entries and caps sizes as above; storage keys are cleaned and refused when they leave the root |
| Malicious state or module content | State is treated as opaque bytes; the registry serves module/provider archives without executing them |
| Clickjacking, MIME sniffing, mixed content | `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `Referrer-Policy` and HSTS over TLS |
| Credential theft through the API from a malicious site | CORS defaults to `*` for compatibility but `TERRAKUBE_CORS_ORIGINS` pins it; bearer tokens are not cookies, so a cross-site request cannot read responses without a token |
| Sandbox escape from Terraform code | Not mitigated: Terraform runs as the executor user on the job host. Operators isolate the executor (container, VM, separate credentials per environment) |

## Accepted risks

| Risk | Why | Owner | Expiry |
| --- | --- | --- | --- |
| `TERRAKUBE_AUTH_TOKEN` accepts a static token | Development convenience | Maintainers | Until #45 lands PAT-only operation for production deployments |
| CORS `Access-Control-Allow-Origin: *` by default | The UI served by this API depends on it; split-origin deployments are rare | Maintainers | Review when the UI is always served from the API origin |
| No content security policy | The UI bundle owns its asset policy and needs a verified policy, not a guessed one | Maintainers | #57, when a UI build is available to test against |
| Rate limiting off by default | A wrong default breaks the UI's fan-out and parallel CLI runs | Maintainers | Review after load testing the UI |
| VCS/notification URLs are not SSRF-checked | There are no outbound calls yet; the check will land with the first caller | Maintainers | #48 and #52 |
| Executor has no network egress policy | Terraform providers need arbitrary egress; isolation is a deployment concern | Operators | Documented in `deploy/`; revisit in #56 |
| Job logs are stored without encryption at rest | Logs are bounded and workspace-scoped; storage encryption is delegated to the operator | Maintainers | Review with #56 backup/restore work |

Exceptions are reviewed when the linked issue lands; a risk that is not
addressed by then is re-evaluated rather than silently renewed.
