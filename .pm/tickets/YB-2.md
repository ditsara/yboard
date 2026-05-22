---
blocks:
    - YB-3
created_at: "2026-05-22T17:00:56Z"
id: YB-2
parent: YB-1
priority: medium
status: backlog
title: Project scaffold & Go module setup
type: task
updated_at: "2026-05-22T17:05:23Z"
---


# Description


[Claude Sonnet 4.6]

## Problem Statement
The project has no Go module or dependency setup. Before any code can be written, the repository needs a `go.mod`, the required dependencies vendored or tidied, and a basic entry point that can be compiled.

## Solution Approach
- Run `go mod init` with an appropriate module path
- Add `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/lipgloss` as direct dependencies
- Create a minimal `main.go` that starts a bubbletea program (empty model is fine at this stage)
- Ensure `go build ./...` succeeds

## Implementation Steps
- [ ] Run `go mod init github.com/ditsara/yboard` (or appropriate module path)
- [ ] `go get github.com/charmbracelet/bubbletea`
- [ ] `go get github.com/charmbracelet/lipgloss`
- [ ] Create `main.go` with a stub bubbletea program (no-op model)
- [ ] Confirm `go build ./...` and `go vet ./...` pass with no errors

## Acceptance Criteria
- [ ] `go.mod` and `go.sum` exist and reference bubbletea and lipgloss
- [ ] `go build ./...` produces a binary without errors
- [ ] Running the binary opens and immediately closes an alternate-screen session (via `tea.WithAltScreen()`) on F10 or Ctrl+C
