## Appendix A: Glossary

Owner: TBD | Status: Draft | Last updated: 2026-10-14 | Jira: TBD

This glossary defines terms used across the SRS. Epic-specific behavior and authoritative status transitions remain in the linked functional requirements.

### A.1 People and organizations

| Term | Definition |
|---|---|
| Candidate | A person who uses AMARA to maintain a professional profile, apply for postings, complete hiring steps, view outcomes, and manage verified skills or a dynamic resume. Use **Candidate** as the current product term. |
| Company | An employer organization represented in AUTH and associated with its postings, members, applications, and hiring workflows. |
| Company Representative | An authenticated user acting for a Company. Company access and available actions are restricted by AUTH role and company scope. |
| Guest / Visitor | A person who has not authenticated. A Guest can use public or registration entry points; protected actions may require sign-in. |
| Hiring Manager | A Company-authorized user who manages postings and the interview workflow, configures pipelines, and makes the manual decisions assigned to that role. |
| Interviewer | A user associated with an approved Company and assigned to conduct a live-interview stage or submit its scorecard. |
| System Administrator (Admin) | A privileged platform operator who performs AUTH-defined account, company, taxonomy, and entity-administration actions. |
| Tenant / company scope | The Company context used to isolate company-owned records and authorize access. A valid role alone does not permit access to another Company's data. |

### A.2 Hiring and application workflow

| Term | Definition |
|---|---|
| Job posting / posting | A Company's published description of a role, its requirements, hiring configuration, open-position count, and application window. JOB owns posting content and lifecycle. |
| Posting states | JOB posting lifecycle values: **Draft**, **Published**, and **Archived**. Only JOB-defined transitions are valid. |
| Application | A Candidate's submission to one posting, including the immutable resume and answers captured at submission. JOB owns submission, initial screening, and pipeline admission; INVW owns progression after admission. |
| Application cap | The maximum number of accepted submissions that JOB will accept for a posting. It limits application intake. It is not the number of open positions. |
| Open-position capacity | The number of job positions configured on a posting by JOB. INVW reads this value when reviewing pending and accepted offers. |
| Finalist target | A planning target for the desired number of candidates still active in the hiring pipeline. It is not a limit on how many Candidates qualify as Finalists. |
| Screening | JOB's initial evaluation of an accepted application against the posting's structured Must-have criteria, before admission to INVW. |
| Screened Out | An application that does not satisfy one or more Must-have criteria and is not admitted to the interview pipeline. JOB records the failed criteria. |
| Admitted application | An accepted application that satisfies the Must-have criteria and is handed from JOB to INVW. If no pipeline is published, it may wait in **Awaiting Pipeline** state. |
| Interview pipeline | A posting-specific, directed acyclic graph (DAG) of assessment and interview stages configured by a Hiring Manager. |
| Pipeline version | A saved configuration of a posting's pipeline. INVW publishes one version at a time and defines migration behavior for in-progress applications. |
| Stage | A named node in a pipeline, with a module or live interview and one advancement rule. A stage may have incoming and outgoing links. |
| Entry stage | A pipeline stage with no incoming links. An admitted application begins at every entry stage. |
| Assessment module | Configured assessment content and scoring behavior associated with an individual assessment stage. |
| Test instance | A Candidate- and stage-assignment-specific execution of a published individual assessment module. Reopening returns the same instance and selected content. |
| Live interview | A scheduled interview stage conducted by assigned Interviewers or Hiring Managers, with a rubric, meeting, outcome, notes, and scorecard. It is distinct from an individual test instance. |
| Advancement rule | The stage policy that determines whether a Candidate passes, is rejected, or is held: manual decision, score threshold, or automatic on completion. |
| Held | A non-terminal application/stage condition awaiting a permitted Hiring Manager decision or required completion. It is not the same as a rejection. |
| Finalist | A Candidate whose application has passed every assessment stage on every path of the Candidate's published pipeline. INVW sets the application status to **Finalist** and adds it to the posting's Finalist list. |

### A.3 Assessment, skills, and analytics

| Term | Definition |
|---|---|
| Assessment score | A numeric result recorded for a fully scored assessment instance or live interview. INVW defines calculation and use for stage advancement. |
| Amara | A named, personified skill-crest definition associated with an assessment's scored evidence. Definitions may be built-in or Company-associated custom definitions. |
| Amarat | Plural of Amara; the skill crests awarded to a Candidate after the corresponding stage is passed. |
| Amara award | A recorded award of an Amara to a Candidate. It is associated with the source assessment stage and the Company that owns the posting. |
| Amara level / level band | A proficiency label and score range for an Amara, when level scoring is configured. Amara bands are distinct from stage advancement thresholds. |
| Crest | Candidate-facing term used in CAND for a displayed skill emblem/Amara award. A Crest is evidence-linked and is not merely a Candidate-entered skill claim. |
| Verified skill | A Candidate skill associated with an eligible assessment outcome and retained according to CAND/INVW rules. It is distinct from an unverified profile claim. |
| Scorecard | A live-interview evaluation submitted by an assigned Interviewer, containing rubric ratings and a recommendation. |
| Candidate insight | A labelled summary of strengths or weaknesses based on recorded assessment performance. It supports Hiring Manager review and does not independently change an application status. |
| Integrity Risk Score | An ANA-produced score and evidence indicating an integrity risk. Where configured, it flags a case for manual review; it is not itself a hiring decision. |
| Telemetry | Interaction or execution events received from a configured assessment environment for ANA analysis. Telemetry is personal/sensitive data and is subject to applicable access, purpose, and retention controls. |

### A.4 Offers, states, and system qualities

| Term | Definition |
|---|---|
| Offer | A structured proposal sent to a Finalist for a posting, containing the job title, compensation summary, proposed start date, response deadline, and optional message. Candidate response is recorded in the platform. |
| Offer record status | The state of an offer record: **Offer Pending**, **Offer Accepted**, **Offer Declined**, **Offer Expired**, or **Withdrawn**, as applicable. Offer-record state and application status are related but distinct; for example, withdrawing a pending offer returns the application to **Finalist**. |
| Anonymity | An INVW configuration that hides specified Candidate identity fields from Hiring Managers and Interviewers until the selected reveal point. It does not mean identity data is deleted. |
| Identity-reveal point | The configured pipeline stage or Finalist-list point after which authorized hiring users may see the Candidate's identity. The reveal is permanent. |
| Feedback mode | A per-stage or per-pipeline configuration for human-written, AI-generated, or AI-generated feedback with human approval, as defined by INVW. |
| Queue / worker | A queue stores asynchronous work for later execution; a worker consumes it, retries according to policy, and records job state. A worker is a runtime role of the backend, not a separate business epic. |
| Role-based access control (RBAC) | AUTH's permission mechanism for deciding which role may perform an operation. Resource access must also enforce Company scope when applicable. |
| RPO | Recovery Point Objective: the maximum targeted amount of data loss, expressed as a time interval. |
| RTO | Recovery Time Objective: the maximum targeted time to restore a specified service after an outage. |
| RPN | Risk Priority Number in the FMEA: Severity × Occurrence × Detectability using the scale stated in `08-risk-fmea.md`. |

### A.5 State ownership reminder

JOB is authoritative for posting states, application submission and initial screening/admission. INVW is authoritative for application states after admission, pipeline progression, finalist/offer transitions, and hiring-workflow notifications. CAND displays those states but does not create a competing state machine. AUTH owns identity and access-control state; ANA owns its analytical outputs and does not write another epic's domain records.
