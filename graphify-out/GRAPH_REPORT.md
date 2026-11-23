# Graph Report - terrakube-go  (2026-11-23)

## Corpus Check
- 387 files · ~179,114 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 3300 nodes · 12826 edges · 158 communities (91 shown, 67 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 2504 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7f0ff46d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- context.Context
- Agent Workflows and Compatibility Docs
- Globalvar
- io.Reader
- Core Database Schema
- tagStub
- collectionStub
- Handler
- TFE Handler Tests
- WriteResource
- State Backend Tests
- graphql/collections.go
- Resource
- steps
- .QueryRow
- Vc
- stubQuerier
- Ssh
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- Identity
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- api/main.go
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
- Pat
- executor-e2e.sh
- logSink
- context.Context
- nestedStub
- .handle
- .streamStepLog
- WriteError
- Querier
- workspaceTagBindings
- stubQuerier
- dex-config.yaml
- NewLayout
- .QueryRow
- models
- github.com/vektah/gqlparser/v2/ast.SelectionSet
- Team
- collections.go
- config.go
- Job
- generated.go
- where_test.go
- streamEvents
- templateStub
- decodeJSON
- graphql/models.go
- newJWKSServer
- StringOrInt
- Migrate
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- Collection
- .updateWorkspaceRelationship
- Action
- .CreateNotificationConfiguration
- Node
- environment
- NewSubscriptions
- Tag
- newService
- .CreateVariable
- NotificationConfigurationConnection
- ui-rest-inventory.sh
- ui-smoke.sh
- SetWorkspaceLastJobParams
- capture.sh
- 0008_action.sql
- graphql/collections.go
- ._Project
- compare.sh
- start.sh
- .QueryRow
- New
- tarGzEntries
- .createVCS
- JobConnection
- net/http.Handler
- run
- ParseSort
- newTestExecutor
- .createPAT
- query.go
- .EnqueueJob
- .cloneConfiguration

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Queries` - 141 edges
3. `Querier` - 139 edges
4. `New()` - 128 edges
5. `API` - 127 edges
6. `WriteError()` - 99 edges
7. `stubQuerier` - 95 edges
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

## Communities (158 total, 67 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (28): Write(), WriteResource(), configurationVersionResource(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource() (+20 more)

### Community 1 - "context.Context"
Cohesion: 0.07
Nodes (4): artifactStore, historyStub, variableStore, stubQuerier

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "Globalvar"
Cohesion: 0.13
Nodes (4): globalVarStub, CreateGlobalVariableParams, UpdateGlobalVariableParams, Globalvar

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (11): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+3 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "tagStub"
Cohesion: 0.09
Nodes (4): tagStub, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, Workspacetag

### Community 7 - "collectionStub"
Cohesion: 0.08
Nodes (5): collectionStub, CreateItemParams, CreateReferenceParams, UpdateItemParams, Item

### Community 8 - "Handler"
Cohesion: 0.08
Nodes (21): Claims, Type, Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken(), claimsFromContext() (+13 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.10
Nodes (37): assertContractJSONAPI(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage(), stringPtr() (+29 more)

### Community 10 - "WriteResource"
Cohesion: 0.12
Nodes (10): API, listQuery, filterSortPage(), New(), param(), parseList(), unknownCollection(), Page (+2 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.10
Nodes (26): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+18 more)

### Community 12 - "graphql/collections.go"
Cohesion: 0.14
Nodes (9): DB, failingDB, listArgs, afterString(), loadPage(), ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage() (+1 more)

### Community 13 - "Resource"
Cohesion: 0.11
Nodes (32): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), jobResource(), stepResource() (+24 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 15 - ".QueryRow"
Cohesion: 0.05
Nodes (4): GetProjectByNameParams, History, Step, UpdateProjectParams

### Community 16 - "Vc"
Cohesion: 0.11
Nodes (4): vcsStub, CreateVCSParams, UpdateVCSParams, Vc

### Community 17 - "stubQuerier"
Cohesion: 0.03
Nodes (4): stubQuerier, CreateOrganizationParams, NotificationConfiguration, NotificationTrigger

### Community 18 - "Ssh"
Cohesion: 0.06
Nodes (8): sshAgentStub, CreateAgentParams, CreateSSHKeyParams, UpdateSSHKeyParams, Agent, Ssh, Template, Workspace

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.13
Nodes (30): replay, upstreamCapture, sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), assertSameDocumentShape(), assertSameResourceShape() (+22 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 25 - "Identity"
Cohesion: 0.09
Nodes (19): Deps, PATService, stubDexVerifier, stubPATService, DexValidator, fakeVerifier, Identity, Verifier (+11 more)

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

### Community 30 - "api/main.go"
Cohesion: 0.06
Nodes (65): actionFilter, agentAttrs, atomicResult, collectionAttrs, createPATRequest, globalVarAttrs, jobAttrs, notificationTriggerAttrs (+57 more)

### Community 31 - ".runApply"
Cohesion: 0.14
Nodes (7): Executor, preparedWorkspace, isCancelledStatus(), stringPointer(), TestBinarySelection(), transient(), contentIDFromSource()

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
Cohesion: 0.12
Nodes (8): NotificationConfigurationResolver, OrganizationResolver, ProjectResolver, QueryResolver, Resolver, ResolverRoot, WorkspaceResolver, idsFrom()

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
Cohesion: 0.04
Nodes (81): newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes(), newTestAPI() (+73 more)

### Community 53 - "Pat"
Cohesion: 0.17
Nodes (5): fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "logSink"
Cohesion: 0.15
Nodes (4): captureLogPublisher, LogPublisher, logSink, baseEnvironment()

### Community 58 - ".handle"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 59 - ".streamStepLog"
Cohesion: 0.17
Nodes (11): Validator, isTerminal(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData() (+3 more)

### Community 60 - "WriteError"
Cohesion: 0.10
Nodes (37): actionAttrs, atomicOperation, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, decodeActionAttrs(), stringOr() (+29 more)

### Community 61 - "Querier"
Cohesion: 0.04
Nodes (8): New(), AddNotificationConfigurationTemplateParams, CreateNotificationTriggerParams, DBTX, Querier, Queries, UpdateHistoryOutputParams, UpdateOrganizationParams

### Community 62 - "workspaceTagBindings"
Cohesion: 0.25
Nodes (7): nestedDocument, RelationshipIdentifiers(), RelationshipList(), identifierStringAttribute(), workspaceTagBindings(), RequestRelationship, ResourceIdentifier

### Community 63 - "stubQuerier"
Cohesion: 0.05
Nodes (6): CreateContentParams, UpdateContentStatusParams, CreateProjectParams, CreateTagParams, Content, stubQuerier

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.06
Nodes (52): redactor, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+44 more)

### Community 67 - ".QueryRow"
Cohesion: 0.05
Nodes (3): Reference, TempArchive, UpdateJobPlanChangesParams

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "github.com/vektah/gqlparser/v2/ast.SelectionSet"
Cohesion: 0.05
Nodes (14): AccessConnection, AccessEdge, NotificationTriggerConnection, NotificationTriggerEdge, OrganizationConnection, OrganizationEdge, PageInfo, Query (+6 more)

### Community 70 - "Team"
Cohesion: 0.13
Nodes (4): teamStub, CreateTeamParams, UpdateTeamParams, Team

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 72 - "config.go"
Cohesion: 0.06
Nodes (29): Config, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1` (+21 more)

### Community 73 - "Job"
Cohesion: 0.06
Nodes (9): jobStub, cancelRunStore, reliabilityStore, statusQuerier, UpdateJobStatusParams, Publisher, publishingStore, stubQuerier (+1 more)

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (5): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), NewExecutableSchema(), newExecutionContext(), TestSchemaLoads()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "streamEvents"
Cohesion: 0.24
Nodes (5): SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 77 - "templateStub"
Cohesion: 0.13
Nodes (3): templateStub, CreateTemplateParams, UpdateTemplateParams

