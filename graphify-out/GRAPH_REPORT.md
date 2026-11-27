# Graph Report - terrakube-go  (2026-11-27)

## Corpus Check
- 509 files · ~265,505 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4675 nodes · 19009 edges · 271 communities (124 shown, 147 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3384 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3233ea50`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- ._Project
- memStorage
- Core Database Schema
- Team
- NewCipher
- net/http.ResponseWriter
- newFixture
- .internalError
- fakeQuerier
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- Service
- steps
- github.com/aminbenmansour/terrakube-go/internal/rsql.Node
- deliver_test.go
- tfe_test.go
- sqlc Generation Config
- serveAccessToken
- Upstream Elide entities (GraphQL surface)
- fakeStore
- templateStub
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
- Content
- dispatcherStore
- graphql/collections.go
- Variable
- CliAuthSession
- JobConnection
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- ._NotificationConfiguration
- ResourceIdentifier
- Tag
- Job
- Queries
- reliabilityStore
- tagStub
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- WorkspaceTagConnection
- models
- Load
- Team
- collections.go
- rotateValue
- Cipher
- generated.go
- where_test.go
- streamEvents
- .prepareWorkspace
- decodeJSON
- graphql/models.go
- newJWKSServer
- newLogSink
- Secret encryption at rest — design
- ui.go
- REST contract, captured from upstream
- pat.py
- Layout
- stubQuerier
- WorkspaceRunTrigger
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
- CreateNotificationConfigurationParams
- time.Duration
- compare.sh
- start.sh
- .CreateGlobalVariable
- Review Focus
- time.Time
- tarGzEntries
- .processJob
- rsql.Node
- cancelRunStore
- clilogin/handler_test.go
- newTestAPI
- Secret key rotation — design
- models_gen.go
- Agent
- API
- log/slog.Logger
- addDestination
- .CreateNotificationTrigger
- NewSubscriptions
- Dispatcher
- stubQuerier
- Executor
- TemplateConnection
- .QueryRow
- .streamStepLog
- Dex
- github.com/aminbenmansour/terrakube-go/internal/store.Step
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
- io.Reader
- serveNested
- RSQL filters in SQL — design
- grantsAPI
- Dispatcher
- Workspace
- .Ping
- .Get
- Schedule
- Review Focus
- .updateWorkspaceRelationship
- .build
- .Query
- Deps
- Object
- query.go
- cleanKey
- .SetGlobalVariableSecret
- .SetNotificationSecret
- .SetSSHKeySecret
- .SetVariableSecret
- Local
- IssueCliAuthCodeParams
- .SetCliAuthSessionIdentity
- .SoftDeletePat
- .ListNotificationSecretsAfter
- .ListSSHKeySecretsAfter

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 236 edges
3. `Queries` - 234 edges
4. `API` - 188 edges
5. `stubQuerier` - 158 edges
6. `New()` - 156 edges
7. `WriteError()` - 135 edges
8. `param()` - 86 edges
9. `ParseFieldsets()` - 72 edges
10. `API` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Secret Encryption at Rest Implementation Plan` --references--> `Querier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/store/querier.go
- `Task 3: Transparent store decorator` --references--> `Querier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/store/querier.go
- `Task 3: Transparent store decorator` --references--> `fakeQuerier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/secretstore/store_test.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (271 total, 147 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (27): Write(), WriteResource(), configurationVersionResource(), projectResource(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID() (+19 more)

### Community 1 - "rbac.go"
Cohesion: 0.19
Nodes (27): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, adminOnly(), CanApproveJob() (+19 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "._Project"
Cohesion: 0.13
Nodes (3): ProjectAccessConnection, ProjectAccessEdge, ProjectEdge

### Community 4 - "memStorage"
Cohesion: 0.14
Nodes (5): failingBackend, newMemStorage(), logBackend, memObject, memStorage

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "Team"
Cohesion: 0.08
Nodes (4): contains(), fakeOrganizationStore, OrganizationStore, Team

### Community 7 - "NewCipher"
Cohesion: 0.16
Nodes (30): TestBackfillEncryptsPlaintextRows(), NewCipher(), TestCipherEmptyAndPlaintext(), TestCipherKeyring(), TestCipherRejectsBadConfiguration(), TestCipherRejectsWrongColumnAndKey(), TestCipherRoundTrip(), testKey() (+22 more)

### Community 8 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (20): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), WithClaims(), runnerClaims(), claimsFromContext() (+12 more)

### Community 9 - "newFixture"
Cohesion: 0.09
Nodes (34): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+26 more)

### Community 10 - ".internalError"
Cohesion: 0.09
Nodes (33): actionAttrs, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, decodeActionAttrs(), stringOr(), boolValue() (+25 more)

### Community 12 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.09
Nodes (3): jobStub, statusQuerier, stubQuerier

### Community 13 - "Service"
Cohesion: 0.15
Nodes (11): Service, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), NewService() (+3 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 17 - "github.com/aminbenmansour/terrakube-go/internal/rsql.Node"
Cohesion: 0.06
Nodes (8): listQuery, stubLists, Task 3: Handler cutover, scopeIDs(), demoOrganizations(), workspaceUpdateParams(), Organization, Workspace

### Community 18 - "deliver_test.go"
Cohesion: 0.12
Nodes (18): deliveryError(), stringPtr(), TestDeliverDoesNotFollowRedirects(), testDeliverer(), TestDeliverIgnoresMalformedRetryAfter(), TestDeliverNetworkErrorIsRetryableAndHidesTheURL(), TestDeliverParsesHTTPDateRetryAfter(), TestDeliverPostsPayload() (+10 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.07
Nodes (47): NewLocal(), TestLocalHealth(), TestLocalLifecycle(), TestLocalRejectsPathTraversal(), NewS3(), New(), assertContractJSONAPI(), TestContractDiscoveryAndPing() (+39 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "serveAccessToken"
Cohesion: 0.12
Nodes (24): expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups(), TestAccessTokenRejectsBadInput() (+16 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "fakeStore"
Cohesion: 0.05
Nodes (16): captureJobs(), captureSteps(), stubQuerier, fakeStore, fakeStore, CreateJobParams, CreateStepParams, FailStuckRunTriggerEventsParams (+8 more)

### Community 24 - "templateStub"
Cohesion: 0.17
Nodes (3): templateStub, CreateTemplateParams, UpdateTemplateParams

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
Nodes (63): actionFilter, agentAttrs, atomicOperation, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs (+55 more)

### Community 31 - "Link"
Cohesion: 0.09
Nodes (30): vcsAttrs, actionResource(), strPtr(), strValue(), collectionResource(), itemResource(), referenceResource(), globalVarResource() (+22 more)

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

### Community 44 - "Content"
Cohesion: 0.29
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 45 - "dispatcherStore"
Cohesion: 0.11
Nodes (12): dispatcherWith(), newDispatcherFixture(), TestDispatcherFailsAfterMaxAttempts(), TestDispatcherFailsPermanentRejection(), TestDispatcherIncludesTheFailureSummary(), TestDispatcherRetriesTransientFailures(), TestDispatcherSendsAndRecordsSuccess(), TestDispatcherSkipsSupersededJobStatus() (+4 more)

### Community 46 - "graphql/collections.go"
Cohesion: 0.18
Nodes (10): page, accessEdges(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), projectAccessEdges(), projectEdges(), templateEdges() (+2 more)

### Community 48 - "CliAuthSession"
Cohesion: 0.18
Nodes (4): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams

### Community 51 - "time.Duration"
Cohesion: 0.10
Nodes (3): captureLogPublisher, fakeQueue, JobQueue

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (112): TestAccessTokenWithoutStore(), newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+104 more)

### Community 53 - "Pat"
Cohesion: 0.23
Nodes (4): fakeStore, Store, CreatePatParams, Pat

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 57 - "ResourceIdentifier"
Cohesion: 0.20
Nodes (10): nestedDocument, RelationshipIdentifiers(), RelationshipList(), TestRelationshipIdentifiers(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship (+2 more)

### Community 60 - "Job"
Cohesion: 0.12
Nodes (4): ApproveJobParams, GetJobByWorkspaceAndTriggerParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams

### Community 61 - "Queries"
Cohesion: 0.02
Nodes (7): AddNotificationConfigurationTemplateParams, DBTX, Querier, Queries, UpdateHistoryOutputParams, UpdateJobOutputParams, UpdateStepOutputParams

### Community 62 - "reliabilityStore"
Cohesion: 0.20
Nodes (3): reliabilityStore, CreateHistoryParams, UpdateStepStatusParams

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.06
Nodes (54): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+46 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "Load"
Cohesion: 0.06
Nodes (50): main(), run(), main(), run(), main(), run(), formatOldKeys(), logReport() (+42 more)

### Community 70 - "Team"
Cohesion: 0.05
Nodes (10): AccessTokenStore, Deps, fakeAccessTokenStore, contains(), captureAccessRow(), captureGrantStore(), captureGrantWorkspaces(), captureProjectAccessRow() (+2 more)

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 72 - "rotateValue"
Cohesion: 0.30
Nodes (12): envelopeKey(), rotateGlobalVariables(), rotateNotifications(), rotatePtr(), rotateSSHKeys(), rotateValue(), rotateVariables(), rotateVCS() (+4 more)

### Community 73 - "Cipher"
Cohesion: 0.15
Nodes (21): Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation (+13 more)

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.19
Nodes (16): column(), predicate(), Field, likePattern(), mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode() (+8 more)

### Community 76 - "streamEvents"
Cohesion: 0.14
Nodes (8): projectResolver, Resolver, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

### Community 77 - ".prepareWorkspace"
Cohesion: 0.15
Nodes (11): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), transient(), contentIDFromSource(), gitEnv(), gitSource(), isConfigurationTarball() (+3 more)

### Community 78 - "decodeJSON"
Cohesion: 0.10
Nodes (42): newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent(), TestGlobalVariableFromAnotherOrganizationIsNotFound(), TestListGlobalVariablesHidesSensitiveValues() (+34 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "newLogSink"
Cohesion: 0.13
Nodes (11): LogPublisher, logSink, redactor, TestLogSinkSplitsLinesAndKeepsTail(), TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues() (+3 more)

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

### Community 87 - "Layout"
Cohesion: 0.11
Nodes (12): Validator, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1` (+4 more)

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
Cohesion: 0.19
Nodes (5): Deps, EventSource, fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (20): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+12 more)

