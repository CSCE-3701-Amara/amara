## 4. Use cases

Owner: TBD | Status: Draft | Last updated: 2026-10-08 | Jira: HIRE-n
---

## 4.1 Actors

// This is only an example.
- Candidate
- Employer
- Administrator
- System

## 4.2 AUTH Use Cases

## 4.3 CAND Use Cases

## 4.4 APP Use Cases

## 4.5 JOB Use Cases

## 4.6 ANA Use Cases
## 4. Use cases

Owner: TBD | Status: Draft | Last updated: 2026-10-08 | Jira: HIRE-n
---

## 4.1 Actors

// This is only an example.
- Candidate
- Employer
- Administrator
- System

## 4.2 AUTH Use Cases

## 4.3 CAND Use Cases

## 4.4 APP Use Cases

## 4.5 JOB Use Cases

## 4.6 ANA Use Cases
### Use Case 1: Stream & Ingest IDE Telemetry

| Field | Description |
| --- | --- |
| **Use Case ID** | UC-ANA-01 |
| **Use Case Name** | Stream & Ingest IDE Telemetry |
| **Actor(s)** | Candidate, IDE |
| **Description** |Stream and analyze telemetry data coming from the IDE that the candidate is using. |
| **Preconditions** | Candidate active session in the IDE environment; telemetry endpoint is established. |
| **Trigger** | An action is triggered through the IDE. |
| **Main Flow** | 1. The IDE captures interaction events in real-time.<br>2. Events are buffered and streamed to the telemetry endpoint. |
| **Alternative Flows** | None. |
| **Exception Flows** | Telemetry connection drops -> IDE plugin queues events locally in buffer until connection is restored. |
| **Postconditions** | Telemetry events successfully connected to the candidate record. |
| **Business Rules** |The telemetry streaming must run flawlessly and not affect the candidate’s experience. |

---

### Use Case 2: Calculate Execution Efficiency Metrics

| Field | Description |
| --- | --- |
| **Use Case ID** | UC-ANA-02 |
| **Use Case Name** | Calculate Efficiency Metrics |
| **Actor(s)** | Sandbo, Analytics Subsystem |
| **Description** |Measures the performance related metrics such as Run Time, Memory Consumption, and complexity. |
| **Preconditions** | Code submitted to the sandbox environment. |
| **Trigger** | The candidate submits code solutions to the sandbox. |
| **Main Flow** | 1. Sandbox executes submitted code against test suites.<br>2. Analytics subsystem calculates execution runtime and peak memory usage.<br>3. Subsystem calculates O-complexity approximation.<br>4. Execution metrics are saved to the candidate record. |
| **Alternative Flows** | None. |
| **Exception Flows** | Execution timeout -> System records partial diagnostic telemetry and returns failure flags. |
| **Postconditions** |The metrics are saved to the candidate’s records. |
| **Business Rules** | Metrics must be computed in an isolated execution sandbox environment. |

---

### Use Case 3: Calculate Technical Pass Rate Score

| Field | Description |
| --- | --- |
| **Use Case ID** | UC-ANA-03 |
| **Use Case Name** | Calculate Technical Pass Rate Score |
| **Actor(s)** | Analytics Subsystem |
| **Description** | Compute candidate technical Pass Rate Score based on automated test suite outcomes. |
| **Preconditions** | Test suite execution completed in a sandbox environment. |
| **Trigger** | Code execution finishes in the sandbox. |
| **Main Flow** | 1. Subsystem retrieves passed and failed test case counts.<br>2. Computes the Pass Rate Score percentage.<br>3. Updates technical scorecard record. |
| **Alternative Flows** | None. |
| **Exception Flows** | Unhandled code runtime crash -> Pass Rate Score recorded as 0% with diagnostic failure details. |
| **Postconditions** | Technical Pass Rate Score updated on candidate evaluation scorecard. |
| **Business Rules** | A zero pass rate score is not an auto rejection |

---

### Use Case 4: Evaluate Integrity Risk & Log Anomalies

| Field | Description |
| --- | --- |
| **Use Case ID** | UC-ANA-04 |
| **Use Case Name** | Evaluate Integrity Risk |
| **Actor(s)** | Analytics Subsystem |
| **Description** | Analyzing the telemetry events, sequence, and frequency and generating a integrity score based on it|
| **Preconditions** | Telemetry streaming active and integrity threshold rules configured. |
| **Trigger** | Telemetry events exceed certain limits (e.g., massive paste size, focus loss). |
| **Main Flow** | 1. Subsystem detects violation in incoming stream.<br>2. Calculates/updates Integrity Risk Score (0–100).<br>3. Attaches evidence log. |
| **Alternative Flows** | None. |
| **Exception Flows** |Issues with the telemetry stream. |
| **Postconditions** | Evidence logs attached and candidate Integrity Risk Score updated. |
| **Business Rules** |A high integrity score requires human revision to avoid auto rejection. |

---

### Use Case 5: Export Candidate Evaluation Report

| Field | Description |
| --- | --- |
| **Use Case ID** | UC-ANA-05 |
| **Use Case Name** | Export Candidate Evaluation Report |
| **Actor(s)** | Recruiter, Hiring Manager, Analytics Subsystem |
| **Description** | Consolidate technical, behavioral, and integrity metrics into a unified exportable candidate report (PDF). |
| **Preconditions** | Technical, behavioral, and integrity score calculations completed. |
| **Trigger** | Recruiter or Hiring Manager requests report export. |
| **Main Flow** | 1. Subsystem aggregates technical pass rate, execution metrics, and integrity risk score.<br>2. Renders report in requested export format (PDF).<br>3. Delivers file for download or platform display. |
| **Alternative Flows** | Automated report generation triggered upon assessment completion. |
| **Exception Flows** | PDF Generation failure -> System retries up to 3 times before queuing error notification. |
| **Postconditions** | Exported PDF candidate evaluation report generated. |
| **Business Rules** |None |




