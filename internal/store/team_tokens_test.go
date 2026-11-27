package store

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
)

// TestTeamTokenQueries covers the metadata round-trip, group listing, the soft
// delete that revocation depends on, and the last-used stamp.
func TestTeamTokenQueries(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	days := int32(0)
	hours := int32(2)
	minutes := int32(30)
	group := "team-" + uuid.NewString()
	description := "ci runner"
	token, err := q.CreateTeamToken(ctx, CreateTeamTokenParams{
		ID: uuid.NewString(), Days: &days, Hours: &hours, Minutes: &minutes,
		Description: &description, GroupName: &group,
	})
	if err != nil {
		t.Fatalf("create team token: %v", err)
	}
	if token.Deleted || token.GroupName == nil || *token.GroupName != group {
		t.Fatalf("created token = %+v", token)
	}

	found, err := q.GetTeamToken(ctx, token.ID)
	if err != nil || found.ID != token.ID {
		t.Fatalf("get team token = %+v/%v", found, err)
	}
	list, err := q.ListTeamTokensByGroups(ctx, []string{group, "someone-else"})
	if err != nil || len(list) != 1 || list[0].ID != token.ID {
		t.Fatalf("list by group = %+v/%v", list, err)
	}
	if other, err := q.ListTeamTokensByGroups(ctx, []string{"someone-else"}); err != nil || len(other) != 0 {
		t.Fatalf("list other group = %+v/%v", other, err)
	}

	if err := q.StampTeamTokenLastUsed(ctx, token.ID); err != nil {
		t.Fatalf("stamp last used: %v", err)
	}
	stamped, err := q.GetTeamToken(ctx, token.ID)
	if err != nil || !stamped.LastUsedAt.Valid {
		t.Fatalf("last used = %+v/%v", stamped.LastUsedAt, err)
	}

	by := "tester"
	if err := q.SoftDeleteTeamToken(ctx, SoftDeleteTeamTokenParams{ID: token.ID, UpdatedBy: &by}); err != nil {
		t.Fatalf("soft delete: %v", err)
	}
	deleted, err := q.GetTeamToken(ctx, token.ID)
	if err != nil || !deleted.Deleted || deleted.UpdatedBy == nil || *deleted.UpdatedBy != by {
		t.Fatalf("deleted token = %+v/%v", deleted, err)
	}

	if _, err := q.GetTeamToken(ctx, uuid.NewString()); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("missing token = %v", err)
	}
}
