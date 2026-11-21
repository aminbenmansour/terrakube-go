package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type workspaceAttrs struct {
	Name                            *string `json:"name"`
	Description                     *string `json:"description"`
	Source                          *string `json:"source"`
	Branch                          *string `json:"branch"`
	Folder                          *string `json:"folder"`
	TerraformVersion                *string `json:"terraformVersion"`
	IacType                         *string `json:"iacType"`
	ExecutionMode                   *string `json:"executionMode"`
	AllowRemoteApply                *bool   `json:"allowRemoteApply"`
	GlobalRemoteState               *bool   `json:"globalRemoteState"`
	DefaultTemplate                 *string `json:"defaultTemplate"`
	RestrictApplyToConfiguredBranch *bool   `json:"restrictApplyToConfiguredBranch"`
	RequireApproveToApply           *bool   `json:"requireApproveToApply"`
}

func (a *API) listWorkspaces(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}

	q := parseList(r, "workspace")
	items, err := a.deps.Store.ListWorkspacesByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list workspaces", err) {
			a.internalError(w, "list workspaces", err)
		}
		return
	}

	page, total := filterSortPage(items, q, func(ws store.Workspace) string { return ws.Name })
	resources := make([]jsonapi.Resource, 0, len(page))
	for _, ws := range page {
		resources = append(resources, workspaceResource(ws))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, total, q.page(), jsonapi.ParseFieldsets(r))
}

func (a *API) listProjectWorkspaces(w http.ResponseWriter, r *http.Request) {
	projectID := param(r, "projectID")
	if _, err := a.deps.Store.GetProject(r.Context(), projectID); err != nil {
		if !a.writeStoreError(w, "get project", err) {
			a.internalError(w, "get project", err)
		}
		return
	}

	items, err := a.deps.Store.ListWorkspacesByProject(r.Context(), &projectID)
	if err != nil {
		if !a.writeStoreError(w, "list project workspaces", err) {
			a.internalError(w, "list project workspaces", err)
		}
		return
	}

	resources := make([]jsonapi.Resource, 0, len(items))
	for _, ws := range items {
		resources = append(resources, workspaceResource(ws))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getWorkspace(w http.ResponseWriter, r *http.Request) {
	ws, err := a.deps.Store.GetWorkspace(r.Context(), param(r, "workspaceID"))
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}

	// The UI loads every workspace tab from one request with include=job,
	// variable,... and reads the list out of `included` (Details.tsx). We serve
	// the relationships that have tables; the rest are ignored rather than
	// failing the whole request.
	resource := workspaceResource(ws)
	var included []jsonapi.Resource
	want := parseInclude(r)

	if want["variable"] {
		variables, err := a.deps.Store.ListVariablesByWorkspace(r.Context(), ws.ID)
		if err != nil {
			if !a.writeStoreError(w, "list variables", err) {
				a.internalError(w, "list variables", err)
			}
			return
		}
		ids := make([]string, 0, len(variables))
		for _, variable := range variables {
			ids = append(ids, variable.ID)
			included = append(included, variableResource(variable))
		}
		resource.Relationships["variable"] = jsonapi.LinkList("variable", ids...)
	}

	if want["organization"] {
		org, err := a.deps.Store.GetOrganization(r.Context(), ws.OrganizationID)
		if err != nil {
			if !a.writeStoreError(w, "get organization", err) {
				a.internalError(w, "get organization", err)
			}
			return
		}
		included = append(included, organizationResource(org))
	}

	if want["project"] && ws.ProjectID != nil && *ws.ProjectID != "" {
		project, err := a.deps.Store.GetProject(r.Context(), *ws.ProjectID)
		if err != nil {
			if !a.writeStoreError(w, "get project", err) {
				a.internalError(w, "get project", err)
			}
			return
		}
		included = append(included, projectResource(project))
	}

	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: resource, Included: included}, jsonapi.ParseFieldsets(r))
}

