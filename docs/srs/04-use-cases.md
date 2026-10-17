## 4. Use cases

Owner: All | Status: Draft | Last updated: 2026-10-18 | Jira: AMARA-32

---

### 4.1 Actors
- Candidate
- System Administrator
- Interviewer
- Hiring Manager
- System

### 4.2 AUTH Use Cases

#### UC-AUTH-001: Candidate or Interviewer Registration
|Field|Description|
|-----|-----------|
|**Actor(s)**| Guest, System, Company Representative (optional) |
|**Preconditions**| The user does not have an account of the selected type. For interviewer registration, the user has a claimed association with a company that may be verified during registration. |
|**Trigger**| A user opens the candidate or interviewer registration page. |
| **Main Flow** | 1. The user selects the account type and enters their personal information, email address, and password.<br>2. For interviewer registration, the user enters company association information.<br>3. The user clicks **Register**.<br>4. The system validates the password, registration information, and email format.<br>5. The system verifies that the email address is not already registered.<br>6. For interviewer registration, the system validates the user's association with an approved company.<br>7. The system creates the account with the collected data and, for interviewers, associates it with the company.<br>8. The system sends a verification code to the user's email and redirects the user to the verification page.<br>9. The user enters the verification code.<br>10. The system checks the code and activates the account.|
|**Alternative Flows**|**Candidate MFA setup**: After activation, the candidate may choose to set up MFA, and the system associates the MFA method with the account.<br>**Interviewer MFA setup**: After email activation, the interviewer must set up MFA before the account can be used to log in.<br>**Interviewer company confirmation**: The system asks a company representative to confirm the interviewer's claimed association before activating the account.|
|**Exception Flows**|**Invalid registration information**: The system displays an error message.<br>**Company association cannot be verified**: The system does not activate the interviewer account and informs the user of the required next step.<br>**The user does not receive the code**: The user requests a new verification code, and the system sends it to the registered email address.<br>**Invalid verification code**: The system displays an error message and allows the user to re-enter the code or request a new verification code.|
|**Postconditions**|The user has an active candidate or interviewer account and can login.|

#### UC-AUTH-002: Company Registration
|Field|Description|
|-----|-----------|
|**Actor(s)**| Guest, System, System Admin |
|**Preconditions**| The user has the information required to register a company and does not already have an approved company account. |
|**Trigger**| A user opens the company registration page. |
| **Main Flow** | 1. The user enters the company information and the requesting person's information.<br>2. The user submits the company registration request.<br>3. The system validates the submitted information.<br>4. The system checks the company against official registries.<br>5. The system verifies the identity of the requesting person when required.<br>6. The system creates a pending company registration request.<br>7. The system notifies the user that the request is pending review.<br>8. The system or a system admin accepts the request after verification.<br>9. The system creates the company account and sends an activation code to the requesting person's email address.<br>10. The user enters the activation code.<br>11. The system activates the account and requires the user to set up MFA.|
|**Exception Flows**|**Invalid company or requester information**: The system displays an error message and asks the user to correct the information.<br>**Company not found in official registries**: The system rejects the request and notifies the user.<br>**Manual review required**: The system keeps the request pending until a system admin accepts or rejects it.<br>**Request rejected**: The system notifies the user and does not create an active company account.<br>**The user does not receive the activation code**: The user requests a new code, and the system sends it to the registered email address.|
|**Postconditions**|The company has an activated account with MFA configured and can login.|

#### UC-AUTH-003: User Login
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, Interviewer, System |
|**Preconditions**| The user has a registered and activated account. Company representatives must have an approved company account with MFA configured. Interviewers must have an active association with an approved company and MFA configured. |
|**Trigger**| The user opens the login page. |
| **Main Flow** | 1. The user enters their email address and password.<br>2. The user clicks **Login**.<br>3. The system validates the credentials.<br>4. The system checks that the account is activated and not deactivated or banned.<br>5. The system applies the authentication requirements for the user's account type.<br>6. The system creates an authenticated session with permissions scoped to the user's role and, where applicable, associated company.<br>7. The system redirects the user to the appropriate page or company workspace.|
|**Alternative Flows**|**MFA verification**: If MFA is configured, the system requests and validates the user's MFA code.<br>**Candidate social login**: The candidate chooses Google, GitHub, or LinkedIn, authenticates with the selected provider, and the system creates an authenticated session.<br>**Company MFA**: The system requires and validates the company representative's MFA code. Login is prevented if MFA is not configured.<br>**Interviewer MFA**: The system requires and validates the interviewer's MFA code. Login is prevented if MFA is not configured.<br>**Interviewer company access**: The system verifies the interviewer's company association before creating the session.|
|**Exception Flows**|**Invalid credentials**: The system displays an error message and does not create a session.<br>**Account is not activated**: The system informs the user that activation is required and offers to resend the activation code.<br>**Invalid MFA code**: The system displays an error message and allows the user to try again.<br>**Company association is no longer valid**: The system prevents interviewer login and informs the interviewer that access must be reviewed.<br>**Forgotten password**: The user requests a password reset, verifies their identity through email, and sets a new password before trying again.|
|**Postconditions**|The user has an authenticated session and can access resources permitted by their role and associated company.|

#### UC-AUTH-004: Submit Report
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, System Admin, System |
|**Preconditions**| The reporting user is authenticated and can view the target information permitted by their role. |
|**Trigger**| The user selects the option to report a company or candidate. |
| **Main Flow** | 1. The user opens the target's available profile or record.<br>2. The user selects **Report Company** or **Report Candidate**.<br>3. The user enters the reason and supporting details.<br>4. The user submits the report.<br>5. The system validates and records the report.<br>6. The system notifies system admins that a report has been submitted.<br>7. A system admin reviews the report and takes the appropriate action.|
|**Alternative Flows**|**Report a company**: A candidate reports a company profile.<br>**Report a candidate**: A company representative reports a candidate profile or test record.|
|**Exception Flows**|**Missing report details**: The system displays an error message and asks the user to complete the report.<br>**Report cannot be submitted**: The system displays an error message and allows the user to retry.|
|**Postconditions**|The report is stored for system-admin review, or the user is informed that the report was not submitted.|

#### UC-AUTH-005: Modify Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate is authenticated. |
|**Trigger**| The candidate opens the profile management page. |
| **Main Flow** | 1. The candidate views their current profile information.<br>2. The candidate edits the profile fields and, if needed, their job preferences.<br>3. The candidate submits the changes.<br>4. The system validates the updated information.<br>5. The system saves the changes and displays a confirmation.|
|**Exception Flows**|**Invalid profile information**: The system displays an error message and identifies the fields that need correction.<br>**Unauthorized access**: The system denies the action and does not change the profile.<br>**The changes cannot be saved**: The system displays an error message and preserves the current profile information.|
|**Postconditions**|The candidate's profile contains the validated changes, or remains unchanged if validation or saving fails.|

#### UC-AUTH-006: Deactivate/Delete Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate is authenticated. |
|**Trigger**| The candidate opens the account settings page and chooses to deactivate or delete the profile. |
| **Main Flow** | 1. The candidate chooses **Deactivate** or **Delete**.<br>2. The system explains the result of the selected action and requests confirmation.<br>3. The candidate confirms the action.<br>4. The system validates the candidate's authenticated session.<br>5. The system deactivates or deletes the profile according to the selected action.<br>6. The system terminates the candidate's authenticated session.<br>7. The system displays a confirmation.|
|**Exception Flows**|**The candidate cancels**: The system keeps the profile active and returns to account settings.<br>**Unauthorized or expired session**: The system denies the action and requests authentication.<br>**The action cannot be completed**: The system displays an error message and leaves the profile unchanged.|
|**Postconditions**|The profile is deactivated or deleted, and the candidate no longer has an active authenticated session.|

#### UC-AUTH-007: Share Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate is authenticated. |
|**Trigger**| The candidate chooses to share their profile. |
| **Main Flow** | 1. The candidate opens the profile sharing controls.<br>2. The candidate chooses to create a shareable link or QR code.<br>3. The system creates a shareable representation of the candidate's profile.<br>4. The system displays the link or QR code.<br>5. The candidate copies the link or displays the QR code to another person.|
|**Exception Flows**|**Profile is inactive**: The system prevents creation of a shareable link or QR code and displays an error message.<br>**The shareable representation cannot be created**: The system displays an error message and allows the candidate to retry.|
|**Postconditions**|The candidate has a shareable link or QR code for their active profile, or no shareable representation is created.|

#### UC-AUTH-008: Upload Supporting Documents
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, Interviewer, System |
|**Preconditions**| The user is authenticated |
|**Trigger**| The user opens the document upload controls. |
| **Main Flow** | 1. The user chooses a supporting document.<br>2. The user uploads the document.<br>3. The system validates that the document is in PDF or DOCX format and does not exceed 10 MB.<br>4. The system stores the document and associates it with the user's account.<br>5. If supported, the system extracts profile information from the document.<br>6. The system displays the extracted information for the user to review.<br>7. The user confirms the information to fill the related profile fields.<br>8. The system saves the confirmed profile information and displays a confirmation.|
|**Alternative Flows**|**No extraction**: If automatic extraction is unavailable or the user declines it, the system stores the document and the user completes the related profile fields manually.|
|**Exception Flows**|**Unsupported file format**: The system rejects the document and displays an error message listing the accepted formats.<br>**Document exceeds size limit**: The system rejects the document and displays an error message stating that the maximum size is 10 MB.<br>**Upload fails**: The system displays an error message and allows the user to retry.<br>**Information extraction fails**: The system stores the document but does not change the profile, and informs the user that the fields must be completed manually.<br>**Unauthorized access**: The system denies the upload and does not store the document.|
|**Postconditions**|The valid document is associated with the user's account, and any extracted profile information is saved only after the user confirms it.|

