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

### 4.3.1 UC-CAND-01 — Browse Job Postings

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to browse available Computer Science-related job opportunities. |
| Preconditions | The candidate is authenticated. Job postings are available in the system. |
| Trigger | The candidate opens the job discovery page. |
| Main Flow | 1. The candidate opens the job discovery page. <br> 2. The system retrieves available job postings. <br> 3. The system displays the available postings and their essential information. <br> 4. The candidate browses the displayed postings. |
| Alternative Flows | 1. If no postings are available, the system informs the candidate that no jobs are currently available. |
| Exception Flows | 1. If the system cannot retrieve the postings, it displays an appropriate error message. |
| Postconditions | Available job postings are displayed to the candidate, or the candidate is informed that none are available. |
| Business Rules | Only job postings available for application shall be displayed as available opportunities. |

### 4.3.2 UC-CAND-02 — Filter Job Postings

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to narrow job search results using selected criteria. |
| Preconditions | The candidate is authenticated and can access job discovery. |
| Trigger | The candidate selects one or more job filters. |
| Main Flow | 1. The candidate opens the job discovery page. <br> 2. The system displays the available filtering criteria. <br> 3. The candidate selects one or more criteria, such as compensation, position, job grade, or company. <br> 4. The system applies the selected criteria. <br> 5. The system displays the matching job postings. |
| Alternative Flows | 1. If no postings match the criteria, the system informs the candidate that no matching jobs were found. <br> 2. The candidate changes or clears the selected criteria and repeats the search. |
| Exception Flows | 1. If the system cannot retrieve filtered results, it displays an appropriate error message. |
| Postconditions | The displayed job postings match the selected criteria, or no matching postings are displayed. |
| Business Rules | Multiple filtering criteria may be combined. Only available postings shall be included in the results. |

### 4.3.3 UC-CAND-03 — View Job Posting Details

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to inspect the details of a job before deciding whether to apply. |
| Preconditions | The candidate is authenticated. The selected job posting exists. |
| Trigger | The candidate selects a job posting. |
| Main Flow | 1. The candidate selects a job posting. <br> 2. The system retrieves the posting details. <br> 3. The system displays the position title, company, job grade, job description, required qualifications, and compensation when provided. <br> 4. The candidate reviews the information. |
| Alternative Flows | 1. If the posting is no longer available, the system informs the candidate and indicates that an application cannot be submitted. |
| Exception Flows | 1. If the posting details cannot be retrieved, the system displays an appropriate error message. |
| Postconditions | The candidate has viewed the available details of the selected job posting. |
| Business Rules | The system shall display compensation only when the employer has provided it. Unavailable postings shall not accept new applications. |

### 4.3.4 UC-CAND-04 — Submit Job Application

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to apply for an available job posting. |
| Preconditions | The candidate is authenticated. The selected job posting is available for application. The candidate does not already have an active application for that posting. |
| Trigger | The candidate selects the option to apply for a job. |
| Main Flow | 1. The candidate opens the job posting. <br> 2. The candidate selects the application option. <br> 3. The system verifies that the posting is available and that no duplicate active application exists. <br> 4. The system creates the application and associates it with the candidate and job posting. <br> 5. The system records the submission date and time. <br> 6. The system confirms successful submission. |
| Alternative Flows | 1. If additional application information is required, the system requests it before submission. <br> 2. If the candidate already has an active application, the system informs the candidate and does not create another application. |
| Exception Flows | 1. If the posting becomes unavailable before submission, the system informs the candidate and rejects the application. <br> 2. If the application cannot be saved, the system displays an error and does not report the submission as successful. |
| Postconditions | A successfully submitted application is stored and associated with the candidate and job posting. |
| Business Rules | A candidate shall not have more than one active application for the same job posting. Only available postings shall accept applications. |

