Owner: Yousef Abood | Status: Draft | Last updated: 2026-10-10 | Jira: AMARA-31

## 5.4 JOB: Jobs, Search & Integrations

### 5.4.1 Overview

This section defines the functional requirements of the JOB epic. JOB covers the life of a job posting as a published offering (creating, publishing, controlling its visibility, archiving, republishing, and deleting it) and how Interviewees discover it (jobs feed, filtering, search, posting details, saved jobs, following companies, recommendations, and sharing). Its actors are the **Hiring Manager** and the **Interviewee**; Admins act on postings through AUTH. Every user must be signed in to browse (FR-JOB-034).

JOB does not define what happens once an Interviewee applies. Hiring configuration of a posting (content fields, deadlines, capacity, criteria, finalist target), application status, pipelines, assessments, interviews, anonymity, feedback, offers, and notifications belong to APP. This section refers to APP requirements by ID and never repeats them. Every notification mentioned here is sent through the notification mechanism owned by APP.

Use cases are not written yet. Every "Related Use Case" value is TBD until `04-use-cases.md` exists. Priority is expressed with **SHALL** (Must), **SHOULD** (Should), and **COULD** (Could).

Diagrams (TBD): Figure 1: JOB use case diagram. Figure 2: job posting lifecycle state diagram (Draft, Published, Archived).

### 5.4.2 Scope

**In scope**
- Posting creation, editing, versioning, and audit (5.4.3.1)
- Publication, visibility, private access, archiving, republishing, deletion, and the company posting list (5.4.3.2)
- Posting tags (5.4.3.3)
- Jobs feed, filtering, sorting, and search (5.4.3.4 to 5.4.3.6)
- Posting details, saved jobs, following companies (5.4.3.7 to 5.4.3.9)
- Recommendations, saved searches, new-jobs tracker, and sharing (5.4.3.10, 5.4.3.11)

**Out of scope (owned elsewhere)**
- Hiring configuration of a posting: content fields, dates, capacity, Must-have and Nice-to-have criteria, application questions, finalist target (APP 3.1, FR-APP-01 to FR-APP-12)
- Applying, withdrawing, screening, application statuses, pipelines, test modules (including sharing them between companies, as an extension of FR-APP-49), live interviews, advancement, anonymity, feedback, offers, Interviewee application tracking, and all notifications (APP)
- Company registration, verification, roles and permissions, company public profile page, job preferences, document upload, reporting, and admin management of users and taxonomies (AUTH)
- Interviewee resume and profile (CAND)
- Company-level analytics and cross-posting (ANA)

### 5.4.3 Features and Functional Requirements

#### 5.4.3.1 Job Posting Creation and Editing

**Feature Description:** A Hiring Manager creates a job posting as a draft, completes its content and hiring configuration, and edits it over time with a record of every change.

**Preconditions:**
* The company is registered and verified (AUTH).
* The Hiring Manager is signed in and permitted to manage postings of the company, per the AUTH permission matrix.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-001 | The system **SHALL** allow a Hiring Manager to create a job posting, which is saved in the Draft state. | TBD | TBD |
| FR-JOB-002 | The system **SHALL** allow a Draft posting to be saved with incomplete fields. | TBD | TBD |
| FR-JOB-003 | The system **SHALL** allow the Hiring Manager to enter the content and hiring configuration defined in FR-APP-01 to FR-APP-12 as part of creating the posting. | TBD | TBD |
| FR-JOB-004 | The system **SHALL** apply FR-APP-04 (deadline extension only before the deadline passes) and FR-APP-10 (Must-have criteria locked after the first application) when a Hiring Manager edits a Published posting. | TBD | TBD |
| FR-JOB-005 | The system **SHALL** record in an audit log every creation, edit, and state change of a posting, with the author, the time, and the changed fields. | TBD | TBD |
| FR-JOB-006 | The system **SHOULD** keep each previous version of a Published posting and allow the Hiring Manager to view it. | TBD | TBD |
| FR-JOB-007 | The system **SHOULD** notify Interviewees with an active application when the title, salary range, location, or work mode of the posting changes. | TBD | TBD |
| FR-JOB-008 | The system **COULD** allow a Hiring Manager to duplicate a posting into a new Draft that copies its content and criteria but not its applications or dates. | TBD | TBD |
| FR-JOB-009 | The system **SHALL** allow only Hiring Managers permitted by the AUTH permission matrix, and Admins (FR-AUTH-036), to create, edit, archive, republish, or delete a posting of a company. | TBD | TBD |

