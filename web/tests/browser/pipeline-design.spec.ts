import { expect, test } from "@playwright/test";
import {
  ids,
  moduleDefinition,
  productFixtures,
} from "./support/product-fixtures";

test("legacy stage editor mounts on opening and retains edits when reopened", async ({ page }) => {
  await productFixtures(page);
  await page.goto(`/employer/pipelines/${ids.job}`);
  const editor = page.getByRole("region", { name: "Hiring pipeline" });
  await expect(editor).toHaveCount(0);
  const toggle = page.getByText("Legacy stage editor", { exact: true });
  await toggle.click();
  const name = editor.getByLabel("Stage name", { exact: true }).first();
  await name.fill("Retained stage draft");
  await toggle.click();
  await expect(name).toBeHidden();
  await toggle.click();
  await expect(name).toBeVisible();
  await expect(name).toHaveValue("Retained stage draft");
});

test("candidate profile reads back saved details after reload", async ({
  page,
}) => {
  await productFixtures(page);
  let profile = {
    id: ids.application,
    name: "Aya",
    phoneNumber: "+201111111111",
    biography: "Original biography",
    nationality: "Egyptian",
    available: true,
  };
  await page.route("**/v1/candidates/me", (route) => {
    if (route.request().method() === "PATCH")
      profile = { ...profile, ...route.request().postDataJSON() };
    return route.fulfill({ json: profile });
  });
  await page.goto("/candidate/profile/edit");
  await expect(page.getByLabel("Biography", { exact: true })).toHaveValue(
    "Original biography",
  );
  await expect(page.getByLabel("Nationality", { exact: true })).toHaveValue(
    "Egyptian",
  );
  await page.getByLabel("Phone number", { exact: true }).fill("+201234567890");
  await page.getByRole("button", { name: "Save changes", exact: true }).click();
  await expect(page.getByRole("status")).toHaveText("Profile saved.");
  await page.reload();
  await expect(page.getByLabel("Phone number", { exact: true })).toHaveValue(
    "+201234567890",
  );
  expect(profile.biography).toBe("Original biography");
});

test("threshold preview is required, echoed on apply, and refreshed after saving", async ({
  page,
}) => {
  await productFixtures(page);
  let threshold = 70;
  const writes: unknown[] = [];
  const preview = [{ candidate: ids.application, status: "passed", score: 80 }];
  await page.route(`**/v1/jobs/${ids.job}/pipeline/dashboard`, (route) =>
    route.fulfill({
      json: {
        active: 1,
        finalistTarget: null,
        links: [],
        stages: [
          {
            id: ids.version,
            name: "Coding",
            rule: "threshold",
            threshold,
            active: 1,
            passed: 0,
            rejected: 0,
            candidates: [],
          },
        ],
      },
    }),
  );
  await page.route(
    `**/v1/pipeline-stages/${ids.version}/threshold**`,
    (route) => {
      const body = route.request().postDataJSON();
      writes.push(body);
      if (route.request().method() === "PUT") threshold = body.threshold;
      return route.fulfill({ json: preview });
    },
  );
  await page.goto(`/employer/pipelines/${ids.job}/candidates`);
  const panel = page.getByRole("region", { name: "Threshold review" });
  await expect(
    panel.getByRole("button", { name: "Apply reviewed threshold" }),
  ).toBeDisabled();
  await panel.getByLabel("Proposed pass threshold").fill("75");
  await panel.getByRole("button", { name: "Preview outcomes" }).click();
  await expect(
    panel.getByRole("cell", { name: "passed", exact: true }),
  ).toBeVisible();
  await panel.getByRole("button", { name: "Apply reviewed threshold" }).click();
  await expect(panel.getByLabel("Proposed pass threshold")).toHaveValue("75");
  await expect(panel).toContainText("Current pass threshold: 75");
  expect(writes).toEqual([
    { threshold: 75, expectedThreshold: 70 },
    { threshold: 75, expectedThreshold: 70, preview },
  ]);
});

