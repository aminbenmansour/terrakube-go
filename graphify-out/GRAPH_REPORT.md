# Graph Report - terrakube-go  (2026-11-28)

## Corpus Check
- 645 files · ~316,250 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 5484 nodes · 23259 edges · 270 communities (137 shown, 133 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 4207 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a610ecb9`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- WriteError
- Agent Workflows and Compatibility Docs
- ._Project
- io.Reader
- Core Database Schema
- registryStub
- run
- Handler
- newFixture
- New
- fakeQuerier
- policyStub
- Executor
- steps
- models_gen.go
- github.com/aminbenmansour/terrakube-go/internal/rsql.Node
- rbac.go
- tfe_test.go
- sqlc Generation Config
- github.com/99designs/gqlgen/graphql.Marshaler
- Upstream Elide entities (GraphQL surface)
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- Webhook
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
- fakeStore
- History
- ResolveOrganization
- store.Step
- CliAuthSession
- Content
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- federatedStub
- newWorkspaceStub
- executionContext
- Queries
- ._NotificationConfiguration
- tagStub
- dex-config.yaml
- NewLayout
- models
- Load
- stubQuerier
- collections.go
- New
- dispatcherStore
- generated.go
- where_test.go
- streamEvents
- vcs/vcs.go
- decodeJSON
- graphql/models.go
- auth/dex_test.go
- Deps
- Secret encryption at rest — design
- Cipher
- REST contract, captured from upstream
- pat.py
- stubQuerier
- doJSON
- .CreateProject
- net/http.ResponseWriter
- environment
- cron.go
- NewSubscriptions
- newTestService
- context.Context
- Claims
- ui-rest-inventory.sh
- ui-smoke.sh
- runTriggerAPI
- .EnqueueJob
- capture.sh
- 0008_action.sql
- JobConnection
- Service
- compare.sh
- start.sh
- WorkspaceTagConnection
- log/slog.Logger
- time.Time
- strPtr
- Identity
- Parse
- store.Action
- clilogin/handler_test.go
- newTestAPI
- ui.go
- newTestExecutor
- addDestination
- net/http.Request
- WebhookEvent
- NewSubscriptions
- nestedStub
- stubQuerier
- Service
- runTriggerAPI
- tarGzEntries
- webhook_test.go
- RateLimit
- net/http.Handler
- notificationSendAPI
- listQuery
- Review Focus
- Metrics
- transientError
- .CreateOrganization
- Threat model
- .cloneConfiguration
- auth.DexValidator
- TestModelColumnOrder
- RSQL filters in SQL — design
- Upstream UI GraphQL operations
- Workspace
- ParsePage
- Repository
- TeamToken
- Rotate
- Team
- graphql/collections.go
- Querier
- RelationshipID
- Layout
- .QueryRow
- newService
- grantsAPI
- New
- discoveryStub
- .server
- Attr
- authzAPI
- .CreateArchive
- TestWebhookQueries
- logSink

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 321 edges
3. `Queries` - 312 edges
4. `API` - 307 edges
5. `New()` - 239 edges
6. `WriteError()` - 192 edges
7. `stubQuerier` - 171 edges
8. `ParseFieldsets()` - 110 edges
9. `param()` - 109 edges
10. `WriteResource()` - 86 edges

## Surprising Connections (you probably didn't know these)
- `Where we stand` --references--> `provider`  [INFERRED]
  docs/upstream/ui-rest-operations.md → internal/vcs/vcs.go
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Secret Encryption at Rest Implementation Plan` --references--> `Querier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/store/querier.go
- `Task 3: Transparent store decorator` --references--> `Querier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/store/querier.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (270 total, 133 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (26): WriteResource(), configurationVersionResource(), projectResource(), parseWorkspaceQuery(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID() (+18 more)

### Community 1 - "WriteError"
Cohesion: 0.06
Nodes (50): actionAttrs, itemAttrs, moduleAttrs, notificationConfigurationAttrs, policySetParameterAttrs, providerAttrs, providerVersionAttrs, teamAttrs (+42 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (12): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+4 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "registryStub"
Cohesion: 0.04
Nodes (17): registryStub, filterVersions(), CreateImplementationParams, CreateModuleParams, CreateModuleVersionParams, CreateProviderParams, CreateProviderVersionParams, Implementation (+9 more)

### Community 7 - "run"
Cohesion: 0.06
Nodes (67): main(), run(), main(), run(), main(), run(), main(), run() (+59 more)

### Community 8 - "Handler"
Cohesion: 0.21
Nodes (8): isExternalUser(), matchesMD5(), New(), newUUID(), pathIDs(), sharedIDsContain(), Deps, Handler

### Community 9 - "newFixture"
Cohesion: 0.10
Nodes (27): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+19 more)

### Community 10 - "New"
Cohesion: 0.10
Nodes (25): globToRegex(), matchBranch(), matchPaths(), matchTrigger(), pathRegex(), branchFromRef(), deliveryID(), nonEmptyPtr() (+17 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.07
Nodes (3): encryptedStore, fakeQuerier, Variable

### Community 12 - "policyStub"
Cohesion: 0.04
Nodes (6): policySetChildren, policyStub, PolicyAttachment, PolicyExemption, PolicySet, PolicySetParameter

### Community 13 - "Executor"
Cohesion: 0.09
Nodes (12): Executor, preparedWorkspace, isCancelledStatus(), isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestBinarySelection() (+4 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 16 - "models_gen.go"
Cohesion: 0.06
Nodes (10): AccessConnection, AccessEdge, ComplexityRoot, PageInfo, ProjectConnection, ProjectEdge, Query, Subscription (+2 more)

### Community 17 - "github.com/aminbenmansour/terrakube-go/internal/rsql.Node"
Cohesion: 0.06
Nodes (6): listQuery, stubLists, Task 3: Handler cutover, Page, ListSpec, Tag

### Community 18 - "rbac.go"
Cohesion: 0.26
Nodes (22): Task 1: Role matrix and permission set, adminOnly(), CanApproveJob(), CanDestroyJob(), CanManageCollection(), CanManageJob(), CanManageModule(), CanManagePolicies() (+14 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.07
Nodes (46): Global Constraints, M4 run logs and artifacts Implementation Plan, Self-review notes, Backend, assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes() (+38 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.05
Nodes (3): NotificationConfigurationEdge, TemplateEdge, WorkspaceEdge

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.27
Nodes (4): NewDispatcher(), strPtr(), Dispatcher, NotificationOutbox

### Community 24 - "Webhook"
Cohesion: 0.09
Nodes (14): atomicOperation, vcsAttrs, webhookEventAttrs, webhookSyncRef, strPtr(), strValue(), atomicWorkspaceHref(), isAtomicRequest() (+6 more)

### Community 25 - "payload_test.go"
Cohesion: 0.07
Nodes (51): channelLabel(), networkCause(), NewDeliverer(), parseRetryAfter(), deliveryError(), stringPtr(), TestDeliverDoesNotFollowRedirects(), testDeliverer() (+43 more)

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
Nodes (65): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, createTeamTokenRequest, discoveryGroup (+57 more)

### Community 31 - "Link"
Cohesion: 0.06
Nodes (47): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), federatedClaimResource(), federatedResource(), globalVarResource() (+39 more)

### Community 32 - "getJSON"
Cohesion: 0.20
Nodes (7): getJSON(), Branch, BranchPage, github, Group, provider, RepositoryPage

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

### Community 44 - "fakeStore"
Cohesion: 0.06
Nodes (9): captureJobs(), Publisher, publishingStore, stubQuerier, fakeStore, CreateJobParams, Job, Template (+1 more)

### Community 45 - "History"
Cohesion: 0.16
Nodes (3): CreateHistoryParams, GetHistoryBeforeParams, History

### Community 46 - "ResolveOrganization"
Cohesion: 0.10
Nodes (12): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), contains(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access (+4 more)

### Community 47 - "store.Step"
Cohesion: 0.07
Nodes (8): cancelRunStore, reliabilityStore, captureSteps(), CreateStepParams, GetHistoryByJobReferenceParams, GetStepByJobAndNumberParams, Step, UpdateStepStatusParams

### Community 48 - "CliAuthSession"
Cohesion: 0.10
Nodes (5): fakeStore, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams, SetCliAuthSessionIdentityParams

### Community 49 - "Content"
Cohesion: 0.29
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 51 - "time.Duration"
Cohesion: 0.09
Nodes (9): fakeQueue, JobQueue, defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), backoff(), backoff(), TestBackoffGrowsAndCaps() (+1 more)

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (135): expectPermissions(), serveAccessToken(), serveTeamTokens(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups() (+127 more)

### Community 53 - "Pat"
Cohesion: 0.17
Nodes (5): fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "federatedStub"
Cohesion: 0.10
Nodes (7): federatedStub, CreateFederatedClaimParams, CreateFederatedParams, Federated, FederatedClaim, UpdateFederatedClaimParams, UpdateFederatedParams

### Community 57 - "newWorkspaceStub"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 59 - "executionContext"
Cohesion: 0.06
Nodes (4): executionContext, NotificationConfigurationConnection, OrganizationConnection, OrganizationEdge

### Community 61 - "Queries"
Cohesion: 0.02
Nodes (10): replay, FailStuckNotificationOutboxParams, FailStuckRunTriggerEventsParams, Querier, ReclaimStuckNotificationOutboxParams, ReclaimStuckRunTriggerEventsParams, RecordWebhookDeliveryParams, UpdateHistoryOutputParams (+2 more)

### Community 63 - "tagStub"
Cohesion: 0.09
Nodes (4): tagStub, matchesTagFilters(), tagFilter, workspaceQuery

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.05
Nodes (66): redactor, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+58 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "Load"
Cohesion: 0.12
Nodes (28): Config, S3, Storage, Review Focus, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph (+20 more)

### Community 70 - "stubQuerier"
Cohesion: 0.06
Nodes (5): jobStub, statusQuerier, stubQuerier, workspaceWithTags, stubQuerier

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - "New"
Cohesion: 0.20
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 73 - "dispatcherStore"
Cohesion: 0.07
Nodes (15): missingOrganizationStub, scopeIDs(), demoOrganizations(), dispatcherWith(), newDispatcherFixture(), TestDispatcherFailsAfterMaxAttempts(), TestDispatcherFailsPermanentRejection(), TestDispatcherIncludesTheFailureSummary() (+7 more)

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
Cohesion: 0.19
Nodes (22): githubRoutes(), TestBranchesGitHub(), TestGitHubConnectionsWithoutAToken(), TestGroupsGitHub(), TestRepositoriesGitHub(), TestTestConnection(), TestGitLabConnectionWithoutAToken(), TestGitLabDiscovery() (+14 more)

### Community 78 - "decodeJSON"
Cohesion: 0.10
Nodes (37): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent() (+29 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "auth/dex_test.go"
Cohesion: 0.17
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "Deps"
Cohesion: 0.15
Nodes (8): Deps, PATService, stubPATService, stubTeamTokens, TeamTokenService, timePtr(), Token, Token

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - "Cipher"
Cohesion: 0.07
Nodes (39): Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation (+31 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 89 - "stubQuerier"
Cohesion: 0.13
Nodes (26): upstreamCapture, federatedCaptureStore(), capturePathLastID(), captureRelationshipID(), policyCaptureStore(), captureID(), registryCaptureStore(), registryProviderStore() (+18 more)

### Community 91 - "doJSON"
Cohesion: 0.24
Nodes (9): TestValidateEndpointRefusesInternalAddresses(), validateEndpoint(), doJSON(), failed(), providerError(), githubOwnerAndRepo(), Error, gitlab (+1 more)

### Community 92 - ".CreateProject"
Cohesion: 0.20
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 94 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (23): federatedAttrs, federatedClaimAttrs, implementationAttrs, policySetAttrs, decodeFederatedAttrs(), decodeFederatedClaimAttrs(), federatedNotBlank(), federatedRequired() (+15 more)

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
Nodes (20): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+12 more)

### Community 103 - "context.Context"
Cohesion: 0.02
Nodes (17): collectionStub, globalVarStub, sshAgentStub, teamStub, templateStub, vcsStub, workspaceStub, artifactStore (+9 more)

### Community 108 - "Claims"
Cohesion: 0.07
Nodes (31): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken() (+23 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.15
Nodes (13): scheduleStub, newScheduleAPI(), newScheduleStub(), scheduleBody(), scheduleStr(), TestCreateScheduleCopiesAttributesAndDefaultsEnabled(), TestCreateScheduleRejectsForeignTemplate(), TestDeleteScheduleAndCrossWorkspaceLookup() (+5 more)

### Community 121 - ".EnqueueJob"
Cohesion: 0.22
Nodes (6): fakeQueue, JobQueue, fakeQueue, JobQueue, stubJobQueue, stubQueue

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 128 - "Service"
Cohesion: 0.13
Nodes (12): Service, Store, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value() (+4 more)

### Community 135 - "log/slog.Logger"
Cohesion: 0.21
Nodes (12): startScheduler(), waitFor(), NewDispatcher(), NewDispatcher(), strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued() (+4 more)

### Community 137 - "time.Time"
Cohesion: 0.24
Nodes (5): policyExemptionAttrs, bucket, limiter, Dispatcher, memObject

### Community 140 - "strPtr"
Cohesion: 0.19
Nodes (13): sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), newTeamStub(), TestTeamJSONAPI(), newTemplateStub(), TestTemplateJSONAPI() (+5 more)

### Community 141 - "Identity"
Cohesion: 0.11
Nodes (13): groupVerifier, stubDexVerifier, stubIdentityVerifier, fakeVerifier, Identity, fakeIssuer, groupVerifier, stubDexVerifier (+5 more)

### Community 142 - "Parse"
Cohesion: 0.12
Nodes (16): Parse(), dump(), TestParse(), TestParseRejects(), ptr(), seedWorkspaces(), TestListPageBetween(), TestListPageFilters() (+8 more)

### Community 143 - "store.Action"
Cohesion: 0.17
Nodes (4): actionStub, CreateActionParams, UpdateActionParams, Action

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.06
Nodes (50): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+42 more)

### Community 146 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 147 - "newTestExecutor"
Cohesion: 0.31
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 149 - "addDestination"
Cohesion: 0.18
Nodes (18): strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased(), TestTickIgnoresPlanOnlyRun() (+10 more)

### Community 150 - "net/http.Request"
Cohesion: 0.06
Nodes (14): API, authorization, New(), param(), parseList(), unknownCollection(), authorizationFrom(), permissionDeniedDetail() (+6 more)

### Community 154 - "NewSubscriptions"
Cohesion: 0.14
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 159 - "stubQuerier"
Cohesion: 0.02
Nodes (14): stubQuerier, CreateCollectionParams, Agent, Collection, Globalvar, Item, NotificationTrigger, Project (+6 more)

### Community 161 - "Service"
Cohesion: 0.39
Nodes (3): stringPtr(), timePtr(), Service

### Community 162 - "runTriggerAPI"
Cohesion: 0.24
Nodes (13): runTriggerStub, newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), runTriggerRBACAPI(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit(), TestCreateRunTriggerRejectsCycle() (+5 more)

### Community 164 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 165 - "webhook_test.go"
Cohesion: 0.36
Nodes (10): DeleteWebhook(), SyncWebhook(), newHookServer(), TestDeleteWebhookTreats404AsRemoved(), TestGitHubWebhookCreate(), TestGitHubWebhookUpdateAndDelete(), TestGitLabWebhookCreate(), TestGitLabWebhookDelete() (+2 more)

### Community 166 - "RateLimit"
Cohesion: 0.29
Nodes (7): key(), newLimiter(), RateLimit(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), TestRateLimitDisabledByDefault(), TestRateLimitMiddleware()

### Community 167 - "net/http.Handler"
Cohesion: 0.08
Nodes (18): Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), isWriteMethod(), NewNotification(), NewPAT(), writePlainJSON() (+10 more)

### Community 168 - "notificationSendAPI"
Cohesion: 0.14
Nodes (23): newDeliveryStub(), TestNotificationDeliveryHistory(), TestNotificationDeliveryHistoryRejectsABadLimit(), TestNotificationDeliveryHistoryUnknownWorkspace(), TestNotificationDeliveryRetry(), TestNotificationDeliveryRetryRejections(), captureJSONServer(), notificationSendAPI() (+15 more)

### Community 169 - "listQuery"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 171 - "Review Focus"
Cohesion: 0.33
Nodes (5): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract

### Community 172 - "Metrics"
Cohesion: 0.10
Nodes (7): metricsSource, Deps, fakePinger, fakeStats, Metrics, Pinger, StatsSource

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - ".cloneConfiguration"
Cohesion: 0.22
Nodes (9): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL(), runGit() (+1 more)

### Community 184 - "auth.DexValidator"
Cohesion: 0.13
Nodes (7): Dex, Exchanger, fakeDex, contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 191 - "TestModelColumnOrder"
Cohesion: 0.25
Nodes (4): failingDB, collectStrings(), snakeCase(), TestModelColumnOrder()

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 197 - "Upstream UI GraphQL operations"
Cohesion: 0.17
Nodes (10): Connection shape, Entity and relationship surface, Filter expressions the UI generates, Other contracts observed, Sort strings, Upstream UI GraphQL operations, Measured route surface, Order that unblocks the most UI (+2 more)

### Community 204 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 229 - "Repository"
Cohesion: 0.40
Nodes (4): optional(), githubRepositoryJSON, gitlabProjectJSON, Repository

### Community 231 - "TeamToken"
Cohesion: 0.17
Nodes (5): CreateTeamTokenParams, SoftDeleteTeamTokenParams, TeamToken, fakeStore, Store

### Community 232 - "Rotate"
Cohesion: 0.08
Nodes (28): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs, Approaches considered (+20 more)

### Community 234 - "Team"
Cohesion: 0.07
Nodes (11): AccessTokenStore, fakeAccessTokenStore, contains(), captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces(), captureProjectAccessRow() (+3 more)

### Community 235 - "graphql/collections.go"
Cohesion: 0.10
Nodes (20): collectionSpec, DB, listArgs, page, accessEdges(), afterString(), jobEdges(), loadPage() (+12 more)

### Community 236 - "Querier"
Cohesion: 0.03
Nodes (17): createRawVariable(), createRawVariableWithID(), CountJobsByStatusRow, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow (+9 more)

### Community 248 - "RelationshipID"
Cohesion: 0.20
Nodes (10): nestedDocument, RelationshipIdentifiers(), RelationshipList(), TestRelationshipIdentifiers(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship (+2 more)

### Community 251 - "Layout"
Cohesion: 0.12
Nodes (8): Validator, Deps, fakeLocker, JobLocker, Layout, Deps, JobQueue, Deps

### Community 255 - ".QueryRow"
Cohesion: 0.02
Nodes (6): UpdateCollectionParams, UpdateItemParams, ApproveJobParams, GetTagByNameParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams

### Community 259 - "newService"
Cohesion: 0.23
Nodes (11): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRequiresMembershipAndAnInteractiveUser(), TestCreateSignsTheUpstreamToken(), TestListAndDeleteAreScopedToTheCallersTeams() (+3 more)

### Community 261 - "grantsAPI"
Cohesion: 0.31
Nodes (10): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+2 more)

### Community 271 - "New"
Cohesion: 0.49
Nodes (11): New(), githubSignature(), newStubStore(), post(), TestBranchAndPathMustMatchTheTrigger(), TestGitHubBadOrMissingSignatureIsRejected(), TestGitHubPushStartsRun(), TestGitLabPushRequiresTheToken() (+3 more)

### Community 273 - ".server"
Cohesion: 0.33
Nodes (7): fakeHooks, TestAtomicWebhookInvalidTriggerChangesNothing(), TestAtomicWebhookSaveRegistersFinalTriggers(), TestCreateWebhookEventRegistersProvider(), TestDeleteWebhookRemovesProviderHook(), TestWebhookEventTemplateOutsideOrganizationIsRejected(), webhookWriteStub()

### Community 276 - "Attr"
Cohesion: 0.43
Nodes (5): exists(), likePattern(), relate(), Attr, Kind

### Community 277 - "authzAPI"
Cohesion: 0.33
Nodes (10): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), strPtr(), TestResolveOrganization() (+2 more)

### Community 280 - "logSink"
Cohesion: 0.15
Nodes (4): captureLogPublisher, LogPublisher, logSink, baseEnvironment()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 600 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **133 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Queries` to `registryStub`, `log/slog.Logger`, `run`, `time.Time`, `New`, `fakeQuerier`, `policyStub`, `store.Action`, `.marshalOBoolean2bool`, `newTestAPI`, `github.com/aminbenmansour/terrakube-go/internal/rsql.Node`, `tfe_test.go`, `addDestination`, `.CreateArchive`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `WebhookEvent`, `.CreateWorkspaceTag`, `TestWebhookQueries`, `nestedStub`, `api/main.go`, `stubQuerier`, `.CreateNotificationConfiguration`, `notificationSendAPI`, `Metrics`, `fakeStore`, `History`, `.CreateOrganization`, `CliAuthSession`, `Content`, `store.Step`, `ResolveOrganization`, `Pat`, `federatedStub`, `.CreateWebhook`, `tagStub`, `NewLayout`, `stubQuerier`, `Workspace`, `dispatcherStore`, `decodeJSON`, `Deps`, `Cipher`, `.CreateProject`, `context.Context`, `TeamToken`, `Claims`, `Querier`, `runTriggerAPI`, `Layout`, `.QueryRow`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07637178051511759 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.Request` to `net/http.ResponseWriter`, `WriteError`, `net/http.Handler`, `Deps`, `api/main.go`, `Webhook`, `net/http.ResponseWriter`, `Link`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `WriteError` be split into smaller, more focused modules?**
  _Cohesion score 0.06251213356629781 - nodes in this community are weakly interconnected._
- **Why does `New()` connect `net/http.Request` to `net/http.ResponseWriter`, `WriteError`, `run`, `strPtr`, `newTestAPI`, `.server`, `newTestExecutor`, `authzAPI`, `Webhook`, `api/main.go`, `Link`, `runTriggerAPI`, `net/http.Handler`, `notificationSendAPI`, `testing.T`, `NewLayout`, `decodeJSON`, `Deps`, `stubQuerier`, `net/http.ResponseWriter`, `runTriggerAPI`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._