## 4. Use cases

Owner: TBD | Status: Draft | Last updated: 2026-10-13 | Jira: HIRE-n
---

## 4.1 Actors

// This is only an example.
- Candidate
- Visitor (a person who is not signed in)
- Hiring Manager (the employer)
- Admin
- System

## 4.2 AUTH Use Cases
## 4.5 JOB Use Cases
### Use Case Index

| Use Case | Actor(s) | Related Requirements |
| --- | --- | --- |
| UC-JOB-01: Create, Configure and Edit Job Posting | Hiring Manager, Admin, System, Candidate (notified) | FR-JOB-001 to FR-JOB-004, FR-JOB-005 to FR-JOB-013, FR-JOB-014, FR-JOB-015, FR-JOB-018, FR-JOB-045, FR-JOB-046, FR-JOB-089 to FR-JOB-098 |
| UC-JOB-02: Preview Posting | Hiring Manager, System | FR-JOB-016, FR-JOB-017 |
| UC-JOB-03: Publish Job Posting | Hiring Manager, System, Candidate (notified) | FR-JOB-011, FR-JOB-012, FR-JOB-023 to FR-JOB-026, FR-JOB-035, FR-JOB-069, FR-JOB-089, FR-JOB-091, FR-JOB-093, FR-JOB-102, FR-JOB-103 |
| UC-JOB-04: Archive and Republish Posting | Hiring Manager, System, Candidate | FR-JOB-027 to FR-JOB-035, FR-JOB-039 |
| UC-JOB-05: Delete Posting and Manage Postings as Admin | Hiring Manager, Admin, System | FR-JOB-018, FR-JOB-040, FR-JOB-041, FR-JOB-042 |
| UC-JOB-06: Handle Postings of a Banned Company | System, Admin, Hiring Manager, Candidate (notified) | FR-JOB-036 to FR-JOB-039 |
| UC-JOB-07: Manage Company Posting List | Hiring Manager, System | FR-JOB-043, FR-JOB-044 |
| UC-JOB-08: Switch Employer Interface Language (Optional) | Hiring Manager, System | FR-JOB-019 to FR-JOB-022 |
| UC-JOB-09: Browse, Filter, Sort, and Search Jobs and Companies | Candidate, Visitor, System | FR-JOB-047 to FR-JOB-050, FR-JOB-053 to FR-JOB-067 |
| UC-JOB-10: View Posting Details and Company Profile | Candidate, System | FR-JOB-026, FR-JOB-068 to FR-JOB-072 |
| UC-JOB-11: Share and Open Posting Link | Candidate, Visitor, System | FR-JOB-047, FR-JOB-088 |
| UC-JOB-12: Manage Saved Postings | Candidate, System | FR-JOB-073 to FR-JOB-076 |
| UC-JOB-13: Follow and Manage Companies | Candidate, System | FR-JOB-077 to FR-JOB-080 |
| UC-JOB-14: Get Recommendations | Candidate, System | FR-JOB-051, FR-JOB-052, FR-JOB-081 to FR-JOB-084 |
| UC-JOB-15: Manage Saved Searches and New-Jobs Tracker | Candidate, System | FR-JOB-085 to FR-JOB-087 |
| UC-JOB-16: Apply to Posting | Candidate, System | FR-JOB-099 to FR-JOB-109, FR-JOB-111 |
| UC-JOB-17: Screen and Admit Applications | System, Hiring Manager, Candidate (informed) | FR-JOB-025, FR-JOB-105, FR-JOB-107 to FR-JOB-111 |

