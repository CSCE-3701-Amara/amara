## 4. Use cases

Owner: TBD | Status: Draft | Last updated: 2026-10-08 | Jira: HIRE-n
---

## 4.1 Actors

// This is only an example.
- Candidate
- Employer
- Administrator
- System

## 4.2 AUTH Use Cases

## 4.3 CAND Use Cases

## 4.4 APP Use Cases

## 4.5 JOB Use Cases

### 4.5.1 JOB Use Case Index

| ID | Use Case | Primary Actor | Requirements |
|---|---|---|---|
| UC-JOB-01 | Create and Edit Job Posting | Hiring Manager | FR-JOB-001 to 009, 014, 015, 020, 021, 024 |
| UC-JOB-02 | Set Posting Language, Currency and Tags | Hiring Manager | FR-JOB-010 to 013, 058, 059 |
| UC-JOB-03 | Preview Job Posting | Hiring Manager | FR-JOB-022, 023 |
| UC-JOB-04 | View and Restore Posting Version | Hiring Manager | FR-JOB-015 to 019 |
| UC-JOB-05 | Publish Job Posting | Hiring Manager | FR-JOB-029 to 034 |
| UC-JOB-06 | Manage Private Access and Invitations | Hiring Manager | FR-JOB-034 to 037 |
| UC-JOB-07 | Change Posting Visibility | Hiring Manager | FR-JOB-038 |
| UC-JOB-08 | Archive Job Posting | Hiring Manager | FR-JOB-039 to 045 |
| UC-JOB-09 | Republish Archived Posting | Hiring Manager | FR-JOB-046, 047 |
| UC-JOB-10 | Delete Job Posting | Hiring Manager, System Admin | FR-JOB-052 to 055 |
| UC-JOB-11 | Handle Postings of a Banned Company | System, System Admin | FR-JOB-048 to 051, 067 |
| UC-JOB-12 | View Company Posting List | Hiring Manager | FR-JOB-056, 057 |
| UC-JOB-13 | Switch Employer Interface Language | Hiring Manager | FR-JOB-025 to 028 |
| UC-JOB-14 | Browse Jobs Feed | Interviewee | FR-JOB-060 to 067 |
| UC-JOB-15 | Filter and Sort Postings | Interviewee | FR-JOB-068 to 075 |
| UC-JOB-16 | Search Postings and Companies | Interviewee | FR-JOB-076 to 081 |
| UC-JOB-17 | View Posting Details | Interviewee | FR-JOB-082 to 086 |
| UC-JOB-18 | Save and Unsave Posting | Interviewee | FR-JOB-087 to 090 |
| UC-JOB-19 | View Received Invitations | Interviewee | FR-JOB-091 to 093 |
| UC-JOB-20 | Follow Company | Interviewee | FR-JOB-094 to 097 |
| UC-JOB-21 | Get Recommendations | Interviewee | FR-JOB-064, 065, 098 to 101 |
| UC-JOB-22 | Manage Saved Searches and New-Jobs Tracker | Interviewee | FR-JOB-102 to 104 |
| UC-JOB-23 | Share Posting | Interviewee | FR-JOB-105 |
| UC-JOB-24 | Award Verified Crest | System, Interviewee | FR-JOB-106 to 109 |
| UC-JOB-25 | Apply Interviewee Anonymity | System, Hiring Manager | FR-JOB-110 to 113 |
| UC-JOB-26 | Reveal Interviewee Identity | System, Hiring Manager | FR-JOB-114 to 116 |

---

### UC-JOB-01: Create and Edit Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-01 |
| Use Case Name     | Create and Edit Job Posting |
| Actor(s)          | Hiring Manager, System, Interviewee (notified) |
| Description       | A Hiring Manager creates a job posting as a Draft, completes its content and hiring configuration, and edits it over time with every change recorded. |
| Preconditions     | The company is registered and verified (AUTH). The Hiring Manager is signed in with an account permitted to manage postings (FR-JOB-024) and the company is not banned. |
| Trigger           | The Hiring Manager chooses to create a posting or to edit an existing one. |
| Main Flow         | 1. The Hiring Manager opens the posting editor.<br>2. The Hiring Manager enters the content and hiring configuration (FR-APP-01 to FR-APP-12, see UC-APP-01).<br>3. The system autosaves the Draft while there are unsaved changes and shows the time of the last save (FR-JOB-003, FR-JOB-004).<br>4. The Hiring Manager saves the posting; the system stores it in the Draft state and records the creation or edit in the audit log (FR-JOB-001, FR-JOB-015). |
| Alternative Flows | 2a. The Hiring Manager saves a Draft with incomplete fields (FR-JOB-002).<br>2b. The Hiring Manager duplicates an existing posting into a new Draft that copies content and criteria but not applications or dates (FR-JOB-021, optional).<br>3a. The posting is Published: autosave does not apply and changes take effect only when the Hiring Manager saves (FR-JOB-006).<br>4a. After saving a Published posting whose title, salary range, location, or work mode changed, the system notifies Interviewees with an active application (FR-JOB-020). |
| Exception Flows   | E1. The Hiring Manager leaves the editor with unsaved changes: the system offers to save, discard, or stay (FR-JOB-007).<br>E2. The Hiring Manager leaves a new posting never saved as a Draft: the system offers to save it as a Draft, discard it, or stay (FR-JOB-008).<br>E3. Autosave fails: the system states that the latest changes are not saved, keeps the entered data on screen (FR-JOB-003).<br>E4. A Published posting is edited against FR-APP-04 or FR-APP-10 (deadline shortened, Must-have changed after the first application): the system rejects that change (FR-JOB-014).<br>E5. Two accounts of the same company edit the same posting: the system saves the first change, rejects the second, and shows the current version.<br>E6. The company is banned: the system refuses the edit and shows that the account is banned (FR-JOB-049).<br>E7. The account is not permitted by the AUTH permission matrix: the system denies the action (FR-JOB-024). |
| Postconditions    | The posting is stored (Draft, or updated Published) and the change is in the audit log. Autosaves do not create versions or audit entries. |
| Business Rules    | Autosave applies to Drafts only (FR-JOB-006). Autosaves create no separate version or audit entry (FR-JOB-005). Hiring configuration is defined once, in APP 3.1. |

