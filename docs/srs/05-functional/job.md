## 5.4 JOB: Jobs, Search & Integrations

Owner: Yousef Abood | Status: In Review | Last updated: 2026-10-14 | Jira: AMARA-31

### 5.4.1 Overview

This section defines the functional requirements of the JOB epic. JOB covers the life of a job posting as a published offering (creating, configuring, publishing, archiving, republishing, and deleting it), the posting data used by candidate discovery, and how a candidate's application to it is submitted, screened against the Must-have criteria, and handed to the interview workflow after admission. Its actors are the **Hiring Manager**, which is the employer (the company, not a single person), and the **candidate**; Admins also act on postings. Posting management also covers previewing postings, autosaving drafts with warnings, handling banned companies, posting language and currency, and optionally switching the employer interface language. Candidate discovery, saved postings, followed companies, candidate-facing posting details, and candidate notifications are owned by CAND or AUTH and are referenced here only where they depend on posting data.

JOB ends at the admission decision. What happens to an admitted application (pipelines, assessments, interviews, feedback, offers, and candidate tracking of later stages) is outside this section and is not repeated here in INVW phase.

### 5.4.2 Scope

**In scope**
- Posting creation, editing, preview, language and currency, autosave and warnings, and optional employer interface language switching
- Publication, archiving, republishing, deletion, handling of banned companies, and the company posting list
- Posting tags
- Posting data supplied to CAND for the jobs feed, filtering, sorting, search, posting details, saved postings, followed companies, recommendations, saved searches, and sharing
- Hiring configuration of a posting: content, dates, capacity, criteria, and application questions
- Submitting an application
- Screening against Must-have criteria and admission to the pipeline

**Out of scope (owned elsewhere)**
- Processing of admitted applications: pipelines and their stages, test modules, live interviews, advancement, feedback, offers, and candidate tracking of later stages
- Company registration and verification, roles and permissions, company public profile page, job preferences, document upload, reporting, and admin management of users and taxonomies (AUTH)
- Candidate resume and profile (CAND)
- Candidate discovery, candidate-facing posting details, saved postings, followed companies, recommendations, saved searches, sharing, and candidate notifications (CAND; AUTH for the public company profile)
- Company-level analytics and cross-posting (ANA)

### 5.4.3 Features and Functional Requirements

#### 5.4.3.1 Job Posting Creation and Editing

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-001 | The system **SHALL** allow a Hiring Manager to create a job posting, which is saved in the Draft state. | UC-JOB-01 | TBD |
| FR-JOB-002 | The system **SHOULD** autosave a Draft posting while the Hiring Manager edits it, at least every 30 seconds while there are unsaved changes. | UC-JOB-01 | TBD |
| FR-JOB-003 | The system **SHOULD** warn the Hiring Manager before leaving the posting editor with unsaved changes, and **SHOULD** offer to save, discard, or stay. | UC-JOB-01 | TBD |
| FR-JOB-004 | The system **SHALL** not restrict the language in which a posting is written and **SHALL** display the posting text as written, without translating it. | UC-JOB-01 | TBD |
| FR-JOB-005 | The system **SHALL** not restrict the currency of a posting's salary range and **SHALL** display the amount in the chosen currency, without converting it. | UC-JOB-01 | TBD |
| FR-JOB-006 | The system **COULD** allow a Hiring Manager to duplicate a posting into a new Draft that copies its content and criteria but not its applications or dates. | UC-JOB-01 | TBD |
| FR-JOB-007 | The system **SHOULD** allow the Hiring Manager to preview a posting, including a Draft and unsaved edits, as a candidate will see it in the jobs feed and on the detail page. | UC-JOB-02 | TBD |
| FR-JOB-008 | The system **SHALL** allow only Hiring Manager accounts permitted to manage postings, and Admins, to create, edit, archive, republish, or delete a posting of a company. | UC-JOB-01, UC-JOB-05 | TBD |
| FR-JOB-009 | The system **COULD** allow the Hiring Manager to switch the language of the employer-side interface at any time among the supported interface languages (NFR-COMP-004; detailed locale rules are tracked in `docs/srs/TBD.md` (TBD-006)). | UC-JOB-08 | TBD |
| FR-JOB-010 | The system **COULD** keep the chosen language for the Hiring Manager's later sessions. | UC-JOB-08 | TBD |

