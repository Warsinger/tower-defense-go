# Cloudflare Pages deployment

Agreed target: the same Cloudflare account as math-visualization, free Pages
Direct Upload project `tower-defense-go`, production branch `main`. The actual
assigned `pages.dev` hostname is read from the project API and printed in the
workflow summary. Do not buy a domain, enable a paid plan, Functions, Workers,
R2, or another backend. Report a free-tier limit rather than enabling billing.

## Owner setup

1. In the existing Cloudflare account, create a **Pages Direct Upload** project
   named `tower-defense-go` with production branch `main`. Do not connect a
   second Cloudflare Git build. Confirm the assigned hostname.
2. Install repository Actions secrets `CLOUDFLARE_ACCOUNT_ID` and
   `CLOUDFLARE_API_TOKEN` in this repository. Use an expiring token with only
   **Account → Cloudflare Pages → Edit** on that account, like math-visualization.
   This permission covers Pages projects in the account, not just this repository.
   Never paste token values into chat, issues, source files, or command output.
3. The builder cannot read/copy secrets from another repository or administer
   these secrets. The owner must add or rotate them directly. An actual upload
   validates account/project/token access; dashboard login or a build alone does not.

## Build and publication

```sh
go run ./cmd/web
node --test scripts/package-pages.test.mjs
node scripts/package-pages.mjs
```

The Pages command requires a fresh output directory (`build/pages` by default).
For repeat local packaging, pass a fresh path such as `build/pages-check-2`.
It never recursively deletes an existing directory. CI always starts clean.
Node 24 is used for publishing; the game itself does not require Node.

The browser artifact now contains the deployable `build/pages` bundle. Main
pushes publish only after coverage, all native builds/tests, and browser checks
pass. PRs build/test/package without deployment secrets and never auto-publish.

For explicit publishing, use **Go build → Run workflow**, select the revision's
branch and enable **deploy**. Main publishes production. Another branch also
requires **preview_pr**, an open same-repository PR targeting main whose current
head matches the run. The preview alias is `pr-<number>.<assigned-hostname>`.
Manual dispatch defaults to checks only. Previews are retained until explicitly
removed in Cloudflare; automatic PR-close cleanup is not included in this issue.

The publishing job downloads the already tested bundle, checks its source SHA,
validates the exact project and its Direct Upload/main configuration, and uses
pinned Wrangler 4.151.1. It does not create projects, alter plans, or enable
services. Publishing jobs share a non-cancelling concurrency group. Superseded
main revisions and stale/closed/fork PR previews fail before upload. GitHub may
replace an older pending job with a newer one; rerun explicitly if needed.

On success, **Open the game** appears in the Actions summary. Record the exact
deployment URL, revision, browser/OS, and verification results in the PR/issue.
Until the workflow changes are merged, manual dispatch on the implementation
branch is used for preview validation; production publishing begins only after
an authorized merge and successful main checks.

## Caching, limits, and behavior

All launcher, game JavaScript, Go runtime JavaScript, playfield HTML, and Wasm
bytes contribute to a single SHA-256 asset directory. The root HTML points to
that directory; its launcher resolves the playfield within the same directory.
The game loads its matching runtime and Wasm from there. Root HTML revalidates;
versioned assets are immutable for one year. A root 404 page prevents missing
assets from falling back to the game's HTML. No cross-origin isolation is needed.

The roughly 35 MB Go Wasm binary is precompressed with gzip and served with
`Content-Encoding: gzip` and `Content-Type: application/wasm`. Browser fetch
decodes it before Wasm instantiation. This keeps each uploaded file below the
Pages limit without a paid storage service. Generic local file servers must
supply this encoding header to serve the Pages bundle; the existing `cmd/web`
local preview continues to serve its ordinary uncompressed bundle.

The packager rejects compressed files larger than the Pages 25 MiB asset limit. Upload only
this static bundle. Native score files and credentials are never part of it.
An old loaded tab keeps its old game. If an old asset URL disappears during an
update race, reload from the home page; it cannot silently load a mismatched pair.
Browser scores are origin-scoped: production, preview aliases, deployment URLs,
and localhost have separate histories. Multiplayer remains native-only.

## Hosted validation and rollback

1. Open the deployed URL in a fresh desktop browser session and launch the game.
   Check title/options, tower placement/combat, pause/resume, reset, audio after
   interaction, and score persistence after returning to title and reloading.
2. Check root HTML cache headers, JavaScript and Wasm MIME types, immutable asset
   headers, and an unknown asset returning 404. Test an ordinary cached reload
   after another deployment. A changed runtime or Wasm must change its asset URL.
3. In Cloudflare Pages **Deployments**, choose a known-good **production**
   deployment and select **Rollback to this deployment**. This restores that
   uploaded bundle without rebuilding; previews cannot be production rollback
   targets. Also revert the bad source through a reviewed PR, since the next main
   publish will otherwise replace the rollback. Do not merge/rollback without
   owner authorization. Retained GitHub artifacts are not permanent backup storage.

To rotate credentials, install an equivalent scoped token in the repository
secret, validate a preview, then revoke the old token. Missing/expired credentials
fail publication without deleting the last good site. Never print secret values.

References: [Direct Upload CI](https://developers.cloudflare.com/pages/how-to/use-direct-upload-with-continuous-integration/),
[Pages limits](https://developers.cloudflare.com/pages/platform/limits/),
[headers](https://developers.cloudflare.com/pages/configuration/headers/),
[rollbacks](https://developers.cloudflare.com/pages/configuration/rollbacks/).