### UC-JOB-02: Set Posting Language, Currency and Tags

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-02 |
| Use Case Name     | Set Posting Language, Currency and Tags |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager states the language and salary currency of a posting and classifies it with tags used by filters, search, and recommendations. |
| Preconditions     | A Draft or Published posting exists and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager edits the language, currency, or tags of a posting. |
| Main Flow         | 1. The Hiring Manager selects the posting language from the standard list of languages (FR-JOB-012).<br>2. If the posting has a salary range, the Hiring Manager selects its currency from the standard list of currencies (FR-JOB-013).<br>3. The Hiring Manager adds or removes tags chosen from the System Admin-managed taxonomy (FR-JOB-058).<br>4. The system saves the choices and records them in the audit log (FR-JOB-015). |
| Alternative Flows | 1. The Hiring Manager writes the posting in any language and uses any currency; the system does not restrict either (FR-JOB-010, FR-JOB-011). |
| Exception Flows   | E1. The Hiring Manager adds an 11th tag: the system rejects it (FR-JOB-059).<br>E2. A tag was removed from the taxonomy by a System Admin: the posting keeps it, and it cannot be added to other postings.<br>E3. Publishing without a language, or with a salary range and no currency: the posting stays in Draft and the missing item is listed (FR-JOB-030, see UC-JOB-05). |
| Postconditions    | The posting has a language, a currency (if it shows a salary range), and 0 to 10 tags. |
| Business Rules    | Posting text is shown as written and amounts are shown in the chosen currency, without translation or conversion (FR-JOB-010, FR-JOB-011). Tags come only from the System Admin taxonomy (FR-JOB-058). |

### UC-JOB-03: Preview Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-03 |
| Use Case Name     | Preview Job Posting |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager sees a posting as Interviewees will see it, before or after publication. |
| Preconditions     | A posting exists, in any state, and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager chooses Preview in the editor. |
| Main Flow         | 1. The Hiring Manager chooses Preview, including with unsaved edits.<br>2. The system shows the posting as it appears in the jobs feed (FR-JOB-062) and on the detail page (FR-JOB-082).<br>3. The Hiring Manager closes the preview and returns to the editor. |
| Alternative Flows | 1. A Draft posting is previewed before it has ever been published. |
| Exception Flows   | E1. The posting cannot be rendered because required fields are missing: the system shows the preview with the missing fields marked. |
| Postconditions    | No posting, application, or notification data changed, and nobody was notified. |
| Business Rules    | Previewing never changes data or notifies anyone (FR-JOB-023). |

### UC-JOB-04: View and Restore Posting Version

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-04 |
| Use Case Name     | View and Restore Posting Version |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager reviews earlier versions of a Published posting and restores one as a new current version, with history kept. |
| Preconditions     | The posting has at least one previous version and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager opens the version history of a posting. |
| Main Flow         | 1. The system lists the previous versions with author and time (FR-JOB-016).<br>2. The Hiring Manager opens a version to view it.<br>3. The Hiring Manager selects Restore.<br>4. The system creates a new current version with the content of the selected version and keeps all other versions (FR-JOB-017).<br>5. The system records the restoration with author, time, and version restored (FR-JOB-019). |
| Alternative Flows | 1. The Hiring Manager only views versions and does not restore. |
| Exception Flows   | E1. The posting has received an application: the system keeps the current Must-have criteria (FR-APP-10).<br>E2. The restored deadline is not later than the current one, or the current one has passed: the system keeps the current deadline (FR-APP-04).<br>E3. In both cases the system restores the other fields and tells the Hiring Manager which fields were not restored (FR-JOB-018).<br>E4. The company is banned: the system refuses (FR-JOB-049). |
| Postconditions    | A new current version exists, all earlier versions remain, and the restoration is in the audit log. |
| Business Rules    | Restoring never overwrites history and never bypasses locked Must-have criteria or the deadline rule (FR-JOB-017, FR-JOB-018). |

