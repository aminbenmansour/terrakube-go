package graphql

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jobstatus"
	"github.com/aminbenmansour/terrakube-go/internal/query"
)

// DB is the slice of pgx the resolvers need, so a test can supply a pool
// without the whole application.
type DB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// collectionSpec describes one Elide collection: where it lives, how it is
// selected, which filter and sort paths it accepts, and its default order.
// The field map is the injection barrier: a path that is not in it never
// reaches the SQL text.
type collectionSpec struct {
	table      string
	selectList string
	fields     map[string]query.Field
	order      []query.Order
}

func column(name string) query.Field { return query.Field{Column: name} }

func predicate(sql string) query.Field { return query.Field{Predicate: sql} }

// relationship filters that the UI sends. The correlated subquery is why a
// workspace can be filtered by a tag or by the status of one of its jobs.
var (
	workspaceTagIDFilter = predicate("EXISTS (SELECT 1 FROM workspacetag wt WHERE wt.workspace_id = workspace.id AND wt.tag_id = ?)")
	workspaceTagFilter   = predicate("EXISTS (SELECT 1 FROM workspacetag wt WHERE wt.workspace_id = workspace.id AND wt.tag_value = ?)")
	workspaceJobFilter   = predicate("EXISTS (SELECT 1 FROM job j WHERE j.workspace_id = workspace.id AND j.status = ?)")
)

var organizationSpec = collectionSpec{
	table:      "organization",
	selectList: "id, name, description, disabled, execution_mode, icon",
	fields: map[string]query.Field{
		"id":                            column("id"),
		"name":                          column("name"),
		"description":                   column("description"),
		"disabled":                      column("disabled"),
		"executionMode":                 column("execution_mode"),
		"icon":                          column("icon"),
		"project.name":                  predicate("EXISTS (SELECT 1 FROM project p WHERE p.organization_id = organization.id AND p.name = ?)"),
		"workspace.name":                predicate("EXISTS (SELECT 1 FROM workspace w WHERE w.organization_id = organization.id AND w.name = ?)"),
		"notification_configuration.id": predicate("EXISTS (SELECT 1 FROM notification_configuration n WHERE n.organization_id = organization.id AND n.id = ?)"),
	},
	order: []query.Order{{Column: "name"}, {Column: "id"}},
}

var workspaceSpec = collectionSpec{
	table: "workspace",
	selectList: `id, name, description, source, branch, folder,
		to_char(last_job_date AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS last_job_date,
		terraform_version, iac_type, last_job_status, locked, policy_compliance_status,
		allow_remote_apply, execution_mode, project_id`,
	fields: map[string]query.Field{
		"id":                     column("id"),
		"name":                   column("name"),
		"description":            column("description"),
		"source":                 column("source"),
		"branch":                 column("branch"),
		"folder":                 column("folder"),
		"terraformVersion":       column("terraform_version"),
		"iacType":                column("iac_type"),
		"lastJobStatus":          column("last_job_status"),
		"lastJobDate":            column("last_job_date"),
		"locked":                 column("locked"),
		"policyComplianceStatus": column("policy_compliance_status"),
		"project.id":             column("project_id"),
		"workspaceTag.tagId":     workspaceTagIDFilter,
		"workspaceTag.value":     workspaceTagFilter,
		"job.status":             workspaceJobFilter,
	},
	order: []query.Order{{Column: "name"}, {Column: "id"}},
}

var projectSpec = collectionSpec{
	table:      "project",
	selectList: "id, name, description",
	fields: map[string]query.Field{
		"id":          column("id"),
		"name":        column("name"),
		"description": column("description"),
	},
	order: []query.Order{{Column: "name"}, {Column: "id"}},
}

var workspaceTagSpec = collectionSpec{
	table:      "workspacetag",
	selectList: "id, tag_id, tag_value",
	fields: map[string]query.Field{
		"id":     column("id"),
		"tagId":  column("tag_id"),
		"value":  column("tag_value"),
		"tag.id": column("tag_id"),
	},
	order: []query.Order{{Column: "id"}},
}

var jobSpec = collectionSpec{
	table: "job",
	selectList: `id, status, command,
		to_char(created_date AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_date`,
	fields: map[string]query.Field{
		"id":      column("id"),
		"status":  column("status"),
		"command": column("command"),
	},
	order: []query.Order{{Column: "id", Desc: true}},
}

var templateSpec = collectionSpec{
	table:      "template",
	selectList: "id, name, description, version",
	fields: map[string]query.Field{
		"id":          column("id"),
		"name":        column("name"),
		"description": column("description"),
		"version":     column("version"),
	},
	order: []query.Order{{Column: "name"}, {Column: "id"}},
}

