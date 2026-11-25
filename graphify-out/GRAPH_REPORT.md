# Graph Report - terrakube-go  (2026-11-25)

## Corpus Check
- 452 files · ~221,469 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 3984 nodes · 16065 edges · 195 communities (112 shown, 83 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 3027 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ead4c833`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- .CreateGlobalVariable
- io.Reader
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
- Service
- fakeQuerier
- stubQuerier
- encryptedStore
- tfe_test.go
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- newFixture
- Identity
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- api/main.go
- .handle
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
- Migrate
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- logSink
- Claims
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
- github.com/aminbenmansour/terrakube-go/internal/store.Project
- Secret encryption at rest — design
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- CreateTagParams
- nestedStub
- authzAPI
- Queries
- ._WorkspaceTag
- environment
- M4 run logs and artifacts Implementation Plan
- NewSubscriptions
- newTestService
- ProjectAccess
- context.Context
- workspaceStub
- ui-rest-inventory.sh
- workspaceTagBindings
- ui-smoke.sh
- newTestAPI
- New
- capture.sh
- 0008_action.sql
- stubStore
- compare.sh
- start.sh
- CreateOrganizationParams
- tarGzEntries
- log/slog.Logger
- Access
- actionStub
- net/http.Handler
- NewCipher
- Secret key rotation — design
- .CreateAgent
- net/http.Handler
- time.Time
- .CreateNotificationTrigger
- NewSubscriptions
- CreateCollectionParams
- reliabilityStore
- .cloneConfiguration
- .QueryRow
- OrganizationConnection
- ui.go
- logSink
- workspaceQuery
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- New
- auth.DexValidator
- github.com/aminbenmansour/terrakube-go/internal/store.Tag
- notificationTestAPI
- .streamStepLog
- serveAccessToken
- log/slog.Logger
- ParsePage
- CreateProjectParams
- StringOrInt
- Deps
- .SetCliAuthSessionIdentity
- serveGlobalVar
- Review Focus
- .updateWorkspaceRelationship
- get
- transientError

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 193 edges
3. `Queries` - 190 edges
4. `API` - 162 edges
5. `New()` - 145 edges
6. `stubQuerier` - 131 edges
7. `WriteError()` - 114 edges
8. `param()` - 78 edges
9. `API` - 67 edges
10. `ParseFieldsets()` - 65 edges

## Surprising Connections (you probably didn't know these)
- `Command` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Report and verification` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Task 5: Backfill implementation` --references--> `needsEncryption()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/secretstore/backfill.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Task 3: Rotation scan and re-encryption` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-key-rotation.md → internal/secretstore/rotate.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (195 total, 83 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (26): workspaceUpdateParams(), Write(), WriteResource(), configurationVersionResource(), approvalGateStatus(), failedJobStatus(), parseRunID(), stepStatusTFE() (+18 more)

### Community 1 - "rbac.go"
Cohesion: 0.19
Nodes (27): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, adminOnly(), CanApproveJob() (+19 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (11): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+3 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "tagStub"
Cohesion: 0.10
Nodes (3): tagStub, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams

### Community 7 - "collectionStub"
Cohesion: 0.07
Nodes (4): collectionStub, CreateItemParams, CreateReferenceParams, UpdateItemParams

### Community 8 - "Handler"
Cohesion: 0.24
Nodes (6): isExternalUser(), matchesMD5(), New(), newUUID(), pathIDs(), Handler

### Community 9 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.04
Nodes (13): AccessConnection, AccessEdge, JobConnection, JobEdge, NotificationTriggerConnection, NotificationTriggerEdge, PageInfo, ProjectConnection (+5 more)

### Community 10 - "net/http.Request"
Cohesion: 0.11
Nodes (10): API, listQuery, filterSortPage(), New(), param(), parseList(), unknownCollection(), Page (+2 more)

### Community 11 - "CliAuthSession"
Cohesion: 0.20
Nodes (4): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams

### Community 12 - ".Query"
Cohesion: 0.06
Nodes (15): ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow, ListSSHKeySecretsAfterParams, ListSSHKeySecretsAfterRow (+7 more)

### Community 13 - "Resource"
Cohesion: 0.10
Nodes (37): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), accessResource(), projectAccessResource() (+29 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 15 - "Service"
Cohesion: 0.16
Nodes (9): Service, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), randomToken(), sanitizeName() (+1 more)

### Community 16 - "fakeQuerier"
Cohesion: 0.08
Nodes (10): fakeQuerier, CreateNotificationConfigurationParams, CreateSSHKeyParams, CreateVariableParams, CreateVCSParams, UpdateNotificationConfigurationParams, UpdateSSHKeyParams, UpdateVariableParams (+2 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (42): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+34 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.10
Nodes (40): replay, upstreamCapture, sshTestAPI(), TestAgentsJSONAPI(), TestSSHKeysJSONAPI(), TestWorkspaceLinksAnSSHKeyOrAgent(), newTeamStub(), TestTeamJSONAPI() (+32 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 24 - "newFixture"
Cohesion: 0.14
Nodes (23): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+15 more)

### Community 25 - "Identity"
Cohesion: 0.08
Nodes (15): groupVerifier, PATService, stubDexVerifier, stubIdentityVerifier, fakeVerifier, Identity, Exchanger, fakeDex (+7 more)

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
Nodes (61): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs, grantAttrs (+53 more)

### Community 31 - ".handle"
Cohesion: 0.18
Nodes (4): isTransient(), parseAttempt(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

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

### Community 48 - "Migrate"
Cohesion: 0.12
Nodes (18): main(), run(), Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope, Task 3: Transparent store decorator (+10 more)

### Community 51 - "time.Duration"
Cohesion: 0.11
Nodes (6): fakeQueue, JobQueue, envDuration(), defaultRetryDelay(), Message, Queue

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (97): TestAccessTokenWithoutStore(), TestAuthMiddleware(), TestCreateOrganization(), newCollectionAPI(), TestCollectionListCarriesItsItemsAndReferences(), TestCreateCollectionDefaultsPriorityToZero(), TestCreateCollectionItemStoresTheValue(), TestCreateCollectionItemValidatesKeyValueAndCategory() (+89 more)

### Community 53 - "Pat"
Cohesion: 0.16
Nodes (6): newFakeStore(), fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "logSink"
Cohesion: 0.14
Nodes (7): Executor, isCancelledStatus(), isTerminalStatus(), stringPointer(), TestBinarySelection(), TestIsTerminalStatus(), baseEnvironment()

### Community 58 - "Claims"
Cohesion: 0.12
Nodes (17): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken() (+9 more)

### Community 59 - "graphql/collections.go"
Cohesion: 0.11
Nodes (20): collectionSpec, listArgs, page, accessEdges(), afterString(), jobEdges(), loadPage(), notificationConfigurationEdges() (+12 more)

### Community 60 - "WriteError"
Cohesion: 0.09
Nodes (41): actionAttrs, atomicOperation, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, vcsAttrs, decodeActionAttrs() (+33 more)

### Community 61 - "Querier"
Cohesion: 0.03
Nodes (7): AddNotificationConfigurationTemplateParams, Querier, SetGlobalVariableSecretParams, SetNotificationSecretParams, SetSSHKeySecretParams, SetVariableSecretParams, SetVCSSecretsParams

### Community 62 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.06
Nodes (55): redactor, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+47 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "config.go"
Cohesion: 0.18
Nodes (20): Config, Task 1: Keyring configuration, env(), envBool(), envFirst(), envInt64(), envList(), Load() (+12 more)

### Community 70 - "store.Team"
Cohesion: 0.18
Nodes (4): contains(), contains(), fakeOrganizationStore, ListTeamsByOrganizationAndNamesParams

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 73 - "Cipher"
Cohesion: 0.21
Nodes (20): Task 3: Rotation scan and re-encryption, backfillGlobalVariables(), backfillNotifications(), backfillSSHKeys(), backfillVariables(), backfillVCS(), encryptIfNeeded(), needsEncryption() (+12 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.19
Nodes (15): column(), predicate(), Field, likePattern(), mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode() (+7 more)

### Community 76 - "Resolver"
Cohesion: 0.09
Nodes (10): DB, failingDB, projectResolver, Resolver, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient() (+2 more)

### Community 77 - "templateStub"
Cohesion: 0.17
Nodes (3): templateStub, CreateTemplateParams, UpdateTemplateParams

### Community 78 - "decodeJSON"
Cohesion: 0.21
Nodes (21): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound(), TestJobGetMatchesTheOrganizationInThePath() (+13 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "github.com/aminbenmansour/terrakube-go/internal/store.Project"
Cohesion: 0.13
Nodes (11): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveOrganization(), TestResolveOrganizationUsesTheRequestedOrganization(), TestResolveWorkspaceIgnoresAnotherOrganization() (+3 more)

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "workspaces_test.go"
Cohesion: 0.16
Nodes (22): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+14 more)

### Community 92 - "authzAPI"
Cohesion: 0.29
Nodes (9): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), strPtr(), TestResolveWorkspaceAddsGrantRows() (+1 more)

### Community 93 - "Queries"
Cohesion: 0.07
Nodes (5): New(), CreateAccessParams, DBTX, Queries, UpdateJobPlanChangesParams

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 99 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.18
Nodes (9): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+1 more)

### Community 100 - "NewSubscriptions"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 102 - "ProjectAccess"
Cohesion: 0.19
Nodes (6): AccessTokenStore, Deps, fakeAccessTokenStore, fakeWorkspaceStore, ListProjectAccessByProjectAndNamesParams, ProjectAccess

### Community 103 - "context.Context"
Cohesion: 0.04
Nodes (7): globalVarStub, sshAgentStub, teamStub, vcsStub, artifactStore, historyStub, variableStore

### Community 108 - "workspaceStub"
Cohesion: 0.12
Nodes (3): workspaceStub, CreateWorkspaceParams, UpdateWorkspaceParams

### Community 118 - "workspaceTagBindings"
Cohesion: 0.29
Nodes (5): nestedDocument, RelationshipList(), identifierStringAttribute(), RequestRelationship, ResourceIdentifier

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 121 - "newTestAPI"
Cohesion: 0.15
Nodes (18): newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes(), newTestAPI() (+10 more)

### Community 123 - "New"
Cohesion: 0.21
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 127 - "stubStore"
Cohesion: 0.18
Nodes (3): GetArchiveRow, UpdateHistoryOutputParams, stubStore

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 141 - "log/slog.Logger"
Cohesion: 0.15
Nodes (7): Validator, Layout, Deps, JobQueue, stubJobQueue, Deps, Deps

### Community 142 - "Access"
Cohesion: 0.26
Nodes (3): Access, ListAccessByWorkspaceAndNamesParams, UpdateAccessParams

### Community 143 - "actionStub"
Cohesion: 0.15
Nodes (3): actionStub, CreateActionParams, UpdateActionParams

### Community 144 - "net/http.Handler"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "NewCipher"
Cohesion: 0.15
Nodes (33): Task 4: `rotatekeys` command, NewPool(), TestBackfillEncryptsPlaintextRows(), NewCipher(), TestCipherEmptyAndPlaintext(), TestCipherKeyring(), TestCipherRejectsBadConfiguration(), TestCipherRejectsWrongColumnAndKey() (+25 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 150 - "net/http.Handler"
Cohesion: 0.09
Nodes (17): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+9 more)

### Community 151 - "time.Time"
Cohesion: 0.39
Nodes (4): stubPATService, timePtr(), Token, memObject

### Community 154 - "NewSubscriptions"
Cohesion: 0.08
Nodes (34): main(), run(), main(), run(), Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials() (+26 more)

### Community 157 - ".cloneConfiguration"
Cohesion: 0.13
Nodes (12): gitConfig, preparedWorkspace, TestGitSourceDetection(), TestSanitizeGitURL(), transient(), contentIDFromSource(), gitEnv(), gitSource() (+4 more)

### Community 158 - ".QueryRow"
Cohesion: 0.05
Nodes (4): CreateContentParams, UpdateContentStatusParams, CreateProjectAccessParams, UpdateProjectAccessParams

### Community 160 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 161 - "logSink"
Cohesion: 0.20
Nodes (3): captureLogPublisher, LogPublisher, logSink

### Community 162 - "workspaceQuery"
Cohesion: 0.32
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 163 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.07
Nodes (8): jobStub, statusQuerier, captureJobs(), SetWorkspaceLastJobParams, UpdateJobStatusParams, Publisher, publishingStore, stubQuerier

### Community 164 - "New"
Cohesion: 0.21
Nodes (13): TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput(), TestGetStepLogPrefersArchivedObject(), TestGetStepLogRejectsWrongScope(), TestGetStepLogRequiresToken() (+5 more)

### Community 165 - "auth.DexValidator"
Cohesion: 0.23
Nodes (5): Dex, contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 169 - ".streamStepLog"
Cohesion: 0.26
Nodes (9): isTerminal(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData(), writeEvent() (+1 more)

### Community 171 - "serveAccessToken"
Cohesion: 0.25
Nodes (11): expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups(), TestAccessTokenRejectsBadInput() (+3 more)

### Community 172 - "log/slog.Logger"
Cohesion: 0.12
Nodes (14): main(), run(), formatOldKeys(), logReport(), main(), run(), Deps, fakePinger (+6 more)

### Community 174 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 181 - "CreateProjectParams"
Cohesion: 0.16
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 183 - "StringOrInt"
Cohesion: 0.07
Nodes (5): ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, WorkspaceConnection, WorkspaceEdge

### Community 184 - "Deps"
Cohesion: 0.29
Nodes (3): Deps, fakeLocker, JobLocker

### Community 190 - "serveGlobalVar"
Cohesion: 0.39
Nodes (9): newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent(), TestGlobalVariableFromAnotherOrganizationIsNotFound(), TestListGlobalVariablesHidesSensitiveValues() (+1 more)

### Community 191 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 193 - "get"
Cohesion: 0.53
Nodes (6): get(), newUI(), TestDeepLinksServeTheIndex(), TestEnvConfigIsNeverCached(), TestMissingAssetsAreNotFound(), TestServesTheIndexAndAssets()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 427 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **83 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `.CreateGlobalVariable`, `Workspace`, `tagStub`, `collectionStub`, `CreateOrganizationParams`, `CliAuthSession`, `.Query`, `log/slog.Logger`, `Access`, `actionStub`, `fakeQuerier`, `stubQuerier`, `NewCipher`, `tfe_test.go`, `encryptedStore`, `.CreateAgent`, `.marshalOString2ᚖstring`, `.CreateArchive`, `.CreateNotificationTrigger`, `CreateCollectionParams`, `.QueryRow`, `api/main.go`, `github.com/aminbenmansour/terrakube-go/internal/store.Job`, `New`, `.CreateTeam`, `github.com/aminbenmansour/terrakube-go/internal/store.Tag`, `Migrate`, `Pat`, `CreateProjectParams`, `.marshalOBoolean2bool`, `Claims`, `.SetCliAuthSessionIdentity`, `.CreateVCS`, `stubQuerier`, `store.Team`, `templateStub`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `CreateTagParams`, `Queries`, `._WorkspaceTag`, `ProjectAccess`, `context.Context`, `workspaceStub`, `newTestAPI`, `Step`, `stubStore`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.08041140719962599 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `github.com/99designs/gqlgen/graphql.CollectedField` to `.___EnumValue`, `._Project`, `idsFrom`, `context.Context`, `.___Field`, `github.com/99designs/gqlgen/graphql.Marshaler`, `generated.go`, `.marshalOString2ᚖstring`, `StringOrInt`, `.marshalOBoolean2bool`, `._WorkspaceTag`, `OrganizationConnection`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `openRotationTest()` connect `NewCipher` to `context.Context`, `Migrate`, `testing.T`, `Queries`, `api/main.go`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Should `io.Reader` be split into smaller, more focused modules?**
  _Cohesion score 0.07372549019607844 - nodes in this community are weakly interconnected._