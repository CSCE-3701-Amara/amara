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
## 4.5 JOB Use Cases
### UC-JOB-01: Create and Edit Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-01 |
| Use Case Name     | Create and Edit Job Posting |
| Actor(s)          | Hiring Manager, System, Interviewee (notified) |
| Description       | A Hiring Manager creates a job posting as a Draft, completes its content and hiring configuration, and edits it over time with every change recorded. |
| Preconditions     | The company is registered and verified (AUTH). The Hiring Manager is signed in with an account permitted to manage postings and the company is not banned. |
| Trigger           | The Hiring Manager chooses to create a posting or to edit an existing one. |
| Main Flow         | 1. The Hiring Manager opens the posting editor.<br>2. The Hiring Manager enters the content and hiring configuration.<br>3. The system autosaves the Draft while there are unsaved changes and shows the time of the last save.<br>4. The Hiring Manager saves the posting; the system stores it in the Draft state and records the creation or edit in the audit log. |
| Alternative Flows | 2a. The Hiring Manager saves a Draft with incomplete fields.<br>2b. The Hiring Manager duplicates an existing posting into a new Draft that copies content and criteria but not applications or dates (optional).<br>3a. The posting is Published: autosave does not apply and changes take effect only when the Hiring Manager saves.<br>4a. After saving a Published posting whose title, salary range, location, or work mode changed, the system notifies Interviewees with an active application. |
| Exception Flows   | E1. The Hiring Manager leaves the editor with unsaved changes: the system offers to save, discard, or stay.<br>E2. The Hiring Manager leaves a new posting never saved as a Draft: the system offers to save it as a Draft, discard it, or stay.<br>E3. Autosave fails: the system states that the latest changes are not saved, keeps the entered data on screen.<br>E4. A Published posting is edited in a way the hiring rules do not allow (deadline shortened after it passed, Must-have criteria changed after the first application): the system rejects that change.<br>E5. Two accounts of the same company edit the same posting: the system saves the first change, rejects the second, and shows the current version.<br>E6. The company is banned: the system refuses the edit and shows that the account is banned.<br>E7. The account is not permitted by the AUTH permission matrix: the system denies the action. |
| Postconditions    | The posting is stored (Draft, or updated Published) and the change is in the audit log. Autosaves do not create versions or audit entries. |
| Business Rules    | Autosave applies to Drafts only. Autosaves create no separate version or audit entry. Hiring configuration is defined once, in the APP epic. |

### UC-JOB-02: Set Posting Language, Currency and Tags

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-02 |
| Use Case Name     | Set Posting Language, Currency and Tags |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager states the language and salary currency of a posting and classifies it with tags used by filters, search, and recommendations. |
| Preconditions     | A Draft or Published posting exists and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager edits the language, currency, or tags of a posting. |
| Main Flow         | 1. The Hiring Manager selects the posting language from the standard list of languages.<br>2. If the posting has a salary range, the Hiring Manager selects its currency from the standard list of currencies.<br>3. The Hiring Manager adds or removes tags chosen from the System Admin-managed taxonomy.<br>4. The system saves the choices and records them in the audit log. |
| Alternative Flows | 1. The Hiring Manager writes the posting in any language and uses any currency; the system does not restrict either. |
| Exception Flows   | E1. A tag was removed from the taxonomy by a System Admin: the posting keeps it, and it cannot be added to other postings.<br>E2. Publishing without a language, or with a salary range and no currency: the posting stays in Draft and the missing item is listed. |
| Postconditions    | The posting has a language, a currency (if it shows a salary range), and its tags. |
| Business Rules    | Posting text is shown as written and amounts are shown in the chosen currency, without translation or conversion. Tags come only from the System Admin taxonomy. |

