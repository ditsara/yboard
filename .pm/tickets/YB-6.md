---
created_at: "2026-05-22T17:00:56Z"
depends_on:
    - YB-4
id: YB-6
parent: YB-1
priority: medium
status: done
title: Visual keyboard grid rendering
type: task
updated_at: "2026-05-22T17:40:43Z"
---




# Description


[Claude Sonnet 4.6]

## Problem Statement
Each language module ships with `KeyboardRows [][]VisualKey`. These rows must be rendered as an onscreen keyboard grid showing both Normal (bright) and Shift (dimmed) characters simultaneously inside per-key rounded boxes.

## Solution Approach
Implement `renderKeyboard(m model) string`. For each row, render each `VisualKey` as a small lipgloss box, then join the row horizontally, then join all rows vertically.

## Rendering Spec (from SPEC §4.5)
- Each `VisualKey` is a single bordered box (rounded border) rendered with `lipgloss.JoinVertical`:
  - Top line: `Shift` char — color `#888888` (dimmed)
  - Bottom line: `Normal` char — color `#FFFFFF`, bold
- Keys in a row are joined with `lipgloss.JoinHorizontal`
- Rows are joined with `lipgloss.JoinVertical`
- When no language is active (all disabled), return an empty string

## Implementation Steps
- [ ] Create `view_keyboard.go` with `renderKeyboard(m model) string`
- [ ] Define lipgloss styles: `shiftStyle` (`#888888`), `normalStyle` (`#FFFFFF` bold), `keyBoxStyle` (rounded border)
- [ ] For each `VisualKey`, build a two-line vertical join (Shift on top, Normal below) inside a rounded box
- [ ] Join keys in each row horizontally
- [ ] Join all rows vertically
- [ ] Return empty string when no active language module is available
- [ ] Integrate return value into `viewTyping` layout (replaces placeholder)

## Acceptance Criteria
- [ ] Each key renders as a rounded-border box with Shift char (dimmed) above Normal char (bold/white)
- [ ] Keys within a row are horizontally adjacent
- [ ] All rows are vertically stacked in correct order (number row → QWERTY → home row → bottom row)
- [ ] Keyboard disappears gracefully when no language is enabled
- [ ] No panic on empty `KeyboardRows` slice

## Edge Cases
- Rows can have different lengths (Thai keyboard rows are uneven) — each row is independently joined
- Multi-byte Thai/Unicode characters must render at correct cell width (lipgloss handles this via `runewidth`)

---

## TUI Design Mockup

### Single Key Box Anatomy

Each `VisualKey` renders as a **5-wide × 4-tall** rounded box. The physical English
key label is embedded in the top border (replacing the center dash). Inside:
Shift on top (dimmed), Normal below (bold).

```
╭─q─╮
│ ม │   ← Shift char  (#888888, normal weight)
│ ท │   ← Normal char (#FFFFFF, bold)
╰───╯
  ↑
  q = #FFD700 (gold) — "the key you press"
```

The trick: `╭─q─╮` is the same 5-char width as `╭───╮` — no layout change.
The `q` color (`#FFD700`) is distinct from both the Shift grey and the Normal white,
so the user instantly reads: gold = physical key, white = what you get, grey = shifted.

For keys whose `Normal` **and** `Shift` are both empty — i.e. `EmptyKey("k")` — the key
is not remapped by this language. Render it as a **muted/dim** box: key label in `#555555`
(not gold), both slots show `·` in `#333333`. This tells the user "this key passes through unchanged":

```
╭─[─╮          ← key label #555555 (not gold — signals "not remapped")
│ · │          ← #333333
│ · │          ← #333333
╰───╯
```

For keys where only `Shift` is empty (Normal is set), render `·` for the Shift slot only:

```
╭─k─╮
│ · │   ← no Shift variant (#444444)
│ า │
╰───╯
```

---

### Thai (Kedmanee) Full Keyboard Grid

Rows are **centered** relative to the widest row (row 2 = 13 keys = 65 chars).
Shorter rows are padded symmetrically with `lipgloss.PlaceHorizontal`.

Row widths: row 1 = 60 chars (pad 2), row 2 = 65 (pad 0), row 3 = 55 (pad 5), row 4 = 50 (pad 7).

```
  ╭─`─╮╭─1─╮╭─2─╮╭─3─╮╭─4─╮╭─5─╮╭─6─╮╭─7─╮╭─8─╮╭─9─╮╭─0─╮╭─-─╮        ← pad 2
  │ + ││ ๑ ││ ๒ ││ ๓ ││ ๔ ││ ู ││ ฿ ││ ๕ ││ ๖ ││ ๗ ││ ๘ ││ ๙ │
  │ ๅ ││ / ││ _ ││ ภ ││ ถ ││ ุ ││ ึ ││ ค ││ ต ││ จ ││ ข ││ ช │
  ╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
