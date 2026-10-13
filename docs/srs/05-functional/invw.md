# 5.3 INVW: Pipeline & Interviews 

**Owner:** Ziad Eliwa | **Status:** In Review | **Last updated:** 2026-10-13 | **Jira:** AMARA-32

---

## 1. Overview

The INVW epic covers how Hiring Managers design and run a multi-stage interview pipeline shaped as a directed acyclic graph (DAG), and how Candidates are assessed, advanced, given feedback, or rejected. Its actors are the **Candidate**, **Interviewer**, **Hiring Manager**, and **System**. The Hiring Manager is the only human role that configures and administers the pipeline and records manual advancement decisions; Interviewers conduct assigned interviews and provide scorecards. This epic owns application status values and transitions and all notifications (see the boundary table in `requirements.md`).

## 2. Scope

**In scope**
- Pipeline design, assessment modules, live interviews, versioning and migration (3.1 to 3.4)
- Stage advancement, dashboard and insights, anonymity (3.5 to 3.7)
- Feedback and notifications (3.8)
- Assessment test execution and scoring (3.9)
- Candidate-facing pipeline progress view (3.10)
- Amara and Amarat skill crests (3.11)

**Out of scope (owned elsewhere)**
- Company registration and verification, roles, and the permission matrix (AUTH)
- The Candidate's platform resume and profile (CAND)
- Initial Candidate screening and the decision to admit a Candidate to the interview pipeline; a separate epic owns this stage. INVW begins after that epic passes a Candidate to the published pipeline.
- Hiring configuration of a job posting: deadlines, capacity, Must-have and Nice-to-have criteria, finalist target (3.1-JOB)
- Cross-posting and company-level analytics (ANA reads INVW data and never writes it)
- Contracts, onboarding, and payroll after an offer is accepted (Won't, future release)

---

## 3. Features and Functional Requirements

### 3.1 Pipeline Design

**Feature Description:** The Hiring Manager designs the interview pipeline of a posting as a DAG on a visual canvas, using stages connected by directed links.

**Preconditions:**
* The Hiring Manager is signed in and permitted to manage the posting (AUTH).
* The posting exists.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-001 | The system *SHALL* provide a canvas on which the Hiring Manager adds stages by drag and drop from a palette of module types, connects stages with directed links, and moves or deletes stages and links. | UC-INVW-001 |
| FR-INVW-002 | The system *SHALL* allow a stage to have several outgoing links, and *SHALL* make the next stages on all outgoing links available to a Candidate in parallel once the Candidate passes the stage. | UC-INVW-001 |
| FR-INVW-003 | The system *SHALL* allow a stage to have several incoming links, and *SHALL* make that stage available to a Candidate only after the Candidate has passed every stage linked into it. | UC-INVW-001 |
| FR-INVW-004 | The system *SHALL* allow a pipeline to have several entry stages (stages with no incoming link). | UC-INVW-001 |
| FR-INVW-005 | The system *SHALL* detect any cycle in a pipeline design and report the stages involved. | UC-INVW-001 |
| FR-INVW-006 | The system *SHALL* detect any stage that is not reachable from an entry stage, and report it. | UC-INVW-001 |
| FR-INVW-007 | The system *SHALL* detect any stage that has no unique name within the pipeline, no configured module, or no advancement rule (FR-INVW-051), and report it. | UC-INVW-001 |

### 3.2 Assessment Modules

**Feature Description:** A stage runs a module. The Hiring Manager chooses a predefined module, configures it, or builds a custom test in the Design Studio with AI assistance.

**Preconditions:**
* A pipeline design exists for the posting (3.1).
* For library-based modules, the problem library contains problems matching the chosen settings.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-008 | The system *SHALL* provide these individual assessment-module types, distinct from live interviews: coding problems, English proficiency exams, programming-language proficiency exams, asynchronous behavioral questions, video-recorded responses, and aptitude questions. | UC-INVW-002 |
| FR-INVW-009 | The system *SHALL* allow the Hiring Manager to configure a coding-problems module with the number of problems (at least 1), the difficulty of each problem or the difficulty mix, a maximum point value for each problem, and the selection mode: specific problems chosen from the library, or random problems drawn from the library. | UC-INVW-002 |
| FR-INVW-010 | The system *SHALL* give every Candidate in the same random-selection stage the same number of problems and the same difficulty mix. | UC-INVW-010 |
| FR-INVW-011 | The system *SHALL* allow the Hiring Manager to set an active time limit in minutes and a completion window in days for each individual assessment module; the completion window starts when the stage becomes available to the Candidate. Live-interview duration is configured separately (FR-INVW-023). | UC-INVW-002 |
| FR-INVW-012 | The system *SHALL* submit a Candidate's answers automatically when the module time limit expires. | UC-INVW-010 |
| FR-INVW-013 | The system *SHALL* email reminders to a Candidate with an unfinished assessment module 48, 24, 12, 6, and 1 hours before the applicable deadline (the earlier of the Candidate's completion-window deadline and the assessment's scheduled end time, as modified under FR-INVW-015), but only when the Candidate has not submitted and the reminder time falls before that deadline. | UC-INVW-010 |
| FR-INVW-014 | The system *SHALL* mark an assessment instance as missed when the Candidate's completion window or the assessment's scheduled end time expires without a started or submitted attempt, notify the Hiring Manager, and hold the Candidate for the Hiring Manager's decision (FR-INVW-057). | UC-INVW-010 |
| FR-INVW-015 | The system *SHALL* allow the Hiring Manager to grant a Candidate with a started assessment instance additional active time or a later individual completion deadline, including an extension beyond the assessment's scheduled end time, and record the change and reason. | UC-INVW-010 |
| FR-INVW-016 | The system *SHALL* record a numeric score from 0 to 100 for every fully scored assessment instance and live interview. | UC-INVW-004, UC-INVW-010 |
| FR-INVW-017 | The system *SHALL* allow the Hiring Manager to build a custom test module in a design studio using multiple-choice, short-text, long-text, and code questions. | UC-INVW-003 |
| FR-INVW-018 | The system *SHALL* offer AI assistance in the design studio that drafts questions from the Hiring Manager's description of the skill and difficulty. | UC-INVW-003 |
| FR-INVW-019 | The system *SHALL* not make AI-drafted questions available to Candidates until the Hiring Manager has reviewed and approved each one. | UC-INVW-003 |
| FR-INVW-020 | The system *SHALL* allow the Hiring Manager to assign points and an answer key to each closed question in an assessment module, and *SHALL* score those questions automatically. | UC-INVW-002, UC-INVW-003, UC-INVW-010 |
| FR-INVW-021 | The system *SHALL* allow the Hiring Manager to assign each open-ended written question to themselves or an Interviewer assigned to the stage, route each submitted answer to that grader, and limit awarded points to the question's configured maximum. | UC-INVW-002, UC-INVW-003, UC-INVW-010 |
| FR-INVW-022 | The system *SHALL* allow the Hiring Manager to save a custom module to the company's library for reuse in other postings. | UC-INVW-003 |

### 3.3 Live Interviews Module

**Feature Description:** A stage can be a live interview with an Interviewer or a Hiring Manager, held online through an external video-meeting provider. The platform schedules it and records who conducted it, when, and the notes and scorecard.

**Preconditions:**
* The stage is configured as a live interview with at least one assigned Interviewer or Hiring Manager.
* The external video-meeting provider is reachable (see EXT).
* The Candidate has reached the stage.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-023 | The system *SHALL* allow the Hiring Manager to add a live-interview stage, specifying its duration in minutes and one or more assigned Interviewers or Hiring Managers; this stage is not an individual assessment module. Its scheduled start and end times are set as part of the pipeline stage configuration (FR-INVW-094). | UC-INVW-004 |
| FR-INVW-024 | The system *SHALL* allow the Hiring Manager to define an evaluation rubric for a live-interview stage as a list of criteria, each scored from 1 to 5. | UC-INVW-004 |
| FR-INVW-025 | The system *SHALL* allow an assigned Interviewer to publish available time slots for a live interview. | UC-INVW-004 |
| FR-INVW-026 | The system *SHALL* allow the Candidate to select one published slot that starts and ends within the live-interview stage's scheduled start and end times, and *SHALL* then create the online meeting through the external provider and confirm the slot. | UC-INVW-004 |
| FR-INVW-027 | The system *SHALL* send the meeting time and join link to the Candidate and the assigned Interviewer by email and in the platform within 10 minutes of confirmation. | UC-INVW-004 |
| FR-INVW-028 | The system *SHALL* remind the Candidate and the Interviewer 24 hours and 1 hour before a live interview starts. | UC-INVW-004 |
| FR-INVW-029 | The system *SHALL* show every meeting time in the time zone of the person viewing it. | UC-INVW-004 |
| FR-INVW-030 | The system *SHALL* allow a Candidate to reschedule a live interview up to 2 times, no later than 24 hours before its start and within the live-interview stage's scheduled date range. | UC-INVW-004 |
| FR-INVW-031 | The system *SHALL* record for each live interview the assigned Interviewer or Interviewers, the scheduled start and end, and the outcome (held, Candidate absent, or Interviewer absent). | UC-INVW-004 |
| FR-INVW-032 | The system *SHALL* allow an assigned Interviewer to write and save timestamped notes on a live interview during and after it. | UC-INVW-004 |
| FR-INVW-033 | The system *SHALL* show Interviewer notes only to the Hiring Manager and the Interviewers assigned to that stage, never to the Candidate. | UC-INVW-004 |
| FR-INVW-034 | The system *SHALL* track whether each assigned Interviewer has submitted a rating for every rubric criterion and a written recommendation within 48 hours of the live interview, remind an Interviewer whose scorecard is missing, and notify the Hiring Manager when a scorecard remains missing at 48 hours. The Hiring Manager may waive a missing scorecard with a recorded reason. | UC-INVW-004 |
| FR-INVW-035 | The system *SHALL* calculate a live-interview score after every assigned scorecard is submitted or waived and at least one scorecard is submitted: it SHALL average all submitted rubric ratings and normalize the mean to 0–100 using ((mean rating − 1) / 4) × 100, rounded to the nearest whole number. If none is submitted, the System *SHALL* leave the score unset and hold the Candidate for the Hiring Manager's decision. | UC-INVW-004 |
| FR-INVW-036 | The system *SHALL* hold a Candidate for the Hiring Manager's decision and notify the Hiring Manager when a live interview outcome is "Candidate absent" or "Interviewer absent". | UC-INVW-004 |

### 3.4 Pipeline Versioning and Migration

**Feature Description:** A posting can have several saved pipeline versions, exactly one of which is published. When a new version is published while Candidates are in progress, the system moves them onto it without losing their results.

**Preconditions:**
* At least one saved pipeline version exists for the posting.
* For migration: Candidates are in progress on the currently published version.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-037 | The system *SHALL* allow the Hiring Manager to save multiple versions of a posting's pipeline, each with a sequential version number. | UC-INVW-005 |
| FR-INVW-038 | The system *SHALL* keep exactly one version published once a first version has been published, and *SHALL* unpublish the previously published version when another is published. | UC-INVW-005 |
| FR-INVW-039 | The system *SHALL* permit publication of a version only if it passes the checks in FR-INVW-005, FR-INVW-006, FR-INVW-007, FR-INVW-072, FR-INVW-082, FR-INVW-094, and FR-INVW-095, and *SHALL* list every failed check. | UC-INVW-005 |
| FR-INVW-040 | The system *SHALL* create a new draft version when the Hiring Manager edits the published version, and *SHALL* keep the published version in force until the new version is published. | UC-INVW-005 |
| FR-INVW-041 | The system *SHALL* record, for every application in the pipeline, the pipeline version it is running on and its current stage or stages. | UC-INVW-005 |
| FR-INVW-042 | The system *SHALL* show the Hiring Manager, for each version, the Candidates running on it and the stage each one is in. | UC-INVW-005 |
| FR-INVW-043 | The system *SHALL* move all in-progress Candidates to the newly published version when it is published. | UC-INVW-005 |
| FR-INVW-044 | The system *SHALL* show the Hiring Manager a migration preview before publication, giving for each stage the number of Candidates affected and where each group will be placed. | UC-INVW-005 |
| FR-INVW-045 | The system *SHALL* keep all completed stage results of a Candidate during migration, *SHALL* not require the Candidate to repeat a stage whose module and configuration are unchanged, and *SHALL* never move a Candidate to an earlier stage. | UC-INVW-005 |
| FR-INVW-046 | The system *SHALL* require the Hiring Manager to choose a destination stage in the new version for every removed stage that holds Candidates, and *SHALL* block publication until each is chosen. | UC-INVW-005 |
| FR-INVW-047 | The system *SHALL* require the Hiring Manager to decide, for each added stage, whether it also applies to in-progress Candidates who are already at or past a stage that follows it, or only to Candidates who have not yet reached that point. | UC-INVW-005 |
| FR-INVW-048 | The system *SHALL* re-evaluate Candidates waiting at a stage whose advancement rule or threshold changed under the new rule, and *SHALL* not reverse any pass already granted. | UC-INVW-005 |
| FR-INVW-049 | The system *SHALL* notify by email each Candidate who, because of a migration, must complete an additional stage. | UC-INVW-005 |
| FR-INVW-050 | The system *SHALL* record for each publication the Hiring Manager, the time, the previous and new version numbers, and the placement of each migrated Candidate. | UC-INVW-005 |

### 3.5 Stage Advancement

**Feature Description:** Each stage decides in one of three ways whether a Candidate advances: a Hiring Manager decision, an automatic score threshold, or automatic completion.

**Preconditions:**
* A pipeline is published and Candidates are in progress.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-051 | The system *SHALL* require every stage to have exactly one advancement rule: manual decision, automatic score threshold, or automatic on completion. | UC-INVW-006 |
| FR-INVW-052 | The system *SHALL*, under the manual rule, hold a Candidate at the stage until the Hiring Manager records a pass or a reject. | UC-INVW-006 |
| FR-INVW-053 | The system *SHALL* allow only the Hiring Manager to record manual pass or reject decisions; Interviewers provide scorecards only (FR-INVW-034). | UC-INVW-006 |
| FR-INVW-054 | The system *SHALL*, under the threshold rule, pass a Candidate automatically when the module score is at or above the stage's current threshold from 0 to 100 set by the Hiring Manager. | UC-INVW-006 |
| FR-INVW-055 | The system *SHALL* require the Hiring Manager to choose, for each threshold stage, whether a Candidate below the threshold is rejected automatically or held for the Hiring Manager's decision, with "held for decision" as the default. | UC-INVW-006 |
| FR-INVW-056 | The system *SHALL*, under the completion rule, pass a Candidate automatically when the module is completed, regardless of score. | UC-INVW-006 |
| FR-INVW-057 | The system *SHALL* allow the Hiring Manager to override an automatic result at a stage, requiring a written reason. | UC-INVW-006 |
| FR-INVW-058 | The system *SHALL* record for every stage decision its type (manual, threshold, completion, or override), the decider, the time, and the module score at that time. | UC-INVW-006 |
| FR-INVW-059 | The system *SHALL* allow the Hiring Manager to pass or reject several selected Candidates of one stage in a single action. | UC-INVW-006 |
| FR-INVW-060 | The system *SHALL*, when a Candidate is rejected at any stage, set the application status to "Rejected", remove the Candidate from all other active stages, and start rejection feedback (FR-INVW-080). | UC-INVW-006 |

### 3.6 Pipeline Dashboard and Insights

**Feature Description:** The Hiring Manager sees the pipeline graphically, with Candidate counts, rankings, and system-generated insights that support, but do not replace, human decisions.

**Preconditions:**
* A pipeline is published and at least one Candidate is active in a stage.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-061 | The system *SHALL* show the Hiring Manager a graphical view of the published pipeline with, for each stage, the number of Candidates currently in it and the numbers who passed and were rejected. | UC-INVW-007 |
| FR-INVW-062 | The system *SHALL* list the Candidates of a selected stage ranked by module score, with filters for score range, Nice-to-have score, and pending decision. | UC-INVW-007 |
| FR-INVW-063 | The system *SHALL* show, for each Candidate, the score at every completed stage and the percentile relative to the other Candidates at the same stage. | UC-INVW-007 |
| FR-INVW-064 | The system *SHALL* generate for each Candidate a summary of strengths and weaknesses based on performance across the sequence of completed modules, labelled as system-generated and shown with the scores it is based on. | UC-INVW-007 |
| FR-INVW-065 | The system *SHALL* exclude name, photo, gender, age, nationality, and contact details from the data used to rank Candidates and generate insights. | UC-INVW-007 |
| FR-INVW-066 | The system *SHALL* allow the Hiring Manager to record pass and reject decisions (FR-INVW-052) directly from the dashboard. | UC-INVW-007 |
| FR-INVW-067 | The system *SHALL* show the number of Candidates still active against the finalist target (FR-JOB-***). | UC-INVW-007 |
| FR-INVW-068 | The system *SHALL* suggest, for a selected stage, a cut-off score expected to bring the number of active Candidates to the finalist target, and show how many would pass at that cut-off. | UC-INVW-007 |
| FR-INVW-069 | The system *SHALL* allow the Hiring Manager to compare up to 4 selected Candidates side by side. | UC-INVW-007 |

### 3.7 Anonymity

**Feature Description:** The Hiring Manager can hide Candidate identity from hiring staff until a chosen point in the pipeline.

**Preconditions:**
* A pipeline is being designed (3.1).

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-070 | The system *SHALL* allow the Hiring Manager to enable anonymity for a pipeline and to select an identity-reveal point: a stage of the pipeline or the Finalist list. | UC-INVW-008 |
| FR-INVW-071 | The system *SHALL*, before a Candidate reaches the reveal point, replace the Candidate's name, photo, email, phone number, personal links, date of birth, gender, nationality, and address with a system-assigned anonymous identifier for Hiring Managers and Interviewers in every view, including resume copies, the dashboard, scorecards, and notes. | UC-INVW-008 |
| FR-INVW-072 | When anonymity is enabled, the system *SHALL* require the identity-reveal point to occur strictly before every live-interview or video-recorded response stage on every path, and *SHALL* reject a pipeline if any path reaches an incompatible stage at or before the reveal point. | UC-INVW-008 |
| FR-INVW-073 | The system *SHALL* reveal a Candidate's identity when the Candidate reaches the reveal point, *SHALL* not hide it again, and *SHALL* record the time of the reveal. | UC-INVW-008 |
| FR-INVW-074 | The system *SHALL* tell the Candidate, in the application view, that identity is hidden from the employer in the early stages, without naming the reveal point. | UC-INVW-008 |

### 3.8 Feedback and Notifications

**Feature Description:** Candidates receive feedback after each module, an email when they advance, and a detailed feedback report if they are rejected. Feedback can be written by a person, generated by AI, or both.

**Preconditions:**
* The Candidate has completed a module, advanced, or been rejected.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-075 | The system *SHALL* require the Hiring Manager to choose, for each stage, a feedback mode: human-written, AI-generated, or AI-generated with human approval before release (default). | UC-INVW-009 |
| FR-INVW-076 | The system *SHALL* release module feedback to the Candidate by email and in the platform within 10 minutes of module completion (AI-generated mode) or of the author's approval or submission (the other modes). | UC-INVW-009 |
| FR-INVW-077 | The system *SHALL* exclude from all feedback the test questions, answer keys, thresholds, and the data or identity of other Candidates. | UC-INVW-009 |
| FR-INVW-078 | The system *SHALL* label AI-generated feedback as AI-generated when shown to the Candidate. | UC-INVW-009 |
| FR-INVW-079 | The system *SHALL* email the Candidate within 10 minutes of being passed to a next stage, naming that stage. | UC-INVW-009 |
| FR-INVW-080 | The system *SHALL*, within 24 hours of a rejection, produce rejection feedback containing an overall summary, the outcome of each completed module with at least one strength or one area for improvement drawn from that module's results, and at least one recommended action. | UC-INVW-009 |
| FR-INVW-081 | The system *SHALL* require the Hiring Manager to choose a rejection-feedback mode per pipeline (the same three modes as FR-INVW-075) and *SHALL* release rejection feedback by email and in the platform according to it. | UC-INVW-009 |

### 3.9 Assessment Test Execution and Scoring

**Feature Description:** The system creates a stable Candidate-specific test instance from the published module configuration, enforces its time limits, records submissions, and produces a score after all required grading is complete.

**Preconditions:**
* The Candidate has reached an individual assessment stage in the published pipeline.
* The stage has a configured module, content, and scoring method.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-082 | For every individual assessment module, the system *SHALL* let the Hiring Manager select or author at least one problem, question, or prompt, configure its supported response mode and a scoring method (answer key, external code evaluator, or assigned human grader), and block pipeline publication with the stage and missing configuration identified if any item is absent. | UC-INVW-001, UC-INVW-002, UC-INVW-003, UC-INVW-005 |
| FR-INVW-083 | When an individual assessment stage becomes available, the system *SHALL* create one test instance for that Candidate and stage assignment from the published module configuration; reopening the assignment *SHALL* return the same instance and selected content. | UC-INVW-010 |
| FR-INVW-084 | For a timed test instance, the system *SHALL* start the active time limit when the Candidate first opens the instance, display the remaining time, and continue the timer while the Candidate navigates away or disconnects. | UC-INVW-010 |
| FR-INVW-085 | When a test instance's active time limit expires or its applicable completion deadline is reached while the instance is in progress and unsubmitted, the system *SHALL* submit the latest saved answers automatically. The applicable deadline is the earlier of the Candidate's completion-window deadline and the scheduled assessment end time, as modified by any individual extension recorded under FR-INVW-015. | UC-INVW-010 |
| FR-INVW-086 | For a point-based assessment module, the system *SHALL* calculate the module score as 100 multiplied by total points earned divided by total points available, rounded to the nearest whole number. | UC-INVW-010 |
| FR-INVW-087 | For each coding-problem answer, the system *SHALL* submit the Candidate's code and the problem's test cases to a configured external code-evaluation provider and use valid test-case outcomes to determine points under the problem's scoring method. If no valid result is returned, the system *SHALL* keep the instance pending and not assign a zero score. The provider interface and failure handling are specified in `docs/srs/06-external-interfaces.md`. | UC-INVW-010 |
| FR-INVW-088 | The system *SHALL* withhold answer keys, coding test cases, and grading rubrics from the Candidate's assessment view while the Candidate can still submit answers. | UC-INVW-010 |
| FR-INVW-089 | A test instance with open-ended written answers *SHALL* remain pending human grading until every answer has been graded; the system *SHALL NOT* mark it fully scored or evaluate an automatic advancement rule before grading is complete. | UC-INVW-010 |
| FR-INVW-090 | The system *SHALL* accept one final submission per test instance and make its responses uneditable after submission or automatic submission. | UC-INVW-010 |
| FR-INVW-091 | The system *SHALL* let the Hiring Manager configure question-based English, programming-language, asynchronous behavioral, and aptitude modules by setting the question count and selecting or authoring each question, its response format, and its scoring method. | UC-INVW-002, UC-INVW-010 |
| FR-INVW-092 | The system *COULD* provide an LLM-paired programming module that presents a configured coding task, lets the Candidate exchange text with the LLM during the active time limit, and saves the conversation and final code with the test instance. The LLM supports the Candidate and does not determine the assessment score. | UC-INVW-013 |
| FR-INVW-093 | For a video-recorded response module, the system *SHALL* present each configured prompt, let the Candidate record and submit a video response for each prompt within the time limits, and attach each recording to the test instance for assigned human grading. | UC-INVW-002, UC-INVW-010 |
| FR-INVW-094 | The system *SHALL* require the Hiring Manager to set a scheduled start and end time for every assessment stage, including live-interview stages, and *SHALL* require the end time to be later than the start time. Candidates may start an assessment or select a live-interview slot only within this interval; a started individual test may finish under its Candidate-specific deadline or an extension recorded under FR-INVW-015. | UC-INVW-001, UC-INVW-004, UC-INVW-005, UC-INVW-010 |
| FR-INVW-095 | The system *SHALL* prevent pipeline publication unless the scheduled start time of each first individual assessment stage on every path through the pipeline is later than the posting's application deadline (FR-JOB-***). Live-interview stages are not individual assessment stages for this check. | UC-INVW-005 |
| FR-INVW-096 | The system *SHALL* let the Hiring Manager apply a new score threshold, including the suggested cut-off in FR-INVW-068, to a threshold-based stage from the dashboard after reviewing Candidate performance insights; the system *SHALL* preview the affected Candidates, record the previous and new threshold and change time, and recalculate outcomes for Candidates still waiting at that stage. The system *SHALL NOT* reverse a previously granted pass or reopen a terminal decision. | UC-INVW-006, UC-INVW-007 |

### 3.10 Candidate Pipeline Progress View

**Feature Description:** Each Candidate can see the full published pipeline and their own progress through it, including the stage or stages they are currently in.

**Preconditions:**
* The Candidate's application has been passed from the separate initial-screening epic into a published INVW pipeline.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-097 | The system *SHALL* show each Candidate the complete stage graph for the published pipeline version their application is running on, identify their current stage or stages, distinguish their completed, active, and upcoming/locked stages, and identify Amarat awarded for passed stages (FR-INVW-103). The view *SHALL* show only that Candidate's progress and *SHALL NOT* expose other Candidates' data, Hiring Manager notes, or advancement thresholds. | UC-INVW-014 |

### 3.11 Amara and Amarat Skill Crests

**Feature Description:** An Amara is a named, personified crest for a defined skill; Amarat are the crests earned by a Candidate on passing an associated assessment stage. Each crest has configurable proficiency levels and evidence-based descriptions of strengths and areas for development.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-098 | The system *SHALL* associate every predefined individual assessment-module type with one or more built-in Amara definitions, each with a unique name, crest representation, skill domain, and description. | UC-INVW-002, UC-INVW-015 |
| FR-INVW-099 | For every custom test module and live-interview rubric, the system *SHALL* let the Hiring Manager select built-in Amarat or create custom Amara definitions with a unique name, crest representation, skill domain, and description; map each scored question or rubric criterion to exactly one selected Amara; require each assigned Amara to have at least one mapped item; and block publication if any mapping is missing. | UC-INVW-003, UC-INVW-004, UC-INVW-005, UC-INVW-015 |
| FR-INVW-100 | For each Amara, the system *SHOULD* provide at least two ordered levels with unique labels and non-overlapping score ranges that cover all scores from 0 to 100; each level *SHOULD* have a skill-specific description, strengths, and areas for development. Built-in level scales are system-defined; the Hiring Manager may configure the number, labels, ranges, and descriptions of levels for custom Amarat. Amara level bands are separate from stage advancement thresholds (FR-INVW-051, FR-INVW-054). | UC-INVW-015 |
| FR-INVW-101 | The system *SHOULD* calculate a separate 0–100 score for each Amara from only its mapped evidence, using earned/available points for point-scored items and the normalized ratings for mapped live-interview criteria (FR-INVW-035, FR-INVW-086), and assign the level whose configured band contains that score. | UC-INVW-016 |
| FR-INVW-102 | When a Candidate passes an assessment stage, the system *SHALL* award each Amara assigned to that stage; it *SHALL NOT* award Amarat for a failed or held stage. If no numeric level is available, the Candidate still receives the Amara crest without a level. | UC-INVW-006, UC-INVW-016 |
| FR-INVW-103 | The system *SHOULD* show each awarded Amara's crest, name, skill domain, level, level description, Candidate-specific strengths and areas for development, and the mapped assessment evidence supporting them. Descriptions *SHOULD* concern assessed skills and *SHOULD NOT* assert unsupported personality traits. | UC-INVW-014, UC-INVW-016 |
| FR-INVW-104 | The system *SHALL* retain each Amara award event with the Candidate, Amara definition, source stage, pass decision, and award time, even when the Candidate earns the same Amara again. | UC-INVW-016 |
| FR-INVW-105 | If the same Amara is awarded more than once, the system *SHOULD* show the highest earned level and preserve the history of all awards. It *SHOULD* retain the definition/level-scale version used for each award and not retroactively change prior award details when a custom definition is edited. | UC-INVW-016 |
| FR-INVW-106 | The system *SHALL* associate every custom Amara definition and every Amara award with the Company that owns the assessment's posting. Candidate views *SHALL* identify the Company for each Amara award, and Company Hiring Managers *SHALL* only access Amara awards for Candidates/applications they are authorized to manage (AUTH). | UC-INVW-015, UC-INVW-016 |

---
## 4. Business Rules

* **Pipeline responsibility:** The Hiring Manager configures and administers the pipeline and makes manual advancement decisions. Interviewers provide assigned notes and scorecards but do not decide advancement (AUTH; FR-INVW-032, FR-INVW-034, FR-INVW-053).
* **Ownership boundaries:** JOB owns posting criteria, including Must-have criterion freeze after the first application (FR-JOB-***). Initial screening and pipeline admission belong to the separate screening epic. Application status names and transitions are defined in the glossary and application state machine.
* **Scoring and advancement:** They are separate decisions. Closed questions use answer keys (FR-INVW-020); coding uses the external evaluator (FR-INVW-087); open-ended responses, video responses, and live interviews use human grading or scorecards (FR-INVW-021, FR-INVW-034, FR-INVW-093). AI drafting, insights, and feedback do not score Candidates or decide outcomes (FR-INVW-018, FR-INVW-064, FR-INVW-076). Advancement uses the rule configured for the stage (FR-INVW-051, FR-INVW-052, FR-INVW-054, FR-INVW-055, FR-INVW-056).
* **Assessment dates:** Each assessment stage has a valid scheduled start/end interval; the first exam on every pipeline path starts after the posting deadline (FR-INVW-039, FR-INVW-094, FR-INVW-095).
* **Pipeline versions:** Exactly one version is published at a time. Migration and preservation of Candidate progress follow FR-INVW-038, FR-INVW-045, FR-INVW-046, FR-INVW-047, and FR-INVW-048.
* **Identity and feedback:** Identity visibility follows the anonymity and reveal rules (FR-INVW-071, FR-INVW-072, FR-INVW-073). Feedback uses the configured mode and respects its content restrictions (FR-INVW-075, FR-INVW-077, FR-INVW-078, FR-INVW-081).

## 5. Error Handling

* **Invalid pipeline:** block publication and list every failed check with the stage or link concerned (FR-INVW-039).
* **Migration conflict:** block publication until every removed stage has a destination and every added stage has a decision (FR-INVW-046, FR-INVW-047).
* **Meeting creation failure:** retry 3 times within 15 minutes. If it still fails, keep the slot unconfirmed and notify the Candidate and Interviewer to pick another time.
* **Code-evaluation provider failure:** retain submitted code and keep the assessment pending; do not interpret an unavailable or invalid provider result as a failed test case (FR-INVW-087; provider details in `docs/srs/06-external-interfaces.md`).
* **Email delivery failure:** retry up to 3 times, show the notification in the platform regardless, and flag the failure to the Hiring Manager when the email concerns a decision.
* **AI generation failure (feedback, insights, design-studio drafts):** do not release partial output. Hold the item and notify the Hiring Manager to write it manually.
* **Two simultaneous decisions on the same Candidate and stage:** record the first, reject the second, and show the second decider the recorded decision.
* **Access to identity before the reveal point:** deny the request and log the attempt (FR-INVW-071).

---
