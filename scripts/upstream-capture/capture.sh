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

# Like request, but with another PAT, for responses that depend on the groups
# claim (team-derived permissions).
request_as() {
  local as="$1" name="$2" path="$3"
  shift 3
  local status
  status="$(curl -g -sS -o "$out/${name}.json" -w '%{http_code}' \
    -H "Authorization: Bearer ${as}" "$@" "${base}${path}")"
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
request organization-get-fields-empty "/api/v1/organization/${org}?fields[organization]="
request organization-get-fields-other "/api/v1/organization/${org}?fields[job]=status"
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
IFS='|' read -r job_org job_ws <<< "$(pg "select organization_id, id from workspace where deleted = false and name not like 'capture%' order by id limit 1")"
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

# Workspace relationship cases need one organization with two projects and a
# project from a different organization. The workspace is seeded with a fixed id
# so the recordings stay stable; upstream requires source, branch and a
# terraform version on the row.
IFS='|' read -r rel_org rel_p1 rel_p2 <<< "$(pg "select organization_id, min(id), max(id) from project
  group by organization_id having count(*) > 1 order by organization_id limit 1")"
rel_foreign="$(pg "select id from project where organization_id <> '${rel_org}' order by id limit 1")"
rel_ws='77777777-7777-4777-8777-777777777777'
pg "delete from job where workspace_id in (select id from workspace where name like 'capture%')" >/dev/null
pg "delete from workspace where id = '${rel_ws}' or name = 'capture-rel-workspace'" >/dev/null
pg "insert into workspace (id, name, source, branch, terraform_version, organization_id)
    values ('${rel_ws}', 'capture-rel-workspace',
            'https://github.com/terrakube-io/terrakube.git', 'main', '1.9.0', '${rel_org}')" >/dev/null

# VCS, SSH and agent rows with fixed ids so the relationship routes can point at
# something that exists; the workspace links have no foreign key on our side yet.
rel_vcs='aaaaaaaa-0000-4000-8000-000000000001'
rel_ssh='aaaaaaaa-0000-4000-8000-000000000002'
rel_agent='aaaaaaaa-0000-4000-8000-000000000003'
pg "insert into vcs (id, name, description, vcs_type, organization_id, status)
    values ('${rel_vcs}', 'capture-vcs', 'capture', 'GITHUB', '${rel_org}', 'COMPLETED') on conflict (id) do nothing" >/dev/null
pg "insert into ssh (id, name, description, organization_id)
    values ('${rel_ssh}', 'capture-ssh', 'capture', '${rel_org}') on conflict (id) do nothing" >/dev/null
pg "insert into agent (id, name, description, url, organization_id)
    values ('${rel_agent}', 'capture-agent', 'capture', 'http://agent:8090', '${rel_org}') on conflict (id) do nothing" >/dev/null

request workspace-patch-attribute "/api/v1/organization/${rel_org}/workspace/${rel_ws}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"workspace\",\"id\":\"${rel_ws}\",\"attributes\":{\"description\":\"captured\"}}}"
request workspace-patch-project "/api/v1/organization/${rel_org}/workspace/${rel_ws}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"workspace\",\"id\":\"${rel_ws}\",\"relationships\":{\"project\":{\"data\":{\"type\":\"project\",\"id\":\"${rel_p1}\"}}}}}"
request workspace-patch-project-moved "/api/v1/organization/${rel_org}/workspace/${rel_ws}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"workspace\",\"id\":\"${rel_ws}\",\"relationships\":{\"project\":{\"data\":{\"type\":\"project\",\"id\":\"${rel_p2}\"}}}}}"
request workspace-patch-project-foreign "/api/v1/organization/${rel_org}/workspace/${rel_ws}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"workspace\",\"id\":\"${rel_ws}\",\"relationships\":{\"project\":{\"data\":{\"type\":\"project\",\"id\":\"${rel_foreign}\"}}}}}"
request workspace-patch-project-unknown "/api/v1/organization/${rel_org}/workspace/${rel_ws}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"workspace\",\"id\":\"${rel_ws}\",\"relationships\":{\"project\":{\"data\":{\"type\":\"project\",\"id\":\"00000000-0000-0000-0000-000000000000\"}}}}}"
request workspace-patch-project-clear "/api/v1/organization/${rel_org}/workspace/${rel_ws}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"workspace\",\"id\":\"${rel_ws}\",\"relationships\":{\"project\":{\"data\":null}}}}"
request workspace-create-project-foreign "/api/v1/organization/${rel_org}/workspace" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"workspace\",\"attributes\":{\"name\":\"capture-rel-rejected\",\"source\":\"https://github.com/terrakube-io/terrakube.git\",\"branch\":\"main\",\"terraformVersion\":\"1.9.0\"},\"relationships\":{\"project\":{\"data\":{\"type\":\"project\",\"id\":\"${rel_foreign}\"}}}}}"

# Relationship PATCH routes the UI's settings screens call. VCS, SSH and agent
# rows were seeded above; clearing sends data:null.
request workspace-relationship-project-set "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/project" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"project\",\"id\":\"${rel_p1}\"}}"
request workspace-relationship-project-clear "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/project" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":null}'
request workspace-relationship-project-foreign "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/project" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"project\",\"id\":\"${rel_foreign}\"}}"
request workspace-relationship-project-unknown "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/project" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"project","id":"00000000-0000-0000-0000-000000000000"}}'
request workspace-relationship-vcs-set "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/vcs" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"vcs\",\"id\":\"${rel_vcs}\"}}"
request workspace-relationship-vcs-clear "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/vcs" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":null}'
request workspace-relationship-ssh-set "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/ssh" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"ssh\",\"id\":\"${rel_ssh}\"}}"
request workspace-relationship-ssh-clear "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/ssh" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":null}'
request workspace-relationship-agent-set "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/agent" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"agent\",\"id\":\"${rel_agent}\"}}"
request workspace-relationship-agent-clear "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/agent" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":null}'
request workspace-relationship-unknown "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/nope" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":null}'

