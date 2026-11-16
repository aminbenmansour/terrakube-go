package tfe

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const maxStateVersionBytes = 100 << 20

type stateVersionAttributes struct {
	State     *string `json:"state"`
	JSONState *string `json:"json-state"`
	Serial    *int32  `json:"serial"`
	Lineage   *string `json:"lineage"`
}

type stateDocumentMetadata struct {
	Serial           int32  `json:"serial"`
	Lineage          string `json:"lineage"`
	TerraformVersion string `json:"terraform_version"`
}

func (a *API) createStateVersion(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}
	if a.deps.Storage == nil {
		a.internalError(w, "create state version", errors.New("storage is not configured"))
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	attrs, err := jsonapi.DecodeAttributes[stateVersionAttributes](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}

	stateBytes, err := decodeStateAttribute(attrs.State)
	if err != nil {
		jsonapi.WriteError(w, http.StatusUnprocessableEntity, "Invalid State", "state must be base64 encoded")
		return
	}
	jsonBytes, err := decodeStateAttribute(attrs.JSONState)
	if err != nil {
		jsonapi.WriteError(w, http.StatusUnprocessableEntity, "Invalid State", "json-state must be base64 encoded")
		return
	}
	if len(jsonBytes) == 0 {
		jsonBytes = []byte("{}")
	}

	var metadata stateDocumentMetadata
	if len(stateBytes) > 0 {
		// The explicit serial/lineage attributes win, but the state document is
		// the fallback used by clients that only upload the raw state.
		_ = json.Unmarshal(stateBytes, &metadata)
	}
	serial := metadata.Serial
	if attrs.Serial != nil {
		serial = *attrs.Serial
	}
	if serial == 0 {
		serial = 1
	}
	lineage := metadata.Lineage
	if attrs.Lineage != nil {
		lineage = *attrs.Lineage
	}
	if lineage == "" {
		lineage = "0"
	}

	job, err := a.deps.Store.CreateJob(r.Context(), store.CreateJobParams{
		OrganizationID: ws.OrganizationID,
		WorkspaceID:    ws.ID,
		Command:        stringPointer("state-upload"),
		Status:         "completed",
		Via:            "API",
		AutoApply:      false,
		Refresh:        true,
		PlanChanges:    true,
		CreatedBy:      stringPointer("tfe-api"),
	})
	if err != nil {
		a.internalError(w, "create state upload job", err)
		return
	}
	if _, err := a.deps.Store.CreateStep(r.Context(), store.CreateStepParams{
		ID:         uuid.NewString(),
		JobID:      job.ID,
		StepNumber: 100,
		Name:       stringPointer("Dummy State Uploaded"),
		Status:     stringPointer("completed"),
		LogStatus:  stringPointer("completed"),
	}); err != nil {
		a.internalError(w, "create state upload step", err)
		return
	}

	historyID := uuid.NewString()
	history, err := a.deps.Store.CreateHistory(r.Context(), store.CreateHistoryParams{
		ID:           historyID,
		WorkspaceID:  ws.ID,
		JobReference: stringPointer(strconv.FormatInt(job.ID, 10)),
		Output:       stringPointer(a.stateJSONURL(ws, historyID)),
		Serial:       serial,
		Md5:          "0",
		Lineage:      lineage,
		CreatedBy:    stringPointer("tfe-api"),
	})
	if err != nil {
		a.internalError(w, "create state history", err)
		return
	}

	if len(stateBytes) > 0 {
		if err := a.deps.Storage.Put(r.Context(), a.deps.Layout.RawState(ws.OrganizationID, ws.ID, history.ID), bytes.NewReader(stateBytes), int64(len(stateBytes)), "application/json"); err != nil {
			a.internalError(w, "store raw state", err)
			return
		}
	} else {
		// Match upstream: when the client defers the state upload, expose pending
		// archive URLs instead of pretending the version is finalized.
		if _, err := a.deps.Store.CreateArchive(r.Context(), store.CreateArchiveParams{
			ID:        uuid.NewString(),
			Type:      stringPointer("RAW"),
			HistoryID: history.ID,
		}); err != nil {
			a.internalError(w, "create raw state archive", err)
			return
		}
		if _, err := a.deps.Store.CreateArchive(r.Context(), store.CreateArchiveParams{
			ID:        uuid.NewString(),
			Type:      stringPointer("JSON"),
			HistoryID: history.ID,
		}); err != nil {
			a.internalError(w, "create json state archive", err)
			return
		}
	}

	if err := a.deps.Storage.Put(r.Context(), a.deps.Layout.JSONState(ws.OrganizationID, ws.ID, history.ID), bytes.NewReader(jsonBytes), int64(len(jsonBytes)), "application/json"); err != nil {
		a.internalError(w, "store json state", err)
		return
	}

	resource, err := a.stateVersionResource(r, history, ws)
	if err != nil {
		a.internalError(w, "load state version", err)
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, resource)
}