╭─q─╮╭─w─╮╭─e─╮╭─r─╮╭─t─╮╭─y─╮╭─u─╮╭─i─╮╭─o─╮╭─p─╮╭─[─╮╭─]─╮╭─\─╮      ← no pad (widest)
│ ๐ ││ " ││ ฎ ││ ฑ ││ ธ ││ ํ ││ ๊ ││ ณ ││ ฯ ││ ญ ││ ฐ ││ , ││ ฅ │
│ ๆ ││ ไ ││ ำ ││ พ ││ ะ ││ ั ││ ี ││ ร ││ น ││ ย ││ บ ││ ล ││ ฃ │
╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
     ╭─a─╮╭─s─╮╭─d─╮╭─f─╮╭─g─╮╭─h─╮╭─j─╮╭─k─╮╭─l─╮╭─;─╮╭─'─╮           ← pad 5
     │ ฤ ││ ฆ ││ ฏ ││ โ ││ ฌ ││ ็ ││ ๋ ││ ษ ││ ศ ││ ซ ││ . │
     │ ฟ ││ ห ││ ก ││ ด ││ เ ││ ้ ││ ่ ││ า ││ ส ││ ว ││ ง │
     ╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
        ╭─z─╮╭─x─╮╭─c─╮╭─v─╮╭─b─╮╭─n─╮╭─m─╮╭─,─╮╭─.─╮╭─/─╮             ← pad 7
        │ ( ││ ) ││ ฉ ││ ฮ ││ ฺ ││ ์ ││ ? ││ ฒ ││ ฬ ││ ฾ │
        │ ผ ││ ป ││ แ ││ อ ││ ิ ││ ื ││ ท ││ ม ││ ใ ││ ฝ │
        ╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
```

---

### Spanish (Standard) Full Keyboard Grid

```
  ╭─1─╮╭─2─╮╭─3─╮╭─4─╮╭─5─╮╭─6─╮╭─7─╮╭─8─╮╭─9─╮╭─0─╮                  ← pad 2
  │ ! ││ @ ││ # ││ $ ││ % ││ ^ ││ & ││ * ││ ( ││ ) │
  │ 1 ││ 2 ││ 3 ││ 4 ││ 5 ││ 6 ││ 7 ││ 8 ││ 9 ││ 0 │
  ╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
╭─q─╮╭─w─╮╭─e─╮╭─r─╮╭─t─╮╭─y─╮╭─u─╮╭─i─╮╭─o─╮╭─p─╮                    ← no pad (widest)
│ Q ││ W ││ E ││ R ││ T ││ Y ││ U ││ I ││ O ││ P │
│ q ││ w ││ e ││ r ││ t ││ y ││ u ││ i ││ o ││ p │
╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
  ╭─a─╮╭─s─╮╭─d─╮╭─f─╮╭─g─╮╭─h─╮╭─j─╮╭─k─╮╭─l─╮╭─;─╮                  ← pad 2
  │ A ││ S ││ D ││ F ││ G ││ H ││ J ││ K ││ L ││ Ñ │
  │ a ││ s ││ d ││ f ││ g ││ h ││ j ││ k ││ l ││ ñ │
  ╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
     ╭─z─╮╭─x─╮╭─c─╮╭─v─╮╭─b─╮╭─n─╮╭─m─╮╭─;─╮                          ← pad 6
     │ Z ││ X ││ C ││ V ││ B ││ N ││ M ││ Ç │
     │ z ││ x ││ c ││ v ││ b ││ n ││ m ││ ç │
     ╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯╰───╯
```

---

### Lipgloss Style Definitions

```go
var (
    shiftStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#888888"))

    normalStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")).
        Bold(true)

    keyLabelStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFD700")) // gold — the physical key to press

    keyBoxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(0, 1). // 1-cell padding each side → 5 chars total width
        Align(lipgloss.Center)
)
```

Rendering each key (English label injected into the top border character):

```go
func renderKey(vk VisualKey) string {
    isEmpty := vk.Normal == "" && vk.Shift == ""

    // Choose styles based on whether the key is remapped by this language.
    labelSty := keyLabelStyle          // #FFD700 gold
    shiftSty := shiftStyle             // #888888
    normSty  := normalStyle            // #FFFFFF bold
    boxSty   := keyBoxStyle
    if isEmpty {
        labelSty = mutedLabelStyle     // #555555 — signals "not remapped"
        shiftSty = emptySlotStyle      // #333333
        normSty  = emptySlotStyle      // #333333
    }

    shiftChar := vk.Shift
    if shiftChar == "" { shiftChar = "·" }
    normalChar := vk.Normal
    if normalChar == "" { normalChar = "·" }

    top    := shiftSty.Render(shiftChar)
    bottom := normSty.Render(normalChar)
    inner  := lipgloss.JoinVertical(lipgloss.Center, top, bottom)
    box    := boxSty.Render(inner)

    // Inject the English key label into the center of the top border line.
    // The top border of a 5-wide rounded box is "╭───╮"; replace index [2] with Key.
    lines     := strings.Split(box, "\n")
    label     := labelSty.Render(vk.Key)
    lines[0]   = lines[0][:2] + label + lines[0][3:]  // replace center dash
    return strings.Join(lines, "\n")
}
```

Centering rows with `lipgloss.PlaceHorizontal`:

```go
func renderKeyboard(m model) string {
    if /* no active language */ { return "" }
    mod := activeLanguage(m)

    maxWidth := 0
    renderedRows := make([]string, len(mod.KeyboardRows))
    for i, row := range mod.KeyboardRows {
        keys := make([]string, len(row))
        for j, vk := range row { keys[j] = renderKey(vk) }
        renderedRows[i] = lipgloss.JoinHorizontal(lipgloss.Top, keys...)
        if w := lipgloss.Width(renderedRows[i]); w > maxWidth { maxWidth = w }
    }
    for i, row := range renderedRows {
        renderedRows[i] = lipgloss.PlaceHorizontal(maxWidth, lipgloss.Center, row)
    }
    return lipgloss.JoinVertical(lipgloss.Left, renderedRows...)
}
```
