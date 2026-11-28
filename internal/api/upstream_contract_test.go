package api

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// upstreamCapturesDir holds responses recorded from a live upstream
// (azbuilder/api-server:2.34.3) by scripts/upstream-capture. They are the
// contract, so they are replayed here rather than trusted to a reading of the
// OpenAPI file.
const upstreamCapturesDir = "../../docs/upstream/rest"

type upstreamCapture struct {
	name   string
	status int
	path   string
	body   []byte
}

func (c upstreamCapture) resourceID(t *testing.T) string {
	t.Helper()
	var doc struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(c.body, &doc); err != nil || doc.Data.ID == "" {
		t.Fatalf("%s: no data.id to replay against: %v", c.name, err)
	}
	return doc.Data.ID
}

func loadUpstreamCaptures(t *testing.T) map[string]upstreamCapture {
	t.Helper()

	index, err := os.ReadFile(filepath.Join(upstreamCapturesDir, "index.tsv"))
	if err != nil {
		t.Fatalf("read captured index: %v", err)
	}

	captures := make(map[string]upstreamCapture)
	for _, line := range strings.Split(strings.TrimSpace(string(index)), "\n") {
		name, rest, ok := strings.Cut(line, "\t")
		rawStatus, path, ok2 := strings.Cut(rest, "\t")
		if !ok || !ok2 {
			t.Fatalf("malformed index line %q, want name<TAB>status<TAB>path", line)
		}
		status, err := strconv.Atoi(rawStatus)
		if err != nil {
			t.Fatalf("%s: status %q is not a number", name, rawStatus)
		}
		body, err := os.ReadFile(filepath.Join(upstreamCapturesDir, name+".json"))
		if err != nil {
			t.Fatalf("%s: read capture: %v", name, err)
		}
		captures[name] = upstreamCapture{name: name, status: status, path: path, body: body}
	}
	return captures
}

// replay is how this API answers one captured request, with stub data standing
// in for the database. target "resource" means "/organization/{captured id}".
type replay struct {
	target string
	// status overrides the captured status when we deliberately answer
	// differently; 0 means it must match.
	status  int
	querier store.Querier
	// authToken enables authentication, for the unauthenticated capture.
	authToken string
	// idFrom names another capture that carries the resource id, for paths
	// built against an organization described by a different request.
	idFrom string
	// method defaults to GET; body is a template in which {id} is replaced with
	// the captured resource id (PATCH cases).
	method string
	body   string
	// accept overrides both media-type headers, for the atomic extension.
	accept string
}

