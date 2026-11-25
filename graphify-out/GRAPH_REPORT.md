# Graph Report - terrakube-go  (2026-11-25)

## Corpus Check
- 457 files · ~225,166 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4057 nodes · 16367 edges · 189 communities (102 shown, 87 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 3062 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9e8229f7`
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
- fakeQuerier
- stubQuerier
- encryptedStore
- tfe_test.go
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- ._Template
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
- NotificationTriggerEdge
- stubStore
- Dex
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- Executor
- OrganizationConnection
- store.ProjectAccess
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
- .ListTeamsByOrganizationAndNames
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
- authzAPI
- serveAccessToken
- environment
- Layout
- NewSubscriptions
- newTestService
- context.Context
- TestTagListMatchesInMemory
- ui-rest-inventory.sh
- .Query
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
- Access
- actionStub
- net/http.Handler
- run
- Secret key rotation — design
- M4 run logs and artifacts Implementation Plan
- net/http.Handler
- .updateWorkspaceRelationship
- NewSubscriptions
- time.Time
- reliabilityStore
- .cloneConfiguration
- .QueryRow
- Review Focus
- newService
- logSink
- ResourceIdentifier
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- .PublishCapped
- auth.DexValidator
- notificationTestAPI
- SetWorkspaceLastJobParams
- Review Focus
- Deps
- .Ping
- store.Workspace
- ParsePage
- .SoftDeletePat
- CreateProjectParams
- models_gen.go

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 202 edges
3. `Queries` - 200 edges
4. `API` - 162 edges
5. `stubQuerier` - 142 edges
6. `New()` - 142 edges
7. `WriteError()` - 114 edges
8. `param()` - 78 edges
9. `API` - 67 edges
10. `ParseFieldsets()` - 65 edges

## Surprising Connections (you probably didn't know these)
- `Task 5: Backfill implementation` --references--> `needsEncryption()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/secretstore/backfill.go
- `Command` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Report and verification` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Secret Encryption at Rest Implementation Plan` --references--> `Querier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/store/querier.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (189 total, 87 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (30): Write(), WriteResource(), configurationVersionResource(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource() (+22 more)

### Community 1 - "rbac.go"
Cohesion: 0.26
Nodes (22): Task 1: Role matrix and permission set, adminOnly(), CanApproveJob(), CanDestroyJob(), CanManageCollection(), CanManageJob(), CanManageModule(), CanManagePolicies() (+14 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (14): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+6 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "tagStub"
Cohesion: 0.08
Nodes (3): tagStub, matchingTags(), CountTagsByNameParams

### Community 7 - "collectionStub"
Cohesion: 0.06
Nodes (6): collectionStub, CreateCollectionParams, CreateItemParams, CreateReferenceParams, UpdateCollectionParams, UpdateItemParams

### Community 8 - "Handler"
Cohesion: 0.11
Nodes (20): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), WithClaims(), runnerClaims(), claimsFromContext() (+12 more)

### Community 9 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.05
Nodes (5): AccessEdge, NotificationConfigurationEdge, ProjectEdge, WorkspaceEdge, WorkspaceTagEdge

### Community 10 - "net/http.Request"
Cohesion: 0.12
Nodes (7): API, New(), param(), parseList(), unknownCollection(), ParseFieldsets(), WriteCollection()

### Community 11 - "CliAuthSession"
Cohesion: 0.05
Nodes (15): fakeStore, Service, Store, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value() (+7 more)

### Community 12 - ".Query"
Cohesion: 0.05
Nodes (16): ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow, ListSSHKeySecretsAfterParams, ListSSHKeySecretsAfterRow (+8 more)

### Community 13 - "Resource"
Cohesion: 0.11
Nodes (35): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), accessResource(), projectAccessResource() (+27 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 16 - "fakeQuerier"
Cohesion: 0.06
Nodes (10): fakeQuerier, CreateNotificationConfigurationParams, CreateSSHKeyParams, CreateVariableParams, CreateVCSParams, UpdateNotificationConfigurationParams, UpdateSSHKeyParams, UpdateVariableParams (+2 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (41): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue() (+33 more)

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
Cohesion: 0.09
Nodes (35): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+27 more)

### Community 25 - "Identity"
Cohesion: 0.11
Nodes (12): groupVerifier, stubDexVerifier, stubIdentityVerifier, fakeVerifier, Identity, fakeIssuer, TokenIssuer, groupVerifier (+4 more)

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
Nodes (57): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs, grantAttrs (+49 more)

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
Cohesion: 0.15
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

### Community 47 - "stubStore"
Cohesion: 0.18
Nodes (3): GetArchiveRow, UpdateHistoryOutputParams, stubStore

### Community 48 - "Dex"
Cohesion: 0.21
Nodes (3): Dex, Exchanger, fakeDex

### Community 51 - "time.Duration"
Cohesion: 0.11
Nodes (7): fakeQueue, JobQueue, defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), Message, Queue

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (112): TestAccessTokenWithoutStore(), newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+104 more)

### Community 53 - "Pat"
Cohesion: 0.18
Nodes (5): newFakeStore(), fakeStore, Store, CreatePatParams, Pat

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Executor"
Cohesion: 0.11
Nodes (9): Executor, isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestContentIDFromSource(), TestIsTerminalStatus(), transient() (+1 more)

### Community 59 - "graphql/collections.go"
Cohesion: 0.09
Nodes (24): collectionSpec, DB, listArgs, page, accessEdges(), afterString(), column(), jobEdges() (+16 more)

### Community 60 - "WriteError"
Cohesion: 0.09
Nodes (41): actionAttrs, atomicOperation, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, vcsAttrs, decodeActionAttrs() (+33 more)

### Community 61 - "Querier"
Cohesion: 0.03
Nodes (4): AddNotificationConfigurationTemplateParams, DBTX, Querier, Queries

### Community 62 - "newTestExecutor"
Cohesion: 0.31
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.06
Nodes (52): redactor, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+44 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "config.go"
Cohesion: 0.14
Nodes (23): Config, Storage, Task 2: Serve archived logs from `/tfoutput/v1`, Task 1: Keyring configuration, env(), envBool(), envDuration(), envFirst() (+15 more)

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 73 - "Cipher"
Cohesion: 0.11
Nodes (26): Task 3: Rotation scan and re-encryption, backfillGlobalVariables(), backfillNotifications(), backfillSSHKeys(), backfillVariables(), backfillVCS(), encryptIfNeeded(), needsEncryption() (+18 more)

### Community 74 - "generated.go"
Cohesion: 0.07
Nodes (7): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext(), parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

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
Cohesion: 0.14
Nodes (22): nestedStub, TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newJobStub(), serveJob(), TestJobApprovalRequiresWaitingApproval(), TestJobCancelIsNoContentAndReachesTheStatusPath(), TestJobGetFromAnotherOrganizationIsNotFound() (+14 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "ResolveOrganization"
Cohesion: 0.26
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
Cohesion: 0.14
Nodes (5): CreateTagParams, GetTagByNameParams, ListTagsPageDescParams, ListTagsPageParams, UpdateTagParams

### Community 92 - "authzAPI"
Cohesion: 0.33
Nodes (10): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), strPtr(), TestResolveOrganization() (+2 more)

### Community 94 - "serveAccessToken"
Cohesion: 0.25
Nodes (11): expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner(), TestAccessTokenOrganizationPermissionsWithoutGroups(), TestAccessTokenRejectsBadInput() (+3 more)

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 99 - "Layout"
Cohesion: 0.13
Nodes (6): Validator, Deps, fakeLocker, JobLocker, Layout, Deps

### Community 100 - "NewSubscriptions"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 103 - "context.Context"
Cohesion: 0.03
Nodes (9): globalVarStub, sshAgentStub, teamStub, vcsStub, workspaceStub, artifactStore, historyStub, variableStore (+1 more)

### Community 108 - "TestTagListMatchesInMemory"
Cohesion: 0.15
Nodes (17): listQuery, pageMeta, filterSortPage(), prefixPattern(), createListOrganization(), getNames(), newListRequest(), seedOrganizations() (+9 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 123 - "grantsAPI"
Cohesion: 0.31
Nodes (10): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+2 more)

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 143 - "actionStub"
Cohesion: 0.15
Nodes (3): actionStub, CreateActionParams, UpdateActionParams

### Community 144 - "net/http.Handler"
Cohesion: 0.12
Nodes (19): CookieSigner, Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), authorize() (+11 more)

### Community 145 - "run"
Cohesion: 0.05
Nodes (78): main(), run(), main(), run(), main(), run(), main(), run() (+70 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 147 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.25
Nodes (7): Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph

### Community 150 - "net/http.Handler"
Cohesion: 0.09
Nodes (17): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+9 more)

### Community 154 - "NewSubscriptions"
Cohesion: 0.15
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 155 - "time.Time"
Cohesion: 0.39
Nodes (4): stubPATService, timePtr(), Token, memObject

### Community 156 - "reliabilityStore"
Cohesion: 0.06
Nodes (14): cancelRunStore, reliabilityStore, isTerminal(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame() (+6 more)

### Community 157 - ".cloneConfiguration"
Cohesion: 0.16
Nodes (10): gitConfig, preparedWorkspace, TestGitSourceDetection(), TestSanitizeGitURL(), gitEnv(), gitSource(), isConfigurationTarball(), normalizeGitURL() (+2 more)

### Community 158 - ".QueryRow"
Cohesion: 0.04
Nodes (7): CreateContentParams, UpdateContentStatusParams, CountProjectsByNameParams, CreateAccessParams, CreateNotificationTriggerParams, CreateProjectAccessParams, UpdateProjectAccessParams

### Community 159 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 160 - "newService"
Cohesion: 0.21
Nodes (14): decodeSecret(), New(), dexUser(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp(), TestDeleteEnforcesOwnership() (+6 more)

### Community 161 - "logSink"
Cohesion: 0.17
Nodes (4): logSink, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

### Community 162 - "ResourceIdentifier"
Cohesion: 0.29
Nodes (5): nestedDocument, RelationshipList(), identifierStringAttribute(), RequestRelationship, ResourceIdentifier

### Community 163 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.08
Nodes (7): jobStub, statusQuerier, captureJobs(), UpdateJobStatusParams, Publisher, publishingStore, stubQuerier

### Community 165 - "auth.DexValidator"
Cohesion: 0.31
Nodes (4): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 170 - "Review Focus"
Cohesion: 0.33
Nodes (5): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract

### Community 171 - "Deps"
Cohesion: 0.40
Nodes (3): AccessTokenStore, Deps, PATService

### Community 174 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 181 - "CreateProjectParams"
Cohesion: 0.15
Nodes (5): CreateProjectParams, GetProjectByNameParams, ListProjectsPageDescParams, ListProjectsPageParams, UpdateProjectParams

### Community 183 - "models_gen.go"
Cohesion: 0.07
Nodes (11): AccessConnection, ComplexityRoot, NotificationConfigurationConnection, NotificationTriggerConnection, PageInfo, ProjectConnection, Query, Subscription (+3 more)

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 439 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **87 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `.CreateAgent`, `net/http.ResponseWriter`, `.CreateGlobalVariable`, `Workspace`, `tagStub`, `collectionStub`, `Handler`, `.CreateOrganization`, `CliAuthSession`, `.Query`, `Access`, `actionStub`, `fakeQuerier`, `stubQuerier`, `run`, `tfe_test.go`, `encryptedStore`, `._Template`, `newFixture`, `.CreateArchive`, `reliabilityStore`, `.QueryRow`, `api/main.go`, `github.com/aminbenmansour/terrakube-go/internal/store.Job`, `.CreateTeam`, `SetWorkspaceLastJobParams`, `Deps`, `NotificationTriggerEdge`, `stubStore`, `.SoftDeletePat`, `testing.T`, `Pat`, `CreateProjectParams`, `store.ProjectAccess`, `.CreateWorkspaceTag`, `stubQuerier`, `Cipher`, `templateStub`, `Tag`, `Layout`, `context.Context`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07675438596491228 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `github.com/99designs/gqlgen/graphql.CollectedField` to `._Project`, `._Organization`, `context.Context`, `.___Field`, `github.com/99designs/gqlgen/graphql.Marshaler`, `generated.go`, `.childFields___InputValue`, `NotificationTriggerEdge`, `Resolver`, `.childFields___Type`, `._Template`, `models_gen.go`, `OrganizationConnection`, `JobConnection`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `stubQuerier` connect `stubQuerier` to `net/http.ResponseWriter`, `Workspace`, `tagStub`, `Access`, `actionStub`, `encryptedStore`, `strPtr`, `api/main.go`, `.QueryRow`, `github.com/aminbenmansour/terrakube-go/internal/store.Job`, `store.Team`, `SetWorkspaceLastJobParams`, `store.Workspace`, `store.ProjectAccess`, `Querier`, `stubQuerier`, `.ListTeamsByOrganizationAndNames`, `ResolveOrganization`, `github.com/aminbenmansour/terrakube-go/internal/store.Organization`, `context.Context`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `io.Reader` be split into smaller, more focused modules?**
  _Cohesion score 0.07058001397624039 - nodes in this community are weakly interconnected._