### UC-JOB-05: Publish Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-05 |
| Use Case Name     | Publish Job Posting |
| Actor(s)          | Hiring Manager, System, Interviewee (notified) |
| Description       | The Hiring Manager moves a Draft posting to Published after the system verifies it is complete, choosing whether it is public or private. |
| Preconditions     | The posting is in Draft. The Hiring Manager is permitted to manage it and the company is not banned. |
| Trigger           | The Hiring Manager chooses Publish. |
| Main Flow         | 1. The Hiring Manager chooses Public or Private visibility (FR-JOB-032).<br>2. The system verifies the mandatory content and dates (FR-APP-01, FR-APP-03), the language (FR-JOB-012), and the currency when a salary range exists (FR-JOB-013).<br>3. The system sets the state to Published and records it in the audit log (FR-JOB-029, FR-JOB-015).<br>4. For a public posting, the system lists it in the feed, search results, and public company profile (FR-JOB-033).<br>5. The system notifies Interviewees who follow the company (FR-JOB-096, see UC-JOB-20). |
| Alternative Flows | 1. The Hiring Manager publishes before a pipeline is published; admitted applications wait as defined in FR-APP-23 (FR-JOB-031).<br>4a. For a private posting, the system generates the private link and hides the posting from the feed, search, and profile (FR-JOB-034, FR-JOB-035); the Hiring Manager continues with UC-JOB-06. |
| Exception Flows   | E1. A check fails: the posting stays in Draft and the system lists every failed check (FR-JOB-030).<br>E2. The company is banned: the system refuses and shows that the account is banned (FR-JOB-049). |
| Postconditions    | The posting is Published with its chosen visibility, or remains Draft with the failed checks listed. |
| Business Rules    | Only the transitions Draft to Published, Published to Archived, and Archived to Published exist (FR-JOB-029). |

### UC-JOB-06: Manage Private Access and Invitations

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-06 |
| Use Case Name     | Manage Private Access and Invitations |
| Actor(s)          | Hiring Manager, System, Interviewee (invited) |
| Description       | The Hiring Manager shares a private posting by private link or email invitation, tracks the invitations sent, and revokes the link when needed. |
| Preconditions     | The posting is Published with Private visibility. |
| Trigger           | The Hiring Manager opens the sharing controls of a private posting. |
| Main Flow         | 1. The system shows the private link (FR-JOB-035).<br>2. The Hiring Manager enters one or more email addresses and sends the invitation containing the private link (FR-JOB-036).<br>3. The system sends the emails and records each invitation with address, time sent, and delivery status (FR-JOB-037).<br>4. The Hiring Manager views the invitations list of the posting. |
| Alternative Flows | 1. The Hiring Manager copies the private link and shares it outside the platform.<br>2. The Hiring Manager revokes the link; the system generates a new one and the revoked link stops giving access (FR-JOB-035). |
| Exception Flows   | E1. An invitation email fails: the system marks the address as failed and keeps the link available to copy.<br>E2. An invalid email address is entered: the system rejects it and identifies it.<br>E3. A revoked, invalid, or unknown link is opened: the system shows "This posting is not available" without revealing whether the posting exists. |
| Postconditions    | The invitations are recorded and visible to the Hiring Manager, and invited Interviewees can reach the posting. |
| Business Rules    | A private posting is reached only through its private link or an email invitation (FR-JOB-034). The Hiring Manager sees the invitations sent for each private posting. |

### UC-JOB-07: Change Posting Visibility

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-07 |
| Use Case Name     | Change Posting Visibility |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager switches a Published posting between Public and Private. |
| Preconditions     | The posting is Published and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager changes the visibility setting. |
| Main Flow         | 1. The Hiring Manager selects the other visibility.<br>2. The system applies it; a posting made private leaves the feed, search, and company profile, and a posting made public enters them (FR-JOB-033, FR-JOB-034).<br>3. The system records the change in the audit log. |
| Alternative Flows | 1. When a posting becomes private, the system generates a private link (FR-JOB-035). |
| Exception Flows   | E1. The company is banned: the system refuses (FR-JOB-049). |
| Postconditions    | The posting has the new visibility and every existing application is unchanged. Saved entries of Interviewees who can no longer see it show "No longer available" (FR-JOB-089). |
| Business Rules    | Changing visibility never changes an existing application (FR-JOB-038). |

