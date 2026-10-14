# AMARA Documentation Guide

This directory contains the product requirements, software design, diagrams, and documentation conventions for AMARA, a SaaS hiring platform for Computer Science-related roles. Use this page to find the right document and understand how the SRS and SDD fit together.

## Where to start

| If you need to… | Read… |
|---|---|
| Understand the product and its intended scope | [SRS introduction](srs/01-introduction.md) and [overall description](srs/02-overall-description.md) |
| Find a user workflow or actor goal | [Consolidated use cases](srs/04-use-cases.md) |
| Find a functional requirement | The matching epic file in [`srs/05-functional/`](srs/05-functional/) |
| Find quality targets such as security, privacy, capacity, reliability, or accessibility | [SRS non-functional requirements](srs/06-nfr.md) |
| Understand how the system is intended to be structured | [SDD system architecture](sdd/03-system-architecture.md) |
| Find a specific epic's components | The matching page in [`sdd/04-component-diagram/`](sdd/04-component-diagram/) |
| Find the design's editable diagrams or rendered figures | [`sdd/src/`](sdd/src/) and [`sdd/img/`](sdd/img/) |
| Find items that are not specified or completed yet | [Milestone TBDs and open items](srs/TBD.md) |
| Check writing, ID, or diagram conventions | [`conventions/`](conventions/) |

## SRS: what the system must do

The Software Requirements Specification is the product's requirements baseline. It defines the system's goals, actors, constraints, workflows, functional behavior, and measurable quality requirements.

| File | Use it for |
|---|---|
| [`srs/01-introduction.md`](srs/01-introduction.md) | SRS purpose, product scope, audience, and document organization. |
| [`srs/02-overall-description.md`](srs/02-overall-description.md) | Product context, user classes, major capabilities, operating environment, and epic ownership. |
| [`srs/03-constraints.md`](srs/03-constraints.md) | Constraints, assumptions, and dependencies. |
| [`srs/04-use-cases.md`](srs/04-use-cases.md) | Actor-based scenarios for AUTH, CAND, JOB, INVW, and ANA. This is the repository's consolidated use-case file. |
| [`srs/05-functional/`](srs/05-functional/) | Functional requirements organized by epic: AUTH, CAND, JOB, INVW, and ANA. |
| [`srs/06-nfr.md`](srs/06-nfr.md) | Non-functional requirements and measurable quality targets. |
| [`srs/07-global-considerations.md`](srs/07-global-considerations.md) | Global, cultural, social, environmental, economic, health, safety, and welfare considerations. |
| [`srs/08-risk-fmea.md`](srs/08-risk-fmea.md) | Risk register and preliminary failure-mode analysis. |
| [`srs/09-scope-feasibility.md`](srs/09-scope-feasibility.md) | Functional implementation scope through December 1, 2026, and its feasibility basis. Open items are tracked in [`srs/TBD.md`](srs/TBD.md). |
| [`srs/appendix-a-glossary.md`](srs/appendix-a-glossary.md) | Shared terminology and state/ownership definitions. |

Functional requirements should be read with their related use cases and the owning epic's business rules. AUTH owns identity and access control; JOB owns posting configuration, submissions, initial screening, and admission; INVW owns post-admission progression, application statuses, offers, and hiring-workflow notifications; CAND owns the Candidate-facing experience and dynamic resume; ANA owns analytical outputs and reads source-epic data.

## SDD: how the system is intended to be designed

The Software Design Document describes the proposed architecture and design choices that realize the SRS. It uses a modular monolith with worker processes; it is a target design, not evidence that each element has been implemented or deployed.

| File or directory | Use it for |
|---|---|
| [`sdd/01-introduction.md`](sdd/01-introduction.md) | SDD purpose, scope, and design objectives. |
| [`sdd/02-system-overview.md`](sdd/02-system-overview.md) | System context, modules, users, and external integrations. |
| [`sdd/03-system-architecture.md`](sdd/03-system-architecture.md) | Modular-monolith structure, workers, module contracts, data ownership, patterns, and tradeoffs. |
| [`sdd/04-component-diagram/`](sdd/04-component-diagram/) | Component responsibilities, interactions, and design decisions for each epic. |
| [`sdd/05-interfaces.md`](sdd/05-interfaces.md) | Human-facing interfaces and the related SRS use-case flows. |
| [`sdd/src/`](sdd/src/) | Editable PlantUML source diagrams. |
| [`sdd/img/`](sdd/img/) | Rendered diagrams embedded in the SDD pages. |

The component pages link to both the corresponding `.puml` source and rendered figure. When a source diagram changes, update its rendered figure as part of the same documentation change. The SDD interfaces section references [`srs/04-use-cases.md`](srs/04-use-cases.md); there is no `srs/03-use-cases.md` in the current layout.

## Conventions and open work

[`conventions/`](conventions/) contains project rules for SRS requirements, diagrams, documentation, code, testing, and development. [`srs/TBD.md`](srs/TBD.md) records specifications, architecture choices, plans, and verification evidence that are not yet written or finalized. Update that list as items are assigned and resolved; do not treat a TBD as an approved design decision.
