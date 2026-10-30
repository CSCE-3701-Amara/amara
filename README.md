# Amara

Amara's backend is a Go modular monolith using Gin. Business modules live under
`internal/`; shared HTTP, configuration, and logging code lives under
`internal/platform/`. See [backend module boundaries](internal/README.md).

## Run the API

Install Go 1.26 or newer, then run from the repository root:

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
| `DB_POOL_MAX_CONNS` | `4` | Connection pool ceiling per process; shared Supabase plans are small |

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
make fmt            # Format Go packages
go test -race ./... # Check for data races
```

Business HTTP adapters register routes under `/v1` through the server's
`RouteRegistrar` functions. Construct their dependencies in `cmd/api/main.go`;
keep business logic and authentication policy in the owning module.
