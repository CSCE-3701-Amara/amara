# Frontend audit and refactoring

Audit date: 2026-11-08. The working tree already contained the preceding Figma implementation; that work is preserved. A source/config snapshot and 26 representative screenshots were copied to `/tmp/amara-refactor-baseline` before refactoring.

## Architecture inspected

The frontend has 138 source files (85 TSX, 40 TS, 13 CSS), including a generated shader. Next.js 16.4.3 and React 19.3.0 use the App Router, Server Component pages/layouts, Cache Components and partial prefetching. Interactive screens use feature-local Client Components. CSS Modules and local fonts/assets define the existing UI. State is component-local; shared hooks handle cancellable reads, edited-field profile PATCH and unsaved-change warnings. There is no global state library.

Same-origin browser requests forward the backend's HttpOnly session through Next rewrites. Keeping this established flow avoids an unjustified server-side auth/data migration. `apiRequest` handles credentials and both backend error envelopes. Generated contract types coexist with deliberately narrower runtime projections. The existing production browser suite has 388 cases, including all inventoried workspace states at three widths.

An import-graph scan found no source cycles. Root pages/layouts already follow the Server Component convention. Effects for media, GPU rendering, timers, focus and abortable requests have concrete browser lifecycle purposes. No blanket hook removal or memoization is warranted.

## Findings and plan presented before code changes

| Priority | Finding | Action |
| --- | --- | --- |
| High | Pending camera acquisition can resolve after cleanup; recorder callbacks can create an object URL during teardown. | Isolate recording lifecycle; stop late streams and dispose callbacks/URLs. Preserve recording UI. |
| High | Collapsed server-supported pipeline editor mounts and duplicates pipeline/stage/module reads. | Mount on first expansion, retain its state on later collapse; measure requests before/after. |
| Medium | `evaluations.tsx` (872 lines) combines six independent route screens; `hiring.tsx` (708) combines unrelated hiring/offer/feedback/team screens. | Move each route-sized component to an explicit module and update page imports without changing markup. |
| Medium | Assessment (955 lines) and pipeline (925) combine state, presentation and browser-resource handling. | Separate substantial existing responsibilities; retain feature ownership and existing data flow. |
| Medium | API error JSON is implicitly `any`; multiple catches use `as Error`. | Narrow error-envelope values from `unknown`; share error normalization while preserving valid backend messages/codes. |
| Medium | Two datetime-local converters repeat timezone logic. | Use one defensive shared converter; preserve valid input behavior. |
| Low | An unreferenced legacy module editor duplicates the live module studio; candidate API wrappers and 18 empty scaffolds have no importers. | Remove only modules verified unreferenced by TypeScript import-graph and repository search. |
| Low | Broad TypeScript include scans ignored Figma export/test artifacts. | Exclude generated diagnostic directories; keep source/tests/config and required Next types checked. |
| Medium | Shared heading rules depend on CSS chunk order; file extraction exposed a visual regression. | Strengthen the existing section/panel selectors so their existing values win consistently. Compare desktop/mobile screenshots. |
| Low | Editor `.list` and `.row` styles have no consumers. | Remove after a parsed CSS-import/property scan confirmed no references or dynamic accesses. |
| Retain | Every declared dependency is used; React/React DOM are deduplicated. Local installation has extraneous transitive WASM packages. | No manifest/library change or environment-dependent package deletion. |

No explicit application `any` types were found outside generated/vendor code. The stricter exploratory `--noUnusedLocals --noUnusedParameters` check flagged generated Next validator and ignored Figma exports, not application source; those flags were not added to hide or bypass failures. The normal lint/type-check/build baseline is green.

## Validation strategy

Keep the existing UI, DOM structure, CSS classes, routes, payloads, auth and permission boundaries. Run lint, TypeScript, the production build and the full browser suite. Add focused checks for resource cleanup and deferred editor mounting. Compare preserved screenshots for moved screens, and report measured request/bundle differences without claiming an unmeasured latency or Core Web Vitals improvement. Backend startup and missing API dependencies remain outside this task.

## Refactoring completed

