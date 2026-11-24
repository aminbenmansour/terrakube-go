// Package tfoutput serves the executor log endpoints at /tfoutput/v1. It
// mirrors the upstream TerraformOutputController: a plain log read for a job
// step and a Server-Sent Events stream that resumes from Last-Event-ID.
//
// Live records are stored in the Redis stream named after the job, with the
// fields "stepId" and "output". The same layout is used by the executor and by
// the upstream Java implementation.
package tfoutput

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	internalauth "github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

const (
	liveTailRecords   int64 = 5000
	sseMaxDuration          = 30 * time.Minute
	sseReadBlock            = 2 * time.Second
	sseHeartbeatEvery       = 8 // empty reads between heartbeats (about 16s)
	sseReadCount      int64 = 100
)

// Deps are the collaborators the log endpoints need.
type Deps struct {
	Store        store.Querier
	Queue        *queue.Queue
	Storage      storage.Backend
	Layout       tfstate.Layout
	Logger       *slog.Logger
	AuthToken    string
	InternalAuth *internalauth.Validator
}

// Handler serves the /tfoutput/v1 routes.
type Handler struct {
	deps Deps
}

// New builds the log router.
func New(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}

	h := &Handler{deps: deps}
	r := chi.NewRouter()
	r.Use(h.authenticate)
	r.Get("/organization/{organizationID}/job/{jobID}/step/{stepID}", h.getStepLog)
	r.Get("/organization/{organizationID}/job/{jobID}/step/{stepID}/stream", h.streamStepLog)
	return r
}

func (h *Handler) authenticate(next http.Handler) http.Handler {
	if h.deps.AuthToken == "" && h.deps.InternalAuth == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := internalauth.ExtractToken(r.Header.Get("Authorization"))
		if raw == "" {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if h.deps.AuthToken != "" && subtle.ConstantTimeCompare([]byte(raw), []byte(h.deps.AuthToken)) == 1 {
			next.ServeHTTP(w, r)
			return
		}
		if h.deps.InternalAuth != nil {
			if claims, err := h.deps.InternalAuth.Parse(raw); err == nil &&
				internalauth.Classify(claims) != internalauth.TypeNotInternal {
				next.ServeHTTP(w, r.WithContext(internalauth.WithClaims(r.Context(), claims)))
				return
			}
		}
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
}

func (h *Handler) getStepLog(w http.ResponseWriter, r *http.Request) {
	job, step, ok := h.resolve(w, r)
	if !ok {
		return
	}

	// A step that is still running is served from the live stream. Finished
	// steps use the archived output written by the executor.
	if !isTerminal(step.Status) && h.deps.Queue != nil {
		messages, err := h.deps.Queue.HeadStream(r.Context(), streamKey(job.ID), liveTailRecords)
		if err == nil {
			var live strings.Builder
			for _, message := range messages {
				if stepOf(message) == step.ID {
					live.WriteString(outputOf(message))
				}
			}
			if live.Len() > 0 {
				w.Header().Set("Content-Type", "application/octet-stream")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(live.String()))
				return
			}
		}
	}

	if h.deps.Storage != nil {
		key := h.deps.Layout.Output(job.OrganizationID, job.ID, step.ID)
		if rc, err := h.deps.Storage.Get(r.Context(), key); err == nil {
			defer rc.Close()
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			if _, err := io.Copy(w, rc); err != nil {
				h.deps.Logger.Warn("stream archived step log failed", "job_id", job.ID, "step_id", step.ID, "error", err)
			}
			return
		} else if !errors.Is(err, storage.ErrNotFound) {
			h.deps.Logger.Warn("read archived step log failed", "job_id", job.ID, "step_id", step.ID, "error", err)
		}
	}

	if step.Output == nil || *step.Output == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(*step.Output))
}