### UC-JOB-03: Preview Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-03 |
| Use Case Name     | Preview Job Posting |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager sees a posting as Interviewees will see it, before or after publication. |
| Preconditions     | A posting exists, in any state, and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager chooses Preview in the editor. |
| Main Flow         | 1. The Hiring Manager chooses Preview, including with unsaved edits.<br>2. The system shows the posting as it appears in the jobs feed and on the detail page.<br>3. The Hiring Manager closes the preview and returns to the editor. |
| Alternative Flows | 1. A Draft posting is previewed before it has ever been published. |
| Exception Flows   | E1. The posting cannot be rendered because required fields are missing: the system shows the preview with the missing fields marked. |
| Postconditions    | No posting, application, or notification data changed, and nobody was notified. |
| Business Rules    | Previewing never changes data or notifies anyone. |

### UC-JOB-04: View and Restore Posting Version

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-04 |
| Use Case Name     | View and Restore Posting Version |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager reviews earlier versions of a Published posting and restores one as a new current version, with history kept. |
| Preconditions     | The posting has at least one previous version and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager opens the version history of a posting. |
| Main Flow         | 1. The system lists the previous versions with author and time.<br>2. The Hiring Manager opens a version to view it.<br>3. The Hiring Manager selects Restore.<br>4. The system creates a new current version with the content of the selected version and keeps all other versions.<br>5. The system records the restoration with author, time, and version restored. |
| Alternative Flows | 1. The Hiring Manager only views versions and does not restore. |
| Exception Flows   | E1. The posting has received an application: the system keeps the current Must-have criteria.<br>E2. The restored deadline is not later than the current one, or the current one has passed: the system keeps the current deadline.<br>E3. In both cases the system restores the other fields and tells the Hiring Manager which fields were not restored.<br>E4. The company is banned: the system refuses. |
| Postconditions    | A new current version exists, all earlier versions remain, and the restoration is in the audit log. |
| Business Rules    | Restoring never overwrites history and never bypasses locked Must-have criteria or the deadline rule. |

### UC-JOB-05: Publish Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-05 |
| Use Case Name     | Publish Job Posting |
| Actor(s)          | Hiring Manager, System, Interviewee (notified) |
| Description       | The Hiring Manager moves a Draft posting to Published after the system verifies it is complete, choosing whether it is public or private. |
| Preconditions     | The posting is in Draft. The Hiring Manager is permitted to manage it and the company is not banned. |
| Trigger           | The Hiring Manager chooses Publish. |
| Main Flow         | 1. The Hiring Manager chooses Public or Private visibility.<br>2. The system verifies the mandatory content and dates, the language, and the currency when a salary range exists.<br>3. The system sets the state to Published and records it in the audit log.<br>4. For a public posting, the system lists it in the feed, search results, and public company profile.<br>5. The system notifies Interviewees who follow the company. |
| Alternative Flows | 1. The Hiring Manager publishes before a pipeline is published; admitted applications wait until a pipeline is published.<br>4a. For a private posting, the system generates the private link and hides the posting from the feed, search, and profile; the Hiring Manager continues with managing private access. |
| Exception Flows   | E1. A check fails: the posting stays in Draft and the system lists every failed check.<br>E2. The company is banned: the system refuses and shows that the account is banned. |
| Postconditions    | The posting is Published with its chosen visibility, or remains Draft with the failed checks listed. |
| Business Rules    | Only the transitions Draft to Published, Published to Archived, and Archived to Published exist. |

### UC-JOB-06: Manage Private Access and Invitations

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-06 |
| Use Case Name     | Manage Private Access and Invitations |
| Actor(s)          | Hiring Manager, System, Interviewee (invited) |
| Description       | The Hiring Manager shares a private posting by private link or email invitation, tracks the invitations sent, and revokes the link when needed. |
| Preconditions     | The posting is Published with Private visibility. |
| Trigger           | The Hiring Manager opens the sharing controls of a private posting. |
| Main Flow         | 1. The system shows the private link.<br>2. The Hiring Manager enters one or more email addresses and sends the invitation containing the private link.<br>3. The system sends the emails and records each invitation with address, time sent, and delivery status.<br>4. The Hiring Manager views the invitations list of the posting. |
| Alternative Flows | 1. The Hiring Manager copies the private link and shares it outside the platform.<br>2. The Hiring Manager revokes the link; the system generates a new one and the revoked link stops giving access. |
| Exception Flows   | E1. An invitation email fails: the system marks the address as failed and keeps the link available to copy.<br>E2. An invalid email address is entered: the system rejects it and identifies it.<br>E3. A revoked, invalid, or unknown link is opened: the system shows "This posting is not available" without revealing whether the posting exists. |
| Postconditions    | The invitations are recorded and visible to the Hiring Manager, and invited Interviewees can reach the posting. |
| Business Rules    | A private posting is reached only through its private link or an email invitation. The Hiring Manager sees the invitations sent for each private posting. |

