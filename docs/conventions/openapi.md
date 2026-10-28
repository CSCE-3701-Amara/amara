# AMARA API Conventions
 
Owner: Ziad Eliwa | Status: In Review | Last updated: 2026-10-29 | Jira: AMARA-1311

**Status:** Draft — team review required  
**Applies to:** All AMARA HTTP API modules  
**API style:** REST-style JSON over HTTPS  
**Specification:** OpenAPI 3.1 (`docs/api/openapi.yaml`)

## 1. Purpose and scope

These conventions keep AMARA's API predictable across modules (CAND, AUTH, JOB, INVW, ANA, and Course Process & Tooling). They define the public HTTP interface, not how domain logic or database tables are structured.

**Normative terms:** **MUST** is required, **SHOULD** is recommended unless there is a documented reason not to, and **MAY** is optional.

Each API change MUST update its OpenAPI contract and relevant tests in the same pull request. If this document conflicts with an approved API design decision, update this document rather than silently diverging.

## 2. Transport and content types

- Production endpoints MUST use HTTPS. Local development MAY use HTTP.
- Requests and successful JSON responses MUST use `application/json` unless a specific endpoint documents another media type (such as file uploads).
- Error responses MUST use `application/problem+json` (RFC 9457).
- UTF-8 is the standard encoding. Clients SHOULD send `Accept: application/json` for ordinary successful responses.
- API endpoints MUST NOT expose internal database schema, stack traces, secrets, or infrastructure details.

## 3. URLs, resources, and versioning

Base prefix: `/api/v1`. The version is part of the public API path. The `/v1` prefix is reserved for backward-incompatible changes; compatible additions do not require a new version.

Rules:

- Use plural nouns for collections: `/jobs`, `/applications`, `/candidates`.
- Use lowercase kebab-case for multiword URL segments: `/interview-stages`.
- Use path parameters for resource identity: `/jobs/{jobId}`.
- Avoid verbs such as `/getJobs` or `/createJob` for normal CRUD operations.
- Use shallow nesting where possible. Nest only when it clarifies ownership or context, e.g. `/jobs/{jobId}/applications`.
- Use query parameters for filtering, sorting, and pagination, not path segments.
- No trailing slash in canonical routes (except `/`).
- A domain operation that is not CRUD MAY use an explicit action subresource, documented in OpenAPI, e.g. `POST /api/v1/jobs/{jobId}/publish`. Such operations MUST define authorization, validation, idempotency expectations, and possible errors.

Examples:

```http
GET    /api/v1/jobs
POST   /api/v1/jobs
GET    /api/v1/jobs/{jobId}
PATCH  /api/v1/jobs/{jobId}
DELETE /api/v1/jobs/{jobId}
GET    /api/v1/jobs/{jobId}/applications
POST   /api/v1/jobs/{jobId}/publish
```

Resource IDs MUST be treated as opaque strings by clients. Their actual format is defined by each resource schema; clients MUST NOT infer business meaning from an ID.

## 4. HTTP methods and success status codes

| Method | Use | Typical success |
| --- | --- | --- |
| `GET` | Read a resource or collection; no state change | `200 OK` |
| `POST` | Create a resource or trigger a documented operation | `201 Created`, `200 OK`, or `202 Accepted` |
| `PUT` | Replace a complete resource where supported | `200 OK` or `204 No Content` |
| `PATCH` | Partially update a resource | `200 OK` or `204 No Content` |
| `DELETE` | Delete a resource where supported | `204 No Content` |

- A successful resource-creation `POST` SHOULD return `201 Created`, a representation of the created resource, and a `Location` header identifying it.
- Use `202 Accepted` only when work is accepted but not yet complete. Document how the client checks its status.
- `204 No Content` responses MUST NOT include a response body.
- `GET` MUST be safe; `GET`, `PUT`, and `DELETE` MUST preserve HTTP idempotency semantics. Repeating an idempotent operation may still produce a different status (for example, a second deletion may return `404`).
- For partial updates, AMARA initially uses `PATCH` with `application/json` and a **documented endpoint-specific update schema**. It is **not** JSON Patch (RFC 6902) or JSON Merge Patch (RFC 7396) by default. Omitted fields mean “unchanged”; each field's behavior when explicitly `null` MUST be documented. Do not accept unknown or read-only fields.

## 5. JSON field conventions

- JSON properties MUST use `camelCase`, e.g. `createdAt`, `jobId`, `firstName`.
- Boolean names SHOULD read naturally: `isPublished`, `hasCompletedTest`.
- Timestamps MUST be RFC 3339 strings with timezone information. The API SHOULD emit UTC using `Z`, e.g. `2026-10-29T10:30:00Z`.
- Date-only values MUST use `YYYY-MM-DD` when appropriate.
- Money MUST NOT be represented as imprecise floating-point values. Use integer minor units plus a currency code, or a documented decimal-string representation.
`
{
  "transaction_1": {
    "amount": "150.50",
    "currency": "EGP"
  },
  "transaction_2": {
    "amount": "10.99",
    "currency": "USD"
  }
}
`
- Enumerated values MUST be explicitly listed in OpenAPI. Prefer stable lowercase snake-case values for machine-readable enum values (e.g. `in_progress`).
- Distinguish absent fields, `null`, empty arrays, and empty strings; document their meanings. Collections SHOULD return `[]` instead of `null` when empty.
- Do not wrap every successful single-resource response in a generic envelope unless an endpoint has a documented reason. Return the resource object directly. Collection responses use the envelope defined below.

