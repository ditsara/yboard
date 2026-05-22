---
created_at: "2026-05-22T17:00:56Z"
depends_on:
    - YB-4
    - YB-12
    - YB-13
id: YB-7
parent: YB-1
priority: medium
status: backlog
title: Direct input mode
type: task
updated_at: "2026-05-22T17:05:23Z"
---




# Description


[Claude Sonnet 4.6]

## Problem Statement
In `DirectMode`, every alphabetical and symbol key press must be looked up in the active language's `DirectMap` (lowercase) or `ShiftDirectMap` (shifted), and the result appended to `wordBuffer`. Unknown keys fall through as literal English characters with a status warning.

## Solution Approach
Implement `handleDirectInput(m model, key tea.KeyMsg) model` inside `updateTyping`. Detect whether the key is shifted, normalize to the lookup key string, probe the map, append to buffer, and set `statusMessage` on miss.

## Key Handling Logic
```
if key is shifted:
    result, ok = ShiftDirectMap[key.String()]   // e.g. "A", "!", "@"
else:
    result, ok = DirectMap[key.String()]         // e.g. "a", "1", ";"
if ok:
    wordBuffer = append(wordBuffer, []rune(result)...)
    statusMessage = ""
else:
    wordBuffer = append(wordBuffer, []rune(key.String())...)
    statusMessage = fmt.Sprintf("⚠ Unknown key: '%s' — passed through", key.String())
```

## Special Keys (handled before DirectMap lookup)
- **Spacebar** → append `' '` rune, clear `searchQuery`
- **Backspace** → pop last rune from `wordBuffer` (if non-empty)
- **Enter** → clipboard copy + clear buffer (see YB-11)
- **Tab** → switch to SearchMode (handled in YB-4 dispatcher)
- **F2/F3/F4/F9** → handled by higher-level dispatcher

## Implementation Steps
- [ ] Implement `handleDirectInput` function
- [ ] Handle Spacebar → append space
- [ ] Handle Backspace → pop last rune
- [ ] On alphabetic/symbol key: detect shift, lookup map, append result or literal
- [ ] Set `statusMessage` for unknown keys
- [ ] Clear `statusMessage` at the start of each key event (before processing)
- [ ] Wire into `updateTyping` when `m.inputMode == DirectMode`

## Acceptance Criteria
- [ ] Pressing `a` in DirectMode appends the mapped Thai/Spanish character
- [ ] Pressing `Shift+A` appends the ShiftDirectMap character
- [ ] An unmapped key appends the literal English char and shows a `⚠` status message
- [ ] Backspace removes the last rune from the buffer
- [ ] Spacebar appends a literal space
- [ ] Status message is cleared on the next key event
- [ ] No crash when `wordBuffer` is empty and Backspace is pressed

## Edge Cases
- Some keys in ShiftDirectMap are symbols (e.g. `"Q": "๐"`); the key string from bubbletea for `Shift+Q` is `"Q"` — confirm this assumption against bubbletea's key reporting
- If no language is enabled, all key input is blocked (show "⚠ No languages enabled")
