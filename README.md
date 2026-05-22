# yboard

A lightweight TUI tool for typing in foreign languages using a standard **US
keyboard** — no OS-level layout switching required. Works in WSL, Linux, and
native Windows terminals.

## Why

Installing OS keyboard layouts causes accidental language switches during
normal work. yboard runs entirely in the terminal, intercepts your keystrokes,
and pipes the result to the clipboard — your system keyboard stays locked to
English at all times.

## Usage

```
./yboard
```

### Keybindings

| Key           | Action                                       |
|---------------|----------------------------------------------|
| **Tab**       | Toggle between Direct and Search mode        |
| **Enter**     | Copy buffer to clipboard and clear it        |
| **F9**        | Copy buffer to clipboard (keep buffer)       |
| **F3 / F4**   | Switch to previous / next enabled language   |
| **F2**        | Open setup screen (enable/disable languages) |
| **Backspace** | Delete last char in buffer (or search bar)   |
| **Space**     | Append a space to the buffer                 |
| **Ctrl+L**    | Force redraw                                 |
| **F10**       | Quit                                         |

## Input Modes

### Direct Mode
Each key maps 1-to-1 to a target character based on the active language layout.
Hold **Shift** for the shifted variant. The on-screen keyboard grid shows both
Normal (bright) and Shift (dimmed) characters for every key.

If a key has no mapping, the literal English character is appended and a
warning appears in the status bar.

### Search Mode (Phonetic)
Type a phonetic string (e.g. `s`, `ch`, `ng`, `ai`) to search the active
language's phonetic map. Matches are shown numbered `1`–`9` (and `0` for the
10th). Press the corresponding number to pick a character, which appends it to
the buffer and clears the search query.

Shift is ignored in Search mode — input is always lowercased.

## Languages

### Thai (Kedmanee)
Uses the standard Thai Kedmanee layout mapped to a US keyboard. Direct mode
follows the Kedmanee assignments. Phonetic search covers:

| Query | Characters       |
|-------|------------------|
| `s`   | ส ษ ศ ซ          |
| `t`   | ต ถ ท ธ ฑ ฒ ฏ ฐ  |
| `k`   | ก ข ค ฆ ฃ ฅ      |
| `p`   | ป ผ พ ภ ฝ ฟ      |
| `b`   | บ                |
| `d`   | ด ฎ              |
| `f`   | ฝ ฟ              |
| `r`   | ร                |
| `l`   | ล ฬ              |
| `w`   | ว                |
| `y`   | ย ญ              |
| `ch`  | จ ฉ ช ฌ          |
| `ng`  | ง                |
| `m`   | ม                |
| `n`   | น ณ              |
| `h`   | ห ฮ              |
| `a`   | ะ า ำ แ โ ใ ไ    |
| `ae`  | แ                |
| `ai`  | ใ ไ              |
| `e`   | เ                |
| `i`   | ิ ี ึ ื              |
| `o`   | โ อ              |
| `u`   | ุ ู                |
| `ue`  | ึ ื                |
| `um`  | ำ                |

### Spanish (Standard)
Direct mode types regular Latin characters with `;` → `ñ`. Phonetic search covers accented vowels and special characters:

| Query | Characters |
|-------|------------|
| `a`   | á ä        |
| `e`   | é          |
| `i`   | í          |
| `o`   | ó          |
| `u`   | ú ü        |
| `n`   | ñ          |
| `c`   | ç          |
| `h`   | ¡ ¿        |

## Clipboard

On **Enter** or **F9**, yboard pipes the buffer to the clipboard using the first available tool:
1. `clip.exe` — WSL / Windows
2. `wl-copy` — Wayland Linux
3. `xclip` — X11 Linux

## Build

```
make build   # outputs ./yboard
make test    # runs tests
```
