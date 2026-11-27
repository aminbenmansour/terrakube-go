package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	testHookID   = "31313131-3131-4131-8131-313131313131"
	testWSID     = "32323232-3232-4232-8232-323232323232"
	testVCSSID   = "33333333-3333-4333-8333-333333333333"
	testOrgID    = "34343434-3434-4434-8434-343434343434"
	testTplID    = "35353535-3535-4535-8535-353535353535"
	testPushBody = `{"ref":"refs/heads/main","after":"abc123",` +
		`"head_commit":{"id":"abc123"},` +
		`"commits":[{"id":"abc123","added":["modules/vpc/main.tf"],"removed":[],"modified":[]}],` +
		`"sender":{"login":"octocat"},` +
		`"repository":{"clone_url":"https://github.com/acme/infra.git","html_url":"https://github.com/acme/infra"}}`
)

type stubStore struct {
	store.Querier
	provider   string
	source     string
	webhook    store.Webhook
	workspace  store.Workspace
	events     []store.WebhookEvent
	template   store.Template
	deliveries map[string]bool
	job        store.CreateJobParams
	steps      int
	queued     []int64
	failed     bool
}

func newStubStore(provider string) *stubStore {
	source := "https://github.com/acme/infra.git"
	if provider == "GITLAB" {
		source = "https://gitlab.com/acme/infra.git"
	}
	vcsID := testVCSSID
	return &stubStore{
		provider: provider,
		source:   source,
		webhook:  store.Webhook{ID: testHookID, WorkspaceID: testWSID},
		workspace: store.Workspace{
			ID: testWSID, OrganizationID: testOrgID, VcsID: &vcsID, Source: &source,
		},
		template: store.Template{ID: testTplID, OrganizationID: testOrgID},
		events: []store.WebhookEvent{{
			ID: "36363636-3636-4636-8636-363636363636", WebhookID: testHookID,
			Branch: "main", Path: "modules/**", PathType: "PATTERN", Event: "PUSH",
			TemplateID: testTplID, Priority: 1,
		}},
		deliveries: map[string]bool{},
	}
}

func (s *stubStore) GetWebhook(_ context.Context, id string) (store.Webhook, error) {
	if id != s.webhook.ID {
		return store.Webhook{}, pgx.ErrNoRows
	}
	return s.webhook, nil
}

func (s *stubStore) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	if id != s.workspace.ID {
		return store.Workspace{}, pgx.ErrNoRows
	}
	return s.workspace, nil
}

func (s *stubStore) GetVCS(_ context.Context, id string) (store.Vc, error) {
	if id != testVCSSID {
		return store.Vc{}, pgx.ErrNoRows
	}
	return store.Vc{ID: testVCSSID, OrganizationID: testOrgID, VcsType: s.provider}, nil
}

func (s *stubStore) ListWebhookEvents(_ context.Context, id string) ([]store.WebhookEvent, error) {
	if id != s.webhook.ID {
		return nil, nil
	}
	return s.events, nil
}

func (s *stubStore) GetTemplate(_ context.Context, id string) (store.Template, error) {
	if id != s.template.ID {
		return store.Template{}, pgx.ErrNoRows
	}
	return s.template, nil
}

func (s *stubStore) CreateJob(_ context.Context, arg store.CreateJobParams) (store.Job, error) {
	s.job = arg
	return store.Job{ID: 77, WorkspaceID: arg.WorkspaceID}, nil
}

func (s *stubStore) CreateStep(_ context.Context, _ store.CreateStepParams) (store.Step, error) {
	s.steps++
	return store.Step{}, nil
}

func (s *stubStore) SetWorkspaceLastJob(_ context.Context, _ store.SetWorkspaceLastJobParams) error {
	return nil
}

func (s *stubStore) UpdateJobStatus(_ context.Context, arg store.UpdateJobStatusParams) (store.Job, error) {
	s.failed = arg.Status == "failed"
	return store.Job{ID: arg.ID, Status: arg.Status}, nil
}

