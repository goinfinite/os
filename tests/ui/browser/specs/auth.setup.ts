import { expect, test as setup } from "@playwright/test";

const username = process.env.OS_TEST_ACCOUNT_USERNAME as string;
const password = process.env.OS_TEST_ACCOUNT_PASSWORD as string;

setup("authenticate", async ({ page }) => {
  await page.goto("login/");
  await page.locator('input[name="username"]').fill(username);
  await page.locator('input[name="password"]').fill(password);
  await page.locator('button[type="submit"]').click();
  await page.waitForURL("**/overview/");
  await expect(page.locator("#marketplace")).toBeVisible();
  await page.context().storageState({ path: "storageState.json" });
});
