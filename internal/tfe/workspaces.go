package tfe

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

var exactTerraformVersion = regexp.MustCompile(`^\d+(\.\d+){0,2}(-[0-9A-Za-z.-]+)?$`)

type workspaceAttrs struct {
	Name             *string `json:"name"`
	TerraformVersion *string `json:"terraform-version"`
	ExecutionMode    *string `json:"execution-mode"`
}

func (a *API) getWorkspace(w http.ResponseWriter, r *http.Request) {
	org, ok := a.organizationByName(w, r)
	if !ok {
		return
	}
	ws, err := a.deps.Store.GetWorkspaceByName(r.Context(), store.GetWorkspaceByNameParams{
		OrganizationID: org.ID,
		Name:           chiURLParam(r, "workspaceName"),
	})
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}
	if !a.authorizeScope(w, r, org.ID, ws.ID) {
		return
	}
	tagNames, err := a.workspaceTagNames(r, org.ID, ws.ID)
	if err != nil {
		a.internalError(w, "load workspace tags", err)
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, workspaceResource(ws, tagNames))
}

func (a *API) listWorkspaces(w http.ResponseWriter, r *http.Request) {
	org, ok := a.organizationByName(w, r)
	if !ok {
		return
	}

	query := parseWorkspaceQuery(r)
	if !query.hasAnyCondition() {
		writeWorkspaceList(w, nil, query, 0)
		return
	}

	items, err := a.deps.Store.ListWorkspacesByOrganization(r.Context(), org.ID)
	if err != nil {
		if !a.writeStoreError(w, "list workspaces", err) {
			a.internalError(w, "list workspaces", err)
		}
		return
	}

	tagNamesByID, err := a.organizationTagNames(r, org.ID)
	if err != nil {
		a.internalError(w, "load organization tags", err)
		return
	}

	matches := make([]workspaceWithTags, 0, len(items))
	for _, ws := range items {
		if query.searchName != "" && !strings.HasPrefix(ws.Name, query.searchName) {
			continue
		}
		if query.projectID != "" && (ws.ProjectID == nil || *ws.ProjectID != query.projectID) {
			continue
		}
		workspaceTags, err := a.deps.Store.ListWorkspaceTagsByWorkspace(r.Context(), ws.ID)
		if err != nil {
			a.internalError(w, "load workspace tags", err)
			return
		}
		if !matchesTagFilters(workspaceTags, tagNamesByID, query) {
			continue
		}
		matches = append(matches, workspaceWithTags{workspace: ws, tags: workspaceTagNames(workspaceTags, tagNamesByID)})
	}

	sort.Slice(matches, func(i, j int) bool {
		return strings.ToLower(matches[i].workspace.Name) < strings.ToLower(matches[j].workspace.Name)
	})

	total := len(matches)
	start := query.offset()
	if start > total {
		start = total
	}
	end := start + query.pageSize
	if end > total {
		end = total
	}

	page := matches[start:end]
	resources := make([]jsonapi.Resource, 0, len(page))
	for _, item := range page {
		resources = append(resources, workspaceResource(item.workspace, item.tags))
	}
	writeWorkspaceList(w, resources, query, total)
}

type workspaceWithTags struct {
	workspace store.Workspace
	tags      []string
}

func (a *API) createWorkspace(w http.ResponseWriter, r *http.Request) {
	org, ok := a.organizationByName(w, r)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if data.Type != "workspaces" {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Type", "data.type must be workspaces.")
		return
	}

	attrs, err := jsonapi.DecodeAttributes[workspaceAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.Name == nil || strings.TrimSpace(*attrs.Name) == "" {
		jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Missing Attribute", "name is required.", "name")
		return
	}
	name := strings.TrimSpace(*attrs.Name)

	// go-tfe treats workspace creation as idempotent: an existing workspace is
	// returned instead of creating a duplicate.
	existing, err := a.deps.Store.GetWorkspaceByName(r.Context(), store.GetWorkspaceByNameParams{
		OrganizationID: org.ID,
		Name:           name,
	})
	if err == nil {
		tagNames, tagErr := a.workspaceTagNames(r, org.ID, existing.ID)
		if tagErr != nil {
			a.internalError(w, "load workspace tags", tagErr)
			return
		}
		jsonapi.WriteResource(w, http.StatusCreated, workspaceResource(existing, tagNames))
		return
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		a.internalError(w, "look up workspace", err)
		return
	}

	tagBindings, err := workspaceTagBindings(data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusUnprocessableEntity, "Invalid Tag Binding", err.Error())
		return
	}

	var projectID *string
	if id, present, relErr := jsonapi.RelationshipID(data, "project"); relErr != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", relErr.Error())
		return
	} else if present && id != "" {
		project, err := a.deps.Store.GetProject(r.Context(), id)
		if err != nil {
			if !a.writeStoreError(w, "get project", err) {
				a.internalError(w, "get project", err)
			}
			return
		}
		if project.OrganizationID != org.ID {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", "The project does not belong to this organization.")
			return
		}
		projectID = &id
	}

	terraformVersion := "1.4.6"
	if attrs.TerraformVersion != nil && strings.TrimSpace(*attrs.TerraformVersion) != "" {
		terraformVersion = strings.TrimSpace(*attrs.TerraformVersion)
	}

	executionMode := org.ExecutionMode
	if attrs.ExecutionMode != nil {
		requested := strings.ToLower(strings.TrimSpace(*attrs.ExecutionMode))
		if requested == "remote" || requested == "local" {
			executionMode = requested
		}
	}
	if executionMode == "" {
		executionMode = "remote"
	}

	source := "empty"
	branch := "remote-content"
	ws, err := a.deps.Store.CreateWorkspace(r.Context(), store.CreateWorkspaceParams{
		ID:                uuid.NewString(),
		Name:              name,
		Source:            &source,
		Branch:            &branch,
		TerraformVersion:  &terraformVersion,
		IacType:           "terraform",
		ExecutionMode:     &executionMode,
		OrganizationID:    org.ID,
		ProjectID:         projectID,
		GlobalRemoteState: true,
	})
	if err != nil {
		if !a.writeStoreError(w, "create workspace", err) {
			a.internalError(w, "create workspace", err)
		}
		return
	}

	if len(tagBindings) > 0 {
		if err := a.applyTagBindings(r, ws, tagBindings, true); err != nil {
			a.internalError(w, "apply workspace tags", err)
			return
		}
	}
	tagNames, err := a.workspaceTagNames(r, org.ID, ws.ID)
	if err != nil {
		a.internalError(w, "load workspace tags", err)
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, workspaceResource(ws, tagNames))
}

