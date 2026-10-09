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
- System Administraion

### 5.1.3 Requirements

#### 5.1.3.1 General Authentication Requirements
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-001 | The system **SHALL** create the account when all validations succeed. | UC-AUTH-01, UC-AUTH-02, UC-AUTH-16 | TBD |
| FR-AUTH-002 | The system **SHALL** create an authenticated session upon successful login. | UC-AUTH-03, UC-AUTH-04 | TDB |
| FR-AUTH-003 | The system **SHALL** display an appropriate error message when authentation, validation, or verification fail. | UC-AUTH-01 to UC-AUTH-20 | TBD |
| FR-AUTH-004 | The system **SHALL** implement a Role-Based Access Control (RBAC) system to avoid unauthorized access to resources and services. | UC-AUTH-03 to UC-AUTH-14, UC-AUTH-16 to UC-AUTH-20 | TBD |
| FR-AUTH-005 | The system **SHALL** validate passwords according to the defined password policy. | UC-AUTH-01, UC-AUTH-13, UC-AUTH-16 | TBD |
| FR-AUTH-006 | The system **SHALL** allow all users to logout at any time. | UC-AUTH-14 | TBD |
| FR-AUTH-007 | The system **SHALL** terminate the user's authenticated session upon logout. | UC-AUTH-14 | TBD |

#### 5.1.3.2 Candidate Registration
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-008 | The system **SHALL** allow usres to create candidate accounts using a valid email address and password. | UC-AUTH-01 | TBD |
| FR-AUTH-009 | The system **SHALL** verify that the email address is not already registered. | UC-AUTH-01 | TBD |
| FR-AUTH-010 | The system **SHOULD** allow candidate accounts to setup MFA. | UC-AUTH-01, UC-AUTH-03 | TBD |


#### 5.1.3.3 Company Registration
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-011 | The system **SHALL** allow users to submit requests to create new company accounts using information about the company and the person who is requesting. | UC-AUTH-02 | TBD |
| FR-AUTH-012 | The system **SHALL** force company accounts to setup MFA. | UC-AUTH-02, UC-AUTH-04 | TBD |

#### 5.1.3.4 Company Verification
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-013 | The system **SHALL** automatically review company account creation requests against official registries to verify the companies are real. | UC-AUTH-02 | TBD |
| FR-AUTH-014 | The system **SHOULD** allow system admins to manually review company account creation requests and accept/reject them. | UC-AUTH-02, UC-AUTH-18 | TBD |
| FR-AUTH-015 | The system **SHOULD** verify the identity of the person requesting to create a company account to avoid impersonation. | UC-AUTH-02 | TBD |

#### 5.1.3.5 Candidate Login
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-016 | The system **SHALL** allow candidates to login using their direct credentials (email address and password). | UC-AUTH-03 | TBD |
| FR-AUTH-017 | The system **SHOULD** allow candidates to login with their Google, GitHub, or LinkedIn accounts. | UC-AUTH-03 | TBD |
| FR-AUTH-018 | The system **SHOULD** allow candidates to login using MFA if setup. | UC-AUTH-03 | TBD |
| FR-AUTH-019 | The system **SHOULD** allow users to verify their identity through email if they forgot their password. | UC-AUTH-13 | TBD |

#### 5.1.3.6 Company Login
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-020 | The system **SHALL** allow companies to login only using MFA. | UC-AUTH-04 | TBD |

#### 5.1.3.7 Reporting and Rating
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-021 | The system **SHOULD** allow candidates to report fake companies to system admins. | UC-AUTH-05 | TBD |
| FR-AUTH-022 | The system **SHOULD** allow companies to report fake candidates to system admins. | UC-AUTH-06 | TBD |
| FR-AUTH-023 | The system **COULD** allow candidates to rate the companies being tested by. | UC-AUTH-11 | TBD |
| FR-AUTH-024 | The system **COULD** allow companies to rate the candidates they test. | UC-AUTH-12 | TBD |

#### 5.1.3.8 Profile Management
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-025 | The system **SHALL** allow candidates to view and edit their profiles at any time. | UC-AUTH-07 | TBD |
| FR-AUTH-026 | The system **SHALL** allow candidates to deactivate or delete their profiles at any time. | UC-AUTH-08 | TBD |
| FR-AUTH-027 | The system **SHOULD** allow candidates to share their profiles via links/QR codes. | UC-AUTH-09 | TBD |
| FR-AUTH-028 | The system **SHOULD** allow users to reset their passwords. | UC-AUTH-13, UC-AUTH-16 | TBD |
| FR-AUTH-029 | The system **SHALL** provide a public company profile page that shows the information related to the company and their job postings. | UC-AUTH-15 | TBD |
| FR-AUTH-030 | The system **COULD** allow candidates to set job preferneces that the system can use to recommend job positions. | UC-AUTH-07 | TBD |
 
#### 5.1.3.9 Account Activation
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-031 | The system **SHALL** support account activation via email for all users. | UC-AUTH-01, UC-AUTH-02, UC-AUTH-03, UC-AUTH-04 | TBD |
| FR-AUTH-032 | The system **SHALL** not allow non-activated accounts to login. | UC-AUTH-03, UC-AUTH-04 | TBD |

#### 5.1.3.10 Document Uploads
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-033 | The system **SHALL** allow all users to upload supplementary documents in PDF and DOCX formats. | UC-AUTH-10 | TBD |
| FR-AUTH-034 | The system **COULD** automatically extract profile information from the uploaded documents and fill the related profile fields accordingly. | UC-AUTH-10 | TBD |

#### 5.1.3.11 System Administration
|     ID      | Requirement | Related Use Case(s) | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-035 | The system **SHALL** provide a user management interface for admins. | UC-AUTH-16 | TBD |
| FR-AUTH-036 | The system **SHALL** allow admins to add, edit, or delete all entities on the system, including candidates, companies, job positions, and tests. | UC-AUTH-17 | TBD |
| FR-AUTH-037 | The system **SHALL** allow admins to view, ban, or deactivate user accounts and reset account credentials. | UC-AUTH-16 | TBD |
| FR-AUTH-038 | The system **SHALL** allow admins to view and accept/reject company registration requests. | UC-AUTH-18 | TBD |
| FR-AUTH-039 | The system **SHALL** allow admins to manage system taxonomies, platform constants, and reference values. | UC-AUTH-19 | TBD |
| FR-AUTH-040 | The system **SHOULD** provide statistics and reporting features for the admins to get aggregate data related to job posting. This data can be used to enhance the platform or for other research purposes. | UC-AUTH-20 | TBD |
