import { expect, test } from "@playwright/test";
import { productFixtures } from "./support/product-fixtures";

// The active side-menu link must stay readable in every workspace. Hiring
// workspaces once kept a dark fill while the shared style switched to dark ink.
for (const path of [
  "/candidate/applications",
  "/employer/jobs",
  "/interviewer/interviews",
]) {
  test(`active menu link on ${path} has readable contrast`, async ({
    page,
  }) => {
    await productFixtures(page);
    await page.goto(path);
    const active = page
      .getByRole("navigation", { name: "Workspace navigation" })
      .locator('a[aria-current="page"]');
    await expect(active).toHaveCount(1);
    const ratio = await active.evaluate((link) => {
      const rgb = (value: string) =>
        (value.match(/[\d.]+/g) || []).slice(0, 3).map(Number);
      const luminance = ([r, g, b]: number[]) => {
        const [R, G, B] = [r, g, b].map((c) => {
          const s = c / 255;
          return s <= 0.03928 ? s / 12.92 : ((s + 0.055) / 1.055) ** 2.4;
        });
        return 0.2126 * R + 0.7152 * G + 0.0722 * B;
      };
      const style = getComputedStyle(link);
      const fg = luminance(rgb(style.color));
      const bg = luminance(rgb(style.backgroundColor));
      return (Math.max(fg, bg) + 0.05) / (Math.min(fg, bg) + 0.05);
    });
    expect(ratio).toBeGreaterThanOrEqual(4.5);
  });
}
