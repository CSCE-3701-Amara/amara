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

## 4.3.1 UC-CAND-01 — Browse Job Postings

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

## 4.3.2 UC-CAND-02 — Filter Job Postings

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

## 4.3.3 UC-CAND-03 — View Job Posting Details

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to view detailed information about a selected job posting. |
| Preconditions | The candidate is authenticated. The selected job posting exists in the system. |
| Trigger | The candidate selects a job posting from the discovery page. |
| Main Flow | 1. The candidate selects a job posting.<br>2. The system retrieves the posting details.<br>3. The system displays the position title and associated company.<br>4. The system displays the job grade, job description, required skills, and qualifications.<br>5. The system displays the offered compensation and hiring process information when provided.<br>6. The candidate reviews the displayed information. |
| Alternative Flows | 1. If compensation or hiring process information is unavailable, the system displays the remaining available details. |
| Exception Flows | 1. If the selected posting cannot be retrieved, the system displays an appropriate error message.<br>2. If the posting is no longer available for application, the system informs the candidate. |
| Postconditions | The candidate has viewed the available details of the selected job posting. |
| Business Rules | The system shall display compensation and hiring process information only when provided. |

## 4.3.4 UC-CAND-04 — Submit Job Application

| Field | Description |
|---|---|
| Actor(s) | Candidate, System |
| Description | Allows candidates to initiate and submit an application for an available job posting. |
| Preconditions | The candidate is authenticated. The selected job posting is available for application. |
| Trigger | The candidate chooses to apply for a job posting. |
| Main Flow | 1. The candidate opens the details of a job posting.<br>2. The candidate selects the application option.<br>3. The system verifies that the job posting is available for application.<br>4. The system checks whether the candidate already has an active application for the posting.<br>5. The system records the application and associates it with the candidate and job posting.<br>6. The system records the date and time of submission.<br>7. The system confirms successful submission to the candidate. |
| Alternative Flows | 1. If the candidate already has an active application, the system follows UC-CAND-05 — Prevent Duplicate Applications.<br>2. If the application process redirects to an external company website, the system directs the candidate to the external application page, where applicable. |
| Exception Flows | 1. If the posting is no longer available, the system informs the candidate that the application cannot proceed.<br>2. If the application cannot be recorded, the system displays an appropriate error message. |
| Postconditions | The application is recorded and associated with the candidate and job posting, or the candidate is informed that the application could not be submitted. |
| Business Rules | A candidate shall not have more than one active application for the same job posting. |

## 4.3.5 UC-CAND-05 — Prevent Duplicate Applications

| Field | Description |
|---|---|
| Actor(s) | Candidate, System |
| Description | Prevents candidates from submitting multiple active applications for the same job posting. |
| Preconditions | The candidate is authenticated and attempts to apply for a job posting. |
| Trigger | The system checks for an existing application when the candidate initiates an application. |
| Main Flow | 1. The candidate initiates an application for a job posting.<br>2. The system retrieves the candidate's existing applications for the selected posting.<br>3. The system checks whether an active application already exists.<br>4. If no active application exists, the system allows the application process to continue. |
| Alternative Flows | 1. If an active application exists, the system informs the candidate and provides access to the existing application. |
| Exception Flows | 1. If the system cannot verify the existing application status, it displays an appropriate error message and does not create another application. |
| Postconditions | No duplicate active application is created for the same candidate and job posting. |
| Business Rules | A candidate shall not have more than one active application for the same job posting. |

## 4.3.6 UC-CAND-06 — View Application History

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

## 4.3.7 UC-CAND-07 — View Hiring Pipeline

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to view the hiring pipeline associated with an application and track their progress through its stages. |
| Preconditions | The candidate is authenticated and has an application associated with a hiring pipeline. |
| Trigger | The candidate opens the pipeline details of an application. |
| Main Flow | 1. The candidate selects an application.<br>2. The system retrieves the associated hiring pipeline.<br>3. The system displays the applicable pipeline stages.<br>4. The system indicates the candidate's current stage.<br>5. The system displays the completion status of applicable stages.<br>6. The system indicates which stages are not yet available to the candidate. |
| Alternative Flows | 1. If the candidate has not reached an assessment stage, the system displays the current progress and indicates that the assessment is not yet available. |
| Exception Flows | 1. If the system cannot retrieve the pipeline information, it displays an appropriate error message.<br>2. If the candidate attempts to access a restricted stage, the system denies access. |
| Postconditions | The candidate can view the applicable pipeline stages and their current completion statuses. |
| Business Rules | Pipeline access shall follow the employer-defined stage availability and progression rules. |

