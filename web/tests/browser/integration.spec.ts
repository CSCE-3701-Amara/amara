import { expect, test, type Page } from "@playwright/test";
const jobId = "11111111-1111-4111-8111-111111111111";
const moduleId = "22222222-2222-4222-8222-222222222222";
const versionId = "33333333-3333-4333-8333-333333333333";
const applicationId = "44444444-4444-4444-8444-444444444444";
const job = {
  id: jobId,
  company_id: jobId,
  title: "Frontend Engineer",
  description: "Build accessible products.",
  location: "Cairo",
  seniority_level: "junior",
  employment_type: "remote full-time",
  salary_visible: false,
  salary_min: 100,
  salary_max: 200,
  currency: "EGP",
  status: "draft",
  requirements: [],
  amarat: [{ title: "React" }],
};
const assessment = {
  id: moduleId,
  name: "Teamwork",
  type: "behavioral",
  published: false,
  company_id: jobId,
  questions: [{ prompt: "How do you collaborate?", format: "short_text" }],
};
async function fillLogin(page: Page) {
  await page.getByLabel("Email", { exact: true }).fill("candidate@example.com");
  await page.getByLabel("Password", { exact: true }).fill("test-password");
  await page.getByLabel("Password", { exact: true }).press("Enter");
}

test("sign-in cookie reaches the candidate feed and logout handles an empty 204", async ({
  page,
}) => {
  await page.route("**/v1/auth/sessions", (route) =>
    route.fulfill({
      json: {
        user: { id: jobId, email: "candidate@example.com", role: "candidate" },
        session: { expiresAt: "2027-01-01T00:00:00Z" },
      },
      headers: {
        "Set-Cookie":
          "amara_session=test-session; Path=/; HttpOnly; SameSite=Lax",
      },
    }),
  );
  await page.route("**/v1/jobs/feed*", async (route) => {
    expect((await route.request().allHeaders()).cookie).toContain(
      "amara_session=test-session",
    );
    await route.fulfill({
      json: { items: [{ ...job, status: "published" }], next_cursor: "" },
    });
  });
  let logout = false;
  await page.route("**/v1/auth/sessions/current", async (route) => {
    expect(route.request().method()).toBe("DELETE");
    expect((await route.request().allHeaders()).cookie).toContain(
      "amara_session=test-session",
    );
    logout = true;
    await route.fulfill({
      status: 204,
      headers: { "Set-Cookie": "amara_session=; Path=/; Max-Age=0; HttpOnly" },
    });
  });
  await page.goto("/login");
  await fillLogin(page);
  await expect(
    page.getByRole("heading", { name: "Frontend Engineer" }),
  ).toBeVisible();
  expect(await page.evaluate(() => document.cookie)).not.toContain(
    "amara_session",
  );
  await page.getByRole("button", { name: "Sign out", exact: true }).click();
  await expect(page).toHaveURL(/\/login$/);
  expect(logout).toBe(true);
  expect(
    (await page.context().cookies()).find(
      (cookie) => cookie.name === "amara_session",
    ),
  ).toBeUndefined();
});

test("feed filters use handler query names and cursor pagination preserves filters", async ({
  page,
}) => {
  const queries: URLSearchParams[] = [];
  const cursor = "2026-01-01T00:00:00Z|" + jobId;
  await page.route("**/v1/jobs/feed*", (route) => {
    const params = new URL(route.request().url()).searchParams;
    queries.push(params);
    return route.fulfill({
      json: {
        items: [job],
        next_cursor: params.get("cursor") ? undefined : cursor,
      },
    });
  });
  await page.goto("/candidate/dashboard");
  await expect(
    page.getByRole("heading", { name: "Frontend Engineer" }),
  ).toBeVisible();
  await page
    .getByRole("searchbox", { name: "Search job titles" })
    .fill("C++ & React");
  await page.getByLabel("Employment type").selectOption("remote full-time");
  await page.getByRole("button", { name: "Search jobs", exact: true }).click();
  await expect.poll(() => queries.at(-1)?.get("title")).toBe("C++ & React");
  await expect(
    page.getByRole("heading", { name: "Frontend Engineer" }),
  ).toBeVisible();
  await page
    .getByRole("button", { name: "Load more opportunities", exact: true })
    .click();
  await expect.poll(() => queries.at(-1)?.get("cursor")).toBe(cursor);
  expect(queries.at(-1)?.get("employment_type")).toBe("remote full-time");
});

