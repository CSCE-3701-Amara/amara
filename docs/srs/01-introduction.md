## 1. Introduction

Owner: TBD | Status: Draft | Last updated: 2026-10-14 | Jira: TBD

### 1.1 Purpose

This Software Requirements Specification (SRS) defines what the AMARA platform is intended to do, the constraints under which it operates, and the quality requirements against which it will be reviewed and tested. It is the shared requirements baseline for product stakeholders, designers, developers, testers, and project reviewers. The SRS describes required behavior; it does not prescribe implementation details that belong in the Software Design Document (SDD).

### 1.2 Product scope

AMARA is a web-based recruitment platform focused on Computer Science-related hiring. It connects Candidates and employers through job discovery and applications, configurable screening and interview pipelines, assessments, interviews, feedback, hiring decisions and offers. It also supports candidate profiles, assessment-verified skills and Amara crests, dynamic resumes, company identity and trust, and hiring analytics.

The system is organized around five functional epics:

| Epic | Responsibility |
|---|---|
| AUTH — Accounts, Identity & Trust | Registration, activation, authentication, role-based access, company verification, profiles, and administration. |
| CAND — Candidates & Dynamic Resume | Candidate-facing discovery and application experiences, application views, verified skills, crests, and dynamic resumes. |
| JOB — Jobs, Search & Integrations | Job-posting lifecycle and configuration, application submission, Must-have screening, and admission to the interview workflow. |
| INVW — Pipeline & Interviews & Tests | Admitted-application progression, assessment and interview execution, scoring, feedback, application decisions, finalists, offers, and notifications. |
| ANA — Candidate Analytics & Evaluation | Assessment telemetry and evaluation metrics, integrity and behavioral analysis, reports, and bias auditing. |

Functional ownership boundaries are authoritative: for example, JOB admits applications to the pipeline, INVW owns status transitions after admission and their notifications, AUTH owns access control, and ANA reads source data without writing to other epics' records. The detailed boundaries are stated in the functional requirements and use cases.

### 1.3 Product goals and expected outcome

The product is intended to help Candidates discover relevant opportunities and understand their application progress, while helping employers define repeatable hiring workflows and review evidence from applications and assessments. Employers configure pipeline stages and advancement policies; Candidates complete the required steps and receive permitted feedback and decisions. Hiring decisions remain associated with the employer's configured process. Analytics and AI-assisted capabilities support specified tasks but do not replace the scoring or decision rules defined by the responsible epic.

The expected project outcome is a requirements baseline that is sufficiently specific to design, implement, and verify the stated workflows. Each functional requirement is traceable to a use case; non-functional requirements provide measurable quality targets; assumptions, dependencies, global considerations, and risks record important context.

### 1.4 Intended audience

- Product stakeholders and teaching staff reviewing scope and acceptance criteria.
- Epic owners and developers implementing AUTH, CAND, JOB, INVW, and ANA.
- UX and accessibility contributors designing candidate, interviewer, hiring-team, and administrative interactions.
- Testers deriving functional, integration, security, usability, and performance checks.
- Operators reviewing deployment, recovery, privacy, and capacity expectations.

### 1.5 Document organization

| Section | Content |
|---|---|
| `02-overall-description.md` | Product context, users, high-level functions, environment, assumptions, and dependencies. |
| `03-constraints.md` | Constraints, assumptions, and dependencies. |
| `04-use-cases.md` | Actor-oriented workflows and alternatives/exceptions. |
| `05-functional/` | Functional requirements, one file per epic. |
| `06-nfr.md` | Measurable security, privacy, capacity, compatibility, reliability, and other quality requirements. |
| `07-global-considerations.md` | Global, cultural, social, environmental, economic, health, safety, and welfare considerations. |
| `08-risk-fmea.md` | Project risk register and preliminary failure-mode analysis. |

The SDD in `docs/sdd/` explains the proposed architecture, modules, component interactions, and user-facing interfaces that realize these requirements.

### 1.6 Conventions and references

Requirement and use-case identifiers follow `docs/conventions/requirements.md`. **SHALL** denotes a mandatory requirement; **SHOULD** and **COULD** retain the priority meanings used in the NFR section. Terms such as Candidate, Hiring Manager, posting, pipeline, stage, Amara, and Amarat are used consistently with `appendix-a-glossary.md` and the functional epics.
