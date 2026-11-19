package executor

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

func TestExtractTarGzRejectsSymlinks(t *testing.T) {
	dest := t.TempDir()
	archive := tarGzEntries(t, []tarEntry{
		{name: "link", typeflag: tar.TypeSymlink, linkname: "/etc/passwd"},
		{name: "ok.txt", typeflag: tar.TypeReg, body: "ok"},
	})

	if err := extractTarGz(bytes.NewReader(archive), dest); err != nil {
		t.Fatalf("extract: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(dest, "link")); err == nil {
		t.Fatal("symlink entry was created")
	}
	if _, err := os.Stat(filepath.Join(dest, "ok.txt")); err != nil {
		t.Fatalf("regular file missing: %v", err)
	}
}

func TestExtractTarGzRejectsHardlinks(t *testing.T) {
	dest := t.TempDir()
	archive := tarGzEntries(t, []tarEntry{
		{name: "link", typeflag: tar.TypeLink, linkname: "/etc/passwd"},
	})

	if err := extractTarGz(bytes.NewReader(archive), dest); err != nil {
		t.Fatalf("extract: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(dest, "link")); err == nil {
		t.Fatal("hardlink entry was created")
	}
}

func TestExtractTarGzEnforcesFileSizeLimit(t *testing.T) {
	dest := t.TempDir()
	archive := tarGzEntries(t, []tarEntry{
		{name: "big.txt", typeflag: tar.TypeReg, body: "0123456789"},
	})

	err := extractTarGzWithLimits(bytes.NewReader(archive), dest, archiveLimits{
		MaxFileBytes:  4,
		MaxTotalBytes: 100,
		MaxEntries:    10,
	})
	if err == nil {
		t.Fatal("extract succeeded past the per-file size limit")
	}
	if _, statErr := os.Stat(filepath.Join(dest, "big.txt")); statErr == nil {
		t.Fatal("oversized file was left in the destination")
	}
}

func TestExtractTarGzEnforcesTotalSizeLimit(t *testing.T) {
	dest := t.TempDir()
	archive := tarGzEntries(t, []tarEntry{
		{name: "a.txt", typeflag: tar.TypeReg, body: "1234"},
		{name: "b.txt", typeflag: tar.TypeReg, body: "5678"},
	})

	err := extractTarGzWithLimits(bytes.NewReader(archive), dest, archiveLimits{
		MaxFileBytes:  10,
		MaxTotalBytes: 5,
		MaxEntries:    10,
	})
	if err == nil {
		t.Fatal("extract succeeded past the total size limit")
	}
}

func TestExtractTarGzEnforcesEntryLimit(t *testing.T) {
	dest := t.TempDir()
	archive := tarGzEntries(t, []tarEntry{
		{name: "a.txt", typeflag: tar.TypeReg, body: "a"},
		{name: "b.txt", typeflag: tar.TypeReg, body: "b"},
	})

	err := extractTarGzWithLimits(bytes.NewReader(archive), dest, archiveLimits{
		MaxFileBytes:  10,
		MaxTotalBytes: 100,
		MaxEntries:    1,
	})
	if err == nil {
		t.Fatal("extract succeeded past the entry limit")
	}
}

func TestGitSourceDetection(t *testing.T) {
	cases := []struct {
		name    string
		job     store.Job
		ws      store.Workspace
		wantOK  bool
		wantURL string
		wantRef string
	}{
		{
			name: "content tarball URL is not a git source",
			job:  store.Job{OverrideSource: ptr("https://host/remote/tfe/v2/configuration-versions/cv-1/terraformContent.tar.gz")},
			ws:   store.Workspace{Branch: ptr("main")},
		},
		{
			name: "storage key is not a git source",
			job:  store.Job{OverrideSource: ptr("content/cv-1/terraformContent.tar.gz")},
			ws:   store.Workspace{Branch: ptr("main")},
		},
		{
			name:    "git override with branch",
			job:     store.Job{OverrideSource: ptr("git::https://github.com/org/repo.git"), OverrideBranch: ptr("dev")},
			ws:      store.Workspace{Branch: ptr("main")},
			wantOK:  true,
			wantURL: "https://github.com/org/repo.git",
			wantRef: "dev",
		},
		{
			name:    "remote-content branch is ignored",
			job:     store.Job{OverrideSource: ptr("https://github.com/org/repo"), OverrideBranch: ptr("remote-content")},
			ws:      store.Workspace{Branch: ptr("main")},
			wantOK:  true,
			wantURL: "https://github.com/org/repo",
			wantRef: "main",
		},
		{
			name:    "workspace source is used when the job has no override",
			job:     store.Job{},
			ws:      store.Workspace{Source: ptr("https://github.com/org/repo.git"), Branch: ptr("main")},
			wantOK:  true,
			wantURL: "https://github.com/org/repo.git",
			wantRef: "main",
		},
		{
			name:    "commit wins over the workspace branch",
			job:     store.Job{CommitID: ptr("0123456789abcdef0123456789abcdef01234567")},
			ws:      store.Workspace{Source: ptr("git@github.com:org/repo.git"), Branch: ptr("main")},
			wantOK:  true,
			wantURL: "git@github.com:org/repo.git",
			wantRef: "0123456789abcdef0123456789abcdef01234567",
		},
		{
			name: "no source at all",
			job:  store.Job{},
			ws:   store.Workspace{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := gitSource(tc.job, tc.ws)
			if ok != tc.wantOK {
				t.Fatalf("gitSource ok = %v, want %v", ok, tc.wantOK)
			}
			if !ok {
				return
			}
			if got.URL != tc.wantURL {
				t.Errorf("URL = %q, want %q", got.URL, tc.wantURL)
			}
			if got.Ref != tc.wantRef {
				t.Errorf("Ref = %q, want %q", got.Ref, tc.wantRef)
			}
		})
	}
}

func TestSanitizeGitURL(t *testing.T) {
	cases := map[string]string{
		"https://x-access-token:secret@github.com/org/repo.git": "https://github.com/org/repo.git",
		"https://user:pass@example.com/org/repo.git":            "https://example.com/org/repo.git",
		"git@github.com:org/repo.git":                           "git@github.com:org/repo.git",
		"https://github.com/org/repo.git":                       "https://github.com/org/repo.git",
	}
	for input, want := range cases {
		if got := sanitizeGitURL(input); got != want {
			t.Errorf("sanitizeGitURL(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestCloneConfigurationFromLocalRepository(t *testing.T) {
	requireGit(t)
	repo := initGitRepo(t, "main.tf", `resource "terraform_data" "x" {}`)
	dest := t.TempDir()

	e := New(Deps{Logger: testLogger()})
	if err := e.cloneConfiguration(context.Background(), gitConfig{URL: repo, Ref: "main"}, dest); err != nil {
		t.Fatalf("clone: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "main.tf"))
	if err != nil {
		t.Fatalf("read cloned file: %v", err)
	}
	if !strings.Contains(string(data), "terraform_data") {
		t.Fatalf("cloned main.tf = %q", data)
	}
}

func TestCloneConfigurationChecksOutCommit(t *testing.T) {
	requireGit(t)
	repo := initGitRepo(t, "main.tf", "first")
	commit := gitOutput(t, repo, "rev-parse", "HEAD")
	if err := os.WriteFile(filepath.Join(repo, "main.tf"), []byte("second"), 0o640); err != nil {
		t.Fatalf("write second revision: %v", err)
	}
	gitOutput(t, repo, "add", "main.tf")
	gitOutput(t, repo, "commit", "-m", "second")

	dest := t.TempDir()
	e := New(Deps{Logger: testLogger()})
	if err := e.cloneConfiguration(context.Background(), gitConfig{URL: repo, Ref: commit}, dest); err != nil {
		t.Fatalf("clone commit: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "main.tf"))
	if err != nil {
		t.Fatalf("read cloned file: %v", err)
	}
	if string(data) != "first" {
		t.Fatalf("checked out %q, want first", data)
	}
}

func TestAcquireConfigurationClonesGitSource(t *testing.T) {
	requireGit(t)
	repo := initGitRepo(t, "main.tf", `resource "terraform_data" "x" {}`)
	dest := t.TempDir()

	e := New(Deps{Logger: testLogger()})
	ws := store.Workspace{Source: ptr(repo), Branch: ptr("main")}
	if err := e.acquireConfiguration(context.Background(), store.Job{}, ws, dest); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "main.tf")); err != nil {
		t.Fatalf("main.tf missing after clone: %v", err)
	}
}

func TestAcquireConfigurationFailsWithoutSource(t *testing.T) {
	e := New(Deps{Logger: testLogger()})
	err := e.acquireConfiguration(context.Background(), store.Job{}, store.Workspace{}, t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "no configuration source") {
		t.Fatalf("error = %v, want missing configuration source", err)
	}
}

func TestAcquireConfigurationRejectsInvalidArchive(t *testing.T) {
	backend := newMemoryBackend()
	layout := tfstate.NewLayout("local")
	contentID := "cv-bad"
	if err := backend.Put(context.Background(), layout.Content(contentID), strings.NewReader("not a gzip"), 10, "application/gzip"); err != nil {
		t.Fatalf("seed configuration: %v", err)
	}

	e := New(Deps{Storage: backend, Layout: layout, Logger: testLogger()})
	job := store.Job{OverrideSource: ptr("https://host/remote/tfe/v2/configuration-versions/" + contentID + "/terraformContent.tar.gz")}
	err := e.acquireConfiguration(context.Background(), job, store.Workspace{}, t.TempDir())
	if err == nil {
		t.Fatal("invalid archive was accepted")
	}
}

func TestAcquireConfigurationExtractsUploadedArchive(t *testing.T) {
	backend := newMemoryBackend()
	layout := tfstate.NewLayout("local")
	contentID := "cv-good"
	archive := tarGz(t, map[string]string{"main.tf": `resource "terraform_data" "x" {}`})
	if err := backend.Put(context.Background(), layout.Content(contentID), bytes.NewReader(archive), int64(len(archive)), "application/gzip"); err != nil {
		t.Fatalf("seed configuration: %v", err)
	}

	dest := t.TempDir()
	e := New(Deps{Storage: backend, Layout: layout, Logger: testLogger()})
	job := store.Job{OverrideSource: ptr("https://host/remote/tfe/v2/configuration-versions/" + contentID + "/terraformContent.tar.gz")}
	if err := e.acquireConfiguration(context.Background(), job, store.Workspace{}, dest); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "main.tf")); err != nil {
		t.Fatalf("main.tf missing: %v", err)
	}
}

type tarEntry struct {
	name     string
	typeflag byte
	linkname string
	body     string
}

func tarGzEntries(t *testing.T, entries []tarEntry) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, entry := range entries {
		header := &tar.Header{
			Name:     entry.name,
			Mode:     0o640,
			Size:     int64(len(entry.body)),
			Typeflag: entry.typeflag,
			Linkname: entry.linkname,
		}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("write header: %v", err)
		}
		if _, err := tw.Write([]byte(entry.body)); err != nil {
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

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
}

func initGitRepo(t *testing.T, file, content string) string {
	t.Helper()
	dir := t.TempDir()
	gitOutput(t, dir, "init", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, file), []byte(content), 0o640); err != nil {
		t.Fatalf("write %s: %v", file, err)
	}
	gitOutput(t, dir, "add", file)
	gitOutput(t, dir, "commit", "-m", "initial")
	return dir
}

func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}
