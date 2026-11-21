# Graph Report - terrakube-go  (2026-11-21)

## Corpus Check
- 175 files · ~129,688 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 19 file(s) not represented in the graph (top: .graphql 10, (none) 5, .example 1)

## Summary
- 2444 nodes · 9148 edges · 126 communities (72 shown, 54 thin omitted)
- Extraction: 80% EXTRACTED · 20% INFERRED · 0% AMBIGUOUS · INFERRED: 1802 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7f77267a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.ResponseWriter
- ._Job
- Agent Workflows and Compatibility Docs
- ParseRSQL
- io.Reader
- Core Database Schema
- stubQuerier
- Handler
- TFE Handler Tests
- WriteResource
- State Backend Tests
- reliabilityStore
- jsonapi.go
- steps
- StringOrInt
- DecodeRequest
- github.com/99designs/gqlgen/graphql.Marshaler
- .QueryRow
- sqlc Generation Config
- tarGzEntries
- Upstream Elide entities (GraphQL surface)
- run
- TemplateConnection
- cancelRunStore
- Local Dev Services Script
- Compose Build Configuration
- Compose Environment Settings
- Postgres Compose Service
- stubStore
- Executor
- Upstream UI GraphQL operations
- OpenTofu State Conformance
- Compose Project Definition
- Compose Healthchecks
- Compose Service Dependencies
- where_test.go
- Redis Compose Service
- M4 run logs and artifacts — design
- Management API Smoke Test
- TFE API Smoke Test
- State Backend Smoke Test
- Go Module Definition
- time.Duration
- testing.T
- api/main.go
- executor-e2e.sh
- Layout
- testLogger
- .build
- .processJob
- sort.go
- .processJob
- Querier
- models_gen.go
- github.com/vektah/gqlparser/v2/ast.SelectionSet
- dex-config.yaml
- executionContext
- .prepareWorkspace
- Deps
- models
- collections.go
- where.go
- AccessConnection
- generated.go
- config.go
- where_test.go
- ._Project
- serveNested
- graphql/models.go
- .streamStepLog
- WorkspaceTagConnection
- newTestExecutor
- cancelRunStore
- REST contract, captured from upstream
- pat.py
- CreateWorkspaceTagParams
- New
- nestedStub
- ui.go
- .listWorkspaces
- stubQuerier
- logSink
- environment
- .EnqueueJob
- streamEvents
- .GetTagByName
- newRedactor
- .GetWorkspaceByName
- ui-rest-inventory.sh
- .UpdateJobPlanChanges
- ui-smoke.sh
- capture.sh
- compare.sh
- start.sh

## God Nodes (most connected - your core abstractions)
1. `executionContext` - 520 edges
2. `Querier` - 86 edges
3. `API` - 66 edges
4. `Queries` - 64 edges
5. `stubQuerier` - 48 edges
6. `Executor` - 46 edges
7. `New()` - 45 edges
8. `WriteError()` - 41 edges
9. `New()` - 39 edges
10. `WriteResource()` - 33 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Record the state URL and fail on history errors` --references--> `TestArchiveHistory()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/executor/executor_test.go
- `M4 run logs and artifacts Implementation Plan` --references--> `Backend`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/storage/storage.go
- `Upstream Elide entities (GraphQL surface)` --references--> `Template`  [INFERRED]
  docs/upstream/graphql-entities.md → internal/store/models.go
- `Task 2: Serve archived logs from `/tfoutput/v1`` --references--> `New()`  [INFERRED]
  docs/superpowers/plans/2026-11-19-m4-run-artifacts.md → internal/tfoutput/handler.go
- `Generated Code and Migration Policy` --semantically_similar_to--> `sqlc Workflow`  [INFERRED] [semantically similar]
  CONTRIBUTING.md → .opencode/skills/sqlc-workflow/SKILL.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Compatibility Review Flow** — _opencode_agents_compat_reviewer_compat_reviewer, _opencode_agents_migration_reviewer_migration_reviewer, _opencode_commands_compat_compat_command, _opencode_skills_terrakube_compat_skill_terrakube_compat, docs_compat_compatibility_contract [EXTRACTED 1.00]
