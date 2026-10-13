## 4. Use cases

Owner: ALL | Status: Draft | Last updated: 2026-10-13 | Jira: PLACEHOLDER

### UC-AUTH-001: Candidate or Interviewer Registration
|Field|Description|
|-----|-----------|
|**Actor(s)**| Guest, System, Company Representative (optional) |
|**Preconditions**| The user does not have an account of the selected type. For interviewer registration, the user has a claimed association with a company that may be verified during registration. |
|**Trigger**| A user opens the candidate or interviewer registration page. |
| **Main Flow** | 1. The user selects the account type and enters their personal information, email address, and password.<br>2. For interviewer registration, the user enters company association information.<br>3. The user clicks **Register**.<br>4. The system validates the password, registration information, and email format.<br>5. The system verifies that the email address is not already registered.<br>6. For interviewer registration, the system validates the user's association with an approved company.<br>7. The system creates the account with the collected data and, for interviewers, associates it with the company.<br>8. The system sends a verification code to the user's email and redirects the user to the verification page.<br>9. The user enters the verification code.<br>10. The system checks the code and activates the account.|
|**Alternative Flows**|**Candidate MFA setup**: After activation, the candidate may choose to set up MFA, and the system associates the MFA method with the account.<br>**Interviewer MFA setup**: After email activation, the interviewer must set up MFA before the account can be used to log in.<br>**Interviewer company confirmation**: The system asks a company representative to confirm the interviewer's claimed association before activating the account.|
|**Exception Flows**|**Invalid registration information**: The system displays an error message.<br>**Company association cannot be verified**: The system does not activate the interviewer account and informs the user of the required next step.<br>**The user does not receive the code**: The user requests a new verification code, and the system sends it to the registered email address.<br>**Invalid verification code**: The system displays an error message and allows the user to re-enter the code or request a new verification code.|
|**Postconditions**|The user has an active candidate or interviewer account and can login.|

### UC-AUTH-002: Company Registration
|Field|Description|
|-----|-----------|
|**Actor(s)**| Guest, System, System Admin |
|**Preconditions**| The user has the information required to register a company and does not already have an approved company account. |
|**Trigger**| A user opens the company registration page. |
| **Main Flow** | 1. The user enters the company information and the requesting person's information.<br>2. The user submits the company registration request.<br>3. The system validates the submitted information.<br>4. The system checks the company against official registries.<br>5. The system verifies the identity of the requesting person when required.<br>6. The system creates a pending company registration request.<br>7. The system notifies the user that the request is pending review.<br>8. The system or a system admin accepts the request after verification.<br>9. The system creates the company account and sends an activation code to the requesting person's email address.<br>10. The user enters the activation code.<br>11. The system activates the account and requires the user to set up MFA.|
|**Exception Flows**|**Invalid company or requester information**: The system displays an error message and asks the user to correct the information.<br>**Company not found in official registries**: The system rejects the request and notifies the user.<br>**Manual review required**: The system keeps the request pending until a system admin accepts or rejects it.<br>**Request rejected**: The system notifies the user and does not create an active company account.<br>**The user does not receive the activation code**: The user requests a new code, and the system sends it to the registered email address.|
|**Postconditions**|The company has an activated account with MFA configured and can login.|

### UC-AUTH-003: User Login
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, Interviewer, System |
|**Preconditions**| The user has a registered and activated account. Company representatives must have an approved company account with MFA configured. Interviewers must have an active association with an approved company and MFA configured. |
|**Trigger**| The user opens the login page. |
| **Main Flow** | 1. The user enters their email address and password.<br>2. The user clicks **Login**.<br>3. The system validates the credentials.<br>4. The system checks that the account is activated and not deactivated or banned.<br>5. The system applies the authentication requirements for the user's account type.<br>6. The system creates an authenticated session with permissions scoped to the user's role and, where applicable, associated company.<br>7. The system redirects the user to the appropriate page or company workspace.|
|**Alternative Flows**|**MFA verification**: If MFA is configured, the system requests and validates the user's MFA code.<br>**Candidate social login**: The candidate chooses Google, GitHub, or LinkedIn, authenticates with the selected provider, and the system creates an authenticated session.<br>**Company MFA**: The system requires and validates the company representative's MFA code. Login is prevented if MFA is not configured.<br>**Interviewer MFA**: The system requires and validates the interviewer's MFA code. Login is prevented if MFA is not configured.<br>**Interviewer company access**: The system verifies the interviewer's company association before creating the session.|
|**Exception Flows**|**Invalid credentials**: The system displays an error message and does not create a session.<br>**Account is not activated**: The system informs the user that activation is required and offers to resend the activation code.<br>**Invalid MFA code**: The system displays an error message and allows the user to try again.<br>**Company association is no longer valid**: The system prevents interviewer login and informs the interviewer that access must be reviewed.<br>**Forgotten password**: The user requests a password reset, verifies their identity through email, and sets a new password before trying again.|
|**Postconditions**|The user has an authenticated session and can access resources permitted by their role and associated company.|

