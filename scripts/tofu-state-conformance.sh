#!/usr/bin/env bash
# Run a real OpenTofu HTTP backend conformance test against the local API.
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
#   TOFU_BIN                         default tofu
#
# Use a dedicated workspace: the test writes and destroys state.
set -euo pipefail

api_url="${TERRAKUBE_API_URL:-http://localhost:8080}"
secret="${TERRAKUBE_INTERNAL_TOKEN_SECRET:?set TERRAKUBE_INTERNAL_TOKEN_SECRET}"
org="${TERRAKUBE_STATE_ORG_ID:?set TERRAKUBE_STATE_ORG_ID}"
workspace="${TERRAKUBE_STATE_WORKSPACE_ID:?set TERRAKUBE_STATE_WORKSPACE_ID}"
job="${TERRAKUBE_STATE_JOB_ID:?set TERRAKUBE_STATE_JOB_ID}"
step="${TERRAKUBE_STATE_STEP_ID:?set TERRAKUBE_STATE_STEP_ID}"
tofu="${TOFU_BIN:-tofu}"

command -v "$tofu" >/dev/null 2>&1 || {
  echo "FAIL: $tofu not found; install OpenTofu or set TOFU_BIN" >&2
  exit 1
}

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
    "exp": now + 900,
}
unsigned = enc({"alg": alg, "typ": "JWT"}) + "." + enc(payload)
signature = hmac.new(key, unsigned.encode(), digest).digest()
print(unsigned + "." + base64.urlsafe_b64encode(signature).rstrip(b"=").decode())
PY
)"

workdir="$(mktemp -d)"
cleanup() {
  rm -rf "$workdir"
}
trap cleanup EXIT

cat > "$workdir/main.tf" <<HCL
terraform {
  backend "http" {
    address = "$api_url/tfstate/v1/organization/$org/workspace/$workspace/backend/terraform.tfstate"
  }
}

resource "terraform_data" "conformance" {
  input = "state backend conformance"
}
HCL

export TF_HTTP_USERNAME=terraform
export TF_HTTP_PASSWORD="$token"
export TF_IN_AUTOMATION=1
export TF_LOG=DEBUG
export TF_LOG_PATH="$workdir/tofu.log"

cd "$workdir"
"$tofu" init -input=false -no-color
"$tofu" apply -auto-approve -input=false -no-color
state="$("$tofu" state pull)"
if [[ "$state" != *"state backend conformance"* ]]; then
  echo "FAIL: state pull did not contain the test value" >&2
  exit 1
fi
echo "ok: init/apply/state pull"

"$tofu" destroy -auto-approve -input=false -no-color
echo "ok: destroy"

if grep -qiE '\b(LOCK|UNLOCK)\b' "$workdir/tofu.log"; then
  echo "FAIL: OpenTofu log contains LOCK/UNLOCK; upstream HTTP backend should not lock" >&2
  grep -inE '\b(LOCK|UNLOCK)\b' "$workdir/tofu.log" | head -20 >&2
  exit 1
fi
echo "ok: no LOCK/UNLOCK requests"

echo "OpenTofu state backend conformance passed"
