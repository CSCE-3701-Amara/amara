## 5.2 CAND: Candidates & Dynamic Resume

Owner: Ali Mohammed | Status: Draft | Last updated: 2026-10-08 | Jira: AMARA-34

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

User account creation, authentication, and account management are outside the scope of this subsystem.

---

### 5.2.2 Scope

The Candidate Subsystem includes the following functional areas:

1. Job Discovery and Filtering
2. Job Posting Details
3. Application History and Details
4. Assessment Progress Tracking
5. Assessment Results and Feedback
6. Application Status and Final Hiring Decisions
7. Offer Acceptance or Refusal
8. Company Following and Job Posting Bookmarks
9. Verified Skills and Crests
10. Dynamic Resume Generation and Viewing
11. Candidate Notifications

---

### 5.2.3 Functional Requirements

#### 5.2.3.1 Browse Job Postings

| ID | Requirement | Use Case | Related Functional Requirement | Jira |
|---|---|---|---|
| **FR-CAND-001** | The system **SHALL** display currently available job postings to candidates. | UC-CAND-001 — Browse Job Postings | FR-JOP-### | TBD |
| **FR-CAND-002** | The system **SHALL** display the essential information associated with each job posting. | UC-CAND-001 — Browse Job Postings | FR-JOP-### | TBD |
| **FR-CAND-003** | The system **SHALL** exclude job postings that are no longer available for application. | UC-CAND-001 — Browse Job Postings | FR-JOP-### | TBD |

#### 5.2.3.2 Filter Job Postings

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-004** | The system **SHALL** allow candidates to filter postings by offered compensation. | UC-CAND-002 — Filter Job Postings | FR-JOP-### | TBD |
| **FR-CAND-005** | The system **SHALL** allow candidates to filter postings by position. | UC-CAND-002 — Filter Job Postings | FR-JOP-### | TBD |
| **FR-CAND-006** | The system **SHALL** allow candidates to filter postings by job grade. | UC-CAND-002 — Filter Job Postings | FR-JOP-### | TBD |
| **FR-CAND-007** | The system **SHALL** allow candidates to filter postings by company. | UC-CAND-002 — Filter Job Postings | FR-JOP-### | TBD |
| **FR-CAND-008** | The system **SHALL** allow candidates to combine multiple filter criteria. | UC-CAND-002 — Filter Job Postings | FR-JOP-### | TBD |
| **FR-CAND-009** | The system **SHALL** update the displayed postings according to the selected criteria. | UC-CAND-002 — Filter Job Postings | FR-JOP-### | TBD |

#### 5.2.3.3 View Application History

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-010** | The system **SHALL** display the jobs to which the candidate has applied. | UC-CAND-003 — View Application History | TBD |
| **FR-CAND-011** | The system **SHALL** display the current status of each application. | UC-CAND-003 — View Application History | TBD |
| **FR-CAND-012** | The system **SHALL** allow the candidate to access the relevant application details. | UC-CAND-003 — View Application History | TBD |

#### 5.2.3.4 Track Assessment Progress

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-013** | The system **SHALL** indicate the candidate's assessment completion status. | UC-CAND-004 — Track Assessment Progress | TBD |
| **FR-CAND-014** | The system **SHALL** preserve assessment progress where the assessment permits interruption. | UC-CAND-004 — Track Assessment Progress | TBD |

#### 5.2.3.5 View Assessment Result and Feedback

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-015** | The system **SHALL** provide the candidate with feedback after the assessment has been evaluated, when feedback is available. | UC-CAND-005 — View Assessment Result and Feedback | TBD |
| **FR-CAND-016** | The system **SHALL** associate the feedback with the corresponding assessment. | UC-CAND-005 — View Assessment Result and Feedback | TBD |
| **FR-CAND-017** | The system **SHALL** make available any applicable areas of improvement identified by the assessment. | UC-CAND-005 — View Assessment Result and Feedback | TBD |
| **FR-CAND-018** | The system **SHALL** make the feedback accessible from the candidate's application. | UC-CAND-005 — View Assessment Result and Feedback | TBD |
| **FR-CAND-019** | The system **SHALL** indicate whether the candidate passed or failed the applicable assessment when the result is disclosed. | UC-CAND-005 — View Assessment Result and Feedback | TBD |
| **FR-CAND-020** | The system **SHALL** display the assessment result together with its associated feedback when available. | UC-CAND-005 — View Assessment Result and Feedback | TBD |

