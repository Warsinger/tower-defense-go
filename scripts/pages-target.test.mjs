import { test } from 'node:test';
import assert from 'node:assert/strict';
import { deploymentBranch } from './pages-target.mjs';

const sha = 'a'.repeat(40);
const env = { GITHUB_REPOSITORY: 'Warsinger/tower-defense-go', GITHUB_SHA: sha,
  GITHUB_EVENT_NAME: 'workflow_dispatch', GITHUB_REF: 'refs/heads/preview', PREVIEW_PR: '13' };
const pr = { state: 'open', base: { ref: 'main' }, head: { sha, ref: 'preview', repo: { full_name: env.GITHUB_REPOSITORY } } };

test('only the current main revision publishes production', async () => {
  assert.equal(await deploymentBranch({ ...env, GITHUB_REF: 'refs/heads/main' }, async () => ({ object: { sha } })), 'main');
  await assert.rejects(deploymentBranch({ ...env, GITHUB_REF: 'refs/heads/main' }, async () => ({ object: { sha: 'b'.repeat(40) } })), /superseded/);
});
test('manual preview uses a PR-specific alias', async () => {
  assert.equal(await deploymentBranch(env, async () => pr), 'pr-13');
});
test('fork, closed, stale, and wrong-base previews are rejected', async () => {
  for (const modified of [
    { ...pr, state: 'closed' }, { ...pr, base: { ref: 'other' } },
    { ...pr, head: { ...pr.head, sha: 'b'.repeat(40) } },
    { ...pr, head: { ...pr.head, repo: { full_name: 'other/fork' } } },
  ]) await assert.rejects(deploymentBranch(env, async () => modified), /Preview must match/);
});
test('PR events and invalid preview numbers cannot publish', async () => {
  await assert.rejects(deploymentBranch({ ...env, GITHUB_EVENT_NAME: 'pull_request' }, async () => pr), /PR events/);
  await assert.rejects(deploymentBranch({ ...env, PREVIEW_PR: '../main' }, async () => pr), /PR number/);
});
