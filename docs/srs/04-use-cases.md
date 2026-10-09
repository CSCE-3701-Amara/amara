## 4. Use cases

Owner: All | Status: Draft | Last updated: 2026-10-09 | Jira: AMARA-32, AMARA-

---

## 4.1 Actors

- Employer (Organization)
- Organization Administrator
- Interviewee (Candidate)
- Interviewers 
- Hiring Managers

## 4.2 AUTH Use Cases

## 4.3 CAND Use Cases

## 4.4 APP Use Cases

### UC-APP-01: Configure Job Posting Hiring Parameters

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-01 |
| Use Case Name     | Configure Job Posting Hiring Parameters |
| Actor(s)          | Hiring Manager |
| Description       | The Hiring Manager sets the posting parameters that control who can apply, how applications are screened, and how many Interviewees should reach the end of the pipeline. |
| Preconditions     | The company is registered and verified (AUTH). The Hiring Manager is signed in and permitted to manage the company's postings (AUTH). |
| Trigger           | The Hiring Manager opens the hiring configuration of a job posting. |
| Main Flow         | 1. The Hiring Manager specifies the job title, job description, and number of open positions (integer >= 1).<br>2. The Hiring Manager optionally specifies seniority level, employment type, work mode, location, and salary range, and whether the range is shown to Interviewees.<br>3. The Hiring Manager specifies an application opening date and a deadline later than the opening date.<br>4. The Hiring Manager optionally sets a maximum number of applications.<br>5. The Hiring Manager classifies each criterion as Must-have or Nice-to-have, keeping Must-have criteria as checkable structured-resume conditions.<br>6. The Hiring Manager assigns each Nice-to-have criterion a weight from 1 to 5.<br>7. The Hiring Manager optionally adds application questions, each marked required or optional.<br>8. The Hiring Manager sets a finalist target as a percentage or an integer multiple of open positions.<br>9. The system saves the configuration and displays Must-have and Nice-to-have criteria in two separately labelled groups. |
| Alternative Flows | 1. The Hiring Manager may extend the application deadline at any time before it passes.<br>8a. The finalist target is expressed either as a percentage (1-100) of admitted Interviewees or as an integer multiple of open positions. |
| Exception Flows   | E1. The system rejects a deadline not later than the opening date.<br>E2. The system rejects a non-integer or zero/negative number of open positions.<br>E3. The system rejects a free-text Must-have criterion (FR-APP-07).<br>E4. The system rejects a Nice-to-have weight outside 1 to 5. |
| Postconditions    | The posting has a valid hiring configuration; Must-have and Nice-to-have criteria are stored and displayed separately. |
| Business Rules    | Each criterion is exactly Must-have or Nice-to-have (FR-APP-06). Must-have criteria cannot change after the first application arrives (FR-APP-10). Nice-to-have criteria never remove an applicant. |

### UC-APP-02: Apply to Job Posting

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-02 |
| Use Case Name     | Apply to Job Posting |
| Actor(s)          | Interviewee |
| Description       | An Interviewee submits an application to an open job posting using their platform resume and required question answers. |
| Preconditions     | The Interviewee is signed in and has a complete platform resume (CAND). The posting is open (between opening date and deadline) and its application cap, if set, has not been reached. |
| Trigger           | The Interviewee clicks "Apply" on an open posting. |
| Main Flow         | 1. The system presents the posting's Must-have and Nice-to-have criteria and application questions.<br>2. The Interviewee answers all required questions and submits using their platform resume.<br>3. The system creates the application with status "Submitted" and retains an unchangeable copy of the resume and answers.<br>4. The system shows a confirmation with an application reference and emails the same confirmation within 10 minutes. |
| Alternative Flows | 1. The Interviewee answers optional questions before submitting. |
| Exception Flows   | E1. Required questions are unanswered: the system blocks submission and lists what is missing (FR-APP-14).<br>E2. The Interviewee already has an active, rejected, or withdrawn application for the posting: submission is rejected (FR-APP-16).<br>E3. Submission after the deadline or after the cap: the system rejects and tells the Interviewee which applies (FR-APP-17).<br>E4. Email delivery fails: retry up to 3 times and show the confirmation in the platform regardless. |
| Postconditions    | The application exists with status "Submitted"; a frozen copy of the resume and answers is stored. |
| Business Rules    | An Interviewee cannot apply again after being rejected or withdrawing (FR-APP-16). Later resume edits do not alter a submitted application (FR-APP-15). |