### UC-JOB-07: Change Posting Visibility

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-07 |
| Use Case Name     | Change Posting Visibility |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager switches a Published posting between Public and Private. |
| Preconditions     | The posting is Published and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager changes the visibility setting. |
| Main Flow         | 1. The Hiring Manager selects the other visibility.<br>2. The system applies it; a posting made private leaves the feed, search, and company profile, and a posting made public enters them.<br>3. The system records the change in the audit log. |
| Alternative Flows | 1. When a posting becomes private, the system generates a private link. |
| Exception Flows   | E1. The company is banned: the system refuses. |
| Postconditions    | The posting has the new visibility and every existing application is unchanged. Saved entries of Interviewees who can no longer see it show "No longer available". |
| Business Rules    | Changing visibility never changes an existing application. |

### UC-JOB-08: Archive Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-08 |
| Use Case Name     | Archive Job Posting |
| Actor(s)          | Hiring Manager, System, Interviewee (notified) |
| Description       | The Hiring Manager archives a Published posting, which stops new applications and suspends existing ones without deleting or deciding them. |
| Preconditions     | The posting is Published and the Hiring Manager is permitted to manage it. |
| Trigger           | The Hiring Manager chooses Archive. |
| Main Flow         | 1. The system shows the number of active applications and asks for confirmation.<br>2. The Hiring Manager confirms.<br>3. The system sets the state to Archived, removes the posting from the feed, search, and company profile, and stops new applications.<br>4. The system suspends all applications: no module can be started or continued, no live interview booked, no decision recorded.<br>5. The system pauses the module completion windows and marks no module as missed.<br>6. The system notifies Interviewees with an active application. |
| Alternative Flows | 1. An Interviewee with an application opens the application tracking view and sees that the posting is archived and cannot be completed.<br>2. The Interviewee withdraws the application. |
| Exception Flows   | E1. An Interviewee tries to start a module or book an interview: the system refuses and shows that the posting is archived.<br>E2. The Hiring Manager cancels the confirmation: the posting stays Published. |
| Postconditions    | The posting is Archived. Its applications, results, statuses, and pipeline data are retained. |
| Business Rules    | Archiving suspends applications but never deletes or decides them. An Archived posting is shown only to Interviewees who applied to or saved it. |

### UC-JOB-09: Republish Archived Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-09 |
| Use Case Name     | Republish Archived Posting |
| Actor(s)          | Hiring Manager, System, Interviewee (notified) |
| Description       | The Hiring Manager returns an Archived posting to Published and the suspended applications resume. |
| Preconditions     | The posting is Archived, the Hiring Manager is permitted to manage it, and the company is not banned. |
| Trigger           | The Hiring Manager chooses Republish. |
| Main Flow         | 1. The Hiring Manager chooses Republish.<br>2. The system returns the posting to Published with its previous visibility.<br>3. The system resumes the suspended applications and extends each module completion window by the time the posting was archived.<br>4. The system notifies Interviewees with an active application.<br>5. The Hiring Manager decides on applications (pass, reject, offer) through APP. |
| Alternative Flows | 3a. The application deadline has passed: the posting stays closed to new applications while existing applications resume. |
| Exception Flows   | E1. The company is banned: the system refuses.<br>E2. Republishing would leave mandatory data invalid: the system lists the failed checks, as when publishing. |
| Postconditions    | The posting is Published. New applications are accepted only if the deadline has not passed. |
| Business Rules    | A posting can be reopened to new applications only before its deadline passes. Republishing does not decide any application. |