var replayedCaptures = map[string]replay{
	"organization-list": {
		target:  "/organization",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-list-sorted": {
		target:  "/organization?sort=-name",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-list-filtered": {
		target:  "/organization?filter[organization]=name==azure",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-get": {
		target: "resource",
	},
	"organization-get-fields": {
		target:  "resource-query",
		querier: &stubQuerier{org: store.Organization{ID: capturedOrgID, Name: "azure", Description: strPtr("azure sample organization")}},
	},
	"organization-get-fields-empty": {
		target:  "resource-query",
		querier: &stubQuerier{org: store.Organization{ID: capturedOrgID, Name: "azure"}},
	},
	"organization-get-fields-other": {
		target:  "resource-query",
		querier: &stubQuerier{org: store.Organization{ID: capturedOrgID, Name: "azure"}},
	},
	"organization-list-fields": {
		target:  "/organization?fields[organization]=name",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-projects": {
		target:  "resource-projects",
		idFrom:  "organization-get",
		querier: &stubQuerier{},
	},
	"organization-missing": {
		target:  "/organization/00000000-0000-0000-0000-000000000000",
		querier: &stubQuerier{err: pgx.ErrNoRows},
	},
	"organization-bad-id": {
		target: "/organization/not-a-uuid",
	},
	"organization-no-token": {
		target:    "/organization",
		authToken: "secret",
	},
	"organization-job-list": {
		target:  "/organization/" + capturedOrgID + "/job",
		querier: &stubQuerier{jobs: captureJobs()},
	},
	"organization-job-get": {
		target:  "job",
		querier: &stubQuerier{jobs: captureJobs()},
	},
	"organization-job-include": {
		target:  "job-include",
		querier: &stubQuerier{jobs: captureJobs(), steps: captureSteps(), workspaces: captureWorkspaces()},
	},
	"organization-job-patch": {
		target:  "job",
		idFrom:  "organization-job-get",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"job","id":"{id}","attributes":{"status":"cancelled"}}}`,
		querier: &stubQuerier{jobs: captureJobs()},
	},
	"organization-job-patch-invalid": {
		target:  "job",
		idFrom:  "organization-job-get",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"job","id":"{id}","attributes":{"status":"nonsense"}}}`,
		querier: &stubQuerier{jobs: captureJobs()},
	},
	"workspace-patch-attribute": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"workspace","id":"` + capturedRelWS + `","attributes":{"description":"captured"}}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-patch-project": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    workspaceProjectBody(capturedRelP1),
		querier: captureWorkspaceStore(),
	},
	"workspace-patch-project-moved": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    workspaceProjectBody(capturedRelP2),
		querier: captureWorkspaceStore(),
	},
	"workspace-patch-project-foreign": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    workspaceProjectBody(capturedRelForeign),
		querier: captureWorkspaceStore(),
	},
	"workspace-patch-project-unknown": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    workspaceProjectBody(capturedRelUnknown),
		querier: captureWorkspaceStore(),
	},
	"workspace-patch-project-clear": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"workspace","id":"` + capturedRelWS + `","relationships":{"project":{"data":null}}}}`,
		querier: captureWorkspaceStore(),
	},
	// The recording was made against a project of another organization; which
	// id it was is not part of the contract, so the replay names its own.
	"workspace-create-project-foreign": {
		target:  "workspace-collection",
		method:  http.MethodPost,
		body:    `{"data":{"type":"workspace","attributes":{"name":"from-vcs","source":"https://github.com/terrakube-io/terrakube.git","branch":"main","terraformVersion":"1.9.0"},"relationships":{"project":{"data":{"type":"project","id":"` + capturedRelForeign + `"}}}}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-project-set": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"project","id":"` + capturedRelP1 + `"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-project-clear": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-project-foreign": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"project","id":"` + capturedRelForeign + `"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-project-unknown": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"project","id":"` + capturedRelUnknown + `"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-vcs-set": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"vcs","id":"aaaaaaaa-0000-4000-8000-000000000001"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-access-list": {
		target:  "resource-query",
		querier: &stubQuerier{workspaces: captureGrantWorkspaces()},
	},
	"workspace-access-create": {
		target:  "resource-query",
		method:  http.MethodPost,
		body:    `{"data":{"type":"access","attributes":{"name":"capture-group","manageState":true,"manageWorkspace":true,"planJob":true,"approveJob":false,"role":"custom"}}}`,
		querier: captureGrantStore(),
	},
	"workspace-access-get": {
		target:  "resource-query",
		querier: captureGrantStore(),
	},
	"workspace-access-patch": {
		target:  "resource-query",
		idFrom:  "workspace-access-get",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"access","id":"{id}","attributes":{"role":"read"}}}`,
		querier: captureGrantStore(),
	},
	"workspace-access-get-after-patch": {
		target:  "resource-query",
		querier: captureGrantStore(),
	},
	"workspace-access-delete": {
		target:  "resource-query",
		method:  http.MethodDelete,
		querier: captureGrantStore(),
	},
	"workspace-access-get-after-delete": {
		target:  "resource-query",
		querier: &stubQuerier{workspaces: captureGrantWorkspaces(), err: pgx.ErrNoRows},
	},
	"project-access-list": {
		target:  "resource-query",
		querier: &stubQuerier{projects: captureGrantProjects()},
	},
	"project-access-create": {
		target:  "resource-query",
		method:  http.MethodPost,
		body:    `{"data":{"type":"project_access","attributes":{"name":"capture-group","manageState":true,"manageWorkspace":true,"planJob":false,"approveJob":false,"role":"custom"}}}`,
		querier: captureGrantStore(),
	},
	"project-access-get": {
		target:  "resource-query",
		querier: captureGrantStore(),
	},
	"project-access-patch": {
		target:  "resource-query",
		idFrom:  "project-access-get",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"project_access","id":"{id}","attributes":{"role":"plan"}}}`,
		querier: captureGrantStore(),
	},
	"project-access-delete": {
		target:  "resource-query",
		method:  http.MethodDelete,
		querier: captureGrantStore(),
	},
	"project-access-get-after-delete": {
		target:  "resource-query",
		querier: &stubQuerier{projects: captureGrantProjects(), err: pgx.ErrNoRows},
	},
	"workspace-relationship-vcs-clear": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-ssh-set": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"ssh","id":"aaaaaaaa-0000-4000-8000-000000000002"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-ssh-clear": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-agent-set": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"agent","id":"aaaaaaaa-0000-4000-8000-000000000003"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-agent-clear": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-unknown": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"organization-atomic-create": {
		target:  "resource-query",
		method:  http.MethodPost,
		body:    `{"atomic:operations":[{"op":"add","href":"/organization/` + capturedRelOrg + `/workspace","data":{"type":"workspace","lid":"11111111-1111-4111-8111-111111111111","attributes":{"source":"https://github.com/terrakube-io/terrakube.git","folder":"/","name":"capture-atomic","terraformVersion":"1.9.0","branch":"main","iacType":"terraform","executionMode":"remote"}}}]}`,
		accept:  atomicMedia,
		querier: captureWorkspaceStore(),
	},
	"organization-atomic-update": {
		target:  "resource-query",
		method:  http.MethodPost,
		body:    `{"atomic:operations":[{"op":"update","href":"/organization/` + capturedRelOrg + `/workspace/` + capturedRelWS + `","data":{"type":"workspace","id":"` + capturedRelWS + `","attributes":{"description":"atomic-after"}}}]}`,
		accept:  atomicMedia,
		querier: captureWorkspaceStore(),
	},
	"organization-globalvar-list": {
		target:  "globalvar-collection",
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-get": {
		target:  "globalvar-item",
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-get-sensitive": {
		target:  "globalvar-secret",
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-create": {
		target:  "globalvar-collection",
		method:  http.MethodPost,
		body:    `{"data":{"type":"globalvar","attributes":{"key":"CAPTURE_NEW","value":"new","description":"capture","sensitive":false,"hcl":false,"category":"ENV"}}}`,
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-patch": {
		target:  "globalvar-item",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"globalvar","id":"` + capturedGvID + `","attributes":{"key":"CAPTURE_ONE","value":"after","description":"after","hcl":true,"category":"TERRAFORM"}}}`,
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-delete": {
		target:  "globalvar-secret",
		method:  http.MethodDelete,
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-invalid-category": {
		target:  "globalvar-collection",
		method:  http.MethodPost,
		body:    `{"data":{"type":"globalvar","attributes":{"key":"CAPTURE_BAD","value":"x","description":"d","sensitive":false,"hcl":false,"category":"env"}}}`,
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-missing-key": {
		target:  "globalvar-collection",
		method:  http.MethodPost,
		body:    `{"data":{"type":"globalvar","attributes":{"value":"x","description":"d","sensitive":false,"hcl":false,"category":"ENV"}}}`,
		querier: captureGlobalVarStore(),
	},
	"action-list": {
		target:  "resource-query",
		querier: captureActionStore(),
	},
	"action-list-filtered": {
		target:  "resource-query",
		querier: captureActionStore(),
	},
	"action-get": {
		target:  "resource-query",
		querier: captureActionStore(),
	},
	"action-fields": {
		target:  "resource-query",
		querier: captureActionStore(),
	},
	"action-missing": {
		target:  "resource-query",
		querier: &stubQuerier{},
	},
	"action-create": {
		target:  "resource-query",
		method:  http.MethodPost,
		body:    `{"data":{"type":"action","id":"capture.temp-action","attributes":{"name":"Capture Temp","label":"Temp","category":"General","description":"capture","type":"Workspace/Action","action":"dGVzdA==","displayCriteria":"[{\"filter\":\"true\"}]","version":"1.0.0","active":true}}}`,
		querier: captureActionStore(),
	},
	"action-patch": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"action","id":"capture.temp-action","attributes":{"label":"Temp after","active":false}}}`,
		querier: captureActionStore(),
	},
	"action-delete": {
		target:  "resource-query",
		method:  http.MethodDelete,
		querier: captureActionStore(),
	},
	"organization-collection-create": {
		target: "collection-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"collection","attributes":{"name":"capture-collection","description":"capture","priority":10}}}`,
	},
	"organization-collection-list": {
		target: "collection-query",
	},
	"organization-collection-get": {
		target: "collection-query",
	},
	"organization-collection-patch": {
		target: "collection-query",
		idFrom: "organization-collection-get",
		method: http.MethodPatch,
		body:   `{"data":{"type":"collection","id":"{id}","attributes":{"name":"capture-collection-after"}}}`,
	},
	"organization-collection-item-create": {
		target: "collection-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"item","attributes":{"key":"CAPTURE_KEY","value":"one","description":"capture","sensitive":false,"hcl":false,"category":"ENV"}}}`,
	},
	"organization-collection-item-list": {
		target: "collection-query",
	},
	"organization-collection-item-patch": {
		target: "collection-query",
		idFrom: "organization-collection-item-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"item","id":"{id}","attributes":{"value":"after"}}}`,
	},
	"organization-collection-item-invalid-category": {
		target: "collection-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"item","attributes":{"key":"CAPTURE_BAD","value":"x","category":"env"}}}`,
	},
	"organization-collection-item-missing-key": {
		target: "collection-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"item","attributes":{"value":"x","category":"ENV"}}}`,
	},
	"organization-collection-create-missing-name": {
		target: "collection-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"collection","attributes":{"priority":5}}}`,
	},
	"organization-collection-missing": {
		target: "collection-query",
	},
	"organization-collection-reference-create": {
		target: "collection-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"reference","attributes":{"description":"capture"},"relationships":{"workspace":{"data":{"type":"workspace","id":"` + capturedRelWS + `"}}}}}`,
	},
	"organization-collection-reference-list": {
		target: "collection-query",
	},
	"reference-collection": {
		target: "collection-query",
	},
	"organization-workspace-reference": {
		target: "collection-query",
	},
	"organization-collection-item-delete": {
		target: "collection-query",
		method: http.MethodDelete,
	},
	"organization-collection-reference-delete": {
		target: "collection-query",
		method: http.MethodDelete,
	},
	"organization-collection-delete": {
		target: "collection-query",
		method: http.MethodDelete,
	},
	"organization-ssh-list": {
		target: "ssh-agent-query",
	},
	"organization-ssh-get": {
		target: "ssh-agent-query",
	},
	"organization-ssh-create": {
		target: "ssh-agent-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"ssh","attributes":{"name":"capture-ssh-created","description":"capture","sshType":"rsa","privateKey":"-----BEGIN PRIVATE KEY-----\ncapture\n-----END PRIVATE KEY-----"}}}`,
	},
	"organization-ssh-patch": {
		target: "ssh-agent-query",
		idFrom: "organization-ssh-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"ssh","id":"{id}","attributes":{"description":"after"}}}`,
	},
	"organization-ssh-create-missing-name": {
		target: "ssh-agent-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"ssh","attributes":{"sshType":"rsa"}}}`,
	},
	"organization-ssh-missing": {
		target: "ssh-agent-query",
	},
	"organization-ssh-delete": {
		target: "ssh-agent-query",
		idFrom: "organization-ssh-create",
		method: http.MethodDelete,
	},
	"organization-ssh-delete-unnamed": {
		target: "ssh-agent-query",
		idFrom: "organization-ssh-create-missing-name",
		method: http.MethodDelete,
	},
	"organization-agent-list": {
		target: "ssh-agent-query",
	},
	"organization-agent-create": {
		target: "ssh-agent-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"agent","attributes":{"name":"capture-agent-created","description":"capture","url":"http://capture:8090"}}}`,
	},
	"organization-agent-delete": {
		target: "ssh-agent-query",
		idFrom: "organization-agent-create",
		method: http.MethodDelete,
	},
	"organization-ssh-link-in-use": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"ssh","id":"` + capturedRelSSH + `"}}`,
		querier: captureWorkspaceStore(),
	},
	"organization-ssh-delete-in-use": {
		target: "ssh-agent-query",
		method: http.MethodDelete,
	},
	"organization-ssh-unlink-in-use": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"organization-agent-link-in-use": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"agent","id":"` + capturedRelAgent + `"}}`,
		querier: captureWorkspaceStore(),
	},
	"organization-agent-delete-in-use": {
		target: "ssh-agent-query",
		method: http.MethodDelete,
	},
	"organization-vcs-list": {
		target: "vcs-query",
	},
	"organization-vcs-create": {
		target: "vcs-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"vcs","attributes":{"name":"capture-vcs-created","description":"capture","vcsType":"GITLAB","connectionType":"STANDALONE","clientId":"client","clientSecret":"secret","privateKey":"key","callback":"callback","endpoint":"https://gitlab.example.com","apiUrl":"https://gitlab.example.com/api/v4","redirectUrl":"https://ui.example.com/settings/vcs","status":"PENDING"}}}`,
	},
	"organization-vcs-get": {
		target: "vcs-query",
	},
	"organization-vcs-patch": {
		target: "vcs-query",
		idFrom: "organization-vcs-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"vcs","id":"{id}","attributes":{"description":"after","status":"COMPLETED"}}}`,
	},
	"organization-vcs-create-missing-required": {
		target: "vcs-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"vcs","attributes":{"name":"capture-vcs-incomplete"}}}`,
	},
	"organization-vcs-create-invalid-type": {
		target: "vcs-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"vcs","attributes":{"name":"capture-vcs-bad","description":"capture","vcsType":"NOPE"}}}`,
	},
	"organization-vcs-missing": {
		target: "vcs-query",
	},
	"workspace-relationship-vcs-link-in-use": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"vcs","id":"` + capturedRelVCS + `"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-include-vcs": {
		target:  "resource-query",
		querier: captureWorkspaceStore(),
	},
	"organization-vcs-delete-in-use": {
		target: "vcs-query",
		method: http.MethodDelete,
	},
	"workspace-get-after-vcs-delete": {
		target:  "workspace",
		querier: captureWorkspaceStore(),
	},
	"organization-vcs-delete": {
		target: "vcs-query",
		idFrom: "organization-vcs-create",
		method: http.MethodDelete,
	},
	"organization-vcs-get-after-delete": {
		target:  "vcs-query",
		querier: &stubQuerier{},
	},
	"organization-team-list": {
		target: "team-query",
	},
	"organization-team-create": {
		target: "team-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"team","attributes":{"name":"capture-team","manageWorkspace":true,"manageJob":true,"planJob":true,"approveJob":true,"role":"write"}}}`,
	},
	"organization-team-get": {
		target: "team-query",
	},
	"organization-team-patch": {
		target: "team-query",
		idFrom: "organization-team-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"team","id":"{id}","attributes":{"role":"admin","manageVcs":true,"destroyJob":false}}}`,
	},
	"organization-team-create-missing-name": {
		target: "team-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"team","attributes":{"role":"read"}}}`,
	},
	"organization-team-missing": {
		target: "team-query",
	},
	"organization-team-bad-id": {
		target: "team-query",
	},
	"organization-team-delete": {
		target: "team-query",
		idFrom: "organization-team-create",
		method: http.MethodDelete,
	},
	"organization-team-get-after-delete": {
		target:  "team-query",
		querier: &stubQuerier{},
	},
	"organization-notification-list": {
		target: "notification-query",
	},
	"organization-notification-create": {
		target: "notification-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"notification_configuration","attributes":{"name":"capture-notify","description":"capture","channelType":"WEBHOOK","destinationUrl":"https://example.com/hook","signingSecret":"capture-secret","active":true,"messageStyle":"SIMPLE"}}}`,
	},
	"organization-notification-get": {
		target: "notification-query",
	},
	"organization-notification-patch": {
		target: "notification-query",
		idFrom: "organization-notification-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"notification_configuration","id":"{id}","attributes":{"description":null,"active":false,"messageStyle":"DETAILED"}}}`,
	},
	"organization-notification-trigger-create": {
		target: "notification-query",
		idFrom: "organization-notification-create",
		method: http.MethodPost,
		body:   `{"data":{"type":"notification_trigger","attributes":{"jobStatus":"failed"}}}`,
	},
	"organization-notification-trigger-missing": {
		target: "notification-query",
		idFrom: "organization-notification-create",
		method: http.MethodPost,
		body:   `{"data":{"type":"notification_trigger","attributes":{}}}`,
	},
	"organization-notification-trigger-invalid": {
		target: "notification-query",
		idFrom: "organization-notification-create",
		method: http.MethodPost,
		body:   `{"data":{"type":"notification_trigger","attributes":{"jobStatus":"nope"}}}`,
	},
	"organization-notification-templates": {
		target: "notification-query",
		idFrom: "organization-notification-create",
		method: http.MethodPatch,
		body:   `{"data":[{"type":"template","id":"` + capturedNotifyTemplate + `"}]}`,
	},
	"organization-notification-templates-unknown": {
		target: "notification-query",
		idFrom: "organization-notification-create",
		method: http.MethodPatch,
		body:   `{"data":[{"type":"template","id":"00000000-0000-4000-8000-000000000099"}]}`,
	},
	"organization-notification-workspace-create": {
		target: "notification-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"notification_configuration","attributes":{"name":"capture-notify-workspace","channelType":"TEAMS","destinationUrl":"https://teams.example.com/hook"}}}`,
	},
	"organization-notification-workspace-list": {
		target: "notification-query",
	},
	"organization-notification-create-missing-name": {
		target: "notification-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"notification_configuration","attributes":{"channelType":"WEBHOOK","destinationUrl":"https://example.com/x"}}}`,
	},
	"organization-notification-create-invalid": {
		target: "notification-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"notification_configuration","attributes":{"name":"capture-notify-bad","channelType":"NOPE","destinationUrl":"https://example.com/x"}}}`,
	},
	"organization-notification-missing": {
		target: "notification-query",
	},
	"organization-notification-bad-id": {
		target: "notification-query",
	},
	"organization-notification-trigger-delete": {
		target: "notification-query",
		idFrom: "organization-notification-trigger-create",
		method: http.MethodDelete,
	},
	"organization-notification-workspace-delete": {
		target: "notification-query",
		idFrom: "organization-notification-workspace-create",
		method: http.MethodDelete,
	},
	"organization-notification-delete": {
		target: "notification-query",
		idFrom: "organization-notification-create",
		method: http.MethodDelete,
	},
	"organization-notification-get-after-delete": {
		target:  "notification-query",
		querier: &stubQuerier{},
	},
	"organization-template-list": {
		target: "template-query",
	},
	"organization-template-create": {
		target: "template-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"template","attributes":{"name":"capture-template","description":"capture","version":"9.9.9","tcl":"ZmxvdzoK"}}}`,
	},
	"organization-template-get": {
		target: "template-query",
	},
	"organization-template-patch": {
		target: "template-query",
		idFrom: "organization-template-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"template","id":"{id}","attributes":{"description":"after","version":"10.0.0","tcl":"ZmxvdzoKICAtIHR5cGU6ICJ0ZXJyYWZvcm1QbGFuIg=="}}}`,
	},
	"organization-template-create-empty": {
		target: "template-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"template","attributes":{}}}`,
	},
	"organization-template-missing": {
		target: "template-query",
	},
	"organization-template-bad-id": {
		target: "template-query",
	},
	"organization-template-delete": {
		target: "template-query",
		idFrom: "organization-template-create",
		method: http.MethodDelete,
	},
	"organization-template-get-after-delete": {
		target:  "template-query",
		querier: &stubQuerier{},
	},
	"organization-template-delete-empty": {
		target: "template-query",
		idFrom: "organization-template-create-empty",
		method: http.MethodDelete,
	},
	"organization-provider-list": {
		target: "registry-query",
	},
	"organization-provider-create": {
		target: "registry-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"provider","attributes":{"name":"capture-provider","description":"capture","imported":true,"registryNamespace":"acme"}}}`,
	},
	"organization-provider-get": {
		target: "registry-query",
	},
	"organization-provider-get-include": {
		target: "registry-provider-query",
	},
	"organization-provider-patch": {
		target: "registry-query",
		idFrom: "organization-provider-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"provider","id":"{id}","attributes":{"description":"after"}}}`,
	},
	"organization-provider-version-create": {
		target: "registry-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"version","attributes":{"versionNumber":"1.2.3","protocols":"5.0"}}}`,
	},
	"organization-provider-version-patch": {
		target: "registry-query",
		idFrom: "organization-provider-version-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"version","id":"{id}","attributes":{"status":"deprecated","deprecationMessage":"use 2.x"}}}`,
	},
	"organization-provider-version-invalid-status": {
		target: "registry-query",
		idFrom: "organization-provider-version-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"version","id":"{id}","attributes":{"status":"nonsense"}}}`,
	},
	"organization-provider-version-missing": {
		target: "registry-provider-query",
	},
	"organization-provider-implementation-create": {
		target: "registry-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"implementation","attributes":{"os":"linux","arch":"amd64","filename":"terraform-provider-capture_1.2.3_linux_amd64.zip","downloadUrl":"https://example.com/p.zip","shasumsUrl":"https://example.com/SHA256SUMS","shasumsSignatureUrl":"https://example.com/SHA256SUMS.sig","shasum":"abc","keyId":"key","asciiArmor":"armor","trustSignature":"trust","source":"https://github.com/example/p","sourceUrl":"https://github.com/example/p"}}}`,
	},
	"organization-provider-get-populated": {
		target: "registry-query",
	},
	"organization-provider-implementation-delete": {
		target: "registry-query",
		method: http.MethodDelete,
	},
	"organization-provider-version-delete": {
		target: "registry-query",
		method: http.MethodDelete,
	},
	"organization-provider-delete": {
		target: "registry-query",
		method: http.MethodDelete,
	},
	"organization-provider-get-after-delete": {
		target:  "resource-query",
		querier: emptyRegistryStore(),
	},
	"organization-provider-missing": {
		target:  "resource-query",
		querier: emptyRegistryStore(),
	},
	"organization-provider-bad-id": {
		target:  "resource-query",
		querier: emptyRegistryStore(),
	},
	"organization-module-create": {
		target: "registry-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"module","attributes":{"name":"capture-module","description":"capture","provider":"aws","source":"https://github.com/example/terraform-aws-capture","folder":null,"tagPrefix":null}}}`,
	},
	"organization-module-include": {
		target: "registry-query",
	},
	"organization-module-list": {
		target: "registry-query",
	},
	"organization-module-list-fields": {
		target: "registry-query",
	},
	"organization-module-get": {
		target: "registry-query",
	},
	"organization-module-get-include": {
		target: "registry-query",
	},
	"organization-module-patch": {
		target: "registry-query",
		idFrom: "organization-module-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"module","id":"{id}","attributes":{"description":"after"}}}`,
	},
	"organization-module-version-patch": {
		target: "registry-query",
		idFrom: "organization-module-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"module_version","id":"bbbbbbbb-0000-4000-8000-000000000001","attributes":{"status":"deprecated","deprecationMessage":"use 2.x"}}}`,
	},
	"organization-module-version-invalid-status": {
		target: "registry-query",
		idFrom: "organization-module-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"module_version","id":"bbbbbbbb-0000-4000-8000-000000000001","attributes":{"status":"nonsense"}}}`,
	},
	"organization-module-missing": {
		target:  "resource-query",
		querier: emptyRegistryStore(),
	},
	"organization-module-bad-id": {
		target:  "resource-query",
		querier: emptyRegistryStore(),
	},
	"organization-module-delete": {
		target: "registry-query",
		method: http.MethodDelete,
	},
	"organization-module-get-after-delete": {
		target:  "resource-query",
		querier: emptyRegistryStore(),
	},
	"policy-set-list": {
		target: "policy-query",
	},
	"policy-set-list-include": {
		target: "policy-query",
	},
	"policy-set-create": {
		target: "policy-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"policy_set","attributes":{"name":"capture-policy","description":"capture","enforcementLevel":"SOFT_MANDATORY","shadowEnforcementLevel":null,"overrideTeam":"TERRAKUBE_ADMIN","global":false,"repository":"https://github.com/example/policies","branch":"main","folder":"/","opaVersion":"0.68.0"},"relationships":{"organization":{"data":{"type":"organization","id":"` + capturedRelOrg + `"}}}}}`,
	},
	"policy-set-get": {
		target: "policy-query",
	},
	"policy-set-patch": {
		target: "policy-query",
		idFrom: "policy-set-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"policy_set","id":"{id}","attributes":{"description":"after"}}}`,
	},
	"policy-set-delete": {
		target: "policy-query",
		idFrom: "policy-set-create",
		method: http.MethodDelete,
	},
	"policy-set-get-after-delete": {
		target:  "resource-query",
		querier: emptyPolicyStore(),
	},
	"policy-set-missing": {
		target:  "resource-query",
		querier: emptyPolicyStore(),
	},
	"policy-set-bad-id": {
		target:  "resource-query",
		querier: emptyPolicyStore(),
	},
	"policy-set-missing-name": {
		target: "policy-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"policy_set","attributes":{"enforcementLevel":"ADVISORY"}}}`,
	},
	"policy-set-bad-enforcement": {
		target: "policy-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"policy_set","attributes":{"name":"capture-policy-bad","enforcementLevel":"NOPE"}}}`,
	},
	"policy-set-attachment-create": {
		target: "policy-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"policy_attachment","relationships":{"workspace":{"data":{"type":"workspace","id":"0beb7db4-b073-4bd4-9a98-94fe7e2d9939"}}}}}`,
	},
	"policy-set-attachment-list": {
		target: "policy-query",
	},
	"policy-set-attachment-get": {
		target: "policy-query",
	},
	"policy-set-attachment-patch": {
		target: "policy-query",
		idFrom: "policy-set-attachment-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"policy_attachment","id":"{id}","relationships":{"tag":{"data":{"type":"tag","id":"0d7e4a6a-560e-40f1-a6dd-c7433a04f088"}}}}}`,
	},
	"policy-set-attachment-delete": {
		target: "policy-query",
		method: http.MethodDelete,
	},
	"policy-set-parameter-create": {
		target: "policy-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"policy_set_parameter","attributes":{"key":"max_length","value":"8","description":"maximum password length"}}}`,
	},
	"policy-set-parameter-list": {
		target: "policy-query",
	},
	"policy-set-parameter-get": {
		target: "policy-query",
	},
	"policy-set-parameter-patch": {
		target: "policy-query",
		idFrom: "policy-set-parameter-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"policy_set_parameter","id":"{id}","attributes":{"value":"12"}}}`,
	},
	"policy-set-parameter-missing-key": {
		target: "policy-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"policy_set_parameter","attributes":{"value":"1"}}}`,
	},
	"policy-set-parameter-delete": {
		target: "policy-query",
		method: http.MethodDelete,
	},
	"policy-exemption-create": {
		target: "policy-query",
		idFrom: "policy-set-create",
		method: http.MethodPost,
		body:   `{"data":{"type":"policy_exemption","attributes":{"ruleId":"rule-1","ticketReference":"CAPTURE-1","justification":"testing","expiresAt":"2030-01-01T00:00:00Z"},"relationships":{"organization":{"data":{"type":"organization","id":"` + capturedRelOrg + `"}},"policySet":{"data":{"type":"policy_set","id":"{id}"}},"workspace":{"data":{"type":"workspace","id":"0beb7db4-b073-4bd4-9a98-94fe7e2d9939"}}}}}`,
	},
	"policy-exemption-list": {
		target: "policy-query",
	},
	"policy-exemption-get": {
		target: "policy-query",
	},
	"policy-exemption-patch": {
		target: "policy-query",
		idFrom: "policy-exemption-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"policy_exemption","id":"{id}","attributes":{"justification":"after"}}}`,
	},
	"policy-exemption-missing-rule": {
		target: "policy-query",
		idFrom: "policy-set-create",
		method: http.MethodPost,
		body:   `{"data":{"type":"policy_exemption","attributes":{"ticketReference":"CAPTURE-2","justification":"x"},"relationships":{"organization":{"data":{"type":"organization","id":"` + capturedRelOrg + `"}},"policySet":{"data":{"type":"policy_set","id":"{id}"}}}}}`,
	},
	"policy-exemption-delete": {
		target: "policy-query",
		idFrom: "policy-exemption-create",
		method: http.MethodDelete,
	},
	"organization-policy-set-list": {
		target: "policy-query",
	},
	"organization-tag-list": {
		target: "tag-query",
		// The recording predates the create, so the list is empty.
		querier: &stubQuerier{},
	},
	"organization-tag-create": {
		target: "tag-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"tag","attributes":{"name":"capture-tag"}}}`,
	},
	"organization-tag-get": {
		target: "tag-query",
	},
	"organization-tag-patch": {
		target: "tag-query",
		idFrom: "organization-tag-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"tag","id":"{id}","attributes":{"name":"capture-tag-after"}}}`,
	},
	"organization-tag-filter": {
		target: "tag-query",
	},
	"organization-tag-create-missing-name": {
		target: "tag-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"tag","attributes":{}}}`,
	},
	"organization-tag-missing": {
		target: "tag-query",
	},
	"workspace-tag-list": {
		target:  "tag-query",
		querier: &stubQuerier{workspaces: []store.Workspace{{ID: capturedRelWS, OrganizationID: capturedRelOrg}}},
	},
	"workspace-tag-create": {
		target: "tag-query",
		idFrom: "organization-tag-create",
		method: http.MethodPost,
		body:   `{"data":{"type":"workspacetag","attributes":{"tagId":"{id}","value":"prod"}}}`,
	},
	"workspace-tag-list-after": {
		target: "tag-query",
	},
	"workspace-tag-patch": {
		target: "tag-query",
		idFrom: "workspace-tag-create",
		method: http.MethodPatch,
		body:   `{"data":{"type":"workspacetag","id":"{id}","attributes":{"value":null}}}`,
	},
	"workspace-tag-missing-tag": {
		target: "tag-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"workspacetag","attributes":{"value":"x"}}}`,
	},
	"workspace-tag-unknown-tag": {
		target: "tag-query",
		method: http.MethodPost,
		body:   `{"data":{"type":"workspacetag","attributes":{"tagId":"00000000-0000-4000-8000-000000000099","value":"x"}}}`,
	},
	"workspace-tag-delete": {
		target: "tag-query",
		idFrom: "workspace-tag-create",
		method: http.MethodDelete,
	},
	"organization-tag-delete": {
		target: "tag-query",
		idFrom: "organization-tag-create",
		method: http.MethodDelete,
	},
	"organization-tag-delete-unnamed": {
		target: "tag-query",
		idFrom: "organization-tag-create-missing-name",
		method: http.MethodDelete,
	},
	"organization-agent-unlink-in-use": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
}

