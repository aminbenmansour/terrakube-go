package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type globalVarAttrs struct {
	Key         *string `json:"key"`
	Value       *string `json:"value"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Sensitive   *bool   `json:"sensitive"`
	Hcl         *bool   `json:"hcl"`
}

func (a *API) listGlobalVariables(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}

	items, err := a.deps.Store.ListGlobalVariablesByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list global variables", err) {
			a.internalError(w, "list global variables", err)
		}
		return
	}

	resources := make([]jsonapi.Resource, 0, len(items))
	for _, variable := range items {
		resources = append(resources, globalVarResource(variable))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getGlobalVariable(w http.ResponseWriter, r *http.Request) {
	variable, ok := a.globalVariableInOrganization(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, globalVarResource(variable), jsonapi.ParseFieldsets(r))
}

func (a *API) createGlobalVariable(w http.ResponseWriter, r *http.Request) {
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
	if !requireType(w, data, "globalvar") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[globalVarAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.Key == nil || *attrs.Key == "" || attrs.Value == nil || attrs.Description == nil || attrs.Category == nil {
		// Elide answers its generic data-integrity conflict when a required
		// column is missing, not a 4xx naming the attribute.
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This operation conflicts with data that already exists.")
		return
	}
	category, ok := globalVarCategory(attrs.Category)
	if !ok {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+*attrs.Category)
		return
	}

	variable, err := a.deps.Store.CreateGlobalVariable(r.Context(), store.CreateGlobalVariableParams{
		ID:                  uuid.NewString(),
		VariableKey:         *attrs.Key,
		VariableValue:       *attrs.Value,
		VariableDescription: *attrs.Description,
		VariableCategory:    category,
		Sensitive:           boolValue(attrs.Sensitive, false),
		Hcl:                 boolValue(attrs.Hcl, false),
		OrganizationID:      orgID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create global variable", err) {
			a.internalError(w, "create global variable", err)
		}
		return
	}

	jsonapi.WriteResource(w, http.StatusCreated, globalVarResource(variable), jsonapi.ParseFieldsets(r))
}

func (a *API) updateGlobalVariable(w http.ResponseWriter, r *http.Request) {
	variable, ok := a.globalVariableInOrganization(w, r)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "globalvar") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[globalVarAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	key := variable.VariableKey
	if jsonapi.HasAttribute(data, "key") {
		if attrs.Key == nil || *attrs.Key == "" {
			jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This operation conflicts with data that already exists.")
			return
		}
		key = *attrs.Key
	}

	category := variable.VariableCategory
	if jsonapi.HasAttribute(data, "category") {
		normalized, ok := globalVarCategory(attrs.Category)
		if !ok {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+stringValue(attrs.Category))
			return
		}
		category = normalized
	}

	value := variable.VariableValue
	if jsonapi.HasAttribute(data, "value") {
		value = stringValue(attrs.Value)
	}
	description := variable.VariableDescription
	if jsonapi.HasAttribute(data, "description") {
		description = stringValue(attrs.Description)
	}

	if _, err := a.deps.Store.UpdateGlobalVariable(r.Context(), store.UpdateGlobalVariableParams{
		ID:                  variable.ID,
		VariableKey:         key,
		VariableValue:       value,
		VariableDescription: description,
		VariableCategory:    category,
		Sensitive:           pickBool(data, "sensitive", attrs.Sensitive, variable.Sensitive),
		Hcl:                 pickBool(data, "hcl", attrs.Hcl, variable.Hcl),
	}); err != nil {
		if !a.writeStoreError(w, "update global variable", err) {
			a.internalError(w, "update global variable", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteGlobalVariable(w http.ResponseWriter, r *http.Request) {
	variable, ok := a.globalVariableInOrganization(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteGlobalVariable(r.Context(), variable.ID); err != nil {
		if !a.writeStoreError(w, "delete global variable", err) {
			a.internalError(w, "delete global variable", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// globalVariableInOrganization is the item guard and the loader: a variable
// from another organization answers 404, like the nested workspace routes.
func (a *API) globalVariableInOrganization(w http.ResponseWriter, r *http.Request) (store.Globalvar, bool) {
	variable, err := a.deps.Store.GetGlobalVariable(r.Context(), param(r, "globalVarID"))
	if err != nil {
		if !a.writeStoreError(w, "get global variable", err) {
			a.internalError(w, "get global variable", err)
		}
		return store.Globalvar{}, false
	}
	if variable.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Globalvar{}, false
	}
	return variable, true
}

// globalVarCategory is stricter than normalizeCategory: upstream's enum is
// case-sensitive and answers 400 for "env".
func globalVarCategory(category *string) (string, bool) {
	if category == nil {
		return "", false
	}
	switch *category {
	case "TERRAFORM", "ENV":
		return *category, true
	default:
		return "", false
	}
}

func globalVarResource(v store.Globalvar) jsonapi.Resource {
	attributes := map[string]any{
		"key":         v.VariableKey,
		"description": v.VariableDescription,
		"category":    v.VariableCategory,
		"sensitive":   v.Sensitive,
		"hcl":         v.Hcl,
	}
	// Upstream omits the value of a sensitive variable rather than nulling it.
	if !v.Sensitive {
		attributes["value"] = v.VariableValue
	}
	return jsonapi.Resource{
		Type:       "globalvar",
		ID:         v.ID,
		Attributes: attributes,
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", v.OrganizationID),
		},
	}
}
