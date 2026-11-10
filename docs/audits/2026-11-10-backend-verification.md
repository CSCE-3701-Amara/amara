# Backend verification — 10 November 2026

Verified on `FIXING` after merge `85632c9` and frontend fix `2a74db3`.

## Fixes

- Removed the duplicate `AutoApprove` field and eighth SQL argument in the
  generated company-registration query. The SQL has seven parameters.
- Added a PostgreSQL regression test for manual and automatic approval, checking
  returned status, review state, and creation of the approved company.
- Documented narrow security-linter exceptions for the fixed Docker command
  with allowlisted configuration and server-generated arguments, and SHA-1
  required by the HIBP prefix lookup protocol. Password storage uses Argon2id.

## Completed checks

| Check | Result |
| --- | --- |
| `go build ./...` | Pass |
| `go test -p=2 -race -count=1 -cover -json ./...` | 518 tests passed; no skipped tests or race warnings |
| All deployment migrations on a fresh PostgreSQL 16 database | Pass |
| PostgreSQL integration tests | Pass; `AMARA_TEST_DATABASE_URL` enabled |
| Real local Judge0 integration tests | Pass; `AMARA_SANDBOX_URL` enabled |
| `gofmt -l .` | No files reported |
| `go vet ./...` | Pass |
| Staticcheck 2026.2.2 | Pass |
| `go mod verify` | All modules verified |
| Gosec with `.gosec.json` | No findings or package errors |
| Semgrep `p/default`, backend Go files | No findings or scanner errors |
| Govulncheck, source and API/worker module scans | No vulnerabilities found |

The test run also reports packages with no test files; none of the individual
tests were skipped. PostgreSQL ran in a separate disposable container, with no
changes to the existing application databases.

## Remaining CodeQL alerts

CodeQL 2.27.4 with `codeql/go-queries@1.6.14` completed the Go security-extended
suite. The email-injection rule reports zero findings. Three separate alerts
remain and were not suppressed or resolved by this verification:

- `go/weak-sensitive-data-hashing`: HIBP's required SHA-1 prefix lookup in
  `internal/auth/adapters/password/hibp.go`.
- `go/log-injection`: registration error logging in
  `internal/auth/adapters/http/registration.go`.
- `go/log-injection`: request logging in `internal/platform/http/middleware.go`.

The CI checks pass; the CodeQL report is not free of findings.
