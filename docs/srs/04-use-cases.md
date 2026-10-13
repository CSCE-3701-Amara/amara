## 4. Use cases

Owner: All | Status: Draft | Last updated: 2026-10-13 | Jira: AMARA-32

---

## 4.1 Actors

- Candidate
- System Administrator
- Interviewer
- Hiring Manager
- System

## 4.2 AUTH Use Cases

## 4.3 CAND Use Cases

## 4.4 INVW Use Cases
### UC-INVW-001: Design Interview Pipeline

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager designs the interview pipeline of a posting as a directed acyclic graph (DAG) on a visual canvas. |
| Preconditions     | The Hiring Manager is signed in and permitted to manage the posting (AUTH). The posting exists. |
| Trigger           | The Hiring Manager opens the pipeline designer. |
| Main Flow         | 1. The Hiring Manager adds uniquely named stages from the module palette and connects them with directed links.<br>2. The System validates that the graph is acyclic and every stage is reachable from an entry stage.<br>3. The Hiring Manager configures a module, scheduled start/end times, and exactly one advancement rule for each assessment stage (FR-INVW-051, FR-INVW-094).<br>4. The Hiring Manager saves the design as a draft pipeline version; the System records its validation result. |
| Alternative Flows | 1. The Hiring Manager moves or deletes stages and links. |
| Exception Flows   | E1. A cycle is detected: the System identifies the stages involved (FR-INVW-005).<br>E2. A stage is unreachable from an entry stage: the System identifies it (FR-INVW-006).<br>E3. A stage has a duplicate name, missing module, or no advancement rule: the System identifies the stage and missing item (FR-INVW-007).<br>E4. An assessment stage's end time is not later than its start time: the System identifies the stage and invalid date range (FR-INVW-094). |
| Postconditions    | The draft pipeline version and its validation result are saved; a draft with failed checks cannot be published until corrected (FR-INVW-039). |
| Business Rules    | A draft is not available to Candidates until published through UC-INVW-005. |

### UC-INVW-002: Configure Assessment Module

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager configures an individual assessment module, its content, scoring method, and timing. Live interviews are configured through UC-INVW-004. |
| Preconditions     | A pipeline design exists for the posting (UC-INVW-001). For library modules, the problem library contains matching problems. |
| Trigger           | The Hiring Manager selects and configures a module for a stage. |
| Main Flow         | 1. The Hiring Manager selects an individual assessment type from FR-INVW-008.<br>2. The Hiring Manager selects or creates its content and configures the supported response mode and scoring method required for that type (FR-INVW-082). For coding problems, they also configure count, difficulty, point values, and selection mode (FR-INVW-009).<br>3. The Hiring Manager sets the active time limit and Candidate completion window (FR-INVW-011).<br>4. The System validates and saves the module configuration. |
| Alternative Flows | 2a. For random coding-problem selection, the Hiring Manager sets the problem count and difficulty mix; the System selects matching problems when it creates the test instance (FR-INVW-009, FR-INVW-010).<br>2b. For question-based English, programming-language, behavioral, or aptitude modules, the Hiring Manager sets the question count and selects or authors each question, response format, and scoring method (FR-INVW-091).<br>2c. For a video-recorded response module, the Hiring Manager configures each prompt and assigns a grader (FR-INVW-093). |
| Exception Flows   | E1. Required content, response mode, scoring method, active time limit, or Candidate completion window is missing: the System identifies the missing item and marks the module incomplete for publication (FR-INVW-011, FR-INVW-082).<br>E2. The problem library has no problems matching the configured count and difficulty mix: the System reports the mismatch and marks the module incomplete for publication. |
| Postconditions    | A complete individual assessment module is saved and available for pipeline publication (FR-INVW-082). |
| Business Rules    | Live interviews use the separate scheduling and scorecard flow in UC-INVW-004; LLM-paired programming is a future Could feature (UC-INVW-014); current test-instance execution is covered by UC-INVW-010. |

