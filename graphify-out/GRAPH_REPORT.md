# Graph Report - terrakube-go  (2026-11-20)

## Corpus Check
- 144 files · ~119,788 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: .graphql 10, (none) 4, .example 1)

## Summary
- 2302 nodes · 8783 edges · 119 communities (61 shown, 58 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 1784 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3d5001b6`
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
- NotificationTriggerEdge
- newTestExecutor
- models_gen.go
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
- Load
- executor-e2e.sh
- Querier
- testLogger
- loadPage
- .processJob
- ProjectConnection
- ParsePage
- Queries
- .UpdateJobPlanChanges
- api/main.go
- executionContext
- Queue
- models
- ProjectEdge
- collections.go
- github.com/99designs/gqlgen/graphql.CollectedField
- Deps
- newExecutionContext
- .streamStepLog
- graphql/models.go
- New
- ._OrganizationConnection_edges
- .prepareWorkspace
- idsFrom
- Resolver
- .CreateVariable
- run
- M4 run logs and artifacts Implementation Plan
- stubQuerier
- .PublishCapped
- .Run
- CreateTagParams
- streamEvents
- AccessConnection
- .UpdateWorkspace
- .GetWorkspaceByName
- models_gen.go

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
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
- `Upstream Elide entities (GraphQL surface)` --references--> `Template`  [INFERRED]
  docs/upstream/graphql-entities.md → internal/store/models.go
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

## Communities (119 total, 58 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (15): WriteError(), WriteResource(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), normalizeTagValue(), tagValuesEqual() (+7 more)

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
Nodes (8): CreateContentParams, UpdateContentStatusParams, CreateProjectParams, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, CreateWorkspaceParams, SetWorkspaceLockParams, stubQuerier

### Community 7 - "DecodeRequest"
Cohesion: 0.18
Nodes (18): pickBool(), pickString(), pickStringValue(), DecodeAttributes(), DecodeRequest(), RequestData, ResourceIdentifier, HasAttribute() (+10 more)

### Community 8 - "Handler"
Cohesion: 0.14
Nodes (16): Claims, Type, Classify(), TestClassify(), TestEffectiveStepIDFromJTI(), claimsFromContext(), isActiveStepStatus(), isExternalUser() (+8 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.07
Nodes (60): newTestAPI(), TestAuthMiddleware(), TestCreateOrganization(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), TestStateJSONURL(), TestAuthenticationIsRequiredWhenConfigured() (+52 more)

### Community 10 - "WriteResource"
Cohesion: 0.18
Nodes (18): API, listQuery, boolValue(), filterSortPage(), New(), param(), parseList(), requireType() (+10 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.15
Nodes (22): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+14 more)

### Community 12 - "History"
Cohesion: 0.07
Nodes (11): artifactStore, historyStub, reliabilityStore, CreateHistoryParams, GetHistoryByJobReferenceParams, UpdateJobStatusParams, UpdateStepStatusParams, Publisher (+3 more)

### Community 13 - "jsonapi.go"
Cohesion: 0.11
Nodes (19): Resource, Link(), LinkList(), TestRelationshipMarshalling(), Write(), configurationVersionResource(), organizationPermissions(), organizationResource() (+11 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (36): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+28 more)

### Community 16 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 17 - "models_gen.go"
Cohesion: 0.16
Nodes (14): AccessConnection, AccessEdge, JobConnection, JobEdge, NotificationTriggerConnection, OrganizationConnection, PageInfo, ProjectAccessConnection (+6 more)

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
Cohesion: 0.24
Nodes (3): Executor, isCancelledStatus(), TestBinarySelection()

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
Cohesion: 0.10
Nodes (9): fakeQueue, JobQueue, defaultRetryDelay(), Message, Queue, New(), TestDeleteRemovesAMessage(), TestHeadStreamReturnsTheOldestEntries() (+1 more)

### Community 52 - "logSink"
Cohesion: 0.17
Nodes (8): logSink, redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor(), baseEnvironment()

### Community 53 - "Load"
Cohesion: 0.11
Nodes (23): main(), run(), main(), run(), Config, Storage, Task 2: Serve archived logs from `/tfoutput/v1`, env() (+15 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Querier"
Cohesion: 0.15
Nodes (11): Deps, Validator, Deps, Deps, Backend, Querier, Layout, Deps (+3 more)

### Community 56 - "testLogger"
Cohesion: 0.08
Nodes (51): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+43 more)

### Community 57 - "loadPage"
Cohesion: 0.13
Nodes (15): collectionSpec, DB, failingDB, listArgs, afterString(), column(), loadPage(), predicate() (+7 more)

### Community 58 - ".processJob"
Cohesion: 0.16
Nodes (5): isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestIsTerminalStatus()

### Community 59 - "ProjectConnection"
Cohesion: 0.14
Nodes (5): ComplexityRoot, NotificationConfigurationConnection, QueryResolver, WorkspaceConnection, WorkspaceResolver

### Community 60 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 61 - "Queries"
Cohesion: 0.07
Nodes (3): New(), DBTX, Queries

### Community 64 - "api/main.go"
Cohesion: 0.05
Nodes (61): organizationAttrs, projectAttrs, variableAttrs, workspaceAttrs, S3, gitConfig, stateMetadata, transientError (+53 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 71 - "collections.go"
Cohesion: 0.14
Nodes (10): NotificationConfigurationEdge, page, accessEdges(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), projectAccessEdges(), templateEdges() (+2 more)

### Community 72 - "github.com/99designs/gqlgen/graphql.CollectedField"
Cohesion: 0.16
Nodes (12): appliedVersions(), applyMigration(), Migrate(), path(), runQuery(), testDeps(), TestFilterByTagAndJobStatus(), TestRelationshipIDsSelectOneRow() (+4 more)

### Community 73 - "Deps"
Cohesion: 0.13
Nodes (16): main(), run(), Locker, fakeLocker, JobLocker, NewPool(), NewLocker(), TestLockerAllowsDifferentJobs() (+8 more)

### Community 78 - ".streamStepLog"
Cohesion: 0.26
Nodes (9): isTerminal(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData(), writeEvent() (+1 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "New"
Cohesion: 0.16
Nodes (14): ExtractToken(), TestExtractToken(), TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput(), TestGetStepLogPrefersArchivedObject() (+6 more)

### Community 83 - ".prepareWorkspace"
Cohesion: 0.19
Nodes (6): preparedWorkspace, TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv(), runGit()

### Community 87 - "idsFrom"
Cohesion: 0.17
Nodes (3): ProjectConnection, idsFrom(), projectEdges()

### Community 88 - "Resolver"
Cohesion: 0.10
Nodes (7): DirectiveRoot, NotificationConfigurationResolver, OrganizationResolver, ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver

### Community 90 - "run"
Cohesion: 0.08
Nodes (31): main(), run(), Deps, corsMiddleware(), NewExecutableSchema(), TestSchemaLoads(), bearerMatches(), New() (+23 more)

### Community 91 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.25
Nodes (7): Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph

### Community 100 - "streamEvents"
Cohesion: 0.14
Nodes (8): EventSource, fakeEvents, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource(), Event, Bus

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 276 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `executionContext` connect `executionContext` to `.___EnumValue`, `.childFields___InputValue`, `.marshalOBoolean2bool`, `ProjectEdge`, `newExecutionContext`, `.childFields_JobStatusEvent`, `.___Directive`, `executionContext`, `NotificationTriggerEdge`, `._OrganizationConnection_edges`, `._WorkspaceTag`, `.marshalOString2ᚖstring`, `github.com/99designs/gqlgen/graphql.CollectedField`, `._Job`, `Resolver`, `github.com/99designs/gqlgen/graphql.Marshaler`, `ProjectConnection`, `.___Type`?**
  _High betweenness centrality (0.053) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.110990990990991 - nodes in this community are weakly interconnected._
- **Why does `Querier` connect `Querier` to `stubQuerier`, `TFE Handler Tests`, `History`, `newTestExecutor`, `.QueryRow`, `cancelRunStore`, `stubStore`, `Queries`, `.UpdateJobPlanChanges`, `api/main.go`, `Deps`, `New`, `.CreateVariable`, `stubQuerier`, `.ApproveJob`, `CreateTagParams`, `.SetJobOverrideSource`, `.UpdateWorkspace`, `.GetWorkspaceByName`, `models_gen.go`, `.CreateStep`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `newLogSink()` connect `testLogger` to `api/main.go`, `logSink`, `Querier`, `Executor`, `.PublishCapped`, `.___Type`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `ParseRSQL` be split into smaller, more focused modules?**
  _Cohesion score 0.10374149659863946 - nodes in this community are weakly interconnected._