# Atomic operations are how the UI creates and updates workspaces. The media
# type must match on Accept as well, or upstream answers 415.
atomic_ct='application/vnd.api+json;ext="https://jsonapi.org/ext/atomic"'
pg "delete from workspace where name = 'capture-atomic'" >/dev/null
request organization-atomic-create "/api/v1/operations" -X POST \
  -H "Content-Type: ${atomic_ct}" -H "Accept: ${atomic_ct}" \
  --data "{\"atomic:operations\":[{\"op\":\"add\",\"href\":\"/organization/${rel_org}/workspace\",\"data\":{\"type\":\"workspace\",\"lid\":\"11111111-1111-4111-8111-111111111111\",\"attributes\":{\"source\":\"https://github.com/terrakube-io/terrakube.git\",\"folder\":\"/\",\"name\":\"capture-atomic\",\"terraformVersion\":\"1.9.0\",\"branch\":\"main\",\"iacType\":\"terraform\",\"executionMode\":\"remote\"}}}]}"
request organization-atomic-update "/api/v1/operations" -X POST \
  -H "Content-Type: ${atomic_ct}" -H "Accept: ${atomic_ct}" \
  --data "{\"atomic:operations\":[{\"op\":\"update\",\"href\":\"/organization/${rel_org}/workspace/${rel_ws}\",\"data\":{\"type\":\"workspace\",\"id\":\"${rel_ws}\",\"attributes\":{\"description\":\"atomic-after\"}}}]}"
request organization-atomic-unknown "/api/v1/operations" -X POST \
  -H "Content-Type: ${atomic_ct}" -H "Accept: ${atomic_ct}" \
  --data "{\"atomic:operations\":[{\"op\":\"update\",\"href\":\"/organization/${rel_org}/workspace/00000000-0000-0000-0000-000000000000\",\"data\":{\"type\":\"workspace\",\"id\":\"00000000-0000-0000-0000-000000000000\",\"attributes\":{\"description\":\"x\"}}}]}"

# Global variables: seeded with fixed ids so the item routes can be replayed.
# Two visible rows keep the list shape identical whichever order upstream picks;
# the sensitive row is fetched on its own because its value is omitted.
gv_org="${rel_org}"
pg "delete from globalvar where variable_key like 'CAPTURE\_%'" >/dev/null
pg "insert into globalvar (id, variable_key, variable_value, variable_description, sensitive, hcl, variable_category, organization_id)
    values ('99999999-9999-4999-8999-999999999990', 'CAPTURE_ONE', 'one', 'capture', false, false, 'ENV', '${gv_org}'),
           ('99999999-9999-4999-8999-999999999992', 'CAPTURE_TWO', 'two', 'capture', false, true, 'TERRAFORM', '${gv_org}'),
           ('99999999-9999-4999-8999-999999999993', 'CAPTURE_SECRET', 'hidden', 'capture', true, false, 'ENV', '${gv_org}')" >/dev/null

