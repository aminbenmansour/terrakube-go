package executor

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

func TestExtractTarGz(t *testing.T) {
	dest := t.TempDir()
	archive := tarGz(t, map[string]string{
		"main.tf":        `resource "terraform_data" "x" {}`,
		"modules/vpc.tf": `variable "cidr" {}`,
	})

	if err := extractTarGz(bytes.NewReader(archive), dest); err != nil {
		t.Fatalf("extract: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "main.tf"))
	if err != nil {
		t.Fatalf("read main.tf: %v", err)
	}
	if !strings.Contains(string(data), "terraform_data") {
		t.Errorf("main.tf = %q", data)
	}
	if _, err := os.Stat(filepath.Join(dest, "modules", "vpc.tf")); err != nil {
		t.Errorf("nested file: %v", err)
	}
}

func TestExtractTarGzRejectsPathTraversal(t *testing.T) {
	dest := t.TempDir()
	parent := filepath.Dir(dest)
	evil := filepath.Join(parent, "evil.txt")
	_ = os.Remove(evil)

	archive := tarGz(t, map[string]string{
		"../evil.txt": "escaped",
		"ok.txt":      "ok",
	})
	if err := extractTarGz(bytes.NewReader(archive), dest); err != nil {
		t.Fatalf("extract: %v", err)
	}
	if _, err := os.Stat(evil); err == nil {
		t.Fatalf("archive escaped the destination: %s", evil)
	}
	if _, err := os.Stat(filepath.Join(dest, "ok.txt")); err != nil {
		t.Errorf("safe entry missing: %v", err)
	}
}

func TestContentIDFromSource(t *testing.T) {
	cases := []struct {
		source *string
		want   string
	}{
		{nil, ""},
		{ptr(""), ""},
		{ptr("https://host/remote/tfe/v2/configuration-versions/cv-1/terraformContent.tar.gz"), "cv-1"},
		{ptr("https://host/remote/tfe/v2/configuration-versions/cv-2"), ""},
		{ptr("git::https://example.com/repo.git"), ""},
	}
	for _, tc := range cases {
		if got := contentIDFromSource(tc.source); got != tc.want {
			t.Errorf("contentIDFromSource(%v) = %q, want %q", tc.source, got, tc.want)
		}
	}
}

func TestFirstFolder(t *testing.T) {
	cases := []struct {
		folder *string
		want   string
	}{
		{nil, ""},
		{ptr(""), ""},
		{ptr("."), ""},
		{ptr("infra"), "infra"},
		{ptr("/infra/"), "infra"},
		{ptr("infra,other"), "infra"},
	}
	for _, tc := range cases {
		if got := firstFolder(tc.folder); got != tc.want {
			t.Errorf("firstFolder(%v) = %q, want %q", tc.folder, got, tc.want)
		}
	}
}

func TestIsTerminalStatus(t *testing.T) {
	for _, status := range []string{"completed", "failed", "cancelled", "canceled", "noChanges", "rejected", "notExecuted", "NeverExecuted"} {
		if !isTerminalStatus(status) {
			t.Errorf("isTerminalStatus(%q) = false", status)
		}
	}
	for _, status := range []string{"pending", "running", "queue", "waitingApproval", "approved", ""} {
		if isTerminalStatus(status) {
			t.Errorf("isTerminalStatus(%q) = true", status)
		}
	}
}

func TestLogSinkSplitsLinesAndKeepsTail(t *testing.T) {
	sink := newLogSink(t.Context(), nil, nil, 1, "step-1")
	if _, err := sink.Write([]byte("first\nsec")); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := sink.Write([]byte("ond\n")); err != nil {
		t.Fatalf("write: %v", err)
	}
	sink.WriteString("trailing")
	sink.Flush()

	if got, want := sink.String(), "first\nsecond\ntrailing"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestBinarySelection(t *testing.T) {
	exec := &Executor{}
	if got := exec.binary(store.Workspace{IacType: "opentofu"}); got != "tofu" {
		t.Errorf("opentofu binary = %q", got)
	}
	exec = &Executor{deps: Deps{Binary: "/opt/tofu"}}
	if got := exec.binary(store.Workspace{IacType: "terraform"}); got != "/opt/tofu" {
		t.Errorf("configured binary = %q", got)
	}
}

func TestWriteBackendOverride(t *testing.T) {
	dir := t.TempDir()
	e := New(Deps{PublicURL: "https://terrakube.example.com/"})
	job := store.Job{ID: 9, OrganizationID: "org-1", WorkspaceID: "ws-1"}
	if err := e.writeBackendOverride(dir, job, store.Workspace{ID: "ws-1"}); err != nil {
		t.Fatalf("write backend override: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, backendOverrideFile))
	if err != nil {
		t.Fatalf("read override: %v", err)
	}
	want := "https://terrakube.example.com/tfstate/v1/organization/org-1/workspace/ws-1/backend/terraform.tfstate"
	if !strings.Contains(string(data), want) {
		t.Fatalf("override = %s, want address %s", data, want)
	}
	if !strings.Contains(string(data), `backend "http"`) {
		t.Fatalf("override does not configure the http backend: %s", data)
	}
}

func TestArchiveHistory(t *testing.T) {
	ctx := context.Background()
	backend := newMemoryBackend()
	layout := tfstate.NewLayout("local")
	state := []byte(`{"version":4,"serial":7,"lineage":"lineage-1","terraform_version":"1.9.0"}`)
	key := layout.CurrentState("org-1", "ws-1")
	if err := backend.Put(ctx, key, bytes.NewReader(state), int64(len(state)), "application/json"); err != nil {
		t.Fatalf("seed state: %v", err)
	}

	q := &historyStub{}
	e := New(Deps{
		Store:   q,
		Storage: backend,
		Layout:  layout,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	job := store.Job{ID: 9, OrganizationID: "org-1", WorkspaceID: "ws-1"}
	ws := store.Workspace{ID: "ws-1"}
	if err := e.archiveHistory(ctx, job, ws); err != nil {
		t.Fatalf("archive history: %v", err)
	}

	if q.created == nil {
		t.Fatal("no history row was created")
	}
	if q.created.Serial != 7 || q.created.Lineage != "lineage-1" {
		t.Errorf("history serial/lineage = %d/%q", q.created.Serial, q.created.Lineage)
	}
	sum := md5.Sum(state)
	if q.created.Md5 != hex.EncodeToString(sum[:]) {
		t.Errorf("history md5 = %q", q.created.Md5)
	}
	if q.created.JobReference == nil || *q.created.JobReference != "9" {
		t.Errorf("history job reference = %v", q.created.JobReference)
	}
	for _, archiveKey := range []string{
		layout.JSONState("org-1", "ws-1", q.created.ID),
		layout.RawState("org-1", "ws-1", q.created.ID),
	} {
		if _, ok := backend.objects[archiveKey]; !ok {
			t.Errorf("missing archived object %q", archiveKey)
		}
	}
}

func TestArchiveHistoryWithoutStateIsNoop(t *testing.T) {
	q := &historyStub{}
	e := New(Deps{
		Store:   q,
		Storage: newMemoryBackend(),
		Layout:  tfstate.NewLayout("local"),
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	if err := e.archiveHistory(context.Background(), store.Job{ID: 9}, store.Workspace{ID: "ws-1"}); err != nil {
		t.Fatalf("archive history: %v", err)
	}
	if q.created != nil {
		t.Error("history row created without state")
	}
}

type historyStub struct {
	store.Querier
	created *store.CreateHistoryParams
}

func (s *historyStub) CreateHistory(_ context.Context, arg store.CreateHistoryParams) (store.History, error) {
	s.created = &arg
	return store.History{ID: arg.ID}, nil
}

func (s *historyStub) GetHistoryByJobReference(context.Context, store.GetHistoryByJobReferenceParams) (store.History, error) {
	return store.History{}, pgx.ErrNoRows
}

type memoryBackend struct {
	objects map[string][]byte
}

func newMemoryBackend() *memoryBackend {
	return &memoryBackend{objects: map[string][]byte{}}
}

func (m *memoryBackend) Put(_ context.Context, key string, r io.Reader, _ int64, _ string) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	m.objects[key] = data
	return nil
}

func (m *memoryBackend) Get(_ context.Context, key string) (io.ReadCloser, error) {
	data, ok := m.objects[key]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (m *memoryBackend) Delete(_ context.Context, key string) error {
	delete(m.objects, key)
	return nil
}

func (m *memoryBackend) Stat(_ context.Context, key string) (storage.Object, error) {
	data, ok := m.objects[key]
	if !ok {
		return storage.Object{}, storage.ErrNotFound
	}
	return storage.Object{Key: key, Size: int64(len(data))}, nil
}

func (m *memoryBackend) Health(context.Context) error { return nil }

var _ storage.Backend = (*memoryBackend)(nil)

func tarGz(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, content := range files {
		header := &tar.Header{
			Name:     name,
			Mode:     0o640,
			Size:     int64(len(content)),
			Typeflag: tar.TypeReg,
		}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("write header: %v", err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatalf("write body: %v", err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return buf.Bytes()
}

func ptr(value string) *string { return &value }
