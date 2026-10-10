Owner: Yousef Abood | Status: Draft | Last updated: 2026-10-10 | Jira: AMARA-31

## 5.4 JOB: Jobs, Search & Integrations

### 5.4.1 Overview

This section defines the functional requirements of the JOB epic. JOB covers the life of a job posting as a published offering (creating, publishing, controlling its visibility, archiving, republishing, and deleting it) and how Interviewees discover it (jobs feed, filtering, search, posting details, saved jobs, following companies, recommendations, and sharing). Its actors are the **Hiring Manager**, which is the employer (the company, not a single person), and the **Interviewee**; Admins act on postings through AUTH. It also covers awarding Verified Crests for successfully completed assessments and applying candidate anonymity and identity reveal. Posting management also covers restoring and previewing postings, autosaving drafts with warnings, invitations lists, handling banned companies, posting language and currency, and the employer interface language. Every user must be signed in to browse (FR-JOB-060).

JOB does not define what happens once an Interviewee applies. Hiring configuration of a posting (content fields, deadlines, capacity, criteria, finalist target), application status, pipelines, assessments, interviews, feedback, offers, and notifications belong to APP. This section refers to APP requirements by ID and never repeats them. Every notification mentioned here is sent through the notification mechanism owned by APP.

Use cases are not written yet. Every "Related Use Case" value is TBD until `04-use-cases.md` exists. Priority is expressed with **SHALL** (Must), **SHOULD** (Should), and **COULD** (Could).

Diagrams (TBD): Figure 1: JOB use case diagram. Figure 2: job posting lifecycle state diagram (Draft, Published, Archived).

### 5.4.2 Scope

**In scope**
- Posting creation, editing, versioning and restore, preview, language and currency, autosave and warnings, audit, and the employer interface language (5.4.3.1)
- Publication, visibility, private access and the invitations list, archiving, republishing, deletion, handling of banned companies, and the company posting list (5.4.3.2)
- Posting tags (5.4.3.3)
- Jobs feed, filtering, sorting, and search (5.4.3.4 to 5.4.3.6)
- Posting details, saved jobs and invitations, following companies (5.4.3.7 to 5.4.3.9)
- Recommendations, saved searches, new-jobs tracker, and sharing (5.4.3.10, 5.4.3.11)
- Verified Crest award, candidate anonymity, and identity reveal (5.4.3.12 to 5.4.3.14)

**Out of scope (owned elsewhere)**
- Hiring configuration of a posting: content fields, dates, capacity, Must-have and Nice-to-have criteria, application questions, finalist target (APP 3.1, FR-APP-01 to FR-APP-12)
- Applying, withdrawing, screening, application statuses, pipelines, test modules (including sharing them between companies, as an extension of FR-APP-49), live interviews, advancement, feedback, offers, Interviewee application tracking, and all notifications (APP)
- Company registration, verification, roles and permissions, company public profile page, job preferences, document upload, reporting, and admin management of users and taxonomies (AUTH)
- Interviewee resume and profile (CAND)
- Company-level analytics and cross-posting (ANA)

### 5.4.3 Features and Functional Requirements

#### 5.4.3.1 Job Posting Creation and Editing

**Feature Description:** A Hiring Manager creates a job posting as a draft, completes its content and hiring configuration, and edits it over time with a record of every change. The posting states its language and currency, earlier versions can be restored, the posting can be previewed as Interviewees will see it, and the Hiring Manager can switch the language of the employer-side interface. Drafts are saved automatically, and the Hiring Manager is warned before leaving with unsaved changes.

