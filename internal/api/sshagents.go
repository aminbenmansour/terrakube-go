package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type sshAttrs struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	SSHType     *string `json:"sshType"`
	PrivateKey  *string `json:"privateKey"`
}

type agentAttrs struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	URL         *string `json:"url"`
}

// sshResource omits privateKey: upstream only reveals it to a registry service,
// and the UI never reads it back.
func sshResource(key store.Ssh) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "ssh",
		ID:   key.ID,
		Attributes: map[string]any{
			"createdBy":   key.CreatedBy,
			"createdDate": timestamp(key.CreatedDate),
			"description": key.Description,
			"name":        key.Name,
			"sshType":     key.SshType,
			"updatedBy":   key.UpdatedBy,
			"updatedDate": timestamp(key.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", key.OrganizationID),
		},
	}
}

func agentResource(agent store.Agent) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "agent",
		ID:   agent.ID,
		Attributes: map[string]any{
			"description": agent.Description,
			"name":        agent.Name,
			"url":         agent.Url,
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", agent.OrganizationID),
		},
	}
}

func (a *API) listSSHKeys(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	keys, err := a.deps.Store.ListSSHKeysByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list ssh keys", err) {
			a.internalError(w, "list ssh keys", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(keys))
	for _, key := range keys {
		resources = append(resources, sshResource(key))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getSSHKey(w http.ResponseWriter, r *http.Request) {
	key, ok := a.sshKeyInOrganization(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, sshResource(key), jsonapi.ParseFieldsets(r))
}

func (a *API) createSSHKey(w http.ResponseWriter, r *http.Request) {
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
	if !requireType(w, data, "ssh") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[sshAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	key, err := a.deps.Store.CreateSSHKey(r.Context(), store.CreateSSHKeyParams{
		ID:             uuid.NewString(),
		Name:           attrs.Name,
		Description:    attrs.Description,
		PrivateKey:     attrs.PrivateKey,
		SshType:        attrs.SSHType,
		OrganizationID: orgID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create ssh key", err) {
			a.internalError(w, "create ssh key", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, sshResource(key), jsonapi.ParseFieldsets(r))
}

func (a *API) updateSSHKey(w http.ResponseWriter, r *http.Request) {
	key, ok := a.sshKeyInOrganization(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "ssh") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[sshAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if _, err := a.deps.Store.UpdateSSHKey(r.Context(), store.UpdateSSHKeyParams{
		ID:          key.ID,
		Name:        pickString(data, "name", attrs.Name, key.Name),
		Description: pickString(data, "description", attrs.Description, key.Description),
		PrivateKey:  pickString(data, "privateKey", attrs.PrivateKey, key.PrivateKey),
		SshType:     pickString(data, "sshType", attrs.SSHType, key.SshType),
	}); err != nil {
		if !a.writeStoreError(w, "update ssh key", err) {
			a.internalError(w, "update ssh key", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteSSHKey(w http.ResponseWriter, r *http.Request) {
	key, ok := a.sshKeyInOrganization(w, r)
	if !ok {
		return
	}
	used, err := a.deps.Store.CountWorkspacesWithSSHKey(r.Context(), strPtr(key.ID))
	if err != nil {
		if !a.writeStoreError(w, "count ssh key uses", err) {
			a.internalError(w, "count ssh key uses", err)
		}
		return
	}
	if used > 0 {
		writeInUseConflict(w)
		return
	}
	if err := a.deps.Store.DeleteSSHKey(r.Context(), key.ID); err != nil {
		if !a.writeStoreError(w, "delete ssh key", err) {
			a.internalError(w, "delete ssh key", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listAgents(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	agents, err := a.deps.Store.ListAgentsByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list agents", err) {
			a.internalError(w, "list agents", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(agents))
	for _, agent := range agents {
		resources = append(resources, agentResource(agent))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) createAgent(w http.ResponseWriter, r *http.Request) {
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
	if !requireType(w, data, "agent") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[agentAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	agent, err := a.deps.Store.CreateAgent(r.Context(), store.CreateAgentParams{
		ID:             uuid.NewString(),
		Name:           attrs.Name,
		Url:            attrs.URL,
		Description:    attrs.Description,
		OrganizationID: orgID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create agent", err) {
			a.internalError(w, "create agent", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, agentResource(agent), jsonapi.ParseFieldsets(r))
}

func (a *API) deleteAgent(w http.ResponseWriter, r *http.Request) {
	agent, ok := a.agentInOrganization(w, r)
	if !ok {
		return
	}
	used, err := a.deps.Store.CountWorkspacesWithAgent(r.Context(), strPtr(agent.ID))
	if err != nil {
		if !a.writeStoreError(w, "count agent uses", err) {
			a.internalError(w, "count agent uses", err)
		}
		return
	}
	if used > 0 {
		writeInUseConflict(w)
		return
	}
	if err := a.deps.Store.DeleteAgent(r.Context(), agent.ID); err != nil {
		if !a.writeStoreError(w, "delete agent", err) {
			a.internalError(w, "delete agent", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// sshKeyInOrganization and agentInOrganization prove the row belongs to the
// organization in the path, as the collection helpers do.
func (a *API) sshKeyInOrganization(w http.ResponseWriter, r *http.Request) (store.Ssh, bool) {
	key, err := a.deps.Store.GetSSHKey(r.Context(), param(r, "sshID"))
	if err != nil {
		if !a.writeStoreError(w, "get ssh key", err) {
			a.internalError(w, "get ssh key", err)
		}
		return store.Ssh{}, false
	}
	if key.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Ssh{}, false
	}
	return key, true
}

func (a *API) agentInOrganization(w http.ResponseWriter, r *http.Request) (store.Agent, bool) {
	agent, err := a.deps.Store.GetAgent(r.Context(), param(r, "agentID"))
	if err != nil {
		if !a.writeStoreError(w, "get agent", err) {
			a.internalError(w, "get agent", err)
		}
		return store.Agent{}, false
	}
	if agent.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Agent{}, false
	}
	return agent, true
}

// checkVCSRelationship validates a workspace's vcs link; a provider from
// another organization is refused like a project.
func (a *API) checkVCSRelationship(w http.ResponseWriter, r *http.Request, vcsID, orgID string) bool {
	provider, err := a.deps.Store.GetVCS(r.Context(), vcsID)
	if err != nil {
		return a.relationshipLoadFailed(w, "get vcs provider", err, vcsID, "vcs")
	}
	if provider.OrganizationID != orgID {
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "VCS provider does not belong to the workspace organization")
		return false
	}
	return true
}

// checkSSHRelationship validates a workspace's ssh link against the table that
// now backs it; a key from another organization is refused like a project.
func (a *API) checkSSHRelationship(w http.ResponseWriter, r *http.Request, sshID, orgID string) bool {
	key, err := a.deps.Store.GetSSHKey(r.Context(), sshID)
	if err != nil {
		return a.relationshipLoadFailed(w, "get ssh key", err, sshID, "ssh")
	}
	if key.OrganizationID != orgID {
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "SSH key does not belong to the workspace organization")
		return false
	}
	return true
}

// checkAgentRelationship validates a workspace's agent link.
func (a *API) checkAgentRelationship(w http.ResponseWriter, r *http.Request, agentID, orgID string) bool {
	agent, err := a.deps.Store.GetAgent(r.Context(), agentID)
	if err != nil {
		return a.relationshipLoadFailed(w, "get agent", err, agentID, "agent")
	}
	if agent.OrganizationID != orgID {
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "Agent does not belong to the workspace organization")
		return false
	}
	return true
}

// writeInUseConflict is upstream's answer when a foreign key still points at
// the row: Elide reports the database's data-integrity conflict as 409.
func writeInUseConflict(w http.ResponseWriter) {
	jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This operation conflicts with data that already exists.")
}

// validateWorkspaceLinks refuses a vcs, ssh or agent id that does not resolve
// in the workspace organization.
func (a *API) validateWorkspaceLinks(w http.ResponseWriter, r *http.Request, orgID string, links map[string]*string) bool {
	if links["vcs"] != nil && !a.checkVCSRelationship(w, r, *links["vcs"], orgID) {
		return false
	}
	if links["ssh"] != nil && !a.checkSSHRelationship(w, r, *links["ssh"], orgID) {
		return false
	}
	if links["agent"] != nil && !a.checkAgentRelationship(w, r, *links["agent"], orgID) {
		return false
	}
	return true
}

func (a *API) relationshipLoadFailed(w http.ResponseWriter, operation string, err error, id, kind string) bool {
	if errors.Is(err, pgx.ErrNoRows) {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier "+id+" for "+kind)
		return false
	}
	if !a.writeStoreError(w, operation, err) {
		a.internalError(w, operation, err)
	}
	return false
}
