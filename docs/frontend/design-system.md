# Frontend design system

Amara's visual system is implemented in `web/src/app/globals.css`,
`web/src/shared/ui/`, and feature CSS modules. The system favors a warm,
high-contrast surface with a bright green action accent and generous rounded
cards.

## Typography

The root layout loads the local Funnel Sans and Funnel Display font files:

- Funnel Sans is the body, form, navigation, and utility font.
- Funnel Display is available for display treatments and prominent headings.

Use the existing font variables (`--font-funnel-sans` and
`--font-funnel-display`) instead of importing a remote font.

## Core tokens

Global tokens in `globals.css` include:

| Token | Purpose |
| --- | --- |
| `--background` | Warm page background |
| `--surface` / `--surface-muted` | Primary and subdued surfaces |
| `--foreground` / `--text-muted` | Main and secondary text |
| `--accent` / `--accent-hover` | Primary action and hover states |
| `--orange` / `--coral` | Brand illustration and attention colors |
| `--border` | Field, card, and divider borders |
| `--radius-card` | Large card radius |
| `--radius-field` | Form control radius |
| `--card-shadow` | Elevated card treatment |

Prefer tokens and existing module styles over new hard-coded values. Preserve
the system's visible focus ring and do not remove it for visual matching.

## Shared components

Use the shared components before creating feature-specific equivalents:

- `Brand`: consistent logo and wordmark treatment.
- `Workspace`: role navigation, responsive sidebar, logout, and error area.
- `ProductHeading`: page title, description, optional artwork, and compact mode.
- `StatusBadge`, `Stats`, and `EmptyState`: common data presentation.
- `BusyState` and `ApiFeedback`: loading and recoverable error states.
- `ConfirmDialog`: native modal, focus return, Escape dismissal, and optional
  reason field.
- `UnavailableNotice`: explicit boundary for unsupported persistence.

Decorative images use empty alt text. Content images need meaningful alt text.
Use semantic headings, labels, table headers, and `aria-current`/`aria-live`
states as demonstrated by the shared components.

## Responsive and motion behavior

Workspace navigation collapses on small screens and can be opened through the
mobile menu. Design for the existing mobile and desktop Playwright viewports;
do not assume the 1728px Figma frame is the only viewport.

Animations use the shared easing and must honor
`prefers-reduced-motion: reduce`. Reduced motion disables or freezes decorative
animation and changes smooth scrolling to immediate scrolling. Interactive
focus, hover, and disabled states must remain understandable without motion.

## Figma assets

Figma-exported SVGs and fonts are stored under `web/public/`. Shader source is
kept separately under `web/src/shared/lib/figma-shaders/`; generated shader
code should not be reformatted or mixed with application logic.