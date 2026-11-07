# Amara

Amara's backend is a Go modular monolith using Gin. Business modules live under
`internal/`; shared HTTP, configuration, and logging code lives under
`internal/platform/`. See [backend module boundaries](internal/README.md).

## Run the API

Install Go 1.26.10 or newer, then run from the repository root:

```sh
go mod download
make run
```

The API defaults to `:8080`. Check liveness with:

```sh
curl http://localhost:8080/healthz
# {"status":"ok"}
```

The API connects to the Amara Supabase PostgreSQL instance at startup and
fails fast when the database is unreachable. `/health` reports process
liveness only. See [the Supabase ADR](docs/sdd/adr/0001-supabase-postgres.md)
for the connection and migration approach.

## Docker Compose

Compose builds the backend image (the API and worker are two execution roles
of one codebase) and the web client image. No database or queue containers:
persistent data stays in the managed Supabase PostgreSQL instance.

```sh
docker compose up --build
```

Host ports default to `API_PORT=8080` and `WEB_PORT=3000`; set them in the
shell environment when the defaults collide.

- API: `http://localhost:8080/health`
- Web: `http://localhost:3000` (Next.js server with standalone runtime, proxying `/v1` to the
  API, so browser traffic is same-origin)
- Worker: same image as the API

Secrets follow the single Infisical path used by the rest of the project: run
Compose under the Infisical CLI with the backend folder — database credentials
are required and Compose refuses to start without them:

```sh
infisical run --env=dev --path=/backend -- docker compose up
```

`docker-compose.yml` interpolates consuming variables from the injected
environment (for example `DATABASE_URL: ${DATABASE_URL}`); a local `.env` file
stays non-secret. Future provider integrations add their variables to the
same environment block.

## Database migrations

Schema lives in `supabase/migrations/` and runs against the linked Supabase
project (`amara-dev`) with the Supabase CLI. Epic schemas are `accounts`,
`cand`, `job`, `invw`, and `ana` — the `auth` schema name is reserved by
Supabase. See [the ADR](docs/sdd/adr/0001-supabase-postgres.md).

```sh
supabase login                      # once per machine
make db-push                        # apply supabase/migrations to amara-dev
supabase migration new <name>       # create the next migration file
```

`supabase link` was already run for this checkout; new checkouts need
`supabase link --project-ref=<ref>` plus the database password. Do not apply
schema changes through the dashboard SQL editor; they must exist as migration
files in this repository.

## Infisical

