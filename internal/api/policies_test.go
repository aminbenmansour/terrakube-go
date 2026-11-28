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
	policyOrgID     = "51515151-5151-4151-8151-515151515151"
	policySetID     = "52525252-5252-4252-8252-525252525252"
	policyAttID     = "53535353-5353-4353-8353-535353535353"
	policyParamID   = "54545454-5454-4454-8454-545454545454"
	policyExemptID  = "55555555-5555-4555-8555-555555555555"
	policyWSID      = "56565656-5656-4656-8656-565656565656"
	policyProjectID = "57575757-5757-4757-8757-575757575757"
	policyTagID     = "58585858-5858-4858-8858-585858585858"
)

type policyStub struct {
	store.Querier
	org         store.Organization
	sets        []store.PolicySet
	attachments []store.PolicyAttachment
	parameters  []store.PolicySetParameter
	exemptions  []store.PolicyExemption
	workspaces  []store.Workspace
	projects    []store.Project
	tags        []store.Tag
	configs     []store.NotificationConfiguration

	createdSet    *store.CreatePolicySetParams
	updatedSet    *store.UpdatePolicySetParams
	createdAtt    *store.CreatePolicyAttachmentParams
	createdParam  *store.CreatePolicySetParameterParams
	createdExempt *store.CreatePolicyExemptionParams
	deleted       []string
}

func (s *policyStub) GetOrganization(_ context.Context, id string) (store.Organization, error) {
	if id != s.org.ID {
		return store.Organization{}, pgx.ErrNoRows
	}
	return s.org, nil
}

func (s *policyStub) ListPolicySets(context.Context) ([]store.PolicySet, error) {
	return s.sets, nil
}

func (s *policyStub) ListPolicySetsByOrganization(_ context.Context, id string) ([]store.PolicySet, error) {
	var out []store.PolicySet
	for _, set := range s.sets {
		if set.OrganizationID == id {
			out = append(out, set)
		}
	}
	return out, nil
}

func (s *policyStub) ListPolicySetsByOrganizations(_ context.Context, ids []string) ([]store.PolicySet, error) {
	var out []store.PolicySet
	for _, set := range s.sets {
		for _, id := range ids {
			if set.OrganizationID == id {
				out = append(out, set)
			}
		}
	}
	return out, nil
}

func (s *policyStub) GetPolicySet(_ context.Context, id string) (store.PolicySet, error) {
	for _, set := range s.sets {
		if set.ID == id {
			return set, nil
		}
	}
	return store.PolicySet{}, pgx.ErrNoRows
}

func (s *policyStub) CreatePolicySet(_ context.Context, arg store.CreatePolicySetParams) (store.PolicySet, error) {
	s.createdSet = &arg
	return store.PolicySet{
		ID: arg.ID, Name: arg.Name, Description: arg.Description,
		EnforcementLevel: arg.EnforcementLevel, ShadowEnforcementLevel: arg.ShadowEnforcementLevel,
		OverrideTeam: arg.OverrideTeam, Global: arg.Global, Repository: arg.Repository,
		Branch: arg.Branch, Folder: arg.Folder, VcsID: arg.VcsID,
		OrganizationID:              arg.OrganizationID,
		NotificationConfigurationID: arg.NotificationConfigurationID,
		OpaVersion:                  arg.OpaVersion, CreatedBy: arg.CreatedBy,
	}, nil
}

func (s *policyStub) UpdatePolicySet(_ context.Context, arg store.UpdatePolicySetParams) (store.PolicySet, error) {
	s.updatedSet = &arg
	return store.PolicySet{ID: arg.ID, Name: arg.Name}, nil
}

func (s *policyStub) DeletePolicySet(_ context.Context, id string) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *policyStub) ListPolicyAttachments(_ context.Context, id string) ([]store.PolicyAttachment, error) {
	var out []store.PolicyAttachment
	for _, attachment := range s.attachments {
		if attachment.PolicySetID == id {
			out = append(out, attachment)
		}
	}
	return out, nil
}

func (s *policyStub) GetPolicyAttachment(_ context.Context, id string) (store.PolicyAttachment, error) {
	for _, attachment := range s.attachments {
		if attachment.ID == id {
			return attachment, nil
		}
	}
	return store.PolicyAttachment{}, pgx.ErrNoRows
}

func (s *policyStub) CreatePolicyAttachment(_ context.Context, arg store.CreatePolicyAttachmentParams) (store.PolicyAttachment, error) {
	s.createdAtt = &arg
	return store.PolicyAttachment{
		ID: arg.ID, PolicySetID: arg.PolicySetID, WorkspaceID: arg.WorkspaceID,
		ProjectID: arg.ProjectID, TagID: arg.TagID, CreatedBy: arg.CreatedBy,
	}, nil
}

