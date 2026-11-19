# Upstream Elide entities (GraphQL surface)

Source: `/home/amin/Documents/terrakube` at commit `c460c846` (2026-11-19), Elide
7.2.0; GraphQL at `elide.graphql.path=/graphql/api/v1`
(`api/src/main/resources/application.properties:29`). Read from the annotations;
nothing inferred from the UI.

Entity names come from `@Entity(name = "...")`: a snake_case name
(`notification_configuration`) is both the GraphQL type name and the root field
name, while relationship field names are the Java field names (camelCase).
`@Include(rootLevel = false)` entities are **not** queryable at the root.

**`notificationTemplate` does not exist in the upstream API.** No
`NotificationTemplate` class, table, migration or GraphQL usage exists at this
commit (grep over `api/src` and `ui/src` is empty). Notification "templates" are
the generic `Template` entity, `@Entity(name = "template")`
(`api/src/main/java/io/terrakube/api/rs/template/Template.java:22`), reached via
`NotificationConfiguration.templates` (field name `templates`).

## GenericAuditFields (shared superclass)

`api/src/main/java/io/terrakube/api/plugin/security/audit/GenericAuditFields.java:18`
— `@MappedSuperclass`, `@EntityListeners(AuditingEntityListener.class)`. Its
four fields are inherited by Project, Workspace, NotificationConfiguration,
NotificationTrigger, WorkspaceTag and Job (not by Organization, ProjectAccess or
Access); they are listed last in each table, marked *(inherited)*, in declaration
order: `createdDate` (Date, `@CreationTimestamp`, `@Temporal(TIMESTAMP)`,
`@Column(name = "created_date", updatable = false)`), `createdBy` (String,
`@CreatedBy`, `@Column(name = "created_by")`), `updatedDate` (Date,
`@UpdateTimestamp`, `@Column(name = "updated_date")`), `updatedBy` (String,
`@LastModifiedBy`, `@Column(name = "updated_by")`).

`EntityPermissionCoverageTest`
(`api/src/test/java/io/terrakube/api/rs/checks/EntityPermissionCoverageTest.java:30`)
scans every `@Include` class under `io.terrakube.api` and fails if any of the
four class-level permission annotations is missing — the reason every permission
block below is complete.

## organization

`api/src/main/java/io/terrakube/api/rs/Organization.java:41` — `@Entity(name = "organization")`

| annotation | expression (verbatim) |
| --- | --- |
| `@ReadPermission` | `user belongs organization OR executor read organization OR user is an executor service OR user is a registry service` |
| `@CreatePermission` | `user is a superuser` |
| `@UpdatePermission` | `user is a superuser` |
| `@DeletePermission` | `user is a superuser` |

| field | Java type | nullable/required | notes |
| --- | --- | --- | --- |
| id | UUID | required (PK) | `@Id`, `@JdbcTypeCode(Types.VARCHAR)`, `@GeneratedValue(strategy = GenerationType.UUID)` |
| name | String | nullable | `@Column(name = "name")` |
| disabled | boolean | primitive | `@Column(name = "disabled")`; class `@SQLRestriction` hides disabled rows |
| description | String | nullable | `@Column(name = "description")` |
| collection | List\<Collection\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| workspace | List\<Workspace\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| project | List\<Project\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| module | List\<Module\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| provider | List\<Provider\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| job | List\<Job\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| team | List\<Team\> to-many | nullable | `@OneToMany(mappedBy = "organization")` (no field permission) |
| vcs | List\<Vcs\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| ssh | List\<Ssh\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| agent | List\<Agent\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| template | List\<Template\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| globalvar | List\<Globalvar\> to-many | nullable | `@OneToMany(mappedBy = "organization")` (no field permission) |
| tag | List\<Tag\> to-many | nullable | `@OneToMany(mappedBy = "organization")` (no field permission) |
| executionMode | ExecutionMode | nullable | `@Column(name = "execution_mode")`, `@Enumerated(EnumType.STRING)`, default `ExecutionMode.remote` |
| icon | String | nullable | `@Column(name = "icon")` |
| notificationConfiguration | List\<NotificationConfiguration\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| policySet | List\<PolicySet\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |
| policyExemption | List\<PolicyExemption\> to-many | nullable | `@OneToMany(mappedBy = "organization")`, `@UpdatePermission("user belongs organization")` |