- **Local Verification and Smoke Test Flow** — _opencode_commands_verify_verify_command, _opencode_commands_smoke_smoke_command, _opencode_skills_dev_stack_skill_local_dev_stack, agents_verification_suite, readme_smoke_scripts [INFERRED 0.85]
- **sqlc Schema-to-Store Regeneration Flow** — _opencode_skills_sqlc_workflow_skill_sqlc_workflow, _opencode_commands_generate_generate_command, agents_generated_code_off_limits, agents_append_only_migrations, docs_compat_database_schema_parity [INFERRED 0.95]

## Communities (126 total, 54 thin omitted)

### Community 0 - "net/http.ResponseWriter"
Cohesion: 0.14
Nodes (12): WriteError(), WriteResource(), parseRunID(), decodeStateAttribute(), stringPointer(), decodeTFEJSON(), normalizeTagValue(), tagValuesEqual() (+4 more)

### Community 2 - "Agent Workflows and Compatibility Docs"
Cohesion: 0.06
Nodes (36): Compatibility Reviewer Subagent, Migration Reviewer Subagent, /compat Command, /generate Command, /pr Command, /smoke Command, /verify Command, Local Dev Stack (+28 more)

### Community 3 - "ParseRSQL"
Cohesion: 0.17
Nodes (12): ParseRSQL(), TestParseRSQLAcceptsUIExpressions(), TestParseRSQLBooleanStructure(), TestParseRSQLComparisons(), TestParseRSQLEmpty(), TestParseRSQLErrors(), And, Compare (+4 more)

### Community 4 - "io.Reader"
Cohesion: 0.07
Nodes (12): failingBackend, failingKeyBackend, memoryBackend, cleanKey(), pathClean(), TestCleanKey(), isS3NotFound(), Object (+4 more)

### Community 5 - "Core Database Schema"
Cohesion: 0.09
Nodes (37): history, history_workspace_idx, job, job_organization_idx, job_workspace_idx, organization, organization_name_uq, project (+29 more)

### Community 6 - "stubQuerier"
Cohesion: 0.06
Nodes (8): CreateContentParams, UpdateContentStatusParams, CreateProjectParams, CreateTagParams, UpdateWorkspaceTagValueParams, CreateWorkspaceParams, SetWorkspaceLockParams, stubQuerier

### Community 8 - "Handler"
Cohesion: 0.11
Nodes (19): Claims, Type, Classify(), ExtractToken(), TestClassify(), TestEffectiveStepIDFromJTI(), TestExtractToken(), claimsFromContext() (+11 more)

### Community 9 - "TFE Handler Tests"
Cohesion: 0.09
Nodes (38): TestContractDiscoveryAndPing(), Discovery(), capabilitySecret(), extractLogReadURL(), newTestAPI(), newTestAPIWithQueue(), newTestAPIWithStorage(), stringPtr() (+30 more)

### Community 10 - "WriteResource"
Cohesion: 0.16
Nodes (19): API, Deps, listQuery, boolValue(), filterSortPage(), New(), param(), parseInclude() (+11 more)

### Community 11 - "State Backend Tests"
Cohesion: 0.14
Nodes (23): decodeSecret(), NewValidator(), Sign(), TestParseRejectsExpiredToken(), TestParseRejectsWrongSecret(), testSecret(), TestSignAndParse(), TestValidatorAlgorithmByKeyLength() (+15 more)

### Community 12 - "reliabilityStore"
Cohesion: 0.09
Nodes (6): artifactStore, historyStub, reliabilityStore, CreateHistoryParams, GetHistoryByJobReferenceParams, UpdateStepStatusParams

### Community 13 - "jsonapi.go"
Cohesion: 0.09
Nodes (26): strValue(), workspaceResource(), Resource, Link(), LinkList(), TestRelationshipMarshalling(), Write(), configurationVersionResource() (+18 more)

### Community 14 - "steps"
Cohesion: 0.06
Nodes (41): POSTGRES_DB, POSTGRES_HOST_AUTH_METHOD, POSTGRES_USER, TERRAKUBE_TEST_DATABASE_URL, TERRAKUBE_TEST_REDIS_URL, env, runs-on, services (+33 more)

### Community 15 - "StringOrInt"
Cohesion: 0.13
Nodes (6): ComplexityRoot, ProjectConnection, QueryResolver, WorkspaceConnection, WorkspaceResolver, idsFrom()

