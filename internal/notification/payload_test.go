package notification

import (
	"encoding/json"
	"strings"
	"testing"
)

const workspaceURL = "https://terrakube.acme.com/organizations/acme/workspaces/networking"

func decodePayload(t *testing.T, payload string) map[string]any {
	t.Helper()
	var decoded map[string]any
	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		t.Fatalf("payload is not JSON: %v\n%s", err, payload)
	}
	return decoded
}

func object(t *testing.T, value any) map[string]any {
	t.Helper()
	decoded, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("value is %T, want object", value)
	}
	return decoded
}

func array(t *testing.T, value any) []any {
	t.Helper()
	decoded, ok := value.([]any)
	if !ok {
		t.Fatalf("value is %T, want array", value)
	}
	return decoded
}

func TestBuildPayloadEveryChannelAndStyle(t *testing.T) {
	for _, channel := range []string{ChannelSlack, ChannelTeams, ChannelWebhook} {
		for _, style := range []string{StyleDetailed, StyleSimple} {
			t.Run(channel+"/"+style, func(t *testing.T) {
				payload, err := BuildPayload(channel, Context{
					OrganizationName: "acme", WorkspaceName: "networking", JobID: 42,
					Status: "completed", MessageStyle: style,
				})
				if err != nil {
					t.Fatalf("BuildPayload() error = %v", err)
				}
				if len(decodePayload(t, payload)) == 0 {
					t.Fatalf("BuildPayload() = %s, want a non-empty object", payload)
				}
			})
		}
	}
}

func TestBuildPayloadUnknownChannel(t *testing.T) {
	if _, err := BuildPayload("EMAIL", Context{}); err == nil {
		t.Fatal("BuildPayload() error = nil, want an error")
	}
}

func TestSlackDetailedPayload(t *testing.T) {
	payload, err := BuildPayload(ChannelSlack, Context{
		OrganizationName: "acme", WorkspaceName: "networking", JobID: 42, Status: "completed",
		RunURL: workspaceURL + "/runs/42", CommitID: "abc123", ConfigurationName: "Prod Alerts",
		WorkspaceURL: workspaceURL, MessageStyle: StyleDetailed,
	})
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	root := decodePayload(t, payload)
	text, _ := root["text"].(string)
	if !strings.Contains(text, "Run notification for") || !strings.Contains(text, workspaceURL) ||
		!strings.Contains(text, "acme/networking") {
		t.Fatalf("text = %q, want the workspace link", text)
	}
	if root["username"] != "networking" {
		t.Fatalf("username = %v, want networking", root["username"])
	}
	attachment := object(t, array(t, root["attachments"])[0])
	if attachment["color"] != "#2eb67d" {
		t.Fatalf("color = %v, want #2eb67d", attachment["color"])
	}
	blocks := array(t, attachment["blocks"])
	header := object(t, blocks[0])
	if header["type"] != "header" || !strings.Contains(object(t, header["text"])["text"].(string), "Completed") {
		t.Fatalf("header block = %v, want a Completed header", header)
	}
	for _, want := range []string{"acme", "42", workspaceURL + "/runs/42", "Prod Alerts"} {
		if !strings.Contains(payload, want) {
			t.Fatalf("payload missing %q: %s", want, payload)
		}
	}
}

func TestSlackFailedPayloadIncludesFailureReason(t *testing.T) {
	payload, err := BuildPayload(ChannelSlack, Context{
		OrganizationName: "acme", WorkspaceName: "networking", JobID: 43, Status: "failed",
		RunURL: workspaceURL + "/runs/43", FailureReason: "apply exited with code 1",
		WorkspaceURL: workspaceURL, MessageStyle: StyleDetailed,
	})
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	if !strings.Contains(payload, "apply exited with code 1") {
		t.Fatalf("payload missing failure reason: %s", payload)
	}
	attachment := object(t, array(t, decodePayload(t, payload)["attachments"])[0])
	if attachment["color"] != "#e01e5a" {
		t.Fatalf("color = %v, want #e01e5a", attachment["color"])
	}
}

func TestSlackDetailedOmitsAbsentLinks(t *testing.T) {
	payload, err := BuildPayload(ChannelSlack, Context{
		OrganizationName: "acme", WorkspaceName: "(test notification)", Status: "completed",
		ConfigurationName: "Test Notification", MessageStyle: StyleDetailed,
	})
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	blocks := array(t, object(t, array(t, decodePayload(t, payload)["attachments"])[0])["blocks"])
	for _, block := range blocks {
		if object(t, block)["type"] == "actions" {
			t.Fatalf("payload has an actions block without a run URL: %s", payload)
		}
	}
	if strings.Contains(payload, "Run notification for") {
		t.Fatalf("payload links a missing workspace: %s", payload)
	}
}

func TestSlackSimpleStyleIsOneLine(t *testing.T) {
	payload, err := BuildPayload(ChannelSlack, Context{
		OrganizationName: "acme", WorkspaceName: "networking", JobID: 42, Status: "running",
		RunURL: workspaceURL + "/runs/42", WorkspaceURL: workspaceURL, MessageStyle: StyleSimple,
	})
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	root := decodePayload(t, payload)
	text, _ := root["text"].(string)
	if !strings.Contains(text, "networking") || !strings.Contains(text, "Running") {
		t.Fatalf("text = %q, want workspace and Running", text)
	}
	if _, ok := root["attachments"]; ok {
		t.Fatalf("simple payload has attachments: %s", payload)
	}
}