#### UC-AUTH-009: Submit Rating
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, System |
|**Preconditions**| The user is authenticated and has completed the testing interaction required to rate the target. |
|**Trigger**| The user chooses to rate a company or candidate. |
| **Main Flow** | 1. The user opens the target's rating form.<br>2. The user selects a rating value and enters optional supporting comments.<br>3. The user submits the rating.<br>4. The system validates that the user is eligible to rate the target.<br>5. The system records the rating and associates it with the target and the relevant testing interaction.<br>6. The system displays a confirmation.|
|**Alternative Flows**|**Rate a company**: A candidate rates a company after completing a test or testing interaction.<br>**Rate a candidate**: A company representative rates a candidate after testing the candidate.|
|**Exception Flows**|**User is not eligible**: The system prevents the rating and displays an error message.<br>**Invalid rating**: The system displays an error message and asks the user to provide a valid rating.<br>**Rating cannot be submitted**: The system displays an error message and allows the user to retry.|
|**Postconditions**|The eligible user's rating is recorded for the target, or no rating is recorded if validation or submission fails.|

#### UC-AUTH-010: Reset Password
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, Interviewer, System |
|**Preconditions**| The user has a registered account and access to its registered email address. |
|**Trigger**| The user selects **Forgot Password** or **Reset Password** on the login page. |
| **Main Flow** | 1. The user enters the email address associated with their account.<br>2. The user submits the password-reset request.<br>3. The system sends a verification code to the registered email address.<br>4. The user enters the verification code.<br>5. The system verifies the user's identity.<br>6. The user enters and confirms a new password.<br>7. The system validates the new password against the password policy.<br>8. The system updates the password, invalidates existing authenticated sessions, and displays a confirmation.|
|**Alternative Flows**|**Authenticated password change**: An authenticated user opens account settings, enters their current password and a new password, and the system updates the password after validation.|
|**Exception Flows**|**Unknown email address**: The system displays an error message and does not send a reset message.<br>**Invalid or expired reset code**: The system displays an error message and allows the user to request a new reset message.<br>**Invalid new password**: The system displays an error message and asks the user to enter a password that satisfies the password policy.<br>**Password reset cannot be completed**: The system displays an error message and does not change the existing password.|
|**Postconditions**|The user has a new valid password and any previous authenticated sessions have been terminated, or the existing password remains unchanged if the reset fails.|

#### UC-AUTH-011: Logout
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, Interviewer, System |
|**Preconditions**| The user has an authenticated session. |
|**Trigger**| The user selects **Logout**. |
| **Main Flow** | 1. The user selects **Logout**.<br>2. The system terminates the user's authenticated session.<br>3. The system redirects the user to the login page.|
|**Exception Flows**|**Session has already expired**: The system redirects the user to the login page and does not create a new session.|
|**Postconditions**|The user's authenticated session is terminated and the user cannot access protected resources without logging in again.|

#### UC-AUTH-012: View Public Company Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Guest, Candidate, Company Representative, System |
|**Preconditions**| The company has a public profile. |
|**Trigger**| A visitor opens a company's public profile page. |
|**Main Flow**| 1. The visitor opens the company profile link.<br>2. The system retrieves the company's public information.<br>3. The system displays the company information and its job postings.<br>4. The visitor selects a job posting to view its details.|
|*Exception Flows|Company profile does not exist: The system displays an error message or not-found page.<br>Company profile is unavailable*: The system displays an error message and does not expose private company information.|
|**Postconditions**|The visitor can view the company's public information and available job postings.|

#### UC-AUTH-013: Manage Users and Accounts
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to manage user accounts. |
|**Trigger**| The system admin opens the user management interface. |
| **Main Flow** | 1. The system admin searches for or selects a user account, or chooses to create an account.<br>2. The system displays the account information and current status, or collects the information required for the new account.<br>3. The system admin chooses an account management action.<br>4. The system admin confirms the action.<br>5. The system validates the account information, applies the action, and displays a confirmation.|
|**Alternative Flows**|**Create account**: The system admin enters the account information, confirms creation, and the system creates the account according to the applicable activation and verification rules.<br>**View account**: The system admin views the account information without changing it.<br>**Ban or deactivate account**: The system admin selects the action, and the system prevents the account from accessing the system.<br>**Reset credentials**: The system admin initiates a credential reset for the selected account.|
|**Exception Flows**|**Unauthorized access**: The system denies access to the user management interface.<br>**Account not found**: The system displays an error message and does not apply an action.<br>**Action cannot be completed**: The system displays an error message and leaves the account unchanged.|
|**Postconditions**|The selected account is viewed or has the confirmed management action applied.|

#### UC-AUTH-014: Manage System Entities
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to manage system entities. |
|**Trigger**| The system admin opens an entity management interface. |
| **Main Flow** | 1. The system admin selects an entity type, such as a candidate, company, job position, or test.<br>2. The system admin selects an existing entity or chooses to create one.<br>3. The system admin enters or edits the entity information.<br>4. The system validates the information.<br>5. The system admin confirms the action.<br>6. The system adds or updates the entity and displays a confirmation.|
|**Alternative Flows**|**Delete entity**: The system admin selects an existing entity, confirms deletion, and the system deletes it according to the platform's data-retention rules.|
|**Exception Flows**|**Unauthorized access**: The system denies the action.<br>**Invalid entity information**: The system displays an error message and asks the system admin to correct the information.<br>**Entity not found**: The system displays an error message and does not apply the action.<br>**Action cannot be completed**: The system displays an error message and preserves the existing entity.|
|**Postconditions**|The selected entity is added, edited, or deleted as confirmed by the system admin.|

#### UC-AUTH-015: Review Company Registration Requests
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System, Company Representative (optional)|
|**Preconditions**| The system admin is authenticated and authorized to review company registration requests. |
|**Trigger**| The system admin opens the company registration request queue. |
| **Main Flow** | 1. The system displays pending company registration requests.<br>2. The system admin selects a request.<br>3. The system displays the submitted company and requester information, including verification results.<br>4. The system admin reviews the request.<br>5. The system admin accepts or rejects the request.<br>6. The system records the decision and notifies the requester.|
|**Alternative Flows**|**Request requires more information**: The system admin marks the request as pending and requests additional information from the requester.|
|**Exception Flows**|**Unauthorized access**: The system denies access to the request queue.<br>**Request not found**: The system displays an error message and returns to the request queue.<br>**Decision cannot be recorded**: The system displays an error message and leaves the request pending.|
|**Postconditions**|The request is accepted, rejected, or remains pending with a recorded request for more information.|

#### UC-AUTH-016: Manage System Taxonomies and Constants
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to manage system reference values. |
|**Trigger**| The system admin opens the taxonomy and platform-configuration interface. |
| **Main Flow** | 1. The system admin selects a taxonomy, platform constant, or reference value.<br>2. The system displays its current values.<br>3. The system admin adds, edits, or removes a value.<br>4. The system validates the change.<br>5. The system admin confirms the change.<br>6. The system saves the change and displays a confirmation.|
|**Exception Flows**|**Unauthorized access**: The system denies access to the configuration interface.<br>**Invalid or conflicting value**: The system displays an error message and does not save the change.<br>**Change cannot be saved**: The system displays an error message and preserves the existing values.|
|**Postconditions**|The selected taxonomy, constant, or reference value is updated, or remains unchanged if the operation fails.|

#### UC-AUTH-017: View Administrative Reports
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to view administrative reports. |
|**Trigger**| The system admin opens the statistics and reporting interface. |
| **Main Flow** | 1. The system admin selects a report or statistics view.<br>2. The system retrieves aggregate data related to job postings.<br>3. The system displays the aggregate data and reporting results.<br>4. The system admin reviews the results for platform improvement or research.|
|**Exception Flows**|**Unauthorized access**: The system denies access to administrative reports.<br>**Data unavailable**: The system displays an error message and indicates that the report cannot currently be generated.<br>**Report generation fails**: The system displays an error message and does not display incomplete results.|
|**Postconditions**|The system admin has viewed the requested aggregate job-posting data, or no incomplete report is presented.|

---

### 4.3 CAND Use Cases

#### 4.3.1 UC-CAND-001 — Browse Job Postings

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to browse available Computer Science-related job opportunities. |
| Preconditions | The candidate is authenticated. Job postings are available in the system. |
| Trigger | The candidate opens the job discovery page. |
| Main Flow | 1. The candidate opens the job discovery page.<br>2. The system retrieves available job postings.<br>3. The system displays the available postings and their essential information.<br>4. The candidate browses the displayed postings. |
| Alternative Flows | 1. If no postings are available, the system informs the candidate that no jobs are currently available. |
| Exception Flows | 1. If the system cannot retrieve the postings, it displays an appropriate error message. |
| Postconditions | Available job postings are displayed to the candidate, or the candidate is informed that none are available. |
| Business Rules | Only job postings available for application shall be displayed as available opportunities. |

#### 4.3.2 UC-CAND-002 — Filter Job Postings

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to filter job postings according to compensation, position, job grade, and company. |
| Preconditions | The candidate is authenticated and has access to the job discovery page. |
| Trigger | The candidate selects one or more filtering criteria. |
| Main Flow | 1. The candidate opens the job discovery page.<br>2. The system displays the available filtering criteria.<br>3. The candidate selects one or more criteria, such as compensation, position, job grade, or company.<br>4. The candidate applies the selected filters.<br>5. The system retrieves postings matching the selected criteria.<br>6. The system displays the matching postings. |
| Alternative Flows | 1. If no postings match the selected criteria, the system informs the candidate that no matching jobs are available.<br>2. The candidate clears the selected filters, and the system displays all available postings. |
| Exception Flows | 1. If the system cannot retrieve or filter the postings, it displays an appropriate error message. |
| Postconditions | Job postings matching the selected criteria are displayed, or the candidate is informed that no matching postings are available. |
| Business Rules | The system shall support combining multiple filtering criteria. Only available job postings shall be displayed. |