### UC-JOB-01: Create, Configure and Edit Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-01 |
| Use Case Name     | Create, Configure and Edit Job Posting |
| Actor(s)          | Hiring Manager, Admin, System, Candidate (notified) |
| Description       | A Hiring Manager creates a posting as a Draft, completes its details, language, currency, tags, application window, criteria, and application questions, and edits it over time. |
| Preconditions     | The company is registered and verified. The Hiring Manager is signed in with an account permitted to manage postings and the company is not banned. |
| Trigger           | The Hiring Manager chooses to create a posting or to edit an existing one. |
| Main Flow         | 1. The Hiring Manager opens the posting editor.<br>2. The Hiring Manager enters the job title, the job description, and the number of open positions (an integer of at least 1), and optionally the seniority level, employment type, work mode, location, and salary range, and whether the salary range is shown to candidates.<br>3. The Hiring Manager selects the posting language from the standard list of languages and, if there is a salary range, its currency from the standard list of currencies.<br>4. The Hiring Manager adds or removes tags chosen from the Admin-managed taxonomy, up to 10.<br>5. The Hiring Manager sets the application opening date and the application deadline (later than the opening date), and optionally a maximum number of applications.<br>6. The Hiring Manager adds criteria and classifies each as exactly one of Must-have or Nice-to-have. Each Must-have is a condition on a structured resume field (minimum years of experience, minimum education level, minimum language proficiency level, or location eligibility).<br>7. The Hiring Manager adds application questions, each marked required or optional.<br>8. The system autosaves the Draft at least every 30 seconds while there are unsaved changes and shows whether the latest changes are saved, with the time of the last save.<br>9. The Hiring Manager saves the posting; the system stores it in the Draft state. |
| Alternative Flows | 1. Incomplete Draft: the Hiring Manager saves a Draft with incomplete fields; mandatory content is checked at publication (UC-JOB-03).<br>2. Any language or currency: the system does not restrict either and shows text and amounts as written, without translation or conversion.<br>3. Duplicate (optional): the Hiring Manager duplicates a posting into a new Draft that copies content and criteria but not applications or dates.<br>4. Extend deadline: before the deadline passes, the Hiring Manager selects a later deadline and the system saves it.<br>5. Published posting: autosave does not apply and changes take effect only when the Hiring Manager saves them.<br>6. Candidate notification: after a Published posting is saved with a changed title, salary range, location, or work mode, the system emails candidates with an active application. |
| Exception Flows   | E1. The Hiring Manager leaves the editor with unsaved changes: the system offers to save, discard, or stay.<br>E2. The Hiring Manager leaves a new posting never saved as a Draft: the system offers to save it as a Draft, discard it, or stay.<br>E3. Autosave fails: the system states that the latest changes are not saved, keeps the entered data on screen, and retries.<br>E4. The number of open positions is not an integer of at least 1, or the deadline is not later than the opening date: the system rejects it and identifies the field.<br>E5. A Must-have is entered as free text: the system rejects it and asks for a condition on a structured resume field.<br>E6. The posting has received its first application and the Hiring Manager changes the Must-have criteria: the system refuses and states that they are locked.<br>E7. The deadline has already passed and the Hiring Manager extends it: the system refuses and states the reason.<br>E8. The Hiring Manager adds an 11th tag: the system rejects it. A tag removed from the taxonomy by an Admin stays on the posting but cannot be added to others.<br>E9. Two accounts of the same company edit the same posting: the system saves the first change, rejects the second, and shows the current posting data.<br>E10. The company is banned: the system refuses the edit and shows that the account is banned.<br>E11. The account is not permitted to manage postings: the system denies the action.<br>E12. A notification email fails: the system retries delivery and shows the notification inside the platform. |
| Postconditions    | The posting is stored (Draft, or updated Published). |
| Business Rules    | Autosave applies to Drafts only. Every criterion is exactly one of Must-have or Nice-to-have; Must-have criteria are checkable conditions on structured resume fields and are locked after the first application. A deadline can be extended only before it passes. A posting has at most 10 tags, taken only from the Admin taxonomy. A salary range that is not shown is never revealed to candidates. Only Hiring Managers permitted to manage postings, and Admins, can create or edit a posting. |
| Related Requirements | FR-JOB-001 to FR-JOB-004, FR-JOB-005 to FR-JOB-013, FR-JOB-014, FR-JOB-015, FR-JOB-018, FR-JOB-045, FR-JOB-046, FR-JOB-089 to FR-JOB-098 |

