package clilogin

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/pat"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type fakeStore struct {
	sessions map[string]store.CliAuthSession
}

func newFakeStore() *fakeStore { return &fakeStore{sessions: map[string]store.CliAuthSession{}} }

func (f *fakeStore) CreateCliAuthSession(_ context.Context, arg store.CreateCliAuthSessionParams) (store.CliAuthSession, error) {
	row := store.CliAuthSession{
		ID:               arg.ID,
		Status:           "PENDING_IDP",
		CliRedirectUri:   arg.CliRedirectUri,
		CliCodeChallenge: arg.CliCodeChallenge,
		CliState:         arg.CliState,
		DexCodeVerifier:  arg.DexCodeVerifier,
		ExpiresAt:        arg.ExpiresAt,
	}
	f.sessions[row.ID] = row
	return row, nil
}

func (f *fakeStore) GetCliAuthSession(_ context.Context, id string) (store.CliAuthSession, error) {
	row, ok := f.sessions[id]
	if !ok {
		return store.CliAuthSession{}, pgErrNoRows
	}
	return row, nil
}

func (f *fakeStore) SetCliAuthSessionIdentity(_ context.Context, arg store.SetCliAuthSessionIdentityParams) error {
	row, ok := f.sessions[arg.ID]
	if !ok || row.Status != "PENDING_IDP" {
		return pgErrNoRows
	}
	row.Status = "PENDING_CONSENT"
	row.IdentityEmail = arg.IdentityEmail
	row.IdentityName = arg.IdentityName
	row.IdentityGroups = arg.IdentityGroups
	row.DexCodeVerifier = nil
	f.sessions[row.ID] = row
	return nil
}

func (f *fakeStore) FailCliAuthSession(_ context.Context, id string) error {
	row := f.sessions[id]
	row.Status = "FAILED"
	row.DexCodeVerifier = nil
	f.sessions[id] = row
	return nil
}

func (f *fakeStore) IssueCliAuthCode(_ context.Context, arg store.IssueCliAuthCodeParams) error {
	row, ok := f.sessions[arg.ID]
	if !ok || row.Status != "PENDING_CONSENT" {
		return pgErrNoRows
	}
	row.Status = "CODE_ISSUED"
	row.AuthCodeHash = arg.AuthCodeHash
	row.ChosenDays = arg.ChosenDays
	row.ChosenName = arg.ChosenName
	row.CodeExpiresAt = arg.CodeExpiresAt
	f.sessions[row.ID] = row
	return nil
}

func (f *fakeStore) DenyCliAuthSession(_ context.Context, id string) error {
	row := f.sessions[id]
	row.Status = "DENIED"
	f.sessions[id] = row
	return nil
}

func (f *fakeStore) GetCliAuthSessionByCodeHash(_ context.Context, hash *string) (store.CliAuthSession, error) {
	for _, row := range f.sessions {
		if row.AuthCodeHash != nil && hash != nil && *row.AuthCodeHash == *hash {
			return row, nil
		}
	}
	return store.CliAuthSession{}, pgErrNoRows
}

func (f *fakeStore) ClaimCliAuthCode(_ context.Context, id string) (int64, error) {
	row, ok := f.sessions[id]
	if !ok || row.Status != "CODE_ISSUED" || !row.CodeExpiresAt.Valid || !row.CodeExpiresAt.Time.After(time.Now()) {
		return 0, nil
	}
	row.Status = "EXCHANGED"
	f.sessions[id] = row
	return 1, nil
}

func (f *fakeStore) DeleteExpiredCliAuthSessions(context.Context) (int64, error) { return 0, nil }

var pgErrNoRows = pgx.ErrNoRows

type fakeDex struct {
	identity  *auth.Identity
	err       error
	issuer    string
	challenge string
}

func (f *fakeDex) Issuer() string {
	if f.issuer == "" {
		return "https://dex.example/dex"
	}
	return f.issuer
}

