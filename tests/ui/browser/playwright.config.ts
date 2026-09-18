import { defineConfig, devices } from "@playwright/test";

const baseURL = process.env.OS_TEST_API_URL || "https://127.0.0.1:1618";
const outputDir = process.env.OS_TEST_BROWSER_ARTIFACTS_DIR || "test-results";

export default defineConfig({
  testDir: "./specs",
  timeout: 60_000,
  expect: { timeout: 20_000 },
  fullyParallel: false,
  retries: 0,
  reporter: [["list"]],
  outputDir,
  use: {
    baseURL,
    ignoreHTTPSErrors: true,
    screenshot: "only-on-failure",
    trace: "retain-on-failure",
  },
  projects: [
    {
      name: "setup",
      testMatch: /auth\.setup\.ts/,
    },
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"], storageState: "storageState.json" },
      dependencies: ["setup"],
    },
  ],
});
