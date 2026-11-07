#!/usr/bin/env bash
set -euo pipefail

URL="${PUBLIC_ORIGIN:?Set PUBLIC_ORIGIN to your public URL}"

# Fetch headers from the final response after redirects.
HEADERS=$(curl -fsSL --max-time 20 -o /dev/null \
  -w '%{header_json}' "$URL" | tr '[:upper:]' '[:lower:]')

check_header() {
  local name="$1"
  local expected="${2:-}"

  if ! grep -Eq "\"${name}\":" <<< "$HEADERS"; then
    echo "FAIL: Missing header: $name"
    return 1
  fi

  if [[ -n "$expected" ]] &&
     ! grep -Eq "\"${name}\":[[:space:]]*\[[[:space:]]*\"[^\"]*${expected}" <<< "$HEADERS"; then
    echo "FAIL: $name does not match expected value: $expected"
    return 1
  fi

  echo "PASS: $name"
}

check_header "strict-transport-security" || exit 1
check_header "content-security-policy" || exit 1
check_header "x-content-type-options" "nosniff" || exit 1
check_header "referrer-policy" || exit 1
check_header "x-frame-options" || exit 1

echo "Security header checks passed."
