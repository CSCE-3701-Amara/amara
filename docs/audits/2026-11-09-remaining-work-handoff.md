# Remaining work handoff — 9 November 2026

Report for the next engineer/agent continuing after the prior session hit its execution limit. Read it together with the [codebase audit](2026-11-09-codebase-audit.md), its [DAG/execution follow-up](2026-11-09-dag-execution-followup.md) and the machine-readable verification records (`*-verification.json` in this directory). The prior session's Docker-verify notes in this report come from its last recorded message: it had verified the lightweight sandbox against real Docker, found the official Judge0 image still missing, and stopped its PostgreSQL/Redis containers cleanly before the limit was reached.

## Start here: repository state

- Branch `FIXING` contains **6 local commits that have not been pushed** (`git push origin FIXING`), including the DAG execution feature, auth fixes, migration/CI fix, job feed fixes and the web proxy fix. Push before doing anything else so the work is not stranded on a laptop.
- Working tree is clean; `go build ./...`, `go vet ./...` pass at HEAD.
- Verification already recorded: complete Playwright run **396 passed / 0 failed**, Go suite with `-race` and real PostgreSQL **PASS**, gosec/staticcheck/govulncheck clean, frontend build/lint/tsc clean (`2026-11-09-dag-execution-verification.json`).

## A. Immediate unfinished task — official Judge0 Compose verification

This is exactly where the previous agent stopped. Status and evidence:

| Item | State |
| --- | --- |
| Lightweight sandbox (`cmd/sandbox`, Python/JS/TypeScript) | **Verified** against real Docker: 6 actual execution cases passed (stdin, runtime diagnostics, infinite-loop termination, blocked networking). |
| Official `judge0/judge0:1.13.1` image | **Never pulled** — the download stalled. Not present in `docker images`. |
| `amara-judge0-db-1` / `amara-judge0-redis-1` | Stopped cleanly (`Exited (0)`), data intact; compose stack exists at `deployments/judge0/compose.yml` and starts via [`../../scripts/start-judge0.sh`](../../scripts/start-judge0.sh). |
| Full Judge0 startup, Go/C++/full-TS compiler execution, host cgroup compatibility | **NOT verified** — explicitly noted in the follow-up report, "Execution environment" section. |

Resume procedure:

1. Check disk space first (≈19 GB free, disk ~90% full). Confirm `docker pull judge0/judge0:1.13.1` completes; if the registry connection stalls again, switch to a registry mirror or retry over a different network.
2. Run `./scripts/start-judge0.sh` (gitignored credentials are generated once in mode 600). Then run `curl http://127.0.0.1:2358/languages` and verify the language IDs match the README defaults (Python 71, JS 63, TS 74, Go 60, C++ 54).
3. Follow the [Judge0 README](../../deployments/judge0/README.md) to point both API **and worker** (`APP_ENV=development`, `JUDGE0_URL`, `JUDGE0_ALLOW_LOCAL_HTTP`, `JUDGE0_LANGUAGE_IDS`) to the same instance; configure them identically.
4. Prove real end-to-end runs of candidate code through `/v1/assessment-attempts/:attemptId/run`:
   - Python, JavaScript, full TypeScript compilation, and Go (60) / C++ (54), which the lightweight runner cannot execute,
   - the submission block on networking is still enforced with the official worker containers,
   - compile/runtime diagnostics actually appear in run output.
5. Explicitly verify Judge0's privileged-container / cgroup compatibility on this host before declaring the check done (the README documents that these requirements are real, not aspirational).

**Definition of done:** flips the follow-up report's "Full Judge0 startup, compiler execution and host cgroup compatibility are **not verified**" and flips the JSON field `"official_judge0": "not verified; image download stalled"` to a passing result (or a precise, reproducible-environment-bug report). Update both files in `docs/audits/` when done.

## B. Confirmed remaining issues from the audit (R03–R05)

These were deliberately left unfinished and are the only remaining defects from the audit verdict "REQUEST CHANGES": R01/R02 are resolved. Here are the remaining:

| ID | Severity | What is wrong | What "done" requires | Owner note |
| --- | --- | --- | --- | --- |
| R03 | Medium | `tests/database/schema_integrity.sql` is written for a database model that no longer applies — the code does not create `roles`, `users`, tenant membership or resume tables. Running it against a fresh database fails at line 27 (`relation "roles" does not exist`). | Reconcile the intended database model with the migration history and rewrite the integrity assertions for the real schema; the fixture must **not** be deleted or watered down to go green. Do not create speculative tables. Failing check today is logged in `2026-11-09-verification.json` (`schema_integrity`). | Best done jointly with the docs/SRS owner; needs a design decision, not a patch. |
| R04 | Moderate (dev dependency) | Installed `minimatch@10.2.6` (via TypeScript ESLint / `eslint-config-next`) is covered by [GHSA-g3hv-75wg-m249](https://github.com/advisories/GHSA-g3hv-75wg-m249); `npm audit` shows 8 affected dev nodes, 0 production nodes. The fix requires the next major, up to the consumer's current range. | Pick a fix: ESLint config upgrade route, an `overrides` pin, or an accepted-risk document. Whichever is chosen must be validated (`npm run lint`, `npx tsc --noEmit`, `npm run build`, full Playwright run) after the change. | Timeboxed; a dependency-only risk (no exploit was shown). |
| R05 | Medium (docs contract) | `docs/api/openapi.yaml` does not define: registered candidate `GET /v1/applications/{application_id}/pipeline-progress`, company `POST /v1/assessment-attempts/{attemptId}/result-disclosure`, nor the legacy login alias. The routes exist in the code. | Add the three operations to the OpenAPI description (real schemas and status codes as implemented). No wire changes. Validate that the YAML parses and that SDK generation does not break. | The easiest of the three; do first. |

Note: the 17 Semgrep findings from the original audit are classified there (14 mutable GitHub Action references + 3 false positives). Pinning the action SHAs is an optional supply-chain hardening chore, not blocking.

## C. Documented functional gaps (not bugs; planned features)

From the audits' "operational limits" sections — the next agent should treat these as backlog, not defects:

1. **Company suspension lifecycle (FR-JOB-018) — the biggest functional debt.** A05 closes discovery into new admissions. Application/pipeline **suspension** when an employer is banned, or a posting is archived, is *not* implemented as a complete lifecycle operation. Requires a defined contract for worker advancement and deadlines during the suspension, plus dedicated integration coverage.
2. **Incomplete workflows** (explicitly named "partly or wholly planned"): offers/awards stage, interview scheduling/meeting integration, scorecard submission, reminders, admin features, the entire **ANA** analytics placeholder, candidate video upload, and the delivery API for the candidate saved-answer reload.
3. **Data reconciliation:** e-mail duplicates created before the A10 deduplication fix are **not** automatically cleaned up; they need a one-off data-decision by a human, not an application fix.
4. **Environment-level verification that has never happened**: real SMTP, deployed Judge0 IDs/sandboxing in production config, Supabase grants/RLS, public TLS/security headers, deployment images, multi-replica behavior, PG 16 (CI) vs PG 17 (Supabase config) vs local PG 18.6, load-related registration/stress testing, distribution of a limiter beyond a single instance.
5. **Coverage is uneven** (documented in the audit): INVW/JOB database package coverage is low and should not be read as system-wide coverage, but new work on those packages should raise the tests rather than the artifacts.

## D. Local environment notes

- Two amara containers were stale at the inspection (2 days old): `amara-web-1` exited (143); `amara-api-1` **in a restart loop** before being inspected. Check `docker logs amara-api-1` for the problem (likely requires newer `JUDGE0_*` / `SMTP_*` configuration that the compose environment selects by default).
- Supabase local services were mostly stopped after the audit session.
- Prior session's raw test logs are available under `/tmp/amara-audit-*` (may disappear on reboot; copy any needed ones first).

## E. Suggested order of work

1. `git push origin FIXING` (unblocks everything else for review).
2. Task **A** — finish official Judge0 verification (small and timeboxed; un-flips the only "not verified" clause in the follow-up report).
3. Task **B: R05** (OpenAPI docs — hours), **R04** (dependency decision — hours), **R03** (schema reconciliation — days, needs design).
4. Anything from **C** in the order the product team ranks: the company suspension lifecycle is the only confirmed backend gap that has a requirement (FR-JOB-018) waiting on it.
5. Re-run the full gate suite (`AMARA_TEST_DATABASE_URL=… go test -race ./...`, complete Playwright run, analyzers) and record results in a new entry in `docs/audits/` before declaring anything done — the previous audit retained failures rather than suppressing them; keep that discipline.

_Produce this report from the audit artifacts; it does not invalidate any check that has already passed._
