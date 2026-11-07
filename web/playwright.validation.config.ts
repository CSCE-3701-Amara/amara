import base from "./playwright.config";

export default {
  ...base,
  use: {
    ...base.use,
    launchOptions: {
      ...base.use.launchOptions,
      executablePath:
        "/home/yusif/.cache/ms-playwright/chromium-1254/chrome-linux64/chrome",
    },
  },
};
