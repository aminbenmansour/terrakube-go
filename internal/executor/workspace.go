package executor

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	internalauth "github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const backendOverrideFile = "terrakube_override.tf"

// prepareWorkspace downloads and extracts the run configuration, writes the
// HTTP backend override and builds the process environment.
func (e *Executor) prepareWorkspace(ctx context.Context, job store.Job, ws store.Workspace, step store.Step) (workDir, runDir string, env []string, err error) {
	workDir, err = os.MkdirTemp(e.deps.WorkDir, fmt.Sprintf("terrakube-job-%d-", job.ID))
	if err != nil {
		return "", "", nil, fmt.Errorf("create workspace: %w", err)
	}

	if err := e.downloadConfiguration(ctx, job, workDir); err != nil {
		_ = os.RemoveAll(workDir)
		return "", "", nil, err
	}

	runDir = workDir
	if folder := firstFolder(ws.Folder); folder != "" {
		candidate := filepath.Join(workDir, folder)
		if info, statErr := os.Stat(candidate); statErr == nil && info.IsDir() {
			runDir = candidate
		}
	}

	if err := e.writeBackendOverride(runDir, job, ws); err != nil {
		_ = os.RemoveAll(workDir)
		return "", "", nil, err
	}

	env, err = e.environment(ctx, job, ws, step)
	if err != nil {
		_ = os.RemoveAll(workDir)
		return "", "", nil, err
	}
	return workDir, runDir, env, nil
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

func extractTarGz(r io.Reader, dest string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("open configuration archive: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read configuration archive: %w", err)
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
		}
	}
}

func (e *Executor) writeBackendOverride(runDir string, job store.Job, ws store.Workspace) error {
	address := fmt.Sprintf("%s/tfstate/v1/organization/%s/workspace/%s/backend/terraform.tfstate",
		strings.TrimRight(e.deps.PublicURL, "/"), job.OrganizationID, job.WorkspaceID)
	content := fmt.Sprintf("terraform {\n  backend \"http\" {\n    address = %q\n  }\n}\n", address)
	return os.WriteFile(filepath.Join(runDir, backendOverrideFile), []byte(content), 0o640)
}

func (e *Executor) environment(ctx context.Context, job store.Job, ws store.Workspace, step store.Step) ([]string, error) {
	env := []string{
		"TF_IN_AUTOMATION=true",
		"TF_INPUT=0",
		"CHECKPOINT_DISABLE=1",
	}

	if e.deps.InternalTokenSecret != "" {
		token, err := internalauth.Sign(e.deps.InternalTokenSecret, internalauth.Claims{
			OrganizationID: job.OrganizationID,
			WorkspaceID:    job.WorkspaceID,
			JobID:          fmt.Sprintf("%d", job.ID),
			StepID:         step.ID,
			ExpiresAt:      jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		})
		if err != nil {
			return nil, fmt.Errorf("mint workspace token: %w", err)
		}
		env = append(env,
			"TF_HTTP_USERNAME=terraform",
			"TF_HTTP_PASSWORD="+token,
		)
	}

	variables, err := e.deps.Store.ListVariablesByWorkspace(ctx, ws.ID)
	if err != nil {
		return nil, transient(fmt.Errorf("load workspace variables: %w", err))
	}
	for _, variable := range variables {
		if variable.Incomplete {
			continue
		}
		switch strings.ToUpper(variable.VariableCategory) {
		case "ENV":
			env = append(env, variable.VariableKey+"="+variable.VariableValue)
		default:
			env = append(env, "TF_VAR_"+variable.VariableKey+"="+variable.VariableValue)
		}
	}
	return env, nil
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
