# Coding Standards

This is the initial project style guide. It reflects the current codebase and should evolve as the game grows.

## Go Version And Formatting

- Target the Go version declared in `go.mod`.
- Run `gofmt` on all Go files before review.
- Prefer standard Go formatting and names over custom style rules.
- Keep imports grouped as standard library, local `tower-defense/...` packages, then third-party packages.
- Use aliases only when they clarify naming conflicts or established local convention, such as `comp "tower-defense/components"`.

## Package Responsibilities

- `main`: CLI flags and top-level Ebiten startup only.
- `game`: game initialization, scene switching, stat handoff, and Ebiten `Game` methods.
- `scenes`: high-level title, battle, viewer, options UI, and network controller flow.
- `components`: Donburi component data, entity constructors, entity update behavior, render behavior, stats, and collision helpers.
- `assets`: embedded image/sound/font loading and lookup.
- `config`: shared runtime configuration component.
- `network`: WebSocket transport, necs sync registration, and network message types.
- `strategy`: computer player behavior.
- `util`: small reusable math, cooldown, and filter helpers.

When adding new behavior, put it in the package that owns the concept rather than routing everything through scenes.

## ECS Conventions

- Define Donburi component data as `type XData struct`.
- Define the component variable near the data type, for example `var X = donburi.NewComponentType[XData]()`.
- Entity constructors should be named `NewX` and should fully initialize required components.
- Networked entities should call `srvsync.NetworkSync` immediately after creation when they must appear in remote worlds.
- Component update methods should usually have the form `func (x *XData) Update(entry *donburi.Entry) error`.
- Render methods should usually have the form `Draw(screen *ebiten.Image, entry *donburi.Entry, ...)`.
- Prefer passing `donburi.World` or `*donburi.Entry` explicitly over adding new global state.
- Avoid new package globals for mutable game state. Existing `gameStats` is a known legacy shortcut, not a pattern to expand.

## Gameplay Code

- Keep gameplay rules close to the entities they affect.
- Make costs, cooldowns, caps, and spawn constants named constants when they are part of game balance.
- Prefer returning `error` from update and constructor paths that can fail.
- Use small callback hooks for post-combat effects only when they keep attack logic reusable.
- Keep debug output behind the debug config flag.
- Sound effects should respect the sound config flag.

## Input And Scene Flow

- Scene-level keys belong in scene `Update` methods.
- Entity-specific player actions belong in `PlayerData.UserSpeedUpdate`.
- Do not trigger gameplay actions while a modal is open.
- Keep title, battle, and viewer responsibilities separate.
- A scene transition should be explicit through a callback rather than hidden in component logic.

## Rendering And UI

- Prefer component renderers for entity-specific drawing.
- Keep scene text overlays in scene draw callbacks.
- Debug rendering should reveal state without changing simulation state.
- Use existing asset lookup names for sprites and sounds.
- Add new assets through the embedded asset package and document their names if gameplay code references them.

## Errors And Logging

- Return errors from constructors, scene switches, and update paths when the caller can handle them.
- Use debug prints for optional diagnostics.
- Avoid `log.Fatal` outside startup or unrecoverable background service failures.
- Include useful context in errors that are likely to be shown while debugging placement, networking, or assets.

## Tests

- Add focused table-driven tests for deterministic helpers and rules.
- Keep test names in the current style, for example `Test_makeDisplayName`.
- Prefer testing pure functions or small methods before adding larger integration tests.
- For gameplay changes, add tests around balance formulas, placement validation, stat updates, or collision helpers when practical.
- For networking changes, isolate message parsing/registration behavior from live sockets where possible.

## Documentation

- Keep `README.md` useful for players and quick start information.
- Keep `docs/PROJECT_SPEC.md` aligned with implemented behavior.
- Update this coding standards document when a new pattern becomes preferred.
- Document known gaps as gaps rather than silently implying they are supported.

## Review Checklist

- `gofmt` changed Go files.
- `go test ./...` passes or the failure is documented.
- New gameplay constants are named and easy to find.
- User-facing controls or flags are reflected in the README or project spec.
- New assets are embedded, loaded, and referenced by stable names.
- Stats changes include persistence/display considerations.
