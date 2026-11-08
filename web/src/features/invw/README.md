# Interview workflow feature

`invw` contains hiring-workflow UI for pipelines, assessment modules, runtime
stages, grading, interviews, offers, and candidate progression. Employer
pipeline pages are mounted under `/employer/pipelines`; related interviewer
pages are mounted under `/interviewer`.

## API boundary

Pipeline drafts and runtime projections are accessed through the helpers in
`api/` and the runtime types in `types/`. The backend owns stage transitions,
threshold validation, candidate status changes, and authorization. The
frontend may preview a change and must send the observed version/threshold
when the API requires optimistic concurrency.

Do not treat unavailable screens as completed workflows. Keep explicit
`UnavailableNotice` states until the corresponding authenticated backend
contract exists, and never turn provider or grading failures into a zero score.