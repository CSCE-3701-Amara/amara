## 5.1 AUTH: Accounts, Identity & Trust

Owner: Hazem Nasr | Status: Draft | Last updated: 2026-10-09 | Jira: HIRE-n

### 5.1.1 Overview
This section defines the functional requirements related to **AMARA-6 AUTH: Accounts, Identity & Trust**.

### 5.1.2 Scope
This subsystem offers the following functionality:
- Registration
- Verification
- Login
- Logout
- Profile Management
- System Administration

### 5.1.3 Requirements

#### 5.1.3.1 General Authentication Requirements
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-001 | The system **SHALL** create the account when all validations succeed. | UC-AUTH-001, UC-AUTH-002, UC-AUTH-013 | TBD |
| FR-AUTH-002 | The system **SHALL** create an authenticated session upon successful login. | UC-AUTH-003 | TBD |
| FR-AUTH-003 | The system **SHALL** display an appropriate error message when authentication, validation, or verification fail. | UC-AUTH-001, UC-AUTH-002, UC-AUTH-003, UC-AUTH-004, UC-AUTH-005, UC-AUTH-006, UC-AUTH-007, UC-AUTH-008, UC-AUTH-009, UC-AUTH-010, UC-AUTH-011, UC-AUTH-012, UC-AUTH-013, UC-AUTH-014, UC-AUTH-015, UC-AUTH-016, UC-AUTH-017 | TBD |
| FR-AUTH-004 | The system **SHALL** implement a Role-Based Access Control (RBAC) system to avoid unauthorized access to resources and services. | UC-AUTH-003 to UC-AUTH-017 | TBD |
| FR-AUTH-005 | The system **SHALL** validate passwords according to the defined password policy. | UC-AUTH-001, UC-AUTH-010, UC-AUTH-013 | TBD |
| FR-AUTH-006 | The system **SHALL** allow all users, including interviewers, to log out at any time, upon which their authenticated session **SHALL** be terminated. | UC-AUTH-011 | TBD |
| FR-AUTH-007 | The system **SHOULD** allow all user accounts to set up MFA. | UC-AUTH-001, UC-AUTH-003 | TBD |
| FR-AUTH-008 | The system **SHOULD** allow users to log in using MFA if it is configured for their account. | UC-AUTH-003 | TBD |

#### 5.1.3.2 Candidate and Interviewer Registration
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-009 | The system **SHALL** allow users to create candidate or interviewer accounts using a valid email address and password; interviewer accounts **SHALL** include an association with the company for which the interviewer works. | UC-AUTH-001 | TBD |
| FR-AUTH-010 | The system **SHALL** verify that the email address is not already registered for candidate or interviewer registration. | UC-AUTH-001 | TBD |


#### 5.1.3.3 Company Registration
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-011 | The system **SHALL** allow users to submit requests to create new company accounts using information about the company and the person who is requesting. | UC-AUTH-002 | TBD |
| FR-AUTH-012 | The system **SHALL** force company accounts to set up MFA. | UC-AUTH-002, UC-AUTH-003 | TBD |

#### 5.1.3.4 Company Verification
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-013 | The system **SHALL** automatically review company account creation requests against official registries to verify the companies are real. | UC-AUTH-002 | TBD |
| FR-AUTH-014 | The system **SHOULD** allow system admins to manually review company account creation requests and accept/reject them. | UC-AUTH-002, UC-AUTH-015 | TBD |
| FR-AUTH-015 | The system **SHOULD** verify the identity of the person requesting to create a company account to avoid impersonation. | UC-AUTH-002 | TBD |

#### 5.1.3.5 Candidate and Interviewer Login
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-016 | The system **SHALL** allow candidates and interviewers to log in using their direct credentials (email address and password), with interviewer sessions scoped to their associated company. | UC-AUTH-003 | TBD |
| FR-AUTH-017 | The system **SHOULD** allow candidates to log in with their Google, GitHub, or LinkedIn accounts. | UC-AUTH-003 | TBD |
| FR-AUTH-018 | The system **SHOULD** allow users to verify their identity through email if they forgot their password. | UC-AUTH-010 | TBD |

