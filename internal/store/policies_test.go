package store

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
)

// TestPolicyQueries covers the policy set, attachment, parameter and exemption
// round trips the Settings > Policies screens need, including the cascade when
// a set is deleted.
func TestPolicyQueries(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
		ID: uuid.NewString(), Name: "policies-" + uuid.NewString()[:8], ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization: %v", err)
	}

	set, err := q.CreatePolicySet(ctx, CreatePolicySetParams{
		ID: uuid.NewString(), Name: "capture-policy", Description: strPtr("capture"),
		EnforcementLevel: "SOFT_MANDATORY", OverrideTeam: strPtr("TERRAKUBE_ADMIN"),
		Global: false, Repository: strPtr("https://github.com/example/policies"),
		Branch: strPtr("main"), Folder: strPtr("/"), OrganizationID: org.ID,
		OpaVersion: strPtr("0.68.0"), CreatedBy: strPtr("admin@example.com"),
	})
	if err != nil {
		t.Fatalf("create policy set: %v", err)
	}
	sets, err := q.ListPolicySetsByOrganization(ctx, org.ID)
	if err != nil || len(sets) != 1 || sets[0].ID != set.ID {
		t.Fatalf("list policy sets = %+v/%v", sets, err)
	}
	sets, err = q.ListPolicySetsByOrganizations(ctx, []string{org.ID})
	if err != nil || len(sets) != 1 {
		t.Fatalf("list policy sets by organizations = %+v/%v", sets, err)
	}

	attachment, err := q.CreatePolicyAttachment(ctx, CreatePolicyAttachmentParams{
		ID: uuid.NewString(), PolicySetID: set.ID, CreatedBy: strPtr("admin@example.com"),
	})
	if err != nil {
		t.Fatalf("create policy attachment: %v", err)
	}
	attachments, err := q.ListPolicyAttachments(ctx, set.ID)
	if err != nil || len(attachments) != 1 || attachments[0].ID != attachment.ID {
		t.Fatalf("list policy attachments = %+v/%v", attachments, err)
	}

	parameter, err := q.CreatePolicySetParameter(ctx, CreatePolicySetParameterParams{
		ID: uuid.NewString(), PolicySetID: set.ID, ParameterKey: "max_length",
		ParameterValue: "8", Description: strPtr("maximum"),
	})
	if err != nil {
		t.Fatalf("create policy set parameter: %v", err)
	}
	if updated, err := q.UpdatePolicySetParameter(ctx, UpdatePolicySetParameterParams{
		ID: parameter.ID, ParameterKey: "max_length", ParameterValue: "12",
	}); err != nil || updated.ParameterValue != "12" {
		t.Fatalf("update policy set parameter = %+v/%v", updated, err)
	}

	expires := pgtype.Timestamptz{Time: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	exemption, err := q.CreatePolicyExemption(ctx, CreatePolicyExemptionParams{
		ID: uuid.NewString(), PolicySetID: set.ID, RuleID: "rule-1",
		TicketReference: "SEC-1", Justification: "test", ExpiresAt: expires,
		OrganizationID: org.ID, CreatedBy: strPtr("admin@example.com"),
	})
	if err != nil {
		t.Fatalf("create policy exemption: %v", err)
	}
	exemptions, err := q.ListPolicyExemptionsByOrganization(ctx, org.ID)
	if err != nil || len(exemptions) != 1 || exemptions[0].ID != exemption.ID {
		t.Fatalf("list policy exemptions = %+v/%v", exemptions, err)
	}
	if updated, err := q.UpdatePolicyExemption(ctx, UpdatePolicyExemptionParams{
		ID: exemption.ID, RuleID: "rule-1", TicketReference: "SEC-1",
		Justification: "after", ExpiresAt: expires,
	}); err != nil || updated.Justification != "after" {
		t.Fatalf("update policy exemption = %+v/%v", updated, err)
	}

	if updated, err := q.UpdatePolicySet(ctx, UpdatePolicySetParams{
		ID: set.ID, Name: "capture-policy", Description: strPtr("after"),
		EnforcementLevel: "ADVISORY", Global: false, Branch: strPtr("main"),
		Folder: strPtr("/"),
	}); err != nil || updated.EnforcementLevel != "ADVISORY" {
		t.Fatalf("update policy set = %+v/%v", updated, err)
	}

	if err := q.DeletePolicySet(ctx, set.ID); err != nil {
		t.Fatalf("delete policy set: %v", err)
	}
	if _, err := q.GetPolicySet(ctx, set.ID); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("get deleted policy set = %v, want no rows", err)
	}
	if attachments, err := q.ListPolicyAttachments(ctx, set.ID); err != nil || len(attachments) != 0 {
		t.Fatalf("attachments after delete = %+v/%v", attachments, err)
	}
	if parameters, err := q.ListPolicySetParameters(ctx, set.ID); err != nil || len(parameters) != 0 {
		t.Fatalf("parameters after delete = %+v/%v", parameters, err)
	}
	if exemptions, err := q.ListPolicyExemptionsByPolicySet(ctx, set.ID); err != nil || len(exemptions) != 0 {
		t.Fatalf("exemptions after delete = %+v/%v", exemptions, err)
	}
}