request organization-globalvar-list      "/api/v1/organization/${gv_org}/globalvar"
request organization-globalvar-get       "/api/v1/organization/${gv_org}/globalvar/99999999-9999-4999-8999-999999999990"
request organization-globalvar-get-sensitive "/api/v1/organization/${gv_org}/globalvar/99999999-9999-4999-8999-999999999993"
request organization-globalvar-create    "/api/v1/organization/${gv_org}/globalvar" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"globalvar","attributes":{"key":"CAPTURE_NEW","value":"new","description":"capture","sensitive":false,"hcl":false,"category":"ENV"}}}'
request organization-globalvar-patch     "/api/v1/organization/${gv_org}/globalvar/99999999-9999-4999-8999-999999999990" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"globalvar","id":"99999999-9999-4999-8999-999999999990","attributes":{"key":"CAPTURE_ONE","value":"after","description":"after","hcl":true,"category":"TERRAFORM"}}}'
request organization-globalvar-delete    "/api/v1/organization/${gv_org}/globalvar/99999999-9999-4999-8999-999999999993" -X DELETE
request organization-globalvar-invalid-category "/api/v1/organization/${gv_org}/globalvar" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"globalvar","attributes":{"key":"CAPTURE_BAD","value":"x","description":"d","sensitive":false,"hcl":false,"category":"env"}}}'
request organization-globalvar-missing-key "/api/v1/organization/${gv_org}/globalvar" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"globalvar","attributes":{"value":"x","description":"d","sensitive":false,"hcl":false,"category":"ENV"}}}'
# Upstream's duplicate-key hook throws and Spring answers a 500; the recording
# documents the bug rather than a contract to copy.
request organization-globalvar-duplicate "/api/v1/organization/${gv_org}/globalvar" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"globalvar","attributes":{"key":"CAPTURE_ONE","value":"dup","description":"d","sensitive":false,"hcl":false,"category":"ENV"}}}'
pg "delete from globalvar where variable_key like 'CAPTURE\_%'" >/dev/null

# The action catalog. Ids are the editor's slugs, not UUIDs, and only a
# superuser may write; the demo PAT is one. The temp row is removed at the end.
pg "delete from action where id = 'capture.temp-action'" >/dev/null
request action-list              '/api/v1/action'
request action-list-filtered     "/api/v1/action?filter[action]=active==true%3Btype=in=('Workspace/ResourceDrawer/Action','Workspace/ResourceDrawer/Tab')"
request action-get               '/api/v1/action/terrakube.open-documentation'
request action-fields            '/api/v1/action?fields[action]=name,active'
request action-missing           '/api/v1/action/nope'
request action-create            '/api/v1/action' -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"action","id":"capture.temp-action","attributes":{"name":"Capture Temp","label":"Temp","category":"General","description":"capture","type":"Workspace/Action","action":"dGVzdA==","displayCriteria":"[{\"filter\":\"true\"}]","version":"1.0.0","active":true}}}'
request action-patch             '/api/v1/action/capture.temp-action' -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"action","id":"capture.temp-action","attributes":{"label":"Temp after","active":false}}}'
request action-delete            '/api/v1/action/capture.temp-action' -X DELETE
pg "delete from action where id = 'capture.temp-action'" >/dev/null

# Variable collections: Settings > Collections, plus the read-only table the
# workspace page fills from a reference. rel_org/rel_ws carry a team with
# manage_collection (upstream gates collection writes on a team permission).
coll_org="${rel_org}"
coll_ws="${rel_ws}"
pg "delete from item where collection_id in (select id from collection where name like 'capture-collection%');
    delete from reference where collection_id in (select id from collection where name like 'capture-collection%');
    delete from collection where name like 'capture-collection%'" >/dev/null
request organization-collection-create "/api/v1/organization/${coll_org}/collection" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"collection","attributes":{"name":"capture-collection","description":"capture","priority":10}}}'
coll_id="$(python3 -c "import json;print(json.load(open('$out/organization-collection-create.json'))['data']['id'])")"
request organization-collection-list   "/api/v1/organization/${coll_org}/collection"
request organization-collection-get    "/api/v1/organization/${coll_org}/collection/${coll_id}"
request organization-collection-patch  "/api/v1/organization/${coll_org}/collection/${coll_id}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"collection\",\"id\":\"${coll_id}\",\"attributes\":{\"name\":\"capture-collection-after\"}}}"
request organization-collection-item-create "/api/v1/organization/${coll_org}/collection/${coll_id}/item" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"item","attributes":{"key":"CAPTURE_KEY","value":"one","description":"capture","sensitive":false,"hcl":false,"category":"ENV"}}}'
coll_item_id="$(python3 -c "import json;print(json.load(open('$out/organization-collection-item-create.json'))['data']['id'])")"
request organization-collection-item-list "/api/v1/organization/${coll_org}/collection/${coll_id}/item"
request organization-collection-item-patch "/api/v1/organization/${coll_org}/collection/${coll_id}/item/${coll_item_id}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"item\",\"id\":\"${coll_item_id}\",\"attributes\":{\"value\":\"after\"}}}"
request organization-collection-item-invalid-category "/api/v1/organization/${coll_org}/collection/${coll_id}/item" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"item","attributes":{"key":"CAPTURE_BAD","value":"x","category":"env"}}}'
request organization-collection-item-missing-key "/api/v1/organization/${coll_org}/collection/${coll_id}/item" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"item","attributes":{"value":"x","category":"ENV"}}}'
request organization-collection-create-missing-name "/api/v1/organization/${coll_org}/collection" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"collection","attributes":{"priority":5}}}'
request organization-collection-missing "/api/v1/organization/${coll_org}/collection/00000000-0000-4000-8000-000000000000"
request organization-collection-reference-create "/api/v1/organization/${coll_org}/collection/${coll_id}/reference" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"reference\",\"attributes\":{\"description\":\"capture\"},\"relationships\":{\"workspace\":{\"data\":{\"type\":\"workspace\",\"id\":\"${coll_ws}\"}}}}}"
coll_ref_id="$(python3 -c "import json;print(json.load(open('$out/organization-collection-reference-create.json'))['data']['id'])")"
request organization-collection-reference-list "/api/v1/organization/${coll_org}/collection/${coll_id}/reference"
request reference-collection "/api/v1/reference/${coll_ref_id}/collection?include=item"
request organization-workspace-reference "/api/v1/organization/${coll_org}/workspace/${coll_ws}?include=reference"
request organization-collection-item-delete "/api/v1/organization/${coll_org}/collection/${coll_id}/item/${coll_item_id}" -X DELETE
request organization-collection-reference-delete "/api/v1/organization/${coll_org}/collection/${coll_id}/reference/${coll_ref_id}" -X DELETE
request organization-collection-delete "/api/v1/organization/${coll_org}/collection/${coll_id}" -X DELETE

