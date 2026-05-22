# YBoard: A Multi-Language TUI Keyboard Engine

## Background

I need a lightweight, non-intrusive way to occasionally type in foreign
languages—like Thai and Spanish—on my Windows work laptop without installing
OS-level keyboard layouts that constantly trigger accidental, frustrating
language switches during my normal work. To solve this, I am building an
isolated, multi-language Terminal User Interface (TUI) tool in Go (running via
WSL) that intercepts standard English keystrokes to construct foreign text. By
providing dual input modes (direct transliteration and phonetic search), a
unified onscreen visual reference for shifted and unshifted keys, and a direct
pipe to the Windows clipboard, this tool eliminates OS-level layout friction
while keeping my primary system keyboard locked safely to English.

## 1. Project Overview & Tech Stack

Build a single-binary, highly portable Terminal User Interface (TUI) application in Go that allows users to type in alternative languages (e.g., Thai, Spanish) using a standard English keyboard.

* **Core Frameworks:** [github.com/charmbracelet/bubbletea]() (State/Event loop), [github.com/charmbracelet/lipgloss]() (UI/Styling).
* **Target Environments:** WSL, Linux, native Windows.
* **Core Philosophy:** Intercept physical keystrokes, translate them in-memory based on active layout maps, visually render the layout, and pipe the accumulated text buffer to the system clipboard.

---

## 2. Architecture & Data Structures

The system handles both layout mapping and visual rendering using a unified VisualKey structure so that both Normal and Shift states are displayed simultaneously on screen.

```go
// Represents a single physical key's output variants for rendering
type VisualKey struct {
	Normal string
	Shift  string
}

// Represents a complete language package
type LanguageModule struct {
	ID             string
	Name           string
	Enabled        bool                // Controlled by the Setup Screen
	DirectMap      map[string]string   // Physical key (lower or symbol) -> Target char; covers a-z and number/symbol row
	ShiftDirectMap map[string]string   // Physical key (shifted) -> Target char
	PhoneticMap    map[string][]string // Search string -> Array of target chars (max 10 shown; keys 1-9 and 0 select)
	KeyboardRows   [][]VisualKey       // Physical rows for visual rendering
}

// Application State Enums
type AppState int
const (
	StateTyping AppState = iota
	StateSetup
)

type InputMode int
const (
	DirectMode InputMode = iota // Direct 1-to-1 transliteration
	SearchMode                  // Phonetic search and select
)

// The word buffer is a shared []rune that persists across language switches,
// allowing mixed-language text accumulation.

```

---

## 3. Keybindings & Interactions

### Global Controls (Work in all states)

* **F10**: Graceful Shutdown (safely exit the alternate screen buffer).
* **Ctrl+L**: Force Screen Redraw (tea.ClearScreen).

### Typing View Controls (StateTyping)

* **F2**: Open Setup Screen (StateSetup).
* **F3**: Select **Previous** *Enabled* Language (Wrap around).
* **F4**: Select **Next** *Enabled* Language (Wrap around).
* **F9**: Export Word Buffer to Clipboard.
* **Enter**: Copy Word Buffer to Clipboard and clear the buffer (shorthand for F9 + clear).
* **Tab**: Toggle between DirectMode and SearchMode.
* **Spacebar**: Append literal space (" ") to buffer.
* **Backspace**:
* If in SearchMode with active search query, delete last query char.
* Otherwise, pop the last UTF-8 rune from the word buffer.



### Input Modes (Typing View)

* **DirectMode**: Alphabetical and symbol/number keys are intercepted. Capitalized inputs (via physical Shift) query the ShiftDirectMap. Lowercase inputs query the DirectMap. Matches append instantly to the buffer. If a key has no entry in DirectMap/ShiftDirectMap, the literal English character is appended and a warning is shown in the status area.
* **SearchMode**: Keys a-z (Shift is ignored; input is normalized to lowercase) build a searchQuery. Matches from the PhoneticMap are displayed with numerical indexes (1:ส 2:ษ 3:ศ). Number keys 1-9 and 0 (representing 10th match) select the match, append it to the buffer, and clear the searchQuery. If a PhoneticMap entry has more than 10 candidates, only the first 10 are shown; the rest are silently ignored.

### Setup View Controls (StateSetup)

* **Up/Down or k/j**: Navigate list of language modules.
* **Spacebar**: Toggle Enabled boolean on the highlighted module.
* **F2 or Enter or Esc**: Return to StateTyping.

---

## 4. UI Rendering Requirements (Lipgloss)

