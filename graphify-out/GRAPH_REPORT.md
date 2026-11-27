# Graph Report - terrakube-go  (2026-11-27)

## Corpus Check
- 584 files · ~300,933 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 5243 nodes · 21982 edges · 284 communities (142 shown, 142 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3929 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3fa64b12`
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
- Job
- logSink
- steps
- github.com/99designs/gqlgen/graphql.Marshaler
- github.com/aminbenmansour/terrakube-go/internal/rsql.Node
- rbac.go
- tfe_test.go
- sqlc Generation Config
- serveAccessToken
- Upstream Elide entities (GraphQL surface)
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- newRedactor
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
- .QueryRow
- ResolveOrganization
- CliAuthSession
- stubQuerier
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- models_gen.go
- newWorkspaceStub
- ._Organization
- Queries
- dispatcherWith
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
- .Ping
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
- log/slog.Logger
- stubQuerier
- doJSON
- .CreateTag
- .createModule
- environment
- cron.go
- NewSubscriptions
- newTestService
- context.Context
- newService
- ui-rest-inventory.sh
- ui-smoke.sh
- runTriggerAPI
- Step
- capture.sh
- 0008_action.sql
- WorkspaceTagConnection
- Service
- compare.sh
- start.sh
- Collection
- .EnqueueJob
- time.Time
- serveAccessToken
- Identity
- Parse
- Item
- clilogin/handler_test.go
- newTestAPI
- TestRotateReencryptsOldKeyAndPlaintext
- validateEndpoint
- ui.go
- API
- WebhookEvent
- NewSubscriptions
- stubStore
- stubQuerier
- Executor
- runTriggerAPI
- cancelRunStore
- tarGzEntries
- webhook_test.go
- Organization
- API
- notificationSendAPI
- github_test.go
- Metrics
- transientError
- Dispatcher
- Threat model
- .cloneConfiguration
- auth.DexValidator
- Webhook
- Resolver
- RSQL filters in SQL — design
- Upstream UI GraphQL operations
- Review Focus
- Workspace
- .build
- Repository
- TeamToken
- Secret key rotation — design
- Team
- .build
- .Query
- RelationshipID
- OrganizationConnection
- workspaceQuery
- .CreateArchive
- Ssh
- Team
- JobConnection
- github.com/aminbenmansour/terrakube-go/internal/store.Step
- newService
- Attr
- grantsAPI
- Content
- New
- discoveryStub
- .server
- WorkspaceRunTrigger
- newScheduleAPI
- newTestExecutor
- authzAPI
- .listVCSRepositories
- .PublishCapped
- .ReclaimStuckNotificationOutbox

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 285 edges
3. `Queries` - 278 edges
4. `API` - 261 edges
5. `New()` - 201 edges
6. `WriteError()` - 172 edges
7. `stubQuerier` - 171 edges
8. `param()` - 99 edges
9. `ParseFieldsets()` - 92 edges
10. `WriteResource()` - 75 edges

## Surprising Connections (you probably didn't know these)
- `Where we stand` --references--> `provider`  [INFERRED]
  docs/upstream/ui-rest-operations.md → internal/vcs/vcs.go
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Command` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Report and verification` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (284 total, 142 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (28): WriteResource(), configurationVersionResource(), projectResource(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource() (+20 more)

### Community 1 - "WriteError"
Cohesion: 0.09
Nodes (32): actionAttrs, itemAttrs, providerAttrs, providerVersionAttrs, templateAttrs, webhookEventAttrs, decodeActionAttrs(), stringOr() (+24 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (15): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+7 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "registryStub"
Cohesion: 0.04
Nodes (17): registryStub, filterVersions(), CreateImplementationParams, CreateModuleParams, CreateModuleVersionParams, CreateProviderParams, CreateProviderVersionParams, Implementation (+9 more)

### Community 7 - "run"
Cohesion: 0.05
Nodes (69): main(), run(), main(), run(), main(), run(), main(), run() (+61 more)

### Community 8 - "Handler"
Cohesion: 0.07
Nodes (37): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken() (+29 more)

### Community 9 - "newFixture"
Cohesion: 0.15
Nodes (23): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+15 more)

### Community 10 - "New"
Cohesion: 0.13
Nodes (20): branchFromRef(), deliveryID(), nonEmptyPtr(), NormalizeRepository(), viaFor(), parseGitHub(), parseGitLab(), trimDedupe() (+12 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.06
Nodes (4): nestedStub, encryptedStore, fakeQuerier, Variable

### Community 12 - "Job"
Cohesion: 0.09
Nodes (9): captureJobs(), Publisher, publishingStore, stubQuerier, dispatcherStore, FailStuckNotificationOutboxParams, Job, NotificationConfiguration (+1 more)

### Community 13 - "logSink"
Cohesion: 0.14
Nodes (4): logSink, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 16 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.05
Nodes (3): AccessConnection, AccessEdge, WorkspaceEdge

### Community 17 - "github.com/aminbenmansour/terrakube-go/internal/rsql.Node"
Cohesion: 0.06
Nodes (5): listQuery, stubLists, Task 3: Handler cutover, ListSpec, Tag

### Community 18 - "rbac.go"
Cohesion: 0.19
Nodes (27): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, adminOnly(), CanApproveJob() (+19 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (42): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+34 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "serveAccessToken"
Cohesion: 0.19
Nodes (13): sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), newTeamStub(), TestTeamJSONAPI(), newTemplateStub(), TestTemplateJSONAPI() (+5 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.21
Nodes (5): NewDispatcher(), strPtr(), Dispatcher, NotificationOutbox, WebhookDelivery

### Community 24 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

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
Nodes (63): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, createTeamTokenRequest, discoveryGroup (+55 more)

### Community 31 - "Link"
Cohesion: 0.09
Nodes (35): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), accessResource(), jobResource() (+27 more)

### Community 32 - "getJSON"
Cohesion: 0.22
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
Cohesion: 0.12
Nodes (8): notificationConfigurationResolver, organizationResolver, queryResolver, ResolverRoot, workspaceResolver, organizationScope(), TestOrganizationScope(), idsFrom()

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
Cohesion: 0.16
Nodes (18): strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased(), TestTickIgnoresPlanOnlyRun() (+10 more)

### Community 45 - ".QueryRow"
Cohesion: 0.04
Nodes (9): reliabilityStore, CreateHistoryParams, CreateProjectParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams, GetProjectByNameParams, History, UpdateProjectParams (+1 more)

### Community 46 - "ResolveOrganization"
Cohesion: 0.26
Nodes (9): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access, OrganizationStore (+1 more)

### Community 48 - "CliAuthSession"
Cohesion: 0.15
Nodes (5): fakeStore, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams, SetCliAuthSessionIdentityParams

### Community 49 - "stubQuerier"
Cohesion: 0.11
Nodes (3): fakeWorkspaceStore, stubQuerier, workspaceWithTags

### Community 51 - "time.Duration"
Cohesion: 0.09
Nodes (8): waitFor(), fakeQueue, JobQueue, defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), backoff(), Message

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (128): serveTeamTokens(), TestAccessTokenWithoutStore(), TestTeamTokenCreatePassesTheRequestAndReturnsTheToken(), TestTeamTokenCreateRejectsNonMembersAndBadInput(), TestTeamTokenDeleteMapsErrors(), TestTeamTokenListIsAccepted(), TestTeamTokenRoutesRefuseOversizedBodies(), TestTeamTokensAnswer503WithoutAService() (+120 more)

### Community 53 - "Pat"
Cohesion: 0.16
Nodes (6): newFakeStore(), fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "models_gen.go"
Cohesion: 0.05
Nodes (11): ComplexityRoot, NotificationTriggerConnection, NotificationTriggerEdge, PageInfo, ProjectConnection, ProjectEdge, Query, Subscription (+3 more)

### Community 57 - "newWorkspaceStub"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 61 - "Queries"
Cohesion: 0.02
Nodes (14): replay, startScheduler(), NewDispatcher(), NewDispatcher(), createRawVariable(), createRawVariableWithID(), TestWebhookQueries(), DeleteWebhookDeliveryParams (+6 more)

### Community 62 - "dispatcherWith"
Cohesion: 0.36
Nodes (8): dispatcherWith(), newDispatcherFixture(), TestDispatcherFailsAfterMaxAttempts(), TestDispatcherFailsPermanentRejection(), TestDispatcherIncludesTheFailureSummary(), TestDispatcherRetriesTransientFailures(), TestDispatcherSendsAndRecordsSuccess(), TestDispatcherSkipsSupersededJobStatus()

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.07
Nodes (50): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+42 more)

### Community 67 - "Action"
Cohesion: 0.15
Nodes (4): actionStub, CreateActionParams, UpdateActionParams, Action

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "Load"
Cohesion: 0.15
Nodes (24): Config, S3, Storage, Task 2: Serve archived logs from `/tfoutput/v1`, Task 1: Keyring configuration, env(), envBool(), envDuration() (+16 more)

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 72 - "New"
Cohesion: 0.21
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.17
Nodes (17): collectionSpec, column(), predicate(), Field, likePattern(), mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards() (+9 more)

### Community 76 - "streamEvents"
Cohesion: 0.24
Nodes (5): subscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 77 - "vcs/vcs.go"
Cohesion: 0.25
Nodes (14): TestGitHubConnectionsWithoutAToken(), TestGitLabConnectionWithoutAToken(), TestGitLabDiscovery(), TestProviderForRefusesALoopbackAPIURL(), baseURL(), Branches(), Test(), Groups() (+6 more)

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
Nodes (10): AccessTokenStore, Deps, PATService, stubPATService, stubTeamTokens, TeamTokenService, timePtr(), timePtr() (+2 more)

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - "Cipher"
Cohesion: 0.08
Nodes (38): Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation (+30 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "log/slog.Logger"
Cohesion: 0.13
Nodes (15): Validator, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph (+7 more)

### Community 89 - "stubQuerier"
Cohesion: 0.14
Nodes (25): missingOrganizationStub, upstreamCapture, captureID(), registryCaptureStore(), registryProviderStore(), assertSameDocumentShape(), assertSameResourceShape(), captureActionStore() (+17 more)

### Community 91 - "doJSON"
Cohesion: 0.28
Nodes (7): doJSON(), failed(), providerError(), providerMessage(), githubOwnerAndRepo(), gitlab, WebhookSpec

### Community 94 - ".createModule"
Cohesion: 0.16
Nodes (11): implementationAttrs, moduleAttrs, callerEmail(), decodeModuleAttrs(), int64Ptr(), moduleLinks(), moduleProviderError(), validModuleProvider() (+3 more)

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
Nodes (20): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+12 more)

### Community 103 - "context.Context"
Cohesion: 0.02
Nodes (18): collectionStub, globalVarStub, sshAgentStub, teamStub, templateStub, vcsStub, workspaceStub, artifactStore (+10 more)

### Community 108 - "newService"
Cohesion: 0.21
Nodes (14): decodeSecret(), New(), dexUser(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp(), TestDeleteEnforcesOwnership() (+6 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.18
Nodes (4): scheduleStub, CreateScheduleParams, Schedule, UpdateScheduleParams

### Community 121 - "Step"
Cohesion: 0.18
Nodes (3): CreateStepParams, GetStepByJobAndNumberParams, Step

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 128 - "Service"
Cohesion: 0.09
Nodes (12): Service, Store, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value() (+4 more)

### Community 133 - "Collection"
Cohesion: 0.16
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 135 - ".EnqueueJob"
Cohesion: 0.22
Nodes (7): fakeQueue, JobQueue, fakeQueue, JobQueue, JobQueue, stubJobQueue, stubQueue

### Community 137 - "time.Time"
Cohesion: 0.24
Nodes (5): bucket, limiter, backoff(), TestBackoffGrowsAndCaps(), Dispatcher

### Community 140 - "serveAccessToken"
Cohesion: 0.21
Nodes (12): stubIdentityVerifier, expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups() (+4 more)

### Community 141 - "Identity"
Cohesion: 0.10
Nodes (14): groupVerifier, stubDexVerifier, fakeVerifier, Identity, fakeIssuer, groupVerifier, stubDexVerifier, stringPtr() (+6 more)

### Community 142 - "Parse"
Cohesion: 0.34
Nodes (3): And, Or, parser

### Community 143 - "Item"
Cohesion: 0.16
Nodes (3): CreateItemParams, UpdateItemParams, Item

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.06
Nodes (43): pageMeta, createListOrganization(), getNames(), seedOrganizations(), TestOrganizationListEscapesWildcards(), TestOrganizationListPrefiltersByAllowedIDs(), TestOrganizationListSortsAndPages(), TestOrganizationListSQLFilters() (+35 more)

### Community 146 - "TestRotateReencryptsOldKeyAndPlaintext"
Cohesion: 0.09
Nodes (15): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs, assertSealedUnder() (+7 more)

### Community 147 - "validateEndpoint"
Cohesion: 0.50
Nodes (4): blockedAddress(), TestBlockedAddressRanges(), TestValidateEndpointRefusesInternalAddresses(), validateEndpoint()

### Community 149 - "ui.go"
Cohesion: 0.28
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 150 - "API"
Cohesion: 0.07
Nodes (20): API, notificationConfigurationAttrs, teamAttrs, New(), param(), parseList(), unknownCollection(), projectAccessResource() (+12 more)

### Community 151 - "WebhookEvent"
Cohesion: 0.11
Nodes (6): globToRegex(), matchBranch(), matchPaths(), matchTrigger(), pathRegex(), WebhookEvent

### Community 154 - "NewSubscriptions"
Cohesion: 0.14
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 157 - "stubStore"
Cohesion: 0.09
Nodes (5): runTriggerStub, fakeStore, CreateJobParams, Template, stubStore

### Community 159 - "stubQuerier"
Cohesion: 0.03
Nodes (9): stubQuerier, CreateReferenceParams, Agent, Globalvar, NotificationTrigger, Project, Reference, Vc (+1 more)

### Community 161 - "Executor"
Cohesion: 0.16
Nodes (6): Executor, isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestIsTerminalStatus()

### Community 162 - "runTriggerAPI"
Cohesion: 0.35
Nodes (11): newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit(), TestCreateRunTriggerRejectsCycle(), TestCreateRunTriggerRejectsStructuralViolations(), TestDeleteRunTriggerAndContentTypeQuirk() (+3 more)

### Community 163 - "cancelRunStore"
Cohesion: 0.15
Nodes (3): cancelRunStore, UpdateJobOutputParams, UpdateStepOutputParams

### Community 164 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 165 - "webhook_test.go"
Cohesion: 0.36
Nodes (10): DeleteWebhook(), SyncWebhook(), newHookServer(), TestDeleteWebhookTreats404AsRemoved(), TestGitHubWebhookCreate(), TestGitHubWebhookUpdateAndDelete(), TestGitLabWebhookCreate(), TestGitLabWebhookDelete() (+2 more)

### Community 166 - "Organization"
Cohesion: 0.28
Nodes (3): scopeIDs(), demoOrganizations(), Organization

### Community 167 - "API"
Cohesion: 0.06
Nodes (26): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+18 more)

### Community 168 - "notificationSendAPI"
Cohesion: 0.14
Nodes (23): newDeliveryStub(), TestNotificationDeliveryHistory(), TestNotificationDeliveryHistoryRejectsABadLimit(), TestNotificationDeliveryHistoryUnknownWorkspace(), TestNotificationDeliveryRetry(), TestNotificationDeliveryRetryRejections(), captureJSONServer(), notificationSendAPI() (+15 more)

### Community 171 - "github_test.go"
Cohesion: 0.44
Nodes (8): githubRoutes(), TestBranchesGitHub(), TestGroupsGitHub(), TestRepositoriesGitHub(), TestTestConnection(), newFakeProvider(), TestRejectedCredentialIsAUseful4xx(), writeJSON()

### Community 172 - "Metrics"
Cohesion: 0.13
Nodes (4): metricsSource, fakeStats, Metrics, StatsSource

### Community 175 - "Dispatcher"
Cohesion: 0.28
Nodes (8): strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued(), TestTickSkipsLockedWorkspaceAndUnparseableCron(), TestTickUsesTheTemplateOfTheWorkspaceOrganization(), tickAt(), Dispatcher

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - ".cloneConfiguration"
Cohesion: 0.13
Nodes (13): gitConfig, preparedWorkspace, TestGitSourceDetection(), TestSanitizeGitURL(), TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv() (+5 more)

### Community 184 - "auth.DexValidator"
Cohesion: 0.11
Nodes (8): Dex, Exchanger, fakeDex, contains(), decodeJWKPart(), rsaPublicKey(), stringSlice(), mutableQuerier

### Community 190 - "Webhook"
Cohesion: 0.11
Nodes (14): atomicOperation, vcsAttrs, webhookSyncRef, strPtr(), strValue(), atomicWorkspaceHref(), isAtomicRequest(), decodeVCSAttrs() (+6 more)

### Community 191 - "Resolver"
Cohesion: 0.15
Nodes (6): DB, failingDB, projectResolver, Resolver, TestOwnerRule(), loadPage()

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 197 - "Upstream UI GraphQL operations"
Cohesion: 0.17
Nodes (10): Connection shape, Entity and relationship surface, Filter expressions the UI generates, Other contracts observed, Sort strings, Upstream UI GraphQL operations, Measured route surface, Order that unblocks the most UI (+2 more)

### Community 198 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 204 - ".build"
Cohesion: 0.16
Nodes (12): listArgs, afterString(), ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset(), Order, OrderBy() (+4 more)

### Community 229 - "Repository"
Cohesion: 0.40
Nodes (4): optional(), githubRepositoryJSON, gitlabProjectJSON, Repository

### Community 231 - "TeamToken"
Cohesion: 0.17
Nodes (5): CreateTeamTokenParams, SoftDeleteTeamTokenParams, TeamToken, fakeStore, Store

### Community 232 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 234 - "Team"
Cohesion: 0.09
Nodes (9): fakeAccessTokenStore, contains(), captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces(), captureProjectAccessRow(), Access (+1 more)

### Community 235 - ".build"
Cohesion: 0.18
Nodes (10): page, accessEdges(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), projectAccessEdges(), projectEdges(), templateEdges() (+2 more)

### Community 236 - ".Query"
Cohesion: 0.06
Nodes (17): CountJobsByStatusRow, DBTX, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow (+9 more)

### Community 248 - "RelationshipID"
Cohesion: 0.21
Nodes (9): nestedDocument, DecodeAttributes(), RelationshipID(), RelationshipIdentifiers(), TestRelationshipIdentifiers(), TestRelationshipIDNull(), identifierStringAttribute(), RequestRelationship (+1 more)

### Community 252 - "workspaceQuery"
Cohesion: 0.40
Nodes (3): matchesTagFilters(), tagFilter, workspaceQuery

### Community 256 - "Team"
Cohesion: 0.09
Nodes (3): contains(), fakeOrganizationStore, Team

### Community 259 - "newService"
Cohesion: 0.23
Nodes (11): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRequiresMembershipAndAnInteractiveUser(), TestCreateSignsTheUpstreamToken(), TestListAndDeleteAreScopedToTheCallersTeams() (+3 more)

### Community 260 - "Attr"
Cohesion: 0.26
Nodes (5): exists(), likePattern(), relate(), Attr, Kind

### Community 261 - "grantsAPI"
Cohesion: 0.31
Nodes (10): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+2 more)

### Community 262 - "Content"
Cohesion: 0.23
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 271 - "New"
Cohesion: 0.49
Nodes (11): New(), githubSignature(), newStubStore(), post(), TestBranchAndPathMustMatchTheTrigger(), TestGitHubBadOrMissingSignatureIsRejected(), TestGitHubPushStartsRun(), TestGitLabPushRequiresTheToken() (+3 more)

### Community 273 - ".server"
Cohesion: 0.33
Nodes (7): fakeHooks, TestAtomicWebhookInvalidTriggerChangesNothing(), TestAtomicWebhookSaveRegistersFinalTriggers(), TestCreateWebhookEventRegistersProvider(), TestDeleteWebhookRemovesProviderHook(), TestWebhookEventTemplateOutsideOrganizationIsRejected(), webhookWriteStub()

### Community 275 - "newScheduleAPI"
Cohesion: 0.36
Nodes (9): newScheduleAPI(), newScheduleStub(), scheduleBody(), scheduleStr(), TestCreateScheduleCopiesAttributesAndDefaultsEnabled(), TestCreateScheduleRejectsForeignTemplate(), TestDeleteScheduleAndCrossWorkspaceLookup(), TestListAndUpdateSchedule() (+1 more)

### Community 276 - "newTestExecutor"
Cohesion: 0.14
Nodes (20): fakeLocker, newFakeLocker(), TestHandleReleasesTheJobLock(), TestHandleSkipsAJobHeldByAnotherExecutor(), TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestStaleClaimIdleOutlivesTheJobTimeout(), New() (+12 more)

### Community 277 - "authzAPI"
Cohesion: 0.33
Nodes (10): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), strPtr(), TestResolveOrganization() (+2 more)

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 583 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **142 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `API` connect `API` to `WriteError`, `API`, `Deps`, `.listVCSRepositories`, `api/main.go`, `RelationshipID`, `Webhook`, `.createModule`, `Link`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07709906199954245 - nodes in this community are weakly interconnected._
- **Why does `Querier` connect `Queries` to `Team`, `github.com/aminbenmansour/terrakube-go/internal/store.Step`, `Collection`, `registryStub`, `run`, `Handler`, `time.Time`, `Content`, `fakeQuerier`, `Job`, `New`, `.CreateVCS`, `Item`, `.marshalOBoolean2bool`, `TestRotateReencryptsOldKeyAndPlaintext`, `tfe_test.go`, `WorkspaceRunTrigger`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `.CreateGlobalVariable`, `.CreateWebhook`, `WebhookEvent`, `.ListNotificationDeliveriesByWorkspace`, `.ReclaimStuckNotificationOutbox`, `stubStore`, `api/main.go`, `stubQuerier`, `cancelRunStore`, `notificationSendAPI`, `.CreateNotificationConfiguration`, `Metrics`, `addDestination`, `.QueryRow`, `Dispatcher`, `CliAuthSession`, `stubQuerier`, `testing.T`, `Pat`, `tagStub`, `NewLayout`, `Action`, `Workspace`, `.CreateAgent`, `decodeJSON`, `Deps`, `Cipher`, `log/slog.Logger`, `.CreateTag`, `context.Context`, `TeamToken`, `.Query`, `runTriggerAPI`, `Step`, `.UpdateWebhook`, `.CreateArchive`, `Ssh`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `WriteError` be split into smaller, more focused modules?**
  _Cohesion score 0.08751902587519025 - nodes in this community are weakly interconnected._
- **Why does `Queries` connect `Queries` to `Team`, `Collection`, `Content`, `run`, `registryStub`, `fakeQuerier`, `Job`, `.CreateVCS`, `Item`, `newTestAPI`, `TestRotateReencryptsOldKeyAndPlaintext`, `github.com/aminbenmansour/terrakube-go/internal/rsql.Node`, `WorkspaceRunTrigger`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `.CreateGlobalVariable`, `.CreateWebhook`, `WebhookEvent`, `.ListNotificationDeliveriesByWorkspace`, `.ReclaimStuckNotificationOutbox`, `api/main.go`, `stubQuerier`, `cancelRunStore`, `.CreateNotificationConfiguration`, `addDestination`, `.QueryRow`, `CliAuthSession`, `Pat`, `Action`, `Workspace`, `.CreateAgent`, `Cipher`, `.CreateTag`, `context.Context`, `TeamToken`, `Team`, `.Query`, `runTriggerAPI`, `Step`, `.UpdateWebhook`, `.CreateArchive`, `Ssh`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._