#### 5.4.3.2 Publication, Visibility, and Lifecycle

**Feature Description:** A posting moves between Draft, Published, and Archived. A Published posting is either public (listed in the feed and search) or private (reachable only through a private link or an email invitation). Archiving suspends the posting's applications without deleting them, and republishing resumes them.

**Preconditions:**
* The posting exists and the Hiring Manager is permitted to manage it (FR-JOB-009).

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-010 | The system **SHALL** give every posting exactly one of the states Draft, Published, or Archived, and **SHALL** allow only these transitions: Draft to Published, Published to Archived, and Archived to Published. | TBD | TBD |
| FR-JOB-011 | The system **SHALL**, when a Hiring Manager publishes a posting, verify the mandatory content and dates of FR-APP-01 and FR-APP-03, **SHALL** keep the posting in Draft if a check fails, and **SHALL** list every failed check. | TBD | TBD |
| FR-JOB-012 | The system **SHALL** allow a posting to be published before a pipeline is published, in which case admitted applications wait as defined in FR-APP-23. | TBD | TBD |
| FR-JOB-013 | The system **SHALL** require the Hiring Manager to choose Public or Private visibility before a posting is published. | TBD | TBD |
| FR-JOB-014 | The system **SHALL** list a public Published posting in the jobs feed, in search results, and on the public company profile page (FR-AUTH-029). | TBD | TBD |
| FR-JOB-015 | The system **SHALL** exclude a private Published posting from the jobs feed, search results, and the public company profile page, and **SHALL** make it reachable only through its private link or an email invitation. | TBD | TBD |
| FR-JOB-016 | The system **SHALL** generate an unguessable private link for a private posting and **SHALL** allow the Hiring Manager to revoke it and generate a new one; a revoked link **SHALL** stop giving access. | TBD | TBD |
| FR-JOB-017 | The system **SHALL** allow the Hiring Manager to send an email invitation containing the private link of a private posting to one or more email addresses. | TBD | TBD |
| FR-JOB-018 | The system **SHALL** allow the Hiring Manager to change the visibility of a Published posting at any time, without changing any existing application. | TBD | TBD |
| FR-JOB-019 | The system **SHALL**, when a Hiring Manager archives a posting, remove it from the feed, search results, and company profile page, stop accepting new applications, and retain the posting, its applications, and its pipeline data. | TBD | TBD |
| FR-JOB-020 | The system **SHALL**, while a posting is Archived, suspend its applications: Interviewees cannot start or continue modules or book live interviews, no stage, pass, reject, or offer decision can be recorded, and every application keeps its submitted data, results, and status. | TBD | TBD |
| FR-JOB-021 | The system **SHALL**, while a posting is Archived, pause the completion windows of its modules (FR-APP-38), **SHALL** not mark any module as missed (FR-APP-41), and **SHALL** extend each window by the time the posting was archived when it is republished. | TBD | TBD |
| FR-JOB-022 | The system **SHALL** show an Interviewee with an application to an Archived posting, in the application tracking view (APP 3.12), that the posting is archived and cannot be completed. | TBD | TBD |
| FR-JOB-023 | The system **SHALL** allow an Interviewee to withdraw an application to an Archived posting (FR-APP-19). | TBD | TBD |
| FR-JOB-024 | The system **SHOULD** show the Hiring Manager the number of active applications and ask for confirmation before archiving a posting. | TBD | TBD |
| FR-JOB-025 | The system **SHOULD** notify Interviewees with an active application when their posting is archived and when it is republished. | TBD | TBD |
| FR-JOB-026 | The system **SHALL** allow the Hiring Manager to republish an Archived posting, which returns it to Published with its previous visibility and resumes its suspended applications, leaving all pass, reject, and offer decisions to the Hiring Manager. | TBD | TBD |
| FR-JOB-027 | The system **SHALL** accept new applications to a republished posting only if its application deadline has not passed (FR-APP-04, FR-APP-17); a posting republished after its deadline **SHALL** remain closed to new applications. | TBD | TBD |
| FR-JOB-028 | The system **SHALL** allow the Hiring Manager to delete a Draft posting that has never received an application, removing it permanently. | TBD | TBD |
| FR-JOB-029 | The system **SHALL** refuse to delete a posting that is Published, Archived, or has applications, and **SHALL** offer archiving instead. | TBD | TBD |
| FR-JOB-030 | The system **SHALL** show the Hiring Manager a list of all postings of the company in every state, with title, state, visibility, number of applications, and application deadline. | TBD | TBD |
| FR-JOB-031 | The system **SHOULD** allow the Hiring Manager to filter the posting list by state and visibility and to search it by title. | TBD | TBD |

