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

# The demo seed has no runs and the run page is a real screen, so seed one job
# and one step in the disposable upstream database. The ids are read back rather
# than guessed.
pg() { podman exec -i tk-upstream-postgres psql -p 5433 -U terrakube -d terrakubedb -tAq -c "$1"; }
IFS='|' read -r job_org job_ws <<< "$(pg "select organization_id, id from workspace where deleted = false order by id limit 1")"
# Repeatable: drop the row a previous run seeded, or its NULL booleans poison
# every organization read (steps go with it via ON DELETE CASCADE).
pg "delete from job where commit_id = 'deaDBEEF'" >/dev/null
# Every primitive boolean must be set: the Java entity maps them to primitives,
# so a NULL column turns the whole job list into a 500.
job_id="$(pg "insert into job (status, organization_id, workspace_id, created_date, updated_date,
                               created_by, commit_id, via, tcl, output, template_reference,
                               refresh, refresh_only, auto_apply, plan_changes, pr_apply_enabled,
                               variables_as_arguments, deleted, cascade_depth, dispatch_failure_count)
              values ('completed', '${job_org}', '${job_ws}', now(), now(),
                      'admin@example.com', 'deaDBEEF', 'UI', '{}', 'ok',
                      (select id from template order by id limit 1),
                      false, false, false, true, false,
                      false, false, 0, 0)
              returning id")"
pg "insert into step (id, step_number, status, job_id, name, output)
    values ('66666666-6666-4666-8666-666666666666', 1, 'completed', ${job_id}, 'Plan', 'ok')" >/dev/null

request organization-job-list         "/api/v1/organization/${job_org}/job"
request organization-job-get          "/api/v1/organization/${job_org}/job/${job_id}"
request organization-job-include      "/api/v1/organization/${job_org}/job/${job_id}?include=step,workspace"
request organization-job-patch        "/api/v1/organization/${job_org}/job/${job_id}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"job\",\"id\":\"${job_id}\",\"attributes\":{\"status\":\"cancelled\"}}}"
request organization-job-patch-invalid "/api/v1/organization/${job_org}/job/${job_id}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"job\",\"id\":\"${job_id}\",\"attributes\":{\"status\":\"nonsense\"}}}"

echo "wrote $(wc -l < "$out/index.tsv") captures to $out"