test("apply displays backend rejection and only reports submission after success", async ({
  page,
}) => {
  await page.route("**/v1/jobs/feed*", (route) =>
    route.fulfill({ json: { items: [job] } }),
  );
  let calls = 0;
  await page.route(`**/v1/jobs/${jobId}/apply`, (route) => {
    expect(route.request().method()).toBe("POST");
    expect(route.request().postData()).toBeNull();
    return ++calls === 1
      ? route.fulfill({
          status: 401,
          json: {
            error: {
              code: "authentication_required",
              message: "Candidate authentication is required.",
            },
          },
        })
      : route.fulfill({
          status: 201,
          json: {
            id: applicationId,
            reference_number: 41,
            posting_id: jobId,
            status: "screened_out",
            screened_out: true,
            failed_amarat: ["React"],
          },
        });
  });
  await page.goto("/candidate/dashboard");
  await page.getByRole("link", { name: "View opportunity" }).click();
  await page.getByRole("button", { name: "Apply now" }).click();
  await expect(page.locator("main").getByRole("alert")).toContainText(
    "Candidate authentication is required.",
  );
  await expect(page.getByText("Reference: 41.", { exact: false })).toHaveCount(
    0,
  );
  await page.getByRole("button", { name: "Apply now" }).click();
  await expect(page.getByRole("status")).toContainText(
    "did not pass initial screening",
  );
});

test("application history reads the camelCase data envelope and private details", async ({
  page,
}) => {
  await page.route("**/v1/candidates/me/applications*", (route) =>
    route.fulfill({
      json: {
        data: [
          {
            id: applicationId,
            jobTitle: "Frontend Engineer",
            companyName: "Acme",
            referenceNumber: 41,
            status: "admitted",
          },
        ],
        pagination: { limit: 20, offset: 0, total: 1 },
      },
    }),
  );
  await page.route(`**/v1/applications/${applicationId}`, (route) =>
    route.fulfill({
      json: {
        id: applicationId,
        status: "admitted",
        closed: false,
        finalFeedback: "Thank you for applying.",
      },
    }),
  );
  await page.goto("/candidate/applications");
  await expect(
    page.getByRole("heading", { name: "Frontend Engineer" }),
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "Next page", exact: true }),
  ).toHaveCount(0);
  await page.getByRole("link", { name: "Continue", exact: true }).click();
  await expect(page.getByText("Thank you for applying.")).toBeVisible();
});