Example:

```json
{
  "id": "job_123",
  "title": "Backend Engineer",
  "isPublished": false,
  "createdAt": "2026-10-29T10:30:00Z"
}
```

## 6. Collections, filtering, sorting, and pagination

Initial strategy: offset pagination for list endpoints.

- Request parameters: `limit` (default `20`, max `100`) and `offset` (default `0`).
- `limit` MUST be a positive integer; `offset` MUST be a nonnegative integer. Invalid values return `400` with a Problem Details body.
- Use explicit, allowlisted filter parameters, e.g. `status=open`.
- Use `sort` for one field, with `-` prefix for descending order: `sort=-createdAt`. Each endpoint MUST document supported sort fields and a stable tie-breaker (usually ID).
- The response MUST include `data` and `pagination`. `total` is included when a reliable total is available; if omitted, the contract MUST say so. For Sprint 1, endpoints using this format SHOULD include `total`.
- Offset pagination can show duplicates or skips when underlying data changes between requests. Use cursor pagination later if a particular endpoint requires stable traversal at scale.

```http
GET /api/v1/jobs?limit=20&offset=0&sort=-createdAt&status=open
```

```json
{
  "data": [
    {
      "id": "job_123",
      "title": "Backend Engineer"
    }
  ],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "total": 145
  }
}
```

## 7. Error responses

All HTTP errors MUST use RFC 9457 Problem Details with `Content-Type: application/problem+json`. The core fields are `type`, `title`, `status`, `detail`, and `instance`. `detail` MUST be safe to show to clients. A stable `code` extension MAY be used for programmatic handling; clients MUST NOT parse human-readable `detail` text.

