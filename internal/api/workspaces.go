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
	ModuleSshKey                    *string `json:"moduleSshKey"`
	SharedIds                       *string `json:"sharedIds"`
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

	// The workspace page reads its collection references and then asks each
	// one for its collection items (workspaceDataUtils.ts, include.REFERENCE).
	if want["reference"] {
		references, err := a.deps.Store.ListReferencesByWorkspace(r.Context(), ws.ID)
		if err != nil {
			if !a.writeStoreError(w, "list references", err) {
				a.internalError(w, "list references", err)
			}
			return
		}
		ids := make([]string, 0, len(references))
		for _, reference := range references {
			ids = append(ids, reference.ID)
			included = append(included, referenceResource(reference))
		}
		resource.Relationships["reference"] = jsonapi.LinkList("reference", ids...)
	}

	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: resource, Included: included}, jsonapi.ParseFieldsets(r))
}

func (a *API) createWorkspace(w http.ResponseWriter, r *http.Request) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "workspace") {
		return
	}

	ws, ok := a.createWorkspaceFromData(w, r, param(r, "organizationID"), data)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, workspaceResource(ws), jsonapi.ParseFieldsets(r))
}

// createWorkspaceFromData is shared by POST /workspace and the atomic
// operations endpoint, which receives the same resource object without a
// wrapping request body.
func (a *API) createWorkspaceFromData(w http.ResponseWriter, r *http.Request, orgID string, data *jsonapi.RequestData) (store.Workspace, bool) {
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return store.Workspace{}, false
	}

	attrs, err := jsonapi.DecodeAttributes[workspaceAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return store.Workspace{}, false
	}
	if attrs.Name == nil || *attrs.Name == "" {
		jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Missing Attribute", "name is required.", "name")
		return store.Workspace{}, false
	}

	projectID, present, err := jsonapi.RelationshipID(data, "project")
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", err.Error())
		return store.Workspace{}, false
	}
	if present && projectID != "" {
		if !a.checkProjectRelationship(w, r, projectID, orgID) {
			return store.Workspace{}, false
		}
	} else {
		projectID = ""
	}

	links, ok := workspaceLinks(w, data, nil)
	if !ok {
		return store.Workspace{}, false
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
		ModuleSshKey:                    attrs.ModuleSshKey,
		SharedIds:                       attrs.SharedIds,
		VcsID:                           links["vcs"],
		SshID:                           links["ssh"],
		AgentID:                         links["agent"],
	})
	if err != nil {
		if !a.writeStoreError(w, "create workspace", err) {
			a.internalError(w, "create workspace", err)
		}
		return store.Workspace{}, false
	}
	return ws, true
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

	if !a.applyWorkspaceUpdate(w, r, ws, data) {
		return
	}

	// Upstream answers a workspace PATCH with 204 and no body, for attributes
	// and relationships alike; the UI refetches what it needs.
	w.WriteHeader(http.StatusNoContent)
}

// applyWorkspaceUpdate is shared by PATCH and the atomic operations endpoint;
// it writes the error itself and reports whether the row was updated.
func (a *API) applyWorkspaceUpdate(w http.ResponseWriter, r *http.Request, ws store.Workspace, data *jsonapi.RequestData) bool {
	attrs, err := jsonapi.DecodeAttributes[workspaceAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return false
	}

	name := ws.Name
	if jsonapi.HasAttribute(data, "name") {
		if attrs.Name == nil || *attrs.Name == "" {
			jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Invalid Attribute", "name must not be empty.", "name")
			return false
		}
		name = *attrs.Name
	}

	params := workspaceUpdateParams(ws)
	params.Name = name
	params.Description = pickString(data, "description", attrs.Description, ws.Description)
	params.Source = pickString(data, "source", attrs.Source, ws.Source)
	params.Branch = pickString(data, "branch", attrs.Branch, ws.Branch)
	params.Folder = pickString(data, "folder", attrs.Folder, ws.Folder)
	params.TerraformVersion = pickString(data, "terraformVersion", attrs.TerraformVersion, ws.TerraformVersion)
	params.IacType = pickStringValue(data, "iacType", attrs.IacType, ws.IacType)
	params.ExecutionMode = pickString(data, "executionMode", attrs.ExecutionMode, ws.ExecutionMode)
	params.AllowRemoteApply = pickBool(data, "allowRemoteApply", attrs.AllowRemoteApply, ws.AllowRemoteApply)
	params.GlobalRemoteState = pickBool(data, "globalRemoteState", attrs.GlobalRemoteState, ws.GlobalRemoteState)
	params.DefaultTemplate = pickString(data, "defaultTemplate", attrs.DefaultTemplate, ws.DefaultTemplate)
	params.RestrictApplyToConfiguredBranch = pickBool(data, "restrictApplyToConfiguredBranch", attrs.RestrictApplyToConfiguredBranch, ws.RestrictApplyToConfiguredBranch)
	params.RequireApproveToApply = pickBool(data, "requireApproveToApply", attrs.RequireApproveToApply, ws.RequireApproveToApply)
	params.ModuleSshKey = pickString(data, "moduleSshKey", attrs.ModuleSshKey, ws.ModuleSshKey)
	params.SharedIds = pickString(data, "sharedIds", attrs.SharedIds, ws.SharedIds)

	if relID, present, relErr := jsonapi.RelationshipID(data, "project"); relErr != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", relErr.Error())
		return false
	} else if present {
		if relID == "" {
			params.ProjectID = nil
		} else {
			if !a.checkProjectRelationship(w, r, relID, ws.OrganizationID) {
				return false
			}
			params.ProjectID = &relID
		}
	}

	links, ok := workspaceLinks(w, data, map[string]*string{
		"vcs": ws.VcsID, "ssh": ws.SshID, "agent": ws.AgentID,
	})
	if !ok {
		return false
	}
	params.VcsID, params.SshID, params.AgentID = links["vcs"], links["ssh"], links["agent"]

	if _, err := a.deps.Store.UpdateWorkspace(r.Context(), params); err != nil {
		if !a.writeStoreError(w, "update workspace", err) {
			a.internalError(w, "update workspace", err)
		}
		return false
	}
	return true
}

