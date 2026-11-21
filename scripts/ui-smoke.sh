#!/usr/bin/env bash
#
# Checks that the API serves a UI build without shadowing the API itself.
# Requires a running api started with TERRAKUBE_UI_DIR set to a built UI.
#
#   TERRAKUBE_UI_DIR=/path/to/build make run-api
#   TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/ui-smoke.sh
#
# The build can come from terrakube-ui (the pinned pipeline) or from
# `yarn build` in an upstream checkout.
set -euo pipefail

BASE="${TERRAKUBE_API_URL:-http://localhost:8080}"
BASE="${BASE%/}"
UI_DIR="${TERRAKUBE_UI_DIR:-}"
TOKEN="${TERRAKUBE_AUTH_TOKEN:-}"

if [[ -z "$UI_DIR" || ! -f "${UI_DIR}/index.html" ]]; then
  echo "set TERRAKUBE_UI_DIR to a built UI directory (it must contain index.html)" >&2
  exit 2
fi

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

# The API guards /api/v1 when TERRAKUBE_AUTH_TOKEN is set; the UI's own files
# need no token, so only the API assertions use it.
auth=()
if [[ -n "$TOKEN" ]]; then
  auth=(-H "Authorization: Bearer ${TOKEN}")
fi

body() { curl -fsS "$@"; }
status() { curl -s -o /dev/null -w '%{http_code}' "$@"; }
header() { curl -sS -o /dev/null -D - "$@" | tr -d '\r'; }

echo "Waiting for ${BASE}..."
for _ in $(seq 1 30); do
  if curl -fsS "${BASE}/healthz" >/dev/null 2>&1; then break; fi
  sleep 0.5
done

echo "The UI is served at the root"
index="$(body "${BASE}/")" || fail "/ did not return 200"
grep -qi '<!doctype html' <<<"$index" || fail "/ is not HTML: ${index:0:80}"
grep -q 'env-config.js' <<<"$index" || fail "index.html does not load env-config.js"

echo "Deep links return the SPA shell, not a 404"
deep="$(body "${BASE}/organizations/some-org/workspaces")" || fail "deep link did not return 200"
grep -qi '<!doctype html' <<<"$deep" || fail "deep link did not return index.html"

echo "env-config.js is generated and never cached"
config="$(body "${BASE}/env-config.js")" || fail "/env-config.js did not return 200"
grep -q 'window._env_' <<<"$config" || fail "env-config.js has no window._env_"
grep -q 'REACT_APP_TERRAKUBE_API_URL' <<<"$config" || fail "env-config.js has no API URL"
header "${BASE}/env-config.js" | grep -qi '^cache-control: .*no-store' || fail "env-config.js is cacheable; a stale copy points at the wrong API"

echo "Fingerprinted assets are immutable, missing ones 404"
asset="$(cd "$UI_DIR" && find assets -type f -name '*.js' | head -1 || true)"
if [[ -n "$asset" ]]; then
  header "${BASE}/${asset}" | grep -qi 'immutable' || fail "/${asset} is not marked immutable"
fi
code="$(status "${BASE}/assets/does-not-exist.js")"
[[ "$code" == "404" ]] || fail "missing asset answered ${code}, want 404"

echo "API paths still reach their own handlers"
code="$(status "${auth[@]}" "${BASE}/api/v1/organization")"
[[ "$code" != "404" ]] || fail "/api/v1/organization was shadowed by the UI"
body "${auth[@]}" "${BASE}/api/v1/organization" | head -c 200 | grep -q '"data"' || fail "/api/v1/organization did not return JSON:API"
code="$(status "${BASE}/healthz")"
[[ "$code" == "200" ]] || fail "/healthz answered ${code}"

echo "OK: the UI is served and the API is not shadowed"