test("changing stages cancels pending threshold previews", async ({ page }) => {
  await productFixtures(page);
  await page.route(`**/v1/jobs/${ids.job}/pipeline/dashboard`, (route) =>
    route.fulfill({
      json: {
        active: 0,
        finalistTarget: null,
        links: [],
        stages: ["first", "second"].map((id) => ({
          id,
          name: id,
          rule: "threshold",
          threshold: 70,
          active: 0,
          passed: 0,
          rejected: 0,
          candidates: [],
        })),
      },
    }),
  );
  let release!: () => void;
  const pending = new Promise<void>((resolve) => {
    release = resolve;
  });
  await page.route(
    "**/v1/pipeline-stages/first/threshold/preview",
    async (route) => {
      await pending;
      await route
        .fulfill({
          json: [{ candidate: "stale candidate", status: "passed", score: 80 }],
        })
        .catch(() => {});
    },
  );
  await page.goto(`/employer/pipelines/${ids.job}/candidates`);
  await page.getByRole("button", { name: "Preview outcomes" }).click();
  await page.getByLabel("Stage", { exact: true }).selectOption("second");
  release();
  await expect(
    page.getByRole("button", { name: "Apply reviewed threshold" }),
  ).toBeDisabled();
  await expect(page.getByText("stale candidate", { exact: true })).toHaveCount(
    0,
  );
});

test("interview outcome uses the registered operation and retains selection after rejection", async ({
  page,
}) => {
  await productFixtures(page);
  const writes: unknown[] = [];
  await page.route(`**/v1/interviews/${ids.interview}/outcome`, (route) => {
    writes.push(route.request().postDataJSON());
    return writes.length === 1
      ? route.fulfill({
          status: 403,
          json: { detail: "This interview belongs to another team." },
        })
      : route.fulfill({ status: 204 });
  });
  await page.goto("/employer/interviews");
  await page.getByRole("button", { name: ids.interview, exact: true }).click();
  await page.getByLabel("Interview outcome").selectOption("candidate_absent");
  await page
    .getByRole("button", { name: "Review outcome", exact: true })
    .last()
    .click();
  const dialog = page.getByRole("dialog");
  await dialog
    .getByRole("button", { name: "Record outcome", exact: true })
    .click();
  await expect(dialog.getByRole("alert")).toContainText("another team");
  await dialog
    .getByRole("button", { name: "Record outcome", exact: true })
    .click();
  await expect(dialog).toHaveCount(0);
  expect(writes).toEqual([
    { outcome: "candidate_absent" },
    { outcome: "candidate_absent" },
  ]);
});

test("local graph undo restores a removed node with all its connections", async ({
  page,
}) => {
  await productFixtures(page);
  await page.goto(`/employer/pipelines/${ids.job}`);
  const node = page.getByRole("button", {
    name: /^Live interview, 3 incoming/,
  });
  await node.click();
  await page.getByRole("button", { name: "Delete", exact: true }).click();
  await expect(node).toHaveCount(0);
  await page.getByRole("button", { name: "Undo graph change" }).click();
  await expect(node).toBeVisible();
  await page.getByRole("button", { name: "Redo graph change" }).click();
  await expect(node).toHaveCount(0);
  page.on("dialog", (dialog) => dialog.dismiss());
});

test("module editing preserves backend behavioral and test-pass settings", async ({
  page,
}) => {
  await productFixtures(page);
  let saved: Record<string, unknown> | undefined;
  const initial = {
    ...moduleDefinition,
    type: "coding",
    threshold: undefined,
    minimum_tests_passed: 1,
    questions: [
      {
        prompt: "Write code",
        format: "code",
        minimum_cases_passed: 1,
        test_cases: [{ input: "a", expected_output: "a", weight: 1 }],
      },
    ],
  };
  await page.route(`**/v1/modules/${ids.module}`, (route) => {
    if (route.request().method() === "PUT")
      saved = route.request().postDataJSON();
    return route.fulfill({ json: { ...initial, ...saved } });
  });
  await page.goto(`/employer/modules/${ids.module}/questions`);
  await page.getByLabel("Module name (required)").fill("Updated module");
  await page.getByRole("button", { name: "Save draft", exact: true }).click();
  await expect.poll(() => saved?.minimum_tests_passed).toBe(1);
  expect(saved?.questions).toEqual(initial.questions);
});
