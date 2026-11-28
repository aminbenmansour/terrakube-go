# Graph Report - terrakube-go  (2026-11-28)

## Corpus Check
- 618 files · ~311,515 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 5403 nodes · 22839 edges · 271 communities (135 shown, 136 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 4129 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `dc01facf`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- .internalError
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
- serveAccessToken
- Upstream Elide entities (GraphQL surface)
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- Variable
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
- NotificationTriggerConnection
- CliAuthSession
- stubQuerier
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- ._NotificationConfiguration
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
- generated.go
- where_test.go
- Resolver
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
- .CreateTag
- net/http.ResponseWriter
- environment
- cron.go
- NewSubscriptions
- newTestService
- context.Context
- newService
- ui-rest-inventory.sh
- ui-smoke.sh
- runTriggerAPI
- Workspacetag
- capture.sh
- 0008_action.sql
- models_gen.go
- Service
- compare.sh
- start.sh
- time.Duration
- .EnqueueJob
- time.Time
- New
- Identity
- Parse
- time.Time
- clilogin/handler_test.go
- newTestAPI
- Review Focus
- validateEndpoint
- ui.go
- New
- WebhookEvent
- .CreateGlobalVariable
- NewSubscriptions
- nestedStub
- stubQuerier
- Executor
- runTriggerAPI
- vcsDiscoveryAPI
- tarGzEntries
- webhook_test.go
- where.go
- API
- notificationSendAPI
- github.com/jackc/pgx/v5/pgxpool.Pool
- Review Focus
- Metrics
- transientError
- Dispatcher
- Threat model
- .cloneConfiguration
- auth.DexValidator
- Review Focus
- loadPage
- RSQL filters in SQL — design
- Upstream UI GraphQL operations
- Workspace
- .CreateAgent
- ParsePage
- Repository
- TeamToken
- Secret key rotation — design
- Team
- .build
- .Query
- RelationshipID
- .QueryRow
- Team
- JobConnection
- newService
- grantsAPI
- .CreateVCS
- New
- discoveryStub
- .server
- newScheduleAPI
- newTestExecutor
- net/http/httptest.ResponseRecorder
- logSink

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 309 edges
3. `Queries` - 301 edges
4. `API` - 294 edges
5. `New()` - 225 edges
6. `WriteError()` - 188 edges
7. `stubQuerier` - 171 edges
8. `param()` - 107 edges
9. `ParseFieldsets()` - 103 edges
10. `WriteResource()` - 81 edges

## Surprising Connections (you probably didn't know these)
- `Where we stand` --references--> `provider`  [INFERRED]
  docs/upstream/ui-rest-operations.md → internal/vcs/vcs.go
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Command` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Report and verification` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (271 total, 136 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (31): configurationVersionResource(), projectResource(), parseWorkspaceQuery(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource() (+23 more)

### Community 1 - ".internalError"
Cohesion: 0.07
Nodes (30): actionAttrs, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, vcsAttrs, decodeActionAttrs(), stringOr() (+22 more)

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
Cohesion: 0.09
Nodes (48): main(), run(), main(), run(), main(), run(), main(), run() (+40 more)

### Community 8 - "Handler"
Cohesion: 0.05
Nodes (43): Claims, Type, Validator, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log (+35 more)

### Community 9 - "newFixture"
Cohesion: 0.09
Nodes (35): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+27 more)

### Community 10 - "New"
Cohesion: 0.10
Nodes (25): globToRegex(), matchBranch(), matchPaths(), matchTrigger(), pathRegex(), branchFromRef(), deliveryID(), nonEmptyPtr() (+17 more)

### Community 12 - "policyStub"
Cohesion: 0.05
Nodes (13): policyStub, CreatePolicyAttachmentParams, CreatePolicyExemptionParams, CreatePolicySetParameterParams, CreatePolicySetParams, PolicyAttachment, PolicyExemption, PolicySet (+5 more)

### Community 13 - "logSink"
Cohesion: 0.22
Nodes (4): Executor, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 16 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.05
Nodes (3): ProjectConnection, ProjectEdge, WorkspaceEdge

### Community 17 - "github.com/aminbenmansour/terrakube-go/internal/rsql.Node"
Cohesion: 0.04
Nodes (12): listQuery, stubLists, Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs (+4 more)

### Community 18 - "rbac.go"
Cohesion: 0.26
Nodes (22): Task 1: Role matrix and permission set, adminOnly(), CanApproveJob(), CanDestroyJob(), CanManageCollection(), CanManageJob(), CanManageModule(), CanManagePolicies() (+14 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (41): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue() (+33 more)

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
Cohesion: 0.08
Nodes (10): strPtr(), Dispatcher, dispatcherStore, FailStuckNotificationOutboxParams, FailStuckRunTriggerEventsParams, NotificationOutbox, ReclaimStuckNotificationOutboxParams, ReclaimStuckRunTriggerEventsParams (+2 more)

### Community 24 - "Variable"
Cohesion: 0.10
Nodes (3): Task 2: Cipher and envelope, Task 3: Rotation scan and re-encryption, Variable

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
Nodes (72): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, createTeamTokenRequest, discoveryGroup (+64 more)

### Community 31 - "Link"
Cohesion: 0.09
Nodes (35): policySetChildren, parseInclude(), collectionResource(), referenceResource(), jobResource(), stepResource(), timestamp(), moduleResource() (+27 more)

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

### Community 44 - "fakeStore"
Cohesion: 0.04
Nodes (32): captureJobs(), captureSteps(), strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob() (+24 more)

### Community 45 - "History"
Cohesion: 0.05
Nodes (10): cancelRunStore, reliabilityStore, CreateHistoryParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams, History, UpdateJobOutputParams, UpdateStepOutputParams (+2 more)

### Community 46 - "ResolveOrganization"
Cohesion: 0.16
Nodes (10): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access, fakeOrganizationStore (+2 more)

### Community 48 - "CliAuthSession"
Cohesion: 0.10
Nodes (6): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams, SetCliAuthSessionIdentityParams

### Community 49 - "stubQuerier"
Cohesion: 0.10
Nodes (6): CreateContentParams, UpdateContentStatusParams, fakeWorkspaceStore, Content, stubQuerier, workspaceWithTags

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (108): expectPermissions(), serveAccessToken(), serveTeamTokens(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups() (+100 more)

### Community 53 - "Pat"
Cohesion: 0.17
Nodes (5): fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 57 - "newWorkspaceStub"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 59 - "._Organization"
Cohesion: 0.07
Nodes (4): ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, WorkspaceConnection

### Community 61 - "Queries"
Cohesion: 0.02
Nodes (18): replay, Rotate(), assertSealedUnder(), createRawVariable(), createRawVariableWithID(), openRotationTest(), rawVariable(), rotationWorkspace() (+10 more)

### Community 62 - "dispatcherWith"
Cohesion: 0.36
Nodes (8): dispatcherWith(), newDispatcherFixture(), TestDispatcherFailsAfterMaxAttempts(), TestDispatcherFailsPermanentRejection(), TestDispatcherIncludesTheFailureSummary(), TestDispatcherRetriesTransientFailures(), TestDispatcherSendsAndRecordsSuccess(), TestDispatcherSkipsSupersededJobStatus()

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.06
Nodes (53): redactor, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+45 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "Load"
Cohesion: 0.15
Nodes (23): Config, S3, Storage, Task 1: Keyring configuration, env(), envBool(), envDuration(), envFirst() (+15 more)

### Community 70 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.07
Nodes (3): jobStub, statusQuerier, stubQuerier

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - "New"
Cohesion: 0.21
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "Resolver"
Cohesion: 0.10
Nodes (9): projectResolver, Resolver, ResolverRoot, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace() (+1 more)

### Community 77 - "vcs/vcs.go"
Cohesion: 0.20
Nodes (21): githubRoutes(), TestBranchesGitHub(), TestGitHubConnectionsWithoutAToken(), TestGroupsGitHub(), TestRepositoriesGitHub(), TestTestConnection(), TestGitLabConnectionWithoutAToken(), TestGitLabDiscovery() (+13 more)

### Community 78 - "decodeJSON"
Cohesion: 0.14
Nodes (30): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent() (+22 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "auth/dex_test.go"
Cohesion: 0.17
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "Deps"
Cohesion: 0.18
Nodes (8): Deps, PATService, stubTeamTokens, TeamTokenService, timePtr(), timePtr(), Token, Token

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - "Cipher"
Cohesion: 0.11
Nodes (25): Task 5: Backfill implementation, backfillGlobalVariables(), backfillNotifications(), backfillSSHKeys(), backfillVariables(), backfillVCS(), encryptIfNeeded(), needsEncryption() (+17 more)

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
Cohesion: 0.26
Nodes (8): doJSON(), failed(), providerError(), providerMessage(), githubOwnerAndRepo(), Error, gitlab, WebhookSpec

### Community 94 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (35): atomicOperation, implementationAttrs, moduleAttrs, policySetAttrs, policySetParameterAttrs, providerAttrs, providerVersionAttrs, callerEmail() (+27 more)

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
Cohesion: 0.02
Nodes (17): collectionStub, globalVarStub, sshAgentStub, teamStub, templateStub, vcsStub, workspaceStub, artifactStore (+9 more)

### Community 108 - "newService"
Cohesion: 0.19
Nodes (15): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+7 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.17
Nodes (4): scheduleStub, CreateScheduleParams, Schedule, UpdateScheduleParams

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 127 - "models_gen.go"
Cohesion: 0.08
Nodes (9): AccessConnection, AccessEdge, OrganizationConnection, OrganizationEdge, PageInfo, Query, Subscription, WorkspaceTagConnection (+1 more)

### Community 128 - "Service"
Cohesion: 0.15
Nodes (10): Service, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), randomToken() (+2 more)

### Community 133 - "time.Duration"
Cohesion: 0.15
Nodes (6): waitFor(), defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), backoff(), Message

### Community 135 - ".EnqueueJob"
Cohesion: 0.22
Nodes (7): fakeQueue, JobQueue, fakeQueue, JobQueue, JobQueue, stubJobQueue, stubQueue

### Community 137 - "time.Time"
Cohesion: 0.32
Nodes (3): backoff(), TestBackoffGrowsAndCaps(), Dispatcher

### Community 140 - "New"
Cohesion: 0.15
Nodes (14): main(), run(), NewMetrics(), requireMetricsToken(), TestMetricsDisabledWithoutToken(), TestMetricsRecordsRequests(), TestMetricsRequiresToken(), TestSampleRefreshesGauges() (+6 more)

### Community 141 - "Identity"
Cohesion: 0.10
Nodes (14): groupVerifier, stubDexVerifier, stubIdentityVerifier, stubPATService, fakeVerifier, Identity, fakeIssuer, groupVerifier (+6 more)

### Community 142 - "Parse"
Cohesion: 0.10
Nodes (21): Parse(), dump(), TestParse(), TestParseRejects(), exists(), likePattern(), relate(), ptr() (+13 more)

### Community 143 - "time.Time"
Cohesion: 0.22
Nodes (10): policyExemptionAttrs, bucket, limiter, key(), newLimiter(), RateLimit(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets() (+2 more)

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.06
Nodes (49): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+41 more)

### Community 146 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 147 - "validateEndpoint"
Cohesion: 0.50
Nodes (4): blockedAddress(), TestBlockedAddressRanges(), TestValidateEndpointRefusesInternalAddresses(), validateEndpoint()

### Community 149 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 150 - "New"
Cohesion: 0.07
Nodes (16): actionResource(), New(), param(), parseList(), unknownCollection(), itemResource(), globalVarResource(), accessResource() (+8 more)

### Community 154 - "NewSubscriptions"
Cohesion: 0.14
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 159 - "stubQuerier"
Cohesion: 0.02
Nodes (20): actionStub, missingOrganizationStub, stubQuerier, scopeIDs(), demoOrganizations(), CreateActionParams, UpdateActionParams, CreateCollectionParams (+12 more)

### Community 161 - "Executor"
Cohesion: 0.16
Nodes (5): isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestIsTerminalStatus()

### Community 162 - "runTriggerAPI"
Cohesion: 0.23
Nodes (12): runTriggerStub, newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit(), TestCreateRunTriggerRejectsCycle(), TestCreateRunTriggerRejectsStructuralViolations() (+4 more)

### Community 163 - "vcsDiscoveryAPI"
Cohesion: 0.46
Nodes (8): decodeBody(), fakeGitHub(), TestVCSDiscoveryDeniedWithoutPermission(), TestVCSDiscoveryGroups(), TestVCSDiscoveryRejectsBadRequests(), TestVCSDiscoveryRepositoriesAndBranches(), TestVCSDiscoveryTestConnection(), vcsDiscoveryAPI()

### Community 164 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 165 - "webhook_test.go"
Cohesion: 0.36
Nodes (10): DeleteWebhook(), SyncWebhook(), newHookServer(), TestDeleteWebhookTreats404AsRemoved(), TestGitHubWebhookCreate(), TestGitHubWebhookUpdateAndDelete(), TestGitLabWebhookCreate(), TestGitLabWebhookDelete() (+2 more)

### Community 167 - "API"
Cohesion: 0.05
Nodes (37): API, authorization, webhookEventAttrs, webhookSyncRef, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams() (+29 more)

### Community 168 - "notificationSendAPI"
Cohesion: 0.14
Nodes (23): newDeliveryStub(), TestNotificationDeliveryHistory(), TestNotificationDeliveryHistoryRejectsABadLimit(), TestNotificationDeliveryHistoryUnknownWorkspace(), TestNotificationDeliveryRetry(), TestNotificationDeliveryRetryRejections(), captureJSONServer(), notificationSendAPI() (+15 more)

### Community 169 - "github.com/jackc/pgx/v5/pgxpool.Pool"
Cohesion: 0.53
Nodes (5): Locker, NewLocker(), TestLockerAllowsDifferentJobs(), TestLockerExcludesASecondHolder(), testPool()

### Community 171 - "Review Focus"
Cohesion: 0.33
Nodes (5): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract

### Community 172 - "Metrics"
Cohesion: 0.08
Nodes (13): metricsSource, startScheduler(), formatOldKeys(), logReport(), Deps, fakePinger, fakeStats, Metrics (+5 more)

### Community 175 - "Dispatcher"
Cohesion: 0.28
Nodes (8): strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued(), TestTickSkipsLockedWorkspaceAndUnparseableCron(), TestTickUsesTheTemplateOfTheWorkspaceOrganization(), tickAt(), Dispatcher

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - ".cloneConfiguration"
Cohesion: 0.13
Nodes (12): gitConfig, preparedWorkspace, TestSanitizeGitURL(), transient(), contentIDFromSource(), firstFolder(), gitEnv(), gitSource() (+4 more)

### Community 184 - "auth.DexValidator"
Cohesion: 0.13
Nodes (7): Dex, Exchanger, fakeDex, contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 190 - "Review Focus"
Cohesion: 0.40
Nodes (5): Review Focus, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 7: Documentation and graph, Task 8: Verification

### Community 191 - "loadPage"
Cohesion: 0.24
Nodes (3): DB, failingDB, loadPage()

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
Cohesion: 0.14
Nodes (8): decodeSecret(), New(), CreateTeamTokenParams, SoftDeleteTeamTokenParams, TeamToken, fakeStore, Service, Store

### Community 232 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 234 - "Team"
Cohesion: 0.10
Nodes (11): AccessTokenStore, fakeAccessTokenStore, contains(), captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces(), captureProjectAccessRow() (+3 more)

### Community 235 - ".build"
Cohesion: 0.10
Nodes (22): collectionSpec, listArgs, page, accessEdges(), afterString(), column(), jobEdges(), notificationConfigurationEdges() (+14 more)

### Community 236 - ".Query"
Cohesion: 0.05
Nodes (19): CountJobsByStatusRow, CreateArchiveParams, DBTX, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow (+11 more)

### Community 248 - "RelationshipID"
Cohesion: 0.15
Nodes (14): nestedDocument, DecodeAttributes(), RelationshipID(), RelationshipIdentifiers(), RelationshipList(), TestDecodeRequest(), TestRelationshipIdentifiers(), TestRelationshipIDNull() (+6 more)

### Community 255 - ".QueryRow"
Cohesion: 0.03
Nodes (4): CreateProjectParams, GetProjectByNameParams, Ssh, UpdateProjectParams

### Community 259 - "newService"
Cohesion: 0.25
Nodes (10): stringSlice(), dexUser(), newFakeStore(), newService(), TestCreateRequiresMembershipAndAnInteractiveUser(), TestCreateSignsTheUpstreamToken(), TestListAndDeleteAreScopedToTheCallersTeams(), testSecret() (+2 more)

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

### Community 276 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 277 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.10
Nodes (25): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), policyServe(), policyTestStore() (+17 more)

### Community 280 - "logSink"
Cohesion: 0.20
Nodes (3): captureLogPublisher, LogPublisher, logSink

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 591 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **136 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `API` connect `API` to `.internalError`, `Deps`, `New`, `api/main.go`, `RelationshipID`, `net/http.ResponseWriter`, `Link`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.06980240549828179 - nodes in this community are weakly interconnected._
- **Why does `Querier` connect `Queries` to `Team`, `net/http.ResponseWriter`, `registryStub`, `run`, `Handler`, `newFixture`, `time.Time`, `fakeQuerier`, `policyStub`, `New`, `.CreateVCS`, `.marshalOBoolean2bool`, `newTestAPI`, `github.com/aminbenmansour/terrakube-go/internal/rsql.Node`, `tfe_test.go`, `net/http/httptest.ResponseRecorder`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `.CreateGlobalVariable`, `Variable`, `WebhookEvent`, `nestedStub`, `api/main.go`, `stubQuerier`, `notificationSendAPI`, `Metrics`, `fakeStore`, `History`, `Dispatcher`, `CliAuthSession`, `stubQuerier`, `NotificationTriggerConnection`, `ResolveOrganization`, `Pat`, `Review Focus`, `tagStub`, `github.com/aminbenmansour/terrakube-go/internal/store.Job`, `Workspace`, `.CreateAgent`, `Deps`, `Cipher`, `.CreateTag`, `context.Context`, `TeamToken`, `.Query`, `runTriggerAPI`, `Workspacetag`, `.QueryRow`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **Should `.internalError` be split into smaller, more focused modules?**
  _Cohesion score 0.07014157014157014 - nodes in this community are weakly interconnected._
- **Why does `New()` connect `New` to `.internalError`, `runTriggerAPI`, `NewLayout`, `API`, `run`, `notificationSendAPI`, `Deps`, `newTestAPI`, `newScheduleAPI`, `testing.T`, `net/http/httptest.ResponseRecorder`, `serveAccessToken`, `api/main.go`, `.server`, `stubQuerier`, `newTestExecutor`, `net/http.ResponseWriter`, `Link`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._