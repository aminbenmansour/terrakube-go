# Graph Report - terrakube-go  (2026-11-27)

## Corpus Check
- 542 files · ~289,590 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 5040 nodes · 20976 edges · 275 communities (133 shown, 142 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3698 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `12618334`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- WriteError
- Agent Workflows and Compatibility Docs
- ._Project
- io.Reader
- Core Database Schema
- nestedStub
- run
- net/http.ResponseWriter
- newFixture
- New
- fakeQuerier
- M4 run logs and artifacts Implementation Plan
- steps
- github.com/99designs/gqlgen/graphql.Marshaler
- github.com/aminbenmansour/terrakube-go/internal/rsql.Node
- rbac.go
- tfe_test.go
- sqlc Generation Config
- serveAccessToken
- Upstream Elide entities (GraphQL surface)
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- logSink
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
- TestRotateReencryptsOldKeyAndPlaintext
- Queries
- deliver_test.go
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
- time.Time
- generated.go
- where_test.go
- streamEvents
- vcs/vcs.go
- decodeJSON
- graphql/models.go
- auth/dex_test.go
- .CreateProject
- Secret encryption at rest — design
- Cipher
- REST contract, captured from upstream
- pat.py
- Layout
- stubQuerier
- getJSON
- .CreateTag
- .Verify
- environment
- cron.go
- NewSubscriptions
- newTestService
- context.Context
- newService
- ui-rest-inventory.sh
- ui-smoke.sh
- runTriggerAPI
- Claims
- capture.sh
- 0008_action.sql
- WorkspaceTagConnection
- Service
- compare.sh
- start.sh
- TemplateConnection
- .EnqueueJob
- time.Time
- serveAccessToken
- Identity
- Parse
- ._NotificationConfiguration
- clilogin/handler_test.go
- newTestAPI
- Rotate
- failed
- ui.go
- API
- WebhookEvent
- NewSubscriptions
- Job
- stubQuerier
- Executor
- runTriggerAPI
- History
- tarGzEntries
- webhook_test.go
- Organization
- API
- notificationSendAPI
- AccessConnection
- github_test.go
- Metrics
- transientError
- Dispatcher
- Threat model
- .cloneConfiguration
- Dex
- store.Workspace
- loadPage
- RSQL filters in SQL — design
- Upstream UI GraphQL operations
- Review Focus
- Workspace
- .build
- Repository
- TeamToken
- Review Focus
- Team
- .build
- Queries
- TestWebhookQueries
- workspaceQuery
- .CreateArchive
- JobConnection
- newService
- Attr
- grantsAPI
- Content
- discoveryStub
- .server
- where.go
- Deps

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 256 edges
3. `Queries` - 252 edges
4. `API` - 230 edges
5. `New()` - 175 edges
6. `stubQuerier` - 171 edges
7. `WriteError()` - 158 edges
8. `param()` - 90 edges
9. `ParseFieldsets()` - 78 edges
10. `WriteResource()` - 69 edges

## Surprising Connections (you probably didn't know these)
- `Order that unblocks the most UI` --references--> `provider`  [INFERRED]
  docs/upstream/ui-rest-operations.md → internal/vcs/vcs.go
- `Where we stand` --references--> `provider`  [INFERRED]
  docs/upstream/ui-rest-operations.md → internal/vcs/vcs.go
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Task 5: Backfill implementation` --references--> `needsEncryption()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/secretstore/backfill.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (275 total, 142 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (32): Write(), WriteResource(), configurationVersionResource(), projectResource(), parseWorkspaceQuery(), activeRunStatus(), approvalGateStatus(), failedJobStatus() (+24 more)

### Community 1 - "WriteError"
Cohesion: 0.08
Nodes (37): actionAttrs, itemAttrs, nestedDocument, notificationConfigurationAttrs, teamAttrs, templateAttrs, decodeActionAttrs(), stringOr() (+29 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "._Project"
Cohesion: 0.13
Nodes (3): ProjectAccessConnection, ProjectAccessEdge, ProjectEdge

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (15): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+7 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 7 - "run"
Cohesion: 0.05
Nodes (74): main(), run(), main(), run(), main(), run(), main(), run() (+66 more)

### Community 8 - "net/http.ResponseWriter"
Cohesion: 0.24
Nodes (6): isExternalUser(), matchesMD5(), New(), newUUID(), pathIDs(), Handler

### Community 9 - "newFixture"
Cohesion: 0.09
Nodes (35): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+27 more)

### Community 10 - "New"
Cohesion: 0.09
Nodes (36): globToRegex(), matchBranch(), matchPaths(), matchTrigger(), pathRegex(), branchFromRef(), deliveryID(), nonEmptyPtr() (+28 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.06
Nodes (7): Task 2: Cipher and envelope, Task 3: Transparent store decorator, encryptedStore, fakeQuerier, CreateNotificationConfigurationParams, UpdateNotificationConfigurationParams, Variable

### Community 12 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.25
Nodes (7): Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 17 - "github.com/aminbenmansour/terrakube-go/internal/rsql.Node"
Cohesion: 0.06
Nodes (6): listQuery, stubLists, Task 3: Handler cutover, ListSpec, Tag, WorkspaceRunTrigger

### Community 18 - "rbac.go"
Cohesion: 0.26
Nodes (22): Task 1: Role matrix and permission set, adminOnly(), CanApproveJob(), CanDestroyJob(), CanManageCollection(), CanManageJob(), CanManageModule(), CanManagePolicies() (+14 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.08
Nodes (43): Backend, assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL() (+35 more)

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
Cohesion: 0.11
Nodes (8): NewDispatcher(), strPtr(), Dispatcher, dispatcherStore, NotificationConfiguration, NotificationOutbox, RecordNotificationOutboxParams, WebhookDelivery

### Community 24 - "logSink"
Cohesion: 0.12
Nodes (9): captureLogPublisher, LogPublisher, logSink, redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues() (+1 more)

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
Cohesion: 0.07
Nodes (59): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, createTeamTokenRequest, discoveryGroup (+51 more)

### Community 31 - "Link"
Cohesion: 0.07
Nodes (32): vcsAttrs, actionResource(), parseInclude(), strPtr(), collectionResource(), itemResource(), referenceResource(), globalVarResource() (+24 more)

### Community 32 - "getJSON"
Cohesion: 0.21
Nodes (6): Branch, BranchPage, github, Group, provider, RepositoryPage

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

### Community 44 - "addDestination"
Cohesion: 0.37
Nodes (15): strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased(), TestTickIgnoresPlanOnlyRun() (+7 more)

### Community 45 - ".QueryRow"
Cohesion: 0.02
Nodes (11): CreateCollectionParams, CreateItemParams, CreateReferenceParams, UpdateCollectionParams, UpdateItemParams, ApproveJobParams, CreateNotificationTriggerParams, GetJobByWorkspaceAndTriggerParams (+3 more)

### Community 46 - "ResolveOrganization"
Cohesion: 0.13
Nodes (11): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), contains(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access (+3 more)

### Community 47 - ".streamStepLog"
Cohesion: 0.15
Nodes (15): compareStreamIDs(), isTerminal(), outputOf(), parseResumeID(), splitStreamID(), stepOf(), streamKey(), TestCompareStreamIDs() (+7 more)

### Community 48 - "CliAuthSession"
Cohesion: 0.11
Nodes (6): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams, SetCliAuthSessionIdentityParams

### Community 49 - "stubQuerier"
Cohesion: 0.08
Nodes (4): fakeWorkspaceStore, stubQuerier, workspaceWithTags, stubQuerier

### Community 51 - "time.Duration"
Cohesion: 0.07
Nodes (12): waitFor(), fakeQueue, JobQueue, defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), backoff(), blockedDelivery() (+4 more)

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (121): serveTeamTokens(), TestAccessTokenWithoutStore(), TestTeamTokenCreatePassesTheRequestAndReturnsTheToken(), TestTeamTokenCreateRejectsNonMembersAndBadInput(), TestTeamTokenDeleteMapsErrors(), TestTeamTokenListIsAccepted(), TestTeamTokenRoutesRefuseOversizedBodies(), TestTeamTokensAnswer503WithoutAService() (+113 more)

### Community 53 - "Pat"
Cohesion: 0.16
Nodes (6): newFakeStore(), fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "models_gen.go"
Cohesion: 0.05
Nodes (10): ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, OrganizationConnection, OrganizationEdge, PageInfo, ProjectConnection, Query (+2 more)

### Community 57 - "newWorkspaceStub"
Cohesion: 0.23
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 61 - "Queries"
Cohesion: 0.03
Nodes (10): AddNotificationConfigurationTemplateParams, FailStuckNotificationOutboxParams, FailStuckRunTriggerEventsParams, Querier, ReclaimStuckNotificationOutboxParams, ReclaimStuckRunTriggerEventsParams, RecordWebhookDeliveryParams, UpdateHistoryOutputParams (+2 more)

### Community 62 - "deliver_test.go"
Cohesion: 0.11
Nodes (29): channelLabel(), networkCause(), NewDeliverer(), parseRetryAfter(), deliveryError(), stringPtr(), TestDeliverDoesNotFollowRedirects(), testDeliverer() (+21 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.06
Nodes (59): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+51 more)

### Community 67 - "Action"
Cohesion: 0.17
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
Cohesion: 0.23
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 73 - "time.Time"
Cohesion: 0.14
Nodes (12): bucket, Deps, fakePinger, limiter, Pinger, key(), newLimiter(), RateLimit() (+4 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "streamEvents"
Cohesion: 0.10
Nodes (9): projectResolver, Resolver, ResolverRoot, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace() (+1 more)

### Community 77 - "vcs/vcs.go"
Cohesion: 0.29
Nodes (14): TestGitLabConnectionWithoutAToken(), TestGitLabDiscovery(), TestProviderForRefusesALoopbackAPIURL(), baseURL(), Branches(), Test(), Groups(), providerError() (+6 more)

### Community 78 - "decodeJSON"
Cohesion: 0.12
Nodes (31): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent() (+23 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "auth/dex_test.go"
Cohesion: 0.08
Nodes (22): jwksServer, cancelRunStore, reliabilityStore, contains(), decodeJWKPart(), NewDexValidator(), rsaPublicKey(), stringSlice() (+14 more)

### Community 81 - ".CreateProject"
Cohesion: 0.20
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - "Cipher"
Cohesion: 0.09
Nodes (29): Task 3: Rotation scan and re-encryption, backfillGlobalVariables(), backfillNotifications(), backfillSSHKeys(), backfillVariables(), backfillVCS(), encryptIfNeeded(), needsEncryption() (+21 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "Layout"
Cohesion: 0.22
Nodes (6): Validator, Layout, Deps, JobQueue, Deps, Deps

### Community 89 - "stubQuerier"
Cohesion: 0.16
Nodes (20): replay, upstreamCapture, assertSameDocumentShape(), assertSameResourceShape(), captureActionStore(), captureCollectionStore(), captureGlobalVarStore(), captureJobStore() (+12 more)

### Community 91 - "getJSON"
Cohesion: 0.42
Nodes (3): doJSON(), getJSON(), gitlab

### Community 94 - ".Verify"
Cohesion: 0.20
Nodes (4): stubDexVerifier, fakeVerifier, stubDexVerifier, stringSlice()

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 99 - "cron.go"
Cohesion: 0.27
Nodes (8): at(), Parse(), parseField(), splitStep(), mustParse(), TestParseMatchesQuartzFields(), TestParseRejectsWhatQuartzRejects(), Cron

### Community 100 - "NewSubscriptions"
Cohesion: 0.33
Nodes (3): fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (20): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+12 more)

### Community 103 - "context.Context"
Cohesion: 0.02
Nodes (14): collectionStub, globalVarStub, sshAgentStub, teamStub, templateStub, vcsStub, workspaceStub, artifactStore (+6 more)

### Community 108 - "newService"
Cohesion: 0.21
Nodes (14): decodeSecret(), New(), dexUser(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp(), TestDeleteEnforcesOwnership() (+6 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.12
Nodes (13): scheduleStub, newScheduleAPI(), newScheduleStub(), scheduleBody(), scheduleStr(), TestCreateScheduleCopiesAttributesAndDefaultsEnabled(), TestCreateScheduleRejectsForeignTemplate(), TestDeleteScheduleAndCrossWorkspaceLookup() (+5 more)

### Community 121 - "Claims"
Cohesion: 0.13
Nodes (16): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken() (+8 more)

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 128 - "Service"
Cohesion: 0.14
Nodes (12): Exchanger, Service, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value() (+4 more)

### Community 135 - ".EnqueueJob"
Cohesion: 0.24
Nodes (5): fakeQueue, fakeQueue, JobQueue, stubJobQueue, stubQueue

### Community 137 - "time.Time"
Cohesion: 0.23
Nodes (7): startScheduler(), formatOldKeys(), logReport(), NewDispatcher(), NewDispatcher(), Dispatcher, JobQueue

### Community 140 - "serveAccessToken"
Cohesion: 0.21
Nodes (12): stubIdentityVerifier, expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups() (+4 more)

### Community 141 - "Identity"
Cohesion: 0.08
Nodes (18): groupVerifier, PATService, stubPATService, stubTeamTokens, TeamTokenService, Identity, fakeIssuer, groupVerifier (+10 more)

### Community 142 - "Parse"
Cohesion: 0.31
Nodes (3): And, Or, parser

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.07
Nodes (41): pageMeta, createListOrganization(), getNames(), seedOrganizations(), TestOrganizationListEscapesWildcards(), TestOrganizationListPrefiltersByAllowedIDs(), TestOrganizationListSortsAndPages(), TestOrganizationListSQLFilters() (+33 more)

### Community 146 - "Rotate"
Cohesion: 0.08
Nodes (28): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs, Approaches considered (+20 more)

### Community 147 - "failed"
Cohesion: 0.20
Nodes (8): blockedAddress(), TestBlockedAddressRanges(), TestValidateEndpointRefusesInternalAddresses(), validateEndpoint(), failed(), githubOwnerAndRepo(), Error, WebhookSpec

### Community 149 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 150 - "API"
Cohesion: 0.09
Nodes (12): API, New(), param(), parseList(), unknownCollection(), projectAccessResource(), tagResource(), workspaceTagResource() (+4 more)

### Community 154 - "NewSubscriptions"
Cohesion: 0.13
Nodes (20): Deps, EventSource, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity() (+12 more)

### Community 157 - "Job"
Cohesion: 0.05
Nodes (18): captureJobs(), captureSteps(), Publisher, publishingStore, stubQuerier, fakeStore, fakeStore, CreateJobParams (+10 more)

### Community 159 - "stubQuerier"
Cohesion: 0.03
Nodes (11): stubQuerier, Agent, Collection, Globalvar, Item, NotificationTrigger, Project, Reference (+3 more)

### Community 161 - "Executor"
Cohesion: 0.08
Nodes (13): Executor, preparedWorkspace, isCancelledStatus(), isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestBinarySelection() (+5 more)

### Community 162 - "runTriggerAPI"
Cohesion: 0.12
Nodes (24): missingOrganizationStub, runTriggerStub, authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite() (+16 more)

### Community 163 - "History"
Cohesion: 0.16
Nodes (4): CreateHistoryParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams, History

### Community 164 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 165 - "webhook_test.go"
Cohesion: 0.36
Nodes (10): DeleteWebhook(), SyncWebhook(), newHookServer(), TestDeleteWebhookTreats404AsRemoved(), TestGitHubWebhookCreate(), TestGitHubWebhookUpdateAndDelete(), TestGitLabWebhookCreate(), TestGitLabWebhookDelete() (+2 more)

### Community 166 - "Organization"
Cohesion: 0.24
Nodes (3): scopeIDs(), demoOrganizations(), Organization

### Community 167 - "API"
Cohesion: 0.06
Nodes (24): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+16 more)

### Community 168 - "notificationSendAPI"
Cohesion: 0.14
Nodes (23): newDeliveryStub(), TestNotificationDeliveryHistory(), TestNotificationDeliveryHistoryRejectsABadLimit(), TestNotificationDeliveryHistoryUnknownWorkspace(), TestNotificationDeliveryRetry(), TestNotificationDeliveryRetryRejections(), captureJSONServer(), notificationSendAPI() (+15 more)

### Community 171 - "github_test.go"
Cohesion: 0.38
Nodes (9): githubRoutes(), TestBranchesGitHub(), TestGitHubConnectionsWithoutAToken(), TestGroupsGitHub(), TestRepositoriesGitHub(), TestTestConnection(), newFakeProvider(), TestRejectedCredentialIsAUseful4xx() (+1 more)

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
Cohesion: 0.22
Nodes (9): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL(), runGit() (+1 more)

### Community 190 - "store.Workspace"
Cohesion: 0.08
Nodes (16): atomicOperation, webhookEventAttrs, webhookSyncRef, strValue(), atomicWorkspaceHref(), isAtomicRequest(), discoveryPage(), writeRawJSON() (+8 more)

### Community 191 - "loadPage"
Cohesion: 0.24
Nodes (3): DB, failingDB, loadPage()

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

### Community 232 - "Review Focus"
Cohesion: 0.33
Nodes (5): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract

### Community 234 - "Team"
Cohesion: 0.07
Nodes (12): AccessTokenStore, Deps, fakeAccessTokenStore, contains(), captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces() (+4 more)

### Community 235 - ".build"
Cohesion: 0.15
Nodes (14): collectionSpec, page, accessEdges(), column(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), predicate() (+6 more)

### Community 236 - "Queries"
Cohesion: 0.04
Nodes (17): createRawVariable(), createRawVariableWithID(), CountJobsByStatusRow, DBTX, ListGlobalVariableSecretsRow, ListNotificationDeliveriesByWorkspaceParams, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow (+9 more)

### Community 252 - "workspaceQuery"
Cohesion: 0.40
Nodes (3): matchesTagFilters(), tagFilter, workspaceQuery

### Community 259 - "newService"
Cohesion: 0.23
Nodes (11): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRequiresMembershipAndAnInteractiveUser(), TestCreateSignsTheUpstreamToken(), TestListAndDeleteAreScopedToTheCallersTeams() (+3 more)

### Community 260 - "Attr"
Cohesion: 0.29
Nodes (5): exists(), likePattern(), relate(), Attr, Kind

### Community 261 - "grantsAPI"
Cohesion: 0.31
Nodes (10): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+2 more)

### Community 262 - "Content"
Cohesion: 0.29
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 273 - ".server"
Cohesion: 0.14
Nodes (13): fakeHooks, dexFixture, TestAtomicWebhookInvalidTriggerChangesNothing(), TestAtomicWebhookSaveRegistersFinalTriggers(), TestCreateWebhookEventRegistersProvider(), TestDeleteWebhookRemovesProviderHook(), TestWebhookEventTemplateOutsideOrganizationIsRejected(), webhookWriteStub() (+5 more)

### Community 276 - "Deps"
Cohesion: 0.29
Nodes (3): Deps, fakeLocker, JobLocker

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 552 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **142 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `API` connect `API` to `WriteError`, `API`, `Team`, `api/main.go`, `store.Workspace`, `Link`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07363770250368189 - nodes in this community are weakly interconnected._
- **Why does `Rotate()` connect `Rotate` to `context.Context`, `run`, `Queries`, `Cipher`, `api/main.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `WriteError` be split into smaller, more focused modules?**
  _Cohesion score 0.08295281582952815 - nodes in this community are weakly interconnected._
- **Why does `Queries` connect `Queries` to `.CreateTeam`, `Content`, `run`, `fakeQuerier`, `.CreateRunTrigger`, `.UpdateTeam`, `newTestAPI`, `Rotate`, `github.com/aminbenmansour/terrakube-go/internal/rsql.Node`, `WebhookEvent`, `.CreateGlobalVariable`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `Job`, `api/main.go`, `History`, `.QueryRow`, `ResolveOrganization`, `CliAuthSession`, `Pat`, `TestRotateReencryptsOldKeyAndPlaintext`, `Queries`, `Action`, `Workspace`, `.CreateAgent`, `.CreateProject`, `Cipher`, `runTriggerAPI`, `.CreateTag`, `context.Context`, `TeamToken`, `Team`, `.CreateProjectAccess`, `TestWebhookQueries`, `.UpdateWebhook`, `.CreateArchive`, `Ssh`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._