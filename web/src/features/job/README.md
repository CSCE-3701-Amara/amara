# Job feature

The `job` feature owns employer posting management and the candidate posting
detail/application entry points. Employer pages are mounted under
`/employer/jobs`; candidate details are mounted under `/candidate/jobs`.

## Connected posting flows

The API helpers in `api/job.api.ts` cover:

- `GET /v1/jobs` and `GET /v1/jobs/{id}` for company-owned postings.
- `POST /v1/jobs` and `PUT /v1/jobs/{id}` for drafts and edits.
- `POST /v1/jobs/{id}/publish`, `/republish`, and `/archive`.
- `DELETE /v1/jobs/{id}` for company managers or system administrators.
- `GET` and `POST /v1/jobs/{id}/preview` for saved or unsaved previews.
- `POST /v1/jobs/{id}/apply` for candidate submission.

Posting forms validate salary and application dates before sending, encode
resource IDs in paths, and show server errors through shared feedback
components. Preview requests must not mutate the saved posting.

The backend scopes company operations from the validated session's company ID;
the browser must never send an account or company ID as an authority claim.