## 4. Use cases

Owner: Hazem Nasr | Status: Draft | Last updated: 2026-10-09 | Jira: AMARA-30

### UC-AUTH-01: Candidate Registration
|Field|Description|
|-----|-----------|
|**Actor(s)**| Guest, System |
|**Preconditions**| The user does not have a candidate account. |
|**Trigger**| A user opens the candidate registration page. |
| **Main Flow** | 1. The user enters their personal infromation, email address, and password.<br>2. The user clicks **Register**.<br>3. The system validates the password according to the set password policy.<br>4. The system validates that the email address is in the correct format and is not already registered.<br>5. The system creates an account with the collected data.<br>6. The system sends a verification code to the user via email, and redirect the user to the verification page.<br>7. The user enters the verification code.<br>8. The user clicks **Verify**.<br>9. The system checks the code and activates the account.|
|**Alternative Flows**|**MFA setup**: After activation, the candidate chooses to set up MFA, and the system associates the MFA method with the account.|
|**Exception Flows**|**Invalid email/password**: The system displays an error message.<br>**The user does not recieve the code**: The user clicks *Resend Verification Code*, and the system sends a new code. The user then continues from Step 7.|
|**Postconditions**|The user has an active account and can login at any time.|

### UC-AUTH-02: Company Registration
|Field|Description|
|-----|-----------|
|**Actor(s)**| Guest, System, System Admin |
|**Preconditions**| The user has the information required to register a company and does not already have an approved company account. |
|**Trigger**| A user opens the company registration page. |
| **Main Flow** | 1. The user enters the company information and the requesting person's information.<br>2. The user submits the company registration request.<br>3. The system validates the submitted information.<br>4. The system checks the company against official registries.<br>5. The system verifies the identity of the requesting person when required.<br>6. The system creates a pending company registration request.<br>7. The system notifies the user that the request is pending review.<br>8. The system or a system admin accepts the request after verification.<br>9. The system creates the company account and sends an activation code to the requesting person's email address.<br>10. The user enters the activation code.<br>11. The system activates the account and requires the user to set up MFA.|
|**Exception Flows**|**Invalid company or requester information**: The system displays an error message and asks the user to correct the information.<br>**Company not found in official registries**: The system rejects the request and notifies the user.<br>**Manual review required**: The system keeps the request pending until a system admin accepts or rejects it.<br>**Request rejected**: The system notifies the user and does not create an active company account.<br>**The user does not receive the activation code**: The user requests a new code, and the system sends it to the registered email address.|
|**Postconditions**|The company has an activated account with MFA configured and can login.|

### UC-AUTH-03: Candidate Login
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate has a registered and activated account. |
|**Trigger**| The candidate opens the login page. |
| **Main Flow** | 1. The candidate enters their email address and password.<br>2. The candidate clicks **Login**.<br>3. The system validates the credentials.<br>4. The system checks that the account is activated and not deactivated or banned.<br>5. If MFA is configured, the system requests the MFA verification code.<br>6. The candidate enters the MFA verification code.<br>7. The system creates an authenticated session and redirects the candidate to the appropriate page.|
|**Alternative Flows**|**Social login**: The candidate chooses Google, GitHub, or LinkedIn, authenticates with the selected provider, and the system creates an authenticated session.<br>**MFA not configured**: The system skips MFA verification and creates the authenticated session after validating the candidate's direct credentials.|
|**Exception Flows**|**Invalid credentials**: The system displays an error message and does not create a session.<br>**Account is not activated**: The system informs the candidate that activation is required and offers to resend the activation code.<br>**Invalid MFA code**: The system displays an error message and allows the candidate to try again.<br>**Forgotten password**: The candidate requests a password reset, verifies their identity through email, and sets a new password before trying again.|
|**Postconditions**|The candidate has an authenticated session and can access resources permitted by their role.|

