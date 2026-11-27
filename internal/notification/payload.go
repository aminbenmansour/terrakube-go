// Package notification renders the provider-specific notification payloads and
// delivers them, mirroring upstream's payload builders and senders.
package notification

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

// Upstream's NotificationChannelType values.
const (
	ChannelSlack   = "SLACK"
	ChannelTeams   = "TEAMS"
	ChannelWebhook = "WEBHOOK"
)

// Upstream's NotificationMessageStyle values.
const (
	StyleDetailed = "DETAILED"
	StyleSimple   = "SIMPLE"
)

// Context is the data every payload builder renders. A test send has no run or
// workspace page to link to, so those URLs stay empty and JobID is zero, as
// upstream does.
type Context struct {
	OrganizationName  string
	WorkspaceName     string
	JobID             int64
	Status            string
	RunURL            string
	CommitID          string
	FailureReason     string
	ConfigurationName string
	WorkspaceURL      string
	MessageStyle      string
}

// BuildPayload renders the JSON payload for channelType. A blank MessageStyle
// is DETAILED, matching the configuration default.
func BuildPayload(channelType string, ctx Context) (string, error) {
	switch channelType {
	case ChannelSlack:
		return buildSlack(ctx)
	case ChannelTeams:
		return buildTeams(ctx)
	case ChannelWebhook:
		return buildWebhook(ctx)
	default:
		return "", fmt.Errorf("notification: no payload builder for channel type %s", channelType)
	}
}

func buildSlack(ctx Context) (string, error) {
	statusLabel := slackStatusLabel(ctx.Status)
	if ctx.MessageStyle == StyleSimple {
		return marshalPayload(map[string]any{
			"text":     ctx.WorkspaceName + " - " + statusLabel,
			"username": ctx.WorkspaceName,
		})
	}

	contextText := "Organization: " + ctx.OrganizationName + " | Job: #" + fmt.Sprint(ctx.JobID)
	if notBlank(ctx.CommitID) {
		contextText += " | Commit: " + ctx.CommitID
	}
	blocks := []any{
		map[string]any{
			"type": "header",
			"text": map[string]any{"type": "plain_text", "text": ctx.WorkspaceName + " - " + statusLabel},
		},
		map[string]any{
			"type":     "context",
			"elements": []any{map[string]any{"type": "mrkdwn", "text": contextText}},
		},
	}
	if ctx.Status == "failed" && notBlank(ctx.FailureReason) {
		blocks = append(blocks, map[string]any{
			"type": "section",
			"text": map[string]any{"type": "mrkdwn", "text": "*Failure reason:* " + ctx.FailureReason},
		})
	}
	if notBlank(ctx.RunURL) {
		blocks = append(blocks, map[string]any{
			"type": "actions",
			"elements": []any{map[string]any{
				"type": "button",
				"text": map[string]any{"type": "plain_text", "text": "View Run"},
				"url":  ctx.RunURL,
			}},
		})
	}
	if notBlank(ctx.ConfigurationName) {
		blocks = append(blocks, map[string]any{
			"type":     "context",
			"elements": []any{map[string]any{"type": "mrkdwn", "text": "Sent by notification: *" + ctx.ConfigurationName + "*"}},
		})
	}

	// Slack only shows the attachment colour bar when blocks sit inside the
	// legacy "attachments" wrapper.
	text := ctx.WorkspaceName + " - " + statusLabel
	if notBlank(ctx.WorkspaceURL) {
		text = "Run notification for <" + ctx.WorkspaceURL + "|" + ctx.OrganizationName + "/" + ctx.WorkspaceName + ">"
	}
	return marshalPayload(map[string]any{
		"text":     text,
		"username": ctx.WorkspaceName,
		"attachments": []any{map[string]any{
			"color":  slackStatusColor(ctx.Status),
			"blocks": blocks,
		}},
	})
}

