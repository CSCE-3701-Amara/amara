# Pipeline Studio

Pipeline Studio saves a directed acyclic graph as one transaction: stage definitions, module snapshots, positions, and connections. A failed save preserves the previous draft. Stage identifiers in the response replace the editor's temporary identifiers.

## API

All operations require a company session and ownership of the posting.

- `GET /v1/pipeline-versions/:versionId/graph`: stages and links. Existing linear pipelines are returned with their sequential connections.
- `PUT /v1/pipeline-versions/:versionId/graph`: `{stages: Stage[], links: {fromStageId, toStageId}[]}`. Includes `x` and `y` on stages. Returns the saved graph with server IDs.
- `POST /v1/pipeline-versions/:versionId/publish`: activates a saved, nonempty version.
- `POST /v1/jobs/:id/pipeline/draft`: copies the latest published version, including module snapshots and connections. Returns an existing unpublished draft if one already exists.

Saving rejects duplicate IDs, invalid references, cycles, duplicate connections, invalid stage configuration, and graphs exceeding 100 stages or 1,000 connections. Multi-stage graphs require explicit connections. Published and superseded versions are immutable. The ordered-stage API cannot overwrite a branching draft.

## Application behavior

Newly admitted applications enter the currently published version through the existing execution worker. Every root stage starts independently; joins require all incoming stages to pass. Existing applications remain pinned to their original version, preserving progress and assessment evidence. Publication does not migrate active applications.

## Editor workflow

Create a pipeline, add assessment stages, select modules and schedules, and connect stages. Save a draft or publish the current graph. Failed requests retain edits and display the server error. Published graphs are read only; create an editable copy to revise them.

Existing live-interview configuration is preserved. Publication does not migrate active applications.

## Validation

Domain tests cover malformed graphs and cycles. PostgreSQL integration tests use an isolated temporary schema and check round trips, ownership, publication locks, branching entry execution, copying, and preservation of application versions. Browser tests cover rejected saves, retained edits, reload, publication, and mobile layout.