### Community 16 - "DecodeRequest"
Cohesion: 0.13
Nodes (21): nestedDocument, normalizeCategory(), stringValue(), pickBool(), pickString(), pickStringValue(), DecodeAttributes(), DecodeRequest() (+13 more)

### Community 18 - ".QueryRow"
Cohesion: 0.07
Nodes (4): UpdateWorkspaceParams, GetProjectByNameParams, UpdateOrganizationParams, UpdateProjectParams

### Community 20 - "sqlc Generation Config"
Cohesion: 0.13
Nodes (14): go, emit_empty_slices, emit_interface, emit_json_tags, emit_pointers_for_null_types, out, package, sql_package (+6 more)

### Community 21 - "tarGzEntries"
Cohesion: 0.23
Nodes (12): archiveLimits, tarEntry, tarGzEntries(), TestExtractTarGzEnforcesEntryLimit(), TestExtractTarGzEnforcesFileSizeLimit(), TestExtractTarGzEnforcesTotalSizeLimit(), TestExtractTarGzRejectsHardlinks(), TestExtractTarGzRejectsSymlinks() (+4 more)

### Community 22 - "Upstream Elide entities (GraphQL surface)"
Cohesion: 0.14
Nodes (13): access (Access), Conventions, GenericAuditFields (shared superclass), job (Job), notification_configuration (NotificationConfiguration), notification_trigger (NotificationTrigger), organization, project (+5 more)

### Community 23 - "run"
Cohesion: 0.06
Nodes (43): main(), run(), main(), run(), main(), run(), main(), run() (+35 more)

### Community 25 - "cancelRunStore"
Cohesion: 0.14
Nodes (6): statusQuerier, UpdateJobStatusParams, Publisher, publishingStore, stubQuerier, mutableQuerier

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

### Community 31 - "Executor"
Cohesion: 0.18
Nodes (5): Executor, preparedWorkspace, isCancelledStatus(), TestBinarySelection(), baseEnvironment()

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

### Community 37 - "where_test.go"
Cohesion: 0.44
Nodes (10): mustParse(), TestWhereCombinesAndOr(), TestWhereEscapesLikeWildcards(), TestWhereNilNode(), TestWhereOperators(), TestWhereRejectsPredicateWithUnsupportedOperator(), TestWhereRejectsUnknownField(), TestWhereRendersCorrelatedPredicates() (+2 more)

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
Cohesion: 0.09
Nodes (7): fakeQueue, JobQueue, fakePinger, envDuration(), defaultRetryDelay(), Message, Queue

### Community 52 - "testing.T"
Cohesion: 0.05
Nodes (62): upstreamCapture, newTestAPI(), TestAuthMiddleware(), TestCreateOrganization(), TestCreateOrganizationRequiresName(), TestGetOrganizationNotFound(), TestListOrganizationsJSONAPI(), TestMalformedIDIsBadRequest() (+54 more)

### Community 53 - "api/main.go"
Cohesion: 0.06
Nodes (52): organizationAttrs, projectAttrs, replay, variableAttrs, workspaceAttrs, S3, gitConfig, stateMetadata (+44 more)

### Community 54 - "executor-e2e.sh"
Cohesion: 0.67
Nodes (4): fail(), req(), executor-e2e.sh script, wait_for()

### Community 55 - "Layout"
Cohesion: 0.13
Nodes (10): Storage, Global Constraints, M4 run logs and artifacts Implementation Plan, Review Focus, Self-review notes, Task 1: Archive the final step log, Task 2: Serve archived logs from `/tfoutput/v1`, Task 3: Record the state URL and fail on history errors (+2 more)

### Community 56 - "testLogger"
Cohesion: 0.08
Nodes (50): gitOutput(), initGitRepo(), requireGit(), TestAcquireConfigurationClonesGitSource(), TestAcquireConfigurationExtractsUploadedArchive(), TestAcquireConfigurationFailsWithoutSource(), TestAcquireConfigurationRejectsInvalidArchive(), TestCloneConfigurationChecksOutCommit() (+42 more)

### Community 57 - ".build"
Cohesion: 0.14
Nodes (9): DB, failingDB, listArgs, afterString(), loadPage(), ParsePage(), TestParsePage(), TestParsePageRejectsAnOversizedPage() (+1 more)

