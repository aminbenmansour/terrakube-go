package clilogin

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/pat"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// Store is the subset of the generated queries the broker needs.
type Store interface {
	CreateCliAuthSession(ctx context.Context, arg store.CreateCliAuthSessionParams) (store.CliAuthSession, error)
	GetCliAuthSession(ctx context.Context, id string) (store.CliAuthSession, error)
	SetCliAuthSessionIdentity(ctx context.Context, arg store.SetCliAuthSessionIdentityParams) error
	FailCliAuthSession(ctx context.Context, id string) error
	IssueCliAuthCode(ctx context.Context, arg store.IssueCliAuthCodeParams) error
	DenyCliAuthSession(ctx context.Context, id string) error
	GetCliAuthSessionByCodeHash(ctx context.Context, authCodeHash *string) (store.CliAuthSession, error)
	ClaimCliAuthCode(ctx context.Context, id string) (int64, error)
	DeleteExpiredCliAuthSessions(ctx context.Context) (int64, error)
}

// TokenIssuer mints the personal access token the CLI receives.
type TokenIssuer interface {
	Issue(ctx context.Context, user auth.Identity, days int, description, source, createdBy string) (string, error)
}

// Service runs one terraform login at a time.
type Service struct {
	store       Store
	dex         Exchanger
	pats        TokenIssuer
	defaultDays int
	maxDays     int
}

// New builds the service.
func New(store Store, dex Exchanger, pats TokenIssuer, defaultDays, maxDays int) *Service {
	if maxDays > 365 {
		maxDays = 365
	}
	if maxDays < 1 {
		maxDays = 1
	}
	if defaultDays < 1 {
		defaultDays = 1
	}
	if defaultDays > maxDays {
		defaultDays = maxDays
	}
	return &Service{
		store:       store,
		dex:         dex,
		pats:        pats,
		defaultDays: defaultDays,
		maxDays:     maxDays,
	}
}

// DefaultDays and MaxDays drive the consent form.
func (s *Service) DefaultDays() int { return s.defaultDays }
func (s *Service) MaxDays() int     { return s.maxDays }

// AuthorizeRequest is the CLI's authorization request.
type AuthorizeRequest struct {
	ClientID            string
	RedirectURI         string
	ResponseType        string
	CodeChallenge       string
	CodeChallengeMethod string
	State               string
}

// StartAuthorization validates the CLI request and returns the Dex redirect.
func (s *Service) StartAuthorization(ctx context.Context, req AuthorizeRequest) (string, string, error) {
	if req.ClientID != ClientID {
		return "", "", fmt.Errorf("%w: unknown client_id", ErrBadRequest)
	}
	if req.ResponseType != "code" {
		return "", "", fmt.Errorf("%w: response_type must be code", ErrBadRequest)
	}
	if method := req.CodeChallengeMethod; method != "" && method != "S256" {
		return "", "", fmt.Errorf("%w: code_challenge_method must be S256", ErrBadRequest)
	}
	if strings.TrimSpace(req.CodeChallenge) == "" {
		return "", "", fmt.Errorf("%w: code_challenge is required", ErrBadRequest)
	}
	if strings.TrimSpace(req.State) == "" {
		return "", "", fmt.Errorf("%w: state is required", ErrBadRequest)
	}
	if err := ValidateLoopbackRedirect(req.RedirectURI); err != nil {
		return "", "", fmt.Errorf("%w: %s", ErrBadRequest, err)
	}

	// Opportunistic cleanup; the table would otherwise grow forever.
	_, _ = s.store.DeleteExpiredCliAuthSessions(ctx)

	verifier, err := NewCodeVerifier()
	if err != nil {
		return "", "", err
	}
	row, err := s.store.CreateCliAuthSession(ctx, store.CreateCliAuthSessionParams{
		ID:               uuid.NewString(),
		CliRedirectUri:   req.RedirectURI,
		CliCodeChallenge: req.CodeChallenge,
		CliState:         req.State,
		DexCodeVerifier:  &verifier,
		ExpiresAt:        timestamp(time.Now().Add(sessionTTL)),
	})
	if err != nil {
		return "", "", fmt.Errorf("cli login: create session: %w", err)
	}
	return row.ID, s.dex.AuthorizeURL(row.ID, ChallengeS256(verifier)), nil
}

