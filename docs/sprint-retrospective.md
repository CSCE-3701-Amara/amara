# Sprint Retrospective



## What went well

- **The project established a useful foundation before expanding features.**
  The team documented requirements, use cases, and architecture, then added
  the Go service structure, Supabase migrations, Docker Compose, and secrets
  configuration. This gave feature work shared technical and product
  conventions to build on.
- **The product advanced across connected workflows.** Commits delivered
  account and authentication flows, job-posting and application features,
  assessment progress and feedback, candidate pipeline progression, and
  candidate and employer-facing screens. Backend, frontend, and API-contract
  work progressed together rather than treating the product as isolated
  components.
- **Quality and security checks became part of the delivery process.** The
  repository added CI workflows, static analysis, dependency vulnerability
  checks, security checks, and frontend interaction/design coverage. The team
  also responded to test failures and pull-request feedback with follow-up
  fixes.
- **The team adapted features to shared system changes.** Branch updates and
  integration commits show active coordination as schemas, routes, and shared
  conventions evolved.

## What could be improved

- **Integrate and validate changes earlier.** Several feature histories include
  repeated merge, compiler, migration, test, and review-fix commits. For
  example, assessment work needed follow-up changes for compiler errors and a
  duplicate table-creation migration, and pipeline progress needed a route
  registration fix. Finding these issues sooner would reduce rework near
  delivery.
- **Coordinate database changes as one shared history.** The duplicate
  `application_assessments` table migration indicates that parallel branches
  need an earlier check for overlapping schema ownership, migration names, and
  replay order.
- **Make feature completion criteria consistently verifiable.** The history
  shows implementation and documentation moving at different times, and the
  testing and Definition of Done convention pages are still drafts. Each
  feature should have clear acceptance criteria and an agreed test plan before
  implementation is considered complete.
- **Leave more time for integration and review.** Many substantial features
  and their corrective commits landed on November 3-7. A steadier flow of small
  pull requests would make conflicts and regressions easier to isolate.

## What we will do next sprint

1. **Agree on sprint scope and acceptance criteria before implementation.**
   Split work into small, demonstrable stories and identify dependencies,
   owners, and required tests during planning.
2. **Integrate continuously.** Keep feature branches current, open smaller pull
   requests early, and run the relevant checks before requesting review. Merge
   completed work throughout the sprint instead of deferring integration to
   the end.
3. **Add migration replay to the feature completion path.** For schema
   changes, review migration names and ownership together and verify that the
   complete migration history applies cleanly to a fresh database in CI.
4. **Use an end-to-end Definition of Done.** Verify the API contract, route
   registration, persistence, authorization where applicable, automated tests,
   and user-visible behavior for each workflow. Update the testing and
   Definition of Done conventions to reflect this checklist.
5. **Reserve capacity for stabilization.** Schedule time for review feedback,
   regression fixes, and documentation, and track planned versus completed
   work and late integration defects so the next retrospective can assess
   whether these changes helped.

