# Graph Report - terrakube-go  (2026-11-27)

## Corpus Check
- 499 files · ~256,066 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4523 nodes · 18309 edges · 235 communities (116 shown, 119 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3346 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2c95f004`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- StringOrInt
- io.Reader
- Core Database Schema
- run
- Handler
- executionContext
- .internalError
- fakeQuerier
- Service
- Service
- steps
- github.com/99designs/gqlgen/graphql.Marshaler
- github.com/aminbenmansour/terrakube-go/internal/rsql.Node
- WriteError
- tfe_test.go
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- fakeStore
- newFixture
- payload_test.go
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- api/main.go
- Link
- Upstream UI GraphQL operations
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
- graphql/collections.go
- JobConnection
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- New
- .applyWorkspaceUpdate
- NewCipher
- Job
- Queries
- History
- stubQuerier
- dex-config.yaml
- NewLayout
- NewPool
- models
- config.go
- Team
- collections.go
- NotificationConfigurationConnection
- Cipher
- generated.go
- where_test.go
- streamEvents
- .prepareWorkspace
- decodeJSON
- graphql/models.go
- newJWKSServer
- logSink
- Secret encryption at rest — design
- ui.go
- REST contract, captured from upstream
- pat.py
- log/slog.Logger
- Parse
- newTestExecutor
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
- newWorkspaceStub
- capture.sh
- 0008_action.sql
- time.Duration
- compare.sh
- start.sh
- .CreateGlobalVariable
- TestRotateReencryptsOldKeyAndPlaintext
- time.Time
- tarGzEntries
- .processJob
- rsql.Node
- where.go
- clilogin/handler_test.go
- newTestAPI
- Secret key rotation — design
- AccessConnection
- .handle
- API
- Sign
- NewSubscriptions
- stubQuerier
- Executor
- .QueryRow
- auth.DexValidator
- notificationTestAPI
- authzAPI
- Metrics
- transientError
- Threat model
- github.com/vektah/gqlparser/v2/ast.SelectionSet
- loadPage
- RSQL filters in SQL — design
- Workspace
- New
- log/slog.Logger
- Review Focus
- Migrate
- .build

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Queries` - 227 edges
3. `Querier` - 225 edges
4. `API` - 182 edges
5. `New()` - 159 edges
6. `stubQuerier` - 154 edges
7. `WriteError()` - 130 edges
8. `param()` - 82 edges
9. `ParseFieldsets()` - 72 edges
10. `API` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Self-review` --references--> `SortKey`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/store/list.go
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

## Communities (235 total, 119 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (28): configurationVersionResource(), projectResource(), parseWorkspaceQuery(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource() (+20 more)

### Community 1 - "rbac.go"
Cohesion: 0.14
Nodes (36): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, ApplyProjectGrants(), ApplyWorkspaceGrants() (+28 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "StringOrInt"
Cohesion: 0.10
Nodes (4): ComplexityRoot, ProjectConnection, ProjectEdge, WorkspaceConnection

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (15): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+7 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 7 - "run"
Cohesion: 0.14
Nodes (26): main(), run(), main(), run(), main(), run(), main(), run() (+18 more)

### Community 8 - "Handler"
Cohesion: 0.08
Nodes (27): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), WithClaims(), isTerminal(), outputOf() (+19 more)

### Community 9 - "executionContext"
Cohesion: 0.06
Nodes (3): executionContext, ProjectAccessConnection, ProjectAccessEdge

### Community 10 - ".internalError"
Cohesion: 0.10
Nodes (7): API, New(), param(), parseList(), unknownCollection(), withAuthorization(), WriteCollection()

### Community 11 - "fakeQuerier"
Cohesion: 0.03
Nodes (13): nestedStub, Task 3: Transparent store decorator, encryptedStore, fakeQuerier, CreateNotificationConfigurationParams, CreateSSHKeyParams, CreateVCSParams, Ssh (+5 more)

### Community 12 - "Service"
Cohesion: 0.18
Nodes (6): PATService, stubPATService, stringPtr(), timePtr(), Service, Token

### Community 13 - "Service"
Cohesion: 0.14
Nodes (13): Exchanger, Service, Store, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode() (+5 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 17 - "github.com/aminbenmansour/terrakube-go/internal/rsql.Node"
Cohesion: 0.09
Nodes (9): listQuery, stubLists, Task 3: Handler cutover, scopeIDs(), Lister, Organization, Project, SortKey (+1 more)

### Community 18 - "WriteError"
Cohesion: 0.12
Nodes (30): actionAttrs, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, decodeActionAttrs(), stringOr(), boolValue() (+22 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (42): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+34 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.07
Nodes (49): replay, upstreamCapture, expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner() (+41 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "fakeStore"
Cohesion: 0.04
Nodes (21): scheduleStub, stubQuerier, fakeStore, fakeStore, CreateJobParams, CreateScheduleParams, CreateStepParams, FailStuckRunTriggerEventsParams (+13 more)

### Community 24 - "newFixture"
Cohesion: 0.15
Nodes (23): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+15 more)

### Community 25 - "payload_test.go"
Cohesion: 0.19
Nodes (27): BuildPayload(), buildSlack(), buildTeams(), buildWebhook(), capitalize(), marshalPayload(), notBlank(), nullIfBlank() (+19 more)

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
Nodes (61): actionFilter, agentAttrs, atomicOperation, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs (+53 more)

### Community 31 - "Link"
Cohesion: 0.07
Nodes (38): vcsAttrs, actionResource(), parseInclude(), strPtr(), strValue(), collectionResource(), itemResource(), referenceResource() (+30 more)

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

### Community 37 - "idsFrom"
Cohesion: 0.16
Nodes (8): notificationConfigurationResolver, organizationResolver, queryResolver, workspaceResolver, organizationScope(), TestOrganizationScope(), idsFrom(), workspaceEdges()

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

### Community 45 - "addDestination"
Cohesion: 0.33
Nodes (15): strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased(), TestTickIgnoresPlanOnlyRun() (+7 more)

### Community 46 - "graphql/collections.go"
Cohesion: 0.16
Nodes (13): collectionSpec, page, accessEdges(), column(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), predicate() (+5 more)

### Community 51 - "time.Duration"
Cohesion: 0.10
Nodes (4): fakeQueue, JobQueue, Message, JobQueue

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (97): TestAccessTokenWithoutStore(), TestAuthMiddleware(), TestCreateOrganization(), newCollectionAPI(), TestCollectionListCarriesItsItemsAndReferences(), TestCreateCollectionDefaultsPriorityToZero(), TestCreateCollectionItemStoresTheValue(), TestCreateCollectionItemValidatesKeyValueAndCategory() (+89 more)

### Community 53 - "Pat"
Cohesion: 0.16
Nodes (6): newFakeStore(), fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "New"
Cohesion: 0.21
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 57 - ".applyWorkspaceUpdate"
Cohesion: 0.09
Nodes (21): nestedDocument, atomicWorkspaceHref(), isAtomicRequest(), pickStringValue(), workspaceLinks(), workspaceUpdateParams(), DecodeAttributes(), RelationshipID() (+13 more)

### Community 59 - "NewCipher"
Cohesion: 0.23
Nodes (21): TestBackfillEncryptsPlaintextRows(), NewCipher(), TestCipherEmptyAndPlaintext(), TestCipherKeyring(), TestCipherRejectsBadConfiguration(), TestCipherRejectsWrongColumnAndKey(), TestCipherRoundTrip(), testKey() (+13 more)

### Community 60 - "Job"
Cohesion: 0.11
Nodes (5): ApproveJobParams, GetJobByWorkspaceAndTriggerParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams, UpdateJobStatusParams

### Community 61 - "Queries"
Cohesion: 0.02
Nodes (34): metricsSource, UpdateItemParams, AddNotificationConfigurationTemplateParams, Agent, CountJobsByStatusRow, CreateAgentParams, CreateArchiveParams, CreateNotificationTriggerParams (+26 more)

### Community 62 - "History"
Cohesion: 0.10
Nodes (6): reliabilityStore, CreateHistoryParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams, History, UpdateStepStatusParams

### Community 63 - "stubQuerier"
Cohesion: 0.03
Nodes (15): jobStub, tagStub, statusQuerier, captureJobs(), captureSteps(), demoOrganizations(), CreateContentParams, UpdateContentStatusParams (+7 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.06
Nodes (57): redactor, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+49 more)

### Community 67 - "NewPool"
Cohesion: 0.17
Nodes (10): Locker, fakeLocker, NewLocker(), TestLockerAllowsDifferentJobs(), TestLockerExcludesASecondHolder(), testPool(), appliedVersions(), applyMigration() (+2 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "config.go"
Cohesion: 0.15
Nodes (24): Config, S3, Storage, Task 2: Serve archived logs from `/tfoutput/v1`, Task 1: Keyring configuration, env(), envBool(), envDuration() (+16 more)

### Community 70 - "Team"
Cohesion: 0.04
Nodes (13): AccessTokenStore, Deps, fakeAccessTokenStore, contains(), contains(), fakeOrganizationStore, fakeWorkspaceStore, Access (+5 more)

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 73 - "Cipher"
Cohesion: 0.21
Nodes (21): Task 3: Rotation scan and re-encryption, backfillGlobalVariables(), backfillNotifications(), backfillSSHKeys(), backfillVariables(), backfillVCS(), encryptIfNeeded(), needsEncryption() (+13 more)

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "streamEvents"
Cohesion: 0.09
Nodes (9): projectResolver, Resolver, ResolverRoot, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace() (+1 more)

### Community 77 - ".prepareWorkspace"
Cohesion: 0.13
Nodes (12): gitConfig, preparedWorkspace, TestGitSourceDetection(), TestSanitizeGitURL(), transient(), contentIDFromSource(), gitEnv(), gitSource() (+4 more)

### Community 78 - "decodeJSON"
Cohesion: 0.12
Nodes (31): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent() (+23 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "logSink"
Cohesion: 0.17
Nodes (4): captureLogPublisher, LogPublisher, logSink, baseEnvironment()

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "log/slog.Logger"
Cohesion: 0.11
Nodes (15): Validator, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph (+7 more)

### Community 89 - "Parse"
Cohesion: 0.19
Nodes (14): Parse(), dump(), TestParse(), TestParseRejects(), ptr(), seedWorkspaces(), TestListPageBetween(), TestListPageFilters() (+6 more)

### Community 91 - "newTestExecutor"
Cohesion: 0.16
Nodes (19): newFakeLocker(), TestHandleReleasesTheJobLock(), TestHandleSkipsAJobHeldByAnotherExecutor(), TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestStaleClaimIdleOutlivesTheJobTimeout(), New(), TestAckDeletesTheJobMessage() (+11 more)

### Community 94 - "Identity"
Cohesion: 0.12
Nodes (10): groupVerifier, stubDexVerifier, stubIdentityVerifier, fakeVerifier, Identity, fakeDex, fakeIssuer, groupVerifier (+2 more)

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 99 - "cron.go"
Cohesion: 0.24
Nodes (8): at(), Parse(), parseField(), splitStep(), mustParse(), TestParseMatchesQuartzFields(), TestParseRejectsWhatQuartzRejects(), Cron

### Community 100 - "NewSubscriptions"
Cohesion: 0.23
Nodes (6): EventSource, fakeEvents, Event, NewBus(), TestBusRoundTrip(), Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (20): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+12 more)

### Community 103 - "context.Context"
Cohesion: 0.03
Nodes (20): collectionStub, globalVarStub, sshAgentStub, teamStub, templateStub, vcsStub, workspaceStub, artifactStore (+12 more)

### Community 108 - "newService"
Cohesion: 0.21
Nodes (14): decodeSecret(), New(), dexUser(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp(), TestDeleteEnforcesOwnership() (+6 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.16
Nodes (22): runTriggerStub, newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), runTriggerRBACAPI(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit(), TestCreateRunTriggerRejectsCycle() (+14 more)

### Community 121 - "newWorkspaceStub"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 128 - "time.Duration"
Cohesion: 0.19
Nodes (6): waitFor(), defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), backoff(), TestBackoffGrowsAndCaps()

### Community 133 - ".CreateGlobalVariable"
Cohesion: 0.20
Nodes (3): CreateGlobalVariableParams, Globalvar, UpdateGlobalVariableParams

### Community 135 - "TestRotateReencryptsOldKeyAndPlaintext"
Cohesion: 0.16
Nodes (16): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs, assertSealedUnder() (+8 more)

### Community 137 - "time.Time"
Cohesion: 0.23
Nodes (6): bucket, limiter, newLimiter(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), Dispatcher

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 142 - "rsql.Node"
Cohesion: 0.13
Nodes (9): exists(), likePattern(), relate(), And, Or, parser, Attr, Kind (+1 more)

### Community 143 - "where.go"
Cohesion: 0.42
Nodes (3): likePattern(), builder, Compare

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.07
Nodes (42): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+34 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 150 - "API"
Cohesion: 0.08
Nodes (19): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+11 more)

### Community 151 - "Sign"
Cohesion: 0.67
Nodes (3): Sign(), TestSignDifferentSecretsDiffer(), TestSignIsStableHexHMAC()

### Community 154 - "NewSubscriptions"
Cohesion: 0.14
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 159 - "stubQuerier"
Cohesion: 0.03
Nodes (11): actionStub, stubQuerier, CreateActionParams, UpdateActionParams, CreateCollectionParams, UpdateCollectionParams, Action, Collection (+3 more)

### Community 161 - "Executor"
Cohesion: 0.20
Nodes (5): Executor, isCancelledStatus(), isTerminalStatus(), TestBinarySelection(), TestIsTerminalStatus()

### Community 163 - ".QueryRow"
Cohesion: 0.05
Nodes (12): fakeStore, CliAuthSession, CountEnabledRunTriggersBySourceParams, CreateCliAuthSessionParams, CreateProjectParams, CreateTagParams, GetProjectByNameParams, GetTagByNameParams (+4 more)

### Community 165 - "auth.DexValidator"
Cohesion: 0.15
Nodes (5): Dex, contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 171 - "authzAPI"
Cohesion: 0.17
Nodes (13): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), strPtr(), TestResolveOrganization() (+5 more)

### Community 172 - "Metrics"
Cohesion: 0.17
Nodes (3): fakeStats, Metrics, StatsSource

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - "github.com/vektah/gqlparser/v2/ast.SelectionSet"
Cohesion: 0.06
Nodes (11): NotificationTriggerConnection, NotificationTriggerEdge, OrganizationConnection, OrganizationEdge, PageInfo, Query, Subscription, TemplateConnection (+3 more)

### Community 190 - "loadPage"
Cohesion: 0.27
Nodes (3): DB, failingDB, loadPage()

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 199 - "Workspace"
Cohesion: 0.11
Nodes (3): CreateWorkspaceParams, GetWorkspaceByNameParams, SetWorkspaceLockParams

### Community 202 - "New"
Cohesion: 0.12
Nodes (13): Deps, fakePinger, Pinger, requireMetricsToken(), key(), RateLimit(), TestRateLimitDisabledByDefault(), TestRateLimitMiddleware() (+5 more)

### Community 229 - "log/slog.Logger"
Cohesion: 0.13
Nodes (16): startScheduler(), NewDispatcher(), NewDispatcher(), strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued(), TestTickSkipsLockedWorkspaceAndUnparseableCron() (+8 more)

### Community 231 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 234 - "Migrate"
Cohesion: 0.20
Nodes (11): Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope, Task 4: Backfill queries, Task 5: Backfill implementation, Task 6: Wiring and round-trip test (+3 more)

### Community 235 - ".build"
Cohesion: 0.16
Nodes (12): listArgs, afterString(), ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset(), Order, OrderBy() (+4 more)

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 512 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **119 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Queries` to `.CreateGlobalVariable`, `Handler`, `time.Time`, `fakeQuerier`, `newTestAPI`, `tfe_test.go`, `fakeStore`, `api/main.go`, `stubQuerier`, `.QueryRow`, `Pat`, `NewCipher`, `Job`, `History`, `stubQuerier`, `NewLayout`, `Team`, `Workspace`, `log/slog.Logger`, `log/slog.Logger`, `context.Context`, `Migrate`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07366735300846489 - nodes in this community are weakly interconnected._
- **Why does `API` connect `.internalError` to `Team`, `WriteError`, `API`, `.applyWorkspaceUpdate`, `api/main.go`, `Link`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `rbac.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13530655391120508 - nodes in this community are weakly interconnected._
- **Why does `openRotationTest()` connect `TestRotateReencryptsOldKeyAndPlaintext` to `run`, `context.Context`, `newTestAPI`, `testing.T`, `Queries`, `api/main.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._