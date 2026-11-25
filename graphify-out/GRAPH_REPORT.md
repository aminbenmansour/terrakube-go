# Graph Report - terrakube-go  (2026-11-25)

## Corpus Check
- 460 files · ~227,091 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4097 nodes · 16490 edges · 190 communities (101 shown, 89 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 3084 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5de8a9de`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- .CreateGlobalVariable
- .Delete
- Core Database Schema
- tagStub
- collectionStub
- Handler
- github.com/99designs/gqlgen/graphql.Marshaler
- net/http.Request
- CliAuthSession
- .Query
- Resource
- steps
- fakeQuerier
- stubQuerier
- encryptedStore
- tfe_test.go
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- StringOrInt
- newFixture
- Identity
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- api/main.go
- JobConnection
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
- NotificationTriggerConnection
- Service
- logSink
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- .handle
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- store.Access
- .CreateWorkspaceTag
- graphql/collections.go
- WriteError
- Querier
- newTestExecutor
- stubQuerier
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- ._Project
- models
- config.go
- store.Team
- collections.go
- Cipher
- generated.go
- where_test.go
- Resolver
- templateStub
- decodeJSON
- graphql/models.go
- newJWKSServer
- ResolveOrganization
- Secret encryption at rest — design
- ui.go
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- Tag
- New
- CreateSSHKeyParams
- environment
- Backend
- NewSubscriptions
- newTestService
- .processJob
- context.Context
- AccessConnection
- ui-rest-inventory.sh
- loadPage
- ui-smoke.sh
- where.go
- grantsAPI
- capture.sh
- 0008_action.sql
- .CreateAgent
- compare.sh
- start.sh
- Workspace
- .CreateOrganization
- tarGzEntries
- .EnqueueJob
- WorkspaceTagConnection
- actionStub
- net/http.Handler
- run
- Secret key rotation — design
- cancelRunStore
- net/http.Handler
- io.Reader
- NewSubscriptions
- .Get
- reliabilityStore
- .cloneConfiguration
- .QueryRow
- Review Focus
- newService
- logSink
- ResourceIdentifier
- memStorage
- .PublishCapped
- auth.DexValidator
- Object
- cleanKey
- notificationTestAPI
- .createVCS
- workspaceQuery
- Local
- .Ping
- transientError
- ParsePage
- CreateProjectParams
- models_gen.go

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 208 edges
3. `Queries` - 206 edges
4. `API` - 162 edges
5. `stubQuerier` - 148 edges
6. `New()` - 145 edges
7. `WriteError()` - 114 edges
8. `param()` - 78 edges
9. `API` - 67 edges
10. `ParseFieldsets()` - 65 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Command` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Report and verification` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Secret Encryption at Rest Implementation Plan` --references--> `Querier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/store/querier.go
- `Task 3: Transparent store decorator` --references--> `Querier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/store/querier.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (190 total, 89 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (24): workspaceUpdateParams(), RelationshipIdentifier(), Write(), WriteResource(), configurationVersionResource(), approvalGateStatus(), failedJobStatus(), parseRunID() (+16 more)

### Community 1 - "rbac.go"
Cohesion: 0.19
Nodes (27): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, adminOnly(), CanApproveJob() (+19 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 7 - "collectionStub"
Cohesion: 0.06
Nodes (6): collectionStub, CreateCollectionParams, CreateItemParams, CreateReferenceParams, UpdateCollectionParams, UpdateItemParams

### Community 8 - "Handler"
Cohesion: 0.06
Nodes (35): Claims, Type, Validator, ClaimsFromContext(), Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI() (+27 more)

### Community 9 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.06
Nodes (3): NotificationConfigurationConnection, NotificationConfigurationEdge, WorkspaceEdge

### Community 10 - "net/http.Request"
Cohesion: 0.12
Nodes (8): API, New(), param(), parseList(), unknownCollection(), writeInUseConflict(), ParseFieldsets(), WriteCollection()

### Community 11 - "CliAuthSession"
Cohesion: 0.13
Nodes (5): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams, SetCliAuthSessionIdentityParams

### Community 12 - ".Query"
Cohesion: 0.05
Nodes (13): ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListSSHKeySecretsAfterParams, ListSSHKeySecretsAfterRow, ListSSHKeySecretsRow (+5 more)

### Community 13 - "Resource"
Cohesion: 0.10
Nodes (37): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), accessResource(), projectAccessResource() (+29 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 16 - "fakeQuerier"
Cohesion: 0.05
Nodes (9): nestedStub, fakeQuerier, CreateNotificationConfigurationParams, CreateVariableParams, CreateVCSParams, UpdateNotificationConfigurationParams, UpdateVariableParams, UpdateVCSParams (+1 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (42): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+34 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.08
Nodes (43): replay, upstreamCapture, expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner() (+35 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 24 - "newFixture"
Cohesion: 0.09
Nodes (35): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+27 more)

### Community 25 - "Identity"
Cohesion: 0.07
Nodes (20): Deps, groupVerifier, PATService, stubDexVerifier, stubIdentityVerifier, stubPATService, fakeVerifier, Identity (+12 more)

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
Nodes (58): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs, grantAttrs (+50 more)

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
Cohesion: 0.14
Nodes (8): notificationConfigurationResolver, organizationResolver, ProjectConnection, queryResolver, workspaceResolver, organizationScope(), TestOrganizationScope(), idsFrom()

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

### Community 48 - "Service"
Cohesion: 0.16
Nodes (9): Service, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), randomToken(), sanitizeName() (+1 more)

### Community 49 - "logSink"
Cohesion: 0.16
Nodes (7): logSink, redactor, TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor(), baseEnvironment()

### Community 51 - "time.Duration"
Cohesion: 0.11
Nodes (7): fakeQueue, JobQueue, envDuration(), defaultRetryDelay(), waitForProcessExit(), Message, Queue

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (121): TestAccessTokenWithoutStore(), newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+113 more)

### Community 53 - "Pat"
Cohesion: 0.17
Nodes (5): fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - ".handle"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 57 - "store.Access"
Cohesion: 0.07
Nodes (10): AccessTokenStore, fakeAccessTokenStore, contains(), captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces(), captureProjectAccessRow() (+2 more)

### Community 59 - "graphql/collections.go"
Cohesion: 0.10
Nodes (22): collectionSpec, listArgs, page, accessEdges(), afterString(), column(), jobEdges(), notificationConfigurationEdges() (+14 more)

### Community 60 - "WriteError"
Cohesion: 0.09
Nodes (39): actionAttrs, atomicOperation, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, vcsAttrs, decodeActionAttrs() (+31 more)

### Community 61 - "Querier"
Cohesion: 0.03
Nodes (4): AddNotificationConfigurationTemplateParams, DBTX, Querier, UpdateHistoryOutputParams

### Community 62 - "newTestExecutor"
Cohesion: 0.31
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 63 - "stubQuerier"
Cohesion: 0.04
Nodes (8): jobStub, statusQuerier, captureJobs(), captureSteps(), SetWorkspaceLastJobParams, UpdateJobStatusParams, stubQuerier, stubQuerier

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.08
Nodes (44): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+36 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "config.go"
Cohesion: 0.17
Nodes (20): Config, Task 1: Keyring configuration, env(), envBool(), envFirst(), envInt64(), envList(), Load() (+12 more)

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 73 - "Cipher"
Cohesion: 0.07
Nodes (35): Review Focus, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation, Task 7: Documentation and graph, Task 8: Verification, Task 3: Rotation scan and re-encryption (+27 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "Resolver"
Cohesion: 0.11
Nodes (9): projectResolver, Resolver, ResolverRoot, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace() (+1 more)

### Community 77 - "templateStub"
Cohesion: 0.17
Nodes (3): templateStub, CreateTemplateParams, UpdateTemplateParams

### Community 78 - "decodeJSON"
Cohesion: 0.20
Nodes (21): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath() (+13 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "ResolveOrganization"
Cohesion: 0.13
Nodes (9): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access, OrganizationStore (+1 more)

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

### Community 87 - "workspaces_test.go"
Cohesion: 0.16
Nodes (22): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+14 more)

### Community 88 - "Tag"
Cohesion: 0.11
Nodes (3): CreateTagParams, GetTagByNameParams, UpdateTagParams

### Community 89 - "New"
Cohesion: 0.21
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 99 - "Backend"
Cohesion: 0.07
Nodes (22): formatOldKeys(), logReport(), Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log (+14 more)

### Community 100 - "NewSubscriptions"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 102 - ".processJob"
Cohesion: 0.22
Nodes (3): stringPointer(), transient(), contentIDFromSource()

### Community 103 - "context.Context"
Cohesion: 0.03
Nodes (8): globalVarStub, sshAgentStub, teamStub, vcsStub, workspaceStub, artifactStore, historyStub, variableStore

### Community 118 - "loadPage"
Cohesion: 0.24
Nodes (3): DB, failingDB, loadPage()

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 123 - "grantsAPI"
Cohesion: 0.31
Nodes (10): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+2 more)

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 133 - "Workspace"
Cohesion: 0.08
Nodes (8): CreateWorkspaceParams, GetWorkspaceByNameParams, ListWorkspacesByProjectPageDescParams, ListWorkspacesByProjectPageParams, ListWorkspacesPageDescParams, ListWorkspacesPageParams, SetWorkspaceLockParams, UpdateWorkspaceParams

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 143 - "actionStub"
Cohesion: 0.15
Nodes (3): actionStub, CreateActionParams, UpdateActionParams

### Community 144 - "net/http.Handler"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "run"
Cohesion: 0.05
Nodes (85): listQuery, pageMeta, main(), run(), main(), run(), main(), run() (+77 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 150 - "net/http.Handler"
Cohesion: 0.10
Nodes (15): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+7 more)

### Community 154 - "NewSubscriptions"
Cohesion: 0.14
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 157 - ".cloneConfiguration"
Cohesion: 0.24
Nodes (8): gitConfig, TestSanitizeGitURL(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL(), runGit(), sanitizeGitURL()

### Community 158 - ".QueryRow"
Cohesion: 0.03
Nodes (10): CreateContentParams, UpdateContentStatusParams, CountWorkspacesByNameParams, CountWorkspacesByProjectNameParams, CreateAccessParams, CreateNotificationTriggerParams, CreateProjectAccessParams, Queries (+2 more)

### Community 159 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 160 - "newService"
Cohesion: 0.19
Nodes (15): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+7 more)

### Community 161 - "logSink"
Cohesion: 0.21
Nodes (5): Executor, preparedWorkspace, isCancelledStatus(), TestBinarySelection(), firstFolder()

### Community 162 - "ResourceIdentifier"
Cohesion: 0.22
Nodes (8): nestedDocument, RelationshipList(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship, ResourceIdentifier, tagBindingInput

### Community 163 - "memStorage"
Cohesion: 0.17
Nodes (6): newMemStorage(), Publisher, publishingStore, stubQuerier, mutableQuerier, memStorage

### Community 165 - "auth.DexValidator"
Cohesion: 0.23
Nodes (5): Dex, contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 167 - "cleanKey"
Cohesion: 0.29
Nodes (5): cleanKey(), pathClean(), TestCleanKey(), NewS3(), S3

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 169 - ".createVCS"
Cohesion: 0.53
Nodes (4): strPtr(), strValue(), validVCSEnum(), vcsEnumOrDefault()

### Community 170 - "workspaceQuery"
Cohesion: 0.33
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 174 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 181 - "CreateProjectParams"
Cohesion: 0.12
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 183 - "models_gen.go"
Cohesion: 0.07
Nodes (9): OrganizationConnection, OrganizationEdge, PageInfo, ProjectAccessConnection, ProjectAccessEdge, Query, Subscription, TemplateConnection (+1 more)

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 447 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **89 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `.CreateAgent`, `.CreateGlobalVariable`, `Workspace`, `tagStub`, `collectionStub`, `Handler`, `.CreateOrganization`, `CliAuthSession`, `.Query`, `actionStub`, `fakeQuerier`, `run`, `stubQuerier`, `tfe_test.go`, `encryptedStore`, `newFixture`, `Identity`, `.CreateArchive`, `.QueryRow`, `api/main.go`, `Access`, `testing.T`, `Pat`, `CreateProjectParams`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `store.Access`, `.CreateWorkspaceTag`, `stubQuerier`, `Cipher`, `templateStub`, `Tag`, `CreateSSHKeyParams`, `Backend`, `context.Context`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07603092783505154 - nodes in this community are weakly interconnected._
- **Why does `New()` connect `net/http.Request` to `net/http.ResponseWriter`, `NewLayout`, `notificationTestAPI`, `.createVCS`, `Resource`, `run`, `testing.T`, `strPtr`, `net/http.Handler`, `newTestExecutor`, `Identity`, `WriteError`, `api/main.go`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `stubQuerier` connect `stubQuerier` to `net/http.ResponseWriter`, `Workspace`, `tagStub`, `store.Team`, `actionStub`, `ResolveOrganization`, `encryptedStore`, `api/main.go`, `store.Access`, `Querier`, `.QueryRow`, `stubQuerier`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `Core Database Schema` be split into smaller, more focused modules?**
  _Cohesion score 0.056051587301587304 - nodes in this community are weakly interconnected._