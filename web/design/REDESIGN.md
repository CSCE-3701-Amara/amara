# Amara frontend refinement — 10 November 2026

Follow-up: the original welcome and login visuals were restored at the user’s
request; registration and workspaces retain the refinement below. Public entry
now resumes validated sessions. See
[`2026-11-10-session-entry-spec.md`](../../docs/audits/2026-11-10-session-entry-spec.md).

## Direction and references

Warm, precise, readable, and confident. Retain Amara’s existing logo variants,
Egyptian origin, orange/coral and lime accents, warm neutrals, and bundled Funnel
Display and Funnel Sans. No new fonts, assets, libraries, or external font requests.
Code inputs retain their existing monospace convention for legibility.

Reference conventions (principles, not copied branding):

- [Linear](https://linear.app/features): compact navigation, efficient list rows, restrained selected states.
- [Stripe](https://stripe.com/payments): clear primary actions, aligned data and metadata, deliberate numeric hierarchy.
- [Notion](https://www.notion.com/product): readable detail views, light separators, content rather than nested cards.

## Architecture and audit

Next.js 16.4.3 / React 19.3; App Router with candidate, employer, interviewer,
and administrator route groups. Local fonts via next/font. Tailwind’s global
entry establishes the reset; CSS Modules own the shared shell, UI, form/table
patterns, and feature-specific presentation. Resource hooks and API modules
continue to drive the existing interactive screens. Native dialogs retain focus
containment, Escape handling, and focus restoration.

The audit found oversized illustrated page banners, excessive pill controls,
large numeric labels, nested cards, distracting perpetual animation, heavy
shadows, and inconsistent control sizes. It also found unstyled code-runner
controls and assessment columns that became cramped on tablets.

## Visual system

- Funnel Display: page and section headings, restrained weights and tighter tracking.
- Funnel Sans: body copy, labels, controls, and tabular numeric metadata.
- Warm paper canvas, white controls, quiet warm gray rules; orange supports the brand;
  lime identifies primary actions. Success, warning, and error tones retain meaning.
- 4/8/12/16/20/24/32/48 spacing rhythm; 7–8px controls, 10–12px raised surfaces.
- Subtle background/border transitions; no perpetual motion. Reduced motion disables transitions.
- Shared compact heading, navigation, badges, statistics strip, forms, tables,
  dialogs, loading feedback, empty states, and workspace signature.
- Sidebar collapses on desktop; existing mobile disclosure and skip link remain keyboard operable.

## Purpose-specific layouts

| Area | Refinement |
| --- | --- |
| Public home | Asymmetric editorial hero; existing artwork; numbered, lightly divided steps; distinct hiring section; restrained footer |
| Login / registration / recovery | Warm brand panel; focused unboxed form; segmented account selection; consistent error and recovery states |
| Job discovery | Filter rail and searchable result rows with aligned salary, location, and actions |
| Applications / offers | Compact statistics and tab filters; row-based history; readable detail and response areas |
| Profiles | Readable overview, quieter supporting contacts, grouped editable fields |
| Postings / modules | List/table management with compact filters and contextual lifecycle actions |
| Editors | Divided form sections, shared controls, compact action rows |
| Pipeline designer / candidate journey | Warm neutral graph canvas, preserved graph geometry, distinct stage states and inspector |
| Assessments | Question navigation and timer; tablet reflow; visually distinct code; properly styled execution inputs |
| Interviews / grading / decisions | Dense tables, clear rubrics, grouped scorecard controls and private notes |
| Administration | Compact review tables and decision sections using the shared system |

## Route inventory

All 50 existing page files consume the refined shell and/or page-specific styles.
Route parameters and query states are preserved. Existing inventory tests also
exercise loading, empty, error, published, assessment, and dialog variants.

| Route | Rendering component |
| --- | --- |
| `/admin/accounts` | Administration |
| `/admin/dashboard` | Administration |
| `/candidate/applications/[id]` | ApplicationDetail |
| `/candidate/applications/[id]/progress` | ApplicationDetail |
| `/candidate/applications` | Applications |
| `/candidate/assessments/[id]/coding` | CandidateAssessment |
| `/candidate/assessments/[id]` | CandidateAssessment |
| `/candidate/assessments/[id]/video` | CandidateAssessment |
| `/candidate/dashboard` | JobFeed |
| `/candidate/interviews/[id]` | CandidateInterview |
| `/candidate/jobs/[id]` | PostingDetail |
| `/candidate/offers/[id]` | CandidateOffers |
| `/candidate/offers` | CandidateOffers |
| `/candidate/profile/edit` | CandidateProfileScreen |
| `/candidate/profile` | CandidateProfileScreen |
| `/employer/applications/[id]/scorecard` | HiringApplications |
| `/employer/assessments` | AssessmentInstances |
| `/employer/feedback/[id]` | FeedbackAuthor |
| `/employer/grading` | AssignedGrading |
| `/employer/interviewers` | InterviewerManagement |
| `/employer/interviews` | Interviews |
| `/employer/jobs/[id]/applications` | HiringApplications |
| `/employer/jobs/[id]/comparison` | HiringApplications |
| `/employer/jobs/[id]/finalists` | HiringApplications |
| `/employer/jobs/[id]` | PostingEditor |
| `/employer/jobs/[id]/preview` | PostingPreview |
| `/employer/jobs/new` | PostingEditor |
| `/employer/jobs` | Postings |
| `/employer/modules/[id]` | ModuleStudio |
| `/employer/modules/[id]/questions` | ModuleStudio |
| `/employer/modules/library` | ModuleLibrary |
| `/employer/modules/library/problems` | ModuleLibrary |
| `/employer/modules/new` | ModuleStudio |
| `/employer/modules` | ModuleLibrary |
| `/employer/offers/[id]` | HiringOffer |
| `/employer/offers/new` | HiringOffer |
| `/employer/offers` | HiringOffer |
| `/employer/pipelines/[id]/candidates` | StageDecisions |
| `/employer/pipelines/[id]/dashboard` | HiringApplications |
| `/employer/pipelines/[id]` | PipelineStudio |
| `/employer/pipelines` | Pipeline selection |
| `/employer/profile` | TeamProfile |
| `/interviewer/availability` | InterviewAvailability |
| `/interviewer/grading` | AssignedGrading |
| `/interviewer/interviews/[id]/scorecard` | InterviewScorecard |
| `/interviewer/interviews` | Interviews |
| `/interviewer/profile` | TeamProfile |
| `/login` | AuthRoute |
| `/register` | AuthRoute |
| `/` | Home |

## Important changed files

- `src/app/globals.css`: semantic colors, typography, radii, focus, reduced motion.
- `src/app/page.tsx`, `welcome.module.css`: full public page composition.
- `src/shared/ui/workspace.tsx`, `workspace.module.css`: workspace identity and responsive shell.
- `src/shared/ui/product-ui.tsx`, `product-ui.module.css`: shared headings, stats, badges, dialogs, footer.
- `src/shared/ui/editor.module.css`, `calendar.module.css`, `brand.module.css`: forms, tables, calendars, brand presentation.
- `src/features/auth/components/auth-shell.tsx`, `auth.module.css`: all auth and recovery views; remove obsolete animation timing effect.
- `src/features/cand/components/candidate.module.css`, `candidate-pipeline.module.css`, `code-runner.tsx`: candidate screens and code execution presentation.
- `src/features/cand/components/assessment.tsx`: presentation class on the existing language field only; preserve earlier uncommitted workflow work.
- `src/features/job/components/postings.module.css`, `postings.tsx`: posting rows and removal of redundant instructional prose.
- `src/features/invw/components/dag-editor.module.css`: canvas, toolbar, stage, and inspector presentation.
- `src/features/management/components/profile.module.css`: company/interviewer profile details.
- `src/app/(employer)/employer/pipelines/page.tsx`: compact heading and pipeline selection rows.
- `tests/browser/auth-polish.spec.ts`: new selected-surface treatment, field-boundary regression, no continuous motion, reduced motion.
- `tests/browser/design-refinement.spec.ts`: five breakpoint scenarios, six representative page families per scenario; mobile menu, skip-link focus, token loading, overflow, screenshots.
- `tests/browser/integration.spec.ts`: explicitly open the existing legacy editor disclosure before testing its save workflow.

## Test intent and review

The new checks cover missing global CSS (field borders and radii disappear),
viewport overflow, the mobile navigation open/close states, and keyboard access
to main content. Existing suites cover request payloads, failures, retries,
unsaved form values, session recovery, module editing, graph execution,
assessment submission, and dialog focus. The old assertion requiring three
perpetual decorative animations was replaced with no-continuous-motion and
reduced-motion checks. A missing stylesheet or removed mobile toggle / skip-link
target would fail the new checks; no business-logic unit tests were added.

Review scope: presentation changes, existing shared component contracts, and
changed test expectations. No new data fetching, dependencies, endpoints, state
management, auth logic, or application mutations. Earlier uncommitted backend
and workflow work was left in place.

## Verification

- TypeScript: `npx tsc --noEmit` passed; the final production build also passed its TypeScript check.
- ESLint: `npm run lint` passed, including the final contrast refinements.
- Unit tests: `npm run test:unit` — 8 passed.
- Full production browser suite: 403 passed, 0 failed, 0 skipped, 0 flaky.
- After the final input-border and code-placeholder contrast changes: 45 affected
  production visual/accessibility checks passed, 0 failed, 0 skipped, 0 flaky.
- Production build: `npm run build` passed after the final source changes.
- Whitespace check: `git diff --check` passed.
- Automated screenshot / overflow inventory: 375, 1024, and 1728px; representative
  page-family checks also cover 768, 1280, and 1440px. Auth includes 320px.
- Manual rendered review: desktop and mobile home/auth, candidate feed and
  applications, postings, pipeline designer, tablet assessments; contact-sheet
  review of profiles, offers, grading, scheduling, scorecards, and administration.
- Sample computed contrast ratios: body 17.16:1, muted body 6.36:1, primary-button
  text 13.15:1, input boundary against white 3.82:1, code placeholder 9.31:1.
  These are targeted checks, not a complete WCAG certification.

The initial development run had stale assertions for the removed animation and
pill selection, plus an existing legacy-editor test that did not open its
current disclosure. Updated expectations now match the refined presentation.
Overlapping Playwright runs also collided in the shared trace directory; the
final production run was isolated and passed. A Tailwind dev-loader cache kept
old global tokens during hot reload; restarting refreshed them. Border/radius
assertions now catch missing global CSS rather than passing an unstyled page.
Screenshots are generated into ignored `tests/artifacts/product`,
`tests/artifacts/refinement`, and the existing auth screenshot locations.

The browser suite uses the repository’s API fixtures; it does not certify a live
backend, live email/MFA delivery, production cookies, camera hardware, or Judge0.
Existing feature-unavailability notices remain visible and accurate.
