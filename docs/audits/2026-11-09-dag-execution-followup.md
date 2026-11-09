# R01/R02 follow-up: DAG pipelines and candidate execution

R01 and R02 are resolved in the implementation. This report supplements the original [codebase audit](2026-11-09-codebase-audit.md); its other findings remain applicable. No database migration, deployment or commit was performed. Existing audit changes are retained.

## Confirmed issues and implemented changes

| Issue | Severity | Root cause and affected files | Change and verification |
| --- | --- | --- | --- |
| R01: legacy coding attempts never advance | High | `internal/invw/adapters/database/advancement.go` understood only minimum-count definitions. `code_grading.go` discarded weighted percentage evidence. | Legacy coding and programming-language modules now retain integer weighted per-question scores and compare their rounded mean with the module threshold. Equal-weight historical counts remain usable. Unequal-weight historical evidence is regraded by the trusted worker instead of inventing a score. Real PostgreSQL regression verifies a 9:1 weighted case produces 90%, advances successfully and recovers historical evidence exactly once. Existing count-based rules remain intact. |
| R02: DAG authoring missing | High functional gap | `web/src/features/invw/components/pipeline-studio.tsx` rendered only the ordered editor, while the database and worker already supported transitions. | Added `dag-editor.tsx` and its stylesheet; authenticated graph handlers, domain validation and transactional persistence in `internal/invw/{domain,adapters/database,adapters/http}/graph.go`. Supports stage configuration, branches/joins, connection removal, cycle rejection, drag and keyboard positioning, zoom, save/reload and publication. Tests cover malformed graphs, ownership, published immutability and actual diamond execution through assessment admission and decisions. Original responsive DAG browser tests are retained. |
| Candidate could write code but could not execute it | High functional gap | Assessment handler had no execution route and UI had no provider-backed run action. | Added `code_run.go`, `code-runner.tsx`, API runner configuration and execution-language discovery. Candidates run their own source/stdin; attempt ownership, status, stage status and deadline are checked. Payloads and concurrency are bounded. Provider failures return 503 and never become grading failures. HTTP regressions cover role/ownership denial, payload and capacity limits, expiry, missing provider, unknown question, success and provider outage; browser regressions cover both candidate layouts. |
| Saved code was absent after reopening an attempt | High | Attempt projection omitted saved responses; frontend held successful saves only in component state. A default interpreter could replace a saved language. | `assessment_outcomes.go` adds owner-only `savedResponses`; `assessment.tsx` restores source and language and preserves language during later saves/runs. HTTP and browser regressions verify restoration. This is an additive response field documented in OpenAPI. |
| Legacy stages overlapped on the graph canvas | Medium usability | Linear drafts stored every coordinate as zero. | `adapters/database/graph.go` projects readable positions for unpositioned legacy drafts. PostgreSQL regression verifies separate positions and adjacent links before graph conversion. |
| A newly created pipeline did not open its graph | Medium | The graph and ordered editor fetched pipeline state independently. | The DAG entry point now creates and reloads its own pipeline. Browser regression verifies graph controls become available without a page reload. |
| Explicit coding stage threshold could be ignored | Medium | Coding advancement checked module minima before the separate scalar stage-threshold branch. | `advancement.go` now requires both module criteria and an explicitly configured stage threshold. Boundary regressions cover exact and below-threshold results. Count-only modules have no scalar score; use completion/manual advancement for these modules. |
| Container could not read generated source on SELinux hosts | Medium runtime | A non-root container reading a private host bind mount received permission denied under enforcing SELinux. | The development coordinator mounts only its generated source read-only with private SELinux relabeling. The parent directory remains mode 0700; no project or secret directories are relabeled. Actual Python/Node container runs reproduce successful execution after the fix. |

## Execution environment

`cmd/sandbox` is a separate loopback-only development coordinator implementing the Judge0 submission/polling protocol. It runs actual pre-pulled Docker containers for Python, JavaScript and native Node TypeScript. Containers have no network, a read-only root, a non-root UID, dropped capabilities, no privilege escalation, memory/CPU/process/output limits and an in-container timeout. Candidate containers never receive the Docker socket. The API and worker share the existing Judge0 adapter.

[Setup and configuration](../../deployments/judge0/README.md) covers both this working development runner and the pinned official Judge0 deployment. `make sandbox`, `make worker` and `make sandbox-judge0` provide entry points. Keep the normal database, SMTP and MFA configuration when starting Amara.

The official Judge0 image download stalled in this environment. Its PostgreSQL/Redis containers and download process were stopped without deleting data. Full Judge0 startup, compiler execution and host cgroup compatibility are **not verified**. The lightweight runner was tested against real Docker; it is development infrastructure, not a production VM boundary. Full TypeScript compilation, Go and C++ require the separately configured Judge0 service.

## Compatibility

- Reuses existing pipeline, transition, coordinate and snapshot tables; no schema changes.
- Graph GET/PUT and candidate run/language endpoints are additive and documented in `docs/api/openapi.yaml`. Existing routes remain available.
- Graph saves atomically replace an unpublished draft and return new persisted stage IDs; clients must use returned IDs. Published versions remain immutable.
- Existing linear drafts project adjacent connections. A multi-node graph without any links is rejected because it is indistinguishable from legacy linear storage. Independent roots are supported in graphs containing explicit edges.
- The ordered editor remains available. Saving there replaces explicit connections with linear order; use the graph editor for branching pipelines. Reload the page after switching editors to refresh their independent draft views.
- Legacy percentage scores are rounded per question and then averaged equally across questions within each module. Test-case weights apply within each question. Historical terminal hiring decisions are not reversed by evidence repair.
- Candidate Run does not save or submit answers; final grading still uses hidden tests in the worker. Configure the same provider and language mapping on API and worker.

## Verification

| Check | Result |
| --- | --- |
| Complete production Playwright suite | **396 passed, 0 failed, 0 skipped, 0 flaky** |
| Go suite with race detector and real PostgreSQL/Docker enabled | PASS; fresh complete run plus final focused graph/HTTP regressions |
| Actual Docker executions through Judge0 client | PASS: Python stdin, JavaScript, TypeScript, runtime diagnostics, infinite-loop termination and blocked networking |
| Go build, vet, staticcheck | PASS |
| Gosec using repository configuration | PASS; zero findings |
| Frontend production build, ESLint, TypeScript | PASS |
| OpenAPI YAML parsing, Judge0 Compose validation, shell syntax, diff whitespace | PASS |

 Browser tests use fixture-backed API responses; PostgreSQL workflow tests and real Docker client tests independently verify persistence and execution boundaries. No production infrastructure or candidate data was used.

An earlier overlapping browser run passed 393 checks but reported three trace-file `ENOENT` errors because both runs used the default artifact directory. These were infrastructure failures during context cleanup. The final production run uses a separate output directory. A production integration run also exposed an ambiguous test locator matching both loading and success statuses; the assertion now targets the exact success message and retains its payload checks.

## Remaining risks

The original audit's R03 (schema-integrity fixture drift), R04 (development dependency advisory) and R05 (other undocumented pre-existing API routes) remain outside this follow-up. Interview scheduling/scorecard integrations and other original operational limitations are unchanged. Passing tests do not establish that the project is bug-free or production-ready.