### UC-INVW-003: Build Custom Test Module with AI Assistance

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager builds a custom test module in the Design Studio, optionally using AI to draft questions, and saves it for reuse. |
| Preconditions     | A pipeline design exists for the posting (UC-INVW-001). |
| Trigger           | The Hiring Manager opens the Design Studio to create a custom test module. |
| Main Flow         | 1. The Hiring Manager adds one or more multiple-choice, short-text, long-text, or code questions.<br>2. The System may draft questions from the Hiring Manager's skill and difficulty description.<br>3. The Hiring Manager reviews and approves every AI-drafted question before saving it (FR-INVW-018, FR-INVW-019).<br>4. The Hiring Manager assigns maximum points and an answer key to closed questions; for open-ended questions, the Hiring Manager assigns themselves or a stage-assigned Interviewer to grade (FR-INVW-020, FR-INVW-021).<br>5. The Hiring Manager optionally saves the complete module to the company library (FR-INVW-022). |
| Alternative Flows | None. |
| Exception Flows   | E1. AI generation fails: the System discards partial output and notifies the Hiring Manager to write the questions manually. |
| Postconditions    | A custom module is configured and, optionally, saved to the company library. |
| Business Rules    | None. |

### UC-INVW-004: Schedule and Conduct Live Interview

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Interviewer, Candidate, System |
| Description       | The System schedules a live-interview stage through the external meeting provider and records its outcome, notes, and scorecards. |
| Preconditions     | The stage is configured as a live interview with scheduled start/end times and at least one assigned Interviewer or Hiring Manager. The external provider is reachable (EXT). The Candidate has reached the stage during its scheduled interval (FR-INVW-094). If anonymity is enabled, the Candidate's identity-reveal point is strictly before this stage (FR-INVW-072). |
| Trigger           | An assigned Interviewer publishes available time slots within the stage's scheduled interval. |
| Main Flow         | 1. The assigned Interviewer publishes available time slots within the stage's scheduled start/end times (FR-INVW-094).<br>2. The Candidate selects a slot that starts and ends within that interval; the System creates and confirms the online meeting through the external provider (FR-INVW-026).<br>3. The System sends the meeting time and join link to the Candidate and assigned Interviewers within 10 minutes (FR-INVW-027).<br>4. The System reminds the Candidate and each assigned Interviewer 24 hours and 1 hour before the interview (FR-INVW-028).<br>5. After the meeting, the System records the assigned Interviewers, scheduled start/end, and outcome (held, Candidate absent, or Interviewer absent) (FR-INVW-031).<br>6. Each assigned Interviewer saves timestamped notes and submits ratings for every rubric criterion and a written recommendation within 48 hours; the System tracks missing scorecards and reminds the Interviewer (FR-INVW-032, FR-INVW-034).<br>7. The System calculates the interview score after all assigned scorecards are submitted or a missing card is waived with a recorded reason; if none are submitted, it leaves the score unset and holds the Candidate for a Hiring Manager decision (FR-INVW-034, FR-INVW-035). |
| Alternative Flows | 2a. The Candidate reschedules up to 2 times, no later than 24 hours before the interview (FR-INVW-030).<br>5a. The Candidate or an assigned Interviewer is absent: the System records the outcome, notifies the Hiring Manager, and holds the Candidate for a decision (FR-INVW-036).<br>6a. A scorecard remains missing at 48 hours: the System notifies the Hiring Manager, who may waive it with a recorded reason (FR-INVW-034).<br>6b. No scorecard is submitted after the Hiring Manager's waiver: the System leaves the module score unset and holds the Candidate for a decision (FR-INVW-035). |
| Exception Flows   | E1. Meeting creation fails: the System retries 3 times within 15 minutes; if it still fails, the slot remains unconfirmed and the System notifies the Candidate and Interviewer to choose another time. |
| Postconditions    | The interview outcome and attendance are recorded; notes and submitted scorecards are stored; a score is calculated when scorecards are complete or waived. |
| Business Rules    | Interviewer notes are visible only to the Hiring Manager and Interviewers assigned to the stage (FR-INVW-033). Meeting times use the viewer's time zone (FR-INVW-029). |

