---
created_at: "2026-05-22T17:00:56Z"
depends_on:
    - YB-4
id: YB-9
parent: YB-1
priority: medium
status: done
title: Setup screen
type: task
updated_at: "2026-05-22T17:42:25Z"
---




# Description


[Claude Sonnet 4.6]

## Problem Statement
Users need to enable/disable individual language modules without restarting the app. The Setup screen provides a simple navigable checklist of all registered `LanguageModule` entries.

## Solution Approach
When `m.state == StateSetup`, `View()` delegates to `viewSetup(m)` and `Update()` delegates to `updateSetup(m, msg)`. The setup screen shows a vertical list of modules with `[X]` / `[ ]` indicators.

## UI Spec
```
Language Modules
────────────────
▶ [X] Thai (Kedmanee)
  [ ] Greek
  [X] Spanish Standard

Space: toggle  ↑/↓ or k/j: navigate  F2/Enter/Esc: return
```

## Keybindings (StateSetup only)
| Key | Action |
|-----|--------|
| Up / `k` | Move selection up (clamp at top) |
| Down / `j` | Move selection down (clamp at bottom) |
| Space | Toggle `Enabled` on highlighted module |
| F2 / Enter / Esc | Return to `StateTyping` |

## Implementation Steps
- [ ] Create `view_setup.go` with `viewSetup(m model) string`
- [ ] Render title + divider
- [ ] Render each language module as `▶ [X] Name` (selected) or `  [X] Name` (not selected)
- [ ] Create `updateSetup(m model, msg tea.Msg) (model, tea.Cmd)` in `update_setup.go`
- [ ] Handle Up/k → decrement `m.setupCursor`, clamp to 0
- [ ] Handle Down/j → increment `m.setupCursor`, clamp to `len(languages)-1`
- [ ] Handle Space → toggle `m.languages[m.setupCursor].Enabled`
- [ ] Handle F2/Enter/Esc → set `m.state = StateTyping`
- [ ] Add `setupCursor int` field to model struct
- [ ] Wire `viewSetup` and `updateSetup` into `View` and `Update`

## Acceptance Criteria
- [ ] Setup screen is accessible via F2 from the typing view
- [ ] All registered language modules appear in the list
- [ ] `[X]` shows for enabled modules, `[ ]` for disabled
- [ ] Cursor wraps correctly (or clamps — pick one and document it; spec says clamp)
- [ ] Space toggles the `Enabled` field on the highlighted module
- [ ] F2 / Enter / Esc returns to typing view
- [ ] After toggling, the typing view reflects the change immediately (F3/F4 skip disabled modules)

## Edge Cases
- If the user disables all modules and returns to typing view, the "⚠ No languages enabled" warning appears (handled in YB-5/YB-7)
- `setupCursor` should persist across Setup↔Typing transitions so the user lands on the same row

---

## TUI Design Mockup

### Setup Screen — default state (Thai selected, Spanish enabled)

```
╭──────────────────────────────────────────────────────────────────────╮
│  Language Modules                                                     │
╰──────────────────────────────────────────────────────────────────────╯

  ▶  [✓]  Thai (Kedmanee)
     [✓]  Spanish Standard
     [ ]  (future module placeholder)



 Space: toggle   ↑/k: up   ↓/j: down   F2/Enter/Esc: return to typing
```

---

### Setup Screen — after disabling Thai (cursor on Thai)

```
╭──────────────────────────────────────────────────────────────────────╮
│  Language Modules                                                     │
╰──────────────────────────────────────────────────────────────────────╯

  ▶  [ ]  Thai (Kedmanee)
     [✓]  Spanish Standard



 Space: toggle   ↑/k: up   ↓/j: down   F2/Enter/Esc: return to typing
```

---

### Setup Screen — all disabled (warning state)

```
╭──────────────────────────────────────────────────────────────────────╮
│  Language Modules                                                     │
╰──────────────────────────────────────────────────────────────────────╯

  ▶  [ ]  Thai (Kedmanee)
     [ ]  Spanish Standard

  ⚠  No languages enabled. Return to typing and configure at least one.



 Space: toggle   ↑/k: up   ↓/j: down   F2/Enter/Esc: return to typing
```

---

### Styling Notes

| Element | Style |
|---------|-------|
| Title box | `lipgloss.RoundedBorder()`, full width |
| `▶` cursor indicator | `#00FF88` (green), bold |
| `[✓]` checked | `#00FF88` (green) |
| `[ ]` unchecked | `#555555` (muted grey) |
| Module name (selected row) | `#FFFFFF`, bold |
| Module name (other rows) | `#AAAAAA`, normal |
| Warning line `⚠ …` | `#FF4444` (red) |
| Footer bar | `#555555`, same style as typing view footer |

The cursor `▶` is left-aligned in a fixed 3-char column. Module rows are indented
by 5 chars (3 for cursor column + 2 padding) to create a clean checklist appearance.

```
 Col: 0    3    7
       ▶  [✓]  Thai (Kedmanee)
          [✓]  Spanish Standard
```
