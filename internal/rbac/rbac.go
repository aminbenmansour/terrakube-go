// Package rbac resolves Terrakube team permissions the way upstream
// RbacV2Service does: a predefined role decides, and "custom" falls back to
// the row's boolean flags.
package rbac

import (
	"strings"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	roleAdmin  = "admin"
	roleWrite  = "write"
	rolePlan   = "plan"
	roleRead   = "read"
	roleCustom = "custom"
)

// Flags are the permission columns of a team, access or project_access row.
// Access rows have no module, provider, VCS, template, collection or policy
// columns, so those stay false for them.
type Flags struct {
	Role             string
	ManageWorkspace  bool
	ManageState      bool
	ManageModule     bool
	ManageProvider   bool
	ManageVcs        bool
	ManageTemplate   bool
	ManageCollection bool
	ManagePolicies   bool
	PlanJob          bool
	ApproveJob       bool
	DestroyJob       bool
}

// FlagsFromTeam reads a team row.
func FlagsFromTeam(team store.Team) Flags {
	return Flags{
		Role:             stringValue(team.Role),
		ManageWorkspace:  team.ManageWorkspace,
		ManageState:      team.ManageState,
		ManageModule:     team.ManageModule,
		ManageProvider:   team.ManageProvider,
		ManageVcs:        team.ManageVcs,
		ManageTemplate:   team.ManageTemplate,
		ManageCollection: team.ManageCollection,
		ManagePolicies:   team.ManagePolicies,
		PlanJob:          team.PlanJob,
		ApproveJob:       team.ApproveJob,
		DestroyJob:       team.DestroyJob,
	}
}

// FlagsFromAccess reads a workspace access row.
func FlagsFromAccess(access store.Access) Flags {
	return Flags{
		Role:            stringValue(access.Role),
		ManageWorkspace: access.ManageWorkspace,
		ManageState:     access.ManageState,
		PlanJob:         access.PlanJob,
		ApproveJob:      access.ApproveJob,
		DestroyJob:      access.DestroyJob,
	}
}

// FlagsFromProjectAccess reads a project access row.
func FlagsFromProjectAccess(access store.ProjectAccess) Flags {
	return Flags{
		Role:            stringValue(access.Role),
		ManageWorkspace: access.ManageWorkspace,
		ManageState:     access.ManageState,
		PlanJob:         access.PlanJob,
		ApproveJob:      access.ApproveJob,
		DestroyJob:      access.DestroyJob,
	}
}

// CanManageWorkspace reports whether the row may change workspace settings.
func CanManageWorkspace(f Flags) bool {
	switch role(f) {
	case roleAdmin, roleWrite:
		return true
	case rolePlan, roleRead:
		return false
	case roleCustom:
		return f.ManageWorkspace
	}
	return false
}

// CanManageState reports whether the row may read and write remote state.
func CanManageState(f Flags) bool {
	switch role(f) {
	case roleAdmin, roleWrite:
		return true
	case rolePlan, roleRead:
		return false
	case roleCustom:
		return f.ManageState
	}
	return false
}

// CanManageModule reports whether the row may manage the module registry.
func CanManageModule(f Flags) bool { return adminOnly(f, f.ManageModule) }

// CanManageProvider reports whether the row may manage the provider registry.
func CanManageProvider(f Flags) bool { return adminOnly(f, f.ManageProvider) }

// CanManageVcs reports whether the row may manage VCS connections.
func CanManageVcs(f Flags) bool { return adminOnly(f, f.ManageVcs) }

// CanManageTemplate reports whether the row may manage templates.
func CanManageTemplate(f Flags) bool { return adminOnly(f, f.ManageTemplate) }

// CanManageCollection reports whether the row may manage collections.
func CanManageCollection(f Flags) bool { return adminOnly(f, f.ManageCollection) }

// CanManagePolicies reports whether the row may manage policy sets.
func CanManagePolicies(f Flags) bool { return adminOnly(f, f.ManagePolicies) }

// CanPlanJob reports whether the row may queue a plan.
func CanPlanJob(f Flags) bool {
	switch role(f) {
	case roleAdmin, roleWrite, rolePlan:
		return true
	case roleRead:
		return false
	case roleCustom:
		return f.PlanJob
	}
	return false
}