#### 5.4.3.3 Posting Tags

**Feature Description:** Tags classify a posting and feed the filters, search, and recommendations.

**Preconditions:**
* A Draft or Published posting exists.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-032 | The system **SHALL** allow the Hiring Manager to add and remove tags on a posting, chosen from the taxonomy managed by Admins (FR-AUTH-039). | TBD | TBD |
| FR-JOB-033 | The system **SHALL** allow at most 10 tags per posting. | TBD | TBD |

#### 5.4.3.4 Jobs Feed

**Feature Description:** The jobs feed lists postings an Interviewee can apply to, with an optional ordering by match to their preferences.

**Preconditions:**
* The user is signed in (FR-JOB-034).

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-034 | The system **SHALL** require a user to be signed in to view the feed, search results, posting details, and companies, **SHALL** send a visitor to sign-in or registration, and **SHALL** open the requested page afterwards. | TBD | TBD |
| FR-JOB-035 | The system **SHALL** list public Published postings that are open or not yet open for applications, newest first, in pages of 20. | TBD | TBD |
| FR-JOB-036 | The system **SHALL** show for each feed item the title, company name, location, work mode, seniority level, employment type, tags, and application deadline, and the salary range only where the posting shows it (FR-APP-02). | TBD | TBD |
| FR-JOB-037 | The system **SHALL** mark feed items that the Interviewee has saved or applied to. | TBD | TBD |
| FR-JOB-038 | The system **SHOULD**, where the Interviewee has set job preferences (FR-AUTH-030), offer a "Recommended" ordering of the feed ranked by match with those preferences. | TBD | TBD |
| FR-JOB-039 | The system **SHOULD** show, for each recommended posting, the preference or tag that caused it to be recommended. | TBD | TBD |
| FR-JOB-040 | The system **SHOULD** allow the Interviewee to limit the feed to postings of followed companies. | TBD | TBD |
| FR-JOB-041 | The system **SHALL** exclude postings of companies that are deactivated or banned (FR-AUTH-037) from the feed and search results. | TBD | TBD |

#### 5.4.3.5 Filtering and Sorting

**Feature Description:** Interviewees narrow and order the feed and search results.

**Preconditions:**
* The user is signed in and the feed or a search result is displayed.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-042 | The system **SHALL** allow filtering by tag, required skill, seniority level, employment type, work mode, location, salary range, company, and posting date range. | TBD | TBD |
| FR-JOB-043 | The system **SHALL** return postings that match every selected filter, where a filter with several selected values matches any one of them. | TBD | TBD |
| FR-JOB-044 | The system **SHALL**, when a salary filter is active, exclude postings whose salary range is not shown to Interviewees, and **SHALL** tell the Interviewee that such postings were excluded. | TBD | TBD |
| FR-JOB-045 | The system **SHALL** allow sorting by newest, nearest deadline, and, for search results, relevance. | TBD | TBD |
| FR-JOB-046 | The system **SHALL** display the active filters and allow removing each one or all at once. | TBD | TBD |
| FR-JOB-047 | The system **SHALL**, when a filter combination returns no postings, display a message and offer to clear the filters. | TBD | TBD |
| FR-JOB-048 | The system **SHALL** apply the same filters and sorting to the feed and to the posting results of a search. | TBD | TBD |

#### 5.4.3.6 Search

**Feature Description:** A keyword search finds postings and companies the user is allowed to see.

