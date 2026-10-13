# Overall architecture

Appletini is a single Go binary. One goroutine polls GitHub's GraphQL API, and
the tray UI renders whatever the poller sends it over channels. A single
`context.Context`, created in `main` and cancelled on Ctrl+C/SIGTERM or when
the tray app quits, stops everything.

```
main ──dig──▶ config ─▶ gitter.GraphQLClient (Querier)
  │                         ▲
  │                         │ Query(ctx, …)
  ├─ go polling.PollPRs ─▶ queries/aggregator ─▶ by_label / by_repo / by_author / personal (or mock)
  │        │                                         │
  │        │                                         └─ queries/templating (text/template + quote)
  │        ├─ hasher (send only when results change)
  │        │
  │        ├── PRChan ─────────┐
  │        └── HasErrChan ─────┤
  │                            ▼
  └─ ui/pages.IndexPage.Run ─ loop ─▶ ui.Systray (fyne) ─▶ menu + tray icon
```

## Packages

| Package | Responsibility | Depends on |
| --- | --- | --- |
| `main` (root) | Wiring with `dig`, signal-aware context, globals and feature flags | everything below |
| `config` | Load `config.json` from the working directory, creating the default if missing, and migrate it to the latest version | `config/migration`, `logging` |
| `config/migration` | Detect the config version and step it forward version by version (optionally keeping backups) | `migration_types`, `v1`, `v2`, `logging` |
| `config/migration/v1`, `v2` | Schema, defaults, load/save and `ToNext()` for each config version. v2 reads the token from the env var it names | `migration_types` (v1 also depends on v2 to migrate into it) |
| `gitter` | `Querier` and `Mutator` interfaces and `GraphQLClient`: one shared HTTP client, per-request timeout, ctx cancellation, non-2xx responses turned into errors. Also the `PullRequest` model | `status` |
| `gitter/gittertest` | Fake `Querier` for tests | none |
| `trackerkey` | Result-key format shared by the poller and UI (`personal`, `label_0`, …), plus parsing and display order | none |
| `queries` | `Query` interface | `gitter` |
| `queries/templating` | Renders `.gql` templates with `text/template`; `quote` produces safe GraphQL string literals | none |
| `queries/by_label`, `by_repo`, `by_author`, `personal` | Build one GraphQL query per tracker type and map the response to `gitter.PullRequest` | `gitter`, `queries/templating` |
| `queries/mock` | Fixed example data (`MockQueries` feature flag) | `gitter` |
| `queries/aggregator` | Runs the queries in order and merges their results; stops at the first error, keeping partial results | `queries`, `queries/mock`, `gitter` |
| `hasher` | Sends the PR map only when its hash changed; honours ctx while sending | `gitter`, `logging` |
| `polling` | Builds the queries from config (using `trackerkey` ids) and runs the poll loop until ctx is cancelled | `config`, `gitter`, `hasher`, `queries/*`, `trackerkey` |
| `status` | Classifies a PR's review, mergeable and CI (checks) state into emoji, message and icon flags | none |
| `reviews` | Approves a PR via the `addPullRequestReview` mutation, with values sent as GraphQL variables | `gitter` |
| `pullrequests` | Merges and closes a PR via the `mergePullRequest` / `closePullRequest` mutations, with the id sent as a GraphQL variable | `gitter` |
| `actions` | Opens a link in the OS browser (`open` / `xdg-open` / `rundll32`) after validating it | none |
| `ui` | Fyne systray wrapper: menu model, `SetMenu`/`SetIcon` guarded by one mutex, `Notify`, `Confirm` popup, `Quit` | `logging`, fyne |
| `ui/components` | Menu items for a PR and its options | `ui`, `status`, `actions` |
| `ui/icons` | Embedded tray icons | fyne |
| `ui/pages` | Index page: one loop that turns poller updates into the menu and the tray icon; runs Approve, and Merge/Close after a confirmation popup, in the background and reports the result as a notification | `ui`, `ui/components`, `ui/icons`, `config`, `gitter`, `pullrequests`, `reviews`, `trackerkey` |
| `logging` | `Logger` interface and charmbracelet/log implementation | none |

There are no import cycles. `trackerkey`, `status`, `logging`,
`queries/templating` and `gitter/gittertest` are leaves.

## Concurrency

- **Poller goroutine** (`polling.PollPRs`): one round per `poll.frequencySeconds`.
  Every send (`HasErrChan`, `PRChan`) and the wait between rounds also select on
  `ctx.Done()`, so cancellation never leaves it blocked.
- **UI loop** (`ui/pages.IndexPage.loop`): the only goroutine that owns the icon
  state. The error flag comes only from `HasErrChan`; the green/red flags come
  only from the user's own PRs in the latest results.
- **`ui.Systray`** has a single mutex guarding the fyne app, menu and icon, so
  there is no lock ordering to get wrong. `Quit` runs at most once.

## Shutdown

Ctrl+C/SIGTERM cancels the context, which makes `IndexPage.Run` call
`Systray.Quit`. The Quit menu item calls `Systray.Quit` directly. Either way
fyne's `Run` returns, `main` exits, and the deferred cancel stops the poller and
any in-flight GraphQL request.

## Testing

`go test -tags ci ./...` (`task test`). The `ci` tag makes Fyne headless.
Queries are tested against `gittertest.Fake`, and the HTTP client against
`httptest` servers.
