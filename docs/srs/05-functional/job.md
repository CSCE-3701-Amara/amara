## 5.4 JOB: Jobs, Search & Integrations

Owner: Yousef Abood | Status: In Review | Last updated: 2026-10-13 | Jira: AMARA-31

### 5.4.1 Overview

This section defines the functional requirements of the JOB epic. JOB covers the life of a job posting as a published offering (creating, configuring, publishing, archiving, republishing, and deleting it), the posting data used by candidate discovery, and how a candidate's application to it is submitted, screened against the Must-have criteria, and handed to the interview workflow after admission. Its actors are the **Hiring Manager**, which is the employer (the company, not a single person), and the **candidate**; Admins also act on postings. Posting management also covers previewing postings, autosaving drafts with warnings, handling banned companies, posting language and currency, and optionally switching the employer interface language. Candidate discovery, saved postings, followed companies, candidate-facing posting details, and candidate notifications are owned by CAND or AUTH and are referenced here only where they depend on posting data.

JOB ends at the admission decision. What happens to an admitted application (pipelines, assessments, interviews, feedback, offers, and candidate tracking of later stages) is outside this section and is not repeated here in INVW phase.

### 5.4.2 Scope

**In scope**
- Posting creation, editing, preview, language and currency, autosave and warnings, and optional employer interface language switching (5.4.3.1)
- Publication, archiving, republishing, deletion, handling of banned companies, and the company posting list (5.4.3.2)
- Posting tags (5.4.3.3)
- Posting data supplied to CAND for the jobs feed, filtering, sorting, search, posting details, saved postings, followed companies, recommendations, saved searches, and sharing (5.4.3.4 to 5.4.3.11)
- Hiring configuration of a posting: content, dates, capacity, criteria, and application questions (5.4.3.12)
- Submitting an application (5.4.3.13)
- Screening against Must-have criteria and admission to the pipeline (5.4.3.14)

**Out of scope (owned elsewhere)**
- Processing of admitted applications: pipelines and their stages, test modules, live interviews, advancement, feedback, offers, and candidate tracking of later stages
- Company registration and verification, roles and permissions, company public profile page, job preferences, document upload, reporting, and admin management of users and taxonomies (AUTH)
- Candidate resume and profile (CAND)
- Candidate discovery, candidate-facing posting details, saved postings, followed companies, recommendations, saved searches, sharing, and candidate notifications (CAND; AUTH for the public company profile)
- Company-level analytics and cross-posting (ANA)

### 5.4.3 Features and Functional Requirements

#### 5.4.3.1 Job Posting Creation and Editing

**Feature Description:** A Hiring Manager creates a job posting as a draft, completes its content and hiring configuration, and edits it over time. The posting states its language and currency, can be previewed as candidates will see it, and the employer interface language can optionally be switched. Drafts are saved automatically, and the Hiring Manager is warned before leaving with unsaved changes.