Infisical is the secrets manager. The CLI injects configuration and secrets as
environment variables before starting the API; Go reads them once at startup.
The API does not load `.env` files or require Infisical credentials in its code.
See the official [CLI setup](https://infisical.com/docs/cli/usage) and
[runtime injection example](https://infisical.com/blog/secure-secrets-management-for-cursor-cloud-agents).

After installing the CLI, log in (select the US Region) and run the API.
The committed [.infisical.yaml](.infisical.yaml) links this checkout to the
Amara project, so `infisical init` is only needed when creating a new project:

```sh
infisical login
make run-secrets
```

All API configuration and secrets live in the `backend` folder of the
project's `dev` environment. The Makefile targets that folder by default;
select another environment with:

```sh
make run-secrets INFISICAL_ENV=staging
```

For CI and deployment, authenticate a machine identity outside the application,
provide its access token as `INFISICAL_TOKEN`, and pass the project explicitly:

```sh
make build
infisical run --projectId=YOUR_PROJECT_ID --env=prod --path=/backend -- ./bin/api
```

Use your actual project ID, environment slug, and folder. Set `APP_ENV=production`
in the production environment. No project ID or credentials are preconfigured.
Restart the API after changing injected values so it reads the new configuration.

## Configuration

[.env.example](.env.example) lists non-secret settings and development values.
Defaults apply when a variable is absent; invalid values fail startup.
An explicitly empty origin or proxy list disables that feature.

| Variable | Default | Purpose |
| --- | --- | --- |
| `APP_ENV` | `development` | `development`, `test`, `staging`, or `production` |
| `LOG_LEVEL` | `info` | Minimum slog level, such as `debug`, `info`, `warn`, or `error` |
| `HTTP_ADDR` | `:8080` | Listen host and port |
| `HTTP_READ_HEADER_TIMEOUT` | `5s` | Time allowed to read request headers |
| `HTTP_READ_TIMEOUT` | `15s` | Time allowed to read the whole request |
| `HTTP_WRITE_TIMEOUT` | `30s` | Response write deadline |
| `HTTP_IDLE_TIMEOUT` | `60s` | Keep-alive idle timeout |
| `HTTP_SHUTDOWN_TIMEOUT` | `10s` | Time allowed to drain requests after SIGINT or SIGTERM |
| `CORS_ALLOWED_ORIGINS` | Local Vite origins in development | Comma-separated exact HTTP(S) origins; empty by default in other environments |
| `HTTP_TRUSTED_PROXIES` | Empty | Comma-separated trusted proxy IP addresses or CIDR ranges |
| `DATABASE_URL` | Required | PostgreSQL connection string; secrets arrive through Infisical. Use the Supavisor session pooler URL for local runs |
| `JUDGE0_URL` | Optional | Judge0 CE base URL. Set on the worker to enable automatic code grading |
| `JUDGE0_AUTH_TOKEN` | Optional | Native Judge0 `X-Auth-Token` credential. Inject through Infisical if the service requires authentication |
| `JUDGE0_LANGUAGE_IDS` | Optional | JSON map from saved language names to Judge0 language IDs. Defaults to JavaScript 63, TypeScript 74, Python 71, Go 60, and C++ (`cpp`) 54 |
| `JUDGE0_ALLOW_LOCAL_HTTP` | `false` | Set to `true` in development to allow unauthenticated HTTP to a loopback IP. Credentialed connections always require HTTPS |
| `AUTH_MFA_ENCRYPTION_KEY` | Required | Exactly 32-byte secret used to encrypt MFA secrets at rest |
| `DB_POOL_MAX_CONNS` | `4` | Connection pool ceiling per process; shared Supabase plans are small |

Apply the database migrations before enabling Judge0. Use a trusted Judge0 CE
instance with sandbox networking disabled and check the language IDs against
that instance's `GET /languages`. The worker grades submitted attempts from
their saved responses and test-case snapshots, including deadline submissions.
It sends candidate source and test inputs to Judge0 and compares stdout with
expected outputs in Amara. Comparison normalizes CRLF and ignores trailing
newlines. Blank answers earn zero passed cases. Provider errors leave the
attempt pending for retry, which reruns the attempt's cases.
For mixed attempts that require human scores, the company submits evidence for
all non-code questions to the assessment evaluation endpoint. The worker then
combines that evidence with Judge0 case counts and finalizes the attempt.

Development allows `http://localhost:5173` and `http://127.0.0.1:5173`.
Cookie credentials are disabled; cross-origin clients can send an `Authorization`
header. Configure trusted proxies explicitly when deploying behind a load balancer.

Logs are JSON by default; set `LOG_FORMAT=text` for a human-readable format on
container-console runs. Each request receives a generated `X-Request-ID`, also available
through `http.RequestID(ctx)` and in JSON error responses. Request logs include
route templates, status, and duration; they omit bodies, headers, query strings,
and path parameter values. Gin prints its route diagnostics in development mode.

## Development commands

```sh
make build          # Build bin/api
make test           # Run Go tests
make vet            # Run go vet
make gosec          # Run gosec with the reviewed false-positive exceptions
make govulncheck    # Check Go code and both executables' dependency modules
make fmt            # Format Go packages
go test -race ./... # Check for data races
```

Business HTTP adapters register routes under `/v1` through the server's
`RouteRegistrar` functions. Construct their dependencies in `cmd/api/main.go`;
keep business logic and authentication policy in the owning module.

`POST /v1/auth/company-registrations` uses the shared database service and
the HIBP password compromise checker. Candidate/interviewer registration at
`POST /v1/auth/registrations` sends verification codes through the configured
SMTP server. Set `SMTP_HOST`, `SMTP_PORT`, and `SMTP_FROM`; authenticated
servers also require `SMTP_USERNAME` and `SMTP_PASSWORD`. Both endpoints use
the shared request ID, logging, recovery, and CORS middleware.

The `gosec` executable must be installed before running `make gosec`.
[.gosec.json](.gosec.json) excludes only G101 (hardcoded credentials) in the
generated `auth_registration.sql.go`, `login.sql.go`, and `mfa.sql.go` files:
the four flagged constants are parameterized SQL statements containing column
names such as `secret` and `password_hash`, with credential values supplied at
runtime. Other rules still scan these files, and G101 remains enabled elsewhere.
The configuration survives `sqlc generate` without editing generated code.
The SHA-1 exceptions are annotated at their exact import/call sites: HIBP uses
SHA-1 hash prefixes for breach lookups, and TOTP uses HMAC-SHA-1 for existing
authenticator enrollments. Password storage continues to use Argon2id.

Install `govulncheck` before running `make govulncheck`. It scans all Go packages,
then checks the API and worker dependency modules, rejecting advisories even
when the executables do not call the affected functions. CI uses the Go version
in `go.mod`, and Docker builds use Go 1.26.10.
The pinned `golang.org/x/crypto` revision removes the unmaintained OpenPGP
package, addressing [GO-2026-5932](https://pkg.go.dev/vuln/GO-2026-5932) in
addition to the security fixes included in its preceding release.

CI runs Go vet, Staticcheck, gosec, govulncheck, and Semgrep (`p/default`) on
pull requests. To run the individual checks locally, install the required tools:

```sh
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
uv tool install semgrep
make vet
staticcheck ./...
make gosec
make govulncheck
semgrep scan --config p/default --error --strict --metrics=off .
```