### UC-JOB-10: Delete Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-10 |
| Use Case Name     | Delete Job Posting |
| Actor(s)          | Hiring Manager, System Admin, System |
| Description       | A posting is permanently removed. A Hiring Manager may delete only an unused Draft, while a System Admin may delete a posting in any state. |
| Preconditions     | The posting exists. The actor is permitted to manage it (Hiring Manager) or is a System Admin. |
| Trigger           | The actor chooses Delete on a posting. |
| Main Flow         | 1. The system checks that the Hiring Manager's posting is a Draft that never received an application.<br>2. The system asks for confirmation and states that deletion is permanent.<br>3. The actor confirms.<br>4. The system deletes the posting permanently and records the deletion in the audit log. |
| Alternative Flows | 1. The System Admin deletes a posting in any state, including one with applications; the confirmation shows the number of applications affected. |
| Exception Flows   | E1. The Hiring Manager requests deletion of a Published or Archived posting, or one with applications: the system refuses and offers archiving instead.<br>E2. The actor cancels: nothing is deleted. |
| Postconditions    | The posting is permanently removed and the deletion is in the audit log. Saved entries and invitations to it show "No longer available". |
| Business Rules    | A Hiring Manager can only delete a Draft that never received an application. A System Admin can delete a posting at any time. |

### UC-JOB-11: Handle Postings of a Banned Company

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-11 |
| Use Case Name     | Handle Postings of a Banned Company |
| Actor(s)          | System, System Admin (bans and resolves through AUTH), Interviewee (notified) |
| Description       | When a company is banned, the system archives its Published postings and suspends their applications until the ban is resolved. |
| Preconditions     | A System Admin bans a company. |
| Trigger           | The ban is applied, or later resolved. |
| Main Flow         | 1. The system archives every Published posting of the company and suspends their applications.<br>2. The system keeps Draft postings as Draft.<br>3. The system excludes the company's postings from the feed and search.<br>4. The system blocks publishing or republishing while the ban lasts.<br>5. The system notifies Interviewees with an active application, without stating the reason. |
| Alternative Flows | 1. When the ban is resolved, the system returns each posting archived because of the ban to Published with its previous visibility and resumes its applications as when a posting is republished.<br>2. Postings that were Archived before the ban stay Archived.<br>3. The system notifies Interviewees that their application resumed. |
| Exception Flows   | E1. A banned company tries to edit or publish: the system refuses and shows that the account is banned. |
| Postconditions    | During the ban no posting of the company is Published. After resolution the affected postings are Published again. |
| Business Rules    | While a company is banned, its postings stay archived and nothing can be published. The Hiring Manager decides on applications after the ban is resolved. |

### UC-JOB-12: View Company Posting List

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-12 |
| Use Case Name     | View Company Posting List |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager sees all postings of the company in every state. |
| Preconditions     | The Hiring Manager is signed in with an account permitted to view postings. |
| Trigger           | The Hiring Manager opens the company's postings page. |
| Main Flow         | 1. The system lists all postings of the company with title, state, visibility, number of applications, and application deadline.<br>2. The Hiring Manager opens a posting to edit, archive, republish, or delete it. |
| Alternative Flows | 1. The Hiring Manager filters the list by state and visibility and searches it by title. |
| Exception Flows   | E1. The company has no postings: the system shows an empty list and offers to create a posting. |
| Postconditions    | The Hiring Manager has an up-to-date view of the company's postings. |
| Business Rules    | The list includes Draft, Published, and Archived postings. |

