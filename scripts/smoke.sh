#!/usr/bin/env bash
#
# End-to-end smoke test for the management API. Requires a running api and jq.
#
#   TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/smoke.sh
#
set -euo pipefail

BASE="${TERRAKUBE_API_URL:-http://localhost:8080}"
BASE="${BASE%/}/api/v1"
TOKEN="${TERRAKUBE_AUTH_TOKEN:-}"

command -v jq >/dev/null || { echo "jq is required" >&2; exit 1; }

auth=()
if [[ -n "$TOKEN" ]]; then
  auth=(-H "Authorization: Bearer ${TOKEN}")
fi
json=(-H "Content-Type: application/vnd.api+json")

req() { curl -fsS "${auth[@]}" "${json[@]}" "$@"; }

status() {
  local code
  code=$(curl -s -o /dev/null -w '%{http_code}' "${auth[@]}" "${json[@]}" "$@")
  echo "$code"
}

expect_status() {
  local want=$1 got=$2 what=$3
  if [[ "$got" != "$want" ]]; then
    echo "FAIL: ${what}: got HTTP ${got}, want ${want}" >&2
    exit 1
  fi
}

echo "Waiting for ${BASE}..."
for _ in $(seq 1 30); do
  if curl -fsS "${BASE%/api/v1}/healthz" >/dev/null 2>&1; then break; fi
  sleep 0.5
done

if [[ -n "$TOKEN" ]]; then
  code=$(curl -s -o /dev/null -w '%{http_code}' "${BASE}/organization")
  expect_status 401 "$code" "request without token"
  echo "ok: auth enforced"
fi

name="smoke-$(date +%s)"

org_id=$(req -X POST "${BASE}/organization" -d "{
  \"data\": {\"type\": \"organization\", \"attributes\": {\"name\": \"${name}\", \"description\": \"smoke test\"}}
}" | jq -r '.data.id')
[[ -n "$org_id" && "$org_id" != "null" ]] || { echo "FAIL: organization id" >&2; exit 1; }

# A failure halfway through would otherwise leave a half-built organisation in
# the database for every run.
cleanup() {
  curl -s -o /dev/null -X DELETE "${auth[@]}" "${json[@]}" "${BASE}/organization/${org_id}" || true
}
trap cleanup EXIT
echo "ok: organization ${org_id}"

got_name=$(req "${BASE}/organization/${org_id}" | jq -r '.data.attributes.name')
[[ "$got_name" == "$name" ]] || { echo "FAIL: organization name ${got_name}" >&2; exit 1; }
echo "ok: get organization"

req "${BASE}/organization?filter%5Borganization%5D=name%3D%3D${name}" >/dev/null
echo "ok: filter organizations"

project_id=$(req -X POST "${BASE}/organization/${org_id}/project" -d '{
  "data": {"type": "project", "attributes": {"name": "smoke-project", "description": "smoke"}}
}' | jq -r '.data.id')
echo "ok: project ${project_id}"

ws_id=$(req -X POST "${BASE}/organization/${org_id}/workspace" -d "{
  \"data\": {
    \"type\": \"workspace\",
    \"attributes\": {\"name\": \"smoke-workspace\", \"description\": \"before\", \"iacType\": \"terraform\"},
    \"relationships\": {\"project\": {\"data\": {\"type\": \"project\", \"id\": \"${project_id}\"}}}
  }
}" | jq -r '.data.id')
echo "ok: workspace ${ws_id}"

gql="${BASE%/api/v1}/graphql/api/v1"
gql_query() {
  curl -fsS "${auth[@]}" -H "Content-Type: application/json" -d "$1" "${gql}"
}
ws_page=$(gql_query '{"query":"query($ids:[String]){ organization(ids:$ids){ edges{ node{ workspace(first:10, sort:\"name\", filter:\"name=ini=*smoke*\"){ edges{ node{ id name workspaceTag{edges{node{tagId}}} project{edges{node{id}}} } } pageInfo{ totalRecords } } } } } }","variables":{"ids":["'"${org_id}"'"]}}')
gql_count=$(echo "$ws_page" | jq -r '.data.organization.edges[0].node.workspace.pageInfo.totalRecords')
gql_id=$(echo "$ws_page" | jq -r '.data.organization.edges[0].node.workspace.edges[0].node.id')
[[ "$gql_count" == "1" && "$gql_id" == "$ws_id" ]] || {
  echo "FAIL: graphql workspace list (${ws_page})" >&2
  exit 1
}
echo "ok: graphql workspace list"

