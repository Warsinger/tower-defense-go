# Tower Defense Project Spec

This document captures the current behavior of the Go tower defense game so future changes can be made against an explicit baseline. It describes what exists today, not every idea in the TODO list.

## Product Goals

Tower Defense is a desktop arcade tower defense game built with Ebiten. The core loop is to protect a base by placing, healing, and upgrading towers while increasingly difficult creep waves move down the board.

The project currently supports:

- Single-player survival play.
- Optional computer-controlled play for simulation-style runs.
- Experimental network play with a local game and remote viewer.
- Persistent high score and aggregate play statistics.

## Runtime

- Language: Go 1.23.1.
- Renderer/input/audio: Ebiten.
- UI widgets: EbitenUI.
- ECS: Donburi.
- Network ECS sync: leap-fish/necs over WebSocket.

The executable starts from `main.go`, parses CLI flags, loads assets, initializes game state, and runs the Ebiten game loop.

## Launch Options

Supported flags:

| Flag | Default | Meaning |
| --- | ---: | --- |
| `-width` | `600` | Board width in pixels. |
| `-height` | `800` | Board height in pixels. |
| `-speed` | `60` | Game update speed in ticks per second. Values are clamped to `0..60` in battle. |
| `-debug` | `false` | Start with debug rendering enabled. |
| `-level` | `0` | Starting tower-level progress. Higher values increase max tower level and creep level. |
| `-computer` | `false` | Enable the computer player strategy instead of direct player placement. |
| `-complevel` | `3` | Computer player action speed from `1` slowest to `5` fastest. |
| `-nosound` | `false` | Start with sound effects disabled. |

## Scene Flow

The game has a simple scene stack:

- Title scene: shows high scores, instructions, title art, and EbitenUI controls.
- Battle scene: runs the active game board.
- Viewer scene: renders a synced remote world next to the local board during multiplayer.

Title scene actions:

- Start Game begins battle mode.
- Game Options opens a modal for multiplayer setup plus debug and grid-line options.
- Space starts the game when no modal is open.

Battle scene end/reset:

- The base dying sets game over, finalizes game time, and leaves the final board visible.
- Pressing `R` returns to the title scene and merges current run stats into persistent stats.
- Pressing `Q` from any non-modal state saves stats and exits.

## Board And Entities

The board is a 2D pixel space with a configurable width and height. Entity collision and placement use image rectangles.

Core entity/component groups:

- Board: board width and height.
- Player/base: position, health, attack, sprite render, info render, score, money, and tower-level progress.
- Tower: position, health, attack, level, sprite render, range render, and info render.
- Creep: position, velocity, health, attack, sprite render, range render, and info render.
- Bullet: position, velocity, attack, bullet render, and launch path metadata.
- Battle state: paused and game-over flags.
- Config: debug, grid lines, stats display, sound, computer, server port, and client address.

## Player/Base Rules

- The base starts near the bottom of the board at `board.Height - 70`.
- The base starts with 100 health.
- The player starts each run with `$500`.
- The base has a short-range attack against creeps.
- The base does not get removed when health reaches zero; it is marked dead and battle enters game-over state.
- Score increases when towers kill creeps.
- Money increases from creep kills and from periodic creep-wave income.

Tower-level progress controls difficulty:

- `TowerLevels` starts from the `-level` flag.
- Creep level is `floor(TowerLevels / 5) + 1`.
- Max tower level is `5 + floor(TowerLevels / 20)`.
- Upgrading a tower increments `TowerLevels`.

## Tower Rules

- Towers cost `$50` to place.
- Towers are centered on the mouse click.
- A tower cannot be placed out of board bounds.
- A tower cannot overlap the base or existing blocking entities.
- New towers start with 20 health, level 1, ranged single-target attack, power 1, range 50, and attack cooldown 30 ticks.
- Tower health also acts as ammo. Each tower shot decrements tower health by 1.
- A tower is removed when its health/ammo reaches zero.
- Healing a tower costs `$25` and restores it to full health.
- Upgrading a tower costs `$50`, heals it to full, and increases level, max health, power, range, and attack speed.
- Each upgrade adds 5 max health, adds 3 range, reduces cooldown by 3 to a minimum of 3, and adds `level / 3` power.
- Towers cannot be upgraded beyond the current max tower level.

## Creep Rules

