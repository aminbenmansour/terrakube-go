package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	teamOrgID     = "41414141-4141-4141-8141-414141414141"
	teamOtherID   = "42424242-4242-4242-8242-424242424242"
	teamID        = "43434343-4343-4343-8343-434343434343"
	teamForeignID = "44444444-4444-4444-8444-444444444444"
)

// teamStub answers only the team methods; the embedded interface panics for
// anything else, which would mean the handler grew a dependency.
type teamStub struct {
	store.Querier
	org     store.Organization
	teams   []store.Team
	created *store.CreateTeamParams
	updated *store.UpdateTeamParams
	deleted string
}

func (s *teamStub) GetOrganization(_ context.Context, id string) (store.Organization, error) {
	if id != s.org.ID {
		return store.Organization{}, pgx.ErrNoRows
	}
	return s.org, nil
}

func (s *teamStub) ListTeamsByOrganization(context.Context, string) ([]store.Team, error) {
	return s.teams, nil
}

func (s *teamStub) GetTeam(_ context.Context, id string) (store.Team, error) {
	for _, team := range s.teams {
		if team.ID == id {
			return team, nil
		}
	}
	return store.Team{}, pgx.ErrNoRows
}

func (s *teamStub) CreateTeam(_ context.Context, arg store.CreateTeamParams) (store.Team, error) {
	s.created = &arg
	return store.Team{
		ID: arg.ID, Name: arg.Name, ManageWorkspace: arg.ManageWorkspace,
		ManageModule: arg.ManageModule, ManageProvider: arg.ManageProvider,
		ManageVcs: arg.ManageVcs, ManageTemplate: arg.ManageTemplate,
		ManageState: arg.ManageState, ManageCollection: arg.ManageCollection,
		ManageJob: arg.ManageJob, ManagePolicies: arg.ManagePolicies,
		PlanJob: arg.PlanJob, ApproveJob: arg.ApproveJob, DestroyJob: arg.DestroyJob,
		Role: arg.Role, OrganizationID: arg.OrganizationID,
	}, nil
}

func (s *teamStub) UpdateTeam(_ context.Context, arg store.UpdateTeamParams) (store.Team, error) {
	s.updated = &arg
	return store.Team{ID: arg.ID, Name: arg.Name, Role: arg.Role}, nil
}

func (s *teamStub) DeleteTeam(_ context.Context, id string) error {
	s.deleted = id
	return nil
}

func newTeamStub() *teamStub {
	return &teamStub{
		org: store.Organization{ID: teamOrgID},
		teams: []store.Team{{
			ID: teamID, Name: "developers", ManageWorkspace: true, ManageJob: true,
			PlanJob: true, ApproveJob: true, DestroyJob: true, Role: strPtr("write"),
			OrganizationID: teamOrgID,
		}},
	}
}

func teamTestAPI(q store.Querier) http.Handler {
	return New(Deps{Store: q})
}

func TestTeamJSONAPI(t *testing.T) {
	q := newTeamStub()
	h := teamTestAPI(q)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+teamOrgID+"/team", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d; body = %s", rec.Code, rec.Body.String())
	}
	var list struct {
		Data []struct {
			Type       string         `json:"type"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Data) != 1 || list.Data[0].Type != "team" {
		t.Fatalf("list = %+v", list.Data)
	}
	if list.Data[0].Attributes["name"] != "developers" || list.Data[0].Attributes["manageWorkspace"] != true {
		t.Errorf("attributes = %+v", list.Data[0].Attributes)
	}

	// Upstream maps a missing primitive to false, except destroyJob which
	// defaults true because a client that does not know it must keep working.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+teamOrgID+"/team",
		strings.NewReader(`{"data":{"type":"team","attributes":{"name":"readers","role":"read"}}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.created == nil || q.created.Name != "readers" || q.created.ManageWorkspace ||
		!q.created.DestroyJob || q.created.Role == nil || *q.created.Role != "read" {
		t.Errorf("create params = %+v", q.created)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+teamOrgID+"/team",
		strings.NewReader(`{"data":{"type":"team","attributes":{"role":"read"}}}`)))
	if rec.Code != http.StatusConflict {
		t.Errorf("missing name = %d, want 409", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/organization/"+teamOrgID+"/team/"+teamID,
		strings.NewReader(`{"data":{"type":"team","id":"`+teamID+`","attributes":{"manageVcs":true,"role":null}}}`)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("patch = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.updated == nil || !q.updated.ManageVcs || q.updated.Role != nil ||
		q.updated.Name != "developers" || !q.updated.ManageWorkspace {
		t.Errorf("patch params = %+v", q.updated)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/organization/"+teamOrgID+"/team/"+teamID,
		strings.NewReader(`{"data":{"type":"team","id":"`+teamID+`","attributes":{"name":null}}}`)))
	if rec.Code != http.StatusConflict {
		t.Errorf("null name = %d, want 409", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/organization/"+teamOrgID+"/team/"+teamID, nil))
	if rec.Code != http.StatusNoContent || q.deleted != teamID {
		t.Fatalf("delete = %d, deleted %q", rec.Code, q.deleted)
	}

	other := newTeamStub()
	other.org = store.Organization{ID: teamOtherID}
	other.teams = append(other.teams, store.Team{ID: teamForeignID, Name: "foreign", OrganizationID: teamOtherID})
	rec = httptest.NewRecorder()
	teamTestAPI(other).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+teamOtherID+"/team/"+teamID, nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("cross-org get = %d, want 404", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+teamOrgID+"/team/not-a-uuid", nil))
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "Invalid value: not-a-uuid") {
		t.Errorf("bad id = %d; body = %s", rec.Code, rec.Body.String())
	}
}
