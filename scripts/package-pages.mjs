import fs from 'node:fs';
import path from 'node:path';
import crypto from 'node:crypto';
import { gzipSync } from 'node:zlib';
import { pathToFileURL } from 'node:url';

export function packagePages(source, destination, revision = 'local') {
  const names = ['game.html', 'launcher.js', 'game.js', 'wasm_exec.js', 'game.wasm'];
  const files = new Map(names.map(name => [name, fs.readFileSync(path.join(source, name))]));
  files.set('game.wasm', gzipSync(files.get('game.wasm'), { level: 9 }));
  const digest = crypto.createHash('sha256');
  for (const [name, content] of files) {
    if (content.length > 25 * 1024 * 1024) throw new Error(`${name} exceeds the Pages 25 MiB asset limit`);
    digest.update(name).update('\0').update(content);
  }
  const hash = digest.digest('hex');
  // Destination must be a fresh directory; never recursively delete caller data.
  fs.mkdirSync(destination, { recursive: false });
  const assets = path.join(destination, 'assets', hash);
  fs.mkdirSync(assets, { recursive: true });
  for (const [name, content] of files) fs.writeFileSync(path.join(assets, name), content);
  const index = fs.readFileSync(path.join(source, 'index.html'), 'utf8');
  if (!index.includes('src="launcher.js"')) throw new Error('Launcher reference missing');
  fs.writeFileSync(path.join(destination, 'index.html'), index.replace('src="launcher.js"', `src="assets/${hash}/launcher.js"`));
  fs.writeFileSync(path.join(destination, '404.html'), '<!doctype html><title>Not found</title><h1>Not found</h1><p>Reload the game from its home page.</p>');
  fs.writeFileSync(path.join(destination, '_headers'), `/*\n  X-Content-Type-Options: nosniff\n/\n  Cache-Control: no-cache\n/index.html\n  Cache-Control: no-cache\n/version.json\n  Cache-Control: no-cache\n/assets/*\n  Cache-Control: public, max-age=31536000, immutable\n/assets/${hash}/game.wasm\n  Content-Encoding: gzip\n  Content-Type: application/wasm\n`);
  fs.writeFileSync(path.join(destination, 'version.json'), JSON.stringify({ revision, assets: `assets/${hash}` }) + '\n');
  return hash;
}

if (process.argv[1] && import.meta.url === pathToFileURL(path.resolve(process.argv[1])).href) {
  const output = process.argv[2] || 'build/pages';
  console.log(`Pages bundle ${packagePages('build/web', output, process.env.GITHUB_SHA)}: ${output}`);
}
