package store

import (
	"slices"
	"testing"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
	"github.com/aminbenmansour/terrakube-go/internal/rsql"
)

func TestRunTriggerListFiltersEitherEndAndScopesOrganizations(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	newOrg := func(name string) Organization {
		org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
			ID: uuid.NewString(), Name: name, ExecutionMode: "remote",
		})
		if err != nil {
			t.Fatalf("create organization %q: %v", name, err)
		}
		return org
	}
	newWorkspace := func(org Organization, name string) Workspace {
		ws, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
			ID: uuid.NewString(), Name: name, IacType: "terraform", OrganizationID: org.ID,
		})
		if err != nil {
			t.Fatalf("create workspace %q: %v", name, err)
		}
		return ws
	}
	newTrigger := func(source, destination Workspace, enabled bool) WorkspaceRunTrigger {
		trigger, err := q.CreateRunTrigger(ctx, CreateRunTriggerParams{
			ID:                     uuid.NewString(),
			SourceWorkspaceID:      source.ID,
			DestinationWorkspaceID: destination.ID,
			OrganizationID:         source.OrganizationID,
			Enabled:                enabled,
		})
		if err != nil {
			t.Fatalf("create run trigger: %v", err)
		}
		return trigger
	}

	org := newOrg("triggers")
	otherOrg := newOrg("triggers-other")
	a := newWorkspace(org, "a")
	b := newWorkspace(org, "b")
	c := newWorkspace(org, "c")
	x := newWorkspace(otherOrg, "x")
	y := newWorkspace(otherOrg, "y")

	incoming := newTrigger(a, b, true)
	outgoing := newTrigger(b, c, false)
	newTrigger(x, y, true)

	filter, err := rsql.Parse("sourceWorkspace.id==" + b.ID + ",destinationWorkspace.id==" + b.ID)
	if err != nil {
		t.Fatalf("parse filter: %v", err)
	}

	assertIDs := func(scope []any, want []string) {
		t.Helper()
		rows, err := q.ListRunTriggers(ctx, scope, filter, nil, nil, 0)
		if err != nil {
			t.Fatalf("ListRunTriggers: %v", err)
		}
		got := make([]string, 0, len(rows))
		for _, row := range rows {
			got = append(got, row.ID)
		}
		if len(got) != len(want) {
			t.Fatalf("ids = %v, want %v", got, want)
		}
		for _, id := range want {
			if !slices.Contains(got, id) {
				t.Errorf("ids = %v, want %v", got, want)
			}
		}

		total, err := q.CountRunTriggers(ctx, scope, filter)
		if err != nil {
			t.Fatalf("CountRunTriggers: %v", err)
		}
		if total != int64(len(want)) {
			t.Errorf("count = %d, want %d", total, len(want))
		}
	}

	// An owner (nil scope) sees both edges around b; the other organization's
	// edge matches neither end.
	assertIDs([]any{nil}, []string{incoming.ID, outgoing.ID})
	assertIDs([]any{[]string{org.ID}}, []string{incoming.ID, outgoing.ID})
	assertIDs([]any{[]string{otherOrg.ID}}, nil)
	assertIDs([]any{[]string{}}, nil)
}

func TestListEnabledRunTriggersBySourceSkipsDisabledAndDeleted(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
		ID: uuid.NewString(), Name: "trigger-enabled", ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization: %v", err)
	}
	newWorkspace := func(name string) Workspace {
		ws, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
			ID: uuid.NewString(), Name: name, IacType: "terraform", OrganizationID: org.ID,
		})
		if err != nil {
			t.Fatalf("create workspace %q: %v", name, err)
		}
		return ws
	}
	source, live, disabledDest, deletedDest := newWorkspace("source"), newWorkspace("live"),
		newWorkspace("disabled"), newWorkspace("deleted")
	if err := q.SoftDeleteWorkspace(ctx, deletedDest.ID); err != nil {
		t.Fatalf("soft-delete destination: %v", err)
	}

	newEdge := func(destination Workspace, enabled bool) WorkspaceRunTrigger {
		edge, err := q.CreateRunTrigger(ctx, CreateRunTriggerParams{
			ID: uuid.NewString(), SourceWorkspaceID: source.ID, DestinationWorkspaceID: destination.ID,
			OrganizationID: org.ID, Enabled: enabled,
		})
		if err != nil {
			t.Fatalf("create run trigger: %v", err)
		}
		return edge
	}
	wanted := newEdge(live, true)
	newEdge(disabledDest, false)
	newEdge(deletedDest, true)

	edges, err := q.ListEnabledRunTriggersBySource(ctx, source.ID)
	if err != nil {
		t.Fatalf("ListEnabledRunTriggersBySource: %v", err)
	}
	if len(edges) != 1 || edges[0].ID != wanted.ID {
		t.Errorf("edges = %+v, want only the enabled edge to a live workspace", edges)
	}
}

func TestRunTriggerRepositoryGraphAndFanOut(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
		ID: uuid.NewString(), Name: "trigger-graph", ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization: %v", err)
	}
	workspaces := make([]Workspace, 0, 3)
	for _, name := range []string{"a", "b", "c"} {
		ws, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
			ID: uuid.NewString(), Name: name, IacType: "terraform", OrganizationID: org.ID,
		})
		if err != nil {
			t.Fatalf("create workspace: %v", err)
		}
		workspaces = append(workspaces, ws)
	}

	edge, err := q.CreateRunTrigger(ctx, CreateRunTriggerParams{
		ID:                     uuid.NewString(),
		SourceWorkspaceID:      workspaces[0].ID,
		DestinationWorkspaceID: workspaces[1].ID,
		OrganizationID:         org.ID,
		Enabled:                true,
	})
	if err != nil {
		t.Fatalf("create run trigger: %v", err)
	}

	edges, err := q.ListRunTriggersByOrganization(ctx, org.ID)
	if err != nil {
		t.Fatalf("ListRunTriggersByOrganization: %v", err)
	}
	if len(edges) != 1 || edges[0].ID != edge.ID {
		t.Errorf("edges = %v, want the created edge", edges)
	}

	count, err := q.CountEnabledRunTriggersBySource(ctx, CountEnabledRunTriggersBySourceParams{
		SourceWorkspaceID: workspaces[0].ID,
		ID:                "",
	})
	if err != nil || count != 1 {
		t.Fatalf("count = %d, err = %v; want 1", count, err)
	}
	// The edge itself is excluded, which is how an update re-checks its source.
	count, err = q.CountEnabledRunTriggersBySource(ctx, CountEnabledRunTriggersBySourceParams{
		SourceWorkspaceID: workspaces[0].ID,
		ID:                edge.ID,
	})
	if err != nil || count != 0 {
		t.Fatalf("count excluding the edge = %d, err = %v; want 0", count, err)
	}
}
