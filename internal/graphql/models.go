package graphql

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// StringOrInt is the scalar Elide uses for the paging arguments. The upstream
// UI sends both ("10000" and 20), so both are accepted.
type StringOrInt string

func (s StringOrInt) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(string(s))))
}

func (s *StringOrInt) UnmarshalGQL(v any) error {
	switch value := v.(type) {
	case string:
		*s = StringOrInt(value)
	case int:
		*s = StringOrInt(strconv.Itoa(value))
	case int64:
		*s = StringOrInt(strconv.FormatInt(value, 10))
	case json.Number:
		*s = StringOrInt(value.String())
	case float64:
		*s = StringOrInt(strconv.FormatFloat(value, 'f', -1, 64))
	default:
		return fmt.Errorf("StringOrInt: unsupported value %T", v)
	}
	return nil
}

// Int converts the scalar for the paging helper. An unparseable value yields
// nil, which means "use the default".
func (s *StringOrInt) Int() *int {
	if s == nil {
		return nil
	}
	value, err := strconv.Atoi(string(*s))
	if err != nil {
		return nil
	}
	return &value
}

// The model types below are scanned straight from SQL rows (pgx matches the
// db tag to the selected column alias), so nullable columns are pointers.

type Organization struct {
	ID            string  `db:"id"`
	Name          string  `db:"name"`
	Description   *string `db:"description"`
	Disabled      bool    `db:"disabled"`
	ExecutionMode string  `db:"execution_mode"`
	Icon          *string `db:"icon"`
}

type Workspace struct {
	ID                     string  `db:"id"`
	Name                   string  `db:"name"`
	Description            *string `db:"description"`
	Source                 *string `db:"source"`
	Branch                 *string `db:"branch"`
	Folder                 *string `db:"folder"`
	TerraformVersion       *string `db:"terraform_version"`
	IacType                string  `db:"iac_type"`
	LastJobStatus          string  `db:"last_job_status"`
	LastJobDate            *string `db:"last_job_date"`
	Locked                 bool    `db:"locked"`
	PolicyComplianceStatus string  `db:"policy_compliance_status"`
	AllowRemoteApply       bool    `db:"allow_remote_apply"`
	ExecutionMode          *string `db:"execution_mode"`
	// ProjectID backs the project relationship; it is not a schema field.
	ProjectID *string `db:"project_id"`
}

type Project struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
}

type WorkspaceTag struct {
	ID    string  `db:"id"`
	TagID string  `db:"tag_id"`
	Value *string `db:"tag_value"`
}

type Job struct {
	ID          int     `db:"id"`
	Status      string  `db:"status"`
	Command     *string `db:"command"`
	CreatedDate *string `db:"created_date"`
}

type Template struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	Version     *string `db:"version"`
}

type ProjectAccess struct {
	ID              string  `db:"id"`
	Name            *string `db:"name"`
	Role            *string `db:"role"`
	ManageWorkspace bool    `db:"manage_workspace"`
	ManageState     bool    `db:"manage_state"`
	ManageJob       bool    `db:"manage_job"`
	PlanJob         bool    `db:"plan_job"`
	ApproveJob      bool    `db:"approve_job"`
	DestroyJob      bool    `db:"destroy_job"`
}

type Access struct {
	ID              string  `db:"id"`
	Name            string  `db:"name"`
	Role            *string `db:"role"`
	ManageWorkspace bool    `db:"manage_workspace"`
	ManageState     bool    `db:"manage_state"`
	PlanJob         bool    `db:"plan_job"`
	ApproveJob      bool    `db:"approve_job"`
	DestroyJob      bool    `db:"destroy_job"`
}

type NotificationConfiguration struct {
	ID             string  `db:"id"`
	Name           string  `db:"name"`
	Description    *string `db:"description"`
	ChannelType    *string `db:"channel_type"`
	DestinationURL *string `db:"destination_url"`
	Active         bool    `db:"active"`
	MessageStyle   *string `db:"message_style"`
	// WorkspaceID backs the workspace relationship; it is not a schema field.
	WorkspaceID *string `db:"workspace_id"`
}

type NotificationTrigger struct {
	ID        string  `db:"id"`
	JobStatus *string `db:"job_status"`
}
