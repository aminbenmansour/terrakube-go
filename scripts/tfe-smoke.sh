#!/usr/bin/env bash
#
# Smoke test for the TFE v2 API at /remote/tfe/v2. Requires a running api and jq.
#
#   TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/tfe-smoke.sh
#
set -euo pipefail

ROOT="${TERRAKUBE_API_URL:-http://localhost:8080}"
ROOT="${ROOT%/}"
API="${ROOT}/api/v1"
TFE="${ROOT}/remote/tfe/v2"
TOKEN="${TERRAKUBE_AUTH_TOKEN:-}"

command -v jq >/dev/null || { echo "jq is required" >&2; exit 1; }

auth=()
if [[ -n "$TOKEN" ]]; then
  auth=(-H "Authorization: Bearer ${TOKEN}")
fi
json=(-H "Content-Type: application/vnd.api+json")

req() { curl -fsS "${auth[@]}" "${json[@]}" "$@"; }

status() {
  curl -s -o /dev/null -w '%{http_code}' "${auth[@]}" "${json[@]}" "$@"
}

expect_status() {
  local want=$1 got=$2 what=$3
  if [[ "$got" != "$want" ]]; then
    echo "FAIL: ${what}: got HTTP ${got}, want ${want}" >&2
    exit 1
  fi
}

echo "Waiting for ${ROOT}..."
for _ in $(seq 1 30); do
  if curl -fsS "${ROOT}/healthz" >/dev/null 2>&1; then break; fi
  sleep 0.5
done

echo "--- discovery ---"
discovery=$(curl -fsS "${ROOT}/.well-known/terraform.json")
[[ "$(jq -r '.["tfe.v2"]' <<<"$discovery")" == "/remote/tfe/v2/" ]] || {
  echo "FAIL: discovery tfe.v2" >&2; exit 1;
}
[[ "$(jq -r '.["state.v2"]' <<<"$discovery")" == "/remote/state/v2/" ]] || {
  echo "FAIL: discovery state.v2" >&2; exit 1;
}
echo "ok: discovery"

echo "--- ping ---"
headers=$(mktemp)
trap 'rm -f "$headers"' EXIT
expect_status 404 "$(curl -sS -D "$headers" -o /dev/null -w '%{http_code}' "${TFE}/ping")" "ping"
grep -qi '^TFP-API-Version: 2.5' "$headers" || { echo "FAIL: ping API version header" >&2; exit 1; }
grep -qi '^TFP-AppName: Terrakube' "$headers" || { echo "FAIL: ping app header" >&2; exit 1; }
echo "ok: ping"

if [[ -n "$TOKEN" ]]; then
  expect_status 401 "$(curl -s -o /dev/null -w '%{http_code}' "${TFE}/organizations/missing")" "TFE auth"
  echo "ok: auth enforced"
fi

name="tfe-smoke-$(date +%s)"
org_id=$(req -X POST "${API}/organization" -d "{
  \"data\": {\"type\": \"organization\", \"attributes\": {\"name\": \"${name}\", \"description\": \"tfe smoke\"}}
}" | jq -r '.data.id')
[[ -n "$org_id" && "$org_id" != "null" ]] || { echo "FAIL: organization id" >&2; exit 1; }
echo "ok: organization ${org_id}"

project_id=$(req -X POST "${API}/organization/${org_id}/project" -d '{
  "data": {"type": "project", "attributes": {"name": "tfe-smoke-project", "description": "smoke"}}
}' | jq -r '.data.id')
ws_id=$(req -X POST "${API}/organization/${org_id}/workspace" -d "{
  \"data\": {
    \"type\": \"workspace\",
    \"attributes\": {\"name\": \"tfe-smoke-workspace\", \"terraformVersion\": \"1.9.0\", \"executionMode\": \"remote\"},
    \"relationships\": {\"project\": {\"data\": {\"type\": \"project\", \"id\": \"${project_id}\"}}}
  }
}" | jq -r '.data.id')
[[ -n "$ws_id" && "$ws_id" != "null" ]] || { echo "FAIL: workspace id" >&2; exit 1; }
echo "ok: workspace ${ws_id}"

org_name=$(req "${TFE}/organizations/${name}" | jq -r '.data.attributes.name')
[[ "$org_name" == "$name" ]] || { echo "FAIL: TFE organization" >&2; exit 1; }
echo "ok: TFE organization"

[[ "$(req "${TFE}/organizations/${name}/entitlement-set" | jq -r '.data.attributes["state-storage"]')" == "true" ]] || {
  echo "FAIL: entitlement set" >&2; exit 1;
}
[[ "$(req "${TFE}/organizations/${name}/capacity" | jq -r '.data.attributes.pending')" == "0" ]] || {
  echo "FAIL: capacity" >&2; exit 1;
}
echo "ok: entitlement and capacity"