Extras: no `@Paginate`; `@SQLRestriction(value = "disabled = false")`; `@LifeCycleHookBinding` present (three bindings, `OrganizationManageHook`, CREATE PRECOMMIT/POSTCOMMIT and UPDATE PRECOMMIT). No audit fields (does not extend `GenericAuditFields`).
Naming: root type, root field `organization`. No parent.

## workspace

`api/src/main/java/io/terrakube/api/rs/workspace/Workspace.java:52` — `@Entity(name = "workspace")`

| annotation | expression (verbatim) |
| --- | --- |
| `@ReadPermission` | `user is a superuser OR workspace read filter OR user is an executor service OR user is a registry service` |
| `@CreatePermission` | `team manage workspace OR team project limited create workspace` |
| `@UpdatePermission` | `team manage workspace OR team project limited manage workspace OR team limited manage workspace` |
| `@DeletePermission` | `team manage workspace` |

| field | Java type | nullable/required | notes |
| --- | --- | --- | --- |
| id | UUID | required (PK) | `@Id`, `@JdbcTypeCode(Types.VARCHAR)`, `@GeneratedValue(UUID)` |
| name | String | nullable | `@Size(max = 128, message = "Workspace name must be 128 characters or fewer")`, `@Column(name = "name")`, `@LifeCycleHookBinding` (UPDATE PRECOMMIT, `WorkspaceManageHook`) |
| description | String | nullable | `@Column(name = "description")` |
| source | String | nullable | `@Column(name = "source")` |
| branch | String | nullable | `@Column(name = "branch")` |
| folder | String | nullable | `@Column(name = "folder")` |
| lastJobStatus | JobStatus | nullable | `@Enumerated(EnumType.STRING)`, `@Column(name = "last_job_status")`, default `NeverExecuted` |
| lastJobDate | Date | nullable | `@Column(name = "last_job_date")` |
| locked | boolean | primitive | `@Column(name = "locked")`, two `@LifeCycleHookBinding` (UPDATE, `WorkspaceUnlockHook`) |
| deleted | boolean | primitive | `@Column(name = "deleted")`, `@LifeCycleHookBinding` (UPDATE PRECOMMIT, `WorkspaceManageHook`) |
| allowRemoteApply | boolean | primitive | `@Column(name = "allow_remote_apply")`, default `false` |
| globalRemoteState | boolean | primitive | `@Column(name = "global_remote_state")`, default `true` |
| sharedIds | String | nullable | `@Column(name = "shared_ids")` |
| defaultTemplate | String | nullable | `@Column(name = "default_template")` |
| restrictApplyToConfiguredBranch | boolean | required | `@Column(name = "restrict_apply_to_configured_branch", nullable = false)`, default `false`, `@UpdatePermission("user is a superuser OR team manage workspace")` |
| requireApproveToApply | boolean | required | `@Column(name = "require_approve_to_apply", nullable = false)`, default `false`, `@UpdatePermission("user is a superuser OR team manage workspace")` |
| lockDescription | String | nullable | `@Column(name = "lock_description")` |
| iacType | String | nullable | `@Column(name = "iac_type")`, default `"terraform"` |
| moduleSshKey | String | nullable | `@Column(name = "module_ssh_key")` |
| terraformVersion | String | nullable | `@Column(name = "terraform_version")` |
| executionMode | ExecutionMode | nullable | `@Column(name = "execution_mode")`, `@Enumerated(EnumType.STRING)` |
| policyComplianceStatus | PolicyComplianceStatus | nullable | `@Column(name = "policy_compliance_status")`, `@Enumerated(EnumType.STRING)`, default `UNKNOWN` |
| organization | Organization to-one | nullable | `@ManyToOne` |
| variable | List\<Variable\> to-many | nullable | `@OneToMany(mappedBy = "workspace")` |
| history | List\<History\> to-many | nullable | `@OneToMany(mappedBy = "workspace")`, `@UpdatePermission("user is a super service OR executor manage workspace history OR user is an executor service")` |
| schedule | List\<Schedule\> to-many | nullable | `@OneToMany(mappedBy = "workspace")` |
| job | List\<Job\> to-many | nullable | `@OneToMany(mappedBy = "workspace")`, `@UpdatePermission("team view workspace OR team project limited view workspace OR team limited view workspace")` |
| content | List\<Content\> to-many | nullable | `@Exclude` (not exposed), `@OneToMany(mappedBy = "workspace")` |
| workspaceTag | List\<WorkspaceTag\> to-many | nullable | `@OneToMany(mappedBy = "workspace")` |
| project | Project to-one | nullable | `@OneToOne`, `@UpdatePermission("team manage workspace OR team project limited reassign workspace")`, `@LifeCycleHookBinding` (UPDATE PRECOMMIT, `WorkspaceManageHook`) |
| vcs | Vcs to-one | nullable | `@ManyToOne`, `@LifeCycleHookBinding` (UPDATE PRECOMMIT, `WorkspaceManageHook`) |
| ssh | Ssh to-one | nullable | `@OneToOne`, `@LifeCycleHookBinding` (UPDATE PRECOMMIT, `WorkspaceManageHook`) |
| agent | Agent to-one | nullable | `@OneToOne`, `@LifeCycleHookBinding` (UPDATE PRECOMMIT, `WorkspaceManageHook`) |
| webhook | Webhook to-one | nullable | `@OneToOne(mappedBy = "workspace", fetch = LAZY, cascade = {PERSIST, MERGE})` |
| notificationConfiguration | List\<NotificationConfiguration\> to-many | nullable | `@OneToMany(mappedBy = "workspace", fetch = LAZY)` |
| reference | List\<Reference\> to-many | nullable | `@OneToMany(mappedBy = "workspace", fetch = LAZY)`, `@UpdatePermission("team view workspace OR team project limited view workspace OR team limited view workspace")` |
| runTriggers | List\<WorkspaceRunTrigger\> to-many | nullable | `@OneToMany(mappedBy = "destinationWorkspace", fetch = LAZY)`, same `@UpdatePermission` |
| sourceRunTriggers | List\<WorkspaceRunTrigger\> to-many | nullable | `@OneToMany(mappedBy = "sourceWorkspace", fetch = LAZY)`, same `@UpdatePermission` |
| access | List\<Access\> to-many | nullable | `@OneToMany(mappedBy = "workspace")`, `@UpdatePermission("user is a superuser OR team manage workspace OR team workspace admin manages access field")` |
| createdDate | Date | nullable | *(inherited)* `@CreationTimestamp` |
| createdBy | String | nullable | *(inherited)* `@CreatedBy` |
| updatedDate | Date | nullable | *(inherited)* `@UpdateTimestamp` |
| updatedBy | String | nullable | *(inherited)* `@LastModifiedBy` |

