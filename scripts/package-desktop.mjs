import fs from 'node:fs';
import path from 'node:path';
import { execFileSync } from 'node:child_process';
import { pathToFileURL } from 'node:url';

// Keep the standalone binary and add an archive that preserves Unix permissions.
export function packageDesktop(target, binary, output = 'build', icons = 'assets/icons') {
  if (!['linux-amd64', 'darwin-amd64', 'darwin-arm64'].includes(target)) throw new Error(`Unsupported desktop package: ${target}`);
  const root = path.join(output, `package-${target}`);
  fs.mkdirSync(root); // Refuse to overwrite an existing package.
  if (target.startsWith('darwin-')) {
    const contents = path.join(root, 'Tower Defense.app', 'Contents');
    fs.mkdirSync(path.join(contents, 'MacOS'), { recursive: true });
    fs.mkdirSync(path.join(contents, 'Resources'));
    fs.copyFileSync(binary, path.join(contents, 'MacOS', 'tower-defense'));
    fs.chmodSync(path.join(contents, 'MacOS', 'tower-defense'), 0o755);
    // Finder does not provide a writable working directory. Keep the existing
    // relative score storage in a stable user directory without changing CLI runs.
    const launcher = path.join(contents, 'MacOS', 'launch');
    fs.writeFileSync(launcher, `#!/bin/sh
set -eu
bundle_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
mkdir -p "$HOME/Library/Application Support/Tower Defense"
cd "$HOME/Library/Application Support/Tower Defense"
exec "$bundle_dir/tower-defense" "$@"
`);
    fs.chmodSync(launcher, 0o755);
    fs.copyFileSync(path.join(icons, 'tower-defense.icns'), path.join(contents, 'Resources', 'tower-defense.icns'));
    fs.writeFileSync(path.join(contents, 'Info.plist'), `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>CFBundleName</key><string>Tower Defense</string>
<key>CFBundleDisplayName</key><string>Tower Defense</string>
<key>CFBundleIdentifier</key><string>io.github.warsinger.tower-defense</string>
<key>CFBundleExecutable</key><string>launch</string>
<key>CFBundleIconFile</key><string>tower-defense.icns</string>
<key>CFBundlePackageType</key><string>APPL</string>
<key>NSHighResolutionCapable</key><true/>
</dict></plist>
`);
  } else {
    fs.copyFileSync(binary, path.join(root, 'tower-defense'));
    fs.chmodSync(path.join(root, 'tower-defense'), 0o755);
    fs.copyFileSync(path.join(icons, 'icon-256.png'), path.join(root, 'tower-defense.png'));
    fs.writeFileSync(path.join(root, 'tower-defense.desktop'), `[Desktop Entry]
Type=Application
Name=Tower Defense
Comment=Protect your base with LEGO towers
Exec=tower-defense
Icon=tower-defense
Terminal=false
Categories=Game;StrategyGame;
`);
    fs.writeFileSync(path.join(root, 'INSTALL.txt'), `Install for your user from this directory:
mkdir -p "$HOME/.local/bin" "$HOME/.local/share/applications" "$HOME/.local/share/icons/hicolor/256x256/apps"
install -m 755 tower-defense "$HOME/.local/bin/tower-defense"
install -m 644 tower-defense.desktop "$HOME/.local/share/applications/tower-defense.desktop"
install -m 644 tower-defense.png "$HOME/.local/share/icons/hicolor/256x256/apps/tower-defense.png"

Ensure $HOME/.local/bin is on your desktop session's PATH, then sign out/in if needed.
The desktop environment may cache icons; refresh its icon cache if needed.
`);
  }
  const archive = path.join(output, `tower-defense-${target}.tar.gz`);
  execFileSync('tar', ['-czf', path.resolve(archive), '-C', root, '.']);
  return archive;
}

if (process.argv[1] && import.meta.url === pathToFileURL(path.resolve(process.argv[1])).href) {
  console.log(packageDesktop(process.argv[2], process.argv[3]));
}
