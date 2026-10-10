## 4. Use cases

Owner: All | Status: Draft | Last updated: 2026-10-10 | Jira: AMARA-32, AMARA-

---

## 4.1 Actors

- Employer (Organization)
- Organization Administrator
- Interviewee (Candidate)
- Interviewers 
- Hiring Managers

## 4.2 AUTH Use Cases

## 4.3 CAND Use Cases

## 4.4 INVW Use Cases
### UC-INVW-001: Design Interview Pipeline

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-INVW-001 |
| Use Case Name     | Design Interview Pipeline |
| Actor(s)          | Hiring Manager |
| Description       | The Hiring Manager designs the interview pipeline of a posting as a directed acyclic graph (DAG) on a visual canvas. |
| Preconditions     | The Hiring Manager is signed in and permitted to manage the posting (AUTH). The posting exists. |
| Trigger           | The Hiring Manager opens the pipeline designer. |
| Main Flow         | 1. The Hiring Manager drags stages from the module palette onto the canvas and connects them with directed links.<br>2. The system enforces DAG semantics: several outgoing links make next stages available in parallel; several incoming links require all linked stages to be passed first.<br>3. The Hiring Manager configures each stage's module and advancement rule (FR-INVW-051).<br>4. The Hiring Manager saves the pipeline as a version. |
| Alternative Flows | 1. The Hiring Manager moves or deletes stages and links. |
| Exception Flows   | E1. A cycle is detected: the system reports the stages involved (FR-INVW-005).<br>E2. A stage is unreachable from an entry stage: reported (FR-INVW-006).<br>E3. A stage lacks a unique name, a configured module, or an advancement rule: reported (FR-INVW-007). |
| Postconditions    | A valid pipeline design is saved as a version; invalid designs are blocked with reasons. |
| Business Rules    | Several entry stages are allowed, all started in parallel (FR-INVW-004). A posting has exactly one published pipeline at a time. |

### UC-INVW-002: Configure Assessment Module

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-INVW-002 |
| Use Case Name     | Configure Assessment Module |
| Actor(s)          | Hiring Manager |
| Description       | The Hiring Manager configures a stage's assessment module, including type, parameters, timing, and scoring. |
| Preconditions     | A pipeline design exists for the posting (3.4). For library modules, the problem library contains matching problems. |
| Trigger           | The Hiring Manager selects and configures a module for a stage. |
| Main Flow         | 1. The Hiring Manager chooses a predefined module type (coding problems, English proficiency, programming-language proficiency, LLM-paired programming, behavioral interview, video-recorded response, or aptitude questions).<br>2. The Hiring Manager configures the module (e.g., number, difficulty, and selection mode for coding problems).<br>3. The Hiring Manager sets a time limit in minutes and a completion window in days.<br>4. The system produces a score from 0 to 100 for every completed module and Interviewee. |
| Alternative Flows | 2. Random selection mode: every Interviewee in the stage receives the same number of problems and difficulty mix. |
| Exception Flows   | E1. Time limit expires: the system auto-submits the Interviewee's answers (FR-INVW-012).<br>E2. Completion window ends without completion: the system marks the module missed, notifies the Hiring Manager, and holds the Interviewee for decision (FR-INVW-014).<br>E3. The system reminds the Interviewee 48h and 24h before the completion window ends (FR-INVW-013). |
| Postconditions    | The stage's module is configured; scoring (0-100) is enabled for completions. |
| Business Rules    | Extra time or a longer window may be granted to an individual Interviewee with a recorded reason (FR-INVW-015). |

### UC-INVW-003: Build Custom Test Module with AI Assistance

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-INVW-003 |
| Use Case Name     | Build Custom Test Module with AI Assistance |
| Actor(s)          | Hiring Manager, System (AI assistance) |
| Description       | The Hiring Manager builds a custom test module in the Design Studio, optionally using AI to draft questions, and saves it for reuse. |
| Preconditions     | A pipeline design exists for the posting (3.4). |
| Trigger           | The Hiring Manager opens the Design Studio to create a custom test module. |
| Main Flow         | 1. The Hiring Manager adds multiple-choice, short-text, long-text, and code questions.<br>2. The system optionally drafts questions from the Hiring Manager's description of the skill and difficulty (AI assistance).<br>3. The Hiring Manager reviews and approves each AI-drafted question.<br>4. The Hiring Manager assigns points and an answer key to each closed question (auto-scored).<br>5. The Hiring Manager optionally saves the module to the company library for reuse. |
| Alternative Flows | 1. Open-ended answers are routed to the Hiring Manager or an assigned Interviewer for grading against assigned points. |
| Exception Flows   | E1. AI-drafted questions are not made available to Interviewees until reviewed and approved (FR-INVW-019).<br>E2. AI generation fails: do not release partial output; notify the Hiring Manager to write manually. |
| Postconditions    | A custom module is configured and, optionally, saved to the company library. |
| Business Rules    | Closed questions are auto-scored against the answer key (FR-INVW-020). Open-ended answers require manual grading (FR-INVW-021). |

