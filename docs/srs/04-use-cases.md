## 4. Use cases

Owner: All | Status: Draft | Last updated: 2026-10-13 | Jira: AMARA-32

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
| Actor(s) | Candidate, Employer |
| Description | Allows candidates to view the final hiring decision submitted by the employer and any accompanying employer feedback or AI-generated feedback in case of rejection. |
| Preconditions | The candidate is authenticated and has an application associated with an employer. |
| Trigger | The employer submits a final hiring decision for the application. |
| Main Flow | 1. The employer submits the final hiring decision.<br>2. The system associates the decision with the corresponding application.<br>3. If the decision is a rejection, the system generates brief AI feedback for the candidate.<br>4. The system makes the decision and any applicable feedback available to the candidate.<br>5. The candidate opens the application.<br>6. The system displays the final hiring decision.<br>7. The system displays any final feedback provided by the employer.<br>8. If the candidate was rejected, the system displays the AI-generated feedback. |
| Alternative Flows | 1. If the employer submits a decision without additional feedback, the system displays the decision alone, along with AI-generated feedback if the decision is a rejection.<br>2. If the candidate does not immediately open the application, the decision and its available feedback remain accessible for later viewing. |
| Exception Flows | 1. If the employer has not submitted a final decision, the system continues to display the current application status.<br>2. If the decision or its associated feedback cannot be retrieved, the system displays an appropriate error message. |
| Postconditions | The candidate can access the final hiring decision and any accompanying employer feedback or AI-generated rejection feedback. |
| Business Rules | The final decision shall be associated with the corresponding application and shall only be displayed once submitted by the employer. AI-generated feedback shall be provided in case of rejection. Employer-provided final feedback shall be displayed when available. |

#### 4.3.8 UC-CAND-008 — Accept Offer

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to accept or refuse an offer associated with an application. |
| Preconditions | The candidate is authenticated and has received an offer associated with an application. |
| Trigger | The candidate opens the offer details for an application. |
| Main Flow | 1. The candidate opens the application containing the offer.<br>2. The system displays the offer and the available response options.<br>3. The candidate chooses to accept or refuse the offer.<br>4. The system records the candidate's response.<br>5. The system reflects acceptance of the offer in the candidate's profile when the offer is accepted. |
| Alternative Flows | 1. If the candidate refuses the offer, the system records the refusal and updates the relevant application information. |
| Exception Flows | 1. If the response cannot be recorded, the system displays an appropriate error message. |
| Postconditions | The candidate's offer response is recorded and the candidate's profile reflects acceptance when applicable. |
| Business Rules | The system shall associate the candidate's response with the corresponding offer and application. |

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
| Main Flow         | 1. The System creates the Candidate's instance from the published module configuration and preserves its selected content for subsequent opens (FR-INVW-083); random-selection stages give every Candidate the same problem count and difficulty mix (FR-INVW-010).<br>2. The Candidate's completion window starts when the stage becomes available; the active timer starts when the Candidate first opens the instance, displays remaining time, and continues across navigation or disconnection (FR-INVW-011, FR-INVW-084). Responses are saved and restored according to NFR-REL-008.<br>3. If the instance remains unfinished, the System sends reminders 48, 24, 12, 6, and 1 hours before the applicable deadline—the earlier of the completion-window deadline and assessment end time, as extended for the Candidate under FR-INVW-015 when applicable (FR-INVW-013).<br>4. The Candidate submits responses once; the System locks them and also locks the latest saved responses when auto-submitting at the active time limit or applicable deadline (FR-INVW-012, FR-INVW-085, FR-INVW-090).<br>5. The System automatically scores closed questions from answer keys (FR-INVW-020), sends coding solutions and test cases to the configured external code-evaluation provider (FR-INVW-087; see `docs/srs/06-external-interfaces.md`), and routes open-ended written answers and video responses to assigned human graders (FR-INVW-021, FR-INVW-093). AI-assisted question drafting does not score Candidate responses (FR-INVW-018, FR-INVW-019).<br>6. The System keeps the instance pending until all required human grading is complete, then records the final 0–100 score (FR-INVW-016, FR-INVW-086, FR-INVW-089).<br>7. Once fully scored, stage advancement follows UC-INVW-006 and module feedback follows UC-INVW-009. |
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
| Main Flow         | 1. The System presents the configured coding task and starts the active time limit.<br>2. The Candidate exchanges text with the LLM and edits a code solution during the session.<br>3. The System saves the conversation and final code with the test instance (FR-INVW-092).<br>4. The System sends the final code to the external code-evaluation provider specified in `docs/srs/06-external-interfaces.md`; scoring uses its valid test-case results (FR-INVW-087). |
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

#### Use Case Index

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

#### UC-JOB-01: Create, Configure and Edit Job Posting

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

#### UC-JOB-02: Preview Posting

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

#### UC-JOB-03: Publish Job Posting

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

#### UC-JOB-04: Archive and Republish Posting

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

#### UC-JOB-05: Delete Posting and Manage Postings as Admin

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

#### UC-JOB-06: Handle Postings of a Banned Company

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

#### UC-JOB-07: Manage Company Posting List

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

#### UC-JOB-08: Switch Employer Interface Language

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

#### UC-JOB-09: Browse, Filter, Sort, and Search Jobs and Companies

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

#### UC-JOB-10: View Posting Details and Company Profile

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

