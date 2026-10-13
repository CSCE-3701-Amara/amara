# Software Requirements Specification (SRS)

Owner: Ziad Eliwa | Status: Draft | Last updated: 2026-10-10 | Jira: AMARA-59

The SRS holds the product backlog, use cases, and all requirements. It lives in `docs/srs/`, one Markdown file per section, so five people can edit at the same time without conflicts.

## 1. Layout

```
docs/srs/
  00-front-matter.md
  01-introduction.md
  02-overall-description.md
  03-constraints.md
  04-use-cases.md
  05-functional/
    auth.md  cand.md  job.md  invw.md  ana.md
  06-external-interfaces.md
  07-nfr.md
  08-data.md
  09-legal.md
  10-i18n.md
  11-global-considerations.md
  12-ai-requirements.md
  13-ai-workflow.md
  14-risk-fmea.md
  15-scope-feasibility.md
  appendix-a-glossary.md
  appendix-b-traceability.md
```

## 2. File header

Most files start with a section heading, then one owner line:

```
## <n>. <Title>

Owner: <name> | Status: Draft | Last updated: <YYYY-MM-DD> | Jira: AMARA-<n>
```

- Status values used in the SRS: Draft, In Review. Reviewed and Approved remain the later states in [requirements.md](requirements.md).
- One owner per file. Others suggest changes through a pull request comment or a Jira comment.
- Optional sections are marked `(Opt)` in the heading until someone fills them in.
- Filled files use `AMARA-<n>`. Some unfilled files still say `Jira: HIRE-n`. New and updated files use `AMARA`.
- `05-functional/invw.md` uses a `#` heading, bold field labels, and Status `In Review`.
- `00-front-matter.md` puts the owner line after the change history.
- `13-ai-workflow.md`, `14-risk-fmea.md`, and `15-scope-feasibility.md` currently have only the owner line.

## 3. What each file contains

| File | Content | Rubric link |
|---|---|---|
| 00 | Title, team (as on Canvas), change history, table of contents | Item 7 |
| 01 | Purpose, overview and motivation, scope (in and out), stakeholders | Criterion 1 |
| 02 | Product perspective, actors, operating environment, epic overview | Criterion 1 |
| 03 | Constraints, assumptions, dependencies | Criterion 4 |
| 04 | Actors, then per-epic use cases. The filled cases are `UC-INVW-nnn` tables (preconditions, trigger, main flow, alternative and exception flows, postconditions, business rules). Other epic headings are empty. | Item 4 |
| 05 | One file per epic. The filled file, `app.md`, has overview, scope, features, functional requirements (ID, requirement, related use case), business rules, and error handling. The other epic files are headings only. | Criterion 2 |
| 06 | AI services, APIs, integrations, user interfaces, error handling | Criterion 2 |
| 07 | Non-functional requirements, each measurable | Criterion 2 |
| 08 | Domain model, data to be stored, database requirements | Criterion 2 |
| 09 | Laws, regulations, and the requirements they create | Criterion 2 |
| 10 | Languages, right-to-left, date, number, currency, time zone, names | Criterion 2 |
| 11 | Global, cultural, social, environmental, economic, health, safety, welfare considerations | Item 8 |
| 12 | AI transparency and fairness. The file heading currently says `11. AI Transperancy & Fairness Requirements`. | — |
| 13 | AI workflow. Not filled beyond the owner line. | — |
| 14 | Risk register and FMEA matrix | Item 5, criterion 4 |
| 15 | Features to fully implement (MVP), deferred features, feasibility to December 1 | Item 6, criterion 4 |
| A | Glossary | All |
| B | Traceability matrix (Opt) | Item 7 |

## 4. Table columns

Use the columns below. Functional requirements and NFRs do not share one table.

- **Functional requirements:** ID | Requirement | Related Use Case
- **Non-functional requirements:** ID | Requirement | Priority | Status | Source story | Jira
- **Interfaces:** ID | Service | Purpose | Inputs / outputs | Fallback if unavailable | Core or stretch
- **Legal:** ID | Law or regulation | What it requires | System requirement | Verified how
- **Global considerations:** Dimension | Issue | Resulting requirement IDs | Mitigation
- **Risk register:** ID | Risk | Category | Likelihood | Impact | Mitigation | Owner | Status
- **FMEA:** ID | Failure mode | Effect | Cause | S | O | D | RPN | Mitigation | Owner | Status
- **Traceability:** a functional requirement cites a use case; an NFR cites source requirement IDs, `cross-cutting`, or `TBD (<EPIC>)`. `appendix-b-traceability.md` is optional and not filled. The SRS does not use a separate test-case ID. A required test is stated in the requirement.

**FMEA scales.** S (severity), O (occurrence), D (detection) are scored 1 to 10, where 10 is the worst (most severe, most frequent, hardest to detect). RPN = S x O x D. Any RPN above the team threshold (TBD, for example 100) needs a mitigation and an owner.

## 5. Writing rules

- Follow [requirements.md](requirements.md) for IDs, wording, and priority.
- Filename numbers are the section numbers. Do not renumber without updating every file and reference. A heading that does not match its filename is an inconsistency in that file, not a new number. Current mismatches: `12-ai-requirements.md` is headed `11.`; `job.md` and `ana.md` are both headed `5.4`; `cand.md` is headed `5.5`.
- Figures: caption `Figure n: <title>`, image from `docs/diagrams/img/`, source committed in `docs/diagrams/src/`, and referenced from the text.
- Link to other sections by file path and ID, not by copying text.

## 6. Versions and change tracking

- `0.x` until the Milestone 1 presentation, `1.0` at Milestone 1, then `1.1`, `1.2`, and so on, one increase per sprint.
- Every pull request that changes the SRS adds a row to the change history in `00-front-matter.md`: Version | Date | Author | Sections | What changed | Jira.
- A pull request without the row is not approved.
- The tracking is graded in Milestones 2 to 4, so keep it complete from the start.

## 7. Updating each sprint

At the end of each sprint, the epic owners update their functional files (new, changed, withdrawn requirements), and the SRS owner updates scope and risks. The sprint task "Update requirements and design specifications with tracking of changes" is not Done until this is finished.

## 8. Review

Before an SRS section goes to Reviewed, run the checklist in [requirements.md](requirements.md) and the consistency check in [documentation.md](documentation.md).

## 9. Export

One PDF is generated from the files in order (tooling TBD, for example Pandoc). One person owns the export and tests it before every milestone, not the night before.