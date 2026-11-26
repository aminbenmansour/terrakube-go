# Graph Report - terrakube-go  (2026-11-26)

## Corpus Check
- 491 files · ~247,700 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4409 nodes · 17817 edges · 252 communities (120 shown, 132 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 3297 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2c290d5a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- StringOrInt
- io.Reader
- Core Database Schema
- Job
- Collection
- Handler
- ._Project
- .internalError
- fakeQuerier
- stubQuerier
- ResolveOrganization
- steps
- github.com/99designs/gqlgen/graphql.Marshaler
- stubQuerier
- Variable
- tfe_test.go
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- Workspace
- newFixture
- .Query
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
- Access
- graphql/collections.go
- ProjectAccess
- logSink
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- New
- .QueryRow
- github.com/aminbenmansour/terrakube-go/internal/store.Step
- NewCipher
- tagStub
- Queries
- History
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- OrganizationConnection
- models
- run
- Team
- collections.go
- Workspacetag
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
- log/slog.Logger
- .handle
- newTestExecutor
- Identity
- environment
- cron.go
- NewSubscriptions
- newTestService
- .CreateProject
- context.Context
- Object
- ui-rest-inventory.sh
- Content
- ui-smoke.sh
- runTriggerAPI
- newWorkspaceStub
- cancelRunStore
- capture.sh
- 0008_action.sql
- .CreateAgent
- compare.sh
- start.sh
- TemplateConnection
- Review Focus
- time.Time
- tarGzEntries
- Dex
- rsql.Node
- where.go
- clilogin/handler_test.go
- newTestAPI
- Secret key rotation — design
- AccessConnection
- WorkspaceTagConnection
- API
- Item
- .CreateArchive
- NewSubscriptions
- .prepareWorkspace
- Node
- Executor
- newService
- CliAuthSession
- Reference
- auth.DexValidator
- .applyWorkspaceUpdate
- M4 run logs and artifacts Implementation Plan
- notificationTestAPI
- nestedStub
- authzAPI
- Metrics
- transientError
- Vc
- Threat model
- models_gen.go
- stubStore
- loadPage
- .processJob
- RSQL filters in SQL — design
- .Get
- .listWorkspaceAccess
- Workspace
- .Ping
- .SoftDeletePat
- log/slog.Logger
- Review Focus
- Migrate
- ParsePage

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Queries` - 219 edges
3. `Querier` - 214 edges
4. `API` - 182 edges
5. `New()` - 161 edges
6. `stubQuerier` - 154 edges
7. `WriteError()` - 130 edges
8. `param()` - 82 edges
9. `ParseFieldsets()` - 72 edges
10. `API` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Self-review` --references--> `SortKey`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/store/list.go
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Command` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Report and verification` --references--> `Rotate()`  [INFERRED]
  docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md → internal/secretstore/rotate.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (252 total, 132 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (28): configurationVersionResource(), projectResource(), parseWorkspaceQuery(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID(), runResource() (+20 more)

### Community 1 - "rbac.go"
Cohesion: 0.19
Nodes (27): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, adminOnly(), CanApproveJob() (+19 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "StringOrInt"
Cohesion: 0.10
Nodes (3): ComplexityRoot, ProjectConnection, ProjectEdge

### Community 4 - "io.Reader"
Cohesion: 0.11
Nodes (5): failingBackend, failingKeyBackend, memoryBackend, S3, logBackend

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "Job"
Cohesion: 0.13
Nodes (4): ApproveJobParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams, UpdateJobStatusParams

### Community 7 - "Collection"
Cohesion: 0.20
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 8 - "Handler"
Cohesion: 0.10
Nodes (27): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), WithClaims(), size(), writeTooLarge() (+19 more)

### Community 10 - ".internalError"
Cohesion: 0.09
Nodes (32): actionAttrs, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, actionResource(), decodeActionAttrs(), stringOr() (+24 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.08
Nodes (8): encryptedStore, fakeQuerier, CreateNotificationConfigurationParams, CreateSSHKeyParams, CreateVCSParams, UpdateNotificationConfigurationParams, UpdateSSHKeyParams, UpdateVCSParams

### Community 13 - "ResolveOrganization"
Cohesion: 0.15
Nodes (10): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access, fakeWorkspaceStore (+2 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 16 - "github.com/99designs/gqlgen/graphql.Marshaler"
Cohesion: 0.06
Nodes (4): NotificationConfigurationConnection, NotificationConfigurationEdge, WorkspaceConnection, WorkspaceEdge

### Community 17 - "stubQuerier"
Cohesion: 0.03
Nodes (7): actionStub, stubQuerier, scopeIDs(), CreateActionParams, UpdateActionParams, Action, Organization

### Community 19 - "tfe_test.go"
Cohesion: 0.09
Nodes (42): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+34 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.08
Nodes (50): replay, upstreamCapture, expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions(), TestAccessTokenOrganizationPermissionsForTheOwner() (+42 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "Workspace"
Cohesion: 0.08
Nodes (12): fakeAccessTokenStore, scheduleStub, stubQuerier, fakeStore, CreateJobParams, CreateScheduleParams, Job, NotificationTrigger (+4 more)

### Community 24 - "newFixture"
Cohesion: 0.08
Nodes (37): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+29 more)

### Community 25 - ".Query"
Cohesion: 0.06
Nodes (16): CountJobsByStatusRow, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow, ListSSHKeySecretsAfterParams (+8 more)

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
Nodes (60): actionFilter, agentAttrs, atomicOperation, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs (+52 more)

### Community 31 - "Link"
Cohesion: 0.06
Nodes (31): parseInclude(), parseList(), strValue(), collectionResource(), itemResource(), referenceResource(), jobResource(), stepResource() (+23 more)

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
Nodes (7): queryResolver, workspaceResolver, organizationScope(), TestOrganizationScope(), idsFrom(), projectEdges(), workspaceEdges()

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

### Community 45 - "Access"
Cohesion: 0.11
Nodes (3): contains(), Access, CreateAccessParams

### Community 46 - "graphql/collections.go"
Cohesion: 0.10
Nodes (22): collectionSpec, DB, listArgs, page, accessEdges(), afterString(), column(), jobEdges() (+14 more)

### Community 47 - "ProjectAccess"
Cohesion: 0.12
Nodes (3): CreateProjectAccessParams, ProjectAccess, UpdateProjectAccessParams

### Community 49 - "logSink"
Cohesion: 0.11
Nodes (10): captureLogPublisher, LogPublisher, logSink, redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues() (+2 more)

### Community 51 - "time.Duration"
Cohesion: 0.10
Nodes (5): fakeQueue, JobQueue, envDuration(), defaultRetryDelay(), Message

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (95): TestAccessTokenWithoutStore(), TestAuthMiddleware(), TestCreateOrganization(), newCollectionAPI(), TestCollectionListCarriesItsItemsAndReferences(), TestCreateCollectionDefaultsPriorityToZero(), TestCreateCollectionItemStoresTheValue(), TestCreateCollectionItemValidatesKeyValueAndCategory() (+87 more)

### Community 53 - "Pat"
Cohesion: 0.23
Nodes (4): fakeStore, Store, CreatePatParams, Pat

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "New"
Cohesion: 0.21
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 56 - ".QueryRow"
Cohesion: 0.05
Nodes (9): CountEnabledRunTriggersBySourceParams, CreateGlobalVariableParams, CreateTagParams, GetTagByNameParams, Globalvar, Ssh, UpdateGlobalVariableParams, UpdateOrganizationParams (+1 more)

### Community 59 - "NewCipher"
Cohesion: 0.13
Nodes (34): formatOldKeys(), logReport(), main(), run(), TestBackfillEncryptsPlaintextRows(), NewCipher(), TestCipherEmptyAndPlaintext(), TestCipherKeyring() (+26 more)

### Community 61 - "Queries"
Cohesion: 0.03
Nodes (7): AddNotificationConfigurationTemplateParams, CreateNotificationTriggerParams, DBTX, Querier, Queries, SetWorkspaceLastJobParams, UpdateHistoryOutputParams

### Community 62 - "History"
Cohesion: 0.14
Nodes (4): reliabilityStore, CreateHistoryParams, GetHistoryByJobReferenceParams, History

### Community 63 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.10
Nodes (3): jobStub, statusQuerier, captureJobs()

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.07
Nodes (48): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+40 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "run"
Cohesion: 0.06
Nodes (52): main(), run(), main(), run(), main(), run(), Config, Task 1: Keyring configuration (+44 more)

### Community 70 - "Team"
Cohesion: 0.09
Nodes (3): contains(), fakeOrganizationStore, Team

### Community 71 - "collections.go"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 72 - "Workspacetag"
Cohesion: 0.17
Nodes (3): CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams, Workspacetag

### Community 73 - "Cipher"
Cohesion: 0.11
Nodes (27): Task 3: Rotation scan and re-encryption, Task 4: `rotatekeys` command, backfillGlobalVariables(), backfillNotifications(), backfillSSHKeys(), backfillVariables(), backfillVCS(), encryptIfNeeded() (+19 more)

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (5): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), NewExecutableSchema(), newExecutionContext(), TestSchemaLoads()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "streamEvents"
Cohesion: 0.10
Nodes (11): notificationConfigurationResolver, organizationResolver, projectResolver, Resolver, ResolverRoot, subscriptionResolver, TestOwnerRule(), streamEvents() (+3 more)

### Community 77 - "templateStub"
Cohesion: 0.13
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
Cohesion: 0.19
Nodes (14): Parse(), dump(), TestParse(), TestParseRejects(), ptr(), seedWorkspaces(), TestListPageBetween(), TestListPageFilters() (+6 more)

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

### Community 87 - "log/slog.Logger"
Cohesion: 0.21
Nodes (8): Validator, Deps, Deps, Backend, Layout, Deps, Deps, Deps

### Community 89 - ".handle"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 91 - "newTestExecutor"
Cohesion: 0.31
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 94 - "Identity"
Cohesion: 0.09
Nodes (16): groupVerifier, PATService, stubDexVerifier, stubIdentityVerifier, stubPATService, fakeVerifier, Identity, fakeIssuer (+8 more)

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
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 102 - ".CreateProject"
Cohesion: 0.20
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 103 - "context.Context"
Cohesion: 0.03
Nodes (9): collectionStub, globalVarStub, sshAgentStub, teamStub, vcsStub, workspaceStub, artifactStore, historyStub (+1 more)

### Community 108 - "Object"
Cohesion: 0.22
Nodes (5): cleanKey(), pathClean(), TestCleanKey(), Object, Local

### Community 118 - "Content"
Cohesion: 0.29
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.10
Nodes (25): runTriggerStub, newRunTriggerStub(), runTriggerAPI(), runTriggerBody(), runTriggerRBACAPI(), TestCreateRunTriggerDerivesOrganizationAndRelationships(), TestCreateRunTriggerEnforcesFanOutLimit(), TestCreateRunTriggerRejectsCycle() (+17 more)

### Community 121 - "newWorkspaceStub"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 123 - "cancelRunStore"
Cohesion: 0.08
Nodes (7): cancelRunStore, CreateStepParams, GetStepByJobAndNumberParams, Step, UpdateJobOutputParams, UpdateStepOutputParams, UpdateStepStatusParams

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 135 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 137 - "time.Time"
Cohesion: 0.36
Nodes (6): bucket, limiter, newLimiter(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), memObject

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 141 - "Dex"
Cohesion: 0.21
Nodes (3): Dex, Exchanger, fakeDex

### Community 142 - "rsql.Node"
Cohesion: 0.13
Nodes (9): exists(), likePattern(), relate(), And, Or, parser, Attr, Kind (+1 more)

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.09
Nodes (37): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+29 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 150 - "API"
Cohesion: 0.07
Nodes (23): API, authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), New(), requireUUIDParams(), unknownCollection() (+15 more)

### Community 151 - "Item"
Cohesion: 0.20
Nodes (3): CreateItemParams, UpdateItemParams, Item

### Community 154 - "NewSubscriptions"
Cohesion: 0.13
Nodes (19): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+11 more)

### Community 157 - ".prepareWorkspace"
Cohesion: 0.15
Nodes (11): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), transient(), contentIDFromSource(), gitEnv(), gitSource(), isConfigurationTarball() (+3 more)

### Community 159 - "Node"
Cohesion: 0.07
Nodes (12): AccessTokenStore, Deps, listQuery, stubLists, Task 3: Handler cutover, CreateRunTriggerParams, Lister, Project (+4 more)

### Community 161 - "Executor"
Cohesion: 0.19
Nodes (5): Executor, preparedWorkspace, isCancelledStatus(), TestBinarySelection(), firstFolder()

### Community 162 - "newService"
Cohesion: 0.19
Nodes (15): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+7 more)

