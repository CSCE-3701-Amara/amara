# Amara frontend

Next.js App Router frontend. Run commands from `web/`:

```sh
npm ci
API_URL=http://127.0.0.1:8080 npm run dev
npm run lint
npm run build
npm run test:e2e
```

`API_URL` is a **server-only origin**, configured when Next.js starts in development or builds for production. The development default is `http://127.0.0.1:8080`; the production default is Docker's `http://api:8080`. Set the origin before building if your production API runs elsewhere. Do not include `/v1` in the origin. No browser-exposed API URL or token is needed.

The browser calls same-origin `/v1/...`; Next.js rewrites these requests to the Go API, forwarding cookies and `Set-Cookie`. `/health` and `/health/*` are proxied too. Preserve the backend's Secure, HttpOnly session cookie; use HTTPS outside localhost. Do not run a development server on port 3000 while Docker publishes that port.

Browser tests run a temporary development server on **3100**, configurable with `PLAYWRIGHT_PORT`. Install Chromium once using `npx playwright install chromium`. Tests mock authenticated API responses; read-only backend probes are described in [INTEGRATION.md](./INTEGRATION.md).

## Connected screens

- `/login`, `/register`: sessions, MFA, registration, and backend error states.
- `/candidate/dashboard`: filtered, cursor-paginated candidate feed and application submission.
- `/candidate/applications`: offset-paginated history and application details.
- `/employer/jobs`, `/employer/jobs/new`, `/employer/jobs/[id]`: posting reads, creation, replacement, publication, archiving, deletion, and per-posting pipeline drafts.
- `/employer/pipelines`: assessment drafts, editing, publication, and copying from the published library.
- Workspace headers: logout through the backend's current-session endpoint.

Current backend blockers and contract discrepancies are recorded in [INTEGRATION.md](./INTEGRATION.md). A connected request can still be rejected by the existing backend; the UI displays the error and supports retry. It never fakes registration, submission, or persistence.

## Figma coverage and production browser checks

The current full-file inventory and per-state implementation/API status are in
[design/COVERAGE.md](design/COVERAGE.md). Regenerate the Markdown table after
editing the audited inventory with `python3 scripts/figma-coverage.py`.

To exercise the built frontend with Playwright owning the server process:

```sh
npm run build
PLAYWRIGHT_SERVER=production PLAYWRIGHT_PORT=3202 npm run test:e2e
```

The default test command still starts the development server on port 3100.
Product-state smoke checks cover 1728px, 1024px, and 375px, with separate DAG
interaction checks at all three widths. Browser API successes use mocks;
[INTEGRATION.md](INTEGRATION.md) records why live API verification is blocked.
