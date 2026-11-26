// Package schedule parses upstream's Quartz cron expressions and turns a
// matching tick into a queued run. Upstream's scheduler is Quartz; the UI
// builds these expressions with cron-to-quartz, so the parser here accepts the
// same six- and seven-field shape.
package schedule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Cron is a parsed Quartz expression. Matching is evaluated against one
// instant, so the dispatcher can ask whether the current minute fires.
type Cron struct {
	seconds  []bool
	minutes  []bool
	hours    []bool
	days     []bool
	months   []bool
	weekdays []bool
	years    []bool

	// A "?" in day-of-month or day-of-week means "not specified": only the
	// other one constrains the day.
	anyDay     bool
	anyWeekday bool
	hasYear    bool
	expression string
}

// Parse parses a Quartz expression: seconds minutes hours day-of-month month
// day-of-week [year]. One of day-of-month and day-of-week must be "?"; Quartz
// refuses to combine them and its parser rejects the expression, which leaves
// the schedule stored but never firing upstream.
//
// Supported syntax is what cron-to-quartz emits and what people write by hand:
// "*", "?", single values, "a-b" ranges, comma lists and "/step" increments.
// The calendar extensions ("L", "W", "#") and month/weekday names are not
// supported and make Parse fail, so the dispatcher skips that schedule instead
// of firing on the wrong days.
func Parse(expression string) (*Cron, error) {
	fields := strings.Fields(expression)
	if len(fields) != 6 && len(fields) != 7 {
		return nil, fmt.Errorf("cron: want 6 or 7 fields, got %d", len(fields))
	}

	seconds, err := parseField(fields[0], 0, 59, false)
	if err != nil {
		return nil, fmt.Errorf("cron: seconds: %w", err)
	}
	minutes, err := parseField(fields[1], 0, 59, false)
	if err != nil {
		return nil, fmt.Errorf("cron: minutes: %w", err)
	}
	hours, err := parseField(fields[2], 0, 23, false)
	if err != nil {
		return nil, fmt.Errorf("cron: hours: %w", err)
	}
	months, err := parseField(fields[4], 1, 12, false)
	if err != nil {
		return nil, fmt.Errorf("cron: month: %w", err)
	}

	domPart, dowPart := fields[3], fields[5]
	if (domPart == "?") == (dowPart == "?") {
		return nil, errors.New("cron: exactly one of day-of-month and day-of-week must be ?")
	}
	days, err := parseField(domPart, 1, 31, true)
	if err != nil {
		return nil, fmt.Errorf("cron: day-of-month: %w", err)
	}
	weekdays, err := parseField(dowPart, 1, 7, true)
	if err != nil {
		return nil, fmt.Errorf("cron: day-of-week: %w", err)
	}

	cron := &Cron{
		seconds:    seconds,
		minutes:    minutes,
		hours:      hours,
		days:       days,
		months:     months,
		weekdays:   weekdays,
		anyDay:     domPart == "?",
		anyWeekday: dowPart == "?",
		expression: expression,
	}
	if len(fields) == 7 && fields[6] != "*" {
		years, err := parseField(fields[6], 1970, 2199, false)
		if err != nil {
			return nil, fmt.Errorf("cron: year: %w", err)
		}
		cron.years = years
		cron.hasYear = true
	}
	return cron, nil
}

// Matches reports whether the expression fires at the given instant, to the
// second.
func (c *Cron) Matches(t time.Time) bool {
	if !at(c.seconds, t.Second()) || !at(c.minutes, t.Minute()) || !at(c.hours, t.Hour()) {
		return false
	}
	if !at(c.months, int(t.Month())) {
		return false
	}
	if c.hasYear && !at(c.years, t.Year()) {
		return false
	}
	switch {
	case c.anyDay:
		return at(c.weekdays, int(t.Weekday())+1) // Quartz: Sunday=1 … Saturday=7
	case c.anyWeekday:
		return at(c.days, t.Day())
	default:
		return at(c.days, t.Day()) && at(c.weekdays, int(t.Weekday())+1)
	}
}

// String returns the expression, which is what logs should name.
func (c *Cron) String() string { return c.expression }

func at(set []bool, value int) bool {
	return value >= 0 && value < len(set) && set[value]
}

func parseField(expr string, min, max int, allowQuestion bool) ([]bool, error) {
	if expr == "?" {
		if !allowQuestion {
			return nil, errors.New("? is only allowed in day-of-month and day-of-week")
		}
		expr = "*"
	}
	set := make([]bool, max+1)
	for _, item := range strings.Split(expr, ",") {
		if item == "" {
			return nil, errors.New("empty list item")
		}
		base, step, stepped, err := splitStep(item)
		if err != nil {
			return nil, err
		}

		var lo, hi int
		switch {
		case base == "*":
			lo, hi = min, max
		case strings.Contains(base, "-"):
			parts := strings.SplitN(base, "-", 2)
			if lo, err = strconv.Atoi(parts[0]); err != nil {
				return nil, fmt.Errorf("invalid range %q", item)
			}
			if hi, err = strconv.Atoi(parts[1]); err != nil {
				return nil, fmt.Errorf("invalid range %q", item)
			}
		default:
			value, err := strconv.Atoi(base)
			if err != nil {
				return nil, fmt.Errorf("invalid value %q", item)
			}
			lo, hi = value, value
			// Quartz reads "a/n" as "a, a+n, … up to the maximum".
			if stepped {
				hi = max
			}
		}
		if lo < min || hi > max || lo > hi {
			return nil, fmt.Errorf("value %q outside %d-%d", item, min, max)
		}
		for value := lo; value <= hi; value += step {
			set[value] = true
		}
	}
	return set, nil
}

func splitStep(item string) (base string, step int, stepped bool, err error) {
	step = 1
	base = item
	if i := strings.Index(item, "/"); i >= 0 {
		base = item[:i]
		step, err = strconv.Atoi(item[i+1:])
		if err != nil || step <= 0 {
			return "", 0, false, fmt.Errorf("invalid step %q", item)
		}
		stepped = true
	}
	if base == "" {
		return "", 0, false, fmt.Errorf("empty value in %q", item)
	}
	return base, step, stepped, nil
}
