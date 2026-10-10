# 5.3 INVW: Pipeline & Interviews 

**Owner:** Ziad Eliwa | **Status:** In Review | **Last updated:** 2026-10-10 | **Jira:** AMARA-32

---

## 1. Overview

The INVW epic covers what happens after a company publishes a job posting: how Hiring Managers design and run a multi-stage interview pipeline shaped as a directed acyclic graph (DAG), and how Interviewees are evaluated, advanced, given feedback, or rejected. Its actors are the **Hiring Manager**, the **Interviewer**, and the **Interviewee**. This epic owns application status values and transitions and all notifications (see the boundary table in `requirements.md`).

## 2. Scope

**In scope**
- Pipeline design, assessment modules, live interviews, versioning and migration (3.1 to 3.7)
- Stage advancement, dashboard and insights, anonymity (3.8 to 3.10)
- Feedback and notifications (3.11)

**Out of scope (owned elsewhere)**
- Company registration and verification, roles, and the permission matrix (AUTH)
- The Interviewee's platform resume and profile (CAND)
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
| FR-INVW-001 | The system *SHALL* provide a canvas on which the Hiring Manager adds stages by drag and drop from a palette of module types, connects stages with directed links, and moves or deletes stages and links. | UC-INVW-006 |
| FR-INVW-002 | The system *SHALL* allow a stage to have several outgoing links, and *SHALL* make the next stages on all outgoing links available to an Interviewee in parallel once the Interviewee passes the stage. | UC-INVW-006 |
| FR-INVW-003 | The system *SHALL* allow a stage to have several incoming links, and *SHALL* make that stage available to an Interviewee only after the Interviewee has passed every stage linked into it. | UC-INVW-006 |
| FR-INVW-004 | The system *SHALL* allow several entry stages (stages with no incoming link), all of which an admitted Interviewee starts in parallel. | UC-INVW-006 |
| FR-INVW-005 | The system *SHALL* detect any cycle in a pipeline design and report the stages involved. | UC-INVW-006 |
| FR-INVW-006 | The system *SHALL* detect any stage that is not reachable from an entry stage, and report it. | UC-INVW-006 |
| FR-INVW-007 | The system *SHALL* detect any stage that has no unique name within the pipeline, no configured module, or no advancement rule (FR-INVW-051), and report it. | UC-INVW-006 |

### 3.2 Assessment Modules

**Feature Description:** The Hiring Manager designs the interview pipeline of a posting as a DAG on a visual canvas, using stages connected by directed links.

**Preconditions:**
* The Hiring Manager is signed in and permitted to manage the posting (AUTH).
* The posting exists.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-001 | The system *SHALL* provide a canvas on which the Hiring Manager adds stages by drag and drop from a palette of module types, connects stages with directed links, and moves or deletes stages and links. | UC-INVW-001 |
| FR-INVW-002 | The system *SHALL* allow a stage to have several outgoing links, and *SHALL* make the next stages on all outgoing links available to an Interviewee in parallel once the Interviewee passes the stage. | UC-INVW-001 |
| FR-INVW-003 | The system *SHALL* allow a stage to have several incoming links, and *SHALL* make that stage available to an Interviewee only after the Interviewee has passed every stage linked into it. | UC-INVW-001 |
| FR-INVW-004 | The system *SHALL* allow several entry stages (stages with no incoming link), all of which an admitted Interviewee starts in parallel. | UC-INVW-001 |
| FR-INVW-005 | The system *SHALL* detect any cycle in a pipeline design and report the stages involved. | UC-INVW-001 |
| FR-INVW-006 | The system *SHALL* detect any stage that is not reachable from an entry stage, and report it. | UC-INVW-001 |
| FR-INVW-007 | The system *SHALL* detect any stage that has no unique name within the pipeline, no configured module, or no advancement rule (FR-INVW-051), and report it. | UC-INVW-001 |

### 3.5 Assessment Modules

**Feature Description:** A stage runs a module. The Hiring Manager chooses a predefined module, configures it, or builds a custom test in the Design Studio with AI assistance.

