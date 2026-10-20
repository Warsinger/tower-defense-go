# Application icons

The icon uses Matt's original high-resolution photograph of Nathan's LEGO tower,
scaled onto the approved pale gray background (`#dbe5eb`) with 6% padding.
The model is not redrawn. Gameplay sprites remain unchanged.

`assets/icons/source.png` is the source photograph. Regenerate PNG, ICO and ICNS
files from the repository root with `go run ./cmd/icons`. The generator uses
Go's PNG encoder and `golang.org/x/image/draw` (already used by the project).
Generated files are checked in so ordinary builds need no image tooling.

Windows amd64 builds automatically embed `resource_windows_amd64.syso`. After
regenerating icons, regenerate that resource with the pinned build tool:

```sh
go run github.com/akavel/rsrc@v0.10.2 -arch amd64 -ico assets/icons/favicon.ico -o resource_windows_amd64.syso
```

The native window manager receives 32, 48 and 256 pixel icons through
`ebiten.SetWindowIcon`. macOS uses its application bundle icon instead.

## Build artifacts

- Windows: the `.exe` includes the icon for Explorer and application windows.
- macOS: extract the `.tar.gz` and launch `Tower Defense.app` for the Finder/Dock
  icon. The standalone executable is also retained. Bundles are unsigned and
  not notarized; this change does not add a signing/distribution service.
  The bundle launcher uses `~/Library/Application Support/Tower Defense` as its
  working directory so Finder launches can save `score/stats.txt`. Existing CLI
  scores are not moved; the standalone executable keeps its existing behavior.
- Linux: the archive contains the executable, desktop entry, PNG and installation
  instructions. Install those into the documented user directories to show the
  launcher icon. Window/taskbar behavior depends on the desktop environment.

For a local macOS/Linux build, build the executable first, then run:

```sh
node scripts/package-desktop.mjs darwin-arm64 build/tower-defense-darwin-arm64
# Or darwin-amd64 / linux-amd64 with the matching native binary.
```

Web builds include ICO, 16/32 pixel favicons, a 180 pixel Apple touch icon and
192/512 pixel PNGs. Both the launcher and standalone playfield declare icons.
Pages packaging gives icons versioned URLs alongside the other assets and
retains `/favicon.ico` for clients that request the conventional path. This does
not introduce a service worker or installable PWA.

CI checks icon file structure, icon URL packaging, native builds, macOS metadata
and ICNS decoding, and the Linux desktop entry. Manually check browser tabs,
Windows Explorer/taskbar, macOS Finder/Dock, and the Linux launcher when reviewing
platform appearance; operating systems and browsers can cache previous icons.