func (s *policyStub) UpdatePolicyAttachment(_ context.Context, arg store.UpdatePolicyAttachmentParams) (store.PolicyAttachment, error) {
	return store.PolicyAttachment{ID: arg.ID, WorkspaceID: arg.WorkspaceID, ProjectID: arg.ProjectID, TagID: arg.TagID}, nil
}

func (s *policyStub) DeletePolicyAttachment(_ context.Context, id string) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *policyStub) ListPolicySetParameters(_ context.Context, id string) ([]store.PolicySetParameter, error) {
	var out []store.PolicySetParameter
	for _, parameter := range s.parameters {
		if parameter.PolicySetID == id {
			out = append(out, parameter)
		}
	}
	return out, nil
}

func (s *policyStub) GetPolicySetParameter(_ context.Context, id string) (store.PolicySetParameter, error) {
	for _, parameter := range s.parameters {
		if parameter.ID == id {
			return parameter, nil
		}
	}
	return store.PolicySetParameter{}, pgx.ErrNoRows
}

func (s *policyStub) CreatePolicySetParameter(_ context.Context, arg store.CreatePolicySetParameterParams) (store.PolicySetParameter, error) {
	s.createdParam = &arg
	return store.PolicySetParameter{
		ID: arg.ID, PolicySetID: arg.PolicySetID, ParameterKey: arg.ParameterKey,
		ParameterValue: arg.ParameterValue, Description: arg.Description, CreatedBy: arg.CreatedBy,
	}, nil
}

func (s *policyStub) UpdatePolicySetParameter(_ context.Context, arg store.UpdatePolicySetParameterParams) (store.PolicySetParameter, error) {
	return store.PolicySetParameter{
		ID: arg.ID, ParameterKey: arg.ParameterKey, ParameterValue: arg.ParameterValue,
		Description: arg.Description,
	}, nil
}

func (s *policyStub) DeletePolicySetParameter(_ context.Context, id string) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *policyStub) ListPolicyExemptionsByOrganization(_ context.Context, id string) ([]store.PolicyExemption, error) {
	var out []store.PolicyExemption
	for _, exemption := range s.exemptions {
		if exemption.OrganizationID == id {
			out = append(out, exemption)
		}
	}
	return out, nil
}

func (s *policyStub) ListPolicyExemptionsByPolicySet(_ context.Context, id string) ([]store.PolicyExemption, error) {
	var out []store.PolicyExemption
	for _, exemption := range s.exemptions {
		if exemption.PolicySetID == id {
			out = append(out, exemption)
		}
	}
	return out, nil
}

func (s *policyStub) GetPolicyExemption(_ context.Context, id string) (store.PolicyExemption, error) {
	for _, exemption := range s.exemptions {
		if exemption.ID == id {
			return exemption, nil
		}
	}
	return store.PolicyExemption{}, pgx.ErrNoRows
}

func (s *policyStub) CreatePolicyExemption(_ context.Context, arg store.CreatePolicyExemptionParams) (store.PolicyExemption, error) {
	s.createdExempt = &arg
	return store.PolicyExemption{
		ID: arg.ID, PolicySetID: arg.PolicySetID, RuleID: arg.RuleID,
		WorkspaceID: arg.WorkspaceID, ProjectID: arg.ProjectID,
		TicketReference: arg.TicketReference, Justification: arg.Justification,
		ExpiresAt: arg.ExpiresAt, OrganizationID: arg.OrganizationID, CreatedBy: arg.CreatedBy,
	}, nil
}

func (s *policyStub) UpdatePolicyExemption(_ context.Context, arg store.UpdatePolicyExemptionParams) (store.PolicyExemption, error) {
	return store.PolicyExemption{
		ID: arg.ID, RuleID: arg.RuleID, WorkspaceID: arg.WorkspaceID, ProjectID: arg.ProjectID,
		TicketReference: arg.TicketReference, Justification: arg.Justification, ExpiresAt: arg.ExpiresAt,
	}, nil
}

func (s *policyStub) DeletePolicyExemption(_ context.Context, id string) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *policyStub) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *policyStub) GetProject(_ context.Context, id string) (store.Project, error) {
	for _, project := range s.projects {
		if project.ID == id {
			return project, nil
		}
	}
	return store.Project{}, pgx.ErrNoRows
}

