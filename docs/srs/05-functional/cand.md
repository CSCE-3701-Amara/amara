## 5.2 CAND: Candidates & Dynamic Resume

Owner: Ali | Status: Draft | Last updated: 2026-10-08 | Jira: HIRE-n

# Candidate Functional Requirements

## 1. Overview

This section defines the functional requirements related to the **Candidate Subsystem** of the platform.

The Candidate Subsystem enables candidates to:

- Discover available Computer Science-related job opportunities
- Search and filter job postings
- View job posting details
- Apply for available positions
- Participate in employer-defined hiring pipelines
- Complete assessment modules
- Track application and assessment progress
- Receive assessment feedback and hiring decisions
- Accumulate verified skills through successfully completed assessments
- Build and maintain a dynamic resume
- Manage professional profile information
- Participate in anonymized hiring processes

User account creation, authentication, and account management are outside the scope of this subsystem.

---

# 2. Scope

The Candidate Subsystem includes the following functional areas:

1. Job Discovery
2. Job Posting Details
3. Job Applications
4. Hiring Pipeline
5. Assessment Completion
6. Assessment Feedback
7. Application Status and Decisions
8. Verified Skills and Crests
9. Dynamic Resume
10. Professional Profile
11. Candidate Anonymity
12. Candidate Notifications

---

# 3. Functional Requirements

## 3.1 Job Discovery

### FR-CAND-001 — Browse Job Postings

**Description:**  
The system shall allow candidates to browse available job postings relevant to Computer Science-related fields.

**Requirements:**
- The system shall display currently available job postings to candidates.
- The system shall display the essential information associated with each job posting.
- The system shall exclude job postings that are no longer available for application.

**Priority:** Highest

**Related Use Case:** UC-CAND-01 — Browse Job Postings

---

### FR-CAND-002 — Filter Job Postings

**Description:**  
The system shall allow candidates to filter available job postings according to supported criteria.

**Requirements:**
- The system shall allow candidates to filter postings by offered compensation.
- The system shall allow candidates to filter postings by position.
- The system shall allow candidates to filter postings by job grade.
- The system shall allow candidates to filter postings by company.
- The system shall allow candidates to combine multiple filter criteria.
- The system shall update the displayed postings according to the selected criteria.

**Priority:** Highest

**Related Use Case:** UC-CAND-02 — Filter Job Postings

---

### FR-CAND-003 — View Job Posting Details

**Description:**  
The system shall allow candidates to view the details of an available job posting.

**Requirements:**
- The system shall display the position title.
- The system shall display the company associated with the posting.
- The system shall display the job grade.
- The system shall display the offered compensation when provided.
- The system shall display the job description.
- The system shall display the required skills and qualifications.
- The system shall display the relevant hiring process information when provided.

**Priority:** Highest

**Related Use Case:** UC-CAND-03 — View Job Posting

---

## 3.2 Job Applications

### FR-CAND-004 — Submit Job Application

**Description:**  
The system shall allow candidates to submit an application for an available job posting.

**Requirements:**
- The system shall allow a candidate to initiate an application from an available job posting.
- The system shall associate the submitted application with the candidate and the selected job posting.
- The system shall record the date and time at which the application was submitted.
- The system shall confirm successful submission of the application to the candidate.

**Priority:** Highest

**Related Use Case:** UC-CAND-04 — Apply for Job

---

### FR-CAND-005 — Prevent Duplicate Applications

**Description:**  
The system shall prevent a candidate from submitting multiple active applications for the same job posting.

**Requirements:**
- The system shall determine whether the candidate already has an active application for the selected job posting.
- The system shall prevent submission of a duplicate active application.
- The system shall inform the candidate when an active application already exists.

**Priority:** High

**Related Use Case:** UC-CAND-04 — Apply for Job

---

### FR-CAND-006 — View Application History

**Description:**  
The system shall allow candidates to view their submitted job applications.

**Requirements:**
- The system shall display the jobs to which the candidate has applied.
- The system shall display the current status of each application.
- The system shall allow the candidate to access the relevant application details.

**Priority:** Highest

**Related Use Case:** UC-CAND-05 — View Applications

---

## 3.3 Hiring Pipeline

### FR-CAND-007 — View Hiring Pipeline