Extras: `@Paginate(defaultPageSize = 10000, maxPageSize = 10000)`; `@SQLRestriction(value = "deleted = false")`; `@LifeCycleHookBinding` present (field-level, `WorkspaceManageHook` / `WorkspaceUnlockHook`).
Naming: field `workspace` on `organization`; also to-one `workspace` on `job` and on `notification_configuration`. Root field `workspace`.

## project

`api/src/main/java/io/terrakube/api/rs/project/Project.java:28` — `@Entity(name = "project")`

| annotation | expression (verbatim) |
| --- | --- |
| `@ReadPermission` | `team view project OR team project access view project` |
| `@CreatePermission` | `team manage project` |
| `@UpdatePermission` | `team manage project OR team project admin manages project access field` |
| `@DeletePermission` | `team manage project` |

| field | Java type | nullable/required | notes |
| --- | --- | --- | --- |
| id | UUID | required (PK) | `@Id`, `@JdbcTypeCode(Types.VARCHAR)`, `@GeneratedValue(UUID)` |
| name | String | nullable | `@Column(name = "name")` |
| description | String | nullable | `@Column(name = "description")` |
| organization | Organization to-one | nullable | `@ManyToOne` |
| projectAccess | List\<ProjectAccess\> to-many | nullable | `@OneToMany(mappedBy = "project", cascade = ALL, orphanRemoval = true)`, `@BatchSize(size = 100)`, `@UpdatePermission("team manage project OR team project admin manages project access field")` |
| createdDate | Date | nullable | *(inherited)* `@CreationTimestamp` |
| createdBy | String | nullable | *(inherited)* `@CreatedBy` |
| updatedDate | Date | nullable | *(inherited)* `@UpdateTimestamp` |
| updatedBy | String | nullable | *(inherited)* `@LastModifiedBy` |

