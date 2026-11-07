# Analytics feature

`ana` is the frontend boundary for analytics and evaluation views. It currently
contains the feature entry point and type placeholders while the analytics API
contract is being completed.

## Boundary

Analytics reads projections owned by other modules; it must not write
candidate, job, or interview-workflow state directly. Add API helpers under
`api/`, response types under `types/`, and presentational components under
`components/` as endpoints become available.

Until a route is backed by an authenticated API contract, keep the screen
explicitly unavailable rather than inventing metrics or using fixture data in
production UI.