1. **Alternate Screen:** Must use tea.WithAltScreen() to prevent terminal scrollback clutter.
2. **Word Buffer:** Displayed prominently in a bordered box.
3. **Mode Badges:** Clearly display whether the active mode is Direct or Search.
4. **Status Message Area:** A single-line area below the buffer box displaying transient messages (e.g., "📋 Copied!", "⚠ Unknown key: '[' — passed through", "⚠ No languages enabled"). Messages are cleared on the next keypress.
5. **Visual Keyboard Grid:**
* Iterate through ActiveModule.KeyboardRows.
* Each VisualKey is rendered using lipgloss.JoinVertical to stack the Shift character above the Normal character inside a single rounded border box.
* The Shift character should be slightly dimmed (e.g., #888888), and the Normal character bold and bright (#FFFFFF).

6. **Setup Screen:** A simple list showing checkboxes. [X] Thai (Kedmanee), [ ] Greek, [X] Spanish. If all modules are disabled, the Typing view shows a "⚠ No languages enabled" warning and blocks all key input until at least one module is re-enabled.
7. **Hotkeys Footer:** Always visible list of shortcuts so the user never has to memorize them.

---

## 5. Clipboard Integration

Triggered by F9 or Enter. Checks for available clipboard tools in the following priority order: `clip.exe` (WSL/Windows), `wl-copy` (Wayland Linux), `xclip` (X11 Linux).

```go
func executeClipboardCopy(buffer string) string {
	if len(buffer) == 0 { return "Buffer empty" }
	
	var cmd *exec.Cmd
	if _, err := exec.LookPath("clip.exe"); err == nil {
		cmd = exec.Command("clip.exe")
	} else if _, err := exec.LookPath("wl-copy"); err == nil {
		cmd = exec.Command("wl-copy")
	} else {
		cmd = exec.Command("xclip", "-selection", "clipboard")
	}

	cmd.Stdin = strings.NewReader(buffer)
	if err := cmd.Run(); err == nil {
		return "📋 Copied text block to clipboard!"
	}
	return "❌ Clipboard process pipeline error"
}

```

---

## 6. Language Module Seed Data

### Module 1: Thai (Kedmanee)

```go
var ThaiModule = LanguageModule{
	ID: "thai", Name: "Thai (Kedmanee)", Enabled: true,
	DirectMap: map[string]string{
		"q": "ๆ", "w": "ไ", "e": "ำ", "r": "พ", "t": "ะ", "y": "ั", "u": "ี", "i": "ร", "o": "น", "p": "ย",
		"a": "ฟ", "s": "ห", "d": "ก", "f": "ด", "g": "เ", "h": "้", "j": "่", "k": "า", "l": "ส",
		"z": "ผ", "x": "ป", "c": "แ", "v": "อ", "b": "ิ", "n": "ื", "m": "ท",
	},
	ShiftDirectMap: map[string]string{
		"Q": "๐", "W": "\"", "E": "ฎ", "R": "ฑ", "T": "ธ", "Y": "ํ", "U": "๊", "I": "ณ", "O": "ฯ", "P": "ญ",
		"A": "ฤ", "S": "ฆ", "D": "ฏ", "F": "โ", "G": "ฌ", "H": "็", "J": "๋", "K": "ษ", "L": "ศ",
		"Z": "(", "X": ")", "C": "ฉ", "V": "ฮ", "B": "ฺ", "N": "์", "M": "?",
	},
	PhoneticMap: map[string][]string{
		"s": {"ส", "ษ", "ศ", "ซ"}, "t": {"ต", "ถ", "ท", "ธ", "ฑ", "ฒ"},
		"k": {"ก", "ข", "ค", "ฆ"}, "p": {"ป", "ผ", "พ", "ภ", "ฝ", "ฟ"},
		"a": {"ะ", "า", "ำ", "แ", "โ", "ใ", "ไ"}, "i": {"ิ", "ี", "ึ", "ื"}, "u": {"ุ", "ู"},
	},
	KeyboardRows: [][]VisualKey{
		{ {Normal: "ๅ", Shift: "+"}, {Normal: "/", Shift: "๑"}, {Normal: "_", Shift: "๒"}, {Normal: "ภ", Shift: "๓"}, {Normal: "ถ", Shift: "๔"}, {Normal: "ุ", Shift: "ู"}, {Normal: "ึ", Shift: "฿"}, {Normal: "ค", Shift: "๕"}, {Normal: "ต", Shift: "๖"}, {Normal: "จ", Shift: "๗"}, {Normal: "ข", Shift: "๘"}, {Normal: "ช", Shift: "๙"} },
		{ {Normal: "ๆ", Shift: "๐"}, {Normal: "ไ", Shift: "\""}, {Normal: "ำ", Shift: "ฎ"}, {Normal: "พ", Shift: "ฑ"}, {Normal: "ะ", Shift: "ธ"}, {Normal: "ั", Shift: "ํ"}, {Normal: "ี", Shift: "๊"}, {Normal: "ร", Shift: "ณ"}, {Normal: "น", Shift: "ฯ"}, {Normal: "ย", Shift: "ญ"}, {Normal: "บ", Shift: "ฐ"}, {Normal: "ล", Shift: ","}, {Normal: "ฃ", Shift: "ฅ"} },
		{ {Normal: "ฟ", Shift: "ฤ"}, {Normal: "ห", Shift: "ฆ"}, {Normal: "ก", Shift: "ฏ"}, {Normal: "ด", Shift: "โ"}, {Normal: "เ", Shift: "ฌ"}, {Normal: "้", Shift: "็"}, {Normal: "่", Shift: "๋"}, {Normal: "า", Shift: "ษ"}, {Normal: "ส", Shift: "ศ"}, {Normal: "ว", Shift: "ซ"}, {Normal: "ง", Shift: "."} },
		{ {Normal: "ผ", Shift: "("}, {Normal: "ป", Shift: ")"}, {Normal: "แ", Shift: "ฉ"}, {Normal: "อ", Shift: "ฮ"}, {Normal: "ิ", Shift: "ฺ"}, {Normal: "ื", Shift: "์"}, {Normal: "ท", Shift: "?"}, {Normal: "ม", Shift: "ฒ"}, {Normal: "ใ", Shift: "ฬ"}, {Normal: "ฝ", Shift: "฾"} },
	},
}

```

### Module 2: Spanish (Standard ISO)

```go
var SpanishModule = LanguageModule{
	ID: "spanish", Name: "Spanish Standard", Enabled: true,
	DirectMap: map[string]string{
		"q": "q", "w": "w", "e": "e", "r": "r", "t": "t", "y": "y", "u": "u", "i": "i", "o": "o", "p": "p",
		"a": "a", "s": "s", "d": "d", "f": "f", "g": "g", "h": "h", "j": "j", "k": "k", "l": "l", ";": "ñ",
		"z": "z", "x": "x", "c": "c", "v": "v", "b": "b", "n": "n", "m": "m", ",": ",", ".": ".", "/": "-",
	},
	ShiftDirectMap: map[string]string{
		"Q": "Q", "W": "W", "E": "E", "R": "R", "T": "T", "Y": "Y", "U": "U", "I": "I", "O": "O", "P": "P",
		"A": "A", "S": "S", "D": "D", "F": "F", "G": "G", "H": "H", "J": "J", "K": "K", "L": "L", ":": "Ñ",
		"Z": "Z", "X": "X", "C": "C", "V": "V", "B": "B", "N": "N", "M": "M", "<": ";", ">": ":", "?": "_",
	},
	PhoneticMap: map[string][]string{
		"a": {"á", "ä"}, "e": {"é"}, "i": {"í"}, "o": {"ó"},
		"u": {"ú", "ü"}, "n": {"ñ"}, "c": {"ç"}, "h": {"¡", "¿"},
	},
	KeyboardRows: [][]VisualKey{
		{ {Normal: "1", Shift: "!"}, {Normal: "2", Shift: "@"}, {Normal: "3", Shift: "#"}, {Normal: "4", Shift: "$"}, {Normal: "5", Shift: "%"}, {Normal: "6", Shift: "^"}, {Normal: "7", Shift: "&"}, {Normal: "8", Shift: "*"}, {Normal: "9", Shift: "("}, {Normal: "0", Shift: ")"} },
		{ {Normal: "q", Shift: "Q"}, {Normal: "w", Shift: "W"}, {Normal: "e", Shift: "E"}, {Normal: "r", Shift: "R"}, {Normal: "t", Shift: "T"}, {Normal: "y", Shift: "Y"}, {Normal: "u", Shift: "U"}, {Normal: "i", Shift: "I"}, {Normal: "o", Shift: "O"}, {Normal: "p", Shift: "P"} },
		{ {Normal: "a", Shift: "A"}, {Normal: "s", Shift: "S"}, {Normal: "d", Shift: "D"}, {Normal: "f", Shift: "F"}, {Normal: "g", Shift: "G"}, {Normal: "h", Shift: "H"}, {Normal: "j", Shift: "J"}, {Normal: "k", Shift: "K"}, {Normal: "l", Shift: "L"}, {Normal: "ñ", Shift: "Ñ"} },
		{ {Normal: "z", Shift: "Z"}, {Normal: "x", Shift: "X"}, {Normal: "c", Shift: "C"}, {Normal: "v", Shift: "V"}, {Normal: "b", Shift: "B"}, {Normal: "n", Shift: "N"}, {Normal: "m", Shift: "M"}, {Normal: "ç", Shift: "Ç"} },
	},
}

```