Extras: `@Paginate(defaultPageSize = 10000, maxPageSize = 10000)`; no `@SQLRestriction`; no `@LifeCycleHookBinding`.
Naming: field `project` on `organization`, and to-one `project` on `workspace`. Root field `project`.

## project_access (ProjectAccess)

`api/src/main/java/io/terrakube/api/rs/project/access/ProjectAccess.java:25` — `@Entity(name = "project_access")`

| annotation | expression (verbatim) |
| --- | --- |
| `@ReadPermission` | `team manage project access OR team project limited manage project access` |
| `@CreatePermission` | `user is a superuser OR team manage project access OR team project limited manage project access` |
| `@UpdatePermission` | `user is a superuser OR team manage project access OR team project limited manage project access` |
| `@DeletePermission` | `user is a superuser OR team manage project access OR team project limited manage project access` |

| field | Java type | nullable/required | notes |
| --- | --- | --- | --- |
| id | UUID | required (PK) | `@Id`, `@JdbcTypeCode(Types.VARCHAR)`, `@GeneratedValue(UUID)` |
| name | String | nullable | `@Column(name = "name")` |
| manageState | boolean | primitive | `@Column(name = "manage_state")` |
| manageJob | boolean | primitive | `@Column(name = "manage_job")` |
| manageWorkspace | boolean | primitive | `@Column(name = "manage_workspace")` |
| planJob | boolean | primitive | `@Column(name = "plan_job")` |
| approveJob | boolean | primitive | `@Column(name = "approve_job")` |
| destroyJob | boolean | primitive | `@Column(name = "destroy_job")`, default `true` |
| role | String | nullable | `@Column(name = "role")` |
| project | Project to-one | nullable | `@ManyToOne` |

Extras: no `@Paginate`; no `@SQLRestriction`; no `@LifeCycleHookBinding`; `@Include(rootLevel = false)` — not root-queryable. No audit fields.
Naming: field `projectAccess` on `project`.

## access (Access)

`api/src/main/java/io/terrakube/api/rs/workspace/access/Access.java:25` — `@Entity(name = "access")`

| annotation | expression (verbatim) |
| --- | --- |
| `@ReadPermission` | `user is a superuser OR team manage workspace access OR team limited manage workspace access OR user is a registry service` |
| `@CreatePermission` | `user is a superuser OR team manage workspace access OR team limited manage workspace access` |
| `@UpdatePermission` | `user is a superuser OR team manage workspace access OR team limited manage workspace access` |
| `@DeletePermission` | `user is a superuser OR team manage workspace access OR team limited manage workspace access` |