func (a *API) updateWorkspace(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if data.Type != "workspaces" {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Type", "data.type must be workspaces.")
		return
	}
	attrs, err := jsonapi.DecodeAttributes[workspaceAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	terraformVersion := ws.TerraformVersion
	if attrs.TerraformVersion != nil {
		value := strings.TrimSpace(*attrs.TerraformVersion)
		terraformVersion = &value
	}

	updated, err := a.deps.Store.UpdateWorkspace(r.Context(), store.UpdateWorkspaceParams{
		ID:                              ws.ID,
		Name:                            ws.Name,
		Description:                     ws.Description,
		Source:                          ws.Source,
		Branch:                          ws.Branch,
		Folder:                          ws.Folder,
		TerraformVersion:                terraformVersion,
		IacType:                         ws.IacType,
		ExecutionMode:                   ws.ExecutionMode,
		AllowRemoteApply:                ws.AllowRemoteApply,
		GlobalRemoteState:               ws.GlobalRemoteState,
		DefaultTemplate:                 ws.DefaultTemplate,
		RestrictApplyToConfiguredBranch: ws.RestrictApplyToConfiguredBranch,
		RequireApproveToApply:           ws.RequireApproveToApply,
		ProjectID:                       ws.ProjectID,
	})
	if err != nil {
		if !a.writeStoreError(w, "update workspace", err) {
			a.internalError(w, "update workspace", err)
		}
		return
	}

	tagNames, err := a.workspaceTagNames(r, updated.OrganizationID, updated.ID)
	if err != nil {
		a.internalError(w, "load workspace tags", err)
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, workspaceResource(updated, tagNames))
}

func (a *API) lockWorkspace(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}
	if ws.Locked {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "Unable to lock workspace. The workspace is already locked.")
		return
	}

	locked, err := a.deps.Store.SetWorkspaceLock(r.Context(), store.SetWorkspaceLockParams{
		ID:              ws.ID,
		Locked:          true,
		LockDescription: ws.LockDescription,
	})
	if err != nil {
		if !a.writeStoreError(w, "lock workspace", err) {
			a.internalError(w, "lock workspace", err)
		}
		return
	}
	tagNames, err := a.workspaceTagNames(r, locked.OrganizationID, locked.ID)
	if err != nil {
		a.internalError(w, "load workspace tags", err)
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, workspaceResource(locked, tagNames))
}

func (a *API) unlockWorkspace(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}
	unlocked, err := a.deps.Store.SetWorkspaceLock(r.Context(), store.SetWorkspaceLockParams{
		ID:              ws.ID,
		Locked:          false,
		LockDescription: nil,
	})
	if err != nil {
		if !a.writeStoreError(w, "unlock workspace", err) {
			a.internalError(w, "unlock workspace", err)
		}
		return
	}
	tagNames, err := a.workspaceTagNames(r, unlocked.OrganizationID, unlocked.ID)
	if err != nil {
		a.internalError(w, "load workspace tags", err)
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, workspaceResource(unlocked, tagNames))
}

