---
blocks:
    - YB-5
    - YB-6
    - YB-7
    - YB-8
    - YB-9
    - YB-10
    - YB-11
created_at: "2026-05-22T17:00:56Z"
depends_on:
    - YB-3
id: YB-4
parent: YB-1
priority: medium
status: backlog
title: BubbleTea model skeleton (Init/Update/View)
type: task
updated_at: "2026-05-22T17:05:23Z"
---









# Description


[Claude Sonnet 4.6]

## Problem Statement
The application needs a working bubbletea `Model` that holds all runtime state and correctly routes `Init`, `Update`, and `View` calls. This is the central integration point for all other features.

## Solution Approach
Define the `model` struct and implement the three bubbletea interface methods. At this stage, `Update` only needs to handle global keys (F10, Ctrl+L) and dispatch to sub-handlers (stubbed). `View` returns a placeholder string. The goal is a compiling, runnable program with the alternate screen active.

## Model State Fields
```go
type model struct {
    state         AppState
    inputMode     InputMode
    languages     []LanguageModule   // all registered modules
    activeIndex   int                // index into enabled languages
    wordBuffer    []rune
    searchQuery   string
    statusMessage string
    termWidth     int
    termHeight    int
}
```

## Keybindings to Wire in This Ticket (Global)
- **F10** → `tea.Quit`
- **Ctrl+L** → `tea.ClearScreen`
- `tea.WindowSizeMsg` → store `termWidth` / `termHeight`

## Implementation Steps
- [ ] Define `model` struct with all fields listed above
- [ ] Implement `Init() tea.Cmd` — return `nil`
- [ ] Implement `Update(msg tea.Msg) (tea.Model, tea.Cmd)`:
  - Handle `tea.WindowSizeMsg`
  - Handle `tea.KeyMsg` for F10 and Ctrl+L
  - Dispatch to `updateTyping` or `updateSetup` stubs based on `m.state`
- [ ] Implement `View() string` — return placeholder `"YBoard loading…"`
- [ ] Wire into `main.go` with `tea.WithAltScreen()`
- [ ] Confirm binary runs and exits cleanly on F10

## Acceptance Criteria
- [ ] Program starts in alternate screen mode
- [ ] F10 exits the program cleanly (restores terminal)
- [ ] Ctrl+L clears the screen without crashing
- [ ] Window resize events are captured and stored
- [ ] Dispatches to per-state update stubs (even if they're no-ops)

## Edge Cases
- If no languages are loaded yet, `activeIndex` stays at 0 and typing view shows a "no languages enabled" warning (handled in YB-5/YB-7, not here)