// CanApproveJob reports whether the row may approve or reject a run.
func CanApproveJob(f Flags) bool {
	switch role(f) {
	case roleAdmin, roleWrite:
		return true
	case rolePlan, roleRead:
		return false
	case roleCustom:
		return f.ApproveJob
	}
	return false
}

// CanManageJob reports whether the row may do anything with jobs.
func CanManageJob(f Flags) bool { return CanPlanJob(f) || CanApproveJob(f) }

// CanDestroyJob reports whether the row may destroy without approval.
func CanDestroyJob(f Flags) bool {
	switch role(f) {
	case roleAdmin, roleWrite:
		return true
	case rolePlan, roleRead:
		return false
	case roleCustom:
		return f.DestroyJob
	}
	return false
}

// Permissions is upstream's PermissionSet, the body of the UI's permission
// endpoints. The field order matches the JSON the UI reads.
type Permissions struct {
	ManageState      bool `json:"manageState"`
	ManageWorkspace  bool `json:"manageWorkspace"`
	ManageModule     bool `json:"manageModule"`
	ManageProvider   bool `json:"manageProvider"`
	ManageVcs        bool `json:"manageVcs"`
	ManageTemplate   bool `json:"manageTemplate"`
	ManageCollection bool `json:"manageCollection"`
	ManageJob        bool `json:"manageJob"`
	PlanJob          bool `json:"planJob"`
	ApproveJob       bool `json:"approveJob"`
	ManagePermission bool `json:"managePermission"`
	ManagePolicies   bool `json:"managePolicies"`
}

// ApplyTeam adds every permission an organization team grants.
func (p *Permissions) ApplyTeam(f Flags) {
	p.ManageState = p.ManageState || CanManageState(f)
	p.ManageWorkspace = p.ManageWorkspace || CanManageWorkspace(f)
	p.ManageModule = p.ManageModule || CanManageModule(f)
	p.ManageProvider = p.ManageProvider || CanManageProvider(f)
	p.ManageVcs = p.ManageVcs || CanManageVcs(f)
	p.ManageTemplate = p.ManageTemplate || CanManageTemplate(f)
	p.ManageCollection = p.ManageCollection || CanManageCollection(f)
	p.ManageJob = p.ManageJob || CanManageJob(f)
	p.PlanJob = p.PlanJob || CanPlanJob(f)
	p.ApproveJob = p.ApproveJob || CanApproveJob(f)
	p.ManagePolicies = p.ManagePolicies || CanManagePolicies(f)
	// Upstream ties the Settings > Teams permission to workspace management.
	p.ManagePermission = p.ManagePermission || CanManageWorkspace(f)
}

// ApplyAccess adds the permissions a workspace access row grants.
func (p *Permissions) ApplyAccess(f Flags) {
	p.ManageState = p.ManageState || CanManageState(f)
	p.ManageWorkspace = p.ManageWorkspace || CanManageWorkspace(f)
	p.ManageJob = p.ManageJob || CanManageJob(f)
	p.PlanJob = p.PlanJob || CanPlanJob(f)
	p.ApproveJob = p.ApproveJob || CanApproveJob(f)
}

// ApplyProjectAccess adds the permissions a project access row grants.
func (p *Permissions) ApplyProjectAccess(f Flags) {
	p.ApplyAccess(f)
}

// GrantOwner is upstream's instance-owner branch: only these two flags.
func (p *Permissions) GrantOwner() {
	p.ManagePermission = true
	p.ManagePolicies = true
}

func adminOnly(f Flags, custom bool) bool {
	switch role(f) {
	case roleAdmin:
		return true
	case roleWrite, rolePlan, roleRead:
		return false
	case roleCustom:
		return custom
	}
	return false
}

// role is upstream's normalizeRole: blank means custom, everything else is
// lower-cased and trimmed, and an unknown name grants nothing.
func role(f Flags) string {
	value := strings.ToLower(strings.TrimSpace(f.Role))
	if value == "" {
		return roleCustom
	}
	return value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