**Preconditions:**
* The company is registered and verified.
* The Hiring Manager is signed in with a company account permitted to manage postings.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-001 | The system **SHALL** allow a Hiring Manager to create a job posting, which is saved in the Draft state. | UC-JOB-01 | TBD |
| FR-JOB-002 | The system **SHALL** allow a Draft posting to be saved with incomplete fields. | UC-JOB-01 | TBD |
| FR-JOB-003 | The system **SHALL** autosave a Draft posting while the Hiring Manager edits it, at least every 30 seconds while there are unsaved changes. | UC-JOB-01 | TBD |
| FR-JOB-004 | The system **SHALL** show the Hiring Manager whether the latest changes to a Draft are saved, with the time of the last save. | UC-JOB-01 | TBD |
| FR-JOB-005 | The system **SHALL** apply autosave to Draft postings only; changes to a Published posting take effect only when the Hiring Manager saves them. | UC-JOB-01 | TBD |
| FR-JOB-006 | The system **SHALL** warn the Hiring Manager before leaving the posting editor with unsaved changes, and **SHALL** offer to save, discard, or stay. | UC-JOB-01 | TBD |
| FR-JOB-007 | The system **SHALL** warn the Hiring Manager before leaving a new posting that has not yet been saved as a Draft, and **SHALL** offer to save it as a Draft, discard it, or stay. | UC-JOB-01 | TBD |
| FR-JOB-008 | The system **SHALL** allow the Hiring Manager to enter the content and hiring configuration defined in FR-JOB-089 to FR-JOB-098 as part of creating the posting. | UC-JOB-01 | TBD |
| FR-JOB-009 | The system **SHALL** not restrict the language in which a posting is written and **SHALL** display the posting text as written, without translating it. | UC-JOB-01 | TBD |
| FR-JOB-010 | The system **SHALL** not restrict the currency of a posting's salary range and **SHALL** display the amount in the chosen currency, without converting it. | UC-JOB-01 | TBD |
| FR-JOB-011 | The system **SHALL** require the Hiring Manager to specify the language of a posting, chosen from a standard list of languages, before the posting is published (FR-JOB-024). | UC-JOB-01, UC-JOB-03 | TBD |
| FR-JOB-012 | The system **SHALL** require the Hiring Manager to specify the currency of the salary range, chosen from a standard list of currencies, before a posting with a salary range is published (FR-JOB-024). | UC-JOB-01, UC-JOB-03 | TBD |
| FR-JOB-013 | The system **SHALL** apply the deadline rule (FR-JOB-092) and the Must-have lock (FR-JOB-097) when a Hiring Manager edits a Published posting. | UC-JOB-01 | TBD |
| FR-JOB-014 | Candidate notification behavior for posting changes is owned by INVW. | UC-JOB-01 | TBD |
| FR-JOB-015 | The system **COULD** allow a Hiring Manager to duplicate a posting into a new Draft that copies its content and criteria but not its applications or dates. | UC-JOB-01 | TBD |
| FR-JOB-016 | The system **SHALL** allow the Hiring Manager to preview a posting, including a Draft and unsaved edits, as a candidate will see it in the jobs feed (FR-JOB-049) and on the detail page (FR-JOB-068). | UC-JOB-02 | TBD |
| FR-JOB-017 | The system **SHALL** not change any posting, application, or notification data, and **SHALL** not notify anyone, when a posting is previewed. | UC-JOB-02 | TBD |
| FR-JOB-018 | The system **SHALL** allow only Hiring Manager accounts permitted to manage postings, and Admins, to create, edit, archive, republish, or delete a posting of a company. | UC-JOB-01, UC-JOB-05 | TBD |
| FR-JOB-019 | The system **COULD** allow the Hiring Manager to switch the language of the employer-side interface at any time, among the interface languages defined in `10-i18n.md`. | UC-JOB-08 | TBD |
| FR-JOB-020 | When employer interface language switching is available, the system **SHALL** apply the chosen language immediately, without requiring the Hiring Manager to sign in again and without losing unsaved changes. | UC-JOB-08 | TBD |
| FR-JOB-021 | The system **SHOULD** keep the chosen language for the Hiring Manager's later sessions. | UC-JOB-08 | TBD |
| FR-JOB-022 | When employer interface language switching is available, the system **SHALL** not change the content of any posting when the interface language is switched. | UC-JOB-08 | TBD |

#### 5.4.3.2 Publication and Lifecycle

**Feature Description:** A posting moves between Draft, Published, and Archived. A Published posting is listed in the jobs feed and in search. Archiving suspends the posting's applications without deleting them, and republishing resumes them. When a company is banned, its Published postings are archived until the ban is resolved. Deleting a posting needs confirmation.

