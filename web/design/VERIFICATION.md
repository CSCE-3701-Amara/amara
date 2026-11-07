# Verification evidence

## Candidate email verification follow-up

`POST /v1/auth/verifications` is wired to the candidate registration pending state. The focused auth, accessibility, and responsive-screen suites passed **38/38** against the updated production frontend, including invalid-code retry, input preservation, successful activation, local code-format validation, and pending-screen layouts at 1728px and 375px. Lint and type-check passed. The production build passed with Next.js webpack; a later Turbopack rebuild hit this sandbox's local-port restriction in its CSS worker. The handler and OpenAPI request/response shape were inspected. Browser API responses are mocked; live activation remains unverified because no development account was used.

## Backend route fix follow-up

The pipeline-progress route now uses `:applicationId`, matching sibling application routes. A regression test covers registration in both orders. `go test ./...` passed. `make run-secrets` connected to the development database and started the API; direct read-only requests returned `200` from `/health/live` and `401 authentication_required` from the pipeline-progress route without a session. The API was stopped cleanly after the probes. Authenticated pipeline progress and other live writes remain unverified.

Final production browser run on 2026-11-07: **380 passed, 0 failed, 0 skipped, 0 flaky** in approximately 99 seconds. Results are recorded in [verification.json](verification.json).

## Final checks

| Check                                                                                               | Result                                                                  |
| --------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| `npm run lint`                                                                                      | Passed                                                                  |
| `tsc --project web/tsconfig.json --noEmit`                                                          | Passed after fresh standard build                                       |
| `npm run build`                                                                                     | Passed, using the actual checkout and default Turbopack build           |
| `go test ./...`                                                                                     | Passed                                                                  |
| `python3 web/scripts/contract-types.py`                                                             | Succeeded; projections regenerated and formatted                        |
| `git diff --check`                                                                                  | Passed                                                                  |
| `PLAYWRIGHT_SERVER=production PLAYWRIGHT_PORT=3202 npm run test:e2e -- --workers=2 --reporter=json` | 380 passed                                                              |
| Documented live backend startup                                                                     | Passed after wildcard fix; development database connected               |
| Direct API / proxied health                                                                         | Direct liveness 200; proxy not repeated after fix                        |
| Live authenticated writes and reload persistence                                                    | Unverified; no development account or database writes used              |

## Browser coverage

All 102 workspace screens/states in the inventory have route-rendering, visible-image loading, page-error, and document-overflow checks at **1728px, 1024px and 375px**. All 35 dialog routes are exercised at these widths with initial Cancel focus and Escape dismissal. These are smoke checks, not assertions that every distinct preview state perfectly matches Figma.

Existing authentication tests cover entry navigation, registration, server-role redirects, MFA retries, field errors, rate limiting, cookie propagation/clearing, keyboard recovery, responsive layouts, and reduced motion. Primary authentication state screenshots are captured at 1728px and 375px; existing responsive interaction tests also cover narrower desktop sizes.

The added integration tests exercise edited-field profile PATCH payloads, preserved values after failure, assessment start identities, ordered save-before-submit behavior, expired-session recovery, submitted receipts, direct-dialog navigation, manual stage decision payloads and conflicts. Existing integration checks cover feed filters/cursors, application envelopes, posting CRUD/lifecycle, modules and ordered pipeline-stage writes. **Successful HTTP responses in these tests are mocks.**

DAG interaction tests at 1728px, 1024px and 375px verify explicit joins, cycle rejection, keyboard node movement, zoom, and disabled graph persistence. The DAG editor was retained.

Screenshots are generated under `web/tests/artifacts/product/<node-id>-<width>.png`, with `-dialog.png` variants before modal dismissal. This directory is intentionally gitignored. The browser report records each passing case without adding binary screenshots or traces to source control.

## Visual and service limits

Figma MCP enumerated every page and collection container. It returned fresh design context and screenshots for 63 product frames before the Education-plan quota was exhausted. The exact 51 missing references are listed in [COVERAGE.md](COVERAGE.md). Separate motion context is also unverified for 11 retrieved frames that flag animation.

Representative screenshot/code comparisons informed posting-card layout, unillustrated hiring headings, profile overview/edit layouts, typography, and responsive behavior. **Full pixel-perfect visual parity remains unverified.** Existing fixtures supply variable product content and cannot substitute for complete visual review or real backend data.

The real API starts after the pipeline-progress wildcard fix. JOB management/application handlers and candidate pipeline-progress authentication remain fail-closed. Other missing operations are documented in [INTEGRATION.md](../INTEGRATION.md). The route name and its regression test were the only backend source changes; no database data was changed.

## Earlier failures and their resolution

- The checkout initially contained inconsistent generated `.next` route types. A fresh standard production build resolved them; the final type-check passes.
- Earlier browser runs exposed the existing implicit textarea-label behavior and an overflow regression introduced by rendering new nonnumeric dashboard statistics. Both were fixed and pass in the final run.
- A new confirmation test initially used “Deactivate account” instead of the Figma dialog's “Deactivate” button label. The test now uses the actual accessible name.
- The isolated development run exposed an intermittent Next.js development-router initialization error. A separately launched production server was later terminated during testing, producing connection-refused failures. The final test configuration lets Playwright own the built production server, and all 380 checks pass. Neither earlier run is counted as final verification.
