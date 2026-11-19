# Graph Report - terrakube-go  (2026-11-19)

## Corpus Check
- 109 files · ~63,991 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 6 file(s) not represented in the graph (top: (none) 4, .example 1, .jsonc 1)

## Summary
- 1260 nodes · 4519 edges · 57 communities (44 shown, 13 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 732 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0f68967b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Management API Handlers
- JSON:API Resource Models
- Agent Workflows and Compatibility Docs
- config.go
- Object Storage Backends
- Core Database Schema
- .QueryRow
- tfoutput/handler_test.go
- Handler
- TFE Handler Tests
- stubQuerier
- State Backend Tests
- Tag
- Organization Store Queries
- steps
- Querier
- newTestExecutor
- transientError
- cancelRunStore
- .cloneConfiguration
- sqlc Generation Config
- tarGzEntries
- Content
- Workspace Tag Store
- newLogSink
- History
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- initGitRepo
- Job
- variableStore
- OpenTofu State Conformance
- Compose Project Definition
- Compose Healthchecks
- Compose Service Dependencies
- MinIO Compose Service
- Redis Compose Service
- TFE Discovery Endpoint
- Management API Smoke Test
- TFE API Smoke Test
- State Backend Smoke Test
- Go Module Definition
- context.Context
- context.Context
- Backend
- executor-e2e.sh
- NewLayout

## God Nodes (most connected - your core abstractions)
1. `Querier` - 80 edges
2. `API` - 66 edges
3. `Queries` - 64 edges
4. `Job` - 61 edges
5. `stubQuerier` - 57 edges
6. `Workspace` - 54 edges
7. `New()` - 44 edges
8. `Executor` - 39 edges
9. `Step` - 37 edges
10. `WriteError()` - 36 edges

## Surprising Connections (you probably didn't know these)
- `Generated Code and Migration Policy` --semantically_similar_to--> `sqlc Workflow`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .opencode/skills/sqlc-workflow/SKILL.md
- `GraphQL UI Requirement` --semantically_similar_to--> `GraphQL and Subscription UI Contract`  [INFERRED] [semantically similar]
  AGENTS.md → docs/compat.md
- `No State LOCK/UNLOCK Endpoint` --semantically_similar_to--> `No LOCK/UNLOCK on /tfstate/v1`  [INFERRED] [semantically similar]
  AGENTS.md → docs/compat.md
- `state.v2 vs /tfstate/v1 Gotcha` --semantically_similar_to--> `Discovery state.v2 vs /tfstate/v1 Mismatch`  [INFERRED] [semantically similar]
  AGENTS.md → docs/compat.md
- `Variable Values Stored Plaintext` --semantically_similar_to--> `Variable Encryption Gap`  [INFERRED] [semantically similar]
  AGENTS.md → docs/compat.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (57 total, 13 thin omitted)

### Community 0 - "Management API Handlers"
Cohesion: 0.05
Nodes (77): API, Deps, listQuery, boolValue(), filterSortPage(), New(), param(), parseList() (+69 more)

### Community 1 - "JSON:API Resource Models"
Cohesion: 0.07
Nodes (40): organizationAttrs, projectAttrs, variableAttrs, workspaceAttrs, main(), run(), main(), run() (+32 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "config.go"
Cohesion: 0.18
Nodes (18): Config, S3, Storage, env(), envBool(), envDuration(), envFirst(), envInt64() (+10 more)

### Community 4 - "Object Storage Backends"
Cohesion: 0.08
Nodes (13): memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3(), Object, newMemStorage() (+5 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.12
Nodes (27): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+19 more)

### Community 6 - ".QueryRow"
Cohesion: 0.13
Nodes (5): GetWorkspaceByNameParams, UpdateWorkspaceParams, ApproveJobParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams

### Community 7 - "tfoutput/handler_test.go"
Cohesion: 0.09
Nodes (29): Validator, corsMiddleware(), ExtractToken(), TestExtractToken(), compareStreamIDs(), isTerminal(), New(), outputOf() (+21 more)

### Community 8 - "Handler"
Cohesion: 0.13
Nodes (17): Claims, Type, Classify(), TestClassify(), TestEffectiveStepIDFromJTI(), claimsFromContext(), isActiveStepStatus(), isExternalUser() (+9 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.09
Nodes (53): newTestAPI(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), TestDecodeRequestRejectsBadJSON(), TestWriteCollectionTotals(), TestWriteError(), TestWriteResource() (+45 more)

### Community 10 - "stubQuerier"
Cohesion: 0.07
Nodes (9): Project, TempArchive, CreateProjectParams, CreateWorkspaceParams, SetWorkspaceLockParams, CreateArchiveParams, GetProjectByNameParams, UpdateProjectParams (+1 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.17
Nodes (20): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+12 more)

### Community 12 - "Tag"
Cohesion: 0.27
Nodes (3): Tag, CreateTagParams, GetTagByNameParams

### Community 13 - "Organization Store Queries"
Cohesion: 0.19
Nodes (4): stubQuerier, Organization, CreateOrganizationParams, UpdateOrganizationParams

### Community 14 - "steps"
Cohesion: 0.06
Nodes (36): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+28 more)

### Community 15 - "Querier"
Cohesion: 0.06
Nodes (10): Querier, UpdateJobOutputParams, CreateJobParams, CreateVariableParams, DBTX, GetArchiveRow, Queries, SetWorkspaceLastJobParams (+2 more)

### Community 16 - "newTestExecutor"
Cohesion: 0.36
Nodes (6): jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob(), TestHandleGivesUpAfterMaxAttempts(), TestHandleLeavesMessagePendingOnShutdown(), TestHandleRetriesTransientFailure()

### Community 18 - "cancelRunStore"
Cohesion: 0.13
Nodes (3): cancelRunStore, UpdateJobStatusParams, UpdateStepOutputParams

### Community 19 - ".cloneConfiguration"
Cohesion: 0.22
Nodes (9): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL(), runGit() (+1 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 22 - "Content"
Cohesion: 0.27
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 23 - "Workspace Tag Store"
Cohesion: 0.27
Nodes (3): Workspacetag, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams

### Community 24 - "newLogSink"
Cohesion: 0.20
Nodes (10): TestLogSinkSplitsLinesAndKeepsTail(), TestLogSinkRedactsSecrets(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), TestRunCommandDoesNotInheritExecutorSecrets(), TestRunCommandKillsProcessGroup(), waitForPID(), waitForProcessExit() (+2 more)

### Community 25 - "History"
Cohesion: 0.15
Nodes (5): historyStub, reliabilityStore, History, CreateHistoryParams, GetHistoryByJobReferenceParams

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

### Community 30 - "initGitRepo"
Cohesion: 0.60
Nodes (5): gitOutput(), initGitRepo(), requireGit(), TestCloneConfigurationChecksOutCommit(), TestCloneConfigurationFromLocalRepository()

### Community 31 - "Job"
Cohesion: 0.11
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

### Community 40 - "Management API Smoke Test"
Cohesion: 0.60
Nodes (3): expect_status(), req(), smoke.sh script

### Community 41 - "TFE API Smoke Test"
Cohesion: 0.70
Nodes (4): expect_status(), req(), tfe-smoke.sh script, status()

### Community 51 - "context.Context"
Cohesion: 0.06
Nodes (15): captureLogPublisher, fakeQueue, JobQueue, LogPublisher, logSink, redactor, fakePinger, Pinger (+7 more)

### Community 52 - "context.Context"
Cohesion: 0.13
Nodes (6): Step, UpdateStepStatusParams, CreateStepParams, GetStepByJobAndNumberParams, stubQuerier, stubStore

### Community 53 - "Backend"
Cohesion: 0.11
Nodes (14): Deps, Deps, New(), requestLogger(), TestHealthz(), TestReadyz(), TestVersion(), writeJSON() (+6 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 56 - "NewLayout"
Cohesion: 0.11
Nodes (31): TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), deref(), TestExecutorEndToEnd(), New(), newMemoryBackend() (+23 more)

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 148 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **13 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `Management API Handlers`, `variableStore`, `JSON:API Resource Models`, `.QueryRow`, `tfoutput/handler_test.go`, `Handler`, `TFE Handler Tests`, `stubQuerier`, `Tag`, `Organization Store Queries`, `newTestExecutor`, `cancelRunStore`, `context.Context`, `Backend`, `Content`, `Workspace Tag Store`, `History`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Management API Handlers` be split into smaller, more focused modules?**
  _Cohesion score 0.050813281355020834 - nodes in this community are weakly interconnected._
- **Why does `Workspace` connect `Job` to `Management API Handlers`, `JSON:API Resource Models`, `.QueryRow`, `stubQuerier`, `Querier`, `cancelRunStore`, `.cloneConfiguration`, `context.Context`, `History`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Should `JSON:API Resource Models` be split into smaller, more focused modules?**
  _Cohesion score 0.07006651884700665 - nodes in this community are weakly interconnected._
- **Why does `API` connect `Management API Handlers` to `JSON:API Resource Models`, `Backend`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._