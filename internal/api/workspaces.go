package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"

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
	jsonapi.WriteCollection(w, http.StatusOK, resources, total, q.Totals)
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
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), false)
}

func (a *API) getWorkspace(w http.ResponseWriter, r *http.Request) {
	ws, err := a.deps.Store.GetWorkspace(r.Context(), param(r, "workspaceID"))
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, workspaceResource(ws))
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
		if err := a.validateProjectInOrganization(r.Context(), projectID, orgID); err != nil {
			if !a.writeStoreError(w, "validate project", err) {
				jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", "The project does not belong to this organization.")
			}
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

	jsonapi.WriteResource(w, http.StatusCreated, workspaceResource(ws))
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
			if err := a.validateProjectInOrganization(r.Context(), relID, ws.OrganizationID); err != nil {
				if !a.writeStoreError(w, "validate project", err) {
					jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", "The project does not belong to this organization.")
				}
				return
			}
			projectID = &relID
		}
	}

	updated, err := a.deps.Store.UpdateWorkspace(r.Context(), store.UpdateWorkspaceParams{
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
	})
	if err != nil {
		if !a.writeStoreError(w, "update workspace", err) {
			a.internalError(w, "update workspace", err)
		}
		return
	}

	jsonapi.WriteResource(w, http.StatusOK, workspaceResource(updated))
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

func (a *API) validateProjectInOrganization(ctx context.Context, projectID, orgID string) error {
	project, err := a.deps.Store.GetProject(ctx, projectID)
	if err != nil {
		return err
	}
	if project.OrganizationID != orgID {
		return errProjectOutsideOrganization
	}
	return nil
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
