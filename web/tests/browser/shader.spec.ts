import { expect, test } from "@playwright/test";

test("the welcome shader renders its authored WebGPU output", async ({
  page,
}) => {
  test.setTimeout(120000);
  const errors: string[] = [];
  page.on("pageerror", (error) => errors.push(error.message));
  page.on("console", (message) => {
    if (message.type() === "error") errors.push(message.text());
  });
  await page.goto("/");
  const canvas = page.locator("canvas");
  const supported = await page.evaluate(
    async () => !!(navigator.gpu && (await navigator.gpu.requestAdapter())),
  );
  test.skip(
    !supported,
    "This browser has no WebGPU adapter; the base orange surface remains visible.",
  );
  await expect(canvas).toHaveAttribute("data-ready", "true", {
    timeout: 90000,
  });
  await expect(canvas).not.toHaveAttribute("data-error");
  await expect.poll(() => errors).toEqual([]);
  await page.screenshot({ path: "tests/artifacts/welcome-1728-viewport.png" });
});