func (a *API) createWorkspace(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "workspace") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[workspaceAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.Name == nil || *attrs.Name == "" {
		jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Missing Attribute", "name is required.", "name")
		return
	}

	projectID, present, err := jsonapi.RelationshipID(data, "project")
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", err.Error())
		return
	}
	if present && projectID != "" {
		if !a.checkProjectRelationship(w, r, projectID, orgID) {
			return
		}
	} else {
		projectID = ""
	}

	iacType := "terraform"
	if attrs.IacType != nil && *attrs.IacType != "" {
		iacType = *attrs.IacType
	}

	ws, err := a.deps.Store.CreateWorkspace(r.Context(), store.CreateWorkspaceParams{
		ID:                              uuid.NewString(),
		Name:                            *attrs.Name,
		Description:                     attrs.Description,
		Source:                          attrs.Source,
		Branch:                          attrs.Branch,
		Folder:                          attrs.Folder,
		TerraformVersion:                attrs.TerraformVersion,
		IacType:                         iacType,
		ExecutionMode:                   attrs.ExecutionMode,
		OrganizationID:                  orgID,
		ProjectID:                       strPtr(projectID),
		AllowRemoteApply:                boolValue(attrs.AllowRemoteApply, false),
		GlobalRemoteState:               boolValue(attrs.GlobalRemoteState, true),
		DefaultTemplate:                 attrs.DefaultTemplate,
		RestrictApplyToConfiguredBranch: boolValue(attrs.RestrictApplyToConfiguredBranch, false),
		RequireApproveToApply:           boolValue(attrs.RequireApproveToApply, false),
		CreatedBy:                       nil,
	})
	if err != nil {
		if !a.writeStoreError(w, "create workspace", err) {
			a.internalError(w, "create workspace", err)
		}
		return
	}

	jsonapi.WriteResource(w, http.StatusCreated, workspaceResource(ws), jsonapi.ParseFieldsets(r))
}

