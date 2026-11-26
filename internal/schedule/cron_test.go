package schedule

import (
	"testing"
	"time"
)

func mustParse(t *testing.T, expression string) *Cron {
	t.Helper()
	cron, err := Parse(expression)
	if err != nil {
		t.Fatalf("Parse(%q): %v", expression, err)
	}
	return cron
}

func TestParseMatchesQuartzFields(t *testing.T) {
	// Monday 2026-11-30 00:00:00 UTC.
	monday := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		expression string
		at         time.Time
		want       bool
	}{
		{"0 * * * * ?", monday.Add(9*time.Hour + 12*time.Minute), true},
		{"0 * * * * ?", monday.Add(9*time.Hour + 12*time.Minute + time.Second), false},
		{"0 0/15 * * * ?", monday.Add(9*time.Hour + 15*time.Minute), true},
		{"0 0/15 * * * ?", monday.Add(9*time.Hour + 16*time.Minute), false},
		{"0 10-20/5 * * * ?", monday.Add(9*time.Hour + 20*time.Minute), true},
		{"0 10-20/5 * * * ?", monday.Add(9*time.Hour + 21*time.Minute), false},
		{"30 5 4 * * ?", monday.Add(4*time.Hour + 5*time.Minute + 30*time.Second), true},
		{"30 5 4 * * ?", monday.Add(4*time.Hour + 5*time.Minute + 31*time.Second), false},
		// Quartz day-of-week: 1 = Sunday, 2 = Monday.
		{"0 0 9 ? * 2", monday.Add(9 * time.Hour), true},
		{"0 0 9 ? * 2", monday.Add(9*time.Hour + 24*time.Hour), false},
		{"0 0 9 ? * 1", monday.Add(9 * time.Hour), false},
		// Day of month constrained, weekday unspecified.
		{"0 0 9 30 * ?", monday.Add(9 * time.Hour), true},
		{"0 0 9 1 * ?", monday.Add(9 * time.Hour), false},
		// Optional year.
		{"0 0 0 1 1 ? 2027", time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC), true},
		{"0 0 0 1 1 ? 2027", time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC), false},
		{"0 0 0 1 1 ? *", time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC), true},
	}
	for _, tc := range cases {
		cron := mustParse(t, tc.expression)
		if got := cron.Matches(tc.at); got != tc.want {
			t.Errorf("%q at %s = %v, want %v", tc.expression, tc.at, got, tc.want)
		}
	}
}

func TestParseRejectsWhatQuartzRejects(t *testing.T) {
	cases := []string{
		"0 * * * *",     // five fields
		"0 0 0 * * *",   // day-of-month and day-of-week both specified
		"0 0 0 ? * ?",   // both unspecified
		"0 0 0 L * ?",   // calendar extension
		"0 0 0 1W * ?",  // nearest weekday
		"0 0 0 ? * 2#1", // nth weekday
		"60 * * * * ?",  // seconds out of range
		"0 0 24 * * ?",  // hours out of range
		"0 0 0 0 * ?",   // day-of-month below 1
		"0 0 0 ? * 0",   // Quartz weekdays start at 1
		"0 0 0 ? * 8",   // and end at 7
		"0 */0 * * * ?", // zero step
		"0 0 0 ? * MON", // names are not supported
		"? 0 0 * * ?",   // ? outside the day fields
		"0 0 0 1- * ?",  // malformed range
	}
	for _, expression := range cases {
		if _, err := Parse(expression); err == nil {
			t.Errorf("Parse(%q) succeeded, want an error", expression)
		}
	}
}
