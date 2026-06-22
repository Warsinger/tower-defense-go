# Testing Guide

This project uses Go's standard `testing` package. The current test suite focuses on deterministic logic that can run without opening an Ebiten window.

## Running Tests

Run all tests from the repository root:

```sh
go test ./...
```

For a quick coverage snapshot:

```sh
go test ./... -cover
```

## Current Coverage Focus

The suite now covers these non-UI behaviors:

- Stats display-name formatting.
- Stats initialization, high-score preservation, aggregation, reset, and output formatting.
- Player difficulty formulas for creep level and max tower level.
- Tower healing, upgrade scaling, max-level blocking, ammo consumption, and ammo-out removal.
- Cooldown timer lifecycle and display behavior.

## Preferred Test Shape

- Keep tests near the package that owns the behavior.
- Prefer table-driven tests for formulas and formatting.
- Build small Donburi worlds directly in tests when ECS state is needed.
- Use component constructors only when the behavior under test needs their side effects.
- Avoid starting the full Ebiten game loop from unit tests.

## Known Test Seams

Some logic is testable but still needs cleaner seams:

- Tower placement currently depends on loaded image assets for sprite bounds.
- Combat targeting and bullet creation depend on render bounds and world/config setup.
- Battle wave spawning uses randomness and scene state, so it needs either injectable randomness or narrower helper functions.
- UI behavior is best verified manually or with screenshot-driven checks until the UI construction is split into smaller testable pieces.
- Networking should be tested around message handling and sync registration before attempting live socket integration tests.

## Next Candidates

Useful next tests would cover:

- Tower placement success and placement errors after extracting or injecting sprite bounds.
- Creep movement collision decisions with small synthetic entities.
- Bullet expiry and hit behavior with deterministic world setup.
- Battle wave spawn count selection after separating random choice from entity creation.
- Computer strategy decisions using small board fixtures.