func (a *API) updateWorkspace(w http.ResponseWriter, r *http.Request) {
	ws, err := a.deps.Store.GetWorkspace(r.Context(), param(r, "workspaceID"))
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "workspace") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[workspaceAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	name := ws.Name
	if jsonapi.HasAttribute(data, "name") {
		if attrs.Name == nil || *attrs.Name == "" {
			jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Invalid Attribute", "name must not be empty.", "name")
			return
		}
		name = *attrs.Name
	}

	projectID := ws.ProjectID
	if relID, present, relErr := jsonapi.RelationshipID(data, "project"); relErr != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", relErr.Error())
		return
	} else if present {
		if relID == "" {
			projectID = nil
		} else {
			if !a.checkProjectRelationship(w, r, relID, ws.OrganizationID) {
				return
			}
			projectID = &relID
		}
	}

	if _, err := a.deps.Store.UpdateWorkspace(r.Context(), store.UpdateWorkspaceParams{
		ID:                              ws.ID,
		Name:                            name,
		Description:                     pickString(data, "description", attrs.Description, ws.Description),
		Source:                          pickString(data, "source", attrs.Source, ws.Source),
		Branch:                          pickString(data, "branch", attrs.Branch, ws.Branch),
		Folder:                          pickString(data, "folder", attrs.Folder, ws.Folder),
		TerraformVersion:                pickString(data, "terraformVersion", attrs.TerraformVersion, ws.TerraformVersion),
		IacType:                         pickStringValue(data, "iacType", attrs.IacType, ws.IacType),
		ExecutionMode:                   pickString(data, "executionMode", attrs.ExecutionMode, ws.ExecutionMode),
		AllowRemoteApply:                pickBool(data, "allowRemoteApply", attrs.AllowRemoteApply, ws.AllowRemoteApply),
		GlobalRemoteState:               pickBool(data, "globalRemoteState", attrs.GlobalRemoteState, ws.GlobalRemoteState),
		DefaultTemplate:                 pickString(data, "defaultTemplate", attrs.DefaultTemplate, ws.DefaultTemplate),
		RestrictApplyToConfiguredBranch: pickBool(data, "restrictApplyToConfiguredBranch", attrs.RestrictApplyToConfiguredBranch, ws.RestrictApplyToConfiguredBranch),
		RequireApproveToApply:           pickBool(data, "requireApproveToApply", attrs.RequireApproveToApply, ws.RequireApproveToApply),
		ProjectID:                       projectID,
	}); err != nil {
		if !a.writeStoreError(w, "update workspace", err) {
			a.internalError(w, "update workspace", err)
		}
		return
	}

	// Upstream answers a workspace PATCH with 204 and no body, for attributes
	// and relationships alike; the UI refetches what it needs.
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteWorkspace(w http.ResponseWriter, r *http.Request) {
	id := param(r, "workspaceID")
	if _, err := a.deps.Store.GetWorkspace(r.Context(), id); err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}
	if err := a.deps.Store.SoftDeleteWorkspace(r.Context(), id); err != nil {
		if !a.writeStoreError(w, "delete workspace", err) {
			a.internalError(w, "delete workspace", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// checkProjectRelationship resolves the project an incoming workspace names.
// Upstream distinguishes an id that does not exist (404) from a project in
// another organization (403); both were a 400 here, and a client cannot tell
// which mistake it made from that.
func (a *API) checkProjectRelationship(w http.ResponseWriter, r *http.Request, projectID, orgID string) bool {
	project, err := a.deps.Store.GetProject(r.Context(), projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier "+projectID+" for project")
			return false
		}
		if !a.writeStoreError(w, "get project", err) {
			a.internalError(w, "get project", err)
		}
		return false
	}
	if project.OrganizationID != orgID {
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "Project does not belong to the workspace organization")
		return false
	}
	return true
}

// workspaceInOrganization guards the organization-scoped workspace routes. A
// workspace of another organization is answered 404, not 403: upstream does not
// confirm that an id exists outside the caller's scope.
func (a *API) workspaceInOrganization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := a.deps.Store.GetWorkspace(r.Context(), param(r, "workspaceID"))
		if err != nil {
			if !a.writeStoreError(w, "get workspace", err) {
				a.internalError(w, "get workspace", err)
			}
			return
		}
		if ws.OrganizationID != param(r, "organizationID") {
			jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func workspaceResource(ws store.Workspace) jsonapi.Resource {
	relationships := map[string]jsonapi.Relationship{
		"organization": jsonapi.Link("organization", ws.OrganizationID),
	}
	if ws.ProjectID != nil && *ws.ProjectID != "" {
		relationships["project"] = jsonapi.Link("project", *ws.ProjectID)
	}

	return jsonapi.Resource{
		Type: "workspace",
		ID:   ws.ID,
		Attributes: map[string]any{
			"name":                            ws.Name,
			"description":                     ws.Description,
			"source":                          strValue(ws.Source),
			"branch":                          strValue(ws.Branch),
			"folder":                          strValue(ws.Folder),
			"terraformVersion":                strValue(ws.TerraformVersion),
			"iacType":                         ws.IacType,
			"executionMode":                   ws.ExecutionMode,
			"locked":                          ws.Locked,
			"deleted":                         ws.Deleted,
			"allowRemoteApply":                ws.AllowRemoteApply,
			"globalRemoteState":               ws.GlobalRemoteState,
			"defaultTemplate":                 ws.DefaultTemplate,
			"restrictApplyToConfiguredBranch": ws.RestrictApplyToConfiguredBranch,
			"requireApproveToApply":           ws.RequireApproveToApply,
			"lastJobStatus":                   ws.LastJobStatus,
			"policyComplianceStatus":          ws.PolicyComplianceStatus,
			"lockDescription":                 ws.LockDescription,
		},
		Relationships: relationships,
	}
}

func pickString(data *jsonapi.RequestData, name string, incoming, existing *string) *string {
	if jsonapi.HasAttribute(data, name) {
		return incoming
	}
	return existing
}

func pickStringValue(data *jsonapi.RequestData, name string, incoming *string, existing string) string {
	if jsonapi.HasAttribute(data, name) && incoming != nil && *incoming != "" {
		return *incoming
	}
	return existing
}

func pickBool(data *jsonapi.RequestData, name string, incoming *bool, existing bool) bool {
	if jsonapi.HasAttribute(data, name) {
		return boolValue(incoming, existing)
	}
	return existing
}
