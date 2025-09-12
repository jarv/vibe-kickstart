import { test, expect } from "@playwright/test";

test.describe("WebSocket Connection Tests", () => {
  test("WebSocket connection is established on page load", async ({ page }) => {
    const wsPromise = page.waitForEvent("websocket");

    await page.goto("/");

    const ws = await wsPromise;

    // Verify WebSocket URL is correct
    expect(ws.url()).toMatch(/ws:\/\/localhost:8910\/ws/);

    // Verify WebSocket is not closed
    expect(ws.isClosed()).toBeFalsy();
  });
});