**Preconditions:**
* A pipeline design exists for the posting (3.4).
* For library-based modules, the problem library contains problems matching the chosen settings.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-008 | The system *SHALL* provide these predefined module types: coding problems, English proficiency exam, programming-language proficiency exam, LLM-paired programming session, behavioral interview, video-recorded response, and aptitude (IQ-style) questions. | UC-INVW-002 |
| FR-INVW-009 | The system *SHALL* allow the Hiring Manager to configure a coding-problems module with the number of problems (at least 1), the difficulty of each problem or the difficulty mix, and the selection mode: specific problems chosen from the library, or random problems drawn from the library. | UC-INVW-002 |
| FR-INVW-010 | The system *SHALL* give every Interviewee in the same random-selection stage the same number of problems and the same difficulty mix. | UC-INVW-002 |
| FR-INVW-011 | The system *SHALL* allow the Hiring Manager to set, for each module, a time limit in minutes and a completion window in days that starts when the stage becomes available to the Interviewee. | UC-INVW-002 |
| FR-INVW-012 | The system *SHALL* submit an Interviewee's answers automatically when the module time limit expires. | UC-INVW-002 |
| FR-INVW-013 | The system *SHALL* remind an Interviewee by email 48 hours and 24 hours before the completion window of an unfinished module ends. | UC-INVW-002 |
| FR-INVW-014 | The system *SHALL* mark a module as missed when its completion window ends without completion, *SHALL* notify the Hiring Manager, and *SHALL* hold the Interviewee for the Hiring Manager's decision (FR-INVW-057). | UC-INVW-002 |
| FR-INVW-015 | The system *SHALL* allow the Hiring Manager to grant an individual Interviewee extra time or a longer completion window for a module and record the reason. | UC-INVW-002 |
| FR-INVW-016 | The system *SHALL* produce a score from 0 to 100 for every completed module and Interviewee. | UC-INVW-002 |
| FR-INVW-017 | The system *SHALL* allow the Hiring Manager to build a custom test module in a design studio using multiple-choice, short-text, long-text, and code questions. | UC-INVW-003 |
| FR-INVW-018 | The system *SHALL* offer AI assistance in the design studio that drafts questions from the Hiring Manager's description of the skill and difficulty. | UC-INVW-003 |
| FR-INVW-019 | The system *SHALL* not make AI-drafted questions available to Interviewees until the Hiring Manager has reviewed and approved each one. | UC-INVW-003 |
| FR-INVW-020 | The system *SHALL* allow the Hiring Manager to assign points and an answer key to each closed question of a custom module, and *SHALL* score those questions automatically. | UC-INVW-003 |
| FR-INVW-021 | The system *SHALL* route open-ended answers of a custom module to the Hiring Manager or an assigned Interviewer for grading against the points assigned. | UC-INVW-003 |
| FR-INVW-022 | The system *SHALL* allow the Hiring Manager to save a custom module to the company's library for reuse in other postings. | UC-INVW-003 |

### 3.6 Live Interviews Module

**Feature Description:** A stage can be a live interview with an Interviewer or a Hiring Manager, held online through an external video-meeting provider. The platform schedules it and records who conducted it, when, and the notes and scorecard.

**Preconditions:**
* The stage is configured as a live interview with at least one assigned Interviewer or Hiring Manager.
* The external video-meeting provider is reachable (see EXT).
* The Interviewee has reached the stage.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-023 | The system *SHALL* allow the Hiring Manager to add a live-interview stage, specifying its duration in minutes and one or more assigned Interviewers or Hiring Managers. | UC-INVW-004 |
| FR-INVW-024 | The system *SHALL* allow the Hiring Manager to define an evaluation rubric for a live-interview stage as a list of criteria, each scored from 1 to 5. | UC-INVW-004 |
| FR-INVW-025 | The system *SHALL* allow an assigned Interviewer to publish available time slots for a live interview. | UC-INVW-004 |
| FR-INVW-026 | The system *SHALL* allow the Interviewee to select one published slot, and *SHALL* then create the online meeting through the external provider and confirm the slot. | UC-INVW-004 |
| FR-INVW-027 | The system *SHALL* send the meeting time and join link to the Interviewee and the assigned Interviewer by email and in the platform within 10 minutes of confirmation. | UC-INVW-004 |
| FR-INVW-028 | The system *SHALL* remind the Interviewee and the Interviewer 24 hours and 1 hour before a live interview starts. | UC-INVW-004 |
| FR-INVW-029 | The system *SHALL* show every meeting time in the time zone of the person viewing it. | UC-INVW-004 |
| FR-INVW-030 | The system *SHALL* allow an Interviewee to reschedule a live interview up to 2 times, no later than 24 hours before its start. | UC-INVW-004 |
| FR-INVW-031 | The system *SHALL* record for each live interview the assigned Interviewer or Interviewers, the scheduled start and end, and the outcome (held, Interviewee absent, or Interviewer absent). | UC-INVW-004 |
| FR-INVW-032 | The system *SHALL* allow an assigned Interviewer to write and save timestamped notes on a live interview during and after it. | UC-INVW-004 |
| FR-INVW-033 | The system *SHALL* show Interviewer notes only to the Hiring Manager and the Interviewers assigned to that stage, never to the Interviewee. | UC-INVW-004 |
| FR-INVW-034 | The system *SHALL* require each assigned Interviewer to submit a scorecard (rubric scores and a written recommendation) within 48 hours of the live interview, and *SHALL* remind the Interviewer if it is missing. | UC-INVW-004 |
| FR-INVW-035 | The system *SHALL* set the module score of a live interview (FR-INVW-016) to the mean of the submitted scorecard scores, normalized to 0 to 100. | UC-INVW-004 |
| FR-INVW-036 | The system *SHALL* hold an Interviewee for the Hiring Manager's decision and notify the Hiring Manager when a live interview outcome is "Interviewee absent" or "Interviewer absent". | UC-INVW-004 |

