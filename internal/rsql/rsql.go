// Package rsql parses the RSQL dialect Elide accepts, which is what the
// upstream UI sends in filter[type] parameters.
package rsql

import (
	"errors"
	"fmt"
	"strings"
)

// ErrSyntax marks a filter the parser rejects. Callers turn it into a 400.
var ErrSyntax = errors.New("invalid filter expression")

// Op is a comparison operator. The values are the RSQL spellings.
type Op string

const (
	OpEq          Op = "=="
	OpNe          Op = "!="
	OpIn          Op = "=in="
	OpOut         Op = "=out="
	OpIsNull      Op = "=isnull="
	OpIni         Op = "=ini="
	OpOuti        Op = "=outi="
	OpLt          Op = "=lt="
	OpLe          Op = "=le="
	OpGt          Op = "=gt="
	OpGe          Op = "=ge="
	OpBetween     Op = "=between="
	OpNotBetween  Op = "=notbetween="
	OpHasMember   Op = "=hasmember="
	OpHasNoMember Op = "=hasnomember="
	OpIsEmpty     Op = "=isempty="
)

// Node is a filter expression.
type Node interface{ node() }

// And is the conjunction written with ';'.
type And struct{ Left, Right Node }

// Or is the disjunction written with ','.
type Or struct{ Left, Right Node }

// Compare is one attribute predicate. Values is a single value except for the
// set operators (=in=, =out=) and =between=, which hold a list.
type Compare struct {
	Attr   string
	Op     Op
	Values []string
}

func (And) node()     {}
func (Or) node()      {}
func (Compare) node() {}

// Parse reads a filter expression. An empty string is an error; callers skip
// the filter entirely when the parameter is absent.
func Parse(raw string) (Node, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("%w: empty expression", ErrSyntax)
	}

	p := &parser{s: raw}
	n, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	p.skipSpace()
	if p.pos != len(p.s) {
		return nil, p.errorf("unexpected %q", p.s[p.pos:])
	}
	return n, nil
}

type parser struct {
	s   string
	pos int
}

// operators are matched longest first, because "=in=" is a prefix of "=ini="
// and "=out=" of "=outi=".
var operators = []Op{
	OpHasNoMember, OpHasMember, OpNotBetween, OpBetween, OpIsEmpty,
	OpIsNull, OpIni, OpOuti, OpOut, OpIn, OpLt, OpLe, OpGt, OpGe, OpEq, OpNe,
}

func (p *parser) parseOr() (Node, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for {
		p.skipSpace()
		if !p.consume(",") {
			return left, nil
		}
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = Or{Left: left, Right: right}
	}
}

func (p *parser) parseAnd() (Node, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for {
		p.skipSpace()
		if !p.consume(";") {
			return left, nil
		}
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		left = And{Left: left, Right: right}
	}
}

func (p *parser) parsePrimary() (Node, error) {
	p.skipSpace()
	if p.consume("(") {
		n, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		p.skipSpace()
		if !p.consume(")") {
			return nil, p.errorf("missing )")
		}
		return n, nil
	}
	return p.parseCompare()
}

func (p *parser) parseCompare() (Node, error) {
	p.skipSpace()
	start := p.pos
	for p.pos < len(p.s) && !strings.ContainsRune("=!<>", rune(p.s[p.pos])) {
		p.pos++
	}
	attr := strings.TrimSpace(p.s[start:p.pos])
	if attr == "" {
		return nil, p.errorf("missing attribute")
	}

	var op Op
	for _, candidate := range operators {
		if strings.HasPrefix(p.s[p.pos:], string(candidate)) {
			op = candidate
			p.pos += len(candidate)
			break
		}
	}
	if op == "" {
		return nil, p.errorf("missing operator after %q", attr)
	}

	switch op {
	case OpIn, OpOut:
		values, err := p.parseList()
		if err != nil {
			return nil, err
		}
		return Compare{Attr: attr, Op: op, Values: values}, nil

	case OpBetween, OpNotBetween:
		values, err := p.parseList()
		if err != nil {
			return nil, err
		}
		if len(values) != 2 {
			return nil, p.errorf("%s takes two values", op)
		}
		return Compare{Attr: attr, Op: op, Values: values}, nil
	}

	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	return Compare{Attr: attr, Op: op, Values: []string{value}}, nil
}

// parseList reads the parenthesised values of =in= and =out=.
func (p *parser) parseList() ([]string, error) {
	p.skipSpace()
	if !p.consume("(") {
		return nil, p.errorf("expected ( after operator")
	}

	var values []string
	for {
		p.skipSpace()
		if p.consume(")") {
			if len(values) == 0 {
				return nil, p.errorf("empty value list")
			}
			return values, nil
		}
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		values = append(values, value)

		p.skipSpace()
		if p.consume(",") {
			continue
		}
		if p.consume(")") {
			return values, nil
		}
		return nil, p.errorf("expected , or ) in value list")
	}
}

// parseValue reads one value: quoted (single or double, with backslash
// escapes) or bare up to the next delimiter.
func (p *parser) parseValue() (string, error) {
	p.skipSpace()
	if p.pos >= len(p.s) {
		return "", p.errorf("missing value")
	}

	if quote := p.s[p.pos]; quote == '\'' || quote == '"' {
		p.pos++
		var b strings.Builder
		for p.pos < len(p.s) {
			c := p.s[p.pos]
			switch {
			case c == '\\' && p.pos+1 < len(p.s):
				b.WriteByte(p.s[p.pos+1])
				p.pos += 2
			case c == quote:
				p.pos++
				return b.String(), nil
			default:
				b.WriteByte(c)
				p.pos++
			}
		}
		return "", p.errorf("unterminated quoted value")
	}

	start := p.pos
	for p.pos < len(p.s) && !strings.ContainsRune(";,) ", rune(p.s[p.pos])) {
		p.pos++
	}
	value := p.s[start:p.pos]
	if value == "" {
		return "", p.errorf("missing value")
	}
	return value, nil
}

func (p *parser) skipSpace() {
	for p.pos < len(p.s) && p.s[p.pos] == ' ' {
		p.pos++
	}
}

func (p *parser) consume(s string) bool {
	if strings.HasPrefix(p.s[p.pos:], s) {
		p.pos += len(s)
		return true
	}
	return false
}

func (p *parser) errorf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrSyntax, fmt.Sprintf(format, args...))
}