### UC-INVW-005: Publish and Migrate Pipeline Version

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Candidate, System |
| Description       | The Hiring Manager publishes a new pipeline version; the System migrates in-progress Candidates without losing results. |
| Preconditions     | At least one saved pipeline version exists. For migration, Candidates are in progress on the currently published version. |
| Trigger           | The Hiring Manager publishes a pipeline version. |
| Main Flow         | 1. The System validates the version against FR-INVW-005, FR-INVW-006, FR-INVW-007, FR-INVW-072, FR-INVW-082, FR-INVW-094, and FR-INVW-095 (FR-INVW-039). If any check fails, it lists every failure and stops publication.<br>2. The System shows a migration preview with the number of affected Candidates per stage and their proposed placements.<br>3. The Hiring Manager confirms publication; the System unpublishes the previous version, migrates in-progress Candidates without moving anyone backward, preserves completed results, and re-evaluates waiting Candidates under changed rules without reversing a pass (FR-INVW-043, FR-INVW-045, FR-INVW-048).<br>4. The System records the publication and placements and emails each Candidate who must complete an additional stage (FR-INVW-049, FR-INVW-050). |
| Alternative Flows | 1. Editing the published version creates a new draft version while the published version stays in force. |
| Exception Flows   | E1. A removed stage holds Candidates: publication is blocked until the Hiring Manager selects a destination (FR-INVW-046).<br>E2. An added stage affects in-progress Candidates: publication is blocked until the Hiring Manager chooses its scope (FR-INVW-047).<br>E3. An assessment stage has an invalid date range or a first exam that starts on or before the posting deadline: publication is blocked and the affected stage/date check is reported (FR-INVW-094, FR-INVW-095). |
| Postconditions    | Exactly one version is published; in-progress Candidates are migrated with completed results preserved; publication and placements are recorded. |
| Business Rules    | A stage with unchanged module and configuration is not repeated during migration (FR-INVW-045). |

### UC-INVW-006: Advance Candidate Between Stages

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Candidate, System |
| Description       | The System applies the stage's configured rule to decide whether a Candidate passes, is rejected, or remains held for a Hiring Manager decision. |
| Preconditions     | A pipeline is published and the Candidate is active in at least one stage. |
| Trigger           | An assessment is fully scored, a live-interview outcome is recorded, or the Hiring Manager records a manual decision. |
| Main Flow         | 1. The System applies the stage's configured advancement rule (FR-INVW-051).<br>2. The System records the decision type, decider, time, and score (FR-INVW-058).<br>3. On pass, the System makes linked stages available only when all incoming stages have passed; on rejection, it marks the application Rejected, removes the Candidate from all active stages, and starts rejection feedback (FR-INVW-002, FR-INVW-003, FR-INVW-060). |
| Alternative Flows | 1a. Manual rule: the Hiring Manager records pass or reject (FR-INVW-052, FR-INVW-053).<br>1b. Threshold rule: the System passes at or above the stage's current threshold; below-threshold handling is automatic rejection or hold, with hold as the default. Threshold changes from the dashboard are handled as defined in UC-INVW-007 (FR-INVW-054, FR-INVW-055, FR-INVW-096).<br>1c. Completion rule: the System passes when the module is complete, regardless of score (FR-INVW-056).<br>2a. The Hiring Manager overrides an automatic result with a written reason (FR-INVW-057).<br>2b. The Hiring Manager passes or rejects several selected Candidates at the same stage in one action (FR-INVW-059). |
| Exception Flows   | E1. Two simultaneous decisions concern the same Candidate and stage: the System records the first and rejects the second, showing the recorded decision. |
| Postconditions    | The decision is recorded and the Candidate is advanced, rejected, or held for a Hiring Manager decision. |
| Business Rules    | Only the Hiring Manager records manual decisions; Interviewers submit scorecards but do not make advancement decisions (FR-INVW-053). Stage advancement is independently configured as manual, threshold-based, or on completion (FR-INVW-051). |

### UC-INVW-007: Review Pipeline Dashboard and Insights

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, System |
| Description       | The Hiring Manager reviews Candidate progress, rankings, and system-generated insights to support decisions. |
| Preconditions     | A pipeline is published and at least one Candidate is active in it. |
| Trigger           | The Hiring Manager opens the pipeline dashboard. |
| Main Flow         | 1. The System shows a graphical view with each stage's active, passed, and rejected Candidate counts (FR-INVW-061).<br>2. The Hiring Manager selects a stage; the System lists its Candidates ranked by module score and provides score-range, Nice-to-have, and pending-decision filters (FR-INVW-062).<br>3. The System shows each Candidate's completed-stage scores and percentile (FR-INVW-063), plus a labelled summary of strengths and weaknesses (FR-INVW-064).<br>4. The Hiring Manager may record pass/reject decisions from the dashboard; decision processing follows UC-INVW-006 (FR-INVW-066). |
| Alternative Flows | 2a. The System suggests a cut-off score expected to reach the finalist target and shows the resulting pass count (FR-INVW-068).<br>3a. The Hiring Manager compares up to 4 selected Candidates side by side (FR-INVW-069).<br>3b. After reviewing the score distribution and insights, the Hiring Manager applies the suggested or another revised threshold; the System previews and recalculates outcomes for Candidates still waiting at that stage (FR-INVW-096). |
| Exception Flows   | None. |
| Postconditions    | The Hiring Manager has an updated pipeline view. |
| Business Rules    | Ranking and insight inputs exclude identity and demographic attributes (FR-INVW-065). Insights support but do not replace human decisions except where an automatic advancement rule is configured (FR-INVW-051). |

