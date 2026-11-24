package rbac_test

import (
	"encoding/json"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/rbac"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// The matrix is upstream RbacV2Service: for each predefined role every check
// has a fixed answer, and "custom" falls back to the row's booleans.
func TestRoleMatrix(t *testing.T) {
	want := func(admin, write, plan, read bool) map[string]bool {
		return map[string]bool{
			"admin": admin, "write": write, "plan": plan, "read": read,
			"custom": false, "": false, "ADMIN": admin, " Write ": write,
		}
	}
	cases := []struct {
		name  string
		check func(rbac.Flags) bool
		want  map[string]bool
	}{
		{"manageWorkspace", rbac.CanManageWorkspace, want(true, true, false, false)},
		{"manageState", rbac.CanManageState, want(true, true, false, false)},
		{"manageModule", rbac.CanManageModule, want(true, false, false, false)},
		{"manageProvider", rbac.CanManageProvider, want(true, false, false, false)},
		{"manageVcs", rbac.CanManageVcs, want(true, false, false, false)},
		{"manageTemplate", rbac.CanManageTemplate, want(true, false, false, false)},
		{"manageCollection", rbac.CanManageCollection, want(true, false, false, false)},
		{"managePolicies", rbac.CanManagePolicies, want(true, false, false, false)},
		{"planJob", rbac.CanPlanJob, want(true, true, true, false)},
		{"approveJob", rbac.CanApproveJob, want(true, true, false, false)},
		{"manageJob", rbac.CanManageJob, want(true, true, true, false)},
		{"destroyJob", rbac.CanDestroyJob, want(true, true, false, false)},
	}
	for _, tc := range cases {
		for role, expectation := range tc.want {
			if got := tc.check(rbac.Flags{Role: role}); got != expectation {
				t.Errorf("%s(%q) = %v, want %v", tc.name, role, got, expectation)
			}
		}
	}
}

func TestCustomRoleUsesEveryFlag(t *testing.T) {
	all := rbac.Flags{
		Role: "custom", ManageWorkspace: true, ManageState: true, ManageModule: true,
		ManageProvider: true, ManageVcs: true, ManageTemplate: true, ManageCollection: true,
		ManagePolicies: true, PlanJob: true, ApproveJob: true, DestroyJob: true,
	}
	for name, check := range map[string]func(rbac.Flags) bool{
		"manageWorkspace": rbac.CanManageWorkspace, "manageState": rbac.CanManageState,
		"manageModule": rbac.CanManageModule, "manageProvider": rbac.CanManageProvider,
		"manageVcs": rbac.CanManageVcs, "manageTemplate": rbac.CanManageTemplate,
		"manageCollection": rbac.CanManageCollection, "managePolicies": rbac.CanManagePolicies,
		"planJob": rbac.CanPlanJob, "approveJob": rbac.CanApproveJob,
		"manageJob": rbac.CanManageJob, "destroyJob": rbac.CanDestroyJob,
	} {
		if !check(all) {
			t.Errorf("%s(custom with every flag) = false, want true", name)
		}
		if check(rbac.Flags{Role: "custom"}) {
			t.Errorf("%s(custom with no flags) = true, want false", name)
		}
	}
}

func TestFlagsFromRows(t *testing.T) {
	team := store.Team{
		ManageWorkspace: true, ManageModule: true, ManageProvider: true, ManageVcs: true,
		ManageTemplate: true, ManageState: true, ManageCollection: true, ManagePolicies: true,
		PlanJob: true, ApproveJob: true, DestroyJob: true, Role: stringPtr("custom"),
	}
	access := store.Access{
		ManageWorkspace: true, ManageState: true, PlanJob: true, ApproveJob: true,
		DestroyJob: true, Role: stringPtr("custom"),
	}
	projectAccess := store.ProjectAccess{
		ManageWorkspace: true, ManageState: true, PlanJob: true, ApproveJob: true,
		DestroyJob: true, Role: stringPtr("custom"),
	}

	for name, flags := range map[string]rbac.Flags{
		"team": rbac.FlagsFromTeam(team), "access": rbac.FlagsFromAccess(access),
		"project access": rbac.FlagsFromProjectAccess(projectAccess),
	} {
		if !rbac.CanManageWorkspace(flags) || !rbac.CanManageState(flags) ||
			!rbac.CanPlanJob(flags) || !rbac.CanApproveJob(flags) || !rbac.CanDestroyJob(flags) {
			t.Errorf("%s: the shared flags were dropped: %+v", name, flags)
		}
	}
	if !rbac.CanManageModule(rbac.FlagsFromTeam(team)) {
		t.Error("team: manageModule was dropped")
	}
	if rbac.CanManageModule(rbac.FlagsFromAccess(access)) {
		t.Error("access: manageModule is not an access column and must stay false")
	}
}

func TestPermissionsAccumulate(t *testing.T) {
	var permissions rbac.Permissions
	permissions.ApplyTeam(rbac.Flags{Role: "read"})
	// A predefined role ignores the flags: plan may queue but not approve.
	permissions.ApplyAccess(rbac.Flags{Role: "plan", ManageWorkspace: true, ApproveJob: true})
	permissions.ApplyProjectAccess(rbac.Flags{Role: "custom", ManageState: true, PlanJob: true})

	if !permissions.PlanJob || !permissions.ManageJob {
		t.Errorf("permissions = %+v, want the plan access to grant planning", permissions)
	}
	if permissions.ApproveJob {
		t.Errorf("permissions = %+v, want the plan role to ignore its approve flag", permissions)
	}
	if permissions.ManageWorkspace {
		t.Errorf("permissions = %+v, want the plan role to ignore its workspace flag", permissions)
	}
	if !permissions.ManageState {
		t.Errorf("permissions = %+v, want the custom project access flag", permissions)
	}
	if permissions.ManageModule || permissions.ManagePolicies || permissions.ManagePermission {
		t.Errorf("permissions = %+v, want the ungranted permissions false", permissions)
	}
}

func TestOwnerAddsOnlyTheTwoPermissionFlags(t *testing.T) {
	var permissions rbac.Permissions
	permissions.GrantOwner()

	if !permissions.ManagePermission || !permissions.ManagePolicies {
		t.Errorf("permissions = %+v, want owner to grant managePermission and managePolicies", permissions)
	}
	if permissions.ManageWorkspace || permissions.ManageState || permissions.ManageJob ||
		permissions.PlanJob || permissions.ApproveJob {
		t.Errorf("permissions = %+v, want owner not to grant operational permissions", permissions)
	}
}

func TestPermissionsJSONUsesUpstreamFieldNames(t *testing.T) {
	raw, err := json.Marshal(rbac.Permissions{ManagePermission: true, ManagePolicies: true, PlanJob: true})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := map[string]bool{
		"manageState": false, "manageWorkspace": false, "manageModule": false,
		"manageProvider": false, "manageVcs": false, "manageTemplate": false,
		"manageCollection": false, "manageJob": false, "planJob": true,
		"approveJob": false, "managePermission": true, "managePolicies": true,
	}
	var got map[string]bool
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("field count = %d, want %d: %s", len(got), len(want), raw)
	}
	for key, value := range want {
		if got[key] != value {
			t.Errorf("%s = %v, want %v", key, got[key], value)
		}
	}
}

func stringPtr(value string) *string { return &value }