### UC-AUTH-004: Submit Report
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, System Admin, System |
|**Preconditions**| The reporting user is authenticated and can view the target information permitted by their role. |
|**Trigger**| The user selects the option to report a company or candidate. |
| **Main Flow** | 1. The user opens the target's available profile or record.<br>2. The user selects **Report Company** or **Report Candidate**.<br>3. The user enters the reason and supporting details.<br>4. The user submits the report.<br>5. The system validates and records the report.<br>6. The system notifies system admins that a report has been submitted.<br>7. A system admin reviews the report and takes the appropriate action.|
|**Alternative Flows**|**Report a company**: A candidate reports a company profile.<br>**Report a candidate**: A company representative reports a candidate profile or test record.|
|**Exception Flows**|**Missing report details**: The system displays an error message and asks the user to complete the report.<br>**Report cannot be submitted**: The system displays an error message and allows the user to retry.|
|**Postconditions**|The report is stored for system-admin review, or the user is informed that the report was not submitted.|

### UC-AUTH-005: Modify Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate is authenticated. |
|**Trigger**| The candidate opens the profile management page. |
| **Main Flow** | 1. The candidate views their current profile information.<br>2. The candidate edits the profile fields and, if needed, their job preferences.<br>3. The candidate submits the changes.<br>4. The system validates the updated information.<br>5. The system saves the changes and displays a confirmation.|
|**Exception Flows**|**Invalid profile information**: The system displays an error message and identifies the fields that need correction.<br>**Unauthorized access**: The system denies the action and does not change the profile.<br>**The changes cannot be saved**: The system displays an error message and preserves the current profile information.|
|**Postconditions**|The candidate's profile contains the validated changes, or remains unchanged if validation or saving fails.|

### UC-AUTH-006: Deactivate/Delete Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate is authenticated. |
|**Trigger**| The candidate opens the account settings page and chooses to deactivate or delete the profile. |
| **Main Flow** | 1. The candidate chooses **Deactivate** or **Delete**.<br>2. The system explains the result of the selected action and requests confirmation.<br>3. The candidate confirms the action.<br>4. The system validates the candidate's authenticated session.<br>5. The system deactivates or deletes the profile according to the selected action.<br>6. The system terminates the candidate's authenticated session.<br>7. The system displays a confirmation.|
|**Exception Flows**|**The candidate cancels**: The system keeps the profile active and returns to account settings.<br>**Unauthorized or expired session**: The system denies the action and requests authentication.<br>**The action cannot be completed**: The system displays an error message and leaves the profile unchanged.|
|**Postconditions**|The profile is deactivated or deleted, and the candidate no longer has an active authenticated session.|

### UC-AUTH-007: Share Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate is authenticated. |
|**Trigger**| The candidate chooses to share their profile. |
| **Main Flow** | 1. The candidate opens the profile sharing controls.<br>2. The candidate chooses to create a shareable link or QR code.<br>3. The system creates a shareable representation of the candidate's profile.<br>4. The system displays the link or QR code.<br>5. The candidate copies the link or displays the QR code to another person.|
|**Exception Flows**|**Profile is inactive**: The system prevents creation of a shareable link or QR code and displays an error message.<br>**The shareable representation cannot be created**: The system displays an error message and allows the candidate to retry.|
|**Postconditions**|The candidate has a shareable link or QR code for their active profile, or no shareable representation is created.|

