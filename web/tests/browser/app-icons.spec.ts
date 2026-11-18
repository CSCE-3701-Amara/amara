import { expect, test } from "@playwright/test";
import { ids, productFixtures } from "./support/product-fixtures";

// The Amara mark is the tab icon on every page: public, candidate and employer
// routes all inherit it from the root app directory.
for (const path of [
  "/login",
  "/candidate/applications",
  `/employer/pipelines/${ids.job}`,
]) {
  test(`${path} uses the Amara tab icon`, async ({ page }) => {
    await productFixtures(page);
    await page.goto(path);
    const icons = await page
      .locator('head link[rel="icon"], head link[rel="apple-touch-icon"]')
      .evaluateAll((links) =>
        links.map((link) => ({
          rel: link.getAttribute("rel"),
          href: link.getAttribute("href") || "",
          type: link.getAttribute("type"),
        })),
      );
    const hrefs = icons.map((icon) => icon.href.split("?")[0]);
    expect(hrefs).toEqual(
      expect.arrayContaining(["/favicon.ico", "/icon.svg", "/apple-icon.png"]),
    );
    expect(icons.find((icon) => icon.href.startsWith("/icon.svg"))?.type).toBe(
      "image/svg+xml",
    );
    for (const { href } of icons) {
      const response = await page.request.get(href);
      expect(response.status(), href).toBe(200);
    }
    const svg = await (await page.request.get("/icon.svg")).text();
    expect(svg).toContain("prefers-color-scheme: dark");
  });
}