**Description:**  
The system shall allow candidates to view the stages of the hiring pipeline associated with an application.

**Requirements:**
- The system shall display the stages applicable to the candidate's application.
- The system shall indicate the candidate's current stage.
- The system shall indicate the completion status of applicable stages.
- The system shall indicate stages that are not yet available to the candidate.

**Priority:** Highest

**Related Use Case:** UC-CAND-06 — View Hiring Pipeline

---

### FR-CAND-008 — Access Assigned Assessment

**Description:**  
The system shall allow candidates to access an assessment when the corresponding hiring pipeline stage becomes available to them.

**Requirements:**
- The system shall make an assessment available when the candidate reaches its corresponding pipeline stage.
- The system shall prevent candidates from accessing assessments that are not yet available to them.
- The system shall provide the candidate with the information necessary to begin the assessment.

**Priority:** Highest

**Related Use Case:** UC-CAND-07 — Take Assessment

---

### FR-CAND-009 — Complete Assessment

**Description:**  
The system shall allow candidates to complete assessment modules assigned to them as part of a hiring pipeline.

**Requirements:**
- The system shall present the assessment content to the candidate.
- The system shall record the candidate's submitted responses.
- The system shall evaluate the candidate's responses according to the assessment's defined evaluation criteria.
- The system shall record the outcome of the assessment.
- The system shall associate the assessment result with the candidate's application.

**Priority:** Highest

**Related Use Case:** UC-CAND-07 — Take Assessment

---

### FR-CAND-010 — Track Assessment Progress

**Description:**  
The system shall allow candidates to track their progress within an available assessment.

**Requirements:**
- The system shall indicate the candidate's assessment completion status.
- The system shall preserve assessment progress where the assessment permits interruption.
- The system shall indicate whether an assessment has been completed or remains incomplete.

**Priority:** Highest

**Related Use Case:** UC-CAND-07 — Take Assessment

---

## 3.4 Assessment Feedback

### FR-CAND-011 — View Assessment Feedback

**Description:**  
The system shall provide candidates with feedback associated with their completed assessments when feedback is available.

**Requirements:**
- The system shall provide the candidate with feedback after the assessment has been evaluated.
- The system shall associate the feedback with the corresponding assessment.
- The system shall make available any applicable areas of improvement identified by the assessment.
- The system shall make the feedback accessible from the candidate's application.

**Priority:** Highest

**Related Use Case:** UC-CAND-08 — View Assessment Feedback

---

### FR-CAND-012 — View Assessment Result

**Description:**  
The system shall allow candidates to view the outcome of completed assessments when the employer has configured the assessment result to be visible.

**Requirements:**
- The system shall indicate whether the candidate passed or failed the applicable assessment when the result is disclosed.
- The system shall display the assessment result together with its associated feedback when available.

**Priority:** Highest

**Related Use Case:** UC-CAND-08 — View Assessment Feedback

---

## 3.5 Application Status and Hiring Decision

### FR-CAND-013 — View Application Status

**Description:**  
The system shall allow candidates to view the current status of their job applications.

**Requirements:**
- The system shall display the current status of each active application.
- The system shall update the application status when the application progresses through the hiring pipeline.
- The system shall indicate when an application has reached a final decision.

**Priority:** Highest

**Related Use Case:** UC-CAND-05 — View Applications

---

### FR-CAND-014 — Receive Final Hiring Decision

**Description:**  
The system shall provide candidates with the final hiring decision associated with their application.

**Requirements:**
- The system shall display the final decision once the employer has submitted it.
- The system shall associate the decision with the corresponding application.
- The system shall provide any final feedback made available by the employer.

**Priority:** Highest

**Related Use Case:** UC-CAND-09 — View Hiring Decision

---

## 3.6 Verified Skills and Crests

### FR-CAND-015 — Award Verified Crest

**Description:**  
The system shall award a Crest to a candidate when the candidate successfully completes an eligible assessment module.

**Requirements:**
- The system shall determine whether the candidate has successfully completed the applicable assessment.
- The system shall award the corresponding Crest when the required criteria are satisfied.
- The system shall associate the Crest with the skill or competency verified by the assessment.
- The system shall associate the Crest with the candidate's profile.

**Priority:** Highest

