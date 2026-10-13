# Appletini

[![Build Status](https://github.com/the-chaos-collective/appletini/actions/workflows/build-check.yml/badge.svg)](https://github.com/the-chaos-collective/appletini/actions/workflows/build-check.yml)

Appletini is a GitHub client that runs in your systray.

It allows you to keep track of Pull Requests assigned to you and others you may be interested in.

<img width="501" height="534" alt="image" src="https://github.com/user-attachments/assets/81e87a6c-c18f-4262-be09-0cc0c7065aa6" />

## Compatibility

Appletini is multi-platform, supporting Windows, Mac OS, and Linux\*

_\* tested with Gnome on Ubuntu 22.04 (both Wayland & XOrg)_

## Running

### Download a pre-built binary

You may find pre-built binaries for your operating system [here](https://github.com/the-chaos-collective/appletini/releases).

### Configuration

Appletini reads `config.json` from the **current working directory**, that is,
the directory you launch it from, not necessarily the directory the binary is in.
If the file doesn't exist there, a default one is generated in that directory.
When starting Appletini from a desktop shortcut or launcher, make sure its
working directory is set to the folder holding your `config.json`.

Create an env variable on your system holding your personal access token (I used
a classic one with repo, user and projects permission).
You may use any name for it, just make sure to use the same name in the config file.

Each PR also shows its CI status (✅ passing, ❌ failing, ⏳ running), taken from the
checks on its latest commit. A classic token with `repo` scope can read it. If you use
a fine-grained token and the CI status never appears, make sure the token can read the
repository's commit statuses and checks.

On your own PRs, failing CI shows the red badge, and a PR only shows as ready to merge
(green badge) when it's approved, has no conflicts, and its checks pass (or it has none).

Each PR in a tracker has an **Approve** option in its menu, which approves it with the
comment "LGTM!". It isn't shown on your own PRs (GitHub doesn't allow approving those),
on PRs you've already approved yourself (other people's approvals don't hide it), or
on red PRs: merge conflicts, failing CI, or changes requested by someone else. Changes
you requested yourself don't hide it, so you can approve once they're addressed.
Approving needs a token that can write pull request reviews; a classic token with
`repo` scope can. The result (approved, or why not) shows as a desktop notification,
and the menu updates on the next poll.

PRs also have **Merge** (only when the PR is green: approved, no conflicts, checks
passing) and **Close…** options. Merge runs as soon as you click it; Close opens a
confirmation popup first, so nothing is closed by accident. Merge uses the repository's
default merge method. Both need a token with write access to the repository, and report
the result as a desktop notification.

- github:
  - gqlAPI: GitHub GraphQL endpoint (`https://api.github.com/graphql`, or your GitHub Enterprise one)
  - token: NAME of the env variable that holds your github personal access token (not the token itself)
- poll:
  - frequencySeconds: seconds between polls. GitHub rate-limits the API, so keep this reasonable
- itemCount: maximum number of PRs fetched per tracker
- tracking: your trackers (different filters to find PRs you're interested in)
  - personal: true if you'd like to see the "My Pull Requests" sub-menu, false otherwise
  - byLabel:
    - title: string that will appear on the systray to hold all the PRs related to that label
    - label: label you want to filter by
    - repo: repo where you want to filter by label
    - owner: repo owner
  - byRepo:
    - title: string that will appear on the systray to hold your PRs
    - repo: repo you want to see all the PRs of
    - owner: repo owner
  - byAuthor:
    - title: string that will appear on the systray to hold the PRs by these authors
    - authors: list of authors you want to filter by
    - repo: repo you want to see the PRs of
    - owner: repo owner
- darkMode: changes the icon color

### Environment variables

| Variable | Required | Default | Purpose | Used in |
| --- | --- | --- | --- | --- |
| _token variable_, named by `github.token` in `config.json` (`GITHUB_ACCESS_TOKEN` in the default config) | Yes | none | Your GitHub personal access token. Appletini exits at startup if it isn't set. | [config/migration/v2/token.go](config/migration/v2/token.go) |
| `APPLETINI_LOG_LEVEL` | No | `info` | Log level: `debug`, `info`, `warn`, `error` or `fatal`. Unrecognised values fall back to `info`. | [globals.go](globals.go) |

## Development

### Pre-requisites

- [x] go
- [x] gcc

Ubuntu:
`sudo apt install gcc golang`

Mac:
`brew install gcc golang`

Windows:
https://go.dev/doc/install

### Building

We use [Taskfile](https://taskfile.dev). You may run the following command to build & run the software:

```
task run
```

### Testing

```
task test
```

Tests run with Fyne's `ci` build tag, so they don't need a display or the
OpenGL/X11 development headers. `task test:race` runs them with the race
detector, and `task lint` runs golangci-lint.

See [Overall_Architect.md](Overall_Architect.md) for how the packages fit together.
