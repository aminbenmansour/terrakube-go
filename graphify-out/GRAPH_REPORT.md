# Graph Report - terrakube-go  (2026-11-19)

## Corpus Check
- 112 files · ~69,258 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 6 file(s) not represented in the graph (top: (none) 4, .example 1, .jsonc 1)

## Summary
- 1311 nodes · 4706 edges · 71 communities (52 shown, 19 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 763 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3f891507`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.Request
- JSON:API Resource Models
- Agent Workflows and Compatibility Docs
- WriteResource
- memStorage
- Core Database Schema
- .QueryRow
- .streamStepLog
- Handler
- TFE Handler Tests
- stubQuerier
- State Backend Tests
- DecodeRequest
- Organization Store Queries
- steps
- Querier
- New
- transientError
- Project
- config_test.go
- sqlc Generation Config
- tarGzEntries
- net/http.Handler
- Resource
- testing.T
- context.Context
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- M4 run logs and artifacts Implementation Plan
- Job
- io.Reader
- OpenTofu State Conformance
- Compose Project Definition
- Compose Healthchecks
- Compose Service Dependencies
- MinIO Compose Service
- Redis Compose Service
- M4 run logs and artifacts — design
- Management API Smoke Test
- TFE API Smoke Test
- State Backend Smoke Test
- Go Module Definition
- context.Context
- context.Context
- Backend
- executor-e2e.sh
- NewLayout
- TempArchive
- net/http.Handler
- Workspacetag
- .cloneConfiguration
- .Get
- newRedactor
- query.go
- Object
- cleanKey
- Layout
- Local
- Discovery
- .GetWorkspaceByName
- .UpdateWorkspace

## God Nodes (most connected - your core abstractions)
1. `Querier` - 81 edges
2. `API` - 66 edges
3. `Queries` - 64 edges
4. `Job` - 64 edges
5. `stubQuerier` - 57 edges
6. `Workspace` - 56 edges
7. `New()` - 44 edges
8. `Executor` - 42 edges
9. `Step` - 40 edges
10. `WriteError()` - 36 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `Layout`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfstate/layout.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `New()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfoutput/handler.go
- `Generated Code and Migration Policy` --semantically_similar_to--> `sqlc Workflow`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .opencode/skills/sqlc-workflow/SKILL.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (71 total, 19 thin omitted)

### Community 0 - "net/http.Request"
Cohesion: 0.09
Nodes (19): WriteError(), configurationVersionResource(), approvalGateStatus(), failedJobStatus(), parseRunID(), stepStatusTFE(), stepStatusValue(), decodeStateAttribute() (+11 more)

### Community 1 - "JSON:API Resource Models"
Cohesion: 0.09
Nodes (23): organizationAttrs, projectAttrs, variableAttrs, workspaceAttrs, S3, stateMetadata, configureProcessGroup(), killProcessGroup() (+15 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "WriteResource"
Cohesion: 0.20
Nodes (17): API, Deps, listQuery, boolValue(), filterSortPage(), New(), param(), parseList() (+9 more)

### Community 4 - "memStorage"
Cohesion: 0.16
Nodes (4): failingBackend, newMemStorage(), logBackend, memStorage

### Community 5 - "Core Database Schema"
Cohesion: 0.12
Nodes (27): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+19 more)

### Community 6 - ".QueryRow"
Cohesion: 0.12
Nodes (5): ApproveJobParams, CreateJobParams, CreateVariableParams, UpdateJobPlanChangesParams, UpdateVariableParams

### Community 7 - ".streamStepLog"
Cohesion: 0.14
Nodes (16): compareStreamIDs(), isTerminal(), outputOf(), parseResumeID(), splitStreamID(), stepOf(), streamKey(), TestCompareStreamIDs() (+8 more)

### Community 8 - "Handler"
Cohesion: 0.14
Nodes (15): Claims, Type, Classify(), TestEffectiveStepIDFromJTI(), claimsFromContext(), isActiveStepStatus(), isExternalUser(), withClaims() (+7 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.12
Nodes (34): TestContractWorkspaceAndRunEnvelopes(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage(), stringPtr(), TestAuth(), TestConfigurationUploadTickets() (+26 more)

### Community 10 - "stubQuerier"
Cohesion: 0.08
Nodes (9): CreateContentParams, UpdateContentStatusParams, Content, Tag, CreateTagParams, GetTagByNameParams, CreateWorkspaceParams, SetWorkspaceLockParams (+1 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.17
Nodes (20): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+12 more)

### Community 12 - "DecodeRequest"
Cohesion: 0.18
Nodes (18): pickBool(), pickString(), pickStringValue(), DecodeAttributes(), DecodeRequest(), RequestData, ResourceIdentifier, HasAttribute() (+10 more)

### Community 13 - "Organization Store Queries"
Cohesion: 0.19
Nodes (4): stubQuerier, Organization, CreateOrganizationParams, UpdateOrganizationParams

### Community 14 - "steps"
Cohesion: 0.06
Nodes (36): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+28 more)

### Community 15 - "Querier"
Cohesion: 0.07
Nodes (9): variableStore, Querier, UpdateJobOutputParams, DBTX, GetArchiveRow, Queries, SetJobOverrideSourceParams, SetWorkspaceLastJobParams (+1 more)

### Community 16 - "New"
Cohesion: 0.11
Nodes (18): New(), contains(), envValue(), TestEnvironmentSeparatesSecrets(), TestJobContextTimesOut(), TestJobContextWithoutTimeoutOnlyCancels(), TestRunCommandDoesNotInheritExecutorSecrets(), TestRunCommandKillsProcessGroup() (+10 more)

### Community 18 - "Project"
Cohesion: 0.18
Nodes (4): Project, CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 19 - "config_test.go"
Cohesion: 0.06
Nodes (45): main(), run(), main(), run(), main(), run(), main(), run() (+37 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 22 - "net/http.Handler"
Cohesion: 0.47
Nodes (3): assertContractJSONAPI(), TestContractOrganizationEnvelope(), fixture

### Community 23 - "Resource"
Cohesion: 0.15
Nodes (17): strValue(), projectResource(), workspaceResource(), Resource, Link(), LinkList(), TestRelationshipMarshalling(), Write() (+9 more)

### Community 24 - "testing.T"
Cohesion: 0.10
Nodes (29): newTestAPI(), TestAuthMiddleware(), TestCreateOrganization(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), TestClassify(), gitOutput() (+21 more)

### Community 25 - "context.Context"
Cohesion: 0.09
Nodes (10): artifactStore, cancelRunStore, historyStub, reliabilityStore, History, CreateHistoryParams, GetHistoryByJobReferenceParams, UpdateJobStatusParams (+2 more)

### Community 26 - "Local Dev Services Script"
Cohesion: 0.47
Nodes (8): container_exists(), down(), ensure_network(), ensure_volume(), dev-services.sh script, status(), up(), wipe()

### Community 27 - "Compose Build Configuration"
Cohesion: 0.25
Nodes (8): build, ports, profiles, SERVICE, args, context, dockerfile, api

### Community 28 - "Compose Environment Settings"
Cohesion: 0.25
Nodes (8): environment, TERRAKUBE_DATABASE_URL, TERRAKUBE_ENV, TERRAKUBE_LOG_LEVEL, TERRAKUBE_PUBLIC_URL, TERRAKUBE_REDIS_URL, TERRAKUBE_STORAGE_LOCAL_PATH, TERRAKUBE_STORAGE_TYPE

### Community 29 - "Postgres Compose Service"
Cohesion: 0.25
Nodes (8): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, environment, image, ports, volumes, postgres

### Community 30 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.25
Nodes (7): Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph

### Community 31 - "Job"
Cohesion: 0.12
Nodes (15): Executor, preparedWorkspace, isCancelledStatus(), isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestBinarySelection() (+7 more)

### Community 33 - "OpenTofu State Conformance"
Cohesion: 0.25
Nodes (6): tofu-state-conformance.sh script, TF_HTTP_PASSWORD, TF_HTTP_USERNAME, TF_IN_AUTOMATION, TF_LOG, TF_LOG_PATH

### Community 34 - "Compose Project Definition"
Cohesion: 0.29
Nodes (6): name, services, volumes, minio-data, postgres-data, redis-data

### Community 35 - "Compose Healthchecks"
Cohesion: 0.53
Nodes (6): interval, retries, test, timeout, healthcheck, healthcheck

### Community 36 - "Compose Service Dependencies"
Cohesion: 0.40
Nodes (5): depends_on, postgres, redis, condition, condition

### Community 37 - "MinIO Compose Service"
Cohesion: 0.40
Nodes (5): command, image, ports, volumes, minio

### Community 38 - "Redis Compose Service"
Cohesion: 0.40
Nodes (5): command, image, ports, volumes, redis

### Community 39 - "M4 run logs and artifacts — design"
Cohesion: 0.25
Nodes (7): Approaches considered, Decisions, Goals, Intent, M4 run logs and artifacts — design, Non-goals, What is left to the builder

### Community 40 - "Management API Smoke Test"
Cohesion: 0.60
Nodes (3): expect_status(), req(), smoke.sh script

### Community 41 - "TFE API Smoke Test"
Cohesion: 0.70
Nodes (4): expect_status(), req(), tfe-smoke.sh script, status()

### Community 51 - "context.Context"
Cohesion: 0.06
Nodes (13): captureLogPublisher, fakeQueue, JobQueue, LogPublisher, logSink, fakePinger, Pinger, envDuration() (+5 more)

### Community 52 - "context.Context"
Cohesion: 0.11
Nodes (5): Step, UpdateStepStatusParams, CreateStepParams, GetStepByJobAndNumberParams, stubQuerier

### Community 53 - "Backend"
Cohesion: 0.16
Nodes (10): Validator, Deps, Deps, Backend, Deps, JobQueue, stubJobQueue, Deps (+2 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 56 - "NewLayout"
Cohesion: 0.18
Nodes (23): TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationRejectsInvalidArchive(), TestArchiveHistoryDoesNotCreateRowWhenUploadFails(), TestCompleteStepArchivesLog(), TestCompleteStepCancelledKeepsStatusWhenLogUploadFails(), TestCompleteStepLogUploadFailureFailsStep(), TestFinishApplyCompletesStep() (+15 more)

### Community 58 - "net/http.Handler"
Cohesion: 0.17
Nodes (14): corsMiddleware(), ExtractToken(), TestExtractToken(), TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput() (+6 more)

### Community 59 - "Workspacetag"
Cohesion: 0.27
Nodes (3): Workspacetag, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams

### Community 60 - ".cloneConfiguration"
Cohesion: 0.22
Nodes (9): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL(), runGit() (+1 more)

### Community 62 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 63 - "query.go"
Cohesion: 0.32
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 65 - "cleanKey"
Cohesion: 0.29
Nodes (5): cleanKey(), pathClean(), TestCleanKey(), NewS3(), S3

### Community 68 - "Discovery"
Cohesion: 0.40
Nodes (3): TestContractDiscoveryAndPing(), Discovery(), TestDiscovery()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 158 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `JSON:API Resource Models`, `WriteResource`, `.GetWorkspaceByName`, `.QueryRow`, `.UpdateWorkspace`, `TFE Handler Tests`, `stubQuerier`, `Organization Store Queries`, `New`, `Project`, `context.Context`, `Backend`, `testing.T`, `context.Context`, `net/http.Handler`, `Workspacetag`, `TempArchive`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.Request` be split into smaller, more focused modules?**
  _Cohesion score 0.0941358024691358 - nodes in this community are weakly interconnected._
- **Why does `Workspace` connect `Job` to `net/http.Request`, `JSON:API Resource Models`, `.GetWorkspaceByName`, `.UpdateWorkspace`, `stubQuerier`, `Querier`, `Resource`, `context.Context`, `.cloneConfiguration`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `JSON:API Resource Models` be split into smaller, more focused modules?**
  _Cohesion score 0.09224137931034483 - nodes in this community are weakly interconnected._
- **Why does `run()` connect `config_test.go` to `net/http.Request`, `JSON:API Resource Models`, `WriteResource`, `Discovery`, `Handler`, `State Backend Tests`, `NewLayout`, `net/http.Handler`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._