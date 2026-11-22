# Graph Report - terrakube-go  (2026-11-22)

## Corpus Check
- 337 files · ~164,320 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 3039 nodes · 11934 edges · 142 communities (79 shown, 63 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 2299 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9c831ff1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- Agent Workflows and Compatibility Docs
- Globalvar
- io.Reader
- Core Database Schema
- Tag
- stubQuerier
- Handler
- TFE Handler Tests
- WriteResource
- State Backend Tests
- graphql/collections.go
- Resource
- steps
- ._NotificationConfiguration
- Vc
- Template
- Ssh
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- context.Context
- .acquireConfiguration
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
- logSink
- testLogger
- nestedStub
- .handle
- New
- WriteError
- Querier
- workspaceTagBindings
- stubQuerier
- dex-config.yaml
- context.Context
- net/http.Handler
- .QueryRow
- models
- github.com/99designs/gqlgen/graphql.Marshaler
- testDeps
- collections.go
- .createVCS
- newRedactor
- generated.go
- where_test.go
- Resolver
- jobStub
- decodeJSON
- graphql/models.go
- serveGlobalVar
- net/http.Handler
- newTestExecutor
- ._Job
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- Workspace
- .updateWorkspaceRelationship
- Action
- .PublishCapped
- Node
- environment
- streamEvents
- Deps
- .EnqueueJob
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

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 153 edges
3. `Queries` - 118 edges
4. `New()` - 109 edges
5. `API` - 103 edges
6. `stubQuerier` - 92 edges
7. `WriteError()` - 85 edges
8. `API` - 66 edges
9. `param()` - 56 edges
10. `stubQuerier` - 55 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `Layout`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfstate/layout.go
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

## Communities (142 total, 63 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (19): configurationVersionResource(), parseWorkspaceQuery(), approvalGateStatus(), failedJobStatus(), parseRunID(), stepStatusTFE(), stepStatusValue(), decodeStateAttribute() (+11 more)

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
Cohesion: 0.07
Nodes (8): tagStub, Tag, Workspacetag, CreateTagParams, CreateWorkspaceTagParams, GetTagByNameParams, UpdateTagParams, UpdateWorkspaceTagValueParams

### Community 7 - "stubQuerier"
Cohesion: 0.04
Nodes (13): collectionStub, stubQuerier, demoOrganizations(), CreateCollectionParams, CreateItemParams, CreateReferenceParams, UpdateCollectionParams, UpdateItemParams (+5 more)

### Community 8 - "Handler"
Cohesion: 0.07
Nodes (29): Claims, Type, Validator, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log (+21 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.10
Nodes (38): assertContractJSONAPI(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), capabilitySecret(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage() (+30 more)

### Community 10 - "WriteResource"
Cohesion: 0.14
Nodes (11): API, Deps, listQuery, filterSortPage(), New(), param(), parseList(), unknownCollection() (+3 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.16
Nodes (21): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+13 more)

### Community 12 - "graphql/collections.go"
Cohesion: 0.11
Nodes (19): collectionSpec, DB, failingDB, listArgs, afterString(), column(), loadPage(), predicate() (+11 more)

### Community 13 - "Resource"
Cohesion: 0.09
Nodes (36): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), jobResource(), stepResource() (+28 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 16 - "Vc"
Cohesion: 0.13
Nodes (4): vcsStub, Vc, CreateVCSParams, UpdateVCSParams

### Community 17 - "Template"
Cohesion: 0.14
Nodes (4): templateStub, Template, CreateTemplateParams, UpdateTemplateParams

### Community 18 - "Ssh"
Cohesion: 0.09
Nodes (6): sshAgentStub, Agent, Ssh, CreateAgentParams, CreateSSHKeyParams, UpdateSSHKeyParams

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.11
Nodes (29): replay, upstreamCapture, sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), newTemplateStub(), TestTemplateJSONAPI() (+21 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 24 - "context.Context"
Cohesion: 0.06
Nodes (4): artifactStore, ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge

### Community 25 - ".acquireConfiguration"
Cohesion: 0.24
Nodes (5): TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv(), runGit()

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
Cohesion: 0.17
Nodes (5): Executor, preparedWorkspace, isCancelledStatus(), stringPointer(), TestBinarySelection()

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
Nodes (13): page, QueryResolver, WorkspaceResolver, accessEdges(), idsFrom(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges() (+5 more)

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
Nodes (68): newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes(), newTestAPI() (+60 more)

### Community 53 - "go_pkg_context"
Cohesion: 0.05
Nodes (69): actionFilter, agentAttrs, atomicResult, collectionAttrs, globalVarAttrs, jobAttrs, organizationAttrs, projectAttrs (+61 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "logSink"
Cohesion: 0.20
Nodes (4): logSink, redactor, TestNilRedactorIsSafe(), baseEnvironment()

### Community 56 - "testLogger"
Cohesion: 0.08
Nodes (45): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+37 more)

### Community 58 - ".handle"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 59 - "New"
Cohesion: 0.09
Nodes (26): Storage, Task 2: Serve archived logs from `/tfoutput/v1`, TestWriteBackendOverride(), New(), isTerminal(), New(), outputOf(), stepOf() (+18 more)

### Community 60 - "WriteError"
Cohesion: 0.11
Nodes (32): actionAttrs, atomicOperation, itemAttrs, templateAttrs, decodeActionAttrs(), stringOr(), boolValue(), requireType() (+24 more)

### Community 61 - "Querier"
Cohesion: 0.03
Nodes (10): historyStub, New(), Querier, DBTX, GetProjectByNameParams, Queries, TempArchive, UpdateProjectParams (+2 more)

### Community 62 - "workspaceTagBindings"
Cohesion: 0.22
Nodes (8): nestedDocument, RelationshipIdentifiers(), TestRelationshipIdentifiers(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship, tagBindingInput

### Community 63 - "stubQuerier"
Cohesion: 0.03
Nodes (17): cancelRunStore, reliabilityStore, statusQuerier, captureJobs(), CreateContentParams, UpdateContentStatusParams, Job, Project (+9 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "net/http.Handler"
Cohesion: 0.05
Nodes (58): main(), run(), main(), run(), main(), run(), main(), run() (+50 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.05
Nodes (11): AccessConnection, AccessEdge, OrganizationConnection, OrganizationEdge, PageInfo, ProjectConnection, ProjectEdge, Query (+3 more)

### Community 70 - "testDeps"
Cohesion: 0.24
Nodes (9): TestUpstreamUIQueriesExecute(), runQuery(), testDeps(), TestFilterByTagAndJobStatus(), TestRelationshipIDsSelectOneRow(), TestWorkspacePageEndToEnd(), TestContractDiscoveryAndPing(), Discovery() (+1 more)

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 72 - ".createVCS"
Cohesion: 0.31
Nodes (8): vcsAttrs, strPtr(), strValue(), writeInUseConflict(), decodeVCSAttrs(), requiredVCSString(), validVCSEnum(), vcsEnumOrDefault()

### Community 73 - "newRedactor"
Cohesion: 0.50
Nodes (4): TestLogSinkRedactsSecrets(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "Resolver"
Cohesion: 0.12
Nodes (10): NotificationConfigurationResolver, OrganizationResolver, ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient() (+2 more)

### Community 78 - "decodeJSON"
Cohesion: 0.21
Nodes (21): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath() (+13 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "serveGlobalVar"
Cohesion: 0.39
Nodes (9): newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent(), TestGlobalVariableFromAnotherOrganizationIsNotFound(), TestListGlobalVariablesHidesSensitiveValues() (+1 more)

### Community 81 - "net/http.Handler"
Cohesion: 0.29
Nodes (3): requireUUIDParams(), templateTestAPI(), ExtractToken()

### Community 82 - "newTestExecutor"
Cohesion: 0.31
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "workspaces_test.go"
Cohesion: 0.16
Nodes (22): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+14 more)

### Community 88 - "Workspace"
Cohesion: 0.08
Nodes (7): workspaceStub, captureWorkspaces(), workspaceUpdateParams(), Workspace, CreateWorkspaceParams, UpdateWorkspaceParams, workspaceWithTags

### Community 92 - "Action"
Cohesion: 0.17
Nodes (4): actionStub, CreateActionParams, UpdateActionParams, Action

### Community 94 - "Node"
Cohesion: 0.42
Nodes (3): likePattern(), builder, Compare

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "Deps"
Cohesion: 0.29
Nodes (3): Deps, fakeLocker, JobLocker

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 133 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.13
Nodes (3): History, NotificationConfiguration, NotificationTrigger

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 352 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **63 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `Globalvar`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `Tag`, `stubQuerier`, `Handler`, `TFE Handler Tests`, `WriteResource`, `Vc`, `Template`, `Ssh`, `strPtr`, `context.Context`, `stubStore`, `testing.T`, `go_pkg_context`, `nestedStub`, `New`, `stubQuerier`, `net/http.Handler`, `.QueryRow`, `jobStub`, `net/http.Handler`, `newTestExecutor`, `Workspace`, `Action`, `Deps`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.08698207464002351 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `context.Context` to `ProjectAccessConnection`, `.marshalOString2ᚖstring`, `github.com/99designs/gqlgen/graphql.Marshaler`, `WorkspaceTagConnection`, `generated.go`, `Resolver`, `._NotificationConfiguration`, `github.com/99designs/gqlgen/graphql.CollectedField`, `._Job`, `.___Field`, `context.Context`, `NotificationTriggerConnection`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `.marshalOString2ᚖstring` be split into smaller, more focused modules?**
  _Cohesion score 0.12380952380952381 - nodes in this community are weakly interconnected._
- **Why does `newLogSink()` connect `testLogger` to `newRedactor`, `go_pkg_context`, `logSink`, `context.Context`, `.PublishCapped`, `.runApply`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._