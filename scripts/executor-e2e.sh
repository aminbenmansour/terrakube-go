#!/usr/bin/env bash
#
# End-to-end executor test: uploads a real OpenTofu configuration, lets the
# executor plan and apply it, and verifies logs, state outputs and history.
#
# Requires a running api and executor, jq and tofu (or terraform):
#
#   TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/executor-e2e.sh
#
# The API must be able to reach itself at TERRAKUBE_PUBLIC_URL, because the
# executor writes that address into the generated HTTP backend override.
set -euo pipefail

ROOT="${TERRAKUBE_API_URL:-http://localhost:8080}"
ROOT="${ROOT%/}"
API="${ROOT}/api/v1"
TFE="${ROOT}/remote/tfe/v2"
TOKEN="${TERRAKUBE_AUTH_TOKEN:-}"
TIMEOUT="${TERRAKUBE_E2E_TIMEOUT:-180}"

command -v jq >/dev/null || { echo "jq is required" >&2; exit 1; }
command -v tar >/dev/null || { echo "tar is required" >&2; exit 1; }

auth=()
if [[ -n "$TOKEN" ]]; then
  auth=(-H "Authorization: Bearer ${TOKEN}")
fi
json=(-H "Content-Type: application/vnd.api+json")

req() { curl -fsS "${auth[@]}" "${json[@]}" "$@"; }
status() { curl -s -o /dev/null -w '%{http_code}' "${auth[@]}" "${json[@]}" "$@"; }

fail() { echo "FAIL: $*" >&2; exit 1; }

wait_for() {
  local what=$1 expression=$2 url=$3
  local deadline=$((SECONDS + TIMEOUT))
  while (( SECONDS < deadline )); do
    local value
    value=$(req "$url" | jq -r "$expression")
    case "$value" in
      finished|errored|failed) echo "$value"; return 0 ;;
    esac
    sleep 1
  done
  fail "timed out waiting for ${what} (${url})"
}

echo "Waiting for ${ROOT}..."
for _ in $(seq 1 30); do
  curl -fsS "${ROOT}/healthz" >/dev/null 2>&1 && break
  sleep 0.5
done

echo "--- create workspace ---"
name="e2e-$(date +%s)"
org_id=$(req -X POST "${API}/organization" -d "{
  \"data\": {\"type\": \"organization\", \"attributes\": {\"name\": \"${name}\"}}
}" | jq -r '.data.id')
project_id=$(req -X POST "${API}/organization/${org_id}/project" -d '{
  "data": {"type": "project", "attributes": {"name": "e2e"}}
}' | jq -r '.data.id')
ws_id=$(req -X POST "${API}/organization/${org_id}/workspace" -d "{
  \"data\": {
    \"type\": \"workspace\",
    \"attributes\": {\"name\": \"executor\", \"iacType\": \"opentofu\"},
    \"relationships\": {\"project\": {\"data\": {\"type\": \"project\", \"id\": \"${project_id}\"}}}
  }
}" | jq -r '.data.id')
req -X POST "${API}/workspace/${ws_id}/variable" -d '{
  "data": {"type": "variable", "attributes": {"key": "message", "value": "hello-executor", "category": "TERRAFORM"}}
}' >/dev/null
echo "ok: organization ${org_id}, workspace ${ws_id}"

echo "--- upload configuration ---"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
cat > "${work}/main.tf" <<'EOF'
variable "message" {
  type = string
}

resource "terraform_data" "example" {
  input = var.message
}

output "message" {
  value = terraform_data.example.output
}
EOF
tar -czf "${work}/config.tar.gz" -C "${work}" main.tf

config_response=$(req -X POST "${TFE}/workspaces/${ws_id}/configuration-versions" -d '{
  "data": {"type": "configuration-versions", "attributes": {"speculative": false, "auto-queue-runs": false}}
}')
config_id=$(jq -r '.data.id' <<<"$config_response")
upload_path=$(jq -r '.data.attributes["upload-url"]' <<<"$config_response" | sed -E 's#^https?://[^/]+##')
curl -fsS -X PUT "${ROOT}${upload_path}" --data-binary @"${work}/config.tar.gz" >/dev/null
config_status=$(req "${TFE}/configuration-versions/${config_id}" | jq -r '.data.attributes.status')
[[ "$config_status" == "uploaded" ]] || fail "configuration status ${config_status}"
echo "ok: configuration ${config_id}"