// Callback exchanges the Dex code and moves the session to consent.
func (s *Service) Callback(ctx context.Context, code, state, iss string) (string, error) {
	row, err := s.active(ctx, state)
	if err != nil {
		return "", err
	}
	if row.Status != "PENDING_IDP" {
		return "", fmt.Errorf("%w: session is not awaiting identity provider", ErrBadRequest)
	}
	if iss != "" && iss != s.dex.Issuer() {
		s.fail(ctx, row.ID)
		return "", fmt.Errorf("%w: issuer mismatch", ErrBadRequest)
	}
	if strings.TrimSpace(code) == "" {
		return "", fmt.Errorf("%w: missing authorization code", ErrBadRequest)
	}
	verifier := ""
	if row.DexCodeVerifier != nil {
		verifier = *row.DexCodeVerifier
	}
	identity, err := s.dex.Exchange(ctx, code, verifier)
	if err != nil {
		s.fail(ctx, row.ID)
		return "", err
	}

	groups, err := json.Marshal(identity.Groups)
	if err != nil {
		return "", fmt.Errorf("cli login: encode groups: %w", err)
	}
	email, name, groupList := identity.Email, identity.Name, string(groups)
	if err := s.store.SetCliAuthSessionIdentity(ctx, store.SetCliAuthSessionIdentityParams{
		ID:             row.ID,
		IdentityEmail:  &email,
		IdentityName:   &name,
		IdentityGroups: &groupList,
	}); err != nil {
		return "", fmt.Errorf("cli login: store identity: %w", err)
	}
	return row.ID, nil
}

// CallbackError handles a denial at the identity provider and returns the CLI
// redirect that makes terraform exit cleanly.
func (s *Service) CallbackError(ctx context.Context, state, oauthError string) (string, error) {
	row, err := s.active(ctx, state)
	if err != nil {
		return "", err
	}
	if row.Status != "PENDING_IDP" {
		return "", fmt.Errorf("%w: session is not awaiting identity provider", ErrBadRequest)
	}
	s.fail(ctx, row.ID)
	code := "server_error"
	if oauthError == "access_denied" {
		code = "access_denied"
	}
	return cliRedirect(row.CliRedirectUri, "error", code, row.CliState), nil
}

// ConsentSession returns the pending session for the consent page.
func (s *Service) ConsentSession(ctx context.Context, sessionID string) (store.CliAuthSession, error) {
	row, err := s.active(ctx, sessionID)
	if err != nil {
		return store.CliAuthSession{}, err
	}
	if row.Status != "PENDING_CONSENT" {
		return store.CliAuthSession{}, fmt.Errorf("%w: session is not awaiting consent", ErrBadRequest)
	}
	return row, nil
}

