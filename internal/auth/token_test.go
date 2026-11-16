package auth

import (
	"bytes"
	"encoding/base64"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func testSecret(n int) string {
	return base64.RawURLEncoding.EncodeToString(bytes.Repeat([]byte{0x42}, n))
}

func TestValidatorAlgorithmByKeyLength(t *testing.T) {
	cases := []struct {
		bytes int
		alg   string
	}{
		{32, "HS256"},
		{48, "HS384"},
		{64, "HS512"},
	}
	for _, tc := range cases {
		v, err := NewValidator(testSecret(tc.bytes))
		if err != nil {
			t.Fatalf("NewValidator(%d bytes) error = %v", tc.bytes, err)
		}
		if got := v.alg.Alg(); got != tc.alg {
			t.Errorf("NewValidator(%d bytes) alg = %q, want %q", tc.bytes, got, tc.alg)
		}
	}
}

func TestSignAndParse(t *testing.T) {
	secret := testSecret(64)
	raw, err := Sign(secret, Claims{
		OrganizationID: "org-1",
		WorkspaceID:    "ws-1",
		JobID:          "42",
		StepID:         "step-1",
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	v, err := NewValidator(secret)
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}
	claims, err := v.Parse(raw)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.OrganizationID != "org-1" || claims.WorkspaceID != "ws-1" {
		t.Fatalf("claims = %#v", claims)
	}
	if got := Classify(claims); got != TypeWorkspaceExecutor {
		t.Errorf("Classify() = %v, want workspace executor", got)
	}
	if got := claims.EffectiveStepID(); got != "step-1" {
		t.Errorf("EffectiveStepID() = %q, want step-1", got)
	}
}

func TestEffectiveStepIDFromJTI(t *testing.T) {
	claims := &Claims{ID: "11111111-1111-1111-1111-111111111111-22222222-2222-2222-2222-222222222222"}
	if got := claims.EffectiveStepID(); got != "22222222-2222-2222-2222-222222222222" {
		t.Errorf("EffectiveStepID() = %q", got)
	}
}

func TestParseRejectsExpiredToken(t *testing.T) {
	secret := testSecret(32)
	past := time.Now().Add(-time.Hour)
	raw, err := Sign(secret, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(past),
			ExpiresAt: jwt.NewNumericDate(past.Add(time.Minute)),
		},
		WorkspaceID: "ws-1",
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	v, err := NewValidator(secret)
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}
	if _, err := v.Parse(raw); err == nil {
		t.Fatal("Parse() error = nil, want expired token error")
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	raw, err := Sign(testSecret(32), Claims{WorkspaceID: "ws-1"})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	v, err := NewValidator(testSecret(64))
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}
	if _, err := v.Parse(raw); err == nil {
		t.Fatal("Parse() error = nil, want signature error")
	}
}

func TestClassify(t *testing.T) {
	cases := []struct {
		name   string
		claims Claims
		want   Type
	}{
		{"not internal", Claims{}, TypeNotInternal},
		{"workspace", Claims{RegisteredClaims: jwt.RegisteredClaims{Issuer: Issuer}, WorkspaceID: "ws"}, TypeWorkspaceExecutor},
		{"executor service", Claims{RegisteredClaims: jwt.RegisteredClaims{Issuer: Issuer, Subject: SubjectExecutor}}, TypeExecutorService},
		{"registry service", Claims{RegisteredClaims: jwt.RegisteredClaims{Issuer: Issuer, Subject: SubjectRegistry}}, TypeRegistryService},
		{"unknown", Claims{RegisteredClaims: jwt.RegisteredClaims{Issuer: Issuer, Subject: "something"}}, TypeUnknownInternal},
	}
	for _, tc := range cases {
		if got := Classify(&tc.claims); got != tc.want {
			t.Errorf("%s: Classify() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestExtractToken(t *testing.T) {
	basic := "Basic " + base64.StdEncoding.EncodeToString([]byte("terraform:abc.def.ghi"))
	if got := ExtractToken(basic); got != "abc.def.ghi" {
		t.Errorf("ExtractToken(Basic) = %q", got)
	}
	if got := ExtractToken("bearer token-value"); got != "token-value" {
		t.Errorf("ExtractToken(Bearer) = %q", got)
	}
	for _, header := range []string{"", "Basic !!!", "Basic " + base64.StdEncoding.EncodeToString([]byte("no-colon")), "Token abc"} {
		if got := ExtractToken(header); got != "" {
			t.Errorf("ExtractToken(%q) = %q, want empty", header, got)
		}
	}
}
