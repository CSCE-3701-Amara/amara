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
