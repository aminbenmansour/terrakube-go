package query

import "testing"

func TestParsePage(t *testing.T) {
	first := func(v int) *int { return &v }
	after := func(v string) *string { return &v }

	cases := []struct {
		name       string
		first      *int
		after      *string
		wantLimit  int64
		wantOffset int64
	}{
		{"defaults", nil, nil, DefaultPageSize, 0},
		{"explicit size", first(20), nil, 20, 0},
		{"offset", first(20), after("40"), 20, 40},
		{"caps the page size", first(20000), nil, MaxPageSize, 0},
		{"ignores a zero size", first(0), nil, DefaultPageSize, 0},
		{"ignores a negative size", first(-5), nil, DefaultPageSize, 0},
		{"empty offset", first(10), after(""), 10, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			limit, offset, err := ParsePage(tc.first, tc.after)
			if err != nil {
				t.Fatalf("ParsePage: %v", err)
			}
			if limit != tc.wantLimit || offset != tc.wantOffset {
				t.Fatalf("ParsePage = (%d, %d), want (%d, %d)", limit, offset, tc.wantLimit, tc.wantOffset)
			}
		})
	}
}

func TestParsePageRejectsBadOffset(t *testing.T) {
	for _, value := range []string{"x", "-1", "1.5"} {
		if _, _, err := ParsePage(nil, &value); err == nil {
			t.Fatalf("ParsePage accepted offset %q", value)
		}
	}
}
