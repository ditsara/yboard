---
assignee: ""
blocks: []
created_at: "2026-05-22T17:00:49Z"
depends_on: []
id: YB-1
labels: []
milestones: []
parent: ""
points: 0
priority: medium
related: []
status: backlog
title: 'YBoard: Multi-Language TUI Keyboard Engine'
type: epic
updated_at: "2026-05-22T17:05:23Z"
---

# Description


[Claude Sonnet 4.6]

Build a single-binary TUI application in Go that lets users type in alternative languages (Thai, Spanish) using a standard English keyboard — without touching OS-level keyboard settings. Intercepts keystrokes, translates them via configurable language maps, renders a visual keyboard reference, and pipes composed text to the system clipboard.

**Tech Stack:** Go · `github.com/charmbracelet/bubbletea` · `github.com/charmbracelet/lipgloss`  
**Target:** WSL, Linux, native Windows

## Child Tickets

| ID | Title |
|----|-------|
| YB-2  | Project scaffold & Go module setup |
| YB-3  | Core data structures & types |
| YB-4  | BubbleTea model skeleton (Init/Update/View) |
| YB-5  | Typing view UI rendering |
| YB-6  | Visual keyboard grid rendering |
| YB-7  | Direct input mode |
| YB-8  | Search input mode |
| YB-9  | Setup screen |
| YB-10 | Language switching (F3/F4) |
| YB-11 | Clipboard integration |
| YB-12 | Thai (Kedmanee) language module |
| YB-13 | Spanish Standard language module |