### 4.3.5 UC-CAND-05 — View Application History

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to review their submitted applications and their current statuses. |
| Preconditions | The candidate is authenticated. |
| Trigger | The candidate opens the application history page. |
| Main Flow | 1. The candidate opens the application history page. <br> 2. The system retrieves the candidate's applications. <br> 3. The system displays the associated job postings and current application statuses. <br> 4. The candidate selects an application to view its details. <br> 5. The system displays the selected application's information. |
| Alternative Flows | 1. If the candidate has not submitted any applications, the system displays an appropriate empty-state message. |
| Exception Flows | 1. If the application history cannot be retrieved, the system displays an error message. |
| Postconditions | The candidate has viewed their application history and, if selected, the details of an application. |
| Business Rules | Candidates shall only be able to access their own application information. Application statuses shall reflect the latest recorded updates. |

### 4.3.6 UC-CAND-06 — View Hiring Pipeline

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to track their progress through the hiring process for a particular application. |
| Preconditions | The candidate is authenticated and has an application associated with a hiring pipeline. |
| Trigger | The candidate opens an application's hiring pipeline. |
| Main Flow | 1. The candidate selects an application. <br> 2. The candidate opens its hiring pipeline. <br> 3. The system retrieves the pipeline stages applicable to the application. <br> 4. The system displays the stages and indicates the current stage. <br> 5. The system indicates completed, current, and not-yet-available stages. |
| Alternative Flows | 1. If the application has not entered a pipeline stage, the system displays the current application status and the available pipeline information. |
| Exception Flows | 1. If the pipeline information cannot be retrieved, the system displays an appropriate error message. |
| Postconditions | The candidate can view the applicable pipeline stages and their recorded progress. |
| Business Rules | Pipeline progression shall follow the employer-defined workflow. Candidates shall not be able to access stages that have not been made available to them. |

### 4.3.7 UC-CAND-07 — Access and Complete Assessment

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to access and complete assessments assigned to their applications. |
| Preconditions | The candidate is authenticated. An assessment is assigned to the application and its pipeline stage is available. |
| Trigger | The candidate selects an available assessment. |
| Main Flow | 1. The candidate opens the relevant application. <br> 2. The candidate selects an available assessment. <br> 3. The system verifies assessment availability. <br> 4. The system displays the assessment instructions and content. <br> 5. The candidate completes the assessment and provides the required responses. <br> 6. The candidate submits the assessment. <br> 7. The system records the submitted responses. <br> 8. The system evaluates the responses according to the assessment's evaluation criteria. <br> 9. The system records the result and updates the assessment completion status. |
| Alternative Flows | 1. If the assessment permits interruptions, the candidate leaves the assessment and later resumes the saved progress. <br> 2. If the assessment permits incomplete submissions, the system processes the submission according to its configured rules. |
| Exception Flows | 1. If the assessment is not yet available, the system prevents access and informs the candidate. <br> 2. If the submission fails, the system informs the candidate and does not mark the assessment as successfully submitted. <br> 3. If the assessment cannot be evaluated, the system records or reports the failure according to the applicable error-handling rules. |
| Postconditions | The assessment is recorded as completed and evaluated when successful. Otherwise, its previous completion status is retained and the candidate is informed of the problem. |
| Business Rules | Assessment access shall follow pipeline availability rules. Evaluation shall follow the criteria defined for the assessment. Progress shall be preserved when interruption and resumption are supported. |

### 4.3.8 UC-CAND-08 — View Assessment Results and Feedback

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to review disclosed assessment results and available feedback. |
| Preconditions | The candidate is authenticated and has an assessment associated with an application. The result or feedback is available for disclosure. |
| Trigger | The candidate opens an assessment's results or feedback. |
| Main Flow | 1. The candidate opens the relevant application. <br> 2. The candidate selects the completed assessment. <br> 3. The system checks whether the result and feedback are available for disclosure. <br> 4. The system displays the assessment result when disclosed. <br> 5. The system displays the available feedback and areas for improvement. |
| Alternative Flows | 1. If the assessment has been submitted but not evaluated, the system indicates that the result is pending. <br> 2. If a result is available but feedback is not, the system displays the result without feedback. |
| Exception Flows | 1. If the assessment result cannot be retrieved, the system displays an error message. <br> 2. If the result is not authorized for disclosure, the system prevents access to it. |
| Postconditions | The candidate has viewed the result and any feedback authorized for disclosure. |
| Business Rules | Assessment results and feedback shall only be disclosed according to the applicable assessment and hiring process rules. |