# SSH keys and agent pools: Settings > SSH Keys / Agents and the workspace
# pickers. Upstream requires no attribute beyond the type, so the missing-name
# create is recorded to keep that permissiveness honest.
request organization-ssh-list   "/api/v1/organization/${rel_org}/ssh"
request organization-ssh-get    "/api/v1/organization/${rel_org}/ssh/${rel_ssh}"
request organization-ssh-create "/api/v1/organization/${rel_org}/ssh" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"ssh","attributes":{"name":"capture-ssh-created","description":"capture","sshType":"rsa","privateKey":"-----BEGIN PRIVATE KEY-----\ncapture\n-----END PRIVATE KEY-----"}}}'
ssh_new="$(python3 -c "import json;print(json.load(open('$out/organization-ssh-create.json'))['data']['id'])")"
request organization-ssh-patch  "/api/v1/organization/${rel_org}/ssh/${ssh_new}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"ssh\",\"id\":\"${ssh_new}\",\"attributes\":{\"description\":\"after\"}}}"
request organization-ssh-create-missing-name "/api/v1/organization/${rel_org}/ssh" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"ssh","attributes":{"sshType":"rsa"}}}'
ssh_unnamed="$(python3 -c "import json;print(json.load(open('$out/organization-ssh-create-missing-name.json'))['data']['id'])")"
request organization-ssh-missing "/api/v1/organization/${rel_org}/ssh/00000000-0000-4000-8000-000000000099"
request organization-ssh-delete "/api/v1/organization/${rel_org}/ssh/${ssh_new}" -X DELETE
request organization-ssh-delete-unnamed "/api/v1/organization/${rel_org}/ssh/${ssh_unnamed}" -X DELETE

request organization-agent-list   "/api/v1/organization/${rel_org}/agent"
request organization-agent-create "/api/v1/organization/${rel_org}/agent" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"agent","attributes":{"name":"capture-agent-created","description":"capture","url":"http://capture:8090"}}}'
agent_new="$(python3 -c "import json;print(json.load(open('$out/organization-agent-create.json'))['data']['id'])")"
request organization-agent-delete "/api/v1/organization/${rel_org}/agent/${agent_new}" -X DELETE

# A key or pool that a workspace still points at is refused with Elide's
# data-integrity 409, so the link is recorded too and cleared afterwards.
request organization-ssh-link-in-use "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/ssh" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"ssh\",\"id\":\"${rel_ssh}\"}}"
request organization-ssh-delete-in-use "/api/v1/organization/${rel_org}/ssh/${rel_ssh}" -X DELETE
request organization-ssh-unlink-in-use "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/ssh" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":null}'
request organization-agent-link-in-use "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/agent" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"agent\",\"id\":\"${rel_agent}\"}}"
request organization-agent-delete-in-use "/api/v1/organization/${rel_org}/agent/${rel_agent}" -X DELETE
request organization-agent-unlink-in-use "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/agent" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":null}'