**Preconditions:**
* The posting exists and the Hiring Manager is permitted to manage it (FR-JOB-018).

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-023 | The system **SHALL** give every posting exactly one of the states Draft, Published, or Archived, and **SHALL** allow only these transitions: Draft to Published, Published to Archived, and Archived to Published. | UC-JOB-03 | TBD |
| FR-JOB-024 | The system **SHALL**, when a Hiring Manager publishes a posting, verify the mandatory content and dates (FR-JOB-089, FR-JOB-091) and the language and currency required by FR-JOB-011 and FR-JOB-012, **SHALL** keep the posting in Draft if a check fails, and **SHALL** list every failed check. | UC-JOB-03 | TBD |
| FR-JOB-025 | The system **SHALL** allow a posting to be published before a pipeline is published, in which case admitted applications are handed to INVW and await its published-pipeline handling (FR-JOB-109). | UC-JOB-03, UC-JOB-17 | TBD |
| FR-JOB-026 | The system **SHALL** make Published posting data available to CAND discovery and the AUTH-owned public company profile; Draft and Archived posting data is not exposed through those surfaces. | UC-JOB-03, UC-JOB-10 | TBD |
| FR-JOB-027 | The system **SHALL**, when a Hiring Manager archives a posting, remove it from the feed, search results, and company profile page, stop accepting new applications, and retain the posting, its applications, and its pipeline data. | UC-JOB-04 | TBD |
| FR-JOB-028 | The system **SHALL**, while a posting is Archived, suspend participation in its downstream interview workflow: candidates cannot start or continue modules or book live interviews, and the posting and submitted application data are retained. Application statuses and transitions remain owned by INVW. | UC-JOB-04 | TBD |
| FR-JOB-029 | The system **SHALL**, while a posting is Archived, pause the completion windows of its modules, **SHALL** not mark any module as missed, and **SHALL** extend each window by the time the posting was archived when it is republished. | UC-JOB-04 | TBD |
| FR-JOB-030 | The system **SHALL** show a candidate with an application to an Archived posting, in the application tracking view, that the posting is archived and cannot be completed. | UC-JOB-04 | TBD |
| FR-JOB-031 | Application status transitions, including withdrawal, are owned by INVW. | UC-JOB-04 | TBD |
| FR-JOB-032 | The system **SHOULD** show the Hiring Manager the number of active applications and ask for confirmation before archiving a posting. | UC-JOB-04 | TBD |
| FR-JOB-033 | Candidate notification behavior for archiving and republishing is owned by INVW. | UC-JOB-04 | TBD |
| FR-JOB-034 | The system **SHALL** allow the Hiring Manager to republish an Archived posting, which returns it to Published and resumes its suspended applications, leaving all pass, reject, and offer decisions to the Hiring Manager. | UC-JOB-04 | TBD |
| FR-JOB-035 | The system **SHALL** accept new applications to a republished posting only if its application deadline has not passed; a posting republished after its deadline **SHALL** remain closed to new applications. | UC-JOB-03, UC-JOB-04 | TBD |
| FR-JOB-036 | The system **SHALL**, when a company is banned, archive every Published posting of the company, which suspends all their applications (FR-JOB-028, FR-JOB-029), and **SHALL** keep Draft postings as Draft. | UC-JOB-06 | TBD |
| FR-JOB-037 | The system **SHALL** not allow a banned company to publish or republish any posting while the ban lasts. | UC-JOB-06 | TBD |
| FR-JOB-038 | The system **SHALL**, when the ban is resolved, return each posting archived because of the ban to Published and resume its applications as in FR-JOB-034, and **SHALL** keep postings that were Archived before the ban as Archived. | UC-JOB-06 | TBD |
| FR-JOB-039 | Candidate notification behavior for company bans and resumptions is owned by INVW. | UC-JOB-04, UC-JOB-06 | TBD |
| FR-JOB-040 | The system **SHALL** allow a Hiring Manager permitted to manage postings to permanently delete any posting of their company, in any state and including postings that have applications. | UC-JOB-05 | TBD |
| FR-JOB-041 | The system **SHALL** allow an Admin to permanently delete any posting at any time, including a posting that has applications. | UC-JOB-05 | TBD |
| FR-JOB-042 | The system **SHALL** ask for confirmation before a posting is deleted, **SHALL** state that deletion is permanent, and **SHALL** show the actor the number of applications affected when the posting has applications. | UC-JOB-05 | TBD |
| FR-JOB-043 | The system **SHALL** show the Hiring Manager a list of all postings of the company in every state, with title, state, number of applications, and application deadline. | UC-JOB-07 | TBD |
| FR-JOB-044 | The system **SHOULD** allow the Hiring Manager to filter the posting list by state and to search it by title. | UC-JOB-07 | TBD |

