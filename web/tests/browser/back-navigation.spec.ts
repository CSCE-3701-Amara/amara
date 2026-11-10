import { expect, test } from "@playwright/test";
import { ids, productFixtures } from "./support/product-fixtures";

for (const width of [375, 1280]) {
  test(`directly opened pages have usable back navigation at ${width}px`, async ({ page }) => {
    await productFixtures(page);
    await page.setViewportSize({ width, height: 900 });
    for (const [path, destination, label] of [
      ["/employer/jobs/new", "/employer/jobs", "Back to job postings"],
      [`/candidate/applications/${ids.application}/progress`, `/candidate/applications/${ids.application}`, "Back to application"],
      ["/interviewer/profile", "/interviewer/interviews", "Back to interviews"],
      ["/admin/accounts", "/admin/dashboard", "Back to company registrations"],
    ]) {
      await page.goto(path);
      const back = page.locator("#workspace-content").getByRole("link", { name: label, exact: true }).first();
      await expect(back).toBeVisible();
      await expect(back).toHaveAttribute("href", destination);
      expect((await back.boundingBox())!.height).toBeGreaterThanOrEqual(44);
      expect(await page.evaluate(() => document.documentElement.scrollWidth)).toBeLessThanOrEqual(width);
      await back.focus();
      await expect(back).toBeFocused();
      await page.keyboard.press("Enter");
      await expect(page).toHaveURL(destination);
    }
  });
}

test("nested back links lead through parents and never to the public welcome page", async ({ page }) => {
  await productFixtures(page);
  await page.goto("/employer/modules/library/problems");
  await page.getByRole("link", { name: "Back to module library", exact: true }).click();
  await expect(page).toHaveURL("/employer/modules/library");
  await page.getByRole("link", { name: "Back to modules", exact: true }).click();
  await expect(page).toHaveURL("/employer/modules");
  await page.getByRole("link", { name: "Back to job postings", exact: true }).click();
  await expect(page).toHaveURL("/employer/jobs");
  await expect(page.getByRole("link", { name: /^Back to / })).toHaveCount(0);
});

test("returning from a coding variant preserves the active assessment attempt", async ({ page }) => {
  await productFixtures(page);
  const context = `attempt=${ids.module}&application=${ids.application}`;
  await page.goto(`/candidate/assessments/${ids.job}/coding?${context}`);
  const back = page.getByRole("link", { name: "Back to assessment", exact: true });
  await expect(back).toHaveAttribute("href", `/candidate/assessments/${ids.job}?${context}`);
  await back.click();
  await expect(page).toHaveURL(`/candidate/assessments/${ids.job}?${context}`);
});
