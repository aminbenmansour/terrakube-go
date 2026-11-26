# Graph Report - terrakube-go  (2026-11-26)

## Corpus Check
- 476 files · ~236,371 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4248 nodes · 16893 edges · 238 communities (112 shown, 126 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3125 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a9e6b785`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- .CreateGlobalVariable
- memStorage
- Core Database Schema
- .QueryRow
- CreateItemParams
- Handler
- github.com/99designs/gqlgen/graphql.Marshaler
- net/http.Request
- encryptedStore
- .Query
- ResolveOrganization
- steps
- WorkspaceTagConnection
- stubQuerier
- fakeQuerier
- tfe_test.go
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- newFixture
- Service
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
- graphql/collections.go
- AccessConnection
- Service
- logSink
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- .handle
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- store.Access
- NewCipher
- WriteError
- Querier
- Backend
- stubQuerier
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- models
- config.go
- store.Team
- collections.go
- Cipher
- generated.go
- where_test.go
- streamEvents
- templateStub
- decodeJSON
- graphql/models.go
- newJWKSServer
- Parse
- Secret encryption at rest — design
- ui.go
- REST contract, captured from upstream
- pat.py
- tagStub
- time.Duration
- Job
- Identity
- environment
- NewSubscriptions
- newTestService
- .CreateProject
- context.Context
- stubStore
- ui-rest-inventory.sh
- Dex
- ui-smoke.sh
- cancelRunStore
- newWorkspaceStub
- reliabilityStore
- capture.sh
- 0008_action.sql
- CreateVCSParams
- compare.sh
- start.sh
- TemplateConnection
- ResourceIdentifier
- time.Time
- tarGzEntries
- newTestExecutor
- parser
- github.com/aminbenmansour/terrakube-go/internal/store.Organization
- clilogin/handler_test.go
- newTestAPI
- Rotate
- io.Reader
- net/http.Handler
- actionStub
- .CreateArchive
- NewSubscriptions
- Workspace
- .CreateAgent
- .prepareWorkspace
- Node
- M4 run logs and artifacts Implementation Plan
- logSink
- newService
- CliAuthSession
- ._NotificationConfiguration
- auth.DexValidator
- .updateWorkspaceRelationship
- notificationTestAPI
- CreateCollectionParams
- Metrics
- transientError
- .CreateReference
- .CountList
- Threat model
- models_gen.go
- .Get
- store.Workspace
- New
- RSQL filters in SQL — design
- cleanKey
- CreateSSHKeyParams
- Object
- Compare
- Review Focus
- Review Focus
- workspaceQuery
- Local
- .TryLock
- ParsePage
- .Ping
- .ListOrganizations

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Queries` - 205 edges
3. `Querier` - 197 edges
4. `API` - 162 edges
5. `New()` - 146 edges
6. `stubQuerier` - 142 edges
7. `WriteError()` - 122 edges
8. `param()` - 78 edges
9. `API` - 67 edges
10. `ParseFieldsets()` - 65 edges

## Surprising Connections (you probably didn't know these)
- `Self-review` --references--> `SortKey`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/store/list.go
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Task 5: Backfill implementation` --references--> `needsEncryption()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/secretstore/backfill.go
- `Command` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (238 total, 126 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (25): Write(), WriteResource(), configurationVersionResource(), projectResource(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID() (+17 more)

### Community 1 - "rbac.go"
Cohesion: 0.19
Nodes (27): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, adminOnly(), CanApproveJob() (+19 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "memStorage"
Cohesion: 0.15
Nodes (5): failingBackend, newMemStorage(), logBackend, memObject, memStorage

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - ".QueryRow"
Cohesion: 0.04
Nodes (10): CreateContentParams, UpdateContentStatusParams, CreateAccessParams, CreateProjectAccessParams, CreateTagParams, GetTagByNameParams, UpdateOrganizationParams, UpdateProjectAccessParams (+2 more)

### Community 8 - "Handler"
Cohesion: 0.08
Nodes (31): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken() (+23 more)

### Community 9 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.05
Nodes (5): JobEdge, NotificationConfigurationEdge, ProjectAccessConnection, ProjectAccessEdge, ProjectEdge

### Community 10 - "net/http.Request"
Cohesion: 0.10
Nodes (18): API, New(), param(), parseInclude(), parseList(), unknownCollection(), referenceResource(), projectAccessResource() (+10 more)

### Community 12 - ".Query"
Cohesion: 0.06
Nodes (17): CountJobsByStatusRow, DBTX, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow (+9 more)

### Community 13 - "ResolveOrganization"
Cohesion: 0.11
Nodes (10): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), contains(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access (+2 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 18 - "fakeQuerier"
Cohesion: 0.06
Nodes (10): Review Focus, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation, Task 7: Documentation and graph, Task 8: Verification, fakeQuerier (+2 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (42): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+34 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.07
Nodes (44): replay, stubIdentityVerifier, upstreamCapture, expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions() (+36 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 24 - "newFixture"
Cohesion: 0.08
Nodes (37): Storage, Task 2: Serve archived logs from `/tfoutput/v1`, decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret() (+29 more)

### Community 25 - "Service"
Cohesion: 0.16
Nodes (7): AccessTokenStore, Deps, PATService, stubPATService, stringPtr(), Service, Token

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
Nodes (58): actionFilter, agentAttrs, atomicOperation, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs (+50 more)

### Community 31 - "Link"
Cohesion: 0.08
Nodes (30): vcsAttrs, actionResource(), strPtr(), strValue(), collectionResource(), itemResource(), globalVarResource(), accessResource() (+22 more)

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
Cohesion: 0.11
Nodes (10): notificationConfigurationResolver, organizationResolver, queryResolver, ResolverRoot, workspaceResolver, organizationScope(), TestOrganizationScope(), idsFrom() (+2 more)

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

### Community 46 - "graphql/collections.go"
Cohesion: 0.12
Nodes (18): collectionSpec, listArgs, page, accessEdges(), afterString(), jobEdges(), loadPage(), notificationConfigurationEdges() (+10 more)

### Community 48 - "Service"
Cohesion: 0.13
Nodes (12): Service, Store, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value() (+4 more)

### Community 49 - "logSink"
Cohesion: 0.19
Nodes (7): logSink, redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 51 - "time.Duration"
Cohesion: 0.09
Nodes (7): captureLogPublisher, fakeQueue, JobQueue, LogPublisher, Queue, JobQueue, stubJobQueue

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (103): TestAccessTokenWithoutStore(), TestAuthMiddleware(), TestCreateOrganization(), authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound() (+95 more)

### Community 53 - "Pat"
Cohesion: 0.17
Nodes (5): fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 56 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.10
Nodes (3): jobStub, statusQuerier, captureJobs()

### Community 59 - "NewCipher"
Cohesion: 0.06
Nodes (67): main(), run(), main(), run(), main(), run(), main(), run() (+59 more)

### Community 60 - "WriteError"
Cohesion: 0.10
Nodes (34): actionAttrs, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, decodeActionAttrs(), stringOr(), boolValue() (+26 more)

### Community 61 - "Querier"
Cohesion: 0.03
Nodes (7): AddNotificationConfigurationTemplateParams, CreateNotificationTriggerParams, Querier, Queries, SetWorkspaceLastJobParams, UpdateJobOutputParams, UpdateStepOutputParams

### Community 62 - "Backend"
Cohesion: 0.19
Nodes (10): Validator, formatOldKeys(), logReport(), Deps, Deps, Backend, Layout, Deps (+2 more)

### Community 63 - "stubQuerier"
Cohesion: 0.06
Nodes (5): captureSteps(), fakeWorkspaceStore, stubQuerier, workspaceWithTags, stubQuerier

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.07
Nodes (48): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+40 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "config.go"
Cohesion: 0.17
Nodes (20): Config, Task 1: Keyring configuration, env(), envBool(), envFirst(), envInt64(), envList(), Load() (+12 more)

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 73 - "Cipher"
Cohesion: 0.11
Nodes (26): Task 3: Rotation scan and re-encryption, backfillGlobalVariables(), backfillNotifications(), backfillSSHKeys(), backfillVariables(), backfillVCS(), encryptIfNeeded(), needsEncryption() (+18 more)

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.19
Nodes (16): column(), predicate(), Field, likePattern(), mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode() (+8 more)

### Community 76 - "streamEvents"
Cohesion: 0.07
Nodes (15): DB, Deps, EventSource, failingDB, projectResolver, Resolver, subscriptionResolver, TestOwnerRule() (+7 more)

### Community 77 - "templateStub"
Cohesion: 0.17
Nodes (3): templateStub, CreateTemplateParams, UpdateTemplateParams

### Community 78 - "decodeJSON"
Cohesion: 0.12
Nodes (31): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent() (+23 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "Parse"
Cohesion: 0.21
Nodes (13): Parse(), dump(), TestParse(), TestParseRejects(), ptr(), seedWorkspaces(), TestListPageBetween(), TestListPageFilters() (+5 more)

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

### Community 89 - "time.Duration"
Cohesion: 0.16
Nodes (5): envDuration(), defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), Message

### Community 91 - "Job"
Cohesion: 0.11
Nodes (5): ApproveJobParams, CreateJobParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams, UpdateJobStatusParams

### Community 94 - "Identity"
Cohesion: 0.17
Nodes (8): groupVerifier, stubDexVerifier, fakeVerifier, Identity, fakeIssuer, groupVerifier, stubDexVerifier, stringSlice()

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "NewSubscriptions"
Cohesion: 0.33
Nodes (3): fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (20): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+12 more)

### Community 102 - ".CreateProject"
Cohesion: 0.18
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 103 - "context.Context"
Cohesion: 0.03
Nodes (9): collectionStub, globalVarStub, sshAgentStub, teamStub, vcsStub, workspaceStub, artifactStore, historyStub (+1 more)

### Community 108 - "stubStore"
Cohesion: 0.18
Nodes (3): GetArchiveRow, UpdateHistoryOutputParams, stubStore

### Community 118 - "Dex"
Cohesion: 0.21
Nodes (3): Dex, Exchanger, fakeDex

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "cancelRunStore"
Cohesion: 0.12
Nodes (5): cancelRunStore, Publisher, publishingStore, stubQuerier, mutableQuerier

### Community 121 - "newWorkspaceStub"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 123 - "reliabilityStore"
Cohesion: 0.07
Nodes (6): reliabilityStore, CreateHistoryParams, CreateStepParams, GetHistoryByJobReferenceParams, GetStepByJobAndNumberParams, UpdateStepStatusParams

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 135 - "ResourceIdentifier"
Cohesion: 0.20
Nodes (10): nestedDocument, RelationshipIdentifiers(), RelationshipList(), TestRelationshipIdentifiers(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship (+2 more)

### Community 137 - "time.Time"
Cohesion: 0.22
Nodes (10): bucket, limiter, key(), newLimiter(), RateLimit(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), TestRateLimitDisabledByDefault() (+2 more)

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 141 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.12
Nodes (19): CookieSigner, Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), authorize() (+11 more)

### Community 145 - "newTestAPI"
Cohesion: 0.10
Nodes (36): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+28 more)

### Community 146 - "Rotate"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 150 - "net/http.Handler"
Cohesion: 0.09
Nodes (18): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+10 more)

### Community 151 - "actionStub"
Cohesion: 0.15
Nodes (3): actionStub, CreateActionParams, UpdateActionParams

### Community 154 - "NewSubscriptions"
Cohesion: 0.15
Nodes (18): Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestSubscriptionAcceptsADexConnectionToken(), NewExecutableSchema(), TestSchemaLoads(), NewSubscriptions() (+10 more)

### Community 155 - "Workspace"
Cohesion: 0.14
Nodes (3): CreateWorkspaceParams, GetWorkspaceByNameParams, SetWorkspaceLockParams

### Community 157 - ".prepareWorkspace"
Cohesion: 0.12
Nodes (13): gitConfig, preparedWorkspace, TestGitSourceDetection(), TestSanitizeGitURL(), TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv() (+5 more)

### Community 159 - "Node"
Cohesion: 0.11
Nodes (3): stubLists, Node, SortKey

### Community 160 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.25
Nodes (7): Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph

### Community 161 - "logSink"
Cohesion: 0.14
Nodes (7): Executor, isCancelledStatus(), isTerminalStatus(), stringPointer(), TestBinarySelection(), TestIsTerminalStatus(), baseEnvironment()

### Community 162 - "newService"
Cohesion: 0.19
Nodes (15): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+7 more)

### Community 163 - "CliAuthSession"
Cohesion: 0.11
Nodes (5): fakeStore, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams, SetCliAuthSessionIdentityParams

### Community 165 - "auth.DexValidator"
Cohesion: 0.31
Nodes (4): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 172 - "Metrics"
Cohesion: 0.13
Nodes (4): metricsSource, fakeStats, Metrics, StatsSource

### Community 175 - ".CountList"
Cohesion: 0.21
Nodes (8): exists(), likePattern(), relate(), Compare, Op, Attr, Kind, ListSpec

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - "models_gen.go"
Cohesion: 0.07
Nodes (11): ComplexityRoot, JobConnection, NotificationConfigurationConnection, OrganizationConnection, OrganizationEdge, PageInfo, ProjectConnection, Query (+3 more)

### Community 190 - "store.Workspace"
Cohesion: 0.20
Nodes (6): captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces(), captureProjectAccessRow(), captureWorkspaces()

### Community 191 - "New"
Cohesion: 0.29
Nodes (6): requireMetricsToken(), SecurityHeaders(), TestSecurityHeaders(), New(), requestLogger(), writeJSON()

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 198 - "cleanKey"
Cohesion: 0.25
Nodes (5): cleanKey(), pathClean(), TestCleanKey(), NewS3(), S3

### Community 230 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 231 - "Review Focus"
Cohesion: 0.17
Nodes (10): listQuery, Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 3: Handler cutover (+2 more)

### Community 232 - "workspaceQuery"
Cohesion: 0.38
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 235 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 495 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **126 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `CreateVCSParams`, `.CreateGlobalVariable`, `.QueryRow`, `CreateItemParams`, `Handler`, `encryptedStore`, `.Query`, `ResolveOrganization`, `github.com/aminbenmansour/terrakube-go/internal/store.Organization`, `.marshalOString2ᚖstring`, `stubQuerier`, `fakeQuerier`, `newTestAPI`, `tfe_test.go`, `actionStub`, `newFixture`, `Service`, `.CreateArchive`, `Workspace`, `.CreateAgent`, `api/main.go`, `CliAuthSession`, `CreateCollectionParams`, `Metrics`, `.CreateReference`, `Pat`, `store.Access`, `NewCipher`, `Backend`, `stubQuerier`, `github.com/99designs/gqlgen/graphql.CollectedField`, `Cipher`, `CreateSSHKeyParams`, `templateStub`, `tagStub`, `Job`, `.CreateProject`, `context.Context`, `stubStore`, `reliabilityStore`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.08205128205128205 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `github.com/99designs/gqlgen/graphql.CollectedField` to `.___Directive`, `.childFields___InputValue`, `._NotificationConfiguration`, `TemplateConnection`, `context.Context`, `._Organization`, `github.com/99designs/gqlgen/graphql.Marshaler`, `generated.go`, `idsFrom`, `.childFields___Type`, `AccessConnection`, `.marshalOString2ᚖstring`, `WorkspaceTagConnection`, `.___EnumValue`, `models_gen.go`, `.___Field`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `New()` connect `net/http.Request` to `NewLayout`, `.updateWorkspaceRelationship`, `notificationTestAPI`, `newTestExecutor`, `newTestAPI`, `testing.T`, `strPtr`, `net/http.Handler`, `Service`, `NewCipher`, `WriteError`, `api/main.go`, `Link`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `Core Database Schema` be split into smaller, more focused modules?**
  _Cohesion score 0.056051587301587304 - nodes in this community are weakly interconnected._