import {env} from 'node:process';
import {test, expect} from '@playwright/test';
import {login, apiCreateRepo, apiCreateIssue, randomString} from './utils.ts';

test('saved reply created in settings is inserted from the comment toolbar', async ({page, request}) => {
  const repoName = `e2e-saved-reply-${randomString(8)}`;
  const owner = env.GITEA_TEST_E2E_USER;
  await apiCreateRepo(request, {name: repoName, autoInit: false});
  await Promise.all([
    apiCreateIssue(request, {owner, repo: repoName, title: 'Saved reply test'}),
    login(page),
  ]);

  const title = `saved reply ${randomString(8)}`;
  const content = `Thanks for reporting ${randomString(8)}`;
  await page.goto('/user/settings/saved_replies');
  await page.getByLabel('Title').fill(title);
  await page.locator('textarea[name="content"]').fill(content);
  await page.getByRole('button', {name: 'Add Saved Reply'}).click();
  await expect(page.locator('.item-title').filter({hasText: title})).toBeVisible();

  await page.goto(`/${owner}/${repoName}/issues/1`);
  await page.getByRole('button', {name: 'Insert a saved reply'}).click();
  await page.getByPlaceholder('Filter saved replies').pressSequentially(title);
  await page.getByRole('button', {name: title}).press('Enter');
  await expect(page.getByPlaceholder('Leave a comment')).toHaveValue(content);
});
