# Database Design

Status: Draft | Last updated: 2026-11-06

## Scope and sources

This document describes AMARA's relational model for accounts, candidate profiles, job applications, hiring pipelines, assessments, interviews, offers, and Amara awards. It is based on the [entity relationship diagram](erd_diagrams/main.mmd) and the SQL in [`supabase/migrations/`](../../../supabase/migrations/), through `20261105232902_auth_mfa.sql`.

The SQL definitions determine types, nullability, defaults, keys, and checks. The diagram provides the overall model but omits some constraints and contains relationship discrepancies, recorded below. This is a repository schema reference, not confirmation of the schema deployed to Supabase.

The dictionary combines the forward changes in the migration files: the original table definitions and subsequent additions. Duplicate November table-creation files have been removed, and all remaining migration files contain forward SQL only. See [migration history and limitations](#migration-history-and-limitations) for the repair and verification approach.

## Storage and schema organization

AMARA uses Supabase PostgreSQL. The [database ADR](../adr/0001-supabase-postgres.md) specifies direct backend access through pgx and the Supavisor session pooler. Application authentication uses AMARA's own account tables; these tables do not reference Supabase `auth.users`.

The ADR proposes epic schemas named `accounts`, `cand`, `job`, `invw`, and `ana`, with `auth` reserved for Supabase. The migrations currently use unqualified table names and contain no `CREATE SCHEMA` or `SET search_path` statements. Their actual namespace depends on the connection's search path, normally `public`; the proposed epic namespaces have not been implemented by these migrations. An `auth_` table prefix is part of a table name, not a separate schema.

| Domain | Tables | Responsibility |
| --- | --- | --- |
| Accounts and registration | `candidates`, `companies`, `interviewers`, `auth_email_verification_challenges`, `auth_company_registration_requests`, `auth_mfa_credentials` | Credentials, verification, company review, and MFA enrollment |
| Candidate profile | `candidate_education`, `candidate_work_experience`, `candidate_languages`, `candidate_documents`, `candidate_certifications`, `candidate_amarat` | Resume details, document references, and earned awards |
| Jobs and applications | `job_postings`, `applications` | Vacancies, submissions, screening outcomes, and final decisions |
| Pipeline definitions | `modules`, `pipelines`, `pipeline_versions`, `pipeline_stages`, `stage_transitions`, `assessment_questions`, `test_configurations` | Reusable assessment content and versioned hiring workflows |
| Execution and assessment | `pipeline_executions`, `stage_executions`, `candidate_stage_progress`, `assessment_instances`, `assessment_attempts`, `assessment_submissions`, `assessment_results` | Candidate progression, answers, grading, and disclosure |
| Interviews and hiring | `interviews`, `interviewer_assignments`, `interview_feedback`, `identity_reveal_records`, `offers` | Scheduling, scorecards, anonymity release, and offers |

There are 33 application tables in this model. These migrations define no analytics tables, views, materialized views, stored functions, or triggers.

## Reading the data dictionary

- **Required** means `NOT NULL` or membership in a primary key. All other columns accept SQL `NULL`.
- **PK**, **FK**, and **UQ** mean primary key, foreign key, and uniqueness constraint or unique index. A composite constraint applies to the combination of its columns.
- A default applies when an insert omits a column; it does not automatically update that column later. Required columns without defaults must be supplied by the writer.
- `timestamp` means PostgreSQL `timestamp without time zone`; `timestamptz` means `timestamp with time zone`. Both occur in the model. `interval` stores a duration.
- Unless listed otherwise, UUID primary keys have no database generation default. Generated UUIDs are explicitly identified below.
- Unless listed otherwise, foreign keys use PostgreSQL's default `NO ACTION` update/delete behavior. None of these migrations specifies cascading deletion.
- `CHECK` constraints permit SQL `NULL` unless a separate `NOT NULL` constraint prevents it. A nullable bounded score, for example, may remain unknown.

## Accounts and registration

### `candidates`

Stores a candidate account and the main profile fields.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `name` | text | Yes | Candidate name |
| `email` | text | Yes | UQ on `lower(email)` within this table |
| `password_hash`, `salt` | bytea | Yes | Application credential material |
| `created_at` | timestamp | Yes | No default |
| `updated_at`, `deleted_at`, `deactivated_at`, `last_login_at` | timestamp | No | Account lifecycle timestamps |
| `phone_number`, `github_profile`, `linkedin_profile`, `portfolio_link`, `profile_picture_url` | text | No | Contact and profile links |
| `gender` | text | No | Profile attribute |
| `birth_date` | date | No | Date of birth |
| `nationality`, `address`, `biography` | text | No | Profile details |
| `years_of_experience` | integer | No | No nonnegative check |
| `available`, `banned` | boolean | No | No defaults; `NULL` is distinct from false |
| `ban_reason` | text | No | No check tying it to `banned` |
| `account_status` | text | Yes | Default `pending_verification`; no allowed-value check |

### `companies`

Stores the employer account, organization profile, and approval flag. Hiring-manager references in the hiring workflow currently point to this table.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `email` | text | Yes | UQ on `lower(email)` within this table |
| `password_hash`, `salt` | bytea | Yes | Application credential material |
| `legal_name` | text | Yes | Legal organization name |
| `display_name`, `description`, `website_link`, `logo_url` | text | No | Organization profile |
| `contact_phone`, `contact_email`, `linkedin_profile` | text | No | Contact information |
| `created_at`, `updated_at`, `deleted_at`, `deactivated_at`, `last_login_at` | timestamp | No | No defaults |
| `banned` | boolean | No | No default |
| `ban_reason` | text | No | Ban explanation |
| `approved` | boolean | Yes | Default false |
| `registration_secret_hash` | bytea | No | Registration secret hash |

### `interviewers`

Stores interviewer accounts independently from company credentials.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `company_id` | uuid | No | FK → `companies.id`; company membership is nullable |
| `name` | text | Yes | Interviewer name |
| `email` | text | Yes | UQ on `lower(email)` within this table |
| `password_hash`, `salt` | bytea | Yes | Application credential material |
| `created_at` | timestamp | Yes | No default |
| `updated_at`, `deleted_at`, `deactivated_at`, `last_login_at` | timestamp | No | Account lifecycle timestamps |
| `account_status` | text | Yes | Default `pending_verification`; no allowed-value check |

### `auth_email_verification_challenges`

Stores expiring verification-code hashes for candidates and interviewers.

| Column | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK; default `gen_random_uuid()` |
| `account_id` | uuid | Yes | Polymorphic account identifier; no FK |
| `account_type` | text | Yes | `candidate` or `interviewer` |
| `code_hash` | bytea | Yes | Verification-code hash |
| `expires_at` | timestamp | Yes | Expiry timestamp |
| `consumed_at` | timestamp | No | Consumption timestamp |
| `created_at` | timestamp | Yes | Default `now()` |

An ordinary index covers `(account_id, account_type)`. Multiple challenges per account are allowed. The database does not enforce account existence, expiry ordering, or single-use consumption.

### `auth_company_registration_requests`

Stores company registration submissions before or during review. It has no FK to `companies` and no column linking an approved request to a created company.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK; default `gen_random_uuid()` |
| `company_name`, `registration_number`, `email` | text | Yes | Submitted company details |
| `password_hash`, `salt` | bytea | Yes | Submitted credential material |
| `status` | text | Yes | Default `pending_review`; no allowed-value check |
| `created_at` | timestamp | Yes | Default `now()` |
| `reviewed_at` | timestamp | No | Review timestamp |
| `review_reason` | text | No | Review explanation |

Separate partial unique indexes on `lower(email)` and `lower(registration_number)` apply only to rows whose status is `pending_review` or `approved`. Rows with other statuses do not participate in those uniqueness checks. These indexes do not enforce uniqueness against the `companies` table.

### `auth_mfa_credentials`

Stores one MFA credential per account/type pair.

| Column | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `account_id` | uuid | Yes | Composite PK with `account_type`; no FK |
| `account_type` | text | Yes | Composite PK; `candidate`, `interviewer`, or `company` |
| `secret` | bytea | Yes | MFA secret material; encryption is not specified by the DDL |
| `enabled` | boolean | Yes | Default false |
| `created_at` | timestamp | Yes | Default `now()` |
| `enabled_at` | timestamp | No | No check tying this to `enabled` |

## Candidate profiles and awards

### Profile detail tables

Each table below has a nullable `cand_id` FK to `candidates.id`. None has a primary key, unique constraint, generated row identifier, or timestamps beyond the fields listed. Duplicate rows and rows without a candidate are permitted by the SQL.

| Table | Column(s) | Type | Required | Meaning |
| --- | --- | --- | --- | --- |
| `candidate_education` | `cand_id` | uuid | No | Candidate FK |
| `candidate_education` | `degree_type`, `provider`, `summary` | text | No | Qualification and institution details |
| `candidate_education` | `start_date`, `end_date` | date | No | Study period |
| `candidate_work_experience` | `cand_id` | uuid | No | Candidate FK |
| `candidate_work_experience` | `position`, `company_name`, `summary` | text | No | Employment details |
| `candidate_work_experience` | `start_date`, `end_date` | date | No | Employment period |
| `candidate_languages` | `cand_id` | uuid | No | Candidate FK |
| `candidate_languages` | `language` | text | Yes | Language name |
| `candidate_languages` | `level` | text | No | Proficiency; unrestricted text |
| `candidate_documents` | `cand_id` | uuid | No | Candidate FK |
| `candidate_documents` | `docu_name`, `doc_type` | text | No | Document name and type; `docu_name` is the actual column name |
| `candidate_documents` | `doc_size` | integer | No | Size; units and valid range are unspecified |
| `candidate_documents` | `doc_url` | text | Yes | Document reference; file contents are not stored here |
| `candidate_documents` | `uploaded_at` | timestamp | No | No default |
| `candidate_certifications` | `cand_id` | uuid | No | Candidate FK |
| `candidate_certifications` | `certification_name` | text | Yes | Certification name |
| `candidate_certifications` | `certification_type`, `verification_link` | text | No | Certification metadata |

There are no date-order, URL-format, or document-size checks on these tables.

### `candidate_amarat`

Stores earned Amara awards, their provenance, and the definition used at award time. Repeated awards are intentionally allowed; there is no uniqueness constraint on candidate, module, title, or source execution.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK; default `gen_random_uuid()` |
| `cand_id` | uuid | No | FK → `candidates.id` |
| `amara_title` | text | Yes | Award name |
| `amara_description` | text | No | Award description |
| `module_id` | uuid | No | FK → `modules.id` |
| `company_id` | uuid | No | FK → `companies.id` |
| `application_id` | uuid | No | FK → `applications.id` |
| `source_stage_execution_id` | uuid | No | FK → `stage_executions.id` |
| `definition_version` | integer | No | Must be greater than zero when present |
| `definition_snapshot` | jsonb | No | Award definition and level scale snapshot |
| `score` | integer | No | Between 0 and 100 |
| `level`, `strengths`, `areas_for_development` | text | No | Award interpretation |
| `pass_decision` | jsonb | No | Decision metadata; no JSON shape check |
| `awarded_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |

An ordinary index covers `(cand_id, awarded_at)`. The separate FKs do not establish that the award's candidate, application, company, module, and source execution belong to the same hiring workflow.

## Jobs and applications

### `job_postings`

Stores vacancy details, publication dates, and a pipeline reference.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `company_id` | uuid | Yes | FK → `companies.id` |
| `title` | text | Yes | Posting title |
| `description`, `language`, `currency`, `seniority_level`, `employment_type`, `location` | text | No | Posting details; no allowed-value checks |
| `salary_min`, `salary_max` | integer | No | No range, ordering, or currency-unit checks |
| `salary_visible` | boolean | No | No default |
| `open_positions`, `max_applications` | integer | No | No positive-value checks |
| `application_opening`, `application_deadline` | timestamp | No | No date-order check |
| `status` | text | No | No default or allowed-value check |
| `last_autosaved_at`, `published_at`, `archived_at`, `created_at`, `updated_at`, `deleted_at` | timestamp | No | No defaults |
| `pipeline_id` | uuid | No | FK → `pipelines.id`; not unique |

### `applications`

Stores the link between a candidate and a posting, plus screening and final-decision fields.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `reference_number` | integer | Yes | UQ; no sequence or generation default |
| `posting_id` | uuid | Yes | FK → `job_postings.id`; individually UQ |
| `cand_id` | uuid | Yes | FK → `candidates.id`; individually UQ |
| `status` | text | No | No default or allowed-value check |
| `created_at`, `updated_at`, `submitted_at`, `withdrawn_at` | timestamp | No | No defaults |
| `final_decision` | boolean | No | Final-decision flag; no encoded decision vocabulary |
| `final_feedback` | text | No | Final feedback |
| `screened_out` | boolean | No | Screening flag; no default |
| `screening_completed_at` | timestamp | No | Screening completion timestamp |

**Current cardinality:** a posting can have at most one application, and a candidate can have at most one application across all postings. This follows from the two individual unique constraints. They do not represent a composite `UNIQUE (posting_id, cand_id)` rule allowing many candidates per posting and many applications per candidate.

## Assessment modules and pipeline definitions

### `modules`

Stores reusable assessment or interview definitions. The database does not restrict `type` to an enumerated set or require a company for company-sourced content.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `name`, `type` | text | Yes | Module name and type |
| `description` | text | No | Module description |
| `company_id` | uuid | No | FK → `companies.id` |
| `active_time_limit` | interval | No | Module-level time limit; no positivity check |
| `published` | boolean | No | No default |
| `source` | text | Yes | Default `company`; `platform` or `company` |
| `completion_window` | interval | No | Greater than zero |
| `selection_mode` | text | Yes | Default `specific`; `specific` or `random` |
| `question_count` | integer | No | Greater than zero |
| `difficulty_mix` | jsonb | Yes | Default `{}`; no shape check |
| `amara_definitions` | jsonb | Yes | Default `[]`; must be an array |
| `rubric` | jsonb | Yes | Default `[]`; must be an array |
| `interview_duration_minutes` | integer | No | Greater than zero |
| `created_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |
| `updated_at` | timestamptz | No | No automatic update |

The migration comment describes Amara definitions containing name, crest, domain, version, and level bands, with rubric items and questions naming one configured Amara. Publication must validate those semantic rules; the SQL only checks the top-level array types.

### `pipelines`

Stores a posting's hiring workflow and identity/feedback policy.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `posting_id` | uuid | Yes | FK → `job_postings.id`; not unique |
| `published_version_id` | uuid | No | FK → `pipeline_versions.id` |
| `status` | text | No | Unrestricted text |
| `anonymity_enabled` | boolean | Yes | Default false |
| `identity_reveal_point` | text | No | `stage` or `finalist` |
| `reveal_stage_id` | uuid | No | FK → `pipeline_stages.id` |
| `rejection_feedback_mode` | text | Yes | Default `ai_with_approval`; `human`, `ai`, or `ai_with_approval` |
| `created_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |
| `updated_at` | timestamptz | No | No automatic update |

When anonymity is enabled, `identity_reveal_point` is required: `stage` requires a non-null `reveal_stage_id`, while `finalist` requires it to be null. When anonymity is disabled, this conditional check does not require either configuration. The FK alone does not ensure the reveal stage belongs to this pipeline.

`job_postings.pipeline_id` and `pipelines.posting_id` form a circular reference. Neither is unique and no constraint ensures the two pointers agree. A writer can insert a posting with a null pipeline, insert the pipeline, then set the posting's pipeline pointer.

### `pipeline_versions`

Separates workflow revisions from candidate executions.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `version_number` | integer | Yes | UQ with `pipeline_id`; no positive-value check |
| `is_published` | boolean | No | No default |
| `pipeline_id` | uuid | No | FK → `pipelines.id` |
| `previous_version_id` | uuid | No | Self-FK → `pipeline_versions.id` |
| `published_by` | uuid | No | FK → `companies.id` |
| `migration_plan` | jsonb | Yes | Default `{}`; no shape check |
| `created_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |
| `published_at` | timestamptz | No | Publication timestamp |

`UNIQUE (pipeline_id, version_number)` distinguishes revisions within a pipeline. A partial unique index on `pipeline_id WHERE is_published IS TRUE` permits at most one published version per non-null pipeline. Nullable pipeline IDs can bypass these per-pipeline uniqueness guarantees under normal PostgreSQL null semantics.

The database does not enforce that `published_version_id` points to a published version of the same pipeline, that the previous version belongs to the same pipeline, or that published definitions are immutable.

### `pipeline_stages`

Defines stages within one pipeline version.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `version_id` | uuid | Yes | FK → `pipeline_versions.id` |
| `name`, `module_type` | text | Yes | Stage name and module type |
| `advancement_rule` | text | Yes | Default `manual`; `manual`, `threshold`, or `completion` |
| `score_threshold` | integer | No | Between 0 and 100; required for threshold advancement |
| `below_threshold_action` | text | Yes | Default `hold`; `hold` or `reject` |
| `feedback_mode` | text | Yes | Default `ai_with_approval`; `human`, `ai`, or `ai_with_approval` |
| `scheduled_start_at`, `scheduled_end_at` | timestamptz | No | End must be after start when both exist |
| `position_x`, `position_y` | numeric | Yes | Default 0; visual editor coordinates |

`UNIQUE (version_id, name)` prevents duplicate stage names in a version. `UNIQUE (version_id, id)` supports the composite foreign keys on transitions. `module_type` is unrestricted text and is not an FK to `modules`.

### `stage_transitions`

Stores directed edges between stages.

| Column | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `pipeline_version_id` | uuid | Yes | FK → `pipeline_versions.id` |
| `from_stage_id` | uuid | Yes | Composite FK with version → `pipeline_stages(version_id, id)` |
| `to_stage_id` | uuid | Yes | Composite FK with version → `pipeline_stages(version_id, id)` |

`UNIQUE (pipeline_version_id, from_stage_id, to_stage_id)` prevents duplicate edges. `from_stage_id <> to_stage_id` prevents self-loops. Both endpoint stages must belong to the specified version. Longer cycles, unreachable stages, and valid entry/exit structure require application validation.

### `assessment_questions`

Stores questions belonging to a test module and their scoring configuration.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `test_module_id` | uuid | Yes | FK → `modules.id` |
| `prompt` | text | Yes | Question text |
| `response_format` | text | Yes | `multiple_choice`, `short_text`, `long_text`, `code`, or `video` |
| `scoring_method` | text | Yes | `answer_key`, `code_evaluator`, or `human` |
| `max_points` | numeric | Yes | Greater than zero |
| `answer_key`, `options`, `test_cases` | jsonb | No | No shape checks |
| `difficulty` | text | No | Unrestricted text |
| `amara_name` | text | Yes | Semantic link to a configured Amara; no FK |
| `grader_id` | uuid | No | FK → `interviewers.id` |
| `grader_company_id` | uuid | No | FK → `companies.id` |
| `ai_generated` | boolean | Yes | Default false |
| `reviewed_by` | uuid | No | FK → `companies.id` |
| `reviewed_at` | timestamptz | No | Review timestamp |
| `position` | integer | Yes | Default 0; at least zero; not unique within a module |
| `configuration` | jsonb | Yes | Default `{}`; no shape check |
| `created_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |
| `updated_at` | timestamptz | No | No automatic update |

At most one of `grader_id` and `grader_company_id` may be populated; both may be null. SQL does not require a grader for human scoring, review for AI-generated questions, or a compatible response-format/scoring-method combination.

### `test_configurations`

Connects a stage to a module and optional stage-specific timing/configuration.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `stage_id` | uuid | Yes | UQ; FK → `pipeline_stages.id` |
| `test_module_id` | uuid | Yes | FK → `modules.id` |
| `active_time_limit`, `completion_window` | interval | No | Each must be greater than zero |
| `configuration` | jsonb | Yes | Default `{}`; no shape check |
| `created_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |
| `updated_at` | timestamptz | No | No automatic update |

A stage has at most one test configuration. Multiple stages may use the same module. Timing precedence between module and stage configuration is not encoded in SQL.

## Pipeline execution and assessments

### `pipeline_executions`

Stores an application's execution against a specific pipeline version.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `application_id` | uuid | Yes | UQ; FK → `applications.id` |
| `pipeline_version_id` | uuid | Yes | FK → `pipeline_versions.id` |
| `status` | text | Yes | Default `active`; `active`, `completed`, `rejected`, `withdrawn`, or `cancelled` |
| `anonymous_identifier` | uuid | Yes | UQ; default `gen_random_uuid()` |
| `started_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |
| `completed_at` | timestamptz | No | Completion timestamp |
| `rejection_feedback` | jsonb | No | No shape check |
| `feedback_is_ai_generated` | boolean | Yes | Default false |
| `feedback_approved_at`, `feedback_released_at`, `updated_at` | timestamptz | No | Feedback and update timestamps |

An application has at most one execution, including historical executions; no attempt/version dimension relaxes that uniqueness. An ordinary index covers `pipeline_version_id`.

### `stage_executions`

Stores one stage's runtime state for a pipeline execution.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `pipeline_execution_id` | uuid | Yes | FK → `pipeline_executions.id` |
| `stage_id` | uuid | Yes | FK → `pipeline_stages.id` |
| `status` | text | Yes | Default `locked`; see allowed values below |
| `available_at`, `deadline`, `started_at`, `completed_at` | timestamptz | No | Runtime timestamps |
| `score` | integer | No | Between 0 and 100 |
| `decision_type` | text | No | `manual`, `threshold`, `completion`, or `override` |
| `decided_by` | uuid | No | FK → `companies.id` |
| `decided_at` | timestamptz | No | Decision timestamp |
| `decision_reason` | text | No | Required and nonblank for an override |
| `score_at_decision` | integer | No | Between 0 and 100 |
| `decision_history` | jsonb | Yes | Default `[]`; must be an array |
| `feedback` | text | No | Stage feedback |
| `feedback_is_ai_generated` | boolean | Yes | Default false |
| `feedback_approved_at`, `feedback_released_at` | timestamptz | No | Approval and release timestamps |

Allowed statuses are `locked`, `available`, `in_progress`, `pending_grading`, `held`, `passed`, `rejected`, `cancelled`, and `missed`. `UNIQUE (pipeline_execution_id, stage_id)` prevents repeated execution rows for the same pair. An ordinary index covers `(stage_id, status)`.

The stage FK does not ensure the stage belongs to the execution's pipeline version. Status checks restrict vocabulary but do not enforce legal transitions or require timestamp/decision fields for a particular status.

### `candidate_stage_progress`

Preserves the earlier per-application progress model and optionally links it to a stage execution.

| Column | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `application_id` | uuid | Yes | FK → `applications.id` |
| `stage_id` | uuid | Yes | FK → `pipeline_stages.id` |
| `status` | text | Yes | No default or allowed-value check |
| `amara_awarded` | text | No | Text field; not an FK to `candidate_amarat` |
| `stage_execution_id` | uuid | No | UQ; FK → `stage_executions.id` |

There is no uniqueness constraint on `(application_id, stage_id)`. SQL does not synchronize this table with `stage_executions` or `candidate_amarat`.

### `assessment_instances`

Stores the earlier assessment scheduling record, retained by the execution model.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `application_id` | uuid | Yes | FK → `applications.id` |
| `stage_id` | uuid | Yes | FK → `pipeline_stages.id` |
| `status` | text | Yes | No default or allowed-value check |
| `started_at`, `deadline` | timestamp | No | No defaults |

`UNIQUE (application_id, stage_id)` permits at most one instance per pair. Its timestamps use `timestamp`, while attempt timestamps use `timestamptz`; writers must handle the distinction explicitly.

### `assessment_attempts`

Stores the selected content and configuration used by a candidate for a stage assessment.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `stage_execution_id` | uuid | Yes | UQ; FK → `stage_executions.id` |
| `candidate_id` | uuid | Yes | FK → `candidates.id` |
| `instance_id` | uuid | Yes | UQ; FK → `assessment_instances.id` |
| `status` | text | Yes | Default `available`; `available`, `in_progress`, `submitted`, `pending_grading`, `scored`, or `missed` |
| `selected_content` | jsonb | Yes | Default `[]`; must be an array |
| `configuration_snapshot` | jsonb | Yes | Default `{}`; no shape check |
| `started_at` | timestamptz | No | Start timestamp |
| `deadline` | timestamptz | Yes | No default or time-order check |
| `submitted_at` | timestamptz | No | Submission timestamp |
| `auto_submitted` | boolean | Yes | Default false |
| `extension_history` | jsonb | Yes | Default `[]`; must be an array |

A stage execution and an assessment instance each have at most one attempt. Despite the table name, the schema does not support multiple attempts for either. The FKs do not ensure that candidate, instance, and stage execution refer to the same application/stage.

### `assessment_submissions`

Stores a candidate's answer and grading data for one question in an attempt.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `attempt_id` | uuid | Yes | FK → `assessment_attempts.id` |
| `question_id` | uuid | Yes | FK → `assessment_questions.id` |
| `response` | jsonb | No | Answer payload; no shape check |
| `max_points` | numeric | Yes | Greater than zero |
| `earned_points` | numeric | No | At least zero and no greater than `max_points` |
| `grading_status` | text | Yes | Default `pending`; `pending`, `evaluation_pending`, or `graded` |
| `grader_id` | uuid | No | FK → `interviewers.id` |
| `grader_company_id` | uuid | No | FK → `companies.id` |
| `graded_at` | timestamptz | No | Grading timestamp |
| `evaluator_result` | jsonb | No | Evaluation output; no shape check |
| `saved_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP`; not automatically refreshed on update |

`UNIQUE (attempt_id, question_id)` permits one stored answer per question/attempt. At most one grader reference may be non-null. The question FK does not require the question to appear in the attempt's selected content or module.

### `assessment_results`

Stores the overall result for an assessment instance, optionally connected to the newer attempt model.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `instance_id` | uuid | Yes | UQ; FK → `assessment_instances.id` |
| `passed` | boolean | Yes | No default |
| `score` | integer | No | Between 0 and 100 |
| `feedback` | text | No | Result feedback |
| `is_disclosed` | boolean | No | No default |
| `attempt_id` | uuid | No | UQ; FK → `assessment_attempts.id` |
| `earned_points` | numeric | No | At least zero; no greater than `max_points` when both exist |
| `max_points` | numeric | No | Greater than zero |
| `feedback_is_ai_generated` | boolean | Yes | Default false |
| `feedback_approved_at`, `disclosed_at`, `scored_at` | timestamptz | No | Approval, disclosure, and scoring timestamps |

There is at most one result per instance and per non-null attempt. No constraint derives the score from points, derives `passed` from a threshold, matches the result's instance to its attempt, or synchronizes `is_disclosed` with `disclosed_at`.

## Interviews, identity reveal, and offers

### `interviews`

Stores one scheduled interview per stage execution.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `stage_execution_id` | uuid | Yes | UQ; FK → `stage_executions.id` |
| `scheduled_start_at`, `scheduled_end_at` | timestamptz | Yes | End must be after start |
| `status` | text | Yes | Default `pending_confirmation`; `pending_confirmation`, `confirmed`, `completed`, or `cancelled` |
| `outcome` | text | No | `held`, `candidate_absent`, or `interviewer_absent` |
| `provider`, `provider_meeting_id`, `join_url` | text | No | Meeting-provider details |
| `reschedule_count` | smallint | Yes | Default 0; between 0 and 2 |
| `reschedule_history`, `rubric_snapshot` | jsonb | Yes | Each defaults to `[]` and must be an array |
| `score` | integer | No | Between 0 and 100 |
| `created_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |
| `updated_at` | timestamptz | No | No automatic update |

The reschedule-count bound does not verify the number of history entries. Interview completion, outcomes, and score calculation are not coordinated by SQL.

### `interviewer_assignments`

Associates interviewers with interviews; its composite PK implements the many-to-many relationship.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `interview_id` | uuid | Yes | Composite PK; FK → `interviews.id` |
| `interviewer_id` | uuid | Yes | Composite PK; FK → `interviewers.id` |
| `assigned_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |
| `available_slots` | jsonb | Yes | Default `[]`; must be an array |
| `scorecard_status` | text | Yes | Default `pending`; `pending`, `submitted`, or `waived` |
| `waived_by` | uuid | No | FK → `companies.id` |
| `waived_at` | timestamptz | No | Waiver timestamp |
| `waiver_reason` | text | No | Waiver explanation |

If `scorecard_status = 'waived'`, the company, timestamp, and a nonblank reason are all required. A submitted status does not itself require a matching feedback row.

### `interview_feedback`

Stores a scorecard for an assigned interviewer. The migration explicitly describes `notes` as private hiring-team notes that must never be included in candidate feedback; this restriction is a data-handling requirement, not a SQL access policy.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `interview_id`, `interviewer_id` | uuid | Yes | Composite FK → `interviewer_assignments(interview_id, interviewer_id)` |
| `notes` | text | No | Private hiring-team notes |
| `notes_updated_at` | timestamp | No | Uses `timestamp`, unlike other feedback timestamps |
| `ratings` | jsonb | Yes | Default `{}`; must be an object |
| `recommendation` | text | No | Unrestricted text |
| `submitted_at` | timestamptz | No | Submission timestamp |
| `created_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |
| `updated_at` | timestamptz | No | No automatic update |

`UNIQUE (interview_id, interviewer_id)` allows at most one feedback row per assignment. The composite FK requires the assignment to exist; there is no direct company-author FK on this table.

### `identity_reveal_records`

Records identity disclosure for an application.

| Column | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `candidate_id` | uuid | Yes | FK → `candidates.id` |
| `revealed_by` | uuid | No | FK → `companies.id` |
| `application_id` | uuid | Yes | UQ; FK → `applications.id` |
| `stage_id` | uuid | No | FK → `pipeline_stages.id` |
| `revealed_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |

An application has at most one reveal record. The database does not ensure its candidate matches the application or that the reveal follows the pipeline's configured policy.

### `offers`

Stores offer history for an application, with at most one pending offer at a time.

| Column(s) | Type | Required | Default / constraints / meaning |
| --- | --- | --- | --- |
| `id` | uuid | Yes | PK |
| `application_id` | uuid | Yes | FK → `applications.id` |
| `job_title`, `compensation_summary` | text | Yes | Offer terms |
| `proposed_start_date` | date | Yes | Proposed employment start |
| `response_deadline` | timestamptz | Yes | Must be after `created_at` |
| `message` | text | No | Offer message |
| `status` | text | Yes | Default `pending`; `pending`, `accepted`, `declined`, `expired`, or `withdrawn` |
| `response` | text | No | `accepted` or `declined` |
| `responded_at`, `withdrawn_at` | timestamptz | No | Response/withdrawal timestamps |
| `deadline_history` | jsonb | Yes | Default `[]`; must be an array |
| `capacity_confirmed_at` | timestamptz | No | Capacity confirmation timestamp |
| `created_at` | timestamptz | Yes | Default `CURRENT_TIMESTAMP` |
| `updated_at` | timestamptz | No | No automatic update |

The response check permits either:

1. No response and no response timestamp, with a status other than `accepted` or `declined`.
2. A response and response timestamp, with `status = response` and `responded_at < response_deadline`.

A partial unique index on `application_id WHERE status = 'pending'` prevents concurrent pending offers for one application. A second partial index on `response_deadline WHERE status = 'pending'` supports finding pending deadlines. Expiry, job-capacity enforcement, deadline-history maintenance, and application-status updates require application or worker logic.

## Relationship and workflow summary

The following summarizes the intended use of the records. It does not claim these operations or transitions are implemented by the migrations.

1. A candidate creates an account and profile. Verification challenges and optional MFA credentials refer to the account by `(account_id, account_type)`. Company review requests are stored separately from company accounts.
2. A company creates a posting and pipeline. Pipeline versions contain stages, with transitions constrained to connect stages in the same version. Test configurations associate stages with reusable modules and questions.
3. An application links a candidate and posting. Screening and final decisions are represented on the application. The current individual uniqueness constraints restrict both candidate and posting to one application.
4. Admission can create a pipeline execution referencing a selected version. Stage executions record progression. The earlier progress and assessment-instance tables remain linked records rather than being replaced.
5. An assessment attempt captures selected content and configuration, with one submission per question. Assessment results hold aggregate points, pass/fail, and disclosure metadata.
6. Interview stages can create an interview, assign interviewers, and collect one scorecard per assignment. Identity reveals record one disclosure per application.
7. Stage results can produce repeated Amara awards with definition snapshots. Offer rows retain history, while the partial unique index restricts each application to one pending offer.

Definition snapshots preserve the information used at execution or award time, but the database does not make snapshots or history arrays immutable. Keeping legacy and newer execution records consistent requires coordinated writes.

## Indexes and query support

Primary keys and unique constraints create indexes automatically. Composite uniqueness is documented with each table. The migrations also explicitly create these indexes:

| Index | Key / expression | Condition | Purpose |
| --- | --- | --- | --- |
| `candidates_email_lower_idx` | `lower(email)` | Unique | Case-insensitive candidate email uniqueness |
| `interviewers_email_lower_idx` | `lower(email)` | Unique | Case-insensitive interviewer email uniqueness |
| `companies_email_lower_idx` | `lower(email)` | Unique | Case-insensitive company email uniqueness |
| `auth_email_verification_challenges_account_idx` | `(account_id, account_type)` | None | Account challenge lookup |
| `auth_company_registration_email_idx` | `lower(email)` | Unique where status is `pending_review` or `approved` | Active registration email uniqueness |
| `auth_company_registration_number_idx` | `lower(registration_number)` | Unique where status is `pending_review` or `approved` | Active registration-number uniqueness |
| `pipeline_versions_one_published_idx` | `pipeline_id` | Unique where `is_published IS TRUE` | One published version per non-null pipeline |
| `pipeline_executions_version_idx` | `pipeline_version_id` | None | Executions using a version |
| `stage_executions_stage_status_idx` | `(stage_id, status)` | None | Stage/status filtering |
| `offers_one_pending_application_idx` | `application_id` | Unique where status is `pending` | One pending offer per application |
| `offers_pending_deadline_idx` | `response_deadline` | Where status is `pending` | Pending-offer deadline lookup |
| `candidate_amarat_candidate_awarded_idx` | `(cand_id, awarded_at)` | None | Candidate award history |

Email uniqueness is per account table, not global across candidates, interviewers, and companies. Lowercasing does not trim whitespace or validate email syntax. Foreign keys do not automatically create indexes on their referencing columns; the migrations do not index every FK.

## Integrity, access control, and lifecycle boundaries

The SQL enforces entity references, the documented uniqueness rules, selected status vocabularies, bounded scores and points, JSON container types where explicitly checked, and the conditional rules for anonymity, overrides, waivers, and offer responses.

The migrations do not define row-level security enablement, policies, grants, or revocations. They therefore establish no application-specific database access rules or tenant isolation policies. This statement does not establish the effective permissions of an existing Supabase project, which may have configuration outside this repository.

Company ownership is represented by FKs, primarily `job_postings.company_id`, `interviewers.company_id`, and `modules.company_id`. Cross-table references usually check that a row exists, not that it belongs to the same company. Backend authorization must check tenant ownership, account approval/status, interviewer assignment, grading authority, and candidate visibility. The schema has no shared `tenant_id` column or composite tenant keys.

Verification codes have hash fields, while MFA has a binary secret field. Hashing algorithms, encryption, challenge consumption, and secret exposure controls are not specified by this DDL. Feedback approval/release timestamps and disclosure flags store state but do not prevent premature reads. Private interview notes need explicit filtering in application responses.

`deleted_at` and `deactivated_at` support lifecycle tracking on account tables, and postings also have `deleted_at`; no trigger or policy automatically filters these rows. Hard deletion can be blocked by referencing rows because no cascading actions are declared. The migrations define no retention jobs, automatic expiry, timestamp-update triggers, audit append-only protections, or backup/restore configuration.

## Migration history and limitations

### Forward schema changes

All paths in the following table are relative to [`supabase/migrations/`](../../../supabase/migrations/).

| Migration(s) | Contribution |
| --- | --- |
| `20261027223035_baseline.sql` | Comments describing future epic schemas; no schema objects |
| `20261029000001`–`20261029000007` | Initial candidate/profile, company/interviewer, posting, and application tables |
| `20261103202052`–`20261103202058` (removed) | Former duplicate table definitions; the October migrations already create these tables |
| `20261103202059_create_modules_table.sql` | Assessment modules |
| `20261103202100_link_amarat_to_modules.sql` | Award-to-module FK |
| `20261103202101_create_pipelines_table.sql` | Pipelines, versions, stages, and version-to-pipeline FK |
| `20261103202102_link_postings_to_pipelines.sql` | Posting-to-pipeline FK |
| `20261103202103_create_candidate_stage_progress_table.sql` | Earlier candidate progression records |
| `20261103202104_create_assessment_instances_table.sql` | Assessment instances |
| `20261103202105_create_assessment_results_table.sql` | Assessment results |
| [`20261105065037_configure_invw_pipeline_definitions.sql`](../../../supabase/migrations/20261105065037_configure_invw_pipeline_definitions.sql) | Module configuration, anonymity, versioning, stage rules, transitions, questions, and test configurations |
| [`20261105065038_create_invw_execution_tables.sql`](../../../supabase/migrations/20261105065038_create_invw_execution_tables.sql) | Executions, attempts, submissions, and links to earlier assessment/progress records |
| [`20261105065039_create_invw_live_interview_tables.sql`](../../../supabase/migrations/20261105065039_create_invw_live_interview_tables.sql) | Interviews, assignments, and feedback |
| [`20261105065040_create_invw_offers_and_awards.sql`](../../../supabase/migrations/20261105065040_create_invw_offers_and_awards.sql) | Identity reveals, extended awards, and offers |
| [`20261105232900_auth_registration.sql`](../../../supabase/migrations/20261105232900_auth_registration.sql) | Account verification status, company approval, email indexes, and verification challenges |
| [`20261105232901_company_registration_requests.sql`](../../../supabase/migrations/20261105232901_company_registration_requests.sql) | Company review requests and partial uniqueness |
| [`20261105232902_auth_mfa.sql`](../../../supabase/migrations/20261105232902_auth_mfa.sql) | MFA credentials |

### Migration compatibility and verification

The November `20261103202052`–`20261103202058` files repeated the October table definitions. This caused `relation "candidates" already exists` when replaying migrations. Their definitions were verified to match the October definitions, and the duplicate files have been removed. The original October migrations remain responsible for creating these tables.

The November files also originally included executable rollback SQL after `-- +goose Down`. Supabase executes the whole SQL file; those markers do not protect rollback statements. The rollback sections and Goose annotations have been removed so migrations contain only forward changes. Reversals must be expressed as separate, deliberate migrations.

Run `make test-migrations` with Python 3 and PostgreSQL's `initdb`, `pg_ctl`, and `psql` on `PATH`. The [migration replay test](../../../tests/migrations/test_migrations.py) creates a temporary local PostgreSQL cluster, executes every complete migration file in filename order, compares all 33 tables and 373 columns with the ER diagram, and verifies that a candidate inserted after the October migrations survives with its later account-status default. It stops the cluster and removes its temporary files afterward. The repaired sequence passed on PostgreSQL 18.6; the Supabase configuration targets PostgreSQL 17, so this does not replace a Supabase reset on the configured version.

No remote database or remote migration history was changed by this repair. Existing deployments require comparison with their actual schema and migration history before applying changes, including whether any removed duplicate versions were recorded. Modifying or removing a file does not reapply or undo an already recorded migration.

### Diagram differences and incomplete guarantees

| Topic | Diagram presentation | SQL definition / consequence |
| --- | --- | --- |
| Posting and pipeline | Optional one-to-one | `pipelines.posting_id` is required, but neither posting/pipeline pointer is unique or checked against the reverse pointer |
| Stage test configuration | One stage to many configurations | `test_configurations.stage_id` is unique, so a stage has at most one configuration |
| Company-authored interview feedback | `COMPANIES` → `INTERVIEW_FEEDBACK` relationship labeled `submits` | No company-author column/FK; feedback references an interviewer assignment |
| Candidate profile ownership | Relationships imply a candidate owner | Profile-detail and award `cand_id` columns are nullable |
| Pipeline version ownership | A version belongs to a pipeline | `pipeline_versions.pipeline_id` is nullable |
| Interviewer company membership | Relationship implies company ownership | `interviewers.company_id` is nullable |
| Authentication relationships | Account-to-challenge/MFA edges | Polymorphic account references have no SQL FK, as the diagram comment notes |
| Composite relationships | Some edges and key labels are simplified | Transition endpoints reference `(version_id, id)`; feedback references an assignment pair; other composite uniqueness rules require the SQL/dictionary |

The diagram correctly reflects the restrictive application cardinalities created by the individual candidate/posting unique constraints. Whether those constraints should be replaced by pairwise uniqueness is a design decision for a later migration, not a behavior assumed by this reference.

Other application responsibilities include graph validation, same-version and same-tenant consistency, valid state transitions, score calculation, grading completeness, feedback approval, automatic submission, offer expiry/capacity, and maintaining history/snapshot data. Adding a timestamp, status field, or JSON history column alone does not enforce those behaviors.

## Maintaining this reference

When a migration changes the schema, update the relevant dictionary entries, relationships, checks, and index list, then reconcile the [ER diagram](erd_diagrams/main.mmd). Describe forward changes separately from rollback instructions. Validate migration order and deployed history independently before claiming a successful database rebuild or deployment.