func (s *stubStore) RecordWebhookDelivery(_ context.Context, arg store.RecordWebhookDeliveryParams) (int64, error) {
	key := arg.WebhookID + "/" + arg.DeliveryID
	if s.deliveries[key] {
		return 0, nil
	}
	s.deliveries[key] = true
	return 1, nil
}

func (s *stubStore) DeleteWebhookDelivery(_ context.Context, arg store.DeleteWebhookDeliveryParams) error {
	delete(s.deliveries, arg.WebhookID+"/"+arg.DeliveryID)
	return nil
}

func (s *stubStore) PruneWebhookDeliveries(context.Context) error { return nil }

type stubQueue struct {
	queued []int64
}

func (q *stubQueue) EnqueueJob(_ context.Context, jobID int64) (string, error) {
	q.queued = append(q.queued, jobID)
	return "1-0", nil
}

func githubSignature(workspaceID string, body string) string {
	secret := base64.StdEncoding.EncodeToString([]byte(workspaceID))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func post(t *testing.T, handler http.Handler, body string, headers map[string]string) int {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/"+testHookID, strings.NewReader(body))
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder.Code
}

func TestGitHubPushStartsRun(t *testing.T) {
	stub := newStubStore("GITHUB")
	queue := &stubQueue{}
	handler := New(Deps{Store: stub, Queue: queue})

	status := post(t, handler, testPushBody, map[string]string{
		"X-GitHub-Event":      "push",
		"X-GitHub-Delivery":   "delivery-1",
		"X-Hub-Signature-256": githubSignature(testWSID, testPushBody),
	})
	if status != http.StatusOK {
		t.Fatalf("status = %d", status)
	}
	if stub.job.Via != "Github" || stub.job.Status != "pending" || !stub.job.Refresh || !stub.job.PlanChanges {
		t.Fatalf("job = %+v", stub.job)
	}
	if stub.job.OverrideBranch == nil || *stub.job.OverrideBranch != "main" {
		t.Fatalf("override branch = %v", stub.job.OverrideBranch)
	}
	if stub.job.CommitID == nil || *stub.job.CommitID != "abc123" {
		t.Fatalf("commit = %v", stub.job.CommitID)
	}
	if stub.job.TemplateReference == nil || *stub.job.TemplateReference != testTplID {
		t.Fatalf("template = %v", stub.job.TemplateReference)
	}
	if stub.job.CreatedBy == nil || *stub.job.CreatedBy != "octocat" {
		t.Fatalf("created by = %v", stub.job.CreatedBy)
	}
	if stub.steps != 3 || len(queue.queued) != 1 || queue.queued[0] != 77 {
		t.Fatalf("steps = %d, queued = %v", stub.steps, queue.queued)
	}
}

func TestGitHubBadOrMissingSignatureIsRejected(t *testing.T) {
	for name, headers := range map[string]map[string]string{
		"missing": {"X-GitHub-Event": "push"},
		"wrong": {
			"X-GitHub-Event":      "push",
			"X-Hub-Signature-256": githubSignature("someone-else", testPushBody),
		},
	} {
		t.Run(name, func(t *testing.T) {
			stub := newStubStore("GITHUB")
			handler := New(Deps{Store: stub, Queue: &stubQueue{}})
			if status := post(t, handler, testPushBody, headers); status != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", status)
			}
			if stub.job.WorkspaceID != "" || len(stub.deliveries) != 0 {
				t.Fatalf("created job or delivery for an unsigned request: %+v / %+v", stub.job, stub.deliveries)
			}
		})
	}
}