### UC-JOB-13: Switch Employer Interface Language

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-13 |
| Use Case Name     | Switch Employer Interface Language |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager changes the language of the employer-side interface. |
| Preconditions     | The Hiring Manager is signed in. |
| Trigger           | The Hiring Manager opens the language selector. |
| Main Flow         | 1. The Hiring Manager selects one of the interface languages supported by the platform.<br>2. The system applies the language immediately, without a new sign-in and without losing unsaved changes.<br>3. The system keeps the choice for later sessions. |
| Alternative Flows | 1. The Hiring Manager switches while editing a posting; the editor content stays as entered. |
| Exception Flows   | E1. The chosen language cannot be loaded: the system keeps the current language and shows an error message. |
| Postconditions    | The employer interface uses the chosen language. No posting content changed. |
| Business Rules    | Switching the interface language never changes posting content. |

### UC-JOB-14: Browse Jobs Feed

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-14 |
| Use Case Name     | Browse Jobs Feed |
| Actor(s)          | Interviewee, System |
| Description       | An Interviewee browses public Published postings, optionally ordered by match with their job preferences. |
| Preconditions     | The user is signed in. |
| Trigger           | The Interviewee opens the jobs feed. |
| Main Flow         | 1. The system lists public Published postings that are open or not yet open for applications, newest first, in pages.<br>2. The system shows for each item the title, company, location, work mode, seniority, employment type, language, tags, deadline, and the salary range only where shown.<br>3. The system marks items the Interviewee saved or applied to.<br>4. The Interviewee scrolls or pages through the feed and opens an item. |
| Alternative Flows | 1. A visitor who is not signed in is sent to sign-in or registration and the requested page opens afterwards.<br>2. The Interviewee selects the Recommended ordering.<br>3. The Interviewee limits the feed to followed companies. |
| Exception Flows   | E1. No postings are available: the system shows an empty-feed message.<br>E2. Postings of deactivated or banned companies are never shown. |
| Postconditions    | The Interviewee has seen the current public postings. |
| Business Rules    | Only public Published postings appear in the feed. Every user must be signed in to browse. |

### UC-JOB-15: Filter and Sort Postings

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-15 |
| Use Case Name     | Filter and Sort Postings |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee narrows and orders the feed or search results. |
| Preconditions     | The user is signed in and the feed or a search result is displayed. |
| Trigger           | The Interviewee selects a filter or a sort option. |
| Main Flow         | 1. The Interviewee selects filters: tag, required skill, seniority, employment type, work mode, location, salary range, salary currency, posting language, company, posting date range.<br>2. The system returns postings matching every selected filter; several values within one filter match any of them.<br>3. The system displays the active filters.<br>4. The Interviewee sorts by newest, nearest deadline, or relevance for search results. |
| Alternative Flows | 1. The Interviewee removes one filter or all filters at once.<br>2. The same filters and sorting apply to the feed and to search results. |
| Exception Flows   | E1. No posting matches: the system shows a message and offers to clear the filters.<br>E2. Salary minimum above maximum: the system identifies the invalid input and keeps the previous results.<br>E3. A salary filter is active: postings that do not show a salary range are excluded and the Interviewee is told. |
| Postconditions    | The displayed postings match the selected filters and order. |
| Business Rules    | A salary filter applies only to postings in the selected currency and never converts amounts. A hidden salary range is never revealed through filters. |

### UC-JOB-16: Search Postings and Companies

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-16 |
| Use Case Name     | Search Postings and Companies |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee finds postings and companies by keyword. |
| Preconditions     | The user is signed in. |
| Trigger           | The Interviewee enters keywords in the search box, available on every page. |
| Main Flow         | 1. The Interviewee enters at least 2 characters and submits.<br>2. The system matches the keywords against posting title, description, tags, required skills, company name, and company names.<br>3. The system returns results grouped by type (postings, companies), ranked by relevance, in pages.<br>4. The Interviewee opens a result. |
| Alternative Flows | 1. The Interviewee applies filters and sorting to the posting results.<br>2. Queries with common spelling variations still return results (optional). |
| Exception Flows   | E1. The query is too short: the system says so and keeps the previous results.<br>E2. No result: the system shows a no-results message.<br>E3. The search service is unavailable: the system shows "Search is temporarily unavailable" and the feed and filters keep working. |
| Postconditions    | The Interviewee sees only postings and companies they are allowed to see. |
| Business Rules    | Search returns only public Published postings and public company profiles, never Draft, Archived, or private postings of other companies. |