### UC-INVW-004: Schedule and Conduct Live Interview

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-INVW-004 |
| Use Case Name     | Schedule and Conduct Live Interview |
| Actor(s)          | Hiring Manager, Interviewer, Interviewee, external video-meeting provider (EXT) |
| Description       | The platform schedules a live-interview stage and records who conducted it, when, and the notes and scorecard. |
| Preconditions     | The stage is configured as a live interview with at least one assigned Interviewer or Hiring Manager. The external provider is reachable (EXT). The Interviewee has reached the stage. |
| Trigger           | An assigned Interviewer publishes available time slots. |
| Main Flow         | 1. The Interviewer publishes available time slots.<br>2. The Interviewee selects a slot; the system creates the online meeting via the external provider and confirms the slot.<br>3. The system sends the meeting time and join link to both parties within 10 minutes.<br>4. The system records for the interview the assigned Interviewer(s), scheduled start/end, and outcome (held, Interviewee absent, or Interviewer absent).<br>5. The Interviewer writes timestamped notes and submits a scorecard (rubric scores 1-5 and a written recommendation) within 48 hours.<br>6. The system sets the module score to the mean of submitted scorecard scores, normalized to 0-100. |
| Alternative Flows | 1. The system reminds both parties 24 hours and 1 hour before the interview.<br>5a. The system reminds an Interviewer if the scorecard is missing. |
| Exception Flows   | E1. The Interviewee reschedules up to 2 times, no later than 24 hours before the start (FR-INVW-030).<br>E2. Outcome is "Interviewee absent" or "Interviewer absent": the system holds the Interviewee for the Hiring Manager's decision and notifies them (FR-INVW-036).<br>E3. Meeting creation fails: retry 3 times within 15 minutes; if still failing, keep the slot unconfirmed and notify both parties to pick another time. |
| Postconditions    | The interview is recorded (assignees, timing, outcome); notes and scorecard are stored; a module score is set. |
| Business Rules    | Interviewer notes are shown only to the Hiring Manager and assigned Interviewers, never to the Interviewee (FR-INVW-033). Meeting times are shown in the viewer's time zone (FR-INVW-029). |

### UC-INVW-005: Publish and Migrate Pipeline Version

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-INVW-005 |
| Use Case Name     | Publish and Migrate Pipeline Version |
| Actor(s)          | Hiring Manager, Interviewee (notified) |
| Description       | The Hiring Manager publishes a new pipeline version; in-progress Interviewees are migrated onto it without losing results. |
| Preconditions     | At least one saved pipeline version exists. For migration, Interviewees are in progress on the currently published version. |
| Trigger           | The Hiring Manager publishes a pipeline version. |
| Main Flow         | 1. The system validates the version against checks FR-INVW-005, FR-INVW-006, FR-INVW-007, and FR-INVW-072, listing every failed check.<br>2. The system shows a migration preview (per stage, the number of Interviewees affected and where each group will be placed).<br>3. The Hiring Manager confirms publication; the system unpublishes the previous version and moves in-progress Interviewees to the new version, keeping completed results and never moving anyone backward.<br>4. The system records the publication (Hiring Manager, time, version numbers, placements) and notifies by email each Interviewee who must complete an additional stage. |
| Alternative Flows | 1. Editing the published version creates a new draft version while the published version stays in force. |
| Exception Flows   | E1. A removed stage holds Interviewees: the system blocks publication until a destination stage is chosen (FR-INVW-046).<br>E2. An added stage applies to in-progress Interviewees: the system blocks until the Hiring Manager decides its scope (FR-INVW-047).<br>E3. A stage's rule or threshold changed: waiting Interviewees are re-evaluated without reversing any already-granted pass (FR-INVW-048). |
| Postconditions    | Exactly one version is published; in-progress Interviewees are migrated with results preserved; the publication is recorded. |
| Business Rules    | Migration never moves an Interviewee backward and never deletes results (FR-INVW-045). A stage is not repeated when its module and configuration are unchanged. |

### UC-INVW-006: Advance Interviewee Between Stages

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-INVW-006 |
| Use Case Name     | Advance Interviewee Between Stages |
| Actor(s)          | Hiring Manager, System (automatic) |
| Description       | Each stage decides, by one of three advancement rules, whether an Interviewee passes to the next stage or is rejected. |
| Preconditions     | A pipeline is published and Interviewees are in progress. |
| Trigger           | A module completes or the Hiring Manager records a decision. |
| Main Flow         | 1. The system applies the stage's advancement rule: manual decision, automatic score threshold, or automatic on completion.<br>2. On a pass, the system moves the Interviewee to the next stage(s) and records the decision (type, decider, time, module score). |
| Alternative Flows | 1. Manual rule: only the Hiring Manager records pass or reject.<br>2b. Threshold rule: the system passes when the score is at or above the set threshold.<br>2c. Completion rule: the system passes when the module completes, regardless of score.<br>2d. The Hiring Manager passes or rejects several selected Interviewees of one stage in a single action. |
| Exception Flows   | E1. Below-threshold handling: rejected automatically or held for decision, per the Hiring Manager's choice (default held).<br>E2. On any rejection: status set to "Rejected", Interviewee removed from all other active stages, rejection feedback started (FR-INVW-060).<br>E3. The Hiring Manager overrides an automatic result with a written reason (FR-INVW-057).<br>E4. Two simultaneous decisions on the same Interviewee/stage: record the first, reject the second, and show the recorded decision. |
| Postconditions    | The stage decision is recorded; the Interviewee is moved to next stage(s) or rejected. |
| Business Rules    | Only Hiring Managers record manual decisions; Interviewers contribute scorecards only (FR-INVW-053). |