### UC-AUTH-04: Company Login
|Field|Description|
|-----|-----------|
|**Actor(s)**| Company Representative, System |
|**Preconditions**| The company account has been approved and activated, and MFA has been configured. |
|**Trigger**| The company representative opens the login page. |
| **Main Flow** | 1. The company representative enters the account credentials.<br>2. The company representative clicks **Login**.<br>3. The system validates the credentials.<br>4. The system checks that the account is activated and not deactivated or banned.<br>5. The system requests the MFA verification code.<br>6. The company representative enters the MFA verification code.<br>7. The system creates an authenticated session and redirects the company representative to the company dashboard.|
|**Exception Flows**|**Invalid credentials**: The system displays an error message and does not create a session.<br>**Account is not activated**: The system informs the company representative that activation is required.<br>**MFA is not configured**: The system prevents login and directs the company representative to complete MFA setup.<br>**Invalid MFA code**: The system displays an error message and allows the company representative to try again.|
|**Postconditions**|The company representative has an authenticated session with permissions assigned to the company account.|

### UC-AUTH-05: Report A Company
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System Admin, System |
|**Preconditions**| The candidate is authenticated. |
|**Trigger**| The candidate selects the option to report a company profile. |
| **Main Flow** | 1. The candidate opens the company's public profile.<br>2. The candidate selects **Report Company**.<br>3. The candidate enters the reason and supporting details.<br>4. The candidate submits the report.<br>5. The system validates the report and records it.<br>6. The system notifies system admins that a report has been submitted.<br>7. A system admin reviews the report and takes the appropriate action.|
|**Exception Flows**|**Missing report details**: The system displays an error message and asks the candidate to complete the report.<br>**Report cannot be submitted**: The system displays an error message and allows the candidate to retry.|
|**Postconditions**|The report is stored for system-admin review, or the candidate is informed that the report was not submitted.|

### UC-AUTH-06: Report A Candidate
|Field|Description|
|-----|-----------|
|**Actor(s)**| Company Representative, System Admin, System |
|**Preconditions**| The company representative is authenticated and can view the candidate information permitted by their role. |
|**Trigger**| The company representative selects the option to report a candidate. |
| **Main Flow** | 1. The company representative opens the candidate's available profile or test record.<br>2. The company representative selects **Report Candidate**.<br>3. The company representative enters the reason and supporting details.<br>4. The company representative submits the report.<br>5. The system validates the report and records it.<br>6. The system notifies system admins that a report has been submitted.<br>7. A system admin reviews the report and takes the appropriate action.|
|**Exception Flows**|**Missing report details**: The system displays an error message and asks the company representative to complete the report.<br>**Report cannot be submitted**: The system displays an error message and allows the company representative to retry.|
|**Postconditions**|The report is stored for system-admin review, or the company representative is informed that the report was not submitted.|

### UC-AUTH-07: Modify Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate is authenticated. |
|**Trigger**| The candidate opens the profile management page. |
| **Main Flow** | 1. The candidate views their current profile information.<br>2. The candidate edits the profile fields and, if needed, their job preferences.<br>3. The candidate submits the changes.<br>4. The system validates the updated information.<br>5. The system saves the changes and displays a confirmation.|
|**Exception Flows**|**Invalid profile information**: The system displays an error message and identifies the fields that need correction.<br>**Unauthorized access**: The system denies the action and does not change the profile.<br>**The changes cannot be saved**: The system displays an error message and preserves the current profile information.|
|**Postconditions**|The candidate's profile contains the validated changes, or remains unchanged if validation or saving fails.|

