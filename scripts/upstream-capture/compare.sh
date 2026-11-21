#!/usr/bin/env bash
#
# Sends the captured requests to this API and to upstream, and prints the
# difference in shape (not values). Needs the stack from start.sh plus a running
# terrakube-go API started with TERRAKUBE_AUTH_TOKEN=dev-token.
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
captures="${CAPTURES:-$here/../../docs/upstream/rest}"
ours="${TERRAKUBE_API_URL:-http://localhost:8080}"
upstream="${TK_UPSTREAM_URL:-http://localhost:18080}"
upstream_token="$(cat "${TK_UPSTREAM_TOKEN:-$here/token.txt}")"

shape() {
  python3 -c '
import json,sys
try: d=json.load(sys.stdin)
except Exception: print("<not json>"); raise SystemExit
def walk(v):
    if isinstance(v,dict): return {k: walk(x) for k,x in sorted(v.items())}
    if isinstance(v,list): return [walk(v[0])] if v else []
    return type(v).__name__
print(json.dumps(walk(d),sort_keys=True))
'
}

while IFS=$'\t' read -r name status path; do
  echo "=== ${name} (captured ${status}) ${path}"
  ours_status="$(curl -g -s -o /tmp/tk-ours.json -w '%{http_code}' -H 'Authorization: Bearer dev-token' "${ours}${path}")"
  theirs_status="$(curl -g -s -o /tmp/tk-theirs.json -w '%{http_code}' -H "Authorization: Bearer ${upstream_token}" "${upstream}${path}")"
  echo "status: ours=${ours_status} upstream=${theirs_status}"
  echo "ours:     $(shape < /tmp/tk-ours.json)"
  echo "upstream: $(shape < /tmp/tk-theirs.json)"
done < "$captures/index.tsv"