func (f *fakeDex) AuthorizeURL(state, challenge string) string {
	f.challenge = challenge
	return "https://dex.example/dex/auth?state=" + state + "&code_challenge=" + challenge
}

func (f *fakeDex) Exchange(context.Context, string, string) (*auth.Identity, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.identity, nil
}

type fakeIssuer struct {
	token   string
	days    int
	desc    string
	source  string
	creator string
	user    auth.Identity
	err     error
}

func (f *fakeIssuer) Issue(_ context.Context, user auth.Identity, days int, description, source, createdBy string) (string, error) {
	f.user, f.days, f.desc, f.source, f.creator = user, days, description, source, createdBy
	return f.token, f.err
}

func newTestService(t *testing.T) (*Service, *fakeStore, *fakeDex, *fakeIssuer) {
	t.Helper()
	store := newFakeStore()
	dex := &fakeDex{identity: &auth.Identity{
		Issuer: "https://dex.example/dex",
		Email:  "user@example.com",
		Name:   "Example User",
		Groups: []string{"developers"},
	}}
	pats := &fakeIssuer{token: "signed-token"}
	return New(store, dex, pats, 30, 365), store, dex, pats
}

const goodVerifier = "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"

func start(t *testing.T, service *Service) (string, string) {
	t.Helper()
	sessionID, location, err := service.StartAuthorization(context.Background(), AuthorizeRequest{
		ClientID:            ClientID,
		RedirectURI:         "http://localhost:10000/login",
		ResponseType:        "code",
		CodeChallenge:       ChallengeS256(goodVerifier),
		CodeChallengeMethod: "S256",
		State:               "cli-state",
	})
	if err != nil {
		t.Fatalf("StartAuthorization: %v", err)
	}
	if !strings.Contains(location, "&code_challenge=") {
		t.Errorf("location = %q, want a Dex code_challenge", location)
	}
	return sessionID, location
}

func TestStartAuthorizationRejectsBadRequests(t *testing.T) {
	service, _, _, _ := newTestService(t)
	base := AuthorizeRequest{
		ClientID:            ClientID,
		RedirectURI:         "http://localhost:10000/login",
		ResponseType:        "code",
		CodeChallenge:       ChallengeS256(goodVerifier),
		CodeChallengeMethod: "S256",
		State:               "state",
	}
	cases := map[string]func(*AuthorizeRequest){
		"client":    func(r *AuthorizeRequest) { r.ClientID = "other" },
		"type":      func(r *AuthorizeRequest) { r.ResponseType = "token" },
		"method":    func(r *AuthorizeRequest) { r.CodeChallengeMethod = "plain" },
		"challenge": func(r *AuthorizeRequest) { r.CodeChallenge = "" },
		"state":     func(r *AuthorizeRequest) { r.State = "" },
		"redirect":  func(r *AuthorizeRequest) { r.RedirectURI = "https://evil.example:10000/login" },
	}
	for name, mutate := range cases {
		request := base
		mutate(&request)
		if _, _, err := service.StartAuthorization(context.Background(), request); !errors.Is(err, ErrBadRequest) {
			t.Errorf("%s: error = %v, want ErrBadRequest", name, err)
		}
	}
}

