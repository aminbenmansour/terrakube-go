package rbac

import (
	"context"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type fakeOrganizationStore struct {
	teams []store.Team
	orgs  []string
}

func (f fakeOrganizationStore) ListTeamsByOrganizationAndNames(_ context.Context, arg store.ListTeamsByOrganizationAndNamesParams) ([]store.Team, error) {
	var out []store.Team
	for _, team := range f.teams {
		if team.OrganizationID == arg.OrganizationID && contains(arg.Names, team.Name) {
			out = append(out, team)
		}
	}
	return out, nil
}

func (f fakeOrganizationStore) ListOrganizationIDsByGroupNames(context.Context, []string) ([]string, error) {
	return f.orgs, nil
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

const testOrganization = "70000000-0000-0000-0000-000000000001"

func TestResolveOrganization(t *testing.T) {
	tests := []struct {
		name       string
		store      fakeOrganizationStore
		groups     []string
		wantMember bool
		wantVcs    bool
	}{
		{
			name: "team membership grants the role",
			store: fakeOrganizationStore{
				orgs: []string{testOrganization},
				teams: []store.Team{{
					Name: "writers", OrganizationID: testOrganization, Role: strPtr("admin"),
				}},
			},
			groups:     []string{"writers"},
			wantMember: true,
			wantVcs:    true,
		},
		{
			name:       "no group matches nothing",
			store:      fakeOrganizationStore{teams: []store.Team{{Name: "writers", OrganizationID: testOrganization, Role: strPtr("write")}}},
			groups:     []string{"strangers"},
			wantMember: false,
			wantVcs:    false,
		},
		{
			name:       "no groups skips the lookups",
			store:      fakeOrganizationStore{orgs: []string{testOrganization}},
			groups:     nil,
			wantMember: false,
			wantVcs:    false,
		},
		{
			name:       "a workspace access row is limited membership",
			store:      fakeOrganizationStore{orgs: []string{testOrganization}},
			groups:     []string{"limited"},
			wantMember: true,
			wantVcs:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			access, err := ResolveOrganization(context.Background(), test.store, testOrganization, test.groups)
			if err != nil {
				t.Fatalf("resolve: %v", err)
			}
			if access.Member != test.wantMember {
				t.Errorf("member = %v, want %v", access.Member, test.wantMember)
			}
			if access.Permissions.ManageVcs != test.wantVcs {
				t.Errorf("manageVcs = %v, want %v", access.Permissions.ManageVcs, test.wantVcs)
			}
		})
	}
}

func TestResolveOrganizationUsesTheRequestedOrganization(t *testing.T) {
	other := "70000000-0000-0000-0000-000000000099"
	fake := fakeOrganizationStore{
		orgs: []string{testOrganization},
		teams: []store.Team{
			{Name: "writers", OrganizationID: other, Role: strPtr("write")},
		},
	}
	access, err := ResolveOrganization(context.Background(), fake, testOrganization, []string{"writers"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if access.Permissions.ManageVcs {
		t.Error("a team of another organization must not grant permissions")
	}
}

func strPtr(value string) *string { return &value }

type fakeWorkspaceStore struct {
	fakeOrganizationStore
	workspace     store.Workspace
	project       store.Project
	access        []store.Access
	projectAccess []store.ProjectAccess
}

func (f fakeWorkspaceStore) GetWorkspace(context.Context, string) (store.Workspace, error) {
	return f.workspace, nil
}

func (f fakeWorkspaceStore) GetProject(context.Context, string) (store.Project, error) {
	return f.project, nil
}

func (f fakeWorkspaceStore) ListAccessByWorkspaceAndNames(_ context.Context, arg store.ListAccessByWorkspaceAndNamesParams) ([]store.Access, error) {
	var out []store.Access
	for _, row := range f.access {
		if row.WorkspaceID == arg.WorkspaceID && contains(arg.Names, row.Name) {
			out = append(out, row)
		}
	}
	return out, nil
}

func (f fakeWorkspaceStore) ListProjectAccessByProjectAndNames(_ context.Context, arg store.ListProjectAccessByProjectAndNamesParams) ([]store.ProjectAccess, error) {
	var out []store.ProjectAccess
	for _, row := range f.projectAccess {
		if row.ProjectID == arg.ProjectID && row.Name != nil && contains(arg.Names, *row.Name) {
			out = append(out, row)
		}
	}
	return out, nil
}

func TestResolveWorkspaceAddsGrantRows(t *testing.T) {
	project := "70000000-0000-0000-0000-000000000010"
	fake := fakeWorkspaceStore{
		fakeOrganizationStore: fakeOrganizationStore{orgs: []string{testOrganization}},
		workspace:             store.Workspace{ID: "ws", OrganizationID: testOrganization, ProjectID: &project},
		project:               store.Project{ID: project, OrganizationID: testOrganization},
		access: []store.Access{{
			WorkspaceID: "ws", Name: "team-x", ManageWorkspace: true,
		}},
		projectAccess: []store.ProjectAccess{{
			ProjectID: project, Name: strPtr("team-x"), ManageState: true,
		}},
	}

	access, err := ResolveWorkspace(context.Background(), fake, testOrganization, "ws", []string{"team-x"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if !access.Member {
		t.Error("member = false, want true")
	}
	if !access.Permissions.ManageWorkspace {
		t.Error("manageWorkspace = false, want the access row to grant it")
	}
	if !access.Permissions.ManageState {
		t.Error("manageState = false, want the project access row to grant it")
	}
}

func TestResolveWorkspaceIgnoresAnotherOrganization(t *testing.T) {
	fake := fakeWorkspaceStore{
		fakeOrganizationStore: fakeOrganizationStore{orgs: []string{testOrganization}},
		workspace:             store.Workspace{ID: "ws", OrganizationID: "other-org"},
		access: []store.Access{{
			WorkspaceID: "ws", Name: "team-x", ManageWorkspace: true,
		}},
	}

	access, err := ResolveWorkspace(context.Background(), fake, testOrganization, "ws", []string{"team-x"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if access.Permissions.ManageWorkspace {
		t.Error("a workspace of another organization must not grant permissions")
	}
}
