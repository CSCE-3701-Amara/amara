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

The scaffold runs without external services. Supabase and business endpoints
will be added separately. `/healthz` reports process liveness only.

## Infisical

Infisical is the secrets manager. The CLI injects configuration and secrets as
environment variables before starting the API; Go reads them once at startup.
The API does not load `.env` files or require Infisical credentials in its code.
See the official [CLI setup](https://infisical.com/docs/cli/usage) and
[runtime injection example](https://infisical.com/blog/secure-secrets-management-for-cursor-cloud-agents).

After installing the CLI, log in and link this checkout to your project, Make sure to select US Region: 

```sh
infisical login
infisical init
make run-secrets
```

Set values in the project's `dev` environment. The Makefile uses the root secret
folder by default; select another environment or folder with:

```sh
make run-secrets INFISICAL_ENV=staging INFISICAL_PATH=/backend
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
| `LOG_FORMAT` | `json` | `json` for aggregation, or `text` for readable runs such as `docker compose logs` |
| `HTTP_ADDR` | `:8080` | Listen host and port |
| `HTTP_READ_HEADER_TIMEOUT` | `5s` | Time allowed to read request headers |
| `HTTP_READ_TIMEOUT` | `15s` | Time allowed to read the whole request |
| `HTTP_WRITE_TIMEOUT` | `30s` | Response write deadline |
| `HTTP_IDLE_TIMEOUT` | `60s` | Keep-alive idle timeout |
| `HTTP_SHUTDOWN_TIMEOUT` | `10s` | Time allowed to drain requests after SIGINT or SIGTERM |
| `CORS_ALLOWED_ORIGINS` | Local Vite origins in development | Comma-separated exact HTTP(S) origins; empty by default in other environments |
| `HTTP_TRUSTED_PROXIES` | Empty | Comma-separated trusted proxy IP addresses or CIDR ranges |

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
