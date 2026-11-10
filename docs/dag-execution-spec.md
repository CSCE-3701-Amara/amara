# DAG pipelines and code execution

Hiring teams design, save and publish acyclic assessment graphs; candidates write and run code through Judge0.
User authorization: implement missing behavior and choose implementation details autonomously.

- Reuse existing pipeline stages, transitions, coordinates, immutable snapshots and worker joins. No schema changes.
- GET/PUT `/v1/pipeline-versions/:versionId/graph`: `{stages,links}`; stage IDs are client keys on save, replaced with persisted IDs in the response. PUT atomically replaces unpublished drafts; published graphs are immutable.
- Validate unique IDs/names, bounded graphs, existing endpoints, no self/duplicate edges and no cycles. Existing linear drafts load as adjacent links. Saving linear stages clears explicit graph links.
- Designer uses the existing typography/palette with a scrollable graph canvas, SVG connections, node inspector, keyboard positioning and zoom. Save errors preserve work.
- Candidate attempt responses include owner-only savedResponses so source/language survive reload; the field is additive and contains no grading evidence.
- Code execution: authenticated candidate POST `/v1/assessment-attempts/:attemptId/run` with questionId, language, source and stdin. Check ownership, deadline and editable status; never expose hidden test cases. Bound concurrency and payloads. Provider errors return 503 without changing answers or scores.
- Judge0 remains a separately operated sandbox. API and worker use configured HTTPS Judge0; development can explicitly use a loopback HTTP service. No Docker socket in the API.
- Grading: worker calculates weighted per-question percentages. Existing minimum-count modules retain count rules. Legacy percentage modules use their module threshold, with equal-weight historical case counts supported; insufficient weighted evidence waits for trusted regrading.
- Cases: reject cyclic/malformed/oversized graphs; deny foreign company/candidate; reject edits after publish and runs after deadline; joins wait for all parents; failed providers never count as candidate failure; failed save retains edits.
- Verify domain validation, real PostgreSQL save/reload/publication, scoring boundaries, provider failure and browser design interactions.
- Rollout: deploy API/worker and frontend together; configure Judge0 on both execution roles. Existing linear endpoints remain available.
