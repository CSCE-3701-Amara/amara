# Milestone TBDs and Open Items

Owner: TBD | Status: Draft | Last updated: 2026-10-14 | Jira: TBD

This is a working register of important requirements, design, and verification work that is not yet present or finalized in the current documentation set. **Open** means that an owner, decision, or deliverable is still needed; it does not imply that a particular implementation has been selected. The register records known gaps and should be updated as items are assigned and resolved.

## 1. SRS documentation gaps

| ID | Open item | Needed deliverable / completion condition | Status |
|---|---|---|---|
| TBD-001 | SRS front matter | Add the SRS title, team details, table of contents, and maintained change history. The current `docs/srs/` tree has no front-matter file. | Open |
| TBD-002 | Domain and data model | Create a conceptual model/ER diagram covering Companies, users and memberships, postings, applications/submissions, pipeline versions/stages, assessment instances/results, Amara definitions/awards, offers, and analytics records. State cardinalities, ownership, tenant scope, and lifecycle. | Open |
| TBD-003 | Data dictionary and retention map | Define stored field meanings, sensitivity, ownership, required/optional rules, identifiers, retention, and deletion behavior for domain rows, resumes, video, telemetry, and derived reports. | Open |
| TBD-004 | External-interface requirements | Specify the referenced external-interface material, including provider contracts, payloads, authentication, timeouts, retries, idempotency, and failure behavior for email, registry, social login, meeting, code evaluation, AI, storage, and queue interactions. | Open |
| TBD-005 | Legal and regulatory mapping | Map applicable Egypt privacy and accessibility obligations to exact system requirements, retention periods, user rights, and verification evidence. Existing NFR references are not a substitute for a reviewed legal mapping. | Open |
| TBD-006 | Localization specification | Define supported locales, translation ownership, Arabic RTL behavior, date/number/currency formatting, fallback behavior, and localized status/error content in a dedicated specification. | Open |
| TBD-007 | AI transparency and workflow | Consolidate AI use cases, input/output data, user disclosures, human review points, provider constraints, and evaluation criteria for question drafts, feedback, insights, and analytics. | Open |
| TBD-009 | Use-case-to-test traceability | Decide whether to add the optional traceability appendix and, if used, map requirements/use cases to planned verification evidence. | Open |
| TBD-010 | Requirement approval and ownership | Confirm owners, review/approval status, unresolved placeholders, and cross-epic agreements for the existing requirements. | Open |

## 2. SDD and architecture gaps

| ID | Open item | Needed deliverable / completion condition | Status |
|---|---|---|---|
| TBD-011 | Physical deployment diagram | Add a deployment view showing browser delivery, API replicas, worker replicas, queue, relational database, object storage, external providers, trust boundaries, and network paths. The current diagrams are logical epic/component and user-flow views, not a physical deployment plan. | Open |
| TBD-012 | Cloud environment plan | Select AWS or Azure and a production region; document account/subscription structure, network ingress/egress, compute, database, queue, object storage, email, secrets, TLS, backups, monitoring, residency, and cost assumptions. | Open |
| TBD-013 | Infrastructure-as-code and environment topology | Define local, CI, staging, and production resources, configuration/secrets management, promotion/rollback, and deployment ownership in reproducible infrastructure code. | Open |
| TBD-014 | Technology decisions / ADRs | Record approved choices for the database engine/schema strategy, queue, object storage, hosting, API/session implementation, and provider adapters, including alternatives and cost/operational tradeoffs. | Open |
| TBD-015 | Physical data design | Convert the conceptual model into schema/migration plans, tenant keys, uniqueness and concurrency constraints, indexes, object metadata, retention jobs, and backup/restore boundaries. | Open |
| TBD-016 | Workflow sequence diagrams | Add sequence/activity views for registration/MFA, application screening/admission, pipeline publication/migration, assessment scoring, anonymity/reveal, finalist/offer/expiry, notification retry, and provider failure. | Open |
| TBD-017 | API and event contracts | Define versioned request/response schemas and internal/event contracts, including `IAccessControl`, `IPostingData`, `IPipelineData`, `IAmaraAwards`, `IOffersAndDecisions`, and `ICandidateScorecard`. | Open |
| TBD-018 | UI wireframes and design system | Define navigation, page states, form validation, empty/error/offline states, responsive layouts, Arabic RTL, and accessible interaction patterns for Candidate, Hiring Manager, Interviewer, and Admin surfaces. | Open |
| TBD-019 | Security threat model and privacy data-flow review | Complete threat models and data-flow reviews for company isolation, Candidate identity/anonymity, resumes/media, telemetry, AI processors, and privileged admin actions; assign mitigation and test owners. | Open |
| TBD-020 | Operations and observability | Specify dashboards, alerts, queue-lag/provider-failure signals, incident severity/response, log/trace redaction, retention, restore exercises, and operational runbooks. | Open |

