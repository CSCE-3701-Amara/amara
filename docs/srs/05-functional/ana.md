# 5.4 ANA: Candidate Analytics & Evaluation Subsystem

**Owner:** Ali Mohamed | **Status:** Draft | **Last updated:** 2026-10-09 | **Jira:** AMARA-33

---

## 1. Overview

The Analytics Subsystem is concerned with receiving real-time telemetry data, code execution metrics, and behavioural statistics to produce technical and behavioural insights, logs, and scorecards for the candidate to be displayed on our platform Amara.

---

## 2. Scope

* **In-Scope:**
  * IDE telemetry ingestion and real-time event streaming.
  * Sandbox execution metrics (runtime, memory, O-complexity approximation) and test pass rate scoring.
  * Multi-source behavioral scoring (peer reviews, PR metrics, task-board activity).
  * Integrity risk scoring.
  * Aggregated PDF/JSON candidate evaluation report generation.
  * Demographic score variance auditing to detect algorithmic bias.
* **Out-of-Scope:**
  * Direct execution sandbox provisioning and virtualization infrastructure.
  * Candidate scheduling and communications.
  * The nature of tests taken by the candidate.

---

## 3. Features and Functional Requirements

### 3.1 Telemetry Ingestion & Real-Time Streaming

**Feature Description:** Ingestion and low-latency processing of candidate interaction events directly within the IDE environment to capture raw behavioral signals without affecting user experience.

#### Preconditions:

* Candidate active session in the web or desktop IDE environment.
* Telemetry streaming endpoint is established.

#### Functional Requirements

| ID | Requirement | Priority | Target Milestone | Related Use Case | Jira |
| --- | --- | --- | --- | --- | --- |
| FR-ANA-01 | The subsystem shall stream and ingest telemetry events (keystrokes, pastes, focus events, cursor moves) at sub-second intervals. | Must Have | M2 | | AMARA-33 |

---

### 3.2 Technical Execution & Code Analysis

**Feature Description:** Evaluation of code submissions within the sandbox, generating performance benchmarks and test outcome scores.

#### Preconditions:

* Code submitted to the sandbox environment.
* Automated tests and benchmarking are developed.

#### Functional Requirements

| ID | Requirement | Priority | Target Milestone | Related Use Case | Jira |
| --- | --- | --- | --- | --- | --- |
| FR-ANA-02 | The subsystem shall calculate individual execution efficiency metrics (runtime, peak memory, O-complexity approximation) upon code submission to the sandbox. | Must Have | M2 | | AMARA-33 |
| FR-ANA-03 | The subsystem shall calculate a Pass Rate Score based on the tests passed. | Must Have | M2 | | AMARA-33 |

---

### 3.3 Behavioral & Integrity Scoring

**Feature Description:** Computation of candidate behavioral profiles from collaborative interactions and detection of proctoring anomalies using heuristic thresholds.

#### Preconditions:

* Integration active with peer review, pull request, and task-board data pipelines.
* Telemetry threshold rules configured for integrity detection.

#### Functional Requirements

| ID | Requirement | Priority | Target Milestone | Related Use Case | Jira |
| --- | --- | --- | --- | --- | --- |
| FR-ANA-04 | The subsystem shall compute a normalized Behavioral Score based on peer review surveys, PR interactions, and task-board allocation metrics. | Should Have | M3 | | AMARA-33 |
| FR-ANA-05 | The subsystem shall generate an Integrity Risk Score (0-100) and attach evidence logs (paste size, focus loss timestamps) whenever threshold heuristics are exceeded. | Should Have | M4 | | AMARA-33 |

---

### 3.4 Candidate Reporting

**Feature Description:** Consolidation of candidate evaluation metrics into unified reports and automated monitoring of evaluation logs for demographic bias.

#### Preconditions:

* Technical, behavioral, and integrity calculations completed for candidates.

#### Functional Requirements

| ID | Requirement | Priority | Target Milestone | Related Use Case | Jira |
| --- | --- | --- | --- | --- | --- |
| FR-ANA-07 | The subsystem shall aggregate technical, behavioral, and integrity metrics into a single exportable candidate report (PDF/JSON). | Could Have | M4 | | AMARA-33 |

---

## 4. Business Rules

* **BR-ANA-01:** Candidates with an Integrity Risk Score exceeding 70 must be flagged for manual review before final report generation.

---

## 5. Error Handling

* **ERR-ANA-01:** Telemetry connection drops must cause local event queuing in the IDE buffer without interrupting developer code input.
* **ERR-ANA-02:** Sandbox timeouts or memory exceptions during code submission must record partial diagnostic telemetry and return execution failure flags.
* **ERR-ANA-03:** Report generation failures during PDF/JSON rendering must trigger up to 3 automatic retries.
