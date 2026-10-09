# 5.3 APP: Applications, Pipeline, Interviews & Offers

**Owner:** Ziad Eliwa | **Status:** In Review | **Last updated:** 2026-10-09 | **Jira:** AMARA-32

---

## 1. Overview

The APP epic covers what happens after a company publishes a job posting: how Interviewees apply, how applications are screened against Must-have criteria, how Hiring Managers design and run a multi-stage interview pipeline shaped as a directed acyclic graph (DAG), how Interviewees are evaluated, advanced, given feedback or rejected, and how Finalists receive and answer offers. Its actors are the **Hiring Manager**, the **Interviewer**, and the **Interviewee**. This epic owns application status values and transitions and all notifications (see the boundary table in `requirements.md`).

## 2. Scope

**In scope**
- Hiring configuration of a job posting: deadlines, capacity, Must-have and Nice-to-have criteria, finalist target (3.1)
- Application submission and initial screening (3.2, 3.3)
- Pipeline design, assessment modules, live interviews, versioning and migration (3.4 to 3.7)
- Stage advancement, dashboard and insights, anonymity (3.8 to 3.10)
- Feedback, notifications, Interviewee tracking, Finalists and offers (3.11 to 3.13)

**Out of scope (owned elsewhere)**
- Company registration and verification, roles, and the permission matrix (AUTH)
- The Interviewee's platform resume and profile (CAND)
- Posting lifecycle outside the hiring configuration in 3.1 (JOB)
- Cross-posting and company-level analytics (ANA reads APP data and never writes it)
- Contracts, onboarding, and payroll after an offer is accepted (Won't, future release)

---

## 3. Features and Functional Requirements

### 3.1 Job Posting Hiring Configuration

**Feature Description:** The Hiring Manager sets, per job posting, the parameters that control who can apply, how applicants are screened, and how many Interviewees should reach the end of the pipeline.

**Preconditions:**
* The company is registered and verified (AUTH).
* The Hiring Manager is signed in and permitted to manage postings of the company, per the AUTH permission matrix.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-01 | The system *SHALL* require the Hiring Manager to specify a job title, a job description, and the number of open positions (an integer of at least 1) for every job posting. | UC-APP-01 |
| FR-APP-02 | The system *SHALL* allow the Hiring Manager to specify a seniority level, employment type, work mode, location, and salary range for a job posting, and whether the salary range is shown to Interviewees. | UC-APP-01 |
| FR-APP-03 | The system *SHALL* require the Hiring Manager to specify an application opening date and an application deadline, with the deadline later than the opening date. | UC-APP-01 |
| FR-APP-04 | The system *SHALL* allow the Hiring Manager to extend the application deadline of a posting at any time before the deadline passes. | UC-APP-01 |
| FR-APP-05 | The system *SHALL* allow the Hiring Manager to set a maximum number of applications for a posting and *SHALL* stop accepting applications once that number has been submitted. | UC-APP-01 |
| FR-APP-06 | The system *SHALL* require the Hiring Manager to classify each criterion of a posting as exactly one of Must-have or Nice-to-have. | UC-APP-01 |
| FR-APP-07 | The system *SHALL* require each Must-have criterion to be a checkable condition on a structured resume field (minimum years of experience, required skill, minimum education level, minimum language proficiency level, or location eligibility) and *SHALL* not accept free-text Must-have criteria. | UC-APP-01 |
| FR-APP-08 | The system *SHALL* allow the Hiring Manager to assign each Nice-to-have criterion a weight from 1 to 5. | UC-APP-01 |
| FR-APP-09 | The system *SHALL* display the Must-have criteria and the Nice-to-have criteria of a posting to Interviewees in two separately labelled groups. | UC-APP-01 |
| FR-APP-10 | The system *SHALL* prevent changes to the Must-have criteria of a posting once the posting has received its first application. | UC-APP-01 |
| FR-APP-11 | The system *SHALL* allow the Hiring Manager to add application questions to a posting, each marked required or optional. | UC-APP-01 |
| FR-APP-12 | The system *SHALL* allow the Hiring Manager to set a finalist target for a posting, expressed either as a percentage (1 to 100) of the Interviewees admitted to the pipeline or as an integer multiple of the number of open positions. | UC-APP-01 |

### 3.2 Application Submission

**Feature Description:** An Interviewee applies to an open job posting using their platform resume and can later withdraw the application.

**Preconditions:**
* The Interviewee is signed in and has a complete platform resume (CAND).
* The posting is open: the current time is between its opening date and deadline, and its application cap (if set) has not been reached.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-13 | The system *SHALL* allow an Interviewee to submit an application to an open job posting using their platform resume, creating the application with status "Submitted". | UC-APP-02 |
| FR-APP-14 | The system *SHALL* require answers to all required application questions (FR-APP-11) before accepting a submission. | UC-APP-02 |
| FR-APP-15 | The system *SHALL* retain an unchangeable copy of the resume and answers as submitted, so that later resume edits do not alter the application. | UC-APP-02 |
| FR-APP-16 | The system *SHALL* reject a submission from an Interviewee who already has an active application, or a rejected or withdrawn application, for the same posting. | UC-APP-02 |
| FR-APP-17 | The system *SHALL* reject a submission made after the deadline or after the application cap is reached, and *SHALL* tell the Interviewee which of the two applies. | UC-APP-02 |
| FR-APP-18 | The system *SHALL* show the Interviewee a confirmation with an application reference and *SHALL* email the same confirmation within 10 minutes of submission. | UC-APP-02 |
| FR-APP-19 | The system *SHALL* allow an Interviewee to withdraw an application at any time before an offer is accepted or declined, setting its status to "Withdrawn" and removing the Interviewee from all active stages. | UC-APP-03 |

### 3.3 Initial Screening

**Feature Description:** Each submitted application is checked automatically against the posting's Must-have criteria. Applications that fail are screened out, and the rest are admitted to the pipeline.

**Preconditions:**
* An application with status "Submitted" exists.
* The posting has its Must-have and Nice-to-have criteria set (FR-APP-06, FR-APP-07).

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-20 | The system *SHALL* evaluate every Must-have criterion of the posting against the retained resume copy (FR-APP-15) within 60 seconds of submission. | UC-APP-04 |
| FR-APP-21 | The system *SHALL* set the status of an application that fails at least one Must-have criterion to "Screened Out" and record each failed criterion. | UC-APP-04 |
| FR-APP-22 | The system *SHALL* admit an application that satisfies all Must-have criteria to the pipeline. | UC-APP-04 |
| FR-APP-23 | The system *SHALL* keep an admitted application with status "Awaiting Pipeline" until a pipeline is published for the posting, and *SHALL* then start the Interviewee at every entry stage (FR-APP-31) without further action. | UC-APP-04 |
| FR-APP-24 | The system *SHALL* calculate a Nice-to-have score from 0 to 100 for each admitted application, as the weighted share of Nice-to-have criteria satisfied (FR-APP-08), and use it only for ordering and filtering. | UC-APP-04 |
| FR-APP-25 | The system *SHALL* notify an Interviewee who is screened out by email and in the platform within 10 minutes, listing the failed Must-have criteria. | UC-APP-04 |
| FR-APP-26 | The system *SHALL* allow the Hiring Manager to view all screened-out applications of a posting together with their failed criteria. | UC-APP-05 |
| FR-APP-27 | The system *SHALL* allow the Hiring Manager to override a screen-out and admit the application, requiring a written reason that is stored with the application. | UC-APP-05 |

### 3.4 Pipeline Design

**Feature Description:** The Hiring Manager designs the interview pipeline of a posting as a DAG on a visual canvas, using stages connected by directed links.

**Preconditions:**
* The Hiring Manager is signed in and permitted to manage the posting (AUTH).
* The posting exists.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-28 | The system *SHALL* provide a canvas on which the Hiring Manager adds stages by drag and drop from a palette of module types, connects stages with directed links, and moves or deletes stages and links. | UC-APP-06 |
| FR-APP-29 | The system *SHALL* allow a stage to have several outgoing links, and *SHALL* make the next stages on all outgoing links available to an Interviewee in parallel once the Interviewee passes the stage. | UC-APP-06 |
| FR-APP-30 | The system *SHALL* allow a stage to have several incoming links, and *SHALL* make that stage available to an Interviewee only after the Interviewee has passed every stage linked into it. | UC-APP-06 |
| FR-APP-31 | The system *SHALL* allow several entry stages (stages with no incoming link), all of which an admitted Interviewee starts in parallel. | UC-APP-06 |
| FR-APP-32 | The system *SHALL* detect any cycle in a pipeline design and report the stages involved. | UC-APP-06 |
| FR-APP-33 | The system *SHALL* detect any stage that is not reachable from an entry stage, or from which the Finalist list (FR-APP-117) cannot be reached, and report it. | UC-APP-06 |
| FR-APP-34 | The system *SHALL* detect any stage that has no unique name within the pipeline, no configured module, or no advancement rule (FR-APP-78), and report it. | UC-APP-06 |

### 3.5 Assessment Modules

**Feature Description:** A stage runs a module. The Hiring Manager chooses a predefined module, configures it, or builds a custom test in the Design Studio with AI assistance.

**Preconditions:**
* A pipeline design exists for the posting (3.4).
* For library-based modules, the problem library contains problems matching the chosen settings.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-35 | The system *SHALL* provide these predefined module types: coding problems, English proficiency exam, programming-language proficiency exam, LLM-paired programming session, behavioral interview, video-recorded response, and aptitude (IQ-style) questions. | UC-APP-07 |
| FR-APP-36 | The system *SHALL* allow the Hiring Manager to configure a coding-problems module with the number of problems (at least 1), the difficulty of each problem or the difficulty mix, and the selection mode: specific problems chosen from the library, or random problems drawn from the library. | UC-APP-07 |
| FR-APP-37 | The system *SHALL* give every Interviewee in the same random-selection stage the same number of problems and the same difficulty mix. | UC-APP-07 |
| FR-APP-38 | The system *SHALL* allow the Hiring Manager to set, for each module, a time limit in minutes and a completion window in days that starts when the stage becomes available to the Interviewee. | UC-APP-07 |
| FR-APP-39 | The system *SHALL* submit an Interviewee's answers automatically when the module time limit expires. | UC-APP-07 |
| FR-APP-40 | The system *SHALL* remind an Interviewee by email 48 hours and 24 hours before the completion window of an unfinished module ends. | UC-APP-07 |
| FR-APP-41 | The system *SHALL* mark a module as missed when its completion window ends without completion, *SHALL* notify the Hiring Manager, and *SHALL* hold the Interviewee for the Hiring Manager's decision (FR-APP-84). | UC-APP-07 |
| FR-APP-42 | The system *SHALL* allow the Hiring Manager to grant an individual Interviewee extra time or a longer completion window for a module and record the reason. | UC-APP-07 |
| FR-APP-43 | The system *SHALL* produce a score from 0 to 100 for every completed module and Interviewee. | UC-APP-07 |
| FR-APP-44 | The system *SHALL* allow the Hiring Manager to build a custom test module in a design studio using multiple-choice, short-text, long-text, and code questions. | UC-APP-08 |
| FR-APP-45 | The system *SHALL* offer AI assistance in the design studio that drafts questions from the Hiring Manager's description of the skill and difficulty. | UC-APP-08 |
| FR-APP-46 | The system *SHALL* not make AI-drafted questions available to Interviewees until the Hiring Manager has reviewed and approved each one. | UC-APP-08 |
| FR-APP-47 | The system *SHALL* allow the Hiring Manager to assign points and an answer key to each closed question of a custom module, and *SHALL* score those questions automatically. | UC-APP-08 |
| FR-APP-48 | The system *SHALL* route open-ended answers of a custom module to the Hiring Manager or an assigned Interviewer for grading against the points assigned. | UC-APP-08 |
| FR-APP-49 | The system *SHALL* allow the Hiring Manager to save a custom module to the company's library for reuse in other postings. | UC-APP-08 |

### 3.6 Live Interviews Module

**Feature Description:** A stage can be a live interview with an Interviewer or a Hiring Manager, held online through an external video-meeting provider. The platform schedules it and records who conducted it, when, and the notes and scorecard.

**Preconditions:**
* The stage is configured as a live interview with at least one assigned Interviewer or Hiring Manager.
* The external video-meeting provider is reachable (see EXT).
* The Interviewee has reached the stage.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-50 | The system *SHALL* allow the Hiring Manager to add a live-interview stage, specifying its duration in minutes and one or more assigned Interviewers or Hiring Managers. | UC-APP-09 |
| FR-APP-51 | The system *SHALL* allow the Hiring Manager to define an evaluation rubric for a live-interview stage as a list of criteria, each scored from 1 to 5. | UC-APP-09 |
| FR-APP-52 | The system *SHALL* allow an assigned Interviewer to publish available time slots for a live interview. | UC-APP-09 |
| FR-APP-53 | The system *SHALL* allow the Interviewee to select one published slot, and *SHALL* then create the online meeting through the external provider and confirm the slot. | UC-APP-09 |
| FR-APP-54 | The system *SHALL* send the meeting time and join link to the Interviewee and the assigned Interviewer by email and in the platform within 10 minutes of confirmation. | UC-APP-09 |
| FR-APP-55 | The system *SHALL* remind the Interviewee and the Interviewer 24 hours and 1 hour before a live interview starts. | UC-APP-09 |
| FR-APP-56 | The system *SHALL* show every meeting time in the time zone of the person viewing it. | UC-APP-09 |
| FR-APP-57 | The system *SHALL* allow an Interviewee to reschedule a live interview up to 2 times, no later than 24 hours before its start. | UC-APP-09 |
| FR-APP-58 | The system *SHALL* record for each live interview the assigned Interviewer or Interviewers, the scheduled start and end, and the outcome (held, Interviewee absent, or Interviewer absent). | UC-APP-09 |
| FR-APP-59 | The system *SHALL* allow an assigned Interviewer to write and save timestamped notes on a live interview during and after it. | UC-APP-09 |
| FR-APP-60 | The system *SHALL* show Interviewer notes only to the Hiring Manager and the Interviewers assigned to that stage, never to the Interviewee. | UC-APP-09 |
| FR-APP-61 | The system *SHALL* require each assigned Interviewer to submit a scorecard (rubric scores and a written recommendation) within 48 hours of the live interview, and *SHALL* remind the Interviewer if it is missing. | UC-APP-09 |
| FR-APP-62 | The system *SHALL* set the module score of a live interview (FR-APP-43) to the mean of the submitted scorecard scores, normalized to 0 to 100. | UC-APP-09 |
| FR-APP-63 | The system *SHALL* hold an Interviewee for the Hiring Manager's decision and notify the Hiring Manager when a live interview outcome is "Interviewee absent" or "Interviewer absent". | UC-APP-09 |

### 3.7 Pipeline Versioning and Migration

**Feature Description:** A posting can have several saved pipeline versions, exactly one of which is published. When a new version is published while Interviewees are in progress, the system moves them onto it without losing their results.

**Preconditions:**
* At least one saved pipeline version exists for the posting.
* For migration: Interviewees are in progress on the currently published version.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-64 | The system *SHALL* allow the Hiring Manager to save multiple versions of a posting's pipeline, each with a sequential version number. | UC-APP-10 |
| FR-APP-65 | The system *SHALL* keep exactly one version published once a first version has been published, and *SHALL* unpublish the previously published version when another is published. | UC-APP-10 |
| FR-APP-66 | The system *SHALL* permit publication of a version only if it passes the checks in FR-APP-32, FR-APP-33, FR-APP-34, and FR-APP-99, and *SHALL* list every failed check. | UC-APP-10 |
| FR-APP-67 | The system *SHALL* create a new draft version when the Hiring Manager edits the published version, and *SHALL* keep the published version in force until the new version is published. | UC-APP-10 |
| FR-APP-68 | The system *SHALL* record, for every application in the pipeline, the pipeline version it is running on and its current stage or stages. | UC-APP-10 |
| FR-APP-69 | The system *SHALL* show the Hiring Manager, for each version, the Interviewees running on it and the stage each one is in. | UC-APP-10 |
| FR-APP-70 | The system *SHALL* move all in-progress Interviewees to the newly published version when it is published. | UC-APP-10 |
| FR-APP-71 | The system *SHALL* show the Hiring Manager a migration preview before publication, giving for each stage the number of Interviewees affected and where each group will be placed. | UC-APP-10 |
| FR-APP-72 | The system *SHALL* keep all completed stage results of an Interviewee during migration, *SHALL* not require the Interviewee to repeat a stage whose module and configuration are unchanged, and *SHALL* never move an Interviewee to an earlier stage. | UC-APP-10 |
| FR-APP-73 | The system *SHALL* require the Hiring Manager to choose a destination stage in the new version for every removed stage that holds Interviewees, and *SHALL* block publication until each is chosen. | UC-APP-10 |
| FR-APP-74 | The system *SHALL* require the Hiring Manager to decide, for each added stage, whether it also applies to in-progress Interviewees who are already at or past a stage that follows it, or only to Interviewees who have not yet reached that point. | UC-APP-10 |
| FR-APP-75 | The system *SHALL* re-evaluate Interviewees waiting at a stage whose advancement rule or threshold changed under the new rule, and *SHALL* not reverse any pass already granted. | UC-APP-10 |
| FR-APP-76 | The system *SHALL* notify by email each Interviewee who, because of a migration, must complete an additional stage. | UC-APP-10 |
| FR-APP-77 | The system *SHALL* record for each publication the Hiring Manager, the time, the previous and new version numbers, and the placement of each migrated Interviewee. | UC-APP-10 |

### 3.8 Stage Advancement

**Feature Description:** Each stage decides in one of three ways whether an Interviewee advances: a Hiring Manager decision, an automatic score threshold, or automatic completion.

**Preconditions:**
* A pipeline is published and Interviewees are in progress.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-78 | The system *SHALL* require every stage to have exactly one advancement rule: manual decision, automatic score threshold, or automatic on completion. | UC-APP-11 |
| FR-APP-79 | The system *SHALL*, under the manual rule, hold an Interviewee at the stage until the Hiring Manager records a pass or a reject. | UC-APP-11 |
| FR-APP-80 | The system *SHALL* allow only the Hiring Manager to record manual pass or reject decisions; Interviewers provide scorecards only (FR-APP-61). | UC-APP-11 |
| FR-APP-81 | The system *SHALL*, under the threshold rule, pass an Interviewee automatically when the module score is at or above a threshold from 0 to 100 set by the Hiring Manager. | UC-APP-11 |
| FR-APP-82 | The system *SHALL* require the Hiring Manager to choose, for each threshold stage, whether an Interviewee below the threshold is rejected automatically or held for the Hiring Manager's decision, with "held for decision" as the default. | UC-APP-11 |
| FR-APP-83 | The system *SHALL*, under the completion rule, pass an Interviewee automatically when the module is completed, regardless of score. | UC-APP-11 |
| FR-APP-84 | The system *SHALL* allow the Hiring Manager to override an automatic result at a stage, requiring a written reason. | UC-APP-11 |
| FR-APP-85 | The system *SHALL* record for every stage decision its type (manual, threshold, completion, or override), the decider, the time, and the module score at that time. | UC-APP-11 |
| FR-APP-86 | The system *SHALL* allow the Hiring Manager to pass or reject several selected Interviewees of one stage in a single action. | UC-APP-11 |
| FR-APP-87 | The system *SHALL*, when an Interviewee is rejected at any stage, set the application status to "Rejected", remove the Interviewee from all other active stages, and start rejection feedback (FR-APP-108). | UC-APP-11 |

### 3.9 Pipeline Dashboard and Insights

**Feature Description:** The Hiring Manager sees the pipeline graphically, with Interviewee counts, rankings, and system-generated insights that support, but do not replace, human decisions.

**Preconditions:**
* A pipeline is published and at least one Interviewee has been admitted.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-88 | The system *SHALL* show the Hiring Manager a graphical view of the published pipeline with, for each stage, the number of Interviewees currently in it and the numbers who passed and were rejected. | UC-APP-12 |
| FR-APP-89 | The system *SHALL* list the Interviewees of a selected stage ranked by module score, with filters for score range, Nice-to-have score, and pending decision. | UC-APP-12 |
| FR-APP-90 | The system *SHALL* show, for each Interviewee, the score at every completed stage and the percentile relative to the other Interviewees at the same stage. | UC-APP-12 |
| FR-APP-91 | The system *SHALL* generate for each Interviewee a summary of strengths and weaknesses based on performance across the sequence of completed modules, labelled as system-generated and shown with the scores it is based on. | UC-APP-12 |
| FR-APP-92 | The system *SHALL* exclude name, photo, gender, age, nationality, and contact details from the data used to rank Interviewees and generate insights. | UC-APP-12 |
| FR-APP-93 | The system *SHALL* allow the Hiring Manager to record pass and reject decisions (FR-APP-79) directly from the dashboard. | UC-APP-12 |
| FR-APP-94 | The system *SHALL* show the number of Interviewees still active against the finalist target (FR-APP-12). | UC-APP-12 |
| FR-APP-95 | The system *SHALL* suggest, for a selected stage, a cut-off score expected to bring the number of active Interviewees to the finalist target, and show how many would pass at that cut-off. | UC-APP-12 |
| FR-APP-96 | The system *SHALL* allow the Hiring Manager to compare up to 4 selected Interviewees side by side. | UC-APP-12 |

### 3.10 Anonymity

**Feature Description:** The Hiring Manager can hide Interviewee identity from hiring staff until a chosen point in the pipeline.

**Preconditions:**
* A pipeline is being designed (3.4).

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-97 | The system *SHALL* allow the Hiring Manager to enable anonymity for a pipeline and to select an identity-reveal point: a stage of the pipeline or the Finalist list. | UC-APP-13 |
| FR-APP-98 | The system *SHALL*, before an Interviewee reaches the reveal point, replace the Interviewee's name, photo, email, phone number, personal links, date of birth, gender, nationality, and address with a system-assigned anonymous identifier for Hiring Managers and Interviewers in every view, including resume copies, the dashboard, scorecards, and notes. | UC-APP-13 |
| FR-APP-99 | The system *SHALL* mark the video-recorded response module and live interviews as incompatible with anonymity, and *SHALL* treat a pipeline as invalid if any path can reach an incompatible module before the reveal point. | UC-APP-13 |
| FR-APP-100 | The system *SHALL* reveal an Interviewee's identity when the Interviewee reaches the reveal point, *SHALL* not hide it again, and *SHALL* record the time of the reveal. | UC-APP-13 |
| FR-APP-101 | The system *SHALL* reveal an Interviewee's identity no later than the moment an offer is created for that Interviewee (FR-APP-119). | UC-APP-13 |
| FR-APP-102 | The system *SHALL* tell the Interviewee, in the application view, that identity is hidden from the employer in the early stages, without naming the reveal point. | UC-APP-13 |

### 3.11 Feedback and Notifications

**Feature Description:** Interviewees receive feedback after each module, an email when they advance, and a detailed feedback report if they are rejected. Feedback can be written by a person, generated by AI, or both.

**Preconditions:**
* The Interviewee has completed a module, advanced, or been rejected.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-103 | The system *SHALL* require the Hiring Manager to choose, for each stage, a feedback mode: human-written, AI-generated, or AI-generated with human approval before release (default). | UC-APP-14 |
| FR-APP-104 | The system *SHALL* release module feedback to the Interviewee by email and in the platform within 10 minutes of module completion (AI-generated mode) or of the author's approval or submission (the other modes). | UC-APP-14 |
| FR-APP-105 | The system *SHALL* exclude from all feedback the test questions, answer keys, thresholds, and the data or identity of other Interviewees. | UC-APP-14 |
| FR-APP-106 | The system *SHALL* label AI-generated feedback as AI-generated when shown to the Interviewee. | UC-APP-14 |
| FR-APP-107 | The system *SHALL* email the Interviewee within 10 minutes of being passed to a next stage, naming the stage. | UC-APP-14 |
| FR-APP-108 | The system *SHALL*, within 24 hours of a rejection, produce rejection feedback containing an overall summary, the outcome of each completed module with at least one strength or one area for improvement drawn from that module's results, and at least one recommended action. | UC-APP-14 |
| FR-APP-109 | The system *SHALL* require the Hiring Manager to choose a rejection-feedback mode per pipeline (the same three modes as FR-APP-103) and *SHALL* release rejection feedback by email and in the platform according to it. | UC-APP-14 |

### 3.12 Interviewee Application Tracking

**Feature Description:** An Interviewee follows their applications, current stage, feedback, and past applications. They see progress only, never the structure or content of the employer's interview design.

**Preconditions:**
* The Interviewee is signed in and has submitted at least one application.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-110 | The system *SHALL* show the Interviewee all active applications with the posting title, company, application status, and current stage name or names. | UC-APP-15 |
| FR-APP-111 | The system *SHALL* show the Interviewee past applications (no longer active) with their final outcome and dates. | UC-APP-15 |
| FR-APP-112 | The system *SHALL* show the Interviewee a graphical view of the application's pipeline displaying stage names, links, passed stages, and current stages. | UC-APP-15 |
| FR-APP-113 | The system *SHALL* not show the Interviewee any stage configuration (module settings, thresholds, advancement rules, question counts, difficulty, assigned Interviewers), any Interviewer notes, or any data of other Interviewees. | UC-APP-15 |
| FR-APP-114 | The system *SHALL* show the Interviewee all module feedback and rejection feedback released for the application. | UC-APP-15 |
| FR-APP-115 | The system *SHALL* show the Interviewee the action currently required (start a module, choose a slot, or join a meeting) with its due date or meeting time. | UC-APP-15 |
| FR-APP-116 | The system *SHALL* show the Interviewee a timeline of status changes with timestamps. | UC-APP-15 |

### 3.13 Finalists and Offers

**Feature Description:** Interviewees who complete the pipeline become Finalists. The Hiring Manager picks from them within the capacity of the job, sends offers, and rejects the rest. Offers are accepted or declined in the platform.

**Preconditions:**
* At least one Interviewee has passed every stage on all paths of the pipeline.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-APP-117 | The system *SHALL* set an Interviewee's status to "Finalist" and add them to the posting's Finalist list when they have passed every stage on all paths of the pipeline. | UC-APP-16 |
| FR-APP-118 | The system *SHALL* order the Finalist list by the mean module score of each Finalist by default, show each Finalist's insights (FR-APP-91), and allow ordering by the score of any single stage. | UC-APP-16 |
| FR-APP-119 | The system *SHALL* allow the Hiring Manager to create an offer for a Finalist containing the job title, compensation summary, proposed start date, response deadline, and an optional message. | UC-APP-16 |
| FR-APP-120 | The system *SHALL* deliver an offer to the Interviewee by email and in the platform. | UC-APP-16 |
| FR-APP-121 | The system *SHALL* allow the Interviewee to accept or decline a pending offer in the platform before its response deadline. | UC-APP-16 |
| FR-APP-122 | The system *SHALL* record the Interviewee's response with a timestamp, set the status to "Offer Accepted" or "Offer Declined", and notify the Hiring Manager within 10 minutes. | UC-APP-16 |
| FR-APP-123 | The system *SHALL* set an offer with no response at its deadline to "Offer Expired" and notify the Interviewee and the Hiring Manager. | UC-APP-16 |
| FR-APP-124 | The system *SHALL* allow the Hiring Manager to extend the response deadline of a pending offer. | UC-APP-16 |
| FR-APP-125 | The system *SHALL* allow the Hiring Manager to withdraw a pending offer and *SHALL* notify the Interviewee. | UC-APP-16 |
| FR-APP-126 | The system *SHALL* warn the Hiring Manager, and require explicit confirmation, when a new offer would make the number of pending and accepted offers exceed the number of open positions. | UC-APP-16 |
| FR-APP-127 | The system *SHALL* allow the Hiring Manager to reject Finalists individually or in bulk, and *SHALL* start rejection feedback (FR-APP-108) for each. | UC-APP-16 |
| FR-APP-128 | The system *SHALL* allow at most one pending offer per Interviewee per posting. | UC-APP-16 |

---

## 4. Business Rules

* A posting has exactly one published pipeline at a time. Interviewees are admitted but do not start any stage until one is published (FR-APP-23, FR-APP-65).
* Nice-to-have criteria never remove an applicant. Only a failed Must-have criterion screens an application out (FR-APP-21, FR-APP-24).
* Must-have criteria cannot change after the first application arrives (FR-APP-10).
* Application status values and their transitions are defined once, in the glossary and the application state machine, not in this section.
* Only Hiring Managers record manual pass or reject decisions. Interviewers contribute scorecards and notes (FR-APP-80).
* The system's insights support decisions but never replace them, except where the Hiring Manager configured an automatic advancement rule (FR-APP-78).
* Interviewees see progress only. Interview structure, content settings, thresholds, notes, and other Interviewees' data stay hidden (FR-APP-113).
* An Interviewee cannot apply again to a posting after being rejected or withdrawing (FR-APP-16).
* Migration never moves an Interviewee backward and never deletes results (FR-APP-72).
* Identity reveal is permanent and happens no later than the offer (FR-APP-100, FR-APP-101).
* Live-interview notes are never shown to Interviewees (FR-APP-60).

## 5. Error Handling

* **Late or capped submission:** reject with the specific reason (deadline passed or cap reached) and keep nothing partial (FR-APP-17).
* **Incomplete resume or missing required answers:** block submission and list what is missing (FR-APP-14).
* **Invalid pipeline:** block publication and list every failed check with the stage or link concerned (FR-APP-66).
* **Migration conflict:** block publication until every removed stage has a destination and every added stage has a decision (FR-APP-73, FR-APP-74).
* **Meeting creation failure:** retry 3 times within 15 minutes. If it still fails, keep the slot unconfirmed and notify the Interviewee and Interviewer to pick another time.
* **Email delivery failure:** retry up to 3 times, show the notification in the platform regardless, and flag the failure to the Hiring Manager when the email concerns a decision.
* **AI generation failure (feedback, insights, design-studio drafts):** do not release partial output. Hold the item and notify the Hiring Manager to write it manually.
* **Two simultaneous decisions on the same Interviewee and stage:** record the first, reject the second, and show the second decider the recorded decision.
* **Response to an expired, withdrawn, or already answered offer:** reject the response and show the current offer status.
* **Access to identity before the reveal point:** deny the request and log the attempt (FR-APP-98).

---
