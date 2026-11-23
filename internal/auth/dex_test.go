package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwksServer is a minimal OIDC issuer: discovery plus a JWKS whose keys the
// test can rotate.
type jwksServer struct {
	t      *testing.T
	server *httptest.Server

	mu       sync.Mutex
	keys     map[string]*rsa.PrivateKey
	discover int
	jwks     int
}

func newJWKSServer(t *testing.T, keys map[string]*rsa.PrivateKey) *jwksServer {
	t.Helper()
	s := &jwksServer{t: t, keys: keys}
	mux := http.NewServeMux()
	mux.HandleFunc("/dex/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		s.discover++
		s.mu.Unlock()
		writeJSON(t, w, map[string]any{
			"issuer":   s.issuer(),
			"jwks_uri": s.issuer() + "/keys",
		})
	})
	mux.HandleFunc("/dex/keys", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		s.jwks++
		published := make([]map[string]any, 0, len(s.keys))
		for kid, key := range s.keys {
			published = append(published, map[string]any{
				"kty": "RSA",
				"kid": kid,
				"use": "sig",
				"alg": "RS256",
				"n":   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
				"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
			})
		}
		s.mu.Unlock()
		writeJSON(t, w, map[string]any{"keys": published})
	})
	s.server = httptest.NewServer(mux)
	t.Cleanup(s.server.Close)
	return s
}

func (s *jwksServer) issuer() string { return s.server.URL + "/dex" }

func (s *jwksServer) setKeys(keys map[string]*rsa.PrivateKey) {
	s.mu.Lock()
	s.keys = keys
	s.mu.Unlock()
}

func (s *jwksServer) fetchCounts() (discovery, jwks int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.discover, s.jwks
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Errorf("encode response: %v", err)
	}
}

func newRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return key
}

func signDexToken(t *testing.T, key *rsa.PrivateKey, kid string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	raw, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return raw
}

func dexClaims(issuer, audience string, expires time.Time) jwt.MapClaims {
	return jwt.MapClaims{
		"iss":   issuer,
		"sub":   "user-123",
		"aud":   audience,
		"exp":   expires.Unix(),
		"iat":   time.Now().Add(-time.Minute).Unix(),
		"email": "user@example.com",
	}
}

func TestDexValidatorAcceptsAValidToken(t *testing.T) {
	key := newRSAKey(t)
	srv := newJWKSServer(t, map[string]*rsa.PrivateKey{"key-a": key})
	validator := NewDexValidator(srv.issuer(), []string{"terrakube"}, false)

	raw := signDexToken(t, key, "key-a", dexClaims(srv.issuer(), "terrakube", time.Now().Add(time.Hour)))
	if err := validator.Verify(context.Background(), raw); err != nil {
		t.Fatalf("Verify(valid) = %v, want nil", err)
	}

	discovery, jwks := srv.fetchCounts()
	if discovery != 1 || jwks != 1 {
		t.Errorf("fetches = discovery %d, jwks %d; want each cached after one fetch", discovery, jwks)
	}
}

func TestDexValidatorRejectsExpiredAndWrongIssuer(t *testing.T) {
	key := newRSAKey(t)
	srv := newJWKSServer(t, map[string]*rsa.PrivateKey{"key-a": key})
	validator := NewDexValidator(srv.issuer(), []string{"terrakube"}, false)

	expired := signDexToken(t, key, "key-a", dexClaims(srv.issuer(), "terrakube", time.Now().Add(-time.Hour)))
	if err := validator.Verify(context.Background(), expired); err == nil {
		t.Error("Verify(expired) = nil, want an error")
	}

	wrongIssuer := signDexToken(t, key, "key-a", dexClaims("https://evil.example", "terrakube", time.Now().Add(time.Hour)))
	if err := validator.Verify(context.Background(), wrongIssuer); err == nil {
		t.Error("Verify(wrong issuer) = nil, want an error")
	}
}

func TestDexValidatorAudienceValidationIsOptIn(t *testing.T) {
	key := newRSAKey(t)
	srv := newJWKSServer(t, map[string]*rsa.PrivateKey{"key-a": key})
	raw := signDexToken(t, key, "key-a", dexClaims(srv.issuer(), "another-client", time.Now().Add(time.Hour)))

	lenient := NewDexValidator(srv.issuer(), []string{"terrakube"}, false)
	if err := lenient.Verify(context.Background(), raw); err != nil {
		t.Fatalf("Verify(unchecked audience) = %v, want nil", err)
	}

	strict := NewDexValidator(srv.issuer(), []string{"terrakube"}, true)
	if err := strict.Verify(context.Background(), raw); err == nil {
		t.Error("Verify(wrong audience) = nil, want an error")
	}

	extra := NewDexValidator(srv.issuer(), []string{"terrakube", "another-client"}, true)
	if err := extra.Verify(context.Background(), raw); err != nil {
		t.Errorf("Verify(extra audience) = %v, want nil", err)
	}
}

func TestDexValidatorRefreshesKeysForANewKid(t *testing.T) {
	oldKey, newKey := newRSAKey(t), newRSAKey(t)
	srv := newJWKSServer(t, map[string]*rsa.PrivateKey{"old": oldKey})
	validator := NewDexValidator(srv.issuer(), []string{"terrakube"}, false)
	validator.refreshInterval = 0

	oldToken := signDexToken(t, oldKey, "old", dexClaims(srv.issuer(), "terrakube", time.Now().Add(time.Hour)))
	if err := validator.Verify(context.Background(), oldToken); err != nil {
		t.Fatalf("Verify(old key) = %v, want nil", err)
	}

	srv.setKeys(map[string]*rsa.PrivateKey{"new": newKey})
	newToken := signDexToken(t, newKey, "new", dexClaims(srv.issuer(), "terrakube", time.Now().Add(time.Hour)))
	if err := validator.Verify(context.Background(), newToken); err != nil {
		t.Fatalf("Verify(rotated key) = %v, want nil after a JWKS refresh", err)
	}
}

func TestDexValidatorRejectsAnUnknownKid(t *testing.T) {
	key, stranger := newRSAKey(t), newRSAKey(t)
	srv := newJWKSServer(t, map[string]*rsa.PrivateKey{"key-a": key})
	validator := NewDexValidator(srv.issuer(), []string{"terrakube"}, false)
	validator.refreshInterval = 0

	raw := signDexToken(t, stranger, "missing", dexClaims(srv.issuer(), "terrakube", time.Now().Add(time.Hour)))
	if err := validator.Verify(context.Background(), raw); err == nil {
		t.Error("Verify(unknown kid) = nil, want an error")
	}
}

func TestDexValidatorRejectsATokenSignedByAnotherKey(t *testing.T) {
	key, impostor := newRSAKey(t), newRSAKey(t)
	srv := newJWKSServer(t, map[string]*rsa.PrivateKey{"key-a": key})
	validator := NewDexValidator(srv.issuer(), []string{"terrakube"}, false)

	raw := signDexToken(t, impostor, "key-a", dexClaims(srv.issuer(), "terrakube", time.Now().Add(time.Hour)))
	if err := validator.Verify(context.Background(), raw); err == nil {
		t.Error("Verify(bad signature) = nil, want an error")
	}
}

func TestDexValidatorReportsDiscoveryFailures(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusServiceUnavailable)
	}))
	t.Cleanup(srv.Close)

	validator := NewDexValidator(srv.URL+"/dex", []string{"terrakube"}, false)
	if err := validator.Verify(context.Background(), "not-a-jwt"); err == nil {
		t.Error("Verify with unreachable discovery = nil, want an error")
	}
}
