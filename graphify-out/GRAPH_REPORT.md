# Graph Report - terrakube-go  (2026-11-27)

## Corpus Check
- 503 files · ~258,982 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 4576 nodes · 18588 edges · 258 communities (127 shown, 131 thin omitted)
- Extraction: 82% EXTRACTED · 18% INFERRED · 0% AMBIGUOUS · INFERRED: 3377 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9b956a99`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- rbac.go
- Agent Workflows and Compatibility Docs
- ._Project
- io.Reader
- Core Database Schema
- github.com/vektah/gqlparser/v2/ast.SelectionSet
- Rotate
- net/http.ResponseWriter
- newFixture
- .internalError
- fakeQuerier
- github.com/aminbenmansour/terrakube-go/internal/store.Job
- Service
- steps
- github.com/aminbenmansour/terrakube-go/internal/rsql.Node
- deliver_test.go
- tfe_test.go
- sqlc Generation Config
- strPtr
- Upstream Elide entities (GraphQL surface)
- fakeStore
- templateStub
- payload_test.go
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
- stubQuerier
- Claims
- graphql/collections.go
- github.com/aminbenmansour/terrakube-go/internal/store.Workspace
- CliAuthSession
- JobConnection
- Go Module Definition
- time.Duration
- testing.T
- Pat
- executor-e2e.sh
- New
- .notificationConfigurationByID
- .QueryRow
- Queries
- History
- tagStub
- dex-config.yaml
- github.com/99designs/gqlgen/graphql.CollectedField
- NewLayout
- WorkspaceTagConnection
- models
- config.go
- Team
- collections.go
- rotateValue
- Cipher
- generated.go
- where_test.go
- streamEvents
- .prepareWorkspace
- decodeJSON
- graphql/models.go
- newJWKSServer
- logSink
- Secret encryption at rest — design
- ui.go
- REST contract, captured from upstream
- pat.py
- log/slog.Logger
- .createVCS
- newTestExecutor
- Workspacetag
- Identity
- environment
- cron.go
- NewSubscriptions
- newTestService
- context.Context
- newService
- ui-rest-inventory.sh
- ui-smoke.sh
- runTriggerAPI
- newWorkspaceStub
- capture.sh
- 0008_action.sql
- time.Duration
- compare.sh
- start.sh
- .CreateGlobalVariable
- Review Focus
- time.Time
- tarGzEntries
- .processJob
- rsql.Node
- where.go
- clilogin/handler_test.go
- newTestAPI
- Secret key rotation — design
- AccessConnection
- API
- actionStub
- addDestination
- run
- Dispatcher
- stubQuerier
- Executor
- Item
- .CreateProject
- .streamStepLog
- Dex
- OrganizationConnection
- NewAccessToken
- notificationSendAPI
- ResolveOrganization
- Collection
- Metrics
- transientError
- Deps
- Threat model
- models_gen.go
- auth.DexValidator
- loadPage
- M4 run logs and artifacts Implementation Plan
- RSQL filters in SQL — design
- Review Focus
- Workspace
- .Ping
- .projectGrantParent
- Schedule
- Review Focus
- Migrate
- .build
- .Query
- .List
- ParsePage
- .validateRunTriggerGraph
- .ListVCSSecretsAfter
- .SetGlobalVariableSecret
- .SetNotificationSecret
- .SetSSHKeySecret
- .SetVariableSecret
- .SetVCSSecrets

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Queries` - 227 edges
3. `Querier` - 226 edges
4. `API` - 186 edges
5. `New()` - 159 edges
6. `stubQuerier` - 154 edges
7. `WriteError()` - 133 edges
8. `param()` - 84 edges
9. `ParseFieldsets()` - 72 edges
10. `API` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Task 2: Store list builder and resource specs` --references--> `Pool()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-rsql-sql-filters.md → internal/dbtest/dbtest.go
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `Task 5: Backfill implementation` --references--> `needsEncryption()`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/secretstore/backfill.go
- `RBAC: permission resolution and `/access-token/v1` Implementation Plan` --references--> `Flags`  [INFERRED]
  docs/superpowers/plans/2026-11-24-rbac-permission-resolution.md → internal/rbac/rbac.go
