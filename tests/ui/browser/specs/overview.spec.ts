import { expect, test } from "./fixtures";

test.describe("overview", () => {
  test("renders the paired cards and the services table", async ({ page }) => {
    await page.goto("overview/");
    await expect(page.locator("#marketplace")).toBeVisible();
    await expect(page.locator("#services")).toBeVisible();
    await expect(page.locator("#terminal-sessions")).toBeVisible();
    await expect(page.locator("#system-info")).toBeVisible();
    await expect(page.locator("#resource-usage")).toBeVisible();
    await expect(page.locator("#cpuAndMemoryUsageChart svg")).toBeVisible();
    await expect(page.locator("#installed-services-table")).toBeVisible();
    await expect(page.locator("#terminal-sessions-table")).toBeVisible();
    await expect(page.locator("#marketplace-items-table")).toBeVisible();
  });

  test("opens the sessions modal from the footer", async ({ page }) => {
    await page.goto("overview/");
    await page.locator("#terminal-sessions-footer-trigger").click();
    await expect(
      page.locator("#terminal-sessions-modal-body aside"),
    ).toBeVisible();
  });

  test("creates a session with defaults from the terminal sessions tile", async ({
    page,
  }) => {
    await page.goto("overview/");
    await page
      .locator('#terminal-sessions button:has-text("new session")')
      .first()
      .click();
    await expect(
      page.locator("#terminal-sessions-modal-body aside"),
    ).toBeVisible();
    await expect(
      page.locator("#terminal-sessions-modal-body .xterm").first(),
    ).toBeVisible();

    const listResponse = await page.request.get("/api/v1/terminal-sessions/");
    const listBody = await listResponse.json();
    for (const session of listBody.body.terminalSessions) {
      await page.request.delete(`/api/v1/terminal-sessions/${session.id}/`);
    }
  });

  test("attaches to a session from the terminal sessions tile row", async ({
    page,
  }) => {
    const createResponse = await page.request.post(
      "/api/v1/terminal-sessions/",
      {
        data: { workingDir: "/app" },
      },
    );
    expect(createResponse.status()).toBe(201);
    const createdBody = await createResponse.json();
    const sessionId = createdBody.body.id as string;

    await page.goto("overview/");
    const sessionRow = page
      .locator("#terminal-sessions-table tbody tr")
      .filter({ hasText: sessionId });
    await expect(sessionRow).toBeVisible();
    await sessionRow.locator('button[aria-label="Attach"]').click();

    await expect(
      page.locator("#terminal-sessions-modal-body aside"),
    ).toBeVisible();
    await expect(
      page.locator("#terminal-sessions-modal-body .xterm").first(),
    ).toBeVisible();

    await page.request.delete(`/api/v1/terminal-sessions/${sessionId}/`);
  });

  test("kills a session from the terminal sessions tile row", async ({
    page,
  }) => {
    const createResponse = await page.request.post(
      "/api/v1/terminal-sessions/",
      {
        data: { workingDir: "/app" },
      },
    );
    expect(createResponse.status()).toBe(201);
    const createdBody = await createResponse.json();
    const sessionId = createdBody.body.id as string;

    const createdListResponse = await page.request.get(
      "/api/v1/terminal-sessions/",
    );
    const createdListBody = await createdListResponse.json();
    expect(
      createdListBody.body.terminalSessions.map(
        (session: { id: string }) => session.id,
      ),
    ).toContain(sessionId);

    await page.goto("overview/");
    const sessionRow = page
      .locator("#terminal-sessions-table tbody tr")
      .filter({ hasText: sessionId });
    await expect(sessionRow).toBeVisible();
    await sessionRow.locator('button[aria-label="Kill"]').click();

    await expect(
      page
        .locator("#terminal-sessions-table tbody tr")
        .filter({ hasText: sessionId }),
    ).toHaveCount(0);

    const listResponse = await page.request.get("/api/v1/terminal-sessions/");
    const listBody = await listResponse.json();
    expect(
      listBody.body.terminalSessions.filter(
        (session: { id: string }) => session.id === sessionId,
      ),
    ).toHaveLength(0);
  });

  test("seeds the table page from the query params on a full load", async ({
    page,
  }) => {
    for (let index = 1; index <= 7; index++) {
      const createResponse = await page.request.post(
        "/api/v1/terminal-sessions/",
        { data: { name: `pgtest-${index}`, workingDir: "/app" } },
      );
      expect(createResponse.status()).toBe(201);
    }

    try {
      await page.goto(
        "overview/?pageNumber=1&itemsPerPage=5&sortBy=name&sortDirection=asc",
      );

      const table = page.locator("#terminal-sessions-table");
      await expect(
        table.locator("tbody tr").filter({ hasText: "pgtest-6" }),
      ).toBeVisible();
      await expect(
        table.locator("tbody tr").filter({ hasText: "pgtest-7" }),
      ).toBeVisible();
      await expect(table.locator("tbody tr")).toHaveCount(2);
      await expect(table.locator("p.select-none")).toHaveText("6–7 of 7");
    } finally {
      const listResponse = await page.request.get(
        "/api/v1/terminal-sessions/?itemsPerPage=50",
      );
      const listBody = await listResponse.json();
      for (const session of listBody.body.terminalSessions) {
        if (session.name?.startsWith("pgtest")) {
          await page.request.delete(`/api/v1/terminal-sessions/${session.id}/`);
        }
      }
    }
  });

  test("falls back to the table defaults on malformed query params", async ({
    page,
  }) => {
    await page.goto(
      "overview/?pageNumber=abc&itemsPerPage=-5&sortBy=!!&sortDirection=sideways",
    );

    const table = page.locator("#terminal-sessions-table");
    await expect(table).toBeVisible();
    await expect(table.locator("p.select-none")).toBeVisible();
    const title = await page.title();
    expect(title).toBe("Infinite OS");
  });

  test("switches the service installation form to custom", async ({ page }) => {
    await page.goto("overview/");
    await page
      .locator('#services button:has-text("install service")')
      .first()
      .click();

    const installableLabel = page
      .locator("label:visible")
      .filter({ hasText: "Installable" });
    const customLabel = page
      .locator("label:visible")
      .filter({ hasText: "Custom" });
    await expect(installableLabel).toBeVisible();
    await expect(customLabel).toBeVisible();

    await customLabel.click();

    const installationForm = page.locator("form:has(#install-service-button)");
    await expect(installationForm.locator('input[name="name"]')).toBeVisible();
    await expect(
      installationForm.locator('input[name="startCmd"]'),
    ).toBeVisible();
  });
});
