---
created_at: "2026-05-22T17:00:56Z"
depends_on:
    - YB-4
    - YB-12
    - YB-13
id: YB-8
parent: YB-1
priority: medium
status: done
title: Search input mode
type: task
updated_at: "2026-05-22T17:50:19Z"
---






# Description


[Claude Sonnet 4.6]

## Problem Statement
In `SearchMode`, alphabetic key presses build a `searchQuery` string. Matching entries from `PhoneticMap` are displayed as a numbered candidate list. Number keys `1`–`9` and `0` (= 10th) select and append the candidate to `wordBuffer`, then clear the query.

## Solution Approach
Implement `handleSearchInput(m model, key tea.KeyMsg) model`. Build the query from `a-z` keys (Shift ignored, normalized lowercase). Query the active language's `PhoneticMap`. Display up to 10 candidates. Number keys commit selection.

## Search Logic
```
if key is a-z (any case):
    searchQuery += strings.ToLower(key.Rune)
    // Look up PhoneticMap[searchQuery] — show candidates in view
elif key is '1'-'9':
    idx = int(key) - 1          // 0-based
    candidates = PhoneticMap[searchQuery]
    if idx < len(candidates):
        wordBuffer = append(wordBuffer, []rune(candidates[idx])...)
    searchQuery = ""
elif key is '0':
    idx = 9
    // same as above
elif key is Backspace:
    if searchQuery != "":
        searchQuery = searchQuery[:len-1]   // pop last char
    else:
        pop last rune from wordBuffer
```

## PhoneticMap Lookup Rules
- Only exact key match (e.g. query `"s"` → `PhoneticMap["s"]`)
- At most 10 candidates shown; extras silently ignored
- If no match exists, show empty candidate list (no error)
- Candidate display format: `1:ส  2:ษ  3:ศ  4:ซ`

## Implementation Steps
- [ ] Implement `handleSearchInput(m model, key tea.KeyMsg) model`
- [ ] Build `searchQuery` from a-z keys (normalize to lowercase)
- [ ] Backspace: pop from `searchQuery` if non-empty, else pop from `wordBuffer`
- [ ] Number keys 1-9: select candidate at index `n-1`, clear query
- [ ] Key `0`: select candidate at index 9
- [ ] Candidate list capped at 10 entries
- [ ] Wire into `updateTyping` when `m.inputMode == SearchMode`
- [ ] Candidate list rendered in `viewTyping` search query area (from YB-5)

## Acceptance Criteria
- [ ] Typing `s` in SearchMode shows candidates `1:ส  2:ษ  3:ศ  4:ซ` (Thai module)
- [ ] Pressing `1` appends `ส` to buffer and clears the query
- [ ] Pressing `0` selects the 10th candidate if available
- [ ] Backspace clears one character from the query; pressing again pops from buffer
- [ ] A query with no PhoneticMap match shows an empty candidate list (no crash)
- [ ] Candidates beyond index 9 are never shown
- [ ] Shift is ignored on alphabetic input (treated as lowercase)

## Edge Cases
- `PhoneticMap` entry with fewer than 10 candidates: number keys for missing slots do nothing
- If query grows beyond any existing key (e.g. `"ss"`), candidate list is empty but query is still shown
- Switching from SearchMode to DirectMode (Tab) should clear `searchQuery`

---

## TUI Design Mockup

### SearchMode — query being built, candidates displayed

The search query and candidate strip live between the mode badge and the status
line. They replace (not supplement) the normal "nothing here" gap.

```
╭──────────────────────────────────────────────────────────────────────╮
│  สวัสดี                                                                │
╰──────────────────────────────────────────────────────────────────────╯

 ╔════════╗  Thai (Kedmanee)
 ║ SEARCH ║
 ╚════════╝

 Search: s▌
 1:ส  2:ษ  3:ศ  4:ซ

                           ← status line (blank)

  ╭───╮╭───╮╭───╮  ...keyboard rows...
```

---

### Candidate line rendering rules

| Scenario | Candidate line |
|----------|----------------|
| Query `"s"` → 4 matches | `1:ส  2:ษ  3:ศ  4:ซ` |
| Query `"t"` → 6 matches | `1:ต  2:ถ  3:ท  4:ธ  5:ฑ  6:ฒ` |
| Query `"a"` → 7 matches | `1:ะ  2:า  3:ำ  4:แ  5:โ  6:ใ  7:ไ` |
| Hypothetical 10 matches  | `1:x  2:x  3:x  4:x  5:x  6:x  7:x  8:x  9:x  0:x` |
| Query `"ss"` → no match  | *(blank candidate line — row still rendered)* |
| Query `""` (empty)       | *(blank candidate line)* |

**Format per candidate:** `N:char` where `N` is `#888888` and `char` is `#FFFFFF` bold.
Candidates separated by two spaces. Key `0` maps to the 10th candidate (index 9).

```
 Search: t▌
 1:ต  2:ถ  3:ท  4:ธ  5:ฑ  6:ฒ
```

---

### Backspace cascade behaviour (visual)

```
Buffer:  [ส][ว][ั][ส][ด][ี]     Query: "st"
                                     ↓ Backspace
Buffer:  [ส][ว][ั][ส][ด][ี]     Query: "s"
                                     ↓ Backspace
Buffer:  [ส][ว][ั][ส][ด][ี]     Query: ""
                                     ↓ Backspace
Buffer:  [ส][ว][ั][ส][ด]        Query: ""
```

---

### Spanish PhoneticMap SearchMode example

```
╭──────────────────────────────────────────────────────────────────────╮
│  Buenos días                                                          │
╰──────────────────────────────────────────────────────────────────────╯

 ╔════════╗  Spanish Standard
 ║ SEARCH ║
 ╚════════╝

 Search: u▌
 1:ú  2:ü

```
