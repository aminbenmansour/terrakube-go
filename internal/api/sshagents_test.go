package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	sshOrgID   = "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	sshOtherID = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	sshKeyID   = "cccccccc-cccc-4ccc-8ccc-cccccccccccc"
	sshAgentID = "dddddddd-dddd-4ddd-8ddd-dddddddddddd"
	sshWSID    = "eeeeeeee-eeee-4eee-8eee-eeeeeeeeeeee"
)

// sshAgentStub answers only the ssh and agent methods; the embedded interface
// panics for anything else, which would mean the handler grew a dependency.
type sshAgentStub struct {
	store.Querier
	org        store.Organization
	keys       []store.Ssh
	agents     []store.Agent
	workspaces []store.Workspace
	createdKey *store.CreateSSHKeyParams
	updatedKey *store.UpdateSSHKeyParams
	deletedKey string
	deletedAgt string
	keyUses    int64
	agentUses  int64
}

func (s *sshAgentStub) GetOrganization(context.Context, string) (store.Organization, error) {
	return s.org, nil
}

func (s *sshAgentStub) ListSSHKeysByOrganization(context.Context, string) ([]store.Ssh, error) {
	return s.keys, nil
}

func (s *sshAgentStub) GetSSHKey(_ context.Context, id string) (store.Ssh, error) {
	for _, key := range s.keys {
		if key.ID == id {
			return key, nil
		}
	}
	return store.Ssh{}, pgx.ErrNoRows
}

func (s *sshAgentStub) CreateSSHKey(_ context.Context, arg store.CreateSSHKeyParams) (store.Ssh, error) {
	s.createdKey = &arg
	return store.Ssh{ID: arg.ID, Name: arg.Name, Description: arg.Description, SshType: arg.SshType, OrganizationID: arg.OrganizationID}, nil
}

func (s *sshAgentStub) UpdateSSHKey(_ context.Context, arg store.UpdateSSHKeyParams) (store.Ssh, error) {
	s.updatedKey = &arg
	return store.Ssh{ID: arg.ID, Name: arg.Name, Description: arg.Description, SshType: arg.SshType}, nil
}

func (s *sshAgentStub) CountWorkspacesWithSSHKey(context.Context, *string) (int64, error) {
	return s.keyUses, nil
}

func (s *sshAgentStub) CountWorkspacesWithAgent(context.Context, *string) (int64, error) {
	return s.agentUses, nil
}

func (s *sshAgentStub) DeleteSSHKey(_ context.Context, id string) error {
	s.deletedKey = id
	return nil
}

func (s *sshAgentStub) ListAgentsByOrganization(context.Context, string) ([]store.Agent, error) {
	return s.agents, nil
}

func (s *sshAgentStub) GetAgent(_ context.Context, id string) (store.Agent, error) {
	for _, agent := range s.agents {
		if agent.ID == id {
			return agent, nil
		}
	}
	return store.Agent{}, pgx.ErrNoRows
}

func (s *sshAgentStub) CreateAgent(_ context.Context, arg store.CreateAgentParams) (store.Agent, error) {
	return store.Agent{ID: arg.ID, Name: arg.Name, Url: arg.Url, Description: arg.Description, OrganizationID: arg.OrganizationID}, nil
}

func (s *sshAgentStub) DeleteAgent(_ context.Context, id string) error {
	s.deletedAgt = id
	return nil
}

func (s *sshAgentStub) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *sshAgentStub) UpdateWorkspace(_ context.Context, arg store.UpdateWorkspaceParams) (store.Workspace, error) {
	return store.Workspace{ID: arg.ID, SshID: arg.SshID, AgentID: arg.AgentID}, nil
}

func sshTestAPI(q store.Querier) http.Handler {
	return New(Deps{Store: q})
}

