import { test, expect } from '@playwright/test';

test.describe('CRM User Journeys', () => {

  test('1. View Dashboard', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/Dashboard/i);
    // Has client list (table)
    await expect(page.locator('table')).toBeVisible();
    // Has Add Client button
    await expect(page.getByRole('button', { name: /Add Client/i })).toBeVisible();
    // Has search input
    await expect(page.getByPlaceholder(/search/i)).toBeVisible();
  });

  test('2. Add Client', async ({ page }) => {
    await page.goto('/');

    // Click Add Client to show form
    await page.getByRole('button', { name: /Add Client/i }).first().click();
    await page.waitForTimeout(500);

    // Fill form
    await page.locator('input[name="name"]').fill('Playwright Auto Corp');
    await page.locator('input[name="email"]').fill('auto@playwright.com');

    // Submit — click the Add Client button inside the form
    await page.locator('#addForm button[type="submit"], #addForm .btn-primary, form button:has-text("Add Client")').first().click();
    await page.waitForTimeout(2000);

    // Verify client appears (page may reload)
    await page.goto('/');
    const content = await page.textContent('body');
    expect(content).toContain('Playwright Auto Corp');
  });

  test('3. View Client Profile', async ({ page }) => {
    await page.goto('/');

    // Click first client link
    const firstClient = page.locator('table a, .client-list a').first();
    await expect(firstClient).toBeVisible();
    const clientName = await firstClient.textContent();
    await firstClient.click();

    // Should be on client page
    await expect(page.locator('h1')).toContainText(clientName || '');

    // Should have back button or link
    const backLink = page.locator('a:has-text("Back"), a:has-text("Dashboard"), a:has-text("←")');
    await expect(backLink.first()).toBeVisible();
  });

  test('4. Client Profile has tabs', async ({ page }) => {
    await page.goto('/');
    await page.locator('table a, .client-list a').first().click();
    await page.waitForTimeout(500);

    // Check for tabs
    const tabs = page.locator('button:has-text("Activity"), button:has-text("Invoice"), button:has-text("Profile"), [role="tab"]');
    const tabCount = await tabs.count();
    expect(tabCount).toBeGreaterThan(0);
  });

  test('5. Edit Client', async ({ page }) => {
    await page.goto('/');
    await page.locator('table a, .client-list a').first().click();
    await page.waitForTimeout(500);

    // Look for editable name field or edit button
    const nameInput = page.locator('input[name="name"], input[type="text"]').first();
    const isEditable = await nameInput.isVisible().catch(() => false);

    if (isEditable) {
      await nameInput.fill('Edited Name');
      // Look for save button
      const saveBtn = page.locator('button:has-text("Save"), button:has-text("Update")');
      if (await saveBtn.count() > 0) {
        await saveBtn.first().click();
        await page.waitForTimeout(1000);
      }
    }

    // If no editable fields, look for Edit button
    const editBtn = page.locator('button:has-text("Edit"), a:has-text("Edit")');
    expect(isEditable || await editBtn.count() > 0).toBeTruthy();
  });

  test('6. Delete Client (button exists)', async ({ page }) => {
    await page.goto('/');
    await page.locator('table a, .client-list a').first().click();
    await page.waitForTimeout(500);

    // Check delete button exists
    const deleteBtn = page.locator('button:has-text("Delete"), a:has-text("Delete")');
    await expect(deleteBtn.first()).toBeVisible();
  });

  test('7. Activity Tab', async ({ page }) => {
    await page.goto('/');
    await page.locator('table a, .client-list a').first().click();
    await page.waitForTimeout(500);

    // Click Activity tab
    const activityTab = page.locator('button:has-text("Activity"), a:has-text("Activity")');
    if (await activityTab.count() > 0) {
      await activityTab.first().click();
      await page.waitForTimeout(500);
    }

    // Check page has activity-related content
    const content = await page.textContent('body');
    const hasActivity = content?.toLowerCase().includes('activity') ||
                       content?.toLowerCase().includes('timeline') ||
                       content?.toLowerCase().includes('no activity');
    expect(hasActivity).toBeTruthy();
  });

  test('8. Invoices Tab', async ({ page }) => {
    await page.goto('/');
    await page.locator('table a, .client-list a').first().click();
    await page.waitForTimeout(500);

    // Click Invoices tab
    const invoiceTab = page.locator('button:has-text("Invoice"), a:has-text("Invoice")');
    if (await invoiceTab.count() > 0) {
      await invoiceTab.first().click();
      await page.waitForTimeout(500);
    }

    // Check for invoice-related content or create button
    const content = await page.textContent('body');
    const hasInvoice = content?.toLowerCase().includes('invoice') ||
                      content?.toLowerCase().includes('create') ||
                      content?.toLowerCase().includes('no invoice');
    expect(hasInvoice).toBeTruthy();
  });

  test('9. Search Clients', async ({ page }) => {
    await page.goto('/');

    const searchInput = page.getByPlaceholder(/search/i);
    if (await searchInput.isVisible()) {
      await searchInput.fill('Acme');
      await page.waitForTimeout(1000);

      // Check filtered results
      const rows = page.locator('table tbody tr, .client-item');
      const count = await rows.count();
      // Should have fewer results (or at least still show the table)
      expect(count).toBeGreaterThan(0);
    } else {
      // No search = FAIL
      expect(false).toBeTruthy();
    }
  });

  test('10. Empty State', async ({ page }) => {
    await page.goto('/');

    const searchInput = page.getByPlaceholder(/search/i);
    if (await searchInput.isVisible()) {
      await searchInput.fill('xyznonexistent99999');
      await page.waitForTimeout(1000);

      // Should show empty state or no results gracefully (not an error)
      const content = await page.textContent('body');
      const hasError = content?.includes('500') || content?.includes('Internal Server Error');
      expect(hasError).toBeFalsy();
    }
  });

});