### Community 58 - ".processJob"
Cohesion: 0.14
Nodes (6): isTerminalStatus(), isTransient(), parseAttempt(), TestIsTerminalStatus(), TestCleanupStaleWorkspaces(), cleanupStaleWorkspaces()

### Community 59 - "sort.go"
Cohesion: 0.33
Nodes (6): Order, OrderBy(), ParseSort(), TestOrderByRendersSQL(), TestParseSort(), TestParseSortRejectsUnknownField()

### Community 61 - "Querier"
Cohesion: 0.07
Nodes (5): New(), Querier, DBTX, Queries, stubQuerier

### Community 62 - "models_gen.go"
Cohesion: 0.08
Nodes (9): NotificationConfigurationConnection, NotificationConfigurationEdge, NotificationTriggerConnection, NotificationTriggerEdge, OrganizationConnection, OrganizationEdge, PageInfo, Query (+1 more)

### Community 64 - "dex-config.yaml"
Cohesion: 0.15
Nodes (12): issuer, oauth2, skipApprovalScreen, staticClients, id, name, redirectURIs, secret (+4 more)

### Community 66 - ".prepareWorkspace"
Cohesion: 0.22
Nodes (7): TestSanitizeGitURL(), TestContentIDFromSource(), transient(), contentIDFromSource(), gitEnv(), runGit(), sanitizeGitURL()

### Community 67 - "Deps"
Cohesion: 0.29
Nodes (3): Deps, fakeLocker, JobLocker

### Community 68 - "models"
Cohesion: 0.05
Nodes (38): model, exec, filename, package, model, model, model, filename (+30 more)

### Community 71 - "collections.go"
Cohesion: 0.15
Nodes (14): collectionSpec, page, accessEdges(), column(), jobEdges(), notificationConfigurationEdges(), notificationTriggerEdges(), predicate() (+6 more)

### Community 74 - "generated.go"
Cohesion: 0.09
Nodes (3): DirectiveRoot, TestUpstreamUIDocumentsMatchTheSchema(), newExecutionContext()

### Community 75 - "config.go"
Cohesion: 0.16
Nodes (16): Config, env(), envBool(), envFirst(), envInt64(), envList(), Load(), ParseLogLevel() (+8 more)

### Community 76 - "where_test.go"
Cohesion: 0.11
Nodes (10): NotificationConfigurationResolver, OrganizationResolver, ProjectResolver, Resolver, ResolverRoot, SubscriptionResolver, streamEvents(), TestJobStatusSubscriptionClosesWithTheClient() (+2 more)

### Community 78 - "serveNested"
Cohesion: 0.24
Nodes (11): decodeJSON(), newNestedStub(), serveNested(), TestNestedVariableCreateUsesTheWorkspaceInThePath(), TestNestedVariableDeleteRemovesOnlyThatVariable(), TestNestedVariableFromAnotherWorkspaceIsNotFound(), TestNestedVariablesListIsScopedToTheWorkspace(), TestNestedWorkspaceGetIncludesVariables() (+3 more)

### Community 79 - "graphql/models.go"
Cohesion: 0.11
Nodes (13): Access, Job, JobStatusEvent, NotificationConfiguration, NotificationTrigger, Organization, Project, ProjectAccess (+5 more)

### Community 80 - ".streamStepLog"
Cohesion: 0.23
Nodes (10): isTerminal(), outputOf(), stepOf(), streamKey(), TestIsTerminal(), TestWriteDataSplitsMultilineOutput(), TestWriteEventFrame(), writeData() (+2 more)

### Community 82 - "newTestExecutor"
Cohesion: 0.27
Nodes (11): TestHandleWithoutALockerRunsTheJob(), TestReclaimStaleHandlesClaimedJobs(), TestAckDeletesTheJobMessage(), TestHandleDropsLiveLogsWhenTheJobEnds(), TestHandleKeepsLiveLogsOnShutdown(), jobMessage(), newTestExecutor(), TestHandleAcksTerminalJob() (+3 more)

### Community 83 - "cancelRunStore"
Cohesion: 0.15
Nodes (3): cancelRunStore, UpdateJobOutputParams, UpdateStepOutputParams

### Community 84 - "REST contract, captured from upstream"
Cohesion: 0.22
Nodes (7): How these were captured, REST contract, captured from upstream, Shapes worth copying when implementing an M6 resource, What a live upstream actually does, Adding a case, Capturing the upstream contract, Why a hand-minted PAT