**Preconditions:**
* The company is registered and verified (AUTH).
* The Hiring Manager is signed in with a company account permitted to manage postings, per the AUTH permission matrix.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-001 | The system **SHALL** allow a Hiring Manager to create a job posting, which is saved in the Draft state. | TBD | TBD |
| FR-JOB-002 | The system **SHALL** allow a Draft posting to be saved with incomplete fields. | TBD | TBD |
| FR-JOB-003 | The system **SHALL** autosave a Draft posting while the Hiring Manager edits it, at least every 30 seconds while there are unsaved changes. | TBD | TBD |
| FR-JOB-004 | The system **SHALL** show the Hiring Manager whether the latest changes to a Draft are saved, with the time of the last save. | TBD | TBD |
| FR-JOB-005 | The system **SHALL** not create a separate version or audit entry (FR-JOB-015) for each autosave. | TBD | TBD |
| FR-JOB-006 | The system **SHALL** apply autosave to Draft postings only; changes to a Published posting take effect only when the Hiring Manager saves them. | TBD | TBD |
| FR-JOB-007 | The system **SHALL** warn the Hiring Manager before leaving the posting editor with unsaved changes, and **SHALL** offer to save, discard, or stay. | TBD | TBD |
| FR-JOB-008 | The system **SHALL** warn the Hiring Manager before leaving a new posting that has not yet been saved as a Draft, and **SHALL** offer to save it as a Draft, discard it, or stay. | TBD | TBD |
| FR-JOB-009 | The system **SHALL** allow the Hiring Manager to enter the content and hiring configuration defined in FR-APP-01 to FR-APP-12 as part of creating the posting. | TBD | TBD |
| FR-JOB-010 | The system **SHALL** not restrict the language in which a posting is written and **SHALL** display the posting text as written, without translating it. | TBD | TBD |
| FR-JOB-011 | The system **SHALL** not restrict the currency of a posting's salary range and **SHALL** display the amount in the chosen currency, without converting it. | TBD | TBD |
| FR-JOB-012 | The system **SHALL** require the Hiring Manager to specify the language of a posting, chosen from a standard list of languages, before the posting is published (FR-JOB-030). | TBD | TBD |
| FR-JOB-013 | The system **SHALL** require the Hiring Manager to specify the currency of the salary range, chosen from a standard list of currencies, before a posting with a salary range is published (FR-JOB-030). | TBD | TBD |
| FR-JOB-014 | The system **SHALL** apply FR-APP-04 (deadline extension only before the deadline passes) and FR-APP-10 (Must-have criteria locked after the first application) when a Hiring Manager edits a Published posting. | TBD | TBD |
| FR-JOB-015 | The system **SHALL** record in an audit log every creation, edit, state change, and deletion of a posting, with the author, the time, and the changed fields. | TBD | TBD |
| FR-JOB-016 | The system **SHALL** keep each previous version of a Published posting and allow the Hiring Manager to view it. | TBD | TBD |
| FR-JOB-017 | The system **SHALL** allow the Hiring Manager to restore an earlier version of a posting, which creates a new current version with the content of the selected version and keeps all other versions. | TBD | TBD |
| FR-JOB-018 | The system **SHALL**, when restoring a version, keep the current Must-have criteria if the posting has received an application (FR-APP-10), keep the current application deadline unless the restored deadline is later and the current one has not passed (FR-APP-04), and tell the Hiring Manager which fields were not restored. | TBD | TBD |
| FR-JOB-019 | The system **SHALL** record each restoration in the audit log with the author, the time, and the version restored (FR-JOB-015). | TBD | TBD |
| FR-JOB-020 | The system **SHOULD** notify Interviewees with an active application when the title, salary range, location, or work mode of the posting changes. | TBD | TBD |
| FR-JOB-021 | The system **COULD** allow a Hiring Manager to duplicate a posting into a new Draft that copies its content and criteria but not its applications or dates. | TBD | TBD |
| FR-JOB-022 | The system **SHALL** allow the Hiring Manager to preview a posting, including a Draft and unsaved edits, as an Interviewee will see it in the jobs feed (FR-JOB-062) and on the detail page (FR-JOB-083). | TBD | TBD |
| FR-JOB-023 | The system **SHALL** not change any posting, application, or notification data, and **SHALL** not notify anyone, when a posting is previewed. | TBD | TBD |
| FR-JOB-024 | The system **SHALL** allow only the Hiring Manager, through accounts permitted by the AUTH permission matrix, and Admins (FR-AUTH-036), to create, edit, archive, republish, or delete a posting of a company. | TBD | TBD |
| FR-JOB-025 | The system **SHALL** allow the Hiring Manager to switch the language of the employer-side interface at any time, among the interface languages defined in `10-i18n.md`. | TBD | TBD |
| FR-JOB-026 | The system **SHALL** apply the chosen language immediately, without requiring the Hiring Manager to sign in again and without losing unsaved changes. | TBD | TBD |
| FR-JOB-027 | The system **SHOULD** keep the chosen language for the Hiring Manager's later sessions. | TBD | TBD |
| FR-JOB-028 | The system **SHALL** not change the content of any posting when the interface language is switched. | TBD | TBD |