### UC-JOB-02: Preview Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-02 |
| Use Case Name     | Preview Posting |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager sees a posting as candidates will see it. |
| Preconditions     | A posting exists and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager chooses Preview in the editor. |
| Main Flow         | 1. The Hiring Manager chooses Preview, including with unsaved edits.<br>2. The system shows the posting as it appears in the jobs feed and on the detail page.<br>3. The Hiring Manager closes the preview. |
| Alternative Flows | 1. A Draft posting is previewed before it has ever been published. |
| Exception Flows   | E1. The posting cannot be fully rendered because required fields are missing: the system shows the preview with the missing fields marked. |
| Postconditions    | Previewing changed nothing and notified nobody. |
| Business Rules    | Previewing never changes data or notifies anyone. |
| Related Requirements | FR-JOB-016, FR-JOB-017 |

### UC-JOB-03: Publish Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-03 |
| Use Case Name     | Publish Job Posting |
| Actor(s)          | Hiring Manager, System, Candidate (notified) |
| Description       | The Hiring Manager moves a Draft posting to Published after the system verifies it is complete. The posting then opens and closes for applications by its dates and cap. |
| Preconditions     | The posting is in Draft. The Hiring Manager is permitted to manage it and the company is not banned. |
| Trigger           | The Hiring Manager chooses Publish. |
| Main Flow         | 1. The Hiring Manager chooses Publish.<br>2. The system checks the title, description, number of open positions, opening date and deadline, language, and the currency when a salary range exists.<br>3. The system sets the state to Published.<br>4. The system lists the posting in the feed, search results, and public company profile.<br>5. The system emails candidates who follow the company within 10 minutes.<br>6. Before the opening date the posting shows as opening on a given date and rejects submissions; from the opening date it accepts applications; when the deadline passes or the cap is reached it stops accepting and shows as closed, with the reason. |
| Alternative Flows | 1. The Hiring Manager publishes before a pipeline is published; admitted applications wait (UC-JOB-17).<br>2. A posting republished after its deadline stays closed to new applications (UC-JOB-04). |
| Exception Flows   | E1. A check fails: the posting stays in Draft and the system lists every failed check.<br>E2. The company is banned: the system refuses and shows that the account is banned.<br>E3. Two submissions compete for the last place under the cap: the system accepts the first and rejects the second.<br>E4. The company is deactivated or banned, or the follower muted the company: no follower notification email is sent. |
| Postconditions    | The posting is Published, or remains Draft with the failed checks listed. |
| Business Rules    | Only the transitions Draft to Published, Published to Archived, and Archived to Published exist. A posting can be published before a pipeline is published. Only accepted submissions count toward the cap. |
| Related Requirements | FR-JOB-011, FR-JOB-012, FR-JOB-023 to FR-JOB-026, FR-JOB-035, FR-JOB-069, FR-JOB-089, FR-JOB-091, FR-JOB-093, FR-JOB-102, FR-JOB-103 |

### UC-JOB-04: Archive and Republish Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-04 |
| Use Case Name     | Archive and Republish Posting |
| Actor(s)          | Hiring Manager, System, Candidate |
| Description       | The Hiring Manager archives a Published posting, which stops new applications and suspends existing ones without deleting or deciding them, and later republishes it. |
| Preconditions     | The posting is Published (to archive) or Archived (to republish) and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager chooses Archive or Republish. |
| Main Flow         | 1. The system shows the number of active applications and asks for confirmation.<br>2. The Hiring Manager confirms.<br>3. The system sets the state to Archived, removes the posting from the feed, search, and company profile, and stops new applications.<br>4. The system suspends all applications: no module can be started or continued, no live interview booked, no decision recorded.<br>5. The system pauses the module completion windows and marks no module as missed.<br>6. The system emails candidates with an active application.<br>7. Later the Hiring Manager chooses Republish and the system returns the posting to Published.<br>8. The system resumes the suspended applications, extends each module completion window by the time the posting was archived, and emails candidates with an active application.<br>9. The Hiring Manager decides on the applications (pass, reject, offer). |
| Alternative Flows | 1. Candidate view: a candidate opens the application tracking view and sees that the posting is archived and cannot be completed, with the submitted data, results, and status unchanged.<br>2. Withdraw: the candidate withdraws the application.<br>3. The application deadline has passed at republication: the posting stays closed to new applications while existing applications resume. |
| Exception Flows   | E1. A candidate tries to start a module or book an interview: the system refuses and shows that the posting is archived.<br>E2. The Hiring Manager cancels the confirmation: the posting stays Published.<br>E3. The company is banned: the system refuses to republish.<br>E4. Republishing would leave mandatory data invalid: the system lists the failed checks, as when publishing.<br>E5. A notification email fails: the system retries delivery and shows the notification inside the platform. |
| Postconditions    | The posting is Archived with its applications, results, statuses, and pipeline data retained, or Published again with applications resumed. |
| Business Rules    | Archiving suspends applications but never deletes or decides them. An Archived posting is shown only to candidates who applied to or saved it. A posting can be reopened to new applications only before its deadline passes. A withdrawn application blocks a new one to the same posting. |
| Related Requirements | FR-JOB-027 to FR-JOB-035, FR-JOB-039 |

