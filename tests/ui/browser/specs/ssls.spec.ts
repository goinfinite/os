import { expect, test } from "@playwright/test";

test.describe("ssls", () => {
  test("keeps the import certificate field labels inside the modal", async ({
    page,
  }) => {
    await page.goto("ssls/");
    await page
      .locator('button:has-text("import ssl certificate")')
      .first()
      .click();

    const hostnamesFieldset = page
      .locator("form#import-ssl-certificate-form fieldset")
      .first();
    await expect(hostnamesFieldset).toBeVisible();

    const modalContent = page
      .locator("form#import-ssl-certificate-form")
      .locator("..");
    const fieldBox = await hostnamesFieldset.boundingBox();
    const contentBox = await modalContent.boundingBox();
    expect(fieldBox).not.toBeNull();
    expect(contentBox).not.toBeNull();
    expect(fieldBox?.y).toBeGreaterThanOrEqual((contentBox?.y ?? 0) - 1);
  });
});
