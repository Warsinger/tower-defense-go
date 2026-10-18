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

For a local coverage report with the same files produced by CI:

```sh
go test -coverprofile coverage.out -covermode atomic ./...
go tool cover -func coverage.out > coverage-summary.txt
go tool cover -html coverage.out -o coverage.html
```

Open `coverage.html` in a browser to inspect uncovered lines. `coverage-summary.txt` contains the package/function totals in plain text.

## Coverage Tracking

Coverage is tracked with Go's built-in coverage tools and GitHub Actions only. The `Coverage` workflow job:

- Runs `go test -coverprofile coverage.out -covermode atomic ./...`.
- Writes the `go tool cover -func` summary into the GitHub Actions job summary.
- Uploads `coverage.out`, `coverage-summary.txt`, and `coverage.html` as the `go-coverage` artifact.

This keeps coverage history available through GitHub workflow runs without depending on a third-party coverage service.

## Current Coverage Focus

The suite now covers these non-UI behaviors:

- Balance config loading from embedded defaults and external JSON files.
- Balance-driven player, tower, normal creep, and super creep initialization.
- Balance-driven tower heal and upgrade costs.
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