| field | Java type | nullable/required | notes |
| --- | --- | --- | --- |
| id | UUID | required (PK) | `@Id`, `@JdbcTypeCode(Types.VARCHAR)`, `@GeneratedValue(UUID)` |
| name | String | nullable | `@Column(name = "name")` |
| manageState | boolean | primitive | `@Column(name = "manage_state")` |
| manageJob | boolean | primitive | `@Column(name = "manage_job")` |
| manageWorkspace | boolean | primitive | `@Column(name = "manage_workspace")` |
| planJob | boolean | primitive | `@Column(name = "plan_job")` |
| approveJob | boolean | primitive | `@Column(name = "approve_job")` |
| destroyJob | boolean | primitive | `@Column(name = "destroy_job")`, default `true` |
| role | String | nullable | `@Column(name = "role")` |
| workspace | Workspace to-one | nullable | `@ManyToOne` |

Extras: no `@Paginate`; no `@SQLRestriction`; no `@LifeCycleHookBinding`; `@Include(rootLevel = false)` — not root-queryable. No audit fields.
Naming: field `access` on `workspace`.

## notification_configuration (NotificationConfiguration)

`api/src/main/java/io/terrakube/api/rs/notification/NotificationConfiguration.java:47` — `@Entity(name = "notification_configuration")`

| annotation | expression (verbatim) |
| --- | --- |
| `@ReadPermission` | `team read notification configuration` |
| `@CreatePermission` | `team manage notification configuration` |
| `@UpdatePermission` | `team manage notification configuration` |
| `@DeletePermission` | `team manage notification configuration` |

| field | Java type | nullable/required | notes |
| --- | --- | --- | --- |
| id | UUID | required (PK) | `@Id`, `@JdbcTypeCode(Types.VARCHAR)`, `@GeneratedValue(UUID)` |
| name | String | nullable | plain field, no `@Column` name override |
| description | String | nullable | plain `@Column` (text); comment explains `@Lob` is avoided |
| organization | Organization to-one | required | `@ManyToOne(fetch = LAZY, optional = false)` |
| workspace | Workspace to-one | nullable | `@ManyToOne(fetch = EAGER)` (EAGER required so GraphQL serializes a real id, not `"null"`) |
| channelType | NotificationChannelType | nullable | `@Column(name = "channel_type")`, `@Enumerated(EnumType.STRING)` |
| destinationUrl | String | nullable | `@Column(name = "destination_url")` |
| signingSecret | String | nullable | `@Column(name = "signing_secret")`, `@ReadPermission("read notification signing secret")` |
| active | boolean | primitive | default `true` |
| messageStyle | NotificationMessageStyle | nullable | `@Column(name = "message_style")`, `@Enumerated(EnumType.STRING)`, default `DETAILED` |
| triggers | List\<NotificationTrigger\> to-many | nullable | `@OneToMany(mappedBy = "configuration", fetch = LAZY, cascade = REMOVE)` |
| templates | List\<Template\> to-many | nullable | `@ManyToMany(fetch = LAZY)`, join table `notification_configuration_template` |
| policySet | List\<PolicySet\> to-many | nullable | `@OneToMany(mappedBy = "notificationConfiguration", fetch = LAZY)` |
| createdDate | Date | nullable | *(inherited)* `@CreationTimestamp` |
| createdBy | String | nullable | *(inherited)* `@CreatedBy` |
| updatedDate | Date | nullable | *(inherited)* `@UpdateTimestamp` |
| updatedBy | String | nullable | *(inherited)* `@LastModifiedBy` |

Extras: no `@Paginate`; no `@SQLRestriction`; no `@LifeCycleHookBinding` (it has `@PrePersist`/`@PreUpdate` callbacks instead). Root-queryable.
Naming: field `notificationConfiguration` on `organization` and on `workspace`. Root field `notification_configuration`. The `templates` relationship targets `template`; there is no `notificationTemplate` type.

## notification_trigger (NotificationTrigger)

`api/src/main/java/io/terrakube/api/rs/notification/NotificationTrigger.java:37` — `@Entity(name = "notification_trigger")`