**Preconditions:**
* The user is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-049 | The system **SHALL** provide a keyword search that can be started from every page of the platform. | TBD | TBD |
| FR-JOB-050 | The system **SHALL** match the keywords against the title, description, tags, required skills, and company name of postings, and against the names of companies. | TBD | TBD |
| FR-JOB-051 | The system **SHALL** return results grouped by type (postings, companies), ranked by relevance, in pages. | TBD | TBD |
| FR-JOB-052 | The system **SHALL** require at least 2 characters in a search and **SHALL** tell the user when the query is too short. | TBD | TBD |
| FR-JOB-053 | The system **SHALL** return only public Published postings and public company profiles, and **SHALL** never return Draft, Archived, or private postings of other companies. | TBD | TBD |
| FR-JOB-054 | The system **SHALL** reflect a new, edited, archived, or visibility-changed posting in search results within 60 seconds. | TBD | TBD |
| FR-JOB-055 | The system **SHOULD** return results for queries with common spelling variations. | TBD | TBD |

#### 5.4.3.7 Posting Details

**Feature Description:** The detail page shows everything an Interviewee needs to decide whether to apply and what action is possible.

**Preconditions:**
* The user is signed in.
* The posting is public Published, or the user holds its private link or an invitation, or the user has applied to or saved it.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-056 | The system **SHALL** display the posting's title, company (linked to its public profile, FR-AUTH-029), description, location, work mode, seniority level, employment type, tags, application dates, the criteria in the two groups of FR-APP-09, and the salary range only where shown (FR-APP-02). | TBD | TBD |
| FR-JOB-057 | The system **SHALL** display whether the posting is open, opens on a given date, closed, or archived, and for a closed posting **SHALL** state whether the deadline passed or the application cap was reached (FR-APP-17). | TBD | TBD |
| FR-JOB-058 | The system **SHALL** offer an Apply action on an open posting, and, if the Interviewee cannot apply, **SHALL** state what is missing (incomplete resume or an existing application, FR-APP-16). | TBD | TBD |
| FR-JOB-059 | The system **SHALL** offer Save (5.4.3.8) and Follow company (5.4.3.9) actions on the detail page. | TBD | TBD |
| FR-JOB-060 | The system **SHALL** show an Archived posting only to Interviewees who applied to or saved it, with a notice that it no longer accepts applications. | TBD | TBD |

#### 5.4.3.8 Saved Jobs

**Feature Description:** An Interviewee keeps a personal list of postings to return to.

**Preconditions:**
* The Interviewee is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-061 | The system **SHALL** allow an Interviewee to save and unsave a posting from the feed, search results, and detail page. | TBD | TBD |
| FR-JOB-062 | The system **SHALL** show the saved list with each posting's title, company, deadline, and whether it is open, closed, or archived. | TBD | TBD |
| FR-JOB-063 | The system **SHALL**, when a saved posting is no longer visible to the Interviewee (made private or removed), keep its entry labelled "No longer available" and show none of its content. | TBD | TBD |
| FR-JOB-064 | The system **COULD** remind an Interviewee 48 hours before the deadline of a saved posting to which they have not applied. | TBD | TBD |

#### 5.4.3.9 Following Companies

**Feature Description:** An Interviewee follows companies to see their new postings.

**Preconditions:**
* The Interviewee is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-065 | The system **SHOULD** allow an Interviewee to follow and unfollow a company from its public profile page (FR-AUTH-029) and from a posting's detail page. | TBD | TBD |
| FR-JOB-066 | The system **SHOULD** show the Interviewee the list of followed companies. | TBD | TBD |
| FR-JOB-067 | The system **SHOULD** notify an Interviewee within 10 minutes when a followed company publishes a public posting. | TBD | TBD |
| FR-JOB-068 | The system **COULD** allow the Interviewee to turn off posting notifications for an individual followed company. | TBD | TBD |

#### 5.4.3.10 Recommendations, Saved Searches, and New-Jobs Tracker

**Feature Description:** The system suggests relevant postings and alerts Interviewees to new ones.