### 4.3.9 UC-CAND-09 — Receive Final Hiring Decision

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to view the final hiring decision for an application. |
| Preconditions | The candidate is authenticated and has an application for which a final decision has been recorded and made available. |
| Trigger | The candidate opens the application or receives a notification that a final decision is available. |
| Main Flow | 1. The candidate opens the relevant application. <br> 2. The system retrieves the final hiring decision. <br> 3. The system displays the decision. <br> 4. The system displays any final feedback made available by the employer. |
| Alternative Flows | 1. If no final decision has been recorded, the system displays the current application status instead. <br> 2. If the decision is available but no feedback has been provided, the system displays the decision without feedback. |
| Exception Flows | 1. If the decision cannot be retrieved, the system displays an appropriate error message. |
| Postconditions | The candidate has viewed the final decision and any available final feedback. |
| Business Rules | The system shall display a final decision only after it has been recorded and made available to the candidate. |

### 4.3.10 UC-CAND-10 — Earn Verified Skills and Crests

| Field | Description |
|---|---|
| Actor(s) | Candidate; Assessment Evaluation System |
| Description | Allows a candidate to obtain verified skills and Crests after satisfying the applicable assessment criteria. |
| Preconditions | The candidate has completed an eligible assessment, and its evaluation criteria and verification rules are defined. |
| Trigger | An eligible assessment is evaluated successfully. |
| Main Flow | 1. The system evaluates the assessment. <br> 2. The system checks the result against the applicable verification criteria. <br> 3. The system determines which skills or competencies have been verified. <br> 4. The system awards the corresponding Crest or verified skill. <br> 5. The system associates the verified information with the candidate's profile. |
| Alternative Flows | 1. If the candidate does not satisfy the required criteria, the system records the assessment outcome without awarding the corresponding Crest. <br> 2. If the assessment verifies multiple competencies, the system records each competency for which the criteria are satisfied. |
| Exception Flows | 1. If evaluation fails, the system does not award a Crest based on an unverified result. <br> 2. If the verification record cannot be saved, the system reports the failure and avoids reporting an award as successful. |
| Postconditions | Eligible verified skills and Crests are associated with the candidate, or no new award is made if the criteria are not satisfied. |
| Business Rules | Verified skills and Crests shall be awarded according to defined assessment and verification criteria. Candidate-provided claims alone shall not establish a verified skill. |

### 4.3.11 UC-CAND-11 — View Earned Crests and Verified Skills

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to review the skills and competencies verified through eligible assessments. |
| Preconditions | The candidate is authenticated. |
| Trigger | The candidate opens the verified skills or Crests section of their profile. |
| Main Flow | 1. The candidate opens the verified skills section. <br> 2. The system retrieves the candidate's verified skills and earned Crests. <br> 3. The system displays each verified skill and its associated Crest, where applicable. <br> 4. The system displays the available verification information. |
| Alternative Flows | 1. If the candidate has no verified skills or Crests, the system displays an appropriate empty-state message. |
| Exception Flows | 1. If the records cannot be retrieved, the system displays an error message. |
| Postconditions | The candidate has viewed their available verified skills and Crests. |
| Business Rules | Only skills supported by the platform's verification process shall be presented as verified. |

### 4.3.12 UC-CAND-12 — View and Generate Dynamic Resume

