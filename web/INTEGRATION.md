# Frontend / backend integration

Checked against the registrars in `cmd/api/main.go`, handler adapters under `internal/`, and `docs/api/openapi.yaml`. Backend code and the OpenAPI file were left unchanged.

## Transport and schemas

`src/shared/lib/api-client.ts` makes same-origin, uncached requests with `credentials: include`. It understands both RFC 9457 Problem Details and the legacy `{error: {code, message}}` format, preserves Retry-After and request IDs, accepts empty 204 responses, and hides non-JSON proxy error bodies. Reads abort when their parameters change or a route is hidden/unmounted. Failed writes keep form input for retry.

The frontend uses registered `/v1` paths. The contract's legacy `/api/v1/auth/*` paths and claimed compatibility aliases are not mounted by this checkout's `internal/platform/http/routes.go`. Do not change frontend prefixes to `/api/v1` without reconciling those registrars first.

| Workflow                              | Registered paths                                                                    | Response shape                                                 |
| ------------------------------------- | ----------------------------------------------------------------------------------- | -------------------------------------------------------------- |
| Login / MFA                           | POST `/v1/auth/sessions`                                                            | `user`, `session.expiresAt`, HttpOnly cookie                   |
| Registration                          | POST `/v1/auth/registrations`, `/v1/auth/company-registrations`                     | RegistrationResponse or Problem Details                        |
| Logout                                | DELETE `/v1/auth/sessions/current`                                                  | 204 + clearing cookie                                          |
| Account actions (client helper only)  | POST `/v1/auth/account-actions`                                                     | 204 + clearing cookie                                          |
| Candidate feed                        | GET `/v1/jobs/feed`                                                                 | `items`, optional `next_cursor`; snake_case filters            |
| Application submission                | POST `/v1/jobs/{id}/apply`                                                          | snake_case screening result                                    |
| Application history / details         | GET `/v1/candidates/me/applications`, `/v1/applications/{id}`                       | `data`, `pagination`; camelCase application fields             |
| Posting management                    | GET/POST `/v1/jobs`, GET/PUT/DELETE `/v1/jobs/{id}`, POST publish/republish/archive | `items` or a snake_case posting; deletion 204                  |
| Company listing (client helper only)  | GET `/v1/company/{company_id}/postings`                                             | `items`, optional `next_cursor`                                |
| Assessment modules                    | `/v1/modules`, `/v1/modules/library`, `/v1/modules/{id}`, publish/copy              | Bare summary arrays or full module definition                  |
| Library problems (client helper only) | GET `/v1/modules/library/problems`                                                  | Bare problem array                                             |
| Pipeline draft                        | GET/POST `/v1/jobs/{id}/pipeline`                                                   | `id`, `postingId`, `versionId`, `versionNumber`                |
| Linear stages                         | GET/PUT `/v1/pipeline-versions/{versionId}/stages`                                  | Ordered bare stage array; camelCase fields, RFC3339 timestamps |

The assessment UI creates behavioral modules and edits existing draft metadata without discarding other questions, rubric, library configuration, or Amarat. Published and copied modules remain locked. Pipeline editing adds/reorders assessment stages and preserves existing live-interview configuration; new live-interview authoring is not supplied by this screen. More advanced planned OpenAPI operations, admin review, profiles, grading, and pipeline publication have no connected UI.

## Backend changes needed before all flows can succeed

These are findings and proposed next steps, not changes made by this task:

1. `cmd/api/main.go` passes `nil` to `NewPostingHandler` and `NewApplicationHandler`. Their registrars deliberately return 401 for all posting management and application submission requests, even with valid sessions. Wire the existing session validator to JOB authentication and company scoping before claiming these workflows work end to end. JOB currently treats its authenticated ID as a company ID for management, whereas AUTH distinguishes AccountID and CompanyID; the adapter must preserve that distinction and use validated company scope, never a client-supplied identity. The frontend preserves these errors rather than automatically redirecting every 401 to login.
2. Candidate/interviewer registration uses `RegisterUnavailableRegistrationRoutes`: it always returns 503 `registration_unavailable` until verification delivery is configured. Company registration has a real handler. The UI only enters a pending state after successful registration.
3. Company listing requires a company ID. Login returns an account ID, email, and role, but no company ID; there is no registered session/profile read endpoint to recover it. The frontend does not guess a company ID, store session secrets, or ask users to paste database IDs. Its company-list client helper can be used when an authenticated identity endpoint provides that scope.
4. OpenAPI still marks company registration, logout, and account actions as planned, although their handlers are now registered. It also claims `/api/v1` compatibility aliases that this router does not register. Reconcile the spec with the code before generating a client directly from it.

## Verification

Browser integration tests use response shapes from the existing handlers and contract. They verify session-cookie propagation and clearing, MFA, feed filter/cursor encoding, application history/detail envelopes, failed and successful submission states, posting payloads and deletion, module creation/publication/copying, ordered pipeline-stage writes, retry, and narrow layouts. Mocked success cannot establish that the blocked backend paths succeed.

For a running backend, safe read-only probes are `GET /health/ready`, `GET /v1/jobs/feed` without a cookie (401), and GET-only checks of documented routes. Run the same probes through the frontend origin to verify rewrite behavior. Authenticated end-to-end verification needs valid test accounts and the backend fixes above; this task does not seed accounts or mutate backend data.

During this run, direct read-only probes to `127.0.0.1:8080` failed with connection refused, and the development proxy consequently could not reach the API. The backend was not started or modified. Live authenticated success and production-proxy cookie behavior therefore remain unverified; the successful workflow tests use mocked handler responses.
