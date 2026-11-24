package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// grantAttrs are the fields the UI's access editor sends for both grant kinds.
type grantAttrs struct {
	Name            *string `json:"name"`
	Role            *string `json:"role"`
	ManageWorkspace *bool   `json:"manageWorkspace"`
	ManageState     *bool   `json:"manageState"`
	ManageJob       *bool   `json:"manageJob"`
	PlanJob         *bool   `json:"planJob"`
	ApproveJob      *bool   `json:"approveJob"`
	DestroyJob      *bool   `json:"destroyJob"`
}

func accessResource(access store.Access) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "access",
		ID:   access.ID,
		Attributes: map[string]any{
			"approveJob":      access.ApproveJob,
			"destroyJob":      access.DestroyJob,
			"manageJob":       access.ManageJob,
			"manageState":     access.ManageState,
			"manageWorkspace": access.ManageWorkspace,
			"name":            access.Name,
			"planJob":         access.PlanJob,
			"role":            access.Role,
		},
		Relationships: map[string]jsonapi.Relationship{
			"workspace": jsonapi.Link("workspace", access.WorkspaceID),
		},
	}
}

func projectAccessResource(access store.ProjectAccess) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "project_access",
		ID:   access.ID,
		Attributes: map[string]any{
			"approveJob":      access.ApproveJob,
			"destroyJob":      access.DestroyJob,
			"manageJob":       access.ManageJob,
			"manageState":     access.ManageState,
			"manageWorkspace": access.ManageWorkspace,
			"name":            access.Name,
			"planJob":         access.PlanJob,
			"role":            access.Role,
		},
		Relationships: map[string]jsonapi.Relationship{
			"project": jsonapi.Link("project", access.ProjectID),
		},
	}
}

// grantParentInOrganization resolves an organization-nested project, the way
// workspaceInOrganization does for workspaces.
func (a *API) projectGrantParent(w http.ResponseWriter, r *http.Request) (store.Project, bool) {
	projectID := param(r, "projectID")
	project, err := a.deps.Store.GetProject(r.Context(), projectID)
	if err == nil && project.OrganizationID == param(r, "organizationID") {
		return project, true
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		a.internalError(w, "get project", err)
		return store.Project{}, false
	}
	jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+projectID+"] for project")
	return store.Project{}, false
}