| Field | Description |
|---|---|
| Actor(s) | Candidate |
| Description | Allows candidates to view a resume that reflects their professional information and verified competencies. |
| Preconditions | The candidate is authenticated. Candidate profile information is available. |
| Trigger | The candidate opens the dynamic resume. |
| Main Flow | 1. The candidate opens the dynamic resume section. <br> 2. The system retrieves the candidate's relevant professional information. <br> 3. The system retrieves the candidate's verified skills and earned Crests. <br> 4. The system generates or updates the resume using the available information. <br> 5. The system displays the dynamic resume to the candidate. |
| Alternative Flows | 1. If some optional profile information is missing, the system displays the available information. <br> 2. If the candidate has earned new Crests or verified skills, the system incorporates them into the resume. |
| Exception Flows | 1. If the resume cannot be generated or retrieved, the system displays an appropriate error message. |
| Postconditions | The candidate can view a resume reflecting the relevant information currently available in the system. |
| Business Rules | Verified skills and Crests shall be represented according to the platform's verification records. The resume shall reflect relevant updates to candidate information. |

### 4.3.13 UC-CAND-13 — Participate in Anonymous Evaluation

| Field | Description |
|---|---|
| Actor(s) | Candidate; Employer or Evaluator; System |
| Description | Supports the evaluation of a candidate without exposing identifying information during designated hiring stages. |
| Preconditions | The candidate has an application in a hiring pipeline. The employer has configured one or more stages for anonymous evaluation. |
| Trigger | The candidate's application enters a stage configured for anonymous evaluation. |
| Main Flow | 1. The system identifies the applicable pipeline stage. <br> 2. The system determines whether anonymous evaluation is enabled for that stage. <br> 3. The system restricts the display of candidate-identifying information to users who are not authorized to access it during that stage. <br> 4. The evaluator reviews the permitted candidate information and assessment results. <br> 5. When the application reaches a stage where identity disclosure is permitted, the system makes the candidate's identity available to authorized users according to the configured rules. |
| Alternative Flows | 1. If the current stage does not require anonymity, the system applies the normal information-access rules for that stage. <br> 2. If identity disclosure is configured for a later stage, the system keeps identifying information restricted until that stage is reached. |
| Exception Flows | 1. If the system cannot determine the applicable disclosure rules, it shall not expose restricted identifying information until the rules can be enforced. <br> 2. If an unauthorized user attempts to access hidden identifying information, the system denies access. |
| Postconditions | Candidate-identifying information remains restricted during anonymous stages and is disclosed only when permitted by the configured rules. |
| Business Rules | Anonymity shall follow the employer-configured pipeline rules. Only authorized users may access candidate-identifying information. Information necessary for evaluating professional qualifications may remain available where permitted by the rules. |

### 4.3.14 UC-CAND-14 — Receive Application Notifications

| Field | Description |
|---|---|
| Actor(s) | Candidate; Notification System |
| Description | Keeps candidates informed about relevant updates to their applications and assessments. |
| Preconditions | The candidate is authenticated. A notification-triggering event has occurred, and the relevant notification mechanism is available. |
| Trigger | An application stage changes, an assessment becomes available, an assessment result is disclosed, or a final hiring decision becomes available. |
| Main Flow | 1. The system detects a notification-triggering event. <br> 2. The system identifies the candidate associated with the relevant application. <br> 3. The system prepares a notification containing the relevant update. <br> 4. The system sends or makes the notification available through the supported notification mechanism. <br> 5. The candidate receives or views the notification. <br> 6. The candidate opens the relevant application when further details are needed. |
| Alternative Flows | 1. If the notification is available within the platform but external delivery is unavailable, the candidate can access the update through the supported in-platform notification mechanism, if implemented. |
| Exception Flows | 1. If notification delivery fails, the system handles the failure according to its notification-delivery rules. <br> 2. If the associated application cannot be retrieved, the system does not provide a broken application link and reports the issue through the supported error-handling mechanism. |
| Postconditions | The notification is made available or delivered according to the supported notification mechanism. |
| Business Rules | Notifications shall correspond to recorded application or assessment events. Notifications shall not disclose information that the candidate is not authorized to access. |

## 4.4 APP Use Cases

## 4.5 JOB Use Cases

## 4.6 ANA Use Cases
