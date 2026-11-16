package tfe

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type projectAttrs struct {
	Name *string `json:"name"`
}

func (a *API) listProjects(w http.ResponseWriter, r *http.Request) {
	org, ok := a.organizationByName(w, r)
	if !ok {
		return
	}
	projects, err := a.deps.Store.ListProjectsByOrganization(r.Context(), org.ID)
	if err != nil {
		if !a.writeStoreError(w, "list projects", err) {
			a.internalError(w, "list projects", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(projects))
	for _, project := range projects {
		resources = append(resources, projectResource(project))
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: resources})
}

func (a *API) createProject(w http.ResponseWriter, r *http.Request) {
	org, ok := a.organizationByName(w, r)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if data.Type != "projects" {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Type", "data.type must be projects.")
		return
	}
	attrs, err := jsonapi.DecodeAttributes[projectAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.Name == nil || strings.TrimSpace(*attrs.Name) == "" {
		jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Missing Attribute", "name is required.", "name")
		return
	}

	project, err := a.deps.Store.CreateProject(r.Context(), store.CreateProjectParams{
		ID:             uuid.NewString(),
		Name:           strings.TrimSpace(*attrs.Name),
		OrganizationID: org.ID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create project", err) {
			a.internalError(w, "create project", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, projectResource(project))
}

func projectResource(project store.Project) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "projects",
		ID:   project.ID,
		Attributes: map[string]any{
			"name": project.Name,
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organizations", project.OrganizationID),
		},
	}
}
