# Browser preview

The desktop browser version reuses the native Go/Ebitengine game and embedded
assets. This version supports single-player gameplay only. See [Deployment](DEPLOYMENT.md)
for Cloudflare hosting.

## Build and run locally

Install the Go version in `go.mod` (currently Go 1.27.2). From the repository root:

```sh
npm ci --ignore-scripts
go run ./cmd/web -serve
```

Open <http://localhost:8080>. The game loads directly on the page with instructions
below the playfield and no separate launch step. Press Space or click Start Game
on the game's title screen. Click inside the game if it needs focus. Stop the local server with Ctrl+C.
Use `-port 8081` to choose another port. The server binds only to this computer
and disables caching for local development. Do not open the HTML through `file://`.

To build without serving:

```sh
npm ci --ignore-scripts
go run ./cmd/web
```

The static bundle is written to `build/web`. The command compiles for `js/wasm`
and copies `wasm_exec.js` from the same Go installation used for compilation.
Node 24+ is required to install the pinned shared browser package. No backend or hosting account is required for local builds. Native builds remain independent of npm.
The command must be run using native Go, without setting GOOS/GOARCH to a
cross-compilation target in the parent shell.

## Shared browser components

The exact Warsinger/web-shared Git commit and integrity hash are locked in package-lock.json. cmd/web copies the shared frame lifecycle, Go loader, simple embed layout and footer stylesheet into build/web/shared/. Packaging hashes all these files alongside the game and icons; deployed pages make no cross-site script/style requests. Update the dependency through a reviewed consumer PR with npm install, committed lockfile, tests and a preview. Do not use a mutable branch dependency or edit node_modules.

The outer page installs its session-validated listener before navigating the iframe, reports real Go readiness/errors and focuses the canvas. The existing towerDefenseReady callback is configured in the shared Go loader. The Go runtime is still copied from the exact compiling toolchain. The game's title/start screen, gameplay, settings, sound and browser-local score storage are unchanged. Common runtime/lifecycle tests live in web-shared; Tower retains bundle/icon/desktop packaging tests and the actual Go/Wasm storage harness.

## Controls and browser differences

- The logical board stays 600×800 and scales with the containing frame. Resizing
  changes its display size, not the simulation or placement coordinates.
- Left-click places a tower. Hover a tower and press H to heal or U to upgrade.
- P/Space pauses, R returns to the title and saves the completed run, S toggles
  sound, D toggles debug information, and L toggles the grid.
- Click the playfield to focus input. Space and arrow keys do not scroll while
  the game has focus. Browser shortcuts remain browser shortcuts.
- Audio requires a user interaction inside the game; browser/site mute settings
  still apply. Automatic loading does not unlock audio by itself.
- Game Options retains debug/grid settings. Native server/client controls are
  unavailable; browser hosting of a WebSocket server and native peer discovery
  are outside this milestone. The native networking/viewer paths are unchanged.
- Q does not terminate the browser runtime. Return to title with R before
  closing/reloading the tab. Refreshing during a run loses that unfinished run.
- Browser startup uses default CLI values and embedded balance configuration;
  local filesystem balance overrides are a native feature.

## Scores

Completed-run totals and high scores save to localStorage under
`tower-defense.stats.v1`, separate from native `score/stats.txt`. Storage is
scoped to the site origin, including port: changing host, port, or deployment URL
starts a separate score history. Clearing site data removes browser scores.

If storage reads/writes are blocked or fail, the game remains usable with scores
in memory for the current page session and logs a warning in the browser console.
No native score file is read, migrated, overwritten, or deleted by the browser.

## Validation

Run `go test ./...`, `go vet ./...`, `go build ./...`, and `go run ./cmd/web`.
Run `node scripts/test-browser-storage.mjs` (Node 22 or later) to exercise the actual
Go/Wasm storage functions with normal storage, a throwing localStorage property
getter, throwing getItem/setItem methods, and missing storage. Each failure case
must preserve an in-memory score round trip without terminating the runtime.
Node also installs the shared browser assets before building or serving; the runtime does not fetch npm packages.

The CI browser job runs this harness, builds and uploads the static bundle; it does not simulate
graphical browser interaction. Test locally using the server procedure above.

Manual checks: launch/loading/title rendering; Game Options save/cancel; starting
a game; placement/combat; healing/upgrading; pause/resume; R/reset; score reload;
keyboard focus and scaling; audio after interaction. Inspect developer console
for errors. A successful Wasm build alone does not establish browser support.
The implementation PR records the tested browser/OS and remaining limitations.
Mobile/touch support, broad browser coverage, and automatic browser graphics tests
remain follow-up work.
