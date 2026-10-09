## 3. Constraints, assumptions, dependencies

Owner: Ali Mohammed | Status: Draft | Last updated: 2026-10-08 | Jira: HIRE-n

### 3.1 Constraints

CON-001 — Project Scope: The platform shall focus on hiring processes for Computer Science-related roles.

CON-002 — SaaS Delivery Model: The platform shall be delivered as a Software-as-a-Service (SaaS) solution.

CON-003 — Hiring Pipeline Configuration: The platform shall support configurable hiring pipelines in which stages and progression rules can be defined by employers.

CON-004 — Candidate Anonymity: The platform shall support anonymous candidate evaluation at designated hiring stages, according to the configured information-disclosure rules.

CON-005 — Assessment-Based Verification: The platform shall award verified skills and Crests based on eligible assessment outcomes rather than solely on candidate-provided claims.

CON-006 — Professional Document Formats: The platform shall restrict candidate document uploads to PDF and DOCX formats.

CON-007 — Maximum Upload Size: The platform shall limit each uploaded professional document to a maximum size of 10 MB.

CON-008 — Candidate Communication: The platform shall not support direct communication between candidates within the current project scope.

CON-009 — Access Control: The platform shall restrict access to information and functionality according to user roles and permissions.

CON-010 — Privacy During Anonymous Evaluation: The platform shall prevent unauthorized disclosure of candidate-identifying information during hiring stages configured for anonymous evaluation.

---

### 3.2 Assumptions

ASM-001 — Internet Connectivity: Candidates and employers are assumed to have internet access when using the platform.

ASM-002 — Device and Browser Compatibility: Users are assumed to access the platform through devices and browsers that support its required functionality.

ASM-003 — Candidate Information Accuracy: Candidates are assumed to provide accurate and up-to-date professional information, including skills, education, and work experience.

ASM-004 — Document Validity: Candidates are assumed to upload readable and legitimate professional documents in the supported formats.

ASM-005 — Employer Information Accuracy: Employers are assumed to provide accurate job descriptions, position requirements, compensation information when applicable, and hiring criteria.

ASM-006 — Employer Pipeline Configuration: Employers are assumed to configure their hiring pipelines, assessment requirements, progression rules, and identity-disclosure settings before using them to evaluate candidates.

ASM-007 — Assessment Definition: Assessment content, evaluation criteria, and passing conditions are assumed to be defined before assessments are made available to candidates.

ASM-008 — Assessment Relevance: Assessments are assumed to measure the skills or competencies they are intended to evaluate.

ASM-009 — External Profile Availability: Candidates who provide GitHub, LinkedIn, or other professional profile links are assumed to maintain valid links. The platform cannot guarantee the availability or accuracy of information hosted on external platforms.

ASM-010 — Employer Decision Responsibility: Employers are assumed to review candidate information and make hiring decisions according to their own hiring criteria.

ASM-011 — Verification Criteria: The criteria for earning each Crest or obtaining a verified skill are assumed to be defined and consistently applied by the platform.

ASM-012 — Employer Participation: Employers are assumed to maintain their job postings and record application status changes and hiring decisions in a timely manner.

ASM-013 — Notification Contact Information: Users are assumed to provide valid contact information if the platform uses that information to deliver notifications.

---

### 3.3 Dependencies

DEP-001 — Authentication and Account Management: The platform depends on account creation, authentication, and account management functionality to identify users and provide access to the appropriate features.

DEP-002 — Job Posting Management: Job discovery and application functionality depend on employers creating and maintaining job postings and specifying their availability.

DEP-003 — Hiring Pipeline Configuration: Application progression depends on the hiring stages, assessment requirements, and progression rules configured by employers.

DEP-004 — Assessment Management and Evaluation: Assessment completion, results, feedback, and assessment-dependent progression depend on available assessment content and evaluation functionality.

DEP-005 — Candidate Profile Data: Dynamic resume generation and candidate profile presentation depend on the availability of candidate-provided professional information, uploaded documents, and external profile links.

DEP-006 — Verified Skill Records: The generation and maintenance of verified skills and Crests depend on assessment outcomes and the corresponding verification criteria.

DEP-007 — Application Status Updates: Application tracking and hiring decision notifications depend on application statuses and hiring decisions being recorded and updated.

DEP-008 — Document Storage and Retrieval: Document upload, viewing, and removal depend on functionality for storing, retrieving, and deleting uploaded files.

DEP-009 — Notification Delivery: Candidate notifications depend on an available notification mechanism capable of delivering updates about application progress, assessment availability, assessment results, and final hiring decisions.

DEP-010 — Analytics Data Availability: Employer analytics depend on the availability and accuracy of application, pipeline, assessment, and hiring outcome data collected by the platform.

DEP-011 — Access Control and Anonymity: Anonymous evaluation depends on user roles, access permissions, and mechanisms that restrict the visibility of candidate-identifying information at designated hiring stages.

DEP-012 — External Professional Platforms: The usefulness of external professional profile links depends on the continued availability of the corresponding third-party platforms. Direct integrations with these platforms are not assumed unless explicitly included in the project requirements.

DEP-013 — Assessment and Hiring Workflow Integration: The platform depends on the correct exchange of information between assessment evaluation, application status tracking, pipeline progression, and candidate notification functionality.