**Related Use Case:** UC-CAND-10 — Earn Crest

---

### FR-CAND-016 — View Earned Crests

**Description:**  
The system shall allow candidates to view the Crests they have earned.

**Requirements:**
- The system shall display the Crests associated with the candidate.
- The system shall display the skill or competency represented by each Crest.
- The system shall display the relevant verification information associated with each Crest.

**Priority:** Highest

**Related Use Case:** UC-CAND-11 — View Crests

---

### FR-CAND-017 — Maintain Verified Skills

**Description:**  
The system shall maintain a record of skills and competencies verified through successfully completed assessment modules.

**Requirements:**
- The system shall associate verified skills with the candidate.
- The system shall retain verified skills across different eligible hiring processes.
- The system shall update the candidate's verified skills when new eligible assessments are successfully completed.

**Priority:** Highest

**Related Use Case:** UC-CAND-10 — Earn Crest

---

## 3.7 Dynamic Resume

### FR-CAND-018 — Generate Dynamic Resume

**Description:**  
The system shall generate a dynamic resume based on the candidate's professional information and verified competencies.

**Requirements:**
- The system shall include the candidate's verified skills in the dynamic resume.
- The system shall include the candidate's earned Crests in the dynamic resume.
- The system shall incorporate newly verified competencies into the dynamic resume.
- The system shall update the dynamic resume when relevant candidate information changes.

**Priority:** Highest

**Related Use Case:** UC-CAND-12 — Manage Dynamic Resume

---

### FR-CAND-019 — View Dynamic Resume

**Description:**  
The system shall allow candidates to view their dynamic resume.

**Requirements:**
- The system shall display the candidate's professional information included in the resume.
- The system shall display the candidate's verified competencies.
- The system shall display the candidate's earned Crests.
- The system shall reflect the candidate's current verified skills.

**Priority:** Highest

**Related Use Case:** UC-CAND-12 — Manage Dynamic Resume

---

## 3.8 Candidate Anonymity

### FR-CAND-020 — Apply Candidate Anonymity

**Description:**  
The system shall hide designated personally identifiable candidate information during hiring pipeline stages configured for anonymous evaluation.

**Requirements:**
- The system shall determine whether anonymity is enabled for the applicable hiring pipeline stage.
- The system shall hide candidate information designated as non-disclosable during anonymous stages.
- The system shall continue to expose information required for evaluating the candidate's professional qualifications.
- The system shall prevent hidden candidate information from being displayed to authorized evaluators during the anonymous stage.

**Priority:** Medium

**Related Use Case:** UC-CAND-13 — Participate in Anonymous Evaluation

---

### FR-CAND-021 — Reveal Candidate Identity

**Description:**  
The system shall reveal the candidate's identity when the candidate reaches a hiring pipeline stage configured for identity disclosure.

**Requirements:**
- The system shall determine when the candidate reaches the configured identity-disclosure stage.
- The system shall make the candidate's identity available to authorized users at that stage.
- The system shall retain the candidate's application and assessment history after identity disclosure.

**Priority:** Medium

**Related Use Case:** UC-CAND-13 — Participate in Anonymous Evaluation

---

## 3.9 Candidate Notifications

### FR-CAND-022 — Receive Application Notifications

**Description:**  
The system shall notify candidates of significant changes to their job applications.

**Requirements:**
- The system shall notify candidates when they advance to a new applicable hiring pipeline stage.
- The system shall notify candidates when an assessment becomes available.
- The system shall notify candidates when an assessment result becomes available.
- The system shall notify candidates when a final hiring decision becomes available.

**Priority:** High

**Related Use Case:** UC-CAND-14 — Receive Notifications

---

# 4. Business Rules

The following business rules govern the behavior of the Candidate Subsystem:

- A candidate shall not be able to participate in an assessment before the corresponding pipeline stage becomes available.
- A Crest shall only be awarded when the candidate satisfies the criteria defined for the corresponding assessment.
- A verified skill shall remain associated with the candidate after the hiring process in which it was earned has ended.
- Candidate-to-candidate communication shall not be provided by the platform.
- Personally identifiable candidate information shall remain hidden during pipeline stages configured for anonymous evaluation.
- Candidate identity shall only become available to authorized employer users at the configured identity-disclosure stage.
- Assessment results and feedback shall only be visible to candidates when disclosure is permitted by the assessment configuration.
- A candidate shall not have more than one active application for the same job posting.

