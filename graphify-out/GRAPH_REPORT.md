# Graph Report - terrakube-go  (2026-11-25)

## Corpus Check
- 437 files · ~203,827 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 3751 nodes · 14892 edges · 183 communities (101 shown, 82 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 2824 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `115e44b9`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- Globalvar
- io.Reader
- Core Database Schema
- tagStub
- collectionStub
- Handler
- WriteResource
- CliAuthSession
- graphql/collections.go
- Resource
- steps
- CreateProjectParams
- Vc
- stubQuerier
- Ssh
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- github.com/99designs/gqlgen/graphql.Marshaler
- newFixture
- Identity
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- api/main.go
- .runApply
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
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- Executor
- nestedStub
- New
- graphql/collections.go
- WriteError
- Querier
- Querier
- stubQuerier
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- ._Project
- models
- Service
- Team
- collections.go
- Access
- reliabilityStore
- generated.go
- where_test.go
- Resolver
- templateStub
- decodeJSON
- graphql/models.go
- newJWKSServer
- TemplateConnection
- store.Job
- REST contract, captured from upstream
- pat.py
- workspaces_test.go
- CreateTagParams
- .CreateAccess
- authzAPI
- .Query
- Node
- environment
- NewSubscriptions
- newTestService
- newService
- github.com/vektah/gqlparser/v2/ast.SelectionSet
- workspaceStub
- ui-rest-inventory.sh
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- ui-smoke.sh
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- New
- capture.sh
- 0008_action.sql
- stubStore
- context.Context
- compare.sh
- start.sh
- .QueryRow
- .CreateProjectAccess
- tarGzEntries
- log/slog.Logger
- Dex
- actionStub
- net/http.Handler
- run
- tfe_test.go
- net/http.Handler
- JobConnection
- ProjectAccess
- CreateCollectionParams
- cancelRunStore
- .cloneConfiguration
- .CreateVariable
- AccessConnection
- ui.go
- Content
- auth.DexValidator
- newTestExecutor
- ResolveOrganization
- notificationTestAPI
- .streamStepLog
- .updateWorkspaceRelationship
- logSink
- ResourceIdentifier
- time.Time
- ParsePage
- .AddNotificationConfigurationTemplate

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 168 edges
3. `Queries` - 164 edges
4. `API` - 162 edges
5. `New()` - 145 edges
6. `stubQuerier` - 131 edges
7. `WriteError()` - 112 edges
8. `param()` - 78 edges
9. `API` - 66 edges
10. `ParseFieldsets()` - 65 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `Layout`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfstate/layout.go
- `Task 3: `/access-token/v1` permission endpoints` --references--> `Authenticate()`  [INFERRED]
  docs/superpowers/plans/2026-11-24-rbac-permission-resolution.md → internal/auth/identity.go
- `Task 3: `/access-token/v1` permission endpoints` --references--> `NewAccessToken()`  [INFERRED]
  docs/superpowers/plans/2026-11-24-rbac-permission-resolution.md → internal/api/accesstoken.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (183 total, 82 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (28): Write(), WriteResource(), configurationVersionResource(), parseWorkspaceQuery(), approvalGateStatus(), failedJobStatus(), parseRunID(), stepStatusTFE() (+20 more)

### Community 1 - "rbac.go"
Cohesion: 0.18
Nodes (28): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 3: `/access-token/v1` permission endpoints, Task 4: Record the contract, adminOnly() (+20 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "Globalvar"
Cohesion: 0.14
Nodes (3): globalVarStub, CreateGlobalVariableParams, UpdateGlobalVariableParams

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (11): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+3 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "tagStub"
Cohesion: 0.11
Nodes (3): tagStub, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams

### Community 7 - "collectionStub"
Cohesion: 0.09
Nodes (4): collectionStub, CreateItemParams, CreateReferenceParams, UpdateItemParams

### Community 8 - "Handler"
Cohesion: 0.13
Nodes (15): Claims, Type, Classify(), ExtractToken(), claimsFromContext(), isActiveStepStatus(), isExternalUser(), withClaims() (+7 more)

### Community 10 - "WriteResource"
Cohesion: 0.11
Nodes (11): API, listQuery, filterSortPage(), New(), param(), parseList(), unknownCollection(), projectAccessResource() (+3 more)

### Community 11 - "CliAuthSession"
Cohesion: 0.10
Nodes (6): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams, SetCliAuthSessionIdentityParams

### Community 12 - "graphql/collections.go"
Cohesion: 0.24
Nodes (3): DB, failingDB, loadPage()

### Community 13 - "Resource"
Cohesion: 0.10
Nodes (37): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), accessResource(), jobResource() (+29 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 15 - "CreateProjectParams"
Cohesion: 0.18
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 16 - "Vc"
Cohesion: 0.12
Nodes (3): vcsStub, CreateVCSParams, UpdateVCSParams

### Community 18 - "Ssh"
Cohesion: 0.07
Nodes (4): sshAgentStub, CreateAgentParams, CreateSSHKeyParams, UpdateSSHKeyParams

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.07
Nodes (53): replay, upstreamCapture, expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner() (+45 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.06
Nodes (3): PageInfo, WorkspaceTagConnection, WorkspaceTagEdge

### Community 24 - "newFixture"
Cohesion: 0.13
Nodes (23): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+15 more)

### Community 25 - "Identity"
Cohesion: 0.11
Nodes (12): groupVerifier, PATService, stubDexVerifier, stubIdentityVerifier, fakeVerifier, Identity, fakeIssuer, groupVerifier (+4 more)

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
Nodes (60): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs, grantAttrs (+52 more)

### Community 31 - ".runApply"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

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

### Community 51 - "time.Duration"
Cohesion: 0.10
Nodes (6): fakeQueue, JobQueue, defaultRetryDelay(), Message, Queue, stubJobQueue

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (107): TestAccessTokenWithoutStore(), newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+99 more)

### Community 53 - "Pat"
Cohesion: 0.17
Nodes (5): fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Executor"
Cohesion: 0.15
Nodes (6): Executor, preparedWorkspace, isCancelledStatus(), stringPointer(), TestBinarySelection(), baseEnvironment()

### Community 58 - "New"
Cohesion: 0.10
Nodes (21): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+13 more)

### Community 59 - "graphql/collections.go"
Cohesion: 0.10
Nodes (22): collectionSpec, listArgs, page, accessEdges(), afterString(), column(), jobEdges(), notificationConfigurationEdges() (+14 more)

### Community 60 - "WriteError"
Cohesion: 0.09
Nodes (43): actionAttrs, atomicOperation, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, vcsAttrs, decodeActionAttrs() (+35 more)

### Community 61 - "Querier"
Cohesion: 0.04
Nodes (3): New(), DBTX, Queries

### Community 62 - "Querier"
Cohesion: 0.07
Nodes (3): CreateOrganizationParams, Querier, UpdateOrganizationParams

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.07
Nodes (47): fakeLocker, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+39 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "Service"
Cohesion: 0.15
Nodes (10): Service, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), randomToken() (+2 more)

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - "Access"
Cohesion: 0.28
Nodes (4): contains(), contains(), Access, ListAccessByWorkspaceAndNamesParams

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "Resolver"
Cohesion: 0.10
Nodes (9): projectResolver, Resolver, ResolverRoot, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace() (+1 more)

