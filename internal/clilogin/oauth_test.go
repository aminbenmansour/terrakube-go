package clilogin

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestPKCEMatchesTheRFCVector(t *testing.T) {
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	if got := ChallengeS256(verifier); got != "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM" {
		t.Fatalf("ChallengeS256 = %q, want the RFC 7636 vector", got)
	}
	if !VerifyS256(verifier, ChallengeS256(verifier)) {
		t.Error("VerifyS256 rejected the matching verifier")
	}
	if VerifyS256(verifier, ChallengeS256("other")) {
		t.Error("VerifyS256 accepted a mismatched verifier")
	}
	if VerifyS256("", "") {
		t.Error("VerifyS256 accepted empty values")
	}
}

func TestNewCodeVerifier(t *testing.T) {
	first, err := NewCodeVerifier()
	if err != nil {
		t.Fatalf("NewCodeVerifier: %v", err)
	}
	second, _ := NewCodeVerifier()
	if len(first) != 64 || first == second {
		t.Fatalf("NewCodeVerifier = %q/%q, want 64 distinct characters", first, second)
	}
	if _, err := base64.RawURLEncoding.DecodeString(first); err != nil {
		t.Errorf("verifier is not base64url: %v", err)
	}
}

func TestValidateLoopbackRedirect(t *testing.T) {
	valid := []string{
		"http://localhost:10000/login",
		"http://127.0.0.1:10001/login",
		"http://[::1]:10010/login",
	}
	for _, raw := range valid {
		if err := ValidateLoopbackRedirect(raw); err != nil {
			t.Errorf("ValidateLoopbackRedirect(%q) = %v, want nil", raw, err)
		}
	}

	invalid := []string{
		"https://localhost:10000/login",
		"http://evil.example:10000/login",
		"http://localhost:9999/login",
		"http://localhost:10011/login",
		"http://localhost/login",
		"http://localhost:10000/callback",
		"http://localhost:10000/login?x=1",
		"http://localhost:10000/login#fragment",
		"http://user:pass@localhost:10000/login",
	}
	for _, raw := range invalid {
		if err := ValidateLoopbackRedirect(raw); err == nil {
			t.Errorf("ValidateLoopbackRedirect(%q) = nil, want rejection", raw)
		}
	}
}

func TestCookieBindsAndDetectsTampering(t *testing.T) {
	signer, err := NewCookieSigner(base64.RawURLEncoding.EncodeToString([]byte(strings.Repeat("k", 32))))
	if err != nil {
		t.Fatalf("NewCookieSigner: %v", err)
	}
	value := signer.Value("session-1")
	if got, ok := signer.SessionID(value); !ok || got != "session-1" {
		t.Fatalf("SessionID = %q/%v, want session-1", got, ok)
	}
	for _, bad := range []string{"", "session-1", "session-1.", ".signature", value + "x", strings.Replace(value, "session-1", "session-2", 1)} {
		if _, ok := signer.SessionID(bad); ok {
			t.Errorf("SessionID(%q) accepted a bad cookie", bad)
		}
	}

	other, _ := NewCookieSigner(base64.RawURLEncoding.EncodeToString([]byte(strings.Repeat("z", 32))))
	if _, ok := other.SessionID(value); ok {
		t.Error("a cookie signed with another secret was accepted")
	}
}

func TestConsentPageEscapesAndShowsLimits(t *testing.T) {
	page := ConsentPage(`<script>alert(1)</script>`, 30, 365, `bad "input"`)
	if strings.Contains(page, "<script>") {
		t.Error("ConsentPage did not escape the email")
	}
	if !strings.Contains(page, `max="365"`) || !strings.Contains(page, `value="30"`) {
		t.Error("ConsentPage does not carry the day limits")
	}
	if strings.Contains(page, `bad "input"`) {
		t.Error("ConsentPage did not escape the problem message")
	}
}