### Community 163 - "CliAuthSession"
Cohesion: 0.07
Nodes (15): fakeStore, Service, Store, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value() (+7 more)

### Community 165 - "auth.DexValidator"
Cohesion: 0.31
Nodes (4): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 166 - ".applyWorkspaceUpdate"
Cohesion: 0.08
Nodes (26): nestedDocument, vcsAttrs, strPtr(), atomicWorkspaceHref(), isAtomicRequest(), decodeVCSAttrs(), requiredVCSString(), validVCSEnum() (+18 more)

### Community 167 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.18
Nodes (9): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+1 more)

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 171 - "authzAPI"
Cohesion: 0.33
Nodes (10): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), strPtr(), TestResolveOrganization() (+2 more)

### Community 172 - "Metrics"
Cohesion: 0.14
Nodes (4): metricsSource, fakeStats, Metrics, StatsSource

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - "models_gen.go"
Cohesion: 0.09
Nodes (7): JobConnection, JobEdge, NotificationTriggerConnection, NotificationTriggerEdge, PageInfo, Query, Subscription

### Community 190 - "loadPage"
Cohesion: 0.25
Nodes (4): failingDB, collectStrings(), snakeCase(), TestModelColumnOrder()

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 199 - "Workspace"
Cohesion: 0.14
Nodes (3): CreateWorkspaceParams, GetWorkspaceByNameParams, SetWorkspaceLockParams

