package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type vcsAttrs struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	VcsType        *string `json:"vcsType"`
	ClientID       *string `json:"clientId"`
	ClientSecret   *string `json:"clientSecret"`
	AccessToken    *string `json:"accessToken"`
	Callback       *string `json:"callback"`
	RedirectURL    *string `json:"redirectUrl"`
	Endpoint       *string `json:"endpoint"`
	APIURL         *string `json:"apiUrl"`
	ConnectionType *string `json:"connectionType"`
	Status         *string `json:"status"`
	PrivateKey     *string `json:"privateKey"`
}

var (
	vcsTypes = map[string]bool{
		"GITHUB": true, "GITLAB": true, "BITBUCKET": true,
		"AZURE_DEVOPS": true, "AZURE_SP_MI": true, "PUBLIC": true,
	}
	vcsConnectionTypes = map[string]bool{"OAUTH": true, "STANDALONE": true}
	vcsStatuses        = map[string]bool{"PENDING": true, "COMPLETED": true, "ERROR": true}
)

// vcsResource hides clientSecret, accessToken, refreshToken and privateKey:
// upstream only exposes them to a registry service or the read-vcs-secret
// permission, and the UI never reads them back.
func (a *API) vcsResource(ctx context.Context, vcs store.Vc) (jsonapi.Resource, error) {
	ids, err := a.deps.Store.ListWorkspaceIDsByVCS(ctx, strPtr(vcs.ID))
	if err != nil {
		return jsonapi.Resource{}, err
	}
	return jsonapi.Resource{
		Type: "vcs",
		ID:   vcs.ID,
		Attributes: map[string]any{
			"apiUrl":         vcs.ApiUrl,
			"callback":       vcs.Callback,
			"clientId":       vcs.ClientID,
			"connectionType": vcs.ConnectionType,
			"createdBy":      vcs.CreatedBy,
			"createdDate":    timestamp(vcs.CreatedDate),
			"description":    vcs.Description,
			"endpoint":       vcs.Endpoint,
			"name":           vcs.Name,
			"redirectUrl":    vcs.RedirectUrl,
			"status":         vcs.Status,
			"updatedBy":      vcs.UpdatedBy,
			"updatedDate":    timestamp(vcs.UpdatedDate),
			"vcsType":        vcs.VcsType,
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", vcs.OrganizationID),
			"workspace":    jsonapi.LinkList("workspace", ids...),
		},
	}, nil
}