var projectAccessSpec = collectionSpec{
	table:      "project_access",
	selectList: "id, name, role, manage_workspace, manage_state, manage_job, plan_job, approve_job, destroy_job",
	fields: map[string]query.Field{
		"id":   column("id"),
		"name": column("name"),
		"role": column("role"),
	},
	order: []query.Order{{Column: "name"}, {Column: "id"}},
}

var accessSpec = collectionSpec{
	table:      "access",
	selectList: "id, name, role, manage_workspace, manage_state, plan_job, approve_job, destroy_job",
	fields: map[string]query.Field{
		"id":   column("id"),
		"name": column("name"),
		"role": column("role"),
	},
	order: []query.Order{{Column: "name"}, {Column: "id"}},
}

var notificationConfigurationSpec = collectionSpec{
	table:      "notification_configuration",
	selectList: "id, name, description, channel_type, destination_url, active, message_style, workspace_id",
	fields: map[string]query.Field{
		"id":             column("id"),
		"name":           column("name"),
		"description":    column("description"),
		"channelType":    column("channel_type"),
		"destinationUrl": column("destination_url"),
		"active":         column("active"),
		"messageStyle":   column("message_style"),
		"workspace.id":   column("workspace_id"),
	},
	order: []query.Order{{Column: "name"}, {Column: "id"}},
}

var notificationTriggerSpec = collectionSpec{
	table:      "notification_trigger",
	selectList: "id, job_status",
	fields: map[string]query.Field{
		"id":        column("id"),
		"jobStatus": column("job_status"),
	},
	order: []query.Order{{Column: "id"}},
}

// listArgs are the arguments every Elide connection accepts.
type listArgs struct {
	IDs    []string
	Filter *string
	First  *StringOrInt
	After  *StringOrInt
	Sort   *string
}

// page is one loaded collection page plus the offset it started at, which is
// what the edge cursors are derived from.
type page[T any] struct {
	Nodes  []T
	Info   *PageInfo
	Offset int64
}

func (p page[T]) edges() []string {
	cursors := make([]string, len(p.Nodes))
	for i := range p.Nodes {
		cursors[i] = strconv.FormatInt(p.Offset+int64(i), 10)
	}
	return cursors
}

// build renders the count and page queries. Filter and sort values go through
// the spec's whitelist; the cursor is an offset, as Elide's is.
func (s collectionSpec) build(scope string, scopeArgs []any, args listArgs) (sqlText, countText string, sqlArgs []any, offset int64, err error) {
	limit, offset, err := query.ParsePage(args.First.Int(), afterString(args.After))
	if err != nil {
		return "", "", nil, 0, err
	}

	var filter query.Node
	if args.Filter != nil && strings.TrimSpace(*args.Filter) != "" {
		filter, err = query.ParseRSQL(*args.Filter)
		if err != nil {
			return "", "", nil, 0, err
		}
	}

	var where []string
	values := make([]any, 0, len(args.IDs)+4)
	if scope != "" {
		where = append(where, scope)
		values = append(values, scopeArgs...)
	}
	if len(args.IDs) > 0 {
		values = append(values, args.IDs)
		where = append(where, fmt.Sprintf("id = ANY($%d)", len(values)))
	}
	filterSQL, filterArgs, err := query.Where(filter, s.fields, len(values))
	if err != nil {
		return "", "", nil, 0, err
	}
	if filterSQL != "" {
		where = append(where, filterSQL)
		values = append(values, filterArgs...)
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = " WHERE " + strings.Join(where, " AND ")
	}

	orders := s.order
	if args.Sort != nil && strings.TrimSpace(*args.Sort) != "" {
		orders, err = query.ParseSort(*args.Sort, s.fields)
		if err != nil {
			return "", "", nil, 0, err
		}
	}

	countText = "SELECT count(*) FROM " + s.table + whereSQL
	sqlText = "SELECT " + s.selectList + " FROM " + s.table + whereSQL
	if orderBy := query.OrderBy(orders); orderBy != "" {
		sqlText += " " + orderBy
	}
	sqlText += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(values)+1, len(values)+2)
	sqlArgs = append(values, limit, offset)
	return sqlText, countText, sqlArgs, offset, nil
}

func loadPage[T any](ctx context.Context, db DB, spec collectionSpec, scope string, scopeArgs []any, args listArgs) (page[T], error) {
	sqlText, countText, sqlArgs, offset, err := spec.build(scope, scopeArgs, args)
	if err != nil {
		return page[T]{}, err
	}

	var total int
	countArgs := sqlArgs[:len(sqlArgs)-2]
	if err := db.QueryRow(ctx, countText, countArgs...).Scan(&total); err != nil {
		return page[T]{}, fmt.Errorf("count %s: %w", spec.table, err)
	}

	rows, err := db.Query(ctx, sqlText, sqlArgs...)
	if err != nil {
		return page[T]{}, fmt.Errorf("list %s: %w", spec.table, err)
	}
	nodes, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[T])
	if err != nil {
		return page[T]{}, fmt.Errorf("scan %s: %w", spec.table, err)
	}

	info := &PageInfo{
		HasNextPage:     offset+int64(len(nodes)) < int64(total),
		HasPreviousPage: offset > 0,
		TotalRecords:    total,
	}
	if len(nodes) > 0 {
		start := strconv.FormatInt(offset, 10)
		end := strconv.FormatInt(offset+int64(len(nodes))-1, 10)
		info.StartCursor = &start
		info.EndCursor = &end
	}
	return page[T]{Nodes: nodes, Info: info, Offset: offset}, nil
}