// deferredCaptures are recorded but cannot be asserted yet. Each needs a reason,
// and the test below fails when a capture is added without one, so a new
// recording cannot be silently ignored.
var deferredCaptures = map[string]string{
	"organization-list-totals":         "we accept paging on organization, upstream answers 400; deliberate superset (docs/upstream/rest/README.md)",
	"organization-list-paged":          "same as organization-list-totals",
	"organization-list-include":        "we ignore include= and emit no `included` member; the UI sends no include on /api/v1",
	"organization-relationships":       "no /relationships/{relation} routes; the UI patches relationships on workspace only (#116)",
	"workspace-list-paged":             "no global /workspace collection; the capture documents meta.page, not a route we serve",
	"workspace-list-totals":            "same as workspace-list-paged",
	"delete-method-not-allowed":        "upstream answers 500 with a servlet error page, which is not a contract worth copying",
	"organization-globalvar-duplicate": "upstream's duplicate-key lifecycle hook throws and Spring answers 500; our unique constraint answers 409",
	"organization-atomic-unknown":      "upstream answers 400 with an array of per-operation results; we answer one JSON:API error object",
	"organization-tag-duplicate":       "upstream's unique-constraint violation surfaces as Spring's 500 page; our unique index answers 409",
	"workspace-tag-duplicate":          "the same Spring 500 on the unique (tag, workspace) constraint; ours answers 409",
	// /access-token/v1 is mounted outside /api/v1, which is all this harness
	// replays; internal/api/accesstoken_test.go asserts the 202 bodies.
	"access-token-current-teams":                    "served by NewAccessToken at /access-token/v1, outside this /api/v1 replay harness; accesstoken_test.go covers it",
	"access-token-no-token":                         "same endpoint; our 401 carries a JSON error where upstream answers an empty body",
	"access-token-permissions-organization":         "same endpoint; accesstoken_test.go covers the owner and team cases",
	"access-token-permissions-organization-team":    "same endpoint; internal/rbac/rbac_test.go covers the role matrix and the capture pins the write row",
	"access-token-permissions-organization-unknown": "same endpoint; accesstoken_test.go covers the owner flags on an unknown organization",
	"access-token-permissions-organization-reader":  "same endpoint; the capture pins an all-false read role",
	"access-token-permissions-organization-bad-id":  "same endpoint; upstream answers 500, we answer 400 for a malformed UUID",
	"access-token-permissions-workspace":            "same endpoint; accesstoken_test.go covers the workspace accumulation",
	"access-token-permissions-workspace-team":       "same endpoint; accesstoken_test.go covers the team half",
	"access-token-permissions-workspace-missing":    "same endpoint; accesstoken_test.go covers the missing-workspace case",
	"organization-provider-duplicate":               "upstream answers Spring's 500 page on the unique (organization, name) constraint; our unique index answers 409",
	"organization-provider-version-duplicate":       "the same Spring 500 on the unique (provider, version) constraint; ours answers 409",
	"organization-module-duplicate":                 "upstream answers 409 through a lifecycle hook; the replay stub cannot raise the database constraint, so internal/api/registry_test.go covers our 409",
}

