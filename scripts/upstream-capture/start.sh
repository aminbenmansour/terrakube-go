#!/usr/bin/env bash
#
# A disposable upstream Terrakube for capturing the real JSON:API contract.
#
# It runs upstream's released Java image so responses can be recorded instead of
# guessed. Everything binds on the host network on non-default ports so it cannot
# collide with the terrakube-go dev stack. Needs podman; 2.34.3 seeds the `demo`
# profile, which is where the recorded ids and names come from.
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
tag="${TK_UPSTREAM_VERSION:-2.34.3}"

# No secrets in the repository: a shared PAT and internal secret are generated
# once per checkout and kept in a git-ignored file so pat.py can sign with them.
secrets="$here/.secrets.env"
if [ ! -f "$secrets" ]; then
  {
    # base64url: upstream decodes PatSecret with a URL decoder and exits on a +.
    echo "PAT_SECRET=$(head -c 32 /dev/urandom | base64 | tr '+/' '-_' | tr -d '=')"
    echo "INTERNAL_SECRET=$(head -c 32 /dev/urandom | base64 | tr '+/' '-_' | tr -d '=')"
  } > "$secrets"
  chmod 600 "$secrets"
fi
# shellcheck disable=SC1090
. "$secrets"
pat="$PAT_SECRET"
internal="$INTERNAL_SECRET"

# Container credentials are generated too: the upstream demo values are public,
# but they read as leaked passwords to a scanner, and this stack binds the host
# network.
pg_pass="$(head -c 24 /dev/urandom | base64 | tr '+/' '-_' | tr -d '=')"
redis_pass="$(head -c 24 /dev/urandom | base64 | tr '+/' '-_' | tr -d '=')"

podman rm -f tk-upstream-postgres tk-upstream-redis tk-upstream-dex tk-upstream-api >/dev/null 2>&1 || true

podman run -d --name tk-upstream-postgres --network host \
  -e POSTGRES_USER=terrakube -e POSTGRES_PASSWORD="$pg_pass" \
  -e POSTGRES_DB=terrakubedb \
  docker.io/library/postgres:16 -p 5433 >/dev/null

podman run -d --name tk-upstream-redis --network host \
  docker.io/library/redis:7 redis-server --port 6380 --requirepass "$redis_pass" >/dev/null

podman run -d --name tk-upstream-dex --network host \
  -v "$here/dex-config.yaml:/etc/dex/config.yaml:ro,z" \
  ghcr.io/dexidp/dex:v2.44.0 dex serve /etc/dex/config.yaml >/dev/null

# Wait for Postgres before the API: Liquibase would exit rather than retry.
for _ in $(seq 1 60); do
  if podman exec tk-upstream-postgres pg_isready -p 5433 -U terrakube -d terrakubedb >/dev/null 2>&1; then break; fi
  sleep 1
done

podman run -d --name tk-upstream-api --network host \
  -e SERVER_PORT=18080 \
  -e ApiDataSourceType=POSTGRESQL \
  -e DatasourceHostname=localhost -e DatasourcePort=5433 \
  -e DatasourceDatabase=terrakubedb -e DatasourceUser=terrakube \
  -e DatasourcePassword="$pg_pass" \
  -e GroupValidationType=DEX -e UserValidationType=DEX -e AuthenticationValidationType=DEX \
  -e DexIssuerUri=http://localhost:5556/dex -e DexClientId=terrakube-app \
  -e TerrakubeHostname=localhost:18080 \
  -e TerrakubeUiURL=http://localhost:3000 \
  -e PatSecret="$pat" -e InternalSecret="$internal" \
  -e StorageType=LOCAL \
  -e AzBuilderExecutorUrl=http://localhost:19090/api/v1/terraform-rs \
  -e TerrakubeToolsRepository=https://github.com/terrakube-io/terrakube-extensions.git \
  -e TerrakubeToolsBranch=main \
  -e CustomTerraformReleasesUrl=https://releases.hashicorp.com/terraform/index.json \
  -e TerrakubeRedisHostname=localhost -e TerrakubeRedisPort=6380 \
  -e TerrakubeRedisPassword="$redis_pass" \
  -e spring_profiles_active=demo \
  "docker.io/azbuilder/api-server:${tag}" >/dev/null

echo "postgres 5433, redis 6380, dex 5556, api 18080 (upstream ${tag})"
echo "watch: podman logs -f tk-upstream-api"
