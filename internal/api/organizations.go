package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type organizationAttrs struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Disabled      *bool   `json:"disabled"`
	ExecutionMode *string `json:"executionMode"`
	Icon          *string `json:"icon"`
}

func (a *API) listOrganizations(w http.ResponseWriter, r *http.Request) {
	q := parseList(r, "organization")

	items, err := a.deps.Store.ListOrganizations(r.Context())
	if err != nil {
		if !a.writeStoreError(w, "list organizations", err) {
			a.internalError(w, "list organizations", err)
		}
		return
	}

	// Upstream filters the list to the caller's organizations; a token with no
	// matching group sees an empty collection, not an error.
	if !a.callerIsOwner(r) {
		identity, _ := auth.IdentityFromContext(r.Context())
		ids, err := a.deps.Store.ListOrganizationIDsByGroupNames(r.Context(), identity.Groups)
		if err != nil {
			a.internalError(w, "list organizations", err)
			return
		}
		visible := make(map[string]bool, len(ids))
		for _, id := range ids {
			visible[id] = true
		}
		visibleItems := items[:0]
		for _, org := range items {
			if visible[org.ID] {
				visibleItems = append(visibleItems, org)
			}
		}
		items = visibleItems
	}

	page, total := filterSortPage(items, q, func(o store.Organization) string { return o.Name })
	resources := make([]jsonapi.Resource, 0, len(page))
	for _, org := range page {
		resources = append(resources, organizationResource(org))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, total, q.page(), jsonapi.ParseFieldsets(r))
}

func (a *API) getOrganization(w http.ResponseWriter, r *http.Request) {
	org, err := a.deps.Store.GetOrganization(r.Context(), param(r, "organizationID"))
	if err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, organizationResource(org), jsonapi.ParseFieldsets(r))
}

func (a *API) createOrganization(w http.ResponseWriter, r *http.Request) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "organization") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[organizationAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.Name == nil || *attrs.Name == "" {
		jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Missing Attribute", "name is required.", "name")
		return
	}

	executionMode := "remote"
	if attrs.ExecutionMode != nil && *attrs.ExecutionMode != "" {
		executionMode = *attrs.ExecutionMode
	}

	org, err := a.deps.Store.CreateOrganization(r.Context(), store.CreateOrganizationParams{
		ID:            uuid.NewString(),
		Name:          *attrs.Name,
		Description:   attrs.Description,
		Disabled:      boolValue(attrs.Disabled, false),
		ExecutionMode: executionMode,
		Icon:          attrs.Icon,
	})
	if err != nil {
		if !a.writeStoreError(w, "create organization", err) {
			a.internalError(w, "create organization", err)
		}
		return
	}

	jsonapi.WriteResource(w, http.StatusCreated, organizationResource(org), jsonapi.ParseFieldsets(r))
}

func (a *API) updateOrganization(w http.ResponseWriter, r *http.Request) {
	org, err := a.deps.Store.GetOrganization(r.Context(), param(r, "organizationID"))
	if err != nil {
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
	if !requireType(w, data, "organization") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[organizationAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	name := org.Name
	if jsonapi.HasAttribute(data, "name") {
		if attrs.Name == nil || *attrs.Name == "" {
			jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Invalid Attribute", "name must not be empty.", "name")
			return
		}
		name = *attrs.Name
	}
	description := org.Description
	if jsonapi.HasAttribute(data, "description") {
		description = attrs.Description
	}
	disabled := org.Disabled
	if jsonapi.HasAttribute(data, "disabled") {
		disabled = boolValue(attrs.Disabled, org.Disabled)
	}
	executionMode := org.ExecutionMode
	if jsonapi.HasAttribute(data, "executionMode") && attrs.ExecutionMode != nil && *attrs.ExecutionMode != "" {
		executionMode = *attrs.ExecutionMode
	}
	icon := org.Icon
	if jsonapi.HasAttribute(data, "icon") {
		icon = attrs.Icon
	}

	updated, err := a.deps.Store.UpdateOrganization(r.Context(), store.UpdateOrganizationParams{
		ID:            org.ID,
		Name:          name,
		Description:   description,
		Disabled:      disabled,
		ExecutionMode: executionMode,
		Icon:          icon,
	})
	if err != nil {
		if !a.writeStoreError(w, "update organization", err) {
			a.internalError(w, "update organization", err)
		}
		return
	}

	jsonapi.WriteResource(w, http.StatusOK, organizationResource(updated), jsonapi.ParseFieldsets(r))
}

func (a *API) deleteOrganization(w http.ResponseWriter, r *http.Request) {
	id := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), id); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	if err := a.deps.Store.DeleteOrganization(r.Context(), id); err != nil {
		if !a.writeStoreError(w, "delete organization", err) {
			a.internalError(w, "delete organization", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func organizationResource(o store.Organization) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "organization",
		ID:   o.ID,
		Attributes: map[string]any{
			"name":          o.Name,
			"description":   o.Description,
			"disabled":      o.Disabled,
			"executionMode": o.ExecutionMode,
			"icon":          o.Icon,
		},
	}
}
