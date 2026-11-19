# Graph Report - terrakube-go  (2026-11-19)

## Corpus Check
- 115 files · ~70,429 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 6 file(s) not represented in the graph (top: (none) 4, .example 1, .jsonc 1)

## Summary
- 1333 nodes · 4799 edges · 73 communities (55 shown, 18 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 784 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `dc7077a7`
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
- context.Context
- State Backend Tests
- .applyTagBindings
- Organization Store Queries
- steps
- Querier
- newTestExecutor
- transientError
- config.go
- config_test.go
- sqlc Generation Config
- tarGzEntries
- net/http.Handler
- Resource
- initGitRepo
- History
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
- time.Duration
- cancelRunStore
- Backend
- executor-e2e.sh
- .planResource
- testLogger
- Tag
- NewLayout
- Workspacetag
- NewLocker
- .Get
- newRedactor
- query.go
- Object
- logSink
- JobQueue
- Local
- Queue
- .PublishCapped
- .UpdateWorkspace
- .Ping

## God Nodes (most connected - your core abstractions)
1. `Querier` - 81 edges
2. `API` - 66 edges
3. `Queries` - 64 edges
4. `Job` - 64 edges
5. `stubQuerier` - 57 edges
6. `Workspace` - 56 edges
7. `Executor` - 45 edges
8. `New()` - 44 edges
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

## Communities (73 total, 18 thin omitted)

### Community 0 - "net/http.Request"
Cohesion: 0.14
Nodes (12): WriteError(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), New(), chiURLParam(), exposedTerraformVersion() (+4 more)

### Community 1 - "JSON:API Resource Models"
Cohesion: 0.08
Nodes (35): organizationAttrs, projectAttrs, variableAttrs, workspaceAttrs, gitConfig, stateMetadata, TestGitSourceDetection(), TestSanitizeGitURL() (+27 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "WriteResource"
Cohesion: 0.14
Nodes (30): API, Deps, listQuery, boolValue(), filterSortPage(), New(), param(), parseList() (+22 more)

### Community 4 - "memStorage"
Cohesion: 0.15
Nodes (4): failingBackend, NewS3(), S3, logBackend

### Community 5 - "Core Database Schema"
Cohesion: 0.12
Nodes (27): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+19 more)

### Community 6 - ".QueryRow"
Cohesion: 0.10
Nodes (7): GetWorkspaceByNameParams, ApproveJobParams, CreateJobParams, CreateStepParams, GetStepByJobAndNumberParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams

### Community 7 - ".streamStepLog"
Cohesion: 0.15
Nodes (15): compareStreamIDs(), isTerminal(), outputOf(), parseResumeID(), splitStreamID(), stepOf(), streamKey(), TestCompareStreamIDs() (+7 more)

### Community 8 - "Handler"
Cohesion: 0.14
Nodes (14): Claims, Type, Classify(), claimsFromContext(), isActiveStepStatus(), isExternalUser(), withClaims(), withExternalUser() (+6 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.09
Nodes (47): TestAuthMiddleware(), TestCreateOrganization(), TestClassify(), TestEffectiveStepIDFromJTI(), TestStateJSONURL(), TestBinarySelection(), TestIsTerminalStatus(), TestDecodeRequestRejectsBadJSON() (+39 more)

### Community 10 - "context.Context"
Cohesion: 0.07
Nodes (13): CreateContentParams, UpdateContentStatusParams, Content, Project, TempArchive, CreateProjectParams, CreateWorkspaceParams, SetWorkspaceLockParams (+5 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.17
Nodes (20): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+12 more)

### Community 12 - ".applyTagBindings"
Cohesion: 0.29
Nodes (6): RelationshipIdentifiers(), normalizeTagValue(), tagValuesEqual(), validateTagBindings(), workspaceTagBindings(), tagBindingInput

### Community 13 - "Organization Store Queries"
Cohesion: 0.19
Nodes (4): stubQuerier, Organization, CreateOrganizationParams, UpdateOrganizationParams

### Community 14 - "steps"
Cohesion: 0.06
Nodes (36): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+28 more)

### Community 15 - "Querier"
Cohesion: 0.06
Nodes (9): Querier, UpdateStepOutputParams, CreateVariableParams, DBTX, GetArchiveRow, Queries, SetWorkspaceLastJobParams, UpdateHistoryOutputParams (+1 more)

### Community 16 - "newTestExecutor"
Cohesion: 0.31
Nodes (8): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob(), TestHandleGivesUpAfterMaxAttempts(), TestHandleLeavesMessagePendingOnShutdown(), TestHandleRetriesTransientFailure()

### Community 18 - "config.go"
Cohesion: 0.17
Nodes (19): Config, S3, Storage, Task 2: Serve archived logs from `/tfoutput/v1`, env(), envBool(), envDuration(), envFirst() (+11 more)

### Community 19 - "config_test.go"
Cohesion: 0.09
Nodes (31): main(), run(), main(), run(), main(), run(), main(), run() (+23 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "tarGzEntries"
Cohesion: 0.22
Nodes (13): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+5 more)

### Community 22 - "net/http.Handler"
Cohesion: 0.16
Nodes (9): corsMiddleware(), newTestAPI(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), ExtractToken(), TestExtractToken(), requestLogger() (+1 more)

### Community 23 - "Resource"
Cohesion: 0.19
Nodes (13): Resource, Link(), Write(), WriteCollection(), configurationVersionResource(), organizationPermissions(), organizationResource(), projectResource() (+5 more)

### Community 24 - "initGitRepo"
Cohesion: 0.53
Nodes (6): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestCloneConfigurationChecksOutCommit(), TestCloneConfigurationFromLocalRepository()

### Community 25 - "History"
Cohesion: 0.10
Nodes (7): artifactStore, historyStub, reliabilityStore, History, CreateHistoryParams, GetHistoryByJobReferenceParams, UpdateStepStatusParams

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
Cohesion: 0.10
Nodes (17): Executor, preparedWorkspace, isCancelledStatus(), isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), transient() (+9 more)

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

