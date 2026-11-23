package clilogin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
)

// ErrBadRequest marks a client error in the login flow.
var ErrBadRequest = errors.New("cli login: bad request")

// ErrUpstream marks an identity-provider failure.
var ErrUpstream = errors.New("cli login: identity provider failed")

// Exchanger is the Dex side of the flow, injectable for tests.
type Exchanger interface {
	Issuer() string
	AuthorizeURL(state, codeChallenge string) string
	Exchange(ctx context.Context, code, verifier string) (*auth.Identity, error)
}

// Dex talks to the configured Dex issuer. The id_token it receives is checked
// by the same validator the UI bearer tokens use, with the CLI client as the
// required audience.
type Dex struct {
	issuer      string
	clientID    string
	callbackURL string
	tokenURL    string
	validator   *auth.DexValidator
	client      *http.Client
}

// NewDex builds the exchanger. issuer is the Dex issuer URL and callbackURL is
// this API's /oauth/callback.
func NewDex(issuer, clientID, callbackURL string) *Dex {
	issuer = strings.TrimRight(strings.TrimSpace(issuer), "/")
	return &Dex{
		issuer:      issuer,
		clientID:    clientID,
		callbackURL: strings.TrimRight(callbackURL, "/"),
		tokenURL:    issuer + "/token",
		validator:   auth.NewDexValidator(issuer, []string{clientID}, true),
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

// Issuer returns the configured issuer, for the callback's iss check.
func (d *Dex) Issuer() string { return d.issuer }

// AuthorizeURL is the Dex authorization request the browser is sent to.
func (d *Dex) AuthorizeURL(state, codeChallenge string) string {
	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("client_id", d.clientID)
	query.Set("redirect_uri", d.callbackURL)
	query.Set("scope", "openid profile email groups")
	query.Set("code_challenge", codeChallenge)
	query.Set("code_challenge_method", "S256")
	query.Set("state", state)
	return d.issuer + "/auth?" + query.Encode()
}

type dexTokenResponse struct {
	IDToken string `json:"id_token"`
}

// Exchange redeems the authorization code at Dex and returns the verified
// identity. Errors never include the code or the tokens.
func (d *Dex) Exchange(ctx context.Context, code, verifier string) (*auth.Identity, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", d.callbackURL)
	form.Set("client_id", d.clientID)
	form.Set("code_verifier", verifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("%w: build token request", ErrUpstream)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: token request", ErrUpstream)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: token endpoint returned %s", ErrUpstream, resp.Status)
	}
	var body dexTokenResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&body); err != nil {
		return nil, fmt.Errorf("%w: decode token response", ErrUpstream)
	}
	if body.IDToken == "" {
		return nil, fmt.Errorf("%w: token response has no id_token", ErrUpstream)
	}

	identity, err := d.validator.Verify(ctx, body.IDToken)
	if err != nil {
		return nil, fmt.Errorf("%w: id_token rejected", ErrUpstream)
	}
	if identity.Email == "" {
		return nil, fmt.Errorf("%w: id_token has no email", ErrUpstream)
	}
	return identity, nil
}
