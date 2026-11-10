"""Render the audited full-file inventory without inferring verification from routes."""
from pathlib import Path
import json

root = Path(__file__).resolve().parents[1]
inventory = json.loads((root / 'design/figma-inventory.json').read_text())
lines = ['# Full-file Figma coverage', '', '**Inventory verified; implementation and live API verification remain incomplete.**', '']
lines += [note + '\n' for note in inventory['notes']]
lines += ['Statuses describe visual implementation separately from API integration. `Blocked` design rows retain existing UI but lack a fresh reference. `Implemented but unverified` rows have UI and retrieved references; browser smoke checks do not establish complete visual parity. No endpoint is labeled Connected and verified without real backend evidence.', '', 'See [verification evidence](VERIFICATION.md), [accessibility and React review](AUDIT.md), and [backend integration](../INTEGRATION.md).', '', '| Figma page | Frame / distinct state | Node | App route / UI state | Method / endpoint | Implementation status | API integration | Verification evidence / dependency |', '| --- | --- | --- | --- | --- | --- | --- | --- |']
def cell(value):
    return str(value).replace('|', '\\|').replace('\n', ' ')
for s in inventory['screens']:
    values = [s['page'], s['name'], f"[{s['id']}]({s['url']})", f"`{s['route']}`", s['operations'], s['designStatus'], s['integrationStatus'], s['evidence']]
    lines.append('| ' + ' | '.join(cell(v) for v in values) + ' |')
lines += ['', '## Exclusions', '', '| Figma page | Frame | Node | Status / reason |', '| --- | --- | --- | --- |']
for s in inventory['exclusions']:
    link = f"https://www.figma.com/design/{inventory['fileKey']}/Amara?node-id={s['id'].replace(':','-')}"
    lines.append('| ' + ' | '.join(cell(v) for v in [s['page'], s['name'], f"[{s['id']}]({link})", 'Excluded with reason: ' + s['reason']]) + ' |')
lines += ['', '## Collection containers', '', 'These organize their child dialogs on the canvas and do not need separate routes.', '']
for s in inventory['collections']:
    lines.append(f"- {s['page']}: {s['name']} (`{s['id']}`)")
lines += ['', '## References still blocked by Figma access', '', 'Provide renewed MCP access or exported design context/screenshots for these exact nodes. Their pages were enumerated; fresh frame references could not be retrieved.', '']
for page in inventory['pages']:
    blocked = [s for s in inventory['screens'] if s['page'] == page['name'] and s['referenceStatus'].startswith('Blocked')]
    if blocked:
        lines += [f"### {page['name']}", '']
        lines += [f"- [{s['name']} ({s['id']})]({s['url']})" for s in blocked]
        lines.append('')
lines += ['## Design decisions and shared code', '', '- Figma remains the visual source: Funnel Sans/Funnel Display, local brand assets and existing orange-to-coral gradient tokens are reused.', '- Updated the DAG palette, toolbar, selected nodes and inspector. Downloaded 16 exact SVG assets; intrinsic dimensions are retained. Dependencies determine node columns without inventing edges. Added bounded undo/redo, fit-to-view and browser leave protection for local graph edits.', '- Graph persistence, publication, privacy policy and migration remain unavailable. Unsupported controls are labeled or disabled; loaded ordered-stage editing stays separate from local DAG edits.', '- Replaced the creation form at `/employer/offers` with the dedicated overview and URL-backed filters/pagination. `/employer/offers/new` remains the creation destination. The offers backend is planned; the contract lacks candidate display identity.', '- Wired candidate profile GET before mounting editable defaults. Existing company/interviewer reads are retained; profile writes patch edited fields. Module edits preserve behavioral and minimum-test settings from the current backend.', '- Added reviewed threshold changes through the registered POST preview / PUT apply operations, including exact preview echo, stale request cancellation, conflict refresh and duplicate-submit prevention. Interview outcome POST is connected when an authorized record is available; discovery is still missing.', '- Removed header illustrations from screens whose retrieved references omit them. Shared primary actions use Figma green with dark text; orange-to-coral remains in the gradient slots.', '- With no mobile frames, palettes/cards/forms stack and tables/graphs retain local scrolling. Toolbar targets stay at least 44px; graph keyboard movement, shared visible focus, reduced-motion behavior and dialog focus are retained.', '- No new dependency, font, icon package, route, backend source change or database write was introduced.', '']
(root / 'design/COVERAGE.md').write_text('\n'.join(lines))
