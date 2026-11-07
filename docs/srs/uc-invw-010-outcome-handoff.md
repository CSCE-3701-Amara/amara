# UC-INVW-010 assessment outcome handoff

Join `assessment_attempts.stage_execution_id` to `stage_executions.id`, then join `stage_executions.pipeline_execution_id` to `pipeline_executions.id`. The pipeline execution identifies the application and published version. The stage execution identifies the stage. The attempt identifies the candidate directly and has a unique `instance_id` referencing `assessment_instances.id`. Use the attempt as the source of truth for assessment outcomes. Its status is mirrored onto the instance for delivery. UC-010 does not write advancement decisions or change the application or pipeline execution to rejected.

| Outcome | Durable signal |
| --- | --- |
| Missed individual assessment | `assessment_attempts.status = 'missed'`, `started_at` and `submitted_at` are null, and `auto_submitted = false`. The worker creates a missed attempt and instance if the assignment first comes up for processing after its deadline. |
| Candidate completed submission | `submitted_at IS NOT NULL`. `auto_submitted = true` means the attempt was started and the deadline finalized the saved responses. Explicit and automatic submission freeze `saved_responses`. |
| Behavioral submission | `submitted_at IS NOT NULL` on a module with definition `type = 'behavioral'` or a recorded-video definition with `behavioral = true`. A wholly behavioral attempt first uses `submitted`; mixed attempts use `pending_grading`. Qualitative evaluation may later set `scored` with a null `final_score`. The completion signal remains `submitted_at`. |
| Grading pending | `status = 'pending_grading'` and `evaluated_at IS NULL`. Invalid or incomplete evidence leaves this state intact. |
| Evaluation complete | `status = 'scored'` and `evaluated_at IS NOT NULL`. Read `final_evidence` in the same row. `final_score` is the rounded mean of question scores when every question carries a score. For mixed attempts, use each module's question evidence. |
| Live interview completed | Join `interviews.stage_execution_id` to the stage execution. `status = 'completed'` and `outcome = 'held'` establish attendance and completion. Absence outcomes are distinct. Read the behavioral flag from `pipeline_stages.interview_configuration.behavioral`. Human scorecards and notes remain in `interview_feedback` and join on `interviewer_assignments (interview_id, interviewer_id)`. For scored held interviews the worker writes `interviews.score` when all assignments are submitted or waived and at least one valid scorecard was submitted. |

## Coding evidence

`selected_content` is an ordered array of objects with `module_id` and an ordered `questions` array. Question identity is the UUID produced by casting the MD5 of `attempt_id:module_ordinal:question_ordinal`, with both ordinals one-based. The ID identifies the selected question within this immutable attempt, including library selections. `configuration_snapshot.modules` has matching module order and stores each published definition at `configuration.definition`.

`final_evidence` is a JSON array with exactly one record per selected question. A coding record has `questionId`, `passedCases`, and `totalCases`. `totalCases` must match the selected question's snapshotted test-case count, and `passedCases` must be between zero and that count. Human records contain `questionId`, optional `notes`, and `score` for scored questions. Behavioral records omit `score`. The evaluator handoff writes evidence, evaluation timestamp, applicable score, and status in a single transaction. Identical retries return success without replacing evidence.

UC-006 can resolve `questionId` against the snapshot, compare `passedCases` with the question's `minimum_cases_passed`, then count passed questions for each module and compare with the corresponding module definition's `minimum_tests_passed`. Missing thresholds on older snapshots remain missing. UC-010 does not infer or apply an advancement rule.

## Entry points

- The worker spawns attempts, records deadline outcomes, and reconciles scores from submitted interview scorecards. When a stage contains multiple modules, the attempt uses the earliest configured completion window and the shortest configured active timer.
- Candidate instance listing, start, attempt retrieval, autosave, and submission routes use the candidate's authenticated session. Start fixes the effective deadline. Candidate responses and timing survive retries, and the candidate question projection excludes private configuration and grading material.
- The owning company can review completed submissions and write final human evidence with GET and POST on `/v1/assessment-attempts/{attemptId}/evaluation`. POST takes `{"evidence":[{"questionId":"<uuid>","score":80,"notes":"Optional notes"}]}` and returns 204. The company route rejects coding evidence.
- The interview outcome route records held or absence for a confirmed interview after its scheduled end. It accepts an assigned interviewer or the owning company.
- Coding finalization is implemented as a trusted repository entry point. This branch has no external evaluator client or callback wired to it, so coding remains pending in deployment until that integration submits evidence. Interview scheduling, scorecard submission and waiver routes, and video upload are also unwired on this branch. The interview outcome and score paths consume interview and scorecard records when present. The assessment autosave route accepts a video reference but does not upload or validate ownership of the video asset.
