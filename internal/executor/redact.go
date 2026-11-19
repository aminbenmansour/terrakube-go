package executor

import (
	"sort"
	"strings"
)

const redactedValue = "[REDACTED]"

// redactor replaces known secret values in command output and error text.
// Terraform already masks values it knows are sensitive, but a command can
// still print an environment variable or a secret in an error message.
type redactor struct {
	values []string
}

// newRedactor builds a redactor from sensitive variable values and internal
// tokens. Empty values and duplicates are dropped, and longer values are
// replaced first so a secret that contains another secret is fully masked.
func newRedactor(secrets []string) *redactor {
	seen := make(map[string]struct{}, len(secrets))
	values := make([]string, 0, len(secrets))
	for _, secret := range secrets {
		if secret == "" {
			continue
		}
		if _, ok := seen[secret]; ok {
			continue
		}
		seen[secret] = struct{}{}
		values = append(values, secret)
	}
	if len(values) == 0 {
		return &redactor{}
	}
	sort.Slice(values, func(i, j int) bool { return len(values[i]) > len(values[j]) })
	return &redactor{values: values}
}

// redact masks every known secret in line. A nil redactor is a no-op.
func (r *redactor) redact(line string) string {
	if r == nil || len(r.values) == 0 || line == "" {
		return line
	}
	for _, value := range r.values {
		line = strings.ReplaceAll(line, value, redactedValue)
	}
	return line
}