### UC-JOB-08: Archive Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-08 |
| Use Case Name     | Archive Job Posting |
| Actor(s)          | Hiring Manager, System, Interviewee (notified) |
| Description       | The Hiring Manager archives a Published posting, which stops new applications and suspends existing ones without deleting or deciding them. |
| Preconditions     | The posting is Published and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager chooses Archive. |
| Main Flow         | 1. The system shows the number of active applications and asks for confirmation (FR-JOB-044).<br>2. The Hiring Manager confirms.<br>3. The system sets the state to Archived, removes the posting from the feed, search, and company profile, and stops new applications (FR-JOB-039).<br>4. The system suspends all applications: no module can be started or continued, no live interview booked, no decision recorded (FR-JOB-040).<br>5. The system pauses the module completion windows and marks no module as missed (FR-JOB-041).<br>6. The system notifies Interviewees with an active application (FR-JOB-045). |
| Alternative Flows | 1. An Interviewee with an application opens the application tracking view and sees that the posting is archived and cannot be completed (FR-JOB-042).<br>2. The Interviewee withdraws the application (FR-JOB-043, FR-APP-19). |
| Exception Flows   | E1. An Interviewee tries to start a module or book an interview: the system refuses and shows that the posting is archived (FR-JOB-040).<br>E2. The Hiring Manager cancels the confirmation: the posting stays Published. |
| Postconditions    | The posting is Archived. Its applications, results, statuses, and pipeline data are retained. |
| Business Rules    | Archiving suspends applications but never deletes or decides them (FR-JOB-039, FR-JOB-040). An Archived posting is shown only to Interviewees who applied to or saved it (FR-JOB-086). |

### UC-JOB-09: Republish Archived Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-09 |
| Use Case Name     | Republish Archived Posting |
| Actor(s)          | Hiring Manager, System, Interviewee (notified) |
| Description       | The Hiring Manager returns an Archived posting to Published and the suspended applications resume. |
| Preconditions     | The posting is Archived, the Hiring Manager is permitted to manage it, and the company is not banned. |
| Trigger           | The Hiring Manager chooses Republish. |
| Main Flow         | 1. The Hiring Manager chooses Republish.<br>2. The system returns the posting to Published with its previous visibility (FR-JOB-046).<br>3. The system resumes the suspended applications and extends each module completion window by the time the posting was archived (FR-JOB-041).<br>4. The system notifies Interviewees with an active application (FR-JOB-045).<br>5. The Hiring Manager decides on applications (pass, reject, offer) through APP. |
| Alternative Flows | 3a. The application deadline has passed: the posting stays closed to new applications while existing applications resume (FR-JOB-047). |
| Exception Flows   | E1. The company is banned: the system refuses (FR-JOB-049).<br>E2. Republishing would leave mandatory data invalid: the system lists the failed checks as in UC-JOB-05. |
| Postconditions    | The posting is Published. New applications are accepted only if the deadline has not passed. |
| Business Rules    | A posting can be reopened to new applications only before its deadline passes (FR-JOB-047, FR-APP-04). Republishing does not decide any application (FR-JOB-046). |

### UC-JOB-10: Delete Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-10 |
| Use Case Name     | Delete Job Posting |
| Actor(s)          | Hiring Manager, System Admin, System |
| Description       | A posting is permanently removed. A Hiring Manager may delete only an unused Draft, while a System Admin may delete a posting in any state. |
| Preconditions     | The posting exists. The actor is permitted to manage it (Hiring Manager, FR-JOB-024) or is a System Admin (FR-AUTH-036). |
| Trigger           | The actor chooses Delete on a posting. |
| Main Flow         | 1. The system checks that the Hiring Manager's posting is a Draft that never received an application (FR-JOB-052).<br>2. The system asks for confirmation and states that deletion is permanent (FR-JOB-055).<br>3. The actor confirms.<br>4. The system deletes the posting permanently and records the deletion in the audit log (FR-JOB-015). |
| Alternative Flows | 1. The System Admin deletes a posting in any state, including one with applications; the confirmation shows the number of applications affected (FR-JOB-054, FR-JOB-055). |
| Exception Flows   | E1. The Hiring Manager requests deletion of a Published or Archived posting, or one with applications: the system refuses and offers archiving instead (FR-JOB-053).<br>E2. The actor cancels: nothing is deleted. |
| Postconditions    | The posting is permanently removed and the deletion is in the audit log. Saved entries and invitations to it show "No longer available" (FR-JOB-089, FR-JOB-093). |
| Business Rules    | A Hiring Manager can only delete a Draft that never received an application. A System Admin can delete a posting at any time (FR-JOB-052 to FR-JOB-054). |

