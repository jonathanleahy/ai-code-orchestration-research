import { test, expect } from '@playwright/test';

test.describe('Console Error Monitoring', () => {

  test('Dashboard has no console errors', async ({ page }) => {
    const errors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') errors.push(msg.text());
    });
    page.on('pageerror', err => errors.push(err.message));

    await page.goto('/');
    await page.waitForTimeout(2000);

    // Filter out favicon 404 (harmless)
    const realErrors = errors.filter(e => !e.includes('favicon'));
    expect(realErrors).toEqual([]);
  });

  test('Add Client form has no console errors', async ({ page }) => {
    const errors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') errors.push(msg.text());
    });
    page.on('pageerror', err => errors.push(err.message));

    await page.goto('/');
    // Click Add Client
    await page.getByRole('button', { name: /Add Client/i }).first().click();
    await page.waitForTimeout(500);

    // Fill and submit
    await page.locator('input[name="name"]').fill('Console Test');
    await page.locator('input[name="email"]').fill('console@test.com');

    // Handle potential alert dialog
    page.on('dialog', async dialog => await dialog.accept());

    await page.locator('#addForm button[type="submit"], #addForm .btn-primary, form button:has-text("Add Client")').first().click();
    await page.waitForTimeout(3000);

    const realErrors = errors.filter(e => !e.includes('favicon'));
    expect(realErrors).toEqual([]);
  });

  test('Client profile has no console errors', async ({ page }) => {
    const errors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') errors.push(msg.text());
    });
    page.on('pageerror', err => errors.push(err.message));

    await page.goto('/');
    await page.locator('table a').first().click();
    await page.waitForTimeout(1000);

    // Click through tabs
    const tabs = page.locator('button:has-text("Activity"), button:has-text("Invoice"), button:has-text("Profile")');
    const count = await tabs.count();
    for (let i = 0; i < count; i++) {
      await tabs.nth(i).click();
      await page.waitForTimeout(500);
    }

    const realErrors = errors.filter(e => !e.includes('favicon'));
    expect(realErrors).toEqual([]);
  });

  test('Search has no console errors', async ({ page }) => {
    const errors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') errors.push(msg.text());
    });
    page.on('pageerror', err => errors.push(err.message));

    await page.goto('/');
    const search = page.getByPlaceholder(/search/i);
    if (await search.isVisible()) {
      await search.fill('test');
      await page.waitForTimeout(1000);
      await search.fill('');
      await page.waitForTimeout(1000);
      await search.fill('nonexistent999');
      await page.waitForTimeout(1000);
    }

    const realErrors = errors.filter(e => !e.includes('favicon'));
    expect(realErrors).toEqual([]);
  });

  test('No network errors on page load', async ({ page }) => {
    const failedRequests: string[] = [];
    page.on('requestfailed', request => {
      if (!request.url().includes('favicon')) {
        failedRequests.push(`${request.method()} ${request.url()} - ${request.failure()?.errorText}`);
      }
    });

    await page.goto('/');
    await page.waitForTimeout(2000);

    expect(failedRequests).toEqual([]);
  });

});
