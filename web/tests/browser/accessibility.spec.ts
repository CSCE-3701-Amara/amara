import { expect, test, type Page } from "@playwright/test";

// Auth navigation now opens connected workspaces; keep these tests isolated.
test.beforeEach(async ({ page }) => {
  await page.route("**/v1/jobs**", (route) =>
    route.fulfill({ json: { items: [] } }),
  );
});
async function credentials(page: Page) {
  await page.getByLabel("Email", { exact: true }).fill("aya@example.com");
  await page.getByLabel("Password", { exact: true }).fill("test-password");
}
test("field errors are described and focused; password can be revealed", async ({
  page,
}) => {
  await page.goto("/register");
  await credentials(page);
  const password = page.getByLabel("Password", { exact: true });
  await expect(password).toHaveAccessibleDescription("Use 8–128 characters.");
  await page.getByRole("button", { name: "Show password" }).click();
  await expect(password).toHaveAttribute("type", "text");
  await page.getByRole("button", { name: "Hide password" }).click();
  await expect(password).toHaveAttribute("type", "password");
  await page.getByLabel("First name", { exact: true }).fill("   ");
  await page.getByLabel("Last name", { exact: true }).fill("Ali");
  await page.getByLabel("Last name", { exact: true }).press("Enter");
  const first = page.getByLabel("First name", { exact: true });
  await expect(first).toBeFocused();
  await expect(first).toHaveAttribute("aria-invalid", "true");
  await expect(first).toHaveAccessibleDescription(
    "Enter a value; spaces alone aren’t accepted.",
  );
});
test("pending registration keeps the workspace consistent", async ({
  page,
}) => {
  let release!: () => void;
  const pending = new Promise<void>((resolve) => {
    release = resolve;
  });
  await page.route("**/v1/auth/registrations", async (route) => {
    await pending;
    await route.fulfill({
      status: 201,
      json: { status: "pending_verification" },
    });
  });
  await page.goto("/register");
  await credentials(page);
  await page.getByLabel("First name", { exact: true }).fill("Aya");
  await page.getByLabel("Last name", { exact: true }).fill("Ali");
  await page
    .getByRole("button", { name: "Create account", exact: true })
    .click();
  const team = page.getByRole("link", { name: "Hiring team", exact: true });
  await expect(team).toHaveAttribute("aria-disabled", "true");
  await team.dispatchEvent("click");
  await expect(page).toHaveURL(/\/register$/);
  release();
  await expect(
    page.getByRole("heading", { name: "Verification pending" }),
  ).toBeFocused();
});
test("MFA and recovery restore keyboard focus", async ({ page }) => {
  await page.route("**/v1/auth/sessions", (route) =>
    route.fulfill({ status: 403, json: { code: "mfa_required" } }),
  );
  await page.goto("/login");
  await credentials(page);
  await page.getByLabel("Password", { exact: true }).press("Enter");
  await expect(page.getByLabel("Multi-factor code")).toBeFocused();
  await page.getByRole("button", { name: "Back to sign-in" }).click();
  await expect(page.getByLabel("Email", { exact: true })).toBeFocused();
});
test("failed sign-in can be retried with Enter", async ({ page }) => {
  let calls = 0;
  await page.route("**/v1/auth/sessions", (route) =>
    ++calls === 1
      ? route.abort()
      : route.fulfill({
          json: {
            user: { role: "candidate" },
            session: { expiresAt: "2027-01-01" },
          },
        }),
  );
  await page.goto("/login");
  await credentials(page);
  const password = page.getByLabel("Password", { exact: true });
  await password.press("Enter");
  await expect(page.locator("main").getByRole("alert")).toContainText(
    "Check your connection",
  );
  await expect(
    page.getByRole("button", { name: "Log in to amara" }),
  ).toBeEnabled();
  await password.press("Enter");
  await expect(page).toHaveURL(/candidate\/dashboard$/);
});
test("mobile actions have sufficient targets and respect reduced motion", async ({
  page,
}) => {
  await page.emulateMedia({ reducedMotion: "no-preference" });
  await page.setViewportSize({ width: 375, height: 812 });
  await page.goto("/login");
  await page
    .getByRole("link", { name: "Create an account", exact: true })
    .click();
  await expect(
    page.getByRole("button", { name: /pause animation|resume animation/i }),
  ).toHaveCount(0);
  await page.emulateMedia({ reducedMotion: "reduce" });
  await expect(
    page.getByRole("heading", { name: "Create a candidate account" }),
  ).toBeVisible();
  expect(
    await page
      .locator("main:visible")
      .evaluate((el) => el.getAnimations({ subtree: true }).length),
  ).toBe(0);
  expect(
    await page
      .locator("main:visible button, main:visible a")
      .evaluateAll((elements) =>
        elements
          .filter(
            (el) =>
              el.getBoundingClientRect().height > 0 &&
              el.getBoundingClientRect().height < 44,
          )
          .map((el) => el.textContent),
      ),
  ).toEqual([]);
});
