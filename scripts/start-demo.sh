#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

if [[ "${APP_ENV:-development}" != development ]]; then
  echo "The demo runner requires APP_ENV=development." >&2
  exit 1
fi
docker compose up -d --build
if ! curl -fsS --retry 30 --retry-delay 2 --retry-all-errors --max-time 3 http://127.0.0.1:2358/languages >/dev/null; then
  echo "Judge0 did not become ready; inspect docker compose logs judge0 judge0-worker." >&2
  exit 1
fi
echo "Demo stack started. Run a candidate code submission to verify execution."
