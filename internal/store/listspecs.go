package store

// The management API's list specs. Attribute names are the JSON:API names the
// UI sends in filter[type]; the SQL expressions are the columns they map to.
// An attribute missing here is one we do not model, and a filter naming it is
// ignored (see ListSpec.where).

// OrganizationList is scoped by the caller's visible organizations: pass the
// ids, or nil for an owner. An empty, non-nil slice matches nothing.
var OrganizationList = ListSpec{
	Table: "organization",
	Base:  "($1::text[] IS NULL OR id = ANY($1::text[]))",
	Attrs: map[string]Attr{
		"name":        {Expr: "name"},
		"description": {Expr: "description"},
		"id":          {Expr: "id"},
	},
	Sorts: map[string]string{"name": "lower(name)", "id": "id"},
}

var ProjectList = ListSpec{
	Table: "project",
	Base:  "organization_id = $1",
	Attrs: map[string]Attr{
		"name":        {Expr: "name"},
		"description": {Expr: "description"},
		"id":          {Expr: "id"},
	},
	Sorts: map[string]string{"name": "lower(name)", "id": "id"},
}

var TagList = ListSpec{
	Table: "tag",
	Base:  "organization_id = $1",
	Attrs: map[string]Attr{
		"name": {Expr: "name"},
		"id":   {Expr: "id"},
	},
	Sorts: map[string]string{"name": "lower(name)", "id": "id"},
}

var WorkspaceList = ListSpec{
	Table: "workspace",
	Base:  "organization_id = $1 AND deleted = false",
	Attrs: map[string]Attr{
		"name":                   {Expr: "name"},
		"description":            {Expr: "description"},
		"id":                     {Expr: "id"},
		"lastJobStatus":          {Expr: "last_job_status"},
		"policyComplianceStatus": {Expr: "policy_compliance_status"},
		"project.id":             {Expr: "project_id"},
		"lastJobDate":            {Expr: "last_job_date", Kind: KindTime},
	},
	Sorts: map[string]string{"name": "lower(name)", "id": "id", "lastJobDate": "last_job_date"},
}

var ProjectWorkspaceList = ListSpec{
	Table: "workspace",
	Base:  "project_id = $1 AND deleted = false",
	Attrs: WorkspaceList.Attrs,
	Sorts: WorkspaceList.Sorts,
}
