#!/usr/bin/env bash
#
# Captures upstream JSON:API responses verbatim, so our API is compared with the
# real thing instead of with a reading of the OpenAPI spec.
#
#   ./capture.sh [outdir]        (default: ../../docs/upstream/rest)
#
# Needs token.txt from pat.py, and the stack from start.sh.
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
out="${1:-$here/../../docs/upstream/rest}"
token="$(cat "$here/token.txt")"
base="${TK_UPSTREAM_URL:-http://localhost:18080}"
mkdir -p "$out"

# request <name> <path> [extra curl args...]
request() {
  local name="$1" path="$2"
  shift 2
  local status
  status="$(curl -g -sS -o "$out/${name}.json" -w '%{http_code}' \
    -H "Authorization: Bearer ${token}" "$@" "${base}${path}")"
  printf '%s\t%s\t%s\n' "$name" "$status" "$path" >> "$out/index.tsv"
}

# Without a token, to record what an unauthenticated caller sees.
unauthenticated() {
  local name="$1" path="$2"
  local status
  status="$(curl -g -sS -o "$out/${name}.json" -w '%{http_code}' "${base}${path}")"
  printf '%s\t%s\t%s\n' "$name" "$status" "$path" >> "$out/index.tsv"
}

: > "$out/index.tsv"

request organization-list            '/api/v1/organization'
request organization-list-totals     '/api/v1/organization?page[totals]=true'
request organization-list-paged      '/api/v1/organization?page[size]=2&page[number]=2'
request organization-list-sorted     '/api/v1/organization?sort=-name'
request organization-list-filtered   '/api/v1/organization?filter[organization]=name==azure'
request organization-list-fields     '/api/v1/organization?fields[organization]=name'
request organization-list-include    '/api/v1/organization?include=project'

org="$(python3 -c "import json;print(json.load(open('$out/organization-list.json'))['data'][0]['id'])")"
request organization-get             "/api/v1/organization/${org}"
request organization-get-fields      "/api/v1/organization/${org}?fields[organization]=name,description"
request organization-relationships   "/api/v1/organization/${org}/relationships/project"
request organization-projects        "/api/v1/organization/${org}/project"

# Elide paginates this collection, unlike organization: the two recordings are
# the reference for meta.page with and without page[totals].
request workspace-list-paged         '/api/v1/workspace?page[size]=1&page[number]=1'
request workspace-list-totals        '/api/v1/workspace?page[totals]=true'

request organization-missing         '/api/v1/organization/00000000-0000-0000-0000-000000000000'
request organization-bad-id          '/api/v1/organization/not-a-uuid'
unauthenticated organization-no-token '/api/v1/organization'

request delete-method-not-allowed    '/api/v1/organization' -X PATCH \
  -H 'Content-Type: application/vnd.api+json' --data '{}'

echo "wrote $(wc -l < "$out/index.tsv") captures to $out"