// The ids upstream's demo seed gave the captures; the job captures were taken
// against these rows, so the replay uses the same ids.
const (
	capturedOrgID = "1c839e0a-ca5d-49a9-9c78-7a569df534b7"
	capturedWSID  = "0beb7db4-b073-4bd4-9a98-94fe7e2d9939"
	// The workspace relationship captures were seeded against this org, which
	// the upstream demo gives two projects.
	capturedRelOrg         = "70000000-0000-0000-0000-000000000001"
	capturedRelWS          = "77777777-7777-4777-8777-777777777777"
	capturedRelP1          = "70000000-0000-0000-0000-000000000010"
	capturedRelP2          = "70000000-0000-0000-0000-000000000020"
	capturedRelVCS         = "aaaaaaaa-0000-4000-8000-000000000001"
	capturedRelSSH         = "aaaaaaaa-0000-4000-8000-000000000002"
	capturedRelAgent       = "aaaaaaaa-0000-4000-8000-000000000003"
	capturedRelForeign     = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	capturedRelUnknown     = "00000000-0000-0000-0000-000000000000"
	capturedForeignOrg     = "d9b58bd3-f3fc-4056-a026-1163297e80a8"
	capturedGvID           = "99999999-9999-4999-8999-999999999990"
	capturedGvSecret       = "99999999-9999-4999-8999-999999999993"
	capturedStepID         = "66666666-6666-4666-8666-666666666666"
	capturedNotifyTemplate = "00000000-0000-4000-8000-000000000095"
	capturedAccessID       = "0a9ea824-20e1-45bd-8d41-5859a16b152c"
	capturedProjectAccess  = "ceeb8172-ab40-4960-b603-f6f2810b3e68"
)

func workspaceProjectBody(projectID string) string {
	return `{"data":{"type":"workspace","id":"` + capturedRelWS +
		`","relationships":{"project":{"data":{"type":"project","id":"` + projectID + `"}}}}}`
}

func captureWorkspaceStore() *stubQuerier {
	return &stubQuerier{
		workspaces: []store.Workspace{{
			ID: capturedRelWS, OrganizationID: capturedRelOrg, Name: "capture-rel-workspace",
			ProjectID: strPtr(capturedRelP1), VcsID: strPtr(capturedRelVCS),
		}},
		vcs: []store.Vc{{
			ID: capturedRelVCS, Name: "capture-vcs", Description: "capture", VcsType: "GITHUB",
			Status: strPtr("COMPLETED"), ConnectionType: strPtr("OAUTH"), OrganizationID: capturedRelOrg,
		}},
		projects: []store.Project{
			{ID: capturedRelP1, Name: "governance-terraform", OrganizationID: capturedRelOrg},
			{ID: capturedRelP2, Name: "governance-tofu", OrganizationID: capturedRelOrg},
			{ID: capturedRelForeign, Name: "demo-project", OrganizationID: capturedForeignOrg},
		},
		sshKeys: []store.Ssh{{ID: capturedRelSSH, Name: strPtr("capture-ssh"), OrganizationID: capturedRelOrg}},
		agents:  []store.Agent{{ID: capturedRelAgent, Name: strPtr("capture-agent"), OrganizationID: capturedRelOrg}},
	}
}

func captureGlobalVarStore() *stubQuerier {
	return &stubQuerier{
		globals: []store.Globalvar{
			{ID: capturedGvID, VariableKey: "CAPTURE_ONE", VariableValue: "one", VariableDescription: "capture", VariableCategory: "ENV", OrganizationID: capturedRelOrg},
			{ID: capturedGvSecret, VariableKey: "CAPTURE_SECRET", VariableValue: "hidden", VariableDescription: "capture", VariableCategory: "ENV", Sensitive: true, OrganizationID: capturedRelOrg},
			{ID: "99999999-9999-4999-8999-999999999992", VariableKey: "CAPTURE_TWO", VariableValue: "two", VariableDescription: "capture", VariableCategory: "TERRAFORM", Hcl: true, OrganizationID: capturedRelOrg},
		},
	}
}

func captureActionStore() *stubQuerier {
	return &stubQuerier{actions: []store.Action{
		{ID: "terrakube.open-documentation", Name: "Open Documentation", Label: "Open Documentation", Category: "General", Type: "Workspace/ResourceDrawer/Action", Action: "dGVzdA==", DisplayCriteria: `[{"filter": "true"}]`, Version: "1.0.0", Active: true},
		{ID: "terrakube.resource-details", Name: "Resource Details", Label: "Details", Category: "General", Type: "Workspace/ResourceDrawer/Tab", Action: "dGVzdA==", DisplayCriteria: `[{"filter": "true"}]`, Version: "1.0.0", Active: true},
		{ID: "terrakube.open-ai", Name: "Open AI", Label: "Open AI", Category: "AI", Type: "Workspace/Action", Action: "dGVzdA==", DisplayCriteria: `[{"filter": "true"}]`, Version: "0.1.0", Active: false},
		{ID: "capture.temp-action", Name: "Capture Temp", Label: "Temp", Category: "General", Type: "Workspace/Action", Action: "dGVzdA==", DisplayCriteria: `[{"filter":"true"}]`, Version: "1.0.0", Active: true},
	}}
}

// captureCollectionStore answers whichever collection, item and reference ids
// the collection captures were taken against; the ids change every run.
func captureCollectionStore(t *testing.T, captures map[string]upstreamCapture) *stubQuerier {
	t.Helper()
	collectionID := captures["organization-collection-get"].resourceID(t)
	itemID := captures["organization-collection-item-create"].resourceID(t)
	referenceID := captures["organization-collection-reference-create"].resourceID(t)

	return &stubQuerier{
		workspaces: []store.Workspace{
			{ID: capturedRelWS, OrganizationID: capturedRelOrg, Name: "capture-rel-workspace"},
			{ID: capturedWSID, OrganizationID: capturedOrgID, Name: "sample_cosmosdb"},
		},
		collections: []store.Collection{{
			ID: collectionID, Name: "capture-collection", Priority: 10, OrganizationID: capturedRelOrg,
		}},
		collectionItems: []store.Item{{
			ID: itemID, ItemKey: "CAPTURE_KEY", ItemValue: "after", ItemDescription: strPtr("capture"),
			ItemCategory: "ENV", CollectionID: collectionID,
		}},
		collectionRefs: []store.Reference{{
			ID: referenceID, Description: strPtr("capture"), WorkspaceID: capturedRelWS, CollectionID: collectionID,
		}},
	}
}