### Community 78 - "decodeJSON"
Cohesion: 0.16
Nodes (22): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath() (+14 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.24
Nodes (14): jwksServer, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken(), TestDexValidatorAudienceValidationIsOptIn() (+6 more)

### Community 81 - "StringOrInt"
Cohesion: 0.08
Nodes (4): ComplexityRoot, ProjectConnection, ProjectEdge, WorkspaceConnection

### Community 82 - "Migrate"
Cohesion: 0.18
Nodes (12): main(), run(), Locker, NewPool(), NewLocker(), TestLockerAllowsDifferentJobs(), TestLockerExcludesASecondHolder(), testPool() (+4 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "workspaces_test.go"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 88 - "Collection"
Cohesion: 0.16
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 92 - "Action"
Cohesion: 0.13
Nodes (4): actionStub, CreateActionParams, UpdateActionParams, Action

### Community 94 - "Node"
Cohesion: 0.42
Nodes (3): likePattern(), builder, Compare

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "NewSubscriptions"
Cohesion: 0.33
Nodes (3): fakeEvents, Event, Bus

### Community 101 - "Tag"
Cohesion: 0.16
Nodes (3): GetTagByNameParams, UpdateTagParams, Tag

### Community 102 - "newService"
Cohesion: 0.23
Nodes (14): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+6 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 127 - "graphql/collections.go"
Cohesion: 0.15
Nodes (14): collectionSpec, page, accessEdges(), column(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), predicate() (+6 more)

### Community 133 - ".QueryRow"
Cohesion: 0.07
Nodes (3): workspaceStub, CreateWorkspaceParams, UpdateWorkspaceParams

### Community 135 - "New"
Cohesion: 0.21
Nodes (12): TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput(), TestGetStepLogPrefersArchivedObject(), TestGetStepLogRejectsWrongScope(), TestGetStepLogRequiresToken() (+4 more)

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 142 - ".createVCS"
Cohesion: 0.39
Nodes (6): vcsAttrs, strPtr(), strValue(), decodeVCSAttrs(), validVCSEnum(), vcsEnumOrDefault()

### Community 144 - "net/http.Handler"
Cohesion: 0.08
Nodes (29): requireUUIDParams(), newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates(), NewPAT(), TestPATRoutesRequireAuthenticationAndASecret() (+21 more)

### Community 145 - "run"
Cohesion: 0.06
Nodes (47): main(), run(), main(), run(), main(), run(), Deps, EventSource (+39 more)

### Community 146 - "ParseSort"
Cohesion: 0.40
Nodes (4): ParseSort(), TestOrderByRendersSQL(), TestParseSort(), TestParseSortRejectsUnknownField()

### Community 147 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 151 - "query.go"
Cohesion: 0.33
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 157 - ".cloneConfiguration"
Cohesion: 0.50
Nodes (4): TestSanitizeGitURL(), gitEnv(), runGit(), sanitizeGitURL()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 367 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **67 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `executionContext` connect `context.Context` to `._Project`, `context.Context`, `context.Context`, `github.com/vektah/gqlparser/v2/ast.SelectionSet`, `StringOrInt`, `generated.go`, `NotificationConfigurationConnection`, `JobConnection`, `StringOrInt`, `.marshalOBoolean2bool`, `.___EnumValue`, `github.com/99designs/gqlgen/graphql.Marshaler`, `Job`, `._Template`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07525195968645017 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.ResponseWriter` to `api/main.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.06604324956165984 - nodes in this community are weakly interconnected._
- **Why does `newLogSink()` connect `NewLayout` to `context.Context`, `run`, `logSink`, `api/main.go`, `.runApply`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._