### 3.7 Pipeline Versioning and Migration

**Feature Description:** A posting can have several saved pipeline versions, exactly one of which is published. When a new version is published while Interviewees are in progress, the system moves them onto it without losing their results.

**Preconditions:**
* At least one saved pipeline version exists for the posting.
* For migration: Interviewees are in progress on the currently published version.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-037 | The system *SHALL* allow the Hiring Manager to save multiple versions of a posting's pipeline, each with a sequential version number. | UC-INVW-005 |
| FR-INVW-038 | The system *SHALL* keep exactly one version published once a first version has been published, and *SHALL* unpublish the previously published version when another is published. | UC-INVW-005 |
| FR-INVW-039 | The system *SHALL* permit publication of a version only if it passes the checks in FR-INVW-005, FR-INVW-006, FR-INVW-007, and FR-INVW-072, and *SHALL* list every failed check. | UC-INVW-005 |
| FR-INVW-040 | The system *SHALL* create a new draft version when the Hiring Manager edits the published version, and *SHALL* keep the published version in force until the new version is published. | UC-INVW-005 |
| FR-INVW-041 | The system *SHALL* record, for every application in the pipeline, the pipeline version it is running on and its current stage or stages. | UC-INVW-005 |
| FR-INVW-042 | The system *SHALL* show the Hiring Manager, for each version, the Interviewees running on it and the stage each one is in. | UC-INVW-005 |
| FR-INVW-043 | The system *SHALL* move all in-progress Interviewees to the newly published version when it is published. | UC-INVW-005 |
| FR-INVW-044 | The system *SHALL* show the Hiring Manager a migration preview before publication, giving for each stage the number of Interviewees affected and where each group will be placed. | UC-INVW-005 |
| FR-INVW-045 | The system *SHALL* keep all completed stage results of an Interviewee during migration, *SHALL* not require the Interviewee to repeat a stage whose module and configuration are unchanged, and *SHALL* never move an Interviewee to an earlier stage. | UC-INVW-005 |
| FR-INVW-046 | The system *SHALL* require the Hiring Manager to choose a destination stage in the new version for every removed stage that holds Interviewees, and *SHALL* block publication until each is chosen. | UC-INVW-005 |
| FR-INVW-047 | The system *SHALL* require the Hiring Manager to decide, for each added stage, whether it also applies to in-progress Interviewees who are already at or past a stage that follows it, or only to Interviewees who have not yet reached that point. | UC-INVW-005 |
| FR-INVW-048 | The system *SHALL* re-evaluate Interviewees waiting at a stage whose advancement rule or threshold changed under the new rule, and *SHALL* not reverse any pass already granted. | UC-INVW-005 |
| FR-INVW-049 | The system *SHALL* notify by email each Interviewee who, because of a migration, must complete an additional stage. | UC-INVW-005 |
| FR-INVW-050 | The system *SHALL* record for each publication the Hiring Manager, the time, the previous and new version numbers, and the placement of each migrated Interviewee. | UC-INVW-005 |

### 3.8 Stage Advancement

