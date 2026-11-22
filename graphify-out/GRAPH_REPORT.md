# Graph Report - terrakube-go  (2026-11-22)

## Corpus Check
- 225 files · ~140,942 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 2615 nodes · 9911 edges · 144 communities (81 shown, 63 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 1905 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a765e312`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- Agent Workflows and Compatibility Docs
- Globalvar
- io.Reader
- Core Database Schema
- stubQuerier
- Handler
- TFE Handler Tests
- WriteResource
- State Backend Tests
- artifactStore
- jsonapi.go
- steps
- StringOrInt
- .build
- executionContext
- Project
- sqlc Generation Config
- page
- Upstream Elide entities (GraphQL surface)
- run
- Job
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- stubStore
- .runApply
- Upstream UI GraphQL operations
- OpenTofu State Conformance
- Compose Project Definition
- Compose Healthchecks
- Compose Service Dependencies
- idsFrom
- Redis Compose Service
- M4 run logs and artifacts — design
- Management API Smoke Test
- TFE API Smoke Test
- State Backend Smoke Test
- Go Module Definition
- time.Duration
- testing.T
- api/main.go
- executor-e2e.sh
- .streamStepLog
- testLogger
- nestedStub
- .handle
- RequestData
- Queries
- reliabilityStore
- dex-config.yaml
- graphql/handler.go
- Layout
- models
- github.com/99designs/gqlgen/graphql.Marshaler
- Layout
- collections.go
- tarGzEntries
- ui.go
- generated.go
- time.Duration
- Resolver
- New
- decodeJSON
- graphql/models.go
- Job
- newTestExecutor
- NotificationTriggerEdge
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- Workspace
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- stubQuerier
- logSink
- environment
- streamEvents
- Job
- ui-rest-inventory.sh
- Deps
- ui-smoke.sh
- loadPage
- .CreateArchive
- .Ping
- capture.sh
- NotificationConfigurationResolver
- sync.Mutex
- TestSubscriptionSeesAStatusWrittenThroughTheStore
- Step
- compare.sh
- serveGlobalVar
- start.sh
- .PublishCapped
- Queries
- .processJob
- AccessEdge
- .EnqueueJob
- ProjectAccessEdge
- query.go
- newRedactor
- listQuery

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 95 edges
3. `Queries` - 70 edges
4. `API` - 66 edges
5. `New()` - 60 edges
6. `WriteError()` - 52 edges
7. `stubQuerier` - 47 edges
8. `Executor` - 46 edges
9. `New()` - 45 edges
10. `API` - 43 edges

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

## Communities (144 total, 63 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (18): isAtomicRequest(), RelationshipIdentifier(), Write(), WriteCollection(), WriteResource(), configurationVersionResource(), parseWorkspaceQuery(), parseRunID() (+10 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "Globalvar"
Cohesion: 0.16
Nodes (5): globalVarStub, variableStore, CreateGlobalVariableParams, UpdateGlobalVariableParams, Globalvar

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (15): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+7 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.08
Nodes (39): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+31 more)

### Community 6 - "stubQuerier"
Cohesion: 0.05
Nodes (8): CreateContentParams, UpdateContentStatusParams, CreateProjectParams, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, Content, Workspacetag, stubQuerier

### Community 8 - "Handler"
Cohesion: 0.11
Nodes (19): Claims, Type, Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken(), claimsFromContext() (+11 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.09
Nodes (41): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+33 more)

### Community 10 - "WriteResource"
Cohesion: 0.15
Nodes (24): API, Deps, boolValue(), filterSortPage(), New(), param(), parseInclude(), parseList() (+16 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.16
Nodes (21): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+13 more)

### Community 13 - "jsonapi.go"
Cohesion: 0.07
Nodes (30): jobStub, statusQuerier, globalVarResource(), jobResource(), stepResource(), projectResource(), TestRestrictKeepsOnlyTheRequestedAttributes(), TestRestrictLeavesDocumentsWithoutFieldsetsAlone() (+22 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 15 - "StringOrInt"
Cohesion: 0.12
Nodes (13): AccessConnection, JobConnection, NotificationConfigurationConnection, OrganizationConnection, PageInfo, ProjectAccessConnection, ProjectConnection, ProjectEdge (+5 more)

### Community 16 - ".build"
Cohesion: 0.16
Nodes (12): listArgs, afterString(), ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset(), Order, OrderBy() (+4 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "page"
Cohesion: 0.16
Nodes (13): collectionSpec, page, accessEdges(), column(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), predicate() (+5 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 23 - "run"
Cohesion: 0.05
Nodes (53): main(), run(), main(), run(), main(), run(), main(), run() (+45 more)

### Community 25 - "Job"
Cohesion: 0.18
Nodes (9): gitConfig, TestGitSourceDetection(), transient(), contentIDFromSource(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL() (+1 more)

### Community 26 - "Local Dev Services Script"
Cohesion: 0.47
Nodes (8): container_exists(), down(), ensure_network(), ensure_volume(), dev-services.sh script, status(), up(), wipe()

### Community 27 - "Compose Build Configuration"
Cohesion: 0.25
Nodes (8): build, ports, profiles, SERVICE, args, context, dockerfile, api

### Community 28 - "Compose Environment Settings"
Cohesion: 0.22
Nodes (9): environment, TERRAKUBE_CORS_ORIGINS, TERRAKUBE_DATABASE_URL, TERRAKUBE_ENV, TERRAKUBE_LOG_LEVEL, TERRAKUBE_PUBLIC_URL, TERRAKUBE_REDIS_URL, TERRAKUBE_STORAGE_LOCAL_PATH (+1 more)

### Community 29 - "Postgres Compose Service"
Cohesion: 0.25
Nodes (8): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, environment, image, ports, volumes, postgres

### Community 30 - "stubStore"
Cohesion: 0.18
Nodes (3): GetArchiveRow, UpdateHistoryOutputParams, stubStore

### Community 31 - ".runApply"
Cohesion: 0.21
Nodes (3): Executor, isCancelledStatus(), TestBinarySelection()

### Community 32 - "Upstream UI GraphQL operations"
Cohesion: 0.17
Nodes (10): Connection shape, Entity and relationship surface, Filter expressions the UI generates, Other contracts observed, Sort strings, Upstream UI GraphQL operations, Measured route surface, Order that unblocks the most UI (+2 more)

### Community 33 - "OpenTofu State Conformance"
Cohesion: 0.25
Nodes (6): tofu-state-conformance.sh script, TF_HTTP_PASSWORD, TF_HTTP_USERNAME, TF_IN_AUTOMATION, TF_LOG, TF_LOG_PATH

### Community 34 - "Compose Project Definition"
Cohesion: 0.17
Nodes (11): command, image, ports, volumes, name, services, minio, volumes (+3 more)

### Community 35 - "Compose Healthchecks"
Cohesion: 0.53
Nodes (6): interval, retries, test, timeout, healthcheck, healthcheck

### Community 36 - "Compose Service Dependencies"
Cohesion: 0.40
Nodes (5): depends_on, postgres, redis, condition, condition

### Community 37 - "idsFrom"
Cohesion: 0.20
Nodes (3): WorkspaceResolver, idsFrom(), projectEdges()

### Community 38 - "Redis Compose Service"
Cohesion: 0.40
Nodes (5): command, image, ports, volumes, redis

### Community 39 - "M4 run logs and artifacts — design"
Cohesion: 0.25
Nodes (7): Approaches considered, Decisions, Goals, Intent, M4 run logs and artifacts — design, Non-goals, What is left to the builder

### Community 40 - "Management API Smoke Test"
Cohesion: 0.43
Nodes (4): expect_status(), gql_query(), req(), smoke.sh script

### Community 41 - "TFE API Smoke Test"
Cohesion: 0.70
Nodes (4): expect_status(), req(), tfe-smoke.sh script, status()

### Community 52 - "testing.T"
Cohesion: 0.05
Nodes (61): upstreamCapture, newTestAPI(), TestAuthMiddleware(), TestCreateOrganization(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), TestMalformedIDIsBadRequest() (+53 more)

### Community 53 - "api/main.go"
Cohesion: 0.06
Nodes (54): atomicResult, globalVarAttrs, jobAttrs, organizationAttrs, projectAttrs, replay, variableAttrs, workspaceAttrs (+46 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - ".streamStepLog"
Cohesion: 0.26
Nodes (9): isTerminal(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData(), writeEvent() (+1 more)

### Community 56 - "testLogger"
Cohesion: 0.10
Nodes (36): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+28 more)

### Community 58 - ".handle"
Cohesion: 0.22
Nodes (4): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus()

### Community 60 - "RequestData"
Cohesion: 0.15
Nodes (20): atomicOperation, nestedDocument, pickBool(), pickString(), pickStringValue(), workspaceLinks(), DecodeAttributes(), RequestData (+12 more)

### Community 61 - "Queries"
Cohesion: 0.07
Nodes (4): historyStub, Querier, SetWorkspaceLastJobParams, stubQuerier

### Community 63 - "reliabilityStore"
Cohesion: 0.15
Nodes (5): preparedWorkspace, reliabilityStore, captureJobs(), Job, UpdateJobStatusParams

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "graphql/handler.go"
Cohesion: 0.21
Nodes (7): Deps, EventSource, NewExecutableSchema(), TestSchemaLoads(), bearerMatches(), New(), NewSubscriptions()

### Community 67 - "Layout"
Cohesion: 0.25
Nodes (7): Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.05
Nodes (3): OrganizationEdge, TemplateEdge, WorkspaceTagEdge

### Community 70 - "Layout"
Cohesion: 0.24
Nodes (3): Validator, Layout, Deps

### Community 71 - "collections.go"
Cohesion: 0.10
Nodes (24): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), likePattern(), mustParse() (+16 more)

### Community 72 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 73 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (4): ComplexityRoot, DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "time.Duration"
Cohesion: 0.31
Nodes (5): fakeQueue, defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), Message

### Community 76 - "Resolver"
Cohesion: 0.13
Nodes (10): OrganizationResolver, ProjectResolver, QueryResolver, Resolver, ResolverRoot, SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient() (+2 more)

### Community 77 - "New"
Cohesion: 0.21
Nodes (12): TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput(), TestGetStepLogPrefersArchivedObject(), TestGetStepLogRejectsWrongScope(), TestGetStepLogRequiresToken() (+4 more)

### Community 78 - "decodeJSON"
Cohesion: 0.16
Nodes (22): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath() (+14 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 82 - "newTestExecutor"
Cohesion: 0.15
Nodes (19): newFakeLocker(), TestHandleReleasesTheJobLock(), TestHandleSkipsAJobHeldByAnotherExecutor(), TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestStaleClaimIdleOutlivesTheJobTimeout(), New(), TestAckDeletesTheJobMessage() (+11 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "workspaces_test.go"
Cohesion: 0.23
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 89 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.24
Nodes (3): CreateTagParams, GetTagByNameParams, Tag

### Community 92 - "stubQuerier"
Cohesion: 0.07
Nodes (12): stubQuerier, workspaceStub, cancelRunStore, captureWorkspaces(), demoOrganizations(), workspaceUpdateParams(), Organization, Project (+4 more)

### Community 93 - "logSink"
Cohesion: 0.20
Nodes (4): logSink, redactor, TestNilRedactorIsSafe(), baseEnvironment()

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.33
Nodes (3): fakeEvents, Event, Bus

### Community 103 - "Job"
Cohesion: 0.08
Nodes (4): ApproveJobParams, CreateJobParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams

### Community 118 - "Deps"
Cohesion: 0.33
Nodes (3): Deps, fakeLocker, JobLocker

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 121 - "loadPage"
Cohesion: 0.24
Nodes (3): DB, failingDB, loadPage()

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 125 - "NotificationConfigurationResolver"
Cohesion: 0.29
Nodes (3): NotificationConfigurationResolver, NotificationTriggerConnection, TemplateConnection

### Community 126 - "sync.Mutex"
Cohesion: 0.25
Nodes (4): Publisher, publishingStore, stubQuerier, mutableQuerier

### Community 127 - "TestSubscriptionSeesAStatusWrittenThroughTheStore"
Cohesion: 0.36
Nodes (8): dialSubscriptions(), readMessage(), TestSubscriptionOverGraphQLWS(), TestSubscriptionRejectsABadConnectionToken(), TestSubscriptionSeesAStatusWrittenThroughTheStore(), writeMessage(), NewBus(), TestBusRoundTrip()

### Community 128 - "Step"
Cohesion: 0.12
Nodes (3): Step, CreateStepParams, GetStepByJobAndNumberParams

### Community 130 - "serveGlobalVar"
Cohesion: 0.39
Nodes (9): newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent(), TestGlobalVariableFromAnotherOrganizationIsNotFound(), TestListGlobalVariablesHidesSensitiveValues() (+1 more)

### Community 133 - "Queries"
Cohesion: 0.09
Nodes (4): New(), DBTX, Queries, UpdateOrganizationParams

### Community 140 - "query.go"
Cohesion: 0.33
Nodes (4): matchesTagFilters(), writeWorkspaceList(), tagFilter, workspaceQuery

### Community 141 - "newRedactor"
Cohesion: 0.50
Nodes (4): TestLogSinkRedactsSecrets(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 328 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **63 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `executionContext` connect `executionContext` to `github.com/99designs/gqlgen/graphql.CollectedField`, `context.Context`, `github.com/99designs/gqlgen/graphql.Marshaler`, `AccessEdge`, `.childFields___Type`, `generated.go`, `ProjectAccessEdge`, `Resolver`, `.childFields___InputValue`, `Job`, `.marshalOBoolean2bool`, `NotificationTriggerEdge`, `.___Field`, `._Organization`, `._JobStatusEvent`, `.___Schema`?**
  _High betweenness centrality (0.047) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.10981012658227848 - nodes in this community are weakly interconnected._
- **Why does `Querier` connect `Queries` to `Step`, `net/http.ResponseWriter`, `Globalvar`, `Queries`, `stubQuerier`, `Handler`, `TFE Handler Tests`, `WriteResource`, `artifactStore`, `jsonapi.go`, `.UpdateVariable`, `Project`, `run`, `stubStore`, `testing.T`, `api/main.go`, `nestedStub`, `reliabilityStore`, `Layout`, `New`, `newTestExecutor`, `Workspace`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `stubQuerier`, `Job`, `Deps`, `.CreateArchive`, `sync.Mutex`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.10304789550072568 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.ResponseWriter` to `Handler`, `jsonapi.go`, `api/main.go`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._