---

# 5. Error Handling

The system shall provide appropriate feedback when:

- A candidate attempts to apply to an unavailable job posting.
- A candidate attempts to submit a duplicate active application.
- A candidate attempts to access an assessment that is not yet available.
- A candidate attempts to submit an incomplete assessment where completion is required.
- An uploaded professional document does not meet the supported file requirements.
- A candidate provides an invalid professional profile link.
- An assessment cannot be submitted successfully.
- Assessment results or feedback are not yet available.
- A candidate attempts to access information that is restricted during an anonymous hiring stage.

---

# 6. Priority Summary

| ID | Requirement | Priority |
|---|---|---|
| FR-CAND-001 | Browse Job Postings | Highest |
| FR-CAND-002 | Filter Job Postings | Highest |
| FR-CAND-003 | View Job Posting Details | Highest |
| FR-CAND-004 | Submit Job Application | Highest |
| FR-CAND-005 | Prevent Duplicate Applications | High |
| FR-CAND-006 | View Application History | Highest |
| FR-CAND-007 | View Hiring Pipeline | Highest |
| FR-CAND-008 | Access Assigned Assessment | Highest |
| FR-CAND-009 | Complete Assessment | Highest |
| FR-CAND-010 | Track Assessment Progress | Highest |
| FR-CAND-011 | View Assessment Feedback | Highest |
| FR-CAND-012 | View Assessment Result | Highest |
| FR-CAND-013 | View Application Status | Highest |
| FR-CAND-014 | Receive Final Hiring Decision | Highest |
| FR-CAND-015 | Award Verified Crest | Highest |
| FR-CAND-016 | View Earned Crests | Highest |
| FR-CAND-017 | Maintain Verified Skills | Highest |
| FR-CAND-018 | Generate Dynamic Resume | Highest |
| FR-CAND-019 | View Dynamic Resume | Highest |
| FR-CAND-020 | Manage Professional Information | Highest |
| FR-CAND-021 | Upload Professional Documents | High |
| FR-CAND-022 | Manage Professional Links | High |
| FR-CAND-023 | Apply Candidate Anonymity | Medium |
| FR-CAND-024 | Reveal Candidate Identity | Medium |
| FR-CAND-025 | Receive Application Notifications | High |

---

# 7. Traceability

The Jira/User Story identifiers shall be assigned after the candidate user stories have been finalized.

| Requirement | Use Case | Jira/User Story |
|---|---|---|
| FR-CAND-001 | UC-CAND-01 | TBD |
| FR-CAND-002 | UC-CAND-02 | TBD |
| FR-CAND-003 | UC-CAND-03 | TBD |
| FR-CAND-004 | UC-CAND-04 | TBD |
| FR-CAND-005 | UC-CAND-04 | TBD |
| FR-CAND-006 | UC-CAND-05 | TBD |
| FR-CAND-007 | UC-CAND-06 | TBD |
| FR-CAND-008 | UC-CAND-07 | TBD |
| FR-CAND-009 | UC-CAND-07 | TBD |
| FR-CAND-010 | UC-CAND-07 | TBD |
| FR-CAND-011 | UC-CAND-08 | TBD |
| FR-CAND-012 | UC-CAND-08 | TBD |
| FR-CAND-013 | UC-CAND-05 | TBD |
| FR-CAND-014 | UC-CAND-09 | TBD |
| FR-CAND-015 | UC-CAND-10 | TBD |
| FR-CAND-016 | UC-CAND-11 | TBD |
| FR-CAND-017 | UC-CAND-10 | TBD |
| FR-CAND-018 | UC-CAND-12 | TBD |
| FR-CAND-019 | UC-CAND-12 | TBD |
| FR-CAND-020 | UC-CAND-13 | TBD |
| FR-CAND-021 | UC-CAND-13 | TBD |
| FR-CAND-022 | UC-CAND-13 | TBD |
| FR-CAND-023 | UC-CAND-14 | TBD |
| FR-CAND-024 | UC-CAND-14 | TBD |
| FR-CAND-025 | UC-CAND-15 | TBD |