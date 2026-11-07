# Authentication screens

The Figma Login Page canvas is implemented at `/`, `/login`, and `/register`.
`?account=company` selects the hiring workspace; the default is candidate.

`AuthShell`, `FormField`, and the shared `Brand` component compose the frames.
The CSS uses Figma's Funnel Sans / Funnel Display fonts, local SVG assets, and
shared color tokens in `app/globals.css`. The fourteen-second gradient and glow
animations stop when reduced motion is requested. Auth content that extends
past Figma's 1117px frame stays reachable by scrolling.

The account selector derives its active pill, text, and icon colors from the
current workspace. The brand badge uses the same placement on every auth screen.
Background animation stays at the same phase between auth routes. Tab switches
keep email and password in memory and clear workspace-specific registration fields.
Motion includes a sliding account indicator, short content entrances, a floating
illustration, and hover/press feedback; reduced motion disables these effects.

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