#### 4.3.3 UC-CAND-003 — View Application History

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to view their previously submitted applications, current statuses, and associated details. |
| Preconditions | The candidate is authenticated. |
| Trigger | The candidate opens the application history page. |
| Main Flow | 1. The candidate opens the application history page.<br>2. The system retrieves applications associated with the candidate.<br>3. The system displays the relevant job postings and their current application statuses.<br>4. The candidate selects an application.<br>5. The system displays the details of the selected application. |
| Alternative Flows | 1. If the candidate has no applications, the system informs the candidate that no application history is available. |
| Exception Flows | 1. If the system cannot retrieve the application history, it displays an appropriate error message. |
| Postconditions | The candidate can view their application history and access the details of their applications. |
| Business Rules | Candidates shall only be able to access their own application information. |

#### 4.3.4 UC-CAND-004 — Track Assessment Progress

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to check the completion status of their assessments and resume their progress when supported. |
| Preconditions | The candidate is authenticated and has an assessment associated with an application. |
| Trigger | The candidate opens an assessment or checks its status. |
| Main Flow | 1. The candidate accesses the assessment status.<br>2. The system retrieves the recorded assessment state.<br>3. The system displays the assessment status, such as not started, in progress, submitted, or completed, as applicable.<br>4. If the assessment supports interruption and remains incomplete, the system allows the candidate to resume it. |
| Alternative Flows | 1. If the assessment has already been submitted, the system displays its submitted status.<br>2. If the assessment does not support interruption, the system displays the current status without offering resume functionality. |
| Exception Flows | 1. If the system cannot retrieve the assessment status, it displays an appropriate error message.<br>2. If saved progress cannot be restored, the system informs the candidate. |
| Postconditions | The candidate can view the recorded assessment status and resume the assessment when permitted. |
| Business Rules | Assessment progress shall be preserved where the assessment permits interruption. |

#### 4.3.5 UC-CAND-005 — View Assessment Result and Feedback

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to view assessment results and any available constructive feedback. |
| Preconditions | The candidate is authenticated and has submitted an assessment. The result or feedback is available for disclosure. |
| Trigger | The candidate opens the assessment results or feedback section of an application. |
| Main Flow | 1. The candidate selects an application.<br>2. The candidate opens the assessment results section.<br>3. The system verifies that the result and feedback may be disclosed.<br>4. The system displays the assessment outcome, including whether the candidate passed or failed, when disclosed.<br>5. The system displays the associated feedback and identified areas for improvement when available. |
| Alternative Flows | 1. If the result is available but feedback is not, the system displays the result alone.<br>2. If only part of the assessment information is available for disclosure, the system displays the permitted information. |
| Exception Flows | 1. If the result is still being evaluated, the system informs the candidate that it is pending.<br>2. If the result cannot be retrieved, the system displays an appropriate error message. |
| Postconditions | The candidate can view the assessment outcome and any feedback they are permitted to access. |
| Business Rules | Assessment results and feedback shall only be visible when disclosure is permitted by the assessment configuration. Feedback shall be associated with its corresponding assessment. |

#### 4.3.6 UC-CAND-006 — View Application Status

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to view the current status of their applications and determine whether a final hiring decision has been reached. |
| Preconditions | The candidate is authenticated and has at least one application recorded in the system. |
| Trigger | The candidate opens an application or checks its status. |
| Main Flow | 1. The candidate opens their application history.<br>2. The candidate selects an application.<br>3. The system retrieves the application's current status.<br>4. The system displays the current status.<br>5. The system indicates whether a final hiring decision has been reached. |
| Alternative Flows | 1. If the application is still being processed, the system displays its current status without presenting a final decision. |
| Exception Flows | 1. If the system cannot retrieve the application status, it displays an appropriate error message. |
| Postconditions | The candidate can view the latest recorded status of the selected application. |
| Business Rules | Application status shall reflect the recorded progress of the application through the hiring process. |

#### 4.3.7 UC-CAND-007 — Receive Final Hiring Decision

| Field | Description |
|---|---|
| Actor(s) | Candidate, Hiring Manager, System |
| Description | Allows a Candidate to view the final hiring decision and any feedback released by INVW according to the configured feedback mode. |
| Preconditions | The Candidate is authenticated and has an application with a final decision recorded by INVW. |
| Trigger | INVW releases a final hiring decision for the application. |
| Main Flow | 1. The Hiring Manager submits the final hiring decision through INVW.<br>2. INVW associates the decision with the corresponding application and applies the configured rejection-feedback mode (UC-INVW-009).<br>3. When INVW releases the decision and any permitted feedback, CAND makes them available to the Candidate.<br>4. The Candidate opens the application.<br>5. The system displays the final hiring decision and any released feedback. |
| Alternative Flows | 1. If feedback has not yet been released under the configured mode, the system displays the decision and current application status until feedback becomes available.<br>2. If the Candidate does not immediately open the application, the decision and released feedback remain accessible for later viewing. |
| Exception Flows | 1. If INVW has not recorded a final decision, the system continues to display the current application status.<br>2. If the decision or its associated feedback cannot be retrieved, the system displays an appropriate error message. |
| Postconditions | The Candidate can access the final hiring decision and any feedback released for that application. |
| Business Rules | INVW owns the final decision, application status, and configured rejection-feedback mode. CAND displays a decision and feedback only after INVW releases them (FR-INVW-080, FR-INVW-081). |

#### 4.3.8 UC-CAND-008 — Accept Offer

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows a Candidate to accept or decline a pending offer associated with an application. |
| Preconditions | The Candidate is authenticated and has a pending offer associated with an application. |
| Trigger | The candidate opens the offer details for an application. |
| Main Flow | 1. The Candidate opens the application containing the pending offer.<br>2. The system displays the offer, response deadline, and available response options.<br>3. Before the response deadline, the Candidate chooses to accept or decline the offer.<br>4. The system submits the response to INVW, which validates and records it and updates the authoritative offer/application status (UC-INVW-017).<br>5. CAND displays the updated application status and reflects acceptance in the Candidate's profile when the offer is accepted. |
| Alternative Flows | 1. If the Candidate declines the offer, INVW records the decline and updates the application status (UC-INVW-017). |
| Exception Flows | 1. If the response deadline has passed, INVW expires the offer and the system does not accept a late response (UC-INVW-017).<br>2. If the response cannot be recorded, the system displays an appropriate error message. |
| Postconditions | INVW has recorded a valid offer response; CAND displays the updated application state and reflects acceptance when applicable. |
| Business Rules | The response is associated with its offer and application. INVW owns response-deadline validation, response recording, and offer/application status transitions (FR-INVW-111–114). |

#### 4.3.9 UC-CAND-009 — View Earned Crests

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to view their earned digital Crests and the competencies they represent. |
| Preconditions | The candidate is authenticated and has access to their professional profile. |
| Trigger | The candidate opens the earned Crests section of their profile. |
| Main Flow | 1. The candidate opens the Crests section.<br>2. The system retrieves the candidate's earned Crests.<br>3. The system displays each earned Crest.<br>4. The system displays the skill or competency represented by each Crest.<br>5. The system displays the available verification information associated with each Crest. |
| Alternative Flows | 1. If the candidate has not earned any Crests, the system informs the candidate that no Crests are currently available. |
| Exception Flows | 1. If the system cannot retrieve Crest information, it displays an appropriate error message. |
| Postconditions | The candidate can view their earned Crests and associated verification information. |
| Business Rules | A Crest shall only be awarded when the candidate satisfies the criteria defined for the corresponding assessment. |

#### 4.3.10 UC-CAND-010 — Maintain Verified Skills

| Field | Description |
|---|---|
| Actor(s) | System |
| Description | Records and retains skills verified through successful assessments across eligible hiring processes. |
| Preconditions | An assessment has been evaluated, and the candidate has satisfied the criteria for verifying a skill. |
| Trigger | The system records a successful assessment outcome that qualifies the candidate for a verified skill. |
| Main Flow | 1. The system records the successful assessment outcome.<br>2. The system identifies the skill or competency verified by the assessment.<br>3. The system associates the verified skill with the candidate.<br>4. The system retains the verified skill in the candidate's professional profile.<br>5. The system makes the verified skill available to eligible platform features, including the dynamic resume. |
| Alternative Flows | 1. If the candidate already has the same verified skill, the system retains the existing skill and records additional verification information when applicable. |
| Exception Flows | 1. If the assessment outcome does not satisfy the verification criteria, the system does not award the corresponding verified skill.<br>2. If the system cannot record the verified skill, it handles the failure without incorrectly marking the skill as verified. |
| Postconditions | The verified skill is associated with the candidate and retained across eligible hiring processes. |
| Business Rules | Verified skills shall remain associated with the candidate after the hiring process in which they were earned has ended. Skills shall only be verified when the applicable assessment criteria are satisfied. |

#### 4.3.11 UC-CAND-011 — Generate Dynamic Resume

| Field | Description |
|---|---|
| Actor(s) | System |
| Description | Generates and updates the candidate's dynamic resume using verified skills, earned Crests, and relevant candidate information. |
| Preconditions | The candidate's professional profile exists, and the system can access the relevant candidate information. |
| Trigger | The candidate's relevant profile information, verified skills, or earned Crests change, or the resume needs to be generated for the first time. |
| Main Flow | 1. The system retrieves the candidate's relevant professional information.<br>2. The system retrieves the candidate's verified skills.<br>3. The system retrieves the candidate's earned Crests.<br>4. The system assembles the available information into the dynamic resume.<br>5. The system saves or updates the generated resume.<br>6. The system makes the updated resume available to the candidate. |
| Alternative Flows | 1. If the candidate has no verified skills or earned Crests, the system generates the resume using the other available relevant information.<br>2. If only some relevant candidate information changes, the system updates the resume to reflect those changes. |
| Exception Flows | 1. If required information cannot be retrieved, the system handles the generation failure and informs the candidate when appropriate.<br>2. If the resume cannot be generated or updated, the system records or reports the failure. |
| Postconditions | A dynamic resume is generated or updated using the candidate's available relevant information. |
| Business Rules | The dynamic resume shall include verified skills and earned Crests. The system shall update the resume when relevant candidate information changes. |