## 4.3.8 UC-CAND-08 — Access and Complete Assessment

| Field | Description |
|---|---|
| Actor(s) | Candidate, System |
| Description | Allows candidates to access available assessments, submit their responses, and have their submissions evaluated. |
| Preconditions | The candidate is authenticated and has an application in a hiring pipeline. The corresponding assessment stage is available. |
| Trigger | The candidate selects an available assessment. |
| Main Flow | 1. The candidate opens an available assessment.<br>2. The system verifies that the assessment is available to the candidate.<br>3. The system displays the assessment instructions and requirements.<br>4. The candidate reviews the instructions and assessment content.<br>5. The candidate completes the assessment.<br>6. The candidate submits their responses.<br>7. The system records the submitted responses.<br>8. The system evaluates the responses according to the assessment's defined criteria.<br>9. The system records the assessment outcome.<br>10. The system associates the outcome with the candidate's application. |
| Alternative Flows | 1. If the assessment permits interruption, the candidate pauses the assessment and resumes it later using the preserved progress.<br>2. If the assessment requires a specific format, such as a coding exercise or recorded interview, the candidate completes the applicable format. |
| Exception Flows | 1. If the assessment is not yet available, the system denies access and informs the candidate.<br>2. If required responses are incomplete, the system prevents submission and identifies the missing information.<br>3. If the assessment submission fails, the system informs the candidate.<br>4. If evaluation is delayed, the system records the submission and indicates that the result is pending. |
| Postconditions | The assessment responses and outcome are recorded and associated with the candidate's application when successfully processed. |
| Business Rules | Candidates shall not access assessments before their corresponding pipeline stages become available. Assessment responses shall be evaluated according to the defined evaluation criteria. |

## 4.3.9 UC-CAND-09 — Track Assessment Progress

| Field | Description |
|---|---|
| Actor(s) | Candidate, System |
| Description | Allows candidates to check the completion status of their assessments and resume their progress when supported. |
| Preconditions | The candidate is authenticated and has an assessment associated with an application. |
| Trigger | The candidate opens an assessment or checks its status. |
| Main Flow | 1. The candidate accesses the assessment status.<br>2. The system retrieves the recorded assessment state.<br>3. The system displays the assessment status, such as not started, in progress, submitted, or completed, as applicable.<br>4. If the assessment supports interruption and remains incomplete, the system allows the candidate to resume it. |
| Alternative Flows | 1. If the assessment has already been submitted, the system displays its submitted status.<br>2. If the assessment does not support interruption, the system displays the current status without offering resume functionality. |
| Exception Flows | 1. If the system cannot retrieve the assessment status, it displays an appropriate error message.<br>2. If saved progress cannot be restored, the system informs the candidate. |
| Postconditions | The candidate can view the recorded assessment status and resume the assessment when permitted. |
| Business Rules | Assessment progress shall be preserved where the assessment permits interruption. |

## 4.3.10 UC-CAND-10 — View Assessment Result and Feedback

| Field | Description |
|---|---|
| Actor(s) | Candidate, System |
| Description | Allows candidates to view assessment results and any available constructive feedback. |
| Preconditions | The candidate is authenticated and has submitted an assessment. The result or feedback is available for disclosure. |
| Trigger | The candidate opens the assessment results or feedback section of an application. |
| Main Flow | 1. The candidate selects an application.<br>2. The candidate opens the assessment results section.<br>3. The system verifies that the result and feedback may be disclosed.<br>4. The system displays the assessment outcome, including whether the candidate passed or failed, when disclosed.<br>5. The system displays the associated feedback and identified areas for improvement when available. |
| Alternative Flows | 1. If the result is available but feedback is not, the system displays the result alone.<br>2. If only part of the assessment information is available for disclosure, the system displays the permitted information. |
| Exception Flows | 1. If the result is still being evaluated, the system informs the candidate that it is pending.<br>2. If disclosure is restricted, the system does not reveal the restricted information.<br>3. If the result cannot be retrieved, the system displays an appropriate error message. |
| Postconditions | The candidate can view the assessment outcome and any feedback they are permitted to access. |
| Business Rules | Assessment results and feedback shall only be visible when disclosure is permitted by the assessment configuration. Feedback shall be associated with its corresponding assessment. |

