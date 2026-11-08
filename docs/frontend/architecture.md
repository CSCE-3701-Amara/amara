# Frontend architecture

The frontend is a Next.js App Router application under `web/`. It uses
feature-oriented code organization, server-rendered route entry points, and
client components only where interaction or browser APIs require them.

## Layers

```text
web/src/
├── app/                 # App Router pages, route groups, and layouts
├── features/
│   ├── auth/            # Sign-in, MFA, and registration
│   ├── cand/            # Candidate experience
│   ├── job/             # Job postings and applications
│   ├── invw/            # Interview workflow
│   ├── ana/             # Analytics boundary
│   └── management/      # Shared employer/admin workflow screens
├── shared/
│   ├── lib/             # API client, contracts, and cross-feature hooks
│   ├── ui/              # Workspace chrome and reusable product components
│   └── ...
└── generated/           # Generated frontend artifacts
```

Feature code owns API helpers, feature types, pages/components, and feature
documentation. Reusable controls belong in `shared/ui`; cross-feature request
and error behavior belongs in `shared/lib`.

## App Router groups

Route groups keep URL paths clean while applying workspace-specific layouts:

- `(public)` renders `/login` and `/register` without workspace navigation.
- `(candidate)` wraps `/candidate/*` in the candidate `Workspace`.
- `(employer)` wraps `/employer/*` in the company `Workspace`.
- `(interviewer)` wraps `/interviewer/*` in the interviewer `Workspace`.
- `(admin)` wraps `/admin/*` in the administration workspace.

The root layout loads local Funnel Sans and Funnel Display fonts, global
tokens, metadata, and the document language. A route page should compose the
feature component rather than duplicate workspace chrome.

## Data flow

Client components use `shared/lib/api-client.ts` for backend requests. The
client:

1. accepts only relative `/v1/...` paths;
2. sends `credentials: "include"` so the HttpOnly session cookie is used;
3. disables browser caching for private data;
4. maps problem responses into `ApiError`;
5. exposes retry and sign-in recovery through shared feedback components.

`useApiResource` cancels stale reads and supports explicit reloads. Mutations
should guard against duplicate submissions, preserve failed form input, and
show the server's error rather than a success-shaped fallback.

## Integration boundary

Next.js rewrites `/v1/*` and `/health/*` to the origin in `API_URL` (development
defaults to `http://127.0.0.1:8080`; standalone production defaults to
`http://api:8080`). The browser therefore talks to one origin and never needs
an exposed API token or backend URL.

Keep backend ownership and authorization in the Go modules. The frontend may
validate form shape and render state, but it must not treat a client-supplied
company, candidate, or account ID as an authority claim.

## Testing

Run from `web/`:

```sh
npm run lint
npm run build
npm run test:e2e
```

Playwright tests mock API responses for UI behavior and run on port 3100 by
default. Backend route and authorization behavior remains covered by Go tests.