#### 5.4.3.3 Posting Tags

**Feature Description:** Tags classify a posting and feed the filters, search, and recommendations.

**Preconditions:**
* A Draft or Published posting exists.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-045 | The system **SHALL** allow the Hiring Manager to add and remove tags on a posting. | UC-JOB-01 | TBD |
| FR-JOB-046 | The system **SHALL** allow at most 10 tags per posting. | UC-JOB-01 | TBD |

#### 5.4.3.4 Jobs Feed

**Feature Description:** JOB supplies posting data to the CAND-owned jobs feed, which applies candidate visibility and ordering rules.

**Preconditions:**
* The user is signed in (FR-JOB-047).

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-047 | The system **SHALL** require a user to be signed in to view the jobs feed, search results, posting details, or shared posting links, **SHALL** send a visitor to sign-in or registration, and **SHALL** open the requested page afterwards. Access to public company profiles follows AUTH. | UC-JOB-09, UC-JOB-10, UC-JOB-11 | TBD |
| FR-JOB-048 | The system **SHALL** provide CAND with Published postings that are currently open for applications, newest first, in pages of 20. | UC-JOB-09 | TBD |
| FR-JOB-049 | Candidate-facing feed fields are owned by CAND; JOB owns the underlying posting data and salary-visibility rule. | UC-JOB-09 | TBD |
| FR-JOB-050 | Saved and applied markers in candidate discovery are owned by CAND. | UC-JOB-09 | TBD |
| FR-JOB-051 | Candidate recommendation behavior is owned by CAND. | UC-JOB-14 | TBD |
| FR-JOB-052 | Candidate recommendation explanations are owned by CAND. | UC-JOB-14 | TBD |
| FR-JOB-053 | Candidate followed-company filtering is owned by CAND. | UC-JOB-09 | TBD |
| FR-JOB-054 | The system **SHALL** exclude postings of companies that are deactivated or banned from the feed and search results. | UC-JOB-09 | TBD |

#### 5.4.3.5 Filtering and Sorting

**Feature Description:** JOB supplies structured posting fields to the CAND-owned filtering and sorting behavior.

**Preconditions:**
* The user is signed in and the feed or a search result is displayed.

| ID | Requirement | Related Use Case | Related Functional Requirement | Jira |
|---|---|---|---|---|
| FR-JOB-055 | Candidate filtering behavior is owned by CAND; JOB owns the structured posting fields used by those filters. | UC-JOB-09 | FR-CAND-004 to FR-CAND-009 | TBD |
| FR-JOB-056 | Candidate filter matching behavior is owned by CAND. | UC-JOB-09 | — | TBD |
| FR-JOB-057 | The system **SHALL** apply a salary range filter (FR-JOB-055) only to postings that show their salary range (FR-JOB-090) and whose salary range is in the currency the candidate selects with the range, and **SHALL** not convert amounts between currencies. | UC-JOB-09 | — | TBD |
| FR-JOB-058 | Candidate sorting behavior is owned by CAND. | UC-JOB-09 | — | TBD |
| FR-JOB-059 | Candidate filter presentation and clearing behavior are owned by CAND. | UC-JOB-09 | — | TBD |
| FR-JOB-060 | Candidate no-results behavior is owned by CAND. | UC-JOB-09 | — | TBD |
| FR-JOB-061 | Candidate filtering and sorting behavior is owned by CAND. | UC-JOB-09 | — | TBD |

#### 5.4.3.6 Search

**Feature Description:** JOB supplies searchable posting fields to the CAND-owned keyword search; AUTH owns public company profile data.

