# Graph Report - terrakube-go  (2026-11-27)

## Corpus Check
- 528 files · ~279,807 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4883 nodes · 20187 edges · 263 communities (120 shown, 143 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3622 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3597c850`
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
- NewCipher
- net/http.ResponseWriter
- newFixture
- WriteError
- fakeQuerier
- Service
- steps
- stubQuerier
- deliver_test.go
- tfe_test.go
- sqlc Generation Config
- serveAccessToken
- Upstream Elide entities (GraphQL surface)
- Step
- authzAPI
- payload_test.go
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- api/main.go
- Link
- vcs/vcs.go
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
- Content
- Tag
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
- Cipher
- generated.go
- where_test.go
- streamEvents
- .prepareWorkspace
- decodeJSON
- graphql/models.go
- newJWKSServer
- Migrate
- Secret encryption at rest — design
- NewPool
- REST contract, captured from upstream
- pat.py
- Layout
- stubQuerier
- Workspacetag
- Agent
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
- JobConnection
- Item
- compare.sh
- start.sh
- ._NotificationConfiguration
- newTestExecutor
- time.Time
- WebhookEvent
- Ssh
- rsql.Node
- run
- clilogin/handler_test.go
- newTestAPI
- Secret key rotation — design
- OrganizationConnection
- ui.go
- API
- Vc
- UpdateNotificationConfigurationParams
- NewSubscriptions
- fakeStore
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- Executor
- grantsAPI
- History
- tarGzEntries
- Dex
- github.com/aminbenmansour/terrakube-go/internal/store.Step
- net/http.Handler
- notificationSendAPI
- AccessConnection
- Collection
- Metrics
- transientError
- log/slog.Logger
- Threat model
- .marshalNString2string
- auth.DexValidator
- store.Workspace
- loadPage
- RSQL filters in SQL — design
- ResourceIdentifier
- TemplateConnection
- Workspace
- dispatcherWith
- ParsePage
- Schedule
- Review Focus
- .createVCS
- Team
- .build
- .Query
- Review Focus
- where.go
- Deps
- NotificationConfigurationEdge
- .CreateArchive
- .SetCliAuthSessionIdentity

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Queries` - 244 edges
3. `API` - 225 edges
4. `New()` - 176 edges
5. `stubQuerier` - 171 edges
6. `WriteError()` - 154 edges
7. `param()` - 90 edges
8. `ParseFieldsets()` - 78 edges
9. `WriteResource()` - 69 edges
10. `API` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Command` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Report and verification` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Task 5: Backfill implementation` --references--> `needsEncryption()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/secretstore/backfill.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (263 total, 143 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (31): configurationVersionResource(), projectResource(), parseWorkspaceQuery(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource() (+23 more)

### Community 1 - "rbac.go"
Cohesion: 0.12
Nodes (36): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, ApplyProjectGrants(), ApplyWorkspaceGrants() (+28 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (14): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+6 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 7 - "NewCipher"
Cohesion: 0.16
Nodes (30): TestBackfillEncryptsPlaintextRows(), NewCipher(), TestCipherEmptyAndPlaintext(), TestCipherKeyring(), TestCipherRejectsBadConfiguration(), TestCipherRejectsWrongColumnAndKey(), TestCipherRoundTrip(), testKey() (+22 more)

### Community 8 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (20): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), TestClassify(), TestExtractToken(), WithClaims() (+12 more)

### Community 9 - "newFixture"
Cohesion: 0.09
Nodes (35): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+27 more)

### Community 10 - "WriteError"
Cohesion: 0.09
Nodes (33): actionAttrs, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, decodeActionAttrs(), stringOr(), boolValue() (+25 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.14
Nodes (3): encryptedStore, fakeQuerier, CreateNotificationConfigurationParams

### Community 13 - "Service"
Cohesion: 0.14
Nodes (12): Exchanger, Service, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value() (+4 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 17 - "stubQuerier"
Cohesion: 0.03
Nodes (11): listQuery, stubLists, stubQuerier, Task 3: Handler cutover, scopeIDs(), demoOrganizations(), ListSpec, NotificationConfiguration (+3 more)

### Community 18 - "deliver_test.go"
Cohesion: 0.17
Nodes (18): deliveryError(), stringPtr(), TestDeliverDoesNotFollowRedirects(), testDeliverer(), TestDeliverIgnoresMalformedRetryAfter(), TestDeliverNetworkErrorIsRetryableAndHidesTheURL(), TestDeliverParsesHTTPDateRetryAfter(), TestDeliverPostsPayload() (+10 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (41): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue() (+33 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "serveAccessToken"
Cohesion: 0.12
Nodes (24): expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups(), TestAccessTokenRejectsBadInput() (+16 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "Step"
Cohesion: 0.08
Nodes (7): captureSteps(), dispatcherStore, FailStuckNotificationOutboxParams, GetStepByJobAndNumberParams, ReclaimStuckNotificationOutboxParams, RecordNotificationOutboxParams, Step

### Community 24 - "authzAPI"
Cohesion: 0.13
Nodes (19): discoveryStub, fakeHooks, missingOrganizationStub, authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound() (+11 more)

### Community 25 - "payload_test.go"
Cohesion: 0.19
Nodes (26): BuildPayload(), buildSlack(), buildTeams(), buildWebhook(), capitalize(), marshalPayload(), notBlank(), nullIfBlank() (+18 more)

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
Nodes (30): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), accessResource(), projectAccessResource() (+22 more)

### Community 32 - "vcs/vcs.go"
Cohesion: 0.05
Nodes (59): Connection shape, Entity and relationship surface, Filter expressions the UI generates, Other contracts observed, Sort strings, Upstream UI GraphQL operations, Measured route surface, Order that unblocks the most UI (+51 more)

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
Cohesion: 0.17
Nodes (7): notificationConfigurationResolver, organizationResolver, queryResolver, workspaceResolver, organizationScope(), TestOrganizationScope(), idsFrom()

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

### Community 44 - "Content"
Cohesion: 0.25
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 47 - ".streamStepLog"
Cohesion: 0.15
Nodes (15): compareStreamIDs(), isTerminal(), outputOf(), parseResumeID(), splitStreamID(), stepOf(), streamKey(), TestCompareStreamIDs() (+7 more)

### Community 48 - "CliAuthSession"
Cohesion: 0.13
Nodes (5): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams

### Community 51 - "time.Duration"
Cohesion: 0.08
Nodes (7): captureLogPublisher, fakeQueue, JobQueue, LogPublisher, defaultRetryDelay(), backoff(), Message

### Community 52 - "testing.T"
Cohesion: 0.04
Nodes (88): TestAccessTokenWithoutStore(), TestAuthMiddleware(), TestCreateOrganization(), newCollectionAPI(), TestCollectionListCarriesItsItemsAndReferences(), TestCreateCollectionDefaultsPriorityToZero(), TestCreateCollectionItemStoresTheValue(), TestCreateCollectionItemValidatesKeyValueAndCategory() (+80 more)

### Community 53 - "Pat"
Cohesion: 0.16
Nodes (6): newFakeStore(), fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "models_gen.go"
Cohesion: 0.06
Nodes (11): ComplexityRoot, JobConnection, NotificationConfigurationConnection, PageInfo, ProjectConnection, ProjectEdge, Query, Subscription (+3 more)

### Community 57 - "newWorkspaceStub"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 59 - ".QueryRow"
Cohesion: 0.03
Nodes (8): ApproveJobParams, CreateProjectParams, GetProjectByNameParams, NotificationOutboxExistsForWorkspaceParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams, UpdateOrganizationParams, UpdateProjectParams

### Community 61 - "Queries"
Cohesion: 0.03
Nodes (8): isActiveStepStatus(), ValidateActiveStep(), AddNotificationConfigurationTemplateParams, DBTX, Queries, UpdateHistoryOutputParams, UpdateJobOutputParams, UpdateStepOutputParams

### Community 62 - ".createVCS"
Cohesion: 0.50
Nodes (3): discoveryConnection(), discoveryPage(), writeRawJSON()

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.06
Nodes (55): redactor, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+47 more)

### Community 67 - "Action"
Cohesion: 0.13
Nodes (4): actionStub, CreateActionParams, UpdateActionParams, Action

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "Load"
Cohesion: 0.15
Nodes (24): Config, S3, Storage, Task 2: Serve archived logs from `/tfoutput/v1`, Task 1: Keyring configuration, env(), envBool(), envDuration() (+16 more)

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - "New"
Cohesion: 0.20
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 73 - "Cipher"
Cohesion: 0.11
Nodes (26): Task 3: Rotation scan and re-encryption, backfillGlobalVariables(), backfillNotifications(), backfillSSHKeys(), backfillVariables(), backfillVCS(), encryptIfNeeded(), needsEncryption() (+18 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "streamEvents"
Cohesion: 0.10
Nodes (9): projectResolver, Resolver, ResolverRoot, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace() (+1 more)

### Community 77 - ".prepareWorkspace"
Cohesion: 0.15
Nodes (11): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), transient(), contentIDFromSource(), gitEnv(), gitSource(), isConfigurationTarball() (+3 more)

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
Cohesion: 0.13
Nodes (16): Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation (+8 more)

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - "NewPool"
Cohesion: 0.15
Nodes (17): main(), run(), main(), run(), formatOldKeys(), logReport(), main(), run() (+9 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "Layout"
Cohesion: 0.12
Nodes (10): Validator, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph (+2 more)

### Community 89 - "stubQuerier"
Cohesion: 0.16
Nodes (20): replay, upstreamCapture, assertSameDocumentShape(), assertSameResourceShape(), captureActionStore(), captureCollectionStore(), captureGlobalVarStore(), captureJobStore() (+12 more)

### Community 94 - "Identity"
Cohesion: 0.10
Nodes (15): groupVerifier, PATService, stubDexVerifier, stubIdentityVerifier, stubPATService, fakeVerifier, Identity, fakeIssuer (+7 more)

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
Cohesion: 0.03
Nodes (14): collectionStub, globalVarStub, sshAgentStub, teamStub, templateStub, vcsStub, workspaceStub, artifactStore (+6 more)

### Community 108 - "newService"
Cohesion: 0.21
Nodes (14): decodeSecret(), New(), dexUser(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp(), TestDeleteEnforcesOwnership() (+6 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.14
Nodes (21): runTriggerStub, newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit(), TestCreateRunTriggerRejectsCycle(), TestCreateRunTriggerRejectsStructuralViolations() (+13 more)

### Community 121 - "Dispatcher"
Cohesion: 0.23
Nodes (4): NewDispatcher(), strPtr(), Dispatcher, NotificationOutbox

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 128 - "Item"
Cohesion: 0.16
Nodes (3): CreateItemParams, UpdateItemParams, Item

### Community 135 - "newTestExecutor"
Cohesion: 0.31
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 137 - "time.Time"
Cohesion: 0.20
Nodes (8): bucket, limiter, newLimiter(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), backoff(), TestBackoffGrowsAndCaps(), Dispatcher

### Community 142 - "rsql.Node"
Cohesion: 0.09
Nodes (21): Parse(), dump(), TestParse(), TestParseRejects(), exists(), likePattern(), relate(), ptr() (+13 more)

### Community 143 - "run"
Cohesion: 0.06
Nodes (26): main(), run(), cancelRunStore, NewMetrics(), TestMetricsDisabledWithoutToken(), TestMetricsRecordsRequests(), TestMetricsRequiresToken(), TestSampleRefreshesGauges() (+18 more)

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.06
Nodes (49): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+41 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 149 - "ui.go"
Cohesion: 0.12
Nodes (18): requireMetricsToken(), key(), RateLimit(), TestRateLimitDisabledByDefault(), TestRateLimitMiddleware(), SecurityHeaders(), TestSecurityHeaders(), New() (+10 more)

### Community 150 - "API"
Cohesion: 0.09
Nodes (10): API, New(), param(), parseList(), unknownCollection(), writeInUseConflict(), tagResource(), teamResource() (+2 more)

### Community 154 - "NewSubscriptions"
Cohesion: 0.14
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 157 - "fakeStore"
Cohesion: 0.05
Nodes (28): captureJobs(), strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased() (+20 more)

### Community 159 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.06
Nodes (7): CreateReferenceParams, CreateNotificationTriggerParams, FailStuckRunTriggerEventsParams, ListNotificationDeliveriesByWorkspaceParams, NotificationTrigger, ReclaimStuckRunTriggerEventsParams, Reference

### Community 161 - "Executor"
Cohesion: 0.08
Nodes (13): Executor, logSink, preparedWorkspace, isCancelledStatus(), isTerminalStatus(), isTransient(), parseAttempt(), stringPointer() (+5 more)

### Community 162 - "grantsAPI"
Cohesion: 0.31
Nodes (10): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+2 more)

### Community 163 - "History"
Cohesion: 0.09
Nodes (6): reliabilityStore, CreateHistoryParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams, History, UpdateStepStatusParams

### Community 164 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 167 - "net/http.Handler"
Cohesion: 0.07
Nodes (27): authorization, main(), run(), Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom() (+19 more)

### Community 168 - "notificationSendAPI"
Cohesion: 0.14
Nodes (23): newDeliveryStub(), TestNotificationDeliveryHistory(), TestNotificationDeliveryHistoryRejectsABadLimit(), TestNotificationDeliveryHistoryUnknownWorkspace(), TestNotificationDeliveryRetry(), TestNotificationDeliveryRetryRejections(), captureJSONServer(), notificationSendAPI() (+15 more)

### Community 171 - "Collection"
Cohesion: 0.16
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 172 - "Metrics"
Cohesion: 0.10
Nodes (7): metricsSource, Deps, fakePinger, fakeStats, Metrics, Pinger, StatsSource

### Community 175 - "log/slog.Logger"
Cohesion: 0.11
Nodes (18): startScheduler(), waitFor(), NewDispatcher(), NewDispatcher(), strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued() (+10 more)

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 184 - "auth.DexValidator"
Cohesion: 0.15
Nodes (8): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice(), channelLabel(), networkCause(), parseRetryAfter(), Configuration

### Community 190 - "store.Workspace"
Cohesion: 0.10
Nodes (12): atomicOperation, webhookEventAttrs, webhookSyncRef, strValue(), atomicWorkspaceHref(), isAtomicRequest(), isWebhookAtomicHref(), parseWebhookHref() (+4 more)

### Community 191 - "loadPage"
Cohesion: 0.24
Nodes (3): DB, failingDB, loadPage()

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 197 - "ResourceIdentifier"
Cohesion: 0.21
Nodes (10): nestedDocument, DecodeAttributes(), RelationshipIdentifiers(), RelationshipList(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship (+2 more)

### Community 202 - "dispatcherWith"
Cohesion: 0.36
Nodes (8): dispatcherWith(), newDispatcherFixture(), TestDispatcherFailsAfterMaxAttempts(), TestDispatcherFailsPermanentRejection(), TestDispatcherIncludesTheFailureSummary(), TestDispatcherRetriesTransientFailures(), TestDispatcherSendsAndRecordsSuccess(), TestDispatcherSkipsSupersededJobStatus()

### Community 204 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 229 - "Schedule"
Cohesion: 0.19
Nodes (4): scheduleStub, CreateScheduleParams, Schedule, UpdateScheduleParams

### Community 231 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 232 - ".createVCS"
Cohesion: 0.39
Nodes (6): vcsAttrs, strPtr(), decodeVCSAttrs(), requiredVCSString(), validVCSEnum(), vcsEnumOrDefault()

### Community 234 - "Team"
Cohesion: 0.04
Nodes (15): AccessTokenStore, Deps, fakeAccessTokenStore, contains(), captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces() (+7 more)

### Community 235 - ".build"
Cohesion: 0.10
Nodes (22): collectionSpec, listArgs, page, accessEdges(), afterString(), column(), jobEdges(), notificationConfigurationEdges() (+14 more)

### Community 236 - ".Query"
Cohesion: 0.07
Nodes (16): CountJobsByStatusRow, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow, ListSSHKeySecretsAfterParams (+8 more)

### Community 248 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 252 - "Deps"
Cohesion: 0.29
Nodes (3): Deps, fakeLocker, JobLocker

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 549 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **143 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `API` connect `API` to `net/http.Handler`, `.createVCS`, `WriteError`, `Team`, `api/main.go`, `.createVCS`, `store.Workspace`, `Link`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.06980240549828179 - nodes in this community are weakly interconnected._
- **Why does `run()` connect `net/http.Handler` to `net/http.ResponseWriter`, `NewCipher`, `net/http.ResponseWriter`, `newFixture`, `Service`, `run`, `clilogin/handler_test.go`, `newTestAPI`, `tfe_test.go`, `ui.go`, `API`, `NewSubscriptions`, `api/main.go`, `log/slog.Logger`, `NewLayout`, `Load`, `New`, `newJWKSServer`, `Migrate`, `NewPool`, `newService`, `Dispatcher`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `rbac.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1196808510638298 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `github.com/99designs/gqlgen/graphql.CollectedField` to `._Project`, `._NotificationConfiguration`, `TemplateConnection`, `context.Context`, `AccessConnection`, `generated.go`, `.childFields_JobStatusEvent`, `streamEvents`, `.marshalOBoolean2bool`, `github.com/99designs/gqlgen/graphql.Marshaler`, `OrganizationConnection`, `models_gen.go`, `.marshalNString2string`, `NotificationConfigurationEdge`, `JobConnection`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._