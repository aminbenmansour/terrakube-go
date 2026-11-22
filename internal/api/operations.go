package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
)

// The UI's workspace create and settings screens use the JSON:API atomic
// operations extension instead of POST /workspace and PATCH /workspace/{id}.
type atomicOperation struct {
	Op   string               `json:"op"`
	Href string               `json:"href"`
	Data *jsonapi.RequestData `json:"data"`
}

type atomicResult struct {
	Data any `json:"data"`
}

// atomicRequest asks for the same media type on Accept as on Content-Type;
// upstream answers 415 when either is missing.
func isAtomicRequest(r *http.Request) bool {
	const marker = "jsonapi.org/ext/atomic"
	return strings.Contains(r.Header.Get("Content-Type"), marker) &&
		strings.Contains(r.Header.Get("Accept"), marker)
}

func (a *API) atomicOperations(w http.ResponseWriter, r *http.Request) {
	if !isAtomicRequest(r) {
		jsonapi.WriteError(w, http.StatusUnsupportedMediaType, "Unsupported Media Type", "Unsupported Media Type")
		return
	}

	var doc struct {
		Operations []atomicOperation `json:"atomic:operations"`
	}
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", "invalid JSON: "+err.Error())
		return
	}

	results := make([]atomicResult, 0, len(doc.Operations))
	for _, operation := range doc.Operations {
		switch operation.Op {
		case "add":
			orgID, _, ok := atomicWorkspaceHref(operation.Href, false)
			if !ok {
				jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "Unsupported href: "+operation.Href)
				return
			}
			if operation.Data == nil || operation.Data.Type != "workspace" {
				jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "data.type must be workspace")
				return
			}
			ws, ok := a.createWorkspaceFromData(w, r, orgID, operation.Data)
			if !ok {
				return
			}
			results = append(results, atomicResult{Data: workspaceResource(ws)})
		case "update":
			orgID, wsID, ok := atomicWorkspaceHref(operation.Href, true)
			if !ok {
				jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "Unsupported href: "+operation.Href)
				return
			}
			ws, err := a.deps.Store.GetWorkspace(r.Context(), wsID)
			if err != nil {
				if !a.writeStoreError(w, "get workspace", err) {
					a.internalError(w, "get workspace", err)
				}
				return
			}
			if ws.OrganizationID != orgID {
				jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier "+wsID+" for workspace")
				return
			}
			if operation.Data == nil || operation.Data.Type != "workspace" {
				jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "data.type must be workspace")
				return
			}
			if !a.applyWorkspaceUpdate(w, r, ws, operation.Data) {
				return
			}
			results = append(results, atomicResult{Data: nil})
		default:
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "Unsupported op: "+operation.Op)
			return
		}
	}

	w.Header().Set("Content-Type", jsonapi.MediaType+"; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"atomic:results": results})
}

// atomicWorkspaceHref reads the href the UI sends: the collection
// /organization/{org}/workspace for an add, the item for an update.
func atomicWorkspaceHref(href string, item bool) (orgID string, wsID string, ok bool) {
	segments := strings.Split(strings.Trim(strings.TrimPrefix(href, "/organization/"), "/"), "/")
	if len(segments) != 2 && len(segments) != 3 {
		return "", "", false
	}
	if segments[1] != "workspace" {
		return "", "", false
	}
	if _, err := uuid.Parse(segments[0]); err != nil {
		return "", "", false
	}
	if len(segments) == 3 {
		if !item {
			return "", "", false
		}
		if _, err := uuid.Parse(segments[2]); err != nil {
			return "", "", false
		}
		return segments[0], segments[2], true
	}
	if item {
		return "", "", false
	}
	return segments[0], "", true
}