### UC-JOB-05: Delete Posting and Manage Postings as Admin

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-05 |
| Use Case Name     | Delete Posting and Manage Postings as Admin |
| Actor(s)          | Hiring Manager, Admin, System |
| Description       | A Hiring Manager or Admin permanently removes a posting. A Hiring Manager may delete any posting of their company; an Admin may delete postings of any company and may create, edit, archive, or republish them. |
| Preconditions     | The posting exists. The actor is permitted to manage it (Hiring Manager) or is an Admin. |
| Trigger           | The actor chooses Delete on a posting, or an Admin opens a posting of a company. |
| Main Flow         | 1. The system asks for confirmation and states that deletion is permanent; if applications exist, it shows how many will be affected.<br>2. The actor confirms.<br>3. The system permanently deletes the posting. |
| Alternative Flows | 1. An Admin creates, edits, archives, or republishes a posting of any company. |
| Exception Flows   | E1. The actor cancels: nothing is deleted. |
| Postconditions    | The posting is permanently removed. Saved entries of it show "No longer available". |
| Business Rules    | A Hiring Manager permitted to manage postings may delete any posting of their company. An Admin may delete any posting. |
| Related Requirements | FR-JOB-018, FR-JOB-040, FR-JOB-041, FR-JOB-042 |

### UC-JOB-06: Handle Postings of a Banned Company

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-06 |
| Use Case Name     | Handle Postings of a Banned Company |
| Actor(s)          | System, Admin, Hiring Manager, Candidate (notified) |
| Description       | When a company is banned, the system archives its Published postings and suspends their applications; when the ban is resolved, it returns them to Published. |
| Preconditions     | An Admin bans a company, or later resolves the ban. |
| Trigger           | The ban is applied or resolved. |
| Main Flow         | 1. The system archives every Published posting of the company and suspends their applications.<br>2. The system keeps Draft postings as Draft.<br>3. The system excludes the company's postings from the feed and search.<br>4. The system blocks publishing or republishing while the ban lasts.<br>5. The system emails candidates with an active application, without stating the reason. |
| Alternative Flows | 1. Resolution: the system returns each posting archived because of the ban to Published, resumes its applications and extends module completion windows by the suspension time, and emails candidates that their application resumed.<br>2. Postings that were Archived before the ban stay Archived.<br>3. A posting whose deadline has passed stays closed to new applications. |
| Exception Flows   | E1. A banned company tries to edit or publish: the system refuses and shows that the account is banned.<br>E2. Resuming a posting would leave mandatory data invalid: the system lists the failed checks. |
| Postconditions    | During the ban no posting of the company is Published. After resolution the affected postings are Published again. |
| Business Rules    | While a company is banned, its postings stay archived and nothing can be published. The reason for a ban is never disclosed to candidates. The Hiring Manager decides on applications after the ban is resolved. |
| Related Requirements | FR-JOB-036 to FR-JOB-039 |

### UC-JOB-07: Manage Company Posting List

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-07 |
| Use Case Name     | Manage Company Posting List |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager sees all postings of the company in every state and narrows the list. |
| Preconditions     | The Hiring Manager is signed in with an account permitted to view postings. |
| Trigger           | The Hiring Manager opens the company's postings page. |
| Main Flow         | 1. The system lists all postings of the company with title, state, number of applications, and application deadline.<br>2. The Hiring Manager opens a posting to edit, archive, republish, or delete it. |
| Alternative Flows | 1. The Hiring Manager filters the list by state and searches it by title. |
| Exception Flows   | E1. The company has no postings: the system shows an empty list and offers to create a posting.<br>E2. No posting matches the filter or search: the system shows an empty list. |
| Postconditions    | The Hiring Manager has an up-to-date view of the company's postings. |
| Business Rules    | The list includes Draft, Published, and Archived postings. Filtering and searching the list are recommended (SHOULD). |
| Related Requirements | FR-JOB-043, FR-JOB-044 |

