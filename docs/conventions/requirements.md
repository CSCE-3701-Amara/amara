# Requirements and User Stories

Owner: Ziad Eliwa | Status: Draft | Last updated: 2026-10-10 | Jira: AMARA-59

## 1. Identifiers

| Type | Format | Example |
|---|---|---|
| User story | `US-<EPIC>-nnn` | `US-INVW-002` |
| Functional requirement | `FR-<EPIC>-nnn` | `FR-JOB-002` |
| Non-functional requirement | `NFR-<AREA>-nnn` | `NFR-SEC-001` |
| Use case | `UC-<EPIC>-nnn` | `UC-JOB-001` |
| Constraint | `CON-nnn` | `CON-004` |
| Assumption | `ASM-nnn` | `ASM-001` |
| Dependency | `DEP-nnn` | `DEP-001` |
| External interface | `EXT-nnn` | `EXT-003` |
| Data / database | `DB-nnn` | `DB-005` |
| Legal / regulatory | `LEG-nnn` | `LEG-002` |
| Internationalization | `I18N-nnn` | `I18N-001` |
| Risk / FMEA row | `RISK-nnn` | `RISK-007` |

Epic prefixes: `AUTH`, `CAND`, `JOB`, `INVW`, `ANA`. NFR areas used in `docs/srs/07-nfr.md`: `SEC`, `PRIV`, `CAP`, `COMP`, `REL`, `MAINT`, `PORT`, `SCAL`, `USAB`, `ACC`, `PERF`, `ENV`. Availability is part of `REL`. There is no `AVAIL` prefix in the SRS.

Every identifier uses three digits, as in `FR-JOB-001` and `CON-001`. Do not mix widths.

The SRS does not assign a separate test-case ID. A required check is written in the requirement that needs it, usually in `docs/srs/07-nfr.md`: an automated test, an integration test that fails on a stated condition, a load test whose result is stored, or a moderated test. Do not add a `TC-` ID until the SRS does.

Rules:
- IDs are assigned by the epic owner, in order, and **never reused or renumbered**.
- A removed requirement stays in the document marked `Withdrawn`, with the reason.
- The ID goes at the start of the Jira title: `US-INVW-002: Apply with saved profile`.
- Several functional requirements may cite one use case. `FR-JOB-001` through `FR-JOB-012` all cite `UC-JOB-001`. A use case number is not required to match a user-story number. The filled SRS has no `US-` IDs.

## 2. User stories

Format:

> As a **[role]**, I want **[capability]**, so that **[benefit]**.

- The role is an actor from the glossary, not "user".
- One capability per story. If it takes more than 8 points, split it.
- Acceptance criteria in Given/When/Then, at least one for the normal path and one for the main failure path:

```
Given a complete profile and an open posting,
When I click Apply,
Then the application is created with status "Submitted" and I receive a confirmation.
```

- A story passes the INVEST check: Independent, Negotiable, Valuable, Estimable, Small, Testable.

## 3. Functional requirements

- Wording in functional requirements and NFRs: "The system *SHALL* ...", one testable statement per requirement. This is the wording in `docs/srs/05-functional/invw.md` and `docs/srs/07-nfr.md`. Constraints use "The platform shall" and are not rewritten here.
- Describe the need, not the design. Write "store user data persistently", not "use PostgreSQL". Technology goes in constraints or the design specification.
- A functional requirement in a filled epic file cites a related use case. It does not cite a `US-` ID. Do not invent a story ID to fill a trace.
- Refer to another requirement by ID. Do not restate it.
- Epic functional-requirement tables use: ID | Requirement | Related Use Case.
- NFR tables use: ID | Requirement | Priority | Status | Source story | Jira. The Source story cell is an FR ID, another NFR ID, `cross-cutting`, or `TBD (<EPIC>)`. It is not a `US-` ID.

## 4. Non-functional requirements

- Must be **measurable**: a number, a condition, or a standard.
- Bad: "The system shall be fast." Good: the wording of `NFR-PERF-003` in `docs/srs/07-nfr.md`.
- Each important NFR names the component or design decision it drives (the **Drives:** line in `docs/srs/07-nfr.md`).
- NFRs are tracked in Jira as Story or Task with the label `nfr`.

## 5. Terms

- Use only terms from the glossary (`docs/srs/appendix-a-glossary.md`).
- Status values (for example, application and posting statuses) are defined once, in the glossary and the matching state machine. Nobody invents another list.
- A term is capitalized and spelled the same everywhere.

## 6. Priority (MoSCoW)

| MoSCoW | Meaning | Jira priority |
|---|---|---|
| Shall | Required for the system to be useful, in the MVP | Highest |
| Should | Important, in scope if time allows | High |
| Could | Nice to have | Medium |
| Won't | Out of this release, kept for future expansion | Low |

Functional requirements in `docs/srs/05-functional/invw.md` have no Priority column. NFR priority is Shall, Should, or Could. In `docs/srs/07-nfr.md`, Shall is the T1 bar, Should is T2, and Could is T3.

Stories planned for full implementation also carry the label `mvp`.

## 7. Lifecycle

| Status | Meaning | Who sets it |
|---|---|---|
| Draft | Written by the owner | Owner |
| Reviewed | A teammate checked it against the checklist below | Reviewer |
| Approved | The Product Owner accepted it | Product Owner |
| Withdrawn | No longer valid, kept for reference | Owner, with reason |

Changing an Approved requirement means a pull request that updates the text, the status, and the change history, plus a comment on the Jira item explaining why.

## 8. Review checklist

- [ ] ID follows the format and is unique
- [ ] Wording is "*SHALL*" for a functional requirement or NFR, one statement, testable
- [ ] No ambiguous words ("fast", "easy", "user-friendly", "etc.") without a measure
- [ ] Describes the need, not the implementation
- [ ] Terms match the glossary
- [ ] Priority and status set
- [ ] A functional requirement is linked to a use case. An NFR names its source requirement, `cross-cutting`, or `TBD (<EPIC>)`. Do not invent a story ID.
- [ ] No conflict with another epic, especially at boundaries: statuses, roles, notifications, permissions

## 9. Boundary ownership

| Topic | Owned by | Others |
|---|---|---|
| Application status values and transitions | INVW | Reference by ID |
| Notifications | INVW | Reference by ID |
| Roles and permission matrix | AUTH | Reference by ID |
| Reading data for analytics | ANA | Reads only, never writes other epics' data |

## 10. Global, cultural, social, and welfare considerations

Every epic owner checks their stories against these dimensions: global, cultural, social, environmental, economic, public health, safety, public welfare. Any issue found becomes a requirement with an ID and is listed in `docs/srs/11-global-considerations.md`.