**Feature Description:** Each stage decides in one of three ways whether an Interviewee advances: a Hiring Manager decision, an automatic score threshold, or automatic completion.

**Preconditions:**
* A pipeline is published and Interviewees are in progress.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-051 | The system *SHALL* require every stage to have exactly one advancement rule: manual decision, automatic score threshold, or automatic on completion. | UC-INVW-006 |
| FR-INVW-052 | The system *SHALL*, under the manual rule, hold an Interviewee at the stage until the Hiring Manager records a pass or a reject. | UC-INVW-006 |
| FR-INVW-053 | The system *SHALL* allow only the Hiring Manager to record manual pass or reject decisions; Interviewers provide scorecards only (FR-INVW-034). | UC-INVW-006 |
| FR-INVW-054 | The system *SHALL*, under the threshold rule, pass an Interviewee automatically when the module score is at or above a threshold from 0 to 100 set by the Hiring Manager. | UC-INVW-006 |
| FR-INVW-055 | The system *SHALL* require the Hiring Manager to choose, for each threshold stage, whether an Interviewee below the threshold is rejected automatically or held for the Hiring Manager's decision, with "held for decision" as the default. | UC-INVW-006 |
| FR-INVW-056 | The system *SHALL*, under the completion rule, pass an Interviewee automatically when the module is completed, regardless of score. | UC-INVW-006 |
| FR-INVW-057 | The system *SHALL* allow the Hiring Manager to override an automatic result at a stage, requiring a written reason. | UC-INVW-006 |
| FR-INVW-058 | The system *SHALL* record for every stage decision its type (manual, threshold, completion, or override), the decider, the time, and the module score at that time. | UC-INVW-006 |
| FR-INVW-059 | The system *SHALL* allow the Hiring Manager to pass or reject several selected Interviewees of one stage in a single action. | UC-INVW-006 |
| FR-INVW-060 | The system *SHALL*, when an Interviewee is rejected at any stage, set the application status to "Rejected", remove the Interviewee from all other active stages, and start rejection feedback (FR-INVW-080). | UC-INVW-006 |

### 3.9 Pipeline Dashboard and Insights

**Feature Description:** The Hiring Manager sees the pipeline graphically, with Interviewee counts, rankings, and system-generated insights that support, but do not replace, human decisions.

**Preconditions:**
* A pipeline is published and at least one Interviewee has been admitted.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-061 | The system *SHALL* show the Hiring Manager a graphical view of the published pipeline with, for each stage, the number of Interviewees currently in it and the numbers who passed and were rejected. | UC-INVW-007 |
| FR-INVW-062 | The system *SHALL* list the Interviewees of a selected stage ranked by module score, with filters for score range, Nice-to-have score, and pending decision. | UC-INVW-007 |
| FR-INVW-063 | The system *SHALL* show, for each Interviewee, the score at every completed stage and the percentile relative to the other Interviewees at the same stage. | UC-INVW-007 |
| FR-INVW-064 | The system *SHALL* generate for each Interviewee a summary of strengths and weaknesses based on performance across the sequence of completed modules, labelled as system-generated and shown with the scores it is based on. | UC-INVW-007 |
| FR-INVW-065 | The system *SHALL* exclude name, photo, gender, age, nationality, and contact details from the data used to rank Interviewees and generate insights. | UC-INVW-007 |
| FR-INVW-066 | The system *SHALL* allow the Hiring Manager to record pass and reject decisions (FR-INVW-052) directly from the dashboard. | UC-INVW-007 |
| FR-INVW-067 | The system *SHALL* show the number of Interviewees still active against the finalist target (FR-JOB-012). | UC-INVW-007 |
| FR-INVW-068 | The system *SHALL* suggest, for a selected stage, a cut-off score expected to bring the number of active Interviewees to the finalist target, and show how many would pass at that cut-off. | UC-INVW-007 |
| FR-INVW-069 | The system *SHALL* allow the Hiring Manager to compare up to 4 selected Interviewees side by side. | UC-INVW-007 |

### 3.10 Anonymity

**Feature Description:** The Hiring Manager can hide Interviewee identity from hiring staff until a chosen point in the pipeline.

