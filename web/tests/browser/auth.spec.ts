import { expect, test, type Page } from "@playwright/test";

// Auth navigation now opens connected workspaces; keep these tests isolated.
test.beforeEach(async ({ page }) => {
  await page.route("**/v1/jobs**", (route) =>
    route.fulfill({ json: { items: [] } }),
  );
});

async function credentials(page: Page) {
  await page.getByLabel("Email", { exact: true }).fill("candidate@example.com");
  await page
    .getByLabel("Password", { exact: true })
    .fill("correct-horse-battery");
}

test("native validation prevents empty or invalid sign-in requests", async ({
  page,
}) => {
  let calls = 0;
  await page.route("**/v1/auth/sessions", (route) => {
    calls++;
    return route.fulfill({ status: 500, json: {} });
  });
  await page.goto("/login");
  await page.getByRole("button", { name: "Log in to amara" }).click();
  await expect(page.getByLabel("Email", { exact: true })).toBeFocused();
  await page.getByLabel("Email", { exact: true }).fill("invalid-email");
  await page.getByRole("button", { name: "Log in to amara" }).click();
  expect(calls).toBe(0);
});

for (const [role, destination] of [
  ["candidate", "/candidate/dashboard"],
  ["company_representative", "/employer/jobs"],
]) {
  test(`successful ${role} login follows the server role`, async ({ page }) => {
    await page.route("**/v1/auth/sessions", async (route) => {
      expect(route.request().postDataJSON()).toEqual({
        email: "candidate@example.com",
        password: "correct-horse-battery",
      });
      await route.fulfill({
        json: {
          user: { id: "123", email: "candidate@example.com", role },
          session: { expiresAt: "2027-01-01T00:00:00Z" },
        },
      });
    });
    await page.goto("/login");
    await credentials(page);
    await page.getByRole("button", { name: "Log in to amara" }).click();
    await expect(page).toHaveURL(new RegExp(`${destination}$`));
    expect(await page.evaluate(() => ({ ...localStorage }))).toEqual({});
  });
}

test("MFA retains credentials in memory, retries a rejected code, and then navigates", async ({
  page,
}) => {
  const requests: Record<string, string>[] = [];
  await page.route("**/v1/auth/sessions", async (route) => {
    const body = route.request().postDataJSON();
    requests.push(body);
    if (!body.mfaCode)
      return route.fulfill({ status: 403, json: { code: "mfa_required" } });
    if (body.mfaCode === "000000")
      return route.fulfill({ status: 401, json: { code: "invalid_mfa_code" } });
    return route.fulfill({
      json: {
        user: { role: "candidate" },
        session: { expiresAt: "2027-01-01T00:00:00Z" },
      },
    });
  });
  await page.goto("/login");
  await credentials(page);
  await page.getByRole("button", { name: "Log in to amara" }).click();
  await expect(
    page.getByRole("heading", { name: "Verify your sign-in" }),
  ).toBeVisible();
  await page.getByLabel("Multi-factor code").fill("000000");
  await page.getByRole("button", { name: "Verify and sign in" }).click();
  await expect(
    page.getByRole("heading", { name: "Try your code again" }),
  ).toBeVisible();
  await page.getByLabel("Multi-factor code").fill("123456");
  await page.getByRole("button", { name: "Verify and sign in" }).click();
  await expect(page).toHaveURL(/candidate\/dashboard$/);
  expect(requests).toHaveLength(3);
  expect(requests[2]).toEqual({
    email: "candidate@example.com",
    password: "correct-horse-battery",
    mfaCode: "123456",
  });
});

for (const [code, status, title] of [
  ["invalid_credentials", 401, "Check your sign-in details"],
  ["login_rate_limited", 429, "Too many sign-in attempts"],
] as const) {
  test(`renders the ${code} frame from the server response`, async ({
    page,
  }) => {
    await page.route("**/v1/auth/sessions", (route) =>
      route.fulfill({
        status,
        json: { code },
        headers: { "Retry-After": "60" },
      }),
    );
    await page.goto("/login");
    await credentials(page);
    await page.getByRole("button", { name: "Log in to amara" }).click();
    await expect(page.getByRole("heading", { name: title })).toBeVisible();
    await page.getByRole("button", { name: "Back to sign-in" }).click();
    await expect(page.getByLabel("Email", { exact: true })).toHaveValue(
      "candidate@example.com",
    );
  });
}

for (const company of [false, true]) {
  test(`${company ? "company" : "candidate"} registration only shows pending after server success`, async ({
    page,
  }) => {
    let calls = 0;
    await page.route(
      `**/v1/auth/${company ? "company-registrations" : "registrations"}`,
      (route) => {
        const body = route.request().postDataJSON();
        if (company)
          expect(body).toMatchObject({
            companyName: "Acme",
            registrationNumber: "123456",
          });
        else
          expect(body).toMatchObject({
            firstName: "Aya",
            lastName: "Ali",
            accountType: "candidate",
          });
        calls++;
        return calls === 1
          ? route.fulfill({
              status: 503,
              json: {
                code: "registration_unavailable",
                detail: "Account registration is temporarily unavailable.",
              },
            })
          : route.fulfill({
              status: 201,
              json: {
                status: company ? "pending_review" : "pending_verification",
              },
            });
      },
    );
    await page.goto(company ? "/register?account=company" : "/register");
    await credentials(page);
    await page
      .getByLabel(company ? "Company name" : "First name", { exact: true })
      .fill(company ? "Acme" : "Aya");
    await page
      .getByLabel(company ? "Tax registration number" : "Last name", {
        exact: true,
      })
      .fill(company ? "123456" : "Ali");
    await page
      .getByRole("button", { name: "Create account", exact: true })
      .click();
    await expect(page.locator("main").getByRole("alert")).toHaveText(
      "Account registration is temporarily unavailable.",
    );
    await page
      .getByRole("button", { name: "Create account", exact: true })
      .click();
    await expect(
      page.getByRole("heading", {
        name: company ? "Registration under review" : "Verification pending",
      }),
    ).toBeVisible();
  });
}

test("account selector and registration navigation preserve the workspace", async ({
  page,
}) => {
  await page.goto("/login");
  await page.getByRole("link", { name: "Hiring team", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "Log in to your company" }),
  ).toBeVisible();
  await page
    .getByRole("link", { name: "Create an account", exact: true })
    .click();
  await expect(page).toHaveURL(/register\?account=company$/);
  await page.getByRole("link", { name: "Candidate", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "Create a candidate account" }),
  ).toBeVisible();
});
