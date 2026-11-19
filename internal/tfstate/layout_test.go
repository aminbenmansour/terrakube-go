package tfstate

import "testing"

func TestLayoutLocalKeys(t *testing.T) {
	layout := NewLayout("local")
	cases := map[string]string{
		layout.CurrentState("org", "ws"):         "backend/org/ws/terraform.tfstate",
		layout.RawState("org", "ws", "history"):  "state/org/ws/state/history.raw.json",
		layout.JSONState("org", "ws", "history"): "state/org/ws/state/history.json",
		layout.Plan("org", "ws", "job", "step"):  "state/org/ws/job/step/terraformLibrary.tfPlan",
		layout.Output("org", 7, "step"):          "tfoutput/org/7/step.tfoutput",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("key = %q, want %q", got, want)
		}
	}
}

func TestLayoutS3Keys(t *testing.T) {
	layout := NewLayout("s3")
	cases := map[string]string{
		layout.CurrentState("org", "ws"):         "tfstate/org/ws/terraform.tfstate",
		layout.RawState("org", "ws", "history"):  "tfstate/org/ws/state/history.raw.json",
		layout.JSONState("org", "ws", "history"): "tfstate/org/ws/state/history.json",
		layout.Plan("org", "ws", "job", "step"):  "tfstate/org/ws/job/step/terraformLibrary.tfPlan",
		layout.Output("org", 7, "step"):          "tfoutput/org/7/step.tfoutput",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("key = %q, want %q", got, want)
		}
	}
}