#### 5.4.3.2 Publication, Visibility, and Lifecycle

**Feature Description:** A posting moves between Draft, Published, and Archived. A Published posting is either public (listed in the feed and search) or private (reachable only through a private link or an email invitation). Archiving suspends the posting's applications without deleting them, and republishing resumes them. When a company is banned, its Published postings are archived until the ban is resolved. The Hiring Manager sees the invitations sent for each private posting, and deleting a posting needs confirmation.

**Preconditions:**
* The posting exists and the Hiring Manager is permitted to manage it (FR-JOB-024).

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-029 | The system **SHALL** give every posting exactly one of the states Draft, Published, or Archived, and **SHALL** allow only these transitions: Draft to Published, Published to Archived, and Archived to Published. | TBD | TBD |
| FR-JOB-030 | The system **SHALL**, when a Hiring Manager publishes a posting, verify the mandatory content and dates of FR-APP-01 and FR-APP-03 and the language and currency required by FR-JOB-012 and FR-JOB-013, **SHALL** keep the posting in Draft if a check fails, and **SHALL** list every failed check. | TBD | TBD |
| FR-JOB-031 | The system **SHALL** allow a posting to be published before a pipeline is published, in which case admitted applications wait as defined in FR-APP-23. | TBD | TBD |
| FR-JOB-032 | The system **SHALL** require the Hiring Manager to choose Public or Private visibility before a posting is published. | TBD | TBD |
| FR-JOB-033 | The system **SHALL** list a public Published posting in the jobs feed, in search results, and on the public company profile page (FR-AUTH-029). | TBD | TBD |
| FR-JOB-034 | The system **SHALL** exclude a private Published posting from the jobs feed, search results, and the public company profile page, and **SHALL** make it reachable only through its private link or an email invitation. | TBD | TBD |
| FR-JOB-035 | The system **SHALL** generate an unguessable private link for a private posting and **SHALL** allow the Hiring Manager to revoke it and generate a new one; a revoked link **SHALL** stop giving access. | TBD | TBD |
| FR-JOB-036 | The system **SHALL** allow the Hiring Manager to send an email invitation containing the private link of a private posting to one or more email addresses. | TBD | TBD |
| FR-JOB-037 | The system **SHALL** show the Hiring Manager, for each private posting, the list of its email invitations with the email address, the time sent, and the delivery status (sent or failed). | TBD | TBD |
| FR-JOB-038 | The system **SHALL** allow the Hiring Manager to change the visibility of a Published posting at any time, without changing any existing application. | TBD | TBD |
| FR-JOB-039 | The system **SHALL**, when a Hiring Manager archives a posting, remove it from the feed, search results, and company profile page, stop accepting new applications, and retain the posting, its applications, and its pipeline data. | TBD | TBD |
| FR-JOB-040 | The system **SHALL**, while a posting is Archived, suspend its applications: Interviewees cannot start or continue modules or book live interviews, no stage, pass, reject, or offer decision can be recorded, and every application keeps its submitted data, results, and status. | TBD | TBD |
| FR-JOB-041 | The system **SHALL**, while a posting is Archived, pause the completion windows of its modules (FR-APP-38), **SHALL** not mark any module as missed (FR-APP-41), and **SHALL** extend each window by the time the posting was archived when it is republished. | TBD | TBD |
| FR-JOB-042 | The system **SHALL** show an Interviewee with an application to an Archived posting, in the application tracking view (APP 3.12), that the posting is archived and cannot be completed. | TBD | TBD |
| FR-JOB-043 | The system **SHALL** allow an Interviewee to withdraw an application to an Archived posting (FR-APP-19). | TBD | TBD |
| FR-JOB-044 | The system **SHOULD** show the Hiring Manager the number of active applications and ask for confirmation before archiving a posting. | TBD | TBD |
| FR-JOB-045 | The system **SHOULD** notify Interviewees with an active application when their posting is archived and when it is republished. | TBD | TBD |
| FR-JOB-046 | The system **SHALL** allow the Hiring Manager to republish an Archived posting, which returns it to Published with its previous visibility and resumes its suspended applications, leaving all pass, reject, and offer decisions to the Hiring Manager. | TBD | TBD |
| FR-JOB-047 | The system **SHALL** accept new applications to a republished posting only if its application deadline has not passed (FR-APP-04, FR-APP-17); a posting republished after its deadline **SHALL** remain closed to new applications. | TBD | TBD |
| FR-JOB-048 | The system **SHALL**, when a company is banned, archive every Published posting of the company, which suspends all their applications (FR-JOB-040, FR-JOB-041), and **SHALL** keep Draft postings as Draft. | TBD | TBD |
| FR-JOB-049 | The system **SHALL** not allow a banned company to publish or republish any posting while the ban lasts. | TBD | TBD |
| FR-JOB-050 | The system **SHALL**, when the ban is resolved, return each posting archived because of the ban to Published with its previous visibility and resume its applications as in FR-JOB-046, and **SHALL** keep postings that were Archived before the ban as Archived. | TBD | TBD |
| FR-JOB-051 | The system **SHOULD** notify Interviewees with an active application when their posting is archived because of a ban and when it is resumed, without stating the reason. | TBD | TBD |
| FR-JOB-052 | The system **SHALL** allow the Hiring Manager to delete a Draft posting that has never received an application, removing it permanently. | TBD | TBD |
| FR-JOB-053 | The system **SHALL** refuse a Hiring Manager's request to delete a posting that is Published, Archived, or has applications, and **SHALL** offer archiving instead. | TBD | TBD |
| FR-JOB-054 | The system **SHALL** allow an Admin (FR-AUTH-036) to delete a posting in any state at any time, including a posting that has applications, and **SHALL** record the deletion in the audit log (FR-JOB-015). | TBD | TBD |
| FR-JOB-055 | The system **SHALL** ask for confirmation before a posting is deleted, **SHALL** state that deletion is permanent, and **SHALL** show an Admin the number of applications affected when the posting has applications. | TBD | TBD |
| FR-JOB-056 | The system **SHALL** show the Hiring Manager a list of all postings of the company in every state, with title, state, visibility, number of applications, and application deadline. | TBD | TBD |
| FR-JOB-057 | The system **SHOULD** allow the Hiring Manager to filter the posting list by state and visibility and to search it by title. | TBD | TBD |

