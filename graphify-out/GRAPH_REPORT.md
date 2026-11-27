# Graph Report - terrakube-go  (2026-11-27)

## Corpus Check
- 523 files · ~273,162 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4809 nodes · 19677 edges · 259 communities (129 shown, 130 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3489 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `815ca67d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- ._Project
- io.Reader
- Core Database Schema
- Team
- run
- net/http.ResponseWriter
- newFixture
- WriteError
- fakeQuerier
- Service
- steps
- github.com/99designs/gqlgen/graphql.Marshaler
- github.com/aminbenmansour/terrakube-go/internal/rsql.Node
- deliver_test.go
- tfe_test.go
- sqlc Generation Config
- serveAccessToken
- Upstream Elide entities (GraphQL surface)
- Step
- templateStub
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
- ._NotificationConfiguration
- serveGlobalVar
- .QueryRow
- Queries
- .createVCS
- tagStub
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- provider
- models
- Load
- Access
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
- logSink
- Secret encryption at rest — design
- .CreateProject
- REST contract, captured from upstream
- pat.py
- log/slog.Logger
- stubQuerier
- vcs/vcs.go
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
- WorkspaceTagConnection
- compare.sh
- start.sh
- NotificationTriggerConnection
- newTestExecutor
- time.Time
- WebhookEvent
- time.Duration
- rsql.Node
- reliabilityStore
- clilogin/handler_test.go
- newTestAPI
- Secret key rotation — design
- models_gen.go
- ui.go
- API
- New
- fakeStore
- NewSubscriptions
- Job
- stubQuerier
- Executor
- grantsAPI
- History
- newFakeProvider
- Dex
- M4 run logs and artifacts Implementation Plan
- net/http.Handler
- notificationSendAPI
- ResolveOrganization
- Collection
- Metrics
- transientError
- .EnqueueJob
- Threat model
- StringOrInt
- auth.DexValidator
- .updateWorkspaceRelationship
- TestModelColumnOrder
- RSQL filters in SQL — design
- ResourceIdentifier
- ProjectAccess
- Workspace
- Review Focus
- ParsePage
- Schedule
- Review Focus
- Repository
- Team
- .build
- .Query
- IssueCliAuthCodeParams

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 246 edges
3. `Queries` - 244 edges
4. `API` - 198 edges
5. `New()` - 165 edges
6. `stubQuerier` - 164 edges
7. `WriteError()` - 140 edges
8. `param()` - 88 edges
9. `ParseFieldsets()` - 74 edges
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
- `Secret Encryption at Rest Implementation Plan` --references--> `Querier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/store/querier.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (259 total, 130 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (28): configurationVersionResource(), projectResource(), parseWorkspaceQuery(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource() (+20 more)

### Community 1 - "rbac.go"
Cohesion: 0.26
Nodes (22): Task 1: Role matrix and permission set, adminOnly(), CanApproveJob(), CanDestroyJob(), CanManageCollection(), CanManageJob(), CanManageModule(), CanManagePolicies() (+14 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (12): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+4 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "Team"
Cohesion: 0.26
Nodes (3): contains(), contains(), fakeOrganizationStore

### Community 7 - "run"
Cohesion: 0.05
Nodes (69): main(), run(), main(), run(), main(), run(), main(), run() (+61 more)

### Community 8 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (19): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), WithClaims(), runnerClaims(), claimsFromContext() (+11 more)

### Community 9 - "newFixture"
Cohesion: 0.13
Nodes (25): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+17 more)

### Community 10 - "WriteError"
Cohesion: 0.09
Nodes (38): actionAttrs, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, vcsAttrs, decodeActionAttrs(), stringOr() (+30 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.07
Nodes (6): nestedStub, encryptedStore, fakeQuerier, CreateNotificationConfigurationParams, UpdateNotificationConfigurationParams, Variable

### Community 13 - "Service"
Cohesion: 0.16
Nodes (9): Service, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), randomToken(), sanitizeName() (+1 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 17 - "github.com/aminbenmansour/terrakube-go/internal/rsql.Node"
Cohesion: 0.07
Nodes (4): listQuery, stubLists, Task 3: Handler cutover, ListSpec

### Community 18 - "deliver_test.go"
Cohesion: 0.12
Nodes (25): NewDeliverer(), deliveryError(), stringPtr(), TestDeliverDoesNotFollowRedirects(), testDeliverer(), TestDeliverIgnoresMalformedRetryAfter(), TestDeliverNetworkErrorIsRetryableAndHidesTheURL(), TestDeliverParsesHTTPDateRetryAfter() (+17 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (42): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+34 more)

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
Nodes (8): captureSteps(), dispatcherStore, CreateStepParams, FailStuckNotificationOutboxParams, GetStepByJobAndNumberParams, ReclaimStuckNotificationOutboxParams, RecordNotificationOutboxParams, Step

### Community 24 - "templateStub"
Cohesion: 0.14
Nodes (3): templateStub, CreateTemplateParams, UpdateTemplateParams

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
Nodes (59): actionFilter, agentAttrs, atomicOperation, atomicResult, authorizationKey, collectionAttrs, createPATRequest, discoveryGroup (+51 more)

### Community 31 - "Link"
Cohesion: 0.08
Nodes (35): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), accessResource(), projectAccessResource() (+27 more)

### Community 32 - "getJSON"
Cohesion: 0.20
Nodes (8): getJSON(), providerMessage(), Branch, BranchPage, github, gitlab, Group, RepositoryPage

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

### Community 44 - "Content"
Cohesion: 0.29
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 47 - ".streamStepLog"
Cohesion: 0.15
Nodes (15): compareStreamIDs(), isTerminal(), outputOf(), parseResumeID(), splitStreamID(), stepOf(), streamKey(), TestCompareStreamIDs() (+7 more)

### Community 48 - "CliAuthSession"
Cohesion: 0.13
Nodes (5): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams, SetCliAuthSessionIdentityParams

### Community 49 - "stubQuerier"
Cohesion: 0.06
Nodes (5): jobStub, statusQuerier, stubQuerier, workspaceWithTags, stubQuerier

### Community 51 - "time.Duration"
Cohesion: 0.10
Nodes (4): fakeQueue, JobQueue, fakePinger, Pinger

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (120): TestAccessTokenWithoutStore(), newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+112 more)

### Community 53 - "Pat"
Cohesion: 0.17
Nodes (5): fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 57 - "serveGlobalVar"
Cohesion: 0.18
Nodes (10): globalVarStub, newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent(), TestGlobalVariableFromAnotherOrganizationIsNotFound() (+2 more)

### Community 59 - ".QueryRow"
Cohesion: 0.04
Nodes (8): ApproveJobParams, CreateWebhookParams, NotificationOutboxExistsForWorkspaceParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams, UpdateOrganizationParams, UpdateWebhookParams, Webhook

### Community 61 - "Queries"
Cohesion: 0.02
Nodes (10): startScheduler(), NewDispatcher(), NewDispatcher(), AddNotificationConfigurationTemplateParams, Querier, Queries, ReclaimStuckRunTriggerEventsParams, UpdateHistoryOutputParams (+2 more)

### Community 62 - ".createVCS"
Cohesion: 0.25
Nodes (7): strPtr(), strValue(), discoveryConnection(), discoveryPage(), writeRawJSON(), validVCSEnum(), vcsEnumOrDefault()

### Community 63 - "tagStub"
Cohesion: 0.08
Nodes (4): tagStub, matchesTagFilters(), tagFilter, workspaceQuery

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.05
Nodes (64): archiveLimits, tarEntry, gitOutput(), initGitRepo(), requireGit(), tarGzEntries(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive() (+56 more)

### Community 67 - "provider"
Cohesion: 0.15
Nodes (11): Connection shape, Entity and relationship surface, Filter expressions the UI generates, Other contracts observed, Sort strings, Upstream UI GraphQL operations, Measured route surface, Order that unblocks the most UI (+3 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "Load"
Cohesion: 0.14
Nodes (22): Config, Storage, Task 2: Serve archived logs from `/tfoutput/v1`, Task 1: Keyring configuration, env(), envBool(), envFirst(), envInt64() (+14 more)

### Community 70 - "Access"
Cohesion: 0.18
Nodes (7): captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces(), captureProjectAccessRow(), fakeWorkspaceStore, Access

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - "New"
Cohesion: 0.23
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 73 - "Cipher"
Cohesion: 0.08
Nodes (32): Review Focus, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation, Task 7: Documentation and graph, Task 8: Verification, Task 3: Rotation scan and re-encryption (+24 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.19
Nodes (15): column(), predicate(), Field, likePattern(), mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode() (+7 more)

### Community 76 - "streamEvents"
Cohesion: 0.14
Nodes (8): projectResolver, Resolver, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 77 - ".prepareWorkspace"
Cohesion: 0.15
Nodes (12): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv(), gitSource() (+4 more)

### Community 78 - "decodeJSON"
Cohesion: 0.09
Nodes (43): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath() (+35 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "logSink"
Cohesion: 0.12
Nodes (9): captureLogPublisher, LogPublisher, logSink, redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues() (+1 more)

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - ".CreateProject"
Cohesion: 0.20
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "log/slog.Logger"
Cohesion: 0.16
Nodes (11): Validator, formatOldKeys(), logReport(), Deps, JobLocker, Deps, Backend, Layout (+3 more)

### Community 89 - "stubQuerier"
Cohesion: 0.16
Nodes (20): replay, upstreamCapture, assertSameDocumentShape(), assertSameResourceShape(), captureActionStore(), captureCollectionStore(), captureGlobalVarStore(), captureJobStore() (+12 more)

### Community 91 - "vcs/vcs.go"
Cohesion: 0.33
Nodes (12): TestGitLabConnectionWithoutAToken(), TestGitLabDiscovery(), baseURL(), Branches(), failed(), Test(), Groups(), providerFor() (+4 more)

### Community 94 - "Identity"
Cohesion: 0.09
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
Cohesion: 0.14
Nodes (22): AuthorizeRequest, TokenIssuer, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect() (+14 more)

### Community 103 - "context.Context"
Cohesion: 0.03
Nodes (10): collectionStub, sshAgentStub, teamStub, vcsStub, workspaceStub, artifactStore, historyStub, variableStore (+2 more)

### Community 108 - "newService"
Cohesion: 0.19
Nodes (15): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+7 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.35
Nodes (11): newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit(), TestCreateRunTriggerRejectsCycle(), TestCreateRunTriggerRejectsStructuralViolations(), TestDeleteRunTriggerAndContentTypeQuirk() (+3 more)

### Community 121 - "Dispatcher"
Cohesion: 0.27
Nodes (4): NewDispatcher(), strPtr(), Dispatcher, NotificationOutbox

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 135 - "newTestExecutor"
Cohesion: 0.12
Nodes (20): fakeLocker, newFakeLocker(), TestHandleReleasesTheJobLock(), TestHandleSkipsAJobHeldByAnotherExecutor(), TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestStaleClaimIdleOutlivesTheJobTimeout(), New() (+12 more)

### Community 137 - "time.Time"
Cohesion: 0.21
Nodes (7): bucket, limiter, newLimiter(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), Dispatcher, memObject

### Community 140 - "WebhookEvent"
Cohesion: 0.26
Nodes (3): CreateWebhookEventParams, UpdateWebhookEventParams, WebhookEvent

### Community 141 - "time.Duration"
Cohesion: 0.10
Nodes (11): waitFor(), envDuration(), defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), backoff(), blockedDelivery(), backoff() (+3 more)

### Community 142 - "rsql.Node"
Cohesion: 0.15
Nodes (8): exists(), likePattern(), relate(), And, Or, parser, Attr, Kind

### Community 143 - "reliabilityStore"
Cohesion: 0.09
Nodes (6): discoveryStub, missingOrganizationStub, cancelRunStore, reliabilityStore, UpdateStepStatusParams, mutableQuerier

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.08
Nodes (39): pageMeta, createListOrganization(), getNames(), seedOrganizations(), TestOrganizationListEscapesWildcards(), TestOrganizationListPrefiltersByAllowedIDs(), TestOrganizationListSortsAndPages(), TestOrganizationListSQLFilters() (+31 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 147 - "models_gen.go"
Cohesion: 0.07
Nodes (9): AccessConnection, AccessEdge, OrganizationConnection, OrganizationEdge, PageInfo, Query, Subscription, TemplateConnection (+1 more)

### Community 149 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 150 - "API"
Cohesion: 0.10
Nodes (9): API, New(), param(), parseList(), unknownCollection(), tagResource(), teamResource(), ParseFieldsets() (+1 more)

### Community 151 - "New"
Cohesion: 0.18
Nodes (10): requireMetricsToken(), key(), RateLimit(), TestRateLimitDisabledByDefault(), TestRateLimitMiddleware(), SecurityHeaders(), TestSecurityHeaders(), New() (+2 more)

### Community 152 - "fakeStore"
Cohesion: 0.08
Nodes (20): strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased(), TestTickIgnoresPlanOnlyRun() (+12 more)

### Community 154 - "NewSubscriptions"
Cohesion: 0.07
Nodes (35): main(), run(), Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext() (+27 more)

### Community 157 - "Job"
Cohesion: 0.06
Nodes (19): captureJobs(), strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued(), TestTickSkipsLockedWorkspaceAndUnparseableCron(), TestTickUsesTheTemplateOfTheWorkspaceOrganization(), tickAt() (+11 more)

### Community 159 - "stubQuerier"
Cohesion: 0.02
Nodes (23): actionStub, runTriggerStub, stubQuerier, scopeIDs(), demoOrganizations(), CreateActionParams, UpdateActionParams, CreateItemParams (+15 more)

### Community 161 - "Executor"
Cohesion: 0.10
Nodes (10): Executor, preparedWorkspace, isCancelledStatus(), isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestBinarySelection() (+2 more)

### Community 162 - "grantsAPI"
Cohesion: 0.31
Nodes (10): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+2 more)

### Community 163 - "History"
Cohesion: 0.16
Nodes (4): CreateHistoryParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams, History

### Community 164 - "newFakeProvider"
Cohesion: 0.38
Nodes (9): githubRoutes(), TestBranchesGitHub(), TestGitHubConnectionsWithoutAToken(), TestGroupsGitHub(), TestRepositoriesGitHub(), TestTestConnection(), newFakeProvider(), TestRejectedCredentialIsAUseful4xx() (+1 more)

### Community 165 - "Dex"
Cohesion: 0.23
Nodes (3): Dex, Exchanger, fakeDex

### Community 166 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.25
Nodes (7): Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph

### Community 167 - "net/http.Handler"
Cohesion: 0.07
Nodes (25): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+17 more)

### Community 168 - "notificationSendAPI"
Cohesion: 0.14
Nodes (23): newDeliveryStub(), TestNotificationDeliveryHistory(), TestNotificationDeliveryHistoryRejectsABadLimit(), TestNotificationDeliveryHistoryUnknownWorkspace(), TestNotificationDeliveryRetry(), TestNotificationDeliveryRetryRejections(), captureJSONServer(), notificationSendAPI() (+15 more)

### Community 169 - "ResolveOrganization"
Cohesion: 0.20
Nodes (9): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access, OrganizationStore (+1 more)

### Community 171 - "Collection"
Cohesion: 0.16
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 172 - "Metrics"
Cohesion: 0.13
Nodes (4): metricsSource, fakeStats, Metrics, StatsSource

### Community 175 - ".EnqueueJob"
Cohesion: 0.22
Nodes (5): fakeQueue, JobQueue, fakeQueue, JobQueue, stubJobQueue

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - "StringOrInt"
Cohesion: 0.10
Nodes (3): ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge

### Community 184 - "auth.DexValidator"
Cohesion: 0.31
Nodes (4): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 191 - "TestModelColumnOrder"
Cohesion: 0.25
Nodes (4): failingDB, collectStrings(), snakeCase(), TestModelColumnOrder()

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 197 - "ResourceIdentifier"
Cohesion: 0.29
Nodes (5): nestedDocument, RelationshipList(), identifierStringAttribute(), RequestRelationship, ResourceIdentifier

### Community 198 - "ProjectAccess"
Cohesion: 0.15
Nodes (4): AccessTokenStore, Deps, fakeAccessTokenStore, ProjectAccess

### Community 202 - "Review Focus"
Cohesion: 0.33
Nodes (5): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract

### Community 204 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 229 - "Schedule"
Cohesion: 0.12
Nodes (13): scheduleStub, newScheduleAPI(), newScheduleStub(), scheduleBody(), scheduleStr(), TestCreateScheduleCopiesAttributesAndDefaultsEnabled(), TestCreateScheduleRejectsForeignTemplate(), TestDeleteScheduleAndCrossWorkspaceLookup() (+5 more)

### Community 231 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 232 - "Repository"
Cohesion: 0.40
Nodes (4): optional(), githubRepositoryJSON, gitlabProjectJSON, Repository

### Community 235 - ".build"
Cohesion: 0.10
Nodes (20): collectionSpec, DB, listArgs, page, accessEdges(), afterString(), jobEdges(), loadPage() (+12 more)

### Community 236 - ".Query"
Cohesion: 0.06
Nodes (19): CountJobsByStatusRow, CreateArchiveParams, DBTX, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow (+11 more)

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 549 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **130 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Queries` to `IssueCliAuthCodeParams`, `Team`, `run`, `net/http.ResponseWriter`, `time.Time`, `fakeQuerier`, `WebhookEvent`, `.marshalOBoolean2bool`, `tfe_test.go`, `Step`, `fakeStore`, `templateStub`, `Job`, `api/main.go`, `stubQuerier`, `History`, `notificationSendAPI`, `Collection`, `Metrics`, `Content`, `Globalvar`, `Tag`, `CliAuthSession`, `stubQuerier`, `testing.T`, `Pat`, `.QueryRow`, `tagStub`, `NewLayout`, `ProjectAccess`, `Workspace`, `Cipher`, `.CreateProject`, `log/slog.Logger`, `Agent`, `Schedule`, `context.Context`, `Team`, `.Query`, `.CreateProjectAccess`, `Dispatcher`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07668418537983755 - nodes in this community are weakly interconnected._
- **Why does `Rotate()` connect `run` to `context.Context`, `Cipher`, `Secret key rotation — design`, `Queries`, `api/main.go`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `Store boundary` connect `Secret encryption at rest — design` to `Cipher`, `NewSubscriptions`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `io.Reader` be split into smaller, more focused modules?**
  _Cohesion score 0.07164404223227752 - nodes in this community are weakly interconnected._