### UC-INVW-008: Enable Anonymous Hiring

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Interviewer, Candidate , System |
| Description       | The Hiring Manager hides Candidate identity from Hiring Managers and Interviewers until a chosen pipeline reveal point. |
| Preconditions     | A pipeline is being designed (UC-INVW-001). |
| Trigger           | The Hiring Manager enables anonymity on a pipeline. |
| Main Flow         | 1. The Hiring Manager enables anonymity and selects an identity-reveal point (a pipeline stage or the Finalist list) (FR-INVW-070).<br>2. Before that point, the System substitutes the anonymous identifier for identity fields in every Hiring Manager and Interviewer view (FR-INVW-071).<br>3. The System tells the Candidate that identity is hidden from the employer during early stages without naming the reveal point (FR-INVW-074).<br>4. When the Candidate reaches the reveal point, the System reveals identity permanently and records the reveal time (FR-INVW-073). |
| Alternative Flows | None. |
| Exception Flows   | E1. A path can reach a video-recorded response or live-interview stage at or before the reveal point: the System marks the pipeline invalid (FR-INVW-072).<br>E2. A Hiring Manager or Interviewer requests identity before the reveal point: the System denies access and logs the attempt (FR-INVW-071). |
| Postconditions    | Anonymity is active up to the reveal point; reveal is permanent and recorded. |
| Business Rules    | Identity reveal is permanent (FR-INVW-073). In an anonymous pipeline, every live-interview and video-recorded response stage must follow the reveal point (FR-INVW-072). |

### UC-INVW-009: Provide Feedback and Notifications

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Hiring Manager, Candidate , System |
| Description       | The System releases assessment feedback, rejection feedback, and advancement notifications to Candidates by email and in the platform. |
| Preconditions     | The Candidate has completed a fully scored module, advanced, or been rejected; the applicable feedback mode is configured. |
| Trigger           | A module completes, an advancement occurs, or a rejection is recorded. |
| Main Flow         | 1. The System releases AI-generated module feedback within 10 minutes of completion, or human-written/AI-approved feedback within 10 minutes of submission or approval (FR-INVW-076).<br>2. The System emails the Candidate within 10 minutes of a pass to a next stage, naming that stage (FR-INVW-079).<br>3. Within 24 hours of rejection, the System produces and releases feedback with an overall summary, each completed module's outcome and at least one strength or improvement area, and at least one recommended action (FR-INVW-080, FR-INVW-081). |
| Alternative Flows | 1. Feedback mode is human-written or AI-generated with human approval before release. |
| Exception Flows   | E1. AI generation fails: the System holds the feedback item and notifies the Hiring Manager to write it manually. |
| Postconditions    | Feedback and notifications are delivered by email and in the platform. |
| Business Rules    | Feedback excludes questions, answer keys, thresholds, and other Candidates' data (FR-INVW-077); AI-generated feedback is labelled (FR-INVW-078). The Hiring Manager configures the feedback modes in FR-INVW-075 and FR-INVW-081. |

