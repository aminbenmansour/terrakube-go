package webhook

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Delivery is a provider webhook payload reduced to the fields a trigger
// matches on. Event is the Terrakube trigger type, not the provider's name.
type Delivery struct {
	Event      string
	Repository string
	Branch     string
	Commit     string
	Actor      string
	Files      []string
	ReleaseTag string
}

const (
	eventPush        = "PUSH"
	eventPullRequest = "PULL_REQUEST"
	eventRelease     = "RELEASE"
	eventPing        = "PING"
)

// NormalizeRepository reduces an HTTPS, SSH or scp-style clone URL to
// host/owner/repo so a workspace source and a payload repository compare equal.
func NormalizeRepository(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	value = strings.TrimSuffix(value, "/")
	value = strings.TrimSuffix(value, ".git")
	if rest, found := strings.CutPrefix(value, "git@"); found {
		value = rest
	}
	if index := strings.Index(value, "://"); index >= 0 {
		value = value[index+3:]
	}
	// Drop credentials and the scp-style colon in one step.
	if index := strings.Index(value, "@"); index >= 0 {
		value = value[index+1:]
	}
	if index := strings.Index(value, ":"); index >= 0 {
		value = value[:index] + "/" + value[index+1:]
	}
	return strings.Trim(value, "/")
}

// branchFromRef turns refs/heads/main into main and refs/tags/v1 into v1,
// keeping any slashes in the name.
func branchFromRef(ref string) (branch, tag string) {
	switch {
	case strings.HasPrefix(ref, "refs/heads/"):
		return strings.TrimPrefix(ref, "refs/heads/"), ""
	case strings.HasPrefix(ref, "refs/tags/"):
		tag = strings.TrimPrefix(ref, "refs/tags/")
		return tag, tag
	default:
		return ref, ""
	}
}

func deliveryID(headerValue string, body []byte) string {
	if headerValue != "" && len(headerValue) <= 128 {
		return headerValue
	}
	sum := sha256.Sum256(append([]byte(headerValue), body...))
	return hex.EncodeToString(sum[:])
}

func nonEmptyPtr(actor string) *string {
	if strings.TrimSpace(actor) == "" {
		return nil
	}
	return &actor
}

// viaFor uses upstream's JobVia values, which the UI matches case-insensitively.
func viaFor(provider string) string {
	switch provider {
	case "GITHUB":
		return "Github"
	case "GITLAB":
		return "Gitlab"
	default:
		return provider
	}
}
