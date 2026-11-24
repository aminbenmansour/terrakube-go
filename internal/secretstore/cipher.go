// Package secretstore encrypts secret columns at rest and wraps the sqlc
// querier so callers keep reading and writing plaintext.
package secretstore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const envelopePrefix = "enc:v1:"

// Cipher is a keyring of AES-256-GCM keys. The active key encrypts new values;
// every key in the ring can decrypt its own ciphertext.
type Cipher struct {
	aeads  map[string]cipher.AEAD
	active string
}

// NewCipher parses a keyring of comma-separated id:base64key entries. An empty
// spec returns nil, nil so callers can run without encryption.
func NewCipher(spec, active string) (*Cipher, error) {
	if strings.TrimSpace(spec) == "" {
		return nil, nil
	}
	c := &Cipher{aeads: make(map[string]cipher.AEAD)}
	first := ""
	for _, entry := range strings.Split(spec, ",") {
		id, encoded, ok := strings.Cut(strings.TrimSpace(entry), ":")
		if !ok {
			return nil, fmt.Errorf("secretstore: key entry %q is not id:base64key", id)
		}
		if !validKeyID(id) {
			return nil, fmt.Errorf("secretstore: key id %q must match [A-Za-z0-9_-]{1,32}", id)
		}
		key, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, fmt.Errorf("secretstore: key %q is not base64: %w", id, err)
		}
		if len(key) != 32 {
			return nil, fmt.Errorf("secretstore: key %q must decode to 32 bytes, got %d", id, len(key))
		}
		if _, dup := c.aeads[id]; dup {
			return nil, fmt.Errorf("secretstore: duplicate key id %q", id)
		}
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, fmt.Errorf("secretstore: key %q: %w", id, err)
		}
		aead, err := cipher.NewGCM(block)
		if err != nil {
			return nil, fmt.Errorf("secretstore: key %q: %w", id, err)
		}
		c.aeads[id] = aead
		if first == "" {
			first = id
		}
	}
	if active == "" {
		active = first
	}
	if _, ok := c.aeads[active]; !ok {
		return nil, fmt.Errorf("secretstore: active key %q is not in the keyring", active)
	}
	c.active = active
	return c, nil
}

// ActiveKeyID is the key new ciphertext is sealed with.
func (c *Cipher) ActiveKeyID() string { return c.active }

// Encrypt seals plaintext for the named column. Empty values stay empty.
func (c *Cipher) Encrypt(column, plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	aead := c.aeads[c.active]
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("secretstore: encrypt %s: %w", column, err)
	}
	sealed := aead.Seal(nil, nonce, []byte(plaintext), []byte(column))
	return envelopePrefix + c.active + ":" +
		base64.RawURLEncoding.EncodeToString(nonce) + ":" +
		base64.RawURLEncoding.EncodeToString(sealed), nil
}

// Decrypt opens a value for the named column. Legacy plaintext passes through
// unchanged; a value that cannot be opened is an error, never ciphertext.
func (c *Cipher) Decrypt(column, value string) (string, error) {
	if !strings.HasPrefix(value, envelopePrefix) {
		return value, nil
	}
	parts := strings.SplitN(strings.TrimPrefix(value, envelopePrefix), ":", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("secretstore: decrypt %s: malformed envelope", column)
	}
	aead, ok := c.aeads[parts[0]]
	if !ok {
		return "", fmt.Errorf("secretstore: decrypt %s: unknown key %q", column, parts[0])
	}
	nonce, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("secretstore: decrypt %s: bad nonce: %w", column, err)
	}
	if len(nonce) != aead.NonceSize() {
		return "", fmt.Errorf("secretstore: decrypt %s: nonce is %d bytes, want %d", column, len(nonce), aead.NonceSize())
	}
	sealed, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return "", fmt.Errorf("secretstore: decrypt %s: bad ciphertext: %w", column, err)
	}
	plain, err := aead.Open(nil, nonce, sealed, []byte(column))
	if err != nil {
		return "", fmt.Errorf("secretstore: decrypt %s: %w", column, errors.New("authentication failed"))
	}
	return string(plain), nil
}

func validKeyID(id string) bool {
	if len(id) == 0 || len(id) > 32 {
		return false
	}
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '-':
		default:
			return false
		}
	}
	return true
}
