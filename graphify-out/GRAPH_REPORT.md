# Graph Report - terrakube-go  (2026-11-22)

## Corpus Check
- 322 files · ~161,964 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 2995 nodes · 11719 edges · 142 communities (76 shown, 66 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 2263 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `cd50624a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- executionContext
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
- Vc
- TemplateConnection
- Ssh
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- context.Context
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
- Collection
- testLogger
- nestedStub
- .processJob
- New
- WriteError
- Querier
- cancelRunStore
- stubQuerier
- dex-config.yaml
- net/http.Handler
- models
- github.com/99designs/gqlgen/graphql.Marshaler
- testDeps
- collections.go
- Load
- newRedactor
- generated.go
- where_test.go
- streamEvents
- jobStub
- decodeJSON
- graphql/models.go
- ui.go
- Step
- newTestExecutor
- models_gen.go
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- .QueryRow
- .Run
- Action
- .PublishCapped
- Node
- environment
- streamEvents
- .Query
- ParsePage
- ui-rest-inventory.sh
- ui-smoke.sh
- .Ping
- capture.sh
- 0008_action.sql
- NotificationTriggerConnection
- ProjectAccessConnection
- compare.sh
- start.sh
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- WorkspaceTagConnection
- tarGzEntries
- AccessConnection
- Workspace

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 146 edges
3. `Queries` - 113 edges
4. `New()` - 103 edges
5. `API` - 97 edges
6. `stubQuerier` - 86 edges
7. `WriteError()` - 83 edges
8. `API` - 66 edges
9. `stubQuerier` - 55 edges
10. `param()` - 53 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Upstream Elide entities (GraphQL surface)` --references--> `Template`  [INFERRED]
  docs/upstream/graphql-entities.md → internal/store/models.go
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
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

## Communities (142 total, 66 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (23): WriteResource(), configurationVersionResource(), projectResource(), parseWorkspaceQuery(), approvalGateStatus(), failedJobStatus(), parseRunID(), stepStatusTFE() (+15 more)

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
Nodes (52): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+44 more)

### Community 6 - "Tag"
Cohesion: 0.08
Nodes (8): tagStub, Tag, Workspacetag, CreateTagParams, CreateWorkspaceTagParams, GetTagByNameParams, UpdateTagParams, UpdateWorkspaceTagValueParams

### Community 7 - "collectionStub"
Cohesion: 0.09
Nodes (6): collectionStub, CreateItemParams, CreateReferenceParams, UpdateItemParams, Item, Reference

