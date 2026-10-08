import { test as base, expect } from "@playwright/test";
import { chromium } from "playwright";

const cdpEndpoint = process.env.OS_TEST_LIGHTPANDA_CDP_URL;

const test = cdpEndpoint
  ? base.extend({
      browser: [
        // biome-ignore lint/correctness/noEmptyPattern: Playwright requires the empty destructuring pattern for a fixture function's first argument
        async ({}, use) => {
          const browser = await chromium.connectOverCDP(cdpEndpoint);
          await use(browser);
          await browser.close();
        },
        { scope: "worker" },
      ],
    })
  : base;

export type { Page } from "@playwright/test";
export { expect, test };