### Community 51 - "time.Duration"
Cohesion: 0.31
Nodes (5): fakeQueue, defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), Message

### Community 52 - "cancelRunStore"
Cohesion: 0.12
Nodes (4): cancelRunStore, variableStore, UpdateJobOutputParams, UpdateJobStatusParams

### Community 53 - "Backend"
Cohesion: 0.14
Nodes (10): Validator, Deps, JobLocker, Backend, Layout, Deps, JobQueue, stubJobQueue (+2 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - ".planResource"
Cohesion: 0.18
Nodes (4): approvalGateStatus(), failedJobStatus(), stepStatusTFE(), stepStatusValue()

### Community 56 - "testLogger"
Cohesion: 0.09
Nodes (37): TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestArchiveHistoryDoesNotCreateRowWhenUploadFails(), TestCompleteStepArchivesLog(), TestCompleteStepCancelledKeepsStatusWhenLogUploadFails(), TestCompleteStepLogUploadFailureFailsStep(), TestFinishApplyCompletesStep() (+29 more)

### Community 57 - "Tag"
Cohesion: 0.27
Nodes (3): Tag, CreateTagParams, GetTagByNameParams

### Community 58 - "NewLayout"
Cohesion: 0.12
Nodes (23): TestArchiveHistory(), TestArchiveHistoryWithoutStateIsNoop(), TestWriteBackendOverride(), newTestAPIWithStorage(), TestCreateAndGetStateVersion(), TestCreateGetUploadConfigurationVersion(), TestCreateStateVersionUpdatesCurrentState(), TestCurrentStateVersion() (+15 more)

### Community 59 - "Workspacetag"
Cohesion: 0.27
Nodes (3): Workspacetag, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams

### Community 60 - "NewLocker"
Cohesion: 0.27
Nodes (6): Locker, fakeLocker, NewLocker(), TestLockerAllowsDifferentJobs(), TestLockerExcludesASecondHolder(), testPool()

### Community 62 - "newRedactor"
Cohesion: 0.50
Nodes (4): TestLogSinkRedactsSecrets(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 63 - "query.go"
Cohesion: 0.32
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 64 - "Object"
Cohesion: 0.23
Nodes (5): Deps, Object, newMemStorage(), memObject, memStorage

### Community 65 - "logSink"
Cohesion: 0.29
Nodes (3): logSink, redactor, TestNilRedactorIsSafe()

### Community 67 - "Local"
Cohesion: 0.33
Nodes (4): cleanKey(), pathClean(), TestCleanKey(), Local

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 159 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **18 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `JSON:API Resource Models`, `WriteResource`, `.QueryRow`, `.UpdateWorkspace`, `TFE Handler Tests`, `context.Context`, `Organization Store Queries`, `newTestExecutor`, `cancelRunStore`, `Backend`, `net/http.Handler`, `History`, `NewLayout`, `Workspacetag`, `Job`, `Tag`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.Request` be split into smaller, more focused modules?**
  _Cohesion score 0.14292175486205336 - nodes in this community are weakly interconnected._
- **Why does `Job` connect `Job` to `net/http.Request`, `JSON:API Resource Models`, `.QueryRow`, `.streamStepLog`, `context.Context`, `Querier`, `cancelRunStore`, `.planResource`, `Resource`, `History`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `JSON:API Resource Models` be split into smaller, more focused modules?**
  _Cohesion score 0.07596758312504734 - nodes in this community are weakly interconnected._
- **Why does `Workspace` connect `Job` to `net/http.Request`, `JSON:API Resource Models`, `WriteResource`, `.QueryRow`, `.UpdateWorkspace`, `context.Context`, `.applyTagBindings`, `Querier`, `cancelRunStore`, `Resource`, `History`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._