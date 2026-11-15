#!/usr/bin/env bash
# Exercise the /tfstate/v1 HTTP backend against a running API.
#
# Required:
#   TERRAKUBE_INTERNAL_TOKEN_SECRET  base64url signing key used by the API
#   TERRAKUBE_STATE_ORG_ID           organization UUID
#   TERRAKUBE_STATE_WORKSPACE_ID     workspace UUID
#   TERRAKUBE_STATE_JOB_ID           job id
#   TERRAKUBE_STATE_STEP_ID          active step UUID
#
# Optional:
#   TERRAKUBE_API_URL                default http://localhost:8080
#
# Use a dedicated workspace: the test overwrites and deletes its current state.
set -euo pipefail

api_url="${TERRAKUBE_API_URL:-http://localhost:8080}"
secret="${TERRAKUBE_INTERNAL_TOKEN_SECRET:?set TERRAKUBE_INTERNAL_TOKEN_SECRET}"
org="${TERRAKUBE_STATE_ORG_ID:?set TERRAKUBE_STATE_ORG_ID}"
workspace="${TERRAKUBE_STATE_WORKSPACE_ID:?set TERRAKUBE_STATE_WORKSPACE_ID}"
job="${TERRAKUBE_STATE_JOB_ID:?set TERRAKUBE_STATE_JOB_ID}"
step="${TERRAKUBE_STATE_STEP_ID:?set TERRAKUBE_STATE_STEP_ID}"

token="$(SECRET="$secret" ORG="$org" WORKSPACE="$workspace" JOB="$job" STEP="$step" python3 - <<'PY'
import base64, hashlib, hmac, json, os, time

secret = os.environ["SECRET"]
key = base64.urlsafe_b64decode(secret + "=" * (-len(secret) % 4))
if len(key) >= 64:
    alg, digest = "HS512", hashlib.sha512
elif len(key) >= 48:
    alg, digest = "HS384", hashlib.sha384
else:
    alg, digest = "HS256", hashlib.sha256

def enc(value):
    raw = json.dumps(value, separators=(",", ":")).encode()
    return base64.urlsafe_b64encode(raw).rstrip(b"=").decode()

now = int(time.time())
payload = {
    "iss": "TerrakubeInternal",
    "sub": "TerrakubeInternal (EXECUTOR)",
    "aud": "TerrakubeInternal",
    "organizationId": os.environ["ORG"],
    "workspaceId": os.environ["WORKSPACE"],
    "jobId": os.environ["JOB"],
    "stepId": os.environ["STEP"],
    "iat": now,
    "exp": now + 600,
}
unsigned = enc({"alg": alg, "typ": "JWT"}) + "." + enc(payload)
signature = hmac.new(key, unsigned.encode(), digest).digest()
print(unsigned + "." + base64.urlsafe_b64encode(signature).rstrip(b"=").decode())
PY
)"

base="$api_url/tfstate/v1/organization/$org/workspace/$workspace"
state='{"version":4,"serial":1,"lineage":"state-smoke","outputs":{}}'

expect() {
  local want="$1" got="$2" label="$3"
  if [[ "$got" != "$want" ]]; then
    echo "FAIL: $label returned $got, want $want" >&2
    exit 1
  fi
  echo "ok: $label ($got)"
}

code="$(curl -sS -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $token" "$base/backend/terraform.tfstate")"
if [[ "$code" != "200" && "$code" != "204" ]]; then
  echo "FAIL: initial state read returned $code, want 200 or 204" >&2
  exit 1
fi
echo "ok: initial state read ($code)"

code="$(curl -sS -o /dev/null -w '%{http_code}' -X PUT -H "Authorization: Bearer $token" -H 'Content-Type: application/json' --data "$state" "$base/backend/terraform.tfstate")"
expect 200 "$code" "state upload"

body="$(curl -sS -H "Authorization: Bearer $token" "$base/backend/terraform.tfstate")"
if [[ "$body" != "$state" ]]; then
  echo "FAIL: state round-trip mismatch: $body" >&2
  exit 1
fi
echo "ok: state round-trip"

code="$(curl -sS -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $token" "$base/state/terraform.tfstate")"
expect 200 "$code" "current state route"

code="$(curl -sS -o /dev/null -w '%{http_code}' -X DELETE -H "Authorization: Bearer $token" "$base/backend/terraform.tfstate")"
expect 200 "$code" "state delete"

code="$(curl -sS -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $token" "$base/backend/terraform.tfstate")"
expect 204 "$code" "state read after delete"

code="$(curl -sS -o /dev/null -w '%{http_code}' -H "Authorization: Bearer invalid" "$base/backend/terraform.tfstate")"
expect 401 "$code" "invalid token rejection"

echo "state backend smoke test passed"
