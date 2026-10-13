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
| FR-ANA-01 | The system **SHALL** stream and ingest keystrokes at sub-second intervals. | UC-ANA-01 | AMARA-33 |
| FR-ANA-02 | The system **SHALL** stream and ingest pastes at sub-second intervals. | UC-ANA-01 | AMARA-33 |
| FR-ANA-03 | The system **SHALL** stream and ingest focus events at sub-second intervals. | UC-ANA-01 | AMARA-33 |

#### 5.4.3.2 Technical Execution & Code Analysis


|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-ANA-04 | The system **SHALL** calculate runtime upon code submission to the sandbox. | UC-ANA-02| AMARA-33 |
| FR-ANA-05 | The system **SHALL** calculate peak memory usage upon code submission to the sandbox. | UC-ANA-02| AMARA-33 |
| FR-ANA-06 | The system **SHALL** calculate O-complexity approximation upon code submission to the sandbox. | UC-ANA-02| AMARA-33 |
| FR-ANA-07 | The system **SHALL** calculate a Pass Rate Score based on the tests passed. | UC-ANA-03 | AMARA-33 |

#### 5.4.3.3 Behavioral & Integrity Scoring


|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-ANA-09 | The system **SHOULD** generate an Integrity Risk Score (0−100) and attach evidence logs (paste size, focus loss timestamps) whenever threshold is exceeded. | UC-ANA-04 | AMARA-33 |

#### 5.4.3.4 Candidate Reporting


|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-ANA-10 | The system **COULD** aggregate technical, behavioral, and integrity metrics into a single exportable candidate report (PDF/JSON). | UC-ANA-05 | AMARA-33 |

#### 5.4.3.5 Business Rules
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| BR-ANA-11 | Candidates with an Integrity Risk Score exceeding 70 **SHALL** be flagged for manual review before final report generation. | UC-ANA-04 | AMARA-33 |

#### 5.4.3.6 Error Handling
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| ERR-ANA-12 | Telemetry connection drops **SHALL** cause local event queuing in the IDE buffer without interrupting developer code input. | UC-ANA-01 | AMARA-33 |
| ERR-ANA-13 | Sandbox timeouts or memory exceptions during code submission **SHALL** record partial diagnostic telemetry and return execution failure flags. | UC-ANA-01 | AMARA-33 |
| ERR-ANA-14 | Report generation failures during PDF rendering **SHALL** trigger up to 3 automatic retries. | UC-ANA-05 | AMARA-33 |