### UC-JOB-11: Handle Postings of a Banned Company

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-11 |
| Use Case Name     | Handle Postings of a Banned Company |
| Actor(s)          | System, System Admin (bans and resolves through AUTH), Interviewee (notified) |
| Description       | When a company is banned, the system archives its Published postings and suspends their applications until the ban is resolved. |
| Preconditions     | A System Admin bans a company (FR-AUTH-037, UC-AUTH-16). |
| Trigger           | The ban is applied, or later resolved. |
| Main Flow         | 1. The system archives every Published posting of the company and suspends their applications (FR-JOB-048, FR-JOB-040, FR-JOB-041).<br>2. The system keeps Draft postings as Draft.<br>3. The system excludes the company's postings from the feed and search (FR-JOB-067).<br>4. The system blocks publishing or republishing while the ban lasts (FR-JOB-049).<br>5. The system notifies Interviewees with an active application, without stating the reason (FR-JOB-051). |
| Alternative Flows | 1. When the ban is resolved, the system returns each posting archived because of the ban to Published with its previous visibility and resumes its applications as in UC-JOB-09 (FR-JOB-050).<br>2. Postings that were Archived before the ban stay Archived (FR-JOB-050).<br>3. The system notifies Interviewees that their application resumed (FR-JOB-051). |
| Exception Flows   | E1. A banned company tries to edit or publish: the system refuses and shows that the account is banned (FR-JOB-049). |
| Postconditions    | During the ban no posting of the company is Published. After resolution the affected postings are Published again. |
| Business Rules    | While a company is banned, its postings stay archived and nothing can be published. The Hiring Manager decides on applications after the ban is resolved (FR-JOB-048 to FR-JOB-050). |

### UC-JOB-12: View Company Posting List

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-12 |
| Use Case Name     | View Company Posting List |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager sees all postings of the company in every state. |
| Preconditions     | The Hiring Manager is signed in with an account permitted to view postings. |
| Trigger           | The Hiring Manager opens the company's postings page. |
| Main Flow         | 1. The system lists all postings of the company with title, state, visibility, number of applications, and application deadline (FR-JOB-056).<br>2. The Hiring Manager opens a posting to edit, archive, republish, or delete it. |
| Alternative Flows | 1. The Hiring Manager filters the list by state and visibility and searches it by title (FR-JOB-057). |
| Exception Flows   | E1. The company has no postings: the system shows an empty list and offers to create a posting. |
| Postconditions    | The Hiring Manager has an up-to-date view of the company's postings. |
| Business Rules    | The list includes Draft, Published, and Archived postings (FR-JOB-056). |

### UC-JOB-13: Switch Employer Interface Language

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-13 |
| Use Case Name     | Switch Employer Interface Language |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager changes the language of the employer-side interface. |
| Preconditions     | The Hiring Manager is signed in. |
| Trigger           | The Hiring Manager opens the language selector. |
| Main Flow         | 1. The Hiring Manager selects one of the interface languages defined in `10-i18n.md` (FR-JOB-025).<br>2. The system applies the language immediately, without a new sign-in and without losing unsaved changes (FR-JOB-026).<br>3. The system keeps the choice for later sessions (FR-JOB-027). |
| Alternative Flows | 1. The Hiring Manager switches while editing a posting; the editor content stays as entered. |
| Exception Flows   | E1. The chosen language cannot be loaded: the system keeps the current language and shows an error message. |
| Postconditions    | The employer interface uses the chosen language. No posting content changed. |
| Business Rules    | Switching the interface language never changes posting content (FR-JOB-028). |

### UC-JOB-14: Browse Jobs Feed

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-14 |
| Use Case Name     | Browse Jobs Feed |
| Actor(s)          | Interviewee, System |
| Description       | An Interviewee browses public Published postings, optionally ordered by match with their job preferences. |
| Preconditions     | The user is signed in (FR-JOB-060). |
| Trigger           | The Interviewee opens the jobs feed. |
| Main Flow         | 1. The system lists public Published postings that are open or not yet open for applications, newest first, in pages (FR-JOB-061).<br>2. The system shows for each item the title, company, location, work mode, seniority, employment type, language, tags, deadline, and the salary range only where shown (FR-JOB-062).<br>3. The system marks items the Interviewee saved or applied to (FR-JOB-063).<br>4. The Interviewee scrolls or pages through the feed and opens an item (UC-JOB-17). |
| Alternative Flows | 1. A visitor who is not signed in is sent to sign-in or registration and the requested page opens afterwards (FR-JOB-060).<br>2. The Interviewee selects the Recommended ordering (UC-JOB-21).<br>3. The Interviewee limits the feed to followed companies (FR-JOB-066). |
| Exception Flows   | E1. No postings are available: the system shows an empty-feed message.<br>E2. Postings of deactivated or banned companies are never shown (FR-JOB-067). |
| Postconditions    | The Interviewee has seen the current public postings. |
| Business Rules    | Only public Published postings appear in the feed. Every user must be signed in to browse (FR-JOB-060). |

