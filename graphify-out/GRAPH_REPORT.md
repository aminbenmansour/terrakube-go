# Graph Report - terrakube-go  (2026-11-22)

## Corpus Check
- 374 files · ~172,239 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 3165 nodes · 12527 edges · 162 communities (88 shown, 74 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 2429 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b7317a1a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- Agent Workflows and Compatibility Docs
- Globalvar
- io.Reader
- Core Database Schema
- Tag
- collectionStub
- Handler
- TFE Handler Tests
- WriteResource
- State Backend Tests
- graphql/collections.go
- Resource
- steps
- Step
- Vc
- stubQuerier
- Ssh
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- github.com/99designs/gqlgen/graphql.Marshaler
- executionContext
- .acquireConfiguration
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- state.sql.go
- .runApply
- Upstream UI GraphQL operations
- OpenTofu State Conformance
- Compose Project Definition
- Compose Healthchecks
- Compose Service Dependencies
- StringOrInt
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
- logSink
- context.Context
- nestedStub
- .handle
- New
- WriteError
- Querier
- workspaceTagBindings
- stubQuerier
- dex-config.yaml
- artifacts_test.go
- .QueryRow
- models
- models_gen.go
- Team
- collections.go
- Load
- isolation_test.go
- generated.go
- where_test.go
- Resolver
- decodeJSON
- graphql/models.go
- Collection
- newTestAPI
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- ui.go
- .updateWorkspaceRelationship
- Action
- .CreateNotificationConfiguration
- Node
- environment
- NewSubscriptions
- log/slog.Logger
- OrganizationConnection
- ui-rest-inventory.sh
- ui-smoke.sh
- SetWorkspaceLastJobParams
- .Ping
- capture.sh
- 0008_action.sql
- ParsePage
- ._Project
- compare.sh
- start.sh
- .QueryRow
- WorkspaceTagConnection
- workspace.go
- TemplateConnection
- JobConnection
- testDeps
- Wrap
- .CreateNotificationTrigger
- serveNested
- serveGlobalVar
- query.go
- .PublishCapped
- notificationTestAPI
- get
- net/http.Handler
- .EnqueueJob
- actionFilter
- .atomicOperations
- transientError

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 174 edges
3. `Queries` - 136 edges
4. `New()` - 126 edges
5. `API` - 123 edges
6. `stubQuerier` - 113 edges
7. `WriteError()` - 95 edges
8. `API` - 66 edges
9. `param()` - 64 edges
10. `ParseFieldsets()` - 59 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `New()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfoutput/handler.go
- `Generated Code and Migration Policy` --semantically_similar_to--> `sqlc Workflow`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .opencode/skills/sqlc-workflow/SKILL.md
- `GraphQL UI Requirement` --semantically_similar_to--> `GraphQL and Subscription UI Contract`  [INFERRED] [semantically similar]
  AGENTS.md → docs/compat.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (162 total, 74 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (25): Write(), WriteResource(), configurationVersionResource(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource() (+17 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "Globalvar"
Cohesion: 0.14
Nodes (5): globalVarStub, variableStore, CreateGlobalVariableParams, UpdateGlobalVariableParams, Globalvar

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (14): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+6 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "Tag"
Cohesion: 0.07
Nodes (8): tagStub, Tag, Workspacetag, CreateTagParams, CreateWorkspaceTagParams, GetTagByNameParams, UpdateTagParams, UpdateWorkspaceTagValueParams

### Community 7 - "collectionStub"
Cohesion: 0.08
Nodes (6): collectionStub, CreateItemParams, CreateReferenceParams, UpdateItemParams, Item, Reference

### Community 8 - "Handler"
Cohesion: 0.08
Nodes (22): Claims, Type, Validator, Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken() (+14 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.10
Nodes (38): assertContractJSONAPI(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), capabilitySecret(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage() (+30 more)

### Community 10 - "WriteResource"
Cohesion: 0.12
Nodes (12): API, Deps, listQuery, filterSortPage(), New(), param(), parseList(), strValue() (+4 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.16
Nodes (21): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+13 more)

### Community 12 - "graphql/collections.go"
Cohesion: 0.13
Nodes (15): collectionSpec, DB, failingDB, listArgs, afterString(), column(), loadPage(), predicate() (+7 more)

### Community 13 - "Resource"
Cohesion: 0.12
Nodes (31): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), jobResource(), stepResource() (+23 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 15 - "Step"
Cohesion: 0.05
Nodes (3): reliabilityStore, Step, History

### Community 16 - "Vc"
Cohesion: 0.14
Nodes (4): vcsStub, Vc, CreateVCSParams, UpdateVCSParams

### Community 17 - "stubQuerier"
Cohesion: 0.05
Nodes (10): stubQuerier, templateStub, demoOrganizations(), Organization, Template, CreateOrganizationParams, CreateTemplateParams, UpdateTemplateParams (+2 more)

### Community 18 - "Ssh"
Cohesion: 0.08
Nodes (6): sshAgentStub, Agent, Ssh, CreateAgentParams, CreateSSHKeyParams, UpdateSSHKeyParams

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.10
Nodes (34): replay, upstreamCapture, sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), newTeamStub(), TestTeamJSONAPI() (+26 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.06
Nodes (3): AccessEdge, NotificationTriggerEdge, WorkspaceEdge

### Community 25 - ".acquireConfiguration"
Cohesion: 0.18
Nodes (7): preparedWorkspace, TestSanitizeGitURL(), transient(), contentIDFromSource(), gitEnv(), runGit(), sanitizeGitURL()

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

### Community 30 - "state.sql.go"
Cohesion: 0.14
Nodes (24): agentAttrs, atomicResult, collectionAttrs, globalVarAttrs, jobAttrs, notificationTriggerAttrs, organizationAttrs, projectAttrs (+16 more)

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

### Community 37 - "StringOrInt"
Cohesion: 0.10
Nodes (12): OrganizationResolver, page, accessEdges(), idsFrom(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), projectAccessEdges() (+4 more)

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
Nodes (62): TestAuthMiddleware(), TestCreateOrganization(), newCollectionAPI(), TestCollectionListCarriesItsItemsAndReferences(), TestCreateCollectionDefaultsPriorityToZero(), TestCreateCollectionItemStoresTheValue(), TestCreateCollectionItemValidatesKeyValueAndCategory(), TestCreateCollectionRequiresName() (+54 more)

### Community 53 - "go_pkg_context"
Cohesion: 0.08
Nodes (20): stateMetadata, failingRow, graphQLResponse, configureProcessGroup(), killProcessGroup(), configureProcessGroup(), killProcessGroup(), ptr() (+12 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "logSink"
Cohesion: 0.25
Nodes (3): logSink, TestNilRedactorIsSafe(), baseEnvironment()

### Community 58 - ".handle"
Cohesion: 0.18
Nodes (4): isTransient(), parseAttempt(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 59 - "New"
Cohesion: 0.20
Nodes (12): TestWriteBackendOverride(), isTerminal(), New(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame() (+4 more)

### Community 60 - "WriteError"
Cohesion: 0.09
Nodes (42): actionAttrs, atomicOperation, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, vcsAttrs, decodeActionAttrs() (+34 more)

### Community 61 - "Querier"
Cohesion: 0.03
Nodes (9): historyStub, AddNotificationConfigurationTemplateParams, Querier, DBTX, Queries, TempArchive, UpdateHistoryOutputParams, UpdateOrganizationParams (+1 more)

### Community 62 - "workspaceTagBindings"
Cohesion: 0.22
Nodes (8): nestedDocument, RelationshipList(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship, ResourceIdentifier, tagBindingInput

### Community 63 - "stubQuerier"
Cohesion: 0.04
Nodes (19): workspaceStub, cancelRunStore, statusQuerier, captureJobs(), workspaceUpdateParams(), UpdateContentStatusParams, Job, Project (+11 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "artifacts_test.go"
Cohesion: 0.04
Nodes (75): main(), run(), main(), run(), main(), run(), main(), run() (+67 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "models_gen.go"
Cohesion: 0.07
Nodes (10): AccessConnection, ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, NotificationTriggerConnection, PageInfo, ProjectConnection, Query (+2 more)

### Community 70 - "Team"
Cohesion: 0.17
Nodes (4): teamStub, Team, CreateTeamParams, UpdateTeamParams

### Community 71 - "collections.go"
Cohesion: 0.26
Nodes (8): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), Op, parser

### Community 72 - "Load"
Cohesion: 0.22
Nodes (17): S3, env(), envBool(), envDuration(), envFirst(), envInt64(), envList(), Load() (+9 more)

### Community 73 - "isolation_test.go"
Cohesion: 0.24
Nodes (9): redactor, contains(), envValue(), TestEnvironmentSeparatesSecrets(), TestLogSinkRedactsSecrets(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), TestRunCommandDoesNotInheritExecutorSecrets() (+1 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "Resolver"
Cohesion: 0.12
Nodes (11): NotificationConfigurationResolver, ProjectResolver, QueryResolver, Resolver, ResolverRoot, SubscriptionResolver, WorkspaceResolver, streamEvents() (+3 more)

### Community 78 - "decodeJSON"
Cohesion: 0.25
Nodes (12): jobStub, newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath(), TestJobIncludeReturnsStepsAndWorkspace() (+4 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "Collection"
Cohesion: 0.20
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 82 - "newTestAPI"
Cohesion: 0.15
Nodes (18): newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes(), newTestAPI() (+10 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "workspaces_test.go"
Cohesion: 0.16
Nodes (22): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+14 more)

### Community 88 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 92 - "Action"
Cohesion: 0.17
Nodes (4): actionStub, CreateActionParams, UpdateActionParams, Action

### Community 94 - "Node"
Cohesion: 0.22
Nodes (6): likePattern(), And, builder, Compare, Node, Or

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "NewSubscriptions"
Cohesion: 0.15
Nodes (10): Deps, EventSource, fakeEvents, NewExecutableSchema(), TestSchemaLoads(), bearerMatches(), New(), NewSubscriptions() (+2 more)

### Community 101 - "log/slog.Logger"
Cohesion: 0.13
Nodes (11): Config, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1` (+3 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 127 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 128 - "._Project"
Cohesion: 0.13
Nodes (3): ProjectAccessConnection, ProjectAccessEdge, ProjectEdge

### Community 133 - ".QueryRow"
Cohesion: 0.06
Nodes (5): CreateContentParams, Content, GetProjectByNameParams, UpdateProjectParams, Variable

### Community 140 - "workspace.go"
Cohesion: 0.11
Nodes (28): archiveLimits, gitConfig, tarEntry, gitOutput(), initGitRepo(), requireGit(), tarGzEntries(), TestAcquireConfigurationClonesGitSource() (+20 more)

### Community 144 - "testDeps"
Cohesion: 0.24
Nodes (9): TestUpstreamUIQueriesExecute(), runQuery(), testDeps(), TestFilterByTagAndJobStatus(), TestRelationshipIDsSelectOneRow(), TestWorkspacePageEndToEnd(), TestContractDiscoveryAndPing(), Discovery() (+1 more)

### Community 145 - "Wrap"
Cohesion: 0.27
Nodes (10): dialSubscriptions(), readMessage(), TestSubscriptionOverGraphQLWS(), TestSubscriptionRejectsABadConnectionToken(), TestSubscriptionSeesAStatusWrittenThroughTheStore(), writeMessage(), TestWrapDoesNotPublishAFailedUpdate(), TestWrapIgnoresAPublishFailure() (+2 more)

### Community 147 - "serveNested"
Cohesion: 0.36
Nodes (10): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newNestedStub(), serveNested(), TestNestedVariableCreateUsesTheWorkspaceInThePath(), TestNestedVariableDeleteRemovesOnlyThatVariable(), TestNestedVariableFromAnotherWorkspaceIsNotFound(), TestNestedVariablesListIsScopedToTheWorkspace() (+2 more)

### Community 150 - "serveGlobalVar"
Cohesion: 0.39
Nodes (9): newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent(), TestGlobalVariableFromAnotherOrganizationIsNotFound(), TestListGlobalVariablesHidesSensitiveValues() (+1 more)

### Community 151 - "query.go"
Cohesion: 0.32
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 153 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 154 - "get"
Cohesion: 0.53
Nodes (6): get(), newUI(), TestDeepLinksServeTheIndex(), TestEnvConfigIsNeverCached(), TestMissingAssetsAreNotFound(), TestServesTheIndexAndAssets()

### Community 155 - "net/http.Handler"
Cohesion: 0.40
Nodes (3): requireUUIDParams(), teamTestAPI(), templateTestAPI()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 358 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **74 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `net/http.ResponseWriter`, `Globalvar`, `.QueryRow`, `Tag`, `collectionStub`, `Handler`, `TFE Handler Tests`, `WriteResource`, `Step`, `Vc`, `stubQuerier`, `Ssh`, `Wrap`, `.CreateNotificationTrigger`, `strPtr`, `notificationTestAPI`, `net/http.Handler`, `state.sql.go`, `testing.T`, `context.Context`, `nestedStub`, `New`, `stubQuerier`, `artifacts_test.go`, `.QueryRow`, `Team`, `decodeJSON`, `Collection`, `newTestAPI`, `Action`, `.CreateNotificationConfiguration`, `SetWorkspaceLastJobParams`?**
  _High betweenness centrality (0.053) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.08058608058608059 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `executionContext` to `._Project`, `context.Context`, `.marshalOString2ᚖstring`, `models_gen.go`, `OrganizationConnection`, `WorkspaceTagConnection`, `generated.go`, `Resolver`, `._NotificationConfiguration`, `TemplateConnection`, `JobConnection`, `._Organization`, `.___EnumValue`, `.marshalOBoolean2bool`, `github.com/99designs/gqlgen/graphql.Marshaler`, `context.Context`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `workspaceQuery` connect `query.go` to `net/http.ResponseWriter`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `Globalvar` be split into smaller, more focused modules?**
  _Cohesion score 0.14492753623188406 - nodes in this community are weakly interconnected._