func afterString(after *StringOrInt) *string {
	if after == nil {
		return nil
	}
	value := string(*after)
	return &value
}

// idsFrom drops the nils a nullable list element allows.
func idsFrom(ids []*string) []string {
	if len(ids) == 0 {
		return nil
	}
	values := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != nil {
			values = append(values, *id)
		}
	}
	return values
}

// The generated edge types differ per entity, so each connection needs its own
// tiny wrapper. Cursors are offsets, which is what the UI passes back as
// `after`.

func workspaceEdges(p page[Workspace]) []WorkspaceEdge {
	cursors := p.edges()
	edges := make([]WorkspaceEdge, len(p.Nodes))
	for i := range p.Nodes {
		edges[i] = WorkspaceEdge{Node: &p.Nodes[i], Cursor: cursors[i]}
	}
	return edges
}

func projectEdges(p page[Project]) []ProjectEdge {
	cursors := p.edges()
	edges := make([]ProjectEdge, len(p.Nodes))
	for i := range p.Nodes {
		edges[i] = ProjectEdge{Node: &p.Nodes[i], Cursor: cursors[i]}
	}
	return edges
}

func workspaceTagEdges(p page[WorkspaceTag]) []WorkspaceTagEdge {
	cursors := p.edges()
	edges := make([]WorkspaceTagEdge, len(p.Nodes))
	for i := range p.Nodes {
		edges[i] = WorkspaceTagEdge{Node: &p.Nodes[i], Cursor: cursors[i]}
	}
	return edges
}

func jobEdges(p page[Job]) []JobEdge {
	cursors := p.edges()
	edges := make([]JobEdge, len(p.Nodes))
	for i := range p.Nodes {
		edges[i] = JobEdge{Node: &p.Nodes[i], Cursor: cursors[i]}
	}
	return edges
}

func templateEdges(p page[Template]) []TemplateEdge {
	cursors := p.edges()
	edges := make([]TemplateEdge, len(p.Nodes))
	for i := range p.Nodes {
		edges[i] = TemplateEdge{Node: &p.Nodes[i], Cursor: cursors[i]}
	}
	return edges
}

func projectAccessEdges(p page[ProjectAccess]) []ProjectAccessEdge {
	cursors := p.edges()
	edges := make([]ProjectAccessEdge, len(p.Nodes))
	for i := range p.Nodes {
		edges[i] = ProjectAccessEdge{Node: &p.Nodes[i], Cursor: cursors[i]}
	}
	return edges
}

func accessEdges(p page[Access]) []AccessEdge {
	cursors := p.edges()
	edges := make([]AccessEdge, len(p.Nodes))
	for i := range p.Nodes {
		edges[i] = AccessEdge{Node: &p.Nodes[i], Cursor: cursors[i]}
	}
	return edges
}

func notificationConfigurationEdges(p page[NotificationConfiguration]) []NotificationConfigurationEdge {
	cursors := p.edges()
	edges := make([]NotificationConfigurationEdge, len(p.Nodes))
	for i := range p.Nodes {
		edges[i] = NotificationConfigurationEdge{Node: &p.Nodes[i], Cursor: cursors[i]}
	}
	return edges
}

func notificationTriggerEdges(p page[NotificationTrigger]) []NotificationTriggerEdge {
	cursors := p.edges()
	edges := make([]NotificationTriggerEdge, len(p.Nodes))
	for i := range p.Nodes {
		edges[i] = NotificationTriggerEdge{Node: &p.Nodes[i], Cursor: cursors[i]}
	}
	return edges
}

// streamEvents forwards the bus events a subscription asked for, and closes the
// stream when the client disconnects.
func streamEvents(ctx context.Context, r *Resolver, keep func(jobstatus.Event) bool) (<-chan *JobStatusEvent, error) {
	if r.events == nil {
		return nil, errors.New("subscriptions are not configured: set a Redis URL")
	}
	events, err := r.events.SubscribeJobStatus(ctx)
	if err != nil {
		return nil, err
	}

	out := make(chan *JobStatusEvent)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-events:
				if !ok {
					return
				}
				if !keep(event) {
					continue
				}
				select {
				case out <- &JobStatusEvent{JobID: int(event.JobID), WorkspaceID: event.WorkspaceID, Status: event.Status}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out, nil
}