### UC-JOB-15: Filter and Sort Postings

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-15 |
| Use Case Name     | Filter and Sort Postings |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee narrows and orders the feed or search results. |
| Preconditions     | The user is signed in and the feed or a search result is displayed. |
| Trigger           | The Interviewee selects a filter or a sort option. |
| Main Flow         | 1. The Interviewee selects filters: tag, required skill, seniority, employment type, work mode, location, salary range, salary currency, posting language, company, posting date range (FR-JOB-068).<br>2. The system returns postings matching every selected filter; several values within one filter match any of them (FR-JOB-069).<br>3. The system displays the active filters (FR-JOB-073).<br>4. The Interviewee sorts by newest, nearest deadline, or relevance for search results (FR-JOB-072). |
| Alternative Flows | 1. The Interviewee removes one filter or all filters at once (FR-JOB-073).<br>2. The same filters and sorting apply to the feed and to search results (FR-JOB-075). |
| Exception Flows   | E1. No posting matches: the system shows a message and offers to clear the filters (FR-JOB-074).<br>E2. Salary minimum above maximum: the system identifies the invalid input and keeps the previous results.<br>E3. A salary filter is active: postings that do not show a salary range are excluded and the Interviewee is told (FR-JOB-070). |
| Postconditions    | The displayed postings match the selected filters and order. |
| Business Rules    | A salary filter applies only to postings in the selected currency and never converts amounts (FR-JOB-071). A hidden salary range is never revealed through filters (FR-JOB-070). |

### UC-JOB-16: Search Postings and Companies

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-16 |
| Use Case Name     | Search Postings and Companies |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee finds postings and companies by keyword. |
| Preconditions     | The user is signed in. |
| Trigger           | The Interviewee enters keywords in the search box, available on every page (FR-JOB-076). |
| Main Flow         | 1. The Interviewee enters at least 2 characters and submits (FR-JOB-079).<br>2. The system matches the keywords against posting title, description, tags, required skills, company name, and company names (FR-JOB-077).<br>3. The system returns results grouped by type (postings, companies), ranked by relevance, in pages (FR-JOB-078).<br>4. The Interviewee opens a result. |
| Alternative Flows | 1. The Interviewee applies filters and sorting to the posting results (UC-JOB-15).<br>2. Queries with common spelling variations still return results (FR-JOB-081, optional). |
| Exception Flows   | E1. The query is too short: the system says so and keeps the previous results (FR-JOB-079).<br>E2. No result: the system shows a no-results message.<br>E3. The search service is unavailable: the system shows "Search is temporarily unavailable" and the feed and filters keep working. |
| Postconditions    | The Interviewee sees only postings and companies they are allowed to see. |
| Business Rules    | Search returns only public Published postings and public company profiles, never Draft, Archived, or private postings of other companies (FR-JOB-080). |

### UC-JOB-17: View Posting Details

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-17 |
| Use Case Name     | View Posting Details |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee reads everything needed to decide whether to apply and sees which action is possible. |
| Preconditions     | The user is signed in. The posting is public Published, or the user holds its private link or an invitation, or the user applied to or saved it. |
| Trigger           | The Interviewee opens a posting from the feed, search, saved list, invitations, or a shared link. |
| Main Flow         | 1. The system displays the title, company (linked to its profile), description, location, work mode, seniority, employment type, language, tags, application dates, the Must-have and Nice-to-have groups, and the salary range only where shown (FR-JOB-082).<br>2. The system displays whether the posting is open, opens on a date, closed, or archived (FR-JOB-083).<br>3. The system offers Apply on an open posting (FR-JOB-084), continuing with UC-APP-02.<br>4. The system offers Save and Follow company (FR-JOB-085). |
| Alternative Flows | 1. The posting is closed: the system states whether the deadline passed or the application cap was reached (FR-JOB-083).<br>2. The posting is Archived: it is shown only to Interviewees who applied to or saved it, with a notice that it no longer accepts applications (FR-JOB-086). |
| Exception Flows   | E1. The Interviewee cannot apply: the system states what is missing, such as an incomplete resume or an existing application (FR-JOB-084, FR-APP-16).<br>E2. The link is revoked, invalid, or unknown: the system shows "This posting is not available" without revealing whether it exists. |
| Postconditions    | The Interviewee is informed and may apply, save, or follow. |
| Business Rules    | Salary is shown only where the posting shows it (FR-APP-02). |

### UC-JOB-18: Save and Unsave Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-18 |
| Use Case Name     | Save and Unsave Posting |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee keeps a personal list of postings to return to. |
| Preconditions     | The Interviewee is signed in. |
| Trigger           | The Interviewee chooses Save or Unsave on a posting. |
| Main Flow         | 1. The Interviewee saves a posting from the feed, search results, or detail page (FR-JOB-087).<br>2. The system adds it to the saved list.<br>3. The Interviewee opens the saved list; the system shows title, company, deadline, and whether each posting is open, closed, or archived (FR-JOB-088). |
| Alternative Flows | 1. The Interviewee unsaves a posting from any of those places.<br>2. The system reminds the Interviewee 48 hours before the deadline of a saved posting not yet applied to (FR-JOB-090, optional). |
| Exception Flows   | E1. A saved posting is no longer visible (made private or removed): the entry is labelled "No longer available" and shows none of its content (FR-JOB-089). |
| Postconditions    | The saved list reflects the Interviewee's choice. |
| Business Rules    | Saved postings belong to the Interviewee; Hiring Managers and System Admins never see who saved a posting. |

