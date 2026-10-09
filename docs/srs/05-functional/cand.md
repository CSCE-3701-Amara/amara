## 5.5 CAND: Candidates & Dynamic Resume

Owner: Ali Mohammed | Status: Draft | Last updated: 2026-10-08 | Jira: HIRE-n

# Candidate Functional Requirements

### 5.2.1 Overview

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

### 5.2.2 Scope

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

### 5.2.3 Functional Requirements

#### 5.2.3.1 Browse Job Postings

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-001** | The system **SHALL** display currently available job postings to candidates. | UC-CAND-01 — Browse Job Postings | TBD |
| **FR-CAND-002** | The system **SHALL** display the essential information associated with each job posting. | UC-CAND-01 — Browse Job Postings | TBD |
| **FR-CAND-003** | The system **SHALL** exclude job postings that are no longer available for application. | UC-CAND-01 — Browse Job Postings | TBD |

#### 5.2.3.2 Filter Job Postings

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-004** | The system **SHALL** allow candidates to filter postings by offered compensation. | UC-CAND-02 — Filter Job Postings | TBD |
| **FR-CAND-005** | The system **SHALL** allow candidates to filter postings by position. | UC-CAND-02 — Filter Job Postings | TBD |
| **FR-CAND-006** | The system **SHALL** allow candidates to filter postings by job grade. | UC-CAND-02 — Filter Job Postings | TBD |
| **FR-CAND-007** | The system **SHALL** allow candidates to filter postings by company. | UC-CAND-02 — Filter Job Postings | TBD |
| **FR-CAND-008** | The system **SHALL** allow candidates to combine multiple filter criteria. | UC-CAND-02 — Filter Job Postings | TBD |
| **FR-CAND-009** | The system **SHALL** update the displayed postings according to the selected criteria. | UC-CAND-02 — Filter Job Postings | TBD |

#### 5.2.3.3 View Job Posting Details

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-010** | The system **SHALL** display the position title. | UC-CAND-03 — View Job Posting Details | TBD |
| **FR-CAND-011** | The system **SHALL** display the company associated with the posting. | UC-CAND-03 — View Job Posting Details | TBD |
| **FR-CAND-012** | The system **SHALL** display the job grade. | UC-CAND-03 — View Job Posting Details | TBD |
| **FR-CAND-013** | The system **SHALL** display the offered compensation when provided. | UC-CAND-03 — View Job Posting Details | TBD |
| **FR-CAND-014** | The system **SHALL** display the job description. | UC-CAND-03 — View Job Posting Details | TBD |
| **FR-CAND-015** | The system **SHALL** display the required skills and qualifications. | UC-CAND-03 — View Job Posting Details | TBD |
| **FR-CAND-016** | The system **SHALL** display the relevant hiring process information when provided. | UC-CAND-03 — View Job Posting Details | TBD |

#### 5.2.3.4 Submit Job Application

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-017** | The system **SHALL** allow a candidate to initiate an application from an available job posting. | UC-CAND-04 — Submit Job Application | TBD |
| **FR-CAND-018** | The system **SHALL** associate the submitted application with the candidate and the selected job posting. | UC-CAND-04 — Submit Job Application | TBD |
| **FR-CAND-019** | The system **SHALL** record the date and time at which the application was submitted. | UC-CAND-04 — Submit Job Application | TBD |
| **FR-CAND-020** | The system **SHALL** confirm successful submission of the application to the candidate. | UC-CAND-04 — Submit Job Application | TBD |

#### 5.2.3.5 Prevent Duplicate Applications

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-021** | The system **SHALL** determine whether the candidate already has an active application for the selected job posting. | UC-CAND-04 — Prevent Duplicate Applications | TBD |
| **FR-CAND-022** | The system **SHALL** prevent submission of a duplicate active application. | UC-CAND-04 — Prevent Duplicate Applications | TBD |
| **FR-CAND-023** | The system **SHALL** inform the candidate when an active application already exists. | UC-CAND-04 — Prevent Duplicate Applications | TBD |

#### 5.2.3.6 View Application History

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-024** | The system **SHALL** display the jobs to which the candidate has applied. | UC-CAND-05 — View Application History | TBD |
| **FR-CAND-025** | The system **SHALL** display the current status of each application. | UC-CAND-05 — View Application History | TBD |
| **FR-CAND-026** | The system **SHALL** allow the candidate to access the relevant application details. | UC-CAND-05 — View Application History | TBD |

