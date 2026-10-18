import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { gunzipSync } from 'node:zlib';
import { packagePages } from './package-pages.mjs';

test('release paths bind launcher, runtime, and Wasm; changed assets change URLs', t => {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), 'tower-pages-'));
  t.after(() => {
    assert.equal(path.dirname(path.resolve(root)), path.resolve(os.tmpdir()));
    assert.ok(path.basename(root).startsWith('tower-pages-'));
    fs.rmSync(root, { recursive: true, force: true });
  });
  const source = path.join(root, 'source'); fs.mkdirSync(source);
  for (const name of ['game.html', 'launcher.js', 'game.js', 'wasm_exec.js', 'game.wasm']) fs.writeFileSync(path.join(source, name), name);
  fs.writeFileSync(path.join(source, 'index.html'), '<script src="launcher.js"></script>');
  const first = path.join(root, 'first');
  const hash = packagePages(source, first, 'abc');
  assert.match(fs.readFileSync(path.join(first, 'index.html'), 'utf8'), new RegExp(`assets/${hash}/launcher.js`));
  assert.equal(gunzipSync(fs.readFileSync(path.join(first, `assets/${hash}/game.wasm`))).toString(), 'game.wasm');
  assert.ok(fs.existsSync(path.join(first, '404.html')));
  assert.match(fs.readFileSync(path.join(first, '_headers'), 'utf8'), /immutable/);
  assert.equal(packagePages(source, path.join(root, 'same')), hash);
  fs.writeFileSync(path.join(source, 'wasm_exec.js'), 'new runtime');
  assert.notEqual(packagePages(source, path.join(root, 'new')), hash);
  assert.throws(() => packagePages(source, first), /EEXIST/);
});
