import { expect, test } from "@playwright/test";
import { workspaceBackLink } from "../../src/shared/lib/workspace-back-link";

const cases = [
  ["/candidate/jobs/job-id", "/candidate/dashboard", "opportunities"],
  ["/candidate/applications/application-id/progress", "/candidate/applications/application-id", "application"],
  ["/candidate/assessments/assessment-id", "/candidate/applications", "applications"],
  ["/candidate/assessments/assessment-id/coding", "/candidate/assessments/assessment-id", "assessment"],
  ["/candidate/assessments/assessment-id/video", "/candidate/assessments/assessment-id", "assessment"],
  ["/candidate/interviews/interview-id", "/candidate/applications", "applications"],
  ["/candidate/profile/edit", "/candidate/profile", "profile"],
  ["/candidate/offers/offer-id", "/candidate/offers", "offers"],
  ["/employer/jobs/new", "/employer/jobs", "job postings"],
  ["/employer/jobs/job-id/applications", "/employer/jobs/job-id", "job posting"],
  ["/employer/jobs/job-id/preview", "/employer/jobs/job-id", "job posting"],
  ["/employer/pipelines/pipeline-id/candidates", "/employer/pipelines/pipeline-id", "pipeline"],
  ["/employer/modules/module-id/questions", "/employer/modules/module-id", "module"],
  ["/employer/modules/library/problems", "/employer/modules/library", "module library"],
  ["/employer/applications/application-id/scorecard", "/employer/jobs", "job postings"],
  ["/employer/feedback/feedback-id", "/employer/grading", "grading"],
  ["/interviewer/interviews/interview-id/scorecard", "/interviewer/interviews", "interviews"],
  ["/admin/accounts", "/admin/dashboard", "company registrations"],
] as const;

for (const [path, href, label] of cases) {
  test(`${path} returns to an existing parent screen`, () => {
    const home = path.startsWith("/candidate/") ? "/candidate/dashboard"
      : path.startsWith("/employer/") ? "/employer/jobs"
      : path.startsWith("/interviewer/") ? "/interviewer/interviews"
      : "/admin/dashboard";
    expect(workspaceBackLink(path, home)).toEqual({ href, label: `Back to ${label}` });
  });
}

test("home screens have no back link and section screens return to their role home", () => {
  for (const home of ["/candidate/dashboard", "/employer/jobs", "/interviewer/interviews", "/admin/dashboard"]) {
    expect(workspaceBackLink(home, home)).toBeNull();
  }
  expect(workspaceBackLink("/candidate/offers", "/candidate/dashboard")?.href).toBe("/candidate/dashboard");
  expect(workspaceBackLink("/employer/interviews", "/employer/jobs")?.href).toBe("/employer/jobs");
  expect(workspaceBackLink("/interviewer/profile", "/interviewer/interviews")?.href).toBe("/interviewer/interviews");
});

test("assessment variants retain attempt context without reopening action dialogs", () => {
  expect(workspaceBackLink(
    "/candidate/assessments/instance-id/coding",
    "/candidate/dashboard",
    "attempt=attempt-id&application=application-id&state=result&dialog=submit",
  )).toEqual({
    href: "/candidate/assessments/instance-id?attempt=attempt-id&application=application-id&state=result",
    label: "Back to assessment",
  });
});
