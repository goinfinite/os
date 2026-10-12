import { expect, type Page, test } from "./fixtures";

// Lightpanda drops punctuation key events inside xterm, so the echo marker stays alphanumeric.
const marker = `osBrowserMarker${Date.now()}`;

const modalBody = (page: Page) => page.locator("#terminal-sessions-modal-body");
const modalCloseButton = (page: Page) =>
  page.locator("#terminal-sessions-modal button:has(i.ph-x)");
const modalRail = (page: Page) => modalBody(page).locator("aside");

async function createSession(rail: ReturnType<typeof modalRail>) {
  await rail.locator('button:has-text("new session")').click();
}

test.describe("terminal", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("overview/");
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
    await expect(modalRail(page)).toBeVisible();

    await createSession(modalRail(page));
    await expect(modalBody(page).locator(".xterm").first()).toBeVisible();

    const terminal = modalBody(page).locator(".xterm").first();
    await terminal.click();
    await page.keyboard.type(`echo ${marker}`);
    await page.keyboard.press("Enter");
    await expect(terminal).toContainText(marker);

    await modalCloseButton(page).click();
    await expect(modalBody(page)).not.toBeVisible();

    await page.locator("#terminal-sessions-footer-trigger").click();
    await expect(modalRail(page)).toBeVisible();

    const reattachedTerminal = modalBody(page).locator(".xterm").first();
    await expect(reattachedTerminal).toBeVisible();
    await expect(reattachedTerminal).toContainText(marker);

    await modalRail(page).locator('button[aria-label="Kill"]').first().click();
    await expect(modalBody(page).locator(".xterm")).toHaveCount(0);
    await expect(
      modalRail(page).locator('button[aria-label="Kill"]'),
    ).toHaveCount(0);
  });

  test("renames a session from the rail", async ({ page }) => {
    await page.goto("overview/");
    await page.locator("#terminal-sessions-footer-trigger").click();
    await createSession(modalRail(page));
    await expect(modalBody(page).locator(".xterm").first()).toBeVisible();

    await modalRail(page)
      .locator('button[aria-label="Rename"]')
      .first()
      .click();
    const renameInput = modalRail(page).locator(
      'input[placeholder="session name"]',
    );
    await expect(renameInput).toBeVisible();
    await expect(renameInput).toBeFocused();
    await renameInput.fill("opencode");
    await renameInput.press("Enter");

    await expect(
      modalRail(page).locator("span.truncate.font-medium").first(),
    ).toHaveText("opencode");

    const listResponse = await page.request.get("/api/v1/terminal-sessions/");
    const listBody = await listResponse.json();
    expect(listBody.body.terminalSessions[0].name).toBe("opencode");
  });

  test("creates a named session from the modal rail and kills it", async ({
    page,
  }) => {
    const uniqueName = `modal-rail-${Date.now()}`;

    await page.goto("overview/");
    await page.locator("#terminal-sessions-footer-trigger").click();
    await expect(modalRail(page)).toBeVisible();

    await modalRail(page)
      .locator('button[aria-label="Custom Session Settings"]')
      .click();
    await modalRail(page).locator('input[name="name"]').fill(uniqueName);
    await modalRail(page)
      .locator('button[type="submit"]:has-text("create session")')
      .click();

    await expect(modalBody(page).locator(".xterm").first()).toBeVisible();

    const sessionLabel = modalRail(page)
      .locator("span.truncate.font-medium")
      .filter({ hasText: uniqueName });
    await expect(sessionLabel).toHaveCount(1);

    await sessionLabel
      .locator("xpath=ancestor::div[contains(@class,'justify-between')][1]")
      .locator('button[aria-label="Kill"]')
      .click();
    await expect(sessionLabel).toHaveCount(0);
  });

  test("creates an account-user session from the custom form", async ({
    page,
  }) => {
    await page.goto("overview/");
    await page.locator("#terminal-sessions-footer-trigger").click();
    await expect(modalRail(page)).toBeVisible();

    await modalRail(page)
      .locator('button[aria-label="Custom Session Settings"]')
      .click();
    const accountUsername = process.env.OS_TEST_ACCOUNT_USERNAME as string;
    expect(accountUsername).toBeTruthy();
    const runAsForm = modalRail(page).locator("form");
    await runAsForm.locator('div[role="button"]').click();
    await runAsForm
      .locator("label")
      .filter({ hasText: accountUsername })
      .click();

    const workingDirInput = modalRail(page).locator('input[name="workingDir"]');
    await expect(workingDirInput).toHaveValue(`/home/${accountUsername}`);

    await modalRail(page)
      .locator('button[type="submit"]:has-text("create session")')
      .click();
    await expect(modalBody(page).locator(".xterm").first()).toBeVisible();

    const listResponse = await page.request.get("/api/v1/terminal-sessions/");
    const listBody = await listResponse.json();
    const createdSession = listBody.body.terminalSessions.find(
      (session: { runAsUsername: string }) =>
        session.runAsUsername === accountUsername,
    );
    expect(createdSession).toBeTruthy();
    expect(createdSession.workingDir).toBe(`/home/${accountUsername}`);
  });

  test("redirects the removed terminal page to the overview", async ({
    page,
  }) => {
    await page.goto("terminal/");
    await expect(page).toHaveURL(/\/overview\/$/);
  });

  test("retries the sessions modal load after a failed fragment fetch", async ({
    page,
  }) => {
    let shouldFail = true;
    await page.route("**/fragment/terminal-sessions/", (route) => {
      if (shouldFail) {
        shouldFail = false;
        return route.abort();
      }
      return route.continue();
    });

    await page.goto("overview/");
    await page.locator("#terminal-sessions-footer-trigger").click();
    await expect(modalRail(page)).toHaveCount(0);

    await modalCloseButton(page).click();
    await page.locator("#terminal-sessions-footer-trigger").click();
    await expect(modalRail(page)).toBeVisible();
  });

  test("closes the sessions modal from the backdrop", async ({ page }) => {
    await page.goto("overview/");
    await page.locator("#terminal-sessions-footer-trigger").click();
    await expect(modalRail(page)).toBeVisible();

    // Lightpanda has no hit-testing, so a real click cannot reach the backdrop overlay; a real click on the backdrop is verified by the browser-usage rendering loop.
    await page
      .locator("#terminal-sessions-modal .fixed.inset-0")
      .dispatchEvent("mousedown");
    await expect(modalBody(page)).not.toBeVisible();
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
