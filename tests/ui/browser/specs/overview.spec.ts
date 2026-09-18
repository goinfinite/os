import { expect, test } from "@playwright/test";

test.describe("overview", () => {
  test("renders the two-column grid tiles", async ({ page }) => {
    await page.goto("overview/");
    await expect(page.locator("#marketplace")).toBeVisible();
    await expect(page.locator("#services")).toBeVisible();
    await expect(page.locator("#terminal-sessions")).toBeVisible();
    await expect(page.locator("#system-info")).toBeVisible();
    await expect(page.locator("#resource-usage")).toBeVisible();
    await expect(page.locator("#cpuAndMemoryUsageChart svg")).toBeVisible();
  });

  test("opens the terminal modal from the footer", async ({ page }) => {
    await page.goto("overview/");
    await page.locator("#terminal-sessions-footer-trigger").click();
    await expect(page.locator("#terminal-sessions-modal-body")).toBeVisible();
    await expect(
      page.locator(
        '#terminal-sessions-modal-body button[type="submit"]:has-text("create session")',
      ),
    ).toBeVisible();
  });
});
