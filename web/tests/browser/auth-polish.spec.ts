import { expect, test } from "@playwright/test";

const variants = [
  { path: "/login", label: "Candidate", heading: "Log in to your account" },
  {
    path: "/login?account=company",
    label: "Hiring team",
    heading: "Log in to your company",
  },
  {
    path: "/register",
    label: "Candidate",
    heading: "Create a candidate account",
  },
  {
    path: "/register?account=company",
    label: "Hiring team",
    heading: "Create a company account",
  },
];

for (const width of [320, 375, 768, 1024, 1280, 1728]) {
  test(`auth selection, logo alignment, and overflow at ${width}px`, async ({
    page,
  }) => {
    await page.setViewportSize({ width, height: width <= 768 ? 812 : 900 });
    let originalLogo: {
      x: number;
      y: number;
      width: number;
      height: number;
    } | null = null;
    for (const variant of variants) {
      await page.goto(variant.path);
      await expect(
        page.getByRole("heading", { name: variant.heading, exact: true }),
      ).toBeVisible();
      await page.evaluate(() => document.fonts.ready);
      const selector = page.getByRole("navigation", { name: "Account type" });
      const selected = selector.getByRole("link", {
        name: variant.label,
        exact: true,
      });
      await expect(selected).toHaveAttribute("aria-current", "page");
      await expect(selected).toHaveCSS("color", "rgb(255, 255, 255)");
      await expect(selected.locator("span")).toHaveCSS(
        "background-color",
        "rgb(102, 254, 76)",
      );
      await expect(selector.locator("a:not([aria-current])")).toHaveCSS(
        "color",
        "rgb(1, 1, 1)",
      );
      const indicator = await selector.evaluate((el) => {
        // The inactive surface and transition curve come from shared theme tokens.
        const style = getComputedStyle(el, "::before");
        return {
          color: style.backgroundColor,
          x: new DOMMatrix(style.transform).m41,
          width: parseFloat(style.width),
        };
      });
      await expect(selector.locator("a:not([aria-current])")).toHaveCSS("background-color", "rgb(241, 238, 232)");
      await expect(page.getByRole("link", { name: "Back to amara" })).toHaveCSS("background-color", "rgb(241, 238, 232)");
      expect(indicator.color).toBe("rgb(1, 1, 1)");
      expect(indicator.x > indicator.width).toBe(variant.label === "Candidate");
      const logo = page.getByRole("link", { name: "amara home" });
      const bounds = await logo.boundingBox();
      expect(bounds).not.toBeNull();
      if (!originalLogo) originalLogo = bounds;
      else {
        expect(bounds!.x).toBeCloseTo(originalLogo.x, 0);
        expect(bounds!.y).toBeCloseTo(originalLogo.y, 0);
        expect(bounds!.width).toBeCloseTo(originalLogo.width, 0);
        expect(bounds!.height).toBeCloseTo(originalLogo.height, 0);
      }
      expect(
        await page.evaluate(() => document.documentElement.scrollWidth),
      ).toBeLessThanOrEqual(width);
      const outside = await page
        .locator("main input, main button, main a, main h1")
        .evaluateAll((elements) =>
          elements
            .filter((el) => {
              const r = el.getBoundingClientRect();
              return r.width > 0 && (r.left < -1 || r.right > innerWidth + 1);
            })
            .map((el) => el.textContent),
        );
      expect(outside).toEqual([]);
    }
  });
}

test("switching accounts moves the active indicator and preserves entered credentials", async ({
  page,
}) => {
  await page.goto("/login");
  await page.getByLabel("Email", { exact: true }).fill("aya@example.com");
  await page.getByLabel("Password", { exact: true }).fill("a-secure-password");
  const selector = page.getByRole("navigation", { name: "Account type" });
  await selector
    .getByRole("link", { name: "Hiring team", exact: true })
    .click();
  await expect(
    page.getByRole("heading", { name: "Log in to your company" }),
  ).toBeVisible();
  await expect(
    selector.getByRole("link", { name: "Hiring team", exact: true }),
  ).toHaveCSS("color", "rgb(255, 255, 255)");
  await expect(
    selector.getByRole("link", { name: "Candidate", exact: true }),
  ).toHaveCSS("color", "rgb(1, 1, 1)");
  await expect(page.getByLabel("Email", { exact: true })).toHaveValue(
    "aya@example.com",
  );
  await expect(page.getByLabel("Password", { exact: true })).toHaveValue(
    "a-secure-password",
  );
  await selector.getByRole("link", { name: "Candidate", exact: true }).click();
  await expect(
    selector.getByRole("link", { name: "Candidate", exact: true }),
  ).toHaveCSS("color", "rgb(255, 255, 255)");
});

test("workspace-specific registration fields reset when switching accounts", async ({
  page,
}) => {
  await page.goto("/register");
  await page.getByLabel("First name", { exact: true }).fill("Aya");
  await page.getByLabel("Last name", { exact: true }).fill("Ali");
  await page.getByRole("link", { name: "Hiring team", exact: true }).click();
  await expect(page.getByLabel("Company name", { exact: true })).toHaveValue(
    "",
  );
  await expect(
    page.getByLabel("Tax registration number", { exact: true }),
  ).toHaveValue("");
});

test("motion is subtle when enabled and stops for reduced motion", async ({
  page,
}) => {
  await page.emulateMedia({ reducedMotion: "no-preference" });
  await page.goto("/login");
  expect(await page.getByRole("navigation", { name: "Account type" }).evaluate(el => getComputedStyle(el, "::before").transitionDuration)).toBe("0.4s");
  const moving = await page
    .locator("aside")
    .evaluate(
      (el) =>
        [...el.querySelectorAll("*")].filter(
          (node) => getComputedStyle(node).animationName !== "none",
        ).length,
    );
  expect(moving).toBeGreaterThanOrEqual(3);
  await page.emulateMedia({ reducedMotion: "reduce" });
  const remaining = await page
    .locator("main")
    .evaluate(
      (el) =>
        [...el.querySelectorAll("*")].filter(
          (node) => getComputedStyle(node).animationName !== "none",
        ).length,
    );
  expect(remaining).toBe(0);
});

for (const width of [320, 768, 1024, 1280]) {
  test(`welcome navigation and actions fit at ${width}px without WebGPU`, async ({
    page,
  }) => {
    // Also verify the supported fallback for browsers with no graphics adapter.
    await page.addInitScript(() =>
      Object.defineProperty(navigator, "gpu", { value: undefined }),
    );
    await page.setViewportSize({ width, height: 900 });
    await page.goto("/");
    await page.evaluate(() => document.fonts.ready);
    const outside = await page
      .locator("main a, main h1")
      .evaluateAll((elements) =>
        elements
          .filter((el) => {
            const r = el.getBoundingClientRect();
            return r.width > 0 && (r.left < -1 || r.right > innerWidth + 1);
          })
          .map((el) => el.textContent),
      );
    expect(outside).toEqual([]);
    expect(
      await page.evaluate(() => document.documentElement.scrollWidth),
    ).toBeLessThanOrEqual(width);
  });
}
