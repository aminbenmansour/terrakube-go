# Graph Report - terrakube-go  (2026-11-22)

## Corpus Check
- 277 files · ~154,316 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 2858 nodes · 11044 edges · 139 communities (81 shown, 58 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 2134 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `fc266810`
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
- collectionStub
- Handler
- TFE Handler Tests
- WriteResource
- State Backend Tests
- artifactStore
- Resource
- steps
- StringOrInt
- executionContext
- Ssh
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- Workspace
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
- Collection
- testLogger
- nestedStub
- .handle
- net/http.Handler
- WriteError
- Querier
- Workspace
- dex-config.yaml
- run
- Layout
- models
- github.com/99designs/gqlgen/graphql.Marshaler
- .EnqueueJob
- collections.go
- tarGzEntries
- testDeps
- newExecutionContext
- where_test.go
- Resolver
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- decodeJSON
- graphql/models.go
- ui.go
- JobConnection
- newTestExecutor
- models_gen.go
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- Workspace
- Tag
- stubQuerier
- logSink
- where.go
- environment
- streamEvents
- graphql/handler.go
- Wrap
- ui-rest-inventory.sh
- ui-smoke.sh
- loadPage
- .CreateArchive
- .Ping
- capture.sh
- 0008_action.sql
- newRedactor
- Step
- compare.sh
- workspaceTagBindings
- start.sh
- Queries
- AccessEdge

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 129 edges
3. `Queries` - 99 edges
4. `New()` - 87 edges
5. `API` - 78 edges
6. `WriteError()` - 72 edges
7. `API` - 66 edges
8. `stubQuerier` - 64 edges
9. `stubQuerier` - 47 edges
10. `Executor` - 46 edges

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

## Communities (139 total, 58 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (16): WriteResource(), configurationVersionResource(), parseWorkspaceQuery(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), normalizeTagValue() (+8 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "Globalvar"
Cohesion: 0.14
Nodes (5): globalVarStub, variableStore, CreateGlobalVariableParams, UpdateGlobalVariableParams, Globalvar

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (15): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+7 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (50): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+42 more)

### Community 6 - "stubQuerier"
Cohesion: 0.05
Nodes (9): CreateContentParams, UpdateContentStatusParams, CreateProjectParams, CreateTagParams, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, Content, Workspacetag (+1 more)

### Community 7 - "collectionStub"
Cohesion: 0.05
Nodes (12): collectionStub, stubQuerier, workspaceStub, demoOrganizations(), CreateItemParams, CreateReferenceParams, UpdateItemParams, Item (+4 more)

