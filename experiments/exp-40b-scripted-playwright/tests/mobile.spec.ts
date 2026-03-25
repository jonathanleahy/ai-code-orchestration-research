import { test, expect } from '@playwright/test';

test.describe('Mobile Responsive Tests', () => {

  test.use({ viewport: { width: 375, height: 812 } }); // iPhone X

  test('Dashboard renders on mobile', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('h1')).toBeVisible();
    // Client list should still be visible
    const content = await page.textContent('body');
    expect(content?.length).toBeGreaterThan(50);
  });

  test('Client list visible on mobile', async ({ page }) => {
    await page.goto('/');
    // Table or list should be scrollable, not hidden
    const table = page.locator('table, .client-list, [class*="client"]');
    await expect(table.first()).toBeVisible();
  });

  test('View button clickable on mobile', async ({ page }) => {
    await page.goto('/');
    // Find a View link/button
    const viewBtn = page.locator('a:has-text("View"), button:has-text("View"), table a').first();
    const isVisible = await viewBtn.isVisible().catch(() => false);

    if (isVisible) {
      // Check it's not clipped
      const box = await viewBtn.boundingBox();
      expect(box).not.toBeNull();
      if (box) {
        expect(box.width).toBeGreaterThan(10);
        expect(box.x + box.width).toBeLessThanOrEqual(375); // Within viewport
      }
    } else {
      // View button not visible on mobile — THIS IS THE BUG
      expect(isVisible).toBeTruthy();
    }
  });

  test('Add Client works on mobile', async ({ page }) => {
    await page.goto('/');
    const addBtn = page.getByRole('button', { name: /Add Client/i }).first();
    await expect(addBtn).toBeVisible();

    // Check button is within viewport
    const box = await addBtn.boundingBox();
    expect(box).not.toBeNull();
    if (box) {
      expect(box.x + box.width).toBeLessThanOrEqual(375);
    }
  });

  test('Client profile on mobile', async ({ page }) => {
    await page.goto('/');
    // Try to navigate to a client
    const link = page.locator('table a, .client-list a').first();
    if (await link.isVisible()) {
      await link.click();
      await page.waitForTimeout(500);
      // Page should load
      const content = await page.textContent('body');
      expect(content?.length).toBeGreaterThan(50);
    }
  });

  test('Table doesn\'t overflow viewport', async ({ page }) => {
    await page.goto('/');
    const table = page.locator('table').first();
    if (await table.isVisible()) {
      const box = await table.boundingBox();
      if (box) {
        // Table wider than viewport = horizontal scroll = bad UX
        // Allow some overflow but flag if way too wide
        expect(box.width).toBeLessThanOrEqual(400); // 375 + some tolerance
      }
    }
  });

});
