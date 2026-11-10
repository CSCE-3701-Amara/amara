import { expect, test } from "@playwright/test";
import { ids, productFixtures } from "./support/product-fixtures";

test("direct assessment confirmation keeps the submitted destination", async ({
  page,
}) => {
  await productFixtures(page);
  let submitted = false;
  await page.route(`**/v1/assessment-attempts/${ids.application}`, (route) =>
    route.fulfill({
      json: {
        id: ids.application,
        instanceId: ids.application,
        stageExecutionId: ids.application,
        status: submitted ? "submitted" : "in_progress",
        deadline: "2027-01-01T00:00:00Z",
        autoSubmitted: false,
        questions: [],
      },
    }),
  );
  await page.route(
    `**/v1/assessment-attempts/${ids.application}/submission`,
    (route) => {
      submitted = true;
      return route.fulfill({
        json: { id: ids.application, status: "submitted" },
      });
    },
  );
  await page.goto(`/candidate/assessments/${ids.application}?dialog=submit`);
  await page
    .getByRole("dialog")
    .getByRole("button", { name: "Submit assessment", exact: true })
    .click();
  await expect(page).toHaveURL(/state=submitted/);
  await expect(page.getByRole("dialog")).toHaveCount(0);
  await expect(
    page.getByRole("heading", {
      name: "Your answers are saved and locked",
      exact: true,
    }),
  ).toBeVisible();
});

test("direct account confirmation keeps the signed-out destination", async ({
  page,
}) => {
  await productFixtures(page);
  await page.route("**/v1/auth/account-actions", (route) =>
    route.fulfill({ status: 204 }),
  );
  await page.goto("/candidate/profile/edit?dialog=deactivate");
  await page
    .getByRole("dialog")
    .getByRole("button", { name: "Deactivate", exact: true })
    .click();
  await expect(page).toHaveURL(/\/login$/);
});

test("company profile patches only edited fields and preserves values after failure", async ({
  page,
}) => {
  await productFixtures(page);
  let body: unknown;
  let writes = 0;
  await page.route("**/v1/companies/me", async (route) => {
    if (route.request().method() === "GET")
      return route.fulfill({
        status: 404,
        json: { detail: "Profile reads are not registered" },
      });
    body = route.request().postDataJSON();
    writes++;
    return writes === 1
      ? route.fulfill({
          status: 422,
          json: { detail: "Check your company details" },
        })
      : route.fulfill({
          json: {
            id: ids.job,
            legalName: "Existing legal name",
            displayName: "Updated display name",
            approved: true,
          },
        });
  });
  await page.goto("/employer/profile");
  await page.getByRole("link", { name: "Edit profile", exact: true }).click();
  await page
    .getByLabel("Display name", { exact: true })
    .fill("Updated display name");
  await page
    .getByRole("button", { name: "Save company profile", exact: true })
    .click();
  await expect(page.locator("main").getByRole("alert")).toContainText(
    "Check your company details",
  );
  await expect(page.getByLabel("Display name", { exact: true })).toHaveValue(
    "Updated display name",
  );
  expect(body).toEqual({ displayName: "Updated display name" });
  await page
    .getByRole("button", { name: "Save company profile", exact: true })
    .click();
  await expect(page.getByRole("status")).toHaveText("Profile saved.");
  expect(writes).toBe(2);
});

test("candidate profile updates do not send unloaded identity or empty URLs", async ({
  page,
}) => {
  await productFixtures(page);
  let body: unknown;
  await page.route("**/v1/candidates/me", async (route) => {
    if (route.request().method() === "GET")
      return route.fulfill({
        json: {
          id: ids.application,
          name: "Existing candidate",
          available: true,
        },
      });
    expect(route.request().method()).toBe("PATCH");
    body = route.request().postDataJSON();
    return route.fulfill({
      json: { name: "Existing candidate", phoneNumber: "+201234567890" },
    });
  });
  await page.goto("/candidate/profile/edit");
  await page.getByLabel("Phone number", { exact: true }).fill("+201234567890");
  await page.getByRole("button", { name: "Save changes", exact: true }).click();
  await expect(page.getByRole("status")).toHaveText("Profile saved.");
  expect(body).toEqual({ phoneNumber: "+201234567890" });
});