// Authorize records the consent choice and returns the CLI redirect with its
// one-time code.
func (s *Service) Authorize(ctx context.Context, sessionID string, days int, name string) (string, error) {
	row, err := s.ConsentSession(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if days < 1 || days > s.maxDays {
		return "", fmt.Errorf("%w: days must be between 1 and %d", ErrBadRequest, s.maxDays)
	}

	code, err := randomToken()
	if err != nil {
		return "", err
	}
	chosen := int32(days)
	hash := hashCode(code)
	tokenName := sanitizeName(name)
	if err := s.store.IssueCliAuthCode(ctx, store.IssueCliAuthCodeParams{
		ID:            row.ID,
		AuthCodeHash:  &hash,
		ChosenDays:    &chosen,
		ChosenName:    tokenName,
		CodeExpiresAt: timestamp(time.Now().Add(codeTTL)),
	}); err != nil {
		return "", fmt.Errorf("cli login: issue code: %w", err)
	}
	return cliRedirect(row.CliRedirectUri, "code", code, row.CliState), nil
}

// Deny records a refusal and returns the CLI redirect.
func (s *Service) Deny(ctx context.Context, sessionID string) (string, error) {
	row, err := s.ConsentSession(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if err := s.store.DenyCliAuthSession(ctx, row.ID); err != nil {
		return "", fmt.Errorf("cli login: deny: %w", err)
	}
	return cliRedirect(row.CliRedirectUri, "error", "access_denied", row.CliState), nil
}

// TokenResponse is the body the CLI expects from the token endpoint.
type TokenResponse struct {
	AccessToken string
	ExpiresIn   int64
}

// Exchange redeems the one-time code for a personal access token.
func (s *Service) Exchange(ctx context.Context, grantType, code, verifier, redirectURI, clientID string) (TokenResponse, error) {
	if grantType != "authorization_code" {
		return TokenResponse{}, fmt.Errorf("%w: invalid_grant", ErrBadRequest)
	}
	hash := hashCode(code)
	row, err := s.store.GetCliAuthSessionByCodeHash(ctx, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return TokenResponse{}, fmt.Errorf("%w: invalid_grant", ErrBadRequest)
	}
	if err != nil {
		return TokenResponse{}, fmt.Errorf("cli login: load code: %w", err)
	}

	expired := !row.CodeExpiresAt.Valid || row.CodeExpiresAt.Time.Before(time.Now())
	if row.Status != "CODE_ISSUED" || expired ||
		row.CliRedirectUri != redirectURI || clientID != ClientID ||
		!VerifyS256(verifier, row.CliCodeChallenge) {
		return TokenResponse{}, fmt.Errorf("%w: invalid_grant", ErrBadRequest)
	}
	claimed, err := s.store.ClaimCliAuthCode(ctx, row.ID)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("cli login: claim code: %w", err)
	}
	if claimed != 1 {
		return TokenResponse{}, fmt.Errorf("%w: invalid_grant", ErrBadRequest)
	}

	identity := auth.Identity{
		Issuer: "cli_login",
		Email:  stringValue(row.IdentityEmail),
		Name:   stringValue(row.IdentityName),
		Groups: decodeGroups(row.IdentityGroups),
	}
	days := int(int32Value(row.ChosenDays))
	description := stringValue(row.ChosenName)
	if description == "" {
		description = "terraform login " + time.Now().UTC().Format(time.DateOnly)
	}
	token, err := s.pats.Issue(ctx, identity, days, description, pat.SourceCLILogin, identity.Email)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("cli login: mint token: %w", err)
	}
	return TokenResponse{AccessToken: token, ExpiresIn: int64(days) * 86400}, nil
}

func (s *Service) active(ctx context.Context, sessionID string) (store.CliAuthSession, error) {
	if _, err := uuid.Parse(sessionID); err != nil {
		return store.CliAuthSession{}, fmt.Errorf("%w: session not found", ErrBadRequest)
	}
	row, err := s.store.GetCliAuthSession(ctx, sessionID)
	if errors.Is(err, pgx.ErrNoRows) {
		return store.CliAuthSession{}, fmt.Errorf("%w: session not found", ErrBadRequest)
	}
	if err != nil {
		return store.CliAuthSession{}, fmt.Errorf("cli login: load session: %w", err)
	}
	if !row.ExpiresAt.Valid || row.ExpiresAt.Time.Before(time.Now()) {
		return store.CliAuthSession{}, fmt.Errorf("%w: session expired", ErrBadRequest)
	}
	return row, nil
}

func (s *Service) fail(ctx context.Context, sessionID string) {
	_ = s.store.FailCliAuthSession(ctx, sessionID)
}

func cliRedirect(base, key, value, state string) string {
	separator := "?"
	if strings.Contains(base, "?") {
		separator = "&"
	}
	return base + separator + key + "=" + url.QueryEscape(value) + "&state=" + url.QueryEscape(state)
}

func randomToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("cli login: generate code: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func hashCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}

func sanitizeName(raw string) *string {
	trimmed := strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, raw)
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return nil
	}
	if len(trimmed) > 64 {
		trimmed = trimmed[:64]
	}
	return &trimmed
}

func decodeGroups(raw *string) []string {
	if raw == nil || *raw == "" {
		return nil
	}
	var groups []string
	if err := json.Unmarshal([]byte(*raw), &groups); err != nil {
		return nil
	}
	return groups
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func int32Value(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
}

func timestamp(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: value, Valid: true}
}
