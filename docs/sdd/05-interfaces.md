## 5. Human interfaces

Owner: TBD | Status: Draft | Last updated: 2026-10-14 | Jira: TBD

### 5.1 Purpose and reference

This section describes the human-facing interaction surfaces and their ownership. It is a concise interaction contract, not a wireframe or visual-design specification. The use cases referenced here are in the repository's consolidated [`docs/srs/04-use-cases.md`](../srs/04-use-cases.md); the request's suggested path `03-use-cases.md` is not present in this repository.

### 5.2 Role-based surfaces

| User / surface | Main interactions | Representative SRS use cases |
|---|---|---|
| Guest / account entry | Register Candidate/Interviewer or Company Representative, activate accounts, enter login, and request password reset. Company verification may remain pending while registry/admin review completes. | `UC-AUTH-001`–`UC-AUTH-003`, `UC-AUTH-010`, `UC-AUTH-015` |
| Candidate workspace | Browse/filter postings; inspect details; apply; track applications and assessments; read permitted feedback and final decisions; respond to offers; view crests, verified skills, and dynamic resume; save postings and follow companies. | `UC-CAND-001`–`UC-CAND-014`; `UC-INVW-014`, `UC-INVW-016`, `UC-INVW-017`; `UC-JOB-16` |
| Hiring Manager workspace | Manage company postings; view applicants and screened-out reasons; design, publish, and migrate pipelines; configure tests, rubrics, and feedback; review dashboard/insights; make permitted advancement decisions; manage Finalists and offers. | `UC-JOB-01`–`UC-JOB-07`, `UC-JOB-16`–`UC-JOB-17`; `UC-INVW-001`–`UC-INVW-009`, `UC-INVW-015`–`UC-INVW-017` |
| Interviewer workspace | Publish live-interview availability, open assigned interview details, write private notes, and submit rubric scorecards. Company association and assignment scope are enforced by AUTH/INVW. | `UC-INVW-004`; `UC-AUTH-001`, `UC-AUTH-003`, `UC-AUTH-011` |
| System Administrator workspace | Review company registration and reports; manage users, entities, taxonomies, constants, and permitted aggregate reports. Administrative actions go through the owning module's validation and audit/operational controls. | `UC-AUTH-013`–`UC-AUTH-017` |

### 5.3 Interaction rules

- **Authentication and return navigation:** Protected actions direct unauthenticated users to login/registration and return them to the requested route after successful authentication, subject to AUTH activation and MFA rules (`UC-AUTH-003`; `UC-JOB-09`–`UC-JOB-11`).
- **Role and company scope:** Navigation and server responses expose only operations and records the user's role and associated company permit. Hiding a button is not an authorization check; the backend remains authoritative (AUTH FR-AUTH-004, FR-AUTH-016, FR-AUTH-019, FR-AUTH-041; `UC-AUTH-003`).
- **Candidate privacy:** Candidate pages show only the signed-in Candidate's applications and permitted results. Hiring-team pages apply INVW anonymity until the reveal point. Interviewer notes are not shown to Candidates (`UC-INVW-008`, `UC-INVW-014`, FR-INVW-033, FR-INVW-071–074, FR-INVW-097).
- **Statuses and decisions:** Display status values from the owning state machine—JOB for posting lifecycle and INVW for admitted-application/offer states. Candidate actions such as offer acceptance submit a command to INVW; the client displays the resulting authoritative state (`UC-JOB-03`–`UC-JOB-06`, `UC-INVW-006`, `UC-INVW-017`).
- **Feedback and errors:** State what happened and the next permitted action. Show pending evaluation/provider results as pending rather than converting them to a failure or score of zero. Notification failure must not hide a committed decision; in-product notices remain available (INVW error-handling rules; `UC-INVW-009`, `UC-INVW-010`, `UC-INVW-017`).
- **Localization and accessibility:** Provide Arabic RTL and English, localize dates/numbers, preserve the user's posting text and currency, and meet the NFR browser, viewport, keyboard, screen-reader, contrast, and notification-announcement targets (`NFR-COMP-001`–`NFR-COMP-015`, `NFR-ACC-001`–`NFR-ACC-004`, `NFR-ACC-006`–`NFR-ACC-018`).
- **Network resilience:** Candidate forms and timed assessments preserve entered answers according to autosave/recovery requirements and expose an offline or retry state instead of appearing frozen (`NFR-REL-008`, `NFR-REL-022`–`NFR-REL-023`, `NFR-ENV-005`, `NFR-ENV-017`).

### 5.4 Interface ownership

| Interaction | Owning epic / backend contract |
|---|---|
| Account entry, session, MFA, roles, and company scope | AUTH; shared `IAccessControl`. |
| Posting editor and company posting workspace | JOB, protected by AUTH. |
| Candidate discovery, profile and dynamic-resume presentation | CAND, with posting records supplied by JOB and account/profile operations supplied by AUTH. |
| Application submission and initial screening | JOB; successful admissions are handed to INVW. |
| Pipeline, assessment, interview, anonymity, post-admission status, finalist, offer, and notification interactions | INVW; `IPipelineData`, `IAmaraAwards`, and `IOffersAndDecisions`. |
| Hiring analytics/evaluation report display | ANA owns analytical output; INVW or other authorized surfaces consume it through `ICandidateScorecard`. |

The epic-level use-case diagrams, where provided, are available in [`src/`](src/) and [`img/`](img/). They complement but do not replace the textual cases and ownership rules in the SRS.
