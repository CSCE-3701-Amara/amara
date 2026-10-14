## 9. Scope and feasibility

Owner: TBD | Status: Draft | Last updated: 2026-10-14 | Jira: TBD

### 9.1 Implementation scope through December 1, 2026

By December 1, 2026, the team will implement and verify every functional requirement (FR) in `docs/srs/05-functional/*.md` whose normative wording contains `SHALL`, including prohibitions phrased as `SHALL NOT`. This includes mandatory requirements in AUTH, CAND, JOB, INVW, and ANA. No SHALL functional requirement is excluded based on epic or feature area.

Functional requirements marked `SHOULD` or `COULD` are not part of the committed December 1 scope unless the project explicitly reprioritizes them. Non-functional requirements remain separately defined in `docs/srs/06-nfr.md` and are not reclassified by this functional scope rule.

### 9.2 Feasibility basis

Feasibility is planned against that complete set of SHALL functional requirements. The project will sequence implementation across the five epics and their dependencies; it will not silently reduce the SHALL scope. Any blocker or proposed scope change must be raised for an explicit requirements decision and reflected in the functional requirements and project tracking.

The architecture baseline for estimating and sequencing this work is the modular monolith with worker processes documented in `docs/sdd/03-system-architecture.md`. Detailed effort, staffing, and delivery estimates are not provided in this document.
