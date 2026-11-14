package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Local stores objects under a directory on disk. It is the development
// backend and mirrors upstream's local storage mode.
type Local struct {
	root string
}

// NewLocal creates the root directory if needed.
func NewLocal(root string) (*Local, error) {
	if root == "" {
		return nil, fmt.Errorf("storage: local root must not be empty")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("storage: resolve local root: %w", err)
	}
	if err := os.MkdirAll(abs, 0o750); err != nil {
		return nil, fmt.Errorf("storage: create local root: %w", err)
	}
	return &Local{root: abs}, nil
}

func (l *Local) Put(ctx context.Context, key string, r io.Reader, _ int64, _ string) error {
	path, err := l.path(key)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("storage: create object directory: %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), ".upload-*")
	if err != nil {
		return fmt.Errorf("storage: create temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if _, err := io.Copy(tmp, r); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("storage: write object: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("storage: sync object: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("storage: close object: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("storage: commit object: %w", err)
	}
	return nil
}

func (l *Local) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	path, err := l.path(key)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("storage: open object: %w", err)
	}
	return f, nil
}

func (l *Local) Delete(ctx context.Context, key string) error {
	path, err := l.path(key)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("storage: delete object: %w", err)
	}
	return nil
}

func (l *Local) Stat(ctx context.Context, key string) (Object, error) {
	path, err := l.path(key)
	if err != nil {
		return Object{}, err
	}
	if err := ctx.Err(); err != nil {
		return Object{}, err
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return Object{}, ErrNotFound
	}
	if err != nil {
		return Object{}, fmt.Errorf("storage: stat object: %w", err)
	}
	return Object{
		Key:        key,
		Size:       info.Size(),
		ModifiedAt: info.ModTime().UTC(),
	}, nil
}

func (l *Local) Health(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	info, err := os.Stat(l.root)
	if err != nil {
		return fmt.Errorf("storage: local root: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("storage: local root %s is not a directory", l.root)
	}
	return nil
}

func (l *Local) path(key string) (string, error) {
	clean, err := cleanKey(key)
	if err != nil {
		return "", err
	}
	path := filepath.Join(l.root, filepath.FromSlash(clean))
	// Defence in depth: ensure the joined path is still under root.
	rel, err := filepath.Rel(l.root, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("storage: key %q escapes the storage root", key)
	}
	return path, nil
}

func cleanKey(key string) (string, error) {
	key = strings.TrimSpace(strings.ReplaceAll(key, "\\", "/"))
	key = strings.TrimPrefix(key, "/")
	if key == "" {
		return "", fmt.Errorf("storage: key must not be empty")
	}
	clean := pathClean(key)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", fmt.Errorf("storage: invalid key %q", key)
	}
	return clean, nil
}

// pathClean is path.Clean without importing path just for one call site; it
// keeps the behaviour obvious.
func pathClean(key string) string {
	parts := strings.Split(key, "/")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		switch part {
		case "", ".":
			continue
		case "..":
			if len(out) == 0 {
				return ".."
			}
			out = out[:len(out)-1]
		default:
			out = append(out, part)
		}
	}
	return strings.Join(out, "/")
}

var _ Backend = (*Local)(nil)