func TestFullLoginFlow(t *testing.T) {
	service, store, _, pats := newTestService(t)
	ctx := context.Background()
	sessionID, _ := start(t, service)

	consentID, err := service.Callback(ctx, "dex-code", sessionID, "")
	if err != nil {
		t.Fatalf("Callback: %v", err)
	}
	if sessionID != consentID {
		t.Errorf("Callback id = %q, want %q", consentID, sessionID)
	}
	row, err := service.ConsentSession(ctx, sessionID)
	if err != nil {
		t.Fatalf("ConsentSession: %v", err)
	}
	if row.IdentityEmail == nil || *row.IdentityEmail != "user@example.com" || row.DexCodeVerifier != nil {
		t.Errorf("row = %+v, want the Dex identity with the verifier cleared", row)
	}

	location, err := service.Authorize(ctx, sessionID, 30, "  my cli  ")
	if err != nil {
		t.Fatalf("Authorize: %v", err)
	}
	code := codeFromLocation(t, location)
	row = store.sessions[sessionID]
	if row.Status != "CODE_ISSUED" || row.ChosenName == nil || *row.ChosenName != "my cli" || row.ChosenDays == nil || *row.ChosenDays != 30 {
		t.Errorf("row = %+v, want the consent choice stored", row)
	}

	response, err := service.Exchange(ctx, "authorization_code", code, goodVerifier, "http://localhost:10000/login", ClientID)
	if err != nil {
		t.Fatalf("Exchange: %v", err)
	}
	if response.AccessToken != "signed-token" || response.ExpiresIn != 30*86400 {
		t.Errorf("response = %+v", response)
	}
	if pats.user.Email != "user@example.com" || pats.days != 30 || pats.desc != "my cli" ||
		pats.source != pat.SourceCLILogin || pats.creator != "user@example.com" {
		t.Errorf("issued = %+v, want the Dex identity and consent choice", pats)
	}

	if _, err := service.Exchange(ctx, "authorization_code", code, goodVerifier, "http://localhost:10000/login", ClientID); !errors.Is(err, ErrBadRequest) {
		t.Errorf("second Exchange = %v, want the code burned", err)
	}
}

func TestExchangeRejectsWrongCredentials(t *testing.T) {
	service, _, _, _ := newTestService(t)
	ctx := context.Background()
	sessionID, _ := start(t, service)
	if _, err := service.Callback(ctx, "dex-code", sessionID, ""); err != nil {
		t.Fatalf("Callback: %v", err)
	}
	location, err := service.Authorize(ctx, sessionID, 30, "")
	if err != nil {
		t.Fatalf("Authorize: %v", err)
	}
	code := codeFromLocation(t, location)

	cases := map[string]struct {
		grant, code, verifier, redirect, client string
	}{
		"grant":    {"password", code, goodVerifier, "http://localhost:10000/login", ClientID},
		"code":     {"authorization_code", "wrong", goodVerifier, "http://localhost:10000/login", ClientID},
		"verifier": {"authorization_code", code, "wrong", "http://localhost:10000/login", ClientID},
		"redirect": {"authorization_code", code, goodVerifier, "http://localhost:10001/login", ClientID},
		"client":   {"authorization_code", code, goodVerifier, "http://localhost:10000/login", "other"},
	}
	for name, tc := range cases {
		if _, err := service.Exchange(ctx, tc.grant, tc.code, tc.verifier, tc.redirect, tc.client); !errors.Is(err, ErrBadRequest) {
			t.Errorf("%s: error = %v, want ErrBadRequest", name, err)
		}
	}
	// None of the attempts may have burned the code.
	if _, err := service.Exchange(ctx, "authorization_code", code, goodVerifier, "http://localhost:10000/login", ClientID); err != nil {
		t.Errorf("Exchange after rejected attempts = %v, want success", err)
	}
}

func TestExchangeRejectsExpiredCode(t *testing.T) {
	service, store, _, _ := newTestService(t)
	ctx := context.Background()
	sessionID, _ := start(t, service)
	if _, err := service.Callback(ctx, "dex-code", sessionID, ""); err != nil {
		t.Fatalf("Callback: %v", err)
	}
	location, _ := service.Authorize(ctx, sessionID, 30, "")
	code := codeFromLocation(t, location)

	row := store.sessions[sessionID]
	row.CodeExpiresAt = timestamp(time.Now().Add(-time.Minute))
	store.sessions[sessionID] = row
	if _, err := service.Exchange(ctx, "authorization_code", code, goodVerifier, "http://localhost:10000/login", ClientID); !errors.Is(err, ErrBadRequest) {
		t.Errorf("Exchange(expired) = %v, want ErrBadRequest", err)
	}
}

