// Package storage abstracts the object store that holds state, run logs and
// uploaded configuration.
package storage

import (
	"context"
	"errors"
	"io"
	"time"
)

// ErrNotFound is returned when a key does not exist.
var ErrNotFound = errors.New("storage: object not found")

// Object describes a stored object without reading its body.
type Object struct {
	Key         string
	Size        int64
	ContentType string
	ModifiedAt  time.Time
}

// Backend is the minimal object-store surface Terrakube needs. Implementations
// must treat keys as slash-separated paths and must never allow escaping the
// configured root.
type Backend interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	Stat(ctx context.Context, key string) (Object, error)
	Health(ctx context.Context) error
}
