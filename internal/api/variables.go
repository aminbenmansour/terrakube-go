package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type variableAttrs struct {
	Key         *string `json:"key"`
	Value       *string `json:"value"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Sensitive   *bool   `json:"sensitive"`
	Hcl         *bool   `json:"hcl"`
	Incomplete  *bool   `json:"incomplete"`
}

func (a *API) listVariables(w http.ResponseWriter, r *http.Request) {
	wsID := param(r, "workspaceID")
	if _, err := a.deps.Store.GetWorkspace(r.Context(), wsID); err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}

	items, err := a.deps.Store.ListVariablesByWorkspace(r.Context(), wsID)
	if err != nil {
		if !a.writeStoreError(w, "list variables", err) {
			a.internalError(w, "list variables", err)
		}
		return
	}

	resources := make([]jsonapi.Resource, 0, len(items))
	for _, variable := range items {
		resources = append(resources, variableResource(variable))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

// variableInWorkspace guards the nested variable routes: the workspace is
// already proven to belong to the organization by the outer guard, this stops a
// variable id from a sibling workspace being read or written through it.
func (a *API) variableInWorkspace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		variable, err := a.deps.Store.GetVariable(r.Context(), param(r, "variableID"))
		if err != nil {
			if !a.writeStoreError(w, "get variable", err) {
				a.internalError(w, "get variable", err)
			}
			return
		}
		if variable.WorkspaceID != param(r, "workspaceID") {
			jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *API) getVariable(w http.ResponseWriter, r *http.Request) {
	variable, err := a.deps.Store.GetVariable(r.Context(), param(r, "variableID"))
	if err != nil {
		if !a.writeStoreError(w, "get variable", err) {
			a.internalError(w, "get variable", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, variableResource(variable), jsonapi.ParseFieldsets(r))
}

func (a *API) createVariable(w http.ResponseWriter, r *http.Request) {
	wsID := param(r, "workspaceID")
	if _, err := a.deps.Store.GetWorkspace(r.Context(), wsID); err != nil {
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
	if !requireType(w, data, "variable") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[variableAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.Key == nil || *attrs.Key == "" {
		jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Missing Attribute", "key is required.", "key")
		return
	}

	category, ok := normalizeCategory(attrs.Category)
	if !ok {
		jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Invalid Attribute", "category must be TERRAFORM or ENV.", "category")
		return
	}

	variable, err := a.deps.Store.CreateVariable(r.Context(), store.CreateVariableParams{
		ID:                  uuid.NewString(),
		VariableKey:         *attrs.Key,
		VariableValue:       stringValue(attrs.Value),
		VariableDescription: attrs.Description,
		VariableCategory:    category,
		Sensitive:           boolValue(attrs.Sensitive, false),
		Hcl:                 boolValue(attrs.Hcl, false),
		Incomplete:          boolValue(attrs.Incomplete, false),
		WorkspaceID:         wsID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create variable", err) {
			a.internalError(w, "create variable", err)
		}
		return
	}

	jsonapi.WriteResource(w, http.StatusCreated, variableResource(variable), jsonapi.ParseFieldsets(r))
}

func (a *API) updateVariable(w http.ResponseWriter, r *http.Request) {
	variable, err := a.deps.Store.GetVariable(r.Context(), param(r, "variableID"))
	if err != nil {
		if !a.writeStoreError(w, "get variable", err) {
			a.internalError(w, "get variable", err)
		}
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "variable") {
		return
	}

	attrs, err := jsonapi.DecodeAttributes[variableAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	key := variable.VariableKey
	if jsonapi.HasAttribute(data, "key") {
		if attrs.Key == nil || *attrs.Key == "" {
			jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Invalid Attribute", "key must not be empty.", "key")
			return
		}
		key = *attrs.Key
	}

	category := variable.VariableCategory
	if jsonapi.HasAttribute(data, "category") {
		normalized, ok := normalizeCategory(attrs.Category)
		if !ok {
			jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Invalid Attribute", "category must be TERRAFORM or ENV.", "category")
			return
		}
		category = normalized
	}

	value := variable.VariableValue
	if jsonapi.HasAttribute(data, "value") {
		value = stringValue(attrs.Value)
	}

	updated, err := a.deps.Store.UpdateVariable(r.Context(), store.UpdateVariableParams{
		ID:                  variable.ID,
		VariableKey:         key,
		VariableValue:       value,
		VariableDescription: pickString(data, "description", attrs.Description, variable.VariableDescription),
		VariableCategory:    category,
		Sensitive:           pickBool(data, "sensitive", attrs.Sensitive, variable.Sensitive),
		Hcl:                 pickBool(data, "hcl", attrs.Hcl, variable.Hcl),
		Incomplete:          pickBool(data, "incomplete", attrs.Incomplete, variable.Incomplete),
	})
	if err != nil {
		if !a.writeStoreError(w, "update variable", err) {
			a.internalError(w, "update variable", err)
		}
		return
	}

	jsonapi.WriteResource(w, http.StatusOK, variableResource(updated), jsonapi.ParseFieldsets(r))
}

func (a *API) deleteVariable(w http.ResponseWriter, r *http.Request) {
	id := param(r, "variableID")
	if _, err := a.deps.Store.GetVariable(r.Context(), id); err != nil {
		if !a.writeStoreError(w, "get variable", err) {
			a.internalError(w, "get variable", err)
		}
		return
	}
	if err := a.deps.Store.DeleteVariable(r.Context(), id); err != nil {
		if !a.writeStoreError(w, "delete variable", err) {
			a.internalError(w, "delete variable", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func variableResource(v store.Variable) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "variable",
		ID:   v.ID,
		Attributes: map[string]any{
			"key":         v.VariableKey,
			"value":       v.VariableValue,
			"description": v.VariableDescription,
			"category":    v.VariableCategory,
			"sensitive":   v.Sensitive,
			"hcl":         v.Hcl,
			"incomplete":  v.Incomplete,
		},
		Relationships: map[string]jsonapi.Relationship{
			"workspace": jsonapi.Link("workspace", v.WorkspaceID),
		},
	}
}

func normalizeCategory(category *string) (string, bool) {
	if category == nil || *category == "" {
		return "TERRAFORM", true
	}
	value := strings.ToUpper(strings.TrimSpace(*category))
	switch value {
	case "TERRAFORM", "ENV":
		return value, true
	default:
		return "", false
	}
}

func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
