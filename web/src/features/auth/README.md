# Authentication screens

The Figma Login Page canvas is implemented at `/`, `/login`, and `/register`.
`?account=company` selects the hiring workspace; the default is candidate.

`AuthShell`, `FormField`, and the shared `Brand` component compose the frames.
Login, registration, and their MFA/recovery states share the refined
`auth.module.css` styling. The shell and fields use the same styles directly;
form state, validation, focus handling, and password visibility are unchanged.
All auth views use local Funnel Sans / Funnel Display fonts and the existing SVG
assets, with no perpetual animation and support for reduced motion. The welcome
page retains its original design.

Account switches keep email and password in memory and clear workspace-specific
registration fields. Content exceeding a viewport stays reachable by scrolling.

`/` and `/login` read cookie presence inside a Suspense boundary. If an
`amara_session` cookie exists, `SessionEntry` validates it through
`GET /v1/auth/sessions/current` before showing public content. Valid sessions go
to the role's workspace; expired/revoked sessions show the public page. Service
failures show a retry message rather than pretending the visitor signed out.
The lookup returns only ID, role, and expiry, uses `Cache-Control: no-store`, and
does not issue/extend a cookie. Roles map only to fixed local destinations.
Workspace logos link directly to their role's home, including administration.

The backend keeps the existing 24-hour HttpOnly session cookie. Hiring-manager
and interviewer MFA requirements remain enforced when creating a session; an
already verified, valid session can be resumed without resubmitting credentials.
Deploy/restart the API with the frontend so the new GET endpoint is available.

Login posts to `/v1/auth/sessions`, including cookies. The server's role controls
the destination. MFA reuses the entered credentials in component memory;
passwords and tokens are never written to browser storage or URL parameters.
Errors are rendered from Problem Details or the existing JSON error shape.
The rate-limit state respects the server's `Retry-After` header.

Registration posts to the candidate or company endpoint. Candidate registration
opens the verification form after API success. The form sends the email and
64-character code to `POST /v1/auth/verifications`, displays rejected or expired
codes beside the field, and opens sign-in after the backend confirms activation.
The code stays in component memory and is cleared after success. A directly
opened pending state asks for the email again. Company registration still waits
for review and does not use the verification endpoint.

The following `state` query values make each designed state directly reviewable:

| Route                                            | State                | Figma frame |
| ------------------------------------------------ | -------------------- | ----------- |
| `/login`                                         | default              | `10:5`      |
| `/login?account=company`                         | default              | `4:6`       |
| `/register`                                      | default              | `11:51`     |
| `/register?account=company`                      | default              | `4:3`       |
| `/login?state=mfa`                               | MFA entry            | `268:14`    |
| `/login?state=code-rejected`                     | MFA retry            | `268:143`   |
| `/login?state=credentials-rejected`              | Credentials rejected | `268:272`   |
| `/login?state=rate-limited`                      | Rate limited         | `268:397`   |
| `/login?state=expired`                           | Session expired      | `268:522`   |
| `/register?state=verification-pending`           | Verification pending | `268:647`   |
| `/register?account=company&state=review-pending` | Review pending       | `268:772`   |

A directly opened MFA view has no stored credentials. Return to sign-in to enter
them before verification. The pending verification URL opens the real candidate
verification form; other URL states are presentation views and do not create a
session or bypass the API.

Run `npm run lint`, `npx tsc --noEmit`, `npm run build`, and `npm run test:e2e`
from `web/`. Install the browser once with `npx playwright install chromium`.
Browser tests mock API responses to cover validation, role navigation, MFA,
errors, and registration outcomes. They also check the twelve screens at
1728px and 375px, verify assets load, and save screenshots in the ignored
`tests/artifacts/` directory. The welcome page uses frame `149:119`.