func captureSSHStore(t *testing.T, captures map[string]upstreamCapture) *stubQuerier {
	t.Helper()
	created := captures["organization-ssh-create"].resourceID(t)
	unnamed := captures["organization-ssh-create-missing-name"].resourceID(t)
	agent := captures["organization-agent-create"].resourceID(t)

	return &stubQuerier{
		sshKeys: []store.Ssh{
			{ID: capturedRelSSH, Name: strPtr("capture-ssh"), Description: strPtr("capture"), OrganizationID: capturedRelOrg},
			{ID: created, Name: strPtr("capture-ssh-created"), Description: strPtr("after"), SshType: strPtr("rsa"), OrganizationID: capturedRelOrg},
			{ID: unnamed, SshType: strPtr("rsa"), OrganizationID: capturedRelOrg},
		},
		agents: []store.Agent{
			{ID: capturedRelAgent, Name: strPtr("capture-agent"), Description: strPtr("capture"), Url: strPtr("http://agent:8090"), OrganizationID: capturedRelOrg},
			{ID: agent, Name: strPtr("capture-agent-created"), Description: strPtr("capture"), Url: strPtr("http://capture:8090"), OrganizationID: capturedRelOrg},
		},
		sshKeyUses: map[string]int64{capturedRelSSH: 1},
		agentUses:  map[string]int64{capturedRelAgent: 1},
	}
}

func captureVCSStore(t *testing.T, captures map[string]upstreamCapture) *stubQuerier {
	t.Helper()
	created := captures["organization-vcs-create"].resourceID(t)

	return &stubQuerier{
		vcs: []store.Vc{
			{
				ID: capturedRelVCS, Name: "capture-vcs", Description: "capture", VcsType: "GITHUB",
				Status: strPtr("COMPLETED"), ConnectionType: strPtr("OAUTH"), OrganizationID: capturedRelOrg,
			},
			{
				ID: created, Name: "capture-vcs-created", Description: "capture", VcsType: "GITLAB",
				ClientID: strPtr("client"), ClientSecret: strPtr("secret"), PrivateKey: strPtr("key"),
				Callback: strPtr("callback"), Endpoint: strPtr("https://gitlab.example.com"),
				ApiUrl: strPtr("https://gitlab.example.com/api/v4"), RedirectUrl: strPtr("https://ui.example.com/settings/vcs"),
				Status: strPtr("PENDING"), ConnectionType: strPtr("STANDALONE"), OrganizationID: capturedRelOrg,
			},
		},
		workspaces: []store.Workspace{{
			ID: capturedRelWS, OrganizationID: capturedRelOrg, Name: "capture-rel-workspace",
			VcsID: strPtr(capturedRelVCS),
		}},
	}
}

func captureTeamStore(t *testing.T, captures map[string]upstreamCapture) *stubQuerier {
	t.Helper()
	created := captures["organization-team-create"].resourceID(t)

	return &stubQuerier{
		teams: []store.Team{{
			ID: created, Name: "capture-team", ManageWorkspace: true, ManageJob: true,
			PlanJob: true, ApproveJob: true, DestroyJob: true, Role: strPtr("write"),
			OrganizationID: capturedRelOrg,
		}},
	}
}

func captureNotificationStore(t *testing.T, captures map[string]upstreamCapture) *stubQuerier {
	t.Helper()
	configID := captures["organization-notification-create"].resourceID(t)
	workspaceConfigID := captures["organization-notification-workspace-create"].resourceID(t)
	triggerID := captures["organization-notification-trigger-create"].resourceID(t)

	return &stubQuerier{
		org:        store.Organization{ID: capturedRelOrg},
		workspaces: []store.Workspace{{ID: capturedRelWS, OrganizationID: capturedRelOrg}},
		templates:  []store.Template{{ID: capturedNotifyTemplate, OrganizationID: capturedRelOrg}},
		notifications: []store.NotificationConfiguration{
			{
				ID: configID, Name: "capture-notify", Description: strPtr("capture"),
				OrganizationID: capturedRelOrg, ChannelType: "WEBHOOK",
				DestinationUrl: "https://example.com/hook", SigningSecret: strPtr("capture-secret"),
				Active: false, MessageStyle: "DETAILED",
			},
			{
				ID: workspaceConfigID, Name: "capture-notify-workspace",
				OrganizationID: capturedRelOrg, WorkspaceID: strPtr(capturedRelWS),
				ChannelType: "TEAMS", DestinationUrl: "https://teams.example.com/hook",
				Active: true, MessageStyle: "DETAILED",
			},
		},
		triggers: []store.NotificationTrigger{{
			ID: triggerID, JobStatus: "failed", ConfigurationID: configID,
		}},
	}
}

func captureTemplateStore(t *testing.T, captures map[string]upstreamCapture) *stubQuerier {
	t.Helper()
	created := captures["organization-template-create"].resourceID(t)
	empty := captures["organization-template-create-empty"].resourceID(t)

	return &stubQuerier{
		templates: []store.Template{
			{
				ID: created, Name: strPtr("capture-template"), Description: strPtr("capture"),
				Version: strPtr("9.9.9"), Tcl: strPtr("ZmxvdzoK"), OrganizationID: capturedRelOrg,
			},
			{ID: empty, OrganizationID: capturedRelOrg},
		},
	}
}

func captureTagStore(t *testing.T, captures map[string]upstreamCapture) *stubQuerier {
	t.Helper()
	tagID := captures["organization-tag-create"].resourceID(t)
	unnamed := captures["organization-tag-create-missing-name"].resourceID(t)
	binding := captures["workspace-tag-create"].resourceID(t)

	return &stubQuerier{
		tags: []store.Tag{
			{ID: tagID, Name: strPtr("capture-tag-after"), OrganizationID: capturedRelOrg},
			{ID: unnamed, OrganizationID: capturedRelOrg},
		},
		workspaceTags: []store.Workspacetag{
			{ID: binding, TagID: tagID, WorkspaceID: capturedRelWS},
		},
		workspaces: []store.Workspace{{ID: capturedRelWS, OrganizationID: capturedRelOrg, Name: "capture-rel-workspace"}},
	}
}

func captureJobs() []store.Job {
	return []store.Job{{
		ID:             2,
		OrganizationID: capturedOrgID,
		WorkspaceID:    capturedWSID,
		Status:         "completed",
		Via:            "UI",
	}}
}

// captureJobStore answers for whichever job id the job capture was taken
// against; the serial changes on every capture run.
func captureJobStore(id string) *stubQuerier {
	jobID, _ := strconv.ParseInt(id, 10, 64)
	return &stubQuerier{
		jobs: []store.Job{{
			ID:             jobID,
			OrganizationID: capturedOrgID,
			WorkspaceID:    capturedWSID,
			Status:         "completed",
			Via:            "UI",
		}},
		steps: []store.Step{{
			ID:         capturedStepID,
			JobID:      jobID,
			StepNumber: 1,
			Name:       strPtr("Plan"),
			Status:     strPtr("completed"),
			Output:     strPtr("ok"),
		}},
		workspaces: captureWorkspaces(),
	}
}

func captureSteps() []store.Step {
	return []store.Step{{
		ID:         capturedStepID,
		JobID:      2,
		StepNumber: 1,
		Name:       strPtr("Plan"),
		Status:     strPtr("completed"),
		Output:     strPtr("ok"),
	}}
}

func captureWorkspaces() []store.Workspace {
	return []store.Workspace{{ID: capturedWSID, OrganizationID: capturedOrgID, Name: "sample_cosmosdb"}}
}

func (s *stubQuerier) ListJobsByOrganization(context.Context, string) ([]store.Job, error) {
	return s.jobs, s.err
}

func (s *stubQuerier) GetJob(_ context.Context, id int64) (store.Job, error) {
	for _, job := range s.jobs {
		if job.ID == id {
			return job, nil
		}
	}
	return store.Job{}, pgx.ErrNoRows
}

func (s *stubQuerier) ListStepsByJob(_ context.Context, jobID int64) ([]store.Step, error) {
	var out []store.Step
	for _, step := range s.steps {
		if step.JobID == jobID {
			out = append(out, step)
		}
	}
	return out, nil
}

func (s *stubQuerier) UpdateJobStatus(_ context.Context, arg store.UpdateJobStatusParams) (store.Job, error) {
	for i, job := range s.jobs {
		if job.ID == arg.ID {
			s.jobs[i].Status = arg.Status
			return s.jobs[i], nil
		}
	}
	return store.Job{}, pgx.ErrNoRows
}

func (s *stubQuerier) SetWorkspaceLastJob(context.Context, store.SetWorkspaceLastJobParams) error {
	return s.err
}

func (s *stubQuerier) ListGlobalVariablesByOrganization(_ context.Context, orgID string) ([]store.Globalvar, error) {
	var out []store.Globalvar
	for _, variable := range s.globals {
		if variable.OrganizationID == orgID {
			out = append(out, variable)
		}
	}
	return out, s.err
}

