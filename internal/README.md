# Backend Modules

The Go backend is a modular monolith. Each directory is a business module with
its own domain rules, application use cases, ports, adapters, and transport
code.

- `auth`: accounts and access control
- `cand`: candidate-facing experiences
- `job`: postings and application admission
- `invw`: interview workflow, assessments, and offers
- `ana`: analytics and evaluation

`platform` contains shared technical infrastructure. `integration` contains
external provider adapters. Neither directory owns business state.

## Package layout

Each epic follows the same dependency direction:

```text
internal/<epic>/
├── domain/             # Business entities and state-transition rules
├── application/        # Use cases and orchestration
├── ports/              # Interfaces owned by the epic
└── adapters/
    ├── http/           # HTTP handlers and transport mapping
    └── database/       # Persistence implementations for the epic
```

The layers are added consistently now so new work has an obvious home. A
package should contain real implementation code before adding more subpackages
or files inside it.

Shared technical concerns belong under `platform`:

```text
internal/platform/
├── config/
├── database/           # Connection and transaction infrastructure only
├── http/               # Server composition and cross-cutting middleware
├── logging/
└── doc.go
```

Provider-specific clients belong under `integration` and implement ports
defined by the consuming epic or by a shared platform contract:

```text
internal/integration/
├── email/
├── llm/
├── meetings/
├── storage/
└── doc.go
```

### Ownership rules

- Epic code stays inside its epic directory.
- An epic owns its domain state and database write boundary.
- `platform` may provide technical infrastructure but must not own business
  state.
- `integration` adapts external providers and must not contain epic business
  rules.
- Cross-epic writes use the owning epic's application contract or an event;
  another epic must not write its tables directly.
- HTTP handlers validate transport concerns and call application operations;
  domain rules do not belong in handlers.

## API composition

`cmd/api/main.go` and `cmd/worker/main.go` load validated configuration,
create the logger, initialize the database `Service` (fail fast on unreachability),
and run until a signal arrives. The API then starts its HTTP server; the
worker currently has no queue consumption and exits after shutdown.

- `platform/config/config.go` reads process environment variables. Infisical
  injects secrets before startup; the API does not load `.env` files.
- `platform/config/database.go` defines database pool limits, lifetimes, and
  health check timeouts.
- `platform/database/database.go` provides `database.Service` for lifecycle
  management, ping connectivity, and pool stats. Repository adapters obtain
  the underlying `*pgxpool.Pool` via `dbSvc.Pool()`.
- `platform/logging/logging.go` creates the JSON slog logger.
- `platform/http/server.go` configures HTTP timeouts and drains active requests
  on shutdown.
- `platform/http/routes.go` registers `/health/live` (liveness), `/health/ready`
  (readiness checking the database), `/health` (legacy liveness alias), the
  `/v1` module group, and JSON errors for unknown routes and unsupported methods.
- `platform/http/middleware.go` supplies request IDs, request logging, panic
  recovery, and CORS. Forwarded client addresses are trusted only for configured
  proxy addresses.

To expose a new module, implement a route registration method in its HTTP
adapter with the signature `func(*gin.RouterGroup)`. Construct the adapter and
its application dependencies in `cmd/api`, then pass the registration method
to `httpserver.New(cfg.HTTP, logger, dbSvc, handler.RegisterRoutes)`. Registrars receive
the `/v1` group, so registering `/jobs` exposes `/v1/jobs`.

The adapter must apply AUTH's authentication and authorization middleware to
protected routes. The shared database connection pool is managed by
`internal/platform/database.Service` and handed to module repository adapters
(`adapters/database`) via `dbSvc.Pool()`. Business services depend strictly on
epic-owned repository ports, never directly on `pgxpool` or database infrastructure.
