# Git Conventions

Owner: Ziad Eliwa | Status: Draft | Last updated: 2026-10-08 | Jira: AMARA-59

## 1. Purpose

This document defines the Git conventions used by the project.

The objectives are to:

* maintain a clean and understandable Git history;
* connect development work to Jira issues;
* maintain traceability between requirements, Jira work items, implementation, and tests;
* make changes easy to review;
* provide a consistent workflow for all contributors.

These conventions apply to all contributors.

---

## 2. Git Repository

The project MUST use Git for version control.

The repository MUST contain only files that belong to the project or are explicitly required for development.

Generated files, build artifacts, secrets, credentials, and temporary files MUST NOT be committed.

A `.gitignore` file MUST be maintained for the technologies used by the project.

---

## 3. Branching Strategy

The project uses a feature-branch workflow.

The primary branches are:

* `main` — stable, production-ready code.
* `develop` — optional integration branch if the project requires one.

If `develop` is not used, feature branches MAY be merged into `main` through Pull Requests.

Contributors MUST NOT directly push feature work to protected branches.

---

## 4. Jira Issue Association

Jira is the project's work-management system.

Each implementation task SHOULD have a corresponding Jira issue.

The Jira issue key MUST be used to connect the Jira task with Git work.

For example:

```text
Jira:
AMARA-50 — Activity diagrams Epic 1

Git branch:
docs/AMARA-50-activity-diagrams-epic-1

Commit:
docs(AMARA-50): add activity diagrams for epic 1

Pull Request:
AMARA-50 Activity diagrams Epic 1
```

The Jira key identifies the **work item**, while the SRS requirement ID identifies the **system requirement**.

For example:

```text
AUTH-FR-001
    ↓
AMARA-101
    ↓
feature/AMARA-101-user-login
    ↓
feat(AMARA-101): implement user login
```

---

## 5. Requirement and Jira IDs

SRS requirements MUST have their own domain-specific identifiers.

Examples:

```text
AUTH-FR-001
AUTH-FR-002

JOB-FR-001
JOB-FR-002

APP-FR-001
APP-FR-002

ANA-FR-001
```

Jira issues use the project's Jira key:

```text
AMARA-101
AMARA-102
AMARA-103
```

The two identifiers serve different purposes:

| Identifier                              | Purpose                         |
| --------------------------------------- | ------------------------------- |
| `AUTH-FR-001`                           | Identifies a system requirement |
| `AMARA-101`                             | Identifies a Jira work item     |
| `feature/AMARA-101-user-login`          | Identifies the Git branch       |
| `feat(AMARA-101): implement user login` | Identifies the Git change       |

The project MUST NOT attempt to create separate Jira projects solely to obtain domain-specific keys such as `AUTH-1`, `JOB-1`, or `APP-1`.

Domain organization SHOULD instead be handled through SRS requirement IDs, Jira Components, labels, and issue summaries.

---

## 6. Branch Naming

Branches MUST follow:

```text
<type>/<JIRA-KEY>-<short-description>
```

Allowed types:

| Type       | Purpose                     |
| ---------- | --------------------------- |
| `feature`  | New functionality           |
| `fix`      | Bug fix                     |
| `refactor` | Code restructuring          |
| `test`     | Testing changes             |
| `docs`     | Documentation changes       |
| `chore`    | General maintenance         |
| `build`    | Build or dependency changes |
| `ci`       | CI/CD changes               |

Examples:

```text
feature/AMARA-101-user-login
feature/AMARA-125-job-search
fix/AMARA-143-duplicate-application
refactor/AMARA-167-application-service
test/AMARA-101-login-tests
docs/AMARA-60-create-srs-skeleton
docs/AMARA-59-documentation-conventions
ci/AMARA-180-github-actions
```

Branch descriptions MUST use lowercase kebab-case.

Branch names MUST NOT contain spaces.

---

## 7. Commit Messages

Commits associated with a Jira issue MUST follow:

```text
<type>(<JIRA-KEY>): <description>
```

Examples:

```text
feat(AMARA-101): implement user login
fix(AMARA-143): prevent duplicate applications
test(AMARA-101): add login integration tests
docs(AMARA-60): create SRS document structure
docs(AMARA-59): add documentation conventions
ci(AMARA-180): add automated test pipeline
```

Allowed commit types:

| Type       | Purpose                  |
| ---------- | ------------------------ |
| `feat`     | New functionality        |
| `fix`      | Bug fix                  |
| `refactor` | Refactoring              |
| `test`     | Tests                    |
| `docs`     | Documentation            |
| `chore`    | Maintenance              |
| `build`    | Build/dependency changes |
| `ci`       | CI/CD changes            |
| `perf`     | Performance improvements |
| `style`    | Formatting/style changes |

Commit descriptions MUST:

* be concise;
* describe the actual change;
* begin with a lowercase letter;
* use imperative language where practical;
* NOT end with a period.

Good:

```text
feat(AMARA-101): implement user login
```

Bad:

```text
feat(AMARA-101): Added the user login functionality.
```

---

## 8. Commits Without Jira Issues

A Jira key SHOULD be included whenever a commit contributes to a specific Jira issue.

A Jira key is not required for genuinely project-wide changes that do not correspond to a Jira issue.

For example:

```text
chore: update gitignore
```

may be acceptable if there is no corresponding Jira issue.

However, project work that is tracked in Jira SHOULD use the Jira key.

---

## 9. Atomic Commits

Each commit SHOULD represent one logical change.

Unrelated changes MUST NOT be combined into a single commit.

For example, these should normally be separate:

```text
feat(AMARA-101): implement user login
fix(AMARA-143): prevent duplicate applications
docs(AMARA-60): update SRS structure
```

