import { pagesProject } from './pages-project.mjs';
import { deploymentBranch } from './pages-target.mjs';
import fs from 'node:fs';
import { spawnSync } from 'node:child_process';

const project = 'tower-def';
const account = process.env.CLOUDFLARE_ACCOUNT_ID;
const token = process.env.CLOUDFLARE_API_TOKEN;
const repo = 'Warsinger/tower-defense-go';
const sha = process.env.GITHUB_SHA;
if (!account || !token) throw new Error('Set CLOUDFLARE_ACCOUNT_ID and CLOUDFLARE_API_TOKEN repository secrets');
if (process.env.GITHUB_REPOSITORY !== repo || !/^[a-f0-9]{40}$/.test(sha || '')) throw new Error('Deployment requires this repository workflow context');
if (!['push', 'workflow_dispatch'].includes(process.env.GITHUB_EVENT_NAME)) throw new Error('PR events cannot publish');

async function github(endpoint) {
  const response = await fetch(`https://api.github.com/repos/${repo}/${endpoint}`, {
    headers: { Authorization: `Bearer ${process.env.GITHUB_TOKEN}`, Accept: 'application/vnd.github+json' },
  });
  if (!response.ok) throw new Error(`GitHub validation failed: HTTP ${response.status}`);
  return response.json();
}

const branch = await deploymentBranch(process.env, github);

const config = await pagesProject(fetch, `https://api.cloudflare.com/client/v4/accounts/${encodeURIComponent(account)}/pages/projects`,
  { Authorization: `Bearer ${token}` }, process.env.GITHUB_EVENT_NAME === 'workflow_dispatch' && process.env.CREATE_PROJECT === 'true');
const hostname = config.subdomain;
if (hostname !== 'tower-def.pages.dev') throw new Error('Unexpected Pages hostname');
const version = JSON.parse(fs.readFileSync('build/pages/version.json', 'utf8'));
if (version.revision !== sha) throw new Error('Artifact revision does not match the checked workflow revision');
const command = process.platform === 'win32' ? 'npx.cmd' : 'npx';
const publish = spawnSync(command, ['--yes', 'wrangler@4.151.1', 'pages', 'deploy', 'build/pages',
  '--project-name', project, '--branch', branch, '--commit-hash', sha, '--commit-dirty=false'],
  { env: process.env, stdio: 'inherit', windowsHide: true });
if (publish.error || publish.status !== 0) throw new Error('Cloudflare upload failed');
const url = branch === 'main' ? `https://${hostname}/` : `https://${branch}.${hostname}/`;
console.log(`Open the game: ${url}`);
if (process.env.GITHUB_STEP_SUMMARY) fs.appendFileSync(process.env.GITHUB_STEP_SUMMARY,
  `## Tower Defense deployment\n\n[Open the game](${url})\n\nRevision: \`${sha}\`\n`);