**Preconditions:**
* The Interviewee is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-069 | The system **COULD** show up to 5 similar postings on a posting's detail page, based on shared tags, required skills, and seniority level. | TBD | TBD |
| FR-JOB-070 | The system **COULD** show similar postings to an Interviewee after an application is submitted (FR-APP-18). | TBD | TBD |
| FR-JOB-071 | The system **SHALL** exclude from recommendations any posting the Interviewee has already applied to. | TBD | TBD |
| FR-JOB-072 | The system **SHALL** base recommendations and the Recommended ordering only on structured posting data and the Interviewee's job preferences, and **SHALL** not use name, photo, gender, age, or nationality. | TBD | TBD |
| FR-JOB-073 | The system **COULD** allow an Interviewee to save a search with its filters under a name, list saved searches, and delete them. | TBD | TBD |
| FR-JOB-074 | The system **COULD** notify an Interviewee when a newly published posting matches a saved search (new-jobs tracker). | TBD | TBD |
| FR-JOB-075 | The system **COULD** allow the Interviewee to receive tracker notifications immediately or as a daily digest. | TBD | TBD |

#### 5.4.3.11 Sharing

**Feature Description:** A posting can be shared outside the platform through its link.

**Preconditions:**
* The posting is Published and the user is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-076 | The system **SHOULD** provide a link to copy and share for a posting, which opens its detail page after sign-in (FR-JOB-034) and, for a private posting, is its private link (FR-JOB-016). | TBD | TBD |

### 5.4.4 Business Rules

* A posting is Draft, Published, or Archived, and only the transitions in FR-JOB-010 exist.
* A posting can be published before a pipeline is published (FR-JOB-012, FR-APP-23).
* Hiring configuration is defined once, in APP 3.1. JOB never redefines it.
* Every user must be signed in to browse. Visitors are sent to sign-in first (FR-JOB-034).
* Only postings that are public and Published appear in the feed, search, and the company profile page (FR-JOB-014, FR-JOB-015, FR-JOB-053).
* A private posting is reached only through its private link or an email invitation (FR-JOB-015 to FR-JOB-017).
* A posting with applications can never be deleted, only archived (FR-JOB-029).
* Archiving suspends applications but never deletes or decides them. Interviewees can still withdraw, and the Hiring Manager decides on the applications after republishing (FR-JOB-020, FR-JOB-023, FR-JOB-026).
* A posting can be reopened to new applications only before its deadline passes (FR-JOB-027, FR-APP-04).
* Changing visibility never changes an existing application (FR-JOB-018).
* A salary range hidden from Interviewees is never revealed through filters or search (FR-JOB-044).
* Saved postings, saved searches, and followed companies belong to the Interviewee. Hiring Managers and Admins never see which Interviewees saved or followed a posting or company.
* Recommendations never use protected or identifying attributes (FR-JOB-072).
* Tags come only from the taxonomy managed by Admins (FR-JOB-032).
* All notifications of this section are sent through the notification mechanism owned by APP.

### 5.4.5 Error Handling

* **Publishing with missing or invalid data:** keep the posting in Draft and list every failed check (FR-JOB-011).
* **Deleting a posting that cannot be deleted:** refuse and offer archiving (FR-JOB-029).
* **Two Hiring Managers editing the same posting:** save the first change, reject the second, and show the second editor the current version.
* **Revoked, invalid, or unknown private link:** show "This posting is not available" without revealing whether the posting exists.
* **Visitor who is not signed in opens a link:** send the visitor to sign-in or registration and open the requested page afterwards (FR-JOB-034).
* **Submission to a posting archived or closed while the Interviewee is filling the form:** reject with the specific reason (FR-APP-17) and keep nothing partial.
* **Interviewee action on an Archived posting (start module, book interview):** refuse and show that the posting is archived (FR-JOB-020).
* **Salary filter with minimum above maximum, or too-short search query:** show a message identifying the invalid input and keep the previous results.
* **Search service unavailable:** show "Search is temporarily unavailable" and keep the feed and filters working.
* **Invitation or notification email fails:** follow the retry and in-platform rules of APP 3.11; for an invitation, tell the Hiring Manager which addresses failed and keep the private link available to copy.
* **Tag removed from the taxonomy by an Admin:** existing postings keep the tag, and it can no longer be added to new ones.