func TestCallbackFailures(t *testing.T) {
	service, store, dex, _ := newTestService(t)
	ctx := context.Background()

	sessionID, _ := start(t, service)
	if _, err := service.Callback(ctx, "code", sessionID, "https://other.example/dex"); !errors.Is(err, ErrBadRequest) {
		t.Errorf("issuer mismatch = %v, want ErrBadRequest", err)
	}
	if store.sessions[sessionID].Status != "FAILED" {
		t.Error("issuer mismatch did not fail the session")
	}

	sessionID, _ = start(t, service)
	dex.err = errors.New("dex down")
	if _, err := service.Callback(ctx, "code", sessionID, ""); err == nil {
		t.Error("Callback with a broken Dex succeeded")
	}
	if store.sessions[sessionID].Status != "FAILED" {
		t.Error("Dex failure did not fail the session")
	}
	dex.err = nil

	if _, err := service.Callback(ctx, "code", "not-a-uuid", ""); !errors.Is(err, ErrBadRequest) {
		t.Errorf("bad state = %v, want ErrBadRequest", err)
	}
	row := store.sessions[sessionID]
	row.ExpiresAt = timestamp(time.Now().Add(-time.Minute))
	store.sessions[sessionID] = row
	if _, err := service.Callback(ctx, "code", sessionID, ""); !errors.Is(err, ErrBadRequest) {
		t.Errorf("expired session = %v, want ErrBadRequest", err)
	}
}

func TestDenyAndCallbackError(t *testing.T) {
	service, store, _, _ := newTestService(t)
	ctx := context.Background()

	sessionID, _ := start(t, service)
	location, err := service.CallbackError(ctx, sessionID, "access_denied")
	if err != nil {
		t.Fatalf("CallbackError: %v", err)
	}
	if !strings.Contains(location, "error=access_denied") || !strings.Contains(location, "state=cli-state") {
		t.Errorf("location = %q, want the CLI error redirect", location)
	}

	sessionID, _ = start(t, service)
	if _, err := service.Callback(ctx, "code", sessionID, ""); err != nil {
		t.Fatalf("Callback: %v", err)
	}
	location, err = service.Deny(ctx, sessionID)
	if err != nil {
		t.Fatalf("Deny: %v", err)
	}
	if !strings.Contains(location, "error=access_denied") || store.sessions[sessionID].Status != "DENIED" {
		t.Errorf("location/status = %q/%s", location, store.sessions[sessionID].Status)
	}
}

func TestConsentDaysAreClamped(t *testing.T) {
	service, _, _, _ := newTestService(t)
	ctx := context.Background()
	sessionID, _ := start(t, service)
	if _, err := service.Callback(ctx, "code", sessionID, ""); err != nil {
		t.Fatalf("Callback: %v", err)
	}
	if _, err := service.Authorize(ctx, sessionID, 0, ""); !errors.Is(err, ErrBadRequest) {
		t.Errorf("Authorize(0 days) = %v, want ErrBadRequest", err)
	}
	if _, err := service.Authorize(ctx, sessionID, 366, ""); !errors.Is(err, ErrBadRequest) {
		t.Errorf("Authorize(366 days) = %v, want ErrBadRequest", err)
	}

	clamped := New(newFakeStore(), &fakeDex{}, &fakeIssuer{}, 30, 999)
	if clamped.MaxDays() != 365 || clamped.DefaultDays() != 30 {
		t.Errorf("MaxDays/DefaultDays = %d/%d, want 365/30", clamped.MaxDays(), clamped.DefaultDays())
	}
}

func codeFromLocation(t *testing.T, location string) string {
	t.Helper()
	index := strings.Index(location, "code=")
	if index < 0 {
		t.Fatalf("location %q has no code", location)
	}
	code := location[index+len("code="):]
	if end := strings.IndexByte(code, '&'); end >= 0 {
		code = code[:end]
	}
	return code
}