```http
HTTP/1.1 404 Not Found
Content-Type: application/problem+json
```

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "The requested job does not exist.",
  "instance": "/api/v1/jobs/job_123",
  "code": "job_not_found"
}
```

| Status | Meaning / when to use |
| --- | --- |
| `400 Bad Request` | Malformed JSON, invalid query syntax, or invalid parameter format |
| `401 Unauthorized` | Missing or invalid authentication credentials |
| `403 Forbidden` | Authenticated but not permitted to perform the operation |
| `404 Not Found` | Resource does not exist or must not be disclosed to this caller |
| `409 Conflict` | State conflict, such as a duplicate or invalid state transition |
| `413 Content Too Large` | Request exceeds a documented size limit |
| `415 Unsupported Media Type` | Unsupported request body media type |
| `422 Unprocessable Content` | Well-formed request fails semantic/domain validation |
| `429 Too Many Requests` | Rate limit exceeded; include `Retry-After` when appropriate |
| `500 Internal Server Error` | Unexpected server failure; generic client-facing detail |
| `503 Service Unavailable` | Temporary service unavailability |

For validation failures, include a documented extension such as `errors`, containing stable field identifiers and safe messages:

```json
{
  "type": "about:blank",
  "title": "Unprocessable Content",
  "status": 422,
  "detail": "One or more fields are invalid.",
  "code": "validation_failed",
  "errors": [
    { "field": "title", "message": "Title is required." }
  ]
}
```

Do not expose whether a private resource exists when doing so would leak information. Authentication and authorization failures MUST be handled consistently.

## 8. Authentication and authorization
TODO Decide on Authentication & Authorization
<!-- - Protected endpoints MUST accept `Authorization: Bearer <access_token>`.
- AMARA uses Supabase Auth for identity. The backend MUST validate JWT signatures against the trusted Supabase issuer's keys, as well as issuer, audience (when configured), expiry, and other required claims. Do not merely decode the token.
- The exact issuer, audience, accepted algorithms, and key-rotation strategy MUST be configured per environment and documented in the security implementation.
- Authentication proves identity; **authorization** checks permissions for the requested operation and resource. Every protected handler MUST enforce its relevant authorization policy, including ownership or tenant boundaries where applicable.
- Never trust a user ID, role, organization ID, or permissions supplied only in request JSON as proof of authority.
- Do not log bearer tokens, passwords, session cookies, or other secrets.
- Document public versus protected operations in OpenAPI. Use an OpenAPI HTTP bearer security scheme with `bearerFormat: JWT`.
- `401` responses SHOULD include an appropriate `WWW-Authenticate: Bearer` header; avoid revealing sensitive token-validation details. -->

## 9. Request validation and security

- Validate path/query parameters and request bodies at the HTTP boundary; validate business rules in the application/domain layer.
- Reject malformed JSON, unexpected JSON types, and disallowed input fields. Define required fields, bounds, formats, and enum constraints in OpenAPI.
- Set reasonable request body size limits and HTTP server timeouts.
- Never build SQL from unsanitized input; use parameterized queries.
- Configure CORS with explicit allowed origins, methods, and headers for the frontend environments. CORS is not an authentication or authorization mechanism.
- Apply rate limits to sensitive endpoints such as authentication, password recovery, and expensive AI operations as those endpoints are introduced.
- For operations that may be retried and create costly or duplicate effects (e.g. payments, interview invitations, AI generation), consider an `Idempotency-Key` design before implementation. Document key scope, lifetime, replay behavior, and conflict handling if adopted.

## 10. Async work and long-running operations

For work that does not complete within a normal HTTP request (such as generating assessments or processing uploaded media), prefer asynchronous processing:

1. `POST` accepts the request and returns `202 Accepted` with a job/operation identifier or status URL.
2. The worker processes the task independently.
3. The client checks a documented status endpoint (or uses another documented notification mechanism).

Do not return `200 OK` implying completion while work is still pending. Define states, failure responses, and retention for each async operation.

## 11. Observability and logging

- Generate or propagate a request ID (e.g. `X-Request-ID`) and return it in responses where practical.
- Log structured events with request ID, HTTP method, route template, status, duration, and safe error category.
- Log route templates (`/jobs/{jobId}`), not unbounded or sensitive path values, for metrics.
- Never log access tokens, passwords, sensitive applicant details, or raw assessment answers by default.
- Expose a simple liveness/readiness strategy separately from business API contracts as deployment needs mature.

## 12. OpenAPI contract and workflow

Canonical file: `docs/api/openapi.yaml` (OpenAPI 3.1).

The contract MUST describe every implemented public endpoint, including:

- Path, method, summary, tags, and operation ID.
- Path/query parameters and request body schema.
- Successful and expected error responses, including media types.
- Authentication requirements and relevant security schemes.
- Shared reusable schemas (`components/schemas`) for common models, pagination, and Problem Details.
- At least one realistic example for nontrivial operations.

**Contract-first workflow:**

1. Propose the endpoint and review naming, authorization, and status codes.
2. Update `docs/api/openapi.yaml`.
3. Implement the Go handler/adapter and application use case.
4. Add handler tests for valid and invalid requests, authentication, authorization, and documented status codes.
5. Lint and validate OpenAPI in CI; add response/schema conformance tests where feasible.
6. Review contract and implementation together in the pull request.

Generated Go types or interfaces MAY be used, but generated code MUST remain separate from domain logic. Avoid duplicating the same schema in multiple files.

## 13. Architecture boundaries (modular monolith + hexagonal)

- HTTP handlers are **inbound adapters**. They parse and validate HTTP requests, invoke application use cases, and translate results into HTTP responses.
- Application services and domain models MUST NOT depend on HTTP-specific types (`http.Request`, status codes, JSON response structures) or OpenAPI-generated transport DTOs.
- Each module owns its endpoints and application logic. The HTTP server composes module route registrars; modules SHOULD NOT reach into each other's persistence internals.
- Shared HTTP response/error helpers MAY live in a small common transport package. Avoid turning that package into a shared domain layer.
- Database entities SHOULD NOT be serialized directly as public API responses. Define explicit request/response DTOs to keep the API stable when persistence changes.

## 14. Backward compatibility

Within `/api/v1`, changes SHOULD remain backward-compatible:

**Usually compatible:** adding an optional response field, adding a new endpoint, adding an optional request parameter with safe defaults.

**Potentially breaking:** removing/renaming a field, changing a field's type or meaning, making an optional request field required, changing authorization requirements, or changing response shapes/status semantics. Adding enum values can also break clients that exhaustively switch over known values; review such changes carefully.

Breaking changes require explicit team review and a migration plan. Do not create `/api/v2` for every minor addition.

## 15. Pull request checklist

- [ ] Endpoint follows URL and method conventions.
- [ ] OpenAPI contract is updated and valid.
- [ ] Request and response schemas are explicit.
- [ ] Success and error statuses are documented and tested.
- [ ] Authentication and resource authorization are implemented where required.
- [ ] Input validation and request limits are appropriate.
- [ ] No secrets or sensitive data leak through responses/logs.
- [ ] Handler stays thin; business logic lives in the application/domain layers.
- [ ] Changes are backward-compatible or have an approved migration plan.

<!-- ## 16. Team decisions to confirm before approval

These defaults are proposals, not claims about already implemented AMARA behavior:

1. Confirm whether every module uses `/api/v1` and which endpoints are public.
2. Confirm Supabase JWT issuer/audience configuration and role/tenant authorization model.
3. Confirm whether PATCH's endpoint-specific JSON update schemas are sufficient, including explicit-null behavior.
4. Confirm `limit`/`offset` defaults, max limit, and whether `total` is affordable for each collection.
5. Confirm which routes are async, need rate limiting, or require idempotency keys.
6. Confirm which OpenAPI linting, documentation, and Go code-generation tools the team will adopt. -->

---

**References:** [OpenAPI 3.1](https://spec.openapis.org/oas/v3.1.0.html) · [RFC 9457 Problem Details](https://www.rfc-editor.org/rfc/rfc9457) · [RFC 9110 HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110) · [RFC 3339 Date and Time](https://www.rfc-editor.org/rfc/rfc3339)