#### 5.4.3.3 Posting Tags

**Feature Description:** Tags classify a posting and feed the filters, search, and recommendations.

**Preconditions:**
* A Draft or Published posting exists.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-058 | The system **SHALL** allow the Hiring Manager to add and remove tags on a posting, chosen from the taxonomy managed by Admins (FR-AUTH-039). | TBD | TBD |
| FR-JOB-059 | The system **SHALL** allow at most 10 tags per posting. | TBD | TBD |

#### 5.4.3.4 Jobs Feed

**Feature Description:** The jobs feed lists postings an Interviewee can apply to, with an optional ordering by match to their preferences.

**Preconditions:**
* The user is signed in (FR-JOB-060).

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-060 | The system **SHALL** require a user to be signed in to view the feed, search results, posting details, and companies, **SHALL** send a visitor to sign-in or registration, and **SHALL** open the requested page afterwards. | TBD | TBD |
| FR-JOB-061 | The system **SHALL** list public Published postings that are open or not yet open for applications, newest first, in pages of 20. | TBD | TBD |
| FR-JOB-062 | The system **SHALL** show for each feed item the title, company name, location, work mode, seniority level, employment type, posting language, tags, and application deadline, and the salary range only where the posting shows it (FR-APP-02). | TBD | TBD |
| FR-JOB-063 | The system **SHALL** mark feed items that the Interviewee has saved or applied to. | TBD | TBD |
| FR-JOB-064 | The system **SHOULD**, where the Interviewee has set job preferences (FR-AUTH-030), offer a "Recommended" ordering of the feed ranked by match with those preferences. | TBD | TBD |
| FR-JOB-065 | The system **SHOULD** show, for each recommended posting, the preference or tag that caused it to be recommended. | TBD | TBD |
| FR-JOB-066 | The system **SHOULD** allow the Interviewee to limit the feed to postings of followed companies. | TBD | TBD |
| FR-JOB-067 | The system **SHALL** exclude postings of companies that are deactivated or banned (FR-AUTH-037) from the feed and search results. | TBD | TBD |