func TestTeamsDetailedPayload(t *testing.T) {
	payload, err := BuildPayload(ChannelTeams, Context{
		OrganizationName: "acme", WorkspaceName: "networking", JobID: 42, Status: "completed",
		RunURL: workspaceURL + "/runs/42", CommitID: "abc123", ConfigurationName: "Prod Alerts",
		WorkspaceURL: workspaceURL, MessageStyle: StyleDetailed,
	})
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	root := decodePayload(t, payload)
	card := object(t, object(t, array(t, root["attachments"])[0])["content"])
	if card["type"] != "AdaptiveCard" {
		t.Fatalf("card type = %v, want AdaptiveCard", card["type"])
	}
	body := array(t, card["body"])
	if text, _ := object(t, body[0])["text"].(string); !strings.Contains(text, "Run notification for") ||
		!strings.Contains(text, workspaceURL) {
		t.Fatalf("workspace link = %q", text)
	}
	if text, _ := object(t, body[1])["text"].(string); !strings.Contains(text, "networking") {
		t.Fatalf("title = %q, want the workspace name", text)
	}
	facts := array(t, object(t, body[2])["facts"])
	if len(facts) != 4 {
		t.Fatalf("facts = %d, want 4 (organization, job, commit, notification)", len(facts))
	}
	if _, ok := card["actions"]; !ok {
		t.Fatalf("card has no View Run action: %s", payload)
	}
}

func TestTeamsOmitsOptionalParts(t *testing.T) {
	payload, err := BuildPayload(ChannelTeams, Context{
		OrganizationName: "acme", WorkspaceName: "networking", JobID: 42, Status: "completed",
		WorkspaceURL: workspaceURL, ConfigurationName: "Prod Alerts", MessageStyle: StyleDetailed,
	})
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	card := object(t, object(t, array(t, decodePayload(t, payload)["attachments"])[0])["content"])
	if _, ok := card["actions"]; ok {
		t.Fatalf("card has actions without a run URL: %s", payload)
	}
	body := array(t, card["body"])
	if len(body) != 3 {
		t.Fatalf("body blocks = %d, want 3 without a failure reason", len(body))
	}
	if strings.Contains(payload, `"title":"Commit"`) {
		t.Fatalf("payload has a commit fact without a commit: %s", payload)
	}
}

func TestTeamsFailedPayloadIncludesFailureReason(t *testing.T) {
	payload, err := BuildPayload(ChannelTeams, Context{
		OrganizationName: "acme", WorkspaceName: "networking", JobID: 43, Status: "failed",
		RunURL: workspaceURL + "/runs/43", FailureReason: "apply exited with code 1",
		WorkspaceURL: workspaceURL, MessageStyle: StyleDetailed,
	})
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	if !strings.Contains(payload, "apply exited with code 1") {
		t.Fatalf("payload missing failure reason: %s", payload)
	}
	card := object(t, object(t, array(t, decodePayload(t, payload)["attachments"])[0])["content"])
	if body := array(t, card["body"]); len(body) != 4 {
		t.Fatalf("body blocks = %d, want 4 with a failure reason", len(body))
	}
}

func TestTeamsSimpleStyleIsOneBlock(t *testing.T) {
	payload, err := BuildPayload(ChannelTeams, Context{
		OrganizationName: "acme", WorkspaceName: "networking", JobID: 42, Status: "running",
		RunURL: workspaceURL + "/runs/42", WorkspaceURL: workspaceURL, MessageStyle: StyleSimple,
	})
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	card := object(t, object(t, array(t, decodePayload(t, payload)["attachments"])[0])["content"])
	body := array(t, card["body"])
	if len(body) != 1 {
		t.Fatalf("body blocks = %d, want 1", len(body))
	}
	if text, _ := object(t, body[0])["text"].(string); !strings.Contains(text, "networking") ||
		!strings.Contains(text, "running") {
		t.Fatalf("text = %q, want workspace and status", text)
	}
	if _, ok := card["actions"]; ok {
		t.Fatalf("simple card has actions: %s", payload)
	}
	if strings.Contains(payload, "FactSet") || strings.Contains(payload, "Run notification for") {
		t.Fatalf("simple card is not minimal: %s", payload)
	}
}

func TestWebhookPayload(t *testing.T) {
	payload, err := BuildPayload(ChannelWebhook, Context{
		OrganizationName: "acme", WorkspaceName: "networking", JobID: 42, Status: "failed",
		RunURL: workspaceURL + "/runs/42", CommitID: "abc123", FailureReason: "apply exited with code 1",
		ConfigurationName: "Prod Alerts", WorkspaceURL: workspaceURL,
	})
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	root := decodePayload(t, payload)
	want := map[string]any{
		"organization":      "acme",
		"workspace":         "networking",
		"jobId":             float64(42),
		"status":            "failed",
		"runUrl":            workspaceURL + "/runs/42",
		"workspaceUrl":      workspaceURL,
		"commitId":          "abc123",
		"failureReason":     "apply exited with code 1",
		"configurationName": "Prod Alerts",
	}
	for key, value := range want {
		if root[key] != value {
			t.Fatalf("%s = %v, want %v", key, root[key], value)
		}
	}
}

func TestWebhookPayloadUsesNullForAbsentFields(t *testing.T) {
	payload, err := BuildPayload(ChannelWebhook, Context{
		OrganizationName: "acme", WorkspaceName: "(test notification)", Status: "completed",
		ConfigurationName: "Test Notification",
	})
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	root := decodePayload(t, payload)
	for _, key := range []string{"runUrl", "workspaceUrl", "commitId", "failureReason"} {
		if value, ok := root[key]; !ok || value != nil {
			t.Fatalf("%s = %v, want JSON null", key, value)
		}
	}
}
