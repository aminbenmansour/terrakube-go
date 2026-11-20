package api

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// upstreamCapturesDir holds responses recorded from a live upstream
// (azbuilder/api-server:2.34.3) by scripts/upstream-capture. They are the
// contract, so they are replayed here rather than trusted to a reading of the
// OpenAPI file.
const upstreamCapturesDir = "../../docs/upstream/rest"

type upstreamCapture struct {
	name   string
	status int
	path   string
	body   []byte
}

func (c upstreamCapture) resourceID(t *testing.T) string {
	t.Helper()
	var doc struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(c.body, &doc); err != nil || doc.Data.ID == "" {
		t.Fatalf("%s: no data.id to replay against: %v", c.name, err)
	}
	return doc.Data.ID
}

func loadUpstreamCaptures(t *testing.T) map[string]upstreamCapture {
	t.Helper()

	index, err := os.ReadFile(filepath.Join(upstreamCapturesDir, "index.tsv"))
	if err != nil {
		t.Fatalf("read captured index: %v", err)
	}

	captures := make(map[string]upstreamCapture)
	for _, line := range strings.Split(strings.TrimSpace(string(index)), "\n") {
		name, rest, ok := strings.Cut(line, "\t")
		rawStatus, path, ok2 := strings.Cut(rest, "\t")
		if !ok || !ok2 {
			t.Fatalf("malformed index line %q, want name<TAB>status<TAB>path", line)
		}
		status, err := strconv.Atoi(rawStatus)
		if err != nil {
			t.Fatalf("%s: status %q is not a number", name, rawStatus)
		}
		body, err := os.ReadFile(filepath.Join(upstreamCapturesDir, name+".json"))
		if err != nil {
			t.Fatalf("%s: read capture: %v", name, err)
		}
		captures[name] = upstreamCapture{name: name, status: status, path: path, body: body}
	}
	return captures
}

// replay is how this API answers one captured request, with stub data standing
// in for the database. target "resource" means "/organization/{captured id}".
type replay struct {
	target string
	// status overrides the captured status when we deliberately answer
	// differently; 0 means it must match.
	status  int
	querier *stubQuerier
	// authToken enables authentication, for the unauthenticated capture.
	authToken string
	// idFrom names another capture that carries the resource id, for paths
	// built against an organization described by a different request.
	idFrom string
}

var replayedCaptures = map[string]replay{
	"organization-list": {
		target:  "/organization",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-list-sorted": {
		target:  "/organization?sort=-name",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-list-filtered": {
		target:  "/organization?filter[organization]=name==azure",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-get": {
		target: "resource",
	},
	"organization-projects": {
		target:  "resource-projects",
		idFrom:  "organization-get",
		querier: &stubQuerier{},
	},
	"organization-missing": {
		target:  "/organization/00000000-0000-0000-0000-000000000000",
		querier: &stubQuerier{err: pgx.ErrNoRows},
	},
	"organization-bad-id": {
		target: "/organization/not-a-uuid",
	},
	"organization-no-token": {
		target:    "/organization",
		authToken: "secret",
	},
}

// deferredCaptures are recorded but cannot be asserted yet. Each needs a reason,
// and the test below fails when a capture is added without one, so a new
// recording cannot be silently ignored.
var deferredCaptures = map[string]string{
	"organization-list-totals":   "we accept paging on organization, upstream answers 400; deliberate superset (docs/upstream/rest/README.md)",
	"organization-list-paged":    "same as organization-list-totals",
	"organization-list-fields":   "sparse fieldsets are not implemented, so we return every attribute (JSON:API says MUST NOT); see #121",
	"organization-get-fields":    "same as organization-list-fields",
	"organization-list-include":  "we ignore include= and emit no `included` member; the UI sends no include on /api/v1",
	"organization-relationships": "no /relationships/{relation} routes; the UI patches relationships on workspace only (#116)",
	"workspace-list-paged":       "no global /workspace collection; the capture documents meta.page, not a route we serve",
	"workspace-list-totals":      "same as workspace-list-paged",
	"delete-method-not-allowed":  "upstream answers 500 with a servlet error page, which is not a contract worth copying",
}

// ListProjectsByOrganization is stubbed here because only the capture replay
// needs it; the organization handlers above it do not.
func (s *stubQuerier) ListProjectsByOrganization(context.Context, string) ([]store.Project, error) {
	return nil, s.err
}

// demoOrganizations stands in for upstream's demo seed; only the names matter,
// and no assertion here reads them.
func demoOrganizations() []store.Organization {
	return []store.Organization{
		{ID: "1", Name: "azure", ExecutionMode: "remote"},
		{ID: "2", Name: "aws", ExecutionMode: "local"},
	}
}

func TestEveryUpstreamCaptureIsAccountedFor(t *testing.T) {
	for name := range loadUpstreamCaptures(t) {
		_, replayed := replayedCaptures[name]
		reason, deferred := deferredCaptures[name]
		if !replayed && !deferred {
			t.Errorf("%s: captured from upstream but neither replayed nor deferred; add it to replayedCaptures or explain it in deferredCaptures", name)
		}
		if deferred && strings.TrimSpace(reason) == "" {
			t.Errorf("%s: deferred without a reason", name)
		}
	}
}

