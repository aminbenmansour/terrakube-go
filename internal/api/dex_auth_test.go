package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
)

type stubDexVerifier struct {
	token string
	err   error
}

func (s stubDexVerifier) Verify(_ context.Context, raw string) (*auth.Identity, error) {
	if s.err != nil {
		return nil, s.err
	}
	if raw != s.token {
		return nil, errors.New("invalid token")
	}
	return &auth.Identity{Issuer: "dex", Email: "user@example.com"}, nil
}

func TestAuthMiddlewareAcceptsDexTokens(t *testing.T) {
	handler := New(Deps{
		Store:     &stubQuerier{},
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		Verifiers: []auth.Verifier{stubDexVerifier{token: "dex-token"}},
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"errors"`) {
		t.Errorf("401 body = %s, want a JSON:API errors document", rec.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/organization", nil)
	req.Header.Set("Authorization", "Bearer not-the-token")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("with a rejected Dex token status = %d, want 401", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/organization", nil)
	req.Header.Set("Authorization", "Bearer dex-token")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("with a valid Dex token status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
}

func TestAuthMiddlewareAcceptsEitherCredential(t *testing.T) {
	handler := New(Deps{
		Store:     &stubQuerier{},
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken: "dev-token",
		Verifiers: []auth.Verifier{stubDexVerifier{token: "dex-token"}},
	})

	for _, token := range []string{"dev-token", "dex-token"} {
		req := httptest.NewRequest(http.MethodGet, "/organization", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("token %q status = %d, want 200", token, rec.Code)
		}
	}
}
