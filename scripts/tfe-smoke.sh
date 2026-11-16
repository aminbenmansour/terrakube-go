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