### Community 86 - "pat.py"
Cohesion: 0.23
Nodes (3): b64url(), main(), pat_secret()

### Community 88 - "New"
Cohesion: 0.16
Nodes (13): Validator, TestWriteBackendOverride(), New(), newHandler(), stringPtr(), TestGetStepLogFallsBackToInlineOutput(), TestGetStepLogPrefersArchivedObject(), TestGetStepLogRejectsWrongScope() (+5 more)

### Community 90 - "ui.go"
Cohesion: 0.29
Nodes (8): envConfig(), exists(), New(), TestNoDirectoryMeansNoHandler(), Get(), Deps, handler, Info

### Community 91 - ".listWorkspaces"
Cohesion: 0.29
Nodes (6): parseWorkspaceQuery(), matchesTagFilters(), workspaceTagNames(), writeWorkspaceList(), tagFilter, workspaceQuery

### Community 92 - "stubQuerier"
Cohesion: 0.24
Nodes (3): stubQuerier, demoOrganizations(), CreateOrganizationParams

### Community 93 - "logSink"
Cohesion: 0.20
Nodes (3): captureLogPublisher, LogPublisher, logSink

### Community 98 - "environment"
Cohesion: 0.17
Nodes (12): REACT_APP_AUTHORITY, REACT_APP_CLIENT_ID, REACT_APP_REDIRECT_URI, REACT_APP_REGISTRY_URI, REACT_APP_SCOPE, REACT_APP_TERRAKUBE_API_URL, REACT_APP_TERRAKUBE_VERSION, ui (+4 more)

### Community 100 - "streamEvents"
Cohesion: 0.29
Nodes (4): EventSource, fakeEvents, Event, Bus

### Community 102 - "newRedactor"
Cohesion: 0.29
Nodes (6): redactor, TestLogSinkRedactsSecrets(), TestNilRedactorIsSafe(), TestRedactorPrefersLongestValue(), TestRedactorReplacesSensitiveValues(), newRedactor()

### Community 119 - "ui-smoke.sh"
Cohesion: 0.60
Nodes (5): body(), fail(), header(), ui-smoke.sh script, status()

### Community 124 - "capture.sh"
Cohesion: 0.83
Nodes (3): request(), capture.sh script, unauthenticated()

## Knowledge Gaps
- **11 isolated node(s):** `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION`, `TF_LOG`, `TF_LOG_PATH` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 313 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **54 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `executionContext` connect `executionContext` to `._Job`, `.marshalOBoolean2bool`, `.___EnumValue`, `.___Field`, `AccessConnection`, `generated.go`, `where_test.go`, `._Project`, `StringOrInt`, `github.com/99designs/gqlgen/graphql.Marshaler`, `WorkspaceTagConnection`, `.unmarshalOString2ᚕᚖstring`, `models_gen.go`, `TemplateConnection`, `context.Context`, `github.com/vektah/gqlparser/v2/ast.SelectionSet`?**
  _High betweenness centrality (0.045) - this node is a cross-community bridge._
- **What connects `TF_HTTP_PASSWORD`, `TF_HTTP_USERNAME`, `TF_IN_AUTOMATION` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.1395961369622476 - nodes in this community are weakly interconnected._
- **Why does `Querier` connect `Querier` to `.GetStepByJobAndNumber`, `net/http.ResponseWriter`, `stubQuerier`, `Handler`, `TFE Handler Tests`, `WriteResource`, `reliabilityStore`, `.QueryRow`, `run`, `cancelRunStore`, `stubStore`, `testing.T`, `api/main.go`, `Deps`, `newTestExecutor`, `cancelRunStore`, `CreateWorkspaceTagParams`, `New`, `nestedStub`, `stubQuerier`, `.GetTagByName`, `.SetJobOverrideSource`, `.GetWorkspaceByName`, `.UpdateJobPlanChanges`, `.ApproveJob`, `.CreateJob`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `Agent Workflows and Compatibility Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.05989110707803993 - nodes in this community are weakly interconnected._
- **Why does `newLogSink()` connect `testLogger` to `newRedactor`, `api/main.go`, `logSink`, `context.Context`, `Executor`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `io.Reader` be split into smaller, more focused modules?**
  _Cohesion score 0.06966618287373004 - nodes in this community are weakly interconnected._