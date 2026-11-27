package webhook

import (
	"regexp"
	"strings"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// matchTrigger returns the highest-priority trigger for the delivery's event
// whose branch and path patterns match, the order upstream checks them in.
func matchTrigger(events []store.WebhookEvent, delivery Delivery) (store.WebhookEvent, bool) {
	for _, event := range events {
		if event.Event != delivery.Event {
			continue
		}
		if !matchBranch(delivery.Branch, event.Branch) {
			continue
		}
		if !matchPaths(delivery.Files, event) {
			continue
		}
		return event, true
	}
	return store.WebhookEvent{}, false
}

// matchBranch treats each comma-separated pattern as a full regular expression,
// as upstream's String.matches does.
func matchBranch(branch, patterns string) bool {
	if branch == "" {
		return false
	}
	for _, pattern := range strings.Split(patterns, ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		re, err := regexp.Compile("^(?:" + pattern + ")$")
		if err != nil {
			continue
		}
		if re.MatchString(branch) {
			return true
		}
	}
	return false
}

// matchPaths applies the trigger's path filter. With no changed files to check
// (a pull request payload carries none) the branch alone decides, like the
// polled-push path upstream; a file list that is present must match.
func matchPaths(files []string, event store.WebhookEvent) bool {
	if len(files) == 0 {
		return true
	}
	for _, pattern := range strings.Split(event.Path, ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		matcher := pathRegex(event.PathType, pattern)
		if matcher == nil {
			continue
		}
		for _, file := range files {
			if matcher.MatchString(file) {
				return true
			}
		}
	}
	return false
}

func pathRegex(pathType, pattern string) *regexp.Regexp {
	if pathType == "PATTERN" {
		pattern = globToRegex(pattern)
	}
	re, err := regexp.Compile("^(?:" + pattern + ")$")
	if err != nil {
		return nil
	}
	return re
}

// globToRegex mirrors upstream's globToSafeRegex: * spans directories, ? is one
// character, and every other regular-expression metacharacter is literal.
func globToRegex(glob string) string {
	var builder strings.Builder
	for _, char := range glob {
		switch char {
		case '*':
			builder.WriteString(".*")
		case '?':
			builder.WriteString(".")
		default:
			if strings.ContainsRune(`()[]{}|+^$\.`, char) {
				builder.WriteByte('\\')
			}
			builder.WriteRune(char)
		}
	}
	return builder.String()
}
