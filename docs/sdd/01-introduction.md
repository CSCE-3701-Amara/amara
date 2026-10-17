## 1. Introduction

Owner: TBD | Status: Draft | Last updated: 2026-10-14 | Jira: TBD

### 1.1 Purpose

This Software Design Document (SDD) describes the proposed architecture of AMARA and the rationale for its major design choices. It turns the requirements and workflows in `docs/srs/` into a component structure, runtime model, internal contracts, external integrations, and human-facing interaction surfaces. It is intended for implementers, reviewers, testers, and operators.

### 1.2 Scope

The design covers the five SRS epics—AUTH, CAND, JOB, INVW, and ANA—and the external services they use. The selected backend style is a **modular monolith with background workers**. The logical domain modules remain in one Go backend codebase; API and worker processes execute separate roles over shared module boundaries and persisted data. The browser client is a React/Next.js application.

This document describes a target design, not evidence that every component has been implemented or deployed. Details that the SRS leaves open—such as the specific relational database, queue product, cloud provider, and vendors—remain deployment decisions and must not be inferred from the diagrams.

### 1.3 Design objectives

- Keep epic ownership and data-write boundaries explicit while using one backend deployment unit at T1.
- Protect tenant and role boundaries through AUTH-owned access-control contracts.
- Keep external providers replaceable behind adapters and avoid blocking committed user work on slow provider calls.
- Support the throughput, localization, recovery, privacy, and accessibility targets in `docs/srs/06-nfr.md`.
- Keep the design operable within the project’s one-team, student-budget constraints while allowing worker and read capacity to grow independently.

### 1.4 Document organization

| Section | Content |
|---|---|
| `02-system-overview.md` | System context, users, major modules, and external systems. |
| `03-system-architecture.md` | Modular-monolith structure, runtime, patterns, data and worker flows, interfaces, and tradeoffs. |
| `04-component-diagram/` | One component/decision page for each SRS epic, with links to its PlantUML source and rendered figure. |
| `05-interfaces.md` | Human-facing interface surfaces and representative SRS use-case flows. |

### 1.5 Design conventions

The diagrams in `src/` are editable PlantUML sources; their corresponding rendered figures are in `img/`. Epic boundaries and interface names follow the SRS. Dotted arrows in an epic diagram denote logical dependencies or data flow, not necessarily network calls: interactions between modules in the monolith are in-process contracts unless explicitly identified as external or asynchronous.

Use-case references in this document point to the repository’s consolidated `docs/srs/04-use-cases.md`. A legacy request path of `docs/srs/03-use-cases.md` does not exist in the current repository.
