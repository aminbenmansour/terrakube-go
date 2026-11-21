package executor

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	internalauth "github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const backendOverrideFile = "terrakube_override.tf"

// preparedWorkspace is a job's isolated working tree and the environment it
// should run with. secrets are the values that must never reach a log.
type preparedWorkspace struct {
	workDir string
	runDir  string
	env     []string
	secrets []string
}

// prepareWorkspace downloads and extracts the run configuration, writes the
// HTTP backend override and builds the process environment.
func (e *Executor) prepareWorkspace(ctx context.Context, job store.Job, ws store.Workspace, step store.Step) (*preparedWorkspace, error) {
	workDir, err := os.MkdirTemp(e.deps.WorkDir, fmt.Sprintf("terrakube-job-%d-", job.ID))
	if err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}

	if err := e.acquireConfiguration(ctx, job, ws, workDir); err != nil {
		_ = os.RemoveAll(workDir)
		return nil, err
	}

	runDir := workDir
	if folder := firstFolder(ws.Folder); folder != "" {
		candidate := filepath.Join(workDir, folder)
		if info, statErr := os.Stat(candidate); statErr == nil && info.IsDir() {
			runDir = candidate
		}
	}

	if err := e.writeBackendOverride(runDir, job, ws); err != nil {
		_ = os.RemoveAll(workDir)
		return nil, err
	}

	env, secrets, err := e.environment(ctx, job, ws, step)
	if err != nil {
		_ = os.RemoveAll(workDir)
		return nil, err
	}
	return &preparedWorkspace{workDir: workDir, runDir: runDir, env: env, secrets: secrets}, nil
}

// gitConfig is a resolved Git source for a job.
type gitConfig struct {
	URL string
	Ref string
}

// acquireConfiguration populates dest from the job's Git source when it has
// one, and otherwise falls back to the uploaded configuration tarball.
func (e *Executor) acquireConfiguration(ctx context.Context, job store.Job, ws store.Workspace, dest string) error {
	if cfg, ok := gitSource(job, ws); ok {
		return e.cloneConfiguration(ctx, cfg, dest)
	}
	return e.downloadConfiguration(ctx, job, dest)
}

// gitSource resolves the repository and revision a job should build from.
// A run created by the CLI or a VCS webhook carries override_source, while a
// run created from the UI uses the workspace source. Uploaded configuration
// tarballs and object-storage keys are not Git sources.
func gitSource(job store.Job, ws store.Workspace) (gitConfig, bool) {
	source := ws.Source
	if job.OverrideSource != nil && *job.OverrideSource != "" && !isConfigurationTarball(*job.OverrideSource) {
		source = job.OverrideSource
	}
	if source == nil || strings.TrimSpace(*source) == "" {
		return gitConfig{}, false
	}

	ref := ""
	if ws.Branch != nil {
		ref = strings.TrimSpace(*ws.Branch)
	}
	// "remote-content" is the synthetic branch written when a run uploads its
	// configuration; it is not a real revision.
	if job.OverrideBranch != nil {
		if branch := strings.TrimSpace(*job.OverrideBranch); branch != "" && branch != "remote-content" {
			ref = branch
		}
	}
	if job.CommitID != nil && strings.TrimSpace(*job.CommitID) != "" {
		ref = strings.TrimSpace(*job.CommitID)
	}

	return gitConfig{URL: normalizeGitURL(*source), Ref: ref}, true
}

// normalizeGitURL strips the Terraform module source prefix some clients send.
func normalizeGitURL(source string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(source), "git::"))
}

func isConfigurationTarball(source string) bool {
	if strings.Contains(source, "/configuration-versions/") {
		return true
	}
	// Object-storage key written by run creation.
	return strings.HasPrefix(source, "content/") && strings.HasSuffix(source, "/terraformContent.tar.gz")
}

// sanitizeGitURL removes embedded credentials so they cannot reach logs.
func sanitizeGitURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.User == nil {
		return raw
	}
	parsed.User = nil
	return parsed.String()
}