func (a *API) listVCS(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	providers, err := a.deps.Store.ListVCSByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list vcs providers", err) {
			a.internalError(w, "list vcs providers", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(providers))
	for _, provider := range providers {
		resource, err := a.vcsResource(r.Context(), provider)
		if err != nil {
			if !a.writeStoreError(w, "list vcs workspaces", err) {
				a.internalError(w, "list vcs workspaces", err)
			}
			return
		}
		resources = append(resources, resource)
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getVCS(w http.ResponseWriter, r *http.Request) {
	provider, ok := a.vcsInOrganization(w, r)
	if !ok {
		return
	}
	resource, err := a.vcsResource(r.Context(), provider)
	if err != nil {
		if !a.writeStoreError(w, "get vcs workspaces", err) {
			a.internalError(w, "get vcs workspaces", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, resource, jsonapi.ParseFieldsets(r))
}

func (a *API) createVCS(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	attrs, data, ok := decodeVCSAttrs(w, r)
	if !ok {
		return
	}
	// Upstream's columns are NOT NULL, so Elide turns a missing required
	// attribute into the same data-integrity 409 as a constraint violation.
	if attrs.Name == nil || attrs.Description == nil || attrs.VcsType == nil {
		writeInUseConflict(w)
		return
	}
	if !validVCSEnum(w, *attrs.VcsType, vcsTypes) {
		return
	}
	connectionType, ok := vcsEnumOrDefault(w, data, "connectionType", attrs.ConnectionType, vcsConnectionTypes, "OAUTH")
	if !ok {
		return
	}
	status, ok := vcsEnumOrDefault(w, data, "status", attrs.Status, vcsStatuses, "PENDING")
	if !ok {
		return
	}

	provider, err := a.deps.Store.CreateVCS(r.Context(), store.CreateVCSParams{
		ID:             uuid.NewString(),
		Name:           *attrs.Name,
		Description:    *attrs.Description,
		VcsType:        *attrs.VcsType,
		ClientID:       attrs.ClientID,
		ClientSecret:   attrs.ClientSecret,
		AccessToken:    attrs.AccessToken,
		Status:         strPtr(status),
		Callback:       attrs.Callback,
		RedirectUrl:    attrs.RedirectURL,
		Endpoint:       attrs.Endpoint,
		ApiUrl:         attrs.APIURL,
		ConnectionType: strPtr(connectionType),
		PrivateKey:     attrs.PrivateKey,
		OrganizationID: orgID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create vcs provider", err) {
			a.internalError(w, "create vcs provider", err)
		}
		return
	}
	resource, err := a.vcsResource(r.Context(), provider)
	if err != nil {
		if !a.writeStoreError(w, "create vcs workspaces", err) {
			a.internalError(w, "create vcs workspaces", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, resource, jsonapi.ParseFieldsets(r))
}

func (a *API) updateVCS(w http.ResponseWriter, r *http.Request) {
	provider, ok := a.vcsInOrganization(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodeVCSAttrs(w, r)
	if !ok {
		return
	}
	name, ok := requiredVCSString(w, data, "name", attrs.Name, provider.Name)
	if !ok {
		return
	}
	description, ok := requiredVCSString(w, data, "description", attrs.Description, provider.Description)
	if !ok {
		return
	}
	vcsType, ok := requiredVCSString(w, data, "vcsType", attrs.VcsType, provider.VcsType)
	if !ok || !validVCSEnum(w, vcsType, vcsTypes) {
		return
	}
	connectionType, ok := vcsEnumOrDefault(w, data, "connectionType", attrs.ConnectionType, vcsConnectionTypes, strValue(provider.ConnectionType))
	if !ok {
		return
	}
	status, ok := vcsEnumOrDefault(w, data, "status", attrs.Status, vcsStatuses, strValue(provider.Status))
	if !ok {
		return
	}

	if _, err := a.deps.Store.UpdateVCS(r.Context(), store.UpdateVCSParams{
		ID:             provider.ID,
		Name:           name,
		Description:    description,
		VcsType:        vcsType,
		ClientID:       pickString(data, "clientId", attrs.ClientID, provider.ClientID),
		ClientSecret:   pickString(data, "clientSecret", attrs.ClientSecret, provider.ClientSecret),
		AccessToken:    pickString(data, "accessToken", attrs.AccessToken, provider.AccessToken),
		Status:         strPtr(status),
		Callback:       pickString(data, "callback", attrs.Callback, provider.Callback),
		RedirectUrl:    pickString(data, "redirectUrl", attrs.RedirectURL, provider.RedirectUrl),
		Endpoint:       pickString(data, "endpoint", attrs.Endpoint, provider.Endpoint),
		ApiUrl:         pickString(data, "apiUrl", attrs.APIURL, provider.ApiUrl),
		ConnectionType: strPtr(connectionType),
		PrivateKey:     pickString(data, "privateKey", attrs.PrivateKey, provider.PrivateKey),
	}); err != nil {
		if !a.writeStoreError(w, "update vcs provider", err) {
			a.internalError(w, "update vcs provider", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// deleteVCS unlinks its workspaces first: upstream's mapped relationship nulls
// workspace.vcs_id on delete, and leaving the id behind would keep a dangling
// link.
func (a *API) deleteVCS(w http.ResponseWriter, r *http.Request) {
	provider, ok := a.vcsInOrganization(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.ClearWorkspaceVCS(r.Context(), strPtr(provider.ID)); err != nil {
		if !a.writeStoreError(w, "clear vcs workspaces", err) {
			a.internalError(w, "clear vcs workspaces", err)
		}
		return
	}
	if err := a.deps.Store.DeleteVCS(r.Context(), provider.ID); err != nil {
		if !a.writeStoreError(w, "delete vcs provider", err) {
			a.internalError(w, "delete vcs provider", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) vcsInOrganization(w http.ResponseWriter, r *http.Request) (store.Vc, bool) {
	provider, err := a.deps.Store.GetVCS(r.Context(), param(r, "vcsID"))
	if err != nil {
		if !a.writeStoreError(w, "get vcs provider", err) {
			a.internalError(w, "get vcs provider", err)
		}
		return store.Vc{}, false
	}
	if provider.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Vc{}, false
	}
	return provider, true
}

func decodeVCSAttrs(w http.ResponseWriter, r *http.Request) (vcsAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return vcsAttrs{}, nil, false
	}
	if !requireType(w, data, "vcs") {
		return vcsAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[vcsAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return vcsAttrs{}, nil, false
	}
	return attrs, data, true
}

// requiredVCSString keeps a NOT NULL attribute when the patch omits it and
// answers upstream's 409 when the patch sends null.
func requiredVCSString(w http.ResponseWriter, data *jsonapi.RequestData, name string, incoming *string, existing string) (string, bool) {
	if !jsonapi.HasAttribute(data, name) {
		return existing, true
	}
	if incoming == nil {
		writeInUseConflict(w)
		return "", false
	}
	return *incoming, true
}

func vcsEnumOrDefault(w http.ResponseWriter, data *jsonapi.RequestData, name string, incoming *string, allowed map[string]bool, fallback string) (string, bool) {
	value := fallback
	if jsonapi.HasAttribute(data, name) {
		if incoming == nil {
			value = fallback
		} else {
			value = *incoming
		}
	}
	if !validVCSEnum(w, value, allowed) {
		return "", false
	}
	return value, true
}

// validVCSEnum matches Elide's 400 for an unknown enum member.
func validVCSEnum(w http.ResponseWriter, value string, allowed map[string]bool) bool {
	if !allowed[value] {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+value)
		return false
	}
	return true
}