#### UC-JOB-11: Share and Open Posting Link

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

#### UC-JOB-12: Manage Saved Postings

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

#### UC-JOB-13: Follow and Manage Companies

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

#### UC-JOB-14: Get Recommendations

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

#### UC-JOB-15: Manage Saved Searches and New-Jobs Tracker

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

#### UC-JOB-16: Apply to Posting

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

#### UC-JOB-17: Screen and Admit Applications

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

### 4.6 ANA Use Cases

#### Use Case 1: Stream & Ingest IDE Telemetry

| Field | Description |
| --- | --- |
| **Use Case ID** | UC-ANA-01 |
| **Use Case Name** | Stream & Ingest IDE Telemetry |
| **Actor(s)** | Candidate, System |
| **Description** | Stream and analyze telemetry data coming from the IDE that the candidate is using. |
| **Preconditions** | Candidate active session in the IDE environment; telemetry endpoint is established. |
| **Trigger** | An action is triggered through the IDE. |
| **Main Flow** | 1. The IDE captures interaction events in real-time.<br>2. Events are buffered and streamed to the telemetry endpoint. |
| **Alternative Flows** | None. |
| **Exception Flows** | Telemetry connection drops -> IDE plugin queues events locally in buffer until connection is restored. |
| **Postconditions** | Telemetry events successfully connected to the candidate record. |
| **Business Rules** | The telemetry streaming must run flawlessly and not affect the candidate’s experience. |

#### Use Case 2: Calculate Execution Efficiency Metrics

| Field | Description |
| --- | --- |
| **Use Case ID** | UC-ANA-02 |
| **Use Case Name** | Calculate Efficiency Metrics |
| **Actor(s)** | Candidate, System |
| **Description** | Measures the performance related metrics such as Run Time, Memory Consumption, and complexity. |
| **Preconditions** | Code submitted to the sandbox environment. |
| **Trigger** | The candidate submits code solutions to the sandbox. |
| **Main Flow** | 1. Sandbox executes submitted code against test suites.<br>2. Analytics subsystem calculates execution runtime and peak memory usage.<br>3. Subsystem calculates O-complexity approximation.<br>4. Execution metrics are saved to the candidate record. |
| **Alternative Flows** | None. |
| **Exception Flows** | Execution timeout -> System records partial diagnostic telemetry and returns failure flags. |
| **Postconditions** | The metrics are saved to the candidate’s records. |
| **Business Rules** | Metrics must be computed in an isolated execution sandbox environment. |

#### Use Case 3: Calculate Technical Pass Rate Score

| Field | Description |
| --- | --- |
| **Use Case ID** | UC-ANA-03 |
| **Use Case Name** | Calculate Technical Pass Rate Score |
| **Actor(s)** | System |
| **Description** | Compute candidate technical Pass Rate Score based on automated test suite outcomes. |
| **Preconditions** | Test suite execution completed in a sandbox environment. |
| **Trigger** | Code execution finishes in the sandbox. |
| **Main Flow** | 1. Subsystem retrieves passed and failed test case counts.<br>2. Computes the Pass Rate Score percentage.<br>3. Updates technical scorecard record. |
| **Alternative Flows** | None. |
| **Exception Flows** | Unhandled code runtime crash -> Pass Rate Score recorded as 0% with diagnostic failure details. |
| **Postconditions** | Technical Pass Rate Score updated on candidate evaluation scorecard. |
| **Business Rules** | A zero pass rate score is not an auto rejection. |

#### Use Case 4: Evaluate Integrity Risk 

| Field | Description |
| --- | --- |
| **Use Case ID** | UC-ANA-04 |
| **Use Case Name** | Evaluate Integrity Risk |
| **Actor(s)** | System |
| **Description** | Analyzing the telemetry events, sequence, and frequency and generating an integrity score based on it. |
| **Preconditions** | Telemetry streaming active and integrity threshold rules configured. |
| **Trigger** | There is detected focus loss or paste action in the telemetry stream. |
| **Main Flow** | 1. Subsystem detects violation in incoming stream.<br>2. Calculates/updates Integrity Risk Score (0–100).<br>3. Attaches evidence log. |
| **Alternative Flows** | No violation is detected. |
| **Exception Flows** | Issues with the telemetry stream. |
| **Postconditions** | Evidence logs attached and candidate Integrity Risk Score updated. |
| **Business Rules** | A high integrity score requires human revision to avoid auto rejection. |

#### Use Case 5: Export Candidate Evaluation Report

| Field | Description |
| --- | --- |
| **Use Case ID** | UC-ANA-05 |
| **Use Case Name** | Export Candidate Evaluation Report |
| **Actor(s)** | Hiring Manager, System |
| **Description** | Consolidate technical, behavioral, and integrity metrics into a unified exportable candidate report (PDF). |
| **Preconditions** | Technical, behavioral, and integrity score calculations completed. |
| **Trigger** | Hiring Manager requests report export. |
| **Main Flow** | 1. Subsystem aggregates technical pass rate, execution metrics, and integrity risk score.<br>2. Renders report in requested export format (PDF).<br>3. Delivers file for download or platform display. |
| **Alternative Flows** | None. |
| **Exception Flows** | PDF Generation failure -> System retries up to 3 times before queuing error notification. |
| **Postconditions** | Exported PDF candidate evaluation report generated. |
| **Business Rules** | None |