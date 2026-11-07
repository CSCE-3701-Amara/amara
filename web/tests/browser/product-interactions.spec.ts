import { test, expect } from "@playwright/test";
import {
  ids,
  moduleDefinition,
  productFixtures,
} from "./support/product-fixtures";
test("DAG joins remain explicit, cycles are rejected, and nodes support keyboard positioning", async ({
  page,
}) => {
  await productFixtures(page);
  await page.goto(`/employer/pipelines/${ids.job}`);
  const node = page.getByRole("button", { name: /^English check, 0 incoming/ });
  await expect(node).toBeVisible();
  const before = await node.boundingBox();
  await node.focus();
  await page.keyboard.press("Alt+ArrowRight");
  await expect
    .poll(async () => Math.round((await node.boundingBox())!.x))
    .toBe(Math.round(before!.x + 20));
  await page.getByRole("button", { name: /^Live interview, 3 incoming/ }).click();
  await expect(
    page.getByRole("group", { name: "Interview setup" }),
  ).toBeVisible();
  await expect(
    page.getByText("No assessment modules are available yet."),
  ).toHaveCount(0);
  await page.getByLabel("Connect to stage").selectOption("stage-1");
  await page
    .getByRole("button", { name: "Add connection", exact: true })
    .click();
  await expect(
    page.getByRole("status").filter({ hasText: "would create a cycle" }),
  ).toBeVisible();
  await page.getByRole("button", { name: "Zoom out", exact: true }).click();
  await expect(page.locator("output")).toHaveText("90%");
  await expect(
    page.getByRole("button", { name: "Save graph draft" }),
  ).toBeDisabled();
});

test("creates a pipeline draft through the supported API from the empty state", async ({
  page,
}) => {
  await productFixtures(page);
  let created = false;
  await page.route(`**/v1/jobs/${ids.job}/pipeline`, async (route) => {
    if (route.request().method() === "POST") {
      created = true;
      return route.fulfill({
        status: 201,
        json: {
          id: ids.job,
          postingId: ids.job,
          versionId: ids.version,
          versionNumber: 2,
        },
      });
    }
    if (created) {
      return route.fulfill({
        json: {
          id: ids.job,
          postingId: ids.job,
          versionId: ids.version,
          versionNumber: 2,
        },
      });
    }
    return route.fulfill({
      status: 404,
      json: { detail: "This job does not have a pipeline draft." },
    });
  });
  await page.goto(`/employer/pipelines/${ids.job}`);
  await expect(page.getByRole("heading", { name: "Start a pipeline draft" }))
    .toBeVisible();
  await page.getByRole("button", { name: "Create pipeline draft" }).click();
  await expect(page.getByText("Draft 2", { exact: true })).toBeVisible();
  await expect(
    page.getByRole("button", { name: /^English check, 0 incoming/ }),
  ).toBeVisible();
  expect(created).toBe(true);
});

test("pipeline lifecycle dialogs support deep links and clear on dismissal", async ({
  page,
}) => {
  await productFixtures(page);
  await page.goto(
    `/employer/pipelines/${ids.job}?state=migration-review&dialog=migrate`,
  );
  const dialog = page.getByRole("dialog");
  await expect(
    dialog.getByRole("heading", { name: "Confirm pipeline migration?" }),
  ).toBeVisible();
  await dialog.getByRole("button", { name: "Cancel", exact: true }).click();
  await expect(dialog).toHaveCount(0);
  await expect(page).toHaveURL(
    `/employer/pipelines/${ids.job}?state=migration-review`,
  );
});

test("published view does not present draft data as a published version", async ({
  page,
}) => {
  await productFixtures(page);
  let pipelineReads = 0;
  await page.route(`**/v1/jobs/${ids.job}/pipeline`, async (route) => {
    pipelineReads += 1;
    return route.fulfill({ json: {} });
  });
  await page.goto(`/employer/pipelines/${ids.job}?state=published`);
  await expect(
    page.getByRole("heading", { name: "Published pipeline version" }),
  ).toBeVisible();
  await expect(
    page.getByText(/No draft data is shown as a published pipeline/),
  ).toBeVisible();
  expect(pipelineReads).toBe(0);
});

test("failed lifecycle confirmation stays open and returns focus on cancellation", async ({
  page,
}) => {
  await productFixtures(page);
  await page.goto(`/employer/modules/${ids.module}/questions`);
  const trigger = page.getByRole("button", {
    name: "Publish module",
    exact: true,
  });
  await trigger.click();
  const dialog = page.getByRole("dialog");
  await expect(
    dialog.getByRole("button", { name: "Cancel", exact: true }),
  ).toBeFocused();
  await dialog
    .getByRole("button", { name: "Publish module", exact: true })
    .click();
  await expect(dialog.getByRole("alert")).toContainText(
    "Fixture does not simulate successful writes",
  );
  await expect(dialog).toBeVisible();
  await page.keyboard.press("Escape");
  await expect(dialog).toHaveCount(0);
  await expect(trigger).toBeFocused();
});
test("desktop collapse does not hide the mobile navigation", async ({
  page,
}) => {
  await productFixtures(page);
  await page.goto("/candidate/applications");
  await page.getByRole("button", { name: "Hide menu", exact: true }).click();
  await page.setViewportSize({ width: 375, height: 812 });
  await page.getByRole("button", { name: "Menu", exact: true }).click();
  await expect(
    page
      .getByRole("navigation", { name: "Workspace navigation" })
      .getByRole("link", { name: "Offers", exact: true }),
  ).toBeVisible();
  expect(
    await page.evaluate(() => document.documentElement.scrollWidth),
  ).toBeLessThanOrEqual(375);
});
test("editing a module sends only definition fields and preserves all questions", async ({
  page,
}) => {
  await productFixtures(page);
  let saved: Record<string, unknown> | undefined;
  await page.route(`**/v1/modules/${ids.module}`, async (route) => {
    if (route.request().method() === "PUT") {
      saved = route.request().postDataJSON();
      return route.fulfill({ json: { ...moduleDefinition, ...saved } });
    }
    return route.fulfill({ json: moduleDefinition });
  });
  await page.goto(`/employer/modules/${ids.module}/questions`);
  await page.getByLabel("Module name (required)").fill("Updated module");
  await page.getByRole("button", { name: "Save draft", exact: true }).click();
  await expect.poll(() => saved?.name).toBe("Updated module");
  expect(saved).not.toHaveProperty("id");
  expect(saved).not.toHaveProperty("company_id");
  expect(saved).not.toHaveProperty("published");
  expect(saved).not.toHaveProperty("threshold");
  expect(saved?.questions).toEqual(moduleDefinition.questions);
});