func TestGitLabPushRequiresTheToken(t *testing.T) {
	body := `{"object_kind":"push","ref":"refs/heads/main","checkout_sha":"def456",` +
		`"user_username":"alice","commits":[{"id":"def456","added":["modules/vpc/main.tf"]}],` +
		`"project":{"git_http_url":"https://gitlab.com/acme/infra.git"}}`
	secret := base64.StdEncoding.EncodeToString([]byte(testWSID))

	stub := newStubStore("GITLAB")
	queue := &stubQueue{}
	handler := New(Deps{Store: stub, Queue: queue})
	status := post(t, handler, body, map[string]string{
		"X-Gitlab-Token":      secret,
		"X-Gitlab-Event-UUID": "event-1",
	})
	if status != http.StatusOK || stub.job.Via != "Gitlab" || len(queue.queued) != 1 {
		t.Fatalf("status = %d, job = %+v, queued = %v", status, stub.job, queue.queued)
	}

	rejected := newStubStore("GITLAB")
	if status := post(t, New(Deps{Store: rejected, Queue: &stubQueue{}}), body,
		map[string]string{"X-Gitlab-Token": "wrong"}); status != http.StatusUnauthorized {
		t.Fatalf("wrong token status = %d, want 401", status)
	}
	if rejected.job.WorkspaceID != "" || len(rejected.deliveries) != 0 {
		t.Fatalf("wrong token created a job or delivery")
	}
}

func TestReplayedDeliveryStartsOneRun(t *testing.T) {
	stub := newStubStore("GITHUB")
	queue := &stubQueue{}
	handler := New(Deps{Store: stub, Queue: queue})
	headers := map[string]string{
		"X-GitHub-Event":      "push",
		"X-GitHub-Delivery":   "delivery-1",
		"X-Hub-Signature-256": githubSignature(testWSID, testPushBody),
	}
	if status := post(t, handler, testPushBody, headers); status != http.StatusOK {
		t.Fatalf("first status = %d", status)
	}
	if status := post(t, handler, testPushBody, headers); status != http.StatusOK {
		t.Fatalf("replay status = %d", status)
	}
	if len(queue.queued) != 1 {
		t.Fatalf("queued = %v, want one run", queue.queued)
	}
}

func TestUnmatchedRepositoryIsAcknowledgedWithoutARun(t *testing.T) {
	stub := newStubStore("GITHUB")
	stub.source = "https://github.com/acme/other.git"
	stub.workspace.Source = &stub.source
	queue := &stubQueue{}
	handler := New(Deps{Store: stub, Queue: queue})

	status := post(t, handler, testPushBody, map[string]string{
		"X-GitHub-Event":      "push",
		"X-GitHub-Delivery":   "delivery-2",
		"X-Hub-Signature-256": githubSignature(testWSID, testPushBody),
	})
	if status != http.StatusOK || len(queue.queued) != 0 {
		t.Fatalf("status = %d, queued = %v", status, queue.queued)
	}
}

func TestBranchAndPathMustMatchTheTrigger(t *testing.T) {
	for name, configure := range map[string]func(*stubStore){
		"branch": func(s *stubStore) { s.events[0].Branch = "release/.*" },
		"path":   func(s *stubStore) { s.events[0].Path = "*.md" },
		"event":  func(s *stubStore) { s.events[0].Event = "PULL_REQUEST" },
	} {
		t.Run(name, func(t *testing.T) {
			stub := newStubStore("GITHUB")
			configure(stub)
			queue := &stubQueue{}
			handler := New(Deps{Store: stub, Queue: queue})
			status := post(t, handler, testPushBody, map[string]string{
				"X-GitHub-Event":      "push",
				"X-GitHub-Delivery":   "delivery-3",
				"X-Hub-Signature-256": githubSignature(testWSID, testPushBody),
			})
			if status != http.StatusOK || len(queue.queued) != 0 {
				t.Fatalf("status = %d, queued = %v", status, queue.queued)
			}
		})
	}
}

func TestPingIsAcknowledgedWithoutARun(t *testing.T) {
	stub := newStubStore("GITHUB")
	body := `{"zen":"Keep it logically awesome."}`
	queue := &stubQueue{}
	handler := New(Deps{Store: stub, Queue: queue})
	status := post(t, handler, body, map[string]string{
		"X-GitHub-Event":      "ping",
		"X-Hub-Signature-256": githubSignature(testWSID, body),
	})
	if status != http.StatusOK || len(queue.queued) != 0 {
		t.Fatalf("status = %d, queued = %v", status, queue.queued)
	}
}