A commit MAY contain multiple files when those files are part of the same logical change.

---

## 10. Commit History

Contributors SHOULD commit regularly during development.

Temporary commits such as:

```text
test
fix
working
asdf
final
final-final
```

MUST NOT remain in the final branch history.

Before merging, contributors SHOULD squash or reorganize temporary commits where appropriate.

---

## 11. Pull Requests

Changes to protected branches MUST be submitted through a Pull Request (PR).

A PR MUST:

* reference the relevant Jira issue;
* have a clear title;
* describe the changes;
* reference relevant SRS requirements where applicable;
* pass required CI checks;
* receive the required code review;
* contain no unresolved merge conflicts;
* satisfy the project's Definition of Done.

PR titles SHOULD follow:

```text
<JIRA-KEY> <short-description>
```

Example:

```text
AMARA-101 Implement user login
```

---

## 12. Pull Request Description

A PR SHOULD contain:

```markdown
## Summary

Brief description of the change.

## Jira

AMARA-101

## Requirements

- AUTH-FR-001

## Changes

- Added login endpoint
- Added credential validation
- Added authentication tests

## Testing

- Unit tests
- Integration tests

## Notes

Any relevant implementation or deployment considerations.
```

---

## 13. Requirement Traceability

Where applicable, Jira issues SHOULD reference the SRS requirements they implement.

A complete traceability chain SHOULD look like:

```text
SRS Requirement
AUTH-FR-001
       ↓
Jira Issue
AMARA-101
       ↓
Git Branch
feature/AMARA-101-user-login
       ↓
Git Commit
feat(AMARA-101): implement user login
       ↓
Pull Request
AMARA-101 Implement user login
       ↓
Tests
Authentication tests
```

This allows the project to trace a change from its original requirement through implementation and verification.

---

## 14. Merging

Pull Requests SHOULD use squash merging unless there is a specific reason to preserve individual commits.

The resulting commit MUST retain the Jira issue key.

Example:

```text
feat(AMARA-101): implement user login
```

Branches SHOULD be deleted after successful merging.

Merge commits SHOULD NOT be created unnecessarily.

---

## 15. Rebasing

Contributors SHOULD keep long-lived branches synchronized with their target branch.

Rebasing MAY be used to maintain a clean history.

Contributors MUST NOT rewrite the history of shared or protected branches.

Force pushing to protected branches MUST be prohibited.

If force pushing is necessary on a personal feature branch, contributors SHOULD use:

```bash
git push --force-with-lease
```

instead of:

```bash
git push --force
```

---

## 16. Conflict Resolution

Merge conflicts MUST be resolved carefully.

The contributor resolving a conflict MUST verify that:

* the intended behavior of both changes is preserved;
* tests still pass;
* no unrelated changes were introduced;
* the resulting code is consistent with the current requirements.

Conflict resolution MUST NOT simply favor one side without understanding the changes.

---

## 17. Protected Branches

Protected branches SHOULD require:

* Pull Requests;
* successful CI checks;
* required code review;
* no unresolved review conversations;
* successful merge requirements.

Direct pushes to protected branches SHOULD be disabled.

---

## 18. Releases and Tags

Release tags SHOULD follow Semantic Versioning:

```text
v<MAJOR>.<MINOR>.<PATCH>
```

Examples:

```text
v1.0.0
v1.1.0
v1.1.1
```

Version numbers follow:

* **MAJOR** — incompatible changes;
* **MINOR** — backward-compatible functionality;
* **PATCH** — backward-compatible fixes.

Release tags MUST reference a stable commit.

---

## 19. Files That Must Not Be Committed

The following MUST NOT be committed:

* passwords;
* API keys;
* access tokens;
* private keys;
* production credentials;
* `.env` files containing secrets;
* sensitive database dumps;
* build artifacts;
* temporary files;
* personal IDE configuration unless explicitly required.

Examples of files/directories that may need to be ignored:

```text
.env
.env.local
*.pem
*.key
node_modules/
__pycache__/
dist/
build/
```

The project's `.gitignore` MUST be adapted to the technologies used by the project.

---

## 20. Standard Git Workflow

The standard workflow is:

```text
Jira Issue
    ↓
Create Branch
    ↓
Implement
    ↓
Test
    ↓
Commit
    ↓
Push
    ↓
Pull Request
    ↓
Code Review
    ↓
CI Validation
    ↓
Merge
    ↓
Delete Branch
```

Example:

```bash
git checkout main
git pull

git checkout -b docs/AMARA-60-create-srs-skeleton

# Make changes

git add .
git commit -m "docs(AMARA-60): create SRS document structure"

git push -u origin docs/AMARA-60-create-srs-skeleton
```

Then create the Pull Request:

```text
AMARA-60 Create SRS skeleton
```

---

## 21. Exceptions

Exceptions to these conventions MUST be documented and agreed upon by the project team.

Exceptions SHOULD be rare and SHOULD have a clear technical justification.

---

## 22. Summary

For Jira-tracked work, the project follows this convention:

```text
Jira Issue
AMARA-XXX
    ↓
Branch
<type>/AMARA-XXX-description
    ↓
Commit
<type>(AMARA-XXX): description
    ↓
Pull Request
AMARA-XXX Description
```

For example:

```text
AMARA-50
Activity diagrams Epic 1

        ↓

docs/AMARA-50-activity-diagrams-epic-1

        ↓

docs(AMARA-50): add activity diagrams for epic 1

        ↓

AMARA-50 Activity diagrams Epic 1
```

SRS requirements remain separately identifiable through domain-specific requirement IDs such as:

```text
AUTH-FR-001
JOB-FR-001
APP-FR-001
ANA-FR-001
```

This separation ensures that **SRS IDs identify requirements while Jira keys identify work items**.