### UC-APP-03: Withdraw Application

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-03 |
| Use Case Name     | Withdraw Application |
| Actor(s)          | Interviewee |
| Description       | An Interviewee withdraws a submitted application, leaving the pipeline and all active stages. |
| Preconditions     | The Interviewee has an active application to the posting. |
| Trigger           | The Interviewee chooses "Withdraw application". |
| Main Flow         | 1. The Interviewee confirms the withdrawal.<br>2. The system sets the application status to "Withdrawn" and removes the Interviewee from all active stages. |
| Alternative Flows | None. |
| Exception Flows   | E1. Withdrawal after an offer is accepted or declined: the system rejects the withdrawal request (FR-APP-19). |
| Postconditions    | The application status is "Withdrawn" and the Interviewee is removed from all active stages. |
| Business Rules    | Withdrawal is allowed only before an offer is accepted or declined (FR-APP-19). |

### UC-APP-04: Automatically Screen Application

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-04 |
| Use Case Name     | Automatically Screen Application |
| Actor(s)          | System (automatic), Interviewee (notified) |
| Description       | The system evaluates a submitted application against the posting's Must-have criteria, screens out failures, and admits the rest to the pipeline. |
| Preconditions     | An application with status "Submitted" exists. The posting has Must-have and Nice-to-have criteria set (FR-APP-06, FR-APP-07). |
| Trigger           | A new application is submitted. |
| Main Flow         | 1. Within 60 seconds of submission, the system evaluates every Must-have criterion against the retained resume copy.<br>2. If all Must-have criteria are satisfied, the system admits the application, computes the Nice-to-have score (0-100), and sets status "Awaiting Pipeline".<br>3. When a pipeline is published, the system starts the Interviewee at every entry stage (FR-APP-31). |
| Alternative Flows | 1. Before a pipeline is published, the application remains "Awaiting Pipeline". |
| Exception Flows   | E1. At least one Must-have criterion fails: the system sets status "Screened Out", records each failed criterion, and notifies the Interviewee by email and in the platform within 10 minutes (FR-APP-21, FR-APP-25). |
| Postconditions    | The application is either "Screened Out" (with failed criteria recorded) or admitted with a Nice-to-have score. |
| Business Rules    | Only a failed Must-have criterion screens an application out (FR-APP-21). The Nice-to-have score is used only for ordering and filtering (FR-APP-24). |

### UC-APP-05: Override Screen-Out

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-05 |
| Use Case Name     | Override Screen-Out |
| Actor(s)          | Hiring Manager |
| Description       | The Hiring Manager reviews screened-out applications and manually admits one, recording a reason. |
| Preconditions     | At least one application has status "Screened Out". |
| Trigger           | The Hiring Manager opens the list of screened-out applications. |
| Main Flow         | 1. The system shows all screened-out applications of the posting with their failed criteria.<br>2. The Hiring Manager selects an application and chooses to override the screen-out.<br>3. The system requires a written reason.<br>4. The Hiring Manager enters the reason and confirms; the system admits the application and stores the reason with it. |
| Alternative Flows | None. |
| Exception Flows   | E1. No reason entered: the system blocks the override. |
| Postconditions    | The application is admitted (no longer "Screened Out") and the override reason is stored. |
| Business Rules    | Overrides require a written reason stored with the application (FR-APP-27). |

### UC-APP-06: Design Interview Pipeline

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-06 |
| Use Case Name     | Design Interview Pipeline |
| Actor(s)          | Hiring Manager |
| Description       | The Hiring Manager designs the interview pipeline of a posting as a directed acyclic graph (DAG) on a visual canvas. |
| Preconditions     | The Hiring Manager is signed in and permitted to manage the posting (AUTH). The posting exists. |
| Trigger           | The Hiring Manager opens the pipeline designer. |
| Main Flow         | 1. The Hiring Manager drags stages from the module palette onto the canvas and connects them with directed links.<br>2. The system enforces DAG semantics: several outgoing links make next stages available in parallel; several incoming links require all linked stages to be passed first.<br>3. The Hiring Manager configures each stage's module and advancement rule (FR-APP-78).<br>4. The Hiring Manager saves the pipeline as a version. |
| Alternative Flows | 1. The Hiring Manager moves or deletes stages and links. |
| Exception Flows   | E1. A cycle is detected: the system reports the stages involved (FR-APP-32).<br>E2. A stage is unreachable from an entry stage or cannot reach the Finalist list: reported (FR-APP-33).<br>E3. A stage lacks a unique name, a configured module, or an advancement rule: reported (FR-APP-34). |
| Postconditions    | A valid pipeline design is saved as a version; invalid designs are blocked with reasons. |
| Business Rules    | Several entry stages are allowed, all started in parallel (FR-APP-31). A posting has exactly one published pipeline at a time. |