test("assessment saves answers before submitting and locks the returned receipt", async ({
  page,
}) => {
  await productFixtures(page);
  const writes: string[] = [];
  await page.route(
    `**/v1/assessment-attempts/${ids.application}/submissions/*`,
    async (route) => {
      expect(route.request().method()).toBe("PUT");
      writes.push("answer");
      expect(route.request().postDataJSON()).toEqual({
        response: "A buffered send blocks when full.",
      });
      return route.fulfill({
        json: {
          id: ids.module,
          attemptId: ids.application,
          questionId: ids.module,
          response: "A buffered send blocks when full.",
          savedAt: "2026-11-07T12:00:00Z",
        },
      });
    },
  );
  await page.route(
    `**/v1/assessment-attempts/${ids.application}/submission`,
    async (route) => {
      expect(route.request().method()).toBe("POST");
      writes.push("submit");
      return route.fulfill({
        json: { id: ids.application, status: "submitted" },
      });
    },
  );
  await page.goto(`/candidate/assessments/${ids.application}`);
  await page
    .getByLabel("Your answer", { exact: true })
    .first()
    .fill("A buffered send blocks when full.");
  await page
    .getByRole("button", { name: "Submit assessment", exact: true })
    .click();
  await page
    .getByRole("dialog")
    .getByRole("button", { name: "Submit assessment", exact: true })
    .click();
  await expect(page).toHaveURL(/state=submitted/);
  await expect(
    page.getByRole("heading", { name: "Your answers are saved and locked" }),
  ).toBeVisible();
  expect(writes).toEqual(["answer", "submit"]);
});

test("assessment save failures preserve the draft and expose session recovery", async ({
  page,
}) => {
  await productFixtures(page);
  await page.route(
    `**/v1/assessment-attempts/${ids.application}/submissions/*`,
    (route) =>
      route.fulfill({
        status: 401,
        json: {
          error: { code: "session_expired", message: "Your session expired." },
        },
      }),
  );
  await page.goto(`/candidate/assessments/${ids.application}`);
  const answer = page.getByLabel("Your answer", { exact: true }).first();
  await answer.fill("Keep this answer");
  await page
    .getByRole("button", { name: "Save progress", exact: true })
    .click();
  await expect(page.locator("main").getByRole("alert")).toContainText(
    "Your session expired.",
  );
  await expect(
    page.getByRole("link", { name: "Sign in again", exact: true }),
  ).toBeVisible();
  await expect(answer).toHaveValue("Keep this answer");
});

test("assessment start follows the returned attempt identity", async ({
  page,
}) => {
  await productFixtures(page);
  await page.route(`**/v1/assessment-instances/${ids.job}/start`, (route) =>
    route.fulfill({
      json: { id: ids.application, instanceId: ids.job, status: "in_progress" },
    }),
  );
  await page.goto(`/candidate/assessments/${ids.job}?state=availability`);
  await page
    .getByRole("button", { name: "Start assessment", exact: true })
    .click();
  await expect(page).toHaveURL(new RegExp(`attempt=${ids.application}`));
  await expect(
    page.getByRole("heading", { name: "Skills assessment", exact: true }),
  ).toBeVisible();
});

test("manual stage decisions include the observed status and retain conflicts", async ({
  page,
}) => {
  await productFixtures(page);
  let body: unknown;
  await page.route(`**/v1/jobs/${ids.job}/pipeline/dashboard`, (route) =>
    route.fulfill({
      json: {
        active: 1,
        finalistTarget: 1,
        links: [],
        stages: [
          {
            id: ids.version,
            name: "Manual review",
            rule: "manual",
            active: 1,
            passed: 0,
            rejected: 0,
            candidates: [
              {
                id: ids.application,
                candidate: "Candidate 042",
                status: "held",
                active: true,
                score: 80,
              },
            ],
          },
        ],
      },
    }),
  );
  await page.route(
    `**/v1/stage-executions/${ids.application}/decisions`,
    async (route) => {
      body = route.request().postDataJSON();
      return route.fulfill({
        status: 409,
        json: {
          error: {
            code: "stage_decision_conflict",
            message: "The candidate state changed. Refresh before deciding.",
          },
        },
      });
    },
  );
  await page.goto(
    `/employer/pipelines/${ids.job}/candidates?stage=${ids.version}`,
  );
  await page
    .getByRole("checkbox", { name: "Candidate 042", exact: true })
    .check();
  await page.getByRole("button", { name: "Pass", exact: true }).click();
  await page.getByLabel("Decision reason").fill("Meets the rubric");
  await page
    .getByRole("button", { name: "Confirm decision", exact: true })
    .click();
  await expect(page.getByRole("dialog").getByRole("alert")).toContainText(
    "The candidate state changed",
  );
  expect(body).toEqual({
    decision: "pass",
    decisionType: "manual",
    reason: "Meets the rubric",
    expectedStatus: "held",
  });
  await page.keyboard.press("Escape");
  await expect(
    page.getByRole("button", { name: "Pass", exact: true }),
  ).toBeFocused();
});
