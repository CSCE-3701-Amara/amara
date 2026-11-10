import { defineConfig } from "@playwright/test";

const port = process.env.PLAYWRIGHT_PORT || "3100";
const baseURL = `http://127.0.0.1:${port}`;
const production = process.env.PLAYWRIGHT_SERVER === "production";

export default defineConfig({
  testDir: "./tests/browser",
  fullyParallel: true,
  workers: 2,
  use: {
    baseURL,
    viewport: { width: 1728, height: 1117 },
    reducedMotion: "reduce",
    launchOptions: {
      args: [
        "--enable-unsafe-webgpu",
        "--enable-features=Vulkan",
        "--use-angle=vulkan",
        "--use-vulkan=swiftshader",
        "--use-webgpu-adapter=swiftshader",
        "--disable-vulkan-surface",
      ],
    },
    trace: "retain-on-failure",
  },
  webServer: {
    env: { AMARA_E2E: "1" },
    command: `npm run ${production ? "start" : "dev"} -- --hostname 127.0.0.1 --port ${port}`,
    url: baseURL,
    reuseExistingServer: !production && !process.env.CI,
  },
});
