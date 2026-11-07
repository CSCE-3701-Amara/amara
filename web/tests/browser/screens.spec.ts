import { expect, test } from "@playwright/test";
import { mkdir } from "node:fs/promises";

const screens = [
  ["login-candidate", "/login", "Log in to your account"],
  ["login-company", "/login?account=company", "Log in to your company"],
  ["register-candidate", "/register", "Create a candidate account"],
  ["register-company", "/register?account=company", "Create a company account"],
  ["mfa", "/login?state=mfa", "Verify your sign-in"],
  ["code-rejected", "/login?state=code-rejected", "Try your code again"],
  [
    "credentials-rejected",
    "/login?state=credentials-rejected",
    "Check your sign-in details",
  ],
  ["rate-limited", "/login?state=rate-limited", "Too many sign-in attempts"],
  ["expired", "/login?state=expired", "Sign in again"],
  [
    "verification-pending",
    "/register?state=verification-pending",
    "Verification pending",
  ],
  [
    "review-pending",
    "/register?account=company&state=review-pending",
    "Registration under review",
  ],
];

for (const width of [1728, 375]) {
  for (const [name, url, heading] of screens) {
    test(`${name} at ${width}px: assets and overflow`, async ({ page }) => {
      test.setTimeout(120000);
      await page.setViewportSize({
        width,
        height: width === 1728 ? 1117 : 812,
      });
      await page.goto(url);
      await expect(
        page
          .getByRole("heading", {
            name: new RegExp(heading.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")),
          })
          .first(),
      ).toBeVisible();
      await page.evaluate(() => document.fonts.ready);
      // A loaded page with missing CSS must not count as a passing visual check.
      expect(
        await page
          .locator("h1")
          .first()
          .evaluate((el) => parseFloat(getComputedStyle(el).fontSize)),
      ).toBeGreaterThanOrEqual(30);
      // Trigger native lazy loading instead of changing server-rendered attributes
      // before hydration, which creates artificial React mismatch warnings.
      for (const img of await page.locator("img").all()) {
        if (await img.isVisible()) {
          await img.scrollIntoViewIfNeeded();
          await expect.poll(() => img.evaluate(el => (el as HTMLImageElement).naturalWidth)).toBeGreaterThan(0);
        }
      }
      await page.evaluate(() => window.scrollTo({ top: 0, behavior: "instant" }));
      const broken = await page
        .locator("img")
        .evaluateAll((images) =>
          images
            .filter((img) => !(img as HTMLImageElement).naturalWidth)
            .map((img) => (img as HTMLImageElement).src),
        );
      expect(broken).toEqual([]);
      expect(
        await page.evaluate(() => document.documentElement.scrollWidth),
      ).toBeLessThanOrEqual(width);
      const overflowing = await page
        .locator("main input, main button, main a, main h1")
        .evaluateAll((elements) =>
          elements
            .filter((el) => {
              const r = el.getBoundingClientRect();
              return r.width > 0 && (r.left < -1 || r.right > innerWidth + 1);
            })
            .map((el) => el.textContent),
        );
      expect(overflowing).toEqual([]);
      await mkdir("tests/artifacts", { recursive: true });
      await page.screenshot({
        path: `tests/artifacts/${name}-${width}.png`,
        fullPage: true,
      });
    });
  }
}