**Preconditions:**
* A pipeline is being designed (3.4).

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-070 | The system *SHALL* allow the Hiring Manager to enable anonymity for a pipeline and to select an identity-reveal point: a stage of the pipeline or the Finalist list. | UC-INVW-008 |
| FR-INVW-071 | The system *SHALL*, before an Interviewee reaches the reveal point, replace the Interviewee's name, photo, email, phone number, personal links, date of birth, gender, nationality, and address with a system-assigned anonymous identifier for Hiring Managers and Interviewers in every view, including resume copies, the dashboard, scorecards, and notes. | UC-INVW-008 |
| FR-INVW-072 | The system *SHALL* mark the video-recorded response module and live interviews as incompatible with anonymity, and *SHALL* treat a pipeline as invalid if any path can reach an incompatible module before the reveal point. | UC-INVW-008 |
| FR-INVW-073 | The system *SHALL* reveal an Interviewee's identity when the Interviewee reaches the reveal point, *SHALL* not hide it again, and *SHALL* record the time of the reveal. | UC-INVW-008 |
| FR-INVW-074 | The system *SHALL* tell the Interviewee, in the application view, that identity is hidden from the employer in the early stages, without naming the reveal point. | UC-INVW-008 |

### 3.11 Feedback and Notifications

**Feature Description:** Interviewees receive feedback after each module, an email when they advance, and a detailed feedback report if they are rejected. Feedback can be written by a person, generated by AI, or both.

**Preconditions:**
* The Interviewee has completed a module, advanced, or been rejected.

#### Functional Requirements

| ID | Requirement | Related Use Case |
|---|---|---|
| FR-INVW-075 | The system *SHALL* require the Hiring Manager to choose, for each stage, a feedback mode: human-written, AI-generated, or AI-generated with human approval before release (default). | UC-INVW-009 |
| FR-INVW-076 | The system *SHALL* release module feedback to the Interviewee by email and in the platform within 10 minutes of module completion (AI-generated mode) or of the author's approval or submission (the other modes). | UC-INVW-009 |
| FR-INVW-077 | The system *SHALL* exclude from all feedback the test questions, answer keys, thresholds, and the data or identity of other Interviewees. | UC-INVW-009 |
| FR-INVW-078 | The system *SHALL* label AI-generated feedback as AI-generated when shown to the Interviewee. | UC-INVW-009 |
| FR-INVW-079 | The system *SHALL* email the Interviewee within 10 minutes of being passed to a next stage, naming the stage. | UC-INVW-009 |
| FR-INVW-080 | The system *SHALL*, within 24 hours of a rejection, produce rejection feedback containing an overall summary, the outcome of each completed module with at least one strength or one area for improvement drawn from that module's results, and at least one recommended action. | UC-INVW-009 |
| FR-INVW-081 | The system *SHALL* require the Hiring Manager to choose a rejection-feedback mode per pipeline (the same three modes as FR-INVW-075) and *SHALL* release rejection feedback by email and in the platform according to it. | UC-INVW-009 |

---

## 4. Business Rules

* A posting has exactly one published pipeline at a time (FR-INVW-038).
* Must-have criteria cannot change after the first application arrives (FR-JOB-010).
* Application status values and their transitions are defined once, in the glossary and the application state machine, not in this section.
* Only Hiring Managers record manual pass or reject decisions. Interviewers contribute scorecards and notes (FR-INVW-053).
* The system's insights support decisions but never replace them, except where the Hiring Manager configured an automatic advancement rule (FR-INVW-051).
* Migration never moves an Interviewee backward and never deletes results (FR-INVW-045).
* Identity reveal is permanent (FR-INVW-073).
* Live-interview notes are never shown to Interviewees (FR-INVW-033).

## 5. Error Handling

* **Invalid pipeline:** block publication and list every failed check with the stage or link concerned (FR-INVW-039).
* **Migration conflict:** block publication until every removed stage has a destination and every added stage has a decision (FR-INVW-046, FR-INVW-047).
* **Meeting creation failure:** retry 3 times within 15 minutes. If it still fails, keep the slot unconfirmed and notify the Interviewee and Interviewer to pick another time.
* **Email delivery failure:** retry up to 3 times, show the notification in the platform regardless, and flag the failure to the Hiring Manager when the email concerns a decision.
* **AI generation failure (feedback, insights, design-studio drafts):** do not release partial output. Hold the item and notify the Hiring Manager to write it manually.
* **Two simultaneous decisions on the same Interviewee and stage:** record the first, reject the second, and show the second decider the recorded decision.
* **Access to identity before the reveal point:** deny the request and log the attempt (FR-INVW-071).

---