### UC-AUTH-08: Deactivate/Delete Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate is authenticated. |
|**Trigger**| The candidate opens the account settings page and chooses to deactivate or delete the profile. |
| **Main Flow** | 1. The candidate chooses **Deactivate** or **Delete**.<br>2. The system explains the result of the selected action and requests confirmation.<br>3. The candidate confirms the action.<br>4. The system validates the candidate's authenticated session.<br>5. The system deactivates or deletes the profile according to the selected action.<br>6. The system terminates the candidate's authenticated session.<br>7. The system displays a confirmation.|
|**Exception Flows**|**The candidate cancels**: The system keeps the profile active and returns to account settings.<br>**Unauthorized or expired session**: The system denies the action and requests authentication.<br>**The action cannot be completed**: The system displays an error message and leaves the profile unchanged.|
|**Postconditions**|The profile is deactivated or deleted, and the candidate no longer has an active authenticated session.|

### UC-AUTH-09: Share Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate is authenticated. |
|**Trigger**| The candidate chooses to share their profile. |
| **Main Flow** | 1. The candidate opens the profile sharing controls.<br>2. The candidate chooses to create a shareable link or QR code.<br>3. The system creates a shareable representation of the candidate's profile.<br>4. The system displays the link or QR code.<br>5. The candidate copies the link or displays the QR code to another person.|
|**Exception Flows**|**Profile is inactive**: The system prevents creation of a shareable link or QR code and displays an error message.<br>**The shareable representation cannot be created**: The system displays an error message and allows the candidate to retry.|
|**Postconditions**|The candidate has a shareable link or QR code for their active profile, or no shareable representation is created.|

### UC-AUTH-10: Upload Supporting Documents
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, System |
|**Preconditions**| The user is authenticated |
|**Trigger**| The user opens the document upload controls. |
| **Main Flow** | 1. The user chooses a supporting document.<br>2. The user uploads the document.<br>3. The system validates that the document is in PDF or DOCX format.<br>4. The system stores the document and associates it with the user's account.<br>5. If supported, the system extracts profile information from the document.<br>6. The system displays the extracted information for the user to review.<br>7. The user confirms the information to fill the related profile fields.<br>8. The system saves the confirmed profile information and displays a confirmation.|
|**Alternative Flows**|**No extraction**: If automatic extraction is unavailable or the user declines it, the system stores the document and the user completes the related profile fields manually.|
|**Exception Flows**|**Unsupported file format**: The system rejects the document and displays an error message listing the accepted formats.<br>**Upload fails**: The system displays an error message and allows the user to retry.<br>**Information extraction fails**: The system stores the document but does not change the profile, and informs the user that the fields must be completed manually.<br>**Unauthorized access**: The system denies the upload and does not store the document.|
|**Postconditions**|The valid document is associated with the user's account, and any extracted profile information is saved only after the user confirms it.|

### UC-AUTH-11: Rate Company
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, System |
|**Preconditions**| The candidate is authenticated and has completed a test or testing interaction with the company. |
|**Trigger**| The candidate chooses to rate the company. |
| **Main Flow** | 1. The candidate opens the company's rating form.<br>2. The candidate selects a rating value and enters optional supporting comments.<br>3. The candidate submits the rating.<br>4. The system validates that the candidate is eligible to rate the company.<br>5. The system records the rating and associates it with the company and the candidate's interaction.<br>6. The system displays a confirmation.|
|**Exception Flows**|**Candidate is not eligible**: The system prevents the rating and displays an error message.<br>**Invalid rating**: The system displays an error message and asks the candidate to provide a valid rating.<br>**Rating cannot be submitted**: The system displays an error message and allows the candidate to retry.|
|**Postconditions**|The eligible candidate's rating is recorded for the company, or no rating is recorded if validation or submission fails.|