#### 5.4.3.5 Filtering and Sorting

**Feature Description:** Interviewees narrow and order the feed and search results.

**Preconditions:**
* The user is signed in and the feed or a search result is displayed.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-068 | The system **SHALL** allow filtering by tag, required skill, seniority level, employment type, work mode, location, salary range, salary currency, posting language, company, and posting date range. | TBD | TBD |
| FR-JOB-069 | The system **SHALL** return postings that match every selected filter, where a filter with several selected values matches any one of them. | TBD | TBD |
| FR-JOB-070 | The system **SHALL**, when a salary filter is active, exclude postings whose salary range is not shown to Interviewees, and **SHALL** tell the Interviewee that such postings were excluded. | TBD | TBD |
| FR-JOB-071 | The system **SHALL** apply a salary range filter (FR-JOB-068) only to postings whose salary range is in the currency the Interviewee selects with the range, and **SHALL** not convert amounts between currencies. | TBD | TBD |
| FR-JOB-072 | The system **SHALL** allow sorting by newest, nearest deadline, and, for search results, relevance. | TBD | TBD |
| FR-JOB-073 | The system **SHALL** display the active filters and allow removing each one or all at once. | TBD | TBD |
| FR-JOB-074 | The system **SHALL**, when a filter combination returns no postings, display a message and offer to clear the filters. | TBD | TBD |
| FR-JOB-075 | The system **SHALL** apply the same filters and sorting to the feed and to the posting results of a search. | TBD | TBD |

#### 5.4.3.6 Search

**Feature Description:** A keyword search finds postings and companies the user is allowed to see.

**Preconditions:**
* The user is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-076 | The system **SHALL** provide a keyword search that can be started from every page of the platform. | TBD | TBD |
| FR-JOB-077 | The system **SHALL** match the keywords against the title, description, tags, required skills, and company name of postings, and against the names of companies. | TBD | TBD |
| FR-JOB-078 | The system **SHALL** return results grouped by type (postings, companies), ranked by relevance, in pages. | TBD | TBD |
| FR-JOB-079 | The system **SHALL** require at least 2 characters in a search and **SHALL** tell the user when the query is too short. | TBD | TBD |
| FR-JOB-080 | The system **SHALL** return only public Published postings and public company profiles, and **SHALL** never return Draft, Archived, or private postings of other companies. | TBD | TBD |
| FR-JOB-081 | The system **SHALL** reflect a new, edited, archived, or visibility-changed posting in search results within 60 seconds. | TBD | TBD |
| FR-JOB-082 | The system **SHOULD** return results for queries with common spelling variations. | TBD | TBD |

#### 5.4.3.7 Posting Details

**Feature Description:** The detail page shows everything an Interviewee needs to decide whether to apply and what action is possible.

**Preconditions:**
* The user is signed in.
* The posting is public Published, or the user holds its private link or an invitation, or the user has applied to or saved it.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-083 | The system **SHALL** display the posting's title, company (linked to its public profile, FR-AUTH-029), description, location, work mode, seniority level, employment type, posting language, tags, application dates, the criteria in the two groups of FR-APP-09, and the salary range only where shown (FR-APP-02). | TBD | TBD |
| FR-JOB-084 | The system **SHALL** display whether the posting is open, opens on a given date, closed, or archived, and for a closed posting **SHALL** state whether the deadline passed or the application cap was reached (FR-APP-17). | TBD | TBD |
| FR-JOB-085 | The system **SHALL** offer an Apply action on an open posting, and, if the Interviewee cannot apply, **SHALL** state what is missing (incomplete resume or an existing application, FR-APP-16). | TBD | TBD |
| FR-JOB-086 | The system **SHALL** offer Save (5.4.3.8) and Follow company (5.4.3.9) actions on the detail page. | TBD | TBD |
| FR-JOB-087 | The system **SHALL** show an Archived posting only to Interviewees who applied to or saved it, with a notice that it no longer accepts applications. | TBD | TBD |