func (s *policyStub) GetTag(_ context.Context, id string) (store.Tag, error) {
	for _, tag := range s.tags {
		if tag.ID == id {
			return tag, nil
		}
	}
	return store.Tag{}, pgx.ErrNoRows
}

func (s *policyStub) GetNotificationConfiguration(_ context.Context, id string) (store.NotificationConfiguration, error) {
	for _, config := range s.configs {
		if config.ID == id {
			return config, nil
		}
	}
	return store.NotificationConfiguration{}, pgx.ErrNoRows
}

func policyTestStore() *policyStub {
	org := store.Organization{ID: policyOrgID, Name: "policies"}
	return &policyStub{
		org:         org,
		sets:        []store.PolicySet{{ID: policySetID, Name: "capture-policy", OrganizationID: policyOrgID, EnforcementLevel: "SOFT_MANDATORY", Global: false}},
		attachments: []store.PolicyAttachment{{ID: policyAttID, PolicySetID: policySetID, WorkspaceID: strPtr(policyWSID)}},
		parameters:  []store.PolicySetParameter{{ID: policyParamID, PolicySetID: policySetID, ParameterKey: "max_length", ParameterValue: "8"}},
		exemptions:  []store.PolicyExemption{{ID: policyExemptID, PolicySetID: policySetID, OrganizationID: policyOrgID, RuleID: "rule-1", TicketReference: "SEC-1", Justification: "test"}},
		workspaces:  []store.Workspace{{ID: policyWSID, OrganizationID: policyOrgID, Name: "ws"}},
		projects:    []store.Project{{ID: policyProjectID, OrganizationID: policyOrgID, Name: "project"}},
		tags:        []store.Tag{{ID: policyTagID, OrganizationID: policyOrgID, Name: strPtr("tag")}},
	}
}