### UC-JOB-19: View Received Invitations

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-19 |
| Use Case Name     | View Received Invitations |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee sees invitations to private postings sent to the email address of their account. |
| Preconditions     | The Interviewee is signed in with a verified email address. |
| Trigger           | The Interviewee opens the invitations list. |
| Main Flow         | 1. The system matches invitations to the account by its verified email address (FR-JOB-092).<br>2. The system lists each invitation with posting title, company, date invited, and whether the posting is open or closed (FR-JOB-091).<br>3. The Interviewee opens a posting from the list (FR-JOB-093) and continues with UC-JOB-17. |
| Alternative Flows | 1. The invitation was sent before the account existed; it still appears once the email is verified (FR-JOB-092). |
| Exception Flows   | E1. The posting was deleted or archived, or its private link revoked: the invitation is labelled "No longer available" and shows none of the posting's content (FR-JOB-093).<br>E2. The Interviewee has no invitations: the system shows an empty list. |
| Postconditions    | The Interviewee sees only invitations sent to their own email address. |
| Business Rules    | Invitations are matched by verified email only (FR-JOB-092). |

### UC-JOB-20: Follow Company

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-20 |
| Use Case Name     | Follow Company |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee follows companies to see and be notified about their new postings. |
| Preconditions     | The Interviewee is signed in. |
| Trigger           | The Interviewee chooses Follow on a company profile or posting detail page. |
| Main Flow         | 1. The Interviewee follows a company (FR-JOB-094).<br>2. The system adds it to the followed list (FR-JOB-095).<br>3. When the company publishes a public posting, the system notifies the Interviewee (FR-JOB-096). |
| Alternative Flows | 1. The Interviewee unfollows a company.<br>2. The Interviewee turns off notifications for one followed company (FR-JOB-097, optional).<br>3. The Interviewee limits the feed to followed companies (FR-JOB-066). |
| Exception Flows   | E1. The company is deactivated or banned: its postings are not shown and no notification is sent (FR-JOB-067). |
| Postconditions    | The followed list is updated. |
| Business Rules    | Followed companies belong to the Interviewee; Hiring Managers and System Admins never see who follows a company. |

### UC-JOB-21: Get Recommendations

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-21 |
| Use Case Name     | Get Recommendations |
| Actor(s)          | Interviewee, System |
| Description       | The system suggests relevant postings based on job preferences and posting data. |
| Preconditions     | The Interviewee is signed in. For the Recommended ordering, job preferences are set (FR-AUTH-030). |
| Trigger           | The Interviewee selects the Recommended ordering, opens a posting, or submits an application. |
| Main Flow         | 1. The Interviewee selects Recommended ordering in the feed (FR-JOB-064).<br>2. The system ranks postings by match with the preferences.<br>3. The system shows for each recommended posting the preference or tag that caused it (FR-JOB-065). |
| Alternative Flows | 1. On a detail page, the system shows up to 5 similar postings based on shared tags, skills, and seniority (FR-JOB-098, optional).<br>2. After an application is submitted, the system shows similar postings (FR-JOB-099, optional). |
| Exception Flows   | E1. No preferences are set: the Recommended ordering is not offered.<br>E2. Postings the Interviewee has applied to are excluded (FR-JOB-100). |
| Postconditions    | The Interviewee sees ranked or similar postings. |
| Business Rules    | Recommendations use only structured posting data and job preferences, never name, photo, gender, age, or nationality (FR-JOB-101). |

### UC-JOB-22: Manage Saved Searches and New-Jobs Tracker

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-22 |
| Use Case Name     | Manage Saved Searches and New-Jobs Tracker |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee saves a search with its filters and is alerted when new postings match. |
| Preconditions     | The Interviewee is signed in. |
| Trigger           | The Interviewee chooses Save search. |
| Main Flow         | 1. The Interviewee names and saves the current search and filters (FR-JOB-102).<br>2. When a newly published posting matches the saved search, the system notifies the Interviewee (FR-JOB-103). |
| Alternative Flows | 1. The Interviewee lists and deletes saved searches (FR-JOB-102).<br>2. The Interviewee chooses immediate notifications or a daily digest (FR-JOB-104). |
| Exception Flows   | E1. A saved search has no name: the system asks for one. |
| Postconditions    | Saved searches and notification mode are stored. |
| Business Rules    | Saved searches belong to the Interviewee. All notifications go through the APP mechanism. This use case is optional (COULD). |