- Creeps spawn near the top of the board at `SpawnBorder` (`60`) and move downward.
- A normal creep has a 30% chance to become the larger variant.
- Creep speed, health, attack power, range, cooldown, and score value scale with creep level and variant.
- Normal creep score value is 10 or 20 depending on variant.
- Super creeps are multiplayer-only creeps with score value 50, health 20, power 8, diagonal velocity, and the `supercreep` sprite.
- Creeps try to move toward their target position each game tick.
- If blocked by another creep, they attempt small sideways movement.
- If blocked by a tower or base, they try to creep forward slightly and attack when in range.
- Creeps attack towers and the base.

## Combat And Bullets

- Attacks use rectangular range checks expanded from the attacker's render bounds.
- Ranged attacks launch bullets toward a target midpoint.
- Bullet speed is currently 8.
- Tower bullets are green and target creeps.
- Creep bullets are red and target towers or the base.
- Bullets lead moving targets unless `noLead` is set, as it is for the base.
- A bullet is removed when it hits an enemy, leaves the board, or travels more than 150% of its original planned path.
- On creep kill, the creep is removed and the player gains money and score equal to the creep score value.

## Wave And Difficulty Rules

- Battle starts with a partially advanced creep timer.
- The creep timer advances by at least 4 each entity-update tick and scales upward with creep level.
- A new wave can spawn when the timer reaches `180 - creepLevel`.
- At most 25 creeps can be on the board before spawning pauses.
- If the board is at the creep cap, the player still receives `$5` income.
- When a wave spawns, the wave count increments, spawn count is selected probabilistically, and higher creep and wave levels bias toward larger waves.
- Every 10 waves adds one extra creep level for spawn calculations.
- The player receives `$5` per spawned creep.

## Controls

Global:

- `Q`: save stats and quit when no modal is open.
- `F`: toggle full screen.

Title:

- Space: start game when no modal is open.
- Start Game button: start game.
- Game Options button: open multiplayer/debug/grid options.

Battle:

- Left mouse click: place a tower.
- Mouse over tower + `H`: heal tower.
- Mouse over tower + `U`: upgrade tower.
- `P` or Space: pause/unpause.
- `R`: return to title and save current run stats.
- `+`: increase game speed by 5, max 60.
- `-`: decrease game speed by 5, min 0.
- `L`: toggle grid lines.
- `D`: toggle debug rendering.
- `S`: toggle sound.
- `T`: toggle stats display.
- Multiplayer only, `C`: send a super creep to the peer when a peer is connected and the cooldown is ready.

Viewer:

- `L`: toggle viewer grid lines.
- `D`: toggle viewer debug rendering.

## Rendering

Rendering is component-based:

- `DrawBoard` clears/draws the board, entities, grid, and scene text.
- Sprite components map entity names to loaded image assets.
- Info render displays entity health/cooldown/level details.
- Range render displays debug range indicators.
- Bullet render draws colored circles and debug trajectory lines.

Debug mode should expose additional targeting and timing information without changing gameplay behavior.

## Stats Persistence

Stats are stored in `score/stats.txt` as `Name=Value` lines.

Persistent high values:

- `HighScore`
- `HighCreepLevel`
- `HighTowerLevel`

Aggregate tracked stats include bullets expired, bullets fired, creeps killed/spawned, creep waves, games played, money spent, player deaths, tower events, and game time.

Stats are loaded at startup and saved on quit or when returning from battle to title.

## Networking And Multiplayer

Network play is experimental and built on necs ECS synchronization.

Modes:

- Server mode starts a WebSocket server for the current world.
- Client mode connects to a server and creates a separate synced client world.
- When two worlds are available, battle mode renders the local board and the remote viewer side by side.

Network messages:

- `StartGameMessage`: tells the peer to enter battle mode.
- `ClientConnectMessage`: lets a server connect back to a client-provided address.
- `CreepMessage`: requests a super creep spawn in the peer world.

Current constraints:

- Networking assumes a direct WebSocket connection between peers.
- Client/server lifecycle cleanup is incomplete.
- Server sync failures currently terminate the process.
- Super creep sending has a cooldown and currently checks for `$50`, but the money is not deducted.

## Known Gaps

- Tower type selection is represented in costs but only the ranged tower path is currently active.
- Difficulty options are mostly CLI-driven and not yet fully exposed in UI.
- Computer strategy exists but needs a clearer behavioral spec and tests.
- Network setup and teardown need hardening.
- Some collision edge cases are acknowledged in code comments.
- Stats storage is plain text and local to the process working directory.
- Test coverage is currently limited to stats display-name formatting.