#### 4.3.12 UC-CAND-012 — View Dynamic Resume

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to view their dynamic resume and review their verified professional competencies. |
| Preconditions | The candidate is authenticated and has access to their professional profile. |
| Trigger | The candidate opens the dynamic resume section of their profile. |
| Main Flow | 1. The candidate opens the dynamic resume section.<br>2. The system retrieves the latest available dynamic resume.<br>3. The system displays the candidate's verified skills.<br>4. The system displays the candidate's earned Crests.<br>5. The candidate reviews the displayed professional information. |
| Alternative Flows | 1. If the resume has not yet been generated, the system initiates UC-CAND-011 — Generate Dynamic Resume. |
| Exception Flows | 1. If the system cannot retrieve or generate the resume, it displays an appropriate error message. |
| Postconditions | The candidate can view the latest available dynamic resume and its verified professional information. |
| Business Rules | The displayed verified skills and Crests shall correspond to the candidate's recorded achievements. |

#### 4.3.13 UC-CAND-013 — Follow Company

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to follow a company to keep track of its later job postings. |
| Preconditions | The candidate is authenticated and the company is available on the platform. |
| Trigger | The candidate chooses to follow a company. |
| Main Flow | 1. The candidate opens the company's page or profile.<br>2. The system displays the option to follow the company.<br>3. The candidate selects the follow option.<br>4. The system records the company as followed by the candidate. |
| Alternative Flows | 1. If the candidate already follows the company, the system indicates that the company is already followed. |
| Exception Flows | 1. If the system cannot record the follow action, it displays an appropriate error message. |
| Postconditions | The company is recorded as followed by the candidate. |
| Business Rules | The system SHOULD allow candidates to follow companies for later job postings. |

#### 4.3.14 UC-CAND-014 — Save Job Posting

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to save a job posting for later application. |
| Preconditions | The candidate is authenticated and the job posting is available for viewing. |
| Trigger | The candidate chooses to save a job posting. |
| Main Flow | 1. The candidate browses available job postings.<br>2. The candidate selects the option to save a posting.<br>3. The system records the posting as saved by the candidate.<br>4. The system confirms that the posting has been saved. |
| Alternative Flows | 1. If the candidate has already saved the posting, the system indicates that the posting is already saved. |
| Exception Flows | 1. If the system cannot save the posting, it displays an appropriate error message. |
| Postconditions | The job posting is recorded as saved by the candidate for later reference or application. |
| Business Rules | The system SHOULD allow candidates to save job postings for later application. |

---

### 4.4 INVW Use Cases
#### UC-INVW-001: Design Interview Pipeline

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager designs the interview pipeline of a posting as a directed acyclic graph (DAG) on a visual canvas. |
| Preconditions     | The Hiring Manager is signed in and permitted to manage the posting (AUTH). The posting exists. |
| Trigger           | The Hiring Manager opens the pipeline designer. |
| Main Flow         | 1. The Hiring Manager adds uniquely named stages from the module palette and connects them with directed links.<br>2. The System validates that the graph is acyclic and every stage is reachable from an entry stage.<br>3. The Hiring Manager configures a module, scheduled start/end times, and exactly one advancement rule for each assessment stage (FR-INVW-051, FR-INVW-094).<br>4. The Hiring Manager saves the design as a draft pipeline version; the System records its validation result. |
| Alternative Flows | 1. The Hiring Manager moves or deletes stages and links. |
| Exception Flows   | E1. A cycle is detected: the System identifies the stages involved (FR-INVW-005).<br>E2. A stage is unreachable from an entry stage: the System identifies it (FR-INVW-006).<br>E3. A stage has a duplicate name, missing module, or no advancement rule: the System identifies the stage and missing item (FR-INVW-007).<br>E4. An assessment stage's end time is not later than its start time: the System identifies the stage and invalid date range (FR-INVW-094). |
| Postconditions    | The draft pipeline version and its validation result are saved; a draft with failed checks cannot be published until corrected (FR-INVW-039). |
| Business Rules    | A draft is not available to Candidates until published through UC-INVW-005. |

#### UC-INVW-002: Configure Assessment Module

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager configures an individual assessment module, its content, scoring method, and timing. Live interviews are configured through UC-INVW-004. |
| Preconditions     | A pipeline design exists for the posting (UC-INVW-001). For library modules, the problem library contains matching problems. |
| Trigger           | The Hiring Manager selects and configures a module for a stage. |
| Main Flow         | 1. The Hiring Manager selects an individual assessment type from FR-INVW-008; the System attaches that predefined type's built-in Amara definition (FR-INVW-098).<br>2. The Hiring Manager selects or creates its content and configures the supported response mode and scoring method required for that type (FR-INVW-082). For coding problems, they also configure count, difficulty, point values, and selection mode (FR-INVW-009).<br>3. The Hiring Manager sets the active time limit and Candidate completion window (FR-INVW-011).<br>4. The System validates and saves the module configuration. |
| Alternative Flows | 2a. For random coding-problem selection, the Hiring Manager sets the problem count and difficulty mix; the System selects matching problems when it creates the test instance (FR-INVW-009, FR-INVW-010).<br>2b. For question-based English, programming-language, behavioral, or aptitude modules, the Hiring Manager sets the question count and selects or authors each question, response format, and scoring method (FR-INVW-091).<br>2c. For a video-recorded response module, the Hiring Manager configures each prompt and assigns a grader (FR-INVW-093). |
| Exception Flows   | E1. Required content, response mode, scoring method, active time limit, or Candidate completion window is missing: the System identifies the missing item and marks the module incomplete for publication (FR-INVW-011, FR-INVW-082).<br>E2. The problem library has no problems matching the configured count and difficulty mix: the System reports the mismatch and marks the module incomplete for publication. |
| Postconditions    | A complete individual assessment module is saved and available for pipeline publication (FR-INVW-082). |
| Business Rules    | Live interviews use the separate scheduling and scorecard flow in UC-INVW-004; LLM-paired programming is a future Could feature (UC-INVW-013); custom Amara mapping is configured through UC-INVW-015; test-instance execution is covered by UC-INVW-010. |

#### UC-INVW-003: Build Custom Test Module with AI Assistance

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager builds a custom test module in the Design Studio, optionally using AI to draft questions, and saves it for reuse. |
| Preconditions     | A pipeline design exists for the posting (UC-INVW-001). |
| Trigger           | The Hiring Manager opens the Design Studio to create a custom test module. |
| Main Flow         | 1. The Hiring Manager adds one or more multiple-choice, short-text, long-text, or code questions.<br>2. The System may draft questions from the Hiring Manager's skill and difficulty description.<br>3. The Hiring Manager reviews and approves every AI-drafted question before saving it (FR-INVW-018, FR-INVW-019).<br>4. The Hiring Manager assigns maximum points and an answer key to closed questions; for open-ended questions, the Hiring Manager assigns themselves or a stage-assigned Interviewer to grade (FR-INVW-020, FR-INVW-021).<br>5. The Hiring Manager assigns one or more Amarat to the module and maps every scored question to an Amara (UC-INVW-015, FR-INVW-099).<br>6. The Hiring Manager optionally saves the complete module to the company library (FR-INVW-022). |
| Alternative Flows | None. |
| Exception Flows   | E1. AI generation fails: the System discards partial output and notifies the Hiring Manager to write the questions manually. |
| Postconditions    | A custom module is configured and, optionally, saved to the company library. |
| Business Rules    | None. |

#### UC-INVW-004: Schedule and Conduct Live Interview

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Interviewer, Candidate, System |
| Description       | The System schedules a live-interview stage through the external meeting provider and records its outcome, notes, and scorecards. |
| Preconditions     | The stage is configured as a live interview with scheduled start/end times, an Amara mapping for every rubric criterion (UC-INVW-015, FR-INVW-099), and at least one assigned Interviewer or Hiring Manager. The external provider is reachable (EXT). The Candidate has reached the stage during its scheduled interval (FR-INVW-094). If anonymity is enabled, the Candidate's identity-reveal point is strictly before this stage (FR-INVW-072). |
| Trigger           | An assigned Interviewer publishes available time slots within the stage's scheduled interval. |
| Main Flow         | 1. The assigned Interviewer publishes available time slots within the stage's scheduled start/end times (FR-INVW-094).<br>2. The Candidate selects a slot that starts and ends within that interval; the System creates and confirms the online meeting through the external provider (FR-INVW-026).<br>3. The System sends the meeting time and join link to the Candidate and assigned Interviewers within 10 minutes (FR-INVW-027).<br>4. The System reminds the Candidate and each assigned Interviewer 24 hours and 1 hour before the interview (FR-INVW-028).<br>5. After the meeting, the System records the assigned Interviewers, scheduled start/end, and outcome (held, Candidate absent, or Interviewer absent) (FR-INVW-031).<br>6. Each assigned Interviewer saves timestamped notes and submits ratings for every rubric criterion and a written recommendation within 48 hours; the System tracks missing scorecards and reminds the Interviewer (FR-INVW-032, FR-INVW-034).<br>7. The System calculates the interview score after all assigned scorecards are submitted or a missing card is waived with a recorded reason; if none are submitted, it leaves the score unset and holds the Candidate for a Hiring Manager decision (FR-INVW-034, FR-INVW-035). |
| Alternative Flows | 2a. The Candidate reschedules up to 2 times, no later than 24 hours before the interview (FR-INVW-030).<br>5a. The Candidate or an assigned Interviewer is absent: the System records the outcome, notifies the Hiring Manager, and holds the Candidate for a decision (FR-INVW-036).<br>6a. A scorecard remains missing at 48 hours: the System notifies the Hiring Manager, who may waive it with a recorded reason (FR-INVW-034).<br>6b. No scorecard is submitted after the Hiring Manager's waiver: the System leaves the module score unset and holds the Candidate for a decision (FR-INVW-035). |
| Exception Flows   | E1. Meeting creation fails: the System retries 3 times within 15 minutes; if it still fails, the slot remains unconfirmed and the System notifies the Candidate and Interviewer to choose another time. |
| Postconditions    | The interview outcome and attendance are recorded; notes and submitted scorecards are stored; a score is calculated when scorecards are complete or waived. |
| Business Rules    | Interviewer notes are visible only to the Hiring Manager and Interviewers assigned to the stage (FR-INVW-033). Meeting times use the viewer's time zone (FR-INVW-029). |

