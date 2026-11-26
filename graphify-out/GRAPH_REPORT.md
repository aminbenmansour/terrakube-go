# Graph Report - terrakube-go  (2026-11-26)

## Corpus Check
- 481 files · ~238,875 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4293 nodes · 17155 edges · 248 communities (117 shown, 131 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3166 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `83e4955e`
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
- collectionStub
- Handler
- ._Project
- net/http.Request
- fakeQuerier
- .Query
- github.com/aminbenmansour/terrakube-go/internal/store.Project
- steps
- stubQuerier
- newTestExecutor
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- github.com/jackc/pgx/v5/pgtype.Timestamptz
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
- serveGlobalVar
- authzAPI
- logSink
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- .handle
- Tag
- store.Access
- NewCipher
- .applyWorkspaceUpdate
- Querier
- Backend
- stubQuerier
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- OrganizationConnection
- models
- config.go
- Team
- collections.go
- .CreateWorkspaceTag
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
- JobConnection
- Identity
- environment
- NewSubscriptions
- newTestService
- .CreateProject
- context.Context
- stubStore
- ui-rest-inventory.sh
- .ListProjectAccessByProjectAndNames
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
- newScheduleAPI
- time.Time
- tarGzEntries
- .CreateAccess
- parser
- where.go
- clilogin/handler_test.go
- newTestAPI
- Rotate
- io.Reader
- Job
- net/http.Handler
- actionStub
- .CreateArchive
- NewSubscriptions
- Workspace
- .prepareWorkspace
- Node
- M4 run logs and artifacts Implementation Plan
- logSink
- newService
- CliAuthSession
- NotificationTriggerConnection
- auth.DexValidator
- .updateWorkspaceRelationship
- notificationTestAPI
- CreateCollectionParams
- Metrics
- transientError
- .CountList
- Threat model
- models_gen.go
- .Get
- .Query
- .SoftDeletePat
- RSQL filters in SQL — design
- .ListVariableSecretsAfter
- CreateSSHKeyParams
- Object
- .EnqueueJob
- Review Focus
- workspaceQuery
- Local
- .TryLock
- ParsePage
- .Ping

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Queries` - 210 edges
3. `Querier` - 203 edges
4. `API` - 169 edges
5. `New()` - 152 edges
6. `stubQuerier` - 142 edges
7. `WriteError()` - 125 edges
8. `param()` - 81 edges
9. `ParseFieldsets()` - 69 edges
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

## Communities (248 total, 131 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (23): configurationVersionResource(), projectResource(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource(), runStatus() (+15 more)

### Community 1 - "rbac.go"
Cohesion: 0.19
Nodes (27): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, adminOnly(), CanApproveJob() (+19 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - ".CreateGlobalVariable"
Cohesion: 0.24
Nodes (3): CreateGlobalVariableParams, Globalvar, UpdateGlobalVariableParams

### Community 4 - "memStorage"
Cohesion: 0.14
Nodes (4): failingBackend, NewS3(), S3, logBackend

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - ".QueryRow"
Cohesion: 0.05
Nodes (10): CreateContentParams, UpdateContentStatusParams, ApproveJobParams, Content, CreateJobParams, CreateNotificationTriggerParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams (+2 more)

### Community 7 - "collectionStub"
Cohesion: 0.09
Nodes (5): collectionStub, CreateItemParams, CreateReferenceParams, UpdateItemParams, Item

### Community 8 - "Handler"
Cohesion: 0.09
Nodes (28): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), WithClaims(), isTerminal(), outputOf() (+20 more)

### Community 9 - "._Project"
Cohesion: 0.13
Nodes (3): ProjectAccessConnection, ProjectAccessEdge, ProjectEdge

### Community 10 - "net/http.Request"
Cohesion: 0.08
Nodes (37): actionAttrs, API, Deps, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, actionResource() (+29 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.05
Nodes (6): nestedStub, encryptedStore, fakeQuerier, CreateNotificationConfigurationParams, UpdateNotificationConfigurationParams, Variable

### Community 12 - ".Query"
Cohesion: 0.06
Nodes (15): CountJobsByStatusRow, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow, ListSSHKeySecretsAfterParams (+7 more)

### Community 13 - "github.com/aminbenmansour/terrakube-go/internal/store.Project"
Cohesion: 0.11
Nodes (11): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), contains(), Access, fakeOrganizationStore (+3 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 19 - "newTestExecutor"
Cohesion: 0.12
Nodes (24): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+16 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.07
Nodes (52): replay, stubIdentityVerifier, upstreamCapture, expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions() (+44 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.13
Nodes (7): scheduleStub, CreateScheduleParams, IssueCliAuthCodeParams, NotificationTrigger, Schedule, Template, UpdateScheduleParams

### Community 24 - "newFixture"
Cohesion: 0.13
Nodes (25): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+17 more)

### Community 25 - "Service"
Cohesion: 0.20
Nodes (6): PATService, stubPATService, stringPtr(), timePtr(), Service, Token

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
Nodes (28): parseInclude(), collectionResource(), itemResource(), referenceResource(), jobResource(), stepResource(), timestamp(), notificationConfigurationResource() (+20 more)

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

### Community 46 - "graphql/collections.go"
Cohesion: 0.09
Nodes (24): collectionSpec, DB, listArgs, page, accessEdges(), afterString(), column(), jobEdges() (+16 more)

### Community 47 - "serveGlobalVar"
Cohesion: 0.18
Nodes (10): globalVarStub, newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent(), TestGlobalVariableFromAnotherOrganizationIsNotFound() (+2 more)

### Community 48 - "authzAPI"
Cohesion: 0.17
Nodes (13): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), strPtr(), TestResolveOrganization() (+5 more)

### Community 49 - "logSink"
Cohesion: 0.19
Nodes (7): logSink, redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 51 - "time.Duration"
Cohesion: 0.10
Nodes (5): captureLogPublisher, fakeQueue, JobQueue, LogPublisher, Queue

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (118): TestAccessTokenWithoutStore(), TestAuthMiddleware(), TestCreateOrganization(), newCollectionAPI(), TestCollectionListCarriesItsItemsAndReferences(), TestCreateCollectionDefaultsPriorityToZero(), TestCreateCollectionItemStoresTheValue(), TestCreateCollectionItemValidatesKeyValueAndCategory() (+110 more)

### Community 53 - "Pat"
Cohesion: 0.23
Nodes (4): fakeStore, Store, CreatePatParams, Pat

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - ".handle"
Cohesion: 0.22
Nodes (4): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus()

### Community 56 - "Tag"
Cohesion: 0.17
Nodes (4): CreateTagParams, GetTagByNameParams, Tag, UpdateTagParams

### Community 59 - "NewCipher"
Cohesion: 0.06
Nodes (56): main(), run(), main(), run(), main(), run(), formatOldKeys(), logReport() (+48 more)

### Community 60 - ".applyWorkspaceUpdate"
Cohesion: 0.08
Nodes (30): nestedDocument, vcsAttrs, strPtr(), strValue(), atomicWorkspaceHref(), isAtomicRequest(), writeInUseConflict(), decodeVCSAttrs() (+22 more)

### Community 61 - "Querier"
Cohesion: 0.04
Nodes (9): createRawVariable(), createRawVariableWithID(), AddNotificationConfigurationTemplateParams, Agent, CreateAgentParams, DBTX, Querier, Queries (+1 more)

### Community 62 - "Backend"
Cohesion: 0.21
Nodes (8): Validator, Deps, Deps, Backend, Layout, Deps, Deps, Deps

### Community 63 - "stubQuerier"
Cohesion: 0.05
Nodes (6): jobStub, statusQuerier, captureJobs(), stubQuerier, workspaceWithTags, stubQuerier

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.05
Nodes (69): main(), run(), gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource() (+61 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "config.go"
Cohesion: 0.14
Nodes (20): Config, Task 1: Keyring configuration, env(), envBool(), envFirst(), envInt64(), envList(), Load() (+12 more)

### Community 70 - "Team"
Cohesion: 0.09
Nodes (6): AccessTokenStore, fakeAccessTokenStore, contains(), ProjectAccess, Team, Workspace

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - ".CreateWorkspaceTag"
Cohesion: 0.21
Nodes (3): CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, Workspacetag

### Community 73 - "Cipher"
Cohesion: 0.08
Nodes (37): Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation (+29 more)

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "streamEvents"
Cohesion: 0.14
Nodes (8): projectResolver, Resolver, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace(), TestSubscriptionsNeedAnEventSource()

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

### Community 94 - "Identity"
Cohesion: 0.13
Nodes (8): groupVerifier, stubDexVerifier, fakeVerifier, Identity, fakeIssuer, groupVerifier, stubDexVerifier, stringSlice()

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "NewSubscriptions"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (20): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+12 more)

### Community 102 - ".CreateProject"
Cohesion: 0.20
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 103 - "context.Context"
Cohesion: 0.04
Nodes (7): sshAgentStub, teamStub, vcsStub, workspaceStub, artifactStore, historyStub, variableStore

### Community 108 - "stubStore"
Cohesion: 0.18
Nodes (3): GetArchiveRow, UpdateHistoryOutputParams, stubStore

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "cancelRunStore"
Cohesion: 0.15
Nodes (3): cancelRunStore, UpdateJobOutputParams, UpdateStepOutputParams

### Community 121 - "newWorkspaceStub"
Cohesion: 0.16
Nodes (22): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+14 more)

### Community 123 - "reliabilityStore"
Cohesion: 0.07
Nodes (8): reliabilityStore, CreateHistoryParams, CreateStepParams, GetHistoryByJobReferenceParams, GetStepByJobAndNumberParams, History, Step, UpdateStepStatusParams

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 128 - "CreateVCSParams"
Cohesion: 0.24
Nodes (3): CreateVCSParams, UpdateVCSParams, Vc

### Community 135 - "newScheduleAPI"
Cohesion: 0.36
Nodes (9): newScheduleAPI(), newScheduleStub(), scheduleBody(), scheduleStr(), TestCreateScheduleCopiesAttributesAndDefaultsEnabled(), TestCreateScheduleRejectsForeignTemplate(), TestDeleteScheduleAndCrossWorkspaceLookup(), TestListAndUpdateSchedule() (+1 more)

### Community 137 - "time.Time"
Cohesion: 0.22
Nodes (10): bucket, limiter, key(), newLimiter(), RateLimit(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), TestRateLimitDisabledByDefault() (+2 more)

### Community 140 - "tarGzEntries"
Cohesion: 0.29
Nodes (10): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+2 more)

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.12
Nodes (19): CookieSigner, Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), authorize() (+11 more)

### Community 145 - "newTestAPI"
Cohesion: 0.09
Nodes (38): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+30 more)

### Community 146 - "Rotate"
Cohesion: 0.08
Nodes (27): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs, Approaches considered (+19 more)

### Community 149 - "Job"
Cohesion: 0.33
Nodes (4): Publisher, publishingStore, stubQuerier, Job

### Community 150 - "net/http.Handler"
Cohesion: 0.09
Nodes (18): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+10 more)

### Community 151 - "actionStub"
Cohesion: 0.15
Nodes (4): actionStub, CreateActionParams, UpdateActionParams, Action

### Community 154 - "NewSubscriptions"
Cohesion: 0.10
Nodes (25): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+17 more)

### Community 155 - "Workspace"
Cohesion: 0.14
Nodes (3): CreateWorkspaceParams, GetWorkspaceByNameParams, SetWorkspaceLockParams

### Community 157 - ".prepareWorkspace"
Cohesion: 0.13
Nodes (12): gitConfig, preparedWorkspace, TestGitSourceDetection(), TestSanitizeGitURL(), transient(), contentIDFromSource(), gitEnv(), gitSource() (+4 more)

### Community 159 - "Node"
Cohesion: 0.09
Nodes (5): listQuery, stubLists, Task 3: Handler cutover, And, Or

### Community 160 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.18
Nodes (9): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+1 more)

### Community 161 - "logSink"
Cohesion: 0.16
Nodes (5): Executor, isCancelledStatus(), stringPointer(), TestBinarySelection(), baseEnvironment()

### Community 162 - "newService"
Cohesion: 0.19
Nodes (15): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+7 more)

### Community 163 - "CliAuthSession"
Cohesion: 0.07
Nodes (17): Exchanger, fakeStore, Service, Store, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups() (+9 more)

### Community 165 - "auth.DexValidator"
Cohesion: 0.15
Nodes (6): Dex, fakeDex, contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 171 - "CreateCollectionParams"
Cohesion: 0.20
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

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
Cohesion: 0.05
Nodes (13): AccessConnection, AccessEdge, ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, PageInfo, ProjectConnection, Query (+5 more)

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 202 - "CreateSSHKeyParams"
Cohesion: 0.24
Nodes (3): CreateSSHKeyParams, Ssh, UpdateSSHKeyParams

### Community 231 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 232 - "workspaceQuery"
Cohesion: 0.32
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 233 - "Local"
Cohesion: 0.33
Nodes (4): cleanKey(), pathClean(), TestCleanKey(), Local

### Community 235 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 510 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **131 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `CreateVCSParams`, `.CreateGlobalVariable`, `.QueryRow`, `collectionStub`, `Handler`, `net/http.Request`, `fakeQuerier`, `.Query`, `.CreateAccess`, `.marshalOBoolean2bool`, `github.com/99designs/gqlgen/graphql.Marshaler`, `stubQuerier`, `newTestAPI`, `newTestExecutor`, `.CreateNotificationConfiguration`, `Rotate`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `actionStub`, `.CreateArchive`, `Workspace`, `api/main.go`, `CliAuthSession`, `CreateCollectionParams`, `Metrics`, `testing.T`, `Pat`, `Tag`, `store.Access`, `NewCipher`, `Backend`, `.SoftDeletePat`, `stubQuerier`, `NewLayout`, `Team`, `.ListVariableSecretsAfter`, `.CreateWorkspaceTag`, `Cipher`, `CreateSSHKeyParams`, `templateStub`, `tagStub`, `.CreateProject`, `context.Context`, `stubStore`, `.ListProjectAccessByProjectAndNames`, `cancelRunStore`, `reliabilityStore`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.08340016038492382 - nodes in this community are weakly interconnected._
- **Why does `New()` connect `net/http.Request` to `NewLayout`, `.updateWorkspaceRelationship`, `newScheduleAPI`, `notificationTestAPI`, `authzAPI`, `newTestAPI`, `newTestExecutor`, `testing.T`, `strPtr`, `net/http.Handler`, `.applyWorkspaceUpdate`, `api/main.go`, `Link`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `openRotationTest()` connect `Rotate` to `context.Context`, `newTestAPI`, `testing.T`, `Querier`, `api/main.go`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `memStorage` be split into smaller, more focused modules?**
  _Cohesion score 0.14166666666666666 - nodes in this community are weakly interconnected._