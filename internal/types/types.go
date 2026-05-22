package types

import (
	"fmt"
	"strings"
)

// VisualKey represents a single physical key's output variants for rendering.
// Key is the physical English label (e.g. "q", "1", ";") shown in the key box border.
// Normal and Shift are the language output characters; empty string means not remapped.
type VisualKey struct {
	Key    string
	Normal string
	Shift  string
}

// LanguageModule represents a complete language package.
// KeyboardRows must match StandardRowLengths exactly; use EmptyKey() for unremapped positions.
// DirectMap and ShiftDirectMap are derived at startup via BuildDirectMaps — do not set manually.
type LanguageModule struct {
	ID             string
	Name           string
	Enabled        bool
	DirectMap      map[string]string
	ShiftDirectMap map[string]string
	PhoneticMap    map[string][]string
	KeyboardRows   [][]VisualKey
}

// AppState represents which screen is currently active.
type AppState int

const (
	StateTyping AppState = iota
	StateSetup
)

// InputMode represents the active typing mode.
type InputMode int

const (
	DirectMode InputMode = iota
	SearchMode
)

// StandardRowLengths defines the required number of keys per row for all language modules.
// Rows correspond to: number row, QWERTY row, home row, bottom row.
// Physical key positions: [` 1 2 3 4 5 6 7 8 9 0 - =] [q…\] [a…'] [z…/]
var StandardRowLengths = []int{13, 13, 11, 10}

// USKeyShifted maps each base symbol key to the string BubbleTea sends when Shift is held.
// Alpha keys are not listed here — use strings.ToUpper for those.
var USKeyShifted = map[string]string{
	"`": "~",
	"1": "!", "2": "@", "3": "#", "4": "$", "5": "%",
	"6": "^", "7": "&", "8": "*", "9": "(", "0": ")",
	"-": "_", "=": "+",
	"[": "{", "]": "}", `\`: "|",
	";": ":", "'": `"`,
	",": "<", ".": ">", "/": "?",
}

// shiftedKeyString returns the key string BubbleTea sends for Shift+baseKey.
func shiftedKeyString(baseKey string) string {
	if len(baseKey) == 1 && baseKey[0] >= 'a' && baseKey[0] <= 'z' {
		return strings.ToUpper(baseKey)
	}
	if s, ok := USKeyShifted[baseKey]; ok {
		return s
	}
	return ""
}

// BuildDirectMaps derives DirectMap and ShiftDirectMap from a module's KeyboardRows.
// EmptyKeys (Normal == "") are skipped. Keys with no known shifted representation are
// added to DirectMap only.
func BuildDirectMaps(rows [][]VisualKey) (direct, shift map[string]string) {
	direct = make(map[string]string)
	shift = make(map[string]string)
	for _, row := range rows {
		for _, key := range row {
			if key.Normal == "" {
				continue // EmptyKey
			}
			direct[key.Key] = key.Normal
			if key.Shift != "" {
				if sk := shiftedKeyString(key.Key); sk != "" {
					shift[sk] = key.Shift
				}
			}
		}
	}
	return direct, shift
}

// EmptyKey creates a VisualKey for a physical key this language does not remap.
// The renderer draws it as a muted/dim box — only the key label is visible.
func EmptyKey(physicalKey string) VisualKey {
	return VisualKey{Key: physicalKey}
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