func (e *Executor) cloneConfiguration(ctx context.Context, cfg gitConfig, dest string) error {
	clone := exec.CommandContext(ctx, "git", "clone", "--quiet", "--no-checkout", "--", cfg.URL, dest)
	clone.Env = gitEnv()
	if out, err := clone.CombinedOutput(); err != nil {
		return transient(fmt.Errorf("clone %s: %w: %s", sanitizeGitURL(cfg.URL), err, strings.TrimSpace(string(out))))
	}

	ref := cfg.Ref
	if ref == "" || ref == "HEAD" {
		if out, err := runGit(ctx, dest, "checkout", "--quiet"); err != nil {
			return fmt.Errorf("checkout default branch: %w: %s", err, out)
		}
		return nil
	}

	// A commit id is not a ref name, but Git accepts it as a fetch refspec
	// source. Shallow-fetch it first and fall back to a full fetch when the
	// server refuses.
	out, err := runGit(ctx, dest, "fetch", "--quiet", "--depth", "1", "origin", ref)
	if err != nil {
		fullOut, fullErr := runGit(ctx, dest, "fetch", "--quiet", "origin")
		if fullErr != nil {
			return transient(fmt.Errorf("fetch %s from %s: %w: %s", ref, sanitizeGitURL(cfg.URL), err, strings.TrimSpace(out+fullOut)))
		}
		if checkoutOut, checkoutErr := runGit(ctx, dest, "checkout", "--quiet", "--detach", ref); checkoutErr != nil {
			return fmt.Errorf("checkout %s: %w: %s", ref, checkoutErr, checkoutOut)
		}
		return nil
	}
	if checkoutOut, checkoutErr := runGit(ctx, dest, "checkout", "--quiet", "--detach", "FETCH_HEAD"); checkoutErr != nil {
		return fmt.Errorf("checkout %s: %w: %s", ref, checkoutErr, checkoutOut)
	}
	return nil
}

// runGit runs a git command in dir and returns its combined output.
func runGit(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Env = gitEnv()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// gitEnv disables interactive credential prompts so a missing credential fails
// the job instead of hanging the executor.
func gitEnv() []string {
	return append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_ASKPASS=true")
}

func (e *Executor) downloadConfiguration(ctx context.Context, job store.Job, dest string) error {
	if contentID := contentIDFromSource(job.OverrideSource); contentID != "" && e.deps.Storage != nil {
		rc, err := e.deps.Storage.Get(ctx, e.deps.Layout.Content(contentID))
		if err == nil {
			defer rc.Close()
			return extractTarGz(rc, dest)
		}
		if !errors.Is(err, storage.ErrNotFound) {
			return transient(fmt.Errorf("read configuration tarball: %w", err))
		}
	}

	if job.OverrideSource == nil || *job.OverrideSource == "" {
		return errors.New("job has no configuration source")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, *job.OverrideSource, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return transient(fmt.Errorf("download configuration: %w", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return transient(fmt.Errorf("download configuration: unexpected status %s", resp.Status))
	}
	return extractTarGz(resp.Body, dest)
}

// archiveLimits bounds how much a configuration archive may write. The limits
// protect the executor host from malformed or hostile uploads.
type archiveLimits struct {
	MaxFileBytes  int64
	MaxTotalBytes int64
	MaxEntries    int
}

var defaultArchiveLimits = archiveLimits{
	MaxFileBytes:  128 << 20, // 128 MiB per file
	MaxTotalBytes: 1 << 30,   // 1 GiB per archive
	MaxEntries:    100000,
}

func extractTarGz(r io.Reader, dest string) error {
	return extractTarGzWithLimits(r, dest, defaultArchiveLimits)
}

func extractTarGzWithLimits(r io.Reader, dest string, limits archiveLimits) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("open configuration archive: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	entries := 0
	var total int64
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read configuration archive: %w", err)
		}

		entries++
		if limits.MaxEntries > 0 && entries > limits.MaxEntries {
			return fmt.Errorf("configuration archive has more than %d entries", limits.MaxEntries)
		}

		name := filepath.Clean(header.Name)
		if name == "." || filepath.IsAbs(name) || strings.HasPrefix(name, ".."+string(filepath.Separator)) || name == ".." {
			continue
		}
		target := filepath.Join(dest, name)
		if !strings.HasPrefix(target, filepath.Clean(dest)+string(filepath.Separator)) && target != filepath.Clean(dest) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o750); err != nil {
				return err
			}
		case tar.TypeReg:
			if limits.MaxFileBytes > 0 && header.Size > limits.MaxFileBytes {
				_ = os.Remove(target)
				return fmt.Errorf("configuration file %s exceeds %d bytes", name, limits.MaxFileBytes)
			}
			if limits.MaxTotalBytes > 0 {
				total += header.Size
				if total > limits.MaxTotalBytes {
					_ = os.Remove(target)
					return fmt.Errorf("configuration archive exceeds %d bytes", limits.MaxTotalBytes)
				}
			}
			if err := os.MkdirAll(filepath.Dir(target), 0o750); err != nil {
				return err
			}
			mode := os.FileMode(header.Mode).Perm()
			if mode == 0 {
				mode = 0o640
			}
			file, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, tr); err != nil {
				_ = file.Close()
				return err
			}
			if err := file.Close(); err != nil {
				return err
			}
		case tar.TypeSymlink, tar.TypeLink:
			// Symlinks and hardlinks could escape the workspace or point at
			// host files, so they are dropped rather than extracted.
			continue
		}
	}
}