### UC-INVW-007: Review Pipeline Dashboard and Insights

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-INVW-007 |
| Use Case Name     | Review Pipeline Dashboard and Insights |
| Actor(s)          | Hiring Manager |
| Description       | The Hiring Manager reviews the pipeline graphically with counts, rankings, and system-generated insights that support, but do not replace, decisions. |
| Preconditions     | A pipeline is published and at least one Interviewee has been admitted. |
| Trigger           | The Hiring Manager opens the pipeline dashboard. |
| Main Flow         | 1. The system shows a graphical view with, per stage, the count of Interviewees in it and the numbers passed and rejected.<br>2. The Hiring Manager selects a stage; the system lists Interviewees ranked by module score with filters (score range, Nice-to-have score, pending decision).<br>3. The system shows each Interviewee's scores and percentile, plus a system-generated summary of strengths and weaknesses.<br>4. The Hiring Manager records pass/reject decisions directly from the dashboard. |
| Alternative Flows | 1. The Hiring Manager compares up to 4 selected Interviewees side by side.<br>4b. The system suggests a cut-off score to reach the finalist target and shows how many would pass. |
| Exception Flows   | E1. Ranking/insight data excludes name, photo, gender, age, nationality, and contact details (FR-INVW-065). |
| Postconditions    | The Hiring Manager has an updated view of pipeline progress and may record decisions. |
| Business Rules    | Insights support but never replace human decisions, except where an automatic rule is configured (FR-INVW-051). |

### UC-INVW-008: Enable Anonymous Hiring

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-INVW-008 |
| Use Case Name     | Enable Anonymous Hiring |
| Actor(s)          | Hiring Manager, Interviewer, Interviewee (informed) |
| Description       | The Hiring Manager hides Interviewee identity from hiring staff until a chosen reveal point in the pipeline. |
| Preconditions     | A pipeline is being designed (3.4). |
| Trigger           | The Hiring Manager enables anonymity on a pipeline. |
| Main Flow         | 1. The Hiring Manager enables anonymity and selects an identity-reveal point (a stage or the Finalist list).<br>2. Before the reveal point, the system replaces identity fields (name, photo, email, phone, personal links, date of birth, gender, nationality, address) with a system-assigned anonymous identifier in every view for Hiring Managers and Interviewers.<br>3. The system informs the Interviewee that identity is hidden in early stages, without naming the reveal point.<br>4. When an Interviewee reaches the reveal point, the system reveals identity permanently and records the time. |
| Alternative Flows | None. |
| Exception Flows   | E1. A path can reach a video-recorded response or live-interview stage before the reveal point: the pipeline is invalid (FR-INVW-072).<br>E2. Access to identity before the reveal point: deny and log the attempt (FR-INVW-071). |
| Postconditions    | Anonymity is active up to the reveal point; reveal is permanent and recorded. |
| Business Rules    | Identity reveal is permanent (FR-INVW-073). Anonymity is incompatible with video-recorded response and live-interview modules before the reveal point. |

### UC-INVW-009: Provide Feedback and Notifications

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-INVW-009 |
| Use Case Name     | Provide Feedback and Notifications |
| Actor(s)          | Hiring Manager, Interviewee, System (AI) |
| Description       | The system releases module and rejection feedback, and advance notifications, to Interviewees by email and in the platform. |
| Preconditions     | The Interviewee has completed a module, advanced, or been rejected. A feedback mode is chosen per stage and per pipeline. |
| Trigger           | A module completes, an advancement occurs, or a rejection is recorded. |
| Main Flow         | 1. The system releases module feedback within 10 minutes of completion (AI-generated mode) or of approval/submission (other modes).<br>2. The system emails the Interviewee within 10 minutes of being passed to a next stage, naming the stage.<br>3. Within 24 hours of a rejection, the system produces rejection feedback (overall summary, per-module outcome with at least one strength or improvement area, and at least one recommended action) and releases it per the pipeline mode. |
| Alternative Flows | 1. Feedback mode is human-written or AI-generated with human approval before release. |
| Exception Flows   | E1. Feedback excludes test questions, answer keys, thresholds, and other Interviewees' data or identity (FR-INVW-077).<br>E2. AI-generated feedback is labelled as AI-generated (FR-INVW-078).<br>E3. AI generation fails: hold the item and notify the Hiring Manager to write manually. |
| Postconditions    | Feedback and notifications are delivered by email and in the platform. |
| Business Rules    | Rejection feedback is required within 24 hours of rejection (FR-INVW-080). Feedback modes follow the same three options as FR-INVW-075. |

## 4.5 JOB Use Cases

## 4.6 ANA Use Cases
