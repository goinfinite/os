import { expect, test } from "@playwright/test";

test.describe("accounts", () => {
  test("keeps the create account field labels inside the modal", async ({
    page,
  }) => {
    await page.goto("accounts/");
    await page.locator('button:has-text("create account")').first().click();

    const usernameInput = page.locator('input[name="username"]');
    await expect(usernameInput).toBeVisible();
    await usernameInput.fill("hello");

    const usernameLabel = page
      .locator("label:visible")
      .filter({ hasText: "Username" });
    await expect(usernameLabel).toBeVisible();

    const modalContent = page
      .locator('form:has(input[name="username"])')
      .locator("..");
    const labelBox = await usernameLabel.boundingBox();
    const contentBox = await modalContent.boundingBox();
    expect(labelBox).not.toBeNull();
    expect(contentBox).not.toBeNull();
    expect(labelBox?.y).toBeGreaterThanOrEqual((contentBox?.y ?? 0) - 1);
  });
});