### UC-APP-07: Configure Assessment Module

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-07 |
| Use Case Name     | Configure Assessment Module |
| Actor(s)          | Hiring Manager |
| Description       | The Hiring Manager configures a stage's assessment module, including type, parameters, timing, and scoring. |
| Preconditions     | A pipeline design exists for the posting (3.4). For library modules, the problem library contains matching problems. |
| Trigger           | The Hiring Manager selects and configures a module for a stage. |
| Main Flow         | 1. The Hiring Manager chooses a predefined module type (coding problems, English proficiency, programming-language proficiency, LLM-paired programming, behavioral interview, video-recorded response, or aptitude questions).<br>2. The Hiring Manager configures the module (e.g., number, difficulty, and selection mode for coding problems).<br>3. The Hiring Manager sets a time limit in minutes and a completion window in days.<br>4. The system produces a score from 0 to 100 for every completed module and Interviewee. |
| Alternative Flows | 2. Random selection mode: every Interviewee in the stage receives the same number of problems and difficulty mix. |
| Exception Flows   | E1. Time limit expires: the system auto-submits the Interviewee's answers (FR-APP-39).<br>E2. Completion window ends without completion: the system marks the module missed, notifies the Hiring Manager, and holds the Interviewee for decision (FR-APP-41).<br>E3. The system reminds the Interviewee 48h and 24h before the completion window ends (FR-APP-40). |
| Postconditions    | The stage's module is configured; scoring (0-100) is enabled for completions. |
| Business Rules    | Extra time or a longer window may be granted to an individual Interviewee with a recorded reason (FR-APP-42). |

### UC-APP-08: Build Custom Test Module with AI Assistance

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-08 |
| Use Case Name     | Build Custom Test Module with AI Assistance |
| Actor(s)          | Hiring Manager, System (AI assistance) |
| Description       | The Hiring Manager builds a custom test module in the Design Studio, optionally using AI to draft questions, and saves it for reuse. |
| Preconditions     | A pipeline design exists for the posting (3.4). |
| Trigger           | The Hiring Manager opens the Design Studio to create a custom test module. |
| Main Flow         | 1. The Hiring Manager adds multiple-choice, short-text, long-text, and code questions.<br>2. The system optionally drafts questions from the Hiring Manager's description of the skill and difficulty (AI assistance).<br>3. The Hiring Manager reviews and approves each AI-drafted question.<br>4. The Hiring Manager assigns points and an answer key to each closed question (auto-scored).<br>5. The Hiring Manager optionally saves the module to the company library for reuse. |
| Alternative Flows | 1. Open-ended answers are routed to the Hiring Manager or an assigned Interviewer for grading against assigned points. |
| Exception Flows   | E1. AI-drafted questions are not made available to Interviewees until reviewed and approved (FR-APP-46).<br>E2. AI generation fails: do not release partial output; notify the Hiring Manager to write manually. |
| Postconditions    | A custom module is configured and, optionally, saved to the company library. |
| Business Rules    | Closed questions are auto-scored against the answer key (FR-APP-47). Open-ended answers require manual grading (FR-APP-48). |

