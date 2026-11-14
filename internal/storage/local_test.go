package storage

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestLocalLifecycle(t *testing.T) {
	ctx := context.Background()
	store, err := NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}

	if err := store.Put(ctx, "org/ws/terraform.tfstate", strings.NewReader(`{"version":4}`), 13, "application/json"); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	obj, err := store.Stat(ctx, "org/ws/terraform.tfstate")
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if obj.Size != 13 {
		t.Errorf("Stat().Size = %d, want 13", obj.Size)
	}

	rc, err := store.Get(ctx, "org/ws/terraform.tfstate")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	body, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	_ = rc.Close()
	if string(body) != `{"version":4}` {
		t.Errorf("body = %q", body)
	}

	if err := store.Delete(ctx, "org/ws/terraform.tfstate"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := store.Get(ctx, "org/ws/terraform.tfstate"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get() after delete error = %v, want ErrNotFound", err)
	}
	if err := store.Delete(ctx, "org/ws/terraform.tfstate"); err != nil {
		t.Errorf("Delete() missing object error = %v, want nil", err)
	}
}

func TestLocalRejectsPathTraversal(t *testing.T) {
	store, err := NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	ctx := context.Background()

	keys := []string{"", "..", "../etc/passwd", "a/../../b", "/"}
	for _, key := range keys {
		if err := store.Put(ctx, key, strings.NewReader("x"), 1, ""); err == nil {
			t.Errorf("Put(%q) error = nil, want rejection", key)
		}
	}
}

func TestLocalHealth(t *testing.T) {
	store, err := NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	if err := store.Health(context.Background()); err != nil {
		t.Errorf("Health() error = %v", err)
	}
}

func TestCleanKey(t *testing.T) {
	cases := map[string]string{
		"/org/ws/state":       "org/ws/state",
		"org//ws/./state":     "org/ws/state",
		"org/ws/../ws2/state": "org/ws2/state",
	}
	for input, want := range cases {
		got, err := cleanKey(input)
		if err != nil {
			t.Fatalf("cleanKey(%q) error = %v", input, err)
		}
		if got != want {
			t.Errorf("cleanKey(%q) = %q, want %q", input, got, want)
		}
	}
}
