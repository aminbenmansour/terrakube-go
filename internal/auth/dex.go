package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Verifier checks a user bearer token.
type Verifier interface {
	Verify(ctx context.Context, raw string) (*Identity, error)
}

// DexValidator verifies the Dex JWT the upstream UI sends as a bearer token.
// It mirrors DexAuthenticationManagerResolver: issuer metadata is discovered
// once, signing keys are cached, and an unknown kid triggers one refetch so
// Dex key rotation does not lock users out.
type DexValidator struct {
	issuer           string
	audiences        []string
	validateAudience bool
	client           *http.Client

	// refreshInterval bounds how often an unknown kid can trigger a JWKS
	// fetch; zero refreshes on every miss (tests).
	refreshInterval time.Duration

	mu          sync.Mutex
	jwksURI     string
	keys        map[string]*rsa.PublicKey
	lastRefresh time.Time
}

// NewDexValidator builds a validator for issuerURL. No network call is made
// until the first token is verified, so a temporarily unavailable Dex does not
// stop the API from starting.
func NewDexValidator(issuerURL string, audiences []string, validateAudience bool) *DexValidator {
	return &DexValidator{
		issuer:           strings.TrimRight(strings.TrimSpace(issuerURL), "/"),
		audiences:        audiences,
		validateAudience: validateAudience,
		client:           &http.Client{Timeout: 10 * time.Second},
		refreshInterval:  time.Minute,
	}
}

// Verify checks the signature against the issuer's JWKS, then iss, exp and,
// when enabled, aud. The returned error never contains the token.
func (v *DexValidator) Verify(ctx context.Context, raw string) (*Identity, error) {
	if v == nil || v.issuer == "" {
		return nil, errors.New("auth: dex validator is not configured")
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("auth: dex token is empty")
	}

	claims := jwt.MapClaims{}
	// shortcut: Dex's local signer is RS256; add ES256/EdDSA when a Vault signer is used.
	_, err := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		jwt.WithIssuer(v.issuer),
		// exp is checked when present, like Spring's JwtTimestampValidator.
		jwt.WithLeeway(60*time.Second),
	).ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		kid, _ := token.Header["kid"].(string)
		return v.key(ctx, kid)
	})
	if err != nil {
		return nil, fmt.Errorf("auth: verify dex token: %w", err)
	}

	if v.validateAudience {
		audiences, err := claims.GetAudience()
		if err != nil {
			return nil, errors.New("auth: dex token audience is not accepted")
		}
		accepted := false
		for _, audience := range v.audiences {
			accepted = accepted || contains(audiences, audience)
		}
		if !accepted {
			return nil, errors.New("auth: dex token audience is not accepted")
		}
	}

	identity := &Identity{Issuer: v.issuer}
	identity.Subject, _ = claims["sub"].(string)
	identity.Email, _ = claims["email"].(string)
	identity.Name, _ = claims["name"].(string)
	identity.Groups = stringSlice(claims["groups"])
	if identity.Name == "" {
		identity.Name = identity.Email
	}
	return identity, nil
}

func stringSlice(value any) []string {
	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok {
				values = append(values, text)
			}
		}
		return values
	default:
		return nil
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func (v *DexValidator) key(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	if kid == "" {
		return nil, errors.New("auth: dex token has no kid")
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	if key, ok := v.keys[kid]; ok {
		return key, nil
	}
	now := time.Now()
	if v.keys == nil || now.Sub(v.lastRefresh) >= v.refreshInterval {
		if err := v.refreshLocked(ctx); err != nil {
			return nil, err
		}
		if key, ok := v.keys[kid]; ok {
			return key, nil
		}
	}
	return nil, fmt.Errorf("auth: dex signing key %q is not published", kid)
}

func (v *DexValidator) refreshLocked(ctx context.Context) error {
	if v.jwksURI == "" {
		var discovery struct {
			Issuer  string `json:"issuer"`
			JWKSURI string `json:"jwks_uri"`
		}
		if err := v.getJSON(ctx, v.issuer+"/.well-known/openid-configuration", &discovery); err != nil {
			return err
		}
		if discovery.Issuer != v.issuer {
			return fmt.Errorf("auth: dex discovery issuer %q does not match %q", discovery.Issuer, v.issuer)
		}
		if discovery.JWKSURI == "" {
			return errors.New("auth: dex discovery has no jwks_uri")
		}
		v.jwksURI = discovery.JWKSURI
	}

	var set struct {
		Keys []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			Use string `json:"use"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := v.getJSON(ctx, v.jwksURI, &set); err != nil {
		return err
	}

	keys := make(map[string]*rsa.PublicKey, len(set.Keys))
	for _, jwk := range set.Keys {
		if jwk.Kty != "RSA" || jwk.Kid == "" || (jwk.Use != "" && jwk.Use != "sig") ||
			(jwk.Alg != "" && jwk.Alg != jwt.SigningMethodRS256.Alg()) {
			continue
		}
		key, err := rsaPublicKey(jwk.N, jwk.E)
		if err != nil {
			continue
		}
		keys[jwk.Kid] = key
	}
	if len(keys) == 0 {
		return errors.New("auth: dex JWKS has no usable RSA signing key")
	}
	v.keys = keys
	v.lastRefresh = time.Now()
	return nil
}

func (v *DexValidator) getJSON(ctx context.Context, url string, value any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("auth: build dex metadata request: %w", err)
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return fmt.Errorf("auth: fetch dex metadata: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth: dex metadata returned %s", resp.Status)
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(value); err != nil {
		return fmt.Errorf("auth: decode dex metadata: %w", err)
	}
	return nil
}

func rsaPublicKey(n, e string) (*rsa.PublicKey, error) {
	modulus, err := decodeJWKPart(n)
	if err != nil {
		return nil, err
	}
	exponent, err := decodeJWKPart(e)
	if err != nil {
		return nil, err
	}
	exp := new(big.Int).SetBytes(exponent)
	if !exp.IsInt64() || exp.Int64() <= 1 {
		return nil, errors.New("auth: invalid RSA exponent")
	}
	return &rsa.PublicKey{N: new(big.Int).SetBytes(modulus), E: int(exp.Int64())}, nil
}

func decodeJWKPart(value string) ([]byte, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err == nil {
		return decoded, nil
	}
	return base64.URLEncoding.DecodeString(value)
}
