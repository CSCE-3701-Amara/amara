# Routing and authentication

## Public routes

The public route group contains:

- `/`: welcome/landing screen.
- `/login`: candidate sign-in; use `?account=company` for company sign-in.
- `/register`: candidate registration; use `?account=company` for company
  registration.

Auth screens also expose reviewable presentation states through an allowlisted
`state` query parameter: `mfa`, `code-rejected`, `credentials-rejected`,
`rate-limited`, `expired`, `verification-pending`, and `review-pending`.
These states do not create sessions, verify accounts, or bypass the API.

## Role destinations

After a successful login, the API response role determines the destination:

| Session role | Destination |
| --- | --- |
| Candidate | `/candidate/dashboard` |
| Company representative | `/employer/jobs` |
| Interviewer | `/interviewer/interviews` |

Administration routes are available under `/admin`, but they are not selected
by the public login redirect in the current client.

## Workspace navigation

Each protected route group applies a `Workspace` layout:

| Workspace | Main navigation |
| --- | --- |
| Candidate | Feed, applications, offers, profile |
| Employer | Job postings, pipelines, modules, interviews, grading, offers, profile |
| Interviewer | Interviews, assigned grading, availability, profile |
| Admin | Company registrations, accounts |

The workspace requests `DELETE /v1/auth/sessions/current` for logout. A
successful logout redirects to sign-in; a 401 redirects to the expired-session
view. Other failures remain visible in the workspace.

## Session and request rules

The browser does not read or write the session token. The backend sets an
HttpOnly cookie, and `apiRequest` sends it with `credentials: "include"`.
Passwords, MFA codes, and tokens stay in component memory and are never placed
in local storage, URLs, or analytics payloads.

All private API calls must use `/v1/...` relative paths. Next.js rewrites them
server-side using `API_URL`, preserving cookies and keeping the backend origin
out of browser code. Do not construct request URLs from untrusted API data.

The frontend may display a 401 recovery link, but authentication and
authorization remain backend responsibilities. A UI route being reachable
does not grant access to its data or mutation endpoints.

## Error and loading behavior

Private data uses no-store requests and displays a `BusyState` while loading.
`ApiFeedback` shows the server message, offers retry where available, and
offers “Sign in again” only for `ApiError` status 401. Rate-limit responses
respect `Retry-After`.