func buildTeams(ctx Context) (string, error) {
	if ctx.MessageStyle == StyleSimple {
		return teamsPayload(map[string]any{
			"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
			"type":    "AdaptiveCard",
			"version": "1.4",
			"body": []any{map[string]any{
				"type": "TextBlock",
				"wrap": true,
				"text": ctx.WorkspaceName + " - " + ctx.Status,
			}},
		})
	}

	facts := []any{
		map[string]any{"title": "Organization", "value": ctx.OrganizationName},
		map[string]any{"title": "Job", "value": "#" + fmt.Sprint(ctx.JobID)},
	}
	if notBlank(ctx.CommitID) {
		facts = append(facts, map[string]any{"title": "Commit", "value": ctx.CommitID})
	}
	if notBlank(ctx.ConfigurationName) {
		facts = append(facts, map[string]any{"title": "Notification", "value": ctx.ConfigurationName})
	}

	var body []any
	if notBlank(ctx.WorkspaceURL) {
		body = append(body, map[string]any{
			"type":     "TextBlock",
			"isSubtle": true,
			"text": "Run notification for [" + ctx.OrganizationName + "/" + ctx.WorkspaceName + "](" +
				ctx.WorkspaceURL + ")",
		})
	}
	body = append(body,
		map[string]any{
			"type": "TextBlock", "size": "Large", "weight": "Bolder",
			"text": ctx.WorkspaceName + " - " + ctx.Status,
		},
		map[string]any{"type": "FactSet", "facts": facts},
	)
	if ctx.Status == "failed" && notBlank(ctx.FailureReason) {
		body = append(body, map[string]any{
			"type": "TextBlock",
			"wrap": true,
			"text": "Failure reason: " + ctx.FailureReason,
		})
	}

	card := map[string]any{
		"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
		"type":    "AdaptiveCard",
		"version": "1.4",
		"body":    body,
	}
	if notBlank(ctx.RunURL) {
		card["actions"] = []any{map[string]any{
			"type":  "Action.OpenUrl",
			"title": "View Run",
			"url":   ctx.RunURL,
		}}
	}
	return teamsPayload(card)
}

func teamsPayload(card map[string]any) (string, error) {
	return marshalPayload(map[string]any{
		"type": "message",
		"attachments": []any{map[string]any{
			"contentType": "application/vnd.microsoft.card.adaptive",
			"content":     card,
		}},
	})
}

func buildWebhook(ctx Context) (string, error) {
	// Upstream passes the optional fields straight through, so an absent value
	// is JSON null rather than an empty string.
	return marshalPayload(map[string]any{
		"organization":      ctx.OrganizationName,
		"workspace":         ctx.WorkspaceName,
		"jobId":             ctx.JobID,
		"status":            ctx.Status,
		"runUrl":            nullIfBlank(ctx.RunURL),
		"workspaceUrl":      nullIfBlank(ctx.WorkspaceURL),
		"commitId":          nullIfBlank(ctx.CommitID),
		"failureReason":     nullIfBlank(ctx.FailureReason),
		"configurationName": nullIfBlank(ctx.ConfigurationName),
	})
}

func marshalPayload(payload any) (string, error) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("notification: render payload: %w", err)
	}
	return string(encoded), nil
}

func slackStatusLabel(status string) string {
	switch status {
	case "completed", "noChanges":
		return "✅ Completed"
	case "waitingApproval":
		return "⚠️ Needs Approval"
	case "failed", "rejected", "cancelled":
		return "❌ " + capitalize(status)
	default:
		return capitalize(status)
	}
}

func slackStatusColor(status string) string {
	switch status {
	case "completed", "noChanges":
		return "#2eb67d"
	case "waitingApproval":
		return "#ecb22e"
	case "failed", "rejected", "cancelled":
		return "#e01e5a"
	default:
		return "#868686"
	}
}

func capitalize(value string) string {
	if value == "" {
		return value
	}
	runes := []rune(value)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func notBlank(value string) bool { return strings.TrimSpace(value) != "" }

func nullIfBlank(value string) any {
	if notBlank(value) {
		return value
	}
	return nil
}