### UC-AUTH-12: Rate Candidate
|Field|Description|
|-----|-----------|
|**Actor(s)**| Company Representative, System |
|**Preconditions**| The company representative is authenticated and has tested the candidate. |
|**Trigger**| The company representative chooses to rate the candidate. |
| **Main Flow** | 1. The company representative opens the candidate's rating form.<br>2. The company representative selects a rating value and enters optional supporting comments.<br>3. The company representative submits the rating.<br>4. The system validates that the company is eligible to rate the candidate.<br>5. The system records the rating and associates it with the candidate and the company's testing interaction.<br>6. The system displays a confirmation.|
|**Exception Flows**|**Company is not eligible**: The system prevents the rating and displays an error message.<br>**Invalid rating**: The system displays an error message and asks the company representative to provide a valid rating.<br>**Rating cannot be submitted**: The system displays an error message and allows the company representative to retry.|
|**Postconditions**|The eligible company's rating is recorded for the candidate, or no rating is recorded if validation or submission fails.|

### UC-AUTH-13: Reset Password
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, System |
|**Preconditions**| The user has a registered account and access to its registered email address. |
|**Trigger**| The user selects **Forgot Password** or **Reset Password** on the login page. |
| **Main Flow** | 1. The user enters the email address associated with their account.<br>2. The user submits the password-reset request.<br>3. The system sends a verification code to the registered email address.<br>4. The user enters the verification code.<br>5. The system verifies the user's identity.<br>6. The user enters and confirms a new password.<br>7. The system validates the new password against the password policy.<br>8. The system updates the password, invalidates existing authenticated sessions, and displays a confirmation.|
|**Alternative Flows**|**Authenticated password change**: An authenticated user opens account settings, enters their current password and a new password, and the system updates the password after validation.|
|**Exception Flows**|**Unknown email address**: The system displays an error message and does not send a reset message.<br>**Invalid or expired reset code**: The system displays an error message and allows the user to request a new reset message.<br>**Invalid new password**: The system displays an error message and asks the user to enter a password that satisfies the password policy.<br>**Password reset cannot be completed**: The system displays an error message and does not change the existing password.|
|**Postconditions**|The user has a new valid password and any previous authenticated sessions have been terminated, or the existing password remains unchanged if the reset fails.|

### UC-AUTH-14: Logout
|Field|Description|
|-----|-----------|
|**Actor(s)**| Candidate, Company Representative, System |
|**Preconditions**| The user has an authenticated session. |
|**Trigger**| The user selects **Logout**. |
| **Main Flow** | 1. The user selects **Logout**.<br>2. The system terminates the user's authenticated session.<br>3. The system redirects the user to the login page.|
|**Exception Flows**|**Session has already expired**: The system redirects the user to the login page and does not create a new session.|
|**Postconditions**|The user's authenticated session is terminated and the user cannot access protected resources without logging in again.|

### UC-AUTH-15: View Public Company Profile
|Field|Description|
|-----|-----------|
|**Actor(s)**| Guest, Candidate, Company Representative, System |
|**Preconditions**| The company has a public profile. |
|**Trigger**| A visitor opens a company's public profile page. |
| **Main Flow** | 1. The visitor opens the company profile link.<br>2. The system retrieves the company's public information.<br>3. The system displays the company information and its job postings.<br>4. The visitor selects a job posting to view its details.|
|**Exception Flows**|**Company profile does not exist**: The system displays an error message or not-found page.<br>**Company profile is unavailable**: The system displays an error message and does not expose private company information.|
|**Postconditions**|The visitor can view the company's public information and available job postings.|

### UC-AUTH-16: Manage Users and Accounts
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to manage user accounts. |
|**Trigger**| The system admin opens the user management interface. |
| **Main Flow** | 1. The system admin searches for or selects a user account.<br>2. The system displays the account information and current status.<br>3. The system admin chooses an account management action.<br>4. The system admin confirms the action.<br>5. The system applies the action and displays a confirmation.|
|**Alternative Flows**|**View account**: The system admin views the account information without changing it.<br>**Ban or deactivate account**: The system admin selects the action, and the system prevents the account from accessing the system.<br>**Reset credentials**: The system admin initiates a credential reset for the selected account.|
|**Exception Flows**|**Unauthorized access**: The system denies access to the user management interface.<br>**Account not found**: The system displays an error message and does not apply an action.<br>**Action cannot be completed**: The system displays an error message and leaves the account unchanged.|
|**Postconditions**|The selected account is viewed or has the confirmed management action applied.|

