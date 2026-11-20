package query

import (
	"fmt"
	"strings"
)

// Order is one sort key resolved to a column.
type Order struct {
	Column string
	Desc   bool
}

// ParseSort parses Elide sort syntax ("-lastJobDate,+name") against a field
// whitelist. An empty input yields no ordering.
func ParseSort(input string, fields map[string]Field) ([]Order, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}

	var orders []Order
	for _, part := range strings.Split(input, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		desc := false
		switch part[0] {
		case '-':
			desc = true
			part = part[1:]
		case '+':
			part = part[1:]
		}
		field, ok := fields[part]
		if !ok {
			return nil, fmt.Errorf("query: unknown sort field %q", part)
		}
		if field.Column == "" {
			return nil, fmt.Errorf("query: sort field %q is not sortable", part)
		}
		orders = append(orders, Order{Column: field.Column, Desc: desc})
	}
	return orders, nil
}

// OrderBy renders an ORDER BY clause, or an empty string when there is nothing
// to order by.
func OrderBy(orders []Order) string {
	if len(orders) == 0 {
		return ""
	}
	parts := make([]string, 0, len(orders))
	for _, order := range orders {
		direction := "ASC"
		if order.Desc {
			direction = "DESC"
		}
		parts = append(parts, order.Column+" "+direction)
	}
	return "ORDER BY " + strings.Join(parts, ", ")
}
