import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import vm from 'node:vm';

test('automatic playfield reports ready and errors without needing a launch button', () => {
  let listener;
  let focused = 0;
  const status = { dataset: {}, textContent: 'Loading game…' };
  const frame = { contentWindow: {}, focus: () => focused++ };
  const nodes = { game: frame, status };
  const origin = 'https://game.example';
  vm.runInNewContext(fs.readFileSync(new URL('../web/launcher.js', import.meta.url), 'utf8'), {
    window: { addEventListener: (type, callback) => { assert.equal(type, 'message'); listener = callback; } },
    document: { getElementById: id => nodes[id] }, location: { origin }
  });
  const message = (type, source = frame.contentWindow, from = origin) => listener({ origin: from, source, data: { type, message: '<download failed>' } });
  message('game-ready', {});
  message('game-ready', frame.contentWindow, 'https://other.example');
  assert.equal(status.textContent, 'Loading game…');
  message('game-ready');
  assert.match(status.textContent, /Game ready/);
  assert.equal(focused, 1);
  message('game-error');
  assert.equal(status.dataset.error, 'true');
  assert.match(status.textContent, /Reload to retry.*<download failed>/);
  message('game-ready');
  assert.match(status.textContent, /could not run/);
});