#### 5.4.3.2 Publication and Lifecycle

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-011 | The system **SHALL** give every posting exactly one of the states Draft, Published, or Archived, and **SHALL** allow only these transitions: Draft to Published, Published to Archived, and Archived to Published. | UC-JOB-03 | TBD |
| FR-JOB-012 | The system **SHALL**, when a Hiring Manager publishes a posting, verify the mandatory content and dates and required language and currency fields, **SHALL** keep the posting in Draft if a check fails, and **SHALL** list every failed check. | UC-JOB-03 | TBD |
| FR-JOB-013 | The system **SHALL** make Published posting data available to CAND discovery and the AUTH-owned public company profile; Draft and Archived posting data is not exposed through those surfaces. | UC-JOB-03, UC-JOB-10 | TBD |
| FR-JOB-014 | The system **SHALL**, when a Hiring Manager archives a posting, remove it from the feed, search results, and company profile page, stop accepting new applications, and retain the posting, its applications, and its pipeline data. | UC-JOB-04 | TBD |
| FR-JOB-015 | The system **SHOULD** show the Hiring Manager the number of active applications and ask for confirmation before archiving a posting. | UC-JOB-04 | TBD |
| FR-JOB-016 | The system **SHALL** allow the Hiring Manager to republish an Archived posting, which returns it to Published and resumes its suspended applications, leaving all pass, reject, and offer decisions to the Hiring Manager. | UC-JOB-04 | TBD |
| FR-JOB-017 | The system **SHALL** accept new applications to a republished posting only if its application deadline has not passed; a posting republished after its deadline **SHALL** remain closed to new applications. | UC-JOB-03, UC-JOB-04 | TBD |
| FR-JOB-018 | The system **SHALL**, when a company is banned, archive every Published posting of the company, which suspends all their applications, and **SHALL** keep Draft postings as Draft. | UC-JOB-06 | TBD |
| FR-JOB-019 | The system **SHALL** not allow a banned company to publish or republish any posting while the ban lasts. | UC-JOB-06 | TBD |
| FR-JOB-020 | The system **SHALL**, when the ban is resolved, return each posting archived because of the ban to Published and resume its applications as in FR-JOB-016, and **SHALL** keep postings that were Archived before the ban as Archived. | UC-JOB-06 | TBD |
| FR-JOB-021 | The system **SHALL** allow a Hiring Manager permitted to manage postings to permanently delete any posting of their company, in any state and including postings that have applications. | UC-JOB-05 | TBD |
| FR-JOB-022 | The system **SHALL** ask for confirmation before a posting is deleted, **SHALL** state that deletion is permanent, and **SHALL** show the actor the number of applications affected when the posting has applications. | UC-JOB-05 | TBD |
| FR-JOB-023 | The system **SHOULD** show the Hiring Manager a list of all postings of the company in every state, with title, state, number of applications, and application deadline. | UC-JOB-07 | TBD |
| FR-JOB-024 | The system **COULD** allow the Hiring Manager to filter the posting list by state and to search it by title. | UC-JOB-07 | TBD |

#### 5.4.3.3 Posting Tags

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-025 | The system **SHOULD** allow the Hiring Manager to add and remove tags on a posting. | UC-JOB-01 | TBD |
| FR-JOB-026 | The system **COULD** allow at most 10 tags per posting. | UC-JOB-01 | TBD |

#### 5.4.3.4 Jobs Feed

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-027 | The system **SHALL** provide CAND with Published postings that are currently open for applications, newest first, in pages of 20. | UC-JOB-09 | TBD |
| FR-JOB-028 | The system **SHALL** exclude postings of companies that are deactivated or banned from the feed and search results. | UC-JOB-09 | TBD |

#### 5.4.3.5 Filtering and Sorting

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-029 | The system **COULD** apply a salary range filter only to postings that show their salary range and whose salary range is in the currency the candidate selects with the range, and **COULD** not convert amounts between currencies. | UC-JOB-09 | TBD |

#### 5.4.3.6 Search

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-030 | The system **SHOULD** match the keywords against the title, description, tags, required skills, and company name of postings, and against the names of companies. | UC-JOB-09 | TBD |

