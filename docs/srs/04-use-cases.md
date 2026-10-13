## 4. Use cases

Owner: All | Status: Draft | Last updated: 2026-10-13 | Jira: HIRE-n

---

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

## 4.3.1 UC-CAND-001 — Browse Job Postings

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

## 4.3.2 UC-CAND-002 — Filter Job Postings

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

## 4.3.3 UC-CAND-003 — View Application History

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

## 4.3.4 UC-CAND-004 — Track Assessment Progress

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

## 4.3.5 UC-CAND-005 — View Assessment Result and Feedback

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

## 4.3.6 UC-CAND-006 — View Application Status

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

## 4.3.7 UC-CAND-007 — Receive Final Hiring Decision

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

## 4.3.8 UC-CAND-008 — Accept Offer

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

## 4.3.9 UC-CAND-009 — View Earned Crests

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

## 4.3.10 UC-CAND-010 — Maintain Verified Skills

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

## 4.3.11 UC-CAND-011 — Generate Dynamic Resume

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

## 4.3.12 UC-CAND-012 — View Dynamic Resume

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

## 4.3.14 UC-CAND-013 — Follow Company

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

## 4.3.15 UC-CAND-014 — Save Job Posting

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

## 4.4 APP Use Cases

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

---

## 4.6 ANA Use Cases

### Use Case 1: Stream & Ingest IDE Telemetry

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

---

### Use Case 2: Calculate Execution Efficiency Metrics

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

---

### Use Case 3: Calculate Technical Pass Rate Score

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

---

### Use Case 4: Evaluate Integrity Risk 

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

---

### Use Case 5: Export Candidate Evaluation Report

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