- `Task 3: Transparent store decorator` --references--> `fakeQuerier`  [INFERRED]
  docs/superpowers/plans/2026-11-25-secret-encryption.md → internal/secretstore/store_test.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (258 total, 131 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (31): Write(), configurationVersionResource(), projectResource(), parseWorkspaceQuery(), activeRunStatus(), approvalGateStatus(), failedJobStatus(), parseRunID() (+23 more)

### Community 1 - "rbac.go"
Cohesion: 0.26
Nodes (22): Task 1: Role matrix and permission set, adminOnly(), CanApproveJob(), CanDestroyJob(), CanManageCollection(), CanManageJob(), CanManageModule(), CanManagePolicies() (+14 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (11): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+3 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.06
Nodes (54): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+46 more)

### Community 7 - "Rotate"
Cohesion: 0.20
Nodes (16): Command, Decisions, Operator sequence, Per-value action, Report and verification, Scan, Rotate(), assertSealedUnder() (+8 more)

### Community 8 - "net/http.ResponseWriter"
Cohesion: 0.14
Nodes (11): workspaceUpdateParams(), size(), writeTooLarge(), RelationshipIdentifier(), isExternalUser(), matchesMD5(), New(), newUUID() (+3 more)

### Community 9 - "newFixture"
Cohesion: 0.13
Nodes (24): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+16 more)

### Community 10 - ".internalError"
Cohesion: 0.09
Nodes (41): actionAttrs, itemAttrs, notificationConfigurationAttrs, teamAttrs, templateAttrs, decodeActionAttrs(), stringOr(), boolValue() (+33 more)

### Community 11 - "fakeQuerier"
Cohesion: 0.04
Nodes (11): encryptedStore, fakeQuerier, CreateNotificationConfigurationParams, CreateSSHKeyParams, CreateVCSParams, Ssh, UpdateNotificationConfigurationParams, UpdateSSHKeyParams (+3 more)

### Community 12 - "github.com/aminbenmansour/terrakube-go/internal/store.Job"
Cohesion: 0.08
Nodes (4): jobStub, statusQuerier, captureJobs(), stubQuerier

### Community 13 - "Service"
Cohesion: 0.15
Nodes (10): Service, TokenIssuer, TokenResponse, cliRedirect(), decodeGroups(), hashCode(), int32Value(), randomToken() (+2 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (44): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+36 more)

### Community 17 - "github.com/aminbenmansour/terrakube-go/internal/rsql.Node"
Cohesion: 0.10
Nodes (4): listQuery, stubLists, Task 3: Handler cutover, Tag

### Community 18 - "deliver_test.go"
Cohesion: 0.17
Nodes (18): deliveryError(), stringPtr(), TestDeliverDoesNotFollowRedirects(), testDeliverer(), TestDeliverIgnoresMalformedRetryAfter(), TestDeliverNetworkErrorIsRetryableAndHidesTheURL(), TestDeliverParsesHTTPDateRetryAfter(), TestDeliverPostsPayload() (+10 more)