#### UC-INVW-005: Publish and Migrate Pipeline Version

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Candidate, System |
| Description       | The Hiring Manager publishes a new pipeline version; the System migrates in-progress Candidates without losing results. |
| Preconditions     | At least one saved pipeline version exists. For migration, Candidates are in progress on the currently published version. |
| Trigger           | The Hiring Manager publishes a pipeline version. |
| Main Flow         | 1. The System validates the version against FR-INVW-005, FR-INVW-006, FR-INVW-007, FR-INVW-072, FR-INVW-082, FR-INVW-094, and FR-INVW-095 (FR-INVW-039). If any check fails, it lists every failure and stops publication.<br>2. The System shows a migration preview with the number of affected Candidates per stage and their proposed placements.<br>3. The Hiring Manager confirms publication; the System unpublishes the previous version, migrates in-progress Candidates without moving anyone backward, preserves completed results, and re-evaluates waiting Candidates under changed rules without reversing a pass (FR-INVW-043, FR-INVW-045, FR-INVW-048).<br>4. The System records the publication and placements and emails each Candidate who must complete an additional stage (FR-INVW-049, FR-INVW-050). |
| Alternative Flows | 1. Editing the published version creates a new draft version while the published version stays in force. |
| Exception Flows   | E1. A removed stage holds Candidates: publication is blocked until the Hiring Manager selects a destination (FR-INVW-046).<br>E2. An added stage affects in-progress Candidates: publication is blocked until the Hiring Manager chooses its scope (FR-INVW-047).<br>E3. An assessment stage has an invalid date range or a first exam that starts on or before the posting deadline: publication is blocked and the affected stage/date check is reported (FR-INVW-094, FR-INVW-095). |
| Postconditions    | Exactly one version is published; in-progress Candidates are migrated with completed results preserved; publication and placements are recorded. |
| Business Rules    | A stage with unchanged module and configuration is not repeated during migration (FR-INVW-045). |

#### UC-INVW-006: Advance Candidate Between Stages

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Candidate, System |
| Description       | The System applies the stage's configured rule to decide whether a Candidate passes, is rejected, or remains held for a Hiring Manager decision. |
| Preconditions     | A pipeline is published and the Candidate is active in at least one stage. |
| Trigger           | An assessment is fully scored, a live-interview outcome is recorded, or the Hiring Manager records a manual decision. |
| Main Flow         | 1. The System applies the stage's configured advancement rule (FR-INVW-051).<br>2. The System records the decision type, decider, time, and score (FR-INVW-058).<br>3. On pass, the System makes linked stages available only when all incoming stages have passed and awards the Amarat associated with the completed assessment (FR-INVW-002, FR-INVW-003, FR-INVW-102).<br>4. On rejection, the System marks the application Rejected, removes the Candidate from all active stages, and starts rejection feedback (FR-INVW-060). |
| Alternative Flows | 1a. Manual rule: the Hiring Manager records pass or reject (FR-INVW-052, FR-INVW-053).<br>1b. Threshold rule: the System passes at or above the stage's current threshold; below-threshold handling is automatic rejection or hold, with hold as the default. Threshold changes from the dashboard are handled as defined in UC-INVW-007 (FR-INVW-054, FR-INVW-055, FR-INVW-096).<br>1c. Completion rule: the System passes when the module is complete, regardless of score (FR-INVW-056).<br>2a. The Hiring Manager overrides an automatic result with a written reason (FR-INVW-057).<br>2b. The Hiring Manager passes or rejects several selected Candidates at the same stage in one action (FR-INVW-059).<br>3a. The Candidate passes the final outstanding assessment stage on every pipeline path: the System identifies the Candidate as a Finalist and adds them to the Finalist list through UC-INVW-017 (FR-INVW-107). |
| Exception Flows   | E1. Two simultaneous decisions concern the same Candidate and stage: the System records the first and rejects the second, showing the recorded decision. |
| Postconditions    | The decision is recorded and the Candidate is advanced, rejected, or held for a Hiring Manager decision. |
| Business Rules    | Only the Hiring Manager records manual decisions; Interviewers submit scorecards but do not make advancement decisions (FR-INVW-053). Stage advancement is independently configured as manual, threshold-based, or on completion (FR-INVW-051). |

#### UC-INVW-007: Review Pipeline Dashboard and Insights

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager reviews Candidate progress, rankings, and system-generated insights to support decisions. |
| Preconditions     | A pipeline is published and at least one Candidate is active in it. |
| Trigger           | The Hiring Manager opens the pipeline dashboard. |
| Main Flow         | 1. The System shows a graphical view with each stage's active, passed, and rejected Candidate counts (FR-INVW-061).<br>2. The Hiring Manager selects a stage; the System lists its Candidates ranked by module score and provides score-range, Nice-to-have, and pending-decision filters (FR-INVW-062).<br>3. The System shows each Candidate's completed-stage scores and percentile (FR-INVW-063), plus a labelled summary of strengths and weaknesses (FR-INVW-064).<br>4. The Hiring Manager may record pass/reject decisions from the dashboard; decision processing follows UC-INVW-006 (FR-INVW-066). |
| Alternative Flows | 2a. The System suggests a cut-off score expected to reach the finalist target and shows the resulting pass count (FR-INVW-068).<br>3a. The Hiring Manager compares up to 4 selected Candidates side by side (FR-INVW-069).<br>3b. After reviewing the score distribution and insights, the Hiring Manager applies the suggested or another revised threshold; the System previews and recalculates outcomes for Candidates still waiting at that stage (FR-INVW-096). |
| Exception Flows   | None. |
| Postconditions    | The Hiring Manager has an updated pipeline view. |
| Business Rules    | Ranking and insight inputs exclude identity and demographic attributes (FR-INVW-065). Insights support but do not replace human decisions except where an automatic advancement rule is configured (FR-INVW-051). |

#### UC-INVW-008: Enable Anonymous Hiring

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Interviewer, Candidate , System |
| Description       | The Hiring Manager hides Candidate identity from Hiring Managers and Interviewers until a chosen pipeline reveal point. |
| Preconditions     | A pipeline is being designed (UC-INVW-001). |
| Trigger           | The Hiring Manager enables anonymity on a pipeline. |
| Main Flow         | 1. The Hiring Manager enables anonymity and selects an identity-reveal point (a pipeline stage or the Finalist list) (FR-INVW-070).<br>2. Before that point, the System substitutes the anonymous identifier for identity fields in every Hiring Manager and Interviewer view (FR-INVW-071).<br>3. The System tells the Candidate that identity is hidden from the employer during early stages without naming the reveal point (FR-INVW-074).<br>4. When the Candidate reaches the reveal point, the System reveals identity permanently and records the reveal time (FR-INVW-073). |
| Alternative Flows | None. |
| Exception Flows   | E1. A path can reach a video-recorded response or live-interview stage at or before the reveal point: the System marks the pipeline invalid (FR-INVW-072).<br>E2. A Hiring Manager or Interviewer requests identity before the reveal point: the System denies access and logs the attempt (FR-INVW-071). |
| Postconditions    | Anonymity is active up to the reveal point; reveal is permanent and recorded. |
| Business Rules    | Identity reveal is permanent (FR-INVW-073). In an anonymous pipeline, every live-interview and video-recorded response stage must follow the reveal point strictly (FR-INVW-072). |

#### UC-INVW-009: Provide Feedback and Notifications

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Candidate , System |
| Description       | The System releases assessment feedback, rejection feedback, and advancement notifications to Candidates by email and in the platform. |
| Preconditions     | The Candidate has completed a fully scored module, advanced, or been rejected; the applicable feedback mode is configured. |
| Trigger           | A module completes, an advancement occurs, or a rejection is recorded. |
| Main Flow         | 1. The System releases AI-generated module feedback within 10 minutes of completion, or human-written/AI-approved feedback within 10 minutes of submission or approval (FR-INVW-076).<br>2. The System emails the Candidate within 10 minutes of a pass to a next stage, naming that stage (FR-INVW-079).<br>3. Within 24 hours of rejection, the System produces and releases feedback with an overall summary, each completed module's outcome and at least one strength or improvement area, and at least one recommended action (FR-INVW-080, FR-INVW-081). |
| Alternative Flows | 1. Feedback mode is human-written or AI-generated with human approval before release. |
| Exception Flows   | E1. AI generation fails: the System holds the feedback item and notifies the Hiring Manager to write it manually. |
| Postconditions    | Feedback and notifications are delivered by email and in the platform. |
| Business Rules    | Feedback excludes questions, answer keys, thresholds, and other Candidates' data (FR-INVW-077); AI-generated feedback is labelled (FR-INVW-078). The Hiring Manager configures the feedback modes in FR-INVW-075 and FR-INVW-081. |