func (a *API) listWorkspaceAccess(w http.ResponseWriter, r *http.Request) {
	items, err := a.deps.Store.ListAccessByWorkspace(r.Context(), param(r, "workspaceID"))
	if err != nil {
		if !a.writeStoreError(w, "list workspace access", err) {
			a.internalError(w, "list workspace access", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(items))
	for _, item := range items {
		resources = append(resources, accessResource(item))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

// workspaceAccess loads the grant and checks it belongs to the workspace in the
// path; both misses answer upstream's "Unknown identifier".
func (a *API) workspaceAccess(w http.ResponseWriter, r *http.Request) (store.Access, bool) {
	id := param(r, "accessID")
	access, err := a.deps.Store.GetAccess(r.Context(), id)
	if err == nil && access.WorkspaceID == param(r, "workspaceID") {
		return access, true
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		a.internalError(w, "get workspace access", err)
		return store.Access{}, false
	}
	jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+id+"] for access")
	return store.Access{}, false
}

func (a *API) getWorkspaceAccess(w http.ResponseWriter, r *http.Request) {
	access, ok := a.workspaceAccess(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, accessResource(access), jsonapi.ParseFieldsets(r))
}

func (a *API) createWorkspaceAccess(w http.ResponseWriter, r *http.Request) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "access") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[grantAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.Name == nil {
		// Elide answers its data-integrity conflict for a missing required
		// column, as it does for collections and global variables. An empty
		// string is a value, and the column accepts it.
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This operation conflicts with data that already exists.")
		return
	}

	access, err := a.deps.Store.CreateAccess(r.Context(), store.CreateAccessParams{
		ID:              uuid.NewString(),
		Name:            *attrs.Name,
		ManageWorkspace: boolValue(attrs.ManageWorkspace, false),
		ManageState:     boolValue(attrs.ManageState, false),
		ManageJob:       boolValue(attrs.ManageJob, false),
		PlanJob:         boolValue(attrs.PlanJob, false),
		ApproveJob:      boolValue(attrs.ApproveJob, false),
		DestroyJob:      boolValue(attrs.DestroyJob, true),
		Role:            attrs.Role,
		WorkspaceID:     param(r, "workspaceID"),
	})
	if err != nil {
		if !a.writeStoreError(w, "create workspace access", err) {
			a.internalError(w, "create workspace access", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, accessResource(access), jsonapi.ParseFieldsets(r))
}

func (a *API) updateWorkspaceAccess(w http.ResponseWriter, r *http.Request) {
	access, ok := a.workspaceAccess(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "access") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[grantAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	if _, err := a.deps.Store.UpdateAccess(r.Context(), store.UpdateAccessParams{
		ID:              access.ID,
		Name:            pickStringValue(data, "name", attrs.Name, access.Name),
		ManageWorkspace: pickBool(data, "manageWorkspace", attrs.ManageWorkspace, access.ManageWorkspace),
		ManageState:     pickBool(data, "manageState", attrs.ManageState, access.ManageState),
		ManageJob:       pickBool(data, "manageJob", attrs.ManageJob, access.ManageJob),
		PlanJob:         pickBool(data, "planJob", attrs.PlanJob, access.PlanJob),
		ApproveJob:      pickBool(data, "approveJob", attrs.ApproveJob, access.ApproveJob),
		DestroyJob:      pickBool(data, "destroyJob", attrs.DestroyJob, access.DestroyJob),
		Role:            pickString(data, "role", attrs.Role, access.Role),
	}); err != nil {
		if !a.writeStoreError(w, "update workspace access", err) {
			a.internalError(w, "update workspace access", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteWorkspaceAccess(w http.ResponseWriter, r *http.Request) {
	access, ok := a.workspaceAccess(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteAccess(r.Context(), access.ID); err != nil {
		if !a.writeStoreError(w, "delete workspace access", err) {
			a.internalError(w, "delete workspace access", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listProjectAccess(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.projectGrantParent(w, r); !ok {
		return
	}
	items, err := a.deps.Store.ListProjectAccessByProject(r.Context(), param(r, "projectID"))
	if err != nil {
		if !a.writeStoreError(w, "list project access", err) {
			a.internalError(w, "list project access", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(items))
	for _, item := range items {
		resources = append(resources, projectAccessResource(item))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) projectAccess(w http.ResponseWriter, r *http.Request) (store.ProjectAccess, bool) {
	id := param(r, "projectAccessID")
	access, err := a.deps.Store.GetProjectAccess(r.Context(), id)
	if err == nil && access.ProjectID == param(r, "projectID") {
		return access, true
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		a.internalError(w, "get project access", err)
		return store.ProjectAccess{}, false
	}
	jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+id+"] for projectAccess")
	return store.ProjectAccess{}, false
}

func (a *API) getProjectAccess(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.projectGrantParent(w, r); !ok {
		return
	}
	access, ok := a.projectAccess(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, projectAccessResource(access), jsonapi.ParseFieldsets(r))
}

func (a *API) createProjectAccess(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.projectGrantParent(w, r); !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "project_access") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[grantAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	access, err := a.deps.Store.CreateProjectAccess(r.Context(), store.CreateProjectAccessParams{
		ID:              uuid.NewString(),
		Name:            attrs.Name,
		ManageWorkspace: boolValue(attrs.ManageWorkspace, false),
		ManageState:     boolValue(attrs.ManageState, false),
		ManageJob:       boolValue(attrs.ManageJob, false),
		PlanJob:         boolValue(attrs.PlanJob, false),
		ApproveJob:      boolValue(attrs.ApproveJob, false),
		DestroyJob:      boolValue(attrs.DestroyJob, true),
		Role:            attrs.Role,
		ProjectID:       param(r, "projectID"),
	})
	if err != nil {
		if !a.writeStoreError(w, "create project access", err) {
			a.internalError(w, "create project access", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, projectAccessResource(access), jsonapi.ParseFieldsets(r))
}

func (a *API) updateProjectAccess(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.projectGrantParent(w, r); !ok {
		return
	}
	access, ok := a.projectAccess(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "project_access") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[grantAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	if _, err := a.deps.Store.UpdateProjectAccess(r.Context(), store.UpdateProjectAccessParams{
		ID:              access.ID,
		Name:            pickString(data, "name", attrs.Name, access.Name),
		ManageWorkspace: pickBool(data, "manageWorkspace", attrs.ManageWorkspace, access.ManageWorkspace),
		ManageState:     pickBool(data, "manageState", attrs.ManageState, access.ManageState),
		ManageJob:       pickBool(data, "manageJob", attrs.ManageJob, access.ManageJob),
		PlanJob:         pickBool(data, "planJob", attrs.PlanJob, access.PlanJob),
		ApproveJob:      pickBool(data, "approveJob", attrs.ApproveJob, access.ApproveJob),
		DestroyJob:      pickBool(data, "destroyJob", attrs.DestroyJob, access.DestroyJob),
		Role:            pickString(data, "role", attrs.Role, access.Role),
	}); err != nil {
		if !a.writeStoreError(w, "update project access", err) {
			a.internalError(w, "update project access", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteProjectAccess(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.projectGrantParent(w, r); !ok {
		return
	}
	access, ok := a.projectAccess(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteProjectAccess(r.Context(), access.ID); err != nil {
		if !a.writeStoreError(w, "delete project access", err) {
			a.internalError(w, "delete project access", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
