# Graph Report - terrakube-go  (2026-11-21)

## Corpus Check
- 208 files · ~137,818 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 2588 nodes · 9801 edges · 127 communities (75 shown, 52 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 1887 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `41f20c6c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- context.Context
- Agent Workflows and Compatibility Docs
- stubQuerier
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
- ._NotificationConfiguration
- WorkspaceResolver
- Project
- sqlc Generation Config
- page
- Upstream Elide entities (GraphQL surface)
- run
- TemplateConnection
- Job
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
- .processJob
- RequestData
- Queries
- OrganizationConnection
- ._WorkspaceTag
- dex-config.yaml
- Layout
- models
- github.com/99designs/gqlgen/graphql.Marshaler
- Layout
- collections.go
- tarGzEntries
- ui.go
- generated.go
- config.go
- Resolver
- New
- serveNested
- graphql/models.go
- serveNested
- serveGlobalVar
- newTestExecutor
- query.go
- REST contract, captured from upstream
- pat.py
- New
- .GetWorkspaceByName
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- stubQuerier
- logSink
- environment
- streamEvents
- .Query
- ui-rest-inventory.sh
- ui-smoke.sh
- .CreateArchive
- .Ping
- capture.sh
- .QueryRow
- compare.sh
- start.sh

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 95 edges
3. `Queries` - 70 edges
4. `API` - 66 edges
5. `New()` - 58 edges
6. `stubQuerier` - 48 edges
7. `WriteError()` - 47 edges
8. `Executor` - 46 edges
9. `New()` - 45 edges
10. `API` - 39 edges

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

## Communities (127 total, 52 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (19): WriteResource(), configurationVersionResource(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), identifierStringAttribute(), normalizeTagValue() (+11 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "stubQuerier"
Cohesion: 0.09
Nodes (9): globalVarStub, stubQuerier, demoOrganizations(), CreateGlobalVariableParams, UpdateGlobalVariableParams, Globalvar, Organization, Project (+1 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (14): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+6 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.08
Nodes (39): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+31 more)

### Community 6 - "stubQuerier"
Cohesion: 0.05
Nodes (10): CreateContentParams, UpdateContentStatusParams, CreateProjectParams, CreateTagParams, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, SetWorkspaceLockParams, Content (+2 more)

### Community 8 - "Handler"
Cohesion: 0.13
Nodes (17): Claims, Type, requireUUIDParams(), Classify(), ExtractToken(), claimsFromContext(), isActiveStepStatus(), isExternalUser() (+9 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.08
Nodes (41): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+33 more)

### Community 10 - "WriteResource"
Cohesion: 0.15
Nodes (27): API, Deps, boolValue(), filterSortPage(), New(), param(), parseInclude(), parseList() (+19 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.16
Nodes (21): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+13 more)

### Community 13 - "jsonapi.go"
Cohesion: 0.08
Nodes (25): listQuery, statusQuerier, jobResource(), stepResource(), TestRestrictKeepsOnlyTheRequestedAttributes(), Resource, Link(), LinkList() (+17 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 15 - "StringOrInt"
Cohesion: 0.14
Nodes (5): ComplexityRoot, NotificationConfigurationConnection, ProjectConnection, QueryResolver, WorkspaceConnection

### Community 17 - "WorkspaceResolver"
Cohesion: 0.17
Nodes (3): WorkspaceResolver, jobEdges(), workspaceTagEdges()

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "page"
Cohesion: 0.17
Nodes (7): page, accessEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), projectAccessEdges(), projectEdges(), templateEdges()

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 23 - "run"
Cohesion: 0.06
Nodes (44): main(), run(), main(), run(), main(), run(), main(), run() (+36 more)

### Community 25 - "Job"
Cohesion: 0.06
Nodes (20): cancelRunStore, gitConfig, preparedWorkspace, reliabilityStore, captureJobs(), captureWorkspaces(), TestSanitizeGitURL(), transient() (+12 more)

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

### Community 31 - "Executor"
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
Cohesion: 0.24
Nodes (4): NotificationConfigurationResolver, OrganizationResolver, idsFrom(), workspaceEdges()

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
Cohesion: 0.05
Nodes (64): upstreamCapture, newTestAPI(), TestAuthMiddleware(), TestCreateOrganization(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), TestMalformedIDIsBadRequest() (+56 more)

### Community 53 - "api/main.go"
Cohesion: 0.06
Nodes (48): globalVarAttrs, jobAttrs, organizationAttrs, projectAttrs, replay, variableAttrs, workspaceAttrs, S3 (+40 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - ".streamStepLog"
Cohesion: 0.23
Nodes (10): isTerminal(), outputOf(), stepOf(), streamKey(), TestIsTerminal(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData() (+2 more)

### Community 56 - "testLogger"
Cohesion: 0.09
Nodes (40): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+32 more)

### Community 58 - ".processJob"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 60 - "RequestData"
Cohesion: 0.19
Nodes (11): nestedDocument, pickStringValue(), DecodeAttributes(), RequestData, RelationshipID(), RelationshipIdentifiers(), TestDecodeRequest(), TestRelationshipIdentifiers() (+3 more)

### Community 61 - "Queries"
Cohesion: 0.05
Nodes (7): historyStub, Querier, SetWorkspaceLastJobParams, Publisher, publishingStore, UpdateOrganizationParams, stubQuerier

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 67 - "Layout"
Cohesion: 0.18
Nodes (9): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+1 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.04
Nodes (13): AccessConnection, AccessEdge, JobConnection, JobEdge, NotificationTriggerConnection, NotificationTriggerEdge, PageInfo, ProjectAccessConnection (+5 more)

### Community 70 - "Layout"
Cohesion: 0.13
Nodes (8): Validator, Deps, JobLocker, Layout, Deps, JobQueue, stubJobQueue, Deps

### Community 71 - "collections.go"
Cohesion: 0.06
Nodes (43): collectionSpec, DB, failingDB, listArgs, afterString(), column(), loadPage(), predicate() (+35 more)

### Community 72 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 73 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "config.go"
Cohesion: 0.20
Nodes (14): env(), envBool(), envDuration(), envFirst(), envInt64(), envList(), Load(), clearEnv() (+6 more)

### Community 76 - "Resolver"
Cohesion: 0.12
Nodes (8): ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 77 - "New"
Cohesion: 0.24
Nodes (11): TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput(), TestGetStepLogPrefersArchivedObject(), TestGetStepLogRejectsWrongScope(), TestGetStepLogRequiresToken() (+3 more)

### Community 78 - "serveNested"
Cohesion: 0.14
Nodes (22): jobStub, newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath(), TestJobIncludeReturnsStepsAndWorkspace() (+14 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "serveNested"
Cohesion: 0.36
Nodes (10): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newNestedStub(), serveNested(), TestNestedVariableCreateUsesTheWorkspaceInThePath(), TestNestedVariableDeleteRemovesOnlyThatVariable(), TestNestedVariableFromAnotherWorkspaceIsNotFound(), TestNestedVariablesListIsScopedToTheWorkspace() (+2 more)

### Community 81 - "serveGlobalVar"
Cohesion: 0.39
Nodes (9): newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent(), TestGlobalVariableFromAnotherOrganizationIsNotFound(), TestListGlobalVariablesHidesSensitiveValues() (+1 more)

### Community 82 - "newTestExecutor"
Cohesion: 0.12
Nodes (20): fakeLocker, newFakeLocker(), TestHandleReleasesTheJobLock(), TestHandleSkipsAJobHeldByAnotherExecutor(), TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestStaleClaimIdleOutlivesTheJobTimeout(), New() (+12 more)

### Community 83 - "query.go"
Cohesion: 0.38
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "New"
Cohesion: 0.33
Nodes (6): NewLocal(), TestLocalHealth(), TestLocalLifecycle(), TestLocalRejectsPathTraversal(), NewS3(), New()

### Community 89 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.20
Nodes (5): timestamp(), GetTagByNameParams, NotificationConfiguration, NotificationTrigger, Tag

### Community 92 - "stubQuerier"
Cohesion: 0.11
Nodes (4): nestedStub, workspaceStub, CreateWorkspaceParams, UpdateWorkspaceParams

### Community 93 - "logSink"
Cohesion: 0.12
Nodes (9): captureLogPublisher, LogPublisher, logSink, redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues() (+1 more)

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 103 - ".Query"
Cohesion: 0.08
Nodes (7): New(), ApproveJobParams, CreateJobParams, DBTX, Queries, SetJobOverrideSourceParams, UpdateJobPlanChangesParams

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 123 - ".Ping"
Cohesion: 0.33
Nodes (3): Deps, fakePinger, Pinger

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 128 - ".QueryRow"
Cohesion: 0.09
Nodes (3): History, CreateStepParams, GetStepByJobAndNumberParams

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 324 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **52 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `executionContext` connect `context.Context` to `github.com/99designs/gqlgen/graphql.CollectedField`, `.___EnumValue`, `github.com/99designs/gqlgen/graphql.Marshaler`, `.___Type`, `generated.go`, `Resolver`, `StringOrInt`, `._NotificationConfiguration`, `.marshalOBoolean2bool`, `OrganizationConnection`, `TemplateConnection`, `.___Field`, `.___Directive`, `._WorkspaceTag`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.10047846889952153 - nodes in this community are weakly interconnected._
- **Why does `Querier` connect `Queries` to `.QueryRow`, `stubQuerier`, `stubQuerier`, `Handler`, `TFE Handler Tests`, `WriteResource`, `artifactStore`, `jsonapi.go`, `Project`, `run`, `Job`, `stubStore`, `testing.T`, `api/main.go`, `nestedStub`, `Layout`, `New`, `serveNested`, `newTestExecutor`, `.GetWorkspaceByName`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `stubQuerier`, `.Query`, `.CreateArchive`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.08428446005267778 - nodes in this community are weakly interconnected._
- **Why does `newLogSink()` connect `testLogger` to `newTestExecutor`, `.marshalOBoolean2bool`, `api/main.go`, `logSink`, `Executor`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._