import { test, expect } from "@playwright/test";
import inventory from "../../design/figma-inventory.json";
import { productFixtures, ids } from "./support/product-fixtures";
const screens = inventory.screens.filter((row) => row.page !== "Login Page");
for (const width of [1728, 375])
  for (const screen of screens) {
    test(`${screen.id} ${screen.name} at ${width}px`, async ({ page }) => {
      await page.setViewportSize({ width, height: 1117 });
      await productFixtures(page);
      const url = screen.route.replaceAll(
        "[id]",
        screen.route.startsWith("/candidate/jobs/")
          ? ids.job
          : screen.route.includes("/modules/")
            ? ids.module
            : screen.route.includes("/interviews/")
              ? ids.interview
              : ids.application,
      );
      const errors: string[] = [];
      page.on("pageerror", (e) => errors.push(e.message));
      const response = await page.goto(url);
      expect(response?.status()).toBe(200);
      await expect(page.locator("main h1, main h2").first()).toBeVisible();
      await page.evaluate(() => document.fonts.ready);
      await expect(
        page.getByText("Loading your workspace…", { exact: true }),
      ).toHaveCount(0);
      await expect(page.locator("main")).toHaveCSS(
        "--workspace-ink",
        "#211813",
      );
      if (screen.name.startsWith("Dialog")) {
        await expect(page.getByRole("dialog")).toBeVisible();
        await expect(
          page
            .getByRole("dialog")
            .getByRole("button", { name: "Cancel", exact: true }),
        ).toBeFocused();
        await page.screenshot({
          path: `tests/artifacts/product/${screen.id.replace(":", "-")}-${width}-dialog.png`,
          animations: "disabled",
        });
        await page.keyboard.press("Escape");
        await expect(page.getByRole("dialog")).toHaveCount(0);
      }
      const overflow = await page.evaluate(
        () => document.documentElement.scrollWidth > innerWidth + 1,
      );
      expect(overflow).toBe(false);
      expect(errors).toEqual([]);
      await page.screenshot({
        path: `tests/artifacts/product/${screen.id.replace(":", "-")}-${width}.png`,
        fullPage: true,
        animations: "disabled",
      });
    });
  }
