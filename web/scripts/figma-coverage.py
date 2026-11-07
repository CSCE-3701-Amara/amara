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
lines += ['## Design decisions and shared code', '', '- Figma remains the visual source of truth: Funnel Sans/Funnel Display, existing local brand assets, warm white workspace surfaces, and existing orange-to-coral gradient tokens were reused.', '- Updated the unified posting feed to its current two-column card layout with state-specific lifecycle actions. Removed illustrations absent from the retrieved posting/editor/detail references.', '- Added company overview and edit states on the existing profile route (`?edit=true`); shared edited-field PATCH behavior prevents unloaded values from being erased.', '- Kept the DAG editor and its explicit transitions, branches, joins, zoom, dragging and keyboard controls. Graph persistence remains unavailable because its handlers are missing.', '- Responsive choices where Figma has no mobile frame: stack cards and form columns; retain labeled controls and horizontal table/graph scrolling; keep the existing mobile navigation and dialog focus behavior.', '- Added stable explicit textarea labels, unsaved-change protection, save/submission guards and readable unavailable statistics. No new UI dependency, font, icon package, or backend change was introduced.', '- The PATCH response updates the current profile screen; missing read endpoints prevent complete overview/reload behavior. Unsupported video upload, offers, grading, administration, feedback and other mutations remain distinguishable from supported actions.', '']
(root / 'design/COVERAGE.md').write_text('\n'.join(lines))