test("posting writes use backend field names and lifecycle deletion accepts 204", async ({
  page,
}) => {
  let saved: Record<string, unknown> = {};
  await page.route("**/v1/jobs", async (route) => {
    if (route.request().method() === "POST") {
      saved = route.request().postDataJSON();
      await route.fulfill({ status: 201, json: { ...job, ...saved } });
    } else await route.fulfill({ json: { items: [] } });
  });
  await page.route(`**/v1/jobs/${jobId}`, (route) =>
    route.request().method() === "DELETE"
      ? route.fulfill({ status: 204 })
      : route.fulfill({ json: { ...job, ...saved } }),
  );
  await page.route(`**/v1/jobs/${jobId}/pipeline`, (route) =>
    route.fulfill({
      status: 404,
      json: { detail: "No pipeline", code: "pipeline_not_found" },
    }),
  );
  await page.goto("/employer/jobs/new");
  await page.getByLabel("Job title", { exact: true }).fill("Frontend Engineer");
  await page.getByLabel("Seniority level").fill("junior");
  await page.getByLabel("Employment type").selectOption("remote full-time");
  await page
    .getByLabel("Required skills (amarat), one per line")
    .fill("React\nTypeScript");
  await page.getByRole("button", { name: "Create draft", exact: true }).click();
  await expect(page).toHaveURL(new RegExp(`/employer/jobs/${jobId}$`));
  expect(saved).toMatchObject({
    seniority_level: "junior",
    employment_type: "remote full-time",
    amarat: [{ title: "React" }, { title: "TypeScript" }],
  });
  expect(saved).not.toHaveProperty("company_id");
  expect(saved).not.toHaveProperty("seniorityLevel");
  page.on("dialog", (dialog) => dialog.accept());
  await page
    .getByRole("button", { name: "Delete posting", exact: true })
    .click();
  await page
    .getByRole("dialog")
    .getByRole("button", { name: "Delete posting", exact: true })
    .click();
  await expect(page).toHaveURL(/\/employer\/jobs$/);
});

test("failed company requests stay recoverable and do not pretend the session expired", async ({
  page,
}) => {
  let recovered = false;
  await page.route("**/v1/jobs", (route) =>
    !recovered
      ? route.fulfill({
          status: 401,
          json: {
            error: {
              code: "authentication_required",
              message: "Company manager authentication is required.",
            },
          },
        })
      : route.fulfill({ json: { items: [job] } }),
  );
  await page.goto("/employer/jobs");
  await expect(page.locator("main").getByRole("alert")).toContainText(
    "Company manager authentication is required.",
  );
  await expect(page).toHaveURL(/\/employer\/jobs$/);
  recovered = true;
  await page.getByRole("button", { name: "Try again", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "Frontend Engineer" }),
  ).toBeVisible();
});

test("assessment forms use bare array responses and publish/copy existing modules", async ({
  page,
}) => {
  let created = false;
  let published = false;
  let copied = false;
  await page.route("**/v1/modules", (route) => {
    if (route.request().method() === "POST") {
      expect(route.request().postDataJSON()).toMatchObject({
        name: "Teamwork",
        type: "behavioral",
        questions: [
          { format: "short_text", prompt: "How do you collaborate?" },
        ],
      });
      created = true;
      return route.fulfill({ status: 201, json: assessment });
    }
    return route.fulfill({
      json: created ? [{ ...assessment, published }] : [],
    });
  });
  await page.route("**/v1/modules/library", (route) =>
    route.fulfill({ json: [{ ...assessment, published: true }] }),
  );
  await page.route(`**/v1/modules/${moduleId}/publish`, (route) => {
    published = true;
    return route.fulfill({ json: { ...assessment, published: true } });
  });
  await page.route(`**/v1/modules/${moduleId}/copy`, (route) => {
    copied = true;
    return route.fulfill({ status: 201, json: assessment });
  });
  await page.route(`**/v1/modules/${moduleId}`, (route) =>
    route.fulfill({ json: { ...assessment, published } }),
  );
  await page.goto("/employer/modules/new");
  await page.getByLabel("Module name (required)").fill("Teamwork");
  await page.getByLabel("Prompt (required)").fill("How do you collaborate?");
  await page.getByRole("button", { name: "Save draft", exact: true }).click();
  await expect(page).toHaveURL(
    new RegExp(`/employer/modules/${moduleId}/questions$`),
  );
  await page
    .getByRole("button", { name: "Publish module", exact: true })
    .click();
  await page
    .getByRole("dialog")
    .getByRole("button", { name: "Publish module", exact: true })
    .click();
  await expect.poll(() => published).toBe(true);
  await page.goto("/employer/modules/library");
  await page.getByRole("button", { name: "Teamwork", exact: true }).click();
  await page.getByRole("button", { name: "Copy module", exact: true }).click();
  await page
    .getByRole("dialog")
    .getByRole("button", { name: "Copy module", exact: true })
    .click();
  await expect.poll(() => copied).toBe(true);
});

