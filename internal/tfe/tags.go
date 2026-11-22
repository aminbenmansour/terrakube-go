package tfe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	maxTagsPerWorkspace = 10
	maxTagKeyLength     = 128
	maxTagValueLength   = 256
	maxTFEJSONBody      = 8 << 20
)

type tagBindingInput struct {
	Key   string
	Value *string
}

func (a *API) listTagBindings(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}
	bindings, err := a.tagBindings(r, ws)
	if err != nil {
		a.internalError(w, "list tag bindings", err)
		return
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: bindings})
}

func (a *API) updateTagBindings(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}

	var doc struct {
		Data []struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			Attributes struct {
				Key   string  `json:"key"`
				Value *string `json:"value"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := decodeTFEJSON(r, &doc); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}

	bindings := make([]tagBindingInput, 0, len(doc.Data))
	for _, item := range doc.Data {
		bindings = append(bindings, tagBindingInput{Key: item.Attributes.Key, Value: item.Attributes.Value})
	}
	if err := validateTagBindings(bindings); err != nil {
		jsonapi.WriteError(w, http.StatusUnprocessableEntity, "Invalid Tag Binding", err.Error())
		return
	}
	if err := a.applyTagBindings(r, ws, bindings, true); err != nil {
		a.internalError(w, "update tag bindings", err)
		return
	}

	result, err := a.tagBindings(r, ws)
	if err != nil {
		a.internalError(w, "list tag bindings", err)
		return
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: result})
}

func (a *API) updateWorkspaceTags(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}

	var doc struct {
		Data []struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			Attributes struct {
				Name string `json:"name"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := decodeTFEJSON(r, &doc); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}

	bindings := make([]tagBindingInput, 0, len(doc.Data))
	for _, item := range doc.Data {
		if strings.TrimSpace(item.Attributes.Name) == "" {
			continue
		}
		bindings = append(bindings, tagBindingInput{Key: strings.TrimSpace(item.Attributes.Name)})
	}
	if err := validateTagBindings(bindings); err != nil {
		jsonapi.WriteError(w, http.StatusUnprocessableEntity, "Invalid Tag", err.Error())
		return
	}
	if err := a.applyTagBindings(r, ws, bindings, false); err != nil {
		a.internalError(w, "update workspace tags", err)
		return
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) tagBindings(r *http.Request, ws store.Workspace) ([]jsonapi.Resource, error) {
	tagNames, err := a.organizationTagNames(r, ws.OrganizationID)
	if err != nil {
		return nil, err
	}
	workspaceTags, err := a.deps.Store.ListWorkspaceTagsByWorkspace(r.Context(), ws.ID)
	if err != nil {
		return nil, err
	}

	resources := make([]jsonapi.Resource, 0, len(workspaceTags))
	for _, workspaceTag := range workspaceTags {
		key, ok := tagNames[workspaceTag.TagID]
		if !ok || key == "" {
			continue
		}
		value := ""
		if workspaceTag.TagValue != nil {
			value = *workspaceTag.TagValue
		}
		resources = append(resources, jsonapi.Resource{
			Type: "tag-bindings",
			ID:   workspaceTag.ID,
			Attributes: map[string]any{
				"key":   key,
				"value": value,
			},
		})
	}
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Attributes["key"].(string) < resources[j].Attributes["key"].(string)
	})
	return resources, nil
}

// applyTagBindings upserts key/value bindings. overwriteValue is false for the
// legacy key-only tags API, which must not clear the value of an existing key.
func (a *API) applyTagBindings(r *http.Request, ws store.Workspace, bindings []tagBindingInput, overwriteValue bool) error {
	current, err := a.deps.Store.ListWorkspaceTagsByWorkspace(r.Context(), ws.ID)
	if err != nil {
		return err
	}
	current = append([]store.Workspacetag(nil), current...)

	for _, binding := range bindings {
		tag, err := a.ensureTag(r, ws.OrganizationID, binding.Key)
		if err != nil {
			return err
		}

		existingIndex := -1
		for i := range current {
			if current[i].TagID == tag.ID {
				existingIndex = i
				break
			}
		}

		normalized := normalizeTagValue(binding.Value)
		if existingIndex >= 0 {
			if !overwriteValue {
				continue
			}
			existing := current[existingIndex]
			if tagValuesEqual(existing.TagValue, normalized) {
				continue
			}
			updated, err := a.deps.Store.UpdateWorkspaceTagValue(r.Context(), store.UpdateWorkspaceTagValueParams{
				ID:       existing.ID,
				TagValue: normalized,
			})
			if err != nil {
				return err
			}
			current[existingIndex] = updated
			continue
		}

		created, err := a.deps.Store.CreateWorkspaceTag(r.Context(), store.CreateWorkspaceTagParams{
			ID:          uuid.NewString(),
			TagID:       tag.ID,
			TagValue:    normalized,
			WorkspaceID: ws.ID,
		})
		if err != nil {
			return err
		}
		current = append(current, created)
	}
	return nil
}