### UC-APP-09: Schedule and Conduct Live Interview

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-09 |
| Use Case Name     | Schedule and Conduct Live Interview |
| Actor(s)          | Hiring Manager, Interviewer, Interviewee, external video-meeting provider (EXT) |
| Description       | The platform schedules a live-interview stage and records who conducted it, when, and the notes and scorecard. |
| Preconditions     | The stage is configured as a live interview with at least one assigned Interviewer or Hiring Manager. The external provider is reachable (EXT). The Interviewee has reached the stage. |
| Trigger           | An assigned Interviewer publishes available time slots. |
| Main Flow         | 1. The Interviewer publishes available time slots.<br>2. The Interviewee selects a slot; the system creates the online meeting via the external provider and confirms the slot.<br>3. The system sends the meeting time and join link to both parties within 10 minutes.<br>4. The system records for the interview the assigned Interviewer(s), scheduled start/end, and outcome (held, Interviewee absent, or Interviewer absent).<br>5. The Interviewer writes timestamped notes and submits a scorecard (rubric scores 1-5 and a written recommendation) within 48 hours.<br>6. The system sets the module score to the mean of submitted scorecard scores, normalized to 0-100. |
| Alternative Flows | 1. The system reminds both parties 24 hours and 1 hour before the interview.<br>5a. The system reminds an Interviewer if the scorecard is missing. |
| Exception Flows   | E1. The Interviewee reschedules up to 2 times, no later than 24 hours before the start (FR-APP-57).<br>E2. Outcome is "Interviewee absent" or "Interviewer absent": the system holds the Interviewee for the Hiring Manager's decision and notifies them (FR-APP-63).<br>E3. Meeting creation fails: retry 3 times within 15 minutes; if still failing, keep the slot unconfirmed and notify both parties to pick another time. |
| Postconditions    | The interview is recorded (assignees, timing, outcome); notes and scorecard are stored; a module score is set. |
| Business Rules    | Interviewer notes are shown only to the Hiring Manager and assigned Interviewers, never to the Interviewee (FR-APP-60). Meeting times are shown in the viewer's time zone (FR-APP-56). |

### UC-APP-10: Publish and Migrate Pipeline Version

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-10 |
| Use Case Name     | Publish and Migrate Pipeline Version |
| Actor(s)          | Hiring Manager, Interviewee (notified) |
| Description       | The Hiring Manager publishes a new pipeline version; in-progress Interviewees are migrated onto it without losing results. |
| Preconditions     | At least one saved pipeline version exists. For migration, Interviewees are in progress on the currently published version. |
| Trigger           | The Hiring Manager publishes a pipeline version. |
| Main Flow         | 1. The system validates the version against checks FR-APP-32, FR-APP-33, FR-APP-34, and FR-APP-99, listing every failed check.<br>2. The system shows a migration preview (per stage, the number of Interviewees affected and where each group will be placed).<br>3. The Hiring Manager confirms publication; the system unpublishes the previous version and moves in-progress Interviewees to the new version, keeping completed results and never moving anyone backward.<br>4. The system records the publication (Hiring Manager, time, version numbers, placements) and notifies by email each Interviewee who must complete an additional stage. |
| Alternative Flows | 1. Editing the published version creates a new draft version while the published version stays in force. |
| Exception Flows   | E1. A removed stage holds Interviewees: the system blocks publication until a destination stage is chosen (FR-APP-73).<br>E2. An added stage applies to in-progress Interviewees: the system blocks until the Hiring Manager decides its scope (FR-APP-74).<br>E3. A stage's rule or threshold changed: waiting Interviewees are re-evaluated without reversing any already-granted pass (FR-APP-75). |
| Postconditions    | Exactly one version is published; in-progress Interviewees are migrated with results preserved; the publication is recorded. |
| Business Rules    | Migration never moves an Interviewee backward and never deletes results (FR-APP-72). A stage is not repeated when its module and configuration are unchanged. |

### UC-APP-11: Advance Interviewee Between Stages

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-11 |
| Use Case Name     | Advance Interviewee Between Stages |
| Actor(s)          | Hiring Manager, System (automatic) |
| Description       | Each stage decides, by one of three advancement rules, whether an Interviewee passes to the next stage or is rejected. |
| Preconditions     | A pipeline is published and Interviewees are in progress. |
| Trigger           | A module completes or the Hiring Manager records a decision. |
| Main Flow         | 1. The system applies the stage's advancement rule: manual decision, automatic score threshold, or automatic on completion.<br>2. On a pass, the system moves the Interviewee to the next stage(s) and records the decision (type, decider, time, module score). |
| Alternative Flows | 1. Manual rule: only the Hiring Manager records pass or reject.<br>2b. Threshold rule: the system passes when the score is at or above the set threshold.<br>2c. Completion rule: the system passes when the module completes, regardless of score.<br>2d. The Hiring Manager passes or rejects several selected Interviewees of one stage in a single action. |
| Exception Flows   | E1. Below-threshold handling: rejected automatically or held for decision, per the Hiring Manager's choice (default held).<br>E2. On any rejection: status set to "Rejected", Interviewee removed from all other active stages, rejection feedback started (FR-APP-87).<br>E3. The Hiring Manager overrides an automatic result with a written reason (FR-APP-84).<br>E4. Two simultaneous decisions on the same Interviewee/stage: record the first, reject the second, and show the recorded decision. |
| Postconditions    | The stage decision is recorded; the Interviewee is moved to next stage(s) or rejected. |
| Business Rules    | Only Hiring Managers record manual decisions; Interviewers contribute scorecards only (FR-APP-80). |

