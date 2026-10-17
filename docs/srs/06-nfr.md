## 6. Non-functional requirements

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

Objective: one sentence — what this category means for Amara.
Standards/refs: OWASP ASVS §2 | GDPR Art.17 | WCAG 2.2 | ...
Key questions (from research): 2-4 bullets.

| ID | Requirement | Jira |
|----|-------------|------|
| NFR-...-001 | The system *SHALL* ... | TBD | -->
### 1. Security Requirements

#### Objective

Stop unauthorized access, data theft, and tampering of Company, Candidate, and pipeline data at a level a student team can test in CI.

#### Standards/Refs

OWASP Top 10:2021 (A01–A10). OWASP ASVS 4.0 Level 1 for the product, Level 2 for authentication and for any endpoint that can return identity fields. OWASP password guidance (Argon2id).

#### Key Questions

- Which Top 10 item does this requirement close, and how do we fail a build if it regresses?
- Who is allowed to see a Candidate's identity, and when (`FR-INVW-071`, `FR-INVW-073`)?
- What limits protect sign-in from repeated failed attempts?
- Which automated security checks must pass before release?

#### Requirements
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-SEC-001 | The system *SHALL* deny access by default. | TBD |
| NFR-SEC-002 | The system *SHALL* serve all browser and API traffic over TLS 1.2 or later. | TBD |
| NFR-SEC-003 | The system *SHALL* use parameterized queries for all database access. | TBD |
| NFR-SEC-004 | Each epic *SHALL* have a written threat model covering its assets, actors, abuse cases, and the NFR that addresses each abuse case before its first production deployment. | TBD |
| NFR-SEC-005 | Production containers *SHALL* run as non-root. | TBD |
| NFR-SEC-006 | CI *SHALL* run `govulncheck` on the Go module. | TBD |
| NFR-SEC-007 | The system *SHALL* lock an account for 15 minutes after 5 failed sign-in attempts from that account. | TBD |
| NFR-SEC-008 | Container images *SHALL* be built only in CI from the Go and frontend lockfiles (`go.sum`, `package-lock.json`). | TBD |
| NFR-SEC-013 | The system *SHALL* check every API caller's role against the AUTH permission matrix before returning protected data. | TBD |
| NFR-SEC-014 | An automated test *SHALL* exercise every API route with an anonymous caller, a caller with the wrong role, and a caller with the correct role. | TBD |
| NFR-SEC-015 | A release *SHALL* have zero confirmed OWASP A01 (Broken Access Control) findings. | TBD |
| NFR-SEC-016 | The system *SHALL* reject plaintext HTTP requests, except to redirect them to HTTPS. | TBD |
| NFR-SEC-017 | The system *SHALL* encrypt Candidate personal data and secrets at rest using AES-256 or the cloud provider's equivalent. | TBD |
| NFR-SEC-018 | The system *SHALL* hash passwords with Argon2id using at least 64 MiB of memory, at least 3 iterations, and parallelism 1. | TBD |
| NFR-SEC-019 | The system *SHALL* store passwords in a non-reversible form. | TBD |
| NFR-SEC-020 | The system *SHALL NOT* interpolate untrusted input into a shell command string. | TBD |
| NFR-SEC-021 | The system *SHALL* contextually encode untrusted input before rendering it as HTML. | TBD |
| NFR-SEC-022 | A release *SHALL* have zero confirmed OWASP A03 (Injection) findings from SAST and DAST scans run in CI. | TBD |
| NFR-SEC-023 | The sign-in flow *SHALL* be covered by at least one abuse-case test before the first production deployment. | TBD |
| NFR-SEC-024 | Production containers *SHALL NOT* include default credentials. | TBD |
| NFR-SEC-025 | The public origin *SHALL* send the `Strict-Transport-Security` response header. | TBD |
| NFR-SEC-026 | The public origin *SHALL* send the `Content-Security-Policy` response header. | TBD |
| NFR-SEC-027 | The public origin *SHALL* send the `X-Content-Type-Options: nosniff` response header. | TBD |
| NFR-SEC-028 | The public origin *SHALL* send the `Referrer-Policy` response header. | TBD |
| NFR-SEC-029 | A CI header check *SHALL* fail when any required security header is missing from the public origin. | TBD |
| NFR-SEC-030 | CI *SHALL* run `npm audit` against the frontend dependencies and report high- and critical-severity findings. | TBD |
| NFR-SEC-031 | A dependency with a known high- or critical-severity CVE *SHALL* be patched or explicitly waived within 14 days of the advisory. | TBD |
| NFR-SEC-032 | A waiver for a high- or critical-severity dependency CVE *SHALL* name an owner and include a written reason. | TBD |
| NFR-SEC-033 | The system *SHALL* rate-limit sign-in to 10 attempts per IP address per minute. | TBD |
| NFR-SEC-034 | An authenticated Hiring Manager session *SHALL* end after 30 minutes of inactivity. | TBD |
| NFR-SEC-035 | An authenticated Hiring Manager session *SHALL* end after 12 hours, regardless of activity. | TBD |
| NFR-SEC-036 | The system *SHALL* regenerate session identifiers when a user signs in. | TBD |
| NFR-SEC-037 | CI *SHALL* tag each container image with the Git commit that produced it. | TBD |
| NFR-SEC-038 | The deployment process *SHALL NOT* deploy a container image built outside CI. | TBD |
| NFR-SEC-039 | The system *SHALL NOT* deserialize untrusted data using a Go decoder that executes code. | TBD |
| NFR-SEC-040 | CI *SHALL* run SAST and DAST scans. | TBD |
| NFR-SEC-041 | The pipeline anonymity flow *SHALL* be covered by at least one abuse-case test before the first production deployment. | TBD |
| NFR-SEC-042 | The offer flow *SHALL* be covered by at least one abuse-case test before the first production deployment. | TBD |
| NFR-SEC-043 | Passwords *SHALL NOT* be written to logs. | TBD |
| NFR-SEC-044 | Each API route that returns company-scoped data *SHALL* enforce the caller's company scope defined by AUTH. | TBD |
| NFR-SEC-045 | An automated test *SHALL* verify that a user associated with one company cannot access company-scoped data belonging to another company. | TBD |

