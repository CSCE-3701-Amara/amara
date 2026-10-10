## 7. Non-functional requirements

Owner: Ziad Eliwa | Status: Draft | Last updated: 2026-10-10 | Jira: AMARA-17 

Each NFR *SHALL* be measurable (a number, a condition, or a standard) and is identified by `NFR-<AREA>-nnn` (see `docs/conventions/requirements.md`).

Planning horizons used below. Numbers are targets to validate with real usage, not a claim about the Egyptian market size.

| Horizon | When | What it means |
|---|---|---|
| T1 | 6 months | Launch in Egypt. Must run on a student budget and student operations. |
| T2 | Year 1 | Paying small companies and startups. Still one team, still cost-conscious. |
| T3 | Year 2+ | Large Egyptian employers as well as small ones. National footprint. |
| T4 | After T3 | Expansion to the MENA region. |
| T5 | After T4 | International expansion beyond MENA. |

A *Shall* is required for T1. A *Should* is the T2 bar. A *Could* is the T3 bar. T4 and T5 are expansion horizons only. They are not yet written into the requirements below. The architecture *SHALL* not make a later horizon impossible (see `NFR-SCAL-001`).

Stack these requirements assume, and that the design must not fight: Go backend, React.js and Next.js frontend, Docker, AWS or Azure, structured logs and distributed traces.

--- 
<!-- Template for a category -->
<!-- ### n. <Category> Requirements

Objective:  one sentence — what this category means for Amara.
Standards/refs:  OWASP ASVS §2 | GDPR Art.17 | WCAG 2.2 | ...
Key questions (from research):  2-4 bullets.

| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-...-001 | The system shall ... | Shall | Draft | FR-...-... | AMARA-.. | -->
### 1. Security Requirements

#### Objective

Stop unauthorized access, data theft, and tampering of company, Interviewee, and pipeline data, at a bar a student team can test in CI.

#### Standards/Refs

OWASP Top 10:2021 (A01–A10). OWASP ASVS 4.0 Level 1 for the product, Level 2 for authentication and for any endpoint that can return identity fields. OWASP password guidance (Argon2id).

#### Key Questions

