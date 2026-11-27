# Graph Report - terrakube-go  (2026-11-27)

## Corpus Check
- 534 files · ~284,134 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4953 nodes · 20545 edges · 263 communities (125 shown, 138 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3657 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4acb2940`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- ._Project
- io.Reader
- Core Database Schema
- Variable
- run
- net/http.ResponseWriter
- newFixture
- New
- fakeQuerier
- New
- WorkspaceRunTrigger
- steps
- github.com/aminbenmansour/terrakube-go/internal/rsql.Node
- Parse
- tfe_test.go
- sqlc Generation Config
- serveAccessToken
- Upstream Elide entities (GraphQL surface)
- dispatcherStore
- authzAPI
- payload_test.go
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- api/main.go
- Link
- getJSON
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
- addDestination
- .CreateProject
- Globalvar
- .streamStepLog
- CliAuthSession
- stubQuerier
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- models_gen.go
- newWorkspaceStub
- .QueryRow
- Job
- Queries
- .createVCS
- tagStub
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- Action
- models
- Load
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- collections.go
- New
- needsEncryption
- generated.go
- where_test.go
- streamEvents
- vcs/vcs.go
- decodeJSON
- graphql/models.go
- newJWKSServer
- Migrate
- Secret encryption at rest — design
- Cipher
- REST contract, captured from upstream
- pat.py
- Layout
- stubQuerier
- doJSON
- .EnqueueJob
- Identity
- environment
- cron.go
- NewSubscriptions
- newTestService
- context.Context
- newService
- ui-rest-inventory.sh
- ui-smoke.sh
- runTriggerAPI
- Dispatcher
- capture.sh
- 0008_action.sql
- github.com/vektah/gqlparser/v2/ast.SelectionSet
- compare.sh
- start.sh
- ._NotificationConfiguration
- newTestExecutor
- time.Time
- WebhookEvent
- CreateTemplateParams
- rsql.Node
- github.com/aminbenmansour/terrakube-go/internal/store.Workspace
- clilogin/handler_test.go
- newTestAPI
- Secret key rotation — design
- OrganizationConnection
- ui.go
- API
- UpdateNotificationConfigurationParams
- NewSubscriptions
- fakeStore
- stubQuerier
- Executor
- provider
- History
- tarGzEntries
- webhook_test.go
- github.com/aminbenmansour/terrakube-go/internal/store.Step
- net/http.Handler
- notificationSendAPI
- AccessConnection
- github_test.go
- Metrics
- transientError
- log/slog.Logger
- Threat model
- WorkspaceTagConnection
- auth.DexValidator
- store.Workspace
- TestModelColumnOrder
- RSQL filters in SQL — design
- Upstream UI GraphQL operations
- Deps
- Workspace
- dispatcherWith
- ParsePage
- Repository
- Review Focus
- net/http.Request
- Team
- .build
- .Query
- .Encrypt
- workspaceQuery
- .SoftDeletePat
- NotificationConfigurationEdge
- .CreateArchive
- .SetNotificationSecret

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 251 edges
3. `Queries` - 247 edges
4. `API` - 225 edges
5. `New()` - 175 edges
6. `stubQuerier` - 171 edges
7. `WriteError()` - 154 edges
8. `param()` - 90 edges
9. `ParseFieldsets()` - 78 edges
10. `WriteResource()` - 69 edges

## Surprising Connections (you probably didn't know these)
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Command` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Report and verification` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Secret Encryption at Rest Implementation Plan` --references--> `Querier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/store/querier.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (263 total, 138 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (27): Write(), WriteResource(), configurationVersionResource(), projectResource(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID() (+19 more)

### Community 1 - "rbac.go"
Cohesion: 0.19
Nodes (27): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, adminOnly(), CanApproveJob() (+19 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (12): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+4 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 7 - "run"
Cohesion: 0.04
Nodes (84): main(), run(), main(), run(), main(), run(), main(), run() (+76 more)

### Community 8 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (20): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), WithClaims(), runnerClaims(), claimsFromContext() (+12 more)

### Community 9 - "newFixture"
Cohesion: 0.13
Nodes (24): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+16 more)

### Community 10 - "New"
Cohesion: 0.09
Nodes (36): globToRegex(), matchBranch(), matchPaths(), matchTrigger(), pathRegex(), branchFromRef(), deliveryID(), nonEmptyPtr() (+28 more)

### Community 12 - "New"
Cohesion: 0.10
Nodes (22): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+14 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 17 - "github.com/aminbenmansour/terrakube-go/internal/rsql.Node"
Cohesion: 0.05
Nodes (6): listQuery, stubLists, Task 3: Handler cutover, scopeIDs(), GetTagByNameParams, Tag

### Community 18 - "Parse"
Cohesion: 0.19
Nodes (14): Parse(), dump(), TestParse(), TestParseRejects(), ptr(), seedWorkspaces(), TestListPageBetween(), TestListPageFilters() (+6 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (43): Backend, assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL() (+35 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "serveAccessToken"
Cohesion: 0.12
Nodes (24): expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups(), TestAccessTokenRejectsBadInput() (+16 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "dispatcherStore"
Cohesion: 0.13
Nodes (6): demoOrganizations(), dispatcherStore, FailStuckNotificationOutboxParams, Organization, ReclaimStuckNotificationOutboxParams, RecordNotificationOutboxParams

### Community 24 - "authzAPI"
Cohesion: 0.10
Nodes (22): discoveryStub, fakeHooks, missingOrganizationStub, authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound() (+14 more)

### Community 25 - "payload_test.go"
Cohesion: 0.08
Nodes (49): channelLabel(), networkCause(), NewDeliverer(), parseRetryAfter(), deliveryError(), stringPtr(), TestDeliverDoesNotFollowRedirects(), testDeliverer() (+41 more)

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
Nodes (59): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, discoveryGroup, discoveryRepositoryPage (+51 more)

### Community 31 - "Link"
Cohesion: 0.10
Nodes (24): parseInclude(), collectionResource(), itemResource(), referenceResource(), jobResource(), stepResource(), timestamp(), notificationConfigurationResource() (+16 more)

### Community 32 - "getJSON"
Cohesion: 0.28
Nodes (4): getJSON(), github, Group, RepositoryPage

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
Cohesion: 0.11
Nodes (9): notificationConfigurationResolver, organizationResolver, queryResolver, ResolverRoot, workspaceResolver, organizationScope(), TestOrganizationScope(), idsFrom() (+1 more)

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

### Community 44 - "addDestination"
Cohesion: 0.37
Nodes (15): strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased(), TestTickIgnoresPlanOnlyRun() (+7 more)

### Community 45 - ".CreateProject"
Cohesion: 0.20
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 47 - ".streamStepLog"
Cohesion: 0.15
Nodes (15): compareStreamIDs(), isTerminal(), outputOf(), parseResumeID(), splitStreamID(), stepOf(), streamKey(), TestCompareStreamIDs() (+7 more)

### Community 48 - "CliAuthSession"
Cohesion: 0.07
Nodes (15): fakeStore, Service, Store, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value() (+7 more)

### Community 49 - "stubQuerier"
Cohesion: 0.14
Nodes (4): CreateContentParams, UpdateContentStatusParams, Content, stubQuerier

### Community 51 - "time.Duration"
Cohesion: 0.07
Nodes (13): waitFor(), captureLogPublisher, fakeQueue, JobQueue, LogPublisher, envDuration(), defaultRetryDelay(), TestCleanupStaleWorkspaces() (+5 more)

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (100): TestAccessTokenWithoutStore(), TestAuthMiddleware(), TestCreateOrganization(), newCollectionAPI(), TestCollectionListCarriesItsItemsAndReferences(), TestCreateCollectionDefaultsPriorityToZero(), TestCreateCollectionItemStoresTheValue(), TestCreateCollectionItemValidatesKeyValueAndCategory() (+92 more)

### Community 53 - "Pat"
Cohesion: 0.21
Nodes (5): newFakeStore(), fakeStore, Store, CreatePatParams, Pat

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "models_gen.go"
Cohesion: 0.06
Nodes (11): ComplexityRoot, JobConnection, NotificationConfigurationConnection, NotificationTriggerConnection, NotificationTriggerEdge, PageInfo, ProjectConnection, Query (+3 more)

### Community 57 - "newWorkspaceStub"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 59 - ".QueryRow"
Cohesion: 0.03
Nodes (5): Agent, CreateWebhookParams, NotificationOutboxExistsForWorkspaceParams, UpdateOrganizationParams, UpdateWebhookParams

### Community 60 - "Job"
Cohesion: 0.10
Nodes (4): ApproveJobParams, GetJobByWorkspaceAndTriggerParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams

### Community 61 - "Queries"
Cohesion: 0.02
Nodes (12): TestWebhookQueries(), AddNotificationConfigurationTemplateParams, CreateNotificationTriggerParams, DeleteWebhookDeliveryParams, FailStuckRunTriggerEventsParams, Querier, Queries, ReclaimStuckRunTriggerEventsParams (+4 more)

### Community 62 - ".createVCS"
Cohesion: 0.50
Nodes (3): discoveryConnection(), discoveryPage(), writeRawJSON()

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.06
Nodes (49): cancelRunStore, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+41 more)

### Community 67 - "Action"
Cohesion: 0.13
Nodes (4): actionStub, CreateActionParams, UpdateActionParams, Action

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "Load"
Cohesion: 0.17
Nodes (20): Config, Task 1: Keyring configuration, env(), envBool(), envFirst(), envInt64(), envList(), Load() (+12 more)

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - "New"
Cohesion: 0.21
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 73 - "needsEncryption"
Cohesion: 0.13
Nodes (12): Task 5: Backfill implementation, backfillGlobalVariables(), backfillNotifications(), backfillSSHKeys(), backfillVariables(), backfillVCS(), encryptIfNeeded(), needsEncryption() (+4 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.19
Nodes (15): column(), predicate(), Field, likePattern(), mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode() (+7 more)

### Community 76 - "streamEvents"
Cohesion: 0.14
Nodes (8): projectResolver, Resolver, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 77 - "vcs/vcs.go"
Cohesion: 0.30
Nodes (12): TestGitLabConnectionWithoutAToken(), TestGitLabDiscovery(), baseURL(), Branches(), Groups(), providerFor(), Repositories(), Test() (+4 more)

### Community 78 - "decodeJSON"
Cohesion: 0.12
Nodes (31): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent() (+23 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "Migrate"
Cohesion: 0.22
Nodes (9): Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 7: Documentation and graph, Task 8: Verification (+1 more)

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - "Cipher"
Cohesion: 0.29
Nodes (14): formatOldKeys(), logReport(), envelopeKey(), rotateGlobalVariables(), rotateNotifications(), rotatePtr(), rotateSSHKeys(), rotateValue() (+6 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "Layout"
Cohesion: 0.12
Nodes (8): Validator, Deps, fakeLocker, JobLocker, Layout, Deps, JobQueue, Deps

### Community 89 - "stubQuerier"
Cohesion: 0.16
Nodes (20): replay, upstreamCapture, assertSameDocumentShape(), assertSameResourceShape(), captureActionStore(), captureCollectionStore(), captureGlobalVarStore(), captureJobStore() (+12 more)

### Community 91 - "doJSON"
Cohesion: 0.32
Nodes (7): doJSON(), failed(), providerError(), providerMessage(), githubOwnerAndRepo(), gitlab, WebhookSpec

### Community 92 - ".EnqueueJob"
Cohesion: 0.21
Nodes (5): fakeQueue, JobQueue, fakeQueue, stubJobQueue, stubQueue

### Community 94 - "Identity"
Cohesion: 0.09
Nodes (17): groupVerifier, PATService, stubDexVerifier, stubIdentityVerifier, stubPATService, fakeVerifier, Identity, fakeDex (+9 more)

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 99 - "cron.go"
Cohesion: 0.24
Nodes (8): at(), Parse(), parseField(), splitStep(), mustParse(), TestParseMatchesQuartzFields(), TestParseRejectsWhatQuartzRejects(), Cron

### Community 100 - "NewSubscriptions"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 103 - "context.Context"
Cohesion: 0.03
Nodes (12): collectionStub, globalVarStub, sshAgentStub, teamStub, templateStub, vcsStub, workspaceStub, artifactStore (+4 more)

### Community 108 - "newService"
Cohesion: 0.21
Nodes (14): decodeSecret(), New(), dexUser(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp(), TestDeleteEnforcesOwnership() (+6 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.08
Nodes (26): runTriggerStub, scheduleStub, newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), runTriggerRBACAPI(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit() (+18 more)

### Community 121 - "Dispatcher"
Cohesion: 0.27
Nodes (4): NewDispatcher(), strPtr(), Dispatcher, NotificationOutbox

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 135 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 137 - "time.Time"
Cohesion: 0.21
Nodes (7): bucket, limiter, newLimiter(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), Dispatcher, memObject

### Community 140 - "WebhookEvent"
Cohesion: 0.19
Nodes (3): CreateWebhookEventParams, UpdateWebhookEventParams, WebhookEvent

### Community 142 - "rsql.Node"
Cohesion: 0.13
Nodes (9): exists(), likePattern(), relate(), And, Or, parser, Attr, Kind (+1 more)

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.07
Nodes (44): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+36 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 149 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 150 - "API"
Cohesion: 0.07
Nodes (36): actionAttrs, itemAttrs, teamAttrs, templateAttrs, actionResource(), decodeActionAttrs(), stringOr(), boolValue() (+28 more)

### Community 154 - "NewSubscriptions"
Cohesion: 0.13
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 157 - "fakeStore"
Cohesion: 0.05
Nodes (16): captureJobs(), captureSteps(), stubQuerier, fakeStore, fakeStore, CreateJobParams, CreateStepParams, GetStepByJobAndNumberParams (+8 more)

### Community 159 - "stubQuerier"
Cohesion: 0.02
Nodes (17): stubQuerier, CreateCollectionParams, CreateItemParams, CreateReferenceParams, UpdateCollectionParams, UpdateItemParams, Collection, Item (+9 more)

### Community 161 - "Executor"
Cohesion: 0.05
Nodes (29): Executor, gitConfig, logSink, preparedWorkspace, redactor, TestGitSourceDetection(), TestSanitizeGitURL(), isCancelledStatus() (+21 more)

### Community 162 - "provider"
Cohesion: 0.22
Nodes (7): Measured route surface, Order that unblocks the most UI, UI REST operations, Where we stand, Branch, BranchPage, provider

### Community 163 - "History"
Cohesion: 0.09
Nodes (6): reliabilityStore, CreateHistoryParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams, History, UpdateStepStatusParams

### Community 164 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 165 - "webhook_test.go"
Cohesion: 0.36
Nodes (10): DeleteWebhook(), SyncWebhook(), newHookServer(), TestDeleteWebhookTreats404AsRemoved(), TestGitHubWebhookCreate(), TestGitHubWebhookUpdateAndDelete(), TestGitLabWebhookCreate(), TestGitLabWebhookDelete() (+2 more)

### Community 167 - "net/http.Handler"
Cohesion: 0.06
Nodes (33): authorization, nestedDocument, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), parseList(), requireUUIDParams(), authorizationFrom() (+25 more)

### Community 168 - "notificationSendAPI"
Cohesion: 0.14
Nodes (23): newDeliveryStub(), TestNotificationDeliveryHistory(), TestNotificationDeliveryHistoryRejectsABadLimit(), TestNotificationDeliveryHistoryUnknownWorkspace(), TestNotificationDeliveryRetry(), TestNotificationDeliveryRetryRejections(), captureJSONServer(), notificationSendAPI() (+15 more)

### Community 171 - "github_test.go"
Cohesion: 0.38
Nodes (9): githubRoutes(), TestBranchesGitHub(), TestGitHubConnectionsWithoutAToken(), TestGroupsGitHub(), TestRepositoriesGitHub(), TestTestConnection(), newFakeProvider(), TestRejectedCredentialIsAUseful4xx() (+1 more)

### Community 172 - "Metrics"
Cohesion: 0.13
Nodes (4): metricsSource, fakeStats, Metrics, StatsSource

### Community 175 - "log/slog.Logger"
Cohesion: 0.20
Nodes (12): startScheduler(), NewDispatcher(), NewDispatcher(), strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued(), TestTickSkipsLockedWorkspaceAndUnparseableCron() (+4 more)

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 184 - "auth.DexValidator"
Cohesion: 0.15
Nodes (6): Dex, Exchanger, contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 190 - "store.Workspace"
Cohesion: 0.09
Nodes (18): API, atomicOperation, Deps, webhookEventAttrs, webhookSyncRef, strValue(), atomicWorkspaceHref(), isAtomicRequest() (+10 more)

### Community 191 - "TestModelColumnOrder"
Cohesion: 0.25
Nodes (4): failingDB, collectStrings(), snakeCase(), TestModelColumnOrder()

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 197 - "Upstream UI GraphQL operations"
Cohesion: 0.25
Nodes (6): Connection shape, Entity and relationship surface, Filter expressions the UI generates, Other contracts observed, Sort strings, Upstream UI GraphQL operations

### Community 198 - "Deps"
Cohesion: 0.33
Nodes (3): Deps, fakePinger, Pinger

### Community 202 - "dispatcherWith"
Cohesion: 0.36
Nodes (8): dispatcherWith(), newDispatcherFixture(), TestDispatcherFailsAfterMaxAttempts(), TestDispatcherFailsPermanentRejection(), TestDispatcherIncludesTheFailureSummary(), TestDispatcherRetriesTransientFailures(), TestDispatcherSendsAndRecordsSuccess(), TestDispatcherSkipsSupersededJobStatus()

### Community 204 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 229 - "Repository"
Cohesion: 0.40
Nodes (4): optional(), githubRepositoryJSON, gitlabProjectJSON, Repository

### Community 231 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 232 - "net/http.Request"
Cohesion: 0.09
Nodes (14): notificationConfigurationAttrs, vcsAttrs, strPtr(), decodeNotificationConfigurationAttrs(), writeInUseConflict(), decodeVCSAttrs(), requiredVCSString(), validVCSEnum() (+6 more)

### Community 234 - "Team"
Cohesion: 0.04
Nodes (23): AccessTokenStore, fakeAccessTokenStore, contains(), captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces(), captureProjectAccessRow() (+15 more)

### Community 235 - ".build"
Cohesion: 0.10
Nodes (20): collectionSpec, DB, listArgs, page, accessEdges(), afterString(), jobEdges(), loadPage() (+12 more)

### Community 236 - ".Query"
Cohesion: 0.06
Nodes (17): CountJobsByStatusRow, DBTX, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow (+9 more)

### Community 248 - ".Encrypt"
Cohesion: 0.15
Nodes (9): Task 2: Cipher and envelope, Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 3: Rotation scan and re-encryption, Task 5: Operator documentation (+1 more)

### Community 251 - "workspaceQuery"
Cohesion: 0.40
Nodes (3): matchesTagFilters(), tagFilter, workspaceQuery

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 549 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **138 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Queries` to `.CreateProjectAccess`, `Variable`, `run`, `net/http.ResponseWriter`, `time.Time`, `New`, `fakeQuerier`, `New`, `WorkspaceRunTrigger`, `CreateTemplateParams`, `.marshalOBoolean2bool`, `WebhookEvent`, `newTestAPI`, `github.com/aminbenmansour/terrakube-go/internal/rsql.Node`, `tfe_test.go`, `dispatcherStore`, `UpdateNotificationConfigurationParams`, `fakeStore`, `api/main.go`, `stubQuerier`, `History`, `github.com/aminbenmansour/terrakube-go/internal/store.Step`, `notificationSendAPI`, `Metrics`, `.CreateProject`, `Globalvar`, `log/slog.Logger`, `CliAuthSession`, `stubQuerier`, `Pat`, `.QueryRow`, `Job`, `store.Workspace`, `tagStub`, `Action`, `Workspace`, `needsEncryption`, `Migrate`, `Layout`, `context.Context`, `Team`, `.Query`, `runTriggerAPI`, `Dispatcher`, `.SoftDeletePat`, `.CreateArchive`, `.SetNotificationSecret`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.08542319749216301 - nodes in this community are weakly interconnected._
- **Why does `API` connect `store.Workspace` to `net/http.Handler`, `net/http.Request`, `API`, `api/main.go`, `.createVCS`, `Link`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `openRotationTest()` connect `run` to `context.Context`, `newTestAPI`, `testing.T`, `Queries`, `api/main.go`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Should `io.Reader` be split into smaller, more focused modules?**
  _Cohesion score 0.07164404223227752 - nodes in this community are weakly interconnected._