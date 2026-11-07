# Candidate feature

The `cand` feature owns candidate-facing job discovery, applications, profile
views, assessments, interviews, offers, and application progress. Candidate
pages are mounted under `/candidate`.

## API boundary

The current connected flows use:

- `GET /v1/jobs/feed` for authenticated, filtered, cursor-paginated jobs.
- `POST /v1/jobs/{id}/apply` for application submission.
- Candidate application, assessment, interview, and progress endpoints defined
  in `api/` and `shared/lib/contract.types.ts`.

Use the candidate session middleware on the backend for ownership. The client
must use the shared API client, preserve server error messages, and retain
unsent form data when a request fails.

Keep candidate-only UI in this feature; reusable controls belong under
`src/shared/ui/`.