### Community 77 - "templateStub"
Cohesion: 0.17
Nodes (3): templateStub, CreateTemplateParams, UpdateTemplateParams

### Community 78 - "decodeJSON"
Cohesion: 0.11
Nodes (31): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent() (+23 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 82 - "store.Job"
Cohesion: 0.12
Nodes (4): UpdateJobStatusParams, Publisher, publishingStore, stubQuerier

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "workspaces_test.go"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 88 - "CreateTagParams"
Cohesion: 0.15
Nodes (3): CreateTagParams, GetTagByNameParams, UpdateTagParams

### Community 92 - "authzAPI"
Cohesion: 0.17
Nodes (13): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), strPtr(), TestResolveOrganization() (+5 more)

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "NewSubscriptions"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 102 - "newService"
Cohesion: 0.19
Nodes (15): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+7 more)

### Community 103 - "github.com/vektah/gqlparser/v2/ast.SelectionSet"
Cohesion: 0.04
Nodes (13): ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, NotificationTriggerConnection, NotificationTriggerEdge, OrganizationConnection, OrganizationEdge, ProjectConnection (+5 more)

### Community 108 - "workspaceStub"
Cohesion: 0.10
Nodes (4): fakeAccessTokenStore, workspaceStub, CreateWorkspaceParams, UpdateWorkspaceParams

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 122 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.16
Nodes (4): jobStub, statusQuerier, captureJobs(), SetWorkspaceLastJobParams

