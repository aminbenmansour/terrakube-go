package api

import (
	"net/http"
	"strings"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// actionAttrs is what the Settings > Actions editor sends. `action` is
// base64-encoded React source: stored and returned verbatim, never executed
// here.
type actionAttrs struct {
	Name            *string `json:"name"`
	Label           *string `json:"label"`
	Category        *string `json:"category"`
	Description     *string `json:"description"`
	Type            *string `json:"type"`
	Action          *string `json:"action"`
	DisplayCriteria *string `json:"displayCriteria"`
	Version         *string `json:"version"`
	Active          *bool   `json:"active"`
}

func actionResource(action store.Action) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "action",
		ID:   action.ID,
		Attributes: map[string]any{
			"action":          action.Action,
			"active":          action.Active,
			"category":        action.Category,
			"createdBy":       action.CreatedBy,
			"createdDate":     timestamp(action.CreatedDate),
			"description":     action.Description,
			"displayCriteria": action.DisplayCriteria,
			"label":           action.Label,
			"name":            action.Name,
			"type":            action.Type,
			"updatedBy":       action.UpdatedBy,
			"updatedDate":     timestamp(action.UpdatedDate),
			"version":         action.Version,
		},
	}
}

// actionFilter is the UI's filter[action] expression, `active==true;type=in=('A','B')`.
// Unrecognised terms are ignored, like the other in-memory filters: the drawer
// sends only these two.
type actionFilter struct {
	active *bool
	types  []string
	typeOf bool
}

func parseActionFilter(raw string) actionFilter {
	var f actionFilter
	for _, term := range strings.Split(raw, ";") {
		term = strings.TrimSpace(term)
		switch {
		case strings.HasPrefix(term, "active=="):
			active := strings.TrimPrefix(term, "active==") == "true"
			f.active = &active
		case strings.HasPrefix(term, "type=in=("):
			f.typeOf = true
			inner := strings.TrimSuffix(strings.TrimPrefix(term, "type=in=("), ")")
			for _, name := range strings.Split(inner, ",") {
				f.types = append(f.types, strings.Trim(strings.TrimSpace(name), "'"))
			}
		}
	}
	return f
}

func (f actionFilter) match(action store.Action) bool {
	if f.active != nil && action.Active != *f.active {
		return false
	}
	if f.typeOf {
		for _, name := range f.types {
			if action.Type == name {
				return true
			}
		}
		return false
	}
	return true
}

func (a *API) listActions(w http.ResponseWriter, r *http.Request) {
	items, err := a.deps.Store.ListActions(r.Context())
	if err != nil {
		if !a.writeStoreError(w, "list actions", err) {
			a.internalError(w, "list actions", err)
		}
		return
	}

	filter := parseActionFilter(r.URL.Query().Get("filter[action]"))
	resources := make([]jsonapi.Resource, 0, len(items))
	for _, action := range items {
		if filter.match(action) {
			resources = append(resources, actionResource(action))
		}
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getAction(w http.ResponseWriter, r *http.Request) {
	action, err := a.deps.Store.GetAction(r.Context(), param(r, "actionID"))
	if err != nil {
		if !a.writeStoreError(w, "get action", err) {
			a.internalError(w, "get action", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, actionResource(action), jsonapi.ParseFieldsets(r))
}

func (a *API) createAction(w http.ResponseWriter, r *http.Request) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "action") {
		return
	}
	// Upstream's ids are slugs the editor picks ("terrakube.open-ai"), so there
	// is nothing to generate.
	if data.ID == "" {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", "No id provided, cannot persist action")
		return
	}

	attrs, ok := decodeActionAttrs(w, data)
	if !ok {
		return
	}
	action, err := a.deps.Store.CreateAction(r.Context(), store.CreateActionParams{
		ID:              data.ID,
		Name:            *attrs.Name,
		Label:           *attrs.Label,
		Category:        *attrs.Category,
		Description:     attrs.Description,
		Type:            *attrs.Type,
		Action:          *attrs.Action,
		DisplayCriteria: *attrs.DisplayCriteria,
		Version:         *attrs.Version,
		Active:          boolValue(attrs.Active, true),
	})
	if err != nil {
		if !a.writeStoreError(w, "create action", err) {
			a.internalError(w, "create action", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, actionResource(action), jsonapi.ParseFieldsets(r))
}

// updateAction is a partial update, as Elide's PATCH is: an omitted attribute
// keeps its current value. Upstream answers 204 with no body.
func (a *API) updateAction(w http.ResponseWriter, r *http.Request) {
	existing, err := a.deps.Store.GetAction(r.Context(), param(r, "actionID"))
	if err != nil {
		if !a.writeStoreError(w, "get action", err) {
			a.internalError(w, "get action", err)
		}
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "action") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[actionAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if _, err := a.deps.Store.UpdateAction(r.Context(), store.UpdateActionParams{
		ID:              existing.ID,
		Name:            stringOr(attrs.Name, existing.Name),
		Label:           stringOr(attrs.Label, existing.Label),
		Category:        stringOr(attrs.Category, existing.Category),
		Description:     attrs.Description,
		Type:            stringOr(attrs.Type, existing.Type),
		Action:          stringOr(attrs.Action, existing.Action),
		DisplayCriteria: stringOr(attrs.DisplayCriteria, existing.DisplayCriteria),
		Version:         stringOr(attrs.Version, existing.Version),
		Active:          boolValue(attrs.Active, existing.Active),
	}); err != nil {
		if !a.writeStoreError(w, "update action", err) {
			a.internalError(w, "update action", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteAction(w http.ResponseWriter, r *http.Request) {
	if err := a.deps.Store.DeleteAction(r.Context(), param(r, "actionID")); err != nil {
		if !a.writeStoreError(w, "delete action", err) {
			a.internalError(w, "delete action", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// decodeActionAttrs validates the create payload: every NOT NULL column must be
// present. The boolean defaults to true, matching the column default.
func decodeActionAttrs(w http.ResponseWriter, data *jsonapi.RequestData) (actionAttrs, bool) {
	attrs, err := jsonapi.DecodeAttributes[actionAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return actionAttrs{}, false
	}
	for _, required := range []struct {
		name  string
		value *string
	}{
		{"name", attrs.Name},
		{"label", attrs.Label},
		{"category", attrs.Category},
		{"type", attrs.Type},
		{"action", attrs.Action},
		{"displayCriteria", attrs.DisplayCriteria},
		{"version", attrs.Version},
	} {
		if required.value == nil || *required.value == "" {
			jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Missing Attribute", required.name+" is required.", required.name)
			return actionAttrs{}, false
		}
	}
	return attrs, true
}

func stringOr(value *string, fallback string) string {
	if value == nil {
		return fallback
	}
	return *value
}
