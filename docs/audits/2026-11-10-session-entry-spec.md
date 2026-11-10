# Public entry and persistent sessions

Goal: restore the previous welcome/login visuals while retaining form behavior,
and return authenticated visitors to their existing workspace.

Scope: welcome and login presentation, workspace logo destination, a read-only
current-session endpoint, public entry handling, and regression tests. Registration
and workspace visuals, credential validation, session lifetime, and MFA remain intact.

Contract: GET /v1/auth/sessions/current validates the existing HttpOnly cookie and
returns {user: {id, role}, session: {expiresAt}} with Cache-Control: no-store.
No token is returned. Missing/expired/revoked sessions return 401; storage failures
return 500. Existing POST and DELETE session contracts remain unchanged.

| Input | Outcome |
| --- | --- |
| No cookie | Welcome/login remain available |
| Valid candidate cookie | /candidate/dashboard |
| Valid company cookie issued after MFA | /employer/jobs |
| Valid interviewer cookie | /interviewer/interviews |
| Workspace logo | Current role's home, including /admin/dashboard |
| Expired/revoked/invalid cookie | Login/welcome available; no authenticated redirect |
| Malformed role or expiry | Session check error with retry, no guessed destination |
| Session service unavailable | Retry message; do not interpret as signed out |
| New company login without MFA | Existing MFA challenge; no session issued |
| Navigation races | Abort stale lookup; use replace to avoid redirect loops |

Implementation: reuse SessionValidator and its account-status checks. Public routes
read cookie presence inside Suspense; only session-bearing entries need a client
lookup. This is navigation UX, not authorization: protected APIs keep enforcing
sessions independently. Use existing role destinations, no browser-stored tokens
or passwords, and no persistent client-side identity cache.

Test scenarios: endpoint identity/privacy and cache headers; missing/invalid cookies;
service failure; candidate/interviewer/company resume across browser contexts;
expired session; role-aware logo; logout; MFA before manager session creation;
login visual variants and unchanged registration, responsive layouts and focus.

Rollout: backend endpoint and frontend ship together. No database migration,
new dependency, cookie-format change, or session extension.

## Implementation and verification

- Restored the welcome page's previous composition and shader; retained whitespace,
  responsive, focus, and navigation fixes.
- Added a scoped original login theme. Registration and workspaces keep the refined
  theme; shared form state, validation, password visibility, and MFA are preserved.
- Added and documented GET /v1/auth/sessions/current, reusing SessionValidator and
  the existing account-status/expiry/revocation checks. The route is wired in cmd/api.
- Public entry with a cookie validates before displaying welcome/login. All redirect
  destinations are fixed local routes. Storage failures offer retry; 401 shows guest UI.
- Workspace logo destinations are role-specific, including administration.
- Confirmed the initial logo regression (all three roles linked to `/`); the new tests
  now pass. Backend current-session tests cover all three supported roles, token
  privacy, no cookie renewal, cache headers, missing/revoked cookies, and service failure.
- Browser regressions cover restored login versus refined registration, persisted
  cookies in a new browser context, public entry for each role, browser Back after
  sign-in, expired/revoked responses, malformed roles, retry, and manager MFA issuance.
- `go test ./internal/auth/... ./cmd/api`: passed.
- `go vet ./internal/auth/... ./cmd/api`: passed.
- `npm run lint`, `npx tsc --noEmit`, and `npm run build`: passed.
- Full production browser suite: 418 passed, 0 failed, 0 skipped, 0 flaky.
- Desktop/mobile screenshots reviewed for the restored login and welcome views.

The browser suite uses API fixtures; Go route tests execute the real Gin middleware
with a controlled session validator. Live database/email/MFA delivery was not exercised.
The API and frontend must run the updated versions together to expose the new endpoint.