### UC-AUTH-17: Manage System Entities
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to manage system entities. |
|**Trigger**| The system admin opens an entity management interface. |
| **Main Flow** | 1. The system admin selects an entity type, such as a candidate, company, job position, or test.<br>2. The system admin selects an existing entity or chooses to create one.<br>3. The system admin enters or edits the entity information.<br>4. The system validates the information.<br>5. The system admin confirms the action.<br>6. The system adds or updates the entity and displays a confirmation.|
|**Alternative Flows**|**Delete entity**: The system admin selects an existing entity, confirms deletion, and the system deletes it according to the platform's data-retention rules.|
|**Exception Flows**|**Unauthorized access**: The system denies the action.<br>**Invalid entity information**: The system displays an error message and asks the system admin to correct the information.<br>**Entity not found**: The system displays an error message and does not apply the action.<br>**Action cannot be completed**: The system displays an error message and preserves the existing entity.|
|**Postconditions**|The selected entity is added, edited, or deleted as confirmed by the system admin.|

### UC-AUTH-18: Review Company Registration Requests
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System, Company Representative (optional)|
|**Preconditions**| The system admin is authenticated and authorized to review company registration requests. |
|**Trigger**| The system admin opens the company registration request queue. |
| **Main Flow** | 1. The system displays pending company registration requests.<br>2. The system admin selects a request.<br>3. The system displays the submitted company and requester information, including verification results.<br>4. The system admin reviews the request.<br>5. The system admin accepts or rejects the request.<br>6. The system records the decision and notifies the requester.|
|**Alternative Flows**|**Request requires more information**: The system admin marks the request as pending and requests additional information from the requester.|
|**Exception Flows**|**Unauthorized access**: The system denies access to the request queue.<br>**Request not found**: The system displays an error message and returns to the request queue.<br>**Decision cannot be recorded**: The system displays an error message and leaves the request pending.|
|**Postconditions**|The request is accepted, rejected, or remains pending with a recorded request for more information.|

### UC-AUTH-19: Manage System Taxonomies and Constants
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to manage system reference values. |
|**Trigger**| The system admin opens the taxonomy and platform-configuration interface. |
| **Main Flow** | 1. The system admin selects a taxonomy, platform constant, or reference value.<br>2. The system displays its current values.<br>3. The system admin adds, edits, or removes a value.<br>4. The system validates the change.<br>5. The system admin confirms the change.<br>6. The system saves the change and displays a confirmation.|
|**Exception Flows**|**Unauthorized access**: The system denies access to the configuration interface.<br>**Invalid or conflicting value**: The system displays an error message and does not save the change.<br>**Change cannot be saved**: The system displays an error message and preserves the existing values.|
|**Postconditions**|The selected taxonomy, constant, or reference value is updated, or remains unchanged if the operation fails.|

### UC-AUTH-20: View Administrative Reports
|Field|Description|
|-----|-----------|
|**Actor(s)**| System Admin, System |
|**Preconditions**| The system admin is authenticated and authorized to view administrative reports. |
|**Trigger**| The system admin opens the statistics and reporting interface. |
| **Main Flow** | 1. The system admin selects a report or statistics view.<br>2. The system retrieves aggregate data related to job postings.<br>3. The system displays the aggregate data and reporting results.<br>4. The system admin reviews the results for platform improvement or research.|
|**Exception Flows**|**Unauthorized access**: The system denies access to administrative reports.<br>**Data unavailable**: The system displays an error message and indicates that the report cannot currently be generated.<br>**Report generation fails**: The system displays an error message and does not display incomplete results.|
|**Postconditions**|The system admin has viewed the requested aggregate job-posting data, or no incomplete report is presented.|
