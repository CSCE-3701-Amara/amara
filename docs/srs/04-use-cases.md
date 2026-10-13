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

## 4.5 JOB Use Cases

## 4.6 ANA Use Cases
