package clilogin

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type dexFixture struct {
	server   *httptest.Server
	issuer   string
	key      *rsa.PrivateKey
	received url.Values
	idToken  string
	status   int
}

func newDexFixture(t *testing.T) *dexFixture {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	fixture := &dexFixture{key: key, status: http.StatusOK}
	mux := http.NewServeMux()
	mux.HandleFunc("/dex/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"issuer":   fixture.issuer,
			"jwks_uri": fixture.issuer + "/keys",
		})
	})
	mux.HandleFunc("/dex/keys", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"keys": []map[string]string{{
				"kty": "RSA",
				"kid": "test-key",
				"use": "sig",
				"alg": "RS256",
				"n":   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
				"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
			}},
		})
	})
	mux.HandleFunc("/dex/token", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		fixture.received = r.PostForm
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(fixture.status)
		if fixture.status != http.StatusOK {
			_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"id_token": fixture.idToken})
	})
	fixture.server = httptest.NewServer(mux)
	fixture.issuer = fixture.server.URL + "/dex"
	t.Cleanup(fixture.server.Close)
	return fixture
}

func (f *dexFixture) sign(t *testing.T, audience string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":    f.issuer,
		"aud":    audience,
		"sub":    "user-123",
		"exp":    time.Now().Add(time.Hour).Unix(),
		"iat":    time.Now().Unix(),
		"email":  "user@example.com",
		"name":   "Example User",
		"groups": []string{"developers"},
	})
	token.Header["kid"] = "test-key"
	signed, err := token.SignedString(f.key)
	if err != nil {
		t.Fatalf("sign id_token: %v", err)
	}
	return signed
}

func TestDexExchangeSendsTheAuthorizationCodeGrant(t *testing.T) {
	fixture := newDexFixture(t)
	fixture.idToken = fixture.sign(t, "terrakube")
	dex := NewDex(fixture.issuer, "terrakube", "http://localhost:8080/oauth/callback")

	identity, err := dex.Exchange(context.Background(), "dex-code", "the-verifier")
	if err != nil {
		t.Fatalf("Exchange: %v", err)
	}
	if identity.Email != "user@example.com" || identity.Name != "Example User" || len(identity.Groups) != 1 {
		t.Errorf("identity = %+v", identity)
	}
	want := map[string]string{
		"grant_type":    "authorization_code",
		"code":          "dex-code",
		"redirect_uri":  "http://localhost:8080/oauth/callback",
		"client_id":     "terrakube",
		"code_verifier": "the-verifier",
	}
	for key, value := range want {
		if got := fixture.received.Get(key); got != value {
			t.Errorf("form[%s] = %q, want %q", key, got, value)
		}
	}
}

func TestDexExchangeRejectsBadResponses(t *testing.T) {
	fixture := newDexFixture(t)
	dex := NewDex(fixture.issuer, "terrakube", "http://localhost:8080/oauth/callback")

	fixture.idToken = fixture.sign(t, "another-client")
	if _, err := dex.Exchange(context.Background(), "code", "verifier"); err == nil {
		t.Error("Exchange accepted an id_token for another audience")
	}

	fixture.idToken = ""
	if _, err := dex.Exchange(context.Background(), "code", "verifier"); err == nil {
		t.Error("Exchange accepted a response without id_token")
	}

	fixture.status = http.StatusBadRequest
	if _, err := dex.Exchange(context.Background(), "code", "verifier"); err == nil {
		t.Error("Exchange accepted a failed token response")
	}
}

func TestAuthorizeURLCarriesPKCE(t *testing.T) {
	dex := NewDex("https://dex.example/dex/", "terrakube", "https://terrakube.example/oauth/callback")
	location := dex.AuthorizeURL("state-1", "challenge-1")
	parsed, err := url.Parse(location)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.Path != "/dex/auth" {
		t.Errorf("path = %q, want /dex/auth", parsed.Path)
	}
	query := parsed.Query()
	for key, value := range map[string]string{
		"response_type":         "code",
		"client_id":             "terrakube",
		"redirect_uri":          "https://terrakube.example/oauth/callback",
		"code_challenge":        "challenge-1",
		"code_challenge_method": "S256",
		"state":                 "state-1",
	} {
		if got := query.Get(key); got != value {
			t.Errorf("%s = %q, want %q", key, got, value)
		}
	}
	if query.Get("scope") == "" {
		t.Error("authorize URL has no scope")
	}
}