### UC-AUTH-008: Upload Supporting Documents
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, Interviewer, System |
|**Preconditions**| The user is authenticated |
|**Trigger**| The user opens the document upload controls. |
| **Main Flow** | 1. The user chooses a supporting document.<br>2. The user uploads the document.<br>3. The system validates that the document is in PDF or DOCX format and does not exceed 10 MB.<br>4. The system stores the document and associates it with the user's account.<br>5. If supported, the system extracts profile information from the document.<br>6. The system displays the extracted information for the user to review.<br>7. The user confirms the information to fill the related profile fields.<br>8. The system saves the confirmed profile information and displays a confirmation.|
|**Alternative Flows**|**No extraction**: If automatic extraction is unavailable or the user declines it, the system stores the document and the user completes the related profile fields manually.|
|**Exception Flows**|**Unsupported file format**: The system rejects the document and displays an error message listing the accepted formats.<br>**Document exceeds size limit**: The system rejects the document and displays an error message stating that the maximum size is 10 MB.<br>**Upload fails**: The system displays an error message and allows the user to retry.<br>**Information extraction fails**: The system stores the document but does not change the profile, and informs the user that the fields must be completed manually.<br>**Unauthorized access**: The system denies the upload and does not store the document.|
|**Postconditions**|The valid document is associated with the user's account, and any extracted profile information is saved only after the user confirms it.|

### UC-AUTH-009: Submit Rating
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, System |
|**Preconditions**| The user is authenticated and has completed the testing interaction required to rate the target. |
|**Trigger**| The user chooses to rate a company or candidate. |
| **Main Flow** | 1. The user opens the target's rating form.<br>2. The user selects a rating value and enters optional supporting comments.<br>3. The user submits the rating.<br>4. The system validates that the user is eligible to rate the target.<br>5. The system records the rating and associates it with the target and the relevant testing interaction.<br>6. The system displays a confirmation.|
|**Alternative Flows**|**Rate a company**: A candidate rates a company after completing a test or testing interaction.<br>**Rate a candidate**: A company representative rates a candidate after testing the candidate.|
|**Exception Flows**|**User is not eligible**: The system prevents the rating and displays an error message.<br>**Invalid rating**: The system displays an error message and asks the user to provide a valid rating.<br>**Rating cannot be submitted**: The system displays an error message and allows the user to retry.|
|**Postconditions**|The eligible user's rating is recorded for the target, or no rating is recorded if validation or submission fails.|

### UC-AUTH-010: Reset Password
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, Interviewer, System |
|**Preconditions**| The user has a registered account and access to its registered email address. |
|**Trigger**| The user selects **Forgot Password** or **Reset Password** on the login page. |
| **Main Flow** | 1. The user enters the email address associated with their account.<br>2. The user submits the password-reset request.<br>3. The system sends a verification code to the registered email address.<br>4. The user enters the verification code.<br>5. The system verifies the user's identity.<br>6. The user enters and confirms a new password.<br>7. The system validates the new password against the password policy.<br>8. The system updates the password, invalidates existing authenticated sessions, and displays a confirmation.|
|**Alternative Flows**|**Authenticated password change**: An authenticated user opens account settings, enters their current password and a new password, and the system updates the password after validation.|
|**Exception Flows**|**Unknown email address**: The system displays an error message and does not send a reset message.<br>**Invalid or expired reset code**: The system displays an error message and allows the user to request a new reset message.<br>**Invalid new password**: The system displays an error message and asks the user to enter a password that satisfies the password policy.<br>**Password reset cannot be completed**: The system displays an error message and does not change the existing password.|
|**Postconditions**|The user has a new valid password and any previous authenticated sessions have been terminated, or the existing password remains unchanged if the reset fails.|

### UC-AUTH-011: Logout
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, Interviewer, System |
|**Preconditions**| The user has an authenticated session. |
|**Trigger**| The user selects **Logout**. |
| **Main Flow** | 1. The user selects **Logout**.<br>2. The system terminates the user's authenticated session.<br>3. The system redirects the user to the login page.|
|**Exception Flows**|**Session has already expired**: The system redirects the user to the login page and does not create a new session.|
|**Postconditions**|The user's authenticated session is terminated and the user cannot access protected resources without logging in again.|

### UC-AUTH-012: View Public Company Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Guest, Candidate, Company Representative, System |
|**Preconditions**| The company has a public profile. |
|**Trigger**| A visitor opens a company's public profile page. |
| **Main Flow** | 1. The visitor opens the company profile link.<br>2. The system retrieves the company's public information.<br>3. The system displays the company information and its job postings.<br>4. The visitor selects a job posting to view its details.|
|**Exception Flows**|**Company profile does not exist**: The system displays an error message or not-found page.<br>**Company profile is unavailable**: The system displays an error message and does not expose private company information.|
|**Postconditions**|The visitor can view the company's public information and available job postings.|

