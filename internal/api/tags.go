package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type tagAttrs struct {
	Name *string `json:"name"`
}

type workspaceTagAttrs struct {
	TagID *string `json:"tagId"`
	Value *string `json:"value"`
}

func tagResource(tag store.Tag) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "tag",
		ID:   tag.ID,
		Attributes: map[string]any{
			"createdBy":   tag.CreatedBy,
			"createdDate": timestamp(tag.CreatedDate),
			"name":        tag.Name,
			"updatedBy":   tag.UpdatedBy,
			"updatedDate": timestamp(tag.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", tag.OrganizationID),
		},
	}
}

func workspaceTagResource(tag store.Workspacetag) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "workspacetag",
		ID:   tag.ID,
		Attributes: map[string]any{
			"createdBy":   tag.CreatedBy,
			"createdDate": timestamp(tag.CreatedDate),
			"tagId":       tag.TagID,
			"updatedBy":   tag.UpdatedBy,
			"updatedDate": timestamp(tag.UpdatedDate),
			"value":       tag.TagValue,
		},
		Relationships: map[string]jsonapi.Relationship{
			"workspace": jsonapi.Link("workspace", tag.WorkspaceID),
		},
	}
}

func (a *API) listTags(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	items, err := a.deps.Store.ListTagsByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list tags", err) {
			a.internalError(w, "list tags", err)
		}
		return
	}
	q := parseList(r, "tag")
	page, total := filterSortPage(items, q, func(tag store.Tag) string { return strValue(tag.Name) })
	resources := make([]jsonapi.Resource, 0, len(page))
	for _, tag := range page {
		resources = append(resources, tagResource(tag))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, total, q.page(), jsonapi.ParseFieldsets(r))
}

func (a *API) createTag(w http.ResponseWriter, r *http.Request) {
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
	if !requireType(w, data, "tag") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[tagAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	tag, err := a.deps.Store.CreateTag(r.Context(), store.CreateTagParams{
		ID:             uuid.NewString(),
		Name:           attrs.Name,
		OrganizationID: orgID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create tag", err) {
			a.internalError(w, "create tag", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, tagResource(tag), jsonapi.ParseFieldsets(r))
}

func (a *API) getTag(w http.ResponseWriter, r *http.Request) {
	tag, ok := a.tagInOrganization(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, tagResource(tag), jsonapi.ParseFieldsets(r))
}

func (a *API) updateTag(w http.ResponseWriter, r *http.Request) {
	tag, ok := a.tagInOrganization(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "tag") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[tagAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if _, err := a.deps.Store.UpdateTag(r.Context(), store.UpdateTagParams{
		ID:   tag.ID,
		Name: pickString(data, "name", attrs.Name, tag.Name),
	}); err != nil {
		if !a.writeStoreError(w, "update tag", err) {
			a.internalError(w, "update tag", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteTag(w http.ResponseWriter, r *http.Request) {
	tag, ok := a.tagInOrganization(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteTag(r.Context(), tag.ID); err != nil {
		if !a.writeStoreError(w, "delete tag", err) {
			a.internalError(w, "delete tag", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listWorkspaceTags(w http.ResponseWriter, r *http.Request) {
	items, err := a.deps.Store.ListWorkspaceTagsByWorkspace(r.Context(), param(r, "workspaceID"))
	if err != nil {
		if !a.writeStoreError(w, "list workspace tags", err) {
			a.internalError(w, "list workspace tags", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(items))
	for _, item := range items {
		resources = append(resources, workspaceTagResource(item))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) createWorkspaceTag(w http.ResponseWriter, r *http.Request) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "workspacetag") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[workspaceTagAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.TagID == nil || *attrs.TagID == "" {
		writeInUseConflict(w)
		return
	}
	// A tag id that does not exist is a foreign key violation upstream, which
	// Elide reports as the same 409 as a missing one.
	if _, err := a.deps.Store.GetTag(r.Context(), *attrs.TagID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeInUseConflict(w)
			return
		}
		if !a.writeStoreError(w, "get tag", err) {
			a.internalError(w, "get tag", err)
		}
		return
	}
	tag, err := a.deps.Store.CreateWorkspaceTag(r.Context(), store.CreateWorkspaceTagParams{
		ID:          uuid.NewString(),
		TagID:       *attrs.TagID,
		TagValue:    attrs.Value,
		WorkspaceID: param(r, "workspaceID"),
	})
	if err != nil {
		if !a.writeStoreError(w, "create workspace tag", err) {
			a.internalError(w, "create workspace tag", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, workspaceTagResource(tag), jsonapi.ParseFieldsets(r))
}

func (a *API) updateWorkspaceTag(w http.ResponseWriter, r *http.Request) {
	tag, ok := a.workspaceTagInWorkspace(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "workspacetag") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[workspaceTagAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if _, err := a.deps.Store.UpdateWorkspaceTagValue(r.Context(), store.UpdateWorkspaceTagValueParams{
		ID:       tag.ID,
		TagValue: pickString(data, "value", attrs.Value, tag.TagValue),
	}); err != nil {
		if !a.writeStoreError(w, "update workspace tag", err) {
			a.internalError(w, "update workspace tag", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteWorkspaceTag(w http.ResponseWriter, r *http.Request) {
	tag, ok := a.workspaceTagInWorkspace(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteWorkspaceTag(r.Context(), tag.ID); err != nil {
		if !a.writeStoreError(w, "delete workspace tag", err) {
			a.internalError(w, "delete workspace tag", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// tagInOrganization and workspaceTagInWorkspace prove the row belongs to the
// organization or workspace in the path, as the ssh and collection helpers do.
func (a *API) tagInOrganization(w http.ResponseWriter, r *http.Request) (store.Tag, bool) {
	tag, err := a.deps.Store.GetTag(r.Context(), param(r, "tagID"))
	if err != nil {
		if !a.writeStoreError(w, "get tag", err) {
			a.internalError(w, "get tag", err)
		}
		return store.Tag{}, false
	}
	if tag.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Tag{}, false
	}
	return tag, true
}

func (a *API) workspaceTagInWorkspace(w http.ResponseWriter, r *http.Request) (store.Workspacetag, bool) {
	tag, err := a.deps.Store.GetWorkspaceTag(r.Context(), param(r, "workspaceTagID"))
	if err != nil {
		if !a.writeStoreError(w, "get workspace tag", err) {
			a.internalError(w, "get workspace tag", err)
		}
		return store.Workspacetag{}, false
	}
	if tag.WorkspaceID != param(r, "workspaceID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Workspacetag{}, false
	}
	return tag, true
}