- Which Top 10 item does this requirement close, and how do we fail a build if it regresses?
- Who is allowed to see an Interviewee's identity, and when (`FR-INVW-071`, `FR-INVW-073`)?
- What do we log when access is denied, and who can delete those logs?
- Which outbound calls accept a URL we did not choose (`FR-INVW-026`, job-board integrations)?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-SEC-001 | The system *SHALL* deny access by default. Every API route *SHALL* check the caller's role against the AUTH permission matrix before returning data, and an automated test *SHALL* cover each route with an anonymous caller, a wrong-role caller, and a correct-role caller. A release *SHALL* have zero confirmed A01 (Broken Access Control) findings. | Shall | Draft | FR-INVW-033, FR-INVW-053, FR-INVW-071 | TBD |
| NFR-SEC-002 | The system *SHALL* serve all browser and API traffic over TLS 1.2 or later, *SHALL* reject plaintext HTTP except for a redirect to HTTPS, and *SHALL* encrypt Interviewee PII and secrets at rest with AES-256 or the cloud provider's equivalent. Passwords *SHALL* be hashed with Argon2id (memory at least 64 MiB, iterations at least 3, parallelism 1) and *SHALL* never be logged or stored reversible. | Shall | Draft | TBD (AUTH) | TBD |
| NFR-SEC-003 | The system *SHALL* use parameterized queries for all database access and *SHALL* reject or encode untrusted input before it reaches a query, a shell, or HTML. A release *SHALL* have zero confirmed A03 (Injection) findings from the SAST and DAST checks in CI. | Shall | Draft | FR-JOB-011, FR-INVW-017 | TBD |
| NFR-SEC-004 | Each epic *SHALL* have a written threat model (assets, actors, abuse cases, and the NFR that closes each abuse case) before its first production deploy. The anonymity path, the offer path, and sign-in *SHALL* each have at least one abuse-case test. | Shall | Draft | FR-INVW-071, REMOVED-FR-119 | TBD |
| NFR-SEC-005 | Production containers *SHALL* run as non-root, *SHALL* ship no default credentials, and *SHALL* send `Strict-Transport-Security`, `Content-Security-Policy`, `X-Content-Type-Options: nosniff`, and `Referrer-Policy`. A header check in CI *SHALL* fail the build if any of these is missing on the public origin. | Shall | Draft | cross-cutting | TBD |
| NFR-SEC-006 | CI *SHALL* run `govulncheck` on the Go module and `npm audit` at high severity on the frontend. A dependency with a known high or critical CVE *SHALL* be patched or explicitly waived, with a written reason and an owner, within 14 days of the advisory. | Shall | Draft | cross-cutting | TBD |
| NFR-SEC-007 | The system *SHALL* lock an account for 15 minutes after 5 failed sign-in attempts from that account, *SHALL* rate-limit sign-in to 10 attempts per IP per minute, and *SHALL* end a Hiring Manager session after 30 minutes idle or 12 hours absolute, whichever comes first. Session identifiers *SHALL* be regenerated on sign-in. | Shall | Draft | TBD (AUTH) | TBD |
| NFR-SEC-008 | Container images *SHALL* be built only in CI from a lockfile (`go.sum`, `package-lock.json`), *SHALL* be tagged with the git commit, and *SHALL* not be deployable from a laptop. The system *SHALL* not deserialize untrusted data into Go objects with a decoder that executes code. | Shall | Draft | cross-cutting | TBD |
| NFR-SEC-009 | The system *SHALL* write an append-only audit record for sign-in success and failure, role change, identity reveal, denied access to an anonymized field, offer create/withdraw, and permission override. The application database role *SHALL* not be able to update or delete audit records. Audit records *SHALL* be kept for at least 90 days. | Shall | Draft | FR-INVW-058, FR-INVW-071, FR-INVW-073 | TBD |
| NFR-SEC-010 | The server *SHALL* not fetch a URL supplied by a user unless that URL's host is on an allowlist for that feature (meeting provider, job board). Requests to link-local, loopback, and private address ranges *SHALL* be refused. A test *SHALL* cover a private-IP and a metadata-IP payload. | Shall | Draft | FR-INVW-026 | TBD |
| NFR-SEC-011 | Secrets (database passwords, API keys, signing keys) *SHALL* live only in the cloud secret store or in local env files that are gitignored. A CI secret scan *SHALL* fail the pull request if a secret pattern is committed. | Shall | Draft | cross-cutting | TBD |
| NFR-SEC-012 | Hiring Manager accounts *SHALL* support a second factor (TOTP or the provider's equivalent). T1 may ship with it optional; from T2 it *SHALL* be required for every Hiring Manager account. | Should | Draft | TBD (AUTH) | TBD |

**Drives:** AUTH permission middleware, Argon2id hasher, CI security jobs, audit table with a restricted DB role, outbound HTTP allowlist.

### 2. Privacy Requirements

#### Objective

Handle Interviewee and company personal data lawfully for an Egypt launch, and keep a path to stricter rules if the product later serves users abroad. Anonymity is a disclosure rule on a pipeline, not a reason to omit identity from the account. Photo, gender, and the other identity fields are stored on the account and stay in the system. Before the reveal point they are hidden from Hiring Managers and Interviewers only (`CON-004`, `CON-010`, `FR-INVW-071`).

#### Standards/Refs

[Egypt Personal Data Protection Law No. 151 of 2020](https://mcit.gov.eg/Upcont/Documents/Reports%20and%20Documents_1232021000_Law_No_151_2020_Personal_Data_Protection.pdf). Where that law does not publish a numeric threshold, GDPR is the measurable stand-in: Art. 5 (minimization), Art. 15 (access), Art. 17 (erasure), Art. 20 (portability), Art. 22 (automated decisions), Art. 33 (breach notice, 72 hours). Pipeline anonymity is disclosure control, not deletion: `CON-004`, `CON-010`, `FR-INVW-070` to `FR-INVW-074`.

#### Key Questions

- Which identity fields does every account store, and which views must hide them before the reveal point?
- How long do we keep a rejected Interviewee, and who can order deletion?
- Can a Hiring Manager or Interviewer see stored identity fields in a pipeline view, export, or scorecard before the reveal point?
- Which third parties (email, meeting, LLM) receive personal data, and may they train on it?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-PRIV-001 | The system *SHALL* store identity fields on the account, including name, photo, email, phone, date of birth, gender, nationality, and address, and *SHALL* keep them after anonymity is enabled. Ranking and insight generation *SHALL* not use name, photo, gender, age, nationality, or contact details as inputs (`FR-INVW-065`). A field inventory *SHALL* list every stored personal field, its purpose, and its retention. | Shall | Draft | FR-INVW-065, FR-INVW-071 | TBD |
| NFR-PRIV-002 | Before creating an account, the system *SHALL* show a privacy notice stating what is collected, why, who it is shared with, and how long it is kept, in Arabic and English. Account creation *SHALL* require an explicit consent action, stored with a timestamp and the notice version. | Shall | Draft | TBD (AUTH) | TBD |
| NFR-PRIV-003 | The system *SHALL* let an Interviewee download their profile, applications, and released feedback as JSON within 7 days of the request, and *SHALL* complete the download in the product without emailing a file that contains another person's data. | Shall | Draft | REMOVED-FR-110, TBD (CAND) | TBD |
| NFR-PRIV-004 | The system *SHALL* delete an Interviewee's personal data within 30 days of a verified erasure request, except records the company must keep to show a hiring decision was made. Those exception records *SHALL* drop name, contact details, photo, and free-text answers, and *SHALL* keep only the decision, time, and role of the decider. This deletion is not pipeline anonymity. | Shall | Draft | REMOVED-FR-015, FR-INVW-058 | TBD |
| NFR-PRIV-005 | The system *SHALL* delete personal data of an application whose final status is Rejected, Withdrawn, Screened Out, or Offer Declined 24 months after that status was set, unless a legal hold is recorded. A scheduled job *SHALL* do this. Pipeline anonymity *SHALL* not count as this deletion, and *SHALL* not remove the stored account fields. | Shall | Draft | REMOVED-FR-019, REMOVED-FR-021, FR-INVW-060 | TBD |
| NFR-PRIV-006 | When anonymity is enabled, the system *SHALL* keep the identity fields stored and *SHALL* hide them from Hiring Managers and Interviewers until the reveal point. Hiring views, exports, resume copies, the dashboard, scorecards, and notes *SHALL* show the anonymous identifier instead of name, photo, email, phone, personal links, date of birth, gender, nationality, and address (`FR-INVW-071`). The Interviewee's own account view *SHALL* still show the stored fields. At the reveal point the same hiring views *SHALL* show the stored fields, and *SHALL* not hide them again (`FR-INVW-073`). An integration test *SHALL* fail if a pre-reveal hiring response contains those fields, and *SHALL* fail if the stored account record no longer has them. | Shall | Draft | FR-INVW-070, FR-INVW-071, FR-INVW-073 | TBD |
| NFR-PRIV-007 | Where a rule can reject or screen out an Interviewee without a person deciding (`REMOVED-FR-021`, `FR-INVW-055`), the system *SHALL* tell the Interviewee that the outcome was automatic, *SHALL* let them request a human review, and *SHALL* record the human decision within 5 business days. The request path *SHALL* exist even when the default is automatic rejection. | Shall | Draft | REMOVED-FR-021, REMOVED-FR-027, FR-INVW-055 | TBD |
| NFR-PRIV-008 | The system *SHALL* detect a personal-data breach (confirmed unauthorized access to PII or to pre-reveal identity) and *SHALL* produce a notice containing what happened, which data, and how many people, within 72 hours of confirmation. The notice *SHALL* be deliverable to the team lead and, where the law requires it, to the affected person. | Shall | Draft | NFR-SEC-009 | TBD |
| NFR-PRIV-009 | Production personal data *SHALL* stay in one chosen region (see `NFR-ENV-002`). A transfer of personal data outside that region *SHALL* require a written decision naming the recipient, the purpose, and the safeguard, before the transfer is enabled. | Shall | Draft | cross-cutting | TBD |
| NFR-PRIV-010 | A third-party processor (email, meeting, object storage, LLM) *SHALL* not be enabled in production until a record names what data it receives and states that the vendor will not train models on Amara personal data. AI prompts *SHALL* not include name, email, phone, or national ID. If a vendor cannot accept that term, that vendor *SHALL* not receive those fields. | Shall | Draft | FR-INVW-018, FR-INVW-064, FR-INVW-078 | TBD |
| NFR-PRIV-011 | Logs and traces *SHALL* not contain passwords, session tokens, national IDs, or full resumes. A log scrubber *SHALL* redact email and phone patterns. A test *SHALL* fail if a fixture secret appears in a captured log line. | Shall | Draft | NFR-SEC-009 | TBD |

**Drives:** account profile that always stores identity fields, a hiring-view mapper that substitutes the anonymous identifier before the reveal point, erasure job, processor register, log redaction.

### 3. Capacity Requirements

#### Objective

Hold the data and the concurrent use of an Egypt-wide hiring product, from a pilot a student team can pay for up to large and small Egyptian employers, without a rewrite.

#### Standards/Refs

Planning targets below. Egypt launch context: formal private-sector employers plus startups, candidates mostly on phones, deadline-day bursts (`FR-JOB-003`, `FR-JOB-005`). Video-recorded responses (`FR-INVW-008`) dominate storage cost and are called out separately.

#### Key Questions

- What is the T1 number we must load-test before Milestone delivery, and what is only a growth target?
- Which feature blows the student budget if we store it naively (video)?
- How many applications must one posting be able to store, separate from the cap the Hiring Manager chooses (`FR-JOB-005`)?
- What happens to screening and email when a deadline lands on the hour?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-CAP-001 | At T1 the system *SHALL* store, without manual sharding, at least 50 companies, 250 Hiring Manager and Interviewer accounts, 25,000 Interviewee accounts, 500 active postings, and 250,000 applications, on a relational database of at most 50 GB. | Shall | Draft | REMOVED-FR-013, REMOVED-FR-015 | TBD |
| NFR-CAP-002 | The data model and tenant key *SHALL* allow growth to T2 (1,000 companies, 500,000 Interviewee accounts, 10,000,000 applications) and T3 (10,000 companies, 2,000,000 Interviewee accounts, 100,000,000 applications) by adding capacity, not by changing identifiers or dropping tenant isolation. | Shall | Draft | cross-cutting | TBD |
| NFR-CAP-003 | The system *SHALL* store at least 2,000 applications on one posting at T1, 10,000 at T2, and 25,000 at T3, including each retained resume copy. This is a storage ceiling. It is not the application cap in `FR-JOB-005`. That cap is a number the Hiring Manager chooses, and it may be any integer of at least 1, including a number below this ceiling. | Shall | Draft | FR-JOB-005, REMOVED-FR-015 | TBD |
| NFR-CAP-004 | The system *SHALL* keep a pipeline of at least 50 stages and 100 links per posting at T1, and *SHALL* render and validate it (`FR-INVW-005`, `FR-INVW-006`) without timing out the request. | Shall | Draft | FR-INVW-001, FR-INVW-005 | TBD |
| NFR-CAP-005 | Object storage *SHALL* hold retained resume copies and video responses for the retention window in `NFR-PRIV-005`. T1 planning size is 5 TB. Video objects older than 30 days *SHALL* move to a cheaper storage tier, and objects past retention *SHALL* be deleted by the same job as `NFR-PRIV-005`. | Shall | Draft | REMOVED-FR-015, FR-INVW-008 | TBD |
| NFR-CAP-006 | At T1 the system *SHALL* sustain 500 concurrent signed-in users and a burst of 10 application submissions per second for 15 minutes, with no lost submission and with the latency in `NFR-PERF-004`. | Shall | Draft | REMOVED-FR-013 | TBD |
| NFR-CAP-007 | The notification sender *SHALL* absorb T1 volume of 100,000 emails per month, and *SHALL* queue rather than drop when the provider is slow, so that the 10-minute delivery rules (`REMOVED-FR-018`, `REMOVED-FR-025`, `FR-INVW-027`) still have a recorded attempt. | Shall | Draft | REMOVED-FR-018, REMOVED-FR-025 | TBD |
| NFR-CAP-008 | Hot logs and traces *SHALL* be kept 14 days at T1. A cost cap *SHALL* stop ingest, rather than the product, if a day exceeds 2 GB. Searchable archive beyond 14 days is not required at T1. | Shall | Draft | NFR-SEC-009 | TBD |

**Drives:** tenant-keyed tables, object-storage lifecycle rules, submission queue, email queue, log sampling or a hard ingest cap.

### 4. Compatibility Requirements

#### Objective

Work in the browsers, networks, and languages an Egyptian candidate and a Hiring Manager actually use, and with the external services the product already calls.

#### Standards/Refs

Browser support: current and previous major version, matching common SaaS practice. Core Web Vitals and a 3G-class profile for the apply flow. Arabic and English per `docs/srs/10-i18n.md`. HTTP APIs as JSON, versioned.

#### Key Questions

- Which browsers must pass the apply flow, and which do we explicitly refuse?
- Does the apply flow finish on a cheap Android phone over a slow mobile link?
- What happens when the meeting provider or the email provider is down?
- How long do we keep an old API version alive for a mobile client we already shipped?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-COMP-001 | The apply flow, sign-in, and the Hiring Manager posting form *SHALL* work on the current and previous major versions of Chrome, Edge, Firefox, and Safari, and on Chrome for Android and Safari for iOS. Internet Explorer *SHALL* not be supported. A browser matrix *SHALL* be re-run before each milestone. | Shall | Draft | REMOVED-FR-013 | TBD |
| NFR-COMP-002 | Those flows *SHALL* be usable from a viewport width of 360 px to 1920 px. The apply flow *SHALL* not require hover or a horizontal scroll at 360 px. | Shall | Draft | REMOVED-FR-013, REMOVED-FR-115 | TBD |
| NFR-COMP-003 | The apply flow *SHALL* become usable (form visible and submittable) within 8 seconds on a throttled profile of 1.5 Mbps down, 750 Kbps up, and 300 ms round trip. Features that need a fast link (live video, large uploads) *SHALL* fail with a named reason, not a blank page. | Shall | Draft | REMOVED-FR-013 | TBD |
| NFR-COMP-004 | The UI *SHALL* ship Arabic (`ar-EG`, right to left) and English (`en`). Dates, times, and numbers *SHALL* follow the viewer's locale. Stored timestamps *SHALL* be UTC. Display *SHALL* use the viewer's time zone, defaulting to `Africa/Cairo` (`FR-INVW-029`). | Shall | Draft | FR-INVW-029 | TBD |
| NFR-COMP-005 | The HTTP API *SHALL* be JSON under a version prefix (`/v1`). A breaking change *SHALL* increment the version. The previous version *SHALL* keep working for at least 6 months after the new version is announced. | Shall | Draft | cross-cutting | TBD |
| NFR-COMP-006 | Failure of the meeting provider or the email provider *SHALL* not block application submission, screening, or a Hiring Manager decision. The system *SHALL* retry as specified in the INVW error-handling rules and *SHALL* show the in-product notification even when email fails. | Shall | Draft | FR-INVW-026, FR-INVW-027 | TBD |
| NFR-COMP-007 | The system *SHALL* accept a resume upload as PDF or DOCX up to 10 MB, and *SHALL* reject other types with a message that names the allowed types. | Should | Draft | TBD (CAND) | TBD |

**Drives:** browser test matrix in CI, locale middleware, `/v1` router, provider interfaces with a timeout and a fallback.

### 5. Reliability, Availability & Recoverability Requirements

#### Objective

Stay up often enough for a hiring deadline, lose no submitted application, and recover in a time the team can actually meet — tighter each horizon, not all at once.

#### Standards/Refs

Availability is monthly, measured by an external check, excluding a maintenance window announced 72 hours ahead. Egyptian work week is Sunday–Thursday; the window is Friday 22:00 to Saturday 06:00 `Africa/Cairo`. RPO is the maximum data loss. RTO is the maximum time to restore the apply and sign-in paths.

#### Key Questions

- What uptime can a student team honestly promise in year 1, on one region and no night rota?
- What must never be lost even if a container dies mid-submit?
- How do we prove a backup restores, not only that a backup job ran?
- Which dependency may fail without taking submissions down?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-REL-001 | At T1 the apply path and the sign-in path *SHALL* be available 99.0% of each calendar month (at most about 7.3 hours unavailable), measured by an external HTTPS check every 60 seconds from outside the hosting account. A failed check is three consecutive failures. | Shall | Draft | REMOVED-FR-013 | TBD |
| NFR-REL-002 | At T2 those paths *SHALL* be available 99.5% of each calendar month. At T3 they *SHALL* be available 99.9%. The measurement method *SHALL* stay the one in `NFR-REL-001`. | Should | Draft | REMOVED-FR-013 | TBD |
| NFR-REL-003 | At T1 the system *SHALL* recover the apply and sign-in paths within 24 hours of a confirmed outage (RTO 24 h) and *SHALL* lose at most 24 hours of data (RPO 24 h), using an automated daily backup of the database and of object storage. | Shall | Draft | REMOVED-FR-015 | TBD |
| NFR-REL-004 | At T2 the system *SHALL* meet RPO 1 hour and RTO 4 hours, using a standby database in a second availability zone and continuous backup. At T3 the system *SHALL* meet RPO 5 minutes and RTO 1 hour. | Should | Draft | cross-cutting | TBD |
| NFR-REL-005 | A restore from backup *SHALL* be executed and timed at least once per academic term at T1, once every 6 months at T2, and once per quarter at T3. The run *SHALL* record the time to a working sign-in and one readable application. A backup that has not been restored does not satisfy `NFR-REL-003`. | Shall | Draft | NFR-REL-003 | TBD |
| NFR-REL-006 | Every production container *SHALL* expose a health check. The platform *SHALL* restart a container that fails 3 consecutive checks. A deploy *SHALL* not send traffic to a container that is not healthy. | Shall | Draft | cross-cutting | TBD |
| NFR-REL-007 | The system *SHALL* acknowledge an application submission only after the application row and the retained resume copy are committed. A retry of the same submission *SHALL* not create a second application (`REMOVED-FR-016`). A client retry after a timeout *SHALL* be safe. | Shall | Draft | REMOVED-FR-013, REMOVED-FR-015, REMOVED-FR-016 | TBD |
| NFR-REL-008 | Module answers *SHALL* be saved at least every 30 seconds and on a connectivity drop. After a reconnect, the Interviewee *SHALL* resume with the last saved answers. An expired time limit *SHALL* still submit saved answers (`FR-INVW-012`) even if the browser has closed. | Shall | Draft | FR-INVW-011, FR-INVW-012 | TBD |
| NFR-REL-009 | If the screening queue is older than 60 seconds, the system *SHALL* keep accepting submissions, *SHALL* not drop queued work, and *SHALL* show the Hiring Manager a delayed-screening state. The 60-second rule in `REMOVED-FR-020` *SHALL* hold for 95% of submissions at T1 load (`NFR-CAP-006`); the rest *SHALL* complete, late, and be counted. | Shall | Draft | REMOVED-FR-020 | TBD |
| NFR-REL-010 | Planned maintenance *SHALL* fall in Friday 22:00–Saturday 06:00 `Africa/Cairo`, *SHALL* be announced in the product at least 72 hours ahead, and *SHALL* not be scheduled across a posting deadline that is already stored. | Shall | Draft | FR-JOB-003 | TBD |
| NFR-REL-011 | At T1 a Sev-1 incident (apply or sign-in down) *SHALL* be acknowledged within 4 business hours. At T2, within 30 minutes during Sunday–Thursday 09:00–18:00 `Africa/Cairo`. An incident note *SHALL* record start, detection, mitigation, and data loss. | Shall | Draft | NFR-REL-001 | TBD |

**Drives:** external uptime check, daily backup job, health checks on the container platform, idempotency key on submit, answer autosave, screening queue.

### 6. Maintainability Requirements

#### Objective

Let a rotating student team change, test, and hand over the system without a private laptop setup or an undocumented production.

#### Standards/Refs

Go and TypeScript community defaults (`gofmt`, `golangci-lint`, ESLint, Prettier, `tsc` strict). Twelve-factor config. OpenTelemetry for traces. Versioned schema migrations. Team conventions in `docs/conventions/`.

#### Key Questions

- Can a new teammate run the stack on a student laptop without a cloud account?
- What coverage is worth enforcing, and on which packages is a miss dangerous?
- How do we change the schema without a manual production step?
- Can we diagnose a failed submission from logs alone?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-MAINT-001 | A new developer *SHALL* be able to run the backend, frontend, and database with one documented command (`docker compose`) and sign in to a seeded local account within 30 minutes of a clean clone, on an 8 GB RAM laptop, with no cloud account. | Shall | Draft | cross-cutting | TBD |
| NFR-MAINT-002 | CI *SHALL* run format, lint, typecheck, and tests on every pull request. A pull request with a failing check *SHALL* not be merged. The CI run *SHALL* finish within 15 minutes on the student runner quota. | Shall | Draft | cross-cutting | TBD |
| NFR-MAINT-003 | Automated tests *SHALL* cover at least 70% of statements in backend packages overall, and at least 85% in the packages that implement screening, scoring, anonymity, and offer state. Coverage *SHALL* be reported in CI. A drop below the threshold *SHALL* fail the pull request. | Shall | Draft | REMOVED-FR-020, FR-INVW-016, FR-INVW-071 | TBD |
| NFR-MAINT-004 | Schema changes *SHALL* be versioned migration files, applied by the deploy, and reversible by one down migration for one release. A production change *SHALL* not depend on a person running SQL by hand. | Shall | Draft | cross-cutting | TBD |
| NFR-MAINT-005 | Every HTTP request *SHALL* carry a trace identifier, returned in the response header and written on every log line for that request. Logs *SHALL* be structured JSON. Given a submission id, an operator *SHALL* be able to find the request, the decision, and any provider error without reading source code. | Shall | Draft | REMOVED-FR-018, NFR-SEC-009 | TBD |
| NFR-MAINT-006 | The HTTP API *SHALL* have an OpenAPI description checked into the repo. A pull request that changes a request or response shape *SHALL* update that description in the same pull request. | Shall | Draft | NFR-COMP-005 | TBD |
| NFR-MAINT-007 | A significant design choice (datastore, queue, auth mechanism, hosting) *SHALL* have a short architecture decision record in the repo before the code that depends on it is merged. The record *SHALL* name the options considered and why the chosen one fits a student budget. | Shall | Draft | cross-cutting | TBD |
| NFR-MAINT-008 | Production config *SHALL* come from environment variables, not from values compiled into an image. The same image *SHALL* run in local, CI, staging, and production. | Shall | Draft | cross-cutting | TBD |
| NFR-MAINT-009 | ANA code *SHALL* not write AUTH, CAND, JOB, or INVW tables. A review check or an import rule *SHALL* fail if an analytics package imports a write repository of another epic. | Shall | Draft | TBD (ANA) | TBD |

**Drives:** `docker-compose` dev stack, CI workflow, coverage gate, migration tool, OpenTelemetry middleware, OpenAPI file, ADR folder.

### 7. Portability Requirements

#### Objective

Run the same product on a student laptop, in CI, and on either AWS or Azure, so a hosting choice or a vendor price change is not a rewrite.

#### Standards/Refs

OCI container images. Twelve-factor app. Terraform (or an equivalent that targets both AWS and Azure). S3-compatible object storage. No GPU on developer machines; model calls go to an API.

#### Key Questions

- What breaks if we move from AWS to Azure in a term?
- Can we develop through exams with no cloud bill?
- How hard is it to swap the LLM vendor if the first one is too expensive or refuses the privacy term?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-PORT-001 | The backend and the frontend *SHALL* ship as OCI images that run unmodified on the local Docker engine, in CI, and on the chosen cloud container service. Images *SHALL* not depend on a host path or a host-installed toolchain. | Shall | Draft | cross-cutting | TBD |
| NFR-PORT-002 | Cloud resources *SHALL* be described as code that can target AWS or Azure. A resource created only in a cloud console, with no code equivalent, *SHALL* not be required for production. | Shall | Draft | cross-cutting | TBD |
| NFR-PORT-003 | Object storage, email, the meeting provider, and the LLM *SHALL* sit behind an interface in the Go module. Switching provider *SHALL* be a configuration change plus an adapter, not a change to screening, pipeline, or offer code. | Shall | Draft | FR-INVW-026, FR-INVW-018 | TBD |
| NFR-PORT-004 | A full logical export of one tenant (company, postings, applications, decisions, files) *SHALL* be producible as SQL or JSON plus object files, within 7 days of a request, so a company can leave. The export *SHALL* not include another tenant's rows. | Shall | Draft | NFR-PRIV-003 | TBD |
| NFR-PORT-005 | Images *SHALL* be built for `linux/amd64` and `linux/arm64`, so they run on student ARM laptops and on cheaper ARM cloud instances. | Should | Draft | cross-cutting | TBD |
| NFR-PORT-006 | Replacing the LLM vendor *SHALL* not require a schema change. Prompts and model names *SHALL* live in configuration. A second adapter *SHALL* be demonstrable in one sprint if the privacy term in `NFR-PRIV-010` fails. | Should | Draft | FR-INVW-018, FR-INVW-078 | TBD |

**Drives:** Dockerfiles, Terraform modules with an AWS and an Azure target, provider interfaces, tenant export job, multi-arch CI build.

### 8. Scalability Requirements

#### Objective

Grow from the T1 pilot to Egyptian companies of every size by adding instances and queue workers, not by redesigning the request path.

#### Standards/Refs

Stateless app tier. Horizontal scaling. Queue for work that has a deadline measured in seconds or minutes, not milliseconds. Capacity numbers in section 3.

#### Key Questions

- Which work must stay on the request, and which must move to a queue before a deadline day?
- How do analytics reads avoid slowing submissions?
- What is the monthly cost ceiling at T1, and what do we turn off first if we hit it?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-SCAL-001 | The app tier *SHALL* be stateless. Session state *SHALL* not live in process memory. Adding a second instance *SHALL* not require a code change and *SHALL* not break sign-in or submission. | Shall | Draft | cross-cutting | TBD |
| NFR-SCAL-002 | Screening, email, and AI feedback generation *SHALL* run on workers fed by a queue. The web process *SHALL* not perform those jobs inline. Worker count *SHALL* be raisable without redeploying the web process. | Shall | Draft | REMOVED-FR-020, FR-INVW-076 | TBD |
| NFR-SCAL-003 | The system *SHALL* absorb 10 times the trailing-hour submission rate for 60 minutes by queueing, without dropping a committed submission and without taking sign-in down. Excess *SHALL* wait in the queue, not fail the HTTP call that already committed. | Shall | Draft | FR-JOB-005, NFR-CAP-006 | TBD |
| NFR-SCAL-004 | Analytics and dashboard reads *SHALL* be able to use a read replica or a read model. An analytics query *SHALL* not take a write lock on the application submission path. | Should | Draft | FR-INVW-061, TBD (ANA) | TBD |
| NFR-SCAL-005 | At T1, monthly cloud spend for production *SHALL* stay within US$150, or within the student-credit balance if that is lower, unless an ADR records the overrun and what was turned off or resized. Video lifecycle (`NFR-CAP-005`) and log ingest cap (`NFR-CAP-008`) are the first levers. | Shall | Draft | NFR-CAP-005, NFR-CAP-008 | TBD |
| NFR-SCAL-006 | Before each milestone that claims a capacity target, a load test at that target's concurrency *SHALL* be run and its result stored in the repo. A missed latency target *SHALL* be either fixed or written down as a known gap with an owner. | Shall | Draft | NFR-CAP-006, NFR-PERF-001 | TBD |
| NFR-SCAL-007 | At T2 the web tier *SHALL* add an instance within 3 minutes of CPU above 70% for 5 minutes, and *SHALL* scale in when CPU stays below 30% for 15 minutes. Scale-in *SHALL* not drop an in-flight request. | Should | Draft | NFR-SCAL-001 | TBD |

**Drives:** JWT or shared session store, queue and worker process, read replica for ANA, budget alarm, load-test script in CI or a documented manual run.

### 9. Usability Requirements

#### Objective

Let a first-time Interviewee apply, and a first-time Hiring Manager post and build a short pipeline, without training and on a phone.

#### Standards/Refs

ISO 9241-11 (effectiveness, efficiency, satisfaction). System Usability Scale (SUS), average benchmark 68. Nielsen heuristics for error recovery and consistency. Glossary terms only (`docs/srs/appendix-a-glossary.md`).

#### Key Questions

- How long does a first application take when the resume is already on the platform?
- Can the Hiring Manager undo a bad drag on the pipeline canvas?
- Are status words the same in email, dashboard, and the Interviewee view?
- What does the empty state tell a new company to do next?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-USAB-001 | A first-time Interviewee with a complete profile *SHALL* be able to submit an application in 5 minutes or less in a moderated test of at least 5 people. At least 4 of the 5 *SHALL* finish without help from the facilitator. | Shall | Draft | REMOVED-FR-013 | TBD |
| NFR-USAB-002 | A first-time Hiring Manager *SHALL* be able to create a posting with title, description, deadline, and one Must-have criterion in 10 minutes or less, in a moderated test of at least 5 people. | Shall | Draft | FR-JOB-001, FR-JOB-003, FR-JOB-007 | TBD |
| NFR-USAB-003 | A Hiring Manager *SHALL* be able to build a 5-stage pipeline, connect it, and reach a publishable state in 10 minutes or less. The canvas *SHALL* support undo of at least 10 editing actions and *SHALL* offer a non-drag way to add and connect a stage. | Shall | Draft | FR-INVW-001, FR-INVW-039 | TBD |
| NFR-USAB-004 | Every destructive action (reject, withdraw, override, publish a migration) *SHALL* ask for confirmation and *SHALL* name the consequence. Every error *SHALL* state what failed and the next action, in the same words as the INVW error-handling list. | Shall | Draft | REMOVED-FR-017, REMOVED-FR-019, FR-INVW-057 | TBD |
| NFR-USAB-005 | Status labels shown to an Interviewee, a Hiring Manager, and in email *SHALL* be the glossary status values. A pull request that introduces a new status string *SHALL* fail review unless the glossary changes in the same pull request. | Shall | Draft | REMOVED-FR-013, REMOVED-FR-110 | TBD |
| NFR-USAB-006 | Any user action *SHALL* show a visible acknowledgement within 100 ms. An operation that takes more than 1 second *SHALL* show progress or a busy state and *SHALL* not look frozen. | Shall | Draft | cross-cutting | TBD |
| NFR-USAB-007 | A moderated SUS test, at least 5 Interviewees and 5 Hiring Managers, *SHALL* score at least 68 on each role before the milestone that declares the apply and posting flows done. | Should | Draft | NFR-USAB-001, NFR-USAB-002 | TBD |
| NFR-USAB-008 | An empty company account *SHALL* show the next required action (verify company, create a posting, publish a pipeline) and *SHALL* not show an empty chart as the only content. | Should | Draft | REMOVED-FR-023 | TBD |

**Drives:** usability test script before the milestone, canvas undo stack, shared status enum with the glossary, empty-state components.

### 10. Accessibility Requirements

#### Objective

Make the apply, posting, and pipeline flows usable with a keyboard, a screen reader, and low vision, to WCAG 2.2 AA, including on the drag-and-drop canvas. None of these are T1. *Should* is the T2 bar. *Could* is later.

#### Standards/Refs

WCAG 2.2 Level AA. axe-core for the automated slice. Egypt Law No. 10 of 2018 on the rights of persons with disabilities is the local reason to keep these requirements, but they are not in the T1 build. Timing rules interact with `FR-INVW-011` and `FR-INVW-015`.

#### Key Questions

- Can every canvas action be done from the keyboard?
- Is pass/reject/current stage visible without color?
- Can an Interviewee extend a timed module, and are they warned before it ends?
- What do we run in CI, and what still needs a person with a screen reader?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-ACC-001 | The apply flow, sign-in, posting form, dashboard, and Interviewee application view *SHALL* conform to WCAG 2.2 Level AA. An axe-core (or equivalent) scan in CI *SHALL* report zero critical or serious violations on those routes. A manual pass with a keyboard and one screen reader (NVDA or VoiceOver) *SHALL* be recorded before each milestone. | Should | Draft | REMOVED-FR-013, REMOVED-FR-110 | TBD |
| NFR-ACC-002 | Every action, including adding, connecting, and deleting a pipeline stage, *SHALL* be operable by keyboard alone. Focus *SHALL* be visible. There *SHALL* be no keyboard trap. Drag and drop *SHALL* have the non-drag alternative required by `NFR-USAB-003`. | Should | Draft | FR-INVW-001 | TBD |
| NFR-ACC-003 | Text *SHALL* have a contrast ratio of at least 4.5:1, and large text and UI components at least 3:1. Pipeline state (passed, rejected, current, held) *SHALL* use an icon or text as well as color. | Should | Draft | FR-INVW-061, REMOVED-FR-112 | TBD |
| NFR-ACC-004 | Pages *SHALL* remain usable at 200% text zoom and *SHALL* reflow at a 320 px width without a two-dimensional scroll for the apply flow. | Could | Draft | NFR-COMP-002 | TBD |
| NFR-ACC-005 | A module time limit *SHALL* be extendable by the Interviewee at least once before it expires, unless the Hiring Manager has marked the module as non-extendable for integrity, in which case the Interviewee *SHALL* be told the limit before they start. The system *SHALL* warn at least 60 seconds before an in-session limit expires. | Could | Draft | FR-INVW-011, FR-INVW-015 | TBD |
| NFR-ACC-006 | The UI *SHALL* honor `prefers-reduced-motion` and *SHALL* not play motion longer than 5 seconds without a pause control. | Could | Draft | cross-cutting | TBD |
| NFR-ACC-007 | Interactive controls *SHALL* have an accessible name. Status changes that happen without a page load (stage advance, notification arrival, screening result) *SHALL* be announced to assistive technology. The page language *SHALL* be set, and *SHALL* switch with the locale, including `dir="rtl"` for Arabic. | Should | Draft | FR-INVW-079, NFR-COMP-004 | TBD |
| NFR-ACC-008 | Touch and click targets on the apply flow *SHALL* be at least 24 by 24 CSS pixels. | Could | Draft | NFR-COMP-002 | TBD |

**Drives:** axe-core CI job, keyboard stage editor, status icon set, locale `lang` and `dir` on the document.

### 11. Performance Requirements 

#### Objective

Keep the apply and hiring flows fast on a phone in Egypt, and keep the deadlines already promised in the functional requirements under T1 load.

#### Standards/Refs

Core Web Vitals (LCP, INP, CLS). Latency is server time plus network, at the percentile stated, at the load in `NFR-CAP-006` unless a row says otherwise. Functional deadlines already fixed: screening 60 seconds (`REMOVED-FR-020`), several notifications 10 minutes.

#### Key Questions

- Which percentile do we promise, and at which concurrency?
- What is the budget for the first load on a slow mobile link, separate from API latency?
- Which functional deadline is actually a performance NFR we must load-test?
- How do we stop a regression shipping?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-PERF-001 | On the apply page and the sign-in page, the 75th percentile of field loads *SHALL* meet LCP at most 2.5 s, INP at most 200 ms, and CLS at most 0.1, on desktop broadband. The same pages *SHALL* meet the 8-second usable-time rule in `NFR-COMP-003` on the slow profile. | Shall | Draft | REMOVED-FR-013 | TBD |
| NFR-PERF-002 | Read APIs used by the posting form and the Interviewee application view *SHALL* respond in at most 300 ms for 95% of requests and at most 800 ms for 99%, at T1 concurrency, excluding the client's network. | Shall | Draft | REMOVED-FR-110, REMOVED-FR-115 | TBD |
| NFR-PERF-003 | Job search *SHALL* return the first page of results within 2 seconds for 95% of requests at 500 concurrent users. | Shall | Draft | TBD (JOB) | TBD |
| NFR-PERF-004 | Application submission *SHALL* return an acknowledgement within 1 second for 95% of requests at the burst in `NFR-CAP-006`, after the durable write in `NFR-REL-007`. | Shall | Draft | REMOVED-FR-013, REMOVED-FR-018 | TBD |
| NFR-PERF-005 | Must-have screening *SHALL* finish within 30 seconds for 95% of submissions at T1 load, and within 60 seconds for 100% when the queue is not in the delayed state of `NFR-REL-009`. | Shall | Draft | REMOVED-FR-020 | TBD |
| NFR-PERF-006 | Email and in-product notifications that the functional requirements bound by 10 minutes *SHALL* be attempted within 2 minutes for 95% of events at T1 volume, and within 10 minutes for 100%, provider outage excepted and then retried per the INVW error-handling rules. | Shall | Draft | REMOVED-FR-018, REMOVED-FR-025, FR-INVW-027 | TBD |
| NFR-PERF-007 | The Hiring Manager dashboard *SHALL* show stage counts for a posting of 1,000 applications within 3 seconds for 95% of loads. The pipeline canvas *SHALL* render 50 stages within 2 seconds. | Shall | Draft | FR-INVW-061, NFR-CAP-004 | TBD |
| NFR-PERF-008 | The gzipped JavaScript for the apply route *SHALL* be at most 200 KB. List endpoints *SHALL* be paginated, default page size at most 50, and *SHALL* not return an unbounded collection. | Shall | Draft | REMOVED-FR-026, FR-INVW-062 | TBD |
| NFR-PERF-009 | A Lighthouse performance run on the apply route *SHALL* be stored each milestone. A drop of more than 5 points from the previous milestone *SHALL* be fixed or accepted in writing by the epic owner before release. | Should | Draft | NFR-PERF-001 | TBD |

**Drives:** CDN or static hosting for Next.js assets, pagination on list queries, screening worker pool sized for T1, bundle-size check.

### 12. Environmental Requirements 

Operating/deployment environment (platforms, browsers, runtime). Sustainability concerns live in `docs/srs/11-global-considerations.md`.

#### Objective

Pin the runtimes, the region, and the cost envelope so the team can build, deploy, and survive a bad Cairo network without a surprise bill.

#### Standards/Refs

Go and Node.js release lines current as of this draft. Container platform on AWS or Azure. Region chosen for latency from Cairo and for the residency rule in `NFR-PRIV-009`. Browser and network bar from section 4.

#### Key Questions

- Which Go and Node versions does CI enforce, and who bumps them?
- Which region keeps Cairo latency acceptable without moving personal data casually?
- What is the cheapest production shape that still meets T1?
- What must the client do when the network drops mid-module?

#### Requirements
| ID | Requirement | Priority | Status | Source story | Jira |
|----|-------------|----------|--------|--------------|------|
| NFR-ENV-001 | CI *SHALL* build the backend with a pinned Go 1.23 or later patch, and the frontend with a pinned Node.js 22 LTS. A version bump *SHALL* be a pull request that also updates the Dockerfiles and the CI image. Local development *SHALL* use the same major versions. | Shall | Draft | cross-cutting | TBD |
| NFR-ENV-002 | Production *SHALL* run in one region chosen for Cairo users, on AWS or Azure container hosting. The region *SHALL* show a median HTTPS round trip from Cairo of at most 100 ms in a pre-deploy measurement. Candidates to measure: Azure UAE North, Azure Qatar Central, AWS `me-central-1`, AWS `me-south-1`. Personal data *SHALL* not be replicated to a second region at T1. | Shall | Draft | NFR-PRIV-009 | TBD |
| NFR-ENV-003 | Production *SHALL* run only as containers. There *SHALL* be no manually configured virtual machine that the deploy does not recreate from code. Base images *SHALL* be pinned by digest in the Dockerfile. | Shall | Draft | NFR-PORT-001, NFR-PORT-002 | TBD |
| NFR-ENV-004 | The repo *SHALL* define four environments with the same image: local, CI, staging, production. Staging *SHALL* use synthetic data only. Production personal data *SHALL* not be copied to local or CI. | Shall | Draft | NFR-MAINT-008, NFR-PRIV-011 | TBD |
| NFR-ENV-005 | The client *SHALL* keep working through a dropped connection on the module-answer flow by the autosave rule in `NFR-REL-008`, and *SHALL* show an offline state within 5 seconds of a failed request rather than spinning with no message. | Shall | Draft | FR-INVW-012, NFR-REL-008 | TBD |
| NFR-ENV-006 | T1 production *SHALL* be sized to the spend cap in `NFR-SCAL-005`. AI inference *SHALL* be an external API. No developer machine and no T1 server *SHALL* be required to host a GPU. | Shall | Draft | FR-INVW-018, NFR-SCAL-005 | TBD |
| NFR-ENV-007 | All timestamps *SHALL* be stored in UTC. The default display zone *SHALL* be `Africa/Cairo`. A future change of Egypt's daylight-saving rule *SHALL* be picked up from the platform zone database, not from a hardcoded offset. | Shall | Draft | FR-INVW-029, NFR-COMP-004 | TBD |

**Drives:** pinned Dockerfiles, one-region Terraform stack, staging project with seed data, offline banner, zone-database time formatting.

### 13. Per-Epic NF Requirements 
#### 13.1 AUTH
#### 13.2 CAND
#### 13.3 JOB
#### 13.4 INVW
#### 13.5 ANA
