import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { execFileSync } from 'node:child_process';
import { packageDesktop } from './package-desktop.mjs';

test('desktop archives include the actual binary and platform icon metadata', t => {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), 'tower-icons-'));
  t.after(() => {
    assert.equal(path.dirname(path.resolve(root)), path.resolve(os.tmpdir()));
    assert.ok(path.basename(root).startsWith('tower-icons-'));
    fs.rmSync(root, { recursive: true, force: true });
  });
  const binary = path.join(root, 'binary');
  fs.writeFileSync(binary, 'fixture executable');
  for (const target of ['linux-amd64', 'darwin-arm64']) {
    const archive = packageDesktop(target, binary, root);
    const listing = execFileSync('tar', ['-tzf', archive], { encoding: 'utf8' });
    const folder = path.join(root, `package-${target}`);
    if (target.startsWith('darwin-')) {
      assert.match(listing, /Tower Defense.app\/Contents\/Resources\/tower-defense.icns/);
      const contents = path.join(folder, 'Tower Defense.app', 'Contents');
      assert.equal(fs.readFileSync(path.join(contents, 'MacOS', 'tower-defense'), 'utf8'), 'fixture executable');
      assert.match(fs.readFileSync(path.join(contents, 'Info.plist'), 'utf8'), /CFBundleIconFile<\/key><string>tower-defense.icns/);
      assert.match(fs.readFileSync(path.join(contents, 'Info.plist'), 'utf8'), /CFBundleExecutable<\/key><string>launch/);
      assert.match(fs.readFileSync(path.join(contents, 'MacOS', 'launch'), 'utf8'), /Library\/Application Support\/Tower Defense/);
    } else {
      assert.match(listing, /tower-defense.desktop/);
      assert.equal(fs.readFileSync(path.join(folder, 'tower-defense'), 'utf8'), 'fixture executable');
      assert.match(fs.readFileSync(path.join(folder, 'tower-defense.desktop'), 'utf8'), /Icon=tower-defense/);
    }
    assert.throws(() => packageDesktop(target, binary, root), /EEXIST/);
  }
  assert.throws(() => packageDesktop('unknown', binary, root), /Unsupported/);
});
