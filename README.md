# Fjira - the fuzziest Jira command line tool in the world.

<img src="fjira.png" alt="drawing" width="256"/>

![Test](https://github.com/mk-5/fjira/actions/workflows/tests.yml/badge.svg)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://github.com/mk-5/fjira/blob/master/LICENSE)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/mk-5/fjira)
[![Go Report Card](https://goreportcard.com/badge/github.com/mk-5/fjira)](https://goreportcard.com/report/github.com/mk-5/fjira)
[![Go Reference](https://pkg.go.dev/badge/github.com/mk-5/fjira.svg)](https://pkg.go.dev/github.com/mk-5/fjira)

## Demo

![Fjira Demo](demo.gif)

## Features

- Search for Jira Projects, and Issues
- Changes Jira Issue assignee
- Changes Jira Issue status
- Supports multiple workspaces
- Runs on MacOS and Linux

## Install

### Mac OS

```shell
brew tap mk-5/mk-5
brew install fjira
```

### Linux

Go to [https://github.com/mk-5/fjira/releases/latest](https://github.com/mk-5/fjira/releases/latest), and check the
latest release version.

#### Deb

```shell
sudo dpkg -i fjira_0.4.0_linux_amd64.deb
```

#### Binary

```shell
tar -xvzf fjira_0.4.0_Linux_x86_64.tar.gz
cp fjira /usr/local/bin/fjira
```

### Build from sources

```shell
make
./out/bin/fjira
```

## Usage

```text
Usage:
    fjira
    fjira [command]
    fjira [command] [flags]
    fjira [flags]
    fjira [jira-issue] [flags]

Available Commands:
    workspace               Switch fjira workspace
    help               	    Help
    version                 Show version

Flags:
    -p, --project             Open project issues search directly from CLI, example: -p GEN.
    -w, --workspace           Use different fjira workspace without switching it globally, example: -w myworkspace
    -nw, --new-workspace      Create new workspace, example: fjira --new-workspace=abc
```

## Getting Started

The CLI Jira interface is pretty straightforward. Just run `fjira` in your terminal.

```shell
fjira
```

## Workspaces

Fjira will ask you about Jira API url, and token if you run fjira for the very first time.

![Fjira First Run](demo_first_run.gif)

Fjira workspace is a set of jira configuration data, and it's kept in simple json file under `~/.fjira` directory.
You can switch within multiple workspaces using `fjira workspace` command.

```shell
fjira workspace
```

It will open a fuzzy finder with all available workspaces.
In order to create a new workspace you need to use following command:

```shell
fjira --new-workspace abc
```

You can edit existing workspace using `--edit-workspace` flag.

```shell
fjira --edit-workspace abc
```

## Projects search

The projects search is a default view, just run `fjra` in order to open it.

```shell
fjira
```

## Open project directly from cli

You can open a project directly from cli.

```shell
fjira --project=PROJ
```

The fjira app will skip projects search screen, and it'll go to the next screen with issues search.

## Open issue directly from cli

```shell
fjira PROJ-123
```

The app will skip all the screens, and it'll go to the issue view directly.

## Board View

![Fjira Board View](demo_board_view.png)

You can open board-like view using the navigation buttons from the project menu.
First open the project, and then press F4.

## The Future (TODO)

- More docs
- Windows support
- Support Linux packages managers nonsense aka. Snapcraft, Deb, AUR
- More Jira features ;)

#### Motivation

I've created this tool for myself, and the only motivation behind it was laziness (maybe plus fact that I like terminal
tools).
It's really common that from time to time you need to do something like: "I'd just need to move issue 123 to the next
status".
Opening Jira, finding ticket inside the board, opening the Jira issue modal... is okay, but it takes some time.
I found it really handy to do it from terminal, where probably I do something anyway 😝

Sooo, if you feel it the same way as me - I'd love to get some star from you 🤜🤛