### UC-INVW-010: Spawn Test Instance

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Candidate, Hiring Manager, System |
| Description       | The System creates and times an individual assessment instance; the Candidate submits responses and the System scores it after any required manual grading. |
| Preconditions     | The Candidate has reached an individual assessment stage in the published pipeline and its scheduled interval has begun; its module content and scoring method are valid (UC-INVW-002, UC-INVW-003, FR-INVW-082, FR-INVW-094). Live interviews and multi-person tests use UC-INVW-004 and UC-INVW-012, respectively. |
| Trigger           | The scheduled start time has arrived and the Candidate is eligible for the assessment stage (FR-INVW-094). |
| Main Flow         | 1. The System creates the Candidate's instance from the published module configuration and preserves its selected content for subsequent opens (FR-INVW-083); random-selection stages give every Candidate the same problem count and difficulty mix (FR-INVW-010).<br>2. The Candidate's completion window starts when the stage becomes available; the active timer starts when the Candidate first opens the instance, displays remaining time, and continues across navigation or disconnection (FR-INVW-011, FR-INVW-084). Responses are saved and restored according to NFR-REL-008.<br>3. If the instance remains unfinished, the System sends reminders 48, 24, 12, 6, and 1 hours before the applicable deadline—the earlier of the completion-window deadline and assessment end time, as extended for the Candidate under FR-INVW-015 when applicable (FR-INVW-013).<br>4. The Candidate submits responses once; the System locks them and also locks the latest saved responses when auto-submitting at the active time limit or applicable deadline (FR-INVW-012, FR-INVW-085, FR-INVW-090).<br>5. The System automatically scores closed questions from answer keys (FR-INVW-020), sends coding solutions and test cases to the configured external code-evaluation provider (FR-INVW-087; see `docs/srs/06-external-interfaces.md`), and routes open-ended written answers and video responses to assigned human graders (FR-INVW-021, FR-INVW-093). AI-assisted question drafting does not score Candidate responses (FR-INVW-018, FR-INVW-019).<br>6. The System keeps the instance pending until all required human grading is complete, then records the final 0–100 score (FR-INVW-016, FR-INVW-086, FR-INVW-089).<br>7. Once fully scored, stage advancement follows UC-INVW-006 and module feedback follows UC-INVW-009. |
| Alternative Flows | 2a. The Candidate leaves and resumes the same instance; saved answers are restored and the active timer continues (FR-INVW-083, FR-INVW-084, NFR-REL-008).<br>2b. The Hiring Manager grants additional active time or a later individual completion deadline, including beyond the scheduled stage end for this started instance, and records the change and reason (FR-INVW-015).<br>4a. In a question-based module, the System presents each configured question using its response format and scoring method (FR-INVW-091).<br>4b. In a video-recorded response module, the Candidate records and submits a video response for each prompt; the System attaches each recording to the instance for human grading (FR-INVW-093). |
| Exception Flows   | E1. The earlier of the completion-window deadline and scheduled assessment end time arrives before an attempt has started or been submitted: the System marks the instance missed, notifies the Hiring Manager, and holds the Candidate for a decision (FR-INVW-014).<br>E2. The external code-evaluation provider is unavailable or returns an invalid result: the System preserves the code submission, leaves the instance pending, and does not assign a zero score (FR-INVW-087). |
| Postconditions    | The instance is fully scored and ready for advancement, remains pending manual grading, or is marked missed; advancement and feedback follow UC-INVW-006 and UC-INVW-009. |
| Business Rules    | Scoring and advancement are distinct: answer-key questions are auto-scored (FR-INVW-020), coding tests use external provider results (FR-INVW-087), and open-ended/video/live-interview responses use human grading or scorecards (FR-INVW-021, FR-INVW-034, FR-INVW-093). No assessment score or pass/reject decision is assigned by AI; AI question drafting, insights, and feedback do not assign scores (FR-INVW-018, FR-INVW-064, FR-INVW-076). The System withholds answer keys and rubrics from the active Candidate view (FR-INVW-088). A test instance is individual to one Candidate and stage assignment (FR-INVW-083). |

### UC-INVW-011: Detect Cheating on Test Instances

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Priority          | Should (long-term goal, not in the MVP) |
| Actor(s)          | Candidate, Hiring Manager, System |
| Description       | The system monitors spawned test instances for signs of cheating and flags suspicious instances with evidence for the Hiring Manager to review; a flag never decides an outcome by itself. |
| Preconditions     | Test instances exist for the stage (UC-INVW-010). |
| Trigger           | A test instance is started or completed. |
| Main Flow         | 1. During an assessment, the System records configured session events.<br>2. After completion, the System compares answers and timing across Candidates in the same stage.<br>3. The System flags suspicious instances and records the evidence that raised each flag.<br>4. The Hiring Manager reviews flags and records any advancement decision through UC-INVW-006. |
| Alternative Flows | 1. The Hiring Manager dismisses a flag with a written reason; the dismissal is recorded with the flag. |
| Exception Flows   | E1. The monitoring data is missing or inconclusive: no flag is raised.<br>E2. The session ends abnormally: the events recorded up to that point are kept, and missing data alone does not raise a flag. |
| Postconditions    | Suspicious instances carry their flags and evidence; only a Hiring Manager decision changes a Candidate's status (UC-INVW-006). |
| Business Rules    | Flags support but never replace human decisions; only the Hiring Manager records manual decisions (FR-INVW-053). Flags and evidence are shown only to the Hiring Manager. Functional requirements for this long-term goal are to be added when it is scheduled. |

