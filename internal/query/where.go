package query

import (
	"fmt"
	"strconv"
	"strings"
)

// Where renders a filter as a parameterised SQL predicate. fields maps an RSQL
// dot path to a column expression; a path that is not in the map is rejected,
// which is what keeps a filter string out of the SQL text. argOffset is the
// number of placeholders the caller already used.
func Where(node Node, fields map[string]string, argOffset int) (string, []any, error) {
	if node == nil {
		return "", nil, nil
	}
	b := &builder{fields: fields, args: argOffset}
	sql, err := b.render(node)
	if err != nil {
		return "", nil, err
	}
	return sql, b.values, nil
}

type builder struct {
	fields map[string]string
	args   int
	values []any
}

func (b *builder) placeholder(value any) string {
	b.args++
	b.values = append(b.values, value)
	return "$" + strconv.Itoa(b.args)
}

func (b *builder) render(node Node) (string, error) {
	switch n := node.(type) {
	case *Compare:
		return b.compare(n)
	case *And:
		return b.join(n.Nodes, "AND")
	case *Or:
		return b.join(n.Nodes, "OR")
	default:
		return "", fmt.Errorf("query: unsupported filter node %T", node)
	}
}

func (b *builder) join(nodes []Node, operator string) (string, error) {
	if len(nodes) == 0 {
		return "", fmt.Errorf("query: empty %s group", operator)
	}
	parts := make([]string, 0, len(nodes))
	for _, node := range nodes {
		sql, err := b.render(node)
		if err != nil {
			return "", err
		}
		parts = append(parts, sql)
	}
	return "(" + strings.Join(parts, " "+operator+" ") + ")", nil
}

func (b *builder) compare(c *Compare) (string, error) {
	column, ok := b.fields[c.Field]
	if !ok {
		return "", fmt.Errorf("query: unknown filter field %q", c.Field)
	}

	switch c.Op {
	case OpIsNull, OpNotNull:
		wantNull := c.Op == OpIsNull
		if c.Op == OpIsNull {
			if len(c.Values) != 1 {
				return "", fmt.Errorf("query: =isnull= needs a value")
			}
			switch strings.ToLower(strings.TrimSpace(c.Values[0])) {
			case "true":
				wantNull = true
			case "false":
				wantNull = false
			default:
				return "", fmt.Errorf("query: =isnull= expects true or false, got %q", c.Values[0])
			}
		}
		if wantNull {
			return column + " IS NULL", nil
		}
		return column + " IS NOT NULL", nil

	case OpIn, OpOut:
		if len(c.Values) == 0 {
			return "", fmt.Errorf("query: %s needs at least one value", c.Op)
		}
		placeholders := make([]string, 0, len(c.Values))
		for _, value := range c.Values {
			placeholders = append(placeholders, b.placeholder(value))
		}
		operator := "IN"
		if c.Op == OpOut {
			operator = "NOT IN"
		}
		return column + " " + operator + " (" + strings.Join(placeholders, ", ") + ")", nil

	case OpIni:
		if len(c.Values) != 1 {
			return "", fmt.Errorf("query: =ini= needs exactly one value")
		}
		return "lower(" + column + ") LIKE " + b.placeholder(likePattern(c.Values[0])) + ` ESCAPE '\'`, nil

	case OpEq, OpNotEq, OpGt, OpGe, OpLt, OpLe:
		if len(c.Values) != 1 {
			return "", fmt.Errorf("query: %s needs exactly one value", c.Op)
		}
		operator := map[Op]string{
			OpEq:    "=",
			OpNotEq: "<>",
			OpGt:    ">",
			OpGe:    ">=",
			OpLt:    "<",
			OpLe:    "<=",
		}[c.Op]
		return column + " " + operator + " " + b.placeholder(c.Values[0]), nil

	default:
		return "", fmt.Errorf("query: unsupported operator %q", c.Op)
	}
}

// likePattern turns an =ini= value into a LIKE pattern: * is a wildcard, and
// the characters that are special to LIKE are escaped so a literal % or _ in a
// workspace name still matches itself.
func likePattern(value string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(value) {
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
