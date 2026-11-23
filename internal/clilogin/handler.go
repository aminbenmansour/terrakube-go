package clilogin

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Deps are the collaborators the broker handler needs.
type Deps struct {
	Service   *Service
	Cookie    *CookieSigner
	PublicURL string
	Logger    *slog.Logger
}

// New builds the /oauth router. It is mounted outside /api/v1 because upstream
// serves the broker there, and it answers HTML or redirects rather than
// JSON:API.
func New(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	deps.PublicURL = strings.TrimRight(deps.PublicURL, "/")

	r := chi.NewRouter()
	r.Get("/authorize", func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		sessionID, location, err := deps.Service.StartAuthorization(req.Context(), AuthorizeRequest{
			ClientID:            query.Get("client_id"),
			RedirectURI:         query.Get("redirect_uri"),
			ResponseType:        query.Get("response_type"),
			CodeChallenge:       query.Get("code_challenge"),
			CodeChallengeMethod: query.Get("code_challenge_method"),
			State:               query.Get("state"),
		})
		if err != nil {
			deps.badRequest(w, "authorize", err)
			return
		}
		http.SetCookie(w, deps.Cookie.Cookie(sessionID))
		http.Redirect(w, req, location, http.StatusFound)
	})

	r.Get("/callback", func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		state := query.Get("state")
		if !deps.boundToBrowser(req, state) {
			deps.htmlForbidden(w, "This login was not started from this browser. Run terraform login again.")
			return
		}
		if oauthError := query.Get("error"); oauthError != "" {
			location, err := deps.Service.CallbackError(req.Context(), state, oauthError)
			if err != nil {
				deps.badRequest(w, "callback error", err)
				return
			}
			http.SetCookie(w, deps.Cookie.Clear())
			http.Redirect(w, req, location, http.StatusFound)
			return
		}
		sessionID, err := deps.Service.Callback(req.Context(), query.Get("code"), state, query.Get("iss"))
		if err != nil {
			deps.badRequest(w, "callback", err)
			return
		}
		http.SetCookie(w, deps.Cookie.Cookie(sessionID))
		http.Redirect(w, req, deps.PublicURL+"/oauth/consent", http.StatusFound)
	})

	r.Get("/consent", func(w http.ResponseWriter, req *http.Request) {
		sessionID, ok := deps.cookieSession(req)
		if !ok {
			deps.htmlForbidden(w, "Missing or invalid session.")
			return
		}
		session, err := deps.Service.ConsentSession(req.Context(), sessionID)
		if err != nil {
			deps.badRequest(w, "consent", err)
			return
		}
		writeHTML(w, http.StatusOK, ConsentPage(stringValue(session.IdentityEmail),
			deps.Service.DefaultDays(), deps.Service.MaxDays(), ""))
	})

	r.Post("/consent", func(w http.ResponseWriter, req *http.Request) {
		sessionID, ok := deps.cookieSession(req)
		if !ok {
			deps.htmlForbidden(w, "Missing or invalid session.")
			return
		}
		if !originMatches(req, deps.PublicURL) {
			deps.htmlForbidden(w, "Bad request origin.")
			return
		}
		if err := req.ParseForm(); err != nil {
			deps.badRequest(w, "consent form", err)
			return
		}
		days, _ := strconv.Atoi(req.PostForm.Get("days"))
		var location string
		var err error
		if req.PostForm.Get("decision") == "deny" {
			location, err = deps.Service.Deny(req.Context(), sessionID)
		} else {
			location, err = deps.Service.Authorize(req.Context(), sessionID, days, req.PostForm.Get("name"))
		}
		if err != nil {
			if errors.Is(err, ErrBadRequest) {
				if session, loadErr := deps.Service.ConsentSession(req.Context(), sessionID); loadErr == nil {
					writeHTML(w, http.StatusBadRequest, ConsentPage(stringValue(session.IdentityEmail),
						deps.Service.DefaultDays(), deps.Service.MaxDays(), err.Error()))
					return
				}
			}
			deps.badRequest(w, "consent", err)
			return
		}
		http.SetCookie(w, deps.Cookie.Clear())
		http.Redirect(w, req, location, http.StatusFound)
	})

	r.Post("/token", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_request"})
			return
		}
		response, err := deps.Service.Exchange(req.Context(),
			req.PostForm.Get("grant_type"),
			req.PostForm.Get("code"),
			req.PostForm.Get("code_verifier"),
			req.PostForm.Get("redirect_uri"),
			req.PostForm.Get("client_id"))
		if err != nil {
			if errors.Is(err, ErrBadRequest) {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_grant"})
				return
			}
			deps.Logger.Error("terraform login token exchange failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server_error"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"access_token": response.AccessToken,
			"token_type":   "Bearer",
			"expires_in":   response.ExpiresIn,
		})
	})
	return r
}

func (deps Deps) cookieSession(req *http.Request) (string, bool) {
	cookie, err := req.Cookie(CookieName)
	if err != nil {
		return "", false
	}
	return deps.Cookie.SessionID(cookie.Value)
}

func (deps Deps) boundToBrowser(req *http.Request, state string) bool {
	sessionID, ok := deps.cookieSession(req)
	return ok && state != "" && sessionID == state
}

func (deps Deps) badRequest(w http.ResponseWriter, operation string, err error) {
	if errors.Is(err, ErrBadRequest) {
		writeHTML(w, http.StatusBadRequest, ErrorPage(err.Error()))
		return
	}
	deps.Logger.Error("terraform login broker failed", "operation", operation, "error", err)
	writeHTML(w, http.StatusInternalServerError, ErrorPage("The login could not be completed."))
}

func (deps Deps) htmlForbidden(w http.ResponseWriter, message string) {
	writeHTML(w, http.StatusForbidden, ErrorPage(message))
}

// originMatches is upstream's CSRF check on the consent form: the browser must
// have loaded the page from this API's origin.
func originMatches(req *http.Request, publicURL string) bool {
	expected, err := url.Parse(publicURL)
	if err != nil {
		return false
	}
	for _, header := range []string{"Origin", "Referer"} {
		raw := req.Header.Get(header)
		if raw == "" {
			continue
		}
		actual, err := url.Parse(raw)
		if err != nil {
			return false
		}
		return strings.EqualFold(expected.Scheme, actual.Scheme) &&
			strings.EqualFold(expected.Hostname(), actual.Hostname()) &&
			effectivePort(expected) == effectivePort(actual)
	}
	return false
}

func effectivePort(uri *url.URL) string {
	if uri.Port() != "" {
		return uri.Port()
	}
	if strings.EqualFold(uri.Scheme, "https") {
		return "443"
	}
	return "80"
}

func writeHTML(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
