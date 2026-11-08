# Accessibility and React review — 2026-11-08

Applied `web-design-guidelines`, then `vercel-react-best-practices`. The [upstream interface guidelines](https://raw.githubusercontent.com/vercel-labs/web-interface-guidelines/main/command.md) were fetched on 2026-11-08. Figma remains the visual source. `frontend-skill` was not installed; the available `frontend-design` skill supplied responsive/interactivity guidance.

Fixed findings:

```text
src/features/cand/components/profile.tsx:186 - populated biography had unstable implicit label; explicit accessible name added.
src/features/management/components/evaluations.tsx:699 - stage selector label included option content; explicit Stage accessible name added.
src/features/invw/components/pipeline-studio.module.css:88 - shared flex button styling compressed node text into a row; node-specific block layout restores readable cards.
src/features/invw/components/threshold-review.tsx:79 - canceled preview could leave pending state retained; request-identity cleanup now clears it.
src/features/management/components/profile.tsx:139 - async profile defaults could replace edits; form mounts after read completes.
src/features/invw/components/module-studio.tsx:294 - current behavioral/minimum-test settings could be dropped by PUT; preserved with authored questions.
src/features/invw/components/modules.tsx:52 - existing module editor now preserves the same current backend fields.
src/features/invw/components/pipeline-studio.tsx:191 - repeated edge/node scans during dragging replaced with a node map and memoized edge counts.
```

All new icon controls have accessible names and intrinsic-size decorative images. Toolbar targets remain at least 44px, even though the Figma toolbar is denser. Graphs retain local scrolling and Alt+arrow movement, with bounded undo/redo. Shared keyboard focus, mobile navigation, reduced-motion CSS, modal focus/dismissal and return behavior remain intact. No added animation, dependency or server credential enters the client.

Threshold requests abort on stage changes, require exact preview review, prevent duplicate submits, preserve errors and refresh after success/conflict. Profile reads mount initial fields once; edits stay local until PATCH. No profile or answer data is placed in browser storage. Offers filters/pagination use URLs and semantic links. Primary actions use Figma green with dark text; gradient slots retain the approved orange-to-coral tokens.

Final checks: lint and TypeScript pass; production build passes; all 388 browser checks pass. Backend-dependent paths remain unverified because the current API cannot start. No complete WCAG or pixel-parity certification is claimed for unread Figma states or unavailable backend flows.

---

# Earlier accessibility and React review

Reviewed the changed frontend code using the requested `web-design-guidelines`, followed by `vercel-react-best-practices`. The current interface guidelines were fetched from [the upstream source](https://raw.githubusercontent.com/vercel-labs/web-interface-guidelines/main/command.md) on 2026-11-07. Figma controls visual decisions; the audit does not replace missing design references.

## Fixed findings

- `src/features/cand/components/assessment.tsx:224` — explicit textarea labels keep accessible names stable after entering a controlled value; code inputs disable spell checking and browser completion.
- `src/shared/ui/product-ui.tsx:81` — nonnumeric statistic values use a readable text size; numeric values use tabular numerals. This fixes verified dashboard overflow at 375px.
- `src/shared/lib/use-profile-update.ts:31` — only edited fields are sent; forms preserve values on failure and guard concurrent submissions. Blank required names are rejected and focused.
- `src/shared/lib/use-unsaved-changes.ts:5` — unsaved profile/assessment drafts register browser leave protection; profile Cancel links request discard confirmation.
- `src/features/job/components/postings.tsx:17` — lifecycle actions use buttons and confirmation dialogs; navigation uses links; empty and real error states provide a next action.
- `src/features/job/components/postings.module.css:1` — cards collapse on narrow screens, long text wraps, focus remains visible through shared styles, and hover states provide feedback.
- `src/features/management/components/profile.tsx:147` — fields retain labels and correct input types; pending saves disable editable controls and expose `aria-busy`; results use live status messages.
- `src/shared/ui/route-dialog.tsx:17` — successful account, deletion and assessment navigation cannot be overwritten by the modal's close callback. Subsequent dialogs can still dismiss normally.

## React / Next.js review

- `src/features/management/components/hiring.tsx:82` — dashboard reads the registered posting-scoped endpoint directly, eliminating an unnecessary preliminary pipeline lookup and using the actual response schema.
- `src/features/cand/components/assessment.tsx:224` — independent answer writes run concurrently with `Promise.allSettled`; all writes finish before retry/submission becomes available. Successful responses are retained when another save fails. Submission waits for save completion and immediately uses the backend receipt.
- `src/shared/lib/use-api-resource.ts:7` — existing cancellation of stale reads is preserved; an optional primitive reload key permits a submitted-state transition to refetch the attempt.
- `src/shared/lib/use-unsaved-changes.ts:5` — the leave listener exists only while dirty, depends on a boolean, and is removed on cleanup.
- No new runtime dependency, icon package, request-global private state, browser credential storage, or duplicated API transport was added. Existing local fonts, CSS Modules and server-only API configuration remain in use.
- Regenerated the existing OpenAPI projections with `python3 web/scripts/contract-types.py`, then formatted them. Runtime dashboard fields use an explicit feature type because the registered handler differs from the planned contract projection.

## Verified boundaries

Browser tests exercise keyboard focus, modal entry/Escape/return, DAG keyboard movement and cycle rejection, mobile navigation, viewport overflow, loaded images, API error recovery, draft preservation, and response/payload schemas using mocks. No automated accessibility score or exhaustive assistive-technology certification is claimed. Pixel-perfect parity for every frame remains unverified; 51 fresh frame references were inaccessible after the Figma quota limit.

Some retrieved Figma frames also flag animation. Their separate motion context
could not be retrieved after quota exhaustion. Existing animation behavior was
preserved; reduced-motion browser checks do not establish Figma motion fidelity.