func policyServe(t *testing.T, q *policyStub, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	handler := New(Deps{Store: q})
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestPolicySetRoutes(t *testing.T) {
	q := policyTestStore()
	base := "/policy_set"

	body := `{"data":{"type":"policy_set","attributes":{"name":"capture-policy","description":"capture","enforcementLevel":"SOFT_MANDATORY","shadowEnforcementLevel":null,"overrideTeam":"TERRAKUBE_ADMIN","global":false,"repository":"https://github.com/example/policies","branch":"main","folder":"/","opaVersion":"0.68.0"},"relationships":{"organization":{"data":{"type":"organization","id":"` + policyOrgID + `"}}}}}`
	rec := policyServe(t, q, http.MethodPost, base, body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create policy set = %d/%s", rec.Code, rec.Body.String())
	}
	if q.createdSet == nil || q.createdSet.Name != "capture-policy" || q.createdSet.OrganizationID != policyOrgID {
		t.Fatalf("created policy set = %+v", q.createdSet)
	}
	if !strings.Contains(rec.Body.String(), `"type":"policy_set"`) {
		t.Fatalf("create body = %s", rec.Body.String())
	}

	rec = policyServe(t, q, http.MethodPost, base, `{"data":{"type":"policy_set","attributes":{"enforcementLevel":"ADVISORY"}}}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("missing name = %d/%s, want 409", rec.Code, rec.Body.String())
	}
	rec = policyServe(t, q, http.MethodPost, base, `{"data":{"type":"policy_set","attributes":{"name":"x","enforcementLevel":"NOPE"}}}`)
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "Invalid value: NOPE") {
		t.Fatalf("invalid enforcement = %d/%s", rec.Code, rec.Body.String())
	}

	rec = policyServe(t, q, http.MethodGet, base, "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"id":"`+policySetID+`"`) {
		t.Fatalf("list policy sets = %d/%s", rec.Code, rec.Body.String())
	}

	rec = policyServe(t, q, http.MethodGet, base+"/"+policySetID, "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"type":"policy_set"`) {
		t.Fatalf("get policy set = %d/%s", rec.Code, rec.Body.String())
	}

	rec = policyServe(t, q, http.MethodPatch, base+"/"+policySetID, `{"data":{"type":"policy_set","id":"`+policySetID+`","attributes":{"description":"after"}}}`)
	if rec.Code != http.StatusNoContent || q.updatedSet == nil || q.updatedSet.ID != policySetID {
		t.Fatalf("patch policy set = %d/%v", rec.Code, q.updatedSet)
	}

	rec = policyServe(t, q, http.MethodDelete, base+"/"+policySetID, "")
	if rec.Code != http.StatusNoContent || len(q.deleted) != 1 || q.deleted[0] != policySetID {
		t.Fatalf("delete policy set = %d/%v", rec.Code, q.deleted)
	}
}

func TestPolicyAttachmentAndParameterRoutes(t *testing.T) {
	q := policyTestStore()
	base := "/policy_set/" + policySetID

	body := `{"data":{"type":"policy_attachment","relationships":{"workspace":{"data":{"type":"workspace","id":"` + policyWSID + `"}}}}}`
	rec := policyServe(t, q, http.MethodPost, base+"/attachments", body)
	if rec.Code != http.StatusCreated || q.createdAtt == nil || q.createdAtt.WorkspaceID == nil || *q.createdAtt.WorkspaceID != policyWSID {
		t.Fatalf("create attachment = %d/%v", rec.Code, q.createdAtt)
	}

	rec = policyServe(t, q, http.MethodGet, base+"/attachments", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"id":"`+policyAttID+`"`) {
		t.Fatalf("list attachments = %d/%s", rec.Code, rec.Body.String())
	}
	rec = policyServe(t, q, http.MethodGet, base+"/attachments/"+policyAttID, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("get attachment = %d/%s", rec.Code, rec.Body.String())
	}
	rec = policyServe(t, q, http.MethodPatch, base+"/attachments/"+policyAttID,
		`{"data":{"type":"policy_attachment","id":"`+policyAttID+`","relationships":{"tag":{"data":{"type":"tag","id":"`+policyTagID+`"}}}}}`)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("patch attachment = %d/%s", rec.Code, rec.Body.String())
	}
	rec = policyServe(t, q, http.MethodDelete, base+"/attachments/"+policyAttID, "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete attachment = %d/%s", rec.Code, rec.Body.String())
	}

	rec = policyServe(t, q, http.MethodPost, base+"/parameters", `{"data":{"type":"policy_set_parameter","attributes":{"value":"8"}}}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("missing parameter key = %d/%s, want 409", rec.Code, rec.Body.String())
	}
	rec = policyServe(t, q, http.MethodPost, base+"/parameters",
		`{"data":{"type":"policy_set_parameter","attributes":{"key":"max_length","value":"8","description":"maximum"}}}`)
	if rec.Code != http.StatusCreated || q.createdParam == nil || q.createdParam.ParameterKey != "max_length" {
		t.Fatalf("create parameter = %d/%v", rec.Code, q.createdParam)
	}
	rec = policyServe(t, q, http.MethodGet, base+"/parameters", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"key":"max_length"`) {
		t.Fatalf("list parameters = %d/%s", rec.Code, rec.Body.String())
	}
	rec = policyServe(t, q, http.MethodGet, base+"/parameters/"+policyParamID, "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"key":"max_length"`) {
		t.Fatalf("get parameter = %d/%s", rec.Code, rec.Body.String())
	}
	rec = policyServe(t, q, http.MethodPatch, base+"/parameters/"+policyParamID,
		`{"data":{"type":"policy_set_parameter","id":"`+policyParamID+`","attributes":{"value":"12"}}}`)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("patch parameter = %d/%s", rec.Code, rec.Body.String())
	}
	rec = policyServe(t, q, http.MethodDelete, base+"/parameters/"+policyParamID, "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete parameter = %d/%s", rec.Code, rec.Body.String())
	}
}

func TestPolicyExemptionRoutes(t *testing.T) {
	q := policyTestStore()

	rec := policyServe(t, q, http.MethodGet, "/organization/"+policyOrgID+"/policyExemption", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"id":"`+policyExemptID+`"`) {
		t.Fatalf("list exemptions = %d/%s", rec.Code, rec.Body.String())
	}

	rec = policyServe(t, q, http.MethodPost, "/policy_exemption",
		`{"data":{"type":"policy_exemption","attributes":{"ticketReference":"SEC-2","justification":"x"},"relationships":{"organization":{"data":{"type":"organization","id":"`+policyOrgID+`"}},"policySet":{"data":{"type":"policy_set","id":"`+policySetID+`"}}}}}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("missing rule id = %d/%s, want 409", rec.Code, rec.Body.String())
	}

	body := `{"data":{"type":"policy_exemption","attributes":{"ruleId":"rule-1","ticketReference":"SEC-1","justification":"test","expiresAt":"2030-01-01T00:00:00Z"},"relationships":{"organization":{"data":{"type":"organization","id":"` + policyOrgID + `"}},"policySet":{"data":{"type":"policy_set","id":"` + policySetID + `"}},"workspace":{"data":{"type":"workspace","id":"` + policyWSID + `"}}}}}`
	rec = policyServe(t, q, http.MethodPost, "/policy_exemption", body)
	if rec.Code != http.StatusCreated || q.createdExempt == nil || q.createdExempt.RuleID != "rule-1" || !q.createdExempt.ExpiresAt.Valid {
		t.Fatalf("create exemption = %d/%v", rec.Code, q.createdExempt)
	}

	rec = policyServe(t, q, http.MethodGet, "/policy_exemption/"+policyExemptID, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("get exemption = %d/%s", rec.Code, rec.Body.String())
	}
	rec = policyServe(t, q, http.MethodPatch, "/policy_exemption/"+policyExemptID,
		`{"data":{"type":"policy_exemption","id":"`+policyExemptID+`","attributes":{"justification":"after"}}}`)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("patch exemption = %d/%s", rec.Code, rec.Body.String())
	}
	rec = policyServe(t, q, http.MethodDelete, "/policy_exemption/"+policyExemptID, "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete exemption = %d/%s", rec.Code, rec.Body.String())
	}
}