### Community 103 - "context.Context"
Cohesion: 0.03
Nodes (11): collectionStub, globalVarStub, sshAgentStub, teamStub, vcsStub, workspaceStub, artifactStore, historyStub (+3 more)

### Community 108 - "newService"
Cohesion: 0.19
Nodes (15): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+7 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.14
Nodes (22): runTriggerStub, newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), runTriggerRBACAPI(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit(), TestCreateRunTriggerRejectsCycle() (+14 more)

### Community 121 - "Dispatcher"
Cohesion: 0.27
Nodes (4): NewDispatcher(), strPtr(), Dispatcher, NotificationOutbox

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 128 - "time.Duration"
Cohesion: 0.13
Nodes (9): waitFor(), envDuration(), defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), backoff(), backoff(), TestBackoffGrowsAndCaps() (+1 more)

### Community 135 - "Review Focus"
Cohesion: 0.20
Nodes (8): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 4: `rotatekeys` command, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 137 - "time.Time"
Cohesion: 0.24
Nodes (9): bucket, limiter, key(), newLimiter(), RateLimit(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), TestRateLimitDisabledByDefault() (+1 more)

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 141 - ".processJob"
Cohesion: 0.16
Nodes (5): isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestIsTerminalStatus()

### Community 142 - "rsql.Node"
Cohesion: 0.08
Nodes (23): Parse(), dump(), TestParse(), TestParseRejects(), exists(), likePattern(), relate(), ptr() (+15 more)

