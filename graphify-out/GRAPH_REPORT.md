# Graph Report - terrakube-go  (2026-11-28)

## Corpus Check
- 646 files · ~316,384 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 5485 nodes · 23298 edges · 283 communities (140 shown, 143 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 4204 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c5cb99fb`
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
- AccessConnection
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
- store.Step
- ResolveOrganization
- History
- CliAuthSession
- stubQuerier
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- federatedStub
- newWorkspaceStub
- Queries
- tagStub
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- workspaceStub
- models
- Load
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- collections.go
- deliver_test.go
- Organization
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
- store.Team
- stubQuerier
- doJSON
- time.Duration
- .createFederated
- environment
- cron.go
- NewSubscriptions
- newTestService
- context.Context
- Claims
- ui-rest-inventory.sh
- ui-smoke.sh
- runTriggerAPI
- StringOrInt
- capture.sh
- 0008_action.sql
- JobConnection
- Service
- compare.sh
- start.sh
- WorkspaceTagConnection
- Dispatcher
- time.Time
- strPtr
- Identity
- rsql.Node
- store.Action
- NotificationConfigurationConnection
- newTestAPI
- ui.go
- newTestExecutor
- addDestination
- net/http.Request
- WebhookEvent
- .processJob
- NewSubscriptions
- github.com/aminbenmansour/terrakube-go/internal/store.Workspace
- stubQuerier
- New
- runTriggerAPI
- tarGzEntries
- webhook_test.go
- PolicySetParameter
- net/http.Handler
- notificationSendAPI
- listQuery
- Review Focus
- Metrics
- transientError
- Threat model
- .cloneConfiguration
- auth.DexValidator
- fakeStore
- TestModelColumnOrder
- RSQL filters in SQL — design
- Upstream UI GraphQL operations
- Workspace
- OrganizationConnection
- ParsePage
- Repository
- TeamToken
- Service
- Team
- graphql/collections.go
- Querier
- RelationshipID
- Layout
- M4 run logs and artifacts Implementation Plan
- Deps
- .QueryRow
- auth.DexValidator
- newScheduleAPI
- newService
- grantsAPI
- .createVCS
- Deps
- New
- discoveryStub
- .server
- workspaceQuery
- WorkspaceRunTrigger
- Attr
- .CreateArchive
- TestWebhookQueries
- logSink
- .SoftDeleteTeamToken

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 322 edges
3. `Queries` - 312 edges
4. `API` - 307 edges
5. `New()` - 238 edges
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
- `RBAC: permission resolution and `/access-token/v1` Implementation Plan` --references--> `Flags`  [INFERRED]
  docs/superpowers/plans/2026-11-24-rbac-permission-resolution.md → internal/rbac/rbac.go
- `Task 3: Transparent store decorator` --references--> `fakeQuerier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/secretstore/store_test.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (283 total, 143 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (27): WriteResource(), configurationVersionResource(), projectResource(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource() (+19 more)

### Community 1 - "WriteError"
Cohesion: 0.05
Nodes (66): actionAttrs, atomicOperation, federatedClaimAttrs, implementationAttrs, itemAttrs, moduleAttrs, notificationConfigurationAttrs, policySetAttrs (+58 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (14): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+6 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "registryStub"
Cohesion: 0.04
Nodes (17): registryStub, filterVersions(), CreateImplementationParams, CreateModuleParams, CreateModuleVersionParams, CreateProviderParams, CreateProviderVersionParams, Implementation (+9 more)

### Community 7 - "run"
Cohesion: 0.05
Nodes (71): main(), run(), main(), run(), main(), run(), main(), run() (+63 more)

### Community 8 - "Handler"
Cohesion: 0.24
Nodes (6): isExternalUser(), matchesMD5(), New(), newUUID(), pathIDs(), Handler

### Community 9 - "newFixture"
Cohesion: 0.15
Nodes (23): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+15 more)

### Community 10 - "New"
Cohesion: 0.10
Nodes (23): globToRegex(), matchBranch(), matchPaths(), matchTrigger(), pathRegex(), branchFromRef(), deliveryID(), nonEmptyPtr() (+15 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.07
Nodes (3): encryptedStore, fakeQuerier, Variable

### Community 12 - "policyStub"
Cohesion: 0.06
Nodes (5): policyStub, PolicyAttachment, PolicyExemption, PolicySet, WebhookDelivery

### Community 13 - "Executor"
Cohesion: 0.19
Nodes (5): Executor, preparedWorkspace, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 17 - "github.com/aminbenmansour/terrakube-go/internal/rsql.Node"
Cohesion: 0.07
Nodes (5): listQuery, stubLists, Task 3: Handler cutover, Page, ListSpec

### Community 18 - "rbac.go"
Cohesion: 0.26
Nodes (22): Task 1: Role matrix and permission set, adminOnly(), CanApproveJob(), CanDestroyJob(), CanManageCollection(), CanManageJob(), CanManageModule(), CanManagePolicies() (+14 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (42): Backend, assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), extractLogReadURL(), newTestAPI() (+34 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.04
Nodes (11): NotificationTriggerConnection, NotificationTriggerEdge, PageInfo, ProjectConnection, ProjectEdge, Query, Subscription, TemplateConnection (+3 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.23
Nodes (4): NewDispatcher(), strPtr(), Dispatcher, NotificationOutbox

### Community 24 - "Webhook"
Cohesion: 0.09
Nodes (13): webhookEventAttrs, webhookSyncRef, strValue(), atomicWorkspaceHref(), isAtomicRequest(), discoveryPage(), writeRawJSON(), isWebhookAtomicHref() (+5 more)

### Community 25 - "payload_test.go"
Cohesion: 0.15
Nodes (30): channelLabel(), networkCause(), parseRetryAfter(), BuildPayload(), buildSlack(), buildTeams(), buildWebhook(), capitalize() (+22 more)

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
Nodes (72): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, createTeamTokenRequest, discoveryGroup (+64 more)

### Community 31 - "Link"
Cohesion: 0.07
Nodes (45): policySetChildren, actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), federatedClaimResource(), federatedResource() (+37 more)

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
Cohesion: 0.07
Nodes (8): captureJobs(), Publisher, publishingStore, stubQuerier, fakeStore, CreateJobParams, GetJobByWorkspaceAndTriggerParams, stubStore

### Community 45 - "store.Step"
Cohesion: 0.08
Nodes (7): captureSteps(), dispatcherStore, CreateStepParams, FailStuckNotificationOutboxParams, GetStepByJobAndNumberParams, ReclaimStuckNotificationOutboxParams, RecordNotificationOutboxParams

### Community 46 - "ResolveOrganization"
Cohesion: 0.19
Nodes (12): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveOrganization(), TestResolveOrganizationUsesTheRequestedOrganization(), TestResolveWorkspaceAddsGrantRows() (+4 more)

### Community 47 - "History"
Cohesion: 0.05
Nodes (10): cancelRunStore, reliabilityStore, CreateHistoryParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams, History, UpdateJobOutputParams, UpdateStepOutputParams (+2 more)

### Community 48 - "CliAuthSession"
Cohesion: 0.10
Nodes (6): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams, SetCliAuthSessionIdentityParams

### Community 49 - "stubQuerier"
Cohesion: 0.09
Nodes (4): CreateContentParams, UpdateContentStatusParams, Content, stubQuerier

### Community 52 - "testing.T"
Cohesion: 0.02
Nodes (141): expectPermissions(), serveAccessToken(), serveTeamTokens(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups() (+133 more)

### Community 53 - "Pat"
Cohesion: 0.16
Nodes (6): newFakeStore(), fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "federatedStub"
Cohesion: 0.10
Nodes (7): federatedStub, CreateFederatedClaimParams, CreateFederatedParams, Federated, FederatedClaim, UpdateFederatedClaimParams, UpdateFederatedParams

### Community 57 - "newWorkspaceStub"
Cohesion: 0.09
Nodes (37): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+29 more)

### Community 61 - "Queries"
Cohesion: 0.02
Nodes (11): replay, startScheduler(), NewDispatcher(), NewDispatcher(), FailStuckRunTriggerEventsParams, Querier, ReclaimStuckRunTriggerEventsParams, RecordWebhookDeliveryParams (+3 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.07
Nodes (46): redactor, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+38 more)

### Community 67 - "workspaceStub"
Cohesion: 0.07
Nodes (5): workspaceStub, Agent, Project, Ssh, Vc

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "Load"
Cohesion: 0.17
Nodes (20): Config, Task 1: Keyring configuration, env(), envBool(), envFirst(), envInt64(), envList(), Load() (+12 more)

### Community 70 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.11
Nodes (3): jobStub, statusQuerier, stubQuerier

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - "deliver_test.go"
Cohesion: 0.12
Nodes (27): NewDeliverer(), deliveryError(), stringPtr(), TestDeliverDoesNotFollowRedirects(), testDeliverer(), TestDeliverIgnoresMalformedRetryAfter(), TestDeliverNetworkErrorIsRetryableAndHidesTheURL(), TestDeliverParsesHTTPDateRetryAfter() (+19 more)

### Community 73 - "Organization"
Cohesion: 0.18
Nodes (4): missingOrganizationStub, scopeIDs(), demoOrganizations(), Organization

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
Cohesion: 0.18
Nodes (23): githubRoutes(), TestBranchesGitHub(), TestGitHubConnectionsWithoutAToken(), TestGroupsGitHub(), TestRepositoriesGitHub(), TestTestConnection(), TestGitLabConnectionWithoutAToken(), TestGitLabDiscovery() (+15 more)

### Community 78 - "decodeJSON"
Cohesion: 0.14
Nodes (30): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent() (+22 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "auth/dex_test.go"
Cohesion: 0.07
Nodes (39): jwksServer, CookieSigner, Deps, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey() (+31 more)

### Community 81 - "Deps"
Cohesion: 0.24
Nodes (6): stubPATService, stubTeamTokens, timePtr(), timePtr(), Token, Token

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - "Cipher"
Cohesion: 0.04
Nodes (67): Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation (+59 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 89 - "stubQuerier"
Cohesion: 0.11
Nodes (33): upstreamCapture, authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), federatedCaptureStore() (+25 more)

### Community 91 - "doJSON"
Cohesion: 0.26
Nodes (8): doJSON(), failed(), providerError(), providerMessage(), githubOwnerAndRepo(), Error, gitlab, WebhookSpec

### Community 92 - "time.Duration"
Cohesion: 0.11
Nodes (9): waitFor(), envDuration(), defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), backoff(), blockedDelivery(), DeliveryError (+1 more)

### Community 94 - ".createFederated"
Cohesion: 0.47
Nodes (4): federatedAttrs, decodeFederatedAttrs(), federatedNotBlank(), federatedRequired()

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 99 - "cron.go"
Cohesion: 0.27
Nodes (8): at(), Parse(), parseField(), splitStep(), mustParse(), TestParseMatchesQuartzFields(), TestParseRejectsWhatQuartzRejects(), Cron

### Community 100 - "NewSubscriptions"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 103 - "context.Context"
Cohesion: 0.02
Nodes (18): collectionStub, globalVarStub, sshAgentStub, teamStub, templateStub, vcsStub, artifactStore, historyStub (+10 more)

### Community 108 - "Claims"
Cohesion: 0.15
Nodes (14): Claims, Type, ClaimsFromContext(), Classify(), TestClassify(), TestEffectiveStepIDFromJTI(), WithClaims(), runnerClaims() (+6 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.17
Nodes (4): scheduleStub, CreateScheduleParams, Schedule, UpdateScheduleParams

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 128 - "Service"
Cohesion: 0.16
Nodes (9): Service, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), randomToken(), sanitizeName() (+1 more)

### Community 135 - "Dispatcher"
Cohesion: 0.12
Nodes (16): strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued(), TestTickSkipsLockedWorkspaceAndUnparseableCron(), TestTickUsesTheTemplateOfTheWorkspaceOrganization(), tickAt(), fakeQueue (+8 more)

### Community 137 - "time.Time"
Cohesion: 0.18
Nodes (9): policyExemptionAttrs, bucket, limiter, newLimiter(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), backoff(), TestBackoffGrowsAndCaps() (+1 more)

### Community 140 - "strPtr"
Cohesion: 0.19
Nodes (13): sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), newTeamStub(), TestTeamJSONAPI(), newTemplateStub(), TestTemplateJSONAPI() (+5 more)

### Community 141 - "Identity"
Cohesion: 0.12
Nodes (12): groupVerifier, stubDexVerifier, stubIdentityVerifier, fakeVerifier, Identity, fakeIssuer, TokenIssuer, groupVerifier (+4 more)

### Community 142 - "rsql.Node"
Cohesion: 0.34
Nodes (3): And, Or, parser

### Community 143 - "store.Action"
Cohesion: 0.17
Nodes (3): actionStub, CreateActionParams, UpdateActionParams

### Community 145 - "newTestAPI"
Cohesion: 0.05
Nodes (60): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+52 more)

### Community 146 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 147 - "newTestExecutor"
Cohesion: 0.10
Nodes (22): Deps, fakeLocker, JobLocker, newFakeLocker(), TestHandleReleasesTheJobLock(), TestHandleSkipsAJobHeldByAnotherExecutor(), TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs() (+14 more)

### Community 149 - "addDestination"
Cohesion: 0.37
Nodes (15): strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased(), TestTickIgnoresPlanOnlyRun() (+7 more)

### Community 150 - "net/http.Request"
Cohesion: 0.07
Nodes (11): API, New(), param(), parseList(), unknownCollection(), projectAccessResource(), projectResource(), tagResource() (+3 more)

### Community 152 - ".processJob"
Cohesion: 0.16
Nodes (5): isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestIsTerminalStatus()

### Community 154 - "NewSubscriptions"
Cohesion: 0.14
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 157 - "github.com/aminbenmansour/terrakube-go/internal/store.Workspace"
Cohesion: 0.10
Nodes (3): nestedStub, fakeWorkspaceStore, workspaceWithTags

### Community 159 - "stubQuerier"
Cohesion: 0.03
Nodes (3): stubQuerier, UpdateCollectionParams, UpdateItemParams

### Community 161 - "New"
Cohesion: 0.17
Nodes (14): ExtractToken(), TestExtractToken(), TestWriteBackendOverride(), isTerminal(), New(), outputOf(), stepOf(), streamKey() (+6 more)

### Community 162 - "runTriggerAPI"
Cohesion: 0.24
Nodes (13): runTriggerStub, newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), runTriggerRBACAPI(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit(), TestCreateRunTriggerRejectsCycle() (+5 more)

### Community 164 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 165 - "webhook_test.go"
Cohesion: 0.36
Nodes (10): DeleteWebhook(), SyncWebhook(), newHookServer(), TestDeleteWebhookTreats404AsRemoved(), TestGitHubWebhookCreate(), TestGitHubWebhookUpdateAndDelete(), TestGitLabWebhookCreate(), TestGitLabWebhookDelete() (+2 more)

### Community 167 - "net/http.Handler"
Cohesion: 0.07
Nodes (23): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+15 more)

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
Cohesion: 0.13
Nodes (4): metricsSource, fakeStats, Metrics, StatsSource

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - ".cloneConfiguration"
Cohesion: 0.14
Nodes (12): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv(), gitSource() (+4 more)

### Community 184 - "auth.DexValidator"
Cohesion: 0.21
Nodes (3): Dex, Exchanger, fakeDex

### Community 190 - "fakeStore"
Cohesion: 0.19
Nodes (3): fakeStore, RecordRunTriggerEventParams, RunTriggerEvent

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
Cohesion: 0.23
Nodes (4): CreateTeamTokenParams, TeamToken, fakeStore, Store

### Community 232 - "Service"
Cohesion: 0.24
Nodes (4): contains(), stringPtr(), stringSlice(), Service

### Community 234 - "Team"
Cohesion: 0.09
Nodes (9): AccessTokenStore, fakeAccessTokenStore, contains(), captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces(), captureProjectAccessRow() (+1 more)

### Community 235 - "graphql/collections.go"
Cohesion: 0.10
Nodes (20): collectionSpec, DB, listArgs, page, accessEdges(), afterString(), jobEdges(), loadPage() (+12 more)

### Community 236 - "Querier"
Cohesion: 0.03
Nodes (18): rotationWorkspace(), CountJobsByStatusRow, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow (+10 more)

### Community 248 - "RelationshipID"
Cohesion: 0.22
Nodes (9): nestedDocument, RelationshipIdentifiers(), RelationshipList(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship, ResourceIdentifier (+1 more)

### Community 251 - "Layout"
Cohesion: 0.24
Nodes (3): Validator, Layout, Deps

### Community 252 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.18
Nodes (9): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+1 more)

### Community 254 - "Deps"
Cohesion: 0.29
Nodes (3): Deps, PATService, TeamTokenService

### Community 255 - ".QueryRow"
Cohesion: 0.02
Nodes (6): CreateCollectionParams, CreateProjectParams, GetProjectByNameParams, GetTagByNameParams, UpdateOrganizationParams, UpdateProjectParams

### Community 256 - "auth.DexValidator"
Cohesion: 0.31
Nodes (4): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 257 - "newScheduleAPI"
Cohesion: 0.36
Nodes (9): newScheduleAPI(), newScheduleStub(), scheduleBody(), scheduleStr(), TestCreateScheduleCopiesAttributesAndDefaultsEnabled(), TestCreateScheduleRejectsForeignTemplate(), TestDeleteScheduleAndCrossWorkspaceLookup(), TestListAndUpdateSchedule() (+1 more)

### Community 259 - "newService"
Cohesion: 0.23
Nodes (11): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRequiresMembershipAndAnInteractiveUser(), TestCreateSignsTheUpstreamToken(), TestListAndDeleteAreScopedToTheCallersTeams() (+3 more)

### Community 261 - "grantsAPI"
Cohesion: 0.31
Nodes (10): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+2 more)

### Community 262 - ".createVCS"
Cohesion: 0.48
Nodes (5): vcsAttrs, strPtr(), decodeVCSAttrs(), validVCSEnum(), vcsEnumOrDefault()

### Community 270 - "Deps"
Cohesion: 0.33
Nodes (3): Deps, fakePinger, Pinger

### Community 271 - "New"
Cohesion: 0.49
Nodes (11): New(), githubSignature(), newStubStore(), post(), TestBranchAndPathMustMatchTheTrigger(), TestGitHubBadOrMissingSignatureIsRejected(), TestGitHubPushStartsRun(), TestGitLabPushRequiresTheToken() (+3 more)

### Community 273 - ".server"
Cohesion: 0.33
Nodes (7): fakeHooks, TestAtomicWebhookInvalidTriggerChangesNothing(), TestAtomicWebhookSaveRegistersFinalTriggers(), TestCreateWebhookEventRegistersProvider(), TestDeleteWebhookRemovesProviderHook(), TestWebhookEventTemplateOutsideOrganizationIsRejected(), webhookWriteStub()

### Community 274 - "workspaceQuery"
Cohesion: 0.33
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 276 - "Attr"
Cohesion: 0.26
Nodes (5): exists(), likePattern(), relate(), Attr, Kind

### Community 280 - "logSink"
Cohesion: 0.20
Nodes (3): captureLogPublisher, LogPublisher, logSink

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 599 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **143 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `API` connect `net/http.Request` to `WriteError`, `.createVCS`, `net/http.Handler`, `api/main.go`, `Deps`, `Webhook`, `.createFederated`, `Link`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07760635811126694 - nodes in this community are weakly interconnected._
- **Why does `Querier` connect `Queries` to `registryStub`, `run`, `Dispatcher`, `time.Time`, `fakeQuerier`, `policyStub`, `store.Action`, `.marshalOBoolean2bool`, `newTestAPI`, `tfe_test.go`, `WorkspaceRunTrigger`, `.UpdateWebhook`, `.CreateArchive`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `WebhookEvent`, `TestWebhookQueries`, `.SoftDeleteTeamToken`, `github.com/aminbenmansour/terrakube-go/internal/store.Workspace`, `api/main.go`, `stubQuerier`, `New`, `.CreateNotificationConfiguration`, `PolicySetParameter`, `notificationSendAPI`, `Metrics`, `store.Step`, `fakeStore`, `History`, `CliAuthSession`, `stubQuerier`, `._NotificationTrigger`, `testing.T`, `Pat`, `federatedStub`, `newWorkspaceStub`, `fakeStore`, `tagStub`, `workspaceStub`, `github.com/aminbenmansour/terrakube-go/internal/store.Job`, `Workspace`, `Cipher`, `store.Team`, `context.Context`, `TeamToken`, `Claims`, `Querier`, `runTriggerAPI`, `.CreateGlobalVariable`, `Deps`, `.QueryRow`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Should `WriteError` be split into smaller, more focused modules?**
  _Cohesion score 0.05371184859389577 - nodes in this community are weakly interconnected._
- **Why does `Layout` connect `Layout` to `New`, `NewLayout`, `newTestExecutor`, `M4 run logs and artifacts Implementation Plan`, `Queries`, `api/main.go`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._