| annotation | expression (verbatim) |
| --- | --- |
| `@ReadPermission` | `team read notification trigger` |
| `@CreatePermission` | `team manage notification trigger` |
| `@UpdatePermission` | `team manage notification trigger` |
| `@DeletePermission` | `team manage notification trigger` |

| field | Java type | nullable/required | notes |
| --- | --- | --- | --- |
| id | UUID | required (PK) | `@Id`, `@JdbcTypeCode(Types.VARCHAR)`, `@GeneratedValue(UUID)` |
| jobStatus | JobStatus | nullable | `@Column(name = "job_status")`, `@Enumerated(EnumType.STRING)` |
| configuration | NotificationConfiguration to-one | required | `@ManyToOne(fetch = LAZY, optional = false)` |
| createdDate | Date | nullable | *(inherited)* `@CreationTimestamp` |
| createdBy | String | nullable | *(inherited)* `@CreatedBy` |
| updatedDate | Date | nullable | *(inherited)* `@UpdateTimestamp` |
| updatedBy | String | nullable | *(inherited)* `@LastModifiedBy` |

Extras: no `@Paginate`; no `@SQLRestriction`; no `@LifeCycleHookBinding`. Root-queryable.
Naming: the parent field is `triggers` on `notification_configuration` — **not** `notificationTrigger`. Root field `notification_trigger`.

## workspacetag (WorkspaceTag)

`api/src/main/java/io/terrakube/api/rs/workspace/tag/WorkspaceTag.java:24` — `@Entity(name = "workspacetag")`

| annotation | expression (verbatim) |
| --- | --- |
| `@ReadPermission` | `user is a superuser OR team view workspace tag` |
| `@CreatePermission` | `user is a superuser OR team manage workspace tag` |
| `@UpdatePermission` | `user is a superuser OR team manage workspace tag` |
| `@DeletePermission` | `user is a superuser OR team manage workspace tag` |

| field | Java type | nullable/required | notes |
| --- | --- | --- | --- |
| id | UUID | required (PK) | `@Id`, `@JdbcTypeCode(Types.VARCHAR)`, `@GeneratedValue(UUID)` |
| tagId | String | nullable | `@Column(name = "tag_id")` |
| value | String | nullable | `@Column(name = "tag_value")` |
| workspace | Workspace to-one | nullable | `@ManyToOne`, `@JoinColumn(name = "workspace_id")` |
| createdDate | Date | nullable | *(inherited)* `@CreationTimestamp` |
| createdBy | String | nullable | *(inherited)* `@CreatedBy` |
| updatedDate | Date | nullable | *(inherited)* `@UpdateTimestamp` |
| updatedBy | String | nullable | *(inherited)* `@LastModifiedBy` |

Extras: `@Paginate(defaultPageSize = 10000, maxPageSize = 10000)`; no `@SQLRestriction`; no `@LifeCycleHookBinding`; `@Include(rootLevel = false)` — not root-queryable.
Naming: field `workspaceTag` on `workspace` (type name is `workspacetag`, no underscore).

## job (Job)

`api/src/main/java/io/terrakube/api/rs/job/Job.java:61` — `@Entity(name = "job")`

| annotation | expression (verbatim) |
| --- | --- |
| `@ReadPermission` | `team view job OR team project limited view job OR team limited view job OR executor read job OR user is an executor service OR user is a registry service` |
| `@CreatePermission` | `team plan job OR team limited plan job OR team project limited plan job` |
| `@UpdatePermission` | `Job.UPDATE` |
| `@DeletePermission` | `user is a super service OR user is a superuser` |

`Job.UPDATE` is a compile-time constant in the class
(`Job.java:65-66`), value: `team manage job OR team limited manage job OR team
project limited manage job OR user is a super service OR executor manage job OR
user is an executor service`.

