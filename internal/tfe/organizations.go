package tfe

import (
	"net/http"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

func (a *API) ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("TFP-API-Version", "2.5")
	w.Header().Set("TFP-AppName", "Terrakube")
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusNotFound)
}

func (a *API) getOrganization(w http.ResponseWriter, r *http.Request) {
	org, ok := a.organizationByName(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, organizationResource(org))
}

func (a *API) getEntitlementSet(w http.ResponseWriter, r *http.Request) {
	org, ok := a.organizationByName(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, jsonapi.Resource{
		Type: "entitlement-sets",
		ID:   "org-" + org.Name,
		Attributes: map[string]any{
			"operations":              true,
			"private-module-registry": true,
			"sentinel":                false,
			"run-tasks":               false,
			"state-storage":           true,
			"teams":                   false,
			"vcs-integrations":        true,
			"usage-reporting":         false,
			"user-limit":              5,
			"self-serve-billing":      true,
			"audit-logging":           false,
			"agents":                  false,
			"sso":                     false,
		},
	})
}

func (a *API) getCapacity(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.organizationByName(w, r); !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, jsonapi.Resource{
		Type: "organization-capacity",
		Attributes: map[string]any{
			"pending": 0,
			"running": 0,
		},
	})
}

func organizationResource(org store.Organization) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "organizations",
		ID:   org.Name,
		Attributes: map[string]any{
			"name":        org.Name,
			"permissions": organizationPermissions(),
		},
	}
}

func organizationPermissions() map[string]bool {
	return map[string]bool{
		"can-update":                  true,
		"can-destroy":                 true,
		"can-access-via-teams":        false,
		"can-create-module":           false,
		"can-create-team":             false,
		"can-create-workspace":        true,
		"can-manage-users":            false,
		"can-manage-subscription":     false,
		"can-manage-sso":              false,
		"can-update-oauth":            false,
		"can-update-sentinel":         false,
		"can-update-ssh-keys":         false,
		"can-update-api-token":        false,
		"can-traverse":                false,
		"can-start-trial":             false,
		"can-update-agent-pools":      false,
		"can-manage-tags":             true,
		"can-manage-public-modules":   false,
		"can-manage-public-providers": false,
		"can-manage-run-tasks":        true,
		"can-read-run-tasks":          true,
		"can-create-provider":         false,
	}
}
