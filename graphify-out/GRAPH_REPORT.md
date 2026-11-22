# Graph Report - terrakube-go  (2026-11-22)

## Corpus Check
- 304 files · ~157,864 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 2934 nodes · 11432 edges · 150 communities (85 shown, 65 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 2205 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a64a7523`
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
- executor/main.go
- TemplateConnection
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
- Collection
- testLogger
- nestedStub
- .handle
- net/http.Handler
- WriteError
- Querier
- executor_test.go
- stubQuerier
- dex-config.yaml
- net/http.Handler
- models
- github.com/99designs/gqlgen/graphql.Marshaler
- .EnqueueJob
- collections.go
- api/main.go
- graphql_test.go
- generated.go
- where_test.go
- Resolver
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
- Workspace
- Layout
- stubQuerier
- logSink
- where.go
- environment
- streamEvents
- loadPage
- Deps
- ui-rest-inventory.sh
- ui-smoke.sh
- sort.go
- .Ping
- capture.sh
- 0008_action.sql
- NotificationTriggerConnection
- ProjectAccessConnection
- compare.sh
- workspaceTagBindings
- start.sh
- Queries
- WorkspaceTagConnection
- Claims
- M4 run logs and artifacts Implementation Plan
- tarGzEntries
- AccessConnection
- .streamStepLog
- Classify
- .updateWorkspaceRelationship
- query.go
- actionFilter

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 137 edges
3. `Queries` - 106 edges
4. `New()` - 97 edges
5. `API` - 89 edges
6. `WriteError()` - 79 edges
7. `stubQuerier` - 78 edges
8. `API` - 66 edges
9. `stubQuerier` - 58 edges
10. `param()` - 50 edges

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

## Communities (150 total, 65 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (15): WriteResource(), configurationVersionResource(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), normalizeTagValue(), tagValuesEqual() (+7 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "Globalvar"
Cohesion: 0.14
Nodes (5): globalVarStub, variableStore, CreateGlobalVariableParams, UpdateGlobalVariableParams, Globalvar

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (11): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+3 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (50): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+42 more)

### Community 6 - "Tag"
Cohesion: 0.08
Nodes (8): tagStub, Tag, Workspacetag, CreateTagParams, CreateWorkspaceTagParams, GetTagByNameParams, UpdateTagParams, UpdateWorkspaceTagValueParams

### Community 7 - "collectionStub"
Cohesion: 0.08
Nodes (8): collectionStub, CreateItemParams, CreateReferenceParams, UpdateItemParams, Item, Reference, NotificationConfiguration, NotificationTrigger

### Community 8 - "Handler"
Cohesion: 0.25
Nodes (5): matchesMD5(), New(), newUUID(), pathIDs(), Handler

### Community 9 - "TFE Handler Tests"
Cohesion: 0.09
Nodes (41): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+33 more)

### Community 10 - "WriteResource"
Cohesion: 0.14
Nodes (12): API, Deps, listQuery, filterSortPage(), New(), param(), parseList(), unknownCollection() (+4 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.14
Nodes (23): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+15 more)

### Community 12 - "graphql/collections.go"
Cohesion: 0.15
Nodes (14): collectionSpec, page, accessEdges(), column(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), predicate() (+6 more)

### Community 13 - "Resource"
Cohesion: 0.08
Nodes (38): statusQuerier, actionResource(), parseInclude(), strValue(), collectionResource(), itemResource(), referenceResource(), globalVarResource() (+30 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 16 - "executor/main.go"
Cohesion: 0.10
Nodes (18): transientError, TestSanitizeGitURL(), contains(), envValue(), TestEnvironmentSeparatesSecrets(), configureProcessGroup(), killProcessGroup(), configureProcessGroup() (+10 more)

### Community 18 - "Ssh"
Cohesion: 0.09
Nodes (6): sshAgentStub, Agent, Ssh, CreateAgentParams, CreateSSHKeyParams, UpdateSSHKeyParams

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.17
Nodes (22): replay, upstreamCapture, sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), assertSameDocumentShape(), assertSameResourceShape() (+14 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 24 - "context.Context"
Cohesion: 0.06
Nodes (3): ComplexityRoot, ProjectConnection, WorkspaceConnection

### Community 25 - ".acquireConfiguration"
Cohesion: 0.19
Nodes (8): gitConfig, transient(), contentIDFromSource(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL(), runGit()

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
Cohesion: 0.17
Nodes (5): NotificationConfigurationResolver, OrganizationResolver, QueryResolver, WorkspaceResolver, idsFrom()

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
Nodes (77): redactor, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+69 more)

### Community 53 - "go_pkg_context"
Cohesion: 0.15
Nodes (19): agentAttrs, atomicResult, collectionAttrs, globalVarAttrs, jobAttrs, organizationAttrs, projectAttrs, sshAttrs (+11 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Collection"
Cohesion: 0.17
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 56 - "testLogger"
Cohesion: 0.08
Nodes (47): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+39 more)

### Community 58 - ".handle"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 59 - "net/http.Handler"
Cohesion: 0.21
Nodes (12): TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput(), TestGetStepLogPrefersArchivedObject(), TestGetStepLogRejectsWrongScope(), TestGetStepLogRequiresToken() (+4 more)

### Community 60 - "WriteError"
Cohesion: 0.13
Nodes (29): actionAttrs, atomicOperation, itemAttrs, decodeActionAttrs(), stringOr(), boolValue(), requireType(), strPtr() (+21 more)

### Community 61 - "Querier"
Cohesion: 0.05
Nodes (4): artifactStore, historyStub, Querier, stubQuerier

### Community 62 - "executor_test.go"
Cohesion: 0.12
Nodes (10): stateMetadata, NewS3(), Backend, New(), configurationVersionAttributes, runAttributes, stateDocumentMetadata, stateOutput (+2 more)

### Community 63 - "stubQuerier"
Cohesion: 0.04
Nodes (15): cancelRunStore, reliabilityStore, captureJobs(), CreateContentParams, UpdateContentStatusParams, Content, History, Job (+7 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "net/http.Handler"
Cohesion: 0.07
Nodes (35): main(), run(), Deps, Deps, requireUUIDParams(), NewExecutableSchema(), TestSchemaLoads(), bearerMatches() (+27 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - "api/main.go"
Cohesion: 0.11
Nodes (20): main(), run(), Config, S3, env(), envBool(), envFirst(), envInt64() (+12 more)

### Community 73 - "graphql_test.go"
Cohesion: 0.11
Nodes (20): main(), run(), main(), run(), Locker, failingRow, graphQLResponse, NewPool() (+12 more)

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "Resolver"
Cohesion: 0.13
Nodes (8): ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 78 - "decodeJSON"
Cohesion: 0.11
Nodes (31): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent() (+23 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

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

### Community 88 - "Workspace"
Cohesion: 0.06
Nodes (6): workspaceStub, Project, Workspace, CreateWorkspaceParams, UpdateWorkspaceParams, workspaceWithTags

### Community 89 - "Layout"
Cohesion: 0.20
Nodes (6): Validator, Layout, Deps, JobQueue, Deps, Deps

### Community 92 - "stubQuerier"
Cohesion: 0.07
Nodes (8): actionStub, stubQuerier, demoOrganizations(), CreateActionParams, UpdateActionParams, Action, Organization, CreateOrganizationParams

### Community 93 - "logSink"
Cohesion: 0.15
Nodes (4): captureLogPublisher, LogPublisher, logSink, baseEnvironment()

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "loadPage"
Cohesion: 0.18
Nodes (5): DB, failingDB, listArgs, afterString(), loadPage()

### Community 103 - "Deps"
Cohesion: 0.29
Nodes (3): Deps, fakeLocker, JobLocker

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 121 - "sort.go"
Cohesion: 0.33
Nodes (6): Order, OrderBy(), ParseSort(), TestOrderByRendersSQL(), TestParseSort(), TestParseSortRejectsUnknownField()

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 130 - "workspaceTagBindings"
Cohesion: 0.22
Nodes (8): nestedDocument, RelationshipIdentifiers(), TestRelationshipIdentifiers(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship, tagBindingInput

### Community 133 - "Queries"
Cohesion: 0.04
Nodes (9): New(), CreateProjectParams, DBTX, GetProjectByNameParams, Queries, UpdateJobPlanChangesParams, UpdateOrganizationParams, UpdateProjectParams (+1 more)

### Community 136 - "Claims"
Cohesion: 0.24
Nodes (6): Claims, TestEffectiveStepIDFromJTI(), claimsFromContext(), isActiveStepStatus(), isExternalUser(), sharedIDsContain()

### Community 139 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.18
Nodes (9): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+1 more)

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 142 - ".streamStepLog"
Cohesion: 0.26
Nodes (9): isTerminal(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData(), writeEvent() (+1 more)

### Community 143 - "Classify"
Cohesion: 0.22
Nodes (7): Type, Classify(), ExtractToken(), TestClassify(), TestExtractToken(), withClaims(), withExternalUser()

### Community 145 - "query.go"
Cohesion: 0.38
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 349 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **65 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `Globalvar`, `Queries`, `Tag`, `collectionStub`, `TFE Handler Tests`, `WriteResource`, `Resource`, `Ssh`, `strPtr`, `stubStore`, `testing.T`, `go_pkg_context`, `Collection`, `nestedStub`, `net/http.Handler`, `stubQuerier`, `net/http.Handler`, `jobStub`, `Step`, `newTestExecutor`, `Workspace`, `Layout`, `stubQuerier`, `Deps`?**
  _High betweenness centrality (0.043) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.09769555339175592 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.ResponseWriter` to `Layout`, `Resource`, `executor_test.go`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `executionContext` be split into smaller, more focused modules?**
  _Cohesion score 0.08926553672316384 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `executionContext` to `ProjectAccessConnection`, `context.Context`, `Job`, `github.com/99designs/gqlgen/graphql.Marshaler`, `WorkspaceTagConnection`, `generated.go`, `Resolver`, `AccessConnection`, `StringOrInt`, `TemplateConnection`, `github.com/99designs/gqlgen/graphql.CollectedField`, `models_gen.go`, `.childFields___InputValue`, `context.Context`, `NotificationTriggerConnection`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._