### Community 123 - "New"
Cohesion: 0.18
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 127 - "stubStore"
Cohesion: 0.18
Nodes (3): GetArchiveRow, UpdateHistoryOutputParams, stubStore

### Community 128 - "context.Context"
Cohesion: 0.05
Nodes (5): teamStub, artifactStore, historyStub, variableStore, stubQuerier

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 141 - "log/slog.Logger"
Cohesion: 0.09
Nodes (17): Validator, main(), run(), Deps, JobLocker, Deps, fakePinger, Pinger (+9 more)

### Community 142 - "Dex"
Cohesion: 0.21
Nodes (3): Dex, Exchanger, fakeDex

### Community 143 - "actionStub"
Cohesion: 0.15
Nodes (3): actionStub, CreateActionParams, UpdateActionParams

### Community 144 - "net/http.Handler"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "run"
Cohesion: 0.04
Nodes (65): main(), run(), main(), run(), main(), run(), Config, Locker (+57 more)

### Community 147 - "tfe_test.go"
Cohesion: 0.09
Nodes (42): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+34 more)

### Community 150 - "net/http.Handler"
Cohesion: 0.11
Nodes (12): authorization, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail(), withAuthorization() (+4 more)

### Community 154 - "ProjectAccess"
Cohesion: 0.25
Nodes (4): AccessTokenStore, Deps, ListProjectAccessByProjectAndNamesParams, ProjectAccess

### Community 157 - ".cloneConfiguration"
Cohesion: 0.15
Nodes (11): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), transient(), contentIDFromSource(), gitEnv(), gitSource(), isConfigurationTarball() (+3 more)

### Community 160 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 162 - "Content"
Cohesion: 0.25
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 165 - "auth.DexValidator"
Cohesion: 0.31
Nodes (4): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 166 - "newTestExecutor"
Cohesion: 0.31
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 167 - "ResolveOrganization"
Cohesion: 0.16
Nodes (10): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access, fakeWorkspaceStore (+2 more)

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 169 - ".streamStepLog"
Cohesion: 0.33
Nodes (7): isTerminal(), outputOf(), stepOf(), streamKey(), writeData(), writeEvent(), Handler

### Community 171 - "logSink"
Cohesion: 0.12
Nodes (9): captureLogPublisher, LogPublisher, logSink, redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues() (+1 more)

### Community 172 - "ResourceIdentifier"
Cohesion: 0.29
Nodes (5): nestedDocument, RelationshipList(), identifierStringAttribute(), RequestRelationship, ResourceIdentifier

### Community 174 - "time.Time"
Cohesion: 0.39
Nodes (4): stubPATService, timePtr(), Token, memObject

### Community 176 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 389 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **82 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `context.Context`, `Globalvar`, `.QueryRow`, `tagStub`, `collectionStub`, `.CreateProjectAccess`, `CliAuthSession`, `actionStub`, `CreateProjectParams`, `stubQuerier`, `Ssh`, `tfe_test.go`, `Vc`, `github.com/99designs/gqlgen/graphql.Marshaler`, `ProjectAccess`, `CreateCollectionParams`, `api/main.go`, `.CreateVariable`, `Content`, `.GetHistoryByJobReference`, `.AddNotificationConfigurationTemplate`, `testing.T`, `Pat`, `Querier`, `stubQuerier`, `Team`, `Access`, `templateStub`, `store.Job`, `CreateTagParams`, `.CreateAccess`, `.Query`, `workspaceStub`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `stubStore`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07435610302351624 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `github.com/99designs/gqlgen/graphql.CollectedField` to `context.Context`, `._JobStatusEvent`, `.___Type`, `._Project`, `github.com/vektah/gqlparser/v2/ast.SelectionSet`, `.childFields___InputValue`, `generated.go`, `Resolver`, `TemplateConnection`, `.___Directive`, `.marshalOBoolean2bool`, `._Organization`, `github.com/99designs/gqlgen/graphql.Marshaler`, `JobConnection`, `AccessConnection`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `API` connect `WriteResource` to `.updateWorkspaceRelationship`, `Resource`, `net/http.Handler`, `ProjectAccess`, `WriteError`, `api/main.go`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `Globalvar` be split into smaller, more focused modules?**
  _Cohesion score 0.1437908496732026 - nodes in this community are weakly interconnected._