| field | Java type | nullable/required | notes |
| --- | --- | --- | --- |
| id | int | required (PK) | `@Id`, `@GeneratedValue(strategy = GenerationType.IDENTITY)` — **int, not UUID** |
| variables | Map\<String, String\> | n/a | `@Transient`, `@ComputedAttribute`, `@ReadPermission("NONE")`, `@UpdatePermission("NONE")`; create-only input |
| encryptedVariables | String | nullable | `@Column(name = "encrypted_variables", columnDefinition = "text")`, Read/Create/Update `@...Permission("NONE")` |
| variableNames | List\<String\> | nullable | `@Convert(StringListConverter)`, `@Column(name = "variable_names", columnDefinition = "text")`, Create/Update `"NONE"` |
| variablesAsArguments | boolean | required | `@Column(name = "variables_as_arguments", nullable = false)`, default `false`, Create/Update `"NONE"` |
| comments | String | nullable | `@Column(name = "comments")` |
| status | JobStatus | nullable | `@Enumerated(EnumType.STRING)`, default `pending`; `@UpdatePermission("team approve job OR team approve job rbac OR team limited approve job OR team project limited approve job OR user is a super service OR executor manage job OR user is an executor service")`; four `@LifeCycleHookBinding` |
| approvedBy | String | nullable | `@Column(name = "approved_by")`, Create/Update `"NONE"` |
| approvedAt | Date | nullable | `@Column(name = "approved_at")`, Create/Update `"NONE"` |
| output | String | nullable | `@Column(name = "output")`, `@ReadPermission("team view job OR team project limited view job OR team limited view job OR executor read job OR user is an executor service")` |
| commitId | String | nullable | `@Column(name = "commit_id")`, `@UpdatePermission(BRANCH_FIELDS_UPDATE)` |
| autoApply | boolean | primitive | `@Exclude`, `@Column(name = "auto_apply")`, default `false` |
| deleted | boolean | primitive | `@Exclude`, `@Column(name = "deleted")`, default `false`; class `@SQLRestriction` filters it |
| terraformPlan | String | nullable | `@Column(name = "terraform_plan")`, same read permission as `output` |
| terraformPlanArtifacts | String | nullable | `@Column(name = "terraform_plan_artifacts")`, `@CreatePermission("user is a super service")`, `@UpdatePermission("user is a super service OR executor manage job OR user is an executor service")` |
| terraformPlanArtifactsChecksum | String | nullable | `@Column(name = "terraform_plan_artifacts_checksum")`, same create/update as above |
| approvalTeam | String | nullable | `@Column(name = "approval_team")`, same create/update as above |
| tcl | String | nullable | `@Exclude`, `@Column(name = "tcl")` |
| overrideSource | String | nullable | `@Column(name = "override_source")`, Create/Update `"user is a super service"` |
| overrideBranch | String | nullable | `@Column(name = "override_branch")`, `@UpdatePermission(BRANCH_FIELDS_UPDATE)` |
| templateReference | String | nullable | `@Column(name = "template_reference")`, `@UpdatePermission("NONE")` |
| via | String | nullable | `@Column(name = "via")`, default `JobVia.UI.getValue()` |
| refresh | boolean | primitive | `@Column(name = "refresh")`, default `true` |
| planChanges | boolean | primitive | `@Column(name = "plan_changes")`, default `true` |
| refreshOnly | boolean | primitive | `@Column(name = "refresh_only")`, default `false` |
| resolvedTerraformVersion | String | nullable | `@Column(name = "resolved_terraform_version")`, Create/Update `"NONE"` |
| prNumber | Integer | nullable | `@Column(name = "pr_number")`, Create/Update `"user is a super service"` |
| prCommentId | String | nullable | `@Exclude`, `@Column(name = "pr_comment_id")` |
| commandCommentId | String | nullable | `@Exclude`, `@Column(name = "command_comment_id")` |
| prApplyEnabled | boolean | primitive | `@Exclude`, `@Column(name = "pr_apply_enabled")`, default `false` |
| prCommentError | String | nullable | `@Column(name = "pr_comment_error")`, Create/Update `"user is a super service"` |
| triggeredByJobId | Integer | nullable | `@Column(name = "triggered_by_job_id")`, Create/Update `"user is a super service"` |
| cascadeDepth | int | required | `@Column(name = "cascade_depth", nullable = false)`, default `0`, Create/Update `"user is a super service"` |
| dispatchFailureCount | int | required | `@Column(name = "dispatch_failure_count", nullable = false)`, default `0`, Create/Update `"user is a super service"` |
| dispatchFirstFailureAt | Date | nullable | `@Column(name = "dispatch_first_failure_at")`, Create/Update `"user is a super service"` |
| dispatchLastAttemptAt | Date | nullable | `@Column(name = "dispatch_last_attempt_at")`, Create/Update `"user is a super service"` |
| dispatchNextRetryAt | Date | nullable | `@Column(name = "dispatch_next_retry_at")`, Create/Update `"user is a super service"` |
| dispatchLastError | String | nullable | `@Column(name = "dispatch_last_error", length = 1024)`, Create/Update `"user is a super service"` |
| organization | Organization to-one | nullable | `@ManyToOne`, `@UpdatePermission("user is a superuser")` |
| workspace | Workspace to-one | nullable | `@ManyToOne`, `@UpdatePermission("user is a superuser")` |
| step | List\<Step\> to-many | nullable | `@OneToMany(mappedBy = "job", cascade = REMOVE, orphanRemoval = true)`, `@UpdatePermission("user is a super service OR executor manage job OR user is an executor service")` |
| address | List\<Address\> to-many | nullable | `@OneToMany(mappedBy = "job", cascade = REMOVE, orphanRemoval = true)` |
| targetAddrs | List\<String\> | nullable | `@Convert(StringListConverter)`, `@Column(name = "target_addrs")` |
| replaceAddrs | List\<String\> | nullable | `@Convert(StringListConverter)`, `@Column(name = "replace_addrs")` |
| policyEvaluation | List\<PolicyEvaluation\> to-many | nullable | `@OneToMany(mappedBy = "job", cascade = REMOVE, orphanRemoval = true)` |
| createdDate | Date | nullable | *(inherited)* `@CreationTimestamp` |
| createdBy | String | nullable | *(inherited)* `@CreatedBy` |
| updatedDate | Date | nullable | *(inherited)* `@UpdateTimestamp` |
| updatedBy | String | nullable | *(inherited)* `@LastModifiedBy` |