func (s *stubQuerier) GetGlobalVariable(_ context.Context, id string) (store.Globalvar, error) {
	for _, variable := range s.globals {
		if variable.ID == id {
			return variable, nil
		}
	}
	return store.Globalvar{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateGlobalVariable(_ context.Context, arg store.CreateGlobalVariableParams) (store.Globalvar, error) {
	return store.Globalvar{
		ID: arg.ID, VariableKey: arg.VariableKey, VariableValue: arg.VariableValue,
		VariableDescription: arg.VariableDescription, VariableCategory: arg.VariableCategory,
		Sensitive: arg.Sensitive, Hcl: arg.Hcl, OrganizationID: arg.OrganizationID,
	}, s.err
}

func (s *stubQuerier) UpdateGlobalVariable(_ context.Context, arg store.UpdateGlobalVariableParams) (store.Globalvar, error) {
	return store.Globalvar{
		ID: arg.ID, VariableKey: arg.VariableKey, VariableValue: arg.VariableValue,
		VariableDescription: arg.VariableDescription, VariableCategory: arg.VariableCategory,
		Sensitive: arg.Sensitive, Hcl: arg.Hcl, OrganizationID: capturedRelOrg,
	}, s.err
}

func (s *stubQuerier) DeleteGlobalVariable(context.Context, string) error {
	return s.err
}

func (s *stubQuerier) GetProject(_ context.Context, id string) (store.Project, error) {
	for _, project := range s.projects {
		if project.ID == id {
			return project, nil
		}
	}
	return store.Project{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateWorkspace(_ context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
	return store.Workspace{
		ID: arg.ID, Name: arg.Name, OrganizationID: arg.OrganizationID, ProjectID: arg.ProjectID,
	}, s.err
}

func (s *stubQuerier) UpdateWorkspace(_ context.Context, arg store.UpdateWorkspaceParams) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == arg.ID {
			ws.ProjectID = arg.ProjectID
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *stubQuerier) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

// ListProjectsByOrganization is stubbed here because only the capture replay
// needs it; the organization handlers above it do not.
func (s *stubQuerier) ListProjectsByOrganization(context.Context, string) ([]store.Project, error) {
	return nil, s.err
}

// demoOrganizations stands in for upstream's demo seed; only the names matter,
// and no assertion here reads them.
func demoOrganizations() []store.Organization {
	return []store.Organization{
		{ID: "1", Name: "azure", ExecutionMode: "remote"},
		{ID: "2", Name: "aws", ExecutionMode: "local"},
	}
}

func TestEveryUpstreamCaptureIsAccountedFor(t *testing.T) {
	for name := range loadUpstreamCaptures(t) {
		_, replayed := replayedCaptures[name]
		reason, deferred := deferredCaptures[name]
		if !replayed && !deferred {
			t.Errorf("%s: captured from upstream but neither replayed nor deferred; add it to replayedCaptures or explain it in deferredCaptures", name)
		}
		if deferred && strings.TrimSpace(reason) == "" {
			t.Errorf("%s: deferred without a reason", name)
		}
	}
}

func TestCapturedShapesAreWhatWeReturn(t *testing.T) {
	captures := loadUpstreamCaptures(t)

	for name, plan := range replayedCaptures {
		t.Run(name, func(t *testing.T) {
			capture, ok := captures[name]
			if !ok {
				t.Fatalf("no capture named %s", name)
			}

			target, querier := plan.target, plan.querier
			idCapture := capture
			if plan.idFrom != "" {
				idCapture = captures[plan.idFrom]
			}
			switch target {
			case "resource":
				id := idCapture.resourceID(t)
				target = "/organization/" + id
				querier = &stubQuerier{org: store.Organization{ID: id, Name: "azure", ExecutionMode: "remote"}}
			case "resource-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
			case "resource-projects":
				target = "/organization/" + idCapture.resourceID(t) + "/project"
			case "job":
				target = "/organization/" + capturedOrgID + "/job/" + idCapture.resourceID(t)
				querier = captureJobStore(idCapture.resourceID(t))
			case "job-include":
				target = "/organization/" + capturedOrgID + "/job/" + idCapture.resourceID(t) + "?include=step,workspace"
				querier = captureJobStore(idCapture.resourceID(t))
			case "workspace":
				target = "/organization/" + capturedRelOrg + "/workspace/" + capturedRelWS
			case "workspace-collection":
				target = "/organization/" + capturedRelOrg + "/workspace"
			case "globalvar-collection":
				target = "/organization/" + capturedRelOrg + "/globalvar"
			case "globalvar-item":
				target = "/organization/" + capturedRelOrg + "/globalvar/" + capturedGvID
			case "globalvar-secret":
				target = "/organization/" + capturedRelOrg + "/globalvar/" + capturedGvSecret
			case "collection-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
				querier = captureCollectionStore(t, captures)
			case "ssh-agent-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
				querier = captureSSHStore(t, captures)
			case "tag-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
				if querier == nil {
					querier = captureTagStore(t, captures)
				}
			case "vcs-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
				if querier == nil {
					querier = captureVCSStore(t, captures)
				}
			case "notification-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
				if querier == nil {
					querier = captureNotificationStore(t, captures)
				}
			case "team-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
				if querier == nil {
					querier = captureTeamStore(t, captures)
				}
			case "template-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
				if querier == nil {
					querier = captureTemplateStore(t, captures)
				}
			case "registry-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
				if querier == nil {
					querier = registryCaptureStore(captures)
				}
			case "registry-provider-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
				querier = registryProviderStore(t, captures)
			case "policy-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
				if querier == nil {
					querier = policyCaptureStore(captures)
				}
			}
			if querier == nil {
				querier = &stubQuerier{}
			}

			handler := New(Deps{
				Store:     querier,
				Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
				AuthToken: plan.authToken,
			})
			method := plan.method
			if method == "" {
				method = http.MethodGet
			}
			body := plan.body
			if strings.Contains(body, "{id}") {
				body = strings.ReplaceAll(body, "{id}", idCapture.resourceID(t))
			}
			req := httptest.NewRequest(method, target, strings.NewReader(body))
			if body != "" {
				req.Header.Set("Content-Type", "application/vnd.api+json")
			}
			if plan.accept != "" {
				req.Header.Set("Content-Type", plan.accept)
				req.Header.Set("Accept", plan.accept)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			wantStatus := capture.status
			if plan.status != 0 {
				wantStatus = plan.status
			}
			if rec.Code != wantStatus {
				t.Fatalf("%s: status = %d, want %d (captured); body = %s", capture.path, rec.Code, wantStatus, rec.Body.String())
			}
			assertSameDocumentShape(t, capture, rec.Body.Bytes())
		})
	}
}

// assertSameDocumentShape compares the members a client can see, not values: an
// attribute rename or a missing member must fail, a different id or name must
// not.
func assertSameDocumentShape(t *testing.T, capture upstreamCapture, ours []byte) {
	t.Helper()

	if len(strings.TrimSpace(string(capture.body))) == 0 {
		// Upstream's 401 has an empty body; only the status is comparable.
		return
	}

	var theirs, mine map[string]any
	if err := json.Unmarshal(capture.body, &theirs); err != nil {
		t.Fatalf("%s: captured body is not JSON: %v", capture.name, err)
	}
	if err := json.Unmarshal(ours, &mine); err != nil {
		t.Fatalf("%s: our body is not JSON (%v): %s", capture.name, err, ours)
	}

	if theirErrors, ok := theirs["errors"].([]any); ok {
		// We add status and title; JSON:API allows it and the UI ignores them.
		myErrors, ok := mine["errors"].([]any)
		if !ok || len(myErrors) != len(theirErrors) {
			t.Fatalf("%s: errors = %v, want %d error object(s)", capture.name, mine["errors"], len(theirErrors))
		}
		for i, raw := range theirErrors {
			theirsErr, _ := raw.(map[string]any)
			mineErr, _ := myErrors[i].(map[string]any)
			if theirsErr["detail"] != nil && mineErr["detail"] == nil {
				t.Errorf("%s: errors[%d] has no detail; upstream says %v", capture.name, i, theirsErr["detail"])
			}
		}
		return
	}

	if got, want := sortedKeys(mine), sortedKeys(theirs); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("%s: top-level members = %v, want %v", capture.name, got, want)
	}

	switch their := theirs["data"].(type) {
	case []any:
		my, ok := mine["data"].([]any)
		if !ok {
			t.Fatalf("%s: data = %T, want array", capture.name, mine["data"])
		}
		if len(their) == 0 {
			if len(my) != 0 {
				t.Fatalf("%s: upstream returned an empty collection, we returned %d", capture.name, len(my))
			}
			return
		}
		if len(my) == 0 {
			t.Fatalf("%s: we returned an empty collection, upstream returned %d", capture.name, len(their))
		}
		assertSameResourceShape(t, capture.name, their[0], my[0])
	case map[string]any:
		assertSameResourceShape(t, capture.name, their, mine["data"])
	}
}

func assertSameResourceShape(t *testing.T, name string, theirs, mine any) {
	t.Helper()

	theirResource, ok := theirs.(map[string]any)
	if !ok {
		t.Fatalf("%s: captured data is %T, want a resource object", name, theirs)
	}
	myResource, ok := mine.(map[string]any)
	if !ok {
		t.Fatalf("%s: our data is %T, want a resource object", name, mine)
	}

	if myResource["type"] != theirResource["type"] {
		t.Errorf("%s: type = %v, want %v", name, myResource["type"], theirResource["type"])
	}
	if _, ok := myResource["id"].(string); !ok {
		t.Errorf("%s: id = %v, want a string", name, myResource["id"])
	}

	theirAttrs, _ := theirResource["attributes"].(map[string]any)
	myAttrs, _ := myResource["attributes"].(map[string]any)
	if got, want := sortedKeys(myAttrs), sortedKeys(theirAttrs); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("%s: attributes = %v, want %v", name, got, want)
	}

	// Most handlers omit relationships upstream serialises; grants must not,
	// because the UI reads the parent from them.
	if resourceType, _ := myResource["type"].(string); resourceType == "access" || resourceType == "project_access" {
		theirRels, _ := theirResource["relationships"].(map[string]any)
		myRels, _ := myResource["relationships"].(map[string]any)
		if got, want := sortedKeys(myRels), sortedKeys(theirRels); strings.Join(got, ",") != strings.Join(want, ",") {
			t.Errorf("%s: relationships = %v, want %v", name, got, want)
		}
	}
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (s *stubQuerier) ListActions(context.Context) ([]store.Action, error) {
	return s.actions, s.err
}

func (s *stubQuerier) ListCollectionsByOrganization(context.Context, string) ([]store.Collection, error) {
	return s.collections, s.err
}

func (s *stubQuerier) GetCollection(_ context.Context, id string) (store.Collection, error) {
	for _, collection := range s.collections {
		if collection.ID == id {
			return collection, nil
		}
	}
	return store.Collection{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateCollection(_ context.Context, arg store.CreateCollectionParams) (store.Collection, error) {
	return store.Collection{ID: arg.ID, Name: arg.Name, Description: arg.Description, Priority: arg.Priority, OrganizationID: arg.OrganizationID}, s.err
}

func (s *stubQuerier) UpdateCollection(_ context.Context, arg store.UpdateCollectionParams) (store.Collection, error) {
	return store.Collection{ID: arg.ID, Name: arg.Name, Description: arg.Description, Priority: arg.Priority}, s.err
}

func (s *stubQuerier) DeleteCollection(context.Context, string) error {
	return s.err
}

func (s *stubQuerier) ListItemsByCollection(context.Context, string) ([]store.Item, error) {
	return s.collectionItems, s.err
}

func (s *stubQuerier) GetItem(_ context.Context, id string) (store.Item, error) {
	for _, item := range s.collectionItems {
		if item.ID == id {
			return item, nil
		}
	}
	return store.Item{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateItem(_ context.Context, arg store.CreateItemParams) (store.Item, error) {
	return store.Item{
		ID: arg.ID, ItemDescription: arg.ItemDescription, ItemKey: arg.ItemKey,
		ItemValue: arg.ItemValue, Sensitive: arg.Sensitive, Hcl: arg.Hcl,
		ItemCategory: arg.ItemCategory, CollectionID: arg.CollectionID,
	}, s.err
}

func (s *stubQuerier) UpdateItem(_ context.Context, arg store.UpdateItemParams) (store.Item, error) {
	return store.Item{
		ID: arg.ID, ItemDescription: arg.ItemDescription, ItemKey: arg.ItemKey,
		ItemValue: arg.ItemValue, Sensitive: arg.Sensitive, Hcl: arg.Hcl,
		ItemCategory: arg.ItemCategory,
	}, s.err
}

func (s *stubQuerier) DeleteItem(context.Context, string) error {
	return s.err
}

func (s *stubQuerier) ListReferencesByCollection(context.Context, string) ([]store.Reference, error) {
	return s.collectionRefs, s.err
}

func (s *stubQuerier) ListReferencesByWorkspace(_ context.Context, workspaceID string) ([]store.Reference, error) {
	var out []store.Reference
	for _, reference := range s.collectionRefs {
		if reference.WorkspaceID == workspaceID {
			out = append(out, reference)
		}
	}
	return out, s.err
}

func (s *stubQuerier) GetReference(_ context.Context, id string) (store.Reference, error) {
	for _, reference := range s.collectionRefs {
		if reference.ID == id {
			return reference, nil
		}
	}
	return store.Reference{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateReference(_ context.Context, arg store.CreateReferenceParams) (store.Reference, error) {
	return store.Reference{ID: arg.ID, Description: arg.Description, WorkspaceID: arg.WorkspaceID, CollectionID: arg.CollectionID}, s.err
}

func (s *stubQuerier) DeleteReference(context.Context, string) error {
	return s.err
}

func (s *stubQuerier) ListSSHKeysByOrganization(context.Context, string) ([]store.Ssh, error) {
	return s.sshKeys, s.err
}

func (s *stubQuerier) GetSSHKey(_ context.Context, id string) (store.Ssh, error) {
	for _, key := range s.sshKeys {
		if key.ID == id {
			return key, nil
		}
	}
	return store.Ssh{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateSSHKey(_ context.Context, arg store.CreateSSHKeyParams) (store.Ssh, error) {
	return store.Ssh{
		ID: arg.ID, Name: arg.Name, Description: arg.Description,
		PrivateKey: arg.PrivateKey, SshType: arg.SshType, OrganizationID: arg.OrganizationID,
	}, nil
}

func (s *stubQuerier) UpdateSSHKey(_ context.Context, arg store.UpdateSSHKeyParams) (store.Ssh, error) {
	return store.Ssh{
		ID: arg.ID, Name: arg.Name, Description: arg.Description,
		PrivateKey: arg.PrivateKey, SshType: arg.SshType,
	}, nil
}

func (s *stubQuerier) CountWorkspacesWithSSHKey(_ context.Context, id *string) (int64, error) {
	return s.sshKeyUses[*id], s.err
}

func (s *stubQuerier) DeleteSSHKey(context.Context, string) error {
	return nil
}

func (s *stubQuerier) ListAgentsByOrganization(context.Context, string) ([]store.Agent, error) {
	return s.agents, s.err
}

func (s *stubQuerier) GetAgent(_ context.Context, id string) (store.Agent, error) {
	for _, agent := range s.agents {
		if agent.ID == id {
			return agent, nil
		}
	}
	return store.Agent{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateAgent(_ context.Context, arg store.CreateAgentParams) (store.Agent, error) {
	return store.Agent{
		ID: arg.ID, Name: arg.Name, Url: arg.Url, Description: arg.Description, OrganizationID: arg.OrganizationID,
	}, nil
}

func (s *stubQuerier) CountWorkspacesWithAgent(_ context.Context, id *string) (int64, error) {
	return s.agentUses[*id], s.err
}

func (s *stubQuerier) DeleteAgent(context.Context, string) error {
	return nil
}

func (s *stubQuerier) ListTagsByOrganization(context.Context, string) ([]store.Tag, error) {
	return s.tags, s.err
}

func (s *stubQuerier) GetTag(_ context.Context, id string) (store.Tag, error) {
	for _, tag := range s.tags {
		if tag.ID == id {
			return tag, nil
		}
	}
	return store.Tag{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateTag(_ context.Context, arg store.CreateTagParams) (store.Tag, error) {
	return store.Tag{ID: arg.ID, Name: arg.Name, OrganizationID: arg.OrganizationID}, s.err
}

func (s *stubQuerier) UpdateTag(_ context.Context, arg store.UpdateTagParams) (store.Tag, error) {
	return store.Tag{ID: arg.ID, Name: arg.Name, OrganizationID: capturedRelOrg}, s.err
}

func (s *stubQuerier) DeleteTag(context.Context, string) error {
	return s.err
}

func (s *stubQuerier) ListWorkspaceTagsByWorkspace(context.Context, string) ([]store.Workspacetag, error) {
	return s.workspaceTags, s.err
}

func (s *stubQuerier) GetWorkspaceTag(_ context.Context, id string) (store.Workspacetag, error) {
	for _, binding := range s.workspaceTags {
		if binding.ID == id {
			return binding, nil
		}
	}
	return store.Workspacetag{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateWorkspaceTag(_ context.Context, arg store.CreateWorkspaceTagParams) (store.Workspacetag, error) {
	return store.Workspacetag{
		ID: arg.ID, TagID: arg.TagID, TagValue: arg.TagValue, WorkspaceID: arg.WorkspaceID,
	}, s.err
}

func (s *stubQuerier) UpdateWorkspaceTagValue(_ context.Context, arg store.UpdateWorkspaceTagValueParams) (store.Workspacetag, error) {
	return store.Workspacetag{ID: arg.ID, TagValue: arg.TagValue, WorkspaceID: capturedRelWS}, s.err
}

func (s *stubQuerier) DeleteWorkspaceTag(context.Context, string) error {
	return s.err
}

func (s *stubQuerier) ListVCSByOrganization(context.Context, string) ([]store.Vc, error) {
	return s.vcs, s.err
}

func (s *stubQuerier) GetVCS(_ context.Context, id string) (store.Vc, error) {
	for _, provider := range s.vcs {
		if provider.ID == id {
			return provider, nil
		}
	}
	return store.Vc{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateVCS(_ context.Context, arg store.CreateVCSParams) (store.Vc, error) {
	return store.Vc{
		ID: arg.ID, Name: arg.Name, Description: arg.Description, VcsType: arg.VcsType,
		ClientID: arg.ClientID, ClientSecret: arg.ClientSecret, AccessToken: arg.AccessToken,
		Status: arg.Status, Callback: arg.Callback, RedirectUrl: arg.RedirectUrl,
		Endpoint: arg.Endpoint, ApiUrl: arg.ApiUrl, ConnectionType: arg.ConnectionType,
		PrivateKey: arg.PrivateKey, OrganizationID: arg.OrganizationID,
	}, s.err
}

func (s *stubQuerier) UpdateVCS(_ context.Context, arg store.UpdateVCSParams) (store.Vc, error) {
	for i := range s.vcs {
		if s.vcs[i].ID == arg.ID {
			s.vcs[i].Name = arg.Name
			s.vcs[i].Description = arg.Description
			s.vcs[i].VcsType = arg.VcsType
			s.vcs[i].ClientID = arg.ClientID
			s.vcs[i].ClientSecret = arg.ClientSecret
			s.vcs[i].AccessToken = arg.AccessToken
			s.vcs[i].Status = arg.Status
			s.vcs[i].Callback = arg.Callback
			s.vcs[i].RedirectUrl = arg.RedirectUrl
			s.vcs[i].Endpoint = arg.Endpoint
			s.vcs[i].ApiUrl = arg.ApiUrl
			s.vcs[i].ConnectionType = arg.ConnectionType
			s.vcs[i].PrivateKey = arg.PrivateKey
			return s.vcs[i], s.err
		}
	}
	return store.Vc{}, pgx.ErrNoRows
}

func (s *stubQuerier) ClearWorkspaceVCS(_ context.Context, id *string) error {
	for i := range s.workspaces {
		if s.workspaces[i].VcsID != nil && id != nil && *s.workspaces[i].VcsID == *id {
			s.workspaces[i].VcsID = nil
		}
	}
	return s.err
}

func (s *stubQuerier) DeleteVCS(_ context.Context, id string) error {
	for i := range s.vcs {
		if s.vcs[i].ID == id {
			s.vcs = append(s.vcs[:i], s.vcs[i+1:]...)
			return s.err
		}
	}
	return s.err
}

func (s *stubQuerier) ListWorkspaceIDsByVCS(_ context.Context, id *string) ([]string, error) {
	ids := []string{}
	if id == nil {
		return ids, s.err
	}
	for _, ws := range s.workspaces {
		if ws.VcsID != nil && *ws.VcsID == *id {
			ids = append(ids, ws.ID)
		}
	}
	return ids, s.err
}

func (s *stubQuerier) ListTemplatesByOrganization(context.Context, string) ([]store.Template, error) {
	return s.templates, s.err
}

func (s *stubQuerier) GetTemplate(_ context.Context, id string) (store.Template, error) {
	for _, template := range s.templates {
		if template.ID == id {
			return template, nil
		}
	}
	return store.Template{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateTemplate(_ context.Context, arg store.CreateTemplateParams) (store.Template, error) {
	return store.Template{
		ID: arg.ID, Name: arg.Name, Description: arg.Description, Version: arg.Version,
		Tcl: arg.Tcl, OrganizationID: arg.OrganizationID,
	}, s.err
}

func (s *stubQuerier) UpdateTemplate(_ context.Context, arg store.UpdateTemplateParams) (store.Template, error) {
	for i := range s.templates {
		if s.templates[i].ID == arg.ID {
			s.templates[i].Name = arg.Name
			s.templates[i].Description = arg.Description
			s.templates[i].Version = arg.Version
			s.templates[i].Tcl = arg.Tcl
			return s.templates[i], s.err
		}
	}
	return store.Template{}, pgx.ErrNoRows
}

func (s *stubQuerier) DeleteTemplate(_ context.Context, id string) error {
	for i := range s.templates {
		if s.templates[i].ID == id {
			s.templates = append(s.templates[:i], s.templates[i+1:]...)
			return s.err
		}
	}
	return s.err
}

func (s *stubQuerier) ListTeamsByOrganization(context.Context, string) ([]store.Team, error) {
	return s.teams, s.err
}

func (s *stubQuerier) GetTeam(_ context.Context, id string) (store.Team, error) {
	for _, team := range s.teams {
		if team.ID == id {
			return team, nil
		}
	}
	return store.Team{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateTeam(_ context.Context, arg store.CreateTeamParams) (store.Team, error) {
	return store.Team{
		ID: arg.ID, Name: arg.Name, ManageWorkspace: arg.ManageWorkspace,
		ManageModule: arg.ManageModule, ManageProvider: arg.ManageProvider,
		ManageVcs: arg.ManageVcs, ManageTemplate: arg.ManageTemplate,
		ManageState: arg.ManageState, ManageCollection: arg.ManageCollection,
		ManageJob: arg.ManageJob, ManagePolicies: arg.ManagePolicies,
		PlanJob: arg.PlanJob, ApproveJob: arg.ApproveJob, DestroyJob: arg.DestroyJob,
		Role: arg.Role, OrganizationID: arg.OrganizationID,
	}, s.err
}

func (s *stubQuerier) UpdateTeam(_ context.Context, arg store.UpdateTeamParams) (store.Team, error) {
	for i := range s.teams {
		if s.teams[i].ID == arg.ID {
			s.teams[i].Name = arg.Name
			s.teams[i].ManageWorkspace = arg.ManageWorkspace
			s.teams[i].ManageModule = arg.ManageModule
			s.teams[i].ManageProvider = arg.ManageProvider
			s.teams[i].ManageVcs = arg.ManageVcs
			s.teams[i].ManageTemplate = arg.ManageTemplate
			s.teams[i].ManageState = arg.ManageState
			s.teams[i].ManageCollection = arg.ManageCollection
			s.teams[i].ManageJob = arg.ManageJob
			s.teams[i].ManagePolicies = arg.ManagePolicies
			s.teams[i].PlanJob = arg.PlanJob
			s.teams[i].ApproveJob = arg.ApproveJob
			s.teams[i].DestroyJob = arg.DestroyJob
			s.teams[i].Role = arg.Role
			return s.teams[i], s.err
		}
	}
	return store.Team{}, pgx.ErrNoRows
}

func (s *stubQuerier) DeleteTeam(_ context.Context, id string) error {
	for i := range s.teams {
		if s.teams[i].ID == id {
			s.teams = append(s.teams[:i], s.teams[i+1:]...)
			return s.err
		}
	}
	return s.err
}

func (s *stubQuerier) ListNotificationConfigurationsByOrganization(_ context.Context, orgID string) ([]store.NotificationConfiguration, error) {
	var out []store.NotificationConfiguration
	for _, config := range s.notifications {
		if config.OrganizationID == orgID {
			out = append(out, config)
		}
	}
	return out, s.err
}

func (s *stubQuerier) ListNotificationConfigurationsByWorkspace(_ context.Context, workspaceID *string) ([]store.NotificationConfiguration, error) {
	var out []store.NotificationConfiguration
	for _, config := range s.notifications {
		if workspaceID != nil && config.WorkspaceID != nil && *config.WorkspaceID == *workspaceID {
			out = append(out, config)
		}
	}
	return out, s.err
}

func (s *stubQuerier) GetNotificationConfiguration(_ context.Context, id string) (store.NotificationConfiguration, error) {
	for _, config := range s.notifications {
		if config.ID == id {
			return config, nil
		}
	}
	return store.NotificationConfiguration{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateNotificationConfiguration(_ context.Context, arg store.CreateNotificationConfigurationParams) (store.NotificationConfiguration, error) {
	config := store.NotificationConfiguration{
		ID: arg.ID, Name: arg.Name, Description: arg.Description,
		OrganizationID: arg.OrganizationID, WorkspaceID: arg.WorkspaceID,
		ChannelType: arg.ChannelType, DestinationUrl: arg.DestinationUrl,
		SigningSecret: arg.SigningSecret, Active: arg.Active, MessageStyle: arg.MessageStyle,
	}
	s.notifications = append(s.notifications, config)
	return config, s.err
}

func (s *stubQuerier) UpdateNotificationConfiguration(_ context.Context, arg store.UpdateNotificationConfigurationParams) (store.NotificationConfiguration, error) {
	for i := range s.notifications {
		if s.notifications[i].ID == arg.ID {
			s.notifications[i].Name = arg.Name
			s.notifications[i].Description = arg.Description
			s.notifications[i].ChannelType = arg.ChannelType
			s.notifications[i].DestinationUrl = arg.DestinationUrl
			s.notifications[i].SigningSecret = arg.SigningSecret
			s.notifications[i].Active = arg.Active
			s.notifications[i].MessageStyle = arg.MessageStyle
			return s.notifications[i], s.err
		}
	}
	return store.NotificationConfiguration{}, pgx.ErrNoRows
}

func (s *stubQuerier) DeleteNotificationConfiguration(_ context.Context, id string) error {
	for i := range s.notifications {
		if s.notifications[i].ID == id {
			s.notifications = append(s.notifications[:i], s.notifications[i+1:]...)
			break
		}
	}
	var kept []store.NotificationTrigger
	for _, trigger := range s.triggers {
		if trigger.ConfigurationID != id {
			kept = append(kept, trigger)
		}
	}
	s.triggers = kept
	return s.err
}

func (s *stubQuerier) ListNotificationTriggersByConfiguration(_ context.Context, configurationID string) ([]store.NotificationTrigger, error) {
	var out []store.NotificationTrigger
	for _, trigger := range s.triggers {
		if trigger.ConfigurationID == configurationID {
			out = append(out, trigger)
		}
	}
	return out, s.err
}

func (s *stubQuerier) GetNotificationTrigger(_ context.Context, id string) (store.NotificationTrigger, error) {
	for _, trigger := range s.triggers {
		if trigger.ID == id {
			return trigger, nil
		}
	}
	return store.NotificationTrigger{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateNotificationTrigger(_ context.Context, arg store.CreateNotificationTriggerParams) (store.NotificationTrigger, error) {
	trigger := store.NotificationTrigger{
		ID: arg.ID, JobStatus: arg.JobStatus, ConfigurationID: arg.ConfigurationID,
	}
	s.triggers = append(s.triggers, trigger)
	return trigger, s.err
}

func (s *stubQuerier) DeleteNotificationTrigger(_ context.Context, id string) error {
	for i := range s.triggers {
		if s.triggers[i].ID == id {
			s.triggers = append(s.triggers[:i], s.triggers[i+1:]...)
			return s.err
		}
	}
	return s.err
}

func (s *stubQuerier) ListNotificationConfigurationTemplateIDs(_ context.Context, configurationID string) ([]string, error) {
	return s.notificationTemplates[configurationID], s.err
}

func (s *stubQuerier) DeleteNotificationConfigurationTemplates(_ context.Context, configurationID string) error {
	delete(s.notificationTemplates, configurationID)
	return s.err
}

func (s *stubQuerier) AddNotificationConfigurationTemplate(_ context.Context, arg store.AddNotificationConfigurationTemplateParams) error {
	if s.notificationTemplates == nil {
		s.notificationTemplates = map[string][]string{}
	}
	s.notificationTemplates[arg.ConfigurationID] = append(s.notificationTemplates[arg.ConfigurationID], arg.TemplateID)
	return s.err
}

func (s *stubQuerier) GetAction(_ context.Context, id string) (store.Action, error) {
	for _, action := range s.actions {
		if action.ID == id {
			return action, nil
		}
	}
	return store.Action{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateAction(_ context.Context, arg store.CreateActionParams) (store.Action, error) {
	return store.Action{
		ID: arg.ID, Name: arg.Name, Label: arg.Label, Category: arg.Category,
		Description: arg.Description, Type: arg.Type, Action: arg.Action,
		DisplayCriteria: arg.DisplayCriteria, Version: arg.Version, Active: arg.Active,
	}, s.err
}

func (s *stubQuerier) UpdateAction(_ context.Context, arg store.UpdateActionParams) (store.Action, error) {
	for _, action := range s.actions {
		if action.ID == arg.ID {
			return action, s.err
		}
	}
	return store.Action{}, pgx.ErrNoRows
}

func (s *stubQuerier) DeleteAction(context.Context, string) error {
	return s.err
}

// captureGrantWorkspaces and captureGrantProjects give the grant replays the
// parents the captured paths address; the capture request itself is built from
// the path, so no other attribute matters.
func captureGrantWorkspaces() []store.Workspace {
	return []store.Workspace{{ID: capturedRelWS, OrganizationID: capturedRelOrg}}
}

func captureGrantProjects() []store.Project {
	return []store.Project{{ID: capturedRelP1, OrganizationID: capturedRelOrg}}
}

func captureGrantStore() *stubQuerier {
	return &stubQuerier{
		workspaces:    captureGrantWorkspaces(),
		projects:      captureGrantProjects(),
		access:        []store.Access{captureAccessRow()},
		projectAccess: []store.ProjectAccess{captureProjectAccessRow()},
	}
}

func captureAccessRow() store.Access {
	return store.Access{
		ID: capturedAccessID, Name: "capture-group", ManageWorkspace: true, ManageState: true,
		PlanJob: true, DestroyJob: true, Role: strPtr("custom"), WorkspaceID: capturedRelWS,
	}
}

func captureProjectAccessRow() store.ProjectAccess {
	return store.ProjectAccess{
		ID: capturedProjectAccess, Name: strPtr("capture-group"), ManageWorkspace: true,
		ManageState: true, DestroyJob: true, Role: strPtr("custom"), ProjectID: capturedRelP1,
	}
}

func (s *stubQuerier) ListAccessByWorkspace(_ context.Context, workspaceID string) ([]store.Access, error) {
	if s.err != nil {
		return nil, s.err
	}
	var out []store.Access
	for _, access := range s.access {
		if access.WorkspaceID == workspaceID {
			out = append(out, access)
		}
	}
	return out, nil
}

func (s *stubQuerier) GetAccess(_ context.Context, id string) (store.Access, error) {
	if s.err != nil {
		return store.Access{}, s.err
	}
	for _, access := range s.access {
		if access.ID == id {
			return access, nil
		}
	}
	return store.Access{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateAccess(_ context.Context, arg store.CreateAccessParams) (store.Access, error) {
	return store.Access{
		ID: arg.ID, Name: arg.Name, ManageWorkspace: arg.ManageWorkspace,
		ManageState: arg.ManageState, ManageJob: arg.ManageJob, PlanJob: arg.PlanJob,
		ApproveJob: arg.ApproveJob, DestroyJob: arg.DestroyJob, Role: arg.Role,
		WorkspaceID: arg.WorkspaceID,
	}, s.err
}

func (s *stubQuerier) UpdateAccess(_ context.Context, arg store.UpdateAccessParams) (store.Access, error) {
	s.updatedAccess = &arg
	return store.Access{
		ID: arg.ID, Name: arg.Name, ManageWorkspace: arg.ManageWorkspace,
		ManageState: arg.ManageState, ManageJob: arg.ManageJob, PlanJob: arg.PlanJob,
		ApproveJob: arg.ApproveJob, DestroyJob: arg.DestroyJob, Role: arg.Role,
	}, s.err
}

func (s *stubQuerier) DeleteAccess(context.Context, string) error {
	return s.err
}

func (s *stubQuerier) ListProjectAccessByProject(_ context.Context, projectID string) ([]store.ProjectAccess, error) {
	if s.err != nil {
		return nil, s.err
	}
	var out []store.ProjectAccess
	for _, access := range s.projectAccess {
		if access.ProjectID == projectID {
			out = append(out, access)
		}
	}
	return out, nil
}

func (s *stubQuerier) GetProjectAccess(_ context.Context, id string) (store.ProjectAccess, error) {
	if s.err != nil {
		return store.ProjectAccess{}, s.err
	}
	for _, access := range s.projectAccess {
		if access.ID == id {
			return access, nil
		}
	}
	return store.ProjectAccess{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateProjectAccess(_ context.Context, arg store.CreateProjectAccessParams) (store.ProjectAccess, error) {
	return store.ProjectAccess{
		ID: arg.ID, Name: arg.Name, ManageWorkspace: arg.ManageWorkspace,
		ManageState: arg.ManageState, ManageJob: arg.ManageJob, PlanJob: arg.PlanJob,
		ApproveJob: arg.ApproveJob, DestroyJob: arg.DestroyJob, Role: arg.Role,
		ProjectID: arg.ProjectID,
	}, s.err
}

func (s *stubQuerier) UpdateProjectAccess(_ context.Context, arg store.UpdateProjectAccessParams) (store.ProjectAccess, error) {
	return store.ProjectAccess{
		ID: arg.ID, Name: arg.Name, ManageWorkspace: arg.ManageWorkspace,
		ManageState: arg.ManageState, ManageJob: arg.ManageJob, PlanJob: arg.PlanJob,
		ApproveJob: arg.ApproveJob, DestroyJob: arg.DestroyJob, Role: arg.Role,
	}, s.err
}

func (s *stubQuerier) DeleteProjectAccess(context.Context, string) error {
	return s.err
}

func (s *stubQuerier) ListNotificationDeliveriesByWorkspace(_ context.Context, arg store.ListNotificationDeliveriesByWorkspaceParams) ([]store.ListNotificationDeliveriesByWorkspaceRow, error) {
	rows := s.deliveries
	if int64(len(rows)) > int64(arg.Limit) {
		rows = rows[:arg.Limit]
	}
	return rows, s.err
}

func (s *stubQuerier) NotificationOutboxExistsForWorkspace(_ context.Context, arg store.NotificationOutboxExistsForWorkspaceParams) (bool, error) {
	return s.deliveryExists[arg.ID], s.err
}

func (s *stubQuerier) RearmNotificationOutbox(_ context.Context, id string) (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	if !s.rearmable[id] {
		return 0, nil
	}
	s.rearmed = append(s.rearmed, id)
	return 1, nil
}
