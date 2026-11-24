package secretstore

import (
	"encoding/base64"
	"strings"
	"testing"
)

func testKey(seed byte) string {
	key := make([]byte, 32)
	for i := range key {
		key[i] = seed
	}
	return base64.StdEncoding.EncodeToString(key)
}

func TestCipherRoundTrip(t *testing.T) {
	c, err := NewCipher("k1:"+testKey(1), "")
	if err != nil {
		t.Fatalf("NewCipher() = %v", err)
	}
	sealed, err := c.Encrypt("variable.variable_value", "hello")
	if err != nil {
		t.Fatalf("Encrypt() = %v", err)
	}
	if !strings.HasPrefix(sealed, "enc:v1:k1:") {
		t.Fatalf("envelope = %q, want enc:v1:k1: prefix", sealed)
	}
	got, err := c.Decrypt("variable.variable_value", sealed)
	if err != nil {
		t.Fatalf("Decrypt() = %v", err)
	}
	if got != "hello" {
		t.Fatalf("Decrypt() = %q, want hello", got)
	}
}

func TestCipherRejectsWrongColumnAndKey(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	sealed, _ := c.Encrypt("variable.variable_value", "hello")
	if _, err := c.Decrypt("globalvar.variable_value", sealed); err == nil {
		t.Fatal("Decrypt() with another column = nil error, want error")
	} else if strings.Contains(err.Error(), "hello") {
		t.Fatalf("Decrypt() error leaks the value: %v", err)
	}
	other, _ := NewCipher("k2:"+testKey(2), "")
	if _, err := other.Decrypt("variable.variable_value", sealed); err == nil {
		t.Fatal("Decrypt() with another key = nil error, want error")
	}
}

func TestCipherEmptyAndPlaintext(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	if got, err := c.Encrypt("ssh.private_key", ""); err != nil || got != "" {
		t.Fatalf("Encrypt(empty) = %q, %v; want empty, nil", got, err)
	}
	if got, err := c.Decrypt("ssh.private_key", "plain legacy"); err != nil || got != "plain legacy" {
		t.Fatalf("Decrypt(plaintext) = %q, %v; want passthrough", got, err)
	}
	if _, err := c.Decrypt("ssh.private_key", "enc:v1:k1:notbase64:notbase64"); err == nil {
		t.Fatal("Decrypt(malformed) = nil error, want error")
	}
	if _, err := c.Decrypt("ssh.private_key", "enc:v1:unknown:AAAA:AAAA"); err == nil {
		t.Fatal("Decrypt(unknown key) = nil error, want error")
	}
}

func TestCipherKeyring(t *testing.T) {
	c, err := NewCipher("old:"+testKey(1)+",new:"+testKey(2), "new")
	if err != nil {
		t.Fatalf("NewCipher() = %v", err)
	}
	if c.ActiveKeyID() != "new" {
		t.Fatalf("ActiveKeyID() = %q, want new", c.ActiveKeyID())
	}
	old, _ := NewCipher("old:"+testKey(1), "")
	sealed, _ := old.Encrypt("vcs.client_secret", "s3cret")
	if got, err := c.Decrypt("vcs.client_secret", sealed); err != nil || got != "s3cret" {
		t.Fatalf("Decrypt(old key) = %q, %v", got, err)
	}
}

func TestCipherRejectsBadConfiguration(t *testing.T) {
	if c, err := NewCipher("", ""); c != nil || err != nil {
		t.Fatalf("NewCipher(empty) = %v, %v; want nil, nil", c, err)
	}
	cases := []string{
		"nokey",
		"k1:short",
		"k1:" + testKey(1) + ",k1:" + testKey(2),
		"k1:" + testKey(1) + ",bad id:" + testKey(2),
	}
	for _, spec := range cases {
		if _, err := NewCipher(spec, ""); err == nil {
			t.Errorf("NewCipher(%q) = nil error, want error", spec)
		}
	}
	if _, err := NewCipher("k1:"+testKey(1), "missing"); err == nil {
		t.Fatal("NewCipher(active missing) = nil error, want error")
	}
}