### UC-JOB-08: Switch Employer Interface Language

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-08 |
| Use Case Name     | Switch Employer Interface Language (Optional) |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager optionally changes the language of the employer-side interface. |
| Preconditions     | The Hiring Manager is signed in. |
| Trigger           | The Hiring Manager opens the language selector. |
| Main Flow         | 1. The Hiring Manager selects one of the interface languages supported by the platform.<br>2. The system applies the language immediately, without a new sign-in and without losing unsaved changes.<br>3. The system may keep the choice for later sessions. |
| Alternative Flows | 1. The Hiring Manager switches while editing a posting; the editor content stays as entered. |
| Exception Flows   | E1. The chosen language cannot be loaded: the system keeps the current language and shows an error message. |
| Postconditions    | The employer interface uses the chosen language. No posting content changed. |
| Business Rules    | Switching the interface language never changes posting content. |
| Related Requirements | FR-JOB-019 to FR-JOB-022 |

### UC-JOB-09: Browse, Filter, Sort, and Search Jobs and Companies

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-09 |
| Use Case Name     | Browse, Filter, Sort, and Search Jobs and Companies |
| Actor(s)          | Candidate, Visitor, System |
| Description       | A candidate browses, searches, filters, and sorts Published postings and company profiles, optionally ordering postings by match with job preferences or limiting the feed to followed companies. A visitor who is not signed in is first sent to sign-in. |
| Preconditions     | The candidate is signed in, or a Visitor requests a discovery page and is redirected to sign-in or registration. |
| Trigger           | The candidate opens the feed or enters keywords in the search box, available on every page. |
| Main Flow         | 1. The system lists Published postings that are open or not yet open for applications, newest first, in pages of 20.<br>2. The system shows for each item the title, company, location, work mode, seniority level, employment type, posting language, tags, application deadline, and the salary range only where the Hiring Manager chose to show it.<br>3. The system marks items the candidate saved or applied to.<br>4. The candidate enters at least 2 characters in the search box and submits; the system matches the keywords against posting titles, descriptions, tags, required skills, and associated company names, as well as company names.<br>5. The system returns posting and company results grouped by type, ranked by relevance, and in pages.<br>6. The candidate optionally filters by required skill, seniority, employment type, work mode, location, salary range and currency, posting language, company, or posting date range. The system matches every selected filter and any selected value within a filter, and displays active filters.<br>7. The candidate optionally sorts results by newest, nearest deadline, relevance (for search results), or Recommended (UC-JOB-14). The candidate may remove individual filters or clear all filters.<br>8. The candidate scrolls or pages through results and opens an item. |
| Alternative Flows | 1. Sign-in: a Visitor is sent to sign-in or registration and the requested page opens afterwards.<br>2. The candidate limits the feed to followed companies, and removes the limit to see all postings.<br>3. Search queries with common spelling variations still return results (recommended). |
| Exception Flows   | E1. No postings are available, or the candidate follows no company with Published postings: the system shows an empty-feed message.<br>E2. Postings of deactivated or banned companies are never shown.<br>E3. Sign-in fails: the system keeps the Visitor on the sign-in page and shows the error.<br>E4. The query is too short: the system says so and keeps the previous results.<br>E5. No result or no posting matches: the system shows a no-results message and offers to clear filters where relevant.<br>E6. Salary minimum is above maximum: the system identifies the invalid input and keeps the previous results. Postings without a visible salary range are excluded from salary-filter results, without revealing the hidden range.<br>E7. The search service is unavailable: the system shows "Search is temporarily unavailable" and the feed and filters keep working. |
| Postconditions    | The candidate sees Published postings and company profiles allowed to them, matching the active search, filters, and ordering. |
| Business Rules    | Only Published postings appear in the feed and search. Every user must be signed in to browse. The system never converts amounts between currencies or reveals a hidden salary range through filters or search. |
| Related Requirements | FR-JOB-047 to FR-JOB-050, FR-JOB-053 to FR-JOB-067 |