test("pipeline saves an ordered stage array with version id and RFC3339 times", async ({
  page,
}) => {
  const draft = { id: jobId, postingId: jobId, versionId, versionNumber: 1 };
  let created = false;
  let saved: unknown;
  await page.route(`**/v1/jobs/${jobId}`, (route) =>
    route.fulfill({ json: job }),
  );
  await page.route(`**/v1/jobs/${jobId}/pipeline`, (route) => {
    if (route.request().method() === "POST") {
      created = true;
      return route.fulfill({ status: 201, json: draft });
    }
    return created
      ? route.fulfill({ json: draft })
      : route.fulfill({
          status: 404,
          json: { code: "pipeline_not_found", detail: "No pipeline" },
        });
  });
  await page.route(`**/v1/pipeline-versions/${versionId}/stages`, (route) => {
    if (route.request().method() === "PUT")
      saved = route.request().postDataJSON();
    return route.fulfill({ json: saved || [] });
  });
  await page.route("**/v1/modules", (route) =>
    route.fulfill({ json: [assessment] }),
  );
  await page.goto(`/employer/pipelines/${jobId}`);
  await page.getByText("Legacy stage editor", { exact: true }).click();
  await page.getByRole("button", { name: "Create pipeline draft" }).click();
  await page.getByRole("button", { name: "Add assessment stage" }).click();
  await page.getByLabel("Stage name").fill("Skills assessment");
  await page.getByLabel("Stage opens").fill("2027-01-01T09:00");
  await page.getByLabel("Stage deadline").fill("2027-01-02T09:00");
  await page.getByRole("checkbox", { name: "Teamwork" }).check();
  await page.getByRole("button", { name: "Save pipeline stages" }).click();
  await expect(
    page
      .getByRole("form", { name: "Pipeline stages" })
      .getByRole("status")
      .filter({ hasText: "Pipeline stages saved." }),
  ).toContainText("Pipeline stages saved.");
  expect(saved).toEqual([
    expect.objectContaining({
      name: "Skills assessment",
      type: "assessment",
      moduleIds: [moduleId],
      scheduledStartAt: expect.stringMatching(/Z$/),
      scheduledEndAt: expect.stringMatching(/Z$/),
    }),
  ]);
});

for (const width of [375, 1280])
  test(`connected workspaces fit at ${width}px`, async ({ page }) => {
    await page.setViewportSize({ width, height: 900 });
    await page.route("**/v1/**", (route) =>
      route.fulfill({
        status: 403,
        json: {
          detail: "This account cannot access the requested resource.",
          code: "forbidden",
        },
      }),
    );
    for (const path of [
      "/candidate/dashboard",
      "/candidate/applications",
      "/employer/jobs",
      "/employer/jobs/new",
      "/employer/pipelines",
    ]) {
      await page.goto(path);
      await expect(page.locator("h1")).toBeVisible();
      expect(
        await page.evaluate(() => document.documentElement.scrollWidth),
      ).toBeLessThanOrEqual(width);
    }
  });

test("non-JSON proxy errors stay readable and can be retried", async ({
  page,
}) => {
  let recovered = false;
  await page.route("**/v1/jobs/feed*", (route) =>
    recovered
      ? route.fulfill({ json: { items: [job] } })
      : route.fulfill({
          status: 502,
          contentType: "text/html",
          body: "<html>Internal proxy implementation details</html>",
        }),
  );
  await page.goto("/candidate/dashboard");
  await expect(page.locator("main").getByRole("alert")).toContainText(
    "We couldn’t complete your request. Please try again.",
  );
  await expect(
    page.getByText("Internal proxy implementation details"),
  ).toHaveCount(0);
  recovered = true;
  await page.getByRole("button", { name: "Try again", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "Frontend Engineer" }),
  ).toBeVisible();
});