Extras: no `@Paginate`; `@SQLRestriction(value = "deleted = false")`; `@LifeCycleHookBinding` present (five class-level bindings plus five field-level ones on `status`); `@Include(rootLevel = false)` — not root-queryable. `BRANCH_FIELDS_UPDATE` (`Job.java:70-71`) resolves to `(<Job.UPDATE>) AND (user is a super service OR user is an executor service OR executor manage job OR job branch fields editable)`.
Naming: field `job` on `organization` and on `workspace`.

## Conventions

- Root fields take `ids: [String]` (e.g. `organization(ids: ["<uuid>"])`,
  `notification_configuration(ids: [...])`). Root collections also accept
  `filter`. Root field names are the `@Entity(name = ...)` values verbatim,
  including snake_case. Only entities without `@Include(rootLevel = false)`
  exist at the root.
- Relationship collections take `first`, `after`, `sort`, `filter`.
- Every relationship is a Relay connection, including to-one relationships
  (`workspace { edges { node { id } } }`, not `workspace { id }`).
- `pageInfo` exposes `startCursor`, `endCursor`, `hasNextPage`,
  `hasPreviousPage`, `totalRecords`.
- `after` is an offset string (Elide encodes the cursor as the row offset).
- Default page size is Elide's 500-row ceiling unless `@Paginate` overrides it.
  Upstream raises it to 10000 for `project`, `workspace` and `workspacetag`
  (`defaultPageSize = 10000, maxPageSize = 10000`); asking for more must be
  rejected, not silently clamped.
- `first` and `after` are `StringOrInt` scalars: both integer and string
  literals must be accepted (the UI passes `String(request.first)`).