// workspaceUpdateParams seeds an update with the stored row so an attribute the
// request did not send keeps its value.
func workspaceUpdateParams(ws store.Workspace) store.UpdateWorkspaceParams {
	return store.UpdateWorkspaceParams{
		ID:                              ws.ID,
		Name:                            ws.Name,
		Description:                     ws.Description,
		Source:                          ws.Source,
		Branch:                          ws.Branch,
		Folder:                          ws.Folder,
		TerraformVersion:                ws.TerraformVersion,
		IacType:                         ws.IacType,
		ExecutionMode:                   ws.ExecutionMode,
		AllowRemoteApply:                ws.AllowRemoteApply,
		GlobalRemoteState:               ws.GlobalRemoteState,
		DefaultTemplate:                 ws.DefaultTemplate,
		RestrictApplyToConfiguredBranch: ws.RestrictApplyToConfiguredBranch,
		RequireApproveToApply:           ws.RequireApproveToApply,
		ModuleSshKey:                    ws.ModuleSshKey,
		SharedIds:                       ws.SharedIds,
		ProjectID:                       ws.ProjectID,
		VcsID:                           ws.VcsID,
		SshID:                           ws.SshID,
		AgentID:                         ws.AgentID,
	}
}

// workspaceLinks reads the vcs, ssh and agent relationships: absent keeps the
// current id, null clears it. Upstream resolves these ids against its tables;
// those tables do not exist here yet (#48, #120), so an id is stored as sent.
func workspaceLinks(w http.ResponseWriter, data *jsonapi.RequestData, current map[string]*string) (map[string]*string, bool) {
	links := map[string]*string{"vcs": nil, "ssh": nil, "agent": nil}
	for name, id := range current {
		links[name] = id
	}
	for _, name := range []string{"vcs", "ssh", "agent"} {
		id, present, err := jsonapi.RelationshipID(data, name)
		if err != nil {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", err.Error())
			return nil, false
		}
		if !present {
			continue
		}
		if id == "" {
			links[name] = nil
			continue
		}
		value := id
		links[name] = &value
	}
	return links, true
}

// updateWorkspaceRelationship is the PATCH .../relationships/{relation} route
// the UI uses for the project, VCS, SSH and agent pickers.
func (a *API) updateWorkspaceRelationship(w http.ResponseWriter, r *http.Request) {
	relation := param(r, "relation")
	switch relation {
	case "project", "vcs", "ssh", "agent":
	default:
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown collection "+relation)
		return
	}

	ws, err := a.deps.Store.GetWorkspace(r.Context(), param(r, "workspaceID"))
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}

	typ, id, sent, err := jsonapi.RelationshipIdentifier(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", err.Error())
		return
	}
	if !sent {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", "data is required")
		return
	}
	if id != "" && typ != relation {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", "data.type must be "+relation)
		return
	}

	params := workspaceUpdateParams(ws)
	switch relation {
	case "project":
		if id == "" {
			params.ProjectID = nil
		} else {
			if !a.checkProjectRelationship(w, r, id, ws.OrganizationID) {
				return
			}
			params.ProjectID = &id
		}
	case "vcs", "ssh", "agent":
		var value *string
		if id != "" {
			value = &id
		}
		switch relation {
		case "vcs":
			params.VcsID = value
		case "ssh":
			params.SshID = value
		case "agent":
			params.AgentID = value
		}
	}

	if _, err := a.deps.Store.UpdateWorkspace(r.Context(), params); err != nil {
		if !a.writeStoreError(w, "update workspace", err) {
			a.internalError(w, "update workspace", err)
		}
		return
	}
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
	if ws.VcsID != nil && *ws.VcsID != "" {
		relationships["vcs"] = jsonapi.Link("vcs", *ws.VcsID)
	}
	if ws.SshID != nil && *ws.SshID != "" {
		relationships["ssh"] = jsonapi.Link("ssh", *ws.SshID)
	}
	if ws.AgentID != nil && *ws.AgentID != "" {
		relationships["agent"] = jsonapi.Link("agent", *ws.AgentID)
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
			"moduleSshKey":                    strValue(ws.ModuleSshKey),
			"sharedIds":                       strValue(ws.SharedIds),
			"lastJobDate":                     timestamp(ws.LastJobDate),
			"createdBy":                       ws.CreatedBy,
			"createdDate":                     timestamp(ws.CreatedDate),
			"updatedBy":                       ws.UpdatedBy,
			"updatedDate":                     timestamp(ws.UpdatedDate),
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
