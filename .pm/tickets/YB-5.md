---
created_at: "2026-05-22T17:00:56Z"
depends_on:
    - YB-4
id: YB-5
parent: YB-1
priority: medium
status: done
title: Typing view UI rendering
type: task
updated_at: "2026-05-22T17:38:36Z"
---




# Description


[Claude Sonnet 4.6]

## Problem Statement
The typing view (`StateTyping`) needs a complete UI layout: word buffer box, mode badge, status message area, and hotkeys footer. The visual keyboard grid is handled separately in YB-6.

## Solution Approach
Implement `viewTyping(m model) string` using lipgloss to compose the layout. All elements are assembled with `lipgloss.JoinVertical`. Sizes should respect `m.termWidth`.

## Layout (top to bottom)
1. **Word Buffer Box** — bordered box, prominent, shows `string(m.wordBuffer)` or placeholder `"Start typing…"` when empty
2. **Mode Badge** — inline badge: `[DIRECT]` or `[SEARCH]` + current language name (e.g. `Thai (Kedmanee)`)
3. **Search Query Line** (SearchMode only) — shows `"Search: " + m.searchQuery` and numbered candidates from PhoneticMap lookup
4. **Status Message** — single line, dimmed; shows `m.statusMessage` (cleared on next keypress)
5. **Visual Keyboard Grid** — placeholder call to `renderKeyboard(m)` (implemented in YB-6)
6. **Hotkeys Footer** — always visible: `F2:Setup  F3:← Lang  F4:→ Lang  F9:Copy  Enter:Copy+Clear  Tab:Mode  F10:Quit`

## "No Languages Enabled" State
When all `LanguageModule.Enabled` are false, the buffer box shows `"⚠ No languages enabled — press F2 to configure"` and the mode badge is hidden.

## Implementation Steps
- [ ] Create `view_typing.go` with `viewTyping(m model) string`
- [ ] Implement word buffer box with `lipgloss` border style
- [ ] Implement mode badge (Direct / Search) + active language name
- [ ] Implement search query + candidate list (visible in SearchMode)
- [ ] Implement status message area (single line, dimmed color)
- [ ] Add placeholder call `renderKeyboard(m)` (returns `""` until YB-6 is done)
- [ ] Implement hotkeys footer bar
- [ ] Handle "no languages enabled" state in buffer area
- [ ] Wire `viewTyping` into `View()` when `m.state == StateTyping`

## Acceptance Criteria
- [ ] Word buffer content is displayed inside a bordered box
- [ ] Mode badge correctly reflects DirectMode / SearchMode
- [ ] Active language name is shown next to the mode badge
- [ ] Hotkeys footer is always visible at the bottom
- [ ] Status message area appears between the keyboard and the footer
- [ ] "No languages enabled" warning renders when appropriate
- [ ] Layout does not crash on small terminal sizes (gracefully truncates)

## Edge Cases
- Empty `wordBuffer` shows a placeholder, not an empty box
- `statusMessage` is empty string when no message — area still reserves one line to avoid layout jumps

---

## TUI Design Mockup

### State 1 — DirectMode, buffer has text, status message visible

```
╭──────────────────────────────────────────────────────────────────────╮
│  สวัสดีครับ ยินดีที่ได้รู้จัก                                          │
╰──────────────────────────────────────────────────────────────────────╯

 ╔════════╗  Thai (Kedmanee)
 ║ DIRECT ║
 ╚════════╝

 📋 Copied text block to clipboard!

  ╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮
  │ + ││ ๑ ││ ๒ ││ ๓ ││ ๔ ││ ู ││ ฿ ││ ๕ ││ ๖ ││ ๗ ││ ๘ ││ ๙ │
  │ ๅ ││ / ││ _ ││ ภ ││ ถ ││ ุ ││ ึ ││ ค ││ ต ││ จ ││ ข ││ ช │
  ╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
  ╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮
  │ ๐ ││ " ││ ฎ ││ ฑ ││ ธ ││ ํ ││ ๊ ││ ณ ││ ฯ ││ ญ ││ ฐ ││ , ││ ฅ │
  │ ๆ ││ ไ ││ ำ ││ พ ││ ะ ││ ั ││ ี ││ ร ││ น ││ ย ││ บ ││ ล ││ ฃ │
  ╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
  ╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮
  │ ฤ ││ ฆ ││ ฏ ││ โ ││ ฌ ││ ็ ││ ๋ ││ ษ ││ ศ ││ ซ ││ . │
  │ ฟ ││ ห ││ ก ││ ด ││ เ ││ ้ ││ ่ ││ า ││ ส ││ ว ││ ง │
  ╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
  ╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮╭───╮
  │ ( ││ ) ││ ฉ ││ ฮ ││ ฺ ││ ์ ││ ? ││ ฒ ││ ฬ ││ ฾ │
  │ ผ ││ ป ││ แ ││ อ ││ ิ ││ ื ││ ท ││ ม ││ ใ ││ ฝ │
  ╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯

 F2:Setup  F3:← Lang  F4:→ Lang  F9:Copy  Enter:Copy+Clear  Tab:Mode  F10:Quit
```

**Styling notes:**
- Word buffer box: `lipgloss.RoundedBorder()`, full terminal width, min 3 lines tall
- Mode badge `║ DIRECT ║`: `lipgloss.DoubleBorder()`, background `#3C3C3C`, foreground `#00FF88` (green)
  — SearchMode badge uses `#FF8800` (amber) instead
- Language name: plain text, right of the mode badge, `#AAAAAA`
- Status message line: `#888888` (dimmed); always occupies one line even when empty (prevents layout jump)
- Keyboard grid: indented 2 spaces, see YB-6 for key-box detail
- Footer: `#555555` (muted), single line

---

### State 2 — SearchMode, query active

```
╭──────────────────────────────────────────────────────────────────────╮
│  สวัสดี                                                                │
╰──────────────────────────────────────────────────────────────────────╯

 ╔════════╗  Thai (Kedmanee)
 ║ SEARCH ║
 ╚════════╝

 Search: s▌
 1:ส  2:ษ  3:ศ  4:ซ

                                          ← (status line, empty = 1 blank line)

  ╭─q─╮╭─w─╮  ...keyboard rows (English key in gold on top border, centered)...
  │ ฤ ││ ฆ │
  │ ฟ ││ ห │
  ╰───╯╰───╯

 F2:Setup  F3:← Lang  F4:→ Lang  F9:Copy  Enter:Copy+Clear  Tab:Mode  F10:Quit
```

**Styling notes:**
- `Search: s▌` — `s` in bright white, `▌` is the cursor (blinking block or `_`)
- Candidate line: each candidate `N:char` separated by two spaces; numbers in `#888888`, chars in `#FFFFFF` bold
- When query has no PhoneticMap match, candidate line is blank (still reserved)

---

### State 3 — No languages enabled

```
╭──────────────────────────────────────────────────────────────────────╮
│  ⚠ No languages enabled — press F2 to configure                      │
╰──────────────────────────────────────────────────────────────────────╯

                                          ← (mode badge hidden)

                                          ← (status line)

                                          ← (keyboard hidden)

 F2:Setup  F3:← Lang  F4:→ Lang  F9:Copy  Enter:Copy+Clear  Tab:Mode  F10:Quit
```

**Styling notes:**
- Buffer box text `⚠ …` in `#FF4444` (red/warning)
- Mode badge, search line, keyboard: all hidden (no empty boxes — they simply don't render)
- Footer always visible
