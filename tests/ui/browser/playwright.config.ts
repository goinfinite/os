import { defineConfig, devices } from "@playwright/test";

const baseURL = process.env.OS_TEST_API_URL || "https://127.0.0.1:1618";
const outputDir = process.env.OS_TEST_BROWSER_ARTIFACTS_DIR || "test-results";

export default defineConfig({
  testDir: "./specs",
  timeout: 60_000,
  expect: { timeout: 20_000 },
  fullyParallel: false,
  workers: 1,
  retries: 0,
  reporter: [["list"]],
  outputDir,
  use: {
    baseURL,
    ignoreHTTPSErrors: true,
    screenshot: "off",
    trace: "off",
  },
  projects: [
    {
      name: "setup",
      testMatch: /auth\.setup\.ts/,
    },
    {
      name: "lightpanda",
      use: {
        ...devices["Desktop Chrome"],
        // Lightpanda has no rendering or scroll support, so the tall viewport keeps click targets inside the interaction frame.
        viewport: { width: 1280, height: 100000 },
      },
      testIgnore: /visual\.spec\.ts/,
      dependencies: ["setup"],
    },
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
        storageState: "storageState.json",
        screenshot: "only-on-failure",
        trace: "retain-on-failure",
      },
      dependencies: ["setup"],
    },
    {
      name: "firefox",
      use: {
        ...devices["Desktop Firefox"],
        storageState: "storageState.json",
        screenshot: "only-on-failure",
        trace: "retain-on-failure",
      },
      dependencies: ["setup"],
    },
  ],
});