#### 5.4.3.8 Saved Jobs and Invitations

**Feature Description:** An Interviewee keeps a personal list of postings to return to and sees the invitations received to private postings.

**Preconditions:**
* The Interviewee is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-088 | The system **SHALL** allow an Interviewee to save and unsave a posting from the feed, search results, and detail page. | TBD | TBD |
| FR-JOB-089 | The system **SHALL** show the saved list with each posting's title, company, deadline, and whether it is open, closed, or archived. | TBD | TBD |
| FR-JOB-090 | The system **SHALL**, when a saved posting is no longer visible to the Interviewee (made private or removed), keep its entry labelled "No longer available" and show none of its content. | TBD | TBD |
| FR-JOB-091 | The system **COULD** remind an Interviewee 48 hours before the deadline of a saved posting to which they have not applied. | TBD | TBD |
| FR-JOB-092 | The system **SHALL** show an Interviewee the list of invitations received to private postings, with the posting title, company, date invited, and whether the posting is open or closed. | TBD | TBD |
| FR-JOB-093 | The system **SHALL** match an invitation to an Interviewee by the verified email address of the Interviewee's account, including an account created after the invitation was sent. | TBD | TBD |
| FR-JOB-094 | The system **SHALL** let the Interviewee open a posting from the invitations list, and **SHALL** label an invitation "No longer available", showing none of the posting's content, when the posting was deleted or archived or its private link was revoked. | TBD | TBD |

#### 5.4.3.9 Following Companies

**Feature Description:** An Interviewee follows companies to see their new postings.

**Preconditions:**
* The Interviewee is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-095 | The system **SHOULD** allow an Interviewee to follow and unfollow a company from its public profile page (FR-AUTH-029) and from a posting's detail page. | TBD | TBD |
| FR-JOB-096 | The system **SHOULD** show the Interviewee the list of followed companies. | TBD | TBD |
| FR-JOB-097 | The system **SHOULD** notify an Interviewee within 10 minutes when a followed company publishes a public posting. | TBD | TBD |
| FR-JOB-098 | The system **COULD** allow the Interviewee to turn off posting notifications for an individual followed company. | TBD | TBD |

#### 5.4.3.10 Recommendations, Saved Searches, and New-Jobs Tracker

**Feature Description:** The system suggests relevant postings and alerts Interviewees to new ones.

**Preconditions:**
* The Interviewee is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-099 | The system **COULD** show up to 5 similar postings on a posting's detail page, based on shared tags, required skills, and seniority level. | TBD | TBD |
| FR-JOB-100 | The system **COULD** show similar postings to an Interviewee after an application is submitted (FR-APP-18). | TBD | TBD |
| FR-JOB-101 | The system **SHALL** exclude from recommendations any posting the Interviewee has already applied to. | TBD | TBD |
| FR-JOB-102 | The system **SHALL** base recommendations and the Recommended ordering only on structured posting data and the Interviewee's job preferences, and **SHALL** not use name, photo, gender, age, or nationality. | TBD | TBD |
| FR-JOB-103 | The system **COULD** allow an Interviewee to save a search with its filters under a name, list saved searches, and delete them. | TBD | TBD |
| FR-JOB-104 | The system **COULD** notify an Interviewee when a newly published posting matches a saved search (new-jobs tracker). | TBD | TBD |
| FR-JOB-105 | The system **COULD** allow the Interviewee to receive tracker notifications immediately or as a daily digest. | TBD | TBD |

#### 5.4.3.11 Sharing

**Feature Description:** A posting can be shared outside the platform through its link.