**Preconditions:**
* The user is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-062 | Candidate search entry points are owned by CAND. | UC-JOB-09 | TBD |
| FR-JOB-063 | The system **SHALL** match the keywords against the title, description, tags, required skills, and company name of postings, and against the names of companies. | UC-JOB-09 | TBD |
| FR-JOB-064 | Candidate search result grouping and ranking are owned by CAND. | UC-JOB-09 | TBD |
| FR-JOB-065 | Candidate search validation is owned by CAND. | UC-JOB-09 | TBD |
| FR-JOB-066 | The system **SHALL** make only Published posting data available to CAND search and AUTH public company profiles, and **SHALL** never expose Draft or Archived postings outside the owning company's management views. | UC-JOB-09 | TBD |
| FR-JOB-067 | Candidate search matching behavior is owned by CAND. | UC-JOB-09 | TBD |

#### 5.4.3.7 Posting Details

**Feature Description:** JOB supplies posting data and lifecycle information to the CAND-owned detail view.

**Preconditions:**
* The user is signed in.
* The posting is Published, or the user has applied to or saved it.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-068 | Candidate-facing posting details are owned by CAND; JOB owns the underlying posting data and salary-visibility rule. | UC-JOB-10 | TBD |
| FR-JOB-069 | Candidate-facing posting-state presentation is owned by CAND; JOB owns the lifecycle state and application-window rules. | UC-JOB-03, UC-JOB-10 | TBD |
| FR-JOB-070 | Candidate apply-action presentation is owned by CAND; JOB owns application validation and screening. | UC-JOB-10 | TBD |
| FR-JOB-071 | Candidate save and follow actions are owned by CAND. | UC-JOB-10 | TBD |
| FR-JOB-072 | Candidate access to archived postings is owned by CAND; JOB owns archival lifecycle rules. | UC-JOB-10 | TBD |

#### 5.4.3.8 Saved Jobs

**Feature Description:** JOB supplies posting identity and lifecycle data to the CAND-owned saved-posting list.

**Preconditions:**
* The candidate is signed in.

| ID | Requirement | Related Use Case | Related Functional Requirement | Jira |
|---|---|---|---|---|
| FR-JOB-073 | Candidate saved-posting actions are owned by CAND. | UC-JOB-12 | FR-CAND-044 | TBD |
| FR-JOB-074 | Candidate saved-list presentation is owned by CAND. | UC-JOB-12 | — | TBD |
| FR-JOB-075 | Candidate handling of unavailable saved postings is owned by CAND. | UC-JOB-12 | — | TBD |
| FR-JOB-076 | Saved-posting reminders are owned by CAND or INVW notification processing. | UC-JOB-12 | — | TBD |

#### 5.4.3.9 Following Companies

**Feature Description:** JOB supplies posting data to the CAND-owned followed-company behavior.

**Preconditions:**
* The candidate is signed in.

| ID | Requirement | Related Use Case | Related Functional Requirement | Jira |
|---|---|---|---|---|
| FR-JOB-077 | Candidate follow and unfollow actions are owned by CAND; public company profiles are owned by AUTH. | UC-JOB-13 | FR-CAND-043 | TBD |
| FR-JOB-078 | Candidate followed-company list presentation is owned by CAND. | UC-JOB-13 | — | TBD |
| FR-JOB-079 | Candidate notification behavior for followed companies is owned by INVW. | UC-JOB-13 | — | TBD |
| FR-JOB-080 | Candidate notification preferences are owned by INVW. | UC-JOB-13 | — | TBD |

#### 5.4.3.10 Recommendations, Saved Searches, and New-Jobs Tracker

**Feature Description:** JOB supplies structured posting data to CAND recommendations and to INVW notification processing.