func workspaceResource(ws store.Workspace, tagNames []string) jsonapi.Resource {
	if tagNames == nil {
		tagNames = []string{}
	}
	attributes := map[string]any{
		"name":                ws.Name,
		"terraform-version":   exposedTerraformVersion(ws.TerraformVersion),
		"locked":              ws.Locked,
		"auto-apply":          false,
		"execution-mode":      ws.ExecutionMode,
		"global-remote-state": ws.GlobalRemoteState,
		"tag-names":           tagNames,
		"permissions":         workspacePermissions(),
	}
	if ws.Folder != nil {
		folder := strings.TrimSpace(*ws.Folder)
		if folder != "" && strings.Split(folder, ",")[0] != "/" {
			attributes["working-directory"] = strings.Split(folder, ",")[0]
		}
	}

	relationships := map[string]jsonapi.Relationship{}
	if ws.ProjectID != nil && *ws.ProjectID != "" {
		relationships["project"] = jsonapi.Link("projects", *ws.ProjectID)
	}
	if !ws.GlobalRemoteState {
		relationships["remote-state-consumers"] = jsonapi.Relationship{
			Links: map[string]any{
				"related": fmt.Sprintf("/remote/tfe/v2/workspaces/%s/relationships/remote-state-consumers", ws.ID),
			},
		}
	}

	resource := jsonapi.Resource{
		Type:       "workspaces",
		ID:         ws.ID,
		Attributes: attributes,
	}
	if len(relationships) > 0 {
		resource.Relationships = relationships
	}
	return resource
}

func exposedTerraformVersion(version *string) any {
	if version == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*version)
	if trimmed == "" {
		return ""
	}
	if exactTerraformVersion.MatchString(trimmed) {
		return trimmed
	}
	return "latest"
}

func workspacePermissions() map[string]bool {
	return map[string]bool{
		"can-create-state-versions":  true,
		"can-destroy":                true,
		"can-force-unlock":           true,
		"can-lock":                   true,
		"can-manage-run-tasks":       true,
		"can-manage-tags":            true,
		"can-queue-apply":            true,
		"can-queue-destroy":          true,
		"can-queue-run":              true,
		"can-read-settings":          true,
		"can-read-state-versions":    true,
		"can-read-variable":          true,
		"can-unlock":                 true,
		"can-update":                 true,
		"can-update-variable":        true,
		"can-read-assessment-result": true,
		"can-force-delete":           true,
	}
}

func writeWorkspaceList(w http.ResponseWriter, resources []jsonapi.Resource, query workspaceQuery, total int) {
	if resources == nil {
		resources = []jsonapi.Resource{}
	}
	totalPages := 1
	if total > 0 {
		totalPages = (total + query.pageSize - 1) / query.pageSize
	}
	var prevPage any
	if query.pageNumber > 1 {
		prevPage = query.pageNumber - 1
	}
	var nextPage any
	if query.pageNumber < totalPages {
		nextPage = query.pageNumber + 1
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{
		Data: resources,
		Meta: map[string]any{
			"pagination": map[string]any{
				"current-page": query.pageNumber,
				"page-size":    query.pageSize,
				"prev-page":    prevPage,
				"next-page":    nextPage,
				"total-pages":  totalPages,
				"total-count":  total,
			},
		},
	})
}

func (a *API) organizationTagNames(r *http.Request, organizationID string) (map[string]string, error) {
	tags, err := a.deps.Store.ListTagsByOrganization(r.Context(), organizationID)
	if err != nil {
		return nil, err
	}
	names := make(map[string]string, len(tags))
	for _, tag := range tags {
		if tag.Name != nil {
			names[tag.ID] = *tag.Name
		}
	}
	return names, nil
}

func (a *API) workspaceTagNames(r *http.Request, organizationID, workspaceID string) ([]string, error) {
	tagNames, err := a.organizationTagNames(r, organizationID)
	if err != nil {
		return nil, err
	}
	workspaceTags, err := a.deps.Store.ListWorkspaceTagsByWorkspace(r.Context(), workspaceID)
	if err != nil {
		return nil, err
	}
	return workspaceTagNames(workspaceTags, tagNames), nil
}

func workspaceTagNames(workspaceTags []store.Workspacetag, tagNames map[string]string) []string {
	seen := map[string]struct{}{}
	names := make([]string, 0, len(workspaceTags))
	for _, workspaceTag := range workspaceTags {
		name, ok := tagNames[workspaceTag.TagID]
		if !ok || name == "" {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func matchesTagFilters(workspaceTags []store.Workspacetag, tagNames map[string]string, query workspaceQuery) bool {
	if len(query.searchTags) == 0 && len(query.tagged) == 0 {
		return true
	}
	values := make(map[string]string, len(workspaceTags))
	for _, workspaceTag := range workspaceTags {
		key, ok := tagNames[workspaceTag.TagID]
		if !ok || key == "" {
			continue
		}
		if workspaceTag.TagValue == nil {
			values[key] = ""
		} else {
			values[key] = *workspaceTag.TagValue
		}
	}
	for _, key := range query.searchTags {
		if _, ok := values[key]; !ok {
			return false
		}
	}
	for _, filter := range query.tagged {
		value, ok := values[filter.key]
		if !ok {
			return false
		}
		if filter.value != "" && filter.value != value {
			return false
		}
	}
	return true
}

// chiURLParam keeps the URL parameter lookup in one place so the handlers stay
// readable next to the upstream controller methods they mirror.
func chiURLParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}
