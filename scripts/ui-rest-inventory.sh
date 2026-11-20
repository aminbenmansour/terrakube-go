#!/usr/bin/env bash
#
# Lists the REST calls the upstream React UI makes, normalised so they can be
# compared with this API. Needs an upstream checkout, which is not vendored:
#
#   ./scripts/ui-rest-inventory.sh ~/Documents/terrakube/ui
#
# Output is "METHOD /path" with interpolated ids replaced by :id, unique and
# sorted, which is exactly the layout of docs/upstream/ui-rest-operations.txt.
set -euo pipefail

src="${1:-}"
if [[ -z "$src" ]]; then
  echo "usage: $0 <path-to-upstream-ui>" >&2
  exit 2
fi
if [[ ! -d "$src/src" ]]; then
  echo "error: ${src}/src does not exist" >&2
  exit 2
fi

find "$src/src" -name '*.ts' -o -name '*.tsx' | grep -v '__tests__' | sort | xargs \
  perl -0777 -ne 'while (/axios\w*\.(get|post|patch|put|delete)\(\s*[`"'"'"']([^`"'"'"']*)/g) { print uc($1), " ", $2, "\n" }' |
  sed -E 's/\$\{[^}]*\}/:id/g' |
  awk '{ path = $2; sub(/\?.*/, "", path); if (path !~ /^\//) path = "/" path; gsub(/\/+/, "/", path); print $1, path }' |
  sort -u
