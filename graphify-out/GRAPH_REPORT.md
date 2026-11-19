# Graph Report - terrakube-go  (2026-11-19)

## Corpus Check
- 128 files · ~79,668 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 16 file(s) not represented in the graph (top: .graphql 10, (none) 4, .example 1)

## Summary
- 1434 nodes · 5087 edges · 66 communities (54 shown, 12 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 858 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7cd3130b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- JSON:API Resource Models
- Agent Workflows and Compatibility Docs
- Node
- io.Reader
- Core Database Schema
- stubQuerier
- .streamStepLog
- Handler
- TFE Handler Tests
- NewLayout
- State Backend Tests
- History
- Organization Store Queries
- steps
- Querier
- newTestExecutor
- transientError
- .QueryRow
- run
- sqlc Generation Config
- tarGzEntries
- Upstream Elide entities (GraphQL surface)
- .cloneConfiguration
- NewPool
- reliabilityStore
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- context.Context
- Job
- Upstream UI GraphQL operations
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
- newRedactor
- config_test.go
- executor-e2e.sh
- Config
- testLogger
- Tag
- testing.T
- New
- ParsePage
- fixture
- Migrate
- Discovery

## God Nodes (most connected - your core abstractions)
1. `Querier` - 81 edges
2. `API` - 66 edges
3. `Queries` - 64 edges
4. `Job` - 64 edges
5. `stubQuerier` - 58 edges
6. `Workspace` - 56 edges
7. `Executor` - 46 edges
8. `New()` - 45 edges
9. `Step` - 40 edges
10. `WriteError()` - 36 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Upstream Elide entities (GraphQL surface)` --references--> `Template`  [INFERRED]
  docs/upstream/graphql-entities.md → internal/store/models.go
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `Layout`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfstate/layout.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `New()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfoutput/handler.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (66 total, 12 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.05
Nodes (77): API, Deps, listQuery, boolValue(), filterSortPage(), New(), param(), parseList() (+69 more)

### Community 1 - "JSON:API Resource Models"
Cohesion: 0.09
Nodes (23): organizationAttrs, projectAttrs, variableAttrs, workspaceAttrs, S3, stateMetadata, configureProcessGroup(), killProcessGroup() (+15 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "Node"
Cohesion: 0.08
Nodes (28): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), OrderBy(), ParseSort() (+20 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (16): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+8 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.12
Nodes (27): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+19 more)

### Community 6 - "stubQuerier"
Cohesion: 0.08
Nodes (10): CreateContentParams, UpdateContentStatusParams, Content, TempArchive, Workspacetag, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, CreateWorkspaceParams (+2 more)

### Community 7 - ".streamStepLog"
Cohesion: 0.15
Nodes (15): compareStreamIDs(), isTerminal(), outputOf(), parseResumeID(), splitStreamID(), stepOf(), streamKey(), TestCompareStreamIDs() (+7 more)

### Community 8 - "Handler"
Cohesion: 0.08
Nodes (22): Claims, Type, Validator, Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken() (+14 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.13
Nodes (29): TestContractWorkspaceAndRunEnvelopes(), capabilitySecret(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), stringPtr(), TestCapabilitySecretPrefersSharedValues(), TestConfigurationUploadTickets() (+21 more)

### Community 10 - "NewLayout"
Cohesion: 0.12
Nodes (24): corsMiddleware(), TestArchiveHistory(), TestArchiveHistoryWithoutStateIsNoop(), TestWriteBackendOverride(), newTestAPIWithStorage(), TestCreateAndGetStateVersion(), TestCreateGetUploadConfigurationVersion(), TestCreateStateVersionUpdatesCurrentState() (+16 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.15
Nodes (22): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+14 more)

### Community 12 - "History"
Cohesion: 0.16
Nodes (5): artifactStore, historyStub, History, CreateHistoryParams, GetHistoryByJobReferenceParams

### Community 13 - "Organization Store Queries"
Cohesion: 0.19
Nodes (4): stubQuerier, Organization, CreateOrganizationParams, UpdateOrganizationParams

### Community 14 - "steps"
Cohesion: 0.06
Nodes (36): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+28 more)

### Community 15 - "Querier"
Cohesion: 0.06
Nodes (5): Querier, CreateVariableParams, DBTX, Queries, UpdateVariableParams

### Community 16 - "newTestExecutor"
Cohesion: 0.24
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 18 - ".QueryRow"
Cohesion: 0.16
Nodes (5): Project, CreateProjectParams, GetProjectByNameParams, UpdateJobPlanChangesParams, UpdateProjectParams

### Community 19 - "run"
Cohesion: 0.21
Nodes (11): main(), run(), main(), run(), main(), run(), ParseLogLevel(), TestParseLogLevel() (+3 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "tarGzEntries"
Cohesion: 0.13
Nodes (19): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+11 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 23 - ".cloneConfiguration"
Cohesion: 0.22
Nodes (9): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL(), runGit() (+1 more)

### Community 24 - "NewPool"
Cohesion: 0.33
Nodes (8): main(), run(), Locker, NewPool(), NewLocker(), TestLockerAllowsDifferentJobs(), TestLockerExcludesASecondHolder(), testPool()

### Community 25 - "reliabilityStore"
Cohesion: 0.11
Nodes (5): cancelRunStore, reliabilityStore, variableStore, UpdateJobStatusParams, UpdateStepStatusParams

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

### Community 30 - "context.Context"
Cohesion: 0.08
Nodes (6): UpdateJobOutputParams, UpdateStepOutputParams, GetWorkspaceByNameParams, UpdateWorkspaceParams, GetArchiveRow, stubStore

### Community 31 - "Job"
Cohesion: 0.09
Nodes (15): Executor, logSink, preparedWorkspace, isCancelledStatus(), isTerminalStatus(), isTransient(), parseAttempt(), transient() (+7 more)

### Community 32 - "Upstream UI GraphQL operations"
Cohesion: 0.29
Nodes (6): Connection shape, Entity and relationship surface, Filter expressions the UI generates, Other contracts observed, Sort strings, Upstream UI GraphQL operations

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
Cohesion: 0.05
Nodes (17): captureLogPublisher, Deps, fakeLocker, fakeQueue, JobLocker, JobQueue, LogPublisher, Deps (+9 more)

### Community 52 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 53 - "config_test.go"
Cohesion: 0.26
Nodes (13): env(), envBool(), envDuration(), envFirst(), envInt64(), Load(), clearEnv(), TestDatabaseURLFallback() (+5 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Config"
Cohesion: 0.17
Nodes (10): Config, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1` (+2 more)

### Community 56 - "testLogger"
Cohesion: 0.09
Nodes (37): TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestArchiveHistoryDoesNotCreateRowWhenUploadFails(), TestCompleteStepArchivesLog(), TestCompleteStepCancelledKeepsStatusWhenLogUploadFails(), TestCompleteStepLogUploadFailureFailsStep(), TestFinishApplyCompletesStep() (+29 more)

### Community 57 - "Tag"
Cohesion: 0.27
Nodes (3): Tag, CreateTagParams, GetTagByNameParams

### Community 58 - "testing.T"
Cohesion: 0.11
Nodes (27): newTestAPI(), TestAuthMiddleware(), TestCreateOrganization(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), gitOutput(), initGitRepo() (+19 more)

### Community 59 - "New"
Cohesion: 0.22
Nodes (8): New(), requestLogger(), TestHealthz(), TestReadyz(), TestVersion(), writeJSON(), Get(), Info

### Community 60 - "ParsePage"
Cohesion: 0.67
Nodes (3): ParsePage(), TestParsePage(), TestParsePageRejectsBadOffset()

### Community 61 - "fixture"
Cohesion: 0.47
Nodes (3): assertContractJSONAPI(), TestContractOrganizationEnvelope(), fixture

### Community 63 - "Migrate"
Cohesion: 0.47
Nodes (4): appliedVersions(), applyMigration(), Migrate(), path()

### Community 66 - "Discovery"
Cohesion: 0.40
Nodes (3): TestContractDiscoveryAndPing(), Discovery(), TestDiscovery()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 178 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **12 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `net/http.ResponseWriter`, `JSON:API Resource Models`, `.ApproveJob`, `stubQuerier`, `Handler`, `TFE Handler Tests`, `NewLayout`, `History`, `Organization Store Queries`, `newTestExecutor`, `.QueryRow`, `time.Duration`, `.CreateStep`, `reliabilityStore`, `testing.T`, `Job`, `context.Context`, `Tag`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.050114942528735634 - nodes in this community are weakly interconnected._
- **Why does `Template` connect `Upstream Elide entities (GraphQL surface)` to `JSON:API Resource Models`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `JSON:API Resource Models` be split into smaller, more focused modules?**
  _Cohesion score 0.08722245496439045 - nodes in this community are weakly interconnected._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Should `Node` be split into smaller, more focused modules?**
  _Cohesion score 0.08484848484848485 - nodes in this community are weakly interconnected._