#### 5.2.3.7 View Hiring Pipeline

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-027** | The system **SHALL** display the stages applicable to the candidate's application. | UC-CAND-06 — View Hiring Pipeline | TBD |
| **FR-CAND-028** | The system **SHALL** indicate the candidate's current stage. | UC-CAND-06 — View Hiring Pipeline | TBD |
| **FR-CAND-029** | The system **SHALL** indicate the completion status of applicable stages. | UC-CAND-06 — View Hiring Pipeline | TBD |
| **FR-CAND-030** | The system **SHALL** indicate stages that are not yet available to the candidate. | UC-CAND-06 — View Hiring Pipeline | TBD |

#### 5.2.3.8 Access Assigned Assessment

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-031** | The system **SHALL** make an assessment available when the candidate reaches its corresponding pipeline stage. | UC-CAND-07 — Access Assigned Assessment | TBD |
| **FR-CAND-032** | The system **SHALL** prevent candidates from accessing assessments that are not yet available to them. | UC-CAND-07 — Access Assigned Assessment | TBD |
| **FR-CAND-033** | The system **SHALL** provide the candidate with the information necessary to begin the assessment. | UC-CAND-07 — Access Assigned Assessment | TBD |

#### 5.2.3.9 Complete Assessment

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-034** | The system **SHALL** present the assessment content to the candidate. | UC-CAND-07 — Complete Assessment | TBD |
| **FR-CAND-035** | The system **SHALL** record the candidate's submitted responses. | UC-CAND-07 — Complete Assessment | TBD |
| **FR-CAND-036** | The system **SHALL** evaluate the candidate's responses according to the assessment's defined evaluation criteria. | UC-CAND-07 — Complete Assessment | TBD |
| **FR-CAND-037** | The system **SHALL** record the outcome of the assessment. | UC-CAND-07 — Complete Assessment | TBD |
| **FR-CAND-038** | The system **SHALL** associate the assessment result with the candidate's application. | UC-CAND-07 — Complete Assessment | TBD |

#### 5.2.3.10 Track Assessment Progress

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-039** | The system **SHALL** indicate the candidate's assessment completion status. | UC-CAND-07 — Track Assessment Progress | TBD |
| **FR-CAND-040** | The system **SHALL** preserve assessment progress where the assessment permits interruption. | UC-CAND-07 — Track Assessment Progress | TBD |
| **FR-CAND-041** | The system **SHALL** indicate whether an assessment has been completed or remains incomplete. | UC-CAND-07 — Track Assessment Progress | TBD |

#### 5.2.3.11 View Assessment Feedback

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-042** | The system **SHALL** provide the candidate with feedback after the assessment has been evaluated, when feedback is available. | UC-CAND-08 — View Assessment Feedback | TBD |
| **FR-CAND-043** | The system **SHALL** associate the feedback with the corresponding assessment. | UC-CAND-08 — View Assessment Feedback | TBD |
| **FR-CAND-044** | The system **SHALL** make available any applicable areas of improvement identified by the assessment. | UC-CAND-08 — View Assessment Feedback | TBD |
| **FR-CAND-045** | The system **SHALL** make the feedback accessible from the candidate's application. | UC-CAND-08 — View Assessment Feedback | TBD |

#### 5.2.3.12 View Assessment Result

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-046** | The system **SHALL** indicate whether the candidate passed or failed the applicable assessment when the result is disclosed. | UC-CAND-08 — View Assessment Result | TBD |
| **FR-CAND-047** | The system **SHALL** display the assessment result together with its associated feedback when available. | UC-CAND-08 — View Assessment Result | TBD |

#### 5.2.3.13 View Application Status

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-048** | The system **SHALL** display the current status of each active application. | UC-CAND-05 — View Application Status | TBD |
| **FR-CAND-049** | The system **SHALL** update the application status when the application progresses through the hiring pipeline. | UC-CAND-05 — View Application Status | TBD |
| **FR-CAND-050** | The system **SHALL** indicate when an application has reached a final decision. | UC-CAND-05 — View Application Status | TBD |

#### 5.2.3.14 Receive Final Hiring Decision

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-051** | The system **SHALL** display the final decision once the employer has submitted it. | UC-CAND-09 — Receive Final Hiring Decision | TBD |
| **FR-CAND-052** | The system **SHALL** associate the decision with the corresponding application. | UC-CAND-09 — Receive Final Hiring Decision | TBD |
| **FR-CAND-053** | The system **SHALL** provide any final feedback made available by the employer. | UC-CAND-09 — Receive Final Hiring Decision | TBD |

#### 5.2.3.15 Award Verified Crest

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-054** | The system **SHALL** determine whether the candidate has successfully completed the applicable assessment. | UC-CAND-10 — Award Verified Crest | TBD |
| **FR-CAND-055** | The system **SHALL** award the corresponding Crest when the required criteria are satisfied. | UC-CAND-10 — Award Verified Crest | TBD |
| **FR-CAND-056** | The system **SHALL** associate the Crest with the skill or competency verified by the assessment. | UC-CAND-10 — Award Verified Crest | TBD |
| **FR-CAND-057** | The system **SHALL** associate the Crest with the candidate's profile. | UC-CAND-10 — Award Verified Crest | TBD |

