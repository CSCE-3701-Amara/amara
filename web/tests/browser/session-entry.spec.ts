import { expect, test } from "@playwright/test";
import { productFixtures } from "./support/product-fixtures";

const destinations = [
  ["candidate", "/candidate/dashboard"],
  ["company_representative", "/employer/jobs"],
  ["interviewer", "/interviewer/interviews"],
] as const;

for (const [role, destination] of destinations) {
  test(`${role} workspace logo returns directly to its home`, async ({ page }) => {
    await productFixtures(page);
    await page.goto(destination);
    await expect(page.getByRole("link", { name: "amara home" })).toHaveAttribute("href", destination);
    await page.getByRole("link", { name: "amara home" }).click();
    await expect(page).toHaveURL(destination);
  });
  for (const entry of ["/", "/login"]) {
    test(`${role} resumes a persisted session from ${entry}`, async ({ browser, baseURL }) => {
      const context = await browser.newContext({
        storageState: { cookies: [{ name: "amara_session", value: "persisted-session", domain: "127.0.0.1", path: "/", httpOnly: true, secure: false, sameSite: "Lax", expires: Date.now() / 1000 + 3600 }], origins: [] },
      });
      try {
        const page = await context.newPage();
        await productFixtures(page);
        let checks = 0;
        await page.route("**/v1/auth/sessions/current", async route => {
          expect(route.request().method()).toBe("GET");
          expect((await route.request().allHeaders()).cookie).toContain("amara_session=persisted-session");
          checks++;
          await route.fulfill({ json: { user: { id: "user-id", role }, session: { expiresAt: "2027-01-01T00:00:00Z" } } });
        });
        await page.goto(`${baseURL}${entry}`);
        await expect(page).toHaveURL(`${baseURL}${destination}`);
        expect(checks).toBeGreaterThanOrEqual(1);
        await expect(page.getByRole("heading", { name: "Beyond Resumes. Towards Skills." })).toHaveCount(0);
      } finally {
        await context.close();
      }
    });
  }
}

test("revoked session exposes sign-in and failed lookup offers retry without showing welcome", async ({ page, context, baseURL }) => {
  await context.addCookies([{ name: "amara_session", value: "old-session", url: baseURL!, httpOnly: true }]);
  await page.route("**/v1/auth/sessions/current", route => route.fulfill({ status: 503, json: { detail: "Authentication is temporarily unavailable." } }));
  await page.goto("/");
  await expect(page.locator("main").getByRole("alert")).toContainText("Authentication is temporarily unavailable");
  await expect(page.getByRole("heading", { name: /Beyond Resumes/ })).toHaveCount(0);
  await page.route("**/v1/auth/sessions/current", route => route.fulfill({ status: 401, json: { error: { code: "authentication_required", message: "Authentication is required." } } }));
  await page.getByRole("button", { name: "Try again", exact: true }).click();
  await expect(page.getByRole("heading", { name: /Beyond Resumes/ })).toBeVisible();
  await page.goto("/login");
  await expect(page.getByRole("heading", { name: "Log in to your account", exact: true })).toBeVisible();
});

test("expired session responses are not used as authenticated navigation", async ({ page, context, baseURL }) => {
  await context.addCookies([{ name: "amara_session", value: "expired-session", url: baseURL!, httpOnly: true }]);
  await page.route("**/v1/auth/sessions/current", route => route.fulfill({ json: { user: { id: "user-id", role: "candidate" }, session: { expiresAt: "2020-01-01T00:00:00Z" } } }));
  await page.goto("/login");
  await expect(page.getByRole("heading", { name: "Log in to your account", exact: true })).toBeVisible();
  await expect(page).toHaveURL(/\/login$/);
});

test("browser Back after sign-in rechecks the cookie instead of showing the cached welcome page", async ({ page }) => {
  await productFixtures(page);
  await page.route("**/v1/auth/sessions/current", route => route.fulfill({ json: { user: { id: "candidate-id", role: "candidate" }, session: { expiresAt: "2027-01-01T00:00:00Z" } } }));
  await page.route("**/v1/auth/sessions", route => route.fulfill({
    json: { user: { id: "candidate-id", email: "candidate@example.com", role: "candidate" }, session: { expiresAt: "2027-01-01T00:00:00Z" } },
    headers: { "Set-Cookie": "amara_session=active-session; Path=/; HttpOnly; SameSite=Lax" },
  }));
  await page.goto("/");
  await page.getByRole("link", { name: "Log in", exact: true }).click();
  await page.getByLabel("Email", { exact: true }).fill("candidate@example.com");
  await page.getByLabel("Password", { exact: true }).fill("test-password");
  await page.getByRole("button", { name: "Log in to amara", exact: true }).click();
  await expect(page).toHaveURL(/\/candidate\/dashboard$/);
  await page.goBack();
  await expect(page).toHaveURL(/\/candidate\/dashboard$/);
  await expect(page.getByRole("heading", { name: /Beyond Resumes/ })).toHaveCount(0);
});

test("unexpected session roles cannot supply a redirect destination", async ({ page, context, baseURL }) => {
  await context.addCookies([{ name: "amara_session", value: "invalid-response", url: baseURL!, httpOnly: true }]);
  await page.route("**/v1/auth/sessions/current", route => route.fulfill({ json: { user: { id: "user-id", role: "https://untrusted.example" }, session: { expiresAt: "2027-01-01T00:00:00Z" } } }));
  await page.goto("/login");
  await expect(page.locator("main").getByRole("alert")).toContainText("unexpected response");
  await expect(page).toHaveURL(/\/login$/);
});

test("hiring-manager sign-in requires MFA before issuing a resumable session", async ({ page, context }) => {
  await productFixtures(page);
  let issued = false;
  const identity = { user: { id: "manager-id", email: "manager@example.com", role: "company_representative" }, session: { expiresAt: "2027-01-01T00:00:00Z" } };
  await page.route("**/v1/auth/sessions/current", route => route.fulfill({ json: identity }));
  await page.route("**/v1/auth/sessions", async route => {
    const input = route.request().postDataJSON();
    if (!input.mfaCode) return route.fulfill({ status: 403, json: { code: "email_mfa_required" } });
    expect(input.mfaCode).toBe("123456");
    issued = true;
    await route.fulfill({ json: identity, headers: { "Set-Cookie": "amara_session=mfa-verified-session; Path=/; HttpOnly; SameSite=Lax" } });
  });
  await page.goto("/login?account=company");
  await page.getByLabel("Email", { exact: true }).fill("manager@example.com");
  await page.getByLabel("Password", { exact: true }).fill("test-password");
  await page.getByRole("button", { name: "Log in to amara", exact: true }).click();
  await expect(page.getByRole("heading", { name: "Verify your sign-in", exact: true })).toBeVisible();
  expect(issued).toBe(false);
  expect((await context.cookies()).some(cookie => cookie.name === "amara_session")).toBe(false);
  await page.getByLabel("Multi-factor code", { exact: true }).fill("123456");
  await page.getByRole("button", { name: "Verify and sign in", exact: true }).click();
  await expect(page).toHaveURL(/\/employer\/jobs$/);
  expect(issued).toBe(true);
  await page.goto("/login?account=company");
  await expect(page).toHaveURL(/\/employer\/jobs$/);
});

test("administrator workspace logo stays in administration", async ({ page }) => {
  await productFixtures(page);
  await page.goto("/admin/dashboard");
  await expect(page.getByRole("link", { name: "amara home" })).toHaveAttribute("href", "/admin/dashboard");
});