- Split six evaluation screens into `interviews`, `assigned-grading`, `interview-scorecard`, `interview-availability`, `assessment-instances` and `stage-decisions`. Split hiring screens into `hiring-applications`, `hiring-offer`, `feedback-author` and `interviewer-management`; `planned-read` shares the existing read-state behavior. Pages import the focused module directly. All 12 moved hiring/evaluation function bodies are unchanged.
- Reduced route orchestration in `pipeline-studio.tsx` from 925 to 233 lines by moving the existing graph to `pipeline-graph.tsx`. The ordered editor loads on first expansion and stays mounted to preserve edits. Explicit DAG state and interactions are unchanged.
- Reduced `assessment.tsx` from 955 to 547 lines. The existing clock, input, advanced workspace and video recorder have focused feature-local files. AST comparisons confirm the non-media functions are unchanged apart from error normalization.
- Hardened video teardown against pending camera permission, late playback resolution and recorder stop callbacks. Media tracks and object URLs are disposed; duplicate permission requests are blocked. Activity hide/show is accounted for without adding persistence.
- Narrowed API error-envelope JSON from `unknown`, normalized thrown errors centrally, and represented optional HTTP metadata accurately in read-hook errors. Existing valid backend messages, codes, status and request IDs are preserved.
- Consolidated the two timezone-aware datetime-local converters. Invalid input no longer throws during posting rendering; valid date conversion is unchanged.
- Added focus and error association for a blank confirmation reason. Modal markup, copy, dismissal and return behavior are retained.
- Removed 18 empty TS scaffolds, the unreferenced 270-line legacy module editor and unused candidate API wrappers. The live module studio and candidate request flows remain. Removed unused editor list/row CSS after checking all import aliases and finding no dynamic class accesses.
- Excluded ignored diagnostic artifact directories from the TypeScript source set. No application/test errors are suppressed, and no lint/TypeScript rules were disabled.
- Kept all dependencies, fonts, assets, API payloads, routes and backend files unchanged. No new state-management/data-fetch library, server-data migration or design system was introduced.

## Validation results

The final production browser run passed **392/392**, with no failures, skips or flaky cases, in 131 seconds. Lint passed without warnings; `tsc --noEmit`, the standard Turbopack production build and `git diff --check` passed. Four new browser checks cover deferred reads/edit retention, malformed API error recovery, focus/error association and late camera permission cleanup. API responses are mocks. The camera lifecycle check uses controlled browser APIs; physical camera/microphone recording was not exercised.

Twenty-six representative screenshots were compared at 1728px and 375px: **25 pixel-identical**, with the remaining graph screenshot differing by 21 pixels around SVG ports (0.00052% of pixels). All dimensions match. The extraction initially exposed a heading-cascade regression; focused comparisons and the final full suite were repeated after fixing selector specificity. The component font/spacing values remain the baseline values. Only editor CSS changed: deterministic section/panel heading precedence and removal of the proven unused list/row rules.

AST comparisons confirmed unchanged bodies for all 12 moved hiring/evaluation functions and unchanged non-media assessment/graph functions, except shared error normalization. The import graph has no cycles; the ordered editor is reached through an explicit dynamic import. Package manifests, lockfile, Next configuration, fonts/assets and backend files are unchanged.

[Machine-readable evidence](refactoring-verification.json) includes test statistics, all screenshot comparison results and before/after measurements.

### Performance measurements

On a fresh production pipeline page with its stage-configuration disclosure closed:

| Request | Before | After |
| --- | ---: | ---: |
| Job pipeline | 2 | 1 |
| Version stages | 2 | 1 |
| Modules | 2 | 1 |
| Transitions | 1 | 1 |
| Total API reads | **7** | **4** |

This removes three initial API requests (43%). Expanding the ordered editor still loads its current server data; collapsing/reopening preserves edits and does not repeat those reads.

| Route | Fetched JS body bytes, before → after | Estimated gzip bytes, before → after | Script requests, before → after |
| --- | --- | --- | --- |
| Pipeline | 632,695 → 623,190 | 193,032 → 193,054 | 14 → 17 |
| Interviews | 607,221 → 597,137 | 185,154 → 184,702 | 13 → 16 |
| Assigned grading | 584,226 → 574,117 | 177,028 → 176,556 | 10 → 13 |

Measured uncompressed JavaScript decreases by about 9.5–10.1 KB; compressed size is essentially unchanged. Route separation creates three additional script chunks in these samples. The defensible performance improvement is reduced redundant API work, not a claimed loading-time/Core Web Vitals gain. Samples use identical fixtures and include shared framework code; gzip totals are local estimates, not captured wire-transfer sizes.

## Remaining technical debt

- The cohesive graph renderer (~716 lines), module authoring and multi-view hiring application screen remain substantial. Further decomposition should follow real editing/API boundaries rather than a line-count target.
- Backend startup/auth wiring and planned endpoints documented in `INTEGRATION.md` still block live verification; browser API tests are mocked.
- Successful API bodies still rely on contract projections, with explicit auth response checks. Broader runtime schema validation would need an agreed endpoint-by-endpoint strategy; no validation dependency was added speculatively.
- Client-side reads remain the established authenticated transport pattern. A server-first migration would need verified session forwarding, cache boundaries and mutation invalidation, and is intentionally outside this behavior-preserving cleanup.
- Some shared workspace typography remains broad. The verified screen regression is fixed; a full CSS token/cascade redesign is deferred to avoid changing the product UI.
- Local `next start` emits the existing standalone-output warning; Docker already uses the standalone server correctly. The existing Playwright production harness is retained.
- Installed transitive WASM packages were reported as extraneous by `npm ls`; they are not declared dependencies. A clean CI `npm ci` is preferable to changing the lockfile for one local installation.
