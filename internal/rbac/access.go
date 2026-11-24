package rbac

import (
	"context"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// Access is the caller's standing in one organization. Member is true when a
// team, workspace access row or project access row names one of the caller's
// groups; Permissions accumulates what those rows allow.
type Access struct {
	Member      bool
	Permissions Permissions
}

// OrganizationStore is the slice of the generated queries the resolver needs.
type OrganizationStore interface {
	ListTeamsByOrganizationAndNames(ctx context.Context, arg store.ListTeamsByOrganizationAndNamesParams) ([]store.Team, error)
	ListOrganizationIDsByGroupNames(ctx context.Context, names []string) ([]string, error)
}

// WorkspaceStore adds the workspace lookups to OrganizationStore.
type WorkspaceStore interface {
	OrganizationStore
	GetWorkspace(ctx context.Context, id string) (store.Workspace, error)
	GetProject(ctx context.Context, id string) (store.Project, error)
	ListAccessByWorkspaceAndNames(ctx context.Context, arg store.ListAccessByWorkspaceAndNamesParams) ([]store.Access, error)
	ListProjectAccessByProjectAndNames(ctx context.Context, arg store.ListProjectAccessByProjectAndNamesParams) ([]store.ProjectAccess, error)
}

// ResolveOrganization answers the caller's membership and team permissions.
// Upstream has no membership table: a team is a named permission record, so
// the token's groups are compared with the team names.
func ResolveOrganization(ctx context.Context, q OrganizationStore, organizationID string, groups []string) (Access, error) {
	var access Access
	if len(groups) == 0 {
		return access, nil
	}

	organizations, err := q.ListOrganizationIDsByGroupNames(ctx, groups)
	if err != nil {
		return Access{}, err
	}
	if !slices.Contains(organizations, organizationID) {
		return access, nil
	}
	access.Member = true

	teams, err := q.ListTeamsByOrganizationAndNames(ctx, store.ListTeamsByOrganizationAndNamesParams{
		OrganizationID: organizationID,
		Names:          groups,
	})
	if err != nil {
		return Access{}, err
	}
	for _, team := range teams {
		access.Permissions.ApplyTeam(FlagsFromTeam(team))
	}
	return access, nil
}

// ResolveWorkspace adds the workspace and project grant rows to the
// organization access.
func ResolveWorkspace(ctx context.Context, q WorkspaceStore, organizationID, workspaceID string, groups []string) (Access, error) {
	access, err := ResolveOrganization(ctx, q, organizationID, groups)
	if err != nil || !access.Member {
		return access, err
	}

	workspace, err := q.GetWorkspace(ctx, workspaceID)
	if errors.Is(err, pgx.ErrNoRows) {
		return access, nil
	}
	if err != nil {
		return Access{}, err
	}
	return ApplyWorkspaceGrants(ctx, q, access, organizationID, workspace, groups)
}

// ApplyWorkspaceGrants adds one workspace's grant rows to an access already
// resolved for its organization. A workspace of another organization
// contributes nothing; the parent guard answers that request as a 404.
func ApplyWorkspaceGrants(ctx context.Context, q WorkspaceStore, access Access, organizationID string, workspace store.Workspace, groups []string) (Access, error) {
	if !access.Member || workspace.OrganizationID != organizationID {
		return access, nil
	}

	rows, err := q.ListAccessByWorkspaceAndNames(ctx, store.ListAccessByWorkspaceAndNamesParams{
		WorkspaceID: workspace.ID,
		Names:       groups,
	})
	if err != nil {
		return Access{}, err
	}
	for _, row := range rows {
		access.Permissions.ApplyAccess(FlagsFromAccess(row))
	}

	if workspace.ProjectID != nil && *workspace.ProjectID != "" {
		projectRows, err := q.ListProjectAccessByProjectAndNames(ctx, store.ListProjectAccessByProjectAndNamesParams{
			ProjectID: *workspace.ProjectID,
			Names:     groups,
		})
		if err != nil {
			return Access{}, err
		}
		for _, row := range projectRows {
			access.Permissions.ApplyProjectAccess(FlagsFromProjectAccess(row))
		}
	}
	return access, nil
}

// ResolveProject adds the project access rows to the organization access.
func ResolveProject(ctx context.Context, q WorkspaceStore, organizationID, projectID string, groups []string) (Access, error) {
	access, err := ResolveOrganization(ctx, q, organizationID, groups)
	if err != nil || !access.Member {
		return access, err
	}

	project, err := q.GetProject(ctx, projectID)
	if errors.Is(err, pgx.ErrNoRows) {
		return access, nil
	}
	if err != nil {
		return Access{}, err
	}
	return ApplyProjectGrants(ctx, q, access, organizationID, project, groups)
}

// ApplyProjectGrants adds one project's grant rows to an access already
// resolved for its organization.
func ApplyProjectGrants(ctx context.Context, q WorkspaceStore, access Access, organizationID string, project store.Project, groups []string) (Access, error) {
	if !access.Member || project.OrganizationID != organizationID {
		return access, nil
	}

	rows, err := q.ListProjectAccessByProjectAndNames(ctx, store.ListProjectAccessByProjectAndNamesParams{
		ProjectID: project.ID,
		Names:     groups,
	})
	if err != nil {
		return Access{}, err
	}
	for _, row := range rows {
		access.Permissions.ApplyProjectAccess(FlagsFromProjectAccess(row))
	}
	return access, nil
}