### UC-JOB-10: View Posting Details and Company Profile

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-10 |
| Use Case Name     | View Posting Details and Company Profile |
| Actor(s)          | Candidate, System |
| Description       | The candidate reads everything needed to decide whether to apply, sees which action is possible, and sees the Published postings of a company on its public profile. |
| Preconditions     | The user is signed in. The posting is Published, or the user applied to or saved it. |
| Trigger           | The candidate opens a posting from the feed, search, saved list, or a shared link, or opens a company profile. |
| Main Flow         | 1. The system displays the title, company (linked to its profile), description, number of open positions, location, work mode, seniority level, employment type, posting language, tags, application opening date and deadline, the Must-have and Nice-to-have criteria in two separately labelled groups, and the salary range only where the Hiring Manager chose to show it.<br>2. The system displays whether the posting is open, opens on a date, closed, or archived.<br>3. The system offers Apply on an open posting (UC-JOB-16).<br>4. The system offers Save (UC-JOB-12) and Follow company (UC-JOB-13). |
| Alternative Flows | 1. Closed posting: the system states whether the deadline passed or the application cap was reached.<br>2. Archived posting: it is shown only to candidates who applied to or saved it, with a notice that it no longer accepts applications.<br>3. Company profile: the candidate opens a company profile and the system lists its Published postings; the candidate opens one. |
| Exception Flows   | E1. The candidate cannot apply: the system states what is missing, such as an incomplete resume or an existing application.<br>E2. The posting was deleted or the link is invalid: the system shows "This posting is not available" without revealing whether it exists.<br>E3. The company has no Published posting, or is deactivated or banned: the system shows an empty list. |
| Postconditions    | The candidate is informed and may apply, save, or follow. |
| Business Rules    | The salary range is shown only where the Hiring Manager chose to show it. Draft and Archived postings never appear on the company profile. |
| Related Requirements | FR-JOB-026, FR-JOB-068 to FR-JOB-072 |

### UC-JOB-11: Share and Open Posting Link

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-11 |
| Use Case Name     | Share and Open Posting Link |
| Actor(s)          | Candidate, Visitor, System |
| Description       | The candidate copies a link to a Published posting to share it outside the platform, and a recipient opens it. |
| Preconditions     | The posting is Published and the sharing user is signed in. |
| Trigger           | The candidate chooses Share on a posting, or a recipient opens a shared link. |
| Main Flow         | 1. The system provides a link to copy.<br>2. The candidate copies and shares it.<br>3. A recipient opens the link; if not signed in, the system sends them to sign-in or registration.<br>4. After sign-in the system opens the detail page of the posting. |
| Alternative Flows | 1. The recipient is already signed in: the detail page opens at once. |
| Exception Flows   | E1. The posting was deleted or is unavailable: the system shows "This posting is not available" without revealing whether it exists. |
| Postconditions    | The candidate has a shareable link and the recipient sees the posting or the unavailable message. |
| Business Rules    | Sharing never bypasses sign-in. |
| Related Requirements | FR-JOB-047, FR-JOB-088 |

### UC-JOB-12: Manage Saved Postings

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-12 |
| Use Case Name     | Manage Saved Postings |
| Actor(s)          | Candidate, System |
| Description       | The candidate keeps a personal list of postings to return to. |
| Preconditions     | The candidate is signed in. |
| Trigger           | The candidate chooses Save or Unsave on a posting, or opens the saved list. |
| Main Flow         | 1. The candidate saves a posting from the feed, search results, or detail page.<br>2. The system adds it to the saved list.<br>3. The candidate opens the saved list; the system shows title, company, deadline, and whether each posting is open, closed, or archived. |
| Alternative Flows | 1. The candidate unsaves a posting from any of those places or from the list.<br>2. Reminder (optional): the system reminds the candidate 48 hours before the deadline of a saved posting not yet applied to; no reminder is sent once the candidate has applied. |
| Exception Flows   | E1. A saved posting is no longer visible to the candidate (for example, it was deleted): the entry is labelled "No longer available" and shows none of its content. |
| Postconditions    | The saved list reflects the candidate's choice. |
| Business Rules    | Saved postings belong to the candidate; Hiring Managers and Admins never see who saved a posting. |
| Related Requirements | FR-JOB-073 to FR-JOB-076 |