#### 5.2.3.6 View Application Status

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-021** | The system **SHALL** display the current status of each active application. | UC-CAND-006 — View Application Status | TBD |
| **FR-CAND-022** | The system **SHALL** update the application status when the application progresses through the hiring pipeline. | UC-CAND-006 — View Application Status | TBD |
| **FR-CAND-023** | The system **SHALL** indicate when an application has reached a final decision. | UC-CAND-006 — View Application Status | TBD |

#### 5.2.3.7 Receive Final Hiring Decision

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-024** | The system **SHALL** display the final decision once the employer has submitted it. | UC-CAND-007 — Receive Final Hiring Decision | TBD |
| **FR-CAND-025** | The system **SHALL** associate the decision with the corresponding application. | UC-CAND-007 — Receive Final Hiring Decision | TBD |
| **FR-CAND-026** | The system **SHALL** provide any final feedback made available by the employer. | UC-CAND-007 — Receive Final Hiring Decision | TBD |

#### 5.2.3.8 Accept Offer

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-027** | The candidate **SHALL** accept or refuse offer. | UC-CAND-008 — Accept Offer | TBD |
| **FR-CAND-028** | The system **SHALL** reflect accepting the offer in the candidate's profile. | UC-CAND-008 — Accept Offer | TBD |

#### 5.2.3.9 View Earned Crests

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-029** | The system **SHALL** display the Crests associated with the candidate. | UC-CAND-009 — View Earned Crests | TBD |
| **FR-CAND-030** | The system **SHALL** display the skill or competency represented by each Crest. | UC-CAND-009 — View Earned Crests | TBD |
| **FR-CAND-031** | The system **SHALL** display the relevant verification information associated with each Crest. | UC-CAND-009 — View Earned Crests | TBD |

#### 5.2.3.10 Maintain Verified Skills

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-032** | The system **SHALL** associate verified skills with the candidate. | UC-CAND-010 — Maintain Verified Skills | TBD |
| **FR-CAND-033** | The system **SHALL** retain verified skills across different eligible hiring processes. | UC-CAND-010 — Maintain Verified Skills | TBD |

#### 5.2.3.11 Generate Dynamic Resume

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-034** | The system **SHALL** include the candidate's verified skills in the dynamic resume. | UC-CAND-011 — Generate Dynamic Resume | TBD |
| **FR-CAND-035** | The system **SHALL** include the candidate's earned Crests in the dynamic resume. | UC-CAND-011 — Generate Dynamic Resume | TBD |
| **FR-CAND-036** | The system **SHALL** update the dynamic resume when relevant candidate information changes. | UC-CAND-011 — Generate Dynamic Resume | TBD |

#### 5.2.3.12 View Dynamic Resume

| ID | Requirement | Use Case | Jira |
|---|---|---|---|
| **FR-CAND-037** | The system **SHALL** display the candidate's verified skills. | UC-CAND-012 — View Dynamic Resume | TBD |
| **FR-CAND-038** | The system **SHALL** display the candidate's earned Crests. | UC-CAND-012 — View Dynamic Resume | TBD |

#### 5.2.3.13 Following Company

| ID | Requirement | Use Case | Related Functional Requirement | Jira |
|---|---|---|---|---|
| **FR-CAND-043** | The candidate **SHOULD** be able to follow a company for later postings. | UC-CAND-013 — Following Company | FR-JOP-### | TBD |

#### 5.2.3.14 Saving Posting

| ID | Requirement | Use Case | Related Functional Requirement | Jira |
|---|---|---|---|---|
| **FR-CAND-044** | The candidate **SHOULD** be able save a posting for later application. | UC-CAND-014 — Save Jop Posting | FR-JOP-### | TBD |

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