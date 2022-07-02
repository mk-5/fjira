# Fjira - the fuzziest Jira command line tool in the world.

![Fira](fjira.png)

![Test](https://github.com/mk-5/fjira/actions/workflows/tests.yml/badge.svg)

## Demo

## Features

- Search for Jira Projects, and Issues
- Changes Jira Issue assignee
- Changes Jira Issue status
- Supports multiple workspaces
- Runs on MacOS and Linux

## Getting started

```shell
brew install fjira
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
