import { expect, test } from "./fixtures";

test.describe("visual", () => {
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

  test("keeps the marketplace installation modal free of scroll overflow", async ({
    page,
  }) => {
    await page.goto("overview/");
    await page
      .locator('#marketplace button:has-text("install catalog item")')
      .first()
      .click();

    const modalContent = page
      .locator("#schedule-marketplace-item-installation-form")
      .locator("..");
    const lastVisibleCard = page
      .locator("#schedule-marketplace-item-installation-form .group:visible")
      .last();
    await lastVisibleCard.hover();

    await expect(page.locator("body > [role=tooltip]:visible")).toHaveCount(1);
    await expect
      .poll(() =>
        modalContent.evaluate(
          (element) => element.scrollWidth - element.clientWidth,
        ),
      )
      .toBe(0);
    await expect
      .poll(() =>
        modalContent.evaluate(
          (element) => element.scrollHeight - element.clientHeight,
        ),
      )
      .toBe(0);
  });
});