func (a *API) getStateVersion(w http.ResponseWriter, r *http.Request) {
	history, err := a.deps.Store.GetHistory(r.Context(), chi.URLParam(r, "historyID"))
	if err != nil {
		if !a.writeStoreError(w, "get state version", err) {
			a.internalError(w, "get state version", err)
		}
		return
	}
	ws, err := a.deps.Store.GetWorkspace(r.Context(), history.WorkspaceID)
	if err != nil {
		if !a.writeStoreError(w, "get state version workspace", err) {
			a.internalError(w, "get state version workspace", err)
		}
		return
	}
	resource, err := a.stateVersionResource(r, history, ws)
	if err != nil {
		a.internalError(w, "load state version", err)
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, resource)
}

func (a *API) getCurrentStateVersion(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}
	if a.deps.Storage == nil {
		a.internalError(w, "get current state version", errors.New("storage is not configured"))
		return
	}

	key := a.deps.Layout.CurrentState(ws.OrganizationID, ws.ID)
	obj, err := a.deps.Storage.Stat(r.Context(), key)
	if errors.Is(err, storage.ErrNotFound) || (err == nil && obj.Size == 0) {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The workspace has no current state version.")
		return
	}
	if err != nil {
		a.internalError(w, "stat current state", err)
		return
	}
	rc, err := a.deps.Storage.Get(r.Context(), key)
	if errors.Is(err, storage.ErrNotFound) {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The workspace has no current state version.")
		return
	}
	if err != nil {
		a.internalError(w, "read current state", err)
		return
	}
	defer rc.Close()
	body, err := io.ReadAll(io.LimitReader(rc, maxStateVersionBytes))
	if err != nil {
		a.internalError(w, "read current state", err)
		return
	}
	if len(body) == 0 {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The workspace has no current state version.")
		return
	}

	var metadata stateDocumentMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		a.internalError(w, "parse current state", err)
		return
	}
	terraformVersion := any(nil)
	if metadata.TerraformVersion != "" {
		terraformVersion = metadata.TerraformVersion
	} else if ws.TerraformVersion != nil {
		terraformVersion = *ws.TerraformVersion
	}

	jsonapi.WriteResource(w, http.StatusOK, jsonapi.Resource{
		Type: "state-versions",
		ID:   uuid.NewString(),
		Attributes: map[string]any{
			"vcs-commit-url":            nil,
			"vcs-commit-sha":            nil,
			"hosted-state-download-url": a.stateURL(ws, "terraform.tfstate"),
			"serial":                    metadata.Serial,
			"terraform-version":         terraformVersion,
			"status":                    "finalized",
		},
	})
}

func (a *API) stateVersionResource(r *http.Request, history store.History, ws store.Workspace) (jsonapi.Resource, error) {
	archives, err := a.deps.Store.ListArchivesByHistory(r.Context(), history.ID)
	if err != nil {
		return jsonapi.Resource{}, err
	}
	attributes := map[string]any{
		"vcs-commit-sha":                 nil,
		"vcs-commit-url":                 nil,
		"hosted-state-download-url":      a.stateURL(ws, "terraform.tfstate"),
		"hosted-json-state-download-url": a.stateJSONURL(ws, history.ID),
		"status":                         "finalized",
	}
	if len(archives) > 0 {
		attributes["status"] = "pending"
		for _, archive := range archives {
			if archive.Type == nil {
				continue
			}
			switch *archive.Type {
			case "RAW":
				attributes["hosted-state-upload-url"] = a.archiveURL(archive.ID, "terraform.tfstate")
			case "JSON":
				attributes["hosted-json-state-upload-url"] = a.archiveURL(archive.ID, "terraform.json.tfstate")
			}
		}
	}
	return jsonapi.Resource{Type: "state-versions", ID: history.ID, Attributes: attributes}, nil
}

func (a *API) stateURL(ws store.Workspace, filename string) string {
	return fmt.Sprintf("%s/tfstate/v1/organization/%s/workspace/%s/state/%s",
		strings.TrimRight(a.deps.PublicURL, "/"), ws.OrganizationID, ws.ID, filename)
}

func (a *API) stateJSONURL(ws store.Workspace, historyID string) string {
	return a.stateURL(ws, historyID+".json")
}

func (a *API) archiveURL(archiveID, filename string) string {
	return fmt.Sprintf("%s/tfstate/v1/archive/%s/%s",
		strings.TrimRight(a.deps.PublicURL, "/"), archiveID, filename)
}

func decodeStateAttribute(value *string) ([]byte, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	return io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(*value)))
}

func stringPointer(value string) *string {
	return &value
}
