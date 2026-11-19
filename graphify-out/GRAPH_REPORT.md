# Graph Report - terrakube-go  (2026-11-19)

## Corpus Check
- 116 files · ~70,757 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 6 file(s) not represented in the graph (top: (none) 4, .example 1, .jsonc 1)

## Summary
- 1337 nodes · 4827 edges · 65 communities (50 shown, 15 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 789 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `60a36431`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.Request
- JSON:API Resource Models
- Agent Workflows and Compatibility Docs
- WriteResource
- memStorage
- Core Database Schema
- Job
- New
- Handler
- TFE Handler Tests
- .QueryRow
- State Backend Tests
- DecodeRequest
- Organization Store Queries
- steps
- Querier
- newTestExecutor
- transientError
- Step
- config.go
- sqlc Generation Config
- tarGzEntries
- New
- jsonapi.go
- initGitRepo
- context.Context
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- Content
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
- cleanKey
- Discovery
- executor-e2e.sh
- testLogger
- Tag
- testing.T
- Workspacetag
- .Get
- newRedactor
- Object
- Local
- Workspace

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
- `Generated Code and Migration Policy` --semantically_similar_to--> `sqlc Workflow`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .opencode/skills/sqlc-workflow/SKILL.md
- `GraphQL UI Requirement` --semantically_similar_to--> `GraphQL and Subscription UI Contract`  [INFERRED] [semantically similar]
  AGENTS.md → docs/compat.md
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `Storage`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/config/config.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (65 total, 15 thin omitted)

### Community 0 - "net/http.Request"
Cohesion: 0.14
Nodes (11): WriteError(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), normalizeTagValue(), tagValuesEqual(), New() (+3 more)

### Community 1 - "JSON:API Resource Models"
Cohesion: 0.08
Nodes (32): organizationAttrs, projectAttrs, variableAttrs, workspaceAttrs, gitConfig, stateMetadata, TestGitSourceDetection(), TestSanitizeGitURL() (+24 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "WriteResource"
Cohesion: 0.19
Nodes (19): API, Deps, listQuery, boolValue(), filterSortPage(), New(), param(), parseList() (+11 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.12
Nodes (27): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+19 more)

### Community 6 - "Job"
Cohesion: 0.09
Nodes (9): Job, TempArchive, UpdateJobStatusParams, ApproveJobParams, CreateArchiveParams, CreateJobParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams (+1 more)

### Community 7 - "New"
Cohesion: 0.06
Nodes (32): Validator, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1` (+24 more)

### Community 8 - "Handler"
Cohesion: 0.13
Nodes (15): Claims, Type, Classify(), claimsFromContext(), isActiveStepStatus(), isExternalUser(), withClaims(), withExternalUser() (+7 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.14
Nodes (28): TestContractWorkspaceAndRunEnvelopes(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), stringPtr(), TestAuth(), TestConfigurationUploadTickets(), TestCreateGetAndDiscardRun() (+20 more)

### Community 10 - ".QueryRow"
Cohesion: 0.20
Nodes (4): Project, CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 11 - "State Backend Tests"
Cohesion: 0.15
Nodes (22): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+14 more)

### Community 12 - "DecodeRequest"
Cohesion: 0.15
Nodes (18): pickBool(), pickString(), pickStringValue(), DecodeAttributes(), DecodeRequest(), RequestData, ResourceIdentifier, HasAttribute() (+10 more)

### Community 13 - "Organization Store Queries"
Cohesion: 0.18
Nodes (6): stubQuerier, Organization, CreateOrganizationParams, organizationPermissions(), organizationResource(), UpdateOrganizationParams

### Community 14 - "steps"
Cohesion: 0.06
Nodes (36): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+28 more)

### Community 15 - "Querier"
Cohesion: 0.07
Nodes (6): Querier, CreateVariableParams, DBTX, Queries, SetWorkspaceLastJobParams, UpdateVariableParams

### Community 16 - "newTestExecutor"
Cohesion: 0.23
Nodes (11): newFakeLocker(), TestHandleReleasesTheJobLock(), TestHandleSkipsAJobHeldByAnotherExecutor(), TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 18 - "Step"
Cohesion: 0.12
Nodes (5): Step, UpdateStepStatusParams, CreateStepParams, GetStepByJobAndNumberParams, stubQuerier

### Community 19 - "config.go"
Cohesion: 0.05
Nodes (53): main(), run(), main(), run(), main(), run(), main(), run() (+45 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 22 - "New"
Cohesion: 0.14
Nodes (13): TestStaleClaimIdleOutlivesTheJobTimeout(), New(), contains(), envValue(), TestEnvironmentSeparatesSecrets(), TestJobContextTimesOut(), TestJobContextWithoutTimeoutOnlyCancels(), TestRunCommandDoesNotInheritExecutorSecrets() (+5 more)

### Community 23 - "jsonapi.go"
Cohesion: 0.10
Nodes (23): Resource, Link(), LinkList(), TestRelationshipMarshalling(), Write(), WriteCollection(), configurationVersionResource(), projectResource() (+15 more)

### Community 24 - "initGitRepo"
Cohesion: 0.60
Nodes (5): gitOutput(), initGitRepo(), requireGit(), TestCloneConfigurationChecksOutCommit(), TestCloneConfigurationFromLocalRepository()

### Community 25 - "context.Context"
Cohesion: 0.08
Nodes (12): artifactStore, cancelRunStore, historyStub, reliabilityStore, variableStore, History, CreateHistoryParams, GetHistoryByJobReferenceParams (+4 more)

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

### Community 30 - "Content"
Cohesion: 0.27
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 31 - "Job"
Cohesion: 0.08
Nodes (13): Executor, logSink, preparedWorkspace, isCancelledStatus(), isTerminalStatus(), isTransient(), parseAttempt(), stringPointer() (+5 more)

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
Cohesion: 0.07
Nodes (13): captureLogPublisher, fakeQueue, JobQueue, LogPublisher, fakePinger, Pinger, defaultRetryDelay(), TestCleanupStaleWorkspaces() (+5 more)

### Community 52 - "cleanKey"
Cohesion: 0.29
Nodes (4): cleanKey(), pathClean(), TestCleanKey(), S3

### Community 53 - "Discovery"
Cohesion: 0.40
Nodes (3): TestContractDiscoveryAndPing(), Discovery(), TestDiscovery()

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 56 - "testLogger"
Cohesion: 0.10
Nodes (35): corsMiddleware(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestArchiveHistoryDoesNotCreateRowWhenUploadFails(), TestCompleteStepArchivesLog(), TestCompleteStepCancelledKeepsStatusWhenLogUploadFails() (+27 more)

### Community 57 - "Tag"
Cohesion: 0.27
Nodes (3): Tag, CreateTagParams, GetTagByNameParams

### Community 58 - "testing.T"
Cohesion: 0.07
Nodes (37): newTestAPI(), TestAuthMiddleware(), TestCreateOrganization(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), TestClassify(), TestEffectiveStepIDFromJTI() (+29 more)

### Community 59 - "Workspacetag"
Cohesion: 0.24
Nodes (3): Workspacetag, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams

### Community 62 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 64 - "Object"
Cohesion: 0.29
Nodes (4): Object, newMemStorage(), memObject, memStorage

### Community 70 - "Workspace"
Cohesion: 0.13
Nodes (6): Workspace, CreateWorkspaceParams, GetWorkspaceByNameParams, SetWorkspaceLockParams, UpdateWorkspaceParams, workspaceWithTags

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 159 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `JSON:API Resource Models`, `WriteResource`, `Job`, `Workspace`, `New`, `TFE Handler Tests`, `.QueryRow`, `Handler`, `Organization Store Queries`, `newTestExecutor`, `Step`, `config.go`, `testLogger`, `context.Context`, `testing.T`, `Workspacetag`, `Content`, `Tag`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.Request` be split into smaller, more focused modules?**
  _Cohesion score 0.1360842844600527 - nodes in this community are weakly interconnected._
- **Why does `Workspace` connect `Workspace` to `net/http.Request`, `JSON:API Resource Models`, `WriteResource`, `Job`, `Querier`, `jsonapi.go`, `context.Context`, `Workspacetag`, `Job`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `JSON:API Resource Models` be split into smaller, more focused modules?**
  _Cohesion score 0.07675791733762748 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.Request` to `JSON:API Resource Models`, `DecodeRequest`, `New`, `jsonapi.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._