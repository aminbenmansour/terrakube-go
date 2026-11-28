# Graph Report - terrakube-go  (2026-11-28)

## Corpus Check
- 619 files · ~311,644 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 5404 nodes · 22839 edges · 283 communities (137 shown, 146 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 4131 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a0d106cc`
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
- logSink
- steps
- github.com/99designs/gqlgen/graphql.Marshaler
- github.com/aminbenmansour/terrakube-go/internal/rsql.Node
- rbac.go
- tfe_test.go
- sqlc Generation Config
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
- stubQuerier
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- ._Template
- newWorkspaceStub
- ._Organization
- Queries
- dispatcherWith
- tagStub
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- models
- Load
- github.com/aminbenmansour/terrakube-go/internal/store.Job
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
- Deps
- stubQuerier
- doJSON
- Tag
- .atomicOperations
- environment
- cron.go
- NewSubscriptions
- newTestService
- context.Context
- Claims
- ui-rest-inventory.sh
- ui-smoke.sh
- runTriggerAPI
- store.Job
- capture.sh
- 0008_action.sql
- OrganizationConnection
- Service
- compare.sh
- start.sh
- time.Duration
- Dispatcher
- time.Time
- Identity
- Parse
- store.Action
- clilogin/handler_test.go
- newTestAPI
- New
- addDestination
- net/http.Request
- WebhookEvent
- ModuleVersion
- NewSubscriptions
- nestedStub
- stubQuerier
- Executor
- runTriggerAPI
- tarGzEntries
- webhook_test.go
- ._Project
- API
- notificationSendAPI
- listQuery
- Review Focus
- Metrics
- transientError
- .CreateOrganization
- Threat model
- .cloneConfiguration
- auth.DexValidator
- store.Collection
- Resolver
- RSQL filters in SQL — design
- Upstream UI GraphQL operations
- Organization
- Workspace
- M4 run logs and artifacts Implementation Plan
- .build
- Repository
- TeamToken
- Secret key rotation — design
- store.Access
- .build
- Querier
- RelationshipID
- Layout
- github_test.go
- WorkspaceRunTrigger
- seedWorkspaces
- .QueryRow
- Team
- store.Reference
- Project
- newService
- .listVCSRepositories
- grantsAPI
- store.Item
- New
- discoveryStub
- .server
- newScheduleAPI
- Attr
- authzAPI
- .CreateArchive
- TestWebhookQueries
- logSink
- FailStuckRunTriggerEventsParams

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 309 edges
3. `Queries` - 301 edges
4. `API` - 295 edges
5. `New()` - 226 edges
6. `WriteError()` - 188 edges
7. `stubQuerier` - 171 edges
8. `param()` - 107 edges
9. `ParseFieldsets()` - 104 edges
10. `WriteResource()` - 82 edges

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

## Communities (283 total, 146 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (29): WriteResource(), configurationVersionResource(), projectResource(), parseWorkspaceQuery(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID() (+21 more)

### Community 1 - "WriteError"
Cohesion: 0.07
Nodes (56): actionAttrs, implementationAttrs, itemAttrs, moduleAttrs, policySetAttrs, policySetParameterAttrs, providerAttrs, providerVersionAttrs (+48 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (16): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+8 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "registryStub"
Cohesion: 0.05
Nodes (14): registryStub, filterVersions(), CreateImplementationParams, CreateModuleParams, CreateProviderParams, CreateProviderVersionParams, Implementation, Module (+6 more)

### Community 7 - "run"
Cohesion: 0.07
Nodes (53): main(), run(), main(), run(), main(), run(), formatOldKeys(), logReport() (+45 more)

### Community 8 - "Handler"
Cohesion: 0.24
Nodes (6): isExternalUser(), matchesMD5(), New(), newUUID(), pathIDs(), Handler

### Community 9 - "newFixture"
Cohesion: 0.09
Nodes (33): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+25 more)

### Community 10 - "New"
Cohesion: 0.14
Nodes (18): branchFromRef(), deliveryID(), nonEmptyPtr(), NormalizeRepository(), viaFor(), parseGitHub(), parseGitLab(), trimDedupe() (+10 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.07
Nodes (3): encryptedStore, fakeQuerier, Variable

### Community 12 - "policyStub"
Cohesion: 0.05
Nodes (5): policyStub, PolicyAttachment, PolicyExemption, PolicySet, PolicySetParameter

### Community 13 - "logSink"
Cohesion: 0.19
Nodes (5): Executor, preparedWorkspace, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 16 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.05
Nodes (11): AccessConnection, AccessEdge, JobConnection, JobEdge, NotificationTriggerConnection, NotificationTriggerEdge, PageInfo, Query (+3 more)

### Community 18 - "rbac.go"
Cohesion: 0.26
Nodes (22): Task 1: Role matrix and permission set, adminOnly(), CanApproveJob(), CanDestroyJob(), CanManageCollection(), CanManageJob(), CanManageModule(), CanManagePolicies() (+14 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.08
Nodes (42): Backend, assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), extractLogReadURL(), newTestAPI() (+34 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.21
Nodes (5): NewDispatcher(), strPtr(), Dispatcher, NotificationOutbox, WebhookDelivery

### Community 24 - "Webhook"
Cohesion: 0.10
Nodes (11): atomicOperation, vcsAttrs, webhookEventAttrs, webhookSyncRef, strPtr(), strValue(), decodeVCSAttrs(), requiredVCSString() (+3 more)

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
Nodes (72): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, createTeamTokenRequest, discoveryGroup (+64 more)

### Community 31 - "Link"
Cohesion: 0.08
Nodes (41): policySetChildren, actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), accessResource() (+33 more)

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

### Community 44 - "fakeStore"
Cohesion: 0.08
Nodes (7): fakeStore, fakeStore, CreateJobParams, CreateStepParams, RecordRunTriggerEventParams, RunTriggerEvent, stubStore

### Community 45 - "History"
Cohesion: 0.11
Nodes (5): reliabilityStore, CreateHistoryParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams, History

### Community 46 - "ResolveOrganization"
Cohesion: 0.12
Nodes (11): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access, fakeWorkspaceStore (+3 more)

### Community 47 - "store.Step"
Cohesion: 0.07
Nodes (6): cancelRunStore, captureSteps(), GetStepByJobAndNumberParams, UpdateJobOutputParams, UpdateStepOutputParams, UpdateStepStatusParams

### Community 48 - "CliAuthSession"
Cohesion: 0.11
Nodes (6): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams, SetCliAuthSessionIdentityParams

### Community 49 - "stubQuerier"
Cohesion: 0.10
Nodes (4): CreateContentParams, UpdateContentStatusParams, Content, stubQuerier

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (128): expectPermissions(), serveAccessToken(), serveTeamTokens(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups() (+120 more)

### Community 53 - "Pat"
Cohesion: 0.16
Nodes (6): newFakeStore(), fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 57 - "newWorkspaceStub"
Cohesion: 0.10
Nodes (34): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+26 more)

### Community 59 - "._Organization"
Cohesion: 0.07
Nodes (6): ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, ProjectConnection, WorkspaceConnection, WorkspaceEdge

### Community 61 - "Queries"
Cohesion: 0.03
Nodes (7): createRawVariable(), createRawVariableWithID(), Queries, ReclaimStuckRunTriggerEventsParams, RecordWebhookDeliveryParams, SoftDeleteTeamTokenParams, UpdateHistoryOutputParams

### Community 62 - "dispatcherWith"
Cohesion: 0.36
Nodes (8): dispatcherWith(), newDispatcherFixture(), TestDispatcherFailsAfterMaxAttempts(), TestDispatcherFailsPermanentRejection(), TestDispatcherIncludesTheFailureSummary(), TestDispatcherRetriesTransientFailures(), TestDispatcherSendsAndRecordsSuccess(), TestDispatcherSkipsSupersededJobStatus()

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.05
Nodes (65): redactor, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+57 more)

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
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 72 - "New"
Cohesion: 0.21
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 73 - "dispatcherStore"
Cohesion: 0.09
Nodes (5): dispatcherStore, FailStuckNotificationOutboxParams, NotificationConfiguration, ReclaimStuckNotificationOutboxParams, RecordNotificationOutboxParams

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
Cohesion: 0.31
Nodes (13): TestGitHubConnectionsWithoutAToken(), TestGitLabConnectionWithoutAToken(), TestGitLabDiscovery(), TestProviderForRefusesALoopbackAPIURL(), baseURL(), Branches(), Test(), Groups() (+5 more)

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
Cohesion: 0.29
Nodes (5): stubTeamTokens, timePtr(), timePtr(), Token, Token

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - "Cipher"
Cohesion: 0.06
Nodes (53): Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation (+45 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "Deps"
Cohesion: 0.29
Nodes (3): Deps, fakeLocker, JobLocker

### Community 89 - "stubQuerier"
Cohesion: 0.14
Nodes (25): upstreamCapture, capturePathLastID(), captureRelationshipID(), policyCaptureStore(), captureID(), registryCaptureStore(), registryProviderStore(), assertSameDocumentShape() (+17 more)

### Community 91 - "doJSON"
Cohesion: 0.22
Nodes (10): TestValidateEndpointRefusesInternalAddresses(), validateEndpoint(), doJSON(), failed(), providerError(), providerMessage(), githubOwnerAndRepo(), Error (+2 more)

### Community 94 - ".atomicOperations"
Cohesion: 0.18
Nodes (5): atomicWorkspaceHref(), isAtomicRequest(), isWebhookAtomicHref(), parseWebhookHref(), RelationshipIdentifier()

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
Nodes (18): collectionStub, globalVarStub, sshAgentStub, teamStub, templateStub, vcsStub, workspaceStub, artifactStore (+10 more)

### Community 108 - "Claims"
Cohesion: 0.13
Nodes (16): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken() (+8 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.18
Nodes (4): scheduleStub, CreateScheduleParams, Schedule, UpdateScheduleParams

### Community 121 - "store.Job"
Cohesion: 0.11
Nodes (6): captureJobs(), Publisher, publishingStore, stubQuerier, GetJobByWorkspaceAndTriggerParams, mutableQuerier

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 128 - "Service"
Cohesion: 0.15
Nodes (10): Service, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), randomToken() (+2 more)

### Community 133 - "time.Duration"
Cohesion: 0.13
Nodes (9): waitFor(), envDuration(), defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), backoff(), backoff(), TestBackoffGrowsAndCaps() (+1 more)

### Community 135 - "Dispatcher"
Cohesion: 0.12
Nodes (16): strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued(), TestTickSkipsLockedWorkspaceAndUnparseableCron(), TestTickUsesTheTemplateOfTheWorkspaceOrganization(), tickAt(), fakeQueue (+8 more)

### Community 137 - "time.Time"
Cohesion: 0.22
Nodes (6): policyExemptionAttrs, bucket, limiter, pgtypeTime(), TestNotificationOutboxClaimRecordAndSweep(), Dispatcher

### Community 141 - "Identity"
Cohesion: 0.09
Nodes (16): AccessTokenStore, Deps, groupVerifier, PATService, stubDexVerifier, stubIdentityVerifier, stubPATService, TeamTokenService (+8 more)

### Community 142 - "Parse"
Cohesion: 0.18
Nodes (7): Parse(), dump(), TestParse(), TestParseRejects(), And, Or, parser

### Community 143 - "store.Action"
Cohesion: 0.17
Nodes (3): actionStub, CreateActionParams, UpdateActionParams

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.05
Nodes (50): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+42 more)

### Community 146 - "New"
Cohesion: 0.20
Nodes (12): TestWriteBackendOverride(), isTerminal(), New(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame() (+4 more)

### Community 149 - "addDestination"
Cohesion: 0.37
Nodes (15): strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased(), TestTickIgnoresPlanOnlyRun() (+7 more)

### Community 150 - "net/http.Request"
Cohesion: 0.07
Nodes (12): API, New(), param(), parseList(), unknownCollection(), projectAccessResource(), projectResource(), writeInUseConflict() (+4 more)

### Community 151 - "WebhookEvent"
Cohesion: 0.11
Nodes (6): globToRegex(), matchBranch(), matchPaths(), matchTrigger(), pathRegex(), WebhookEvent

### Community 152 - "ModuleVersion"
Cohesion: 0.22
Nodes (3): CreateModuleVersionParams, ModuleVersion, UpdateModuleVersionParams

### Community 154 - "NewSubscriptions"
Cohesion: 0.08
Nodes (33): Deps, decodeBody(), fakeGitHub(), TestVCSDiscoveryDeniedWithoutPermission(), TestVCSDiscoveryGroups(), TestVCSDiscoveryRejectsBadRequests(), TestVCSDiscoveryRepositoriesAndBranches(), TestVCSDiscoveryTestConnection() (+25 more)

### Community 161 - "Executor"
Cohesion: 0.16
Nodes (5): isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestIsTerminalStatus()

### Community 162 - "runTriggerAPI"
Cohesion: 0.23
Nodes (12): runTriggerStub, newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit(), TestCreateRunTriggerRejectsCycle(), TestCreateRunTriggerRejectsStructuralViolations() (+4 more)

### Community 164 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 165 - "webhook_test.go"
Cohesion: 0.36
Nodes (10): DeleteWebhook(), SyncWebhook(), newHookServer(), TestDeleteWebhookTreats404AsRemoved(), TestGitHubWebhookCreate(), TestGitHubWebhookUpdateAndDelete(), TestGitLabWebhookCreate(), TestGitLabWebhookDelete() (+2 more)

### Community 167 - "API"
Cohesion: 0.06
Nodes (29): authorization, main(), run(), Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom() (+21 more)

### Community 168 - "notificationSendAPI"
Cohesion: 0.14
Nodes (23): newDeliveryStub(), TestNotificationDeliveryHistory(), TestNotificationDeliveryHistoryRejectsABadLimit(), TestNotificationDeliveryHistoryUnknownWorkspace(), TestNotificationDeliveryRetry(), TestNotificationDeliveryRetryRejections(), captureJSONServer(), notificationSendAPI() (+15 more)

### Community 169 - "listQuery"
Cohesion: 0.17
Nodes (9): listQuery, Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 3: Handler cutover (+1 more)

### Community 171 - "Review Focus"
Cohesion: 0.33
Nodes (5): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract

### Community 172 - "Metrics"
Cohesion: 0.06
Nodes (21): metricsSource, Deps, fakePinger, fakeStats, Metrics, Pinger, StatsSource, requireMetricsToken() (+13 more)

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - ".cloneConfiguration"
Cohesion: 0.14
Nodes (12): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv(), gitSource() (+4 more)

### Community 184 - "auth.DexValidator"
Cohesion: 0.13
Nodes (7): Dex, Exchanger, fakeDex, contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 191 - "Resolver"
Cohesion: 0.15
Nodes (6): DB, failingDB, projectResolver, Resolver, TestOwnerRule(), loadPage()

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 197 - "Upstream UI GraphQL operations"
Cohesion: 0.17
Nodes (10): Connection shape, Entity and relationship surface, Filter expressions the UI generates, Other contracts observed, Sort strings, Upstream UI GraphQL operations, Measured route surface, Order that unblocks the most UI (+2 more)

### Community 198 - "Organization"
Cohesion: 0.20
Nodes (4): missingOrganizationStub, scopeIDs(), demoOrganizations(), Organization

### Community 202 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.18
Nodes (9): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+1 more)

### Community 204 - ".build"
Cohesion: 0.16
Nodes (12): listArgs, afterString(), ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset(), Order, OrderBy() (+4 more)

### Community 229 - "Repository"
Cohesion: 0.40
Nodes (4): optional(), githubRepositoryJSON, gitlabProjectJSON, Repository

### Community 231 - "TeamToken"
Cohesion: 0.21
Nodes (5): contains(), CreateTeamTokenParams, TeamToken, fakeStore, Store

### Community 232 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 235 - ".build"
Cohesion: 0.18
Nodes (10): page, accessEdges(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), projectAccessEdges(), projectEdges(), templateEdges() (+2 more)

### Community 236 - "Querier"
Cohesion: 0.03
Nodes (23): replay, startScheduler(), NewDispatcher(), NewDispatcher(), CountJobsByStatusRow, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow (+15 more)

### Community 248 - "RelationshipID"
Cohesion: 0.17
Nodes (12): nestedDocument, notificationConfigurationAttrs, decodeNotificationConfigurationAttrs(), DecodeAttributes(), RelationshipIdentifiers(), RelationshipList(), identifierStringAttribute(), validateTagBindings() (+4 more)

### Community 251 - "Layout"
Cohesion: 0.24
Nodes (3): Validator, Layout, Deps

### Community 252 - "github_test.go"
Cohesion: 0.35
Nodes (10): githubRoutes(), TestBranchesGitHub(), TestGroupsGitHub(), TestRepositoriesGitHub(), TestTestConnection(), gitlabRoutes(), newFakeProvider(), requireBearer() (+2 more)

### Community 254 - "seedWorkspaces"
Cohesion: 0.31
Nodes (9): ptr(), seedWorkspaces(), TestListPageBetween(), TestListPageFilters(), TestListPagePages(), TestListPageRejectsBadValues(), TestListPageRelationshipFilters(), TestListPageSortsByJobDate() (+1 more)

### Community 255 - ".QueryRow"
Cohesion: 0.02
Nodes (4): CreateProjectParams, GetProjectByNameParams, GetTagByNameParams, UpdateProjectParams

### Community 258 - "Project"
Cohesion: 0.22
Nodes (6): captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces(), captureProjectAccessRow(), Project

### Community 259 - "newService"
Cohesion: 0.14
Nodes (15): contains(), decodeSecret(), New(), stringPtr(), stringSlice(), dexUser(), newFakeStore(), newService() (+7 more)

### Community 261 - "grantsAPI"
Cohesion: 0.31
Nodes (10): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+2 more)

### Community 271 - "New"
Cohesion: 0.49
Nodes (11): New(), githubSignature(), newStubStore(), post(), TestBranchAndPathMustMatchTheTrigger(), TestGitHubBadOrMissingSignatureIsRejected(), TestGitHubPushStartsRun(), TestGitLabPushRequiresTheToken() (+3 more)

### Community 273 - ".server"
Cohesion: 0.33
Nodes (7): fakeHooks, TestAtomicWebhookInvalidTriggerChangesNothing(), TestAtomicWebhookSaveRegistersFinalTriggers(), TestCreateWebhookEventRegistersProvider(), TestDeleteWebhookRemovesProviderHook(), TestWebhookEventTemplateOutsideOrganizationIsRejected(), webhookWriteStub()

### Community 275 - "newScheduleAPI"
Cohesion: 0.42
Nodes (9): newScheduleAPI(), newScheduleStub(), scheduleBody(), scheduleStr(), TestCreateScheduleCopiesAttributesAndDefaultsEnabled(), TestCreateScheduleRejectsForeignTemplate(), TestDeleteScheduleAndCrossWorkspaceLookup(), TestListAndUpdateSchedule() (+1 more)

### Community 276 - "Attr"
Cohesion: 0.43
Nodes (5): exists(), likePattern(), relate(), Attr, Kind

### Community 277 - "authzAPI"
Cohesion: 0.33
Nodes (10): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), strPtr(), TestResolveOrganization() (+2 more)

### Community 280 - "logSink"
Cohesion: 0.20
Nodes (3): captureLogPublisher, LogPublisher, logSink

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 591 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **146 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `Team`, `store.Reference`, `registryStub`, `run`, `Dispatcher`, `newFixture`, `time.Time`, `fakeQuerier`, `policyStub`, `Identity`, `.CreateVCS`, `store.Action`, `.marshalOBoolean2bool`, `newTestAPI`, `New`, `tfe_test.go`, `.CreateArchive`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `ModuleVersion`, `WebhookEvent`, `TestWebhookQueries`, `FailStuckRunTriggerEventsParams`, `nestedStub`, `api/main.go`, `stubQuerier`, `.CreateNotificationConfiguration`, `store.Item`, `notificationSendAPI`, `Metrics`, `fakeStore`, `History`, `.CreateOrganization`, `CliAuthSession`, `stubQuerier`, `store.Step`, `Pat`, `Queries`, `store.Collection`, `tagStub`, `github.com/aminbenmansour/terrakube-go/internal/store.Job`, `Workspace`, `dispatcherStore`, `decodeJSON`, `Cipher`, `context.Context`, `TeamToken`, `Claims`, `runTriggerAPI`, `store.Job`, `WorkspaceRunTrigger`, `.QueryRow`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07346491228070176 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.Request` to `WriteError`, `.listVCSRepositories`, `API`, `Identity`, `api/main.go`, `Webhook`, `.atomicOperations`, `Link`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `WriteError` be split into smaller, more focused modules?**
  _Cohesion score 0.06516056758775206 - nodes in this community are weakly interconnected._
- **Why does `New()` connect `net/http.Request` to `WriteError`, `runTriggerAPI`, `NewLayout`, `API`, `notificationSendAPI`, `Identity`, `decodeJSON`, `newTestAPI`, `.server`, `newScheduleAPI`, `testing.T`, `authzAPI`, `api/main.go`, `Webhook`, `stubQuerier`, `newWorkspaceStub`, `.atomicOperations`, `Link`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._