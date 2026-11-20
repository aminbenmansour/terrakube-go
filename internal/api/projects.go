package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type projectAttrs struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (a *API) listProjects(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}

	q := parseList(r, "project")
	items, err := a.deps.Store.ListProjectsByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list projects", err) {
			a.internalError(w, "list projects", err)
		}
		return
	}

	page, total := filterSortPage(items, q, func(p store.Project) string { return p.Name })
	resources := make([]jsonapi.Resource, 0, len(page))
	for _, project := range page {
		resources = append(resources, projectResource(project))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, total, q.page())
}

func (a *API) getProject(w http.ResponseWriter, r *http.Request) {
	project, err := a.deps.Store.GetProject(r.Context(), param(r, "projectID"))
	if err != nil {
		if !a.writeStoreError(w, "get project", err) {
			a.internalError(w, "get project", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, projectResource(project))
}

func (a *API) createProject(w http.ResponseWriter, r *http.Request) {
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
	if !requireType(w, data, "project") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[projectAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.Name == nil || *attrs.Name == "" {
		jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Missing Attribute", "name is required.", "name")
		return
	}

	project, err := a.deps.Store.CreateProject(r.Context(), store.CreateProjectParams{
		ID:             uuid.NewString(),
		Name:           *attrs.Name,
		Description:    attrs.Description,
		OrganizationID: orgID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create project", err) {
			a.internalError(w, "create project", err)
		}
		return
	}

	jsonapi.WriteResource(w, http.StatusCreated, projectResource(project))
}

func (a *API) updateProject(w http.ResponseWriter, r *http.Request) {
	project, err := a.deps.Store.GetProject(r.Context(), param(r, "projectID"))
	if err != nil {
		if !a.writeStoreError(w, "get project", err) {
			a.internalError(w, "get project", err)
		}
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "project") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[projectAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	name := project.Name
	if jsonapi.HasAttribute(data, "name") {
		if attrs.Name == nil || *attrs.Name == "" {
			jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Invalid Attribute", "name must not be empty.", "name")
			return
		}
		name = *attrs.Name
	}
	description := project.Description
	if jsonapi.HasAttribute(data, "description") {
		description = attrs.Description
	}

	updated, err := a.deps.Store.UpdateProject(r.Context(), store.UpdateProjectParams{
		ID:          project.ID,
		Name:        name,
		Description: description,
	})
	if err != nil {
		if !a.writeStoreError(w, "update project", err) {
			a.internalError(w, "update project", err)
		}
		return
	}

	jsonapi.WriteResource(w, http.StatusOK, projectResource(updated))
}

func (a *API) deleteProject(w http.ResponseWriter, r *http.Request) {
	id := param(r, "projectID")
	if _, err := a.deps.Store.GetProject(r.Context(), id); err != nil {
		if !a.writeStoreError(w, "get project", err) {
			a.internalError(w, "get project", err)
		}
		return
	}
	if err := a.deps.Store.DeleteProject(r.Context(), id); err != nil {
		if !a.writeStoreError(w, "delete project", err) {
			a.internalError(w, "delete project", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func projectResource(p store.Project) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "project",
		ID:   p.ID,
		Attributes: map[string]any{
			"name":        p.Name,
			"description": p.Description,
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", p.OrganizationID),
		},
	}
}