### Community 143 - "cancelRunStore"
Cohesion: 0.15
Nodes (3): missingOrganizationStub, cancelRunStore, mutableQuerier

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.12
Nodes (19): CookieSigner, Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), authorize() (+11 more)

### Community 145 - "newTestAPI"
Cohesion: 0.11
Nodes (27): pageMeta, createListOrganization(), getNames(), seedOrganizations(), TestOrganizationListEscapesWildcards(), TestOrganizationListPrefiltersByAllowedIDs(), TestOrganizationListSortsAndPages(), TestOrganizationListSQLFilters() (+19 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (14): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+6 more)

### Community 147 - "models_gen.go"
Cohesion: 0.07
Nodes (9): AccessConnection, AccessEdge, OrganizationConnection, OrganizationEdge, PageInfo, Query, Subscription, WorkspaceConnection (+1 more)

### Community 150 - "API"
Cohesion: 0.09
Nodes (14): API, New(), param(), parseInclude(), parseList(), unknownCollection(), accessResource(), projectAccessResource() (+6 more)

### Community 151 - "log/slog.Logger"
Cohesion: 0.19
Nodes (11): startScheduler(), Deps, Pinger, requireMetricsToken(), SecurityHeaders(), TestSecurityHeaders(), New(), requestLogger() (+3 more)

