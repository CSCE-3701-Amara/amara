import { expect, test } from "@playwright/test";
import { ids, productFixtures } from "./support/product-fixtures";

for (const kind of ["coding", "programming_language", "recorded_video"]) {
  test(`editing ${kind} preserves its grading settings`, async ({ page }) => {
    await productFixtures(page);
    const definition = {
      id: ids.module,
      company_id: ids.job,
      published: false,
      name: "Saved assessment",
      type: kind,
      behavioral: kind === "recorded_video" ? true : undefined,
      active_minutes: kind === "recorded_video" ? undefined : 30,
      minimum_tests_passed: kind === "recorded_video" ? undefined : 1,
      questions: [{
        prompt: "Explain your solution",
        format: kind === "recorded_video" ? "video" : "code",
        amara: "Reasoning",
        difficulty: "easy",
        minimum_cases_passed: kind === "recorded_video" ? undefined : 1,
        test_cases: kind === "recorded_video" ? undefined : [{ input: "", expected_output: "", weight: 1 }],
      }],
      amarat: [{ name: "Reasoning", crest: "speech", domain: "Reasoning", description: "Explain clearly" }],
      rubric: kind === "recorded_video" ? [{ description: "Clarity", amara: "Reasoning" }] : undefined,
    };
    let saved: Record<string, unknown> | undefined;
    await page.route(`**/v1/modules/${ids.module}`, async (route) => {
      if (route.request().method() === "PUT") {
        saved = route.request().postDataJSON();
        return route.fulfill({ json: { ...definition, ...saved } });
      }
      return route.fulfill({ json: definition });
    });
    await page.goto(`/employer/modules/${ids.module}/questions`);
    await page.getByLabel("Module name (required)").fill("Renamed assessment");
    await page.getByRole("button", { name: "Save draft", exact: true }).click();
    await expect.poll(() => saved?.name).toBe("Renamed assessment");
    expect(saved?.minimum_tests_passed).toBe(definition.minimum_tests_passed);
    expect(saved?.behavioral).toBe(definition.behavioral);
    expect(saved?.questions).toEqual(JSON.parse(JSON.stringify(definition.questions)));
    expect(saved?.rubric).toEqual(definition.rubric);
    expect(saved).not.toHaveProperty("id");
    expect(saved).not.toHaveProperty("company_id");
    expect(saved).not.toHaveProperty("published");
  });
}

test("editing a posting preserves its saved finalist target", async ({ page }) => {
  await productFixtures(page);
  const { posting } = await import("./support/product-fixtures");
  const initial = { ...posting, finalist_target: 3 };
  let saved: Record<string, unknown> | undefined;
  await page.route(`**/v1/jobs/${ids.job}`, async (route) => {
    if (route.request().method() === "PUT") {
      saved = route.request().postDataJSON();
      return route.fulfill({ json: { ...initial, ...saved } });
    }
    return route.fulfill({ json: initial });
  });
  await page.goto(`/employer/jobs/${ids.job}`);
  await page.getByLabel("Job title", { exact: true }).fill("Updated role");
  await page.getByRole("button", { name: "Save posting", exact: true }).click();
  await expect.poll(() => saved?.title).toBe("Updated role");
  expect(saved?.finalist_target).toBe(3);
});
