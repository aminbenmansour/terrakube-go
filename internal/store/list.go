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
		_, ok := s.Attrs[v.Attr]
		return ok
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
		not := ""
		if c.Op == rsql.OpOut {
			not = "NOT "
		}
		return attr.Expr + " " + not + "IN (" + strings.Join(placeholders, ", ") + ")", nil

	case rsql.OpEq:
		return s.match(attr, c.Values[0], args, false, false)
	case rsql.OpNe:
		return s.match(attr, c.Values[0], args, false, true)
	case rsql.OpIni:
		return s.match(attr, c.Values[0], args, true, false)
	case rsql.OpOuti:
		return s.match(attr, c.Values[0], args, true, true)
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