# VCS providers: Settings > VCS and the workspace source picker. The seeded
# vcs row is used for the link/include/delete-in-use captures, then removed by
# the delete itself (the cleanup at the end tolerates the missing row).
request organization-vcs-list "/api/v1/organization/${rel_org}/vcs"
request organization-vcs-create "/api/v1/organization/${rel_org}/vcs" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"vcs","attributes":{"name":"capture-vcs-created","description":"capture","vcsType":"GITLAB","connectionType":"STANDALONE","clientId":"client","clientSecret":"secret","privateKey":"key","callback":"callback","endpoint":"https://gitlab.example.com","apiUrl":"https://gitlab.example.com/api/v4","redirectUrl":"https://ui.example.com/settings/vcs","status":"PENDING"}}}'
vcs_new="$(python3 -c "import json;print(json.load(open('$out/organization-vcs-create.json'))['data']['id'])")"
request organization-vcs-get "/api/v1/organization/${rel_org}/vcs/${vcs_new}"
request organization-vcs-patch "/api/v1/organization/${rel_org}/vcs/${vcs_new}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"vcs\",\"id\":\"${vcs_new}\",\"attributes\":{\"description\":\"after\",\"status\":\"COMPLETED\"}}}"
request organization-vcs-create-missing-required "/api/v1/organization/${rel_org}/vcs" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"vcs","attributes":{"name":"capture-vcs-incomplete"}}}'
request organization-vcs-create-invalid-type "/api/v1/organization/${rel_org}/vcs" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"vcs","attributes":{"name":"capture-vcs-bad","description":"capture","vcsType":"NOPE"}}}'
request organization-vcs-missing "/api/v1/organization/${rel_org}/vcs/00000000-0000-4000-8000-000000000095"
request workspace-relationship-vcs-link-in-use "/api/v1/organization/${rel_org}/workspace/${rel_ws}/relationships/vcs" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"vcs\",\"id\":\"${rel_vcs}\"}}"
request workspace-include-vcs "/api/v1/organization/${rel_org}/workspace/${rel_ws}?include=vcs"
request organization-vcs-delete-in-use "/api/v1/organization/${rel_org}/vcs/${rel_vcs}" -X DELETE
request workspace-get-after-vcs-delete "/api/v1/organization/${rel_org}/workspace/${rel_ws}"
request organization-vcs-delete "/api/v1/organization/${rel_org}/vcs/${vcs_new}" -X DELETE
request organization-vcs-get-after-delete "/api/v1/organization/${rel_org}/vcs/${vcs_new}"

# Teams: Settings > Teams. Permissions are booleans on the row, so the create
# also records the primitive defaults upstream applies.
pg "delete from team where name like 'capture-team%' and organization_id = '${rel_org}'" >/dev/null
request organization-team-list "/api/v1/organization/${rel_org}/team"
request organization-team-create "/api/v1/organization/${rel_org}/team" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"team","attributes":{"name":"capture-team","manageWorkspace":true,"manageJob":true,"planJob":true,"approveJob":true,"role":"write"}}}'
team_new="$(python3 -c "import json;print(json.load(open('$out/organization-team-create.json'))['data']['id'])")"
request organization-team-get "/api/v1/organization/${rel_org}/team/${team_new}"

# Permission endpoints the UI calls from every screen. They sit outside
# /api/v1 and answer plain JSON with 202, not JSON:API. The group token is
# minted while the team still has role=write, so the team answers show a
# non-admin row and no instance-owner flags.
group_token="$(PAT_GROUPS='capture-team' "$here/pat.py")"
request access-token-current-teams '/access-token/v1/teams/current-teams'
request access-token-permissions-organization "/access-token/v1/teams/permissions/organization/${rel_org}"
request access-token-permissions-organization-unknown '/access-token/v1/teams/permissions/organization/00000000-0000-4000-8000-0000000000ff'
request access-token-permissions-organization-bad-id '/access-token/v1/teams/permissions/organization/not-a-uuid'
request_as "$group_token" access-token-permissions-organization-team "/access-token/v1/teams/permissions/organization/${rel_org}"
request_as "$group_token" access-token-permissions-workspace-team "/access-token/v1/teams/permissions/organization/${rel_org}/workspace/${rel_ws}"
request access-token-permissions-workspace "/access-token/v1/teams/permissions/organization/${rel_org}/workspace/${rel_ws}"
request access-token-permissions-workspace-missing "/access-token/v1/teams/permissions/organization/${rel_org}/workspace/00000000-0000-4000-8000-0000000000ff"
unauthenticated access-token-no-token '/access-token/v1/teams/current-teams'

request organization-team-patch "/api/v1/organization/${rel_org}/team/${team_new}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"team\",\"id\":\"${team_new}\",\"attributes\":{\"role\":\"admin\",\"manageVcs\":true,\"destroyJob\":false}}}"
request organization-team-create-missing-name "/api/v1/organization/${rel_org}/team" -X POST \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":{"type":"team","attributes":{"role":"read"}}}'
request organization-team-missing "/api/v1/organization/${rel_org}/team/00000000-0000-4000-8000-000000000097"
request organization-team-bad-id "/api/v1/organization/${rel_org}/team/not-a-uuid"
request organization-team-delete "/api/v1/organization/${rel_org}/team/${team_new}" -X DELETE
request organization-team-get-after-delete "/api/v1/organization/${rel_org}/team/${team_new}"

