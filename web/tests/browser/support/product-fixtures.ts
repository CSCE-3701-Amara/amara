import type { Page } from "@playwright/test";
export const ids = {
  job: "11111111-1111-4111-8111-111111111111",
  module: "22222222-2222-4222-8222-222222222222",
  version: "33333333-3333-4333-8333-333333333333",
  application: "44444444-4444-4444-8444-444444444444",
  interview: "55555555-5555-4555-8555-555555555555",
};
export const posting = {
  id: ids.job,
  company_id: ids.job,
  title: "Frontend Engineer",
  description:
    "Build accessible product experiences and work with engineering to bring them to life.",
  location: "Cairo, Egypt",
  seniority_level: "Senior",
  employment_type: "remote full-time",
  salary_visible: true,
  salary_min: 30000,
  salary_max: 40000,
  currency: "EGP",
  open_positions: 2,
  max_applications: 100,
  status: "draft",
  requirements: [
    {
      text: "Experience designing accessible digital products.",
      requirement_type: "must_have",
    },
  ],
  amarat: [
    {
      title: "Product design",
      description: "Evidence of product design skills.",
    },
  ],
  application_opening: "2026-11-10T07:00:00Z",
  application_deadline: "2026-12-10T15:00:00Z",
};
export const moduleDefinition = {
  id: ids.module,
  company_id: ids.job,
  name: "Go fundamentals",
  type: "behavioral",
  description: "Question-based assessment",
  active_minutes: 30,
  threshold: 70,
  published: false,
  questions: [
    {
      prompt: "Explain when a buffered channel blocks.",
      format: "short_text",
      amara: "Go concurrency",
      difficulty: "medium",
    },
    {
      prompt: "Which synchronization primitive protects shared state?",
      format: "multiple_choice",
      options: ["Mutex", "Channel", "WaitGroup", "Context"],
    },
    {
      prompt: "Implement a worker pool with bounded concurrency.",
      format: "code",
      test_cases: [{ input: "1", expected_output: "1", weight: 1 }],
    },
  ],
};
const pagination = { limit: 20, offset: 0, total: 1 };
export const application = {
  id: ids.application,
  postingId: ids.job,
  referenceNumber: 1042,
  jobTitle: "Frontend Engineer",
  companyName: "Acme Corp",
  status: "in_progress",
  closed: false,
  submittedAt: "2026-11-03T09:00:00Z",
  screeningCompletedAt: "2026-11-04T09:00:00Z",
  screenedOut: false,
  updatedAt: "2026-11-04T09:00:00Z",
};
export const interview = {
  id: ids.interview,
  stageExecutionId: ids.application,
  scheduledStartAt: "2026-11-12T11:00:00Z",
  scheduledEndAt: "2026-11-12T11:45:00Z",
  status: "confirmed",
  rescheduleCount: 0,
  createdAt: "2026-11-01T09:00:00Z",
  joinUrl: "https://example.com/meeting",
};
export const offer = {
  id: ids.application,
  applicationId: ids.application,
  jobTitle: "Frontend Engineer",
  compensationSummary: "30,000 EGP per month",
  proposedStartDate: "2026-12-01",
  responseDeadline: "2026-11-20T08:00:00Z",
  message: "We are excited to welcome you to the team.",
  status: "pending",
  createdAt: "2026-11-10T08:00:00Z",
};
export const stages = [
  {
    id: "stage-1",
    name: "English check",
    type: "assessment",
    scheduledStartAt: "2026-11-05T07:00:00Z",
    scheduledEndAt: "2026-11-08T07:00:00Z",
    moduleIds: [ids.module],
  },
  {
    id: "stage-2",
    name: "Coding challenge",
    type: "assessment",
    scheduledStartAt: "2026-11-05T07:00:00Z",
    scheduledEndAt: "2026-11-08T07:00:00Z",
    moduleIds: [ids.module],
  },
  {
    id: "stage-3",
    name: "Behavioral response",
    type: "assessment",
    scheduledStartAt: "2026-11-05T07:00:00Z",
    scheduledEndAt: "2026-11-08T07:00:00Z",
    moduleIds: [ids.module],
  },
  {
    id: "stage-4",
    name: "Live interview",
    type: "live_interview",
    scheduledStartAt: "2026-11-10T07:00:00Z",
    scheduledEndAt: "2026-11-13T15:00:00Z",
  },
];
const scorecard = {
  applicationId: ids.application,
  anonymousIdentifier: "042",
  meanScore: 82,
  stages: stages.map((s, i) => ({
    stageId: s.id,
    stageName: s.name,
    status: i ? "in_progress" : "passed",
    score: i ? null : 82,
  })),
  insightsAreSystemGenerated: false,
};
/** Browser-only fixtures. No fixture data is shipped in product routes or written to the API. */
export async function productFixtures(page: Page) {
  await page.route("**/v1/**", async (route) => {
    const request = route.request();
    const u = new URL(request.url());
    const p = u.pathname;
    const pageState = new URL(page.url()).searchParams.get("state");
    if (request.method() !== "GET")
      return route.fulfill({
        status: 501,
        json: { detail: "Fixture does not simulate successful writes" },
      });
    let body: unknown;
    if (p === "/v1/jobs/feed")
      body = {
        items: [{ ...posting, company_name: "Acme Corp", status: "published" }],
        next_cursor: "",
      };
    else if (p === "/v1/jobs")
      body = {
        items: [
          {
            ...posting,
            status: pageState === "published" ? "published" : "draft",
          },
        ],
      };
    else if (p.endsWith("/pipeline"))
      body = {
        id: ids.job,
        postingId: ids.job,
        versionId: ids.version,
        versionNumber: 2,
      };
    else if (p.includes("/pipeline-versions/") && p.endsWith("/stages"))
      body = stages;
    else if (p.endsWith("/transitions"))
      body = [
        {
          id: "edge-1",
          fromStageId: "stage-1",
          toStageId: "stage-4",
          pipelineVersionId: ids.version,
        },
        {
          id: "edge-2",
          fromStageId: "stage-2",
          toStageId: "stage-4",
          pipelineVersionId: ids.version,
        },
        {
          id: "edge-3",
          fromStageId: "stage-3",
          toStageId: "stage-4",
          pipelineVersionId: ids.version,
        },
      ];
    else if (p === "/v1/candidates/me/applications")
      body = { data: pageState === "empty" ? [] : [application], pagination };
    else if (p.endsWith("/stage-progress"))
      body = stages.map((s, i) => ({
        id: s.id,
        stageId: s.id,
        name: s.name,
        status: ["passed", "in_progress", "available", "locked"][i],
      }));
    else if (p.endsWith("/scorecard")) body = scorecard;
    else if (p.endsWith("/finalists")) body = { data: [scorecard], pagination };
    else if (p.endsWith("/comparison"))
      body = [
        scorecard,
        {
          ...scorecard,
          applicationId: ids.job,
          anonymousIdentifier: "057",
          meanScore: 78,
        },
      ];
    else if (p.includes("/pipelines/") && p.endsWith("/dashboard"))
      body = {
        pipelineId: ids.job,
        publishedVersionId: ids.version,
        activeCandidates: 40,
        pendingOffers: 2,
        acceptedOffers: 1,
        stages: stages.map((s) => ({
          stageId: s.id,
          name: s.name,
          activeCount: 8,
          passedCount: 32,
          rejectedCount: 1,
        })),
      };
    else if (p.includes("/pipeline-stages/") && p.endsWith("/executions"))
      body = {
        data: [
          {
            id: ids.application,
            pipelineExecutionId: ids.application,
            stageId: "stage-1",
            status: "held",
            score: 80,
            decisionHistory: [],
          },
        ],
        pagination,
      };
    else if (p.match(/^\/v1\/jobs\/[^/]+\/applications$/))
      body = {
        data: [
          {
            application,
            anonymousIdentifier: "042",
            screeningSummary: "Admitted",
          },
        ],
        pagination,
      };
    else if (p.match(/^\/v1\/jobs\/[^/]+$/))
      body = {
        ...posting,
        status: pageState === "published" ? "published" : "draft",
      };
    else if (p.match(/^\/v1\/applications\/[^/]+$/)) body = application;
    else if (p === "/v1/modules" || p === "/v1/modules/library")
      body = [{ ...moduleDefinition, published: p.endsWith("library") }];
    else if (p === "/v1/modules/library/problems")
      body = [
        {
          module_id: ids.module,
          position: 0,
          difficulty: "medium",
          prompt: "Implement a bounded worker pool.",
        },
      ];
    else if (p.match(/^\/v1\/modules\/[^/]+$/)) body = moduleDefinition;
    else if (p === "/v1/offers") body = { data: [offer], pagination };
    else if (p.match(/^\/v1\/offers\/[^/]+$/)) body = offer;
    else if (p === "/v1/interviews") body = { data: [interview], pagination };
    else if (p.endsWith("/availability"))
      body = [
        { startAt: "2026-11-12T11:00:00Z", endAt: "2026-11-12T11:45:00Z" },
        { startAt: "2026-11-13T12:00:00Z", endAt: "2026-11-13T12:45:00Z" },
      ];
    else if (p.endsWith("/feedback/me"))
      body = {
        id: ids.interview,
        interviewId: ids.interview,
        interviewerId: ids.job,
        ratings: {
          "Technical reasoning": 4,
          "Frontend implementation": 4,
          Communication: 3,
          Collaboration: 3,
        },
        notes: "Private interview notes",
        recommendation: "Evidence supports advancement.",
        createdAt: "2026-11-12T12:00:00Z",
      };
    else if (p.match(/^\/v1\/interviews\/[^/]+$/)) body = interview;
    else if (p === "/v1/grading/submissions")
      body = {
        data: [
          {
            id: ids.application,
            attemptId: ids.application,
            questionId: ids.module,
            anonymousIdentifier: "042",
            prompt: "When does a buffered channel block?",
            response:
              "A send blocks when the buffer is full until a receiver makes room.",
            maxPoints: 10,
            gradingStatus: "pending",
          },
        ],
        pagination,
      };
    else if (p.includes("/assessment-attempts/") && p.endsWith("/submissions"))
      body = [];
    else if (p.match(/^\/v1\/assessment-attempts\/[^/]+$/))
      body = {
        id: ids.application,
        instanceId: ids.application,
        stageExecutionId: ids.application,
        status: pageState === "submitted" ? "submitted" : "in_progress",
        deadline: new Date(Date.now() + 20 * 60000).toISOString(),
        autoSubmitted: false,
        questions: [
          {
            id: ids.module,
            prompt: "Explain when a buffered channel blocks.",
            responseFormat: "short_text",
            position: 0,
          },
          {
            id: ids.job,
            prompt: "Choose a synchronization primitive.",
            responseFormat: "multiple_choice",
            options: ["Mutex", "Channel", "Context"],
            position: 1,
          },
          {
            id: ids.version,
            prompt: "Implement a bounded worker pool.",
            responseFormat: "code",
            position: 2,
          },
        ],
      };
    else if (p.endsWith("/result"))
      body = {
        id: ids.application,
        instanceId: ids.application,
        passed: pageState !== "failed",
        score: 82,
        feedback: "Clear reasoning and strong implementation.",
        feedbackIsAiGenerated: false,
        disclosedAt: "2026-11-12T12:00:00Z",
        scoredAt: "2026-11-12T12:00:00Z",
      };
    else if (p.match(/^\/v1\/assessment-instances\/[^/]+$/))
      body = {
        id: ids.application,
        applicationId: ids.application,
        stageId: "stage-1",
        status: "available",
      };
    else if (p.endsWith("/rejection-feedback") || p.endsWith("/feedback"))
      body = {
        feedback:
          "Your implementation covered the main requirements. Add detail about edge cases and trade-offs.",
        feedbackIsAiGenerated: false,
      };
    else if (p === "/v1/companies/me")
      body = {
        id: ids.job,
        legalName: "Acme Corp",
        displayName: "Acme",
        approved: true,
        contactEmail: "hello@example.com",
        description:
          "We build thoughtful products with a collaborative team based in Egypt.",
      };
    else if (p === "/v1/companies/me/interviewers")
      body = {
        data: [
          {
            id: ids.job,
            companyId: ids.job,
            name: "Daniel Ruiz",
            email: "daniel@example.com",
            accountStatus: "active",
          },
          {
            id: ids.module,
            companyId: ids.job,
            name: "Priya Shah",
            email: "priya@example.com",
            accountStatus: "active",
          },
        ],
        pagination: { ...pagination, total: 2 },
      };
    else if (p === "/v1/interviewers/me")
      body = {
        id: ids.job,
        companyId: ids.job,
        name: "Daniel Ruiz",
        email: "daniel@example.com",
        accountStatus: "active",
      };
    else if (p === "/v1/admin/accounts")
      body = {
        data: [
          {
            id: ids.job,
            accountType: "company",
            name: "Acme Retail",
            email: "company@example.com",
            accountStatus: "active",
            approved: true,
          },
        ],
        pagination,
      };
    else if (p === "/v1/admin/company-registration-requests")
      body = {
        data: [
          {
            id: ids.job,
            companyName: "Acme Retail",
            registrationNumber: "1042",
            email: "company@example.com",
            status: "pending_review",
            createdAt: "2026-11-01T09:00:00Z",
          },
        ],
        pagination,
      };
    if (pageState === "error")
      return route.fulfill({
        status: 503,
        json: { detail: "Service temporarily unavailable" },
      });
    return body === undefined
      ? route.fulfill({
          status: 404,
          json: { detail: "Operation unavailable" },
        })
      : route.fulfill({ json: body });
  });
}
