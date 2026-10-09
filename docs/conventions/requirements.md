# Requirements and User Stories

Owner: Ziad Eliwa | Status: Draft | Last updated: 2026-10-08 | Jira: AMARA-59

## 1. Identifiers

| Type | Format | Example |
|---|---|---|
| User story | `US-<EPIC>-nn` | `US-APP-02` |
| Functional requirement | `FR-<EPIC>-nn` | `FR-APP-02` |
| Non-functional requirement | `NFR-<AREA>-nn` | `NFR-SEC-01` |
| Use case | `UC-<EPIC>-nn` (same number as its story) | `UC-APP-02` |
| External interface | `EXT-nn` | `EXT-03` |
| Data / database | `DB-nn` | `DB-05` |
| Legal / regulatory | `LEG-nn` | `LEG-02` |
| Internationalization | `I18N-nn` | `I18N-01` |
| Risk / FMEA row | `RISK-nn` | `RISK-07` |
| Test case | `TC-<EPIC>-nn` | `TC-APP-04` |

Epic prefixes: `AUTH`, `CAND`, `JOB`, `APP`, `ANA`. NFR areas: `PERF`, `SEC`, `PRIV`, `USAB`, `ACC`, `REL`, `AVAIL`, `MAINT`, `PORT`, `SCAL`, `CAP`, `COMP`, `ENV`.

Rules:
- IDs are assigned by the epic owner, in order, and **never reused or renumbered**.
- A removed requirement stays in the document marked `Withdrawn`, with the reason.
- The ID goes at the start of the Jira title: `US-APP-02: Apply with saved profile`.

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

- Wording: "The system shall ...", one testable statement per requirement.
- Describe the need, not the design. Write "store user data persistently", not "use PostgreSQL". Technology goes in constraints or the design specification.
- Every functional requirement traces to at least one story. Every story has at least one requirement.
- Refer to another requirement by ID. Do not restate it.

## 4. Non-functional requirements

- Must be **measurable**: a number, a condition, or a standard.
- Bad: "The system shall be fast." Good: "Search results shall load within 2 seconds for 95% of requests under 500 concurrent users."
- Each important NFR names the component or design decision it drives (filled in the SRS and the design specification).
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
- [ ] Wording is "shall", one statement, testable
- [ ] No ambiguous words ("fast", "easy", "user-friendly", "etc.") without a measure
- [ ] Describes the need, not the implementation
- [ ] Terms match the glossary
- [ ] Priority and status set
- [ ] Linked to a story, a Jira key, and (if it applies) an NFR or use case
- [ ] No conflict with another epic, especially at boundaries: statuses, roles, notifications, permissions

## 9. Boundary ownership

| Topic | Owned by | Others |
|---|---|---|
| Application status values and transitions | APP | Reference by ID |
| Notifications | APP | Reference by ID |
| Roles and permission matrix | AUTH | Reference by ID |
| Reading data for analytics | ANA | Reads only, never writes other epics' data |

## 10. Global, cultural, social, and welfare considerations

Every epic owner checks their stories against these dimensions: global, cultural, social, environmental, economic, public health, safety, public welfare. Any issue found becomes a requirement with an ID and is listed in `docs/srs/11-global-considerations.md`.