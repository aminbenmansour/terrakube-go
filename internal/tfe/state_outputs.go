package tfe

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
)

type stateOutput struct {
	Value     any  `json:"value"`
	Type      any  `json:"type"`
	Sensitive bool `json:"sensitive"`
}

type stateWithOutputs struct {
	Outputs map[string]stateOutput `json:"outputs"`
}

// getCurrentStateVersionOutputs reads the current backend state and exposes its
// root module outputs as state-version-outputs. Upstream returns an empty list
// when the workspace has no state or the state cannot be parsed.
func (a *API) getCurrentStateVersionOutputs(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}
	if a.deps.Storage == nil {
		a.internalError(w, "get current state outputs", errors.New("storage is not configured"))
		return
	}

	rc, err := a.deps.Storage.Get(r.Context(), a.deps.Layout.CurrentState(ws.OrganizationID, ws.ID))
	if errors.Is(err, storage.ErrNotFound) {
		jsonapi.WriteCollection(w, http.StatusOK, nil, 0, false)
		return
	}
	if err != nil {
		a.internalError(w, "read current state outputs", err)
		return
	}
	defer rc.Close()

	body, err := io.ReadAll(io.LimitReader(rc, maxStateVersionBytes))
	if err != nil {
		a.internalError(w, "read current state outputs", err)
		return
	}
	if len(body) == 0 {
		jsonapi.WriteCollection(w, http.StatusOK, nil, 0, false)
		return
	}

	var state stateWithOutputs
	if err := json.Unmarshal(body, &state); err != nil {
		a.deps.Logger.Warn("current state outputs are not valid JSON", "error", err, "workspace_id", ws.ID)
		jsonapi.WriteCollection(w, http.StatusOK, nil, 0, false)
		return
	}

	names := make([]string, 0, len(state.Outputs))
	for name := range state.Outputs {
		names = append(names, name)
	}
	sort.Strings(names)

	resources := make([]jsonapi.Resource, 0, len(names))
	for _, name := range names {
		output := state.Outputs[name]
		outputType := "object"
		detailedType := output.Type
		if _, ok := output.Type.(string); ok {
			outputType = "string"
			detailedType = "string"
		}
		resources = append(resources, jsonapi.Resource{
			Type: "state-version-outputs",
			ID:   uuid.NewString(),
			Attributes: map[string]any{
				"name":          name,
				"value":         output.Value,
				"sensitive":     output.Sensitive,
				"type":          outputType,
				"detailed-type": detailedType,
			},
		})
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), false)
}
