## 5.1 AUTH: Accounts, Identity & Trust

Owner: Hazem Nasr | Status: Draft | Last updated: 2026-10-08 | Jira: HIRE-n

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
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-001 | The system shall create the account when all validations succeed. | TBD | TBD |
| FR-AUTH-002 | The system shall create an authenticated session upon successful login. | TBD | TDB |
| FR-AUTH-003 | The system shall display an appropriate error message when authentation, validation, or verification fail. | TBD | TBD |
| FR-AUTH-004 | The system shall implement a Role-Based Access Control (RBAC) system to avoid unauthorized access to resources and services. | TBD | TBD |
| FR-AUTH-005 | The system shall validate passwords according to the defined password policy. | TBD | TBD |
| FR-AUTH-006 | The system shall allow all users to logout at any time. | TBD | TBD |
| FR-AUTH-007 | The system shall terminate the user's authenticated session upon logout. | TBD | TBD |

#### 5.1.3.2 Candidate Registration
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-008 | The system shall allow usres to create candidate accounts using a valid email address and password. | TBD | TBD |
| FR-AUTH-009 | The system shall verify that the email address is not already registered. | TBD | TBD |
| FR-AUTH-010 | The system should allow candidate accounts to setup MFA. | TBD | TBD |


#### 5.1.3.3 Company Registration
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-011 | The system shall allow users to submit requests to create new company accounts using information about the company and the person who is requesting. | TBD | TBD |
| FR-AUTH-012 | The system shall force company accounts to setup MFA. | TBD | TBD |

#### 5.1.3.4 Company Verification
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-013 | The system shall automatically review company account creation requests against official registries to verify the companies are real. | TBD | TBD |
| FR-AUTH-014 | The system should allow system admins to manually review company account creation requests and accept/reject them. | TBD | TBD |
| FR-AUTH-015 | The system should verify the identity of the person requesting to create a company account to avoid impersonation. | TBD | TBD |

#### 5.1.3.5 Candidate Login
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-016 | The system shall allow candidates to login using their direct credentials (email address and password). | TBD | TBD |
| FR-AUTH-017 | The system should allow candidates to login with their Google, GitHub, or LinkedIn accounts. | TBD | TBD |
| FR-AUTH-018 | The system should allow candidates to login using MFA if setup. | TBD | TBD |
| FR-AUTH-019 | The system should allow users to verify their identity through email if they forgot their password. | TBD | TBD | 

#### 5.1.3.6 Company Login
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-020 | The system shall allow companies to login only using MFA. | TBD | TBD |

#### 5.1.3.7 Reporting and Rating
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-021 | The system should allow candidates to report fake companies to system admins. | TBD | TBD |
| FR-AUTH-022 | The system should allow companies to report fake candidates to system admins. | TBD | TBD |
| FR-AUTH-023 | The system could allow candidates to rate the companies being tested by. | TBD | TBD |
| FR-AUTH-024 | The system could allow companies to rate the candidates they test. | TBD | TBD |

#### 5.1.3.8 Profile Management
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-025 | The system shall allow candidates to view and edit their profiles at any time. | TBD | TBD |
| FR-AUTH-026 | The system shall allow candidates to deactivate or delete their profiles at any time. | TBD | TBD |
| FR-AUTH-027 | The system should allow candidates to share their profiles via links/QR codes. | TBD | TBD |
| FR-AUTH-028 | The system should allow users to reset their passwords. | TBD | TBD |
| FR-AUTH-029 | The system shall provide a public company profile page that shows the information related to the company and their job postings. | TBD | TBD |
| FR-AUTH-030 | The system could allow candidates to set job preferneces that the system can use to recommend job positions. | TBD | TBD |
 
#### 5.1.3.9 Account Activation
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-031 | The system shall support account activation via email for all users. | TBD | TBD |
| FR-AUTH-032 | The system shall not allow non-activated accounts to login. | TBD | TBD |

#### 5.1.3.10 Document Uploads
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-033 | The system shall allow all users to upload supplementary documents in PDF and DOCX formats. | TBD | TBD |
| FR-AUTH-034 | The system could automatically extract profile information from the uploaded documents and fill the related profile fields accordingly. | TBD | TBD |

#### 5.1.3.11 System Administration
|     ID      | Requirement | Related Use Case | Jira |
|-------------|-------------|------------------|------|
| FR-AUTH-035 | The system shall provide a user management interface for admins. | TBD | TBD |
| FR-AUTH-036 | The system shall allow admins to add, edit, or delete all entities on the system, including candidates, companies, job positions, and tests. | TBD | TBD |
| FR-AUTH-037 | The system shall allow admins to view, ban, or deactivate user accounts and reset account credentials. | TBD | TBD |
| FR-AUTH-038 | The system shall allow admins to view and accept/reject company registration requests. | TBD | TBD |
| FR-AUTH-039 | The system shall allow admins to manage system taxonomies, platform constants, and reference values. | TBD | TBD |
| FR-AUTH-040 | The system should provide statistics and reporting features for the admins to get aggregate data related to job posting. This data can be used to enhance the platform or for other research purposes. | TBD | TBD |