### UC-INVW-012: Run Multi-Person Tests

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Priority          | Could (long-term goal, not in the MVP) |
| Actor(s)          | Hiring Manager, Candidates, System |
| Description       | A stage's test module is taken by several Candidates together (a pair or a team) in one shared session, with an individual score recorded for each member. |
| Preconditions     | A pipeline design exists for the posting (UC-INVW-001). The stage is configured with a multi-person module type and a team size. |
| Trigger           | The first Candidate reaches a stage configured with a multi-person module. |
| Main Flow         | 1. The System groups Candidates who reach the stage into pairs or teams of the configured size.<br>2. The System creates one shared test instance for each group from the stage's configured module.<br>3. Group members complete the test together within the configured completion window and active time limit.<br>4. The System records a 0–100 score for each member (FR-INVW-016).<br>5. Advancement of each member follows UC-INVW-006. |
| Alternative Flows | 1. Pair type: groups of two.<br>2. Team type: groups of the configured size, optionally with a team-level result in addition to the individual scores. |
| Exception Flows   | E1. A group member never joins or drops out: the session continues for the remaining members, and the absent member is held for the Hiring Manager's decision.<br>E2. The time limit expires: the System auto-submits the group's answers (FR-INVW-012).<br>E3. Not enough Candidates remain to complete a group: the System holds the ungrouped Candidates for the Hiring Manager's decision. |
| Postconditions    | Each group member has a completed module score or a held status; advancement follows UC-INVW-006. |
| Business Rules    | Multi-person modules are incompatible with anonymity before the reveal point. They pair Candidates with one another and are distinct from the future LLM-paired programming module (UC-INVW-014). Functional requirements for this long-term goal are to be added when it is scheduled. |

### UC-INVW-014: Use LLM-Paired Programming Module

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Priority          | Could (future feature; not in the current release) |
| Actor(s)          | Candidate, Hiring Manager, System |
| Description       | The Candidate works on a configured coding task with an LLM partner; the LLM assists during the session but does not assess or score the Candidate. |
| Preconditions     | The pipeline includes a configured LLM-paired programming module (FR-INVW-092), and the Candidate has reached the stage. |
| Trigger           | The Candidate opens the LLM-paired programming test instance. |
| Main Flow         | 1. The System presents the configured coding task and starts the active time limit.<br>2. The Candidate exchanges text with the LLM and edits a code solution during the session.<br>3. The System saves the conversation and final code with the test instance (FR-INVW-092).<br>4. The System sends the final code to the external code-evaluation provider specified in `docs/srs/06-external-interfaces.md`; scoring uses its valid test-case results (FR-INVW-087). |
| Alternative Flows | None. |
| Exception Flows   | E1. The LLM service is unavailable: the System preserves the Candidate's existing code and reports that the LLM partner is unavailable.<br>E2. The external code-evaluation provider is unavailable or returns an invalid result: the System preserves the final code and leaves the instance pending without assigning a score (FR-INVW-087). |
| Postconditions    | The conversation and final code are saved; the code is scored by the external provider, not by the LLM. |
| Business Rules    | The LLM assists the Candidate and does not make a pass/reject decision or produce the assessment score. |

### UC-INVW-015: View Candidate Pipeline Progress

| Field             | Description                                  |
| ----------------- | -------------------------------------------- |
| Actor(s)          | Candidate, System |
| Description       | The Candidate views the complete published pipeline and their own progress, including their current stage or stages. |
| Preconditions     | The Candidate's application is assigned to a published INVW pipeline after initial screening by the separate screening epic. |
| Trigger           | The Candidate opens the pipeline progress view for their application. |
| Main Flow         | 1. The System loads the published pipeline version assigned to the Candidate's application (FR-INVW-041).<br>2. The System displays the full stage graph, including every stage name and directed link in that version.<br>3. The System highlights the Candidate's current stage or stages and marks the Candidate's other stages as completed, active, or upcoming/locked (FR-INVW-097). |
| Alternative Flows | 1a. After a pipeline migration, the System shows the version the Candidate is now running and preserves completed-stage progress (FR-INVW-041, FR-INVW-045). |
| Exception Flows   | None. |
| Postconditions    | The Candidate can identify their current position and the remaining stages in the published pipeline. |
| Business Rules    | The view contains only the Candidate's own progress; it does not show other Candidates' data, Hiring Manager notes, or advancement thresholds (FR-INVW-097). |

## 4.5 JOB Use Cases

## 4.6 ANA Use Cases