### Community 152 - "addDestination"
Cohesion: 0.37
Nodes (15): strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased(), TestTickIgnoresPlanOnlyRun() (+7 more)

### Community 154 - "NewSubscriptions"
Cohesion: 0.11
Nodes (22): Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestSubscriptionAcceptsADexConnectionToken(), NewExecutableSchema(), TestSchemaLoads(), NewSubscriptions() (+14 more)

### Community 157 - "Dispatcher"
Cohesion: 0.28
Nodes (8): strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued(), TestTickSkipsLockedWorkspaceAndUnparseableCron(), TestTickUsesTheTemplateOfTheWorkspaceOrganization(), tickAt(), Dispatcher

### Community 159 - "stubQuerier"
Cohesion: 0.03
Nodes (20): actionStub, stubQuerier, captureGrantProjects(), CreateActionParams, UpdateActionParams, CreateItemParams, CreateReferenceParams, UpdateItemParams (+12 more)

### Community 161 - "Executor"
Cohesion: 0.20
Nodes (5): Executor, preparedWorkspace, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

### Community 163 - ".QueryRow"
Cohesion: 0.05
Nodes (8): CreateProjectParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams, GetProjectByNameParams, History, NotificationOutboxExistsForWorkspaceParams, UpdateOrganizationParams, UpdateProjectParams

