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
    fjira [OPTIONS]
    fjira JIRA-TICKET [OPTIONS]
    fjira workspace
    fjira version

Optional options:
    -p, --project               Search for issues withing project, example: GEN.
    -i, --issue                 Open Jira Issue, example: GEN-123.
    -w, --workspace             Use fjira workspace, example: myworkspace
```
