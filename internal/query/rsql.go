// Package query turns the client query strings the upstream UI sends (Elide
// RSQL filters, sort lists and offset paging) into parameterised SQL.
package query

import (
	"fmt"
	"strings"
)

// Op is an RSQL comparison operator.
type Op string

// The operators Elide documents for the Terrakube UI. Anything else is
// rejected rather than guessed at.
const (
	OpEq      Op = "=="
	OpNotEq   Op = "!="
	OpGt      Op = "=gt="
	OpGe      Op = "=ge="
	OpLt      Op = "=lt="
	OpLe      Op = "=le="
	OpIn      Op = "=in="
	OpOut     Op = "=out="
	OpIsNull  Op = "=isnull="
	OpNotNull Op = "=notnull="
	OpIni     Op = "=ini="
)

// Node is a parsed filter expression.
type Node interface{ isNode() }

// Compare is one field comparison.
type Compare struct {
	Field  string
	Op     Op
	Values []string
}

// And requires every child to match.
type And struct{ Nodes []Node }

// Or requires one child to match.
type Or struct{ Nodes []Node }

func (*Compare) isNode() {}
func (*And) isNode()     {}
func (*Or) isNode()      {}

// operators ordered longest first so a prefix never wins by accident.
var operators = []Op{OpNotNull, OpIsNull, OpIni, OpOut, OpIn, OpGt, OpGe, OpLt, OpLe, OpEq, OpNotEq}

// ParseRSQL parses an Elide RSQL expression. An empty expression yields a nil
// Node, which callers treat as "no filter".
func ParseRSQL(input string) (Node, error) {
	p := &parser{input: input}
	p.skipSpace()
	if p.eof() {
		return nil, nil
	}
	node, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	p.skipSpace()
	if !p.eof() {
		return nil, fmt.Errorf("rsql: unexpected %q at offset %d", p.rest(), p.pos)
	}
	return node, nil
}

type parser struct {
	input string
	pos   int
}

func (p *parser) eof() bool { return p.pos >= len(p.input) }

func (p *parser) rest() string { return p.input[p.pos:] }

func (p *parser) skipSpace() {
	for !p.eof() && (p.input[p.pos] == ' ' || p.input[p.pos] == '\t' || p.input[p.pos] == '\n' || p.input[p.pos] == '\r') {
		p.pos++
	}
}

func (p *parser) peek() byte {
	if p.eof() {
		return 0
	}
	return p.input[p.pos]
}

func (p *parser) parseOr() (Node, error) {
	first, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	nodes := []Node{first}
	for {
		p.skipSpace()
		if p.peek() != ',' {
			break
		}
		p.pos++
		next, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, next)
	}
	if len(nodes) == 1 {
		return nodes[0], nil
	}
	return &Or{Nodes: nodes}, nil
}

func (p *parser) parseAnd() (Node, error) {
	first, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	nodes := []Node{first}
	for {
		p.skipSpace()
		if p.peek() != ';' {
			break
		}
		p.pos++
		next, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, next)
	}
	if len(nodes) == 1 {
		return nodes[0], nil
	}
	return &And{Nodes: nodes}, nil
}

func (p *parser) parsePrimary() (Node, error) {
	p.skipSpace()
	if p.eof() {
		return nil, fmt.Errorf("rsql: unexpected end of expression")
	}
	if p.peek() == '(' {
		p.pos++
		node, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		p.skipSpace()
		if p.peek() != ')' {
			return nil, fmt.Errorf("rsql: missing ')' at offset %d", p.pos)
		}
		p.pos++
		return node, nil
	}
	return p.parseComparison()
}

func (p *parser) parseComparison() (Node, error) {
	field := p.readField()
	if field == "" {
		return nil, fmt.Errorf("rsql: expected a field name at offset %d", p.pos)
	}

	op, err := p.readOperator()
	if err != nil {
		return nil, err
	}

	values, err := p.readValues(op)
	if err != nil {
		return nil, err
	}
	return &Compare{Field: field, Op: op, Values: values}, nil
}

func (p *parser) readField() string {
	start := p.pos
	for !p.eof() {
		c := p.input[p.pos]
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '.' {
			p.pos++
			continue
		}
		break
	}
	return p.input[start:p.pos]
}

func (p *parser) readOperator() (Op, error) {
	for _, op := range operators {
		if strings.HasPrefix(p.rest(), string(op)) {
			p.pos += len(op)
			return op, nil
		}
	}
	if p.eof() {
		return "", fmt.Errorf("rsql: missing operator at offset %d", p.pos)
	}
	return "", fmt.Errorf("rsql: unknown operator %q at offset %d", p.rest(), p.pos)
}

func (p *parser) readValues(op Op) ([]string, error) {
	if op == OpIn || op == OpOut {
		p.skipSpace()
		if p.peek() != '(' {
			return nil, fmt.Errorf("rsql: %s needs a list at offset %d", op, p.pos)
		}
		p.pos++
		var values []string
		for {
			p.skipSpace()
			if p.eof() {
				return nil, fmt.Errorf("rsql: missing ')' at offset %d", p.pos)
			}
			if p.peek() == ')' {
				if len(values) == 0 {
					return nil, fmt.Errorf("rsql: empty list at offset %d", p.pos)
				}
				p.pos++
				return values, nil
			}
			value, err := p.readValue()
			if err != nil {
				return nil, err
			}
			values = append(values, value)
			p.skipSpace()
			switch p.peek() {
			case ',':
				p.pos++
			case ')':
			default:
				return nil, fmt.Errorf("rsql: expected ',' or ')' at offset %d", p.pos)
			}
		}
	}

	value, err := p.readValue()
	if err != nil {
		return nil, err
	}
	return []string{value}, nil
}

func (p *parser) readValue() (string, error) {
	p.skipSpace()
	if p.eof() {
		return "", fmt.Errorf("rsql: missing value at offset %d", p.pos)
	}

	if quote := p.peek(); quote == '\'' || quote == '"' {
		p.pos++
		var b strings.Builder
		for {
			if p.eof() {
				return "", fmt.Errorf("rsql: unterminated value at offset %d", p.pos)
			}
			c := p.input[p.pos]
			if c == quote {
				// A doubled quote is an escaped quote inside the value.
				if p.pos+1 < len(p.input) && p.input[p.pos+1] == quote {
					b.WriteByte(quote)
					p.pos += 2
					continue
				}
				p.pos++
				return b.String(), nil
			}
			b.WriteByte(c)
			p.pos++
		}
	}

	start := p.pos
	for !p.eof() {
		switch p.input[p.pos] {
		case ',', ';', ')', '(', ' ', '\t', '\n', '\r':
			value := p.input[start:p.pos]
			if value == "" {
				return "", fmt.Errorf("rsql: missing value at offset %d", p.pos)
			}
			return value, nil
		}
		p.pos++
	}
	if start == p.pos {
		return "", fmt.Errorf("rsql: missing value at offset %d", p.pos)
	}
	return p.input[start:p.pos], nil
}