## 3. Verification and rollout gaps

| ID | Open item | Needed deliverable / completion condition | Status |
|---|---|---|---|
| TBD-021 | Functional and integration test plan | Map use cases and cross-epic handoffs to automated tests, including authorization/company scope, application admission, pipeline progression, offers, and notification failure. | Open |
| TBD-022 | Capacity and performance evidence | Run and retain load-test results for the T1 application, concurrency, pipeline, notification, storage, and latency targets in `docs/srs/06-nfr.md`. | Open |
| TBD-023 | Recovery evidence | Execute and record database/object-storage restore exercises against the defined RPO/RTO targets. | Open |
| TBD-024 | Accessibility and usability evidence | Conduct the specified browser, keyboard, screen-reader, Candidate, and Hiring Manager usability checks; record findings and remediation. | Open |
| TBD-025 | Fairness and analytics validation | Define validation data, review method, subgroup/variance thresholds, and human-review procedure for ANA outputs before they influence hiring workflows. | Open |
| TBD-026 | Cloud cost validation | Validate T1 budget assumptions against selected providers and document alerts, ownership, and actions for projected overspend. | Open |
| TBD-027 | JOB finalist target | Decide whether JOB has a configurable finalist target, define who sets it and how it is used, and add its functional requirement/use case. INVW references this value in FR-INVW-067 and FR-INVW-068, but no JOB functional requirement currently defines it. | Open |
| TBD-028 | Cross-epic notification ownership | Reconcile JOB use cases that send posting/application/following notifications with INVW's notification ownership. Define which epic emits each event, which module chooses recipients/channels, and which requirements specify timing/retries; separately confirm AUTH account-activation and password-reset notification ownership. | Open |
| TBD-029 | Privacy-rights workflows | Assign an epic owner and functional use cases for privacy notice/consent, Candidate data export, erasure requests, and breach-notice workflows currently expressed in NFR privacy requirements without matching functional requirements/use cases. | Open |
| TBD-030 | Identifier-width policy | Reconcile the three-digit ID convention with existing two-digit JOB/ANA use-case IDs and ANA functional IDs. Preserve current references until an approved migration plan updates all requirements, use cases, diagrams, and trace links together. | Open |
| TBD-031 | Posting deletion and retention | Define how JOB's permanent posting deletion behavior interacts with NFR privacy retention, legal holds, and decision records for postings with applications; specify whether any minimal decision record survives deletion. | Open |
| TBD-032 | Assessment telemetry notice and lifecycle | Specify what IDE/assessment telemetry is collected, how Candidates are informed and consented where required, what data is minimized before ANA processing, who may access it, and how long raw events and derived scores are retained. Trace the decisions to ANA functional requirements and privacy controls. | Open |

## 4. Deferred product decisions

- Confirm which optional capabilities are included in the current release, such as social login, AI-assisted generation/feedback, recommendations, and long-term multi-person or cheating-detection use cases.
- Confirm whether regional deployments beyond Egypt are in scope after T1; no production multi-region topology is specified yet.
- Confirm any commercial/pricing model if it is in the product scope. The NFR cloud-spend limit is an engineering budget constraint, not a user pricing decision.

## 5. Maintenance

Assign each open item an owner and target milestone when it is scheduled. When completed, add the resulting document or decision, update SRS/SDD references, and mark or remove the item here so this remains a list of genuinely unresolved work.