### UC-JOB-13: Follow and Manage Companies

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-13 |
| Use Case Name     | Follow and Manage Companies |
| Actor(s)          | Candidate, System |
| Description       | The candidate follows companies, manages the followed list, and is notified when a followed company publishes a posting. |
| Preconditions     | The candidate is signed in. |
| Trigger           | The candidate chooses Follow or Unfollow on a company profile or posting detail page, or opens the followed list. |
| Main Flow         | 1. The candidate follows a company.<br>2. The system adds it to the followed list.<br>3. When the company publishes a posting, the system emails the candidate within 10 minutes.<br>4. The candidate opens the followed list and opens a company. |
| Alternative Flows | 1. The candidate unfollows a company from its profile, a detail page, or the list; the system stops its notification emails.<br>2. Mute (optional): the candidate turns off notification emails for one followed company and keeps following it, and can turn them back on. |
| Exception Flows   | E1. The company is deactivated or banned: its postings are not shown and no notification email is sent.<br>E2. The candidate follows no company: the system shows an empty list.<br>E3. A notification email fails: the system retries delivery and shows the notification inside the platform. |
| Postconditions    | The followed list and email notification settings are updated. |
| Business Rules    | Followed companies belong to the candidate; Hiring Managers and Admins never see who follows a company. Showing the followed list and email notifications are recommended (SHOULD). |
| Related Requirements | FR-JOB-077 to FR-JOB-080 |

### UC-JOB-14: Get Recommendations

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-14 |
| Use Case Name     | Get Recommendations |
| Actor(s)          | Candidate, System |
| Description       | The system suggests relevant postings based on job preferences and posting data, in the feed and on posting pages. |
| Preconditions     | The candidate is signed in. For the Recommended ordering, job preferences are set. |
| Trigger           | The candidate selects the Recommended ordering, opens a posting, or submits an application. |
| Main Flow         | 1. The candidate selects the Recommended ordering in the feed.<br>2. The system ranks postings by match with the preferences.<br>3. The system shows for each recommended posting the preference or tag that caused it. |
| Alternative Flows | 1. Similar postings (optional): on a detail page, and after an application is submitted, the system shows up to 5 postings with shared tags, required skills, and seniority level. |
| Exception Flows   | E1. No preferences are set: the Recommended ordering is not offered.<br>E2. Postings the candidate has already applied to are excluded.<br>E3. No similar posting exists: the system shows none. |
| Postconditions    | The candidate sees ranked or similar postings. |
| Business Rules    | Recommendations use only structured posting data and job preferences, never name, photo, gender, age, or nationality. |
| Related Requirements | FR-JOB-051, FR-JOB-052, FR-JOB-081 to FR-JOB-084 |

### UC-JOB-15: Manage Saved Searches and New-Jobs Tracker

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-15 |
| Use Case Name     | Manage Saved Searches and New-Jobs Tracker |
| Actor(s)          | Candidate, System |
| Description       | The candidate saves a search with its filters and is alerted when new postings match. |
| Preconditions     | The candidate is signed in. |
| Trigger           | The candidate chooses Save search. |
| Main Flow         | 1. The candidate names and saves the current search and filters.<br>2. When a newly published posting matches a saved search, the system emails the candidate. |
| Alternative Flows | 1. The candidate lists and deletes saved searches.<br>2. The candidate chooses immediate email notifications or a daily email digest. |
| Exception Flows   | E1. A saved search has no name: the system asks for one.<br>E2. The posting is from a banned or deactivated company: no email notification is sent. |
| Postconditions    | Saved searches and email notification mode are stored. |
| Business Rules    | Saved searches belong to the candidate. This use case is optional (COULD). |
| Related Requirements | FR-JOB-085 to FR-JOB-087 |