### Community 164 - ".streamStepLog"
Cohesion: 0.23
Nodes (10): isTerminal(), outputOf(), stepOf(), streamKey(), TestIsTerminal(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData() (+2 more)

### Community 165 - "Dex"
Cohesion: 0.23
Nodes (3): Dex, Exchanger, fakeDex

### Community 167 - "net/http.Handler"
Cohesion: 0.07
Nodes (31): authorization, main(), run(), Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom() (+23 more)

### Community 168 - "notificationSendAPI"
Cohesion: 0.13
Nodes (23): newDeliveryStub(), TestNotificationDeliveryHistory(), TestNotificationDeliveryHistoryRejectsABadLimit(), TestNotificationDeliveryHistoryUnknownWorkspace(), TestNotificationDeliveryRetry(), TestNotificationDeliveryRetryRejections(), captureJSONServer(), notificationSendAPI() (+15 more)

### Community 169 - "ResolveOrganization"
Cohesion: 0.12
Nodes (9): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveProject(), ResolveWorkspace(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access, fakeWorkspaceStore, WorkspaceStore (+1 more)

### Community 171 - "Collection"
Cohesion: 0.16
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 172 - "Metrics"
Cohesion: 0.13
Nodes (4): metricsSource, fakeStats, Metrics, StatsSource

### Community 175 - ".EnqueueJob"
Cohesion: 0.23
Nodes (6): fakeQueue, JobQueue, fakeQueue, JobQueue, JobQueue, stubJobQueue

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - "StringOrInt"
Cohesion: 0.07
Nodes (4): ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, ProjectConnection

### Community 184 - "auth.DexValidator"
Cohesion: 0.15
Nodes (8): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice(), channelLabel(), networkCause(), parseRetryAfter(), Configuration

### Community 191 - "serveNested"
Cohesion: 0.33
Nodes (10): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newNestedStub(), serveNested(), TestNestedVariableCreateUsesTheWorkspaceInThePath(), TestNestedVariableDeleteRemovesOnlyThatVariable(), TestNestedVariableFromAnotherWorkspaceIsNotFound(), TestNestedVariablesListIsScopedToTheWorkspace() (+2 more)

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 197 - "grantsAPI"
Cohesion: 0.31
Nodes (10): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+2 more)

### Community 229 - "Schedule"
Cohesion: 0.18
Nodes (4): scheduleStub, CreateScheduleParams, Schedule, UpdateScheduleParams

### Community 231 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 235 - ".build"
Cohesion: 0.10
Nodes (16): collectionSpec, DB, failingDB, listArgs, afterString(), loadPage(), ParsePage(), TestParsePage() (+8 more)

### Community 236 - ".Query"
Cohesion: 0.07
Nodes (14): CountJobsByStatusRow, CreateArchiveParams, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsRow, ListSSHKeySecretsRow, ListVariableSecretsAfterParams (+6 more)

### Community 248 - "Deps"
Cohesion: 0.29
Nodes (3): Deps, fakeLocker, JobLocker

### Community 251 - "query.go"
Cohesion: 0.38
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 252 - "cleanKey"
Cohesion: 0.33
Nodes (4): cleanKey(), pathClean(), TestCleanKey(), S3

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 543 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **147 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Queries` to `.SetVariableSecret`, `net/http.ResponseWriter`, `IssueCliAuthCodeParams`, `.SetCliAuthSessionIdentity`, `.SetSSHKeySecret`, `.CreateGlobalVariable`, `Team`, `NewCipher`, `net/http.ResponseWriter`, `newFixture`, `.ListNotificationSecretsAfter`, `fakeQuerier`, `.ListSSHKeySecretsAfter`, `.SoftDeletePat`, `github.com/aminbenmansour/terrakube-go/internal/store.Job`, `.marshalOBoolean2bool`, `tfe_test.go`, `Agent`, `fakeStore`, `log/slog.Logger`, `.CreateNotificationTrigger`, `templateStub`, `Dispatcher`, `api/main.go`, `stubQuerier`, `.QueryRow`, `github.com/aminbenmansour/terrakube-go/internal/store.Step`, `notificationSendAPI`, `Collection`, `Metrics`, `dispatcherStore`, `Content`, `Variable`, `CliAuthSession`, `testing.T`, `Pat`, `Tag`, `Job`, `tagStub`, `Team`, `Dispatcher`, `Workspace`, `Cipher`, `rotateValue`, `Layout`, `WorkspaceRunTrigger`, `Schedule`, `context.Context`, `.Query`, `Dispatcher`, `.SetGlobalVariableSecret`, `.SetNotificationSecret`, `CreateNotificationConfigurationParams`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.0802675585284281 - nodes in this community are weakly interconnected._
- **Why does `openRotationTest()` connect `NewCipher` to `context.Context`, `newTestAPI`, `testing.T`, `Queries`, `api/main.go`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `github.com/99designs/gqlgen/graphql.CollectedField` to `TemplateConnection`, `._Project`, `WorkspaceTagConnection`, `idsFrom`, `context.Context`, `generated.go`, `.marshalOBoolean2bool`, `github.com/99designs/gqlgen/graphql.Marshaler`, `JobConnection`, `models_gen.go`, `._NotificationConfiguration`, `StringOrInt`, `.unmarshalOString2ᚕᚖstring`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Should `._Project` be split into smaller, more focused modules?**
  _Cohesion score 0.12535612535612536 - nodes in this community are weakly interconnected._