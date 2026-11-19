package tfoutput

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

type stubQuerier struct {
	store.Querier
	job  store.Job
	step store.Step
}

func (s *stubQuerier) GetJob(_ context.Context, id int64) (store.Job, error) {
	if s.job.ID != id {
		return store.Job{}, pgx.ErrNoRows
	}
	return s.job, nil
}

func (s *stubQuerier) GetStep(_ context.Context, id string) (store.Step, error) {
	if s.step.ID != id {
		return store.Step{}, pgx.ErrNoRows
	}
	return s.step, nil
}

func newHandler(q store.Querier, token string) http.Handler {
	return New(Deps{Store: q, AuthToken: token})
}

func TestGetStepLogReturnsArchivedOutput(t *testing.T) {
	output := "Executing terraform plan\nPlan: 1 to add\n"
	q := &stubQuerier{
		job:  store.Job{ID: 7, OrganizationID: "org-1"},
		step: store.Step{ID: "step-1", JobID: 7, Output: &output, Status: stringPtr("completed")},
	}
	handler := newHandler(q, "secret")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/organization/org-1/job/7/step/step-1", nil)
	req.Header.Set("Authorization", "Bearer secret")
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != output {
		t.Fatalf("body = %q, want %q", rec.Body.String(), output)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/octet-stream" {
		t.Errorf("content type = %q", got)
	}
}

func TestGetStepLogRejectsWrongScope(t *testing.T) {
	output := "log"
	q := &stubQuerier{
		job:  store.Job{ID: 7, OrganizationID: "org-1"},
		step: store.Step{ID: "step-1", JobID: 7, Output: &output},
	}
	handler := newHandler(q, "")

	cases := []struct {
		name string
		path string
	}{
		{"unknown organization", "/organization/org-2/job/7/step/step-1"},
		{"unknown job", "/organization/org-1/job/8/step/step-1"},
		{"unknown step", "/organization/org-1/job/7/step/step-9"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404", rec.Code)
			}
		})
	}
}

func TestGetStepLogRequiresToken(t *testing.T) {
	output := "log"
	q := &stubQuerier{
		job:  store.Job{ID: 7, OrganizationID: "org-1"},
		step: store.Step{ID: "step-1", JobID: 7, Output: &output},
	}
	handler := newHandler(q, "secret")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/org-1/job/7/step/step-1", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/organization/org-1/job/7/step/step-1", nil)
	req.Header.Set("Authorization", "Bearer secret")
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestGetStepLogWithoutOutputIsNotFound(t *testing.T) {
	q := &stubQuerier{
		job:  store.Job{ID: 7, OrganizationID: "org-1"},
		step: store.Step{ID: "step-1", JobID: 7, Status: stringPtr("completed")},
	}
	rec := httptest.NewRecorder()
	newHandler(q, "").ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/org-1/job/7/step/step-1", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestStreamWithoutQueueIsUnavailable(t *testing.T) {
	q := &stubQuerier{
		job:  store.Job{ID: 7, OrganizationID: "org-1"},
		step: store.Step{ID: "step-1", JobID: 7},
	}
	rec := httptest.NewRecorder()
	newHandler(q, "").ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/org-1/job/7/step/step-1/stream", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
}

func TestWriteDataSplitsMultilineOutput(t *testing.T) {
	rec := httptest.NewRecorder()
	writeData(rec, "1-0", "first\nsecond\n")
	want := "id: 1-0\ndata: first\ndata: second\n\n"
	if rec.Body.String() != want {
		t.Fatalf("frame = %q, want %q", rec.Body.String(), want)
	}
}

func TestParseResumeID(t *testing.T) {
	cases := map[string]string{
		"":             "0-0",
		"*":            "0-0",
		"bad":          "0-0",
		"12":           "0-0",
		"1-2":          "1-2",
		" 17-0 ":       "17-0",
		"1795019731-4": "1795019731-4",
	}
	for input, want := range cases {
		if got := parseResumeID(input); got != want {
			t.Errorf("parseResumeID(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestCompareStreamIDs(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1-0", "1-0", 0},
		{"1-0", "2-0", -1},
		{"2-0", "1-9", 1},
		{"1-1", "1-2", -1},
		{"1-2", "1-1", 1},
		{"bad", "1-1", 0},
	}
	for _, tc := range cases {
		if got := compareStreamIDs(tc.a, tc.b); got != tc.want {
			t.Errorf("compareStreamIDs(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestIsTerminal(t *testing.T) {
	for _, status := range []string{"completed", "failed", "cancelled", "noChanges", "rejected", "neverExecuted"} {
		if !isTerminal(stringPtr(status)) {
			t.Errorf("isTerminal(%q) = false", status)
		}
	}
	for _, status := range []string{"pending", "running", "queue", "waitingApproval", ""} {
		if isTerminal(stringPtr(status)) {
			t.Errorf("isTerminal(%q) = true", status)
		}
	}
	if isTerminal(nil) {
		t.Error("isTerminal(nil) = true")
	}
}

func TestWriteEventFrame(t *testing.T) {
	rec := httptest.NewRecorder()
	writeEvent(rec, "end", `{"reason":"complete"}`)
	if !strings.Contains(rec.Body.String(), "event: end\n") || !strings.Contains(rec.Body.String(), `data: {"reason":"complete"}`) {
		t.Fatalf("frame = %q", rec.Body.String())
	}
}

func TestGetStepLogPrefersArchivedObject(t *testing.T) {
	inline := "inline text"
	q := &stubQuerier{
		job:  store.Job{ID: 7, OrganizationID: "org-1"},
		step: store.Step{ID: "step-1", JobID: 7, Output: &inline, Status: stringPtr("completed")},
	}
	backend := &logBackend{objects: map[string][]byte{}}
	layout := tfstate.NewLayout("local")
	backend.objects[layout.Output("org-1", 7, "step-1")] = []byte("archived text")

	rec := httptest.NewRecorder()
	handler := New(Deps{Store: q, Storage: backend, Layout: layout})
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/org-1/job/7/step/step-1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != "archived text" {
		t.Fatalf("body = %q, want archived text", rec.Body.String())
	}
}

func TestGetStepLogFallsBackToInlineOutput(t *testing.T) {
	inline := "inline text"
	q := &stubQuerier{
		job:  store.Job{ID: 7, OrganizationID: "org-1"},
		step: store.Step{ID: "step-1", JobID: 7, Output: &inline, Status: stringPtr("completed")},
	}
	handler := New(Deps{Store: q, Storage: &logBackend{objects: map[string][]byte{}}, Layout: tfstate.NewLayout("local")})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/org-1/job/7/step/step-1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != inline {
		t.Fatalf("body = %q, want inline text", rec.Body.String())
	}
}

type logBackend struct{ objects map[string][]byte }

func (b *logBackend) Put(_ context.Context, key string, r io.Reader, _ int64, _ string) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	b.objects[key] = data
	return nil
}

func (b *logBackend) Get(_ context.Context, key string) (io.ReadCloser, error) {
	data, ok := b.objects[key]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (b *logBackend) Delete(_ context.Context, key string) error { delete(b.objects, key); return nil }
func (b *logBackend) Stat(_ context.Context, key string) (storage.Object, error) {
	data, ok := b.objects[key]
	if !ok {
		return storage.Object{}, storage.ErrNotFound
	}
	return storage.Object{Key: key, Size: int64(len(data))}, nil
}
func (b *logBackend) Health(context.Context) error { return nil }

var _ storage.Backend = (*logBackend)(nil)

func stringPtr(value string) *string { return &value }
