import assert from 'node:assert/strict';
import { readFile, readdir } from 'node:fs/promises';
import { join } from 'node:path';
import { gunzipSync } from 'node:zlib';

const assets = join('build', 'pages', 'assets');
const releases = await readdir(assets, { withFileTypes: true });
const directories = releases.filter(entry => entry.isDirectory());
assert.equal(directories.length, 1, 'expected one versioned Pages asset directory');
const packaged = gunzipSync(await readFile(join(assets, directories[0].name, 'game.wasm.gz')));
const original = await readFile(join('build', 'web', 'game.wasm'));
assert.ok(packaged.equals(original), 'packaged Wasm must match the compiled module');
assert.ok(WebAssembly.validate(packaged), 'packaged Wasm must be valid WebAssembly');
console.log('Packaged WebAssembly is valid and matches the compiled module.');