func TestSSHKeysJSONAPI(t *testing.T) {
	q := &sshAgentStub{
		org:  store.Organization{ID: sshOrgID},
		keys: []store.Ssh{{ID: sshKeyID, Name: strPtr("deploy"), PrivateKey: strPtr("secret"), OrganizationID: sshOrgID}},
	}
	h := sshTestAPI(q)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+sshOrgID+"/ssh", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "secret") {
		t.Errorf("list leaked the private key: %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"sshType"`) {
		t.Errorf("list missing sshType: %s", rec.Body.String())
	}

	body := `{"data":{"type":"ssh","attributes":{"name":"ci","description":"d","sshType":"rsa","privateKey":"PRIVATE"}}}`
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+sshOrgID+"/ssh", strings.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.createdKey == nil || q.createdKey.PrivateKey == nil || *q.createdKey.PrivateKey != "PRIVATE" {
		t.Errorf("private key not stored: %+v", q.createdKey)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/organization/"+sshOrgID+"/ssh/"+sshKeyID,
		strings.NewReader(`{"data":{"type":"ssh","id":"`+sshKeyID+`","attributes":{"name":"renamed"}}}`)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("patch status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.updatedKey == nil || q.updatedKey.Name == nil || *q.updatedKey.Name != "renamed" {
		t.Errorf("patch did not rename: %+v", q.updatedKey)
	}
	if q.updatedKey.PrivateKey == nil || *q.updatedKey.PrivateKey != "secret" {
		t.Errorf("patch dropped the stored key: %+v", q.updatedKey)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/organization/"+sshOrgID+"/ssh/"+sshKeyID, nil))
	if rec.Code != http.StatusNoContent || q.deletedKey != sshKeyID {
		t.Errorf("delete = %d, deleted %q", rec.Code, q.deletedKey)
	}

	// A key still attached to a workspace is refused, as upstream's foreign key
	// does, so a workspace cannot keep a dangling id.
	inUse := &sshAgentStub{
		org:     store.Organization{ID: sshOrgID},
		keys:    []store.Ssh{{ID: sshKeyID, OrganizationID: sshOrgID}},
		keyUses: 1,
	}
	rec = httptest.NewRecorder()
	sshTestAPI(inUse).ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/organization/"+sshOrgID+"/ssh/"+sshKeyID, nil))
	if rec.Code != http.StatusConflict {
		t.Errorf("in-use delete = %d, want 409; body = %s", rec.Code, rec.Body.String())
	}

	// A key from another organization is answered 404, like other resources.
	other := &sshAgentStub{org: store.Organization{ID: sshOtherID},
		keys: []store.Ssh{{ID: sshKeyID, OrganizationID: sshOrgID}}}
	rec = httptest.NewRecorder()
	sshTestAPI(other).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+sshOtherID+"/ssh/"+sshKeyID, nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("cross-org get = %d, want 404", rec.Code)
	}
}

func TestAgentsJSONAPI(t *testing.T) {
	q := &sshAgentStub{
		org:    store.Organization{ID: sshOrgID},
		agents: []store.Agent{{ID: sshAgentID, Name: strPtr("pool"), Url: strPtr("http://agent:8090"), OrganizationID: sshOrgID}},
	}
	h := sshTestAPI(q)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+sshOrgID+"/agent", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "http://agent:8090") {
		t.Fatalf("list = %d; body = %s", rec.Code, rec.Body.String())
	}

	body := `{"data":{"type":"agent","attributes":{"name":"ci","description":"d","url":"http://ci:8090"}}}`
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+sshOrgID+"/agent", strings.NewReader(body)))
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), "http://ci:8090") {
		t.Fatalf("create = %d; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/organization/"+sshOrgID+"/agent/"+sshAgentID, nil))
	if rec.Code != http.StatusNoContent || q.deletedAgt != sshAgentID {
		t.Errorf("delete = %d, deleted %q", rec.Code, q.deletedAgt)
	}

	inUse := &sshAgentStub{
		org:       store.Organization{ID: sshOrgID},
		agents:    []store.Agent{{ID: sshAgentID, OrganizationID: sshOrgID}},
		agentUses: 1,
	}
	rec = httptest.NewRecorder()
	sshTestAPI(inUse).ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/organization/"+sshOrgID+"/agent/"+sshAgentID, nil))
	if rec.Code != http.StatusConflict || inUse.deletedAgt != "" {
		t.Errorf("in-use delete = %d, deleted %q; body = %s", rec.Code, inUse.deletedAgt, rec.Body.String())
	}
}

func TestWorkspaceLinksAnSSHKeyOrAgent(t *testing.T) {
	q := &sshAgentStub{
		org:        store.Organization{ID: sshOrgID},
		keys:       []store.Ssh{{ID: sshKeyID, OrganizationID: sshOrgID}},
		agents:     []store.Agent{{ID: sshAgentID, OrganizationID: sshOrgID}},
		workspaces: []store.Workspace{{ID: sshWSID, OrganizationID: sshOrgID, Name: "ws"}},
	}
	h := sshTestAPI(q)

	for relation, id := range map[string]string{"ssh": sshKeyID, "agent": sshAgentID} {
		body := `{"data":{"type":"` + relation + `","id":"` + id + `"}}`
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch,
			"/organization/"+sshOrgID+"/workspace/"+sshWSID+"/relationships/"+relation, strings.NewReader(body)))
		if rec.Code != http.StatusNoContent {
			t.Errorf("%s link = %d; body = %s", relation, rec.Code, rec.Body.String())
		}
	}

	// A link from another organization is refused, so a workspace cannot point
	// at a key the caller cannot see.
	other := &sshAgentStub{
		org:        store.Organization{ID: sshOtherID},
		keys:       []store.Ssh{{ID: sshKeyID, OrganizationID: sshOrgID}},
		workspaces: []store.Workspace{{ID: sshWSID, OrganizationID: sshOtherID}},
	}
	rec := httptest.NewRecorder()
	sshTestAPI(other).ServeHTTP(rec, httptest.NewRequest(http.MethodPatch,
		"/organization/"+sshOtherID+"/workspace/"+sshWSID+"/relationships/ssh",
		strings.NewReader(`{"data":{"type":"ssh","id":"`+sshKeyID+`"}}`)))
	if rec.Code != http.StatusForbidden {
		t.Errorf("cross-org ssh link = %d, want 403", rec.Code)
	}

	// Creating a workspace with an unknown key fails before the row is written.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+sshOrgID+"/workspace",
		strings.NewReader(`{"data":{"type":"workspace","attributes":{"name":"ws2"},"relationships":{"ssh":{"data":{"type":"ssh","id":"`+sshOtherID+`"}}}}}`)))
	if rec.Code != http.StatusNotFound {
		t.Errorf("unknown ssh link = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
}