func TestCapturedShapesAreWhatWeReturn(t *testing.T) {
	captures := loadUpstreamCaptures(t)

	for name, plan := range replayedCaptures {
		t.Run(name, func(t *testing.T) {
			capture, ok := captures[name]
			if !ok {
				t.Fatalf("no capture named %s", name)
			}

			target, querier := plan.target, plan.querier
			idCapture := capture
			if plan.idFrom != "" {
				idCapture = captures[plan.idFrom]
			}
			switch target {
			case "resource":
				id := idCapture.resourceID(t)
				target = "/organization/" + id
				querier = &stubQuerier{org: store.Organization{ID: id, Name: "azure", ExecutionMode: "remote"}}
			case "resource-projects":
				target = "/organization/" + idCapture.resourceID(t) + "/project"
			}
			if querier == nil {
				querier = &stubQuerier{}
			}

			handler := New(Deps{
				Store:     querier,
				Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
				AuthToken: plan.authToken,
			})
			req := httptest.NewRequest(http.MethodGet, target, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			wantStatus := capture.status
			if plan.status != 0 {
				wantStatus = plan.status
			}
			if rec.Code != wantStatus {
				t.Fatalf("%s: status = %d, want %d (captured); body = %s", capture.path, rec.Code, wantStatus, rec.Body.String())
			}
			assertSameDocumentShape(t, capture, rec.Body.Bytes())
		})
	}
}

// assertSameDocumentShape compares the members a client can see, not values: an
// attribute rename or a missing member must fail, a different id or name must
// not.
func assertSameDocumentShape(t *testing.T, capture upstreamCapture, ours []byte) {
	t.Helper()

	if len(strings.TrimSpace(string(capture.body))) == 0 {
		// Upstream's 401 has an empty body; only the status is comparable.
		return
	}

	var theirs, mine map[string]any
	if err := json.Unmarshal(capture.body, &theirs); err != nil {
		t.Fatalf("%s: captured body is not JSON: %v", capture.name, err)
	}
	if err := json.Unmarshal(ours, &mine); err != nil {
		t.Fatalf("%s: our body is not JSON (%v): %s", capture.name, err, ours)
	}

	if theirErrors, ok := theirs["errors"].([]any); ok {
		// We add status and title; JSON:API allows it and the UI ignores them.
		myErrors, ok := mine["errors"].([]any)
		if !ok || len(myErrors) != len(theirErrors) {
			t.Fatalf("%s: errors = %v, want %d error object(s)", capture.name, mine["errors"], len(theirErrors))
		}
		for i, raw := range theirErrors {
			theirsErr, _ := raw.(map[string]any)
			mineErr, _ := myErrors[i].(map[string]any)
			if theirsErr["detail"] != nil && mineErr["detail"] == nil {
				t.Errorf("%s: errors[%d] has no detail; upstream says %v", capture.name, i, theirsErr["detail"])
			}
		}
		return
	}

	if got, want := sortedKeys(mine), sortedKeys(theirs); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("%s: top-level members = %v, want %v", capture.name, got, want)
	}

	switch their := theirs["data"].(type) {
	case []any:
		my, ok := mine["data"].([]any)
		if !ok {
			t.Fatalf("%s: data = %T, want array", capture.name, mine["data"])
		}
		if len(their) == 0 {
			if len(my) != 0 {
				t.Fatalf("%s: upstream returned an empty collection, we returned %d", capture.name, len(my))
			}
			return
		}
		if len(my) == 0 {
			t.Fatalf("%s: we returned an empty collection, upstream returned %d", capture.name, len(their))
		}
		assertSameResourceShape(t, capture.name, their[0], my[0])
	case map[string]any:
		assertSameResourceShape(t, capture.name, their, mine["data"])
	}
}

func assertSameResourceShape(t *testing.T, name string, theirs, mine any) {
	t.Helper()

	theirResource, ok := theirs.(map[string]any)
	if !ok {
		t.Fatalf("%s: captured data is %T, want a resource object", name, theirs)
	}
	myResource, ok := mine.(map[string]any)
	if !ok {
		t.Fatalf("%s: our data is %T, want a resource object", name, mine)
	}

	if myResource["type"] != theirResource["type"] {
		t.Errorf("%s: type = %v, want %v", name, myResource["type"], theirResource["type"])
	}
	if _, ok := myResource["id"].(string); !ok {
		t.Errorf("%s: id = %v, want a string", name, myResource["id"])
	}

	theirAttrs, _ := theirResource["attributes"].(map[string]any)
	myAttrs, _ := myResource["attributes"].(map[string]any)
	if got, want := sortedKeys(myAttrs), sortedKeys(theirAttrs); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("%s: attributes = %v, want %v", name, got, want)
	}
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
