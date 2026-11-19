# Upstream UI GraphQL operations

Extracted from `ui/src` at upstream `c460c846` (2026-11-19). Ten production
operations; the rest of the UI uses JSON:API REST (`axiosInstance`, 161 call
sites). The `.graphql` files in `ui-operations/` are verbatim copies of the
template literals, with `${...}` left where the UI interpolates and the two
dynamic documents rendered for a representative request.

| File | Call site | Root field | Entities | Arguments | Aliases |
| --- | --- | --- | --- | --- | --- |
| `notifications-list.graphql` | `domain/Notifications/NotificationConfigurationList.tsx:72` | `organization` | organization, notificationConfiguration, workspace | `ids` | – |
| `notifications-get.graphql` | `domain/Notifications/EditNotificationConfiguration.tsx:113` | `notification_configuration` | notificationConfiguration, workspace, triggers, templates | `ids` | – |
| `projects-list.graphql` | `modules/projects/projectService.ts:29` | `organization` | organization, project | `ids` | – |
| `projects-get.graphql` | `modules/projects/projectService.ts:114` | `organization` | organization, project | `ids` | – |
| `project-access-list.graphql` | `modules/projects/projectService.ts:169` | `project` | project, projectAccess | `ids` | – |
| `workspace-access-list.graphql` | `modules/workspaces/workspaceAccessService.ts:29` | `workspace` | workspace, access | `ids` | – |
| `workspaces-list.graphql` | `modules/workspaces/workspaceService.ts:17` | `organization` | organization, workspace, workspaceTag, project | `ids`, `sort` | – |
| `workspaces-tag-matches.graphql` | `modules/workspaces/workspaceService.ts:236` | `organization` | organization, workspace | `ids`, `first`, `filter` | one per tag filter |
| `workspaces-page.graphql` | `modules/workspaces/workspaceService.ts:324` | `organization` | organization, workspace, workspaceTag, project | `ids`, `first`, `after`, `filter`, `sort` | one per facet count |
| `runs-resolve.graphql` | `domain/Workspaces/RunRedirectRoute.tsx:27` | `organization` | organization, workspace, job | `filter` | – |

## Entity and relationship surface

Distinct types the UI reads: `organization`, `workspace`, `project`,
`projectAccess`, `access`, `notificationConfiguration`, `notificationTrigger`
(`triggers`), `notificationTemplate` (`templates`), `workspaceTag`, `job`.

Traversed paths:

```
organization.notificationConfiguration.workspace
notification_configuration.workspace
notification_configuration.triggers
notification_configuration.templates
organization.project
organization.workspace.workspaceTag
organization.workspace.project
organization.workspace.job
project.projectAccess
workspace.access
```

Root fields are the Elide `@Entity(name = ...)` values, so they are snake_case
where the entity name is (`notification_configuration`); relationship fields are
the Java field names in camelCase (`notificationConfiguration`).

## Connection shape

Every relationship is a Relay connection, including to-one relationships. The UI
comments on this explicitly: `workspace { id }` fails, `workspace { edges { node
{ id } } }` is required. `pageInfo` is read for `endCursor`, `hasNextPage` and
`totalRecords`; `workspaces-page` also uses `first`, `after`, `sort`, `filter`
and field aliasing to collect facet counts in one request.

`first` and `after` are declared as `StringOrInt` and are passed as strings
(`String(request.first)`), so the scalar must accept both.

## Filter expressions the UI generates

Verbatim from `workspaceService.ts`. These are the acceptance cases for the
RSQL parser; `quoteRsql` wraps a value in single quotes and escapes embedded
quotes.

```
(name=ini='*prod*',description=ini='*prod*')
project.id=isnull=true
project.id=='<uuid>'
(lastJobStatus=isnull=true,lastJobStatus=='NeverExecuted')
(lastJobStatus=='waitingApproval',job.status=='waitingApproval')
lastJobStatus=='running'
(policyComplianceStatus=isnull=true,policyComplianceStatus=='UNKNOWN')
workspaceTag.tagId=='<uuid>'
workspaceTag.tagId=='<uuid>';workspaceTag.value=='prod'
id=in=('<uuid1>','<uuid2>')
id=isnull=true
name=="<org>"          (runs-resolve, double quotes)
```

Filters are combined with `;` (AND) across dimensions and `,` (OR) inside a
dimension, with parentheses grouping.

## Sort strings

```
name,id          -name,-id
lastJobDate,id   -lastJobDate,-id
lastJobStatus,id source,id   -source,-id
```

Compound sort, ascending by default, `-` for descending, ties broken by `id`.

## Other contracts observed

- Root collections accept `filter` as well as `ids` (`organization(filter:)` in
  `runs-resolve`).
- `job(filter: "id==...")` filters a relationship collection by `id`.
- Computed attributes the UI filters and sorts on (`lastJobStatus`,
  `lastJobDate`, `policyComplianceStatus`, `locked`) exist as real columns in our
  schema, so they need no special evaluation.
- `first: "10000"` is requested for tag matching, which is the `@Paginate`
  maximum on `workspace`; asking for more must be rejected, not clamped
  silently.