### UC-JOB-23: Share Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-23 |
| Use Case Name     | Share Posting |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee copies a link to a Published posting to share it outside the platform. |
| Preconditions     | The posting is Published and the user is signed in. |
| Trigger           | The Interviewee chooses Share on a posting. |
| Main Flow         | 1. The system provides a link to copy (FR-JOB-105).<br>2. The Interviewee copies and shares it.<br>3. A recipient opens the link, signs in, and the detail page opens (FR-JOB-060). |
| Alternative Flows | 1. For a private posting, the shared link is its private link (FR-JOB-035). |
| Exception Flows   | E1. The recipient is not signed in: the system sends them to sign-in or registration and opens the page afterwards.<br>E2. The link is revoked or the posting no longer available: the system shows "This posting is not available". |
| Postconditions    | The Interviewee has a shareable link. |
| Business Rules    | Sharing never bypasses sign-in (FR-JOB-060). |

### UC-JOB-24: Award Verified Crest

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-24 |
| Use Case Name     | Award Verified Crest |
| Actor(s)          | System, Interviewee |
| Description       | When an Interviewee successfully completes an assessment, the system awards a Verified Crest tied to the skill or competency the assessment verified. |
| Preconditions     | The Interviewee has completed an assessment module (APP 3.5) for which a Crest can be awarded. |
| Trigger           | An assessment module is completed. |
| Main Flow         | 1. The system determines whether the assessment was completed successfully, according to the criteria required for its Crest (FR-JOB-106).<br>2. If the criteria are satisfied, the system awards the Crest (FR-JOB-107).<br>3. The system associates the Crest with the verified skill or competency (FR-JOB-108).<br>4. The system associates the Crest with the Interviewee's profile (FR-JOB-109). |
| Alternative Flows | 1. The criteria are not satisfied: no Crest is awarded and the assessment result is unaffected. |
| Exception Flows   | E1. The Crest cannot be awarded or attached to the profile (service or data error): the system keeps the assessment result and does not mark the Crest as awarded until it is attached. |
| Postconditions    | The Interviewee has the Crest on their profile, or none if the criteria were not met. |
| Business Rules    | A Crest is awarded only for the successful completion of an assessment and is tied to the skill or competency it verified (FR-JOB-106 to FR-JOB-108). |

### UC-JOB-25: Apply Interviewee Anonymity

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-25 |
| Use Case Name     | Apply Interviewee Anonymity |
| Actor(s)          | System, Hiring Manager |
| Description       | While a stage is anonymous, the system hides identifying information about the Interviewee from the Hiring Manager and keeps qualification information visible. |
| Preconditions     | The Hiring Manager enabled anonymity for the pipeline and chose the reveal point (FR-APP-97, UC-APP-13). The Interviewee was admitted to the pipeline. |
| Trigger           | The Hiring Manager opens an application. |
| Main Flow         | 1. The system determines whether anonymity applies to the application's current stage (FR-JOB-110).<br>2. The system hides the information designated as non-disclosable (FR-JOB-111).<br>3. The system continues to display scores, module results, skills, and Crests (FR-JOB-112).<br>4. The system applies the hiding in every view, including resume copies, the dashboard, scorecards, and notes (FR-JOB-113). |
| Alternative Flows | 1. Anonymity does not apply to the stage: the system displays the full application. |
| Exception Flows   | E1. Access to hidden information is attempted before the reveal point: the system denies it and logs the attempt (FR-APP-98). |
| Postconditions    | Identifying information stays hidden and qualification information stays visible during anonymous stages. |
| Business Rules    | While anonymity applies, qualification information stays visible and identifying information stays hidden in every view (FR-JOB-112, FR-JOB-113). Optional (COULD). |

### UC-JOB-26: Reveal Interviewee Identity

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-26 |
| Use Case Name     | Reveal Interviewee Identity |
| Actor(s)          | System, Hiring Manager |
| Description       | When an Interviewee reaches the configured reveal point, the system makes the hidden identity available to authorized users and keeps the application history. |
| Preconditions     | Anonymity applies to the pipeline (UC-JOB-25) and an identity-reveal point is configured (FR-APP-97). |
| Trigger           | The Interviewee reaches the identity-disclosure stage. |
| Main Flow         | 1. The system determines that the Interviewee reached the configured stage (FR-JOB-114).<br>2. The system makes the identity available to authorized users (FR-JOB-115).<br>3. The system retains the Interviewee's application and assessment history (FR-JOB-116). |
| Alternative Flows | 1. Identity is revealed at the latest when an offer is created (FR-APP-101). |
| Exception Flows   | E1. The reveal cannot be applied (service error): the system keeps the identity hidden and does not create an offer until the reveal succeeds. |
| Postconditions    | The identity is visible to authorized users and the history is intact. |
| Business Rules    | Identity disclosure never removes application or assessment history (FR-JOB-116). Reveal is permanent (FR-APP-100). Optional (COULD). |


## 4.6 ANA Use Cases
