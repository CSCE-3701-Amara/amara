# Frontend sitemap

This inventory reflects the current files under `web/src/app`. Dynamic
segments use the `[id]` notation from the App Router.

## Public

| Path | Purpose |
| --- | --- |
| `/` | Welcome screen |
| `/login` | Candidate sign-in |
| `/register` | Candidate registration |

Both auth pages accept `?account=company`; auth state variants are documented
in [routing-and-auth.md](routing-and-auth.md).

## Candidate workspace

| Path | Purpose |
| --- | --- |
| `/candidate/dashboard` | Job feed |
| `/candidate/jobs/[id]` | Job details |
| `/candidate/applications` | Application history |
| `/candidate/applications/[id]` | Application details |
| `/candidate/applications/[id]/progress` | Pipeline progress |
| `/candidate/assessments/[id]` | Assessment attempt |
| `/candidate/assessments/[id]/coding` | Coding assessment |
| `/candidate/assessments/[id]/video` | Video assessment |
| `/candidate/interviews/[id]` | Candidate interview |
| `/candidate/offers` | Offer list |
| `/candidate/offers/[id]` | Offer details |
| `/candidate/profile` | Candidate profile |
| `/candidate/profile/edit` | Edit candidate profile |

## Employer workspace

| Path | Purpose |
| --- | --- |
| `/employer/jobs` | Posting list |
| `/employer/jobs/new` | Create posting |
| `/employer/jobs/[id]` | Edit posting |
| `/employer/jobs/[id]/preview` | Posting preview |
| `/employer/jobs/[id]/applications` | Posting applicants |
| `/employer/jobs/[id]/comparison` | Candidate comparison |
| `/employer/jobs/[id]/finalists` | Finalist list |
| `/employer/pipelines` | Pipeline list |
| `/employer/pipelines/[id]` | Pipeline studio |
| `/employer/pipelines/[id]/candidates` | Stage candidates |
| `/employer/pipelines/[id]/dashboard` | Pipeline dashboard |
| `/employer/modules` | Assessment modules |
| `/employer/modules/new` | Create module |
| `/employer/modules/[id]` | Edit module |
| `/employer/modules/[id]/questions` | Module questions |
| `/employer/modules/library` | Module library |
| `/employer/modules/library/problems` | Problem library |
| `/employer/interviews` | Interview scheduling/outcomes |
| `/employer/interviewers` | Interviewer management |
| `/employer/grading` | Grading workspace |
| `/employer/feedback/[id]` | Candidate feedback |
| `/employer/applications/[id]/scorecard` | Application scorecard |
| `/employer/offers` | Offer list |
| `/employer/offers/new` | Create offer |
| `/employer/offers/[id]` | Offer details |
| `/employer/assessments` | Assessment instances |
| `/employer/profile` | Company profile |

## Interviewer workspace

| Path | Purpose |
| --- | --- |
| `/interviewer/interviews` | Assigned interviews |
| `/interviewer/interviews/[id]/scorecard` | Interview scorecard |
| `/interviewer/grading` | Assigned grading |
| `/interviewer/availability` | Availability |
| `/interviewer/profile` | Interviewer profile |

## Administration workspace

| Path | Purpose |
| --- | --- |
| `/admin/dashboard` | Company registrations |
| `/admin/accounts` | Account administration |

Route presence is not an authorization guarantee. Protected data and writes
remain subject to backend session and role checks. Update this inventory when
a `page.tsx` route is added, removed, or renamed.