import { expect, test } from "@playwright/test";
import { ids, productFixtures, stages } from "./support/product-fixtures";
test("graph saves branches and reloads persisted connections", async ({
  page,
}) => {
  await productFixtures(page);
  let graph = {
    stages: stages.map((s, i) => ({ ...s, x: 40 + i * 240, y: 40 })),
    links: [{ fromStageId: "stage-1", toStageId: "stage-4" }],
  };
  let saves = 0;
  await page.route(
    `**/v1/pipeline-versions/${ids.version}/graph`,
    async (route) => {
      if (route.request().method() === "PUT") {
        graph = route.request().postDataJSON();
        saves++;
      }
      await route.fulfill({ json: graph });
    },
  );
  await page.goto(`/employer/pipelines/${ids.job}`);
  const undo = page.getByRole("button", { name: "Undo graph change" });
  const redo = page.getByRole("button", { name: "Redo graph change" });
  await expect(undo).toBeDisabled();
  await expect(redo).toBeDisabled();
  await page
    .getByRole("button", { name: /^English check, 0 incoming/ })
    .click();
  await page.getByLabel("Connect to stage").selectOption("stage-2");
  await page
    .getByRole("button", { name: "Add connection", exact: true })
    .click();
  await expect(undo).toBeEnabled();
  await page
    .getByRole("button", { name: "Save graph draft", exact: true })
    .click();
  await expect.poll(() => saves).toBe(1);
  await expect(undo).toBeDisabled();
  await expect(redo).toBeDisabled();
  expect(graph.links).toContainEqual({
    fromStageId: "stage-1",
    toStageId: "stage-2",
  });
  await page.reload();
  await expect(page.getByRole("button", { name: /, 1 incoming/ })).toHaveCount(
    2,
  );
});

test("a new graph edit clears redo without writing the draft to the server", async ({ page }) => {
  await productFixtures(page);
  const writes: string[] = [];
  page.on("request", (request) => {
    if (new URL(request.url()).pathname.endsWith("/graph") && request.method() !== "GET") {
      writes.push(request.method());
    }
  });
  await page.goto(`/employer/pipelines/${ids.job}`);
  const stage = page.getByRole("button", { name: /^English check, 0 incoming/ });
  await expect(stage).toBeVisible();
  await page.getByRole("button", { name: "Add stage", exact: true }).click();
  const added = page.getByRole("button", { name: /^New assessment, 0 incoming/ });
  await expect(added).toBeVisible();
  await page.getByRole("button", { name: "Undo graph change" }).click();
  await expect(added).toHaveCount(0);
  await expect(page.getByRole("button", { name: "Redo graph change" })).toBeEnabled();
  await stage.focus();
  const before = await stage.evaluate((element) => element.style.left);
  await page.keyboard.press("Alt+ArrowRight");
  await expect.poll(() => stage.evaluate((element) => element.style.left)).not.toBe(before);
  await expect(page.getByRole("button", { name: "Redo graph change" })).toBeDisabled();
  expect(writes).toEqual([]);
});
for (const variant of ["/coding", ""]) {
  test(`candidate runs code with custom input without submitting ${variant}`, async ({
    page,
  }) => {
    await productFixtures(page);
    const question = ids.module;
    let body: Record<string, unknown> | undefined;
    await page.route(`**/v1/assessment-attempts/${ids.application}`, (route) =>
      route.fulfill({
        json: {
          id: ids.application,
          instanceId: ids.application,
          status: "in_progress",
          deadline: "2027-01-01T00:00:00Z",
          savedResponses: {
            [question]: { language: "python", code: "print(42)" },
          },
          questions: [
            { id: question, prompt: "Print input", responseFormat: "code" },
          ],
        },
      }),
    );
    await page.route(
      `**/v1/assessment-attempts/${ids.application}/run`,
      (route) => {
        body = route.request().postDataJSON();
        return route.fulfill({ json: { stdout: "hello\n", succeeded: true } });
      },
    );
    await page.goto(`/candidate/assessments/${ids.application}${variant}`);
    await expect(
      page.getByRole("combobox", { name: "Language", exact: true }),
    ).toHaveValue("python");
    await expect(
      page.getByRole("textbox", { name: "Your answer", exact: true }),
    ).toHaveValue("print(42)");
    await page
      .getByRole("textbox", { name: "Your answer", exact: true })
      .fill("print(input())");
    await page.getByLabel("Standard input").fill("hello");
    await page.getByRole("button", { name: "Run code", exact: true }).click();
    await expect(
      page.getByRole("status").filter({ hasText: "hello" }),
    ).toBeVisible();
    expect(body).toMatchObject({
      source: "print(input())",
      language: "python",
      stdin: "hello",
      questionId: question,
    });
  });
}

test("new posting can create a graph without reloading the page", async ({
  page,
}) => {
  await productFixtures(page);
  let created = false;
  await page.route(`**/v1/jobs/${ids.job}/pipeline`, async (route) => {
    if (route.request().method() === "POST") created = true;
    await route.fulfill(
      created
        ? { json: { id: ids.job, versionId: ids.version, published: false } }
        : { status: 404, json: { code: "not_found" } },
    );
  });
  await page.route(`**/v1/pipeline-versions/${ids.version}/graph`, (route) =>
    route.fulfill({ json: { stages: [], links: [] } }),
  );
  await page.goto(`/employer/pipelines/${ids.job}`);
  await page
    .getByRole("button", { name: "Create graph pipeline", exact: true })
    .click();
  await expect(
    page.getByRole("region", { name: "Pipeline graph", exact: true }),
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "Save graph draft", exact: true }),
  ).toBeEnabled();
});
