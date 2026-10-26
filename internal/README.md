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