### Community 8 - "Handler"
Cohesion: 0.12
Nodes (18): Claims, Type, Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken(), claimsFromContext() (+10 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.10
Nodes (38): assertContractJSONAPI(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), capabilitySecret(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage() (+30 more)

### Community 10 - "WriteResource"
Cohesion: 0.16
Nodes (11): API, listQuery, filterSortPage(), New(), param(), parseList(), unknownCollection(), projectResource() (+3 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.16
Nodes (21): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+13 more)

### Community 12 - "artifactStore"
Cohesion: 0.18
Nodes (4): Global Constraints, M4 run logs and artifacts Implementation Plan, Self-review notes, artifactStore

### Community 13 - "Resource"
Cohesion: 0.09
Nodes (33): actionResource(), parseInclude(), strValue(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), jobResource() (+25 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 17 - "executionContext"
Cohesion: 0.07
Nodes (5): executionContext, NotificationConfigurationConnection, NotificationConfigurationEdge, OrganizationConnection, OrganizationEdge

### Community 18 - "Ssh"
Cohesion: 0.10
Nodes (6): sshAgentStub, Agent, Ssh, CreateAgentParams, CreateSSHKeyParams, UpdateSSHKeyParams

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.14
Nodes (19): replay, upstreamCapture, sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), assertSameDocumentShape(), assertSameResourceShape() (+11 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 25 - "Workspace"
Cohesion: 0.13
Nodes (11): gitConfig, preparedWorkspace, captureWorkspaces(), transient(), contentIDFromSource(), gitEnv(), gitSource(), isConfigurationTarball() (+3 more)

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
Cohesion: 0.15
Nodes (6): Executor, isCancelledStatus(), isTerminalStatus(), stringPointer(), TestBinarySelection(), TestIsTerminalStatus()

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
Cohesion: 0.10
Nodes (12): page, ProjectConnection, accessEdges(), idsFrom(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), projectAccessEdges() (+4 more)

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
Cohesion: 0.12
Nodes (5): fakeQueue, JobQueue, defaultRetryDelay(), Message, Queue

### Community 52 - "testing.T"
Cohesion: 0.06
Nodes (60): newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes(), newTestAPI() (+52 more)

### Community 53 - "api/main.go"
Cohesion: 0.06
Nodes (58): agentAttrs, atomicResult, collectionAttrs, globalVarAttrs, jobAttrs, organizationAttrs, projectAttrs, sshAttrs (+50 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Collection"
Cohesion: 0.17
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 56 - "testLogger"
Cohesion: 0.08
Nodes (47): fakeLocker, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+39 more)

### Community 58 - ".handle"
Cohesion: 0.18
Nodes (4): isTransient(), parseAttempt(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 59 - "net/http.Handler"
Cohesion: 0.11
Nodes (23): requireUUIDParams(), TestWriteBackendOverride(), CORS(), TestCORSAllowsAnyOriginByDefault(), TestCORSAllowsOnlyListedOrigins(), TestCORSPreflightDoesNotReachTheHandler(), New(), newHandler() (+15 more)

### Community 60 - "WriteError"
Cohesion: 0.09
Nodes (36): actionAttrs, atomicOperation, itemAttrs, nestedDocument, decodeActionAttrs(), stringOr(), boolValue(), requireType() (+28 more)

### Community 61 - "Querier"
Cohesion: 0.05
Nodes (13): Deps, Validator, Deps, historyStub, Deps, Backend, Querier, Layout (+5 more)

### Community 63 - "Workspace"
Cohesion: 0.07
Nodes (9): cancelRunStore, reliabilityStore, captureJobs(), TestWatchJobStopsWhenJobCompletes(), Job, UpdateJobStatusParams, publishingStore, stubQuerier (+1 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "run"
Cohesion: 0.05
Nodes (54): main(), run(), main(), run(), main(), run(), main(), run() (+46 more)

### Community 67 - "Layout"
Cohesion: 0.50
Nodes (4): Review Focus, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 71 - "collections.go"
Cohesion: 0.27
Nodes (8): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), Node, parser

### Community 72 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 73 - "testDeps"
Cohesion: 0.24
Nodes (9): TestUpstreamUIQueriesExecute(), runQuery(), testDeps(), TestFilterByTagAndJobStatus(), TestRelationshipIDsSelectOneRow(), TestWorkspacePageEndToEnd(), TestContractDiscoveryAndPing(), Discovery() (+1 more)

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "Resolver"
Cohesion: 0.09
Nodes (13): DirectiveRoot, NotificationConfigurationResolver, OrganizationResolver, ProjectResolver, QueryResolver, Resolver, ResolverRoot, SubscriptionResolver (+5 more)

### Community 77 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.20
Nodes (6): jobStub, statusQuerier, approvalGateStatus(), failedJobStatus(), stepStatusTFE(), stepStatusValue()

### Community 78 - "decodeJSON"
Cohesion: 0.16
Nodes (22): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath() (+14 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 82 - "newTestExecutor"
Cohesion: 0.31
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

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
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 88 - "Workspace"
Cohesion: 0.10
Nodes (4): CreateWorkspaceParams, UpdateWorkspaceParams, GetWorkspaceByNameParams, SetWorkspaceLockParams

### Community 89 - "Tag"
Cohesion: 0.24
Nodes (4): GetTagByNameParams, NotificationConfiguration, NotificationTrigger, Tag

### Community 92 - "stubQuerier"
Cohesion: 0.14
Nodes (6): actionFilter, actionStub, parseActionFilter(), CreateActionParams, UpdateActionParams, Action

### Community 93 - "logSink"
Cohesion: 0.17
Nodes (4): captureLogPublisher, LogPublisher, logSink, baseEnvironment()

### Community 94 - "where.go"
Cohesion: 0.19
Nodes (6): likePattern(), And, builder, Compare, Op, Or

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.33
Nodes (3): fakeEvents, Event, Bus

### Community 101 - "graphql/handler.go"
Cohesion: 0.21
Nodes (7): Deps, EventSource, NewExecutableSchema(), TestSchemaLoads(), bearerMatches(), New(), NewSubscriptions()

### Community 103 - "Wrap"
Cohesion: 0.24
Nodes (11): dialSubscriptions(), readMessage(), TestSubscriptionOverGraphQLWS(), TestSubscriptionRejectsABadConnectionToken(), TestSubscriptionSeesAStatusWrittenThroughTheStore(), writeMessage(), TestWrapDoesNotPublishAFailedUpdate(), TestWrapIgnoresAPublishFailure() (+3 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 121 - "loadPage"
Cohesion: 0.11
Nodes (19): collectionSpec, DB, failingDB, listArgs, afterString(), column(), loadPage(), predicate() (+11 more)

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 127 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 128 - "Step"
Cohesion: 0.11
Nodes (4): captureSteps(), Step, CreateStepParams, GetStepByJobAndNumberParams

### Community 130 - "workspaceTagBindings"
Cohesion: 0.67
Nodes (4): identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), tagBindingInput

### Community 133 - "Queries"
Cohesion: 0.04
Nodes (11): SetWorkspaceLastJobParams, ApproveJobParams, CreateJobParams, GetProjectByNameParams, History, Queries, SetJobOverrideSourceParams, UpdateJobPlanChangesParams (+3 more)

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 336 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `Step`, `Globalvar`, `Queries`, `stubQuerier`, `collectionStub`, `TFE Handler Tests`, `artifactStore`, `Ssh`, `strPtr`, `stubStore`, `testing.T`, `api/main.go`, `Collection`, `nestedStub`, `net/http.Handler`, `Workspace`, `github.com/aminbenmansour/terrakube-go/internal/store.Job`, `newTestExecutor`, `Workspace`, `Tag`, `stubQuerier`, `Wrap`, `.CreateArchive`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.09588607594936709 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `executionContext` to `context.Context`, `._Project`, `github.com/99designs/gqlgen/graphql.Marshaler`, `AccessEdge`, `newExecutionContext`, `Resolver`, `StringOrInt`, `.childFields___Type`, `JobConnection`, `.marshalOBoolean2bool`, `models_gen.go`, `.childFields___InputValue`, `.___Field`, `.___EnumValue`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `newLogSink()` connect `testLogger` to `context.Context`, `Querier`, `api/main.go`, `.runApply`, `logSink`, `newRedactor`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `Globalvar` be split into smaller, more focused modules?**
  _Cohesion score 0.13666666666666666 - nodes in this community are weakly interconnected._