### UC-AUTH-013: Manage Users and Accounts
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to manage user accounts. |
|**Trigger**| The system admin opens the user management interface. |
| **Main Flow** | 1. The system admin searches for or selects a user account, or chooses to create an account.<br>2. The system displays the account information and current status, or collects the information required for the new account.<br>3. The system admin chooses an account management action.<br>4. The system admin confirms the action.<br>5. The system validates the account information, applies the action, and displays a confirmation.|
|**Alternative Flows**|**Create account**: The system admin enters the account information, confirms creation, and the system creates the account according to the applicable activation and verification rules.<br>**View account**: The system admin views the account information without changing it.<br>**Ban or deactivate account**: The system admin selects the action, and the system prevents the account from accessing the system.<br>**Reset credentials**: The system admin initiates a credential reset for the selected account.|
|**Exception Flows**|**Unauthorized access**: The system denies access to the user management interface.<br>**Account not found**: The system displays an error message and does not apply an action.<br>**Action cannot be completed**: The system displays an error message and leaves the account unchanged.|
|**Postconditions**|The selected account is viewed or has the confirmed management action applied.|

### UC-AUTH-014: Manage System Entities
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to manage system entities. |
|**Trigger**| The system admin opens an entity management interface. |
| **Main Flow** | 1. The system admin selects an entity type, such as a candidate, company, job position, or test.<br>2. The system admin selects an existing entity or chooses to create one.<br>3. The system admin enters or edits the entity information.<br>4. The system validates the information.<br>5. The system admin confirms the action.<br>6. The system adds or updates the entity and displays a confirmation.|
|**Alternative Flows**|**Delete entity**: The system admin selects an existing entity, confirms deletion, and the system deletes it according to the platform's data-retention rules.|
|**Exception Flows**|**Unauthorized access**: The system denies the action.<br>**Invalid entity information**: The system displays an error message and asks the system admin to correct the information.<br>**Entity not found**: The system displays an error message and does not apply the action.<br>**Action cannot be completed**: The system displays an error message and preserves the existing entity.|
|**Postconditions**|The selected entity is added, edited, or deleted as confirmed by the system admin.|

### UC-AUTH-015: Review Company Registration Requests
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System, Company Representative (optional)|
|**Preconditions**| The system admin is authenticated and authorized to review company registration requests. |
|**Trigger**| The system admin opens the company registration request queue. |
| **Main Flow** | 1. The system displays pending company registration requests.<br>2. The system admin selects a request.<br>3. The system displays the submitted company and requester information, including verification results.<br>4. The system admin reviews the request.<br>5. The system admin accepts or rejects the request.<br>6. The system records the decision and notifies the requester.|
|**Alternative Flows**|**Request requires more information**: The system admin marks the request as pending and requests additional information from the requester.|
|**Exception Flows**|**Unauthorized access**: The system denies access to the request queue.<br>**Request not found**: The system displays an error message and returns to the request queue.<br>**Decision cannot be recorded**: The system displays an error message and leaves the request pending.|
|**Postconditions**|The request is accepted, rejected, or remains pending with a recorded request for more information.|

### UC-AUTH-016: Manage System Taxonomies and Constants
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to manage system reference values. |
|**Trigger**| The system admin opens the taxonomy and platform-configuration interface. |
| **Main Flow** | 1. The system admin selects a taxonomy, platform constant, or reference value.<br>2. The system displays its current values.<br>3. The system admin adds, edits, or removes a value.<br>4. The system validates the change.<br>5. The system admin confirms the change.<br>6. The system saves the change and displays a confirmation.|
|**Exception Flows**|**Unauthorized access**: The system denies access to the configuration interface.<br>**Invalid or conflicting value**: The system displays an error message and does not save the change.<br>**Change cannot be saved**: The system displays an error message and preserves the existing values.|
|**Postconditions**|The selected taxonomy, constant, or reference value is updated, or remains unchanged if the operation fails.|

### UC-AUTH-017: View Administrative Reports
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to view administrative reports. |
|**Trigger**| The system admin opens the statistics and reporting interface. |
| **Main Flow** | 1. The system admin selects a report or statistics view.<br>2. The system retrieves aggregate data related to job postings.<br>3. The system displays the aggregate data and reporting results.<br>4. The system admin reviews the results for platform improvement or research.|
|**Exception Flows**|**Unauthorized access**: The system denies access to administrative reports.<br>**Data unavailable**: The system displays an error message and indicates that the report cannot currently be generated.<br>**Report generation fails**: The system displays an error message and does not display incomplete results.|
|**Postconditions**|The system admin has viewed the requested aggregate job-posting data, or no incomplete report is presented.|
