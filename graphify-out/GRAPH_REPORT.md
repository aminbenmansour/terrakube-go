# Graph Report - terrakube-go  (2026-11-26)

## Corpus Check
- 485 files · ~242,344 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4339 nodes · 17381 edges · 236 communities (113 shown, 123 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3193 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ef819953`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- StringOrInt
- io.Reader
- Core Database Schema
- fakeStore
- collectionStub
- Handler
- ._Project
- .internalError
- fakeQuerier
- stubQuerier
- store.Workspace
- steps
- github.com/99designs/gqlgen/graphql.Marshaler
- stubQuerier
- fakeQuerier
- tfe_test.go
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
- .build
- newRedactor
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- New
- Tag
- github.com/aminbenmansour/terrakube-go/internal/store.Step
- NewCipher
- ResourceIdentifier
- Queries
- reliabilityStore
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- OrganizationConnection
- models
- config.go
- Team
- collections.go
- github.com/jackc/pgx/v5/pgtype.Timestamptz
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
- History
- time.Duration
- newTestExecutor
- JobConnection
- Identity
- environment
- cron.go
- NewSubscriptions
- newTestService
- .CreateProject
- context.Context
- New
- ui-rest-inventory.sh
- Content
- ui-smoke.sh
- cancelRunStore
- newWorkspaceStub
- .QueryRow
- capture.sh
- 0008_action.sql
- .CreateAgent
- compare.sh
- start.sh
- ._NotificationConfiguration
- Review Focus
- time.Time
- tarGzEntries
- ProjectEdge
- parser
- query/rsql.go
- clilogin/handler_test.go
- newTestAPI
- Secret key rotation — design
- .projectGrantParent
- .SetCliAuthSessionIdentity
- net/http.Handler
- actionStub
- .CreateArchive
- NewSubscriptions
- .prepareWorkspace
- Node
- logSink
- newService
- CliAuthSession
- auth.DexValidator
- net/http.ResponseWriter
- notificationTestAPI
- Metrics
- transientError
- Threat model
- models_gen.go
- loadPage
- RSQL filters in SQL — design
- log/slog.Logger
- Review Focus
- .TryLock
- ParsePage

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Queries` - 211 edges
3. `Querier` - 208 edges
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

## Communities (236 total, 123 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (25): WriteResource(), configurationVersionResource(), projectResource(), parseWorkspaceQuery(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID() (+17 more)

### Community 1 - "rbac.go"
Cohesion: 0.19
Nodes (27): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 1: Role matrix and permission set, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract, adminOnly(), CanApproveJob() (+19 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "StringOrInt"
Cohesion: 0.09
Nodes (4): ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, WorkspaceConnection

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (15): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), NewS3() (+7 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 6 - "fakeStore"
Cohesion: 0.06
Nodes (10): Publisher, publishingStore, stubQuerier, fakeStore, ApproveJobParams, CreateJobParams, SetJobOverrideSourceParams, SetWorkspaceLastJobParams (+2 more)

### Community 7 - "collectionStub"
Cohesion: 0.10
Nodes (4): collectionStub, CreateCollectionParams, CreateItemParams, CreateReferenceParams

### Community 8 - "Handler"
Cohesion: 0.05
Nodes (35): Claims, Type, Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log (+27 more)

### Community 10 - ".internalError"
Cohesion: 0.10
Nodes (33): actionAttrs, itemAttrs, notificationConfigurationAttrs, teamAttrs, decodeActionAttrs(), stringOr(), boolValue(), requireType() (+25 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.05
Nodes (7): encryptedStore, CreateGlobalVariableParams, CreateVCSParams, Globalvar, UpdateGlobalVariableParams, UpdateVCSParams, Vc

### Community 13 - "store.Workspace"
Cohesion: 0.05
Nodes (23): AccessTokenStore, Deps, fakeAccessTokenStore, contains(), captureAccessRow(), captureGrantProjects(), captureGrantStore(), captureGrantWorkspaces() (+15 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 17 - "stubQuerier"
Cohesion: 0.03
Nodes (7): stubQuerier, tagStub, scopeIDs(), demoOrganizations(), matchesTagFilters(), tagFilter, workspaceQuery

### Community 18 - "fakeQuerier"
Cohesion: 0.06
Nodes (8): nestedStub, fakeQuerier, CreateNotificationConfigurationParams, CreateSSHKeyParams, Ssh, UpdateNotificationConfigurationParams, UpdateSSHKeyParams, Variable

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

### Community 23 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.13
Nodes (13): scheduleStub, newScheduleAPI(), newScheduleStub(), scheduleBody(), scheduleStr(), TestCreateScheduleCopiesAttributesAndDefaultsEnabled(), TestCreateScheduleRejectsForeignTemplate(), TestDeleteScheduleAndCrossWorkspaceLookup() (+5 more)

### Community 24 - "newFixture"
Cohesion: 0.06
Nodes (46): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+38 more)

### Community 25 - "Service"
Cohesion: 0.16
Nodes (9): Service, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), randomToken(), sanitizeName() (+1 more)

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
Nodes (68): actionFilter, agentAttrs, atomicOperation, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs (+60 more)

### Community 31 - "Link"
Cohesion: 0.07
Nodes (30): actionResource(), parseInclude(), parseList(), strValue(), collectionResource(), itemResource(), referenceResource(), jobResource() (+22 more)

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

### Community 46 - "graphql/collections.go"
Cohesion: 0.18
Nodes (10): page, accessEdges(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), projectAccessEdges(), projectEdges(), templateEdges() (+2 more)

### Community 47 - ".build"
Cohesion: 0.18
Nodes (12): collectionSpec, listArgs, afterString(), column(), predicate(), Order, OrderBy(), ParseSort() (+4 more)

### Community 49 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 51 - "time.Duration"
Cohesion: 0.08
Nodes (5): captureLogPublisher, fakeQueue, JobQueue, LogPublisher, fakePinger

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (101): TestAccessTokenWithoutStore(), TestAuthMiddleware(), TestCreateOrganization(), authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound() (+93 more)

### Community 53 - "Pat"
Cohesion: 0.16
Nodes (6): newFakeStore(), fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "New"
Cohesion: 0.21
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 56 - "Tag"
Cohesion: 0.15
Nodes (4): CreateTagParams, GetTagByNameParams, Tag, UpdateTagParams

### Community 59 - "NewCipher"
Cohesion: 0.06
Nodes (69): main(), run(), main(), run(), main(), run(), main(), run() (+61 more)

### Community 60 - "ResourceIdentifier"
Cohesion: 0.22
Nodes (9): nestedDocument, RelationshipIdentifiers(), RelationshipList(), identifierStringAttribute(), validateTagBindings(), workspaceTagBindings(), RequestRelationship, ResourceIdentifier (+1 more)

### Community 61 - "Queries"
Cohesion: 0.02
Nodes (27): UpdateCollectionParams, AddNotificationConfigurationTemplateParams, Collection, CountJobsByStatusRow, CreateNotificationTriggerParams, CreateWorkspaceParams, DBTX, GetWorkspaceByNameParams (+19 more)

### Community 62 - "reliabilityStore"
Cohesion: 0.19
Nodes (3): reliabilityStore, CreateHistoryParams, UpdateStepStatusParams

### Community 63 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.10
Nodes (3): jobStub, statusQuerier, captureJobs()

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.07
Nodes (49): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+41 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "config.go"
Cohesion: 0.17
Nodes (20): Config, Task 1: Keyring configuration, env(), envBool(), envFirst(), envInt64(), envList(), Load() (+12 more)

### Community 71 - "collections.go"
Cohesion: 0.27
Nodes (8): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), Node, parser

### Community 72 - "github.com/jackc/pgx/v5/pgtype.Timestamptz"
Cohesion: 0.08
Nodes (8): timePtr(), UpdateItemParams, CreateWorkspaceTagParams, IssueCliAuthCodeParams, Item, NotificationTrigger, UpdateWorkspaceTagValueParams, Workspacetag

### Community 73 - "Cipher"
Cohesion: 0.08
Nodes (35): Review Focus, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation, Task 7: Documentation and graph, Task 8: Verification, Task 3: Rotation scan and re-encryption (+27 more)

### Community 74 - "generated.go"
Cohesion: 0.10
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "streamEvents"
Cohesion: 0.10
Nodes (9): projectResolver, Resolver, ResolverRoot, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace() (+1 more)

### Community 77 - "templateStub"
Cohesion: 0.14
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

### Community 91 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 94 - "Identity"
Cohesion: 0.09
Nodes (16): groupVerifier, PATService, stubDexVerifier, stubIdentityVerifier, stubPATService, fakeVerifier, Identity, fakeDex (+8 more)

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 99 - "cron.go"
Cohesion: 0.27
Nodes (8): at(), Parse(), parseField(), splitStep(), mustParse(), TestParseMatchesQuartzFields(), TestParseRejectsWhatQuartzRejects(), Cron

### Community 100 - "NewSubscriptions"
Cohesion: 0.43
Nodes (3): EventSource, fakeEvents, Event

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 102 - ".CreateProject"
Cohesion: 0.18
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 103 - "context.Context"
Cohesion: 0.03
Nodes (8): globalVarStub, sshAgentStub, teamStub, vcsStub, workspaceStub, artifactStore, historyStub, variableStore

### Community 108 - "New"
Cohesion: 0.18
Nodes (10): requireMetricsToken(), key(), RateLimit(), TestRateLimitDisabledByDefault(), TestRateLimitMiddleware(), SecurityHeaders(), TestSecurityHeaders(), New() (+2 more)

### Community 118 - "Content"
Cohesion: 0.29
Nodes (3): CreateContentParams, UpdateContentStatusParams, Content

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 121 - "newWorkspaceStub"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 123 - ".QueryRow"
Cohesion: 0.06
Nodes (6): CreateAccessParams, CreateProjectAccessParams, CreateStepParams, GetStepByJobAndNumberParams, UpdateOrganizationParams, UpdateProjectAccessParams

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 135 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 137 - "time.Time"
Cohesion: 0.43
Nodes (5): bucket, limiter, newLimiter(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets()

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 143 - "query/rsql.go"
Cohesion: 0.19
Nodes (6): likePattern(), And, builder, Compare, Op, Or

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.10
Nodes (36): pageMeta, newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+28 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.14
Nodes (13): Approaches considered, Command, Compatibility and docs, Decisions, Goals, Intent, Non-goals, Operator sequence (+5 more)

### Community 150 - "net/http.Handler"
Cohesion: 0.10
Nodes (21): authorization, Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), New(), requireUUIDParams(), unknownCollection(), authorizationFrom(), isWriteMethod() (+13 more)

### Community 151 - "actionStub"
Cohesion: 0.15
Nodes (4): actionStub, CreateActionParams, UpdateActionParams, Action

### Community 154 - "NewSubscriptions"
Cohesion: 0.11
Nodes (23): Deps, Authenticate(), TestAuthenticatePrefersTheDevelopmentToken(), TestAuthenticateRejectsMissingOrUnknownCredentials(), TestAuthenticateTriesEachVerifier(), TestIdentityRoundTripsThroughContext(), WithIdentity(), TestDexTokenIsAcceptedOnGraphQL() (+15 more)

### Community 157 - ".prepareWorkspace"
Cohesion: 0.14
Nodes (12): gitConfig, preparedWorkspace, TestGitSourceDetection(), TestSanitizeGitURL(), transient(), contentIDFromSource(), gitEnv(), gitSource() (+4 more)

### Community 159 - "Node"
Cohesion: 0.06
Nodes (13): listQuery, stubLists, Task 3: Handler cutover, exists(), likePattern(), relate(), And, Compare (+5 more)

### Community 161 - "logSink"
Cohesion: 0.10
Nodes (10): Executor, logSink, isCancelledStatus(), isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestBinarySelection() (+2 more)

### Community 162 - "newService"
Cohesion: 0.21
Nodes (14): decodeSecret(), New(), dexUser(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp(), TestDeleteEnforcesOwnership() (+6 more)

### Community 163 - "CliAuthSession"
Cohesion: 0.18
Nodes (4): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams

### Community 165 - "auth.DexValidator"
Cohesion: 0.15
Nodes (6): Dex, Exchanger, contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 166 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (19): API, templateAttrs, vcsAttrs, writeAccessTokenJSON(), param(), strPtr(), globalVarResource(), accessResource() (+11 more)

### Community 168 - "notificationTestAPI"
Cohesion: 0.47
Nodes (5): newNotificationStub(), notificationTestAPI(), TestDeleteNotificationConfigurationCascadesTriggers(), TestNotificationConfigurationJSONAPI(), TestNotificationTriggersAndTemplates()

### Community 172 - "Metrics"
Cohesion: 0.14
Nodes (4): metricsSource, fakeStats, Metrics, StatsSource

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - "models_gen.go"
Cohesion: 0.07
Nodes (9): AccessConnection, AccessEdge, NotificationTriggerConnection, NotificationTriggerEdge, PageInfo, Query, Subscription, WorkspaceTagConnection (+1 more)

### Community 190 - "loadPage"
Cohesion: 0.18
Nodes (6): DB, failingDB, loadPage(), collectStrings(), snakeCase(), TestModelColumnOrder()

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 229 - "log/slog.Logger"
Cohesion: 0.13
Nodes (16): startScheduler(), waitFor(), formatOldKeys(), logReport(), NewDispatcher(), strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond() (+8 more)

### Community 231 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 234 - ".TryLock"
Cohesion: 0.29
Nodes (3): Deps, fakeLocker, JobLocker

### Community 235 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 511 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **123 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Querier` connect `Queries` to `.CreateAgent`, `fakeStore`, `collectionStub`, `Handler`, `fakeQuerier`, `stubQuerier`, `store.Workspace`, `.marshalOBoolean2bool`, `github.com/99designs/gqlgen/graphql.Marshaler`, `stubQuerier`, `fakeQuerier`, `newTestAPI`, `tfe_test.go`, `.SetCliAuthSessionIdentity`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `newFixture`, `actionStub`, `.CreateArchive`, `api/main.go`, `CliAuthSession`, `Metrics`, `Pat`, `Tag`, `github.com/aminbenmansour/terrakube-go/internal/store.Step`, `NewCipher`, `reliabilityStore`, `Team`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `Cipher`, `templateStub`, `History`, `log/slog.Logger`, `.CreateProject`, `context.Context`, `Content`, `.QueryRow`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.08146201624462494 - nodes in this community are weakly interconnected._
- **Why does `New()` connect `net/http.Handler` to `NewLayout`, `net/http.ResponseWriter`, `notificationTestAPI`, `.internalError`, `newTestExecutor`, `store.Workspace`, `newTestAPI`, `.projectGrantParent`, `testing.T`, `strPtr`, `github.com/jackc/pgx/v5/pgtype.Timestamptz`, `NewCipher`, `api/main.go`, `Link`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `openRotationTest()` connect `NewCipher` to `context.Context`, `newTestAPI`, `testing.T`, `Queries`, `api/main.go`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `StringOrInt` be split into smaller, more focused modules?**
  _Cohesion score 0.08771929824561403 - nodes in this community are weakly interconnected._