func (e *Executor) writeBackendOverride(runDir string, job store.Job, ws store.Workspace) error {
	address := fmt.Sprintf("%s/tfstate/v1/organization/%s/workspace/%s/backend/terraform.tfstate",
		strings.TrimRight(e.deps.PublicURL, "/"), job.OrganizationID, job.WorkspaceID)
	content := fmt.Sprintf("terraform {\n  backend \"http\" {\n    address = %q\n  }\n}\n", address)
	return os.WriteFile(filepath.Join(runDir, backendOverrideFile), []byte(content), 0o640)
}

func (e *Executor) environment(ctx context.Context, job store.Job, ws store.Workspace, step store.Step) ([]string, []string, error) {
	env := []string{
		"TF_IN_AUTOMATION=true",
		"TF_INPUT=0",
		"CHECKPOINT_DISABLE=1",
	}
	var secrets []string

	if e.deps.InternalTokenSecret != "" {
		token, err := internalauth.Sign(e.deps.InternalTokenSecret, internalauth.Claims{
			OrganizationID: job.OrganizationID,
			WorkspaceID:    job.WorkspaceID,
			JobID:          fmt.Sprintf("%d", job.ID),
			StepID:         step.ID,
			ExpiresAt:      jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		})
		if err != nil {
			return nil, nil, fmt.Errorf("mint workspace token: %w", err)
		}
		env = append(env,
			"TF_HTTP_USERNAME=terraform",
			"TF_HTTP_PASSWORD="+token,
		)
		secrets = append(secrets, token)
	}

	variables, err := e.deps.Store.ListVariablesByWorkspace(ctx, ws.ID)
	if err != nil {
		return nil, nil, transient(fmt.Errorf("load workspace variables: %w", err))
	}
	// A workspace variable shadows an organization global with the same key in
	// the same category; upstream's loadDefault uses putIfAbsent per category.
	defined := map[string]bool{}
	for _, variable := range variables {
		if variable.Incomplete {
			continue
		}
		defined[strings.ToUpper(variable.VariableCategory)+"\x00"+variable.VariableKey] = true
		if variable.Sensitive && variable.VariableValue != "" {
			secrets = append(secrets, variable.VariableValue)
		}
		switch strings.ToUpper(variable.VariableCategory) {
		case "ENV":
			env = append(env, variable.VariableKey+"="+variable.VariableValue)
		default:
			env = append(env, "TF_VAR_"+variable.VariableKey+"="+variable.VariableValue)
		}
	}

	globals, err := e.deps.Store.ListGlobalVariablesByOrganization(ctx, job.OrganizationID)
	if err != nil {
		return nil, nil, transient(fmt.Errorf("load global variables: %w", err))
	}
	for _, variable := range globals {
		if defined[variable.VariableCategory+"\x00"+variable.VariableKey] {
			continue
		}
		if variable.Sensitive && variable.VariableValue != "" {
			secrets = append(secrets, variable.VariableValue)
		}
		switch variable.VariableCategory {
		case "ENV":
			env = append(env, variable.VariableKey+"="+variable.VariableValue)
		default:
			env = append(env, "TF_VAR_"+variable.VariableKey+"="+variable.VariableValue)
		}
	}
	return env, secrets, nil
}

// cleanupStaleWorkspaces removes job directories left behind by a crashed or
// killed executor. Recent directories are kept because another executor may
// still be using them.
func cleanupStaleWorkspaces(root string, maxAge time.Duration) error {
	if root == "" || maxAge <= 0 {
		return nil
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	cutoff := time.Now().Add(-maxAge)
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "terrakube-job-") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(cutoff) {
			continue
		}
		if err := os.RemoveAll(filepath.Join(root, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func contentIDFromSource(source *string) string {
	if source == nil {
		return ""
	}
	const marker = "/configuration-versions/"
	idx := strings.Index(*source, marker)
	if idx < 0 {
		return ""
	}
	rest := (*source)[idx+len(marker):]
	end := strings.IndexByte(rest, '/')
	if end <= 0 {
		return ""
	}
	return rest[:end]
}

func firstFolder(folder *string) string {
	if folder == nil {
		return ""
	}
	first := strings.TrimSpace(strings.Split(*folder, ",")[0])
	first = strings.TrimPrefix(first, "/")
	first = strings.TrimSuffix(first, "/")
	if first == "." {
		return ""
	}
	return first
}
