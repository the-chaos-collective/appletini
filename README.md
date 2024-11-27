# Appletini
[![Build Status](https://github.com/the-chaos-collective/appletini/actions/workflows/build-check.yml/badge.svg)](https://github.com/the-chaos-collective/appletini/actions/workflows/build-check.yml)

Appletini is a GitHub client that runs in your systray.

It allows you to keep track of Pull Requests assigned to you and others you may be interested in.

## Compatibility

Appletini is multi-platform, supporting Windows, Mac OS, and Linux\*

_\* tested with Gnome on Ubuntu 22.04 (both Wayland & XOrg)_

## Running

### Download a pre-built binary

You may find pre-built binaries for your operating system [here](https://github.com/the-chaos-collective/appletini/releases).

### Configuration

You will need a `config.json` file present next to the binary. You may copy [the example file](config.example.json).

Create a env Variable on your system holding your personal access token (I used
a classic one with repo, user and projects permission).
You may use any name for it, just make sure to save the name of this env var, as you will need it later in the config.

- token: has the NAME of the env variable that holds your github personal access token
- trackers: hold your trackers (different filters to find PRs you're interested in)
  - byLabel:
    - title: string that will appear on the systray to hold all the PRs related to that label
    - label: label you want to filter by
    - repo: repo where you want to filter by label
    - owner: repo owner
  - byRepo:
    - title: string that will appear on the systray to hold your PRs
    - repo: repo you want to see all the PRs of
    - owner: repo owner
- darkMode: changes the icon color

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

_Don't forget to copy the `config.example.json` file and rename it to `config.json` in order to configure your GitHub access token, among other things._

We use [Taskfile](https://taskfile.dev). You may run the following command to build & run the software:

```
task run
```