**Preconditions:**
* The candidate is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-081 | Candidate similar-posting recommendations are owned by CAND. | UC-JOB-14 | TBD |
| FR-JOB-082 | Candidate post-application recommendations are owned by CAND. | UC-JOB-14 | TBD |
| FR-JOB-083 | Candidate recommendation exclusions are owned by CAND. | UC-JOB-14 | TBD |
| FR-JOB-084 | Candidate recommendation inputs and fairness constraints are owned by CAND. | UC-JOB-14 | TBD |
| FR-JOB-085 | Candidate saved-search actions are owned by CAND. | UC-JOB-15 | TBD |
| FR-JOB-086 | Candidate notification behavior for saved searches is owned by INVW. | UC-JOB-15 | TBD |
| FR-JOB-087 | Candidate saved-search notification preferences are owned by INVW. | UC-JOB-15 | TBD |

#### 5.4.3.11 Sharing

**Feature Description:** JOB supplies the posting identity and data addressed by the CAND-owned sharing link.

**Preconditions:**
* The posting is Published and the user is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-088 | Candidate posting-link sharing behavior is owned by CAND; JOB owns the posting data addressed by the link. | UC-JOB-11 | TBD |

#### 5.4.3.12 Hiring Configuration

**Feature Description:** The Hiring Manager defines what the posting offers, when and how many applications it accepts, which conditions a candidate must meet (Must-have), which qualities are preferred (Nice-to-have), and which questions candidates answer when applying.

**Preconditions:**
* A Draft or Published posting exists and the Hiring Manager is permitted to manage it (FR-JOB-018).

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-089 | The system **SHALL** require the Hiring Manager to specify a job title, a job description, and the number of open positions (an integer of at least 1) for every job posting. | UC-JOB-01, UC-JOB-03 | TBD |
| FR-JOB-090 | The system **SHALL** allow the Hiring Manager to specify a seniority level, employment type, work mode, location, and salary range for a posting, and whether the salary range is shown to candidates. A salary range that is not shown **SHALL** never be revealed to candidates, including through filters and search (FR-JOB-057). | UC-JOB-01 | TBD |
| FR-JOB-091 | The system **SHALL** require the Hiring Manager to specify an application opening date and an application deadline, with the deadline later than the opening date. | UC-JOB-01, UC-JOB-03 | TBD |
| FR-JOB-092 | The system **SHALL** allow the Hiring Manager to extend the application deadline of a posting at any time before the deadline passes, and **SHALL** not allow the deadline to be extended once it has passed. | UC-JOB-01 | TBD |
| FR-JOB-093 | The system **SHALL** allow the Hiring Manager to set a maximum number of applications for a posting and **SHALL** stop accepting applications once that number has been accepted. | UC-JOB-01, UC-JOB-03 | TBD |
| FR-JOB-094 | The system **SHALL** require the Hiring Manager to classify each criterion of a posting as exactly one of Must-have or Nice-to-have. | UC-JOB-01 | TBD |
| FR-JOB-095 | The system **SHALL** require each Must-have criterion to be a checkable condition on a structured resume field (minimum years of experience, minimum education level, minimum language proficiency level, or location eligibility) and **SHALL** not accept free-text Must-have criteria. | UC-JOB-01 | TBD |
| FR-JOB-096 | The system **SHALL** allow the Hiring Manager to set a list of Nice-to-have criteria for a posting; Nice-to-have criteria **SHALL** not affect whether an application is admitted. | UC-JOB-01 | TBD |
| FR-JOB-097 | The system **SHALL** not allow the Must-have criteria of a posting to be added, changed, or removed once the posting has received its first application. | UC-JOB-01 | TBD |
| FR-JOB-098 | The system **SHALL** allow the Hiring Manager to add application questions to a posting, each marked required or optional. | UC-JOB-01 | TBD |

#### 5.4.3.13 Submitting an Application

**Feature Description:** A candidate applies to an open posting with their resume and answers to the application questions. The submission is validated, stored unchanged, and confirmed.

