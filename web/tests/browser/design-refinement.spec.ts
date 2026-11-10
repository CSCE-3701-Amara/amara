import { expect, test } from "@playwright/test";
import { ids, productFixtures } from "./support/product-fixtures";

const screens = [
  ["welcome", "/"],
  ["feed", "/candidate/dashboard"],
  ["posting-editor", "/employer/jobs/new"],
  ["pipeline", `/employer/pipelines/${ids.job}`],
  ["assessment", `/candidate/assessments/${ids.application}`],
  ["administration", "/admin/dashboard"],
] as const;

// Complement the full inventory's 375/1024/1728 coverage with the design breakpoints.
for (const width of [375, 768, 1280, 1440, 1728]) {
  test(`core layouts retain readable controls and contain overflow at ${width}px`, async ({ page }) => {
    test.setTimeout(120_000);
    await page.setViewportSize({ width, height: 1000 });
    await productFixtures(page);
    for (const [name, path] of screens) {
      await page.goto(path);
      await expect(page.locator("main h1").first()).toBeVisible();
      await expect(page.getByText("Loading your workspace…", { exact: true })).toHaveCount(0);
      await page.evaluate(() => document.fonts.ready);
      expect(await page.evaluate(() => document.documentElement.scrollWidth)).toBeLessThanOrEqual(width);
      // A missing stylesheet must not turn this into a false-positive layout check.
      await expect(page.locator("body")).toHaveCSS("--radius-control", "7px");
      if (path !== "/") {
        const menu = page.getByRole("button", { name: "Menu", exact: true });
        if (width <= 760) {
          await expect(page.getByRole("navigation", { name: "Workspace navigation" })).toBeHidden();
          await menu.click();
          await expect(page.getByRole("button", { name: "Close menu", exact: true })).toHaveAttribute("aria-expanded", "true");
          await expect(page.getByRole("navigation", { name: "Workspace navigation" })).toBeVisible();
          await page.getByRole("button", { name: "Close menu", exact: true }).click();
        }
        const skip = page.getByRole("link", { name: "Skip to content" });
        await skip.focus();
        await page.keyboard.press("Enter");
        await expect(page.locator("#workspace-content")).toBeFocused();
        await page.locator("#workspace-content").blur();
      }
      await page.evaluate(() => window.scrollTo({ top: 0, behavior: "instant" }));
      await page.screenshot({ path: `tests/artifacts/refinement/${name}-${width}.png`, fullPage: true, animations: "disabled" });
    }
  });
}