listed=$(req "${TFE}/organizations/${name}/workspaces?search%5Bname%5D=tfe-smoke" | jq -r '.data[0].id')
[[ "$listed" == "$ws_id" ]] || { echo "FAIL: list workspaces (${listed})" >&2; exit 1; }
echo "ok: list workspaces"

got=$(req "${TFE}/organizations/${name}/workspaces/tfe-smoke-workspace" | jq -r '.data.id')
[[ "$got" == "$ws_id" ]] || { echo "FAIL: get workspace" >&2; exit 1; }
echo "ok: get workspace"

created=$(req -X POST "${TFE}/organizations/${name}/workspaces" -d "{
  \"data\": {\"type\": \"workspaces\", \"attributes\": {\"name\": \"tfe-created\", \"terraform-version\": \"1.9.0\"}}
}" | jq -r '.data.attributes.name')
[[ "$created" == "tfe-created" ]] || { echo "FAIL: create workspace" >&2; exit 1; }
echo "ok: create workspace"

binding=$(req -X PATCH "${TFE}/workspaces/${ws_id}/tag-bindings" -d '{
  "data": [{"type": "tag-bindings", "attributes": {"key": "env", "value": "smoke"}}]
}' | jq -r '.data[0].attributes.key')
[[ "$binding" == "env" ]] || { echo "FAIL: tag binding write" >&2; exit 1; }
req -X POST "${TFE}/workspaces/${ws_id}/relationships/tags" -d '{
  "data": [{"type": "tags", "attributes": {"name": "team"}}]
}' >/dev/null
count=$(req "${TFE}/workspaces/${ws_id}/tag-bindings" | jq '.data | length')
[[ "$count" == "2" ]] || { echo "FAIL: tag binding count (${count})" >&2; exit 1; }
echo "ok: tag bindings"

state_b64=$(printf '{"version":4,"terraform_version":"1.9.0","serial":3,"lineage":"smoke-lineage","outputs":{"smoke":{"value":"ok","type":"string","sensitive":false}}}' | base64 -w0)
state_version=$(req -X POST "${TFE}/workspaces/${ws_id}/state-versions" -d "{
  \"data\": {\"type\": \"state-versions\", \"attributes\": {\"state\": \"${state_b64}\"}}
}" | jq -r '.data.id')
state_status=$(req "${TFE}/state-versions/${state_version}" | jq -r '.data.attributes.status')
[[ "$state_status" == "finalized" ]] || { echo "FAIL: state version status (${state_status})" >&2; exit 1; }
[[ "$(req "${TFE}/workspaces/${ws_id}/current-state-version" | jq -r '.data.attributes.serial')" == "3" ]] || {
  echo "FAIL: current state version" >&2; exit 1;
}
[[ "$(req "${TFE}/workspaces/${ws_id}/current-state-version-outputs" | jq -r '.data[0].attributes.value')" == "ok" ]] || {
  echo "FAIL: current state outputs" >&2; exit 1;
}
echo "ok: state versions"

config_response=$(req -X POST "${TFE}/workspaces/${ws_id}/configuration-versions" -d '{
  "data": {"type": "configuration-versions", "attributes": {"speculative": true, "auto-queue-runs": false}}
}')
config_id=$(jq -r '.data.id' <<<"$config_response")
upload_url=$(jq -r '.data.attributes["upload-url"]' <<<"$config_response")
upload_path=$(printf '%s' "$upload_url" | sed -E 's#^https?://[^/]+##')
printf 'smoke-configuration' | curl -fsS -X PUT "${ROOT}${upload_path}" --data-binary @- >/dev/null
config_status=$(req "${TFE}/configuration-versions/${config_id}" | jq -r '.data.attributes.status')
[[ "$config_status" == "uploaded" ]] || { echo "FAIL: configuration status (${config_status})" >&2; exit 1; }
downloaded=$(curl -fsS "${TFE}/configuration-versions/${config_id}/terraformContent.tar.gz")
[[ "$downloaded" == "smoke-configuration" ]] || { echo "FAIL: configuration download" >&2; exit 1; }
echo "ok: configuration versions"