## 4.3.11 UC-CAND-11 — View Application Status

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

## 4.3.12 UC-CAND-12 — Receive Final Hiring Decision

| Field | Description |
|---|---|
| Actor(s) | Candidate, Employer, System |
| Description | Allows candidates to view the final hiring decision submitted by the employer and any accompanying feedback. |
| Preconditions | The candidate is authenticated and has an application associated with an employer. |
| Trigger | The employer submits a final hiring decision for the application. |
| Main Flow | 1. The employer submits the final hiring decision.<br>2. The system associates the decision with the corresponding application.<br>3. The system makes the decision available to the candidate.<br>4. The candidate opens the application.<br>5. The system displays the final hiring decision.<br>6. The system displays any final feedback provided by the employer. |
| Alternative Flows | 1. If the employer submits a decision without additional feedback, the system displays the decision alone.<br>2. If the candidate does not immediately open the application, the decision remains available for later viewing. |
| Exception Flows | 1. If the employer has not submitted a final decision, the system continues to display the current application status.<br>2. If the decision cannot be retrieved, the system displays an appropriate error message. |
| Postconditions | The candidate can access the final hiring decision and any accompanying feedback made available by the employer. |
| Business Rules | The final decision shall be associated with the corresponding application and shall only be displayed once submitted by the employer. |

## 4.3.13 UC-CAND-13 — View Earned Crests

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

## 4.3.14 UC-CAND-14 — Maintain Verified Skills

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

## 4.3.15 UC-CAND-15 — Generate Dynamic Resume

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

## 4.3.16 UC-CAND-16 — View Dynamic Resume

| Field | Description |
|---|---|
| Actor(s) | Candidate, System |
| Description | Allows candidates to view their dynamic resume and review their verified professional competencies. |
| Preconditions | The candidate is authenticated and has access to their professional profile. |
| Trigger | The candidate opens the dynamic resume section of their profile. |
| Main Flow | 1. The candidate opens the dynamic resume section.<br>2. The system retrieves the latest available dynamic resume.<br>3. The system displays the candidate's verified skills.<br>4. The system displays the candidate's earned Crests.<br>5. The candidate reviews the displayed professional information. |
| Alternative Flows | 1. If the resume has not yet been generated, the system initiates UC-CAND-15 — Generate Dynamic Resume.<br>2. If the candidate has no verified skills or Crests, the system displays the other available resume information. |
| Exception Flows | 1. If the system cannot retrieve or generate the resume, it displays an appropriate error message. |
| Postconditions | The candidate can view the latest available dynamic resume and its verified professional information. |
| Business Rules | The displayed verified skills and Crests shall correspond to the candidate's recorded achievements. |

## 4.3.17 UC-CAND-17 — Receive Application Notifications

| Field | Description |
|---|---|
| Actor(s) | Candidate, System |
| Description | Notifies candidates about relevant application events, including pipeline progression, assessment availability, assessment results, and final hiring decisions. |
| Preconditions | The candidate has an application in the system, and a relevant application event occurs. |
| Trigger | A relevant application event occurs, such as progression to a new pipeline stage, assessment availability, result availability, or a final hiring decision. |
| Main Flow | 1. A relevant application event occurs.<br>2. The system identifies the affected candidate.<br>3. The system determines the applicable notification.<br>4. The system prepares a notification describing the event.<br>5. The system delivers or makes the notification available to the candidate through the supported notification mechanism. |
| Alternative Flows | 1. If multiple relevant events occur, the system processes the corresponding notifications according to its notification rules.<br>2. If the candidate views the application before opening the notification, the candidate can still access the application information directly. |
| Exception Flows | 1. If notification delivery fails, the system handles the failure according to its notification mechanism.<br>2. If the underlying event has not been confirmed, the system does not send a notification that incorrectly represents the event as completed. |
| Postconditions | The relevant notification is made available to the candidate when notification processing succeeds. |
| Business Rules | The system SHOULD notify candidates when they advance to a new applicable pipeline stage, when an assessment becomes available, when an assessment result becomes available, and when a final hiring decision becomes available. |

## 4.4 APP Use Cases

## 4.5 JOB Use Cases

## 4.6 ANA Use Cases