// policyCaptureStore stands in for the rows the policy captures were recorded
// against. Ids come from the captures, so re-recording does not need edits here.
func policyCaptureStore(captures map[string]upstreamCapture) *policyStub {
	setID := captureID(captures, "policy-set-create")
	attachmentID := captureID(captures, "policy-set-attachment-create")
	parameterID := captureID(captures, "policy-set-parameter-create")
	exemptionID := captureID(captures, "policy-exemption-create")
	workspaceID := captureRelationshipID(captures, "policy-set-attachment-create", "workspace")
	deletedAttachmentID := capturePathLastID(captures["policy-set-attachment-delete"].path)
	tagID := "0d7e4a6a-560e-40f1-a6dd-c7433a04f088"
	org := store.Organization{ID: capturedRelOrg, Name: "simple-governance"}
	q := &policyStub{
		org: org,
		sets: []store.PolicySet{{
			ID: setID, Name: "capture-policy", EnforcementLevel: "SOFT_MANDATORY",
			OrganizationID: capturedRelOrg, Branch: strPtr("main"), Folder: strPtr("/"),
		}},
		attachments: []store.PolicyAttachment{
			{ID: attachmentID, PolicySetID: setID, WorkspaceID: &workspaceID},
			{ID: deletedAttachmentID, PolicySetID: setID, TagID: &tagID},
		},
		parameters: []store.PolicySetParameter{{
			ID: parameterID, PolicySetID: setID, ParameterKey: "max_length", ParameterValue: "8",
			Description: strPtr("maximum password length"),
		}},
		exemptions: []store.PolicyExemption{{
			ID: exemptionID, PolicySetID: setID, OrganizationID: capturedRelOrg,
			WorkspaceID: &workspaceID, RuleID: "rule-1", TicketReference: "CAPTURE-1",
			Justification: "testing",
		}},
		workspaces: []store.Workspace{{ID: workspaceID, OrganizationID: capturedRelOrg, Name: "sample_cosmosdb"}},
		tags:       []store.Tag{{ID: tagID, OrganizationID: capturedRelOrg, Name: strPtr("staging")}},
	}
	// The single-parameter capture was recorded against its own set, so both
	// rows come from that path.
	if capture, ok := captures["policy-set-parameter-get"]; ok {
		parts := strings.Split(strings.Trim(capture.path, "/"), "/")
		setID, parameterID := parts[len(parts)-3], parts[len(parts)-1]
		q.sets = append(q.sets, store.PolicySet{
			ID: setID, Name: "capture-policy", EnforcementLevel: "SOFT_MANDATORY",
			OrganizationID: capturedRelOrg,
		})
		q.parameters = append(q.parameters, store.PolicySetParameter{
			ID: parameterID, PolicySetID: setID, ParameterKey: "max_length",
			ParameterValue: "8", Description: strPtr("maximum password length"),
		})
	}
	return q
}

func emptyPolicyStore() *policyStub {
	return &policyStub{org: store.Organization{ID: capturedRelOrg, Name: "simple-governance"}}
}

// capturePathLastID reads the id at the end of a captured path.
func capturePathLastID(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}

// captureRelationshipID reads a to-one id from a captured request or response.
func captureRelationshipID(captures map[string]upstreamCapture, name, relationship string) string {
	var doc struct {
		Data struct {
			Relationships map[string]struct {
				Data struct {
					ID string `json:"id"`
				} `json:"data"`
			} `json:"relationships"`
		} `json:"data"`
	}
	_ = json.Unmarshal(captures[name].body, &doc)
	return doc.Data.Relationships[relationship].Data.ID
}
