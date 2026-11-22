# Graph Report - terrakube-go  (2026-11-22)

## Corpus Check
- 351 files · ~167,141 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 3084 nodes · 12160 edges · 145 communities (80 shown, 65 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 2339 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1924e6e3`
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
- executionContext
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
- net/http.Handler
- .QueryRow
- models
- github.com/99designs/gqlgen/graphql.Marshaler
- Team
- collections.go
- Load
- newRedactor
- generated.go
- where_test.go
- Resolver
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- decodeJSON
- graphql/models.go
- Collection
- newTestExecutor
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- ui.go
- .updateWorkspaceRelationship
- Action
- Node
- environment
- streamEvents
- log/slog.Logger
- loadPage
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
- tarGzEntries

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 160 edges
3. `Queries` - 123 edges
4. `New()` - 115 edges
5. `API` - 109 edges
6. `stubQuerier` - 98 edges
7. `WriteError()` - 87 edges
8. `API` - 66 edges
9. `param()` - 59 edges
10. `stubQuerier` - 55 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `Storage`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/config/config.go
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

## Communities (145 total, 65 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (20): WriteResource(), configurationVersionResource(), parseWorkspaceQuery(), approvalGateStatus(), failedJobStatus(), parseRunID(), stepStatusTFE(), stepStatusValue() (+12 more)

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
Cohesion: 0.13
Nodes (15): Claims, Type, Classify(), ExtractToken(), claimsFromContext(), isActiveStepStatus(), isExternalUser(), withClaims() (+7 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.09
Nodes (41): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+33 more)

### Community 10 - "WriteResource"
Cohesion: 0.11
Nodes (18): API, Deps, listQuery, filterSortPage(), New(), param(), parseList(), strPtr() (+10 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.16
Nodes (21): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+13 more)

### Community 12 - "graphql/collections.go"
Cohesion: 0.11
Nodes (20): collectionSpec, listArgs, page, accessEdges(), afterString(), column(), jobEdges(), notificationConfigurationEdges() (+12 more)

### Community 13 - "Resource"
Cohesion: 0.09
Nodes (38): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), jobResource(), stepResource() (+30 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 15 - "Step"
Cohesion: 0.06
Nodes (6): reliabilityStore, Step, Publisher, publishingStore, stubQuerier, mutableQuerier

### Community 16 - "Vc"
Cohesion: 0.13
Nodes (4): vcsStub, Vc, CreateVCSParams, UpdateVCSParams

### Community 17 - "stubQuerier"
Cohesion: 0.06
Nodes (8): stubQuerier, templateStub, demoOrganizations(), Organization, Template, CreateOrganizationParams, CreateTemplateParams, UpdateTemplateParams

### Community 18 - "Ssh"
Cohesion: 0.08
Nodes (6): sshAgentStub, Agent, Ssh, CreateAgentParams, CreateSSHKeyParams, UpdateSSHKeyParams

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.10
Nodes (33): replay, upstreamCapture, sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), newTeamStub(), TestTeamJSONAPI() (+25 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 24 - "context.Context"
Cohesion: 0.06
Nodes (3): artifactStore, JobConnection, JobEdge

### Community 25 - ".acquireConfiguration"
Cohesion: 0.19
Nodes (6): preparedWorkspace, TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv(), runGit()

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

### Community 37 - "StringOrInt"
Cohesion: 0.16
Nodes (5): QueryResolver, WorkspaceResolver, idsFrom(), projectEdges(), workspaceEdges()

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
Cohesion: 0.04
Nodes (82): newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes(), newTestAPI() (+74 more)

### Community 53 - "go_pkg_context"
Cohesion: 0.06
Nodes (57): agentAttrs, atomicResult, collectionAttrs, globalVarAttrs, jobAttrs, organizationAttrs, projectAttrs, sshAttrs (+49 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "logSink"
Cohesion: 0.20
Nodes (3): captureLogPublisher, LogPublisher, logSink

### Community 56 - "testLogger"
Cohesion: 0.08
Nodes (45): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+37 more)

### Community 58 - ".handle"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 59 - "New"
Cohesion: 0.11
Nodes (23): TestWriteBackendOverride(), isTerminal(), New(), outputOf(), stepOf(), streamKey(), newHandler(), stringPtr() (+15 more)

### Community 60 - "WriteError"
Cohesion: 0.10
Nodes (38): actionAttrs, atomicOperation, itemAttrs, teamAttrs, templateAttrs, vcsAttrs, decodeActionAttrs(), stringOr() (+30 more)

### Community 61 - "Querier"
Cohesion: 0.03
Nodes (9): historyStub, New(), Querier, DBTX, Queries, TempArchive, UpdateOrganizationParams, Variable (+1 more)

### Community 62 - "workspaceTagBindings"
Cohesion: 0.25
Nodes (6): nestedDocument, identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship, tagBindingInput

### Community 63 - "stubQuerier"
Cohesion: 0.04
Nodes (15): workspaceStub, cancelRunStore, captureJobs(), CreateContentParams, UpdateContentStatusParams, Job, Project, Workspace (+7 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "net/http.Handler"
Cohesion: 0.05
Nodes (49): main(), run(), main(), run(), main(), run(), main(), run() (+41 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.04
Nodes (18): AccessConnection, AccessEdge, NotificationConfigurationConnection, NotificationConfigurationEdge, NotificationTriggerConnection, NotificationTriggerEdge, OrganizationConnection, OrganizationEdge (+10 more)

### Community 70 - "Team"
Cohesion: 0.17
Nodes (4): teamStub, Team, CreateTeamParams, UpdateTeamParams

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - "Load"
Cohesion: 0.13
Nodes (18): Config, Storage, env(), envBool(), envDuration(), envFirst(), envInt64(), envList() (+10 more)

### Community 73 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (4): ComplexityRoot, DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

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

### Community 80 - "Collection"
Cohesion: 0.17
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 82 - "newTestExecutor"
Cohesion: 0.27
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

### Community 88 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 92 - "Action"
Cohesion: 0.14
Nodes (6): actionFilter, actionStub, parseActionFilter(), CreateActionParams, UpdateActionParams, Action

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.33
Nodes (3): fakeEvents, Event, Bus

### Community 101 - "log/slog.Logger"
Cohesion: 0.09
Nodes (17): Validator, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+9 more)

### Community 103 - "loadPage"
Cohesion: 0.24
Nodes (3): DB, failingDB, loadPage()

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 127 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 133 - ".QueryRow"
Cohesion: 0.07
Nodes (5): GetProjectByNameParams, History, NotificationConfiguration, NotificationTrigger, UpdateProjectParams

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 354 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **65 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `Globalvar`, `.QueryRow`, `Tag`, `collectionStub`, `TFE Handler Tests`, `WriteResource`, `Step`, `Vc`, `stubQuerier`, `Ssh`, `strPtr`, `context.Context`, `stubStore`, `testing.T`, `go_pkg_context`, `nestedStub`, `New`, `stubQuerier`, `net/http.Handler`, `.QueryRow`, `Team`, `github.com/aminbenmansour/terrakube-go/internal/store.Job`, `Collection`, `newTestExecutor`, `Action`, `log/slog.Logger`, `SetWorkspaceLastJobParams`?**
  _High betweenness centrality (0.055) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.0876750700280112 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `executionContext` to `._Project`, `context.Context`, `.marshalOString2ᚖstring`, `github.com/99designs/gqlgen/graphql.Marshaler`, `WorkspaceTagConnection`, `generated.go`, `Resolver`, `.___Directive`, `github.com/99designs/gqlgen/graphql.CollectedField`, `.___EnumValue`, `context.Context`, `.___InputValue`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Should `.marshalOString2ᚖstring` be split into smaller, more focused modules?**
  _Cohesion score 0.12100840336134454 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.ResponseWriter` to `log/slog.Logger`, `go_pkg_context`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._