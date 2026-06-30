# AGENTS.md

## Purpose

This file gives instructions to AI coding agents working in this repository. Follow these rules to work safely, efficiently, and with minimal token usage.

## Token-Saving Rules

- Read only files relevant to the task. Do not scan the whole repo.
- Use `grep`/`glob` search before opening files.
- Summarize findings; do not paste large file contents.
- Avoid re-reading files unless they changed.
- Prefer minimal diffs. Make the smallest change that solves the problem.
- Do not explain obvious code in comments.
- Do not refactor unrelated code.
- Ask for clarification only when blocked, not speculatively.
- Follow existing patterns in the repo — search for similar code before inventing new patterns.

## Project Overview

**fjira** is a fuzzy-finder TUI (terminal UI) application for Jira, written in Go. It provides interactive fuzzy search for Jira projects and issues, status transitions, assignee changes, comments, JQL searches, and multi-workspace support.

Module: `github.com/mk-5/fjira`

## Repository Structure

```
cmd/fjira-cli/    — CLI entry point (main.go)
internal/         — all application packages
  app/            — core TUI framework (rendering, state, colors)
  jira/           — Jira API client and data models
  fjira/          — main app coordination / wiring
  issues/         — issue list, search, detail views
  projects/       — project selection
  boards/         — board views (scrum/kanban)
  filters/        — saved JQL filters
  statuses/       — status transitions
  labels/         — label management
  comments/       — comment creation/viewing
  users/          — user/assignee operations
  workspaces/     — multi-workspace support
  ui/             — shared UI utilities
  os/             — OS-level helpers
assets/           — embedded assets
docs/             — documentation
features/         — BDD feature specs (bats)
scripts/          — build/release scripts
bats/             — bash test helpers
```

## Common Commands

```bash
# Install / vendor dependencies
make install

# Build binary → out/bin/fjira
make build

# Run tests
make test

# Run tests with coverage report
make test_coverage

# Lint (golangci-lint, matches CI)
golangci-lint run

# Clean build artifacts
make clean
```

There is no separate format command — `gofmt` / `goimports` is expected as part of editor tooling. CI enforces lint via `.golangci.yml`.

## Coding Guidelines

- Follow existing Go code style in the package you are editing.
- Keep changes small and focused. One concern per PR.
- Do not introduce new dependencies unless necessary; prefer stdlib or existing vendored libs.
- Update tests when behavior changes.
- Update docs or comments when public behavior or API changes.
- Colors and theme constants live in `internal/app/colors.go` — add new ones there.
- TUI components follow the pattern in `internal/app/` — read existing views before adding a new one.

## Testing Guidelines

- Run `go test ./internal/<package>/...` for the affected package first.
- Run `make test` (all internal packages) only when your change is cross-cutting.
- Add or update tests for bug fixes and new behavior.
- If tests cannot run (e.g., requires live Jira), note what should be tested manually.

## Agent Workflow

1. Understand the task and identify the affected package(s).
2. Search for relevant files (`grep`, `glob`) before reading.
3. Read only necessary files.
4. Make minimal changes following existing patterns.
5. Run `go build ./...` and `make test` (or scoped test).
6. Summarize: what changed, what was tested, any follow-up risks.

## Response Style for Agents

- Be concise. Lead with the action or finding.
- Mention only files actually changed.
- Do not paste large code blocks unless asked.
- Report commands run and their outcomes.
- Flag any risks or follow-up work needed.

## Do Not

- Do not scan the entire repo without a specific reason.
- Do not perform unrelated cleanup or style fixes.
- Do not rewrite large files unnecessarily.
- Do not add speculative or future-proofing features.
- Do not duplicate documentation that already exists in README or docs/.
- Do not commit generated files (coverage.out, coverage.html, out/, dist/).