### UC-JOB-16: Apply to Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-16 |
| Use Case Name     | Apply to Posting |
| Actor(s)          | Candidate, System |
| Description       | The candidate checks their resume against the Must-have criteria, answers the application questions, and submits. The system validates, stores, and confirms the submission, then hands it to screening. |
| Preconditions     | The candidate is signed in and has a complete resume. The posting is Published and open for applications. |
| Trigger           | The candidate chooses Apply on an open posting. |
| Main Flow         | 1. The system shows which Must-have criteria the candidate's resume does not currently satisfy, if any (recommended).<br>2. The candidate answers the application questions, including all required ones, and submits.<br>3. The system checks that the candidate has no active, rejected, or withdrawn application for the posting, that the opening date has passed, that the deadline has not passed, and that the cap has not been reached.<br>4. The system stores an unchangeable copy of the resume and answers as submitted and counts the application toward the cap.<br>5. The system shows the candidate a confirmation with an application reference and emails the same confirmation within 10 minutes.<br>6. The system passes the application to screening (UC-JOB-17). |
| Alternative Flows | 1. After submitting, the system shows similar postings (UC-JOB-14, optional). |
| Exception Flows   | E1. A required question is unanswered: the system rejects the submission and lists the missing answers.<br>E2. The candidate already has an active, rejected, or withdrawn application: the system rejects the submission and states that an application already exists.<br>E3. The posting is not yet open, the deadline has passed, the cap was reached, or the posting is Archived (including while the candidate fills the form): the system rejects it, states which reason applies, and keeps nothing partial.<br>E4. Two candidates submit for the last place under the cap: the system accepts the first and rejects the second with the cap reason.<br>E5. A resume field needed by a Must-have criterion is empty: the system lists the criterion as not satisfied.<br>E6. The confirmation email fails: the system retries delivery and shows the confirmation inside the platform. |
| Postconditions    | An accepted application is stored unchanged and confirmed. A rejected submission leaves no data. |
| Business Rules    | A candidate has at most one application per posting; a rejected or withdrawn application also blocks a new one. The submitted resume and answers never change after submission. Only accepted submissions count toward the cap. The Must-have warning does not block submission. |
| Related Requirements | FR-JOB-099 to FR-JOB-109, FR-JOB-111 |

### UC-JOB-17: Screen and Admit Applications

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-17 |
| Use Case Name     | Screen and Admit Applications |
| Actor(s)          | System, Hiring Manager, Candidate (informed) |
| Description       | The system checks each accepted application against the Must-have criteria, admits or screens it out, and holds admitted applications until a pipeline is published. The Hiring Manager can review screened-out applications. |
| Preconditions     | A submission was accepted (UC-JOB-16). For reviewing, the Hiring Manager is permitted to manage the posting. |
| Trigger           | The system accepts a submission, a pipeline is published, or the Hiring Manager opens the screened-out applications of a posting. |
| Main Flow         | 1. The system evaluates every Must-have criterion using the structured resume fields as submitted.<br>2. All criteria are satisfied: the system admits the application to the pipeline.<br>3. If a pipeline is published, the system starts the candidate at every entry stage; otherwise the application keeps the status "Awaiting Pipeline" until a pipeline is published, and the system then starts the candidate at every entry stage without further action. |
| Alternative Flows | 1. Screened out: at least one criterion is not satisfied; the system does not admit the application, records each failed criterion, and tells the candidate which criteria were not satisfied.<br>2. Review: the Hiring Manager opens the screened-out applications of a posting and sees each with its failed criteria, its resume, and its answers. |
| Exception Flows   | E1. A resume field needed by a Must-have criterion is empty: the system treats the criterion as not satisfied.<br>E2. The posting has no screened-out applications: the system shows an empty list.<br>E3. The account is not permitted to manage the posting: the system denies the review. |
| Postconditions    | Each application is admitted (possibly "Awaiting Pipeline") or screened out with its failed criteria recorded. Reviewing changes nothing. |
| Business Rules    | Only Must-have criteria decide admission; Nice-to-have criteria never exclude a candidate. A screened-out application is never admitted. A posting can be published before a pipeline is published. |
| Related Requirements | FR-JOB-025, FR-JOB-105, FR-JOB-107 to FR-JOB-111 |

## 4.6 ANA Use Cases