**Preconditions:**
* The posting is Published and the user is signed in.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-106 | The system **SHOULD** provide a link to copy and share for a posting, which opens its detail page after sign-in (FR-JOB-060) and, for a private posting, is its private link (FR-JOB-035). | TBD | TBD |

#### 5.4.3.12 Award Verified Crest

**Feature Description:** An Interviewee who successfully completes an assessment receives a Verified Crest, a credential tied to the skill or competency that the assessment verified.

**Preconditions:**
* The Interviewee has completed an assessment module (APP 3.5).
* The assessment is one for which a Crest can be awarded (see open points).

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-107 | The system **SHALL**, when an Interviewee completes an assessment, determine whether the Interviewee completed it successfully, according to the criteria required for its Crest. | TBD | TBD |
| FR-JOB-108 | The system **SHALL**, when the required criteria are satisfied, award the corresponding Crest to the Interviewee. | TBD | TBD |
| FR-JOB-109 | The system **SHALL**, when a Crest is awarded, associate it with the skill or competency verified by the assessment. | TBD | TBD |
| FR-JOB-110 | The system **SHALL**, when a Crest is awarded, associate it with the Interviewee's profile. | TBD | TBD |

#### 5.4.3.13 Apply Interviewee Anonymity

**Feature Description:** During anonymous stages of a pipeline, identifying information about the Interviewee is hidden from Hiring Managers, while the information needed to evaluate professional qualifications stays visible.

**Preconditions:**
* The Hiring Manager has enabled anonymity for the pipeline and chosen the identity-reveal point (FR-APP-97).
* The Interviewee has been admitted to the pipeline.

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-111 | The system **COULD**, when a Hiring Manager opens an application, determine whether anonymity is enabled for the application's current stage. | TBD | TBD |
| FR-JOB-112 | The system **COULD**, while anonymity applies to a stage, hide the Interviewee information designated as non-disclosable. | TBD | TBD |
| FR-JOB-113 | The system **COULD**, while anonymity applies to a stage, continue to display the information needed to evaluate the Interviewee's professional qualifications, such as scores, module results, skills, and Crests. | TBD | TBD |
| FR-JOB-114 | The system **COULD**, while anonymity applies to a stage, prevent hidden Interviewee information from being displayed to authorized evaluators in any view, including resume copies, the dashboard, scorecards, and notes. | TBD | TBD |

#### 5.4.3.14 Reveal Interviewee Identity

**Feature Description:** When an Interviewee reaches the configured identity-disclosure stage, the hidden identity is made available to authorized users, and the application history is kept.

**Preconditions:**
* Anonymity applies to the pipeline (5.4.3.13).
* An identity-reveal point is configured (FR-APP-97).

| ID | Requirement | Related Use Case | Jira |
|---|---|---|---|
| FR-JOB-115 | The system **COULD** determine when an Interviewee reaches the identity-disclosure stage configured for the pipeline. | TBD | TBD |
| FR-JOB-116 | The system **COULD**, when the Interviewee reaches that stage, make the Interviewee's identity available to authorized users. | TBD | TBD |
| FR-JOB-117 | The system **COULD**, after identity disclosure, retain the Interviewee's application and assessment history. | TBD | TBD |

### 5.4.4 Business Rules

