// Package clilogin implements the terraform login broker: the OAuth
// authorization-code flow that turns a Dex login into a Terrakube personal
// access token. Upstream keeps it in plugin/token/login; the wire contract is
// what matters here.
package clilogin

import (
	"crypto/hkdf"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// ClientID is the public OAuth client the CLI uses, as upstream names it.
	ClientID = "terraform-cli"
	// PortLow and PortHigh bound the loopback ports Terraform may open.
	PortLow  = 10000
	PortHigh = 10010

	// CookieName binds a login to the browser that started it.
	CookieName = "tk_cli_login"
	// CookiePath scopes the cookie to the broker endpoints.
	CookiePath = "/oauth"

	sessionTTL = 10 * time.Minute
	codeTTL    = 60 * time.Second
)

var errInvalidCredential = errors.New("cli login: invalid credential")

// NewCodeVerifier returns a 48-byte PKCE verifier (64 base64url characters).
func NewCodeVerifier() (string, error) {
	raw := make([]byte, 48)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("cli login: generate code verifier: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

// ChallengeS256 is the S256 transformation of a verifier.
func ChallengeS256(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// VerifyS256 compares a presented verifier with the stored challenge.
func VerifyS256(verifier, challenge string) bool {
	if verifier == "" || challenge == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(ChallengeS256(verifier)), []byte(challenge)) == 1
}

// ValidateLoopbackRedirect accepts only http://<loopback>:10000-10010/login,
// which is the redirect Terraform opens for `terraform login`.
func ValidateLoopbackRedirect(raw string) error {
	uri, err := url.Parse(raw)
	if err != nil {
		return errors.New("redirect_uri is not a valid URI")
	}
	if uri.Scheme != "http" {
		return errors.New("redirect_uri must use http on a loopback address")
	}
	switch uri.Hostname() {
	case "localhost", "127.0.0.1", "::1":
	default:
		return errors.New("redirect_uri must target a loopback address")
	}
	if uri.Port() == "" {
		return errors.New("redirect_uri port is outside the allowed range")
	}
	port, err := strconv.Atoi(uri.Port())
	if err != nil || port < PortLow || port > PortHigh {
		return errors.New("redirect_uri port is outside the allowed range")
	}
	if uri.Path != "/login" {
		return errors.New("redirect_uri path must be /login")
	}
	if uri.RawQuery != "" || uri.Fragment != "" || uri.User != nil {
		return errors.New("redirect_uri must not contain a query, fragment, or userinfo")
	}
	return nil
}

// CookieSigner signs the browser-binding cookie with a key derived from the PAT
// secret, so the login cookie and the tokens are independent secrets on the
// wire even though one configuration value feeds both.
type CookieSigner struct {
	key []byte
}

// NewCookieSigner derives the cookie key with HKDF-SHA256, as upstream does.
func NewCookieSigner(patSecret string) (*CookieSigner, error) {
	secret := strings.TrimSpace(patSecret)
	ikm, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil {
		ikm, err = base64.URLEncoding.DecodeString(secret)
	}
	if err != nil {
		return nil, fmt.Errorf("cli login: decode cookie secret: %w", err)
	}
	key, err := hkdf.Key(sha256.New, ikm, nil, "terrakube-cli-login-cookie-v1", 32)
	if err != nil {
		return nil, fmt.Errorf("cli login: derive cookie key: %w", err)
	}
	return &CookieSigner{key: key}, nil
}

// Value returns the signed cookie payload for a session id.
func (c *CookieSigner) Value(sessionID string) string {
	return sessionID + "." + base64.RawURLEncoding.EncodeToString(c.mac(sessionID))
}

// SessionID verifies a cookie payload and returns the session it binds.
func (c *CookieSigner) SessionID(value string) (string, bool) {
	if value == "" {
		return "", false
	}
	dot := strings.LastIndex(value, ".")
	if dot <= 0 || dot == len(value)-1 {
		return "", false
	}
	sessionID := value[:dot]
	signature, err := base64.RawURLEncoding.DecodeString(value[dot+1:])
	if err != nil || !hmac.Equal(signature, c.mac(sessionID)) {
		return "", false
	}
	return sessionID, true
}

// Cookie is the browser-binding cookie for a session.
func (c *CookieSigner) Cookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    c.Value(sessionID),
		Path:     CookiePath,
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

// Clear expires the cookie when the flow ends.
func (c *CookieSigner) Clear() *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     CookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (c *CookieSigner) mac(data string) []byte {
	mac := hmac.New(sha256.New, c.key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}
