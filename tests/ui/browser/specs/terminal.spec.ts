import { expect, test } from "@playwright/test";

const marker = `os-browser-marker-${Date.now()}`;

test.describe("terminal", () => {
  test.beforeEach(async ({ page }) => {
    const listResponse = await page.request.get("/api/v1/terminal-sessions/");
    const listBody = await listResponse.json();
    for (const session of listBody.body.terminalSessions) {
      await page.request.delete(`/api/v1/terminal-sessions/${session.id}/`);
    }
  });

  test("creates, attaches, echoes, persists and kills from the modal", async ({
    page,
  }) => {
    await page.goto("overview/");
    await page.locator("#terminal-sessions-footer-trigger").click();

    const modalBody = page.locator("#terminal-sessions-modal-body");
    await expect(modalBody).toBeVisible();

    await modalBody
      .locator('button[type="submit"]:has-text("create session")')
      .click();

    const terminal = modalBody.locator(".xterm").first();
    await expect(terminal).toBeVisible();
    await terminal.click();
    await page.keyboard.type(`echo ${marker}`);
    await page.keyboard.press("Enter");
    await expect(terminal).toContainText(marker);

    await page.locator("#terminal-sessions-modal-close").click();
    await expect(modalBody).not.toBeVisible();

    await page.locator("#terminal-sessions-footer-trigger").click();
    await expect(modalBody).toBeVisible();
    await modalBody.locator('button:has-text("attach")').first().click();

    const reattachedTerminal = modalBody.locator(".xterm").first();
    await expect(reattachedTerminal).toBeVisible();
    await expect(reattachedTerminal).toContainText(marker);

    await modalBody.locator('button:has-text("kill")').first().click();
    await expect(modalBody.locator(".xterm")).toHaveCount(0);
  });

  test("creates and kills on the terminal page", async ({ page }) => {
    await page.goto("terminal/");
    await page
      .locator('button[type="submit"]:has-text("create session")')
      .click();
    await expect(page.locator(".xterm").first()).toBeVisible();
    await page.locator('button:has-text("kill")').first().click();
    await expect(page.locator(".xterm")).toHaveCount(0);
  });

  test("rejects reading another account's sessions", async ({ page }) => {
    const otherAccountId = process.env.OS_TEST_OTHER_ACCOUNT_ID as string;
    const response = await page.request.get(
      `/api/v1/terminal-sessions/?accountId=${otherAccountId}`,
    );
    expect(response.status()).toBe(400);
    const body = await response.json();
    expect(body.body).toBe("TerminalSessionAccountNotOwned");
  });
});
