# Graph Report - terrakube-go  (2026-11-26)

## Corpus Check
- 466 files · ~231,952 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4173 nodes · 16738 edges · 226 communities (115 shown, 111 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 3116 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `cff4e7be`
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
- net/http.Request
- encryptedStore
- .Query
- Layout
- steps
- CreateNotificationConfigurationParams
- stubQuerier
- encryptedStore
- tfe_test.go
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- Node
- newFixture
- Identity
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- api/main.go
- Resource
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
- ._Job
- graphql/collections.go
- Service
- logSink
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- .handle
- github.com/jackc/pgx/v5/pgtype.Timestamptz
- github.com/aminbenmansour/terrakube-go/internal/store.Project
- loadPage
- WriteError
- Querier
- Claims
- stubQuerier
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- models
- config.go
- store.Team
- collections.go
- executionContext
- Cipher
- generated.go
- where_test.go
- Resolver
- templateStub
- decodeJSON
- graphql/models.go
- newJWKSServer
- Pool
- Secret encryption at rest — design
- ui.go
- REST contract, captured from upstream
- pat.py
- query/rsql.go
- New
- Job
- ._Template
- environment
- M4 run logs and artifacts Implementation Plan
- NewSubscriptions
- newTestService
- ProjectAccessConnection
- context.Context
- ListSpec
- ui-rest-inventory.sh
- grantsAPI
- ui-smoke.sh
- cancelRunStore
- Parse
- .QueryRow
- capture.sh
- 0008_action.sql
- .CreateAgent
- compare.sh
- start.sh
- Workspace
- tarGzEntries
- New
- Parse
- github.com/aminbenmansour/terrakube-go/internal/store.Organization
- net/http.Handler
- run
- Secret key rotation — design
- reliabilityStore
- net/http.Handler
- serveNested
- .CreateArchive
- NewSubscriptions
- .createVCS
- Dex
- .cloneConfiguration
- Review Focus
- newService
- logSink
- ResourceIdentifier
- store.Job
- NotificationTriggerConnection
- auth.DexValidator
- query.go
- notificationTestAPI
- ParsePage
- .Ping
- transientError
- .SoftDeletePat
- newRedactor
- github.com/vektah/gqlparser/v2/ast.SelectionSet
- stubStore
- .CreateProject
- RSQL filters in SQL — design
- Deps
- authzAPI

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Queries` - 204 edges
3. `Querier` - 195 edges
4. `API` - 162 edges
5. `New()` - 145 edges
6. `stubQuerier` - 142 edges
7. `WriteError()` - 119 edges
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
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Task 3: Handler cutover` --references--> `Parse()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/rsql/rsql.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (226 total, 111 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (23): workspaceUpdateParams(), Write(), WriteResource(), configurationVersionResource(), approvalGateStatus(), failedJobStatus(), parseRunID(), stepStatusTFE() (+15 more)

### Community 1 - "rbac.go"
Cohesion: 0.19
Nodes (27): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, adminOnly(), CanApproveJob() (+19 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (15): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+7 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "tagStub"
Cohesion: 0.07
Nodes (3): tagStub, CreateWorkspaceTagParams, UpdateWorkspaceTagValueParams

### Community 7 - "collectionStub"
Cohesion: 0.08
Nodes (4): collectionStub, CreateCollectionParams, CreateItemParams, CreateReferenceParams

### Community 8 - "Handler"
Cohesion: 0.24
Nodes (6): isExternalUser(), matchesMD5(), New(), newUUID(), pathIDs(), Handler

### Community 10 - "net/http.Request"
Cohesion: 0.12
Nodes (8): API, New(), param(), parseList(), unknownCollection(), writeInUseConflict(), ParseFieldsets(), WriteCollection()

### Community 12 - ".Query"
Cohesion: 0.07
Nodes (16): DBTX, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow, ListSSHKeySecretsAfterParams (+8 more)

### Community 13 - "Layout"
Cohesion: 0.21
Nodes (5): Validator, Layout, Deps, JobQueue, Deps

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 17 - "stubQuerier"
Cohesion: 0.04
Nodes (3): stubQuerier, UpdateCollectionParams, UpdateItemParams

### Community 18 - "encryptedStore"
Cohesion: 0.06
Nodes (6): nestedStub, fakeQuerier, CreateSSHKeyParams, CreateVCSParams, UpdateSSHKeyParams, UpdateVCSParams

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

### Community 23 - "Node"
Cohesion: 0.07
Nodes (16): listQuery, stubLists, Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs (+8 more)

### Community 24 - "newFixture"
Cohesion: 0.15
Nodes (23): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+15 more)

### Community 25 - "Identity"
Cohesion: 0.09
Nodes (17): AccessTokenStore, Deps, groupVerifier, PATService, stubDexVerifier, stubIdentityVerifier, stubPATService, fakeVerifier (+9 more)

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
Nodes (59): actionFilter, agentAttrs, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs, grantAttrs (+51 more)

### Community 31 - "Resource"
Cohesion: 0.09
Nodes (37): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), accessResource(), projectAccessResource() (+29 more)

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
Nodes (9): notificationConfigurationResolver, organizationResolver, queryResolver, workspaceResolver, organizationScope(), TestOrganizationScope(), idsFrom(), projectEdges() (+1 more)

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
Cohesion: 0.11
Nodes (20): collectionSpec, listArgs, page, accessEdges(), afterString(), column(), jobEdges(), notificationConfigurationEdges() (+12 more)

### Community 48 - "Service"
Cohesion: 0.15
Nodes (10): Service, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), randomToken() (+2 more)

### Community 49 - "logSink"
Cohesion: 0.15
Nodes (6): captureLogPublisher, LogPublisher, logSink, redactor, TestNilRedactorIsSafe(), baseEnvironment()

### Community 51 - "time.Duration"
Cohesion: 0.10
Nodes (7): fakeQueue, JobQueue, envDuration(), defaultRetryDelay(), Message, Queue, stubJobQueue

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (97): TestAccessTokenWithoutStore(), TestAuthMiddleware(), TestCreateOrganization(), newCollectionAPI(), TestCollectionListCarriesItsItemsAndReferences(), TestCreateCollectionDefaultsPriorityToZero(), TestCreateCollectionItemStoresTheValue(), TestCreateCollectionItemValidatesKeyValueAndCategory() (+89 more)

### Community 53 - "Pat"
Cohesion: 0.21
Nodes (5): newFakeStore(), fakeStore, Store, CreatePatParams, Pat

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - ".handle"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 56 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.08
Nodes (6): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams, SetCliAuthSessionIdentityParams

### Community 57 - "github.com/aminbenmansour/terrakube-go/internal/store.Project"
Cohesion: 0.08
Nodes (14): contains(), ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveOrganization(), TestResolveOrganizationUsesTheRequestedOrganization() (+6 more)

### Community 59 - "loadPage"
Cohesion: 0.18
Nodes (6): DB, failingDB, loadPage(), collectStrings(), snakeCase(), TestModelColumnOrder()

### Community 60 - "WriteError"
Cohesion: 0.10
Nodes (34): actionAttrs, atomicOperation, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, decodeActionAttrs(), stringOr() (+26 more)

### Community 61 - "Querier"
Cohesion: 0.04
Nodes (6): AddNotificationConfigurationTemplateParams, CreateNotificationTriggerParams, Querier, Queries, SetWorkspaceLastJobParams, UpdateHistoryOutputParams

### Community 62 - "Claims"
Cohesion: 0.13
Nodes (16): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken() (+8 more)

### Community 63 - "stubQuerier"
Cohesion: 0.04
Nodes (8): jobStub, statusQuerier, captureJobs(), captureSteps(), CreateContentParams, UpdateContentStatusParams, stubQuerier, stubQuerier

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.08
Nodes (43): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+35 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "config.go"
Cohesion: 0.17
Nodes (20): Config, Task 1: Keyring configuration, env(), envBool(), envFirst(), envInt64(), envList(), Load() (+12 more)

### Community 70 - "store.Team"
Cohesion: 0.07
Nodes (4): fakeAccessTokenStore, captureWorkspaces(), contains(), fakeOrganizationStore

### Community 71 - "collections.go"
Cohesion: 0.27
Nodes (8): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), Node, parser

### Community 73 - "Cipher"
Cohesion: 0.09
Nodes (32): Review Focus, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation, Task 7: Documentation and graph, Task 8: Verification, Task 3: Rotation scan and re-encryption (+24 more)

### Community 74 - "generated.go"
Cohesion: 0.09
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
Cohesion: 0.10
Nodes (42): newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent(), TestGlobalVariableFromAnotherOrganizationIsNotFound(), TestListGlobalVariablesHidesSensitiveValues() (+34 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "Pool"
Cohesion: 0.10
Nodes (34): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+26 more)

### Community 82 - "Secret encryption at rest — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Backfill, Columns, Compatibility and docs, Decisions, Envelope, Goals, Intent (+5 more)

### Community 83 - "ui.go"
Cohesion: 0.18
Nodes (14): envConfig(), exists(), New(), get(), newUI(), TestDeepLinksServeTheIndex(), TestEnvConfigIsNeverCached(), TestMissingAssetsAreNotFound() (+6 more)

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 87 - "query/rsql.go"
Cohesion: 0.19
Nodes (6): likePattern(), And, builder, Compare, Op, Or

### Community 89 - "New"
Cohesion: 0.21
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 91 - "Job"
Cohesion: 0.11
Nodes (5): ApproveJobParams, CreateJobParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams, UpdateJobStatusParams

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 99 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.18
Nodes (9): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+1 more)

### Community 100 - "NewSubscriptions"
Cohesion: 0.33
Nodes (3): fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 103 - "context.Context"
Cohesion: 0.04
Nodes (8): globalVarStub, sshAgentStub, teamStub, vcsStub, workspaceStub, artifactStore, historyStub, variableStore

### Community 108 - "ListSpec"
Cohesion: 0.23
Nodes (8): exists(), likePattern(), relate(), Compare, Op, Attr, Kind, ListSpec

### Community 118 - "grantsAPI"
Cohesion: 0.19
Nodes (15): doGrant(), grantsAPI(), TestProjectAccessCreateAllowsMissingName(), TestProjectAccessDelete(), TestProjectAccessRejectsAnotherOrganization(), TestWorkspaceAccessCreateDefaults(), TestWorkspaceAccessCreateEmptyName(), TestWorkspaceAccessCreateWithoutNameConflicts() (+7 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "cancelRunStore"
Cohesion: 0.15
Nodes (3): cancelRunStore, UpdateJobOutputParams, UpdateStepOutputParams

### Community 123 - ".QueryRow"
Cohesion: 0.04
Nodes (10): actionStub, CreateActionParams, UpdateActionParams, CreateAccessParams, CreateProjectAccessParams, CreateTagParams, GetTagByNameParams, UpdateOrganizationParams (+2 more)

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 133 - "Workspace"
Cohesion: 0.14
Nodes (3): CreateWorkspaceParams, GetWorkspaceByNameParams, SetWorkspaceLockParams

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 141 - "New"
Cohesion: 0.20
Nodes (12): TestWriteBackendOverride(), isTerminal(), New(), outputOf(), stepOf(), streamKey(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame() (+4 more)

### Community 142 - "Parse"
Cohesion: 0.21
Nodes (13): Parse(), dump(), TestParse(), TestParseRejects(), ptr(), seedWorkspaces(), TestListPageBetween(), TestListPageFilters() (+5 more)

### Community 143 - "github.com/aminbenmansour/terrakube-go/internal/store.Organization"
Cohesion: 0.18
Nodes (3): scopeIDs(), demoOrganizations(), CreateOrganizationParams

### Community 144 - "net/http.Handler"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "run"
Cohesion: 0.06
Nodes (72): main(), run(), main(), run(), main(), run(), main(), run() (+64 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 147 - "reliabilityStore"
Cohesion: 0.07
Nodes (6): reliabilityStore, CreateHistoryParams, CreateStepParams, GetHistoryByJobReferenceParams, GetStepByJobAndNumberParams, UpdateStepStatusParams

### Community 150 - "net/http.Handler"
Cohesion: 0.08
Nodes (18): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON(), requireUUIDParams(), authorizationFrom(), isWriteMethod(), permissionDeniedDetail() (+10 more)

### Community 151 - "serveNested"
Cohesion: 0.33
Nodes (10): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newNestedStub(), serveNested(), TestNestedVariableCreateUsesTheWorkspaceInThePath(), TestNestedVariableDeleteRemovesOnlyThatVariable(), TestNestedVariableFromAnotherWorkspaceIsNotFound(), TestNestedVariablesListIsScopedToTheWorkspace() (+2 more)

### Community 154 - "NewSubscriptions"
Cohesion: 0.13
Nodes (20): Deps, EventSource, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity() (+12 more)

### Community 155 - ".createVCS"
Cohesion: 0.39
Nodes (6): vcsAttrs, strPtr(), strValue(), decodeVCSAttrs(), validVCSEnum(), vcsEnumOrDefault()

### Community 156 - "Dex"
Cohesion: 0.21
Nodes (3): Dex, Exchanger, fakeDex

### Community 157 - ".cloneConfiguration"
Cohesion: 0.12
Nodes (13): gitConfig, preparedWorkspace, TestGitSourceDetection(), TestSanitizeGitURL(), TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv() (+5 more)

### Community 159 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 160 - "newService"
Cohesion: 0.21
Nodes (14): decodeSecret(), New(), dexUser(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp(), TestDeleteEnforcesOwnership() (+6 more)

### Community 161 - "logSink"
Cohesion: 0.18
Nodes (4): Executor, isCancelledStatus(), stringPointer(), TestBinarySelection()

### Community 162 - "ResourceIdentifier"
Cohesion: 0.22
Nodes (9): nestedDocument, RelationshipIdentifiers(), RelationshipList(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship, ResourceIdentifier (+1 more)

### Community 163 - "store.Job"
Cohesion: 0.33
Nodes (3): Publisher, publishingStore, stubQuerier

### Community 165 - "auth.DexValidator"
Cohesion: 0.31
Nodes (4): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 166 - "query.go"
Cohesion: 0.32
Nodes (4): parseWorkspaceQuery(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 171 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

### Community 180 - "newRedactor"
Cohesion: 0.50
Nodes (4): TestLogSinkRedactsSecrets(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 183 - "github.com/vektah/gqlparser/v2/ast.SelectionSet"
Cohesion: 0.04
Nodes (15): AccessConnection, AccessEdge, ComplexityRoot, NotificationConfigurationEdge, OrganizationConnection, OrganizationEdge, PageInfo, ProjectConnection (+7 more)

### Community 191 - ".CreateProject"
Cohesion: 0.18
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 202 - "Deps"
Cohesion: 0.29
Nodes (3): Deps, fakeLocker, JobLocker

### Community 204 - "authzAPI"
Cohesion: 0.23
Nodes (10): authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite(), strPtr(), nextFrame() (+2 more)

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 478 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **111 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Querier` to `.CreateAgent`, `.CreateGlobalVariable`, `Workspace`, `tagStub`, `collectionStub`, `encryptedStore`, `.Query`, `Layout`, `New`, `github.com/aminbenmansour/terrakube-go/internal/store.Organization`, `CreateNotificationConfigurationParams`, `run`, `encryptedStore`, `stubQuerier`, `tfe_test.go`, `reliabilityStore`, `._Template`, `.CreateArchive`, `Identity`, `api/main.go`, `.SoftDeletePat`, `.marshalOBoolean2bool`, `testing.T`, `Pat`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `stubStore`, `github.com/aminbenmansour/terrakube-go/internal/store.Project`, `Claims`, `stubQuerier`, `.CreateProject`, `store.Team`, `Cipher`, `templateStub`, `Pool`, `Job`, `store.Access`, `context.Context`, `cancelRunStore`, `.QueryRow`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.08074534161490683 - nodes in this community are weakly interconnected._
- **Why does `executionContext` connect `executionContext` to `github.com/99designs/gqlgen/graphql.CollectedField`, `.childFields___InputValue`, `NotificationTriggerConnection`, `ProjectAccessConnection`, `context.Context`, `github.com/99designs/gqlgen/graphql.Marshaler`, `generated.go`, `Resolver`, `._Job`, `.marshalOBoolean2bool`, `.childFields___Type`, `._Template`, `github.com/vektah/gqlparser/v2/ast.SelectionSet`, `._Template`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `API` connect `net/http.Request` to `net/http.Handler`, `Identity`, `.createVCS`, `WriteError`, `api/main.go`, `Resource`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `io.Reader` be split into smaller, more focused modules?**
  _Cohesion score 0.06753246753246753 - nodes in this community are weakly interconnected._