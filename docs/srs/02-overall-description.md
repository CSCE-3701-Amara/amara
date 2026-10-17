## 2. Overall description

Owner: TBD | Status: Draft | Last updated: 2026-10-14 | Jira: TBD

### 2.1 Product perspective

AMARA is a multi-company SaaS hiring platform with five cooperating functional epics. Candidates use the platform to discover and apply to postings, complete an employer's hiring process, review application outcomes, and maintain a professional profile and dynamic resume. Employers create postings and configure screening and interview workflows. System administrators support company and account trust operations. The system integrates with external services for email, company-registry checks, social login, online meetings, code execution, and selected analytics inputs.

The product is a new system rather than a replacement for an existing applicant-tracking system. Its module boundaries and data ownership are defined in `docs/srs/05-functional/`; its proposed implementation structure is described in `docs/sdd/`.

### 2.2 High-level product functions

| Epic | High-level functions |
|---|---|
| AUTH | Register and activate users; authenticate Candidates, Interviewers, and Company Representatives; enforce role/company-scoped permissions; verify companies; manage profiles, documents, reports, and administrative entities. |
| CAND | Browse, filter, and search opportunities; inspect postings; view submitted applications, progress, outcomes, and feedback; respond to offers; follow companies and save jobs; view earned crests and verified skills; create and view a dynamic resume. |
| JOB | Create and maintain postings; configure capacity, dates, criteria, and application questions; publish, archive, and republish postings; accept and preserve applications; screen them against Must-have criteria; admit qualifying applications to INVW. |
| INVW | Design and version a directed acyclic pipeline; configure assessment and interview stages; execute and score work; apply manual or configured advancement rules; manage anonymity; issue feedback and notifications; determine finalists; manage offers and final application states. |
| ANA | Ingest defined telemetry, calculate execution and pass-rate metrics, support behavioral and integrity analysis, generate evaluation reports, and audit score variance for bias. |

Functions are subject to each epic's SHALL/SHOULD/COULD priority. For example, the LLM-paired programming module is a future capability, while core posting, application, pipeline, and status flows are specified as mandatory.

### 2.3 User classes and characteristics

| User class | Goals and principal interactions |
|---|---|
| Guest | Register or begin company registration; view public information where available; proceed through sign-in gates when a protected action requires an account. |
| Candidate | Maintain a profile; find postings; submit applications; complete assessments and interviews; view progress, feedback, final decisions, offers, crests, verified skills, and the dynamic resume. Most candidate interactions are expected to work on mobile devices and variable networks. |
| Interviewer | Access assigned interviews, publish availability, conduct interviews, and submit notes and scorecards. Access is scoped to the associated company and assigned workflow. |
| Hiring Manager / Company Representative | Manage company postings, configure and operate pipelines, review candidate evidence and analytics, make permitted advancement decisions, and manage finalists and offers. Company access is company-scoped. |
| System Administrator | Review company registration requests and administer permitted accounts, entities, taxonomies, and aggregate administrative reports. |
| System and external actors | Background workers, email and meeting providers, official registries, code evaluators, AI/LLM services, and analytics sources participate through defined system boundaries. They are not interchangeable with human hiring decision-makers. |

### 2.4 Operating environment

- A browser-based user experience built with React.js and Next.js, with Arabic (`ar-EG`, right-to-left) and English (`en`) support.
- A Go backend, organized as a modular monolith, with separate API and background-worker execution roles as described in the SDD.
- OCI/Docker containers deployed to AWS or Azure; production data is planned for one region selected for Cairo latency and residency requirements.
- A relational database for transactional domain records and object storage for documents and assessment media, accessed through provider-neutral interfaces where required.
- A queue and workers for work that must be retried or performed asynchronously, including notification delivery and long-running evaluation tasks.
- External dependencies include email, official company registries, social login providers, a video-meeting provider, a code-evaluation provider, AI/LLM services, and ANA's configured telemetry or review sources.

The exact cloud provider, relational database product, queue product, and external vendors are deployment decisions; the SRS specifies their required behavior and boundaries, not an unapproved vendor choice.

### 2.5 Design and implementation constraints

- The product is scoped to Computer Science-related hiring and is delivered as SaaS (`CON-001`, `CON-002`).
- Employers configure pipelines and assessment rules; access to data and actions is role- and company-scoped (`CON-003`, `CON-009`, AUTH).
- Pipeline anonymity must hide identifying data from hiring staff until the configured reveal point (`CON-004`, `CON-010`, FR-INVW-070–074).
- Skills and crests are assessment-linked rather than awarded solely from Candidate claims (`CON-005`, FR-INVW-098–106).
- Candidate documents are restricted to PDF/DOCX and 10 MB per document (AUTH FR-AUTH-032; `CON-006`, `CON-007`).
- Candidate-to-Candidate communication is outside the stated project scope (`CON-008`).
- Posting lifecycle and application admission belong to JOB; post-admission application states and notifications belong to INVW. Company capacity is configured by JOB and consumed by INVW for finalist/offer decisions.
- Analytics may consume authorized source information but does not own or write another epic's records.

### 2.6 Assumptions and dependencies

The detailed register is maintained in `03-constraints.md`. Important system-level assumptions and dependencies include:

- Candidates and employers have an internet connection and a supported browser, although the product must tolerate the network conditions and device widths specified in `06-nfr.md`.
- Candidate and employer profile, posting, pipeline, and assessment data are supplied and maintained by their responsible users and epic owners.
- Employers configure pipeline stages, assessment content, criteria, and advancement rules before using them. A published posting may exist before its pipeline is published; admitted applications then wait for INVW handling.
- Account activation, authentication, RBAC, and company identity depend on AUTH. Posting creation, deadline/capacity data, application intake, and initial Must-have screening depend on JOB. Candidate profile/resume data depend on CAND. Post-admission status and notification behavior depend on INVW.
- Assessment and hiring workflows depend on external services remaining available. Failure handling is specified so an unavailable provider does not silently become a failed Candidate result or erase committed data.
- Analytics reports depend on the availability and validity of assessment, application, and pipeline data. ANA's contribution is analysis/evaluation; cross-epic source ownership remains with the originating epic.

### 2.7 System boundary and ownership summary

The Candidate-facing journey spans CAND, JOB, and INVW, with AUTH providing identity/access and ANA providing configured analysis. A high-level journey is: browse a published posting → submit an application → JOB screens and admits it → INVW runs the published pipeline → the Candidate sees status and feedback → the Hiring Manager records a final decision or offer → CAND presents the permitted outcome and profile artifacts.

The consolidated workflows are in `04-use-cases.md`. Functional source requirements are in `05-functional/{auth,cand,job,invw,ana}.md`; design decisions and component interactions are in `../sdd/`.
