# Project workflow

- For issue work, use a dedicated branch and create a GitHub pull request targeting `main`; do not push issue implementation directly to `main`.
- Link the relevant issue and PR in GitHub's Development sidebar, not only through a comment or a plain cross-reference. Use a closing reference such as `Closes #10` when the PR completes the issue scope, and verify the PR's `closingIssuesReferences` includes the issue before handoff.
- Automatically assign every issue and pull request you create to `Warsinger`.
- Request a review from `Warsinger` on every new pull request, in addition to assigning it to that account.
- Apply appropriate issue labels. Use `dependencies` for dependency/toolchain work and `build-and-release` for CI, builds, packaging, and releases. Use `priority: high`, `priority: medium`, or `priority: low` when the user explicitly prioritizes the work; do not infer urgency from task order alone.
- Plan requirements and design with the user before implementing gameplay changes. Record user-agreed design decisions and scope changes in the relevant GitHub issue as they are made. Keep the current issue description consistent with those decisions, preserve useful discussion history, and link the implementation PR.
- Run appropriate checks and report any remaining validation limits in the PR. For CI changes, inspect the PR's workflow results and resolve failures before handing it back.
- Do not merge a PR without the user's authorization.

# GitHub identity

- Use `Warsinger Code Builder` for all GitHub interactions on this repository. On this machine, read `C:\Users\eelwe\.codex\github-identities\README.md` and use the shared helper:

  ```powershell
  node C:\Users\eelwe\.codex\github-identities\github-role.mjs run builder --repo Warsinger/tower-defense-go -- gh <arguments>
  ```

- Pass the matching `--repo Warsinger/tower-defense-go` to gh commands that accept a repository option. The helper's repository option selects the installation token scope.
- Do not use connected GitHub tools or the default gh login for this repository. If the helper is unavailable, report the blocker before publishing under another identity.
- Never expose credentials, put tokens in remote URLs or persistent Git configuration, or commit local identity files. Git push authentication and commit authorship must also use the intended builder identity; the helper's API identity does not automatically configure either.

# Game development

- Read `docs/PROJECT_SPEC.md`, `docs/CODING_STANDARDS.md`, and `docs/TESTING.md` before making relevant changes.
- Target the Go version declared in `go.mod`; CI reads that file for its toolchain version.
- Keep gameplay behavior, ECS conventions, formatting, and tests consistent with the project documentation.
- Keep the README and project documentation aligned with changes to controls, gameplay, configuration, or testing.
- For Go changes, run `gofmt` on changed Go files, `go test ./...`, `go vet ./...`, and `go build ./...` as applicable. Report platform or graphics/display limitations instead of claiming unrun checks passed.