### UC-APP-12: Review Pipeline Dashboard and Insights

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-12 |
| Use Case Name     | Review Pipeline Dashboard and Insights |
| Actor(s)          | Hiring Manager |
| Description       | The Hiring Manager reviews the pipeline graphically with counts, rankings, and system-generated insights that support, but do not replace, decisions. |
| Preconditions     | A pipeline is published and at least one Interviewee has been admitted. |
| Trigger           | The Hiring Manager opens the pipeline dashboard. |
| Main Flow         | 1. The system shows a graphical view with, per stage, the count of Interviewees in it and the numbers passed and rejected.<br>2. The Hiring Manager selects a stage; the system lists Interviewees ranked by module score with filters (score range, Nice-to-have score, pending decision).<br>3. The system shows each Interviewee's scores and percentile, plus a system-generated summary of strengths and weaknesses.<br>4. The Hiring Manager records pass/reject decisions directly from the dashboard. |
| Alternative Flows | 1. The Hiring Manager compares up to 4 selected Interviewees side by side.<br>4b. The system suggests a cut-off score to reach the finalist target and shows how many would pass. |
| Exception Flows   | E1. Ranking/insight data excludes name, photo, gender, age, nationality, and contact details (FR-APP-92). |
| Postconditions    | The Hiring Manager has an updated view of pipeline progress and may record decisions. |
| Business Rules    | Insights support but never replace human decisions, except where an automatic rule is configured (FR-APP-78). |

### UC-APP-13: Enable Anonymous Hiring

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-13 |
| Use Case Name     | Enable Anonymous Hiring |
| Actor(s)          | Hiring Manager, Interviewer, Interviewee (informed) |
| Description       | The Hiring Manager hides Interviewee identity from hiring staff until a chosen reveal point in the pipeline. |
| Preconditions     | A pipeline is being designed (3.4). |
| Trigger           | The Hiring Manager enables anonymity on a pipeline. |
| Main Flow         | 1. The Hiring Manager enables anonymity and selects an identity-reveal point (a stage or the Finalist list).<br>2. Before the reveal point, the system replaces identity fields (name, photo, email, phone, personal links, date of birth, gender, nationality, address) with a system-assigned anonymous identifier in every view for Hiring Managers and Interviewers.<br>3. The system informs the Interviewee that identity is hidden in early stages, without naming the reveal point.<br>4. When an Interviewee reaches the reveal point, the system reveals identity permanently and records the time. |
| Alternative Flows | None. |
| Exception Flows   | E1. A path can reach a video-recorded response or live-interview stage before the reveal point: the pipeline is invalid (FR-APP-99).<br>E2. Access to identity before the reveal point: deny and log the attempt (FR-APP-98).<br>E3. Identity is revealed no later than the moment an offer is created (FR-APP-101). |
| Postconditions    | Anonymity is active up to the reveal point; reveal is permanent and recorded. |
| Business Rules    | Identity reveal is permanent (FR-APP-100). Anonymity is incompatible with video-recorded response and live-interview modules before the reveal point. |

