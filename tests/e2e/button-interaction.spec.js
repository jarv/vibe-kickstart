import { test, expect } from "@playwright/test";

test.describe("Button Interaction Tests", () => {
  test("button click sends reset WebSocket message", async ({ page }) => {
    const wsPromise = page.waitForEvent("websocket");

    await page.goto("/");
    await page.waitForLoadState("networkidle");

    const ws = await wsPromise;
    const button = page.locator("#counter-button");

    // Wait for connection to be established
    await expect(button).toHaveText(/\d+ seconds since the last press/);

    // Set up listener for WebSocket messages
    const messagePromise = ws.waitForEvent("framesent");

    // Click the button
    await button.click();

    // Verify WebSocket message was sent
    const frameData = await messagePromise;
    const message = JSON.parse(frameData.payload);
    expect(message).toEqual({ type: "reset" });
  });

  test("counter resets to 0 after button click", async ({ page }) => {
    const wsPromise = page.waitForEvent("websocket");

    await page.goto("/");
    await page.waitForLoadState("networkidle");

    const ws = await wsPromise;
    const button = page.locator("#counter-button");

    // Wait for connection and some time to pass
    await expect(button).toHaveText(/\d+ seconds since the last press/);

    // Wait a bit to ensure counter has incremented
    await page.waitForTimeout(2000);

    // Click button to reset
    await button.click();

    // Counter should reset to 0 (or very low number due to timing)
    await expect(button).toHaveText(/^[0-2] seconds since the last press$/);
  });

  test("multiple button clicks work correctly", async ({ page }) => {
    const wsPromise = page.waitForEvent("websocket");

    await page.goto("/");
    await page.waitForLoadState("networkidle");

    const ws = await wsPromise;
    const button = page.locator("#counter-button");

    // Wait for connection
    await expect(button).toHaveText(/\d+ seconds since the last press/);

    // Click button multiple times
    for (let i = 0; i < 3; i++) {
      const messagePromise = ws.waitForEvent("framesent");
      await button.click();

      // Verify message was sent
      const frameData = await messagePromise;
      const message = JSON.parse(frameData.payload);
      expect(message).toEqual({ type: "reset" });

      // Wait for counter to reset
      await expect(button).toHaveText(/^[0-2] seconds since the last press$/);

      // Wait a bit before next click
      await page.waitForTimeout(500);
    }
  });

  test("WebSocket receives counter updates from server", async ({ page }) => {
    const wsPromise = page.waitForEvent("websocket");

    await page.goto("/");
    await page.waitForLoadState("networkidle");

    const ws = await wsPromise;
    const button = page.locator("#counter-button");

    // Wait for initial connection
    await expect(button).toHaveText(/\d+ seconds since the last press/);

    // Listen for incoming WebSocket messages
    const messagePromise = ws.waitForEvent("framereceived");

    // Wait for server to send an update
    const frameData = await messagePromise;
    const message = JSON.parse(frameData.payload);

    // Should receive update message with counter
    expect(message).toHaveProperty("type", "update");
    expect(message).toHaveProperty("counter");
    expect(typeof message.counter).toBe("number");

    // Button text should reflect the received counter
    await expect(button).toHaveText(
      `${message.counter} seconds since the last press`,
    );
  });
});