**Preconditions:**
* The candidate is signed in and has a complete resume.
* The posting is Published.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-099 | The system **SHALL** require answers to all required application questions before accepting a submission, and **SHALL** list the questions left unanswered. | UC-JOB-16 | TBD |
| FR-JOB-100 | The system **SHALL** retain an unchangeable copy of the resume and answers as submitted, so that later resume edits do not alter the application. | UC-JOB-16 | TBD |
| FR-JOB-101 | The system **SHALL** reject a submission from a candidate who already has an active application for the same posting; application status transitions are owned by INVW. | UC-JOB-16 | TBD |
| FR-JOB-102 | The system **SHALL** reject a submission made before the opening date, after the deadline, after the application cap is reached, or while the posting is Archived, **SHALL** tell the candidate which reason applies, and **SHALL** keep nothing partial. | UC-JOB-03, UC-JOB-16 | TBD |
| FR-JOB-103 | The system **SHALL** count a submission toward the application cap (FR-JOB-093) only once it has been accepted. | UC-JOB-03, UC-JOB-16 | TBD |
| FR-JOB-104 | The system **SHALL** show the candidate a confirmation with an application reference after a submission is accepted; notification delivery is owned by INVW. | UC-JOB-16 | TBD |

#### 5.4.3.14 Screening and Admission

**Feature Description:** Each accepted submission is checked automatically against the posting's Must-have criteria. Applications that satisfy all of them are admitted to the pipeline; the others are screened out and remain visible to the Hiring Manager with the reasons.

**Preconditions:**
* An application submission has been accepted (FR-JOB-104).

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-105 | The system **SHALL** evaluate every accepted submission against each Must-have criterion of the posting, using the structured resume fields as submitted (FR-JOB-100). | UC-JOB-16, UC-JOB-17 | TBD |
| FR-JOB-106 | The system **SHOULD** show the candidate, before the application is submitted, which Must-have criteria their resume does not currently satisfy. | UC-JOB-16 | TBD |
| FR-JOB-107 | The system **SHALL** admit to the pipeline an application that satisfies all Must-have criteria. | UC-JOB-16, UC-JOB-17 | TBD |
| FR-JOB-108 | The system **SHALL** screen out an application that fails any Must-have criterion, **SHALL** not admit it to the pipeline, and **SHALL** record each failed criterion. | UC-JOB-16, UC-JOB-17 | TBD |
| FR-JOB-109 | The system **SHALL** hand an admitted application to INVW, which starts the candidate at every entry stage when a pipeline is published. | UC-JOB-16, UC-JOB-17 | TBD |
| FR-JOB-110 | The system **SHALL** allow the Hiring Manager to view all screened-out applications of a posting together with their failed criteria. | UC-JOB-17 | TBD |
| FR-JOB-111 | The system **SHOULD** tell a candidate whose application was screened out which Must-have criteria were not satisfied. | UC-JOB-16, UC-JOB-17 | TBD |

### 5.4.4 Business Rules

