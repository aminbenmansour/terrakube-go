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

# Leave the disposable database as it was: these rows exist only for captures.
pg "delete from vcs where id = '${rel_vcs}'; delete from ssh where id = '${rel_ssh}'; delete from agent where id = '${rel_agent}'" >/dev/null
pg "delete from workspace where name in ('capture-atomic', 'capture-rel-workspace')" >/dev/null
pg "delete from job where commit_id = 'deaDBEEF'" >/dev/null

echo "wrote $(wc -l < "$out/index.tsv") captures to $out"