project_ref=$(req "${BASE}/workspace/${ws_id}" | jq -r '.data.relationships.project.data.id')
[[ "$project_ref" == "$project_id" ]] || { echo "FAIL: workspace project relationship" >&2; exit 1; }
echo "ok: workspace relationship"

req -X PATCH "${BASE}/workspace/${ws_id}" -d '{
  "data": {"type": "workspace", "id": "'"${ws_id}"'", "attributes": {"description": "after"}}
}' >/dev/null
desc=$(req "${BASE}/workspace/${ws_id}" | jq -r '.data.attributes.description')
[[ "$desc" == "after" ]] || { echo "FAIL: workspace patch (${desc})" >&2; exit 1; }
echo "ok: patch workspace"

var_id=$(req -X POST "${BASE}/workspace/${ws_id}/variable" -d '{
  "data": {"type": "variable", "attributes": {"key": "region", "value": "eu-west-1", "category": "ENV", "sensitive": false}}
}' | jq -r '.data.id')
echo "ok: variable ${var_id}"

count=$(req "${BASE}/workspace/${ws_id}/variable" | jq '.data | length')
[[ "$count" == "1" ]] || { echo "FAIL: variable count ${count}" >&2; exit 1; }
echo "ok: list variables"

req -X PATCH "${BASE}/variable/${var_id}" -d '{
  "data": {"type": "variable", "id": "'"${var_id}"'", "attributes": {"value": "us-east-1"}}
}' >/dev/null
value=$(req "${BASE}/variable/${var_id}" | jq -r '.data.attributes.value')
[[ "$value" == "us-east-1" ]] || { echo "FAIL: variable patch (${value})" >&2; exit 1; }
echo "ok: patch variable"

expect_status 204 "$(status -X DELETE "${BASE}/variable/${var_id}")" "delete variable"
expect_status 404 "$(status "${BASE}/variable/${var_id}")" "get deleted variable"

# The UI's permission endpoints are mounted outside /api/v1 and answer plain
# JSON with 202. They need a bearer token even when the management API is open,
# and the development token is the instance owner: two flags and nothing else
# when no team matches.
if [[ -n "$TOKEN" ]]; then
  root="${BASE%/api/v1}"
  expect_status 202 "$(status "${root}/access-token/v1/teams/current-teams")" "current teams"
  teams=$(curl -fsS "${auth[@]}" "${root}/access-token/v1/teams/current-teams")
  [[ "$(echo "$teams" | jq -r '.groups | length')" == "0" ]] || { echo "FAIL: current teams (${teams})" >&2; exit 1; }
  echo "ok: current teams"

  expect_status 202 "$(status "${root}/access-token/v1/teams/permissions/organization/${org_id}")" "team permissions"
  perms=$(curl -fsS "${auth[@]}" "${root}/access-token/v1/teams/permissions/organization/${org_id}")
  [[ "$(echo "$perms" | jq -r '.managePermission')" == "true" &&
     "$(echo "$perms" | jq -r '.managePolicies')" == "true" &&
     "$(echo "$perms" | jq -r '.manageWorkspace')" == "false" ]] || {
    echo "FAIL: team permissions (${perms})" >&2
    exit 1
  }
  echo "ok: team permissions"
fi
expect_status 204 "$(status -X DELETE "${BASE}/workspace/${ws_id}")" "delete workspace"
expect_status 404 "$(status "${BASE}/workspace/${ws_id}")" "get deleted workspace"
expect_status 204 "$(status -X DELETE "${BASE}/project/${project_id}")" "delete project"
expect_status 204 "$(status -X DELETE "${BASE}/organization/${org_id}")" "delete organization"
expect_status 404 "$(status "${BASE}/organization/${org_id}")" "get deleted organization"

echo "smoke test passed"
