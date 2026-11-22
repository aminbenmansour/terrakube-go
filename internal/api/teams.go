package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type teamAttrs struct {
	Name             *string `json:"name"`
	ManageWorkspace  *bool   `json:"manageWorkspace"`
	ManageModule     *bool   `json:"manageModule"`
	ManageProvider   *bool   `json:"manageProvider"`
	ManageVcs        *bool   `json:"manageVcs"`
	ManageTemplate   *bool   `json:"manageTemplate"`
	ManageState      *bool   `json:"manageState"`
	ManageCollection *bool   `json:"manageCollection"`
	ManageJob        *bool   `json:"manageJob"`
	ManagePolicies   *bool   `json:"managePolicies"`
	PlanJob          *bool   `json:"planJob"`
	ApproveJob       *bool   `json:"approveJob"`
	DestroyJob       *bool   `json:"destroyJob"`
	Role             *string `json:"role"`
}

func teamResource(team store.Team) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "team",
		ID:   team.ID,
		Attributes: map[string]any{
			"approveJob":       team.ApproveJob,
			"destroyJob":       team.DestroyJob,
			"manageCollection": team.ManageCollection,
			"manageJob":        team.ManageJob,
			"manageModule":     team.ManageModule,
			"managePolicies":   team.ManagePolicies,
			"manageProvider":   team.ManageProvider,
			"manageState":      team.ManageState,
			"manageTemplate":   team.ManageTemplate,
			"manageVcs":        team.ManageVcs,
			"manageWorkspace":  team.ManageWorkspace,
			"name":             team.Name,
			"planJob":          team.PlanJob,
			"role":             team.Role,
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", team.OrganizationID),
		},
	}
}

func (a *API) listTeams(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	teams, err := a.deps.Store.ListTeamsByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list teams", err) {
			a.internalError(w, "list teams", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(teams))
	for _, team := range teams {
		resources = append(resources, teamResource(team))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getTeam(w http.ResponseWriter, r *http.Request) {
	team, ok := a.teamInOrganization(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, teamResource(team), jsonapi.ParseFieldsets(r))
}

func (a *API) createTeam(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	attrs, _, ok := decodeTeamAttrs(w, r)
	if !ok {
		return
	}
	if attrs.Name == nil {
		writeInUseConflict(w)
		return
	}
	team, err := a.deps.Store.CreateTeam(r.Context(), store.CreateTeamParams{
		ID:               uuid.NewString(),
		Name:             *attrs.Name,
		ManageWorkspace:  boolValue(attrs.ManageWorkspace, false),
		ManageModule:     boolValue(attrs.ManageModule, false),
		ManageProvider:   boolValue(attrs.ManageProvider, false),
		ManageVcs:        boolValue(attrs.ManageVcs, false),
		ManageTemplate:   boolValue(attrs.ManageTemplate, false),
		ManageState:      boolValue(attrs.ManageState, false),
		ManageCollection: boolValue(attrs.ManageCollection, false),
		ManageJob:        boolValue(attrs.ManageJob, false),
		ManagePolicies:   boolValue(attrs.ManagePolicies, false),
		PlanJob:          boolValue(attrs.PlanJob, false),
		ApproveJob:       boolValue(attrs.ApproveJob, false),
		DestroyJob:       boolValue(attrs.DestroyJob, true),
		Role:             attrs.Role,
		OrganizationID:   orgID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create team", err) {
			a.internalError(w, "create team", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, teamResource(team), jsonapi.ParseFieldsets(r))
}

func (a *API) updateTeam(w http.ResponseWriter, r *http.Request) {
	team, ok := a.teamInOrganization(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodeTeamAttrs(w, r)
	if !ok {
		return
	}
	name := team.Name
	if jsonapi.HasAttribute(data, "name") {
		if attrs.Name == nil {
			writeInUseConflict(w)
			return
		}
		name = *attrs.Name
	}
	if _, err := a.deps.Store.UpdateTeam(r.Context(), store.UpdateTeamParams{
		ID:               team.ID,
		Name:             name,
		ManageWorkspace:  pickBool(data, "manageWorkspace", attrs.ManageWorkspace, team.ManageWorkspace),
		ManageModule:     pickBool(data, "manageModule", attrs.ManageModule, team.ManageModule),
		ManageProvider:   pickBool(data, "manageProvider", attrs.ManageProvider, team.ManageProvider),
		ManageVcs:        pickBool(data, "manageVcs", attrs.ManageVcs, team.ManageVcs),
		ManageTemplate:   pickBool(data, "manageTemplate", attrs.ManageTemplate, team.ManageTemplate),
		ManageState:      pickBool(data, "manageState", attrs.ManageState, team.ManageState),
		ManageCollection: pickBool(data, "manageCollection", attrs.ManageCollection, team.ManageCollection),
		ManageJob:        pickBool(data, "manageJob", attrs.ManageJob, team.ManageJob),
		ManagePolicies:   pickBool(data, "managePolicies", attrs.ManagePolicies, team.ManagePolicies),
		PlanJob:          pickBool(data, "planJob", attrs.PlanJob, team.PlanJob),
		ApproveJob:       pickBool(data, "approveJob", attrs.ApproveJob, team.ApproveJob),
		DestroyJob:       pickBool(data, "destroyJob", attrs.DestroyJob, team.DestroyJob),
		Role:             pickString(data, "role", attrs.Role, team.Role),
	}); err != nil {
		if !a.writeStoreError(w, "update team", err) {
			a.internalError(w, "update team", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteTeam(w http.ResponseWriter, r *http.Request) {
	team, ok := a.teamInOrganization(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteTeam(r.Context(), team.ID); err != nil {
		if !a.writeStoreError(w, "delete team", err) {
			a.internalError(w, "delete team", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) teamInOrganization(w http.ResponseWriter, r *http.Request) (store.Team, bool) {
	id := param(r, "teamID")
	if _, err := uuid.Parse(id); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
		return store.Team{}, false
	}
	team, err := a.deps.Store.GetTeam(r.Context(), id)
	if err != nil {
		if !a.writeStoreError(w, "get team", err) {
			a.internalError(w, "get team", err)
		}
		return store.Team{}, false
	}
	if team.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Team{}, false
	}
	return team, true
}

func decodeTeamAttrs(w http.ResponseWriter, r *http.Request) (teamAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return teamAttrs{}, nil, false
	}
	if !requireType(w, data, "team") {
		return teamAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[teamAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return teamAttrs{}, nil, false
	}
	return attrs, data, true
}
