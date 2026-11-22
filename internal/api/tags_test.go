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
	tagOrgID   = "12121212-1212-4212-8212-121212121212"
	tagOtherID = "13131313-1313-4313-8313-131313131313"
	tagID      = "14141414-1414-4414-8414-141414141414"
	tagBindID  = "15151515-1515-4515-8515-151515151515"
	tagWSID    = "16161616-1616-4616-8616-161616161616"
	tagOtherWS = "17171717-1717-4717-8717-171717171717"
)

// tagStub answers only the tag methods; the embedded interface panics for
// anything else, which would mean the handler grew a dependency.
type tagStub struct {
	store.Querier
	org        store.Organization
	tags       []store.Tag
	bindings   []store.Workspacetag
	workspaces []store.Workspace
	created    *store.CreateWorkspaceTagParams
	deleted    string
}

func (s *tagStub) GetOrganization(context.Context, string) (store.Organization, error) {
	return s.org, nil
}

func (s *tagStub) ListTagsByOrganization(context.Context, string) ([]store.Tag, error) {
	return s.tags, nil
}

func (s *tagStub) GetTag(_ context.Context, id string) (store.Tag, error) {
	for _, tag := range s.tags {
		if tag.ID == id {
			return tag, nil
		}
	}
	return store.Tag{}, pgx.ErrNoRows
}

func (s *tagStub) CreateTag(_ context.Context, arg store.CreateTagParams) (store.Tag, error) {
	return store.Tag{ID: arg.ID, Name: arg.Name, OrganizationID: arg.OrganizationID}, nil
}

func (s *tagStub) UpdateTag(_ context.Context, arg store.UpdateTagParams) (store.Tag, error) {
	for i := range s.tags {
		if s.tags[i].ID == arg.ID {
			s.tags[i].Name = arg.Name
			return s.tags[i], nil
		}
	}
	return store.Tag{}, pgx.ErrNoRows
}

func (s *tagStub) DeleteTag(context.Context, string) error {
	return nil
}

func (s *tagStub) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *tagStub) ListWorkspaceTagsByWorkspace(context.Context, string) ([]store.Workspacetag, error) {
	return s.bindings, nil
}

func (s *tagStub) GetWorkspaceTag(_ context.Context, id string) (store.Workspacetag, error) {
	for _, binding := range s.bindings {
		if binding.ID == id {
			return binding, nil
		}
	}
	return store.Workspacetag{}, pgx.ErrNoRows
}

func (s *tagStub) CreateWorkspaceTag(_ context.Context, arg store.CreateWorkspaceTagParams) (store.Workspacetag, error) {
	s.created = &arg
	return store.Workspacetag{
		ID: arg.ID, TagID: arg.TagID, TagValue: arg.TagValue, WorkspaceID: arg.WorkspaceID,
	}, nil
}

func (s *tagStub) UpdateWorkspaceTagValue(_ context.Context, arg store.UpdateWorkspaceTagValueParams) (store.Workspacetag, error) {
	for i := range s.bindings {
		if s.bindings[i].ID == arg.ID {
			s.bindings[i].TagValue = arg.TagValue
			return s.bindings[i], nil
		}
	}
	return store.Workspacetag{}, pgx.ErrNoRows
}

func (s *tagStub) DeleteWorkspaceTag(_ context.Context, id string) error {
	s.deleted = id
	return nil
}

func newTagStub() *tagStub {
	return &tagStub{
		org: store.Organization{ID: tagOrgID},
		tags: []store.Tag{
			{ID: tagID, Name: strPtr("environment"), OrganizationID: tagOrgID},
			{ID: "18181818-1818-4818-8818-181818181818", Name: strPtr("team"), OrganizationID: tagOrgID},
		},
		bindings: []store.Workspacetag{{
			ID: tagBindID, TagID: tagID, TagValue: strPtr("prod"), WorkspaceID: tagWSID,
		}},
		workspaces: []store.Workspace{
			{ID: tagWSID, OrganizationID: tagOrgID, Name: "ws"},
			{ID: tagOtherWS, OrganizationID: tagOrgID, Name: "other"},
		},
	}
}

func TestTagsJSONAPI(t *testing.T) {
	q := newTagStub()
	h := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+tagOrgID+"/tag?filter[tag]=name==environment", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if strings.Count(rec.Body.String(), `"type":"tag"`) != 1 {
		t.Errorf("filter returned %s, want the one environment tag", rec.Body.String())
	}

	// A name is optional upstream, so a create without one stores null.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+tagOrgID+"/tag",
		strings.NewReader(`{"data":{"type":"tag","attributes":{}}}`)))
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"name":null`) {
		t.Fatalf("nameless create = %d; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/organization/"+tagOrgID+"/tag/"+tagID,
		strings.NewReader(`{"data":{"type":"tag","id":"`+tagID+`","attributes":{"name":"env"}}}`)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("patch = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.tags[0].Name == nil || *q.tags[0].Name != "env" {
		t.Errorf("patch did not rename: %v", q.tags[0].Name)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/organization/"+tagOrgID+"/tag/"+tagID, nil))
	if rec.Code != http.StatusNoContent {
		t.Errorf("delete = %d; body = %s", rec.Code, rec.Body.String())
	}

	other := newTagStub()
	other.org = store.Organization{ID: tagOtherID}
	rec = httptest.NewRecorder()
	newTestAPI(t, other).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+tagOtherID+"/tag/"+tagID, nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("cross-org get = %d, want 404", rec.Code)
	}
}

func TestWorkspaceTagsJSONAPI(t *testing.T) {
	q := newTagStub()
	h := newTestAPI(t, q)
	base := "/organization/" + tagOrgID + "/workspace/" + tagWSID + "/workspaceTag"

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, base, nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"tagId":"`+tagID+`"`) {
		t.Fatalf("list = %d; body = %s", rec.Code, rec.Body.String())
	}

	body := `{"data":{"type":"workspacetag","attributes":{"tagId":"` + tagID + `","value":"prod"}}}`
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, base, strings.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.created == nil || q.created.WorkspaceID != tagWSID || q.created.TagValue == nil || *q.created.TagValue != "prod" {
		t.Errorf("create params = %+v", q.created)
	}

	// A missing or unknown tag id is upstream's data-integrity 409.
	for _, body := range []string{
		`{"data":{"type":"workspacetag","attributes":{"value":"x"}}}`,
		`{"data":{"type":"workspacetag","attributes":{"tagId":"00000000-0000-4000-8000-000000000099","value":"x"}}}`,
	} {
		rec = httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, base, strings.NewReader(body)))
		if rec.Code != http.StatusConflict {
			t.Errorf("create %s = %d, want 409", body, rec.Code)
		}
	}

	// An explicit null clears the value, as the editor's empty input does.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, base+"/"+tagBindID,
		strings.NewReader(`{"data":{"type":"workspacetag","id":"`+tagBindID+`","attributes":{"value":null}}}`)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("patch = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.bindings[0].TagValue != nil {
		t.Errorf("value = %v, want cleared", *q.bindings[0].TagValue)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, base+"/"+tagBindID, nil))
	if rec.Code != http.StatusNoContent || q.deleted != tagBindID {
		t.Errorf("delete = %d, deleted %q", rec.Code, q.deleted)
	}

	// A binding of another workspace in the same organization is not reachable
	// through this workspace's path.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/organization/"+tagOrgID+"/workspace/"+tagOtherWS+"/workspaceTag/"+tagBindID, nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("cross-workspace delete = %d, want 404", rec.Code)
	}
}
