# Graph Report - terrakube-go  (2026-11-23)

## Corpus Check
- 378 files · ~174,571 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 3212 nodes · 12745 edges · 166 communities (89 shown, 77 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 2453 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `31d29a35`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- Agent Workflows and Compatibility Docs
- Globalvar
- .Delete
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
- .processJob
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
- NewLayout
- .QueryRow
- models
- models_gen.go
- Team
- collections.go
- config.go
- Job
- generated.go
- where_test.go
- Resolver
- decodeJSON
- graphql/models.go
- dex_test.go
- api/main.go
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
- M4 run logs and artifacts Implementation Plan
- OrganizationConnection
- ui-rest-inventory.sh
- ui-smoke.sh
- jobStub
- .Ping
- capture.sh
- 0008_action.sql
- graphql/collections.go
- ._Project
- compare.sh
- start.sh
- .QueryRow
- WorkspaceTagConnection
- tarGzEntries
- TemplateConnection
- JobConnection
- testDeps
- net/http.Handler
- Backend
- newTestExecutor
- config_versions.go
- query.go
- stubStore
- notificationTestAPI
- get
- io.Reader
- .EnqueueJob
- .acquireConfiguration
- .Get
- transientError
- cleanKey
- process_unix.go
- Object
- Local

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 174 edges
3. `Queries` - 136 edges
4. `New()` - 128 edges
5. `API` - 124 edges
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
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `Layout`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfstate/layout.go
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

## Communities (166 total, 77 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (26): Write(), WriteResource(), configurationVersionResource(), projectResource(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID() (+18 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "Globalvar"
Cohesion: 0.14
Nodes (5): globalVarStub, variableStore, CreateGlobalVariableParams, UpdateGlobalVariableParams, Globalvar

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "Tag"
Cohesion: 0.07
Nodes (8): tagStub, Tag, Workspacetag, CreateTagParams, CreateWorkspaceTagParams, GetTagByNameParams, UpdateTagParams, UpdateWorkspaceTagValueParams

### Community 7 - "collectionStub"
Cohesion: 0.06
Nodes (9): collectionStub, CreateCollectionParams, CreateItemParams, CreateReferenceParams, UpdateCollectionParams, UpdateItemParams, Collection, Item (+1 more)

### Community 8 - "Handler"
Cohesion: 0.21
Nodes (7): isExternalUser(), matchesMD5(), New(), newUUID(), pathIDs(), sharedIDsContain(), Handler

### Community 9 - "TFE Handler Tests"
Cohesion: 0.10
Nodes (37): assertContractJSONAPI(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage(), stringPtr() (+29 more)

### Community 10 - "WriteResource"
Cohesion: 0.11
Nodes (13): API, Deps, listQuery, filterSortPage(), New(), param(), parseList(), unknownCollection() (+5 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.08
Nodes (34): Claims, Type, Classify(), decodeSecret(), ExtractToken(), NewValidator(), Sign(), TestClassify() (+26 more)

### Community 12 - "graphql/collections.go"
Cohesion: 0.14
Nodes (9): DB, failingDB, listArgs, afterString(), loadPage(), ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage() (+1 more)

### Community 13 - "Resource"
Cohesion: 0.10
Nodes (35): vcsAttrs, actionResource(), parseInclude(), strPtr(), strValue(), collectionResource(), itemResource(), referenceResource() (+27 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 16 - "Vc"
Cohesion: 0.13
Nodes (4): vcsStub, Vc, CreateVCSParams, UpdateVCSParams

### Community 17 - "stubQuerier"
Cohesion: 0.04
Nodes (10): stubQuerier, templateStub, demoOrganizations(), Organization, Template, CreateOrganizationParams, CreateTemplateParams, UpdateTemplateParams (+2 more)

### Community 18 - "Ssh"
Cohesion: 0.08
Nodes (6): sshAgentStub, Agent, Ssh, CreateAgentParams, CreateSSHKeyParams, UpdateSSHKeyParams

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.15
Nodes (30): replay, upstreamCapture, sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), assertSameDocumentShape(), assertSameResourceShape() (+22 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.06
Nodes (3): AccessEdge, NotificationTriggerEdge, WorkspaceEdge

### Community 25 - ".processJob"
Cohesion: 0.20
Nodes (4): stringPointer(), TestContentIDFromSource(), transient(), contentIDFromSource()

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
Cohesion: 0.17
Nodes (17): agentAttrs, atomicResult, collectionAttrs, globalVarAttrs, jobAttrs, notificationTriggerAttrs, organizationAttrs, projectAttrs (+9 more)

### Community 31 - ".runApply"
Cohesion: 0.19
Nodes (5): Executor, preparedWorkspace, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

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
Cohesion: 0.17
Nodes (6): NotificationConfigurationResolver, OrganizationResolver, QueryResolver, WorkspaceResolver, idsFrom(), workspaceEdges()

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
Nodes (73): redactor, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+65 more)

### Community 53 - "go_pkg_context"
Cohesion: 0.11
Nodes (19): stateMetadata, TestSanitizeGitURL(), contains(), envValue(), TestEnvironmentSeparatesSecrets(), sanitizeGitURL(), compareStreamIDs(), parseResumeID() (+11 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "logSink"
Cohesion: 0.20
Nodes (3): captureLogPublisher, LogPublisher, logSink

### Community 58 - ".handle"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 59 - "New"
Cohesion: 0.09
Nodes (27): TestWriteBackendOverride(), New(), TestDeleteRemovesAMessage(), TestHeadStreamReturnsTheOldestEntries(), testQueue(), isTerminal(), New(), outputOf() (+19 more)

### Community 60 - "WriteError"
Cohesion: 0.10
Nodes (36): actionAttrs, atomicOperation, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, decodeActionAttrs(), stringOr() (+28 more)

### Community 61 - "Querier"
Cohesion: 0.03
Nodes (10): historyStub, New(), AddNotificationConfigurationTemplateParams, CreateNotificationTriggerParams, Querier, DBTX, Queries, UpdateOrganizationParams (+2 more)

### Community 62 - "workspaceTagBindings"
Cohesion: 0.20
Nodes (10): nestedDocument, RelationshipIdentifiers(), RelationshipList(), TestRelationshipIdentifiers(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship (+2 more)

### Community 63 - "stubQuerier"
Cohesion: 0.06
Nodes (12): workspaceStub, workspaceUpdateParams(), CreateContentParams, UpdateContentStatusParams, Project, Workspace, CreateProjectParams, CreateWorkspaceParams (+4 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.10
Nodes (37): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+29 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "models_gen.go"
Cohesion: 0.07
Nodes (10): AccessConnection, ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, NotificationTriggerConnection, PageInfo, ProjectConnection, Query (+2 more)

### Community 70 - "Team"
Cohesion: 0.16
Nodes (4): teamStub, Team, CreateTeamParams, UpdateTeamParams

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 72 - "config.go"
Cohesion: 0.08
Nodes (36): main(), run(), main(), run(), main(), run(), Config, S3 (+28 more)

### Community 73 - "Job"
Cohesion: 0.05
Nodes (11): cancelRunStore, reliabilityStore, statusQuerier, captureJobs(), Job, UpdateJobStatusParams, Publisher, publishingStore (+3 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "Resolver"
Cohesion: 0.12
Nodes (8): ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 78 - "decodeJSON"
Cohesion: 0.16
Nodes (22): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath() (+14 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "dex_test.go"
Cohesion: 0.10
Nodes (22): stubDexVerifier, DexValidator, jwksServer, Verifier, stubDexVerifier, contains(), decodeJWKPart(), NewDexValidator() (+14 more)

### Community 82 - "api/main.go"
Cohesion: 0.14
Nodes (4): failingRow, graphQLResponse, ptr(), TestCollectionSpecRejectsInjection()

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "workspaces_test.go"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 88 - "ui.go"
Cohesion: 0.18
Nodes (11): New(), requestLogger(), writeJSON(), envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get() (+3 more)

### Community 92 - "Action"
Cohesion: 0.14
Nodes (6): actionFilter, actionStub, parseActionFilter(), CreateActionParams, UpdateActionParams, Action

### Community 94 - "Node"
Cohesion: 0.42
Nodes (3): likePattern(), builder, Compare

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "NewSubscriptions"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.25
Nodes (7): Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 127 - "graphql/collections.go"
Cohesion: 0.12
Nodes (19): collectionSpec, page, accessEdges(), column(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), predicate() (+11 more)

### Community 128 - "._Project"
Cohesion: 0.13
Nodes (3): ProjectAccessConnection, ProjectAccessEdge, ProjectEdge

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 144 - "testDeps"
Cohesion: 0.15
Nodes (13): appliedVersions(), applyMigration(), Migrate(), path(), TestUpstreamUIQueriesExecute(), runQuery(), testDeps(), TestFilterByTagAndJobStatus() (+5 more)

### Community 145 - "net/http.Handler"
Cohesion: 0.08
Nodes (32): main(), run(), Deps, requireUUIDParams(), newTeamStub(), teamTestAPI(), TestTeamJSONAPI(), newTemplateStub() (+24 more)

### Community 146 - "Backend"
Cohesion: 0.18
Nodes (10): Validator, Deps, Deps, Pinger, Backend, Layout, Deps, JobQueue (+2 more)

### Community 147 - "newTestExecutor"
Cohesion: 0.16
Nodes (19): newFakeLocker(), TestHandleReleasesTheJobLock(), TestHandleSkipsAJobHeldByAnotherExecutor(), TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestStaleClaimIdleOutlivesTheJobTimeout(), New(), TestAckDeletesTheJobMessage() (+11 more)

### Community 150 - "config_versions.go"
Cohesion: 0.17
Nodes (4): configurationVersionAttributes, runAttributes, stateDocumentMetadata, stateVersionAttributes

### Community 151 - "query.go"
Cohesion: 0.32
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 152 - "stubStore"
Cohesion: 0.18
Nodes (3): GetArchiveRow, UpdateHistoryOutputParams, stubStore

### Community 153 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 154 - "get"
Cohesion: 0.53
Nodes (6): get(), newUI(), TestDeepLinksServeTheIndex(), TestEnvConfigIsNeverCached(), TestMissingAssetsAreNotFound(), TestServesTheIndexAndAssets()

### Community 157 - ".acquireConfiguration"
Cohesion: 0.28
Nodes (7): gitConfig, TestGitSourceDetection(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL(), runGit()

### Community 162 - "cleanKey"
Cohesion: 0.25
Nodes (6): cleanKey(), pathClean(), TestCleanKey(), NewS3(), New(), S3

### Community 163 - "process_unix.go"
Cohesion: 0.25
Nodes (4): configureProcessGroup(), killProcessGroup(), configureProcessGroup(), killProcessGroup()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 360 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **77 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `Globalvar`, `.QueryRow`, `Tag`, `collectionStub`, `TFE Handler Tests`, `WriteResource`, `Step`, `Vc`, `stubQuerier`, `Ssh`, `Backend`, `net/http.Handler`, `strPtr`, `newTestExecutor`, `stubStore`, `notificationTestAPI`, `state.sql.go`, `testing.T`, `context.Context`, `nestedStub`, `New`, `stubQuerier`, `.QueryRow`, `Team`, `Job`, `Action`, `.CreateNotificationConfiguration`, `jobStub`?**
  _High betweenness centrality (0.053) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07931199235547062 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `executionContext` to `._Project`, `context.Context`, `.marshalOString2ᚖstring`, `models_gen.go`, `OrganizationConnection`, `WorkspaceTagConnection`, `generated.go`, `Resolver`, `._NotificationConfiguration`, `TemplateConnection`, `JobConnection`, `._Organization`, `.___EnumValue`, `.marshalOBoolean2bool`, `github.com/99designs/gqlgen/graphql.Marshaler`, `context.Context`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `Workspace` connect `stubQuerier` to `net/http.ResponseWriter`, `Tag`, `collectionStub`, `Job`, `WriteResource`, `Resource`, `Vc`, `stubQuerier`, `Ssh`, `strPtr`, `go_pkg_context`, `.processJob`, `WriteError`, `.acquireConfiguration`, `.runApply`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Should `Globalvar` be split into smaller, more focused modules?**
  _Cohesion score 0.13666666666666666 - nodes in this community are weakly interconnected._