#### 5.2.3.16 View Earned Crests

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-058** | The system **SHALL** display the Crests associated with the candidate. | UC-CAND-11 — View Earned Crests | TBD |
| **FR-CAND-059** | The system **SHALL** display the skill or competency represented by each Crest. | UC-CAND-11 — View Earned Crests | TBD |
| **FR-CAND-060** | The system **SHALL** display the relevant verification information associated with each Crest. | UC-CAND-11 — View Earned Crests | TBD |

#### 5.2.3.17 Maintain Verified Skills

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-061** | The system **SHALL** associate verified skills with the candidate. | UC-CAND-10 — Maintain Verified Skills | TBD |
| **FR-CAND-062** | The system **SHALL** retain verified skills across different eligible hiring processes. | UC-CAND-10 — Maintain Verified Skills | TBD |
| **FR-CAND-063** | The system **SHALL** update the candidate's verified skills when new eligible assessments are successfully completed. | UC-CAND-10 — Maintain Verified Skills | TBD |

#### 5.2.3.18 Generate Dynamic Resume

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-064** | The system **SHALL** include the candidate's verified skills in the dynamic resume. | UC-CAND-12 — Generate Dynamic Resume | TBD |
| **FR-CAND-065** | The system **SHALL** include the candidate's earned Crests in the dynamic resume. | UC-CAND-12 — Generate Dynamic Resume | TBD |
| **FR-CAND-066** | The system **SHALL** incorporate newly verified competencies into the dynamic resume. | UC-CAND-12 — Generate Dynamic Resume | TBD |
| **FR-CAND-067** | The system **SHALL** update the dynamic resume when relevant candidate information changes. | UC-CAND-12 — Generate Dynamic Resume | TBD |

#### 5.2.3.19 View Dynamic Resume

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-068** | The system **SHALL** display the candidate's professional information included in the resume. | UC-CAND-12 — View Dynamic Resume | TBD |
| **FR-CAND-069** | The system **SHALL** display the candidate's verified competencies. | UC-CAND-12 — View Dynamic Resume | TBD |
| **FR-CAND-070** | The system **SHALL** display the candidate's earned Crests. | UC-CAND-12 — View Dynamic Resume | TBD |
| **FR-CAND-071** | The system **SHALL** reflect the candidate's current verified skills. | UC-CAND-12 — View Dynamic Resume | TBD |

#### 5.2.3.20 Apply Candidate Anonymity

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-083** | The system **COULD** determine whether anonymity is enabled for the applicable hiring pipeline stage. | UC-CAND-13 — Apply Candidate Anonymity | TBD |
| **FR-CAND-084** | The system **COULD** hide candidate information designated as non-disclosable during anonymous stages. | UC-CAND-13 — Apply Candidate Anonymity | TBD |
| **FR-CAND-085** | The system **COULD** continue to expose information required for evaluating the candidate's professional qualifications. | UC-CAND-13 — Apply Candidate Anonymity | TBD |
| **FR-CAND-086** | The system **COULD** prevent hidden candidate information from being displayed to authorized evaluators during the anonymous stage. | UC-CAND-13 — Apply Candidate Anonymity | TBD |

#### 5.2.3.21 Reveal Candidate Identity

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-087** | The system **COULD** determine when the candidate reaches the configured identity-disclosure stage. | UC-CAND-14 — Reveal Candidate Identity | TBD |
| **FR-CAND-088** | The system **COULD** make the candidate's identity available to authorized users at that stage. | UC-CAND-14 — Reveal Candidate Identity | TBD |
| **FR-CAND-089** | The system **COULD** retain the candidate's application and assessment history after identity disclosure. | UC-CAND-14 — Reveal Candidate Identity | TBD |

#### 5.2.3.22 Receive Application Notifications

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-090** | The system **SHOULD** notify candidates when they advance to a new applicable hiring pipeline stage. | UC-CAND-15 — Receive Application Notifications | TBD |
| **FR-CAND-091** | The system **SHOULD** notify candidates when an assessment becomes available. | UC-CAND-15 — Receive Application Notifications | TBD |
| **FR-CAND-092** | The system **SHOULD** notify candidates when an assessment result becomes available. | UC-CAND-15 — Receive Application Notifications | TBD |
| **FR-CAND-093** | The system **SHOULD** notify candidates when a final hiring decision becomes available. | UC-CAND-15 — Receive Application Notifications | TBD |

---

### 5.2.4 Business Rules

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

### 5.2.5 Error Handling

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