echo "--- plan ---"
run_id=$(req -X POST "${TFE}/runs" -d "{
  \"data\": {
    \"type\": \"runs\",
    \"attributes\": {\"is-destroy\": false, \"auto-apply\": false, \"refresh\": true},
    \"relationships\": {
      \"workspace\": {\"data\": {\"type\": \"workspaces\", \"id\": \"${ws_id}\"}},
      \"configuration-version\": {\"data\": {\"type\": \"configuration-versions\", \"id\": \"${config_id}\"}}
    }
  }
}" | jq -r '.data.id')
run_number=${run_id#run-}

plan_status=$(wait_for "plan" '.data.attributes.status' "${TFE}/plans/${run_number}")
[[ "$plan_status" == "finished" ]] || fail "plan status ${plan_status}"

log_url=$(req "${TFE}/plans/${run_number}" | jq -r '.data.attributes["log-read-url"]')
log_path=$(printf '%s' "$log_url" | sed -E 's#^https?://[^/]+##')
plan_log=$(curl -fsS "${ROOT}${log_path}?offset=0&limit=200000")
grep -q "Executing terraform plan" <<<"$plan_log" || fail "plan log does not contain the plan phase"
grep -q "Plan:" <<<"$plan_log" || fail "plan log does not contain a plan summary"
echo "ok: plan finished and produced a log"

echo "--- apply ---"
req -X POST "${TFE}/runs/${run_id}/actions/apply" >/dev/null
apply_status=$(wait_for "apply" '.data.attributes.status' "${TFE}/applies/${run_number}")
[[ "$apply_status" == "finished" ]] || fail "apply status ${apply_status}"

apply_log_url=$(req "${TFE}/applies/${run_number}" | jq -r '.data.attributes["log-read-url"]')
apply_log_path=$(printf '%s' "$apply_log_url" | sed -E 's#^https?://[^/]+##')
apply_log=$(curl -fsS "${ROOT}${apply_log_path}?offset=0&limit=200000")
grep -q "Executing terraform apply" <<<"$apply_log" || fail "apply log does not contain the apply phase"
grep -q "Apply complete!" <<<"$apply_log" || fail "apply log does not report success"

run_status=$(req "${TFE}/runs/${run_id}" | jq -r '.data.attributes.status')
[[ "$run_status" == "finished" ]] || fail "run status ${run_status}"
echo "ok: apply finished"

echo "--- state and outputs ---"
serial=$(req "${TFE}/workspaces/${ws_id}/current-state-version" | jq -r '.data.attributes.serial')
[[ "$serial" == "1" ]] || fail "current state serial ${serial}"
value=$(req "${TFE}/workspaces/${ws_id}/current-state-version-outputs" | jq -r '.data[0].attributes.value')
[[ "$value" == "hello-executor" ]] || fail "output value ${value}"
echo "ok: state serial ${serial}, output ${value}"

echo "--- destroy ---"
destroy_run=$(req -X POST "${TFE}/runs" -d "{
  \"data\": {
    \"type\": \"runs\",
    \"attributes\": {\"is-destroy\": true, \"auto-apply\": true, \"refresh\": true},
    \"relationships\": {
      \"workspace\": {\"data\": {\"type\": \"workspaces\", \"id\": \"${ws_id}\"}},
      \"configuration-version\": {\"data\": {\"type\": \"configuration-versions\", \"id\": \"${config_id}\"}}
    }
  }
}" | jq -r '.data.id')
destroy_number=${destroy_run#run-}
destroy_status=$(wait_for "destroy" '.data.attributes.status' "${TFE}/applies/${destroy_number}")
[[ "$destroy_status" == "finished" ]] || fail "destroy status ${destroy_status}"
echo "ok: destroy finished"

echo "--- cleanup ---"
[[ "$(status -X DELETE "${API}/workspace/${ws_id}")" == "204" ]] || fail "delete workspace"
[[ "$(status -X DELETE "${API}/project/${project_id}")" == "204" ]] || fail "delete project"
[[ "$(status -X DELETE "${API}/organization/${org_id}")" == "204" ]] || fail "delete organization"

echo "executor e2e test passed"
