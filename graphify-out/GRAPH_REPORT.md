# Graph Report - terrakube-go  (2026-11-23)

## Corpus Check
- 399 files · ~187,321 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 3464 nodes · 13481 edges · 174 communities (94 shown, 80 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 2586 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `be88ebe5`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- Agent Workflows and Compatibility Docs
- Globalvar
- io.Reader
- Core Database Schema
- tagStub
- collectionStub
- Handler
- TFE Handler Tests
- WriteResource
- CliAuthSession
- graphql/collections.go
- Resource
- steps
- .QueryRow
- Vc
- stubQuerier
- Ssh
- context.Context
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- api/main.go
- Identity
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- go_pkg_context
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
- executionContext
- nestedStub
- .handle
- .streamStepLog
- WriteError
- Querier
- state.sql.go
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
- Resolver
- templateStub
- decodeJSON
- graphql/models.go
- newJWKSServer
- TemplateConnection
- config_versions.go
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- Collection
- .updateWorkspaceRelationship
- .Query
- .CreateNotificationConfiguration
- Node
- environment
- NewSubscriptions
- newTestService
- newService
- NewSubscriptions
- api_test.go
- ui-rest-inventory.sh
- ui-smoke.sh
- jobStub
- capture.sh
- 0008_action.sql
- graphql/collections.go
- compare.sh
- start.sh
- .QueryRow
- New
- tarGzEntries
- .createVCS
- net/http.Handler
- run
- Dex
- newTestExecutor
- .createPAT
- query.go
- WorkspaceTagConnection
- time.Time
- .EnqueueJob
- .cloneConfiguration
- ui.go
- AccessConnection
- NotificationTriggerConnection
- Project
- CreateAgentParams
- auth.DexValidator
- .CreateNotificationTrigger
- .Run
- notificationTestAPI
- state_outputs.go
- .Ping
- .SoftDeletePat
- actionFilter
- transientError

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 152 edges
3. `Queries` - 150 edges
4. `API` - 127 edges
5. `New()` - 126 edges
6. `WriteError()` - 99 edges
7. `stubQuerier` - 95 edges
8. `API` - 66 edges
9. `param()` - 64 edges
10. `ParseFieldsets()` - 59 edges

## Surprising Connections (you probably didn't know these)
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `Layout`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfstate/layout.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `New()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfoutput/handler.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Generated Code and Migration Policy` --semantically_similar_to--> `sqlc Workflow`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .opencode/skills/sqlc-workflow/SKILL.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (174 total, 80 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (14): statusQuerier, Write(), WriteResource(), configurationVersionResource(), stepStatusTFE(), stepStatusValue(), stringPointer(), decodeTFEJSON() (+6 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "Globalvar"
Cohesion: 0.13
Nodes (4): globalVarStub, CreateGlobalVariableParams, UpdateGlobalVariableParams, Globalvar

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (12): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+4 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "tagStub"
Cohesion: 0.08
Nodes (5): tagStub, CreateWorkspaceTagParams, UpdateTagParams, UpdateWorkspaceTagValueParams, Workspacetag

### Community 7 - "collectionStub"
Cohesion: 0.11
Nodes (3): collectionStub, CreateItemParams, CreateReferenceParams

### Community 8 - "Handler"
Cohesion: 0.05
Nodes (47): Claims, Type, Validator, Classify(), decodeSecret(), ExtractToken(), NewValidator(), Sign() (+39 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.09
Nodes (36): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+28 more)

### Community 10 - "WriteResource"
Cohesion: 0.12
Nodes (10): API, listQuery, filterSortPage(), New(), param(), parseList(), writeInUseConflict(), Page (+2 more)

### Community 11 - "CliAuthSession"
Cohesion: 0.05
Nodes (24): Deps, fakeStore, Service, Store, TokenResponse, ConsentPage(), ErrorPage(), effectivePort() (+16 more)

### Community 12 - "graphql/collections.go"
Cohesion: 0.14
Nodes (9): DB, failingDB, listArgs, afterString(), loadPage(), ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage() (+1 more)

### Community 13 - "Resource"
Cohesion: 0.08
Nodes (43): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), jobResource(), stepResource() (+35 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 16 - "Vc"
Cohesion: 0.11
Nodes (4): vcsStub, CreateVCSParams, UpdateVCSParams, Vc

### Community 17 - "stubQuerier"
Cohesion: 0.03
Nodes (7): actionStub, stubQuerier, CreateActionParams, UpdateActionParams, CreateOrganizationParams, NotificationConfiguration, NotificationTrigger

### Community 18 - "Ssh"
Cohesion: 0.09
Nodes (4): sshAgentStub, CreateSSHKeyParams, UpdateSSHKeyParams, Ssh

### Community 19 - "context.Context"
Cohesion: 0.05
Nodes (5): artifactStore, historyStub, variableStore, ComplexityRoot, stubQuerier

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.12
Nodes (32): replay, upstreamCapture, newTeamStub(), teamTestAPI(), TestTeamJSONAPI(), newTemplateStub(), templateTestAPI(), TestTemplateJSONAPI() (+24 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 24 - "api/main.go"
Cohesion: 0.15
Nodes (7): atomicOperation, atomicResult, createPATRequest, identityKey, dexTokenResponse, unknownCollection(), atomicWorkspaceHref()

### Community 25 - "Identity"
Cohesion: 0.16
Nodes (9): stubDexVerifier, fakeVerifier, Identity, fakeIssuer, TokenIssuer, stubDexVerifier, stringPtr(), stringSlice() (+1 more)

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

### Community 30 - "go_pkg_context"
Cohesion: 0.08
Nodes (26): S3, gitConfig, stateMetadata, failingRow, graphQLResponse, TestGitSourceDetection(), contains(), envValue() (+18 more)

### Community 31 - ".runApply"
Cohesion: 0.17
Nodes (6): Executor, preparedWorkspace, isCancelledStatus(), TestBinarySelection(), TestFirstFolder(), firstFolder()

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
Cohesion: 0.14
Nodes (7): NotificationConfigurationResolver, OrganizationResolver, ProjectConnection, QueryResolver, WorkspaceConnection, WorkspaceResolver, idsFrom()

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
Nodes (76): redactor, TestAuthMiddleware(), TestCreateOrganization(), newCollectionAPI(), TestCollectionListCarriesItsItemsAndReferences(), TestCreateCollectionDefaultsPriorityToZero(), TestCreateCollectionItemStoresTheValue(), TestCreateCollectionItemValidatesKeyValueAndCategory() (+68 more)

### Community 53 - "Pat"
Cohesion: 0.23
Nodes (4): fakeStore, Store, CreatePatParams, Pat

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "logSink"
Cohesion: 0.17
Nodes (4): captureLogPublisher, LogPublisher, logSink, baseEnvironment()

### Community 58 - ".handle"
Cohesion: 0.16
Nodes (5): isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestIsTerminalStatus()

### Community 59 - ".streamStepLog"
Cohesion: 0.26
Nodes (9): isTerminal(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData(), writeEvent() (+1 more)

### Community 60 - "WriteError"
Cohesion: 0.13
Nodes (29): decodeActionAttrs(), stringOr(), boolValue(), requireType(), decodeItemAttrs(), intValue(), globalVarCategory(), decodeNotificationConfigurationAttrs() (+21 more)

### Community 61 - "Querier"
Cohesion: 0.03
Nodes (6): AddNotificationConfigurationTemplateParams, Querier, Queries, UpdateHistoryOutputParams, UpdateOrganizationParams, Variable

### Community 62 - "state.sql.go"
Cohesion: 0.13
Nodes (31): actionAttrs, agentAttrs, collectionAttrs, globalVarAttrs, itemAttrs, jobAttrs, nestedDocument, notificationConfigurationAttrs (+23 more)

### Community 63 - "stubQuerier"
Cohesion: 0.04
Nodes (18): CreateContentParams, UpdateContentStatusParams, CreateProjectParams, UpdateJobStatusParams, CreateTagParams, GetTagByNameParams, CreateWorkspaceParams, Publisher (+10 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.08
Nodes (46): Task 3: Record the state URL and fail on history errors, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+38 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "github.com/vektah/gqlparser/v2/ast.SelectionSet"
Cohesion: 0.05
Nodes (13): JobConnection, JobEdge, NotificationConfigurationConnection, NotificationConfigurationEdge, OrganizationConnection, OrganizationEdge, PageInfo, ProjectAccessConnection (+5 more)

### Community 70 - "Team"
Cohesion: 0.14
Nodes (4): teamStub, CreateTeamParams, UpdateTeamParams, Team

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 72 - "config.go"
Cohesion: 0.09
Nodes (25): Config, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1` (+17 more)

### Community 73 - "Job"
Cohesion: 0.08
Nodes (3): cancelRunStore, reliabilityStore, mutableQuerier

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "Resolver"
Cohesion: 0.12
Nodes (8): ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 77 - "templateStub"
Cohesion: 0.13
Nodes (3): templateStub, CreateTemplateParams, UpdateTemplateParams

### Community 78 - "decodeJSON"
Cohesion: 0.20
Nodes (21): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath() (+13 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 82 - "config_versions.go"
Cohesion: 0.09
Nodes (10): TestUpstreamUIDocumentsMatchTheSchema(), approvalGateStatus(), failedJobStatus(), parseRunID(), runStatus(), decodeStateAttribute(), configurationVersionAttributes, runAttributes (+2 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "workspaces_test.go"
Cohesion: 0.17
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 88 - "Collection"
Cohesion: 0.16
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 92 - ".Query"
Cohesion: 0.07
Nodes (6): UpdateItemParams, Action, DBTX, Item, Reference, Template

### Community 94 - "Node"
Cohesion: 0.46
Nodes (3): likePattern(), builder, Compare

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "NewSubscriptions"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 102 - "newService"
Cohesion: 0.19
Nodes (15): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+7 more)

### Community 103 - "NewSubscriptions"
Cohesion: 0.15
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 108 - "api_test.go"
Cohesion: 0.20
Nodes (18): newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes(), newTestAPI() (+10 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 124 - "capture.sh"
Cohesion: 0.70
Nodes (4): pg(), request(), capture.sh script, unauthenticated()

### Community 127 - "graphql/collections.go"
Cohesion: 0.11
Nodes (20): collectionSpec, page, accessEdges(), column(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), predicate() (+12 more)

### Community 135 - "New"
Cohesion: 0.19
Nodes (13): TestArchiveHistoryWithoutStateIsNoop(), TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput(), TestGetStepLogPrefersArchivedObject(), TestGetStepLogRejectsWrongScope() (+5 more)

### Community 140 - "tarGzEntries"
Cohesion: 0.22
Nodes (13): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+5 more)

### Community 142 - ".createVCS"
Cohesion: 0.53
Nodes (4): strPtr(), strValue(), validVCSEnum(), vcsEnumOrDefault()

### Community 144 - "net/http.Handler"
Cohesion: 0.16
Nodes (17): CookieSigner, requireUUIDParams(), authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests() (+9 more)

### Community 145 - "run"
Cohesion: 0.05
Nodes (48): main(), run(), main(), run(), main(), run(), main(), run() (+40 more)

### Community 146 - "Dex"
Cohesion: 0.21
Nodes (3): Dex, Exchanger, fakeDex

### Community 147 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 150 - ".createPAT"
Cohesion: 0.36
Nodes (3): NewPAT(), writePlainJSON(), IdentityFromContext()

### Community 151 - "query.go"
Cohesion: 0.33
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 155 - "time.Time"
Cohesion: 0.26
Nodes (6): Deps, PATService, stubPATService, Verifier, timePtr(), Token

### Community 157 - ".cloneConfiguration"
Cohesion: 0.22
Nodes (7): TestSanitizeGitURL(), TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv(), runGit(), sanitizeGitURL()

### Community 158 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 165 - "auth.DexValidator"
Cohesion: 0.31
Nodes (4): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 377 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **80 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `Globalvar`, `.QueryRow`, `tagStub`, `collectionStub`, `Handler`, `TFE Handler Tests`, `CliAuthSession`, `.QueryRow`, `Vc`, `stubQuerier`, `Ssh`, `time.Time`, `Project`, `CreateAgentParams`, `.CreateNotificationTrigger`, `.SoftDeletePat`, `Pat`, `state.sql.go`, `stubQuerier`, `NewLayout`, `.QueryRow`, `Team`, `templateStub`, `Collection`, `.Query`, `.CreateNotificationConfiguration`, `jobStub`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.09753086419753086 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `executionContext` to `.___Field`, `context.Context`, `.___Directive`, `NotificationTriggerConnection`, `github.com/vektah/gqlparser/v2/ast.SelectionSet`, `generated.go`, `Resolver`, `.childFields___InputValue`, `TemplateConnection`, `context.Context`, `.___EnumValue`, `github.com/99designs/gqlgen/graphql.Marshaler`, `WorkspaceTagConnection`, `._Template`, `AccessConnection`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `API` connect `WriteResource` to `Resource`, `.createVCS`, `.createPAT`, `api/main.go`, `.updateWorkspaceRelationship`, `time.Time`, `WriteError`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Should `Globalvar` be split into smaller, more focused modules?**
  _Cohesion score 0.12631578947368421 - nodes in this community are weakly interconnected._