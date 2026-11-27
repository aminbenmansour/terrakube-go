package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/pat"
)

type stubPATService struct {
	token       string
	tokens      []pat.Token
	err         error
	createdDays int
	createdDesc string
	createdUser auth.Identity
	listedEmail string
	deletedID   string
	deletedUser auth.Identity
}

func (s *stubPATService) Verify(context.Context, string) (*auth.Identity, error) {
	return nil, errors.New("not a PAT")
}

func (s *stubPATService) Create(_ context.Context, user auth.Identity, days int, description string) (string, error) {
	s.createdUser, s.createdDays, s.createdDesc = user, days, description
	return s.token, s.err
}

func (s *stubPATService) List(_ context.Context, email string) ([]pat.Token, error) {
	s.listedEmail = email
	return s.tokens, s.err
}

func (s *stubPATService) Delete(_ context.Context, id string, user auth.Identity) error {
	s.deletedID, s.deletedUser = id, user
	return s.err
}

func servePAT(t *testing.T, service *stubPATService, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	router := chi.NewRouter()
	router.Mount("/pat/v1", NewPAT(Deps{
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken: "dev-token",
		PATs:      service,
	}))
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer dev-token")
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestCreatePATReturnsTheTokenOnce(t *testing.T) {
	service := &stubPATService{token: "signed-jwt"}
	rec := servePAT(t, service, http.MethodPost, "/pat/v1", `{"days":30,"description":"ci"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	var response struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if response.Token != "signed-jwt" {
		t.Errorf("token = %q", response.Token)
	}
	if service.createdDays != 30 || service.createdDesc != "ci" || service.createdUser.Email != "dev@localhost" {
		t.Errorf("create args = %d/%q/%+v", service.createdDays, service.createdDesc, service.createdUser)
	}
}

func TestCreatePATRefusesANonInteractiveCaller(t *testing.T) {
	service := &stubPATService{err: pat.ErrNonInteractive}
	rec := servePAT(t, service, http.MethodPost, "/pat/v1", `{"days":1}`)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body.String())
	}
}

func TestListPATsShowsMetadataOnly(t *testing.T) {
	created := time.Date(2026, 11, 23, 0, 39, 9, 598000000, time.UTC)
	description := "ci"
	creator := "dev@localhost"
	service := &stubPATService{tokens: []pat.Token{{
		ID:          "11111111-1111-4111-8111-111111111111",
		Days:        int32Ptr(30),
		Description: &description,
		Source:      "API",
		CreatedDate: &created,
		CreatedBy:   &creator,
	}}}
	rec := servePAT(t, service, http.MethodGet, "/pat/v1", "")

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202; body = %s", rec.Code, rec.Body.String())
	}
	var response []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("items = %d, want 1", len(response))
	}
	item := response[0]
	if item["id"] != "11111111-1111-4111-8111-111111111111" || item["description"] != "ci" || item["createdDate"] != "2026-11-23T00:39:09.598Z" {
		t.Errorf("item = %v", item)
	}
	if _, ok := item["token"]; ok {
		t.Error("list returned a token")
	}
	if service.listedEmail != "dev@localhost" {
		t.Errorf("listed email = %q", service.listedEmail)
	}
}

func TestDeletePATAnswersAccepted(t *testing.T) {
	service := &stubPATService{}
	rec := servePAT(t, service, http.MethodDelete, "/pat/v1/11111111-1111-4111-8111-111111111111", "")

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202; body = %s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
	if service.deletedID != "11111111-1111-4111-8111-111111111111" {
		t.Errorf("deleted id = %q", service.deletedID)
	}
}

func TestDeletePATMapsErrors(t *testing.T) {
	cases := []struct {
		err  error
		want int
	}{
		{pat.ErrNotFound, http.StatusNotFound},
		{pat.ErrForbidden, http.StatusForbidden},
	}
	for _, tc := range cases {
		rec := servePAT(t, &stubPATService{err: tc.err}, http.MethodDelete, "/pat/v1/11111111-1111-4111-8111-111111111111", "")
		if rec.Code != tc.want {
			t.Errorf("Delete(%v) status = %d, want %d", tc.err, rec.Code, tc.want)
		}
	}
}

func TestPATRoutesRequireAuthenticationAndASecret(t *testing.T) {
	handler := chi.NewRouter()
	handler.Mount("/pat/v1", NewPAT(Deps{Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), AuthToken: "dev-token", PATs: &stubPATService{}}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/pat/v1", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", rec.Code)
	}

	disabled := chi.NewRouter()
	disabled.Mount("/pat/v1", NewPAT(Deps{Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), AuthToken: "dev-token"}))
	req := httptest.NewRequest(http.MethodGet, "/pat/v1", nil)
	req.Header.Set("Authorization", "Bearer dev-token")
	rec = httptest.NewRecorder()
	disabled.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("without a PAT secret status = %d, want 503", rec.Code)
	}
}

func int32Ptr(value int32) *int32 { return &value }

func TestPATRoutesRefuseOversizedBodies(t *testing.T) {
	router := NewPAT(Deps{
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken: "dev-token",
		PATs:      &stubPATService{},
	})
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
	request.Header.Set("Authorization", "Bearer dev-token")
	request.ContentLength = jsonapi.MaxBodyBytes + 1
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, request)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413: %s", rec.Code, rec.Body.String())
	}
}
