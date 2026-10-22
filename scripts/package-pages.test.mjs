import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { gunzipSync } from 'node:zlib';
import { packagePages } from './package-pages.mjs';

test('real HTML loads the playfield and runtime from one versioned release', t => {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), 'tower-pages-'));
  t.after(() => {
    assert.equal(path.dirname(path.resolve(root)), path.resolve(os.tmpdir()));
    assert.ok(path.basename(root).startsWith('tower-pages-'));
    fs.rmSync(root, { recursive: true, force: true });
  });
  const source = path.join(root, 'source'); fs.mkdirSync(source);
  const icons = ['favicon.ico', 'icon-16.png', 'icon-32.png', 'icon-180.png', 'icon-192.png', 'icon-512.png'];
  for (const name of ['index.html', 'game.html', 'launcher.js', 'game.js', 'styles.css']) {
    fs.copyFileSync(new URL(`../web/${name}`, import.meta.url), path.join(source, name));
  }
  fs.cpSync(new URL('../node_modules/@warsinger/web-shared/browser', import.meta.url), path.join(source, 'shared'), {recursive:true});
  for (const name of ['wasm_exec.js', 'game.wasm', ...icons]) fs.writeFileSync(path.join(source, name), name);
  const first = path.join(root, 'first');
  const hash = packagePages(source, first, 'abc');
  assert.match(fs.readFileSync(path.join(first, 'index.html'), 'utf8'), new RegExp(`assets/${hash}/launcher.js`));
  for (const folder of [source, first]) {
    const index = fs.readFileSync(path.join(folder, 'index.html'), 'utf8');
    assert.doesNotMatch(index, /id="launch"|Launch game|src="[^\"]*""/);
    const frame = index.match(/<iframe\b[^>]*\bdata-src="([^"]+)"[^>]*><\/iframe>/);
    assert.ok(frame, 'one well-formed iframe is started by the shared session launcher');
    assert.doesNotMatch(index, /<iframe[^>]*\ssrc=/);
    assert.equal((index.match(/<iframe\b/g) || []).length, 1);
    assert.ok(index.indexOf('src="' + (folder === source ? '' : `assets/${hash}/`) + 'launcher.js"') < index.indexOf('<iframe'), 'deferred module owns the session before iframe navigation');
    const gamePath = path.join(folder, frame[1]);
    assert.ok(fs.existsSync(gamePath), 'iframe target exists');
    assert.equal(path.basename(gamePath), 'game.html');
    const game = fs.readFileSync(gamePath, 'utf8');
    assert.doesNotMatch(game, /<iframe\b/);
    assert.match(game, /<script src="wasm_exec.js"><\/script>\s*<script type="module" src="game.js"><\/script>/);
    for (const htmlPath of [path.join(folder, 'index.html'), gamePath]) {
      const html = fs.readFileSync(htmlPath, 'utf8');
      for (const [, reference] of html.matchAll(/(?:src|href)="([^"]+)"/g)) {
        if (/^https:\/\//.test(reference)) continue;
        assert.ok(fs.existsSync(path.resolve(path.dirname(htmlPath), reference)), `missing ${reference}`);
      }
    }
  }
  for (const file of ['frame.mjs','status.mjs','go-runtime.mjs','embed.css','header.css','footer.css']) assert.ok(fs.existsSync(path.join(first, 'assets', hash, 'shared', file)));
  assert.match(fs.readFileSync(path.join(first, 'assets', hash, 'shared/embed.css'), 'utf8'), /@import '.\/footer.css'/);
  assert.equal(gunzipSync(fs.readFileSync(path.join(first, `assets/${hash}/game.wasm.gz`))).toString(), 'game.wasm');
  assert.ok(fs.existsSync(path.join(first, '404.html')));
  assert.match(fs.readFileSync(path.join(first, 'index.html'), 'utf8'), new RegExp(`href="assets/${hash}/icon-32.png"`));
  for (const icon of icons) assert.equal(fs.readFileSync(path.join(first, `assets/${hash}/${icon}`), 'utf8'), icon);
  assert.equal(fs.readFileSync(path.join(first, 'favicon.ico'), 'utf8'), 'favicon.ico');
  assert.match(fs.readFileSync(path.join(first, `assets/${hash}/game.js`), 'utf8'), /game\.wasm\.gz/);
  assert.doesNotMatch(fs.readFileSync(path.join(first, '_headers'), 'utf8'), /Content-Encoding/);
  assert.match(fs.readFileSync(path.join(first, '_headers'), 'utf8'), /immutable/);
  assert.equal(packagePages(source, path.join(root, 'same')), hash);
  fs.appendFileSync(path.join(source, 'shared/footer.css'), '\n/* shared style update */');
  assert.notEqual(packagePages(source, path.join(root, 'new-shared-style')), hash);
  fs.writeFileSync(path.join(source, 'icon-32.png'), 'new icon');
  assert.notEqual(packagePages(source, path.join(root, 'new-icon')), hash);
  fs.writeFileSync(path.join(source, 'wasm_exec.js'), 'new runtime');
  assert.notEqual(packagePages(source, path.join(root, 'new')), hash);
  assert.throws(() => packagePages(source, first), /EEXIST/);
});