#### UC-INVW-010: Spawn Test Instance

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Candidate, Hiring Manager, System |
| Description       | The System creates and times an individual assessment instance; the Candidate submits responses and the System scores it after any required manual grading. |
| Preconditions     | The Candidate has reached an individual assessment stage in the published pipeline and its scheduled interval has begun; its module content and scoring method are valid (UC-INVW-002, UC-INVW-003, FR-INVW-082, FR-INVW-094). Live interviews and multi-person tests use UC-INVW-004 and UC-INVW-012, respectively. |
| Trigger           | The scheduled start time has arrived and the Candidate is eligible for the assessment stage (FR-INVW-094). |
| Main Flow         | 1. The System creates the Candidate's instance from the published module configuration and preserves its selected content for subsequent opens (FR-INVW-083); random-selection stages give every Candidate the same problem count and difficulty mix (FR-INVW-010).<br>2. The Candidate's completion window starts when the stage becomes available; the active timer starts when the Candidate first opens the instance, displays remaining time, and continues across navigation or disconnection (FR-INVW-011, FR-INVW-084). Responses are saved and restored according to NFR-REL-008.<br>3. If the instance remains unfinished, the System sends reminders 48, 24, 12, 6, and 1 hours before the applicable deadline—the earlier of the completion-window deadline and assessment end time, as extended for the Candidate under FR-INVW-015 when applicable (FR-INVW-013).<br>4. The Candidate submits responses once; the System locks them and also locks the latest saved responses when auto-submitting at the active time limit or applicable deadline (FR-INVW-012, FR-INVW-085, FR-INVW-090).<br>5. The System automatically scores closed questions from answer keys (FR-INVW-020), sends coding solutions and test cases to the configured external code-evaluation provider (FR-INVW-087; provider interface details are tracked in `docs/srs/TBD.md`, TBD-004), and routes open-ended written answers and video responses to assigned human graders (FR-INVW-021, FR-INVW-093). AI-assisted question drafting does not score Candidate responses (FR-INVW-018, FR-INVW-019).<br>6. The System keeps the instance pending until all required human grading is complete, then records the final 0–100 score (FR-INVW-016, FR-INVW-086, FR-INVW-089).<br>7. Once fully scored, stage advancement follows UC-INVW-006 and module feedback follows UC-INVW-009. |
| Alternative Flows | 2a. The Candidate leaves and resumes the same instance; saved answers are restored and the active timer continues (FR-INVW-083, FR-INVW-084, NFR-REL-008).<br>2b. The Hiring Manager grants additional active time or a later individual completion deadline, including beyond the scheduled stage end for this started instance, and records the change and reason (FR-INVW-015).<br>4a. In a question-based module, the System presents each configured question using its response format and scoring method (FR-INVW-091).<br>4b. In a video-recorded response module, the Candidate records and submits a video response for each prompt; the System attaches each recording to the instance for human grading (FR-INVW-093). |
| Exception Flows   | E1. The earlier of the completion-window deadline and scheduled assessment end time arrives before an attempt has started or been submitted: the System marks the instance missed, notifies the Hiring Manager, and holds the Candidate for a decision (FR-INVW-014).<br>E2. The external code-evaluation provider is unavailable or returns an invalid result: the System preserves the code submission, leaves the instance pending, and does not assign a zero score (FR-INVW-087). |
| Postconditions    | The instance is fully scored and ready for advancement, remains pending manual grading, or is marked missed; advancement and feedback follow UC-INVW-006 and UC-INVW-009. |
| Business Rules    | Scoring and advancement are distinct: answer-key questions are auto-scored (FR-INVW-020), coding tests use external provider results (FR-INVW-087), and open-ended/video/live-interview responses use human grading or scorecards (FR-INVW-021, FR-INVW-034, FR-INVW-093). No assessment score or pass/reject decision is assigned by AI; AI question drafting, insights, and feedback do not assign scores (FR-INVW-018, FR-INVW-064, FR-INVW-076). The System withholds answer keys and rubrics from the active Candidate view (FR-INVW-088). A test instance is individual to one Candidate and stage assignment (FR-INVW-083). |

#### UC-INVW-011: Detect Cheating on Test Instances

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Priority          | Should (long-term goal, not in the MVP) |
| Actor(s)          | Candidate, Hiring Manager, System |
| Description       | The system monitors spawned test instances for signs of cheating and flags suspicious instances with evidence for the Hiring Manager to review; a flag never decides an outcome by itself. |
| Preconditions     | Test instances exist for the stage (UC-INVW-010). |
| Trigger           | A test instance is started or completed. |
| Main Flow         | 1. During an assessment, the System records configured session events.<br>2. After completion, the System compares answers and timing across Candidates in the same stage.<br>3. The System flags suspicious instances and records the evidence that raised each flag.<br>4. The Hiring Manager reviews flags and records any advancement decision through UC-INVW-006. |
| Alternative Flows | 1. The Hiring Manager dismisses a flag with a written reason; the dismissal is recorded with the flag. |
| Exception Flows   | E1. The monitoring data is missing or inconclusive: no flag is raised.<br>E2. The session ends abnormally: the events recorded up to that point are kept, and missing data alone does not raise a flag. |
| Postconditions    | Suspicious instances carry their flags and evidence; only a Hiring Manager decision changes a Candidate's status (UC-INVW-006). |
| Business Rules    | Flags support but never replace human decisions; only the Hiring Manager records manual decisions (FR-INVW-053). Flags and evidence are shown only to the Hiring Manager. Functional requirements for this long-term goal are to be added when it is scheduled. |

#### UC-INVW-012: Run Multi-Person Tests

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Priority          | Could (long-term goal, not in the MVP) |
| Actor(s)          | Hiring Manager, Candidates, System |
| Description       | A stage's test module is taken by several Candidates together (a pair or a team) in one shared session, with an individual score recorded for each member. |
| Preconditions     | A pipeline design exists for the posting (UC-INVW-001). The stage is configured with a multi-person module type and a team size. |
| Trigger           | The first Candidate reaches a stage configured with a multi-person module. |
| Main Flow         | 1. The System groups Candidates who reach the stage into pairs or teams of the configured size.<br>2. The System creates one shared test instance for each group from the stage's configured module.<br>3. Group members complete the test together within the configured completion window and active time limit.<br>4. The System records a 0–100 score for each member (FR-INVW-016).<br>5. Advancement of each member follows UC-INVW-006. |
| Alternative Flows | 1. Pair type: groups of two.<br>2. Team type: groups of the configured size, optionally with a team-level result in addition to the individual scores. |
| Exception Flows   | E1. A group member never joins or drops out: the session continues for the remaining members, and the absent member is held for the Hiring Manager's decision.<br>E2. The time limit expires: the System auto-submits the group's answers (FR-INVW-012).<br>E3. Not enough Candidates remain to complete a group: the System holds the ungrouped Candidates for the Hiring Manager's decision. |
| Postconditions    | Each group member has a completed module score or a held status; advancement follows UC-INVW-006. |
| Business Rules    | Multi-person modules are incompatible with anonymity before the reveal point. They pair Candidates with one another and are distinct from the future LLM-paired programming module (UC-INVW-013). Functional requirements for this long-term goal are to be added when it is scheduled. |

#### UC-INVW-013: Use LLM-Paired Programming Module

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Priority          | Could (future feature; not in the current release) |
| Actor(s)          | Candidate, Hiring Manager, System |
| Description       | The Candidate works on a configured coding task with an LLM partner; the LLM assists during the session but does not assess or score the Candidate. |
| Preconditions     | The pipeline includes a configured LLM-paired programming module (FR-INVW-092), and the Candidate has reached the stage. |
| Trigger           | The Candidate opens the LLM-paired programming test instance. |
| Main Flow         | 1. The System presents the configured coding task and starts the active time limit.<br>2. The Candidate exchanges text with the LLM and edits a code solution during the session.<br>3. The System saves the conversation and final code with the test instance (FR-INVW-092).<br>4. The System sends the final code to the configured external code-evaluation provider; scoring uses its valid test-case results (FR-INVW-087). Provider interface details are tracked in `docs/srs/TBD.md` (TBD-004). |
| Alternative Flows | None. |
| Exception Flows   | E1. The LLM service is unavailable: the System preserves the Candidate's existing code and reports that the LLM partner is unavailable.<br>E2. The external code-evaluation provider is unavailable or returns an invalid result: the System preserves the final code and leaves the instance pending without assigning a score (FR-INVW-087). |
| Postconditions    | The conversation and final code are saved; the code is scored by the external provider, not by the LLM. |
| Business Rules    | The LLM assists the Candidate and does not make a pass/reject decision or produce the assessment score. |

#### UC-INVW-014: View Candidate Pipeline Progress

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Candidate, System |
| Description       | The Candidate views the complete published pipeline and their own progress, including their current stage or stages. |
| Preconditions     | The Candidate's application is assigned to a published INVW pipeline after initial screening by the separate screening epic. |
| Trigger           | The Candidate opens the pipeline progress view for their application. |
| Main Flow         | 1. The System loads the published pipeline version assigned to the Candidate's application (FR-INVW-041).<br>2. The System displays the full stage graph, including every stage name and directed link in that version.<br>3. The System highlights the Candidate's current stage or stages and marks the Candidate's other stages as completed, active, or upcoming/locked (FR-INVW-097).<br>4. The System displays the Amarat earned at passed stages and identifies the Company associated with each Amara award (FR-INVW-103, FR-INVW-106). The Candidate may open a detailed skill card when that SHOULD capability is available (FR-INVW-103). |
| Alternative Flows | 1a. After a pipeline migration, the System shows the version the Candidate is now running and preserves completed-stage progress (FR-INVW-041, FR-INVW-045). |
| Exception Flows   | None. |
| Postconditions    | The Candidate can identify their current position and the remaining stages in the published pipeline. |
| Business Rules    | The view contains only the Candidate's own progress; it does not show other Candidates' data, Hiring Manager notes, or advancement thresholds (FR-INVW-097). |

#### UC-INVW-015: Configure Amara Mappings and Levels

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager associates a skill crest, called an Amara (plural: Amarat), with the scored questions or rubric criteria of an assessment owned by the posting's Company. |
| Preconditions     | An assessment module or live-interview rubric is being configured for a posting owned by the Hiring Manager's Company. |
| Trigger           | The Hiring Manager opens Amara configuration for an assessment. |
| Main Flow         | 1. For a predefined assessment module, the System supplies its built-in Amara definition (FR-INVW-098).<br>2. For a custom test or live-interview rubric, the Hiring Manager selects a built-in Amara or creates a custom one with a unique name, crest, skill domain, and description; each custom definition is associated with the posting's Company (FR-INVW-099, FR-INVW-106).<br>3. The Hiring Manager maps every scored question or rubric criterion to exactly one Amara and ensures each selected Amara has at least one mapped item (FR-INVW-099).<br>4. The System saves the Amara mappings with the assessment configuration. The Hiring Manager may configure levels for custom Amarat (FR-INVW-100). |
| Alternative Flows | 4a. When level scoring is configured, the Hiring Manager defines ordered levels, score ranges, descriptions, strengths, and areas for development; the System validates that ranges do not overlap and cover 0–100 (FR-INVW-100). |
| Exception Flows   | E1. A scored item or rubric criterion is unmapped, or an Amara has no mapped item: the System identifies the problem and blocks publication (FR-INVW-099).<br>E2. A configured level scale overlaps or does not cover 0–100: the System identifies the problem and does not save the invalid scale (FR-INVW-100). |
| Postconditions    | The assessment is saved with Amara definitions associated with the Company and with complete evidence mappings; any configured level scale is valid. |
| Business Rules    | Amara score bands are distinct from stage pass thresholds (FR-INVW-100, FR-INVW-054). Level configuration is optional; a Candidate may earn an Amara without a numeric level (FR-INVW-102). |

