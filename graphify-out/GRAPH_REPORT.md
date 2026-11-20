# Graph Report - terrakube-go  (2026-11-20)

## Corpus Check
- 142 files · ~118,612 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: .graphql 10, (none) 4, .example 1)

## Summary
- 2286 nodes · 8722 edges · 118 communities (65 shown, 53 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 1781 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e88f3278`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- Agent Workflows and Compatibility Docs
- ParseRSQL
- io.Reader
- Core Database Schema
- stubQuerier
- DecodeRequest
- Handler
- TFE Handler Tests
- WriteResource
- State Backend Tests
- History
- jsonapi.go
- steps
- Queries
- newTestExecutor
- TemplateConnection
- .QueryRow
- sqlc Generation Config
- tarGzEntries
- Upstream Elide entities (GraphQL surface)
- TestSubscriptionOverGraphQLWS
- github.com/99designs/gqlgen/graphql.Marshaler
- cancelRunStore
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- stubStore
- Executor
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
- logSink
- config_test.go
- executor-e2e.sh
- Layout
- testLogger
- .build
- .processJob
- models_gen.go
- ParsePage
- Queries
- .GetStepByJobAndNumber
- executionContext
- api/main.go
- Deps
- models
- ._Project
- collections.go
- Migrate
- Deps
- generated.go
- .planResource
- .streamStepLog
- graphql/models.go
- net/http.Handler
- WorkspaceTagConnection
- .processJob
- .prepareWorkspace
- idsFrom
- Resolver
- .CreateVariable
- New
- initGitRepo
- stubQuerier
- CreateProjectParams
- .GetTagByName
- streamEvents
- .EnqueueJob
- newRedactor
- .GetWorkspaceByName
- .UpdateOrganization

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 84 edges
3. `API` - 66 edges
4. `Queries` - 64 edges
5. `stubQuerier` - 48 edges
6. `Executor` - 46 edges
7. `New()` - 45 edges
8. `WriteError()` - 37 edges
9. `WriteResource()` - 34 edges
10. `New()` - 29 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `Storage`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/config/config.go
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
- `Upstream Elide entities (GraphQL surface)` --references--> `Template`  [INFERRED]
  docs/upstream/graphql-entities.md → internal/store/models.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `New()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfoutput/handler.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (118 total, 53 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.14
Nodes (13): Write(), WriteError(), WriteResource(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), normalizeTagValue() (+5 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "ParseRSQL"
Cohesion: 0.10
Nodes (24): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), likePattern(), mustParse() (+16 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (15): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+7 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.09
Nodes (37): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+29 more)

### Community 6 - "stubQuerier"
Cohesion: 0.06
Nodes (8): CreateContentParams, UpdateContentStatusParams, CreateTagParams, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, CreateWorkspaceParams, SetWorkspaceLockParams, stubQuerier

### Community 7 - "DecodeRequest"
Cohesion: 0.16
Nodes (17): pickBool(), pickString(), pickStringValue(), DecodeAttributes(), RequestData, ResourceIdentifier, HasAttribute(), RelationshipID() (+9 more)

### Community 8 - "Handler"
Cohesion: 0.12
Nodes (17): Claims, Type, Classify(), TestClassify(), TestEffectiveStepIDFromJTI(), claimsFromContext(), isActiveStepStatus(), isExternalUser() (+9 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.07
Nodes (55): newTestAPI(), TestAuthMiddleware(), TestCreateOrganization(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), TestStateJSONURL(), TestAuthenticationIsRequiredWhenConfigured() (+47 more)

### Community 10 - "WriteResource"
Cohesion: 0.18
Nodes (20): API, Deps, listQuery, boolValue(), filterSortPage(), New(), param(), parseList() (+12 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.15
Nodes (22): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+14 more)

### Community 12 - "History"
Cohesion: 0.08
Nodes (7): artifactStore, historyStub, reliabilityStore, CreateHistoryParams, GetHistoryByJobReferenceParams, UpdateJobStatusParams, UpdateStepStatusParams

### Community 13 - "jsonapi.go"
Cohesion: 0.10
Nodes (22): Resource, Link(), LinkList(), TestRelationshipMarshalling(), configurationVersionResource(), organizationPermissions(), organizationResource(), projectResource() (+14 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (36): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+28 more)

### Community 15 - "Queries"
Cohesion: 0.12
Nodes (3): New(), DBTX, Queries

### Community 16 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 18 - ".QueryRow"
Cohesion: 0.07
Nodes (3): UpdateWorkspaceParams, GetProjectByNameParams, UpdateProjectParams

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 23 - "TestSubscriptionOverGraphQLWS"
Cohesion: 0.53
Nodes (5): dialSubscriptions(), readMessage(), TestSubscriptionOverGraphQLWS(), TestSubscriptionRejectsABadConnectionToken(), writeMessage()

### Community 24 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.05
Nodes (9): AccessConnection, AccessEdge, JobConnection, JobEdge, NotificationTriggerConnection, NotificationTriggerEdge, OrganizationConnection, OrganizationEdge (+1 more)

### Community 25 - "cancelRunStore"
Cohesion: 0.12
Nodes (4): cancelRunStore, variableStore, UpdateJobOutputParams, UpdateStepOutputParams

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

### Community 31 - "Executor"
Cohesion: 0.22
Nodes (4): Executor, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

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
Cohesion: 0.53
Nodes (4): expect_status(), gql_query(), req(), smoke.sh script

### Community 41 - "TFE API Smoke Test"
Cohesion: 0.70
Nodes (4): expect_status(), req(), tfe-smoke.sh script, status()

### Community 51 - "time.Duration"
Cohesion: 0.09
Nodes (11): Deps, fakeQueue, JobQueue, envDuration(), defaultRetryDelay(), Message, Queue, New() (+3 more)

### Community 52 - "logSink"
Cohesion: 0.20
Nodes (3): captureLogPublisher, LogPublisher, logSink

### Community 53 - "config_test.go"
Cohesion: 0.06
Nodes (42): main(), run(), main(), run(), main(), run(), Config, Storage (+34 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Layout"
Cohesion: 0.12
Nodes (11): Validator, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+3 more)

### Community 56 - "testLogger"
Cohesion: 0.08
Nodes (45): TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestArchiveHistoryDoesNotCreateRowWhenUploadFails(), TestCompleteStepArchivesLog(), TestCompleteStepCancelledKeepsStatusWhenLogUploadFails(), TestCompleteStepLogUploadFailureFailsStep(), TestFinishApplyCompletesStep() (+37 more)

### Community 57 - ".build"
Cohesion: 0.12
Nodes (12): collectionSpec, DB, failingDB, listArgs, afterString(), loadPage(), Order, OrderBy() (+4 more)

### Community 58 - ".processJob"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 59 - "models_gen.go"
Cohesion: 0.12
Nodes (10): ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, ProjectConnection, ProjectEdge, Query, QueryResolver, Subscription (+2 more)

### Community 60 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 61 - "Queries"
Cohesion: 0.09
Nodes (5): Querier, publishingStore, stubQuerier, Job, stubQuerier

### Community 64 - "api/main.go"
Cohesion: 0.06
Nodes (62): organizationAttrs, projectAttrs, variableAttrs, workspaceAttrs, S3, gitConfig, stateMetadata, transientError (+54 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): page, accessEdges(), column(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), predicate(), projectAccessEdges() (+4 more)

### Community 72 - "Migrate"
Cohesion: 0.28
Nodes (7): runQuery(), testDeps(), TestFilterByTagAndJobStatus(), TestRelationshipIDsSelectOneRow(), TestWorkspacePageEndToEnd(), Discovery(), TestDiscovery()

### Community 73 - "Deps"
Cohesion: 0.14
Nodes (14): main(), run(), Locker, fakeLocker, JobLocker, NewPool(), NewLocker(), TestLockerAllowsDifferentJobs() (+6 more)

### Community 75 - ".planResource"
Cohesion: 0.18
Nodes (4): approvalGateStatus(), failedJobStatus(), stepStatusTFE(), stepStatusValue()

### Community 78 - ".streamStepLog"
Cohesion: 0.26
Nodes (9): isTerminal(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData(), writeEvent() (+1 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "net/http.Handler"
Cohesion: 0.12
Nodes (17): corsMiddleware(), ExtractToken(), TestExtractToken(), TestWriteBackendOverride(), requestLogger(), New(), newHandler(), stringPtr() (+9 more)

### Community 83 - ".prepareWorkspace"
Cohesion: 0.21
Nodes (5): preparedWorkspace, transient(), contentIDFromSource(), gitEnv(), runGit()

### Community 87 - "idsFrom"
Cohesion: 0.18
Nodes (3): WorkspaceResolver, idsFrom(), workspaceEdges()

### Community 88 - "Resolver"
Cohesion: 0.18
Nodes (6): NotificationConfigurationResolver, OrganizationResolver, ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver

### Community 90 - "New"
Cohesion: 0.25
Nodes (7): New(), TestHealthz(), TestReadyz(), TestVersion(), writeJSON(), Get(), Info

### Community 91 - "initGitRepo"
Cohesion: 0.53
Nodes (6): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestCloneConfigurationChecksOutCommit(), TestCloneConfigurationFromLocalRepository()

### Community 100 - "streamEvents"
Cohesion: 0.14
Nodes (8): EventSource, fakeEvents, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource(), Event, Bus

### Community 103 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

## Knowledge Gaps
- **11 isolated node(s):** `github.com/aminbenmansour/terrakube-go`, `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 275 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **53 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `executionContext` connect `executionContext` to `github.com/99designs/gqlgen/graphql.CollectedField`, `.___EnumValue`, `.childFields___InputValue`, `.marshalOBoolean2bool`, `._Project`, `Resolver`, `generated.go`, `.unmarshalOString2ᚕᚖstring`, `.___Field`, `TemplateConnection`, `WorkspaceTagConnection`, `context.Context`, `.marshalOString2ᚖstring`, `._Job`, `github.com/99designs/gqlgen/graphql.Marshaler`, `models_gen.go`?**
  _High betweenness centrality (0.054) - this node is a cross-community bridge._
- **What connects `github.com/aminbenmansour/terrakube-go`, `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.14279624893435636 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.ResponseWriter` to `api/main.go`, `.planResource`, `jsonapi.go`?**
  _High betweenness centrality (0.041) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `newLogSink()` connect `testLogger` to `api/main.go`, `newRedactor`, `context.Context`, `logSink`, `Executor`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `ParseRSQL` be split into smaller, more focused modules?**
  _Cohesion score 0.10374149659863946 - nodes in this community are weakly interconnected._