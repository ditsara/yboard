---
blocks:
    - YB-4
    - YB-12
    - YB-13
created_at: "2026-05-22T17:00:56Z"
depends_on:
    - YB-2
id: YB-3
parent: YB-1
priority: medium
status: done
title: Core data structures & types
type: task
updated_at: "2026-05-22T17:30:52Z"
---







# Description


[Claude Sonnet 4.6]

## Problem Statement
All application logic depends on shared Go types (`VisualKey`, `LanguageModule`, `AppState`, `InputMode`). These must be defined before any other feature ticket can be implemented.

## Solution Approach
Create a `types.go` (or `internal/types/types.go`) file that declares all structs and enums exactly as specified in the SPEC. No logic here — pure data definitions.

Also define the canonical keyboard layout constant and helper functions that enforce a shared contract across all language modules.

## Data Structures to Define

```go
type VisualKey struct {
    Key    string // Physical English key label (e.g. "q", "1", ";") — rendered in top border
    Normal string // Output without Shift; empty string = key is not remapped by this language
    Shift  string // Output with Shift;   empty string = key is not remapped by this language
}

type LanguageModule struct {
    ID             string
    Name           string
    Enabled        bool
    DirectMap      map[string]string
    ShiftDirectMap map[string]string
    PhoneticMap    map[string][]string
    // KeyboardRows must always match StandardRowLengths exactly.
    // Use EmptyKey() for physical keys this language does not remap.
    KeyboardRows   [][]VisualKey
}

type AppState int
const (
    StateTyping AppState = iota
    StateSetup
)

type InputMode int
const (
    DirectMode InputMode = iota
    SearchMode
)
```

## Canonical Keyboard Layout

All `LanguageModule.KeyboardRows` **must** conform to standard US QWERTY dimensions:

```go
// StandardRowLengths defines the required number of keys per row for all language modules.
// Rows correspond to: number row, QWERTY row, home row, bottom row.
// Physical key positions: [` 1 2 3 4 5 6 7 8 9 0 - =] [q…\] [a…'] [z…/]
var StandardRowLengths = []int{13, 13, 11, 10}

// EmptyKey creates a VisualKey for a physical key this language does not remap.
// The renderer draws it as a muted/dim box — only the key label is visible.
func EmptyKey(physicalKey string) VisualKey {
    return VisualKey{Key: physicalKey} // Normal and Shift are zero-value ""
}

// ValidateModule returns an error if a module's KeyboardRows do not match StandardRowLengths.
// Call this once at startup for each registered module.
func ValidateModule(m LanguageModule) error {
    if len(m.KeyboardRows) != len(StandardRowLengths) {
        return fmt.Errorf("%s: expected %d keyboard rows, got %d",
            m.ID, len(StandardRowLengths), len(m.KeyboardRows))
    }
    for i, row := range m.KeyboardRows {
        if len(row) != StandardRowLengths[i] {
            return fmt.Errorf("%s: row %d expected %d keys, got %d",
                m.ID, i, StandardRowLengths[i], len(row))
        }
    }
    return nil
}
```

### Why not fixed-size Go arrays?

`[13]VisualKey` as a field type would guarantee dimensions at compile time, but it
breaks the renderer loop (can't range over a mixed-length array field without
reflection or code duplication). Runtime validation at startup is the idiomatic Go
tradeoff: the constraint is documented and enforced without complicating the hot path.

### Implication for language modules

| Module | Keys to fill with `EmptyKey()` |
|--------|-------------------------------|
| Thai (Kedmanee) | `=` in number row (spec omits it; Kedmanee maps it to ว/ซ — confirm in YB-12) |
| Spanish Standard | `` ` ``, `-` in number row; `[`, `]`, `\` in QWERTY row; `'` in home row; `,`, `.`, `/` in bottom row |

## Implementation Steps
- [ ] Create `types.go` at package root (or an `internal/types` package)
- [ ] Define `VisualKey` struct with `Key`, `Normal`, and `Shift` fields
- [ ] Define `LanguageModule` struct with all fields; document `KeyboardRows` contract in comment
- [ ] Define `StandardRowLengths = []int{13, 13, 11, 10}`
- [ ] Implement `EmptyKey(physicalKey string) VisualKey`
- [ ] Implement `ValidateModule(m LanguageModule) error`
- [ ] Define `AppState` enum (`StateTyping`, `StateSetup`)
- [ ] Define `InputMode` enum (`DirectMode`, `SearchMode`)
- [ ] `go vet ./...` passes

## Acceptance Criteria
- [ ] All four types are exported and correctly defined
- [ ] `StandardRowLengths` is exported and equals `[13, 13, 11, 10]`
- [ ] `EmptyKey("q")` returns `VisualKey{Key: "q", Normal: "", Shift: ""}`
- [ ] `ValidateModule` returns `nil` for a correctly-dimensioned module
- [ ] `ValidateModule` returns a descriptive error for wrong row count or wrong row length
- [ ] File compiles with no errors or warnings
- [ ] No side effects (no `init()`, no globals other than the constants)