### Community 8 - "Handler"
Cohesion: 0.12
Nodes (17): Claims, Type, Classify(), TestClassify(), TestEffectiveStepIDFromJTI(), claimsFromContext(), isActiveStepStatus(), isExternalUser() (+9 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.11
Nodes (37): assertContractJSONAPI(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage(), stringPtr() (+29 more)

### Community 10 - "WriteResource"
Cohesion: 0.12
Nodes (16): API, Deps, listQuery, filterSortPage(), New(), param(), parseList(), requireType() (+8 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.16
Nodes (21): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+13 more)

### Community 12 - "graphql/collections.go"
Cohesion: 0.10
Nodes (22): collectionSpec, DB, listArgs, page, accessEdges(), afterString(), column(), jobEdges() (+14 more)

### Community 13 - "Resource"
Cohesion: 0.08
Nodes (38): vcsAttrs, statusQuerier, actionResource(), parseInclude(), strPtr(), strValue(), collectionResource(), referenceResource() (+30 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 16 - "Vc"
Cohesion: 0.13
Nodes (4): vcsStub, Vc, CreateVCSParams, UpdateVCSParams

### Community 18 - "Ssh"
Cohesion: 0.05
Nodes (10): sshAgentStub, stubQuerier, demoOrganizations(), Agent, Organization, Ssh, CreateOrganizationParams, CreateAgentParams (+2 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.12
Nodes (26): replay, upstreamCapture, sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), assertSameDocumentShape(), assertSameResourceShape() (+18 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 24 - "context.Context"
Cohesion: 0.06
Nodes (3): ComplexityRoot, ProjectConnection, WorkspaceConnection

### Community 25 - "Job"
Cohesion: 0.08
Nodes (13): gitConfig, preparedWorkspace, reliabilityStore, captureJobs(), TestGitSourceDetection(), transient(), gitEnv(), gitSource() (+5 more)

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
Nodes (5): Executor, logSink, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

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
Cohesion: 0.11
Nodes (10): NotificationConfigurationResolver, OrganizationResolver, ProjectResolver, QueryResolver, Resolver, ResolverRoot, WorkspaceResolver, idsFrom() (+2 more)

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
Cohesion: 0.10
Nodes (7): fakeQueue, JobQueue, defaultRetryDelay(), Message, Queue, JobQueue, stubJobQueue

### Community 52 - "testing.T"
Cohesion: 0.05
Nodes (67): newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes(), newTestAPI() (+59 more)

### Community 53 - "go_pkg_context"
Cohesion: 0.06
Nodes (53): agentAttrs, atomicResult, collectionAttrs, globalVarAttrs, jobAttrs, organizationAttrs, projectAttrs, sshAttrs (+45 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Collection"
Cohesion: 0.17
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 56 - "testLogger"
Cohesion: 0.06
Nodes (49): Deps, fakeLocker, JobLocker, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive() (+41 more)

### Community 58 - ".processJob"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestIsTerminalStatus(), contentIDFromSource()

### Community 59 - "New"
Cohesion: 0.09
Nodes (25): Validator, ExtractToken(), TestExtractToken(), TestWriteBackendOverride(), isTerminal(), New(), outputOf(), stepOf() (+17 more)

### Community 60 - "WriteError"
Cohesion: 0.08
Nodes (37): actionAttrs, atomicOperation, itemAttrs, nestedDocument, decodeActionAttrs(), stringOr(), boolValue(), decodeItemAttrs() (+29 more)

### Community 61 - "Querier"
Cohesion: 0.03
Nodes (12): artifactStore, historyStub, New(), Querier, DBTX, GetProjectByNameParams, Queries, TempArchive (+4 more)

### Community 62 - "cancelRunStore"
Cohesion: 0.11
Nodes (5): cancelRunStore, Publisher, publishingStore, stubQuerier, mutableQuerier

### Community 63 - "stubQuerier"
Cohesion: 0.05
Nodes (6): CreateContentParams, UpdateContentStatusParams, Project, CreateProjectParams, Content, stubQuerier

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "net/http.Handler"
Cohesion: 0.05
Nodes (51): main(), run(), main(), run(), main(), run(), main(), run() (+43 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 70 - "testDeps"
Cohesion: 0.24
Nodes (9): TestUpstreamUIQueriesExecute(), runQuery(), testDeps(), TestFilterByTagAndJobStatus(), TestRelationshipIDsSelectOneRow(), TestWorkspacePageEndToEnd(), TestContractDiscoveryAndPing(), Discovery() (+1 more)

### Community 71 - "collections.go"
Cohesion: 0.26
Nodes (8): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), Op, parser

### Community 72 - "Load"
Cohesion: 0.07
Nodes (27): Config, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1` (+19 more)

### Community 73 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "streamEvents"
Cohesion: 0.24
Nodes (5): SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 78 - "decodeJSON"
Cohesion: 0.16
Nodes (22): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath() (+14 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "ui.go"
Cohesion: 0.18
Nodes (14): envConfig(), exists(), New(), get(), newUI(), TestDeepLinksServeTheIndex(), TestEnvConfigIsNeverCached(), TestMissingAssetsAreNotFound() (+6 more)

### Community 82 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 83 - "models_gen.go"
Cohesion: 0.08
Nodes (9): JobConnection, JobEdge, NotificationConfigurationConnection, NotificationConfigurationEdge, OrganizationConnection, OrganizationEdge, PageInfo, Query (+1 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "workspaces_test.go"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 92 - "Action"
Cohesion: 0.14
Nodes (6): actionFilter, actionStub, parseActionFilter(), CreateActionParams, UpdateActionParams, Action

### Community 94 - "Node"
Cohesion: 0.22
Nodes (6): likePattern(), And, builder, Compare, Node, Or

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 103 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 133 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.11
Nodes (3): History, NotificationConfiguration, NotificationTrigger

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 144 - "Workspace"
Cohesion: 0.16
Nodes (6): workspaceStub, workspaceUpdateParams(), Workspace, CreateWorkspaceParams, UpdateWorkspaceParams, workspaceWithTags

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 350 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **66 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `net/http.ResponseWriter`, `Globalvar`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `Tag`, `collectionStub`, `Handler`, `TFE Handler Tests`, `WriteResource`, `Resource`, `Vc`, `Workspace`, `Ssh`, `strPtr`, `Job`, `stubStore`, `testing.T`, `go_pkg_context`, `Collection`, `testLogger`, `nestedStub`, `New`, `cancelRunStore`, `stubQuerier`, `net/http.Handler`, `jobStub`, `Step`, `newTestExecutor`, `.QueryRow`, `Action`?**
  _High betweenness centrality (0.041) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.08307210031347963 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `executionContext` to `ProjectAccessConnection`, `context.Context`, `Job`, `github.com/99designs/gqlgen/graphql.Marshaler`, `StringOrInt`, `WorkspaceTagConnection`, `generated.go`, `AccessConnection`, `StringOrInt`, `TemplateConnection`, `github.com/99designs/gqlgen/graphql.CollectedField`, `models_gen.go`, `.childFields___InputValue`, `context.Context`, `NotificationTriggerConnection`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Should `executionContext` be split into smaller, more focused modules?**
  _Cohesion score 0.08926553672316384 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.ResponseWriter` to `go_pkg_context`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._