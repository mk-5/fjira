# Fjira - the fuzziest Jira command line tool in the world.

<img src="fjira.png" alt="drawing" width="256"/>

![Test](https://github.com/mk-5/fjira/actions/workflows/tests.yml/badge.svg)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://github.com/mk-5/fjira/blob/master/LICENSE)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/mk-5/fjira)

## Demo

![Fjira Demo](demo.gif)

## Features

- Search for Jira Projects, and Issues
- Changes Jira Issue assignee
- Changes Jira Issue status
- Supports multiple workspaces
- Runs on MacOS and Linux

## Getting started

### Mac OS

```shell
brew tap mk-5/mk5
brew install fjira
```

### Build from sources

```shell
git clone git@github.com:mk-5/fjira.git
cd fjira
./scripts/build.sh
cp fjira /usr/local/bin/fjira # or ln -s / or whatever you like
```

## Usage

```shell
Usage:
    fjira [command]
    fjira [flags]
    fjira [jira-issue] [flags]

Available Commands:
    workspace               Switch fjira workspace
    help               	    Help
    version                 Show version

Flags:
    -p, --project           Search for issues withing project, example: -p GEN.
    -w, --workspace         Use different fjira workspace, example: -w myworkspace
```