### UC-JOB-17: View Posting Details

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-17 |
| Use Case Name     | View Posting Details |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee reads everything needed to decide whether to apply and sees which action is possible. |
| Preconditions     | The user is signed in. The posting is public Published, or the user holds its private link or an invitation, or the user applied to or saved it. |
| Trigger           | The Interviewee opens a posting from the feed, search, saved list, invitations, or a shared link. |
| Main Flow         | 1. The system displays the title, company (linked to its profile), description, location, work mode, seniority, employment type, language, tags, application dates, the Must-have and Nice-to-have groups, and the salary range only where shown.<br>2. The system displays whether the posting is open, opens on a date, closed, or archived.<br>3. The system offers Apply on an open posting, continuing to the application.<br>4. The system offers Save and Follow company. |
| Alternative Flows | 1. The posting is closed: the system states whether the deadline passed or the application cap was reached.<br>2. The posting is Archived: it is shown only to Interviewees who applied to or saved it, with a notice that it no longer accepts applications. |
| Exception Flows   | E1. The Interviewee cannot apply: the system states what is missing, such as an incomplete resume or an existing application.<br>E2. The link is revoked, invalid, or unknown: the system shows "This posting is not available" without revealing whether it exists. |
| Postconditions    | The Interviewee is informed and may apply, save, or follow. |
| Business Rules    | Salary is shown only where the posting shows it. |

### UC-JOB-18: Save and Unsave Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-18 |
| Use Case Name     | Save and Unsave Posting |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee keeps a personal list of postings to return to. |
| Preconditions     | The Interviewee is signed in. |
| Trigger           | The Interviewee chooses Save or Unsave on a posting. |
| Main Flow         | 1. The Interviewee saves a posting from the feed, search results, or detail page.<br>2. The system adds it to the saved list.<br>3. The Interviewee opens the saved list; the system shows title, company, deadline, and whether each posting is open, closed, or archived. |
| Alternative Flows | 1. The Interviewee unsaves a posting from any of those places.<br>2. The system reminds the Interviewee 48 hours before the deadline of a saved posting not yet applied to (optional). |
| Exception Flows   | E1. A saved posting is no longer visible (made private or removed): the entry is labelled "No longer available" and shows none of its content. |
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
| Main Flow         | 1. The system matches invitations to the account by its verified email address.<br>2. The system lists each invitation with posting title, company, date invited, and whether the posting is open or closed.<br>3. The Interviewee opens a posting from the list and continues to the posting details. |
| Alternative Flows | 1. The invitation was sent before the account existed; it still appears once the email is verified. |
| Exception Flows   | E1. The posting was deleted or archived, or its private link revoked: the invitation is labelled "No longer available" and shows none of the posting's content.<br>E2. The Interviewee has no invitations: the system shows an empty list. |
| Postconditions    | The Interviewee sees only invitations sent to their own email address. |
| Business Rules    | Invitations are matched by verified email only. |

### UC-JOB-20: Follow Company

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-20 |
| Use Case Name     | Follow Company |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee follows companies to see and be notified about their new postings. |
| Preconditions     | The Interviewee is signed in. |
| Trigger           | The Interviewee chooses Follow on a company profile or posting detail page. |
| Main Flow         | 1. The Interviewee follows a company.<br>2. The system adds it to the followed list.<br>3. When the company publishes a public posting, the system notifies the Interviewee. |
| Alternative Flows | 1. The Interviewee unfollows a company.<br>2. The Interviewee turns off notifications for one followed company (optional).<br>3. The Interviewee limits the feed to followed companies. |
| Exception Flows   | E1. The company is deactivated or banned: its postings are not shown and no notification is sent. |
| Postconditions    | The followed list is updated. |
| Business Rules    | Followed companies belong to the Interviewee; Hiring Managers and System Admins never see who follows a company. |

