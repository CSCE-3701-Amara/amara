## 7. Non-functional requirements

Owner: Hazem Nasr | Status: Draft | Last updated: 2026-10-09 | Jira: AMARA-17

### 7.1 Overview
This section defines the nonfunctional requirements related to the system as a whole in addition to those specific to certain subsystems.

### 7.2 Performance Requirements
#### 7.2.1 Response time
- The system SHALL provide page load times of less than 3 seconds for standard operations under normal load conditions.
- The system shall spawn all kinds of tests in less than 20 seconds under normal load condition.
- The system SHALL maintain response time degradation of no more than 50% during peak load periods.

#### 7.2.2 Throughput
- TODO: 

#### 7.2.3 Resource utilization
- The system SHALL operate within the allocated server resources, utilizing no more than 80% of CPU capacity during normal operations.
- The system SHALL utilize no more than 80% of available memory during normal operations.
- TODO: storage utilization
- TODO: cashing

#### 7.2.4 Scalability
- The system SHALL be designed to scale horizontally by adding more server instances to handle increased load.
- The system SHALL be designed to scale vertically by utilizing additional resources on existing servers.  
- TODO: 


### 7.3 Security Requirements
#### 7.3.1 Authentication and Authorization
- The system SHALL enforce password requirements in accordance with NIST SP 800-63B, including minimum password length and rejection of commonly used or compromised passwords. The system SHALL NOT require periodic password changes unless there is evidence of compromise.
- The system SHALL maintain detailed access logs for all authentication and authorization events. 
- TODO: 

#### 7.3.2 Data Protection
- The system SHALL encrypt all sensitive data at rest using industry-standard encryption algorithms (AES-256 or equivalent). 
- The system SHALL encrypt all data in transit using TLS 1.3 or higher.
- The system SHALL maintain separate environments for development, testing, and production with appropriate data isolation.
- TODO: practices for encryption keys

#### 7.3.3 Privacy and Compliance
- The system SHALL comply with the General Data Protection Regulation (GDPR).
- TODO: by Yousef Abood

#### 7.3.4 Security Monitoring and Incident Response
- TODO:


### 7.4 Reliability and Availability Requirements
#### 7.4.1 Availability
- The system SHALL maintain 99.9% availability, in other words, all the services of system need to have a maximum downtime of 90s per day.
- The system SHALL schedule maintenance windows during periods of lowest expected usage.
- The system SHALL provide advance notice of scheduled maintenance to all users.

#### Fault Tolerance 
- The system SHALL implement database replication to prevent data loss in case of database failures. 
- The system SHALL implement load balancing across multiple servers to distribute traffic and prevent overload. 
- The system SHALL automatically recover from common failure scenarios without manual intervention.

#### Disaster Recovery
- The system SHALL maintain regular backups of all data, with full backups at least weekly and incremental backups daily.
- The system SHALL define and document Recovery Time Objective (RTO) of 4 hours for critical functions and 24 hours for non-critical functions.
- The system SHALL have a documented and tested disaster recovery plan.
- The system SHALL conduct disaster recovery drills at least twice per year.

#### Error Handling 
- The system SHALL provide meaningful error messages to users without exposing sensitive system information.
- The system SHALL log detailed error information for troubleshooting and monitoring.
- TODO:


### 7.5 Usability and Accessibility Requirements
#### 7.5.1 User Interface
- TODO: 

#### Accessibility
- The system SHALL comply with Web Content Accessibility Guidelines (WCAG) 2.1 Level AA standards. 
- The system SHOULD support screen readers.
- The system SHOULD provide keyboard navigation for all functions.
- The system SHOULD provide text alternatives for non-text content.

#### User Experience
- The system SHALL collect and incorporate user feedback for continuous improvement.
- TODO:

### 7.6 Maintainability and Portability Requirements
#### 7.6.1 Maintainability
- The system SHALL implement logging and monitoring to facilitate troubleshooting. 
- The system SHALL implement automated testing with a minimum of 80% code coverage.
- The system SHALL support version control for all system artifacts.

#### Portability
- The system SHALL use containerization technologies to ensure consistent deployment across environments.
- The system SHALL support automated deployment and configuration.

#### Compatibility
- The system SHALL be compatible with the latest versions of major web browsers (Chrome, Firefox, Safari, Edge).
- The system SHALL be compatible with the previous two major versions of supported browsers.
- The system SHALL be compatible with mobile browsers on iOS and Android platforms.
- The system SHALL be compatible with standard email clients for notification delivery. 

### 7.7 Legal and Compliance Requirements
TODO: by Yousef Abood
#### Regulatory Compliance

#### Intellectual Property

#### Service Level Agreements 


### 7.8 Operational Requirements
#### Monitoring and Logging
- The system SHALL provide real-time monitoring of system health and performance.
- The system SHALL generate alerts for critical system events and performance thresholds.

#### Backup and Recovery
- The system SHALL perform automated backups according to defined schedules.
- The system SHALL verify backup integrity through automated testing.
- The system SHALL provide mechanisms for point-in-time recovery. 
- The system SHALL document and test restoration procedures. 
- The system SHALL maintain backup history and audit trails.

#### Documentation
- The system SHALL provide comprehensive user documentation for all user roles.
- The system SHALL provide technical documentation for system
administrators and developers.
- The system SHALL maintain up-to-date system architecture and design documentation.
- The system SHALL provide API documentation for integration partners.
- The system SHALL document all configuration parameters and their effects.
- The system SHALL provide troubleshooting guides and known issue documentation.