# A read role proves the all-false half of the matrix without the owner group.
# Setup and cleanup are not captures themselves.
reader_team="$(curl -sS -H "Authorization: Bearer ${token}" -H 'Content-Type: application/vnd.api+json' \
  -X POST "${base}/api/v1/organization/${rel_org}/team" \
  --data '{"data":{"type":"team","attributes":{"name":"capture-reader","role":"read"}}}' \
  | python3 -c "import json,sys;print(json.load(sys.stdin)['data']['id'])")"
reader_token="$(PAT_GROUPS='capture-reader' "$here/pat.py")"
request_as "$reader_token" access-token-permissions-organization-reader "/access-token/v1/teams/permissions/organization/${rel_org}"
curl -sS -H "Authorization: Bearer ${token}" -X DELETE "${base}/api/v1/organization/${rel_org}/team/${reader_team}" >/dev/null

# Workspace and project access grants: the team pickers on the workspace and
# project screens. `access` and `project_access` are relationship resources
# (rootLevel=false), so every path goes through the parent. A create that omits
# destroyJob records upstream's primitive default of true.
pg "delete from access where name = 'capture-group'; delete from project_access where name = 'capture-group'" >/dev/null
request workspace-access-list "/api/v1/organization/${rel_org}/workspace/${rel_ws}/access"
request workspace-access-create "/api/v1/organization/${rel_org}/workspace/${rel_ws}/access" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"access","attributes":{"name":"capture-group","manageState":true,"manageWorkspace":true,"planJob":true,"approveJob":false,"role":"custom"}}}'
access_new="$(python3 -c "import json;print(json.load(open('$out/workspace-access-create.json'))['data']['id'])")"
request workspace-access-get "/api/v1/organization/${rel_org}/workspace/${rel_ws}/access/${access_new}"
request workspace-access-patch "/api/v1/organization/${rel_org}/workspace/${rel_ws}/access/${access_new}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"access\",\"id\":\"${access_new}\",\"attributes\":{\"role\":\"read\"}}}"
request workspace-access-get-after-patch "/api/v1/organization/${rel_org}/workspace/${rel_ws}/access/${access_new}"
request workspace-access-delete "/api/v1/organization/${rel_org}/workspace/${rel_ws}/access/${access_new}" -X DELETE
request workspace-access-get-after-delete "/api/v1/organization/${rel_org}/workspace/${rel_ws}/access/${access_new}"
request project-access-list "/api/v1/organization/${rel_org}/project/${rel_p1}/projectAccess"
request project-access-create "/api/v1/organization/${rel_org}/project/${rel_p1}/projectAccess" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"project_access","attributes":{"name":"capture-group","manageState":true,"manageWorkspace":true,"planJob":false,"approveJob":false,"role":"custom"}}}'
project_access_new="$(python3 -c "import json;print(json.load(open('$out/project-access-create.json'))['data']['id'])")"
request project-access-get "/api/v1/organization/${rel_org}/project/${rel_p1}/projectAccess/${project_access_new}"
request project-access-patch "/api/v1/organization/${rel_org}/project/${rel_p1}/projectAccess/${project_access_new}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"project_access\",\"id\":\"${project_access_new}\",\"attributes\":{\"role\":\"plan\"}}}"
request project-access-delete "/api/v1/organization/${rel_org}/project/${rel_p1}/projectAccess/${project_access_new}" -X DELETE
request project-access-get-after-delete "/api/v1/organization/${rel_org}/project/${rel_p1}/projectAccess/${project_access_new}"

# Templates: Settings > Templates and the workspace default template picker.
# Upstream requires no attribute, so the empty create is part of the contract.
pg "delete from template where name like 'capture-template%' and organization_id = '${rel_org}'" >/dev/null
request organization-template-list "/api/v1/organization/${rel_org}/template"
request organization-template-create "/api/v1/organization/${rel_org}/template" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"template","attributes":{"name":"capture-template","description":"capture","version":"9.9.9","tcl":"ZmxvdzoK"}}}'
template_new="$(python3 -c "import json;print(json.load(open('$out/organization-template-create.json'))['data']['id'])")"
request organization-template-get "/api/v1/organization/${rel_org}/template/${template_new}"
request organization-template-patch "/api/v1/organization/${rel_org}/template/${template_new}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"template\",\"id\":\"${template_new}\",\"attributes\":{\"description\":\"after\",\"version\":\"10.0.0\",\"tcl\":\"ZmxvdzoKICAtIHR5cGU6ICJ0ZXJyYWZvcm1QbGFuIg==\"}}}"
request organization-template-create-empty "/api/v1/organization/${rel_org}/template" -X POST \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":{"type":"template","attributes":{}}}'
template_empty="$(python3 -c "import json;print(json.load(open('$out/organization-template-create-empty.json'))['data']['id'])")"
request organization-template-missing "/api/v1/organization/${rel_org}/template/00000000-0000-4000-8000-000000000096"
request organization-template-bad-id "/api/v1/organization/${rel_org}/template/not-a-uuid"
request organization-template-delete "/api/v1/organization/${rel_org}/template/${template_new}" -X DELETE
request organization-template-get-after-delete "/api/v1/organization/${rel_org}/template/${template_new}"
request organization-template-delete-empty "/api/v1/organization/${rel_org}/template/${template_empty}" -X DELETE

