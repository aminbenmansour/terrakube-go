package query

import (
	"reflect"
	"testing"
)

func TestParseSort(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Order
	}{
		{"ascending by default", "name", []Order{{Column: "name"}}},
		{"explicit ascending", "+name", []Order{{Column: "name"}}},
		{"descending", "-lastJobDate", []Order{{Column: "last_job_date", Desc: true}}},
		{"compound with tie breaker", "-lastJobDate,+name", []Order{
			{Column: "last_job_date", Desc: true},
			{Column: "name"},
		}},
		{"relationship path", "project.id", []Order{{Column: "project_id"}}},
		{"empty", "", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseSort(tc.input, testFields)
			if err != nil {
				t.Fatalf("ParseSort(%q): %v", tc.input, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("ParseSort(%q) = %#v, want %#v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseSortRejectsUnknownField(t *testing.T) {
	if got, err := ParseSort("password", testFields); err == nil {
		t.Fatalf("ParseSort accepted an unknown field: %#v", got)
	}
}

func TestOrderByRendersSQL(t *testing.T) {
	orders := []Order{{Column: "last_job_date", Desc: true}, {Column: "name"}}
	if got, want := OrderBy(orders), "ORDER BY last_job_date DESC, name ASC"; got != want {
		t.Fatalf("OrderBy = %q, want %q", got, want)
	}
	if got := OrderBy(nil); got != "" {
		t.Fatalf("OrderBy(nil) = %q, want empty", got)
	}
}