func (h *Handler) streamStepLog(w http.ResponseWriter, r *http.Request) {
	job, step, ok := h.resolve(w, r)
	if !ok {
		return
	}
	if h.deps.Queue == nil {
		http.Error(w, "log stream is not configured", http.StatusServiceUnavailable)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming is not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	lastID := parseResumeID(r.Header.Get("Last-Event-ID"))
	stream := streamKey(job.ID)

	if lastID != "0-0" {
		if oldest, err := h.deps.Queue.FirstStreamID(r.Context(), stream); err == nil && oldest != "" && compareStreamIDs(lastID, oldest) < 0 {
			writeEvent(w, "gap", fmt.Sprintf(`{"reason":"trimmed","afterId":%q,"firstAvailableId":%q,"recovery":"step-log"}`, lastID, oldest))
			lastID = oldest
		}
	}
	flusher.Flush()

	deadline := time.Now().Add(sseMaxDuration)
	emptyReads := 0
	for {
		if time.Now().After(deadline) {
			writeEvent(w, "end", `{"reason":"rotate"}`)
			flusher.Flush()
			return
		}

		messages, err := h.deps.Queue.ReadStream(r.Context(), stream, lastID, sseReadCount, sseReadBlock)
		if err != nil {
			if r.Context().Err() != nil {
				return
			}
			h.deps.Logger.Warn("step log stream read failed", "job_id", job.ID, "error", err)
			return
		}

		if len(messages) == 0 {
			emptyReads++
			if current, err := h.deps.Store.GetJob(r.Context(), job.ID); err == nil && isTerminal(&current.Status) {
				writeEvent(w, "end", `{"reason":"complete"}`)
				flusher.Flush()
				return
			}
			if emptyReads%sseHeartbeatEvery == 0 {
				_, _ = fmt.Fprint(w, ": heartbeat\n\n")
				flusher.Flush()
			}
			continue
		}

		emptyReads = 0
		for _, message := range messages {
			lastID = message.ID
			if stepOf(message) != step.ID {
				continue
			}
			writeData(w, message.ID, outputOf(message))
		}
		flusher.Flush()
	}
}

// resolve loads the job and step and verifies that they belong together and to
// the organization in the path.
func (h *Handler) resolve(w http.ResponseWriter, r *http.Request) (store.Job, store.Step, bool) {
	jobID, err := strconv.ParseInt(chi.URLParam(r, "jobID"), 10, 64)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return store.Job{}, store.Step{}, false
	}
	job, err := h.deps.Store.GetJob(r.Context(), jobID)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "not found", http.StatusNotFound)
		return store.Job{}, store.Step{}, false
	}
	if err != nil {
		h.deps.Logger.Error("load job for logs", "job_id", jobID, "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return store.Job{}, store.Step{}, false
	}
	if !strings.EqualFold(job.OrganizationID, chi.URLParam(r, "organizationID")) {
		http.Error(w, "not found", http.StatusNotFound)
		return store.Job{}, store.Step{}, false
	}
	if claims, ok := internalauth.ClaimsFromContext(r.Context()); ok &&
		internalauth.Classify(claims) == internalauth.TypeWorkspaceExecutor &&
		(!strings.EqualFold(claims.OrganizationID, job.OrganizationID) ||
			!strings.EqualFold(claims.WorkspaceID, job.WorkspaceID)) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return store.Job{}, store.Step{}, false
	}

	step, err := h.deps.Store.GetStep(r.Context(), chi.URLParam(r, "stepID"))
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "not found", http.StatusNotFound)
		return store.Job{}, store.Step{}, false
	}
	if err != nil {
		h.deps.Logger.Error("load step for logs", "step_id", chi.URLParam(r, "stepID"), "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return store.Job{}, store.Step{}, false
	}
	if step.JobID != job.ID {
		http.Error(w, "not found", http.StatusNotFound)
		return store.Job{}, store.Step{}, false
	}
	return job, step, true
}

func writeEvent(w http.ResponseWriter, event, data string) {
	_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
}

// writeData writes one SSE data frame. Multi-line output is split into several
// data fields so a raw newline never breaks the framing.
func writeData(w http.ResponseWriter, id, output string) {
	_, _ = fmt.Fprintf(w, "id: %s\n", id)
	for _, line := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		_, _ = fmt.Fprintf(w, "data: %s\n", line)
	}
	_, _ = fmt.Fprint(w, "\n")
}

func streamKey(jobID int64) string {
	return strconv.FormatInt(jobID, 10)
}

func stepOf(message queue.Message) string {
	if value, ok := message.Values["stepId"]; ok {
		return fmt.Sprint(value)
	}
	return ""
}

func outputOf(message queue.Message) string {
	if value, ok := message.Values["output"]; ok {
		return fmt.Sprint(value)
	}
	return ""
}

func parseResumeID(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "*" {
		return "0-0"
	}
	if !validStreamID(raw) {
		return "0-0"
	}
	return raw
}

func validStreamID(id string) bool {
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return false
	}
	for _, part := range parts {
		if part == "" {
			return false
		}
		for _, r := range part {
			if r < '0' || r > '9' {
				return false
			}
		}
	}
	return true
}

// compareStreamIDs orders Redis stream IDs of the form <millis>-<sequence>.
func compareStreamIDs(a, b string) int {
	aMillis, aSeq, aOK := splitStreamID(a)
	bMillis, bSeq, bOK := splitStreamID(b)
	if !aOK || !bOK {
		return 0
	}
	if aMillis != bMillis {
		if aMillis < bMillis {
			return -1
		}
		return 1
	}
	if aSeq < bSeq {
		return -1
	}
	if aSeq > bSeq {
		return 1
	}
	return 0
}

func splitStreamID(id string) (int64, int64, bool) {
	parts := strings.SplitN(id, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	millis, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, false
	}
	seq, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, false
	}
	return millis, seq, true
}

func isTerminal(status *string) bool {
	if status == nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(*status)) {
	case "completed", "failed", "cancelled", "canceled", "nochanges", "rejected", "notexecuted", "neverexecuted":
		return true
	default:
		return false
	}
}