# Notifications: Settings > Notifications. A configuration is org-wide (no
# workspace) or workspace-scoped, gets triggers per job status, and may narrow
# itself to a set of templates.
notify_template='00000000-0000-4000-8000-000000000095'
pg "delete from notification_configuration where name like 'capture-notify%' and organization_id = '${rel_org}';
    insert into template (id, name, organization_id)
    values ('${notify_template}', 'capture-notify-template', '${rel_org}') on conflict (id) do nothing" >/dev/null
request organization-notification-create "/api/v1/organization/${rel_org}/notificationConfiguration" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"notification_configuration","attributes":{"name":"capture-notify","description":"capture","channelType":"WEBHOOK","destinationUrl":"https://example.com/hook","signingSecret":"capture-secret","active":true,"messageStyle":"SIMPLE"}}}'
notify_new="$(python3 -c "import json;print(json.load(open('$out/organization-notification-create.json'))['data']['id'])")"
request organization-notification-list "/api/v1/organization/${rel_org}/notificationConfiguration"
request organization-notification-get "/api/v1/organization/${rel_org}/notificationConfiguration/${notify_new}"
request organization-notification-patch "/api/v1/notification_configuration/${notify_new}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"notification_configuration\",\"id\":\"${notify_new}\",\"attributes\":{\"description\":null,\"active\":false,\"messageStyle\":\"DETAILED\"}}}"
request organization-notification-trigger-create "/api/v1/notification_configuration/${notify_new}/triggers" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"notification_trigger","attributes":{"jobStatus":"failed"}}}'
notify_trigger="$(python3 -c "import json;print(json.load(open('$out/organization-notification-trigger-create.json'))['data']['id'])")"
request organization-notification-trigger-missing "/api/v1/notification_configuration/${notify_new}/triggers" -X POST \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":{"type":"notification_trigger","attributes":{}}}'
request organization-notification-trigger-invalid "/api/v1/notification_configuration/${notify_new}/triggers" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"notification_trigger","attributes":{"jobStatus":"nope"}}}'
request organization-notification-templates "/api/v1/notification_configuration/${notify_new}/relationships/templates" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":[{\"type\":\"template\",\"id\":\"${notify_template}\"}]}"
request organization-notification-templates-unknown "/api/v1/notification_configuration/${notify_new}/relationships/templates" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":[{"type":"template","id":"00000000-0000-4000-8000-000000000099"}]}'
request organization-notification-workspace-create "/api/v1/organization/${rel_org}/workspace/${rel_ws}/notificationConfiguration" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"notification_configuration","attributes":{"name":"capture-notify-workspace","channelType":"TEAMS","destinationUrl":"https://teams.example.com/hook"}}}'
notify_ws="$(python3 -c "import json;print(json.load(open('$out/organization-notification-workspace-create.json'))['data']['id'])")"
request organization-notification-workspace-list "/api/v1/organization/${rel_org}/workspace/${rel_ws}/notificationConfiguration"
request organization-notification-create-missing-name "/api/v1/organization/${rel_org}/notificationConfiguration" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"notification_configuration","attributes":{"channelType":"WEBHOOK","destinationUrl":"https://example.com/x"}}}'
request organization-notification-create-invalid "/api/v1/organization/${rel_org}/notificationConfiguration" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"notification_configuration","attributes":{"name":"capture-notify-bad","channelType":"NOPE","destinationUrl":"https://example.com/x"}}}'
request organization-notification-missing "/api/v1/organization/${rel_org}/notificationConfiguration/00000000-0000-4000-8000-000000000098"
request organization-notification-bad-id "/api/v1/organization/${rel_org}/notificationConfiguration/not-a-uuid"
request organization-notification-trigger-delete "/api/v1/notification_configuration/${notify_new}/triggers/${notify_trigger}" -X DELETE
request organization-notification-workspace-delete "/api/v1/notification_configuration/${notify_ws}" -X DELETE
request organization-notification-delete "/api/v1/notification_configuration/${notify_new}" -X DELETE
request organization-notification-get-after-delete "/api/v1/organization/${rel_org}/notificationConfiguration/${notify_new}"
pg "delete from notification_configuration where name like 'capture-notify%' and organization_id = '${rel_org}';
    delete from template where id = '${notify_template}'" >/dev/null