**Drives:** AUTH permission middleware, Argon2id hasher, sign-in rate limiter, CI security jobs, secure container configuration, and security-header checks.

### 2. Privacy Requirements

#### Objective

Handle Candidate and Company personal data lawfully for an Egypt launch, and keep a path to stricter rules if the product later serves users abroad. Anonymity is a disclosure rule on a pipeline, not a reason to omit identity from the account. Photo, gender, and the other identity fields are stored on the account and stay in the system. Before the reveal point, they are hidden from Hiring Managers and Interviewers (`CON-004`, `CON-010`, `FR-INVW-071`).

#### Standards/Refs

[Egypt Personal Data Protection Law No. 151 of 2020](https://mcit.gov.eg/Upcont/Documents/Reports%20and%20Documents_1232021000_Law_No_151_2020_Personal_Data_Protection.pdf). Where that law does not publish a numeric threshold, GDPR is the measurable stand-in: Art. 5 (minimization), Art. 15 (access), Art. 17 (erasure), Art. 20 (portability), Art. 22 (automated decisions), Art. 33 (breach notice, 72 hours). Pipeline anonymity is disclosure control, not deletion: `CON-004`, `CON-010`, `FR-INVW-070` to `FR-INVW-074`.

#### Key Questions

- Which identity fields does every account store, and which views must hide them before the reveal point?
- How long do we keep a rejected Candidate, and who can order deletion?
- Can a Hiring Manager or Interviewer see stored identity fields in a pipeline view, export, or scorecard before the reveal point?
- Which third parties (email, meeting, LLM) receive personal data, and may they train on it?

#### Requirements
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-PRIV-001 | The system *SHALL* store account identity fields, including name, photo, email, phone, date of birth, gender, nationality, and address, and *SHALL* keep them when pipeline anonymity is enabled. | TBD |
| NFR-PRIV-002 | Before account creation, the system *SHALL* show a privacy notice in Arabic and English that states what personal data is collected, why it is collected, who receives it, and how long it is retained. | TBD |
| NFR-PRIV-003 | The system *SHALL* let a Candidate download their profile, applications, and released feedback as JSON within 7 days of the request. | TBD |
| NFR-PRIV-004 | The system *SHALL* delete a Candidate's personal data within 30 days of a verified erasure request, except records subject to a recorded legal hold or records the Company is required to retain to document a hiring decision. | TBD |
| NFR-PRIV-005 | Records retained to document a hiring decision *SHALL* omit the Candidate's name, contact details, photo, and free-text answers, and *SHALL* retain only the decision, time, and decider role. | TBD |
| NFR-PRIV-006 | Pipeline anonymity *SHALL NOT* be treated as deletion of the Candidate's stored account identity fields. | TBD |
| NFR-PRIV-007 | When anonymity is enabled, hiring views, exports, resume copies, dashboards, scorecards, and notes *SHALL* show the anonymous identifier instead of the Candidate's identity fields before the reveal point (FR-INVW-071). | TBD |
| NFR-PRIV-008 | The Candidate's own account view *SHALL* continue to show the Candidate's stored identity fields while pipeline anonymity is enabled. | TBD |
| NFR-PRIV-009 | At the identity-reveal point, the same hiring views *SHALL* show the Candidate's identity fields and *SHALL NOT* hide them again (FR-INVW-073). | TBD |
| NFR-PRIV-010 | An integration test *SHALL* fail if a pre-reveal hiring response contains a Candidate identity field. | TBD |
| NFR-PRIV-011 | An integration test *SHALL* fail if pipeline anonymity removes a Candidate identity field from the stored account record. | TBD |
| NFR-PRIV-012 | Within 72 hours of confirming a personal-data breach, the system *SHALL* produce a notice stating what happened, which data was affected, and how many people were affected. | TBD |
| NFR-PRIV-013 | A personal-data breach notice *SHALL* be deliverable to the team lead and, where required by law, to affected people. | TBD |
| NFR-PRIV-014 | Production personal data *SHALL* be stored in one selected region (see `NFR-ENV-002`). | TBD |
| NFR-PRIV-015 | Before personal data is transferred outside the selected region, a written decision *SHALL* name the recipient, purpose, and safeguard. | TBD |
| NFR-PRIV-016 | Before a third-party processor is enabled in production, a register *SHALL* identify the data that processor receives. | TBD |
| NFR-PRIV-017 | AI prompts *SHALL NOT* include Candidate identity data, including name, photo, email, phone, personal links, date of birth, gender, nationality, address, or national ID. | TBD |
| NFR-PRIV-018 | A vendor that will not agree to the personal-data training restriction in `NFR-PRIV-026` *SHALL NOT* receive Candidate or Company personal data. | TBD |
| NFR-PRIV-019 | Logs and traces *SHALL NOT* contain session tokens, national IDs, or full resumes. | TBD |
| NFR-PRIV-020 | A log scrubber *SHALL* redact email-address and phone-number patterns. | TBD |
| NFR-PRIV-021 | A test *SHALL* fail if a fixture secret appears in a captured log line. | TBD |
| NFR-PRIV-022 | The system *SHALL* delete personal data for an application with final status Rejected, Withdrawn, Screened Out, Offer Declined, or Offer Expired 24 months after the status was set, unless a legal hold is recorded. | TBD |
| NFR-PRIV-023 | A scheduled job *SHALL* delete application personal data that has reached the retention period. | TBD |
| NFR-PRIV-024 | A third-party processor that receives Candidate or Company personal data *SHALL* agree not to use it to train models. | TBD |
| NFR-PRIV-025 | A field inventory *SHALL* list each personal-data field stored by the system, its purpose, and its retention period. | TBD |
| NFR-PRIV-026 | Account creation *SHALL* require an explicit consent action, stored with a timestamp and the privacy-notice version. | TBD |
| NFR-PRIV-027 | A Candidate data download *SHALL* be completed in the product and *SHALL NOT* be emailed in a file that contains another person's data. | TBD |

**Drives:** account identity fields, a hiring-view mapper that substitutes the anonymous identifier before the reveal point, erasure job, processor register, and log redaction.

### 3. Capacity Requirements

#### Objective

Hold the data and the concurrent use of an Egypt-wide hiring product, from a pilot a student team can pay for up to large and small Egyptian employers, without a rewrite.

#### Standards/Refs

Planning targets below. Egypt launch context: formal private-sector employers plus startups, Candidates mostly on phones, and deadline-day bursts. Posting deadlines and application caps are configured by JOB. Video-recorded responses (`FR-INVW-008`) dominate storage cost and are called out separately.

#### Key Questions

- What is the T1 number we must load-test before Milestone delivery, and what is only a growth target?
- Which feature blows the student budget if we store it naively (video)?
- How many applications must one posting be able to store, separate from the cap the Hiring Manager configures in JOB?
- What happens to screening and email when a deadline lands on the hour?

#### Requirements
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-CAP-001 | At T1, the system *SHALL* store data for at least 50 companies without manual sharding. | TBD |
| NFR-CAP-002 | The data model and tenant key *SHALL* allow growth to T2 (1,000 companies, 500,000 Candidate accounts, 10,000,000 applications) and T3 (10,000 companies, 2,000,000 Candidate accounts, 100,000,000 applications) by adding capacity, without changing identifiers or dropping tenant isolation. | TBD |
| NFR-CAP-003 | The system *SHALL* store at least 2,000 applications on one posting at T1, 10,000 at T2, and 25,000 at T3, including each retained resume copy. This is a storage ceiling, not the Hiring Manager's application cap configured by JOB. | TBD |
| NFR-CAP-004 | The system *SHALL* support a pipeline of at least 50 stages and 100 links per posting at T1. | TBD |
| NFR-CAP-005 | At T1, object storage *SHALL* hold retained resume copies and video responses for the retention window in `NFR-PRIV-024`; the planning size is 5 TB. | TBD |
| NFR-CAP-006 | At T1, the system *SHALL* sustain 500 concurrent signed-in users. | TBD |
| NFR-CAP-007 | The notification sender *SHALL* support a T1 volume of 100,000 emails per month. | TBD |
| NFR-CAP-008 | At T1, hot logs and traces *SHALL* be retained for 14 days. | TBD |
| NFR-CAP-009 | At T1, the system *SHALL* store at least 250 Hiring Manager and Interviewer accounts. | TBD |
| NFR-CAP-010 | At T1, the system *SHALL* store at least 25,000 Candidate accounts. | TBD |
| NFR-CAP-011 | At T1, the system *SHALL* store at least 500 active postings. | TBD |
| NFR-CAP-012 | At T1, the system *SHALL* store at least 250,000 applications. | TBD |
| NFR-CAP-013 | At T1, the relational database *SHALL* store the required data within 50 GB. | TBD |
| NFR-CAP-014 | At T1, the system *SHALL* render and validate a pipeline at the capacity in `NFR-CAP-004` without timing out the request (FR-INVW-005, FR-INVW-006). | TBD |
| NFR-CAP-015 | At T1, video objects older than 30 days *SHALL* move to a lower-cost storage tier. | TBD |
| NFR-CAP-016 | Objects past the retention period in `NFR-PRIV-024` *SHALL* be deleted by a scheduled job. | TBD |
| NFR-CAP-017 | At T1, the system *SHALL* sustain a burst of 10 application submissions per second for 15 minutes. | TBD |
| NFR-CAP-018 | The system *SHALL NOT* lose an application submission during the T1 burst in `NFR-CAP-017`. | TBD |
| NFR-CAP-019 | The notification sender *SHALL* queue messages rather than drop them when the provider is slow, so bounded INVW notifications still receive a recorded delivery attempt. | TBD |
| NFR-CAP-020 | At T1, a log-ingest cost cap *SHALL* stop additional log ingestion, rather than stop the product, if a day's ingest exceeds 2 GB. | TBD |

**Drives:** tenant-keyed tables, object-storage lifecycle rules, submission and email queues, and a log-ingest cap.

### 4. Compatibility Requirements

#### Objective

Work in the browsers, networks, and languages an Egyptian candidate and a Hiring Manager actually use, and with the external services the product already calls.

#### Standards/Refs

Browser support: current and previous major version, matching common SaaS practice. Core Web Vitals and a 3G-class profile for the apply flow. Arabic and English per `NFR-COMP-004`; detailed locale rules are tracked in `docs/srs/TBD.md` (TBD-006). HTTP APIs as JSON, versioned.

#### Key Questions

- Which browsers must pass the apply flow, and which do we explicitly refuse?
- Does the apply flow finish on a cheap Android phone over a slow mobile link?
- What happens when the meeting provider or the email provider is down?
- How long do we keep an old API version alive for a mobile client we already shipped?

#### Requirements
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-COMP-001 | The apply flow, sign-in, and Hiring Manager posting form *SHALL* work on the current and previous major versions of Chrome, Edge, Firefox, and Safari, plus Chrome for Android and Safari for iOS. | TBD |
| NFR-COMP-002 | The apply flow, sign-in, and Hiring Manager posting form *SHALL* be usable at viewport widths from 360 px to 1920 px. | TBD |
| NFR-COMP-003 | On a throttled profile of 1.5 Mbps down, 750 Kbps up, and 300 ms round-trip latency, the apply form *SHALL* be visible and submittable within 8 seconds. | TBD |
| NFR-COMP-004 | The UI *SHALL* be available in Arabic (`ar-EG`, right to left) and English (`en`). | TBD |
| NFR-COMP-005 | The HTTP API *SHALL* return JSON under a version prefix (`/v1`). | TBD |
| NFR-COMP-006 | Failure of the meeting or email provider *SHALL NOT* block application submission, screening, or a Hiring Manager decision. | TBD |
| NFR-COMP-007 | The system *SHOULD* accept resume uploads as PDF or DOCX files up to 10 MB. | TBD |
| NFR-COMP-008 | Internet Explorer *SHALL NOT* be supported. | TBD |
| NFR-COMP-009 | A browser compatibility matrix *SHALL* be rerun before each milestone. | TBD |
| NFR-COMP-010 | The apply flow *SHALL NOT* require hover interactions at a viewport width of 360 px. | TBD |
| NFR-COMP-011 | The apply flow *SHALL NOT* require horizontal scrolling at a viewport width of 360 px. | TBD |
| NFR-COMP-012 | A feature that requires a faster network connection than the profile in `NFR-COMP-003` (such as live video or a large upload) *SHALL* show a named error instead of a blank page. | TBD |
| NFR-COMP-013 | Dates, times, and numbers *SHALL* follow the viewer's locale. | TBD |
| NFR-COMP-014 | The system *SHALL* store timestamps in UTC. | TBD |
| NFR-COMP-015 | The UI *SHALL* display timestamps in the viewer's time zone, defaulting to `Africa/Cairo` (FR-INVW-029). | TBD |
| NFR-COMP-016 | A breaking HTTP API change *SHALL* increment the API version. | TBD |
| NFR-COMP-017 | The previous HTTP API version *SHALL* remain available for at least 6 months after a new version is announced. | TBD |
| NFR-COMP-018 | The system *SHALL* retry meeting-provider and email-provider failures according to the INVW error-handling rules. | TBD |
| NFR-COMP-019 | The system *SHALL* show in-product notifications even when email delivery fails. | TBD |
| NFR-COMP-020 | The system *SHALL* reject resume uploads in unsupported formats and show a message naming the allowed formats. | TBD |

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
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-REL-001 | At T1, the apply and sign-in paths *SHALL* be available 99.0% of each calendar month, measured by an external HTTPS check every 60 seconds from outside the hosting account. Three consecutive failed checks count as an outage. | TBD |
| NFR-REL-002 | At T2, the apply and sign-in paths *SHOULD* be available 99.5% of each calendar month; at T3 they *SHOULD* be available 99.9%, measured as in `NFR-REL-001`. | TBD |
| NFR-REL-003 | At T1, the system *SHALL* recover the apply and sign-in paths within 24 hours of a confirmed outage (RTO 24 h). | TBD |
| NFR-REL-004 | At T2, the system *SHOULD* meet an RPO of 1 hour. | TBD |
| NFR-REL-005 | A restore from backup *SHALL* be executed and timed at least once per academic term at T1, once every 6 months at T2, and once per quarter at T3. The run *SHALL* record the time to a working sign-in and one readable application. | TBD |
| NFR-REL-006 | Every production container *SHALL* expose a health check. | TBD |
| NFR-REL-007 | The system *SHALL* acknowledge an application submission only after the application row and retained resume copy are committed. | TBD |
| NFR-REL-008 | Module answers *SHALL* be saved at least every 30 seconds and when connectivity drops. | TBD |
| NFR-REL-009 | If the screening queue is older than 60 seconds, the system *SHALL* continue accepting submissions and *SHALL NOT* drop queued work. | TBD |
| NFR-REL-010 | Planned maintenance *SHALL* occur Friday 22:00–Saturday 06:00 `Africa/Cairo`. | TBD |
| NFR-REL-011 | At T1, a Sev-1 incident that takes down the apply or sign-in path *SHALL* be acknowledged within 4 business hours. | TBD |
| NFR-REL-012 | At T1, the system *SHALL* have an RPO of no more than 24 hours. | TBD |
| NFR-REL-013 | At T1, the system *SHALL* make automated daily backups of the database and object storage. | TBD |
| NFR-REL-014 | At T2, the system *SHOULD* recover the apply and sign-in paths within 4 hours (RTO 4 h). | TBD |
| NFR-REL-015 | At T2, the system *SHOULD* use a standby database in a second availability zone and continuous backup. | TBD |
| NFR-REL-016 | At T3, the system *SHOULD* have an RPO of no more than 5 minutes. | TBD |
| NFR-REL-017 | At T3, the system *SHOULD* recover the apply and sign-in paths within 1 hour (RTO 1 h). | TBD |
| NFR-REL-018 | The platform *SHALL* restart a production container after it fails 3 consecutive health checks. | TBD |
| NFR-REL-019 | A deployment *SHALL NOT* send traffic to a container that has not passed its health check. | TBD |
| NFR-REL-020 | A retry of the same application submission *SHALL NOT* create a second application. | TBD |
| NFR-REL-021 | A client retry after an application-submission timeout *SHALL* be safe. | TBD |
| NFR-REL-022 | After reconnecting, a Candidate *SHALL* be able to resume an assessment with the last saved answers. | TBD |
| NFR-REL-023 | When an assessment time limit expires, the system *SHALL* submit saved answers even if the browser has closed (FR-INVW-012). | TBD |
| NFR-REL-024 | If the screening queue is older than 60 seconds, the system *SHALL* show the Hiring Manager a delayed-screening state. | TBD |
| NFR-REL-025 | The system *SHALL* record the completion time for every screening job, including jobs completed after the target in `NFR-PERF-005` or `NFR-PERF-011`. | TBD |
| NFR-REL-026 | Planned maintenance *SHALL* be announced in the product at least 72 hours in advance. | TBD |
| NFR-REL-027 | Planned maintenance *SHALL NOT* overlap a stored posting deadline. | TBD |
| NFR-REL-028 | At T2, a Sev-1 incident that takes down the apply or sign-in path *SHALL* be acknowledged within 30 minutes during Sunday–Thursday, 09:00–18:00 `Africa/Cairo`. | TBD |
| NFR-REL-029 | An incident note *SHALL* record the incident start, detection time, mitigation, and any data loss. | TBD |

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
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-MAINT-001 | A new developer *SHALL* be able to run the backend, frontend, and database with one documented command (`docker compose`) and sign in to a seeded local account within 30 minutes of a clean clone, on an 8 GB RAM laptop without a cloud account. | TBD |
| NFR-MAINT-002 | CI *SHALL* run formatting, lint, type checks, and tests on every pull request. | TBD |
| NFR-MAINT-003 | Automated tests *SHALL* cover at least 70% of statements in backend packages overall. | TBD |
| NFR-MAINT-004 | Schema changes *SHALL* be stored as versioned migration files. | TBD |
| NFR-MAINT-005 | Every HTTP request *SHALL* carry a trace identifier that is returned in the response header and written on every log line for that request. | TBD |
| NFR-MAINT-006 | The HTTP API *SHALL* have an OpenAPI description checked into the repository. | TBD |
| NFR-MAINT-007 | A significant design choice (datastore, queue, authentication mechanism, or hosting) *SHALL* have an architecture decision record in the repository before dependent code is merged. | TBD |
| NFR-MAINT-008 | Production configuration *SHALL* come from environment variables rather than values compiled into a container image. | TBD |
| NFR-MAINT-009 | ANA code *SHALL NOT* write AUTH, CAND, JOB, or INVW tables. | TBD |
| NFR-MAINT-010 | A pull request with a failing CI check *SHALL NOT* be merged. | TBD |
| NFR-MAINT-011 | A CI run *SHALL* finish within 15 minutes on the student runner quota. | TBD |
| NFR-MAINT-012 | Automated tests *SHALL* cover at least 85% of statements in backend packages that implement scoring, anonymity, and offer state. | TBD |
| NFR-MAINT-013 | CI *SHALL* report backend test coverage and fail a pull request if either coverage threshold in `NFR-MAINT-003` or `NFR-MAINT-012` is missed. | TBD |
| NFR-MAINT-014 | The deployment process *SHALL* apply versioned schema migrations. | TBD |
| NFR-MAINT-015 | Each schema migration *SHALL* be reversible by one down migration for one release. | TBD |
| NFR-MAINT-016 | A production schema change *SHALL NOT* depend on a person running SQL by hand. | TBD |
| NFR-MAINT-017 | Logs *SHALL* be structured JSON. | TBD |
| NFR-MAINT-018 | Given a submission ID, an operator *SHALL* be able to find its request, decision, and provider errors without reading source code. | TBD |
| NFR-MAINT-019 | A pull request that changes an HTTP request or response shape *SHALL* update the OpenAPI description in the same pull request. | TBD |
| NFR-MAINT-020 | An architecture decision record *SHALL* name the options considered and explain why the selected choice fits the student budget. | TBD |
| NFR-MAINT-021 | The same container image *SHALL* run in local, CI, staging, and production environments. | TBD |
| NFR-MAINT-022 | A review check or import rule *SHALL* fail if an analytics package imports another epic's write repository. | TBD |

**Drives:** `docker-compose` dev stack, CI workflow, coverage gate, migration tool, OpenTelemetry middleware, OpenAPI file, and ADR folder.

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
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-PORT-001 | The backend and frontend *SHALL* ship as OCI container images. | TBD |
| NFR-PORT-002 | Cloud resources *SHALL* be described as code that can target AWS or Azure. | TBD |
| NFR-PORT-003 | Object storage *SHALL* be accessed through an interface in the Go module. | TBD |
| NFR-PORT-004 | A full logical export of one tenant (company, postings, applications, decisions, and files) *SHALL* be producible as SQL or JSON plus object files within 7 days of a request. | TBD |
| NFR-PORT-005 | Container images *SHOULD* be built for `linux/amd64` and `linux/arm64`. | TBD |
| NFR-PORT-006 | Replacing the LLM vendor *SHOULD NOT* require a schema change. | TBD |
| NFR-PORT-007 | OCI images *SHALL* run unmodified on the local Docker engine, in CI, and on the chosen cloud container service. | TBD |
| NFR-PORT-008 | OCI images *SHALL NOT* depend on a host path or a host-installed toolchain. | TBD |
| NFR-PORT-009 | A production cloud resource *SHALL NOT* require manual creation in a cloud console without an equivalent infrastructure-as-code definition. | TBD |
| NFR-PORT-010 | Email delivery *SHALL* be accessed through an interface in the Go module. | TBD |
| NFR-PORT-011 | The meeting provider *SHALL* be accessed through an interface in the Go module. | TBD |
| NFR-PORT-012 | The LLM provider *SHALL* be accessed through an interface in the Go module. | TBD |
| NFR-PORT-013 | Replacing an external provider *SHALL* require only a configuration change and an adapter, without changes to application, pipeline, or offer logic. | TBD |
| NFR-PORT-014 | A tenant export *SHALL NOT* contain rows or files belonging to another tenant. | TBD |
| NFR-PORT-015 | LLM prompts and model names *SHALL* be stored in configuration. | TBD |
| NFR-PORT-016 | A second LLM adapter *SHOULD* be demonstrable in one sprint if a vendor fails the privacy requirement in `NFR-PRIV-026`. | TBD |

**Drives:** Dockerfiles, Terraform modules with AWS and Azure targets, provider interfaces, tenant export job, and multi-architecture CI build.

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
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-SCAL-001 | The application tier *SHALL* be stateless. | TBD |
| NFR-SCAL-002 | Screening work *SHALL* run on workers fed by a queue. | TBD |
| NFR-SCAL-003 | The system *SHALL* absorb 10 times the trailing-hour submission rate for 60 minutes. | TBD |
| NFR-SCAL-004 | Analytics and dashboard reads *SHOULD* be able to use a read replica or a read model. | TBD |
| NFR-SCAL-005 | At T1, monthly cloud spend for production *SHALL* stay within US$150 or within the student-credit balance if that is lower. | TBD |
| NFR-SCAL-006 | Before each milestone that claims a capacity target, a load test at that target's concurrency *SHALL* be run and its result stored in the repository. | TBD |
| NFR-SCAL-007 | At T2, the web tier *SHOULD* add an instance within 3 minutes after CPU usage exceeds 70% for 5 minutes. | TBD |
| NFR-SCAL-008 | Session state *SHALL NOT* live in application-process memory. | TBD |
| NFR-SCAL-009 | Adding a second application instance *SHALL NOT* require a code change or break sign-in or submission. | TBD |
| NFR-SCAL-010 | Email delivery *SHALL* run on workers fed by a queue. | TBD |
| NFR-SCAL-011 | AI feedback generation *SHALL* run on workers fed by a queue. | TBD |
| NFR-SCAL-012 | The web process *SHALL NOT* perform screening, email delivery, or AI feedback generation inline. | TBD |
| NFR-SCAL-013 | Worker count *SHALL* be raisable without redeploying the web process. | TBD |
| NFR-SCAL-014 | The system *SHALL NOT* drop a committed application submission during the load in `NFR-SCAL-003`. | TBD |
| NFR-SCAL-015 | Sign-in *SHALL* remain available during the load in `NFR-SCAL-003`. | TBD |
| NFR-SCAL-016 | Excess application submissions *SHALL* wait in a queue rather than fail an HTTP request that has already committed the submission. | TBD |
| NFR-SCAL-017 | An analytics query *SHALL NOT* take a write lock on the application-submission path. | TBD |
| NFR-SCAL-018 | If the T1 production cloud-spend limit in `NFR-SCAL-005` is exceeded, an ADR *SHALL* record the overrun and the resources turned off or resized. | TBD |
| NFR-SCAL-019 | Video lifecycle management (`NFR-CAP-015`) and log-ingest limits (`NFR-CAP-020`) *SHALL* be the first cost-reduction levers. | TBD |
| NFR-SCAL-020 | A missed latency target *SHALL* be fixed or documented as a known gap with an owner. | TBD |
| NFR-SCAL-021 | At T2, the web tier *SHOULD* scale in when CPU usage stays below 30% for 15 minutes. | TBD |
| NFR-SCAL-022 | Scale-in *SHALL NOT* drop an in-flight request. | TBD |

**Drives:** JWT or shared session store, queue and worker process, read replica for ANA, budget alarm, load-test script in CI or a documented manual run.

### 9. Usability Requirements

#### Objective

Let a first-time Candidate apply, and a first-time Hiring Manager post and build a short pipeline, without training and on a phone.

#### Standards/Refs

ISO 9241-11 (effectiveness, efficiency, satisfaction). System Usability Scale (SUS), average benchmark 68. Nielsen heuristics for error recovery and consistency. Glossary terms only (`docs/srs/appendix-a-glossary.md`).

#### Key Questions

- How long does a first application take when the resume is already on the platform?
- Can the Hiring Manager undo a bad drag on the pipeline canvas?
- Are status words the same in email, dashboard, and the Candidate view?
- What does the empty state tell a new company to do next?

#### Requirements
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-USAB-001 | A first-time Candidate with a complete profile *SHALL* be able to submit an application in 5 minutes or less in a moderated usability test. | TBD |
| NFR-USAB-002 | A first-time Hiring Manager *SHALL* be able to create a posting with a title, description, deadline, and one Must-have criterion in 10 minutes or less. | TBD |
| NFR-USAB-003 | A Hiring Manager *SHALL* be able to build a 5-stage pipeline, connect it, and reach a publishable state in 10 minutes or less. | TBD |
| NFR-USAB-004 | Each destructive action (reject, withdraw, override, or publish a migration) *SHALL* ask for confirmation and name its consequence. | TBD |
| NFR-USAB-005 | Status labels shown to a Candidate, a Hiring Manager, and in email *SHALL* match the status values defined by INVW. | TBD |
| NFR-USAB-006 | A user action *SHALL* show a visible acknowledgement within 100 ms. | TBD |
| NFR-USAB-007 | A System Usability Scale (SUS) test with at least 5 Candidates and 5 Hiring Managers *SHOULD* score at least 68 for each role before the milestone that declares the apply and posting flows complete. | TBD |
| NFR-USAB-008 | An empty Company account *SHOULD* show the next required action (verify Company, create a posting, or publish a pipeline) instead of only an empty chart. | TBD |
| NFR-USAB-009 | Each error *SHALL* state what failed and the next action, using terminology consistent with the INVW error-handling rules. | TBD |
| NFR-USAB-010 | A pull request that introduces an application status *SHALL* update the status definitions in INVW and the glossary in the same pull request. | TBD |
| NFR-USAB-011 | An operation that takes more than 1 second *SHALL* show progress or a busy state. | TBD |
| NFR-USAB-012 | An operation that takes more than 1 second *SHALL NOT* appear frozen. | TBD |
| NFR-USAB-013 | At least 4 of 5 first-time Candidates in a moderated application test *SHALL* complete the test without facilitator assistance. | TBD |
| NFR-USAB-014 | The pipeline canvas *SHALL* support undo for at least 10 editing actions. | TBD |
| NFR-USAB-015 | The pipeline canvas *SHALL* provide a non-drag method to add and connect a stage. | TBD |

**Drives:** usability test script before the milestone, canvas undo stack, shared status enum with the glossary, and empty-state components.

### 10. Accessibility Requirements

#### Objective

Make the apply, posting, and pipeline flows usable with a keyboard, a screen reader, and low vision, to WCAG 2.2 AA, including on the drag-and-drop canvas. None of these are T1. *Should* is the T2 bar. *Could* is later.

#### Standards/Refs

WCAG 2.2 Level AA. axe-core for the automated slice. Egypt Law No. 10 of 2018 on the rights of persons with disabilities is the local reason to keep these requirements, but they are not in the T1 build. Timing rules interact with `FR-INVW-011` and `FR-INVW-015`.

#### Key Questions

- Can every canvas action be done from the keyboard?
- Is pass/reject/current stage visible without color?
- Can a Candidate use a time-limited module with a keyboard and screen reader?
- What do we run in CI, and what still needs a person with a screen reader?

#### Requirements
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-ACC-001 | The apply flow, sign-in, posting form, dashboard, and Candidate application view *SHOULD* conform to WCAG 2.2 Level AA. | TBD |
| NFR-ACC-002 | Every action, including adding, connecting, and deleting a pipeline stage, *SHOULD* be operable by keyboard alone. | TBD |
| NFR-ACC-003 | Text *SHOULD* have a contrast ratio of at least 4.5:1; large text and UI components *SHOULD* have a contrast ratio of at least 3:1. | TBD |
| NFR-ACC-004 | Pages *COULD* remain usable at 200% text zoom. | TBD |
| NFR-ACC-006 | The UI *COULD* honor `prefers-reduced-motion`. | TBD |
| NFR-ACC-007 | Interactive controls *SHOULD* have an accessible name. | TBD |
| NFR-ACC-008 | Touch and click targets on the apply flow *COULD* be at least 24 by 24 CSS pixels. | TBD |
| NFR-ACC-009 | An axe-core (or equivalent) scan in CI *SHOULD* report zero critical or serious accessibility violations on the routes in `NFR-ACC-001`. | TBD |
| NFR-ACC-010 | A manual accessibility test using a keyboard and one screen reader (NVDA or VoiceOver) *SHOULD* be recorded before each milestone. | TBD |
| NFR-ACC-011 | Focus *SHOULD* be visible for all keyboard-operable controls. | TBD |
| NFR-ACC-012 | The user interface *SHOULD NOT* contain a keyboard trap. | TBD |
| NFR-ACC-013 | The pipeline canvas *SHOULD* provide the non-drag alternative in `NFR-USAB-015`. | TBD |
| NFR-ACC-014 | The apply flow *COULD* reflow at a 320 px width without requiring two-dimensional scrolling. | TBD |
| NFR-ACC-015 | Motion that plays longer than 5 seconds *COULD* include a pause control. | TBD |
| NFR-ACC-016 | Status changes that occur without a page load, including stage advances, offer responses, and notification arrivals, *SHOULD* be announced to assistive technology. | TBD |
| NFR-ACC-017 | The page language *SHOULD* match the selected locale and include `dir="rtl"` for Arabic. | TBD |
| NFR-ACC-018 | Pipeline state (passed, rejected, current, and held) *SHOULD* use an icon or text as well as color. | TBD |

**Drives:** axe-core CI job, keyboard stage editor, status icon set, locale `lang` and `dir` on the document.

### 11. Performance Requirements

#### Objective

Keep the apply and hiring flows fast on a phone in Egypt, and keep the deadlines already promised in the functional requirements under T1 load.

#### Standards/Refs

Core Web Vitals (LCP, INP, CLS). Latency is server time plus network, at the percentile stated, at the load in `NFR-CAP-006` and `NFR-CAP-017` unless a row says otherwise. Functional deadlines already set in INVW include 10-minute notifications (for example, FR-INVW-027, FR-INVW-076, and FR-INVW-112).

#### Key Questions

- Which percentile do we promise, and at which concurrency?
- What is the budget for the first load on a slow mobile link, separate from API latency?
- Which functional deadline is actually a performance NFR we must load-test?
- How do we stop a regression shipping?

#### Requirements
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-PERF-001 | On the apply page and sign-in page, the 75th percentile of field loads *SHALL* meet LCP of at most 2.5 s, INP of at most 200 ms, and CLS of at most 0.1 on desktop broadband. | TBD |
| NFR-PERF-002 | Read APIs used by the posting form and Candidate application view *SHALL* respond within 300 ms for 95% of requests and within 800 ms for 99% of requests at T1 concurrency, excluding the client's network. | TBD |
| NFR-PERF-003 | Job search *SHALL* return the first page of results within 2 seconds for 95% of requests at 500 concurrent users. | TBD |
| NFR-PERF-004 | Application submission *SHALL* return an acknowledgement within 1 second for 95% of requests at the burst in `NFR-CAP-017`, after the durable write in `NFR-REL-007`. | TBD |
| NFR-PERF-005 | Must-have screening *SHALL* finish within 30 seconds for 95% of submissions at T1 load. | TBD |
| NFR-PERF-006 | Email and in-product notifications with a 10-minute functional deadline *SHALL* be attempted within 2 minutes for 95% of events at T1 volume. | TBD |
| NFR-PERF-007 | The Hiring Manager dashboard *SHALL* show stage counts for a posting of 1,000 applications within 3 seconds for 95% of loads. | TBD |
| NFR-PERF-008 | The gzipped JavaScript for the apply route *SHALL* be at most 200 KB. | TBD |
| NFR-PERF-009 | A Lighthouse performance run on the apply route *SHOULD* be stored each milestone. A drop of more than 5 points from the previous milestone *SHOULD* be fixed or accepted in writing by the epic owner before release. | TBD |
| NFR-PERF-011 | Must-have screening *SHALL* finish within 60 seconds for 100% of submissions when its queue is not delayed (NFR-REL-009). | TBD |
| NFR-PERF-012 | Email and in-product notifications with a 10-minute functional deadline *SHALL* be attempted within 10 minutes for 100% of events, except during a provider outage; failed attempts *SHALL* be retried per the INVW error-handling rules. | TBD |
| NFR-PERF-013 | The pipeline canvas *SHALL* render 50 stages within 2 seconds. | TBD |
| NFR-PERF-014 | List endpoints *SHALL* be paginated with a default page size of at most 50 and *SHALL NOT* return an unbounded collection. | TBD |

**Drives:** CDN or static hosting for Next.js assets, pagination on list queries, screening worker pool sized for T1, bundle-size check.

### 12. Environmental Requirements

Operating/deployment environment (platforms, browsers, runtime). Sustainability concerns live in `docs/srs/07-global-considerations.md`.

#### Objective

Pin the runtimes, the region, and the cost envelope so the team can build, deploy, and survive a bad Cairo network without a surprise bill.

#### Standards/Refs

Go and Node.js release lines current as of this draft. Container platform on AWS or Azure. Region chosen for latency from Cairo and for the residency rules in `NFR-PRIV-016` and `NFR-PRIV-017`. Browser and network bar from section 4.

#### Key Questions

- Which Go and Node versions does CI enforce, and who bumps them?
- Which region keeps Cairo latency acceptable without moving personal data casually?
- What is the cheapest production shape that still meets T1?
- What must the client do when the network drops mid-module?

#### Requirements
| ID | Requirement | Jira |
|----|-------------|------|
| NFR-ENV-001 | CI *SHALL* build the backend with a pinned Go 1.23 or later patch release. | TBD |
| NFR-ENV-002 | Production *SHALL* run in one region selected for Cairo users, on AWS or Azure container hosting. | TBD |
| NFR-ENV-003 | Production *SHALL* run only as containers. | TBD |
| NFR-ENV-004 | The repository *SHALL* define local, CI, staging, and production environments that use the same container image. | TBD |
| NFR-ENV-005 | The client *SHALL* preserve module answers through a dropped connection using the autosave rule in `NFR-REL-008`. | TBD |
| NFR-ENV-006 | T1 production *SHALL* be sized to the spend cap in `NFR-SCAL-005`. | TBD |
| NFR-ENV-008 | CI *SHALL* build the frontend with a pinned Node.js 22 LTS release. | TBD |
| NFR-ENV-009 | A runtime version bump *SHALL* update the Dockerfiles and CI image in the same pull request. | TBD |
| NFR-ENV-010 | Local development *SHALL* use the same major Go and Node.js versions as CI. | TBD |
| NFR-ENV-011 | A pre-deployment measurement *SHALL* confirm the selected production region has a median HTTPS round trip from Cairo of at most 100 ms. | TBD |
| NFR-ENV-012 | At T1, production personal data *SHALL NOT* be replicated to a second region. | TBD |
| NFR-ENV-013 | Production *SHALL NOT* depend on a manually configured virtual machine that deployment does not recreate from code. | TBD |
| NFR-ENV-014 | Base container images *SHALL* be pinned by digest in the Dockerfile. | TBD |
| NFR-ENV-015 | Staging *SHALL* use synthetic data only. | TBD |
| NFR-ENV-016 | Production personal data *SHALL NOT* be copied to local or CI environments. | TBD |
| NFR-ENV-017 | The client *SHALL* show an offline state within 5 seconds of a failed request. | TBD |
| NFR-ENV-018 | AI inference *SHALL* use an external API. | TBD |
| NFR-ENV-019 | A developer machine or T1 server *SHALL NOT* be required to host a GPU. | TBD |
| NFR-ENV-020 | A future change to Egypt's daylight-saving rule *SHALL* be picked up from the platform time-zone database rather than a hardcoded offset. | TBD |

**Drives:** pinned Dockerfiles, a one-region Terraform stack, a staging project with synthetic data, an offline banner, and zone-database time formatting.