#### UC-INVW-016: Award and Review Amarat

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Candidate, Hiring Manager, System |
| Description       | When a Candidate passes an assessment, the System awards the associated Amarat, retains their Company association, and makes them available for review. |
| Preconditions     | The Candidate has passed a stage with Amara mappings, and the pass decision is recorded. Numeric Amara scoring may be unavailable. |
| Trigger           | The System records a Candidate's pass decision for an Amara-mapped assessment stage (UC-INVW-006). |
| Main Flow         | 1. If numeric Amara scoring and levels are configured, the System may calculate a separate score for each Amara from only its mapped evidence and assign the matching level (FR-INVW-101).<br>2. The System awards each Amara assigned to the passed stage and retains the Candidate, Amara definition, source stage, pass decision, and award time (FR-INVW-102, FR-INVW-104). Each award is associated with the Company that owns the assessment posting (FR-INVW-106).<br>3. The Candidate can see the awarded crest and its Company. When detailed skill-card functionality is available, the Candidate can view its level, description, strengths, areas for development, and mapped results (FR-INVW-103).<br>4. If the Candidate earns the same Amara again, the System retains a separate award history entry; it may show the highest earned level and preserve the definition/level-scale version for each award (FR-INVW-104, FR-INVW-105). |
| Alternative Flows | 1a. No numeric Amara score or level is available: the System still awards the Amara without a level (FR-INVW-102).<br>4a. A later award is at a lower level: the System retains the new award event; if highest-level display is supported, it continues to show the highest earned level (FR-INVW-104, FR-INVW-105). |
| Exception Flows   | E1. The Candidate is held, rejected, or still awaiting required grading: the System does not issue an Amara award (FR-INVW-089, FR-INVW-102). |
| Postconditions    | The Amara award and its evidence are recorded and available to the Candidate in the progress view and Amara Dynamic Resume profile. |
| Business Rules    | Amara descriptions refer only to assessed skills and mapped evidence, not inferred personality traits (FR-INVW-103). A Company Hiring Manager can access awards only for Candidates/applications they are authorized to manage (FR-INVW-106, AUTH). |

#### UC-INVW-017: Manage Finalists and Offers

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Candidate, System |
| Description       | The System identifies Candidates who pass the complete interview pipeline as Finalists. The Hiring Manager reviews Finalists, makes and manages offers against the posting's open positions, and rejects Finalists who will not receive an offer. An offer that exceeds capacity requires explicit confirmation. Candidates respond to offers in the platform. |
| Preconditions     | The Hiring Manager is authorized to manage the posting (AUTH). The posting has a published pipeline and a configured open-position capacity owned by JOB. At least one Candidate has passed every assessment stage on every path of that Candidate's pipeline. |
| Trigger           | The System records a pass for a Candidate's final outstanding assessment stage, or the Hiring Manager opens the Finalist list. |
| Main Flow         | 1. The System sets each eligible Candidate's application status to "Finalist" and adds the Candidate to the posting's Finalist list (FR-INVW-107).<br>2. The System orders the list by descending mean score across each Finalist's completed assessment modules and live interviews with numeric scores, displays the Candidate insights, and allows the Hiring Manager to order by a completed stage's score (FR-INVW-108).<br>3. The Hiring Manager selects a Finalist and creates an offer with the posting's job title, compensation summary, proposed start date, response deadline, and optional message (FR-INVW-109).<br>4. The System verifies that the response deadline is in the future, that the Candidate has no other pending offer for this posting, and that the posting's capacity value is available from JOB (FR-INVW-116, FR-INVW-118).<br>5. The System compares pending and accepted offers, including the proposed offer, with the posting's open-position capacity (FR-INVW-116). If the Hiring Manager confirms any over-capacity warning, the System creates the offer and sets the offer and application status to "Offer Pending", then delivers it by email and in the platform (FR-INVW-109, FR-INVW-110).<br>6. The Candidate accepts or declines the pending offer in the platform before the response deadline. The System records the response and timestamp, updates the offer and application status, and notifies the Hiring Manager within 10 minutes (FR-INVW-111, FR-INVW-112). |
| Alternative Flows | 5a. The proposed offer would exceed capacity: the System warns the Hiring Manager and proceeds only after explicit confirmation (FR-INVW-116).<br>5b. The Hiring Manager cancels or does not confirm an over-capacity offer: the System does not create or send that offer (FR-INVW-116).<br>6a. The Candidate does not respond by the deadline: the System marks the offer and application "Offer Expired" and notifies the Candidate and Hiring Manager (FR-INVW-113).<br>6b. Before the deadline, the Hiring Manager extends the deadline of a pending offer; the System saves the later deadline (FR-INVW-114).<br>6c. Before the deadline, the Hiring Manager withdraws a pending offer; the System marks the offer "Withdrawn", returns the application to "Finalist", and notifies the Candidate (FR-INVW-115).<br>3a. The Hiring Manager rejects one or more selected Finalists individually or in bulk; the System sets each application to "Rejected" and starts rejection feedback (FR-INVW-117). |
| Exception Flows   | E1. Required offer details are missing or the response deadline is not later than the current time: the System identifies the invalid or missing fields and does not create the offer (FR-INVW-109).<br>E2. The Candidate already has a pending offer for the posting: the System blocks creation and shows the existing offer (FR-INVW-118).<br>E3. The response deadline has passed: the System rejects a late Candidate response or Hiring Manager deadline/withdrawal change and expires the offer (FR-INVW-111, FR-INVW-113, FR-INVW-114, FR-INVW-115).<br>E4. The posting's open-position capacity cannot be read: the System does not send the offer and identifies the capacity issue. |
| Postconditions    | Each application reflects the Candidate's finalist, offer, or rejection state. Offer content, delivery, deadline, response or expiry, and any withdrawal are recorded. Rejected Finalists enter the configured rejection-feedback flow. |
| Business Rules    | Only Candidates who pass every assessment stage on every path qualify as Finalists (FR-INVW-107). The JOB finalist target is not a cap on the Finalist list. Capacity is configured by JOB and read by INVW. The default Finalist mean uses available numeric scores from completed assessment modules and live interviews; Candidates with no numeric scores are shown without a mean (FR-INVW-108). A Candidate may have at most one pending offer per posting (FR-INVW-118). An over-capacity offer proceeds only after explicit Hiring Manager confirmation (FR-INVW-116). |

---

### 4.5 JOB Use Cases

#### UC-JOB-01: Create, Configure and Edit Job Posting

| Field | Description |
|---|---|
| Use Case ID | UC-JOB-01 |
| Use Case Name | Create, Configure and Edit Job Posting |
| Actor(s) | Hiring Manager, System |
| Description | A permitted Hiring Manager creates or edits a posting and configures its hiring criteria. |
| Preconditions | The company is verified and the Hiring Manager is authenticated and authorized. |
| Trigger | The Hiring Manager selects Create or Edit. |
| Main Flow | The Hiring Manager opens the editor; enters the job title, description, positions, language, optional salary and currency, dates, tags, criteria and application questions; the system validates and saves the posting as Draft or updates the existing posting. |
| Alternative Flows | The system may autosave a Draft, warn about unsaved edits, or duplicate an existing posting into a new Draft. |
| Exception Flows | Invalid dates, missing mandatory data, unauthorized changes or attempts to modify locked Must-have criteria are rejected. |
| Postconditions | The posting is saved or remains unchanged. |
| Related Requirements | FR-JOB-001, FR-JOB-002, FR-JOB-003, FR-JOB-004, FR-JOB-005, FR-JOB-006, FR-JOB-025, FR-JOB-026, FR-JOB-032, FR-JOB-034, FR-JOB-036, FR-JOB-037, FR-JOB-038, FR-JOB-039, FR-JOB-040 |

#### UC-JOB-02: Preview Posting

| Field | Description |
|---|---|
| Use Case ID | UC-JOB-02 |
| Use Case Name | Preview Posting |
| Actor(s) | Hiring Manager, System |
| Description | The Hiring Manager previews the candidate-facing posting. |
| Preconditions | The Hiring Manager is authorized to manage the posting. |
| Trigger | The Hiring Manager selects Preview. |
| Main Flow | The system displays the posting as candidates would see it, including unsaved edits. |
| Alternative Flows | A Draft may be previewed before publication. |
| Exception Flows | If required fields are missing, the preview may identify them. |
| Postconditions | No posting or application data is changed. |
| Related Requirements | FR-JOB-007 |

#### UC-JOB-03: Publish Job Posting

| Field | Description |
|---|---|
| Use Case ID | UC-JOB-03 |
| Use Case Name | Publish Job Posting |
| Actor(s) | Hiring Manager, System |
| Description | The Hiring Manager publishes a complete Draft. |
| Preconditions | The posting is Draft and the company is not banned. |
| Trigger | The Hiring Manager selects Publish. |
| Main Flow | The system validates required fields and dates; changes the state to Published; makes the posting available to candidate discovery and the public company profile. |
| Alternative Flows | The posting may be published before its interview pipeline is available. |
| Exception Flows | If validation fails, the system lists missing or invalid fields and keeps the posting in Draft. Banned companies cannot publish. |
| Postconditions | The posting is Published or remains Draft. |
| Related Requirements | FR-JOB-011, FR-JOB-012 |

#### UC-JOB-04: Archive and Republish Posting