#### 5.1.3.6 Company Login
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-019 | The system **SHALL** allow company representatives to log in only using MFA. | UC-AUTH-003 | TBD |

#### 5.1.3.7 Reporting and Rating
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-020 | The system **SHOULD** allow candidates to report fake companies to system admins. | UC-AUTH-004 | TBD |
| FR-AUTH-021 | The system **SHOULD** allow companies to report fake candidates to system admins. | UC-AUTH-004 | TBD |
| FR-AUTH-022 | The system **COULD** allow candidates to rate the companies being tested by. | UC-AUTH-009 | TBD |
| FR-AUTH-023 | The system **COULD** allow companies to rate the candidates they test. | UC-AUTH-009 | TBD |

#### 5.1.3.8 Profile Management
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-024 | The system **SHALL** allow candidates to view and edit their profiles at any time. | UC-AUTH-005 | TBD |
| FR-AUTH-025 | The system **SHALL** allow candidates to deactivate or delete their profiles at any time. | UC-AUTH-006 | TBD |
| FR-AUTH-026 | The system **SHOULD** allow candidates to share their profiles via links/QR codes. | UC-AUTH-007 | TBD |
| FR-AUTH-027 | The system **SHOULD** allow users to reset their passwords. | UC-AUTH-010, UC-AUTH-013 | TBD |
| FR-AUTH-028 | The system **SHALL** provide a public company profile page that shows the information related to the company and their job postings. | UC-AUTH-012 | TBD |
| FR-AUTH-029 | The system **COULD** allow candidates to set job preferences that the system can use to recommend job positions. | UC-AUTH-005 | TBD |
 
#### 5.1.3.9 Account Activation
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-030 | The system **SHALL** support account activation via email for all users. | UC-AUTH-001, UC-AUTH-002 | TBD |
| FR-AUTH-031 | The system **SHALL** not allow non-activated accounts to login. | UC-AUTH-003 | TBD |

#### 5.1.3.10 Document Uploads
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-032 | The system **SHALL** allow all users to upload supplementary documents in PDF and DOCX formats with a maximum of 10 MB per document. | UC-AUTH-008 | TBD |
| FR-AUTH-033 | The system **COULD** automatically extract profile information from uploaded documents and populate the related profile fields only after the user reviews and confirms the extracted information. | UC-AUTH-008 | TBD |

#### 5.1.3.11 System Administration
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-034 | The system **SHALL** provide a user management interface for admins. | UC-AUTH-013 | TBD |
| FR-AUTH-035 | The system **SHALL** allow admins to add, edit, or delete all entities on the system, including candidates, companies, job positions, and tests. | UC-AUTH-014 | TBD |
| FR-AUTH-036 | The system **SHALL** allow admins to view, ban, or deactivate user accounts and reset account credentials. | UC-AUTH-013 | TBD |
| FR-AUTH-037 | The system **SHALL** allow admins to view and accept/reject company registration requests. | UC-AUTH-015 | TBD |
| FR-AUTH-038 | The system **SHALL** allow admins to manage system taxonomies, platform constants, and reference values. | UC-AUTH-016 | TBD |
| FR-AUTH-039 | The system **SHOULD** provide statistics and reporting features for admins to obtain aggregate data related to job postings. This data can be used to enhance the platform or for other research purposes. | UC-AUTH-017 | TBD |

#### 5.1.3.12 Interviewer-Specific Requirements
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-040 | The system **SHALL** validate and record the interviewer's association with an approved company before activating the interviewer account. | UC-AUTH-001 | TBD |
| FR-AUTH-041 | The system **SHALL** require interviewers to set up MFA and use MFA verification on every login. | UC-AUTH-001, UC-AUTH-003 | TBD |