#### 5.4.3.12 Hiring Configuration

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-031 | The system **SHALL** require the Hiring Manager to specify a job title, a job description, and the number of open positions (an integer of at least 1) for every job posting. | UC-JOB-01, UC-JOB-03 | TBD |
| FR-JOB-032 | The system **SHALL** allow the Hiring Manager to specify a seniority level, employment type, work mode, location, and salary range for a posting, and whether the salary range is shown to candidates. A salary range that is not shown **SHALL** never be revealed to candidates, including through filters and search. | UC-JOB-01 | TBD |
| FR-JOB-033 | The system **SHALL** require the Hiring Manager to specify an application opening date and an application deadline, with the deadline later than the opening date. | UC-JOB-01, UC-JOB-03 | TBD |
| FR-JOB-034 | The system **SHOULD** allow the Hiring Manager to extend the application deadline of a posting at any time before the deadline passes, and **SHOULD** not allow the deadline to be extended once it has passed. | UC-JOB-01 | TBD |
| FR-JOB-035 | The system **SHALL** allow the Hiring Manager to set a maximum number of applications for a posting and **SHALL** stop accepting applications once that number has been accepted. | UC-JOB-01, UC-JOB-03 | TBD |
| FR-JOB-036 | The system **SHALL** require the Hiring Manager to classify each criterion of a posting as exactly one of Must-have or Nice-to-have. | UC-JOB-01 | TBD |
| FR-JOB-037 | The system **SHALL** require each Must-have criterion to be a checkable condition on a structured resume field (minimum years of experience, minimum education level, minimum language proficiency level, or location eligibility) and **SHALL** not accept free-text Must-have criteria. | UC-JOB-01 | TBD |
| FR-JOB-038 | The system **SHALL** allow the Hiring Manager to set a list of Nice-to-have criteria for a posting; Nice-to-have criteria **SHALL** not affect whether an application is admitted. | UC-JOB-01 | TBD |
| FR-JOB-039 | The system **SHALL** not allow the Must-have criteria of a posting to be added, changed, or removed once the posting has received its first application. | UC-JOB-01 | TBD |
| FR-JOB-040 | The system **SHOULD** allow the Hiring Manager to add application questions to a posting, each marked required or optional. | UC-JOB-01 | TBD |

#### 5.4.3.13 Submitting an Application

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-041 | The system **SHALL** require answers to all required application questions before accepting a submission, and **SHALL** list the questions left unanswered. | UC-JOB-09 | TBD |
| FR-JOB-042 | The system **SHALL** retain an unchangeable copy of the resume and answers as submitted, so that later resume edits do not alter the application. | UC-JOB-09 | TBD |
| FR-JOB-043 | The system **SHALL** reject a submission from a Candidate who already has an active application, or an application with status Rejected or Withdrawn, for the same posting; application status transitions are owned by INVW. | UC-JOB-09 | TBD |
| FR-JOB-044 | The system **SHALL** reject a submission made before the opening date, after the deadline, after the application cap is reached, or while the posting is Archived, **SHALL** tell the candidate which reason applies, and **SHALL** keep nothing partial. | UC-JOB-03, UC-JOB-09 | TBD |
| FR-JOB-045 | The system **SHOULD** show the candidate a confirmation with an application reference after a submission is accepted; notification delivery is owned by INVW. | UC-JOB-09 | TBD |

#### 5.4.3.14 Screening and Admission

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-046 | The system **SHALL** evaluate every accepted submission against each Must-have criterion of the posting, using the structured resume fields as submitted. | UC-JOB-09, UC-JOB-10 | TBD |
| FR-JOB-047 | The system **COULD** show the candidate, before the application is submitted, which Must-have criteria their resume does not currently satisfy. | UC-JOB-09 | TBD |
| FR-JOB-048 | The system **SHALL** admit to the pipeline an application that satisfies all Must-have criteria. | UC-JOB-09, UC-JOB-10 | TBD |
| FR-JOB-049 | The system **SHALL** screen out an application that fails any Must-have criterion, **SHALL** not admit it to the pipeline, and **SHALL** record each failed criterion. | UC-JOB-09, UC-JOB-10 | TBD |
| FR-JOB-050 | The system **SHALL** hand an admitted application to INVW, which starts the candidate at every entry stage when a pipeline is published. | UC-JOB-09, UC-JOB-10 | TBD |
| FR-JOB-051 | The system **SHOULD** allow the Hiring Manager to view all screened-out applications of a posting together with their failed criteria. | UC-JOB-10 | TBD |
| FR-JOB-052 | The system **COULD** tell a candidate whose application was screened out which Must-have criteria were not satisfied. | UC-JOB-09, UC-JOB-10 | TBD |

### 5.4.4 Business Rules

- Postings have Draft, Published, or Archived states. Only Published postings are available for candidate discovery.
- Archiving stops new applications and preserves existing application records.
- Candidates can apply only during the application window and before any application cap is reached.
- Must-have criteria determine screening outcomes; Nice-to-have criteria do not exclude applicants.
- A submitted application retains the resume and answers as submitted.
- CAND owns candidate discovery and saved jobs; AUTH owns company accounts; INVW owns downstream interview stages and notifications.

### 5.4.5 Error Handling

- Invalid or incomplete postings remain Draft and display validation errors.
- Invalid or duplicate applications are rejected with a reason and without partial submission.
- Unsaved posting changes trigger a warning before the Hiring Manager leaves the editor.
- Failed Must-have criteria are recorded for screening review.

