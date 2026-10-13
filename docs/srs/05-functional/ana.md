## 5.4 ANA: Candidate Analytics & Evaluation Subsystem

Owner: Ali Mohamed | Status: Final | Last updated: 2026-10-13 | Jira: AMARA-33

### 5.4.1 Overview
This section defines the functional requirements related to **AMARA-33 ANA: Candidate Analytics & Evaluation Subsystem**. The Analytics Subsystem is concerned with receiving real-time telemetry data, code execution metrics, and behavioural statistics to produce technical and behavioural insights, logs, and scorecards for the candidate to be displayed on our platform Amara.

### 5.4.2 Scope
This subsystem offers the following functionality:
- IDE telemetry ingestion and real-time event streaming.
- Sandbox execution metrics (runtime, memory, O-complexity approximation) and test pass rate scoring.
- Multi-source behavioral scoring (peer reviews, PR metrics, task-board activity).
- Integrity risk scoring.
- Aggregated PDF candidate evaluation report generation.
- Demographic score variance auditing to detect algorithmic bias.

The following functionality is out-of-scope:
- Direct execution sandbox provisioning and virtualization infrastructure.
- Candidate scheduling and communications.
- The nature of tests taken by the candidate.

### 5.4.3 Requirements

#### 5.4.3.1 Telemetry Ingestion & Real-Time Streaming

|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-ANA-01 | The subsystem **SHALL** stream and ingest telemetry events (keystrokes, pastes, focus events, cursor moves) at sub-second intervals. | TBD | AMARA-33 |

#### 5.4.3.2 Technical Execution & Code Analysis


|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-ANA-02 | The subsystem **SHALL** calculate individual execution efficiency metrics (runtime, peak memory, O-complexity approximation) upon code submission to the sandbox. | TBD | AMARA-33 |
| FR-ANA-03 | The subsystem **SHALL** calculate a Pass Rate Score based on the tests passed. | TBD | AMARA-33 |

#### 5.4.3.3 Behavioral & Integrity Scoring


|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-ANA-04 | The subsystem **SHOULD** compute a Behavioral Score based on peer review surveys, PR interactions, and task-board allocation metrics. | TBD | AMARA-33 |
| FR-ANA-05 | The subsystem **SHOULD** generate an Integrity Risk Score (0−100) and attach evidence logs (paste size, focus loss timestamps) whenever threshold heuristics are exceeded. | TBD | AMARA-33 |

#### 5.4.3.4 Candidate Reporting


|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-ANA-07 | The subsystem **COULD** aggregate technical, behavioral, and integrity metrics into a single exportable candidate report (PDF/JSON). | TBD | AMARA-33 |

#### 5.4.3.5 Business Rules
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| BR-ANA-01 | Candidates with an Integrity Risk Score exceeding 70 **MUST** be flagged for manual review before final report generation. | TBD | AMARA-33 |

#### 5.4.3.6 Error Handling
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| ERR-ANA-01 | Telemetry connection drops **MUST** cause local event queuing in the IDE buffer without interrupting developer code input. | TBD | AMARA-33 |
| ERR-ANA-02 | Sandbox timeouts or memory exceptions during code submission **MUST** record partial diagnostic telemetry and return execution failure flags. | TBD | AMARA-33 |
| ERR-ANA-03 | Report generation failures during PDF/JSON rendering **MUST** trigger up to 3 automatic retries. | TBD | AMARA-33 |