| Field | Description |
|---|---|
| Use Case ID | UC-JOB-04 |
| Use Case Name | Archive and Republish Posting |
| Actor(s) | Hiring Manager, System |
| Description | The Hiring Manager archives or republishes a posting. |
| Preconditions | The Hiring Manager is authorized; the posting is Published or Archived. |
| Trigger | The Hiring Manager selects Archive or Republish. |
| Main Flow | For Archive, the system removes the posting from discovery, stops new applications and retains existing data. For Republish, it restores Published status and resumes suspended applications. |
| Alternative Flows | The system should display affected application counts and request confirmation before archiving. |
| Exception Flows | A republished posting past its deadline does not accept new applications. |
| Postconditions | The posting state and application availability are updated. |
| Related Requirements | FR-JOB-014, FR-JOB-015, FR-JOB-016 |

#### UC-JOB-05: Delete Posting

| Field | Description |
|---|---|
| Use Case ID | UC-JOB-05 |
| Use Case Name | Delete Posting |
| Actor(s) | Hiring Manager, Admin, System |
| Description | An authorized Hiring Manager or Admin permanently deletes a posting. |
| Preconditions | The actor is authorized to delete the posting. |
| Trigger | The actor selects Delete. |
| Main Flow | The system displays a permanent-deletion warning, shows the number of affected applications when applicable, requests confirmation and deletes the posting. |
| Alternative Flows | The actor can cancel without changing the posting. |
| Exception Flows | An unauthorized deletion is denied. |
| Postconditions | The posting is deleted or left unchanged. |
| Related Requirements | FR-JOB-021, FR-JOB-022 |

#### UC-JOB-06: Handle Postings of a Banned Company

| Field | Description |
|---|---|
| Use Case ID | UC-JOB-06 |
| Use Case Name | Handle Postings of a Banned Company |
| Actor(s) | Admin, Hiring Manager, System |
| Description | The system handles postings when a company is banned or reinstated. |
| Preconditions | The company has postings and its ban status changes. |
| Trigger | An Admin bans or reinstates the company. |
| Main Flow | On ban, the system archives Published postings and prevents publishing. On reinstatement, it republishes postings archived due to the ban and resumes their applications. |
| Alternative Flows | Existing Drafts and postings archived before the ban retain their states. |
| Exception Flows | A banned company cannot publish or republish postings. |
| Postconditions | Posting states reflect the company ban status. |
| Related Requirements | FR-JOB-018, FR-JOB-019, FR-JOB-020 |

#### UC-JOB-07: Manage Company Posting List

| Field | Description |
|---|---|
| Use Case ID | UC-JOB-07 |
| Use Case Name | Manage Company Posting List |
| Actor(s) | Hiring Manager, System |
| Description | The Hiring Manager views and manages the company posting list. |
| Preconditions | The Hiring Manager is authorized. |
| Trigger | The Hiring Manager opens the posting list. |
| Main Flow | The system displays company postings with title, state, application count and deadline. |
| Alternative Flows | The system may allow searching by title and filtering by state. |
| Exception Flows | If the list is unavailable, the system displays an error. |
| Postconditions | The Hiring Manager can review the posting list. |
| Related Requirements | FR-JOB-023, FR-JOB-024 |

#### UC-JOB-08: Switch Employer Interface Language

| Field | Description |
|---|---|
| Use Case ID | UC-JOB-08 |
| Use Case Name | Switch Employer Interface Language |
| Actor(s) | Hiring Manager, System |
| Description | The Hiring Manager optionally changes the employer interface language. |
| Preconditions | The Hiring Manager is signed in and interface language switching is supported. |
| Trigger | The Hiring Manager selects a language. |
| Main Flow | The system switches the employer interface to the selected supported language. |
| Alternative Flows | The system may remember the language for future sessions. |
| Exception Flows | An unsupported language cannot be selected. |
| Postconditions | The employer interface uses the chosen language without changing posting content. |
| Related Requirements | FR-JOB-009, FR-JOB-010 |

#### UC-JOB-09: Apply to Posting

| Field | Description |
|---|---|
| Use Case ID | UC-JOB-16 |
| Use Case Name | Apply to Posting |
| Actor(s) | Candidate, System |
| Description | The candidate submits an application to an open posting. |
| Preconditions | The candidate is signed in, has a complete resume and the posting is open. |
| Trigger | The candidate selects Apply. |
| Main Flow | The candidate answers required questions and submits; the system checks eligibility, posting dates, application cap and duplicates; saves an immutable copy of the resume and answers; and passes the accepted application to screening. |
| Alternative Flows | The system may warn the candidate about unmet Must-have criteria before submission and should display an application reference after acceptance. |
| Exception Flows | Missing answers, duplicate applications, closed or archived postings, and reached caps cause rejection with the relevant reason. |
| Postconditions | An accepted application is stored unchanged or an invalid submission is rejected. |
| Related Requirements | FR-JOB-041, FR-JOB-042, FR-JOB-043, FR-JOB-045, FR-JOB-047 |

#### UC-JOB-10: Screen and Admit Applications

| Field | Description |
|---|---|
| Use Case ID | UC-JOB-17 |
| Use Case Name | Screen and Admit Applications |
| Actor(s) | System, Hiring Manager |
| Description | The system screens applications against Must-have criteria and admits eligible applications. |
| Preconditions | An application has been accepted. |
| Trigger | The system receives an accepted application. |
| Main Flow | The system evaluates the submitted structured resume against every Must-have criterion; admits applications meeting all criteria to INVW; otherwise screens them out and records failed criteria. |
| Alternative Flows | The Hiring Manager should be able to review screened-out applications; the candidate may be informed of failed criteria. |
| Exception Flows | Missing structured fields fail the relevant Must-have check. |
| Postconditions | The application is admitted or screened out, with the outcome recorded. |
| Related Requirements | FR-JOB-051 |

### 4.6 ANA Use Cases

#### UC-ANA-001: Stream & Ingest IDE Telemetry

| Field | Description |
| --- | --- |
| **Actor(s)** | Candidate, System |
| **Description** | Stream and analyze telemetry data coming from the IDE that the candidate is using. |
| **Preconditions** | Candidate active session in the IDE environment; telemetry endpoint is established. |
| **Trigger** | An action is triggered through the IDE. |
| **Main Flow** | 1. The IDE captures interaction events in real-time.<br>2. Events are buffered and streamed to the telemetry endpoint. |
| **Alternative Flows** | None. |
| **Exception Flows** | Telemetry connection drops: IDE plugin queues events locally in buffer until connection is restored. |
| **Postconditions** | Telemetry events successfully connected to the candidate record. |
| **Business Rules** | The telemetry streaming must run flawlessly and not affect the candidate’s experience. |

#### UC-ANA-002: Calculate Execution Efficiency Metrics

| Field | Description |
| --- | --- |
| **Actor(s)** | Candidate, System |
| **Description** | Measures the performance related metrics such as Run Time, Memory Consumption, and complexity. |
| **Preconditions** | Code submitted to the sandbox environment. |
| **Trigger** | The candidate submits code solutions to the sandbox. |
| **Main Flow** | 1. Sandbox executes submitted code against test suites.<br>2. Analytics subsystem calculates execution runtime and peak memory usage.<br>3. Subsystem calculates O-complexity approximation.<br>4. Execution metrics are saved to the candidate record. |
| **Alternative Flows** | None. |
| **Exception Flows** | Execution timeout: System records partial diagnostic telemetry and returns failure flags. |
| **Postconditions** | The metrics are saved to the candidate’s records. |
| **Business Rules** | Metrics must be computed in an isolated execution sandbox environment. |

#### UC-ANA-003: Calculate Technical Pass Rate Score

| Field | Description |
| --- | --- |
| **Actor(s)** | System |
| **Description** | Compute candidate technical Pass Rate Score based on automated test suite outcomes. |
| **Preconditions** | Test suite execution completed in a sandbox environment. |
| **Trigger** | Code execution finishes in the sandbox. |
| **Main Flow** | 1. Subsystem retrieves passed and failed test case counts.<br>2. Computes the Pass Rate Score percentage.<br>3. Updates technical scorecard record. |
| **Alternative Flows** | None. |
| **Exception Flows** | Unhandled code runtime crash: Pass Rate Score recorded as 0% with diagnostic failure details. |
| **Postconditions** | Technical Pass Rate Score updated on candidate evaluation scorecard. |
| **Business Rules** | A zero pass rate score is not an auto rejection. |

#### UC-ANA-004: Evaluate Integrity Risk

| Field | Description |
| --- | --- |
| **Actor(s)** | System |
| **Description** | Analyzing the telemetry events, sequence, and frequency and generating an integrity score based on it. |
| **Preconditions** | Telemetry streaming active and integrity threshold rules configured. |
| **Trigger** | There is detected focus loss or paste action in the telemetry stream. |
| **Main Flow** | 1. Subsystem detects violation in incoming stream.<br>2. Calculates/updates Integrity Risk Score (0–100).<br>3. Attaches evidence log. |
| **Alternative Flows** | No violation is detected. |
| **Exception Flows** | Issues with the telemetry stream. |
| **Postconditions** | Evidence logs attached and candidate Integrity Risk Score updated. |
| **Business Rules** | A high integrity score requires human revision to avoid auto rejection. |

#### UC-ANA-005: Export Candidate Evaluation Report

| Field | Description |
| --- | --- |
| **Actor(s)** | Hiring Manager, System |
| **Description** | Consolidate technical, behavioral, and integrity metrics into a unified exportable candidate report (PDF). |
| **Preconditions** | Technical, behavioral, and integrity score calculations completed. |
| **Trigger** | Hiring Manager requests report export. |
| **Main Flow** | 1. Subsystem aggregates technical pass rate, execution metrics, and integrity risk score.<br>2. Renders report in requested export format (PDF).<br>3. Delivers file for download or platform display. |
| **Alternative Flows** | None. |
| **Exception Flows** | PDF Generation failure: System retries up to 3 times before queuing error notification. |
| **Postconditions** | Exported PDF candidate evaluation report generated. |
| **Business Rules** | None |