* A posting is Draft, Published, or Archived, and only the transitions in FR-JOB-023 exist.
* A posting can be published before a pipeline is published; admitted applications are handed to INVW for handling until a pipeline is published (FR-JOB-025, FR-JOB-109).
* Hiring configuration of a posting is defined once, in 5.4.3.12.
* Candidate access to the feed, search, posting details, and shared posting links follows CAND; access to public company profiles follows AUTH (FR-JOB-047).
* Only Published postings appear in the feed, search, and the company profile page (FR-JOB-026, FR-JOB-066).
* A Hiring Manager permitted to manage postings and an Admin can permanently delete any posting, in any state and including postings with applications (FR-JOB-040, FR-JOB-041).
* Archiving suspends downstream interview-workflow participation but never deletes application data; application statuses and transitions follow INVW (FR-JOB-028, FR-JOB-034).
* A posting can be reopened to new applications only before its deadline passes (FR-JOB-035, FR-JOB-092).
* A salary range hidden from candidates is never revealed through filters or search (FR-JOB-057, FR-JOB-090).
* Saved postings, saved searches, and followed companies belong to the candidate. The Hiring Manager and Admins never see which candidates saved or followed a posting or company.
* Recommendations never use protected or identifying attributes (FR-JOB-084).
* Tags come only from the taxonomy managed by Admins (FR-JOB-045).
* Candidate notification behavior and delivery follow INVW; JOB supplies the relevant posting and application events.
* Previewing a posting never changes data or notifies anyone (FR-JOB-017).
* While a company is banned, its postings stay archived and nothing can be published. Applications resume when the ban is resolved, and the Hiring Manager then decides on them (FR-JOB-036 to FR-JOB-038).
* Every posting states its language and, when it shows a salary range, its currency. Postings are shown as written, without translation or conversion, and salary filters never convert currencies (FR-JOB-009 to FR-JOB-012, FR-JOB-057).
* Switching the employer interface language never changes posting content (FR-JOB-022).
* Drafts are saved automatically. Changes to a Published posting are saved only on request, and the Hiring Manager is warned before leaving with unsaved changes, before leaving an unsaved new posting, and before deleting a posting (FR-JOB-003 to FR-JOB-007, FR-JOB-042).
* Every criterion is exactly one of Must-have or Nice-to-have. Must-have criteria are checkable conditions on structured resume fields, never free text, and are locked after the first application (FR-JOB-094, FR-JOB-095, FR-JOB-097).
* Only Must-have criteria decide admission. Nice-to-have criteria never exclude a candidate (FR-JOB-096, FR-JOB-107, FR-JOB-108).
* A candidate cannot submit another application while an active application exists for the same posting (FR-JOB-101).
* The submitted resume and answers never change after submission (FR-JOB-100).
* A deadline can be extended only before it passes, and the application cap closes the posting to new applications once reached (FR-JOB-092, FR-JOB-093, FR-JOB-102).

### 5.4.5 Error Handling

* **Publishing with missing or invalid data:** keep the posting in Draft and list every failed check (FR-JOB-024).
* **Two accounts of the same Hiring Manager editing the same posting:** save the first change, reject the second, and show the second editor the current posting data.
* **Link to a deleted or unavailable posting:** show "This posting is not available" without revealing whether the posting exists.
* **Visitor who is not signed in opens a link:** send the visitor to sign-in or registration and open the requested page afterwards (FR-JOB-047).
* **Submission to a posting archived, closed, not yet open, or at its application cap, including one that changes state while the candidate is filling the form:** reject with the specific reason and keep nothing partial (FR-JOB-102).
* **Submission with unanswered required questions:** reject and list the missing answers (FR-JOB-099).
* **Submission by a candidate who already has an active application to the posting:** reject and state that an application already exists (FR-JOB-101).
* **Two submissions racing for the last place under the application cap:** accept the first, reject the second with the cap reason (FR-JOB-093, FR-JOB-102).
* **Hiring Manager changes Must-have criteria after the first application:** refuse and state that the criteria are locked (FR-JOB-097).
* **Hiring Manager extends a deadline that has already passed, or sets a deadline not later than the opening date:** refuse and identify the invalid date (FR-JOB-091, FR-JOB-092).
* **Candidate action on an Archived posting (start module, book interview):** refuse and show that the posting is archived (FR-JOB-028).
* **Salary filter with minimum above maximum, or too-short search query:** show a message identifying the invalid input and keep the previous results.
* **Search service unavailable:** show "Search is temporarily unavailable" and keep the feed and filters working.
* **Notification or confirmation delivery fails:** apply INVW notification retry and in-platform delivery rules.
* **Tag removed from the taxonomy by an Admin:** existing postings keep the tag, and it can no longer be added to new ones.
* **Publishing without a language, or with a salary range but no currency:** keep the posting in Draft and list it as a failed check (FR-JOB-024, FR-JOB-011, FR-JOB-012).
* **Hiring Manager edits or publishes while the company is banned:** refuse and show that the account is banned (FR-JOB-037).
* **Autosave fails (network or service error):** tell the Hiring Manager that the latest changes are not saved, keep the entered data on screen, and retry (FR-JOB-003).
* **Screening cannot be completed (a structured resume field needed by a Must-have criterion is empty):** treat the criterion as not satisfied and record it as a failed criterion (FR-JOB-108).