package notification

import "testing"

func TestSignIsStableHexHMAC(t *testing.T) {
	const want = "97f6f2aac0738ccf3770bd6e3e7f89851ed0695dddfa0be66b33e448d7ee05ed"
	if got := Sign("my-secret", `{"a":1}`); got != want {
		t.Fatalf("Sign() = %q, want %q", got, want)
	}
}

func TestSignDifferentSecretsDiffer(t *testing.T) {
	if Sign("secret-one", "payload") == Sign("secret-two", "payload") {
		t.Fatal("Sign() produced the same signature for different secrets")
	}
}
