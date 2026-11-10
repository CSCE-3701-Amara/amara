import { expect, test } from "@playwright/test";
import { ids, productFixtures } from "./support/product-fixtures";

for (const width of [1440, 375]) {
  test(`candidate sees parallel assessments and join prerequisites at ${width}px`, async ({
    page,
  }) => {
    await page.setViewportSize({ width, height: 1000 });
    await productFixtures(page);
    await page.route(
      `**/v1/applications/${ids.application}/pipeline-progress`,
      (route) =>
        route.fulfill({
          json: {
            stages: [
              { id: "code", name: "Coding assessment", status: "active" },
              {
                id: "collaboration",
                name: "Collaboration assessment",
                status: "active",
              },
              { id: "review", name: "Final review", status: "upcoming_locked" },
            ],
            links: [
              { from_stage_id: "code", to_stage_id: "review" },
              { from_stage_id: "collaboration", to_stage_id: "review" },
            ],
          },
        }),
    );
    await page.route(
      `**/v1/applications/${ids.application}/assessment-instances`,
      (route) =>
        route.fulfill({
          json: [
            { id: "code-instance", stageId: "code", status: "available" },
            {
              id: "collaboration-instance",
              stageId: "collaboration",
              status: "available",
            },
          ],
        }),
    );
    await page.goto(`/candidate/applications/${ids.application}/progress`);
    const first = page.getByRole("region", {
      name: "Pipeline step 1",
      exact: true,
    });
    await expect(
      first.getByText("2 parallel stages", { exact: true }),
    ).toBeVisible();
    await expect(
      first.getByRole("heading", { name: "Coding assessment", exact: true }),
    ).toBeVisible();
    await expect(
      first.getByRole("heading", {
        name: "Collaboration assessment",
        exact: true,
      }),
    ).toBeVisible();
    const next = page.getByRole("region", {
      name: "Pipeline step 2",
      exact: true,
    });
    await expect(
      next.getByText(
        "Opens after Coding assessment and Collaboration assessment pass.",
      ),
    ).toBeVisible();
    await expect(
      first.getByRole("link", { name: "Open assessment" }).first(),
    ).toHaveAttribute("href", /code-instance/);
    await expect(
      next.getByRole("link", { name: "Open assessment" }),
    ).toHaveCount(0);
    const arrows = page.getByRole("img", { name: "Stage dependency arrows" });
    await expect(
      arrows.locator('path[data-from="code"][data-to="review"]'),
    ).toHaveCount(1);
    await expect(
      arrows.locator('path[data-from="collaboration"][data-to="review"]'),
    ).toHaveCount(1);
    const roots = await first.getByRole("article").all();
    const a = await roots[0].boundingBox();
    const b = await roots[1].boundingBox();
    const join = await next.getByRole("article").boundingBox();
    expect(a!.x).toBe(b!.x);
    expect(join!.x).toBeGreaterThan(a!.x + a!.width);
    expect(
      await page.evaluate(
        () => document.documentElement.scrollWidth <= window.innerWidth,
      ),
    ).toBe(true);
    await page.screenshot({
      path: `/tmp/amara-candidate-parallel-${width}.png`,
      fullPage: true,
    });
  });
}
