import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import vm from 'node:vm';
import { gzipSync } from 'node:zlib';

const source = fs.readFileSync(new URL('../web/game.js', import.meta.url), 'utf8');
const wasm = Buffer.from([0, 97, 115, 109, 1, 0, 0, 0]);
for (const compressed of [false, true]) {
  test(`actual loader instantiates ${compressed ? 'gzip' : 'uncompressed'} Wasm`, async () => {
    let ran = false;
    const messages = [];
    const context = {
      window: { addEventListener() {} },
      parent: { postMessage(message) { messages.push(message); } },
      location: { origin: 'https://example.test' },
      console: { error() {} },
      fetch: async () => new Response(compressed ? gzipSync(wasm) : wasm),
      Go: class { importObject = {}; async run() { ran = true; } },
      WebAssembly, Response, Blob, DecompressionStream, Uint8Array,
    };
    await vm.runInNewContext(source, context);
    assert.equal(ran, true);
    // A successful minimal fixture terminates immediately, unlike the game loop.
    assert.equal(messages.length, 1);
    assert.match(messages[0].message, /Game stopped/);
  });
}
