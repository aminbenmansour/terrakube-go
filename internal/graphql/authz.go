package graphql

import (
	"context"
	"errors"
	"slices"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
)

// errPermissionDenied is answered as a GraphQL error; the presenter hides its
// text, as it does for every resolver failure.
var errPermissionDenied = errors.New("permission denied")

// isOwner mirrors the management API's owner rule: the development token owns
// a local stack, and a token carrying the admin group owns any deployment.
func (r *Resolver) isOwner(identity *auth.Identity) bool {
	if identity.Issuer == auth.DevIssuer {
		return true
	}
	group := r.adminGroup
	if group == "" {
		group = "TERRAKUBE_ADMIN"
	}
	return slices.Contains(identity.Groups, group)
}

// organizationIDs returns the organizations the caller may read. A nil result
// means every organization is readable: authentication is disabled, or the
// caller is the instance owner.
func (r *Resolver) organizationIDs(ctx context.Context) ([]string, error) {
	identity, ok := auth.IdentityFromContext(ctx)
	if !ok || r.isOwner(identity) {
		return nil, nil
	}
	if len(identity.Groups) == 0 {
		return []string{}, nil
	}

	rows, err := r.db.Query(ctx, `
SELECT DISTINCT organization_id FROM team WHERE name = ANY($1::text[])
UNION
SELECT DISTINCT w.organization_id FROM access a
JOIN workspace w ON w.id = a.workspace_id
WHERE a.name = ANY($1::text[])
UNION
SELECT DISTINCT p.organization_id FROM project_access pa
JOIN project p ON p.id = pa.project_id
WHERE pa.name = ANY($1::text[])`, identity.Groups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// organizationScope renders the predicate that limits a collection to the
// caller's organizations. A nil ids leaves the collection unscoped.
func organizationScope(column string, ids []string) (string, []any) {
	switch {
	case ids == nil:
		return "", nil
	case len(ids) == 0:
		return "FALSE", nil
	default:
		return column + " = ANY($1)", []any{ids}
	}
}

// canReadWorkspace backs the job-status subscription, which is opened with a
// workspace id rather than an organization.
func (r *Resolver) canReadWorkspace(ctx context.Context, workspaceID string) (bool, error) {
	ids, err := r.organizationIDs(ctx)
	if err != nil {
		return false, err
	}
	if ids == nil {
		return true, nil
	}
	if len(ids) == 0 {
		return false, nil
	}

	var organizationID string
	if err := r.db.QueryRow(ctx, "SELECT organization_id FROM workspace WHERE id = $1", workspaceID).Scan(&organizationID); err != nil {
		return false, err
	}
	return slices.Contains(ids, organizationID), nil
}

func (r *Resolver) canReadOrganization(ctx context.Context, organizationID string) (bool, error) {
	ids, err := r.organizationIDs(ctx)
	if err != nil {
		return false, err
	}
	if ids == nil {
		return true, nil
	}
	return slices.Contains(ids, organizationID), nil
}
