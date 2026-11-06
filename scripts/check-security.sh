#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: scripts/check-security.sh [--release]

Runs Go vet, Staticcheck, gosec, govulncheck (code and API/worker modules),
and Semgrep using p/default. Each check runs even if another fails.
Missing tools and scanner errors count as failures. No tools are installed.

  --release  Also run the active ZAP scan against STAGING_URL and the HTTP
             header check against PUBLIC_ORIGIN. Both URLs are required.
  -h, --help Show this help.

Required tools: go, staticcheck, gosec, govulncheck, semgrep.
Release checks also require Docker, curl, grep, and tr.
Go tools are located through PATH and the configured GOBIN/GOPATH.

SECURITY_REPORT_DIR sets the output directory. By default a temporary
amara-security.* directory holds per-check logs and summary.txt.
Exit status: 0 if all requested checks pass; 1 on a failed/missing check;
2 for invalid arguments or missing release URLs.
EOF
}

release=false
for arg in "$@"; do
  case "$arg" in
    --release) release=true ;;
    -h|--help) usage; exit 0 ;;
    *) printf 'Unknown argument: %s\n' "$arg" >&2; usage >&2; exit 2 ;;
  esac
done

if "$release"; then
  for name in STAGING_URL PUBLIC_ORIGIN; do
    url="${!name:-}"
    case "$url" in
      http://?*|https://?*) ;;
      *) printf '%s must be an HTTP(S) URL when using --release.\n' "$name" >&2; exit 2 ;;
    esac
  done
fi

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
cd -- "$repo_root"

# Go installs command-line tools outside PATH on some developer machines.
if command -v go >/dev/null 2>&1; then
  go_bin_dir="$(go env GOBIN 2>/dev/null || true)"
  if [[ -z "$go_bin_dir" ]]; then
    go_path="$(go env GOPATH 2>/dev/null || true)"
    if [[ -n "$go_path" ]]; then
      go_bin_dir="${go_path%%:*}/bin"
    fi
  fi
  if [[ -n "$go_bin_dir" ]]; then
    export PATH="$PATH:$go_bin_dir"
  fi
fi

if [[ -n "${SECURITY_REPORT_DIR:-}" ]]; then
  mkdir -p -- "$SECURITY_REPORT_DIR"
  report_dir="$(cd -- "$SECURITY_REPORT_DIR" && pwd)"
else
  report_dir="$(mktemp -d "${TMPDIR:-/tmp}/amara-security.XXXXXX")"
fi

failures=0
results=()

run_check() {
  local name="$1" log_name="$2" status
  shift 2
  printf '\n=== %s ===\n' "$name"
  if ! command -v "$1" >/dev/null 2>&1; then
    printf 'Missing required tool: %s\n' "$1" | tee "$report_dir/$log_name.log"
    results+=("MISSING: $name ($1)")
    failures=$((failures + 1))
  elif "$@" 2>&1 | tee "$report_dir/$log_name.log"; then
    results+=("PASS: $name")
  else
    status=$?
    results+=("FAIL: $name (exit $status)")
    failures=$((failures + 1))
  fi
}

printf 'Repository: %s\nReports: %s\n' "$repo_root" "$report_dir"
run_check 'Go vet' go-vet go vet ./...
run_check 'Staticcheck' staticcheck staticcheck ./...
run_check 'gosec' gosec gosec -conf .gosec.json ./...
run_check 'govulncheck: code' govulncheck-code govulncheck ./...
run_check 'govulncheck: API modules' govulncheck-api govulncheck -C cmd/api -scan module
run_check 'govulncheck: worker modules' govulncheck-worker govulncheck -C cmd/worker -scan module
run_check 'Semgrep' semgrep semgrep scan --config p/default --error --strict --metrics=off .

if "$release"; then
  run_check 'ZAP active scan' zap docker run --rm ghcr.io/zaproxy/zaproxy:stable \
    zap-full-scan.py -t "$STAGING_URL" -m 5 -T 10
  run_check 'HTTP security headers' security-headers bash scripts/check-security-headers.sh
else
  printf '\nRelease checks were not requested; use --release with STAGING_URL and PUBLIC_ORIGIN.\n'
fi

printf '\nSecurity check summary\n'
printf '%s\n' "${results[@]}" | tee "$report_dir/summary.txt"
printf 'Reports: %s\n' "$report_dir"
if ((failures > 0)); then
  printf '%d check(s) failed or could not run.\n' "$failures" >&2
  exit 1
fi