### UC-JOB-21: Get Recommendations

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-21 |
| Use Case Name     | Get Recommendations |
| Actor(s)          | Interviewee, System |
| Description       | The system suggests relevant postings based on job preferences and posting data. |
| Preconditions     | The Interviewee is signed in. For the Recommended ordering, job preferences are set. |
| Trigger           | The Interviewee selects the Recommended ordering, opens a posting, or submits an application. |
| Main Flow         | 1. The Interviewee selects Recommended ordering in the feed.<br>2. The system ranks postings by match with the preferences.<br>3. The system shows for each recommended posting the preference or tag that caused it. |
| Alternative Flows | 1. On a detail page, the system shows up to 5 similar postings based on shared tags, skills, and seniority (optional).<br>2. After an application is submitted, the system shows similar postings (optional). |
| Exception Flows   | E1. No preferences are set: the Recommended ordering is not offered.<br>E2. Postings the Interviewee has applied to are excluded. |
| Postconditions    | The Interviewee sees ranked or similar postings. |
| Business Rules    | Recommendations use only structured posting data and job preferences, never name, photo, gender, age, or nationality. |

### UC-JOB-22: Manage Saved Searches and New-Jobs Tracker

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-22 |
| Use Case Name     | Manage Saved Searches and New-Jobs Tracker |
| Actor(s)          | Interviewee, System |
| Description       | The Interviewee saves a search with its filters and is alerted when new postings match. |
| Preconditions     | The Interviewee is signed in. |
| Trigger           | The Interviewee chooses Save search. |
| Main Flow         | 1. The Interviewee names and saves the current search and filters.<br>2. When a newly published posting matches the saved search, the system notifies the Interviewee. |
| Alternative Flows | 1. The Interviewee lists and deletes saved searches.<br>2. The Interviewee chooses immediate notifications or a daily digest. |
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
| Main Flow         | 1. The system provides a link to copy.<br>2. The Interviewee copies and shares it.<br>3. A recipient opens the link, signs in, and the detail page opens. |
| Alternative Flows | 1. For a private posting, the shared link is its private link. |
| Exception Flows   | E1. The recipient is not signed in: the system sends them to sign-in or registration and opens the page afterwards.<br>E2. The link is revoked or the posting no longer available: the system shows "This posting is not available". |
| Postconditions    | The Interviewee has a shareable link. |
| Business Rules    | Sharing never bypasses sign-in. |

### UC-JOB-24: Award Verified Crest

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-24 |
| Use Case Name     | Award Verified Crest |
| Actor(s)          | System, Interviewee |
| Description       | When an Interviewee successfully completes an assessment, the system awards a Verified Crest tied to the skill or competency the assessment verified. |
| Preconditions     | The Interviewee has completed an assessment module for which a Crest can be awarded. |
| Trigger           | An assessment module is completed. |
| Main Flow         | 1. The system determines whether the assessment was completed successfully, according to the criteria required for its Crest.<br>2. If the criteria are satisfied, the system awards the Crest.<br>3. The system associates the Crest with the verified skill or competency.<br>4. The system associates the Crest with the Interviewee's profile. |
| Alternative Flows | 1. The criteria are not satisfied: no Crest is awarded and the assessment result is unaffected. |
| Exception Flows   | E1. The Crest cannot be awarded or attached to the profile (service or data error): the system keeps the assessment result and does not mark the Crest as awarded until it is attached. |
| Postconditions    | The Interviewee has the Crest on their profile, or none if the criteria were not met. |
| Business Rules    | A Crest is awarded only for the successful completion of an assessment and is tied to the skill or competency it verified. |