func (a *API) ensureTag(r *http.Request, organizationID, name string) (store.Tag, error) {
	tag, err := a.deps.Store.GetTagByName(r.Context(), store.GetTagByNameParams{
		OrganizationID: organizationID,
		Name:           &name,
	})
	if err == nil {
		return tag, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return store.Tag{}, err
	}

	tag, err = a.deps.Store.CreateTag(r.Context(), store.CreateTagParams{
		ID:             uuid.NewString(),
		Name:           &name,
		OrganizationID: organizationID,
	})
	if err == nil {
		return tag, nil
	}

	// A concurrent request may have created the same organization/name pair.
	if existing, getErr := a.deps.Store.GetTagByName(r.Context(), store.GetTagByNameParams{
		OrganizationID: organizationID,
		Name:           &name,
	}); getErr == nil {
		return existing, nil
	}
	return store.Tag{}, err
}

func validateTagBindings(bindings []tagBindingInput) error {
	keys := map[string]struct{}{}
	for _, binding := range bindings {
		key := strings.TrimSpace(binding.Key)
		if key == "" {
			return errors.New("tag key must not be empty")
		}
		if len(key) > maxTagKeyLength {
			return fmt.Errorf("tag key %q exceeds %d characters", key, maxTagKeyLength)
		}
		if binding.Value != nil && len(*binding.Value) > maxTagValueLength {
			return fmt.Errorf("value of tag %q exceeds %d characters", key, maxTagValueLength)
		}
		keys[key] = struct{}{}
	}
	if len(keys) > maxTagsPerWorkspace {
		return fmt.Errorf("a workspace can have at most %d tags; this request would leave it with %d", maxTagsPerWorkspace, len(keys))
	}
	return nil
}

func normalizeTagValue(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := *value
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func tagValuesEqual(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func decodeTFEJSON(r *http.Request, out any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()
	if err := json.NewDecoder(io.LimitReader(r.Body, maxTFEJSONBody)).Decode(out); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

// workspaceTagBindings extracts the tag relationships go-tfe sends on workspace
// creation. It validates before the workspace is inserted so a bad request
// cannot leave a partially-tagged workspace behind.
func workspaceTagBindings(data *jsonapi.RequestData) ([]tagBindingInput, error) {
	bindings := []tagBindingInput{}

	tags, _, err := jsonapi.RelationshipIdentifiers(data, "tags")
	if err != nil {
		return nil, err
	}
	for _, identifier := range tags {
		name, ok := identifierStringAttribute(identifier, "name")
		if !ok || strings.TrimSpace(name) == "" {
			continue
		}
		bindings = append(bindings, tagBindingInput{Key: strings.TrimSpace(name)})
	}

	keyValues, _, err := jsonapi.RelationshipIdentifiers(data, "tag-bindings")
	if err != nil {
		return nil, err
	}
	for _, identifier := range keyValues {
		key, ok := identifierStringAttribute(identifier, "key")
		if !ok || strings.TrimSpace(key) == "" {
			return nil, errors.New("tag key must not be empty")
		}
		value, _ := identifierStringAttribute(identifier, "value")
		var valuePtr *string
		if value != "" {
			valuePtr = &value
		}
		bindings = append(bindings, tagBindingInput{Key: strings.TrimSpace(key), Value: valuePtr})
	}

	if err := validateTagBindings(bindings); err != nil {
		return nil, err
	}
	return bindings, nil
}

func identifierStringAttribute(identifier jsonapi.ResourceIdentifier, name string) (string, bool) {
	raw, ok := identifier.Attributes[name]
	if !ok || len(raw) == 0 {
		return "", false
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", false
	}
	return value, true
}
