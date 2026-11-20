# Graph Report - terrakube-go  (2026-11-20)

## Corpus Check
- 174 files · ~127,849 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 2425 nodes · 9051 edges · 132 communities (79 shown, 53 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 1791 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e745479e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- ._Job
- Agent Workflows and Compatibility Docs
- ParseRSQL
- io.Reader
- Core Database Schema
- stubQuerier
- Handler
- TFE Handler Tests
- WriteResource
- State Backend Tests
- reliabilityStore
- jsonapi.go
- steps
- models_gen.go
- .applyTagBindings
- models_gen.go
- .QueryRow
- sqlc Generation Config
- tarGzEntries
- Upstream Elide entities (GraphQL surface)
- run
- ._Template
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
- tfstate/handler_test.go
- Redis Compose Service
- M4 run logs and artifacts — design
- Management API Smoke Test
- TFE API Smoke Test
- State Backend Smoke Test
- Go Module Definition
- time.Duration
- testing.T
- go_pkg_context
- executor-e2e.sh
- Layout
- testLogger
- .build
- .processJob
- go_pkg_strings
- Querier
- Querier
- executionContext
- dex-config.yaml
- .prepareWorkspace
- Deps
- models
- collections.go
- store/models.go
- WorkspaceResolver
- generated.go
- config.go
- where_test.go
- ._Project
- get
- graphql/models.go
- .streamStepLog
- ._WorkspaceTag
- newTestExecutor
- cancelRunStore
- REST contract, captured from upstream
- pat.py
- idsFrom
- New
- .CreateVariable
- api/main.go
- .listWorkspaces
- stubQuerier
- logSink
- environment
- .EnqueueJob
- streamEvents
- .GetTagByName
- newRedactor
- .Run
- .GetWorkspaceByName
- ui-rest-inventory.sh
- .UpdateOrganization
- ui-smoke.sh
- testDeps
- initGitRepo
- .Ping
- New
- capture.sh
- executor.go
- compare.sh
- .GetProjectByName
- start.sh

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 85 edges
3. `API` - 66 edges
4. `Queries` - 64 edges
5. `stubQuerier` - 48 edges
6. `Executor` - 46 edges
7. `New()` - 45 edges
8. `WriteError()` - 39 edges
9. `New()` - 37 edges
10. `WriteResource()` - 34 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
- `Upstream Elide entities (GraphQL surface)` --references--> `Template`  [INFERRED]
  docs/upstream/graphql-entities.md → internal/store/models.go
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

## Communities (132 total, 53 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.16
Nodes (10): Write(), WriteError(), WriteResource(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), New() (+2 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "ParseRSQL"
Cohesion: 0.10
Nodes (24): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), likePattern(), mustParse() (+16 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (14): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+6 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.09
Nodes (37): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+29 more)

### Community 6 - "stubQuerier"
Cohesion: 0.05
Nodes (9): CreateContentParams, UpdateContentStatusParams, CreateProjectParams, CreateTagParams, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, CreateWorkspaceParams, SetWorkspaceLockParams (+1 more)

### Community 8 - "Handler"
Cohesion: 0.12
Nodes (16): Claims, Type, requireUUIDParams(), Classify(), ExtractToken(), claimsFromContext(), isActiveStepStatus(), isExternalUser() (+8 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.10
Nodes (38): assertContractJSONAPI(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), capabilitySecret(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage() (+30 more)

### Community 10 - "WriteResource"
Cohesion: 0.14
Nodes (27): API, listQuery, boolValue(), filterSortPage(), New(), param(), parseList(), requireType() (+19 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.15
Nodes (22): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+14 more)

### Community 12 - "reliabilityStore"
Cohesion: 0.10
Nodes (6): artifactStore, historyStub, reliabilityStore, CreateHistoryParams, GetHistoryByJobReferenceParams, UpdateStepStatusParams

### Community 13 - "jsonapi.go"
Cohesion: 0.07
Nodes (26): strValue(), variableResource(), workspaceResource(), Resource, Link(), LinkList(), TestRelationshipMarshalling(), configurationVersionResource() (+18 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 15 - "models_gen.go"
Cohesion: 0.12
Nodes (13): JobConnection, NotificationConfigurationConnection, NotificationConfigurationEdge, NotificationTriggerConnection, OrganizationConnection, PageInfo, ProjectAccessConnection, ProjectConnection (+5 more)

### Community 16 - ".applyTagBindings"
Cohesion: 0.15
Nodes (10): RelationshipIdentifiers(), TestRelationshipIdentifiers(), identifierStringAttribute(), normalizeTagValue(), tagValuesEqual(), validateTagBindings(), workspaceTagBindings(), RequestRelationship (+2 more)

### Community 17 - "models_gen.go"
Cohesion: 0.05
Nodes (3): NotificationTriggerEdge, ProjectAccessEdge, WorkspaceEdge

### Community 18 - ".QueryRow"
Cohesion: 0.07
Nodes (3): UpdateWorkspaceParams, UpdateJobPlanChangesParams, UpdateProjectParams

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 23 - "run"
Cohesion: 0.15
Nodes (18): main(), run(), Deps, NewExecutableSchema(), TestSchemaLoads(), bearerMatches(), New(), NewSubscriptions() (+10 more)

### Community 25 - "cancelRunStore"
Cohesion: 0.14
Nodes (6): statusQuerier, UpdateJobStatusParams, Publisher, publishingStore, stubQuerier, mutableQuerier

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

### Community 31 - "Executor"
Cohesion: 0.22
Nodes (4): Executor, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

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

### Community 37 - "tfstate/handler_test.go"
Cohesion: 0.12
Nodes (11): gitConfig, stateMetadata, TestGitSourceDetection(), TestSanitizeGitURL(), gitSource(), isConfigurationTarball(), normalizeGitURL(), sanitizeGitURL() (+3 more)

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
Nodes (47): upstreamCapture, newTestAPI(), TestAuthMiddleware(), TestCreateOrganization(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), TestMalformedIDIsBadRequest() (+39 more)

### Community 53 - "go_pkg_context"
Cohesion: 0.25
Nodes (6): organizationAttrs, projectAttrs, variableAttrs, workspaceAttrs, claimsContextKey, externalUserContextKey

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Layout"
Cohesion: 0.13
Nodes (11): Config, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1` (+3 more)

### Community 56 - "testLogger"
Cohesion: 0.08
Nodes (45): TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestArchiveHistoryDoesNotCreateRowWhenUploadFails(), TestCompleteStepArchivesLog(), TestCompleteStepCancelledKeepsStatusWhenLogUploadFails(), TestCompleteStepLogUploadFailureFailsStep(), TestFinishApplyCompletesStep() (+37 more)

### Community 57 - ".build"
Cohesion: 0.11
Nodes (15): collectionSpec, DB, failingDB, listArgs, afterString(), column(), loadPage(), predicate() (+7 more)

### Community 58 - ".processJob"
Cohesion: 0.16
Nodes (5): isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestIsTerminalStatus()

### Community 59 - "go_pkg_strings"
Cohesion: 0.12
Nodes (11): replay, failingRow, graphQLResponse, ptr(), TestCollectionSpecRejectsInjection(), projectAttrs, runAttributes, stateOutput (+3 more)

### Community 60 - "Querier"
Cohesion: 0.15
Nodes (10): Deps, Validator, Deps, Deps, Backend, Querier, Deps, Deps (+2 more)

### Community 61 - "Querier"
Cohesion: 0.07
Nodes (3): New(), DBTX, Queries

### Community 64 - "dex-config.yaml"
Cohesion: 0.10
Nodes (19): enablePasswordDB, issuer, oauth2, passwordConnector, skipApprovalScreen, staticClients, id, name (+11 more)

### Community 66 - ".prepareWorkspace"
Cohesion: 0.21
Nodes (5): preparedWorkspace, transient(), contentIDFromSource(), gitEnv(), runGit()

### Community 67 - "Deps"
Cohesion: 0.10
Nodes (19): main(), run(), Locker, fakeLocker, JobLocker, NewPool(), NewLocker(), TestLockerAllowsDifferentJobs() (+11 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (10): page, accessEdges(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), projectAccessEdges(), projectEdges(), templateEdges() (+2 more)

### Community 72 - "store/models.go"
Cohesion: 0.17
Nodes (11): Access, Content, History, NotificationConfiguration, NotificationConfigurationTemplate, NotificationTrigger, ProjectAccess, Step (+3 more)

### Community 73 - "WorkspaceResolver"
Cohesion: 0.19
Nodes (4): AccessConnection, AccessEdge, WorkspaceResolver, WorkspaceTagConnection

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (4): ComplexityRoot, DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "config.go"
Cohesion: 0.22
Nodes (17): S3, env(), envBool(), envDuration(), envFirst(), envInt64(), envList(), Load() (+9 more)

### Community 76 - "where_test.go"
Cohesion: 0.12
Nodes (8): ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 78 - "get"
Cohesion: 0.31
Nodes (7): get(), newUI(), TestDeepLinksServeTheIndex(), TestEnvConfigIsNeverCached(), TestMissingAssetsAreNotFound(), TestServesTheIndexAndAssets(), fixture

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - ".streamStepLog"
Cohesion: 0.23
Nodes (10): isTerminal(), outputOf(), stepOf(), streamKey(), TestIsTerminal(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData() (+2 more)

### Community 82 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 83 - "cancelRunStore"
Cohesion: 0.12
Nodes (4): cancelRunStore, variableStore, UpdateJobOutputParams, UpdateStepOutputParams

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 87 - "idsFrom"
Cohesion: 0.30
Nodes (3): NotificationConfigurationResolver, OrganizationResolver, idsFrom()

### Community 88 - "New"
Cohesion: 0.24
Nodes (11): TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput(), TestGetStepLogPrefersArchivedObject(), TestGetStepLogRejectsWrongScope(), TestGetStepLogRequiresToken() (+3 more)

### Community 90 - "api/main.go"
Cohesion: 0.07
Nodes (21): main(), run(), main(), run(), configureProcessGroup(), killProcessGroup(), configureProcessGroup(), killProcessGroup() (+13 more)

### Community 91 - ".listWorkspaces"
Cohesion: 0.29
Nodes (6): parseWorkspaceQuery(), matchesTagFilters(), workspaceTagNames(), writeWorkspaceList(), tagFilter, workspaceQuery

### Community 92 - "stubQuerier"
Cohesion: 0.27
Nodes (3): stubQuerier, demoOrganizations(), CreateOrganizationParams

### Community 93 - "logSink"
Cohesion: 0.20
Nodes (3): captureLogPublisher, LogPublisher, logSink

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 102 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "testDeps"
Cohesion: 0.24
Nodes (9): TestUpstreamUIQueriesExecute(), runQuery(), testDeps(), TestFilterByTagAndJobStatus(), TestRelationshipIDsSelectOneRow(), TestWorkspacePageEndToEnd(), TestContractDiscoveryAndPing(), Discovery() (+1 more)

### Community 121 - "initGitRepo"
Cohesion: 0.53
Nodes (6): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestCloneConfigurationChecksOutCommit(), TestCloneConfigurationFromLocalRepository()

### Community 123 - "New"
Cohesion: 0.33
Nodes (6): NewLocal(), TestLocalHealth(), TestLocalLifecycle(), TestLocalRejectsPathTraversal(), NewS3(), New()

### Community 124 - "capture.sh"
Cohesion: 0.83
Nodes (3): request(), capture.sh script, unauthenticated()

### Community 126 - "executor.go"
Cohesion: 0.19
Nodes (7): transientError, compareStreamIDs(), parseResumeID(), splitStreamID(), TestCompareStreamIDs(), TestParseResumeID(), validStreamID()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 316 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **53 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `executionContext` connect `executionContext` to `executionContext`, `._Job`, `.marshalOBoolean2bool`, `.___EnumValue`, `.___Field`, `generated.go`, `where_test.go`, `._Project`, `models_gen.go`, `._WorkspaceTag`, `context.Context`, `._Template`, `._JobStatusEvent`, `github.com/vektah/gqlparser/v2/ast.SelectionSet`?**
  _High betweenness centrality (0.053) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `Querier` connect `Querier` to `.GetStepByJobAndNumber`, `.GetProjectByName`, `stubQuerier`, `TFE Handler Tests`, `reliabilityStore`, `.QueryRow`, `run`, `cancelRunStore`, `stubStore`, `testing.T`, `go_pkg_context`, `Querier`, `newTestExecutor`, `cancelRunStore`, `New`, `.CreateVariable`, `stubQuerier`, `.GetTagByName`, `.GetWorkspaceByName`, `.UpdateOrganization`, `.ApproveJob`, `.CreateJob`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `ParseRSQL` be split into smaller, more focused modules?**
  _Cohesion score 0.10374149659863946 - nodes in this community are weakly interconnected._
- **Why does `newLogSink()` connect `testLogger` to `newRedactor`, `context.Context`, `go_pkg_context`, `Querier`, `logSink`, `Executor`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `io.Reader` be split into smaller, more focused modules?**
  _Cohesion score 0.06936026936026936 - nodes in this community are weakly interconnected._