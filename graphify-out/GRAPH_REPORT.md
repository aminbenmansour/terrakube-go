# Graph Report - terrakube-go  (2026-11-21)

## Corpus Check
- 182 files · ~132,254 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 2491 nodes · 9380 edges · 128 communities (83 shown, 45 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 1828 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `71e38d47`
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
- artifactStore
- jsonapi.go
- steps
- StringOrInt
- DecodeRequest
- .QueryRow
- api/main.go
- sqlc Generation Config
- tarGzEntries
- Upstream Elide entities (GraphQL surface)
- run
- TemplateConnection
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
- where_test.go
- Redis Compose Service
- M4 run logs and artifacts — design
- Management API Smoke Test
- TFE API Smoke Test
- State Backend Smoke Test
- Go Module Definition
- time.Duration
- testing.T
- state.sql.go
- executor-e2e.sh
- Layout
- testLogger
- .build
- .processJob
- go_pkg_strings
- .processJob
- Queries
- models_gen.go
- .marshalNString2string
- dex-config.yaml
- .prepareWorkspace
- tfoutput/handler.go
- models
- github.com/99designs/gqlgen/graphql.Marshaler
- collections.go
- where.go
- workspace.go
- newExecutionContext
- config.go
- where_test.go
- ._Project
- serveNested
- graphql/models.go
- .streamStepLog
- ._WorkspaceTag
- newTestExecutor
- NewSubscriptions
- REST contract, captured from upstream
- pat.py
- CreateWorkspaceTagParams
- s3.go
- Querier
- ui.go
- .listWorkspaces
- stubQuerier
- logSink
- context.Context
- environment
- streamEvents
- CreateTagParams
- newRedactor
- .Query
- Workspace
- ui-rest-inventory.sh
- ui-smoke.sh
- testDeps
- assertSameDocumentShape
- .CreateArchive
- .Ping
- capture.sh
- Step
- compare.sh
- start.sh

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 88 edges
3. `API` - 66 edges
4. `Queries` - 65 edges
5. `stubQuerier` - 48 edges
6. `Executor` - 46 edges
7. `New()` - 45 edges
8. `WriteError()` - 43 edges
9. `New()` - 42 edges
10. `API` - 33 edges

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

## Communities (128 total, 45 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.14
Nodes (10): WriteResource(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), normalizeTagValue(), tagValuesEqual(), New() (+2 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "ParseRSQL"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (11): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+3 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.09
Nodes (37): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+29 more)

### Community 6 - "stubQuerier"
Cohesion: 0.06
Nodes (7): CreateContentParams, UpdateContentStatusParams, CreateProjectParams, CreateWorkspaceParams, SetWorkspaceLockParams, Content, stubQuerier

### Community 8 - "Handler"
Cohesion: 0.12
Nodes (16): Claims, Type, Classify(), ExtractToken(), claimsFromContext(), isActiveStepStatus(), isExternalUser(), withClaims() (+8 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.10
Nodes (38): assertContractJSONAPI(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), capabilitySecret(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage() (+30 more)

### Community 10 - "WriteResource"
Cohesion: 0.15
Nodes (22): API, Deps, listQuery, boolValue(), filterSortPage(), New(), param(), parseInclude() (+14 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.14
Nodes (23): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+15 more)

### Community 13 - "jsonapi.go"
Cohesion: 0.09
Nodes (28): jobResource(), stepResource(), captureJobs(), Page, Resource, Link(), LinkList(), TestRelationshipMarshalling() (+20 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 15 - "StringOrInt"
Cohesion: 0.16
Nodes (5): QueryResolver, WorkspaceResolver, idsFrom(), projectEdges(), workspaceEdges()

### Community 16 - "DecodeRequest"
Cohesion: 0.14
Nodes (19): nestedDocument, pickBool(), pickString(), pickStringValue(), DecodeAttributes(), DecodeRequest(), RequestData, HasAttribute() (+11 more)

### Community 18 - ".QueryRow"
Cohesion: 0.12
Nodes (5): CreateHistoryParams, GetHistoryByJobReferenceParams, GetProjectByNameParams, History, UpdateProjectParams

### Community 19 - "api/main.go"
Cohesion: 0.14
Nodes (9): transientError, failingRow, graphQLResponse, ptr(), TestCollectionSpecRejectsInjection(), CORS(), TestCORSAllowsAnyOriginByDefault(), TestCORSAllowsOnlyListedOrigins() (+1 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "tarGzEntries"
Cohesion: 0.22
Nodes (13): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+5 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 23 - "run"
Cohesion: 0.06
Nodes (39): main(), run(), main(), run(), main(), run(), main(), run() (+31 more)

### Community 25 - "cancelRunStore"
Cohesion: 0.08
Nodes (10): reliabilityStore, statusQuerier, timestamp(), UpdateJobStatusParams, Publisher, publishingStore, stubQuerier, NotificationConfiguration (+2 more)

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
Cohesion: 0.24
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

### Community 37 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

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
Nodes (6): fakeQueue, JobQueue, envDuration(), defaultRetryDelay(), Message, Queue

### Community 52 - "testing.T"
Cohesion: 0.07
Nodes (48): newTestAPI(), TestAuthMiddleware(), TestCreateOrganization(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), TestMalformedIDIsBadRequest(), TestMountPrefixDoesNotConfuseTheIDCheck() (+40 more)

### Community 53 - "state.sql.go"
Cohesion: 0.19
Nodes (11): jobAttrs, organizationAttrs, projectAttrs, replay, variableAttrs, workspaceAttrs, projectAttrs, stateDocumentMetadata (+3 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Layout"
Cohesion: 0.08
Nodes (16): Validator, Config, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log (+8 more)

### Community 56 - "testLogger"
Cohesion: 0.08
Nodes (45): TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestArchiveHistoryDoesNotCreateRowWhenUploadFails(), TestCompleteStepArchivesLog(), TestCompleteStepCancelledKeepsStatusWhenLogUploadFails(), TestCompleteStepLogUploadFailureFailsStep(), TestFinishApplyCompletesStep() (+37 more)

### Community 57 - ".build"
Cohesion: 0.24
Nodes (3): DB, failingDB, loadPage()

### Community 58 - ".processJob"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 59 - "go_pkg_strings"
Cohesion: 0.18
Nodes (7): stateMetadata, Backend, Access, NotificationConfigurationTemplate, ProjectAccess, claimsContextKey, externalUserContextKey

### Community 60 - ".processJob"
Cohesion: 0.25
Nodes (3): stringPointer(), TestContentIDFromSource(), contentIDFromSource()

### Community 61 - "Queries"
Cohesion: 0.09
Nodes (6): SetWorkspaceLastJobParams, DBTX, Queries, UpdateHistoryOutputParams, UpdateJobOutputParams, UpdateStepOutputParams

### Community 62 - "models_gen.go"
Cohesion: 0.12
Nodes (4): NotificationTriggerConnection, OrganizationConnection, OrganizationEdge, PageInfo

### Community 63 - ".marshalNString2string"
Cohesion: 0.07
Nodes (7): NotificationConfigurationConnection, NotificationConfigurationEdge, ProjectConnection, Query, Subscription, WorkspaceConnection, WorkspaceEdge

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - ".prepareWorkspace"
Cohesion: 0.21
Nodes (6): preparedWorkspace, TestSanitizeGitURL(), transient(), gitEnv(), runGit(), sanitizeGitURL()

### Community 67 - "tfoutput/handler.go"
Cohesion: 0.13
Nodes (8): compareStreamIDs(), parseResumeID(), splitStreamID(), TestCompareStreamIDs(), TestParseResumeID(), validStreamID(), configurationVersionAttributes, runAttributes

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.07
Nodes (4): AccessConnection, AccessEdge, NotificationTriggerEdge, ProjectAccessEdge

### Community 71 - "collections.go"
Cohesion: 0.10
Nodes (24): collectionSpec, listArgs, page, accessEdges(), afterString(), column(), jobEdges(), notificationConfigurationEdges() (+16 more)

### Community 72 - "where.go"
Cohesion: 0.42
Nodes (3): likePattern(), builder, Compare

### Community 73 - "workspace.go"
Cohesion: 0.15
Nodes (9): gitConfig, TestGitSourceDetection(), configureProcessGroup(), killProcessGroup(), configureProcessGroup(), killProcessGroup(), gitSource(), isConfigurationTarball() (+1 more)

### Community 75 - "config.go"
Cohesion: 0.21
Nodes (16): S3, env(), envBool(), envFirst(), envInt64(), envList(), Load(), ParseLogLevel() (+8 more)

### Community 76 - "where_test.go"
Cohesion: 0.08
Nodes (9): ComplexityRoot, DirectiveRoot, NotificationConfigurationResolver, OrganizationResolver, ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver (+1 more)

### Community 78 - "serveNested"
Cohesion: 0.19
Nodes (20): newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath(), TestJobIncludeReturnsStepsAndWorkspace(), TestJobListReturnsTheOrganizationsJobs() (+12 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - ".streamStepLog"
Cohesion: 0.12
Nodes (19): isTerminal(), outputOf(), stepOf(), streamKey(), newHandler(), stringPtr(), TestGetStepLogRejectsWrongScope(), TestGetStepLogRequiresToken() (+11 more)

### Community 82 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 83 - "NewSubscriptions"
Cohesion: 0.22
Nodes (12): Deps, NewExecutableSchema(), TestSchemaLoads(), bearerMatches(), New(), NewSubscriptions(), dialSubscriptions(), readMessage() (+4 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "CreateWorkspaceTagParams"
Cohesion: 0.22
Nodes (3): CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, Workspacetag

### Community 88 - "s3.go"
Cohesion: 0.17
Nodes (3): NewS3(), stateOutput, stateWithOutputs

### Community 89 - "Querier"
Cohesion: 0.09
Nodes (4): historyStub, Querier, UpdateOrganizationParams, stubQuerier

### Community 90 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 91 - ".listWorkspaces"
Cohesion: 0.29
Nodes (6): parseWorkspaceQuery(), matchesTagFilters(), workspaceTagNames(), writeWorkspaceList(), tagFilter, workspaceQuery

### Community 92 - "stubQuerier"
Cohesion: 0.06
Nodes (9): jobStub, nestedStub, stubQuerier, cancelRunStore, variableStore, captureWorkspaces(), demoOrganizations(), CreateOrganizationParams (+1 more)

### Community 93 - "logSink"
Cohesion: 0.17
Nodes (4): captureLogPublisher, LogPublisher, logSink, baseEnvironment()

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.14
Nodes (8): EventSource, fakeEvents, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource(), Event, Bus

### Community 101 - "CreateTagParams"
Cohesion: 0.22
Nodes (3): CreateTagParams, GetTagByNameParams, Tag

### Community 102 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 103 - ".Query"
Cohesion: 0.11
Nodes (4): ApproveJobParams, CreateJobParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "testDeps"
Cohesion: 0.24
Nodes (9): TestUpstreamUIQueriesExecute(), runQuery(), testDeps(), TestFilterByTagAndJobStatus(), TestRelationshipIDsSelectOneRow(), TestWorkspacePageEndToEnd(), TestContractDiscoveryAndPing(), Discovery() (+1 more)

### Community 121 - "assertSameDocumentShape"
Cohesion: 0.32
Nodes (7): upstreamCapture, assertSameDocumentShape(), assertSameResourceShape(), loadUpstreamCaptures(), sortedKeys(), TestCapturedShapesAreWhatWeReturn(), TestEveryUpstreamCaptureIsAccountedFor()

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 128 - "Step"
Cohesion: 0.15
Nodes (6): captureSteps(), strPtr(), CreateStepParams, GetStepByJobAndNumberParams, Step, UpdateStepStatusParams

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 313 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **45 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `executionContext` connect `context.Context` to `executionContext`, `._Job`, `._Job`, `github.com/99designs/gqlgen/graphql.Marshaler`, `.___EnumValue`, `.___Field`, `newExecutionContext`, `where_test.go`, `._Project`, `github.com/99designs/gqlgen/graphql.CollectedField`, `._WorkspaceTag`, `._Template`, `TemplateConnection`, `models_gen.go`, `.marshalNString2string`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.14326923076923076 - nodes in this community are weakly interconnected._
- **Why does `Querier` connect `Querier` to `Step`, `stubQuerier`, `Handler`, `TFE Handler Tests`, `WriteResource`, `artifactStore`, `.QueryRow`, `run`, `cancelRunStore`, `stubStore`, `testing.T`, `state.sql.go`, `Layout`, `Queries`, `.streamStepLog`, `newTestExecutor`, `CreateWorkspaceTagParams`, `stubQuerier`, `CreateTagParams`, `.Query`, `Workspace`, `.CreateArchive`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.ResponseWriter` to `tfoutput/handler.go`, `Handler`, `jsonapi.go`, `DecodeRequest`, `Layout`, `.listWorkspaces`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Should `io.Reader` be split into smaller, more focused modules?**
  _Cohesion score 0.07372549019607844 - nodes in this community are weakly interconnected._