* A posting is Draft, Published, or Archived, and only the transitions in FR-JOB-029 exist.
* A posting can be published before a pipeline is published (FR-JOB-031, FR-APP-23).
* Hiring configuration is defined once, in APP 3.1. JOB never redefines it.
* Every user must be signed in to browse. Visitors are sent to sign-in first (FR-JOB-060).
* Only postings that are public and Published appear in the feed, search, and the company profile page (FR-JOB-033, FR-JOB-034, FR-JOB-080).
* A private posting is reached only through its private link or an email invitation (FR-JOB-034 to FR-JOB-036).
* A Hiring Manager can delete only a Draft posting that never received an application, and can otherwise only archive it. An Admin can delete a posting in any state at any time (FR-JOB-052 to FR-JOB-054).
* Archiving suspends applications but never deletes or decides them. Interviewees can still withdraw, and the Hiring Manager decides on the applications after republishing (FR-JOB-040, FR-JOB-043, FR-JOB-046).
* A posting can be reopened to new applications only before its deadline passes (FR-JOB-047, FR-APP-04).
* Changing visibility never changes an existing application (FR-JOB-038).
* A salary range hidden from Interviewees is never revealed through filters or search (FR-JOB-070).
* Saved postings, saved searches, and followed companies belong to the Interviewee. The Hiring Manager and Admins never see which Interviewees saved or followed a posting or company.
* Recommendations never use protected or identifying attributes (FR-JOB-102).
* Tags come only from the taxonomy managed by Admins (FR-JOB-058).
* All notifications of this section are sent through the notification mechanism owned by APP.
* A Verified Crest is awarded only for the successful completion of an assessment and is tied to the skill or competency it verified (FR-JOB-107 to FR-JOB-109).
* While anonymity applies, qualification information stays visible and identifying information stays hidden in every view (FR-JOB-113, FR-JOB-114).
* Identity disclosure never removes the Interviewee's application or assessment history (FR-JOB-117).
* Restoring a version never overwrites history and never bypasses the locked Must-have criteria or the deadline rule (FR-JOB-017, FR-JOB-018).
* Previewing a posting never changes data or notifies anyone (FR-JOB-023).
* While a company is banned, its postings stay archived and nothing can be published. Applications resume when the ban is resolved, and the Hiring Manager then decides on them (FR-JOB-048 to FR-JOB-050).
* Every posting states its language and, when it shows a salary range, its currency. Postings are shown as written, without translation or conversion, and salary filters never convert currencies (FR-JOB-010 to FR-JOB-013, FR-JOB-071).
* Switching the employer interface language never changes posting content (FR-JOB-028).
* Drafts are saved automatically. Changes to a Published posting are saved only on request, and the Hiring Manager is warned before leaving with unsaved changes, before leaving an unsaved new posting, and before deleting a posting (FR-JOB-003 to FR-JOB-008, FR-JOB-055).
* The Hiring Manager sees the invitations of each private posting, and an Interviewee sees only the invitations sent to the email address of their own account (FR-JOB-037, FR-JOB-092 to FR-JOB-094).

### 5.4.5 Error Handling

* **Publishing with missing or invalid data:** keep the posting in Draft and list every failed check (FR-JOB-030).
* **Hiring Manager deleting a posting that cannot be deleted:** refuse and offer archiving (FR-JOB-053).
* **Two accounts of the same Hiring Manager editing the same posting:** save the first change, reject the second, and show the second editor the current version.
* **Revoked, invalid, or unknown private link:** show "This posting is not available" without revealing whether the posting exists.
* **Visitor who is not signed in opens a link:** send the visitor to sign-in or registration and open the requested page afterwards (FR-JOB-060).
* **Submission to a posting archived or closed while the Interviewee is filling the form:** reject with the specific reason (FR-APP-17) and keep nothing partial.
* **Interviewee action on an Archived posting (start module, book interview):** refuse and show that the posting is archived (FR-JOB-040).
* **Salary filter with minimum above maximum, or too-short search query:** show a message identifying the invalid input and keep the previous results.
* **Search service unavailable:** show "Search is temporarily unavailable" and keep the feed and filters working.
* **Invitation or notification email fails:** follow the retry and in-platform rules of APP 3.11; for an invitation, mark the address as failed in the invitations list (FR-JOB-037) and keep the private link available to copy.
* **Crest cannot be awarded or attached to the profile (service or data error):** keep the assessment result, retry, and do not mark the Crest as awarded until it is attached.
* **Tag removed from the taxonomy by an Admin:** existing postings keep the tag, and it can no longer be added to new ones.
* **Publishing without a language, or with a salary range but no currency:** keep the posting in Draft and list it as a failed check (FR-JOB-030, FR-JOB-012, FR-JOB-013).
* **Restoring a version that conflicts with a locked field:** restore the other fields and tell the Hiring Manager which fields were kept (FR-JOB-018).
* **Hiring Manager edits or publishes while the company is banned:** refuse and show that the account is banned (FR-JOB-049).
* **Autosave fails (network or service error):** tell the Hiring Manager that the latest changes are not saved, keep the entered data on screen, and retry (FR-JOB-003).