### UC-JOB-25: Apply Interviewee Anonymity

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-25 |
| Use Case Name     | Apply Interviewee Anonymity |
| Actor(s)          | System, Hiring Manager |
| Description       | While a stage is anonymous, the system hides identifying information about the Interviewee from the Hiring Manager and keeps qualification information visible. |
| Preconditions     | The Hiring Manager enabled anonymity for the pipeline and chose the reveal point. The Interviewee was admitted to the pipeline. |
| Trigger           | The Hiring Manager opens an application. |
| Main Flow         | 1. The system determines whether anonymity applies to the application's current stage.<br>2. The system hides the information designated as non-disclosable.<br>3. The system continues to display scores, module results, skills, and Crests.<br>4. The system applies the hiding in every view, including resume copies, the dashboard, scorecards, and notes. |
| Alternative Flows | 1. Anonymity does not apply to the stage: the system displays the full application. |
| Exception Flows   | E1. Access to hidden information is attempted before the reveal point: the system denies it and logs the attempt. |
| Postconditions    | Identifying information stays hidden and qualification information stays visible during anonymous stages. |
| Business Rules    | While anonymity applies, qualification information stays visible and identifying information stays hidden in every view. Optional (COULD). |

### UC-JOB-26: Reveal Interviewee Identity

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-26 |
| Use Case Name     | Reveal Interviewee Identity |
| Actor(s)          | System, Hiring Manager |
| Description       | When an Interviewee reaches the configured reveal point, the system makes the hidden identity available to authorized users and keeps the application history. |
| Preconditions     | Anonymity applies to the pipeline and an identity-reveal point is configured. |
| Trigger           | The Interviewee reaches the identity-disclosure stage. |
| Main Flow         | 1. The system determines that the Interviewee reached the configured stage.<br>2. The system makes the identity available to authorized users.<br>3. The system retains the Interviewee's application and assessment history. |
| Alternative Flows | 1. Identity is revealed at the latest when an offer is created. |
| Exception Flows   | E1. The reveal cannot be applied (service error): the system keeps the identity hidden and does not create an offer until the reveal succeeds. |
| Postconditions    | The identity is visible to authorized users and the history is intact. |
| Business Rules    | Identity disclosure never removes application or assessment history. Reveal is permanent. Optional (COULD). |

### UC-JOB-27: Request Documents from Applicants

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-JOB-27 |
| Use Case Name     | Request Documents from Applicants |
| Actor(s)          | Hiring Manager, Interviewee, System |
| Description       | The Hiring Manager asks one or more applicants of a posting for supporting documents, and each Interviewee uploads them to their application. |
| Preconditions     | The posting is Published, the Hiring Manager is permitted to manage it, and the application is active (not withdrawn, rejected, or screened out). The Interviewee is signed in to respond. |
| Trigger           | The Hiring Manager chooses Request documents on an application or on a selection of applications. |
| Main Flow         | 1. The Hiring Manager selects the applications and names each document needed, with an optional description and a due date.<br>2. The system records the request and notifies each Interviewee through the APP notification mechanism.<br>3. The Interviewee opens the request in the application tracking view.<br>4. The Interviewee uploads each document in PDF or DOCX format.<br>5. The system stores the documents with the application and notifies the Hiring Manager.<br>6. The Hiring Manager views the submitted documents. |
| Alternative Flows | 1. The Hiring Manager requests documents from several applicants of the same stage in one action.<br>2. The Hiring Manager extends the due date or cancels a pending request; the system notifies the Interviewee. |
| Exception Flows   | E1. The posting is Archived or the company is banned: the system refuses new requests and suspends pending ones.<br>E2. Unsupported file format: the system rejects the file and lists the accepted formats.<br>E3. The due date passes without upload: the system marks the request as overdue and notifies the Hiring Manager, who decides how to proceed.<br>E4. The Interviewee withdrew or was rejected: the system does not allow the request or the upload.<br>E5. The upload fails: the system shows an error message and allows the Interviewee to retry. |
| Postconditions    | The request is recorded, and the submitted documents are stored with the application and visible to the Hiring Manager and assigned reviewers. |
| Business Rules    | A request never changes the submitted application, its resume copy, or its status. An Interviewee sees only the requests made on their own applications. Documents are shown to evaluators according to the anonymity rules. |

## 4.6 ANA Use Cases
