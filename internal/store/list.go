package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/rsql"
)

// Kind is how an attribute's values are read from an RSQL expression.
type Kind int

const (
	KindText Kind = iota
	KindTime
	KindBool
)

// Attr maps a JSON:API attribute name onto a SQL expression.
type Attr struct {
	Expr string
	Kind Kind
	// From and Where describe an attribute on a related table: the predicate
	// becomes an EXISTS over that table instead of a column on the list table.
	From  string
	Where string
	// Collection marks the relationship itself rather than one of its columns;
	// only =isempty= applies to it.
	Collection bool
}

// allows reports whether the attribute supports the operator.
func (a Attr) allows(op rsql.Op) bool {
	switch op {
	case rsql.OpIsEmpty:
		return a.Collection
	case rsql.OpHasMember, rsql.OpHasNoMember:
		return a.From != "" && !a.Collection
	}
	if a.Collection {
		return false
	}
	if a.From == "" {
		return true
	}
	// Related columns support membership and ranges; =isnull= would need a
	// quantifier choice the UI does not send.
	switch op {
	case rsql.OpEq, rsql.OpNe, rsql.OpIn, rsql.OpOut, rsql.OpIni, rsql.OpOuti,
		rsql.OpBetween, rsql.OpNotBetween:
		return true
	}
	return false
}

// ListSpec is one resource's list query: its table, the scope predicate that
// every request carries, the attributes a filter may name, and the sortable
// attributes.
type ListSpec struct {
	Table string
	// Base is the scope predicate. It may use $1..$n, filled from the scope
	// arguments the caller passes to ListPage.
	Base  string
	Attrs map[string]Attr
	Sorts map[string]string
}

// SortKey is one requested sort attribute.
type SortKey struct {
	Column string
	Desc   bool
}

