package auth

import (
	"context"
	"errors"
	"testing"
)

type fakeVerifier struct {
	identity *Identity
	err      error
}

func (f fakeVerifier) Verify(context.Context, string) (*Identity, error) {
	return f.identity, f.err
}

func TestAuthenticatePrefersTheDevelopmentToken(t *testing.T) {
	identity, ok := Authenticate(context.Background(), "Bearer dev", "dev", fakeVerifier{err: errors.New("no")})
	if !ok || identity == nil || identity.Email != "dev@localhost" {
		t.Fatalf("Authenticate(dev) = %+v/%v, want the dev identity", identity, ok)
	}
}

func TestAuthenticateTriesEachVerifier(t *testing.T) {
	want := &Identity{Issuer: "Terrakube", Email: "user@example.com"}
	identity, ok := Authenticate(context.Background(), "Bearer token", "dev",
		fakeVerifier{err: errors.New("not mine")}, fakeVerifier{identity: want})
	if !ok || identity != want {
		t.Fatalf("Authenticate = %+v/%v, want %+v", identity, ok, want)
	}
}

func TestAuthenticateRejectsMissingOrUnknownCredentials(t *testing.T) {
	for _, header := range []string{"", "Bearer ", "Basic abc", "Bearer wrong"} {
		if identity, ok := Authenticate(context.Background(), header, "dev", fakeVerifier{err: errors.New("no")}); ok || identity != nil {
			t.Errorf("Authenticate(%q) = %+v/%v, want rejection", header, identity, ok)
		}
	}
}

func TestIdentityRoundTripsThroughContext(t *testing.T) {
	want := &Identity{Issuer: "Terrakube", Email: "user@example.com"}
	got, ok := IdentityFromContext(WithIdentity(context.Background(), want))
	if !ok || got != want {
		t.Fatalf("IdentityFromContext = %+v/%v, want %+v", got, ok, want)
	}
	if _, ok := IdentityFromContext(context.Background()); ok {
		t.Error("IdentityFromContext(empty) reported an identity")
	}
}