# Tags: Settings > Tags and the workspace tag editor. The tag delete cascades
# its bindings, so the workspace rows go with it.
pg "delete from workspacetag where tag_id in (select id from tag where name like 'capture-tag%');
    delete from tag where name like 'capture-tag%' and organization_id = '${rel_org}'" >/dev/null
request organization-tag-list "/api/v1/organization/${rel_org}/tag"
request organization-tag-create "/api/v1/organization/${rel_org}/tag" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"tag","attributes":{"name":"capture-tag"}}}'
tag_new="$(python3 -c "import json;print(json.load(open('$out/organization-tag-create.json'))['data']['id'])")"
request organization-tag-get "/api/v1/organization/${rel_org}/tag/${tag_new}"
request organization-tag-patch "/api/v1/organization/${rel_org}/tag/${tag_new}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"tag\",\"id\":\"${tag_new}\",\"attributes\":{\"name\":\"capture-tag-after\"}}}"
request organization-tag-filter "/api/v1/organization/${rel_org}/tag?filter[tag]=name==capture-tag-after"
request organization-tag-create-missing-name "/api/v1/organization/${rel_org}/tag" -X POST \
  -H 'Content-Type: application/vnd.api+json' --data '{"data":{"type":"tag","attributes":{}}}'
tag_unnamed="$(python3 -c "import json;print(json.load(open('$out/organization-tag-create-missing-name.json'))['data']['id'])")"
request organization-tag-missing "/api/v1/organization/${rel_org}/tag/00000000-0000-4000-8000-000000000097"
request workspace-tag-list "/api/v1/organization/${rel_org}/workspace/${rel_ws}/workspaceTag"
request workspace-tag-create "/api/v1/organization/${rel_org}/workspace/${rel_ws}/workspaceTag" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"workspacetag\",\"attributes\":{\"tagId\":\"${tag_new}\",\"value\":\"prod\"}}}"
workspace_tag="$(python3 -c "import json;print(json.load(open('$out/workspace-tag-create.json'))['data']['id'])")"
request workspace-tag-list-after "/api/v1/organization/${rel_org}/workspace/${rel_ws}/workspaceTag"
request workspace-tag-patch "/api/v1/organization/${rel_org}/workspace/${rel_ws}/workspaceTag/${workspace_tag}" -X PATCH \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"workspacetag\",\"id\":\"${workspace_tag}\",\"attributes\":{\"value\":null}}}"
request workspace-tag-missing-tag "/api/v1/organization/${rel_org}/workspace/${rel_ws}/workspaceTag" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"workspacetag","attributes":{"value":"x"}}}'
request workspace-tag-unknown-tag "/api/v1/organization/${rel_org}/workspace/${rel_ws}/workspaceTag" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"workspacetag","attributes":{"tagId":"00000000-0000-4000-8000-000000000099","value":"x"}}}'
request workspace-tag-duplicate "/api/v1/organization/${rel_org}/workspace/${rel_ws}/workspaceTag" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data "{\"data\":{\"type\":\"workspacetag\",\"attributes\":{\"tagId\":\"${tag_new}\",\"value\":\"again\"}}}"
request workspace-tag-delete "/api/v1/organization/${rel_org}/workspace/${rel_ws}/workspaceTag/${workspace_tag}" -X DELETE
request organization-tag-duplicate "/api/v1/organization/${rel_org}/tag" -X POST \
  -H 'Content-Type: application/vnd.api+json' \
  --data '{"data":{"type":"tag","attributes":{"name":"capture-tag-after"}}}'
request organization-tag-delete "/api/v1/organization/${rel_org}/tag/${tag_new}" -X DELETE
request organization-tag-delete-unnamed "/api/v1/organization/${rel_org}/tag/${tag_unnamed}" -X DELETE

# Leave the disposable database as it was: these rows exist only for captures.
pg "delete from vcs where id = '${rel_vcs}'; delete from ssh where id = '${rel_ssh}'; delete from agent where id = '${rel_agent}'" >/dev/null
pg "delete from template where name like 'capture-template%' and organization_id = '${rel_org}'" >/dev/null
pg "delete from team where name like 'capture-team%' and organization_id = '${rel_org}'" >/dev/null
pg "delete from workspace where name in ('capture-atomic', 'capture-rel-workspace')" >/dev/null
pg "delete from job where commit_id = 'deaDBEEF'" >/dev/null

echo "wrote $(wc -l < "$out/index.tsv") captures to $out"