run_response=$(req -X POST "${TFE}/runs" -d "{
  \"data\": {
    \"type\": \"runs\",
    \"attributes\": {\"is-destroy\": false, \"auto-apply\": false, \"refresh\": true},
    \"relationships\": {
      \"workspace\": {\"data\": {\"type\": \"workspaces\", \"id\": \"${ws_id}\"}},
      \"configuration-version\": {\"data\": {\"type\": \"configuration-versions\", \"id\": \"${config_id}\"}}
    }
  }
}")
run_id=$(jq -r '.data.id' <<<"$run_response")
run_number=${run_id#run-}
[[ "$(jq -r '.data.attributes.status' <<<"$run_response")" == "pending" ]] || {
  echo "FAIL: create run" >&2; exit 1;
}
# Active-run listings only contain pending, queue, running and
# waitingApproval jobs. When an executor is running it can pick the run up
# between creation and this check, so accept either "listed" or "advanced".
active_run=$(req "${TFE}/workspaces/${ws_id}/runs" | jq -r '.data[0].id // empty')
if [[ -n "$active_run" ]]; then
  [[ "$active_run" == "$run_id" ]] || { echo "FAIL: list workspace runs" >&2; exit 1; }
else
  run_status=$(req "${TFE}/runs/${run_id}" | jq -r '.data.attributes.status')
  [[ "$run_status" != "pending" ]] || { echo "FAIL: list workspace runs (run still pending but not listed)" >&2; exit 1; }
fi
active_run=$(req "${TFE}/organizations/${name}/runs/queue" | jq -r '.data[0].id // empty')
if [[ -n "$active_run" ]]; then
  [[ "$active_run" == "$run_id" ]] || { echo "FAIL: run queue" >&2; exit 1; }
else
  run_status=$(req "${TFE}/runs/${run_id}" | jq -r '.data.attributes.status')
  [[ "$run_status" != "pending" ]] || { echo "FAIL: run queue (run still pending but not listed)" >&2; exit 1; }
fi
[[ "$(req "${TFE}/runs/${run_id}/run-events" | jq -r '.data | length')" == "0" ]] || {
  echo "FAIL: run events" >&2; exit 1;
}
[[ "$(req "${TFE}/runs/${run_id}" | jq -r '.data.relationships["run-events"].data | length')" == "0" ]] || {
  echo "FAIL: run-events relationship" >&2; exit 1;
}
[[ "$(req "${TFE}/runs/${run_id}?include=plan,apply" | jq -r '.included | length')" == "2" ]] || {
  echo "FAIL: run include plan/apply" >&2; exit 1;
}
[[ "$(req "${TFE}/plans/${run_number}" | jq -r '.data.type')" == "plans" ]] || {
  echo "FAIL: plan" >&2; exit 1;
}
[[ "$(req "${TFE}/applies/${run_number}" | jq -r '.data.type')" == "applies" ]] || {
  echo "FAIL: apply" >&2; exit 1;
}
# Discard only succeeds while the run still holds the queue. When an executor
# is running it may have already picked the run up, so accept the upstream
# conflict response in that case.
discard_code=$(status -X POST "${TFE}/runs/${run_id}/actions/discard")
case "$discard_code" in
  200)
    [[ "$(req "${TFE}/runs/${run_id}" | jq -r '.data.attributes.status')" == "unknown" ]] || {
      echo "FAIL: discard run status" >&2; exit 1;
    }
    ;;
  409)
    run_status=$(req "${TFE}/runs/${run_id}" | jq -r '.data.attributes.status')
    [[ "$run_status" != "pending" ]] || { echo "FAIL: discard run conflict while pending" >&2; exit 1; }
    ;;
  *)
    echo "FAIL: discard run: got HTTP ${discard_code}" >&2; exit 1;
    ;;
esac
echo "ok: runs"

locked=$(req -X POST "${TFE}/workspaces/${ws_id}/actions/lock" | jq -r '.data.attributes.locked')
[[ "$locked" == "true" ]] || { echo "FAIL: lock workspace" >&2; exit 1; }
expect_status 409 "$(status -X POST "${TFE}/workspaces/${ws_id}/actions/lock")" "lock conflict"
unlocked=$(req -X POST "${TFE}/workspaces/${ws_id}/actions/unlock" | jq -r '.data.attributes.locked')
[[ "$unlocked" == "false" ]] || { echo "FAIL: unlock workspace" >&2; exit 1; }
echo "ok: lock and unlock"

projects=$(req "${TFE}/organizations/${name}/projects" | jq -r '.data[0].id')
[[ "$projects" == "$project_id" ]] || { echo "FAIL: list projects" >&2; exit 1; }
echo "ok: list projects"

expect_status 204 "$(status -X DELETE "${API}/workspace/${ws_id}")" "delete workspace"
expect_status 204 "$(status -X DELETE "${API}/project/${project_id}")" "delete project"
expect_status 204 "$(status -X DELETE "${API}/organization/${org_id}")" "delete organization"

echo "tfe smoke test passed"