### UC-APP-14: Provide Feedback and Notifications

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-14 |
| Use Case Name     | Provide Feedback and Notifications |
| Actor(s)          | Hiring Manager, Interviewee, System (AI) |
| Description       | The system releases module and rejection feedback, and advance notifications, to Interviewees by email and in the platform. |
| Preconditions     | The Interviewee has completed a module, advanced, or been rejected. A feedback mode is chosen per stage and per pipeline. |
| Trigger           | A module completes, an advancement occurs, or a rejection is recorded. |
| Main Flow         | 1. The system releases module feedback within 10 minutes of completion (AI-generated mode) or of approval/submission (other modes).<br>2. The system emails the Interviewee within 10 minutes of being passed to a next stage, naming the stage.<br>3. Within 24 hours of a rejection, the system produces rejection feedback (overall summary, per-module outcome with at least one strength or improvement area, and at least one recommended action) and releases it per the pipeline mode. |
| Alternative Flows | 1. Feedback mode is human-written or AI-generated with human approval before release. |
| Exception Flows   | E1. Feedback excludes test questions, answer keys, thresholds, and other Interviewees' data or identity (FR-APP-105).<br>E2. AI-generated feedback is labelled as AI-generated (FR-APP-106).<br>E3. AI generation fails: hold the item and notify the Hiring Manager to write manually. |
| Postconditions    | Feedback and notifications are delivered by email and in the platform. |
| Business Rules    | Rejection feedback is required within 24 hours of rejection (FR-APP-108). Feedback modes follow the same three options as FR-APP-103. |

### UC-APP-15: Track Application as Interviewee

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-15 |
| Use Case Name     | Track Application as Interviewee |
| Actor(s)          | Interviewee |
| Description       | An Interviewee follows their applications, current stage, feedback, and past applications, seeing progress only. |
| Preconditions     | The Interviewee is signed in and has submitted at least one application. |
| Trigger           | The Interviewee opens their application tracking view. |
| Main Flow         | 1. The system shows all active applications with posting title, company, status, and current stage name(s).<br>2. The system shows a graphical view of the application's pipeline (stage names, links, passed and current stages).<br>3. The system shows all released module and rejection feedback.<br>4. The system shows the currently required action (start a module, choose a slot, or join a meeting) with its due date or time.<br>5. The system shows a timeline of status changes with timestamps. |
| Alternative Flows | 1. Past (no longer active) applications are shown with their final outcome and dates. |
| Exception Flows   | E1. The system never shows stage configuration, thresholds, advancement rules, question counts, difficulty, assigned Interviewers, Interviewer notes, or other Interviewees' data (FR-APP-113). |
| Postconditions    | The Interviewee has an up-to-date view of progress and required actions. |
| Business Rules    | Interviewees see progress only (FR-APP-113). |

### UC-APP-16: Manage Finalists and Offers

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Use Case ID       | UC-APP-16 |
| Use Case Name     | Manage Finalists and Offers |
| Actor(s)          | Hiring Manager, Interviewee |
| Description       | Interviewees who complete the pipeline become Finalists; the Hiring Manager sends offers within capacity and rejects the rest. |
| Preconditions     | At least one Interviewee has passed every stage on all paths of the pipeline. |
| Trigger           | An Interviewee completes the pipeline, or the Hiring Manager opens the Finalist list. |
| Main Flow         | 1. The system sets an Interviewee's status to "Finalist" and adds them to the Finalist list when they pass every stage.<br>2. The system orders the list by mean module score by default and shows insights.<br>3. The Hiring Manager creates an offer (job title, compensation summary, start date, response deadline, optional message); the system delivers it by email and in the platform.<br>4. The Interviewee accepts or declines before the deadline; the system records the response, sets status "Offer Accepted" or "Offer Declined", and notifies the Hiring Manager within 10 minutes.<br>5. The Hiring Manager rejects remaining Finalists individually or in bulk, starting rejection feedback for each. |
| Alternative Flows | 1. The Hiring Manager orders the list by the score of any single stage.<br>4a. The Hiring Manager extends the response deadline of a pending offer or withdraws a pending offer (notifying the Interviewee). |
| Exception Flows   | E1. No response by the deadline: status "Offer Expired"; both parties notified (FR-APP-123).<br>E2. A new offer would exceed open positions: the system warns and requires explicit confirmation (FR-APP-126).<br>E3. A response to an expired, withdrawn, or already answered offer is rejected with the current status shown.<br>E4. More than one pending offer per Interviewee per posting is not allowed (FR-APP-128). |
| Postconditions    | Offers are accepted, declined, expired, or withdrawn; Finalists are rejected with feedback started. |
| Business Rules    | At most one pending offer per Interviewee per posting (FR-APP-128). Identity is revealed no later than offer creation (FR-APP-101). |


## 4.5 JOB Use Cases

## 4.6 ANA Use Cases
