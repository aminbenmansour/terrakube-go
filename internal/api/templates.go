package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type templateAttrs struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Version     *string `json:"version"`
	Tcl         *string `json:"tcl"`
}

func templateResource(template store.Template) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "template",
		ID:   template.ID,
		Attributes: map[string]any{
			"createdBy":   template.CreatedBy,
			"createdDate": timestamp(template.CreatedDate),
			"description": template.Description,
			"name":        template.Name,
			"tcl":         template.Tcl,
			"updatedBy":   template.UpdatedBy,
			"updatedDate": timestamp(template.UpdatedDate),
			"version":     template.Version,
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", template.OrganizationID),
		},
	}
}

func (a *API) listTemplates(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	templates, err := a.deps.Store.ListTemplatesByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list templates", err) {
			a.internalError(w, "list templates", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(templates))
	for _, template := range templates {
		resources = append(resources, templateResource(template))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getTemplate(w http.ResponseWriter, r *http.Request) {
	template, ok := a.templateInOrganization(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, templateResource(template), jsonapi.ParseFieldsets(r))
}

func (a *API) createTemplate(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	attrs, _, ok := decodeTemplateAttrs(w, r)
	if !ok {
		return
	}
	template, err := a.deps.Store.CreateTemplate(r.Context(), store.CreateTemplateParams{
		ID:             uuid.NewString(),
		Name:           attrs.Name,
		Description:    attrs.Description,
		Version:        attrs.Version,
		Tcl:            attrs.Tcl,
		OrganizationID: orgID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create template", err) {
			a.internalError(w, "create template", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, templateResource(template), jsonapi.ParseFieldsets(r))
}

func (a *API) updateTemplate(w http.ResponseWriter, r *http.Request) {
	template, ok := a.templateInOrganization(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodeTemplateAttrs(w, r)
	if !ok {
		return
	}
	if _, err := a.deps.Store.UpdateTemplate(r.Context(), store.UpdateTemplateParams{
		ID:          template.ID,
		Name:        pickString(data, "name", attrs.Name, template.Name),
		Description: pickString(data, "description", attrs.Description, template.Description),
		Version:     pickString(data, "version", attrs.Version, template.Version),
		Tcl:         pickString(data, "tcl", attrs.Tcl, template.Tcl),
	}); err != nil {
		if !a.writeStoreError(w, "update template", err) {
			a.internalError(w, "update template", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	template, ok := a.templateInOrganization(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteTemplate(r.Context(), template.ID); err != nil {
		if !a.writeStoreError(w, "delete template", err) {
			a.internalError(w, "delete template", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) templateInOrganization(w http.ResponseWriter, r *http.Request) (store.Template, bool) {
	id := param(r, "templateID")
	if _, err := uuid.Parse(id); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
		return store.Template{}, false
	}
	template, err := a.deps.Store.GetTemplate(r.Context(), id)
	if err != nil {
		if !a.writeStoreError(w, "get template", err) {
			a.internalError(w, "get template", err)
		}
		return store.Template{}, false
	}
	if template.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Template{}, false
	}
	return template, true
}

func decodeTemplateAttrs(w http.ResponseWriter, r *http.Request) (templateAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return templateAttrs{}, nil, false
	}
	if !requireType(w, data, "template") {
		return templateAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[templateAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return templateAttrs{}, nil, false
	}
	return attrs, data, true
}
