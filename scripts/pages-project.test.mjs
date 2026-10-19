import { test } from 'node:test';
import assert from 'node:assert/strict';
import { pagesProject } from './pages-project.mjs';
const valid = { name: 'tower-def', production_branch: 'main', source: null };
const result = () => new Response(JSON.stringify({ success: true, result: valid }));
test('existing project requires no mutation', async () => {
  let calls = 0;
  assert.deepEqual(await pagesProject(async (_, options) => { calls++; assert.equal(options.method, undefined); return result(); }, 'https://example.test/projects', {}, true), valid);
  assert.equal(calls, 1);
});
test('only explicit setup creates the exact missing Direct Upload project', async () => {
  let calls = 0;
  await pagesProject(async (_, options) => {
    if (++calls === 1) return new Response('', { status: 404 });
    assert.equal(options.method, 'POST');
    assert.deepEqual(JSON.parse(options.body), { name: 'tower-def', production_branch: 'main' });
    return result();
  }, 'https://example.test/projects', {}, true);
  assert.equal(calls, 2);
});
test('missing projects without setup and access failures never trigger creation', async () => {
  for (const [status, allow] of [[404, false], [401, true], [403, true]]) {
    let calls = 0;
    await assert.rejects(pagesProject(async () => { calls++; return new Response('', { status }); }, 'https://example.test/projects', {}, allow), /HTTP/);
    assert.equal(calls, 1);
  }
});