### Community 19 - "tfe_test.go"
Cohesion: 0.08
Nodes (43): assertContractJSONAPI(), TestContractDiscoveryAndPing(), TestContractOrganizationEnvelope(), TestContractWorkspaceAndRunEnvelopes(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI() (+35 more)

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "strPtr"
Cohesion: 0.08
Nodes (51): replay, stubIdentityVerifier, upstreamCapture, expectPermissions(), serveAccessToken(), TestAccessTokenCurrentTeams(), TestAccessTokenCurrentTeamsWithoutGroups(), TestAccessTokenOrganizationPermissions() (+43 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.15
Nodes (12): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+4 more)

### Community 23 - "fakeStore"
Cohesion: 0.04
Nodes (18): stubQuerier, fakeStore, fakeStore, CreateJobParams, CreateStepParams, FailStuckRunTriggerEventsParams, GetHistoryBeforeParams, GetHistoryByJobReferenceParams (+10 more)

### Community 24 - "templateStub"
Cohesion: 0.14
Nodes (3): templateStub, CreateTemplateParams, UpdateTemplateParams

### Community 25 - "payload_test.go"
Cohesion: 0.19
Nodes (26): BuildPayload(), buildSlack(), buildTeams(), buildWebhook(), capitalize(), marshalPayload(), notBlank(), nullIfBlank() (+18 more)

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
Nodes (86): actionFilter, agentAttrs, atomicOperation, atomicResult, authorizationKey, collectionAttrs, createPATRequest, globalVarAttrs (+78 more)

### Community 31 - "Link"
Cohesion: 0.07
Nodes (30): actionResource(), parseInclude(), collectionResource(), itemResource(), referenceResource(), globalVarResource(), jobResource(), stepResource() (+22 more)

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
Cohesion: 0.16
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

### Community 44 - "stubQuerier"
Cohesion: 0.11
Nodes (4): CreateContentParams, UpdateContentStatusParams, Content, stubQuerier

### Community 45 - "Claims"
Cohesion: 0.15
Nodes (15): Claims, Type, ClaimsFromContext(), Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken() (+7 more)

### Community 46 - "graphql/collections.go"
Cohesion: 0.18
Nodes (10): page, accessEdges(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), projectAccessEdges(), projectEdges(), templateEdges() (+2 more)

### Community 47 - "github.com/aminbenmansour/terrakube-go/internal/store.Workspace"
Cohesion: 0.12
Nodes (3): nestedStub, demoOrganizations(), workspaceWithTags

### Community 48 - "CliAuthSession"
Cohesion: 0.11
Nodes (6): fakeStore, Store, CliAuthSession, CreateCliAuthSessionParams, IssueCliAuthCodeParams, SetCliAuthSessionIdentityParams

### Community 52 - "testing.T"
Cohesion: 0.03
Nodes (105): TestAccessTokenWithoutStore(), newActionStub(), TestCreateActionRequiresAnIDAndAttributes(), TestCreateActionStoresTheCatalogRow(), TestDeleteAction(), TestGetActionAcceptsSlugID(), TestListActionsFiltersTheDrawerExpression(), TestUpdateActionKeepsOmittedAttributes() (+97 more)

### Community 53 - "Pat"
Cohesion: 0.17
Nodes (5): fakeStore, Store, CreatePatParams, Pat, SoftDeletePatParams

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "New"
Cohesion: 0.20
Nodes (9): Deps, ConsentPage(), ErrorPage(), effectivePort(), New(), originMatches(), writeHTML(), writeJSON() (+1 more)

### Community 57 - ".notificationConfigurationByID"
Cohesion: 0.14
Nodes (8): nestedDocument, notificationConfigurationResource(), notificationTriggerResource(), RelationshipList(), TestRelationshipList(), identifierStringAttribute(), RequestRelationship, ResourceIdentifier

### Community 60 - ".QueryRow"
Cohesion: 0.04
Nodes (11): ApproveJobParams, CountEnabledRunTriggersBySourceParams, CreateNotificationTriggerParams, CreateTagParams, GetJobByWorkspaceAndTriggerParams, GetTagByNameParams, SetJobOverrideSourceParams, UpdateJobPlanChangesParams (+3 more)

### Community 61 - "Queries"
Cohesion: 0.03
Nodes (11): CreateReferenceParams, AddNotificationConfigurationTemplateParams, Agent, CreateAgentParams, DBTX, Querier, Queries, Reference (+3 more)

### Community 62 - "History"
Cohesion: 0.10
Nodes (5): cancelRunStore, reliabilityStore, CreateHistoryParams, UpdateStepStatusParams, mutableQuerier

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - "NewLayout"
Cohesion: 0.05
Nodes (67): redactor, gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive() (+59 more)

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 69 - "config.go"
Cohesion: 0.15
Nodes (24): Config, S3, Storage, Task 2: Serve archived logs from `/tfoutput/v1`, Task 1: Keyring configuration, env(), envBool(), envDuration() (+16 more)

### Community 70 - "Team"
Cohesion: 0.04
Nodes (11): AccessTokenStore, fakeAccessTokenStore, contains(), contains(), fakeOrganizationStore, Access, CreateAccessParams, CreateProjectAccessParams (+3 more)

### Community 71 - "collections.go"
Cohesion: 0.18
Nodes (11): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Node (+3 more)

### Community 72 - "rotateValue"
Cohesion: 0.42
Nodes (11): envelopeKey(), rotateGlobalVariables(), rotateNotifications(), rotatePtr(), rotateSSHKeys(), rotateValue(), rotateVariables(), rotateVCS() (+3 more)

### Community 73 - "Cipher"
Cohesion: 0.37
Nodes (9): Task 3: Rotation scan and re-encryption, backfillGlobalVariables(), backfillNotifications(), backfillSSHKeys(), backfillVariables(), backfillVCS(), encryptIfNeeded(), needsEncryption() (+1 more)

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

### Community 76 - "streamEvents"
Cohesion: 0.11
Nodes (9): projectResolver, Resolver, ResolverRoot, subscriptionResolver, TestOwnerRule(), streamEvents(), TestJobStatusSubscriptionClosesWithTheClient(), TestJobStatusSubscriptionKeepsOnlyTheRequestedWorkspace() (+1 more)

### Community 77 - ".prepareWorkspace"
Cohesion: 0.16
Nodes (11): gitConfig, TestGitSourceDetection(), TestSanitizeGitURL(), transient(), contentIDFromSource(), gitEnv(), gitSource(), isConfigurationTarball() (+3 more)

### Community 78 - "decodeJSON"
Cohesion: 0.12
Nodes (31): TestSparseFieldsetsForAnotherTypeOmitAttributes(), TestSparseFieldsetsRestrictAttributes(), newGlobalVarStub(), serveGlobalVar(), TestCreateGlobalVariableHidesASensitiveValue(), TestCreateGlobalVariableRejectsALowercaseCategory(), TestCreateGlobalVariableWithoutAKeyConflicts(), TestDeleteGlobalVariableAnswersNoContent() (+23 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - "newJWKSServer"
Cohesion: 0.15
Nodes (20): jwksServer, dexFixture, NewDexValidator(), dexClaims(), newJWKSServer(), newRSAKey(), signDexToken(), TestDexValidatorAcceptsAValidToken() (+12 more)

### Community 81 - "logSink"
Cohesion: 0.15
Nodes (4): captureLogPublisher, LogPublisher, logSink, baseEnvironment()

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
Cohesion: 0.16
Nodes (13): Validator, startScheduler(), waitFor(), Deps, NewDispatcher(), NewDispatcher(), NewS3(), Backend (+5 more)

### Community 89 - ".createVCS"
Cohesion: 0.39
Nodes (6): vcsAttrs, strPtr(), strValue(), decodeVCSAttrs(), validVCSEnum(), vcsEnumOrDefault()

### Community 91 - "newTestExecutor"
Cohesion: 0.31
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 94 - "Identity"
Cohesion: 0.12
Nodes (11): groupVerifier, stubDexVerifier, stubPATService, fakeVerifier, Identity, fakeIssuer, groupVerifier, stubDexVerifier (+3 more)

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 99 - "cron.go"
Cohesion: 0.27
Nodes (8): at(), Parse(), parseField(), splitStep(), mustParse(), TestParseMatchesQuartzFields(), TestParseRejectsWhatQuartzRejects(), Cron

### Community 100 - "NewSubscriptions"
Cohesion: 0.33
Nodes (3): fakeEvents, Event, Bus

### Community 101 - "newTestService"
Cohesion: 0.15
Nodes (21): AuthorizeRequest, ChallengeS256(), NewCodeVerifier(), TestNewCodeVerifier(), TestPKCEMatchesTheRFCVector(), TestValidateLoopbackRedirect(), ValidateLoopbackRedirect(), VerifyS256() (+13 more)

### Community 103 - "context.Context"
Cohesion: 0.03
Nodes (11): collectionStub, globalVarStub, sshAgentStub, teamStub, vcsStub, workspaceStub, artifactStore, historyStub (+3 more)

### Community 108 - "newService"
Cohesion: 0.19
Nodes (15): decodeSecret(), New(), dexUser(), newFakeStore(), newService(), TestCreateRejectsNonInteractiveUsersAndBadLifetimes(), TestCreateSignsTheUpstreamToken(), TestCreateWithoutExpiryHasNoExp() (+7 more)

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 120 - "runTriggerAPI"
Cohesion: 0.12
Nodes (24): missingOrganizationStub, runTriggerStub, authzAPI(), authzRequest(), TestOrganizationListIsScopedToTheCaller(), TestOrganizationPermissionMatrix(), TestProjectGrantOfAnotherOrganizationIsNotFound(), TestWorkspaceGrantAllowsTheWrite() (+16 more)

### Community 121 - "newWorkspaceStub"
Cohesion: 0.19
Nodes (21): decodeAtomicResults(), serveAtomic(), TestAtomicOperationCreatesAWorkspace(), TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(), TestAtomicOperationRejectsAnUnsupportedOp(), TestAtomicOperationRequiresTheAtomicMediaType(), TestAtomicOperationUpdatesAWorkspace(), TestAtomicOperationWithNoOperationsReturnsNoResults() (+13 more)

### Community 124 - "capture.sh"
Cohesion: 0.60
Nodes (5): pg(), request(), request_as(), capture.sh script, unauthenticated()

### Community 128 - "time.Duration"
Cohesion: 0.15
Nodes (6): defaultRetryDelay(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces(), backoff(), TestBackoffGrowsAndCaps(), Message

### Community 133 - ".CreateGlobalVariable"
Cohesion: 0.20
Nodes (3): CreateGlobalVariableParams, Globalvar, UpdateGlobalVariableParams

### Community 135 - "Review Focus"
Cohesion: 0.20
Nodes (8): Global Constraints, Review Focus, Secret Key Rotation Implementation Plan, Task 1: Paged secret scan queries, Task 2: Regenerate sqlc, Task 4: `rotatekeys` command, Task 5: Operator documentation, Task 6: Graph, full verification and stacked PRs

### Community 137 - "time.Time"
Cohesion: 0.17
Nodes (10): bucket, limiter, key(), newLimiter(), RateLimit(), TestLimiterAllowsBurstThenRefills(), TestLimiterDropsIdleBuckets(), TestRateLimitDisabledByDefault() (+2 more)

### Community 140 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 141 - ".processJob"
Cohesion: 0.16
Nodes (5): isTerminalStatus(), isTransient(), parseAttempt(), stringPointer(), TestIsTerminalStatus()

### Community 142 - "rsql.Node"
Cohesion: 0.09
Nodes (22): Parse(), dump(), TestParse(), TestParseRejects(), exists(), likePattern(), relate(), ptr() (+14 more)

### Community 143 - "where.go"
Cohesion: 0.42
Nodes (3): likePattern(), builder, Compare

### Community 144 - "clilogin/handler_test.go"
Cohesion: 0.28
Nodes (10): CookieSigner, authorize(), doRequest(), newTestHandler(), TestBrokerCallbackErrorRedirectsToTheCLI(), TestBrokerFullFlow(), TestBrokerRejectsBadRequests(), TestConsentRejectsInvalidDays() (+2 more)

### Community 145 - "newTestAPI"
Cohesion: 0.08
Nodes (32): pageMeta, Locker, fakeLocker, JobLocker, createListOrganization(), getNames(), seedOrganizations(), TestOrganizationListEscapesWildcards() (+24 more)

### Community 146 - "Secret key rotation — design"
Cohesion: 0.25
Nodes (7): Approaches considered, Compatibility and docs, Goals, Intent, Non-goals, Secret key rotation — design, Tests

### Community 150 - "API"
Cohesion: 0.07
Nodes (19): API, authorization, New(), param(), parseList(), unknownCollection(), authorizationFrom(), isWriteMethod() (+11 more)

### Community 151 - "actionStub"
Cohesion: 0.15
Nodes (4): actionStub, CreateActionParams, UpdateActionParams, Action

### Community 152 - "addDestination"
Cohesion: 0.33
Nodes (15): strptr(), addDestination(), int64ptr(), newDispatcherForTest(), newFakeStore(), seedCompletedJob(), TestTickDispatchesWhenSerialIncreased(), TestTickIgnoresPlanOnlyRun() (+7 more)

### Community 154 - "run"
Cohesion: 0.05
Nodes (54): main(), run(), main(), run(), main(), run(), Deps, EventSource (+46 more)

### Community 157 - "Dispatcher"
Cohesion: 0.13
Nodes (14): strptr(), newDispatcherFixture(), TestTickFiresOncePerMatchingSecond(), TestTickMarksAJobFailedWhenItCannotBeQueued(), TestTickSkipsLockedWorkspaceAndUnparseableCron(), TestTickUsesTheTemplateOfTheWorkspaceOrganization(), tickAt(), fakeQueue (+6 more)

### Community 159 - "stubQuerier"
Cohesion: 0.03
Nodes (5): stubQuerier, scopeIDs(), CreateRunTriggerParams, UpdateRunTriggerParams, WorkspaceRunTrigger

### Community 161 - "Executor"
Cohesion: 0.21
Nodes (4): Executor, preparedWorkspace, isCancelledStatus(), TestBinarySelection()

### Community 162 - "Item"
Cohesion: 0.18
Nodes (3): CreateItemParams, UpdateItemParams, Item

### Community 163 - ".CreateProject"
Cohesion: 0.20
Nodes (3): CreateProjectParams, GetProjectByNameParams, UpdateProjectParams

### Community 164 - ".streamStepLog"
Cohesion: 0.23
Nodes (10): isTerminal(), outputOf(), stepOf(), streamKey(), TestIsTerminal(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData() (+2 more)

### Community 165 - "Dex"
Cohesion: 0.23
Nodes (3): Dex, Exchanger, fakeDex

### Community 167 - "NewAccessToken"
Cohesion: 0.23
Nodes (3): Task 3: `/access-token/v1` permission endpoints, NewAccessToken(), writeAccessTokenJSON()

### Community 168 - "notificationSendAPI"
Cohesion: 0.18
Nodes (17): captureJSONServer(), notificationSendAPI(), TestNotificationAdHocTestSend(), TestNotificationAdHocTestSendRejectsBadRequests(), TestNotificationAdHocTestSendUnknownOrganization(), TestNotificationTestSendAnswersBadGatewayWhenDeliveryFails(), TestNotificationTestSendSavedConfiguration(), TestNotificationTestSendUnknownConfiguration() (+9 more)

### Community 169 - "ResolveOrganization"
Cohesion: 0.15
Nodes (10): ApplyProjectGrants(), ApplyWorkspaceGrants(), ResolveOrganization(), ResolveProject(), ResolveWorkspace(), TestResolveWorkspaceIgnoresAnotherOrganization(), Access, fakeWorkspaceStore (+2 more)

### Community 171 - "Collection"
Cohesion: 0.20
Nodes (3): CreateCollectionParams, UpdateCollectionParams, Collection

### Community 172 - "Metrics"
Cohesion: 0.12
Nodes (5): metricsSource, Deps, fakeStats, Metrics, StatsSource

### Community 175 - "Deps"
Cohesion: 0.18
Nodes (7): Deps, PATService, channelLabel(), networkCause(), NewDeliverer(), parseRetryAfter(), Configuration

### Community 180 - "Threat model"
Cohesion: 0.17
Nodes (10): Accepted risks, Assets, Threat model, Threats and controls, Trust boundaries, Our own checks, Reporting a vulnerability, Scope (+2 more)

### Community 183 - "models_gen.go"
Cohesion: 0.06
Nodes (12): ComplexityRoot, NotificationConfigurationConnection, NotificationConfigurationEdge, NotificationTriggerConnection, NotificationTriggerEdge, PageInfo, ProjectConnection, Query (+4 more)

### Community 184 - "auth.DexValidator"
Cohesion: 0.31
Nodes (4): contains(), decodeJWKPart(), rsaPublicKey(), stringSlice()

### Community 190 - "loadPage"
Cohesion: 0.18
Nodes (6): DB, failingDB, loadPage(), collectStrings(), snakeCase(), TestModelColumnOrder()

### Community 191 - "M4 run logs and artifacts Implementation Plan"
Cohesion: 0.25
Nodes (7): Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 3: Record the state URL and fail on history errors, Task 4: End-to-end coverage, docs and graph

### Community 193 - "RSQL filters in SQL — design"
Cohesion: 0.20
Nodes (9): Approaches considered, Behavior changes, Design, Goals, Intent, Layering, Non-goals, RSQL filters in SQL — design (+1 more)

### Community 198 - "Review Focus"
Cohesion: 0.33
Nodes (5): Global Constraints, RBAC: permission resolution and `/access-token/v1` Implementation Plan, Review Focus, Task 2: Queries for teams, access and project access by name, Task 4: Record the contract

### Community 199 - "Workspace"
Cohesion: 0.14
Nodes (3): CreateWorkspaceParams, GetWorkspaceByNameParams, SetWorkspaceLockParams

### Community 229 - "Schedule"
Cohesion: 0.12
Nodes (13): scheduleStub, newScheduleAPI(), newScheduleStub(), scheduleBody(), scheduleStr(), TestCreateScheduleCopiesAttributesAndDefaultsEnabled(), TestCreateScheduleRejectsForeignTemplate(), TestDeleteScheduleAndCrossWorkspaceLookup() (+5 more)

### Community 231 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Review Focus, RSQL filters in SQL — implementation plan, Self-review, Task 1: RSQL parser, Task 2: Store list builder and resource specs, Task 4: Docs and graph

### Community 234 - "Migrate"
Cohesion: 0.18
Nodes (12): Global Constraints, Review Focus, Secret Encryption at Rest Implementation Plan, Self-Review Notes, Task 2: Cipher and envelope, Task 3: Transparent store decorator, Task 4: Backfill queries, Task 5: Backfill implementation (+4 more)

### Community 235 - ".build"
Cohesion: 0.18
Nodes (12): collectionSpec, listArgs, afterString(), column(), predicate(), Order, OrderBy(), ParseSort() (+4 more)

### Community 236 - ".Query"
Cohesion: 0.06
Nodes (16): CountJobsByStatusRow, CreateArchiveParams, ListGlobalVariableSecretsAfterParams, ListGlobalVariableSecretsAfterRow, ListGlobalVariableSecretsRow, ListNotificationSecretsAfterParams, ListNotificationSecretsAfterRow, ListNotificationSecretsRow (+8 more)

### Community 249 - "ParsePage"
Cohesion: 0.60
Nodes (4): ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage(), TestParsePageRejectsBadOffset()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 517 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **131 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `executionContext` connect `github.com/99designs/gqlgen/graphql.CollectedField` to `._Project`, `WorkspaceTagConnection`, `.___Directive`, `github.com/vektah/gqlparser/v2/ast.SelectionSet`, `context.Context`, `OrganizationConnection`, `generated.go`, `streamEvents`, `.marshalOBoolean2bool`, `github.com/99designs/gqlgen/graphql.Marshaler`, `JobConnection`, `AccessConnection`, `.childFields___InputValue`, `models_gen.go`, `.___Field`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.07109106529209622 - nodes in this community are weakly interconnected._
- **Why does `openRotationTest()` connect `Rotate` to `context.Context`, `newTestAPI`, `testing.T`, `Queries`, `api/main.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `rawVariable()` connect `Rotate` to `context.Context`, `fakeQuerier`, `testing.T`, `Queries`, `api/main.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `io.Reader` be split into smaller, more focused modules?**
  _Cohesion score 0.07372549019607844 - nodes in this community are weakly interconnected._