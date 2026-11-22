package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type collectionAttrs struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Priority    *int    `json:"priority"`
}

type itemAttrs struct {
	Key         *string `json:"key"`
	Value       *string `json:"value"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Sensitive   *bool   `json:"sensitive"`
	Hcl         *bool   `json:"hcl"`
}

// collectionResource carries the audit attributes upstream's Elide serialises
// even though the editor never sends them. itemIDs and referenceIDs are the
// relationships the collection page and the workspace drawer read.
func collectionResource(collection store.Collection, itemIDs, referenceIDs []string) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "collection",
		ID:   collection.ID,
		Attributes: map[string]any{
			"createdBy":   collection.CreatedBy,
			"createdDate": timestamp(collection.CreatedDate),
			"description": collection.Description,
			"name":        collection.Name,
			"priority":    collection.Priority,
			"updatedBy":   collection.UpdatedBy,
			"updatedDate": timestamp(collection.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"item":         jsonapi.LinkList("item", itemIDs...),
			"organization": jsonapi.Link("organization", collection.OrganizationID),
			"reference":    jsonapi.LinkList("reference", referenceIDs...),
		},
	}
}

func itemResource(item store.Item) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "item",
		ID:   item.ID,
		Attributes: map[string]any{
			"category":    item.ItemCategory,
			"description": item.ItemDescription,
			"hcl":         item.Hcl,
			"key":         item.ItemKey,
			"sensitive":   item.Sensitive,
			"value":       item.ItemValue,
		},
		Relationships: map[string]jsonapi.Relationship{
			"collection": jsonapi.Link("collection", item.CollectionID),
		},
	}
}

func referenceResource(reference store.Reference) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "reference",
		ID:   reference.ID,
		Attributes: map[string]any{
			"description": reference.Description,
		},
		Relationships: map[string]jsonapi.Relationship{
			"collection": jsonapi.Link("collection", reference.CollectionID),
			"workspace":  jsonapi.Link("workspace", reference.WorkspaceID),
		},
	}
}

func (a *API) listCollections(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}

	collections, err := a.deps.Store.ListCollectionsByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list collections", err) {
			a.internalError(w, "list collections", err)
		}
		return
	}

	resources := make([]jsonapi.Resource, 0, len(collections))
	for _, collection := range collections {
		resource, ok := a.collectionResource(r.Context(), w, "list collections", collection)
		if !ok {
			return
		}
		resources = append(resources, resource)
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getCollection(w http.ResponseWriter, r *http.Request) {
	collection, ok := a.collectionInOrganization(w, r)
	if !ok {
		return
	}
	resource, ok := a.collectionResource(r.Context(), w, "get collection", collection)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, resource, jsonapi.ParseFieldsets(r))
}

func (a *API) createCollection(w http.ResponseWriter, r *http.Request) {
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
	if !requireType(w, data, "collection") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[collectionAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.Name == nil || *attrs.Name == "" {
		// Elide answers its data-integrity conflict when a required column is
		// missing, as it does for globalvar (see globalvars.go).
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This operation conflicts with data that already exists.")
		return
	}

	collection, err := a.deps.Store.CreateCollection(r.Context(), store.CreateCollectionParams{
		ID:             uuid.NewString(),
		Name:           *attrs.Name,
		Description:    attrs.Description,
		Priority:       int32(intValue(attrs.Priority, 0)),
		OrganizationID: orgID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create collection", err) {
			a.internalError(w, "create collection", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, collectionResource(collection, nil, nil), jsonapi.ParseFieldsets(r))
}

func (a *API) updateCollection(w http.ResponseWriter, r *http.Request) {
	collection, ok := a.collectionInOrganization(w, r)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "collection") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[collectionAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	if _, err := a.deps.Store.UpdateCollection(r.Context(), store.UpdateCollectionParams{
		ID:          collection.ID,
		Name:        stringOr(attrs.Name, collection.Name),
		Description: pickString(data, "description", attrs.Description, collection.Description),
		Priority:    int32(intValue(attrs.Priority, int(collection.Priority))),
	}); err != nil {
		if !a.writeStoreError(w, "update collection", err) {
			a.internalError(w, "update collection", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteCollection(w http.ResponseWriter, r *http.Request) {
	collection, ok := a.collectionInOrganization(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteCollection(r.Context(), collection.ID); err != nil {
		if !a.writeStoreError(w, "delete collection", err) {
			a.internalError(w, "delete collection", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listCollectionItems(w http.ResponseWriter, r *http.Request) {
	collection, ok := a.collectionInOrganization(w, r)
	if !ok {
		return
	}
	items, err := a.deps.Store.ListItemsByCollection(r.Context(), collection.ID)
	if err != nil {
		if !a.writeStoreError(w, "list collection items", err) {
			a.internalError(w, "list collection items", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(items))
	for _, item := range items {
		resources = append(resources, itemResource(item))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) createCollectionItem(w http.ResponseWriter, r *http.Request) {
	collection, ok := a.collectionInOrganization(w, r)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "item") {
		return
	}
	attrs, ok := decodeItemAttrs(w, data)
	if !ok {
		return
	}

	item, err := a.deps.Store.CreateItem(r.Context(), store.CreateItemParams{
		ID:              uuid.NewString(),
		ItemDescription: attrs.Description,
		ItemKey:         *attrs.Key,
		ItemValue:       *attrs.Value,
		Sensitive:       boolValue(attrs.Sensitive, false),
		Hcl:             boolValue(attrs.Hcl, false),
		ItemCategory:    *attrs.Category,
		CollectionID:    collection.ID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create collection item", err) {
			a.internalError(w, "create collection item", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, itemResource(item), jsonapi.ParseFieldsets(r))
}

func (a *API) updateCollectionItem(w http.ResponseWriter, r *http.Request) {
	item, ok := a.itemInCollection(w, r)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "item") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[itemAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	category := item.ItemCategory
	if jsonapi.HasAttribute(data, "category") {
		normalized, ok := globalVarCategory(attrs.Category)
		if !ok {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+stringValue(attrs.Category))
			return
		}
		category = normalized
	}

	if _, err := a.deps.Store.UpdateItem(r.Context(), store.UpdateItemParams{
		ID:              item.ID,
		ItemDescription: pickString(data, "description", attrs.Description, item.ItemDescription),
		ItemKey:         stringOr(attrs.Key, item.ItemKey),
		ItemValue:       stringOr(attrs.Value, item.ItemValue),
		Sensitive:       boolValue(attrs.Sensitive, item.Sensitive),
		Hcl:             boolValue(attrs.Hcl, item.Hcl),
		ItemCategory:    category,
	}); err != nil {
		if !a.writeStoreError(w, "update collection item", err) {
			a.internalError(w, "update collection item", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteCollectionItem(w http.ResponseWriter, r *http.Request) {
	item, ok := a.itemInCollection(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteItem(r.Context(), item.ID); err != nil {
		if !a.writeStoreError(w, "delete collection item", err) {
			a.internalError(w, "delete collection item", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listCollectionReferences(w http.ResponseWriter, r *http.Request) {
	collection, ok := a.collectionInOrganization(w, r)
	if !ok {
		return
	}
	references, err := a.deps.Store.ListReferencesByCollection(r.Context(), collection.ID)
	if err != nil {
		if !a.writeStoreError(w, "list collection references", err) {
			a.internalError(w, "list collection references", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(references))
	for _, reference := range references {
		resources = append(resources, referenceResource(reference))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) createCollectionReference(w http.ResponseWriter, r *http.Request) {
	collection, ok := a.collectionInOrganization(w, r)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "reference") {
		return
	}
	workspaceID, present, err := jsonapi.RelationshipID(data, "workspace")
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", err.Error())
		return
	}
	if !present || workspaceID == "" {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This operation conflicts with data that already exists.")
		return
	}
	ws, err := a.deps.Store.GetWorkspace(r.Context(), workspaceID)
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}
	// A reference may only bind a workspace of the collection's organization.
	if ws.OrganizationID != collection.OrganizationID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier "+workspaceID+" for workspace")
		return
	}

	attrs, err := jsonapi.DecodeAttributes[struct {
		Description *string `json:"description"`
	}](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	reference, err := a.deps.Store.CreateReference(r.Context(), store.CreateReferenceParams{
		ID:           uuid.NewString(),
		Description:  attrs.Description,
		WorkspaceID:  ws.ID,
		CollectionID: collection.ID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create collection reference", err) {
			a.internalError(w, "create collection reference", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, referenceResource(reference), jsonapi.ParseFieldsets(r))
}

func (a *API) deleteCollectionReference(w http.ResponseWriter, r *http.Request) {
	collection, ok := a.collectionInOrganization(w, r)
	if !ok {
		return
	}
	reference, err := a.deps.Store.GetReference(r.Context(), param(r, "referenceID"))
	if err != nil {
		if !a.writeStoreError(w, "get reference", err) {
			a.internalError(w, "get reference", err)
		}
		return
	}
	if reference.CollectionID != collection.ID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return
	}
	if err := a.deps.Store.DeleteReference(r.Context(), reference.ID); err != nil {
		if !a.writeStoreError(w, "delete collection reference", err) {
			a.internalError(w, "delete collection reference", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// getReferenceCollection is the read the workspace page makes for every
// reference it receives: the collection plus its items, which fill the
// read-only "Collection Variables" table.
func (a *API) getReferenceCollection(w http.ResponseWriter, r *http.Request) {
	reference, err := a.deps.Store.GetReference(r.Context(), param(r, "referenceID"))
	if err != nil {
		if !a.writeStoreError(w, "get reference", err) {
			a.internalError(w, "get reference", err)
		}
		return
	}
	collection, err := a.deps.Store.GetCollection(r.Context(), reference.CollectionID)
	if err != nil {
		if !a.writeStoreError(w, "get collection", err) {
			a.internalError(w, "get collection", err)
		}
		return
	}

	resource, ok := a.collectionResource(r.Context(), w, "get reference collection", collection)
	if !ok {
		return
	}

	var included []jsonapi.Resource
	if parseInclude(r)["item"] {
		items, err := a.deps.Store.ListItemsByCollection(r.Context(), collection.ID)
		if err != nil {
			if !a.writeStoreError(w, "list collection items", err) {
				a.internalError(w, "list collection items", err)
			}
			return
		}
		for _, item := range items {
			included = append(included, itemResource(item))
		}
	}

	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: resource, Included: included}, jsonapi.ParseFieldsets(r))
}

// collectionInOrganization loads the collection addressed by the nested route
// and proves it belongs to the organization in the path.
func (a *API) collectionInOrganization(w http.ResponseWriter, r *http.Request) (store.Collection, bool) {
	collection, err := a.deps.Store.GetCollection(r.Context(), param(r, "collectionID"))
	if err != nil {
		if !a.writeStoreError(w, "get collection", err) {
			a.internalError(w, "get collection", err)
		}
		return store.Collection{}, false
	}
	if collection.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Collection{}, false
	}
	return collection, true
}

// itemInCollection loads the item and proves it belongs to the collection in
// the path.
func (a *API) itemInCollection(w http.ResponseWriter, r *http.Request) (store.Item, bool) {
	item, err := a.deps.Store.GetItem(r.Context(), param(r, "itemID"))
	if err != nil {
		if !a.writeStoreError(w, "get collection item", err) {
			a.internalError(w, "get collection item", err)
		}
		return store.Item{}, false
	}
	if item.CollectionID != param(r, "collectionID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Item{}, false
	}
	return item, true
}

// collectionItems resolves the item and reference ids a collection resource
// lists in its relationships.
func (a *API) collectionResource(ctx context.Context, w http.ResponseWriter, operation string, collection store.Collection) (jsonapi.Resource, bool) {
	items, err := a.deps.Store.ListItemsByCollection(ctx, collection.ID)
	if err != nil {
		if !a.writeStoreError(w, operation, err) {
			a.internalError(w, operation, err)
		}
		return jsonapi.Resource{}, false
	}
	references, err := a.deps.Store.ListReferencesByCollection(ctx, collection.ID)
	if err != nil {
		if !a.writeStoreError(w, operation, err) {
			a.internalError(w, operation, err)
		}
		return jsonapi.Resource{}, false
	}
	itemIDs := make([]string, 0, len(items))
	for _, item := range items {
		itemIDs = append(itemIDs, item.ID)
	}
	referenceIDs := make([]string, 0, len(references))
	for _, reference := range references {
		referenceIDs = append(referenceIDs, reference.ID)
	}
	return collectionResource(collection, itemIDs, referenceIDs), true
}

// decodeItemAttrs validates the NOT NULL columns on a new collection item.
func decodeItemAttrs(w http.ResponseWriter, data *jsonapi.RequestData) (itemAttrs, bool) {
	attrs, err := jsonapi.DecodeAttributes[itemAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return itemAttrs{}, false
	}
	if attrs.Key == nil || *attrs.Key == "" || attrs.Value == nil || attrs.Category == nil {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This operation conflicts with data that already exists.")
		return itemAttrs{}, false
	}
	category, ok := globalVarCategory(attrs.Category)
	if !ok {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+stringValue(attrs.Category))
		return itemAttrs{}, false
	}
	attrs.Category = &category
	return attrs, true
}

func intValue(value *int, fallback int) int {
	if value == nil {
		return fallback
	}
	return *value
}