### Community 229 - "log/slog.Logger"
Cohesion: 0.14
Nodes (15): startScheduler(), waitFor(), NewDispatcher(), strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued(), TestTickSkipsLockedWorkspaceAndUnparseableCron() (+7 more)

### Community 231 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 234 - "Migrate"
Cohesion: 0.08
Nodes (26): main(), run(), Locker, Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope (+18 more)

### Community 235 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 512 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **132 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Queries` to `.CreateAgent`, `Job`, `Collection`, `Handler`, `fakeQuerier`, `stubQuerier`, `stubQuerier`, `newTestAPI`, `tfe_test.go`, `Variable`, `Workspace`, `newFixture`, `.Query`, `.CreateArchive`, `Item`, `api/main.go`, `Node`, `CliAuthSession`, `Reference`, `nestedStub`, `Metrics`, `Access`, `ProjectAccess`, `Vc`, `Pat`, `.QueryRow`, `stubStore`, `github.com/aminbenmansour/terrakube-go/internal/store.Step`, `NewCipher`, `tagStub`, `History`, `Team`, `Workspace`, `Workspacetag`, `Cipher`, `.SoftDeletePat`, `templateStub`, `log/slog.Logger`, `log/slog.Logger`, `.CreateProject`, `Migrate`, `Content`, `cancelRunStore`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07692307692307693 - nodes in this community are weakly interconnected._
- **Why does `Layout` connect `log/slog.Logger` to `NewLayout`, `api/main.go`, `M4 run logs and artifacts Implementation Plan`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `New()` connect `API` to `Node`, `NewLayout`, `run`, `.listWorkspaceAccess`, `.applyWorkspaceUpdate`, `notificationTestAPI`, `.internalError`, `authzAPI`, `newTestAPI`, `Parse`, `testing.T`, `strPtr`, `runTriggerAPI`, `newTestExecutor`, `api/main.go`, `Link`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `StringOrInt` be split into smaller, more focused modules?**
  _Cohesion score 0.09982174688057041 - nodes in this community are weakly interconnected._