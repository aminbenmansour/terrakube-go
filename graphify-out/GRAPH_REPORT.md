# Graph Report - terrakube-go  (2026-11-22)

## Corpus Check
- 238 files · ~144,465 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 2666 nodes · 10150 edges · 138 communities (81 shown, 57 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 1945 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2a4ddeb3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- ._Project
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
- sort.go
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
- Load
- RequestData
- Queries
- Workspace
- dex-config.yaml
- run
- Layout
- models
- github.com/99designs/gqlgen/graphql.Marshaler
- Layout
- collections.go
- tarGzEntries
- testDeps
- generated.go
- where_test.go
- Resolver
- New
- decodeJSON
- graphql/models.go
- Job
- newTestExecutor
- models_gen.go
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- Workspace
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- stubQuerier
- logSink
- where.go
- environment
- streamEvents
- assertSameDocumentShape
- History
- ui-rest-inventory.sh
- Migrate
- ui-smoke.sh
- loadPage
- .CreateArchive
- .Ping
- capture.sh
- 0008_action.sql
- sync.Mutex
- .QueryRow
- compare.sh
- serveGlobalVar
- start.sh
- .PublishCapped
- Queries
- AccessEdge

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 101 edges
3. `Queries` - 75 edges
4. `API` - 66 edges
5. `New()` - 65 edges
6. `WriteError()` - 55 edges
7. `API` - 48 edges
8. `stubQuerier` - 47 edges
9. `Executor` - 46 edges
10. `New()` - 45 edges

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

## Communities (138 total, 57 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.12
Nodes (14): WriteResource(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), normalizeTagValue(), tagValuesEqual(), New() (+6 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "Globalvar"
Cohesion: 0.16
Nodes (5): globalVarStub, variableStore, CreateGlobalVariableParams, UpdateGlobalVariableParams, Globalvar

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (16): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+8 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.08
Nodes (39): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+31 more)

### Community 6 - "stubQuerier"
Cohesion: 0.06
Nodes (7): CreateContentParams, UpdateContentStatusParams, CreateProjectParams, CreateTagParams, CreateWorkspaceTagParams, Content, stubQuerier

### Community 8 - "Handler"
Cohesion: 0.12
Nodes (15): Claims, Type, Classify(), ExtractToken(), claimsFromContext(), isActiveStepStatus(), isExternalUser(), withClaims() (+7 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.10
Nodes (38): assertContractJSONAPI(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), capabilitySecret(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage() (+30 more)

### Community 10 - "WriteResource"
Cohesion: 0.12
Nodes (30): API, Deps, listQuery, stringOr(), boolValue(), filterSortPage(), New(), param() (+22 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.16
Nodes (21): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+13 more)

### Community 13 - "jsonapi.go"
Cohesion: 0.07
Nodes (32): statusQuerier, actionResource(), jobResource(), stepResource(), timestamp(), TestRestrictKeepsOnlyTheRequestedAttributes(), Page, Resource (+24 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 16 - "sort.go"
Cohesion: 0.33
Nodes (6): Order, OrderBy(), ParseSort(), TestOrderByRendersSQL(), TestParseSort(), TestParseSortRejectsUnknownField()

### Community 17 - "executionContext"
Cohesion: 0.07
Nodes (5): executionContext, NotificationConfigurationConnection, NotificationConfigurationEdge, OrganizationConnection, OrganizationEdge

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "page"
Cohesion: 0.15
Nodes (14): collectionSpec, page, accessEdges(), column(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), predicate() (+6 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 23 - "run"
Cohesion: 0.14
Nodes (16): main(), run(), main(), run(), main(), run(), Locker, fakeLocker (+8 more)

### Community 25 - "Job"
Cohesion: 0.21
Nodes (8): gitConfig, transient(), contentIDFromSource(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL(), runGit()

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
Cohesion: 0.16
Nodes (5): Executor, isCancelledStatus(), stringPointer(), TestBinarySelection(), baseEnvironment()

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
Cohesion: 0.16
Nodes (6): NotificationConfigurationResolver, OrganizationResolver, ProjectConnection, QueryResolver, WorkspaceResolver, idsFrom()

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

### Community 51 - "time.Duration"
Cohesion: 0.11
Nodes (7): fakeQueue, JobQueue, defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), Message, Queue

### Community 52 - "testing.T"
Cohesion: 0.05
Nodes (64): newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes(), newTestAPI() (+56 more)

### Community 53 - "api/main.go"
Cohesion: 0.05
Nodes (67): atomicResult, globalVarAttrs, jobAttrs, organizationAttrs, projectAttrs, replay, variableAttrs, workspaceAttrs (+59 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - ".streamStepLog"
Cohesion: 0.33
Nodes (7): isTerminal(), outputOf(), stepOf(), streamKey(), writeData(), writeEvent(), Handler

### Community 56 - "testLogger"
Cohesion: 0.10
Nodes (39): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+31 more)

### Community 58 - ".handle"
Cohesion: 0.16
Nodes (4): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus()

### Community 59 - "Load"
Cohesion: 0.20
Nodes (14): env(), envBool(), envDuration(), envFirst(), envInt64(), envList(), Load(), clearEnv() (+6 more)

### Community 60 - "RequestData"
Cohesion: 0.13
Nodes (23): actionAttrs, atomicOperation, nestedDocument, decodeActionAttrs(), pickBool(), pickString(), pickStringValue(), workspaceLinks() (+15 more)

### Community 61 - "Queries"
Cohesion: 0.08
Nodes (4): historyStub, Querier, SetWorkspaceLastJobParams, stubQuerier

### Community 63 - "Workspace"
Cohesion: 0.09
Nodes (9): cancelRunStore, preparedWorkspace, reliabilityStore, captureJobs(), firstFolder(), Job, Step, Workspace (+1 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "run"
Cohesion: 0.13
Nodes (20): main(), run(), Deps, NewExecutableSchema(), TestSchemaLoads(), bearerMatches(), New(), NewSubscriptions() (+12 more)

### Community 67 - "Layout"
Cohesion: 0.13
Nodes (12): Config, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1` (+4 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 70 - "Layout"
Cohesion: 0.14
Nodes (8): Validator, Deps, Layout, Deps, JobQueue, stubJobQueue, Deps, Deps

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 72 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 73 - "testDeps"
Cohesion: 0.24
Nodes (9): TestUpstreamUIQueriesExecute(), runQuery(), testDeps(), TestFilterByTagAndJobStatus(), TestRelationshipIDsSelectOneRow(), TestWorkspacePageEndToEnd(), TestContractDiscoveryAndPing(), Discovery() (+1 more)

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "Resolver"
Cohesion: 0.12
Nodes (8): ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 77 - "New"
Cohesion: 0.21
Nodes (12): TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput(), TestGetStepLogPrefersArchivedObject(), TestGetStepLogRejectsWrongScope(), TestGetStepLogRequiresToken() (+4 more)

### Community 78 - "decodeJSON"
Cohesion: 0.13
Nodes (23): jobStub, TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound() (+15 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 82 - "newTestExecutor"
Cohesion: 0.16
Nodes (19): newFakeLocker(), TestHandleReleasesTheJobLock(), TestHandleSkipsAJobHeldByAnotherExecutor(), TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestStaleClaimIdleOutlivesTheJobTimeout(), New(), TestAckDeletesTheJobMessage() (+11 more)

### Community 83 - "models_gen.go"
Cohesion: 0.08
Nodes (9): NotificationTriggerConnection, NotificationTriggerEdge, PageInfo, Query, Subscription, TemplateConnection, TemplateEdge, WorkspaceTagConnection (+1 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "workspaces_test.go"
Cohesion: 0.17
Nodes (23): workspaceStub, decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace() (+15 more)

### Community 88 - "Workspace"
Cohesion: 0.13
Nodes (3): CreateWorkspaceParams, GetWorkspaceByNameParams, SetWorkspaceLockParams

### Community 89 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.15
Nodes (6): GetTagByNameParams, UpdateWorkspaceTagValueParams, NotificationConfiguration, NotificationTrigger, Tag, Workspacetag

### Community 92 - "stubQuerier"
Cohesion: 0.07
Nodes (11): actionFilter, actionStub, stubQuerier, parseActionFilter(), demoOrganizations(), CreateActionParams, UpdateActionParams, Action (+3 more)

### Community 93 - "logSink"
Cohesion: 0.21
Nodes (6): logSink, redactor, TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 94 - "where.go"
Cohesion: 0.42
Nodes (3): likePattern(), builder, Compare

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "assertSameDocumentShape"
Cohesion: 0.32
Nodes (7): upstreamCapture, assertSameDocumentShape(), assertSameResourceShape(), loadUpstreamCaptures(), sortedKeys(), TestCapturedShapesAreWhatWeReturn(), TestEveryUpstreamCaptureIsAccountedFor()

### Community 118 - "Migrate"
Cohesion: 0.47
Nodes (4): appliedVersions(), applyMigration(), Migrate(), path()

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 121 - "loadPage"
Cohesion: 0.14
Nodes (9): DB, failingDB, listArgs, afterString(), loadPage(), ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage() (+1 more)

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 126 - "sync.Mutex"
Cohesion: 0.22
Nodes (4): Publisher, publishingStore, stubQuerier, mutableQuerier

### Community 128 - ".QueryRow"
Cohesion: 0.06
Nodes (5): CreateJobParams, CreateStepParams, GetStepByJobAndNumberParams, UpdateOrganizationParams, Variable

### Community 130 - "serveGlobalVar"
Cohesion: 0.39
Nodes (9): newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent(), TestGlobalVariableFromAnotherOrganizationIsNotFound(), TestListGlobalVariablesHidesSensitiveValues() (+1 more)

### Community 133 - "Queries"
Cohesion: 0.09
Nodes (6): New(), ApproveJobParams, DBTX, Queries, SetJobOverrideSourceParams, UpdateJobPlanChangesParams

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 330 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **57 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `executionContext` connect `executionContext` to `context.Context`, `._Project`, `github.com/99designs/gqlgen/graphql.Marshaler`, `AccessEdge`, `.childFields___Type`, `generated.go`, `Resolver`, `StringOrInt`, `.childFields___InputValue`, `Job`, `.marshalOBoolean2bool`, `models_gen.go`, `.___Field`, `.___EnumValue`?**
  _High betweenness centrality (0.046) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.12217697149203999 - nodes in this community are weakly interconnected._
- **Why does `Querier` connect `Queries` to `.QueryRow`, `Globalvar`, `Queries`, `stubQuerier`, `TFE Handler Tests`, `WriteResource`, `artifactStore`, `jsonapi.go`, `Project`, `stubStore`, `testing.T`, `api/main.go`, `nestedStub`, `Workspace`, `run`, `Layout`, `New`, `decodeJSON`, `newTestExecutor`, `workspaces_test.go`, `Workspace`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `stubQuerier`, `History`, `.CreateArchive`, `sync.Mutex`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.ResponseWriter` to `Handler`, `jsonapi.go`, `api/main.go`, `Layout`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `io.Reader` be split into smaller, more focused modules?**
  _Cohesion score 0.06578947368421052 - nodes in this community are weakly interconnected._