// ListPage runs the filtered, sorted and paged query and scans the rows by
// position. sqlc emits model fields in table-column order, which ALTER TABLE
// preserves by appending, so SELECT * lines up with the struct; a test guards
// it.
func (q *Queries) ListPage[T any](ctx context.Context, spec ListSpec, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]T, error) {
	where, args, err := spec.where(scope, filter)
	if err != nil {
		return nil, err
	}

	var b strings.Builder
	b.WriteString("SELECT * FROM ")
	b.WriteString(spec.Table)
	if where != "" {
		b.WriteString(" WHERE ")
		b.WriteString(where)
	}
	b.WriteString(" ORDER BY ")
	b.WriteString(spec.orderBy(sorts))

	args = append(args, limit, offset)
	fmt.Fprintf(&b, " LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := q.db.Query(ctx, b.String(), args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[T])
}

// CountList counts the rows the same filter selects.
func (q *Queries) CountList(ctx context.Context, spec ListSpec, scope []any, filter rsql.Node) (int64, error) {
	where, args, err := spec.where(scope, filter)
	if err != nil {
		return 0, err
	}

	query := "SELECT count(*) FROM " + spec.Table
	if where != "" {
		query += " WHERE " + where
	}

	var total int64
	err = q.db.QueryRow(ctx, query, args...).Scan(&total)
	return total, err
}

func (s ListSpec) where(scope []any, filter rsql.Node) (string, []any, error) {
	args := append([]any{}, scope...)
	clauses := make([]string, 0, 2)
	if s.Base != "" {
		clauses = append(clauses, s.Base)
	}

	// A filter naming an attribute we do not model is dropped whole: applying
	// the rest would narrow the result where upstream widens it.
	if filter != nil && s.supported(filter) {
		clause, err := s.render(filter, &args)
		if err != nil {
			return "", nil, err
		}
		clauses = append(clauses, "("+clause+")")
	}
	return strings.Join(clauses, " AND "), args, nil
}

func (s ListSpec) supported(n rsql.Node) bool {
	switch v := n.(type) {
	case rsql.And:
		return s.supported(v.Left) && s.supported(v.Right)
	case rsql.Or:
		return s.supported(v.Left) && s.supported(v.Right)
	case rsql.Compare:
		attr, ok := s.Attrs[v.Attr]
		return ok && attr.allows(v.Op)
	}
	return false
}

func (s ListSpec) render(n rsql.Node, args *[]any) (string, error) {
	switch v := n.(type) {
	case rsql.And:
		left, err := s.render(v.Left, args)
		if err != nil {
			return "", err
		}
		right, err := s.render(v.Right, args)
		if err != nil {
			return "", err
		}
		return "(" + left + " AND " + right + ")", nil
	case rsql.Or:
		left, err := s.render(v.Left, args)
		if err != nil {
			return "", err
		}
		right, err := s.render(v.Right, args)
		if err != nil {
			return "", err
		}
		return "(" + left + " OR " + right + ")", nil
	case rsql.Compare:
		return s.renderCompare(v, args)
	}
	return "", fmt.Errorf("%w: unknown expression", rsql.ErrSyntax)
}

func (s ListSpec) renderCompare(c rsql.Compare, args *[]any) (string, error) {
	attr := s.Attrs[c.Attr]
	if len(c.Values) == 0 {
		return "", fmt.Errorf("%w: %s needs a value", rsql.ErrSyntax, c.Op)
	}

	switch c.Op {
	case rsql.OpIsNull:
		switch strings.ToLower(c.Values[0]) {
		case "true":
			return attr.Expr + " IS NULL", nil
		case "false":
			return attr.Expr + " IS NOT NULL", nil
		}
		return "", fmt.Errorf("%w: =isnull= takes true or false", rsql.ErrSyntax)

	case rsql.OpIsEmpty:
		clause := exists(attr, "")
		if strings.EqualFold(c.Values[0], "true") {
			return "NOT " + clause, nil
		}
		if strings.EqualFold(c.Values[0], "false") {
			return clause, nil
		}
		return "", fmt.Errorf("%w: =isempty= takes true or false", rsql.ErrSyntax)

	case rsql.OpHasMember, rsql.OpHasNoMember:
		predicate, err := s.match(attr, c.Values[0], args, false, false)
		if err != nil {
			return "", err
		}
		return relate(attr, predicate, c.Op == rsql.OpHasNoMember), nil

	case rsql.OpIn, rsql.OpOut:
		placeholders := make([]string, 0, len(c.Values))
		for _, raw := range c.Values {
			value, err := s.value(attr, raw)
			if err != nil {
				return "", err
			}
			*args = append(*args, value)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(*args)))
		}
		predicate := attr.Expr + " IN (" + strings.Join(placeholders, ", ") + ")"
		if attr.From != "" {
			return relate(attr, predicate, c.Op == rsql.OpOut), nil
		}
		not := ""
		if c.Op == rsql.OpOut {
			not = "NOT "
		}
		return attr.Expr + " " + not + "IN (" + strings.Join(placeholders, ", ") + ")", nil

	case rsql.OpEq, rsql.OpNe, rsql.OpIni, rsql.OpOuti:
		insensitive := c.Op == rsql.OpIni || c.Op == rsql.OpOuti
		negate := c.Op == rsql.OpNe || c.Op == rsql.OpOuti
		// On a related attribute the negation belongs outside EXISTS: "has no
		// related row with this value" also matches rows with no relation.
		predicate, err := s.match(attr, c.Values[0], args, insensitive, negate && attr.From == "")
		if err != nil {
			return "", err
		}
		return relate(attr, predicate, negate && attr.From != ""), nil
	}

	if c.Op == rsql.OpBetween || c.Op == rsql.OpNotBetween {
		if len(c.Values) != 2 {
			return "", fmt.Errorf("%w: %s takes two values", rsql.ErrSyntax, c.Op)
		}
		low, err := s.value(attr, c.Values[0])
		if err != nil {
			return "", err
		}
		high, err := s.value(attr, c.Values[1])
		if err != nil {
			return "", err
		}
		*args = append(*args, low)
		lo := len(*args)
		*args = append(*args, high)
		hi := len(*args)
		predicate := fmt.Sprintf("%s >= $%d AND %s <= $%d", attr.Expr, lo, attr.Expr, hi)
		if c.Op == rsql.OpNotBetween {
			predicate = fmt.Sprintf("(%s < $%d OR %s > $%d)", attr.Expr, lo, attr.Expr, hi)
		}
		return relate(attr, predicate, false), nil
	}

	op, ok := map[rsql.Op]string{
		rsql.OpLt: "<", rsql.OpLe: "<=", rsql.OpGt: ">", rsql.OpGe: ">=",
	}[c.Op]
	if !ok {
		return "", fmt.Errorf("%w: unsupported operator %s", rsql.ErrSyntax, c.Op)
	}
	value, err := s.value(attr, c.Values[0])
	if err != nil {
		return "", err
	}
	*args = append(*args, value)
	return fmt.Sprintf("%s %s $%d", attr.Expr, op, len(*args)), nil
}

// exists renders an EXISTS over the attribute's related table. An empty
// predicate only tests that a related row exists, which =isempty= uses.
func exists(attr Attr, predicate string) string {
	clause := "EXISTS (SELECT 1 FROM " + attr.From + " WHERE " + attr.Where
	if predicate != "" {
		clause += " AND " + predicate
	}
	return clause + ")"
}

// relate wraps a predicate on a related attribute in EXISTS. A negated
// predicate becomes NOT EXISTS, which keeps rows without any related row.
func relate(attr Attr, predicate string, negate bool) string {
	if attr.From == "" {
		return predicate
	}
	clause := exists(attr, predicate)
	if negate {
		return "NOT " + clause
	}
	return clause
}

// match renders the string operators. A '*' in the value is a wildcard; every
// other '%', '_' and '\' is a literal.
func (s ListSpec) match(attr Attr, raw string, args *[]any, insensitive, negate bool) (string, error) {
	if insensitive && attr.Kind != KindText {
		insensitive = false
	}

	expr := attr.Expr
	if insensitive {
		expr = "lower(" + attr.Expr + ")"
	}

	if strings.Contains(raw, "*") {
		pattern := raw
		if insensitive {
			pattern = strings.ToLower(pattern)
		}
		*args = append(*args, likePattern(pattern))
		op := "LIKE"
		if negate {
			op = "NOT LIKE"
		}
		return fmt.Sprintf("%s %s $%d ESCAPE '\\'", expr, op, len(*args)), nil
	}

	value, err := s.value(attr, raw)
	if err != nil {
		return "", err
	}
	if insensitive {
		value = strings.ToLower(value.(string))
	}
	*args = append(*args, value)
	op := "="
	if negate {
		op = "<>"
	}
	return fmt.Sprintf("%s %s $%d", expr, op, len(*args)), nil
}

func (s ListSpec) value(attr Attr, raw string) (any, error) {
	switch attr.Kind {
	case KindTime:
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return nil, fmt.Errorf("%w: %q is not a timestamp", rsql.ErrSyntax, raw)
		}
		return t, nil
	case KindBool:
		switch strings.ToLower(raw) {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}
		return nil, fmt.Errorf("%w: %q is not a boolean", rsql.ErrSyntax, raw)
	}
	return raw, nil
}

func likePattern(raw string) string {
	var b strings.Builder
	for _, r := range raw {
		switch r {
		case '\\', '%', '_':
			b.WriteByte('\\')
			b.WriteRune(r)
		case '*':
			b.WriteByte('%')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (s ListSpec) orderBy(sorts []SortKey) string {
	parts := make([]string, 0, len(sorts)+2)
	desc, hasID := false, false
	for _, key := range sorts {
		expr, ok := s.Sorts[key.Column]
		if !ok {
			continue
		}
		dir := "ASC"
		if key.Desc {
			dir = "DESC"
		}
		parts = append(parts, expr+" "+dir)
		desc, hasID = key.Desc, hasID || key.Column == "id"
	}

	// The list's default order is its name, then id as the tiebreaker that
	// makes paging a total order.
	if len(parts) == 0 {
		if expr, ok := s.Sorts["name"]; ok {
			parts = append(parts, expr+" ASC")
		}
	}
	if !hasID {
		dir := "ASC"
		if desc {
			dir = "DESC"
		}
		parts = append(parts, "id "+dir)
	}
	return strings.Join(parts, ", ")
}

// The list handlers reach the store through an interface, and Go does not allow
// a generic method there, so each list gets a concrete wrapper.

// Lister is the list surface the management API depends on. *Queries has these
// methods; tests can fake them.
type Lister interface {
	ListOrganizations(ctx context.Context, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]Organization, error)
	CountOrganizations(ctx context.Context, scope []any, filter rsql.Node) (int64, error)
	ListProjects(ctx context.Context, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]Project, error)
	CountProjects(ctx context.Context, scope []any, filter rsql.Node) (int64, error)
	ListTags(ctx context.Context, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]Tag, error)
	CountTags(ctx context.Context, scope []any, filter rsql.Node) (int64, error)
	ListWorkspaces(ctx context.Context, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]Workspace, error)
	CountWorkspaces(ctx context.Context, scope []any, filter rsql.Node) (int64, error)
	ListProjectWorkspaces(ctx context.Context, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]Workspace, error)
	CountProjectWorkspaces(ctx context.Context, scope []any, filter rsql.Node) (int64, error)
}

func (q *Queries) ListOrganizations(ctx context.Context, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]Organization, error) {
	return q.ListPage[Organization](ctx, OrganizationList, scope, filter, sorts, limit, offset)
}

func (q *Queries) CountOrganizations(ctx context.Context, scope []any, filter rsql.Node) (int64, error) {
	return q.CountList(ctx, OrganizationList, scope, filter)
}

func (q *Queries) ListProjects(ctx context.Context, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]Project, error) {
	return q.ListPage[Project](ctx, ProjectList, scope, filter, sorts, limit, offset)
}

func (q *Queries) CountProjects(ctx context.Context, scope []any, filter rsql.Node) (int64, error) {
	return q.CountList(ctx, ProjectList, scope, filter)
}

func (q *Queries) ListTags(ctx context.Context, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]Tag, error) {
	return q.ListPage[Tag](ctx, TagList, scope, filter, sorts, limit, offset)
}

func (q *Queries) CountTags(ctx context.Context, scope []any, filter rsql.Node) (int64, error) {
	return q.CountList(ctx, TagList, scope, filter)
}

func (q *Queries) ListWorkspaces(ctx context.Context, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]Workspace, error) {
	return q.ListPage[Workspace](ctx, WorkspaceList, scope, filter, sorts, limit, offset)
}

func (q *Queries) CountWorkspaces(ctx context.Context, scope []any, filter rsql.Node) (int64, error) {
	return q.CountList(ctx, WorkspaceList, scope, filter)
}

func (q *Queries) ListProjectWorkspaces(ctx context.Context, scope []any, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]Workspace, error) {
	return q.ListPage[Workspace](ctx, ProjectWorkspaceList, scope, filter, sorts, limit, offset)
}

func (q *Queries) CountProjectWorkspaces(ctx context.Context, scope []any, filter rsql.Node) (int64, error) {
	return q.CountList(ctx, ProjectWorkspaceList, scope, filter)
}
