import { expect, test, type Request } from "@playwright/test";
import { ids, productFixtures } from "./support/product-fixtures";

test("collapsed stage editor defers duplicate reads and retains edits after reopening", async ({
  page,
}) => {
  await productFixtures(page);
  const reads: string[] = [];
  const pendingReads = new Set<Request>();
  page.on("request", (request) => {
    if (new URL(request.url()).pathname.startsWith("/v1/"))
      pendingReads.add(request);
  });
  // Strict Mode replays effects in development and aborts the first fetch.
  // Count completed reads so canceled probes do not look like duplicate data loads.
  page.on("requestfinished", (request) => {
    pendingReads.delete(request);
    const path = new URL(request.url()).pathname;
    if (path.startsWith("/v1/")) reads.push(path);
  });
  page.on("requestfailed", (request) => pendingReads.delete(request));
  await page.goto(`/employer/pipelines/${ids.job}`);
  await expect(
    page.getByRole("button", { name: /^English check, 0 incoming/ }),
  ).toBeVisible();
  await expect(
    page.getByRole("heading", {
      name: "Hiring pipeline",
      exact: true,
      includeHidden: true,
    }),
  ).toHaveCount(0);
  const completedReads = () => ({
    pipeline: reads.filter((path) => path.endsWith("/pipeline")).length,
    stages: reads.filter((path) => path.endsWith("/stages")).length,
    modules: reads.filter((path) => path === "/v1/modules").length,
  });
  await expect
    .poll(completedReads)
    .toEqual({ pipeline: 1, stages: 1, modules: 1 });
  const summary = page.getByText("Legacy stage editor", {
    exact: true,
  });
  await summary.click();
  const editor = page.getByRole("region", { name: "Hiring pipeline" });
  const stage = editor.getByLabel("Stage name", { exact: true }).first();
  await stage.fill("Retained local edit");
  await expect(
    editor
      .getByRole("checkbox", { name: "Go fundamentals · Behavioral", exact: true })
      .first(),
  ).toBeVisible();
  // Opening the lazy component can replay existing effects in development.
  // Wait for its data before measuring whether a later reopen triggers reads.
  await expect.poll(() => pendingReads.size).toBe(0);
  const before = reads.length;
  await summary.click();
  await expect(stage).toBeHidden();
  await summary.click();
  await expect(stage).toBeVisible();
  await expect(stage).toHaveValue("Retained local edit");
  await expect.poll(() => pendingReads.size).toBe(0);
  expect(reads).toHaveLength(before);
});

test("a malformed API error still provides readable recovery", async ({
  page,
}) => {
  await productFixtures(page);
  await page.route("**/v1/jobs", (route) =>
    route.fulfill({
      status: 500,
      json: {
        detail: { private: "unexpected object" },
        code: ["invalid"],
        error: { message: 42 },
      },
    }),
  );
  await page.goto("/employer/jobs");
  const alert = page.locator("main").getByRole("alert");
  await expect(alert).toContainText(
    "We couldn’t complete your request. Please try again.",
  );
  await expect(alert).not.toContainText("[object Object]");
  await expect(
    alert.getByRole("button", { name: "Try again", exact: true }),
  ).toBeVisible();
});

test("blank decision reasons are described and focused", async ({ page }) => {
  await productFixtures(page);
  await page.goto(`/employer/pipelines/${ids.job}/candidates`);
  await page.getByRole("checkbox").first().check();
  await page.getByRole("button", { name: "Pass", exact: true }).click();
  const dialog = page.getByRole("dialog");
  const reason = dialog.getByRole("textbox", {
    name: "Decision reason",
    exact: true,
  });
  await reason.fill("   ");
  await dialog
    .getByRole("button", { name: "Confirm decision", exact: true })
    .click();
  await expect(reason).toBeFocused();
  await expect(reason).toHaveAttribute("aria-invalid", "true");
  await expect(reason).toHaveAccessibleDescription(
    "Enter a reason to continue.",
  );
});

test("late camera permission is released after leaving the video question", async ({
  page,
}) => {
  await productFixtures(page);
  await page.route(`**/v1/assessment-attempts/${ids.application}`, (route) =>
    route.fulfill({
      json: {
        id: ids.application,
        instanceId: ids.application,
        stageExecutionId: ids.application,
        status: "in_progress",
        deadline: "2027-01-01T00:00:00Z",
        autoSubmitted: false,
        questions: [
          {
            id: ids.module,
            prompt: "Describe your approach",
            responseFormat: "video",
            position: 0,
          },
          {
            id: ids.job,
            prompt: "Second prompt",
            responseFormat: "video",
            position: 1,
          },
        ],
      },
    }),
  );
  await page.addInitScript(() => {
    let resolvePermission: ((stream: MediaStream) => void) | undefined;
    let requests = 0,
      stops = 0;
    Object.defineProperty(navigator, "mediaDevices", {
      configurable: true,
      value: {
        getUserMedia: () => {
          requests += 1;
          return new Promise<MediaStream>((resolve) => {
            resolvePermission = resolve;
          });
        },
      },
    });
    // The stream is never attached: it must be stopped when the stale request resolves.
    const stream = {
      getTracks: () => [
        {
          stop: () => {
            stops += 1;
          },
        },
      ],
    } as unknown as MediaStream;
    Object.defineProperty(window, "cameraProbe", {
      value: {
        resolve: () => resolvePermission?.(stream),
        requests: () => requests,
        stops: () => stops,
      },
    });
  });
  await page.goto(`/candidate/assessments/${ids.application}/video`);
  const record = page.getByRole("button", { name: "Record", exact: true });
  await record.click();
  await expect(record).toBeDisabled();
  await page.getByRole("button", { name: "Next prompt", exact: true }).click();
  await page.evaluate(() => {
    const probe = Reflect.get(window, "cameraProbe") as { resolve: () => void };
    probe.resolve();
  });
  await expect
    .poll(() =>
      page.evaluate(() => {
        const probe = Reflect.get(window, "cameraProbe") as {
          requests: () => number;
          stops: () => number;
        };
        return { requests: probe.requests(), stops: probe.stops() };
      }),
    )
    .toEqual({ requests: 1, stops: 1 });
  await expect(
    page.getByRole("button", { name: "Record", exact: true }),
  ).toBeEnabled();
});
