package clilogin

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func newTestHandler(t *testing.T) (http.Handler, *CookieSigner, *fakeDex, *fakeIssuer) {
	t.Helper()
	service, _, dex, pats := newTestService(t)
	signer, err := NewCookieSigner("ejZRSFgheUBOZXAyUURUITUzdmdINDNeUGpSWHlDM1g=")
	if err != nil {
		t.Fatalf("NewCookieSigner: %v", err)
	}
	router := chi.NewRouter()
	router.Mount("/oauth", New(Deps{
		Service:   service,
		Cookie:    signer,
		PublicURL: "http://localhost:8080",
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	}))
	return router, signer, dex, pats
}

func doRequest(handler http.Handler, method, target string, form url.Values, cookie *http.Cookie) *httptest.ResponseRecorder {
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req := httptest.NewRequest(method, target, body)
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func authorize(handler http.Handler, signer *CookieSigner) (string, *http.Cookie) {
	query := url.Values{
		"client_id":             {ClientID},
		"redirect_uri":          {"http://localhost:10000/login"},
		"response_type":         {"code"},
		"code_challenge":        {ChallengeS256(goodVerifier)},
		"code_challenge_method": {"S256"},
		"state":                 {"cli-state"},
	}
	rec := doRequest(handler, http.MethodGet, "/oauth/authorize?"+query.Encode(), nil, nil)
	if rec.Code != http.StatusFound {
		return "", nil
	}
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == CookieName {
			sessionID, ok := signer.SessionID(cookie.Value)
			if !ok {
				return "", nil
			}
			return sessionID, cookie
		}
	}
	return "", nil
}

func TestBrokerFullFlow(t *testing.T) {
	handler, signer, _, pats := newTestHandler(t)

	sessionID, cookie := authorize(handler, signer)
	if sessionID == "" || cookie == nil {
		t.Fatal("authorize did not set the browser cookie")
	}

	callback := doRequest(handler, http.MethodGet, "/oauth/callback?code=dex-code&state="+sessionID, nil, cookie)
	if callback.Code != http.StatusFound || callback.Header().Get("Location") != "http://localhost:8080/oauth/consent" {
		t.Fatalf("callback = %d %q", callback.Code, callback.Header().Get("Location"))
	}

	page := doRequest(handler, http.MethodGet, "/oauth/consent", nil, cookie)
	if page.Code != http.StatusOK || !strings.Contains(page.Body.String(), "user@example.com") {
		t.Fatalf("consent page = %d %q", page.Code, page.Body.String())
	}

	form := url.Values{"decision": {"approve"}, "days": {"7"}, "name": {"cli"}}
	req := httptest.NewRequest(http.MethodPost, "/oauth/consent", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "http://localhost:8080")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("consent = %d %q", rec.Code, rec.Body.String())
	}
	location := rec.Header().Get("Location")
	code := codeFromLocation(t, location)
	if !strings.Contains(location, "state=cli-state") {
		t.Errorf("location = %q, want the CLI state", location)
	}

	tokenForm := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"code_verifier": {goodVerifier},
		"redirect_uri":  {"http://localhost:10000/login"},
		"client_id":     {ClientID},
	}
	token := doRequest(handler, http.MethodPost, "/oauth/token", tokenForm, nil)
	if token.Code != http.StatusOK {
		t.Fatalf("token = %d %q", token.Code, token.Body.String())
	}
	var response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.Unmarshal(token.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if response.AccessToken != "signed-token" || response.TokenType != "Bearer" || response.ExpiresIn != 7*86400 {
		t.Errorf("response = %+v", response)
	}
	if pats.source != "CLI_LOGIN" || pats.days != 7 || pats.desc != "cli" {
		t.Errorf("issued = %+v", pats)
	}

	again := doRequest(handler, http.MethodPost, "/oauth/token", tokenForm, nil)
	if again.Code != http.StatusBadRequest || !strings.Contains(again.Body.String(), "invalid_grant") {
		t.Errorf("second token = %d %q", again.Code, again.Body.String())
	}
}

func TestBrokerRejectsBadRequests(t *testing.T) {
	handler, signer, _, _ := newTestHandler(t)

	bad := doRequest(handler, http.MethodGet, "/oauth/authorize?client_id=other", nil, nil)
	if bad.Code != http.StatusBadRequest {
		t.Errorf("authorize(unknown client) = %d, want 400", bad.Code)
	}

	sessionID, cookie := authorize(handler, signer)
	if sessionID == "" {
		t.Fatal("authorize failed")
	}

	otherCookie := &http.Cookie{Name: CookieName, Value: signer.Value("another-session")}
	if rec := doRequest(handler, http.MethodGet, "/oauth/callback?code=dex-code&state="+sessionID, nil, otherCookie); rec.Code != http.StatusForbidden {
		t.Errorf("callback(foreign cookie) = %d, want 403", rec.Code)
	}
	if rec := doRequest(handler, http.MethodGet, "/oauth/consent", nil, nil); rec.Code != http.StatusForbidden {
		t.Errorf("consent(no cookie) = %d, want 403", rec.Code)
	}

	form := url.Values{"decision": {"approve"}, "days": {"7"}}
	rec := doRequest(handler, http.MethodPost, "/oauth/consent", form, cookie)
	if rec.Code != http.StatusForbidden {
		t.Errorf("consent(no origin) = %d, want 403", rec.Code)
	}

	token := doRequest(handler, http.MethodPost, "/oauth/token", url.Values{"grant_type": {"password"}}, nil)
	if token.Code != http.StatusBadRequest || !strings.Contains(token.Body.String(), "invalid_grant") {
		t.Errorf("token(bad grant) = %d %q", token.Code, token.Body.String())
	}
}

func TestBrokerCallbackErrorRedirectsToTheCLI(t *testing.T) {
	handler, signer, _, _ := newTestHandler(t)
	sessionID, cookie := authorize(handler, signer)
	if sessionID == "" {
		t.Fatal("authorize failed")
	}
	rec := doRequest(handler, http.MethodGet, "/oauth/callback?error=access_denied&state="+sessionID, nil, cookie)
	if rec.Code != http.StatusFound {
		t.Fatalf("callback error = %d %q", rec.Code, rec.Body.String())
	}
	location := rec.Header().Get("Location")
	if !strings.HasPrefix(location, "http://localhost:10000/login?error=access_denied") || !strings.Contains(location, "state=cli-state") {
		t.Errorf("location = %q", location)
	}
}

func TestConsentRejectsInvalidDays(t *testing.T) {
	handler, signer, _, _ := newTestHandler(t)
	sessionID, cookie := authorize(handler, signer)
	if sessionID == "" {
		t.Fatal("authorize failed")
	}
	if rec := doRequest(handler, http.MethodGet, "/oauth/callback?code=dex-code&state="+sessionID, nil, cookie); rec.Code != http.StatusFound {
		t.Fatalf("callback = %d", rec.Code)
	}

	form := url.Values{"decision": {"approve"}, "days": {"0"}}
	req := httptest.NewRequest(http.MethodPost, "/oauth/consent", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "http://localhost:8080")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "days must be between") {
		t.Errorf("consent(0 days) = %d %q", rec.Code, rec.Body.String())
	}
}
