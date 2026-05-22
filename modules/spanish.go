package modules

import "github.com/ditsara/yboard/internal/types"

// SpanishModule is the Spanish Standard ISO keyboard layout mapped onto a US ANSI keyboard.
// Rows 0-1 use EmptyKey for keys not present in the Spanish standard layout (backtick, -, =, [, ], \).
// Row 2 uses EmptyKey for the apostrophe key.
var SpanishModule = types.LanguageModule{
	ID: "spanish", Name: "Spanish Standard", Enabled: true,
	PhoneticMap: map[string][]string{
		"a": {"á", "ä"}, "e": {"é"}, "i": {"í"}, "o": {"ó"},
		"u": {"ú", "ü"}, "n": {"ñ"}, "c": {"ç"}, "h": {"¡", "¿"},
	},
	KeyboardRows: [][]types.VisualKey{
		// Number row: 1-0 remapped; backtick, -, = not present in Spanish standard
		{
			types.EmptyKey("`"),
			{Key: "1", Normal: "1", Shift: "!"},
			{Key: "2", Normal: "2", Shift: "@"},
			{Key: "3", Normal: "3", Shift: "#"},
			{Key: "4", Normal: "4", Shift: "$"},
			{Key: "5", Normal: "5", Shift: "%"},
			{Key: "6", Normal: "6", Shift: "^"},
			{Key: "7", Normal: "7", Shift: "&"},
			{Key: "8", Normal: "8", Shift: "*"},
			{Key: "9", Normal: "9", Shift: "("},
			{Key: "0", Normal: "0", Shift: ")"},
			types.EmptyKey("-"),
			types.EmptyKey("="),
		},
		// QWERTY row: q-p remapped; [, ], \ not present in Spanish standard
		{
			{Key: "q", Normal: "q", Shift: "Q"},
			{Key: "w", Normal: "w", Shift: "W"},
			{Key: "e", Normal: "e", Shift: "E"},
			{Key: "r", Normal: "r", Shift: "R"},
			{Key: "t", Normal: "t", Shift: "T"},
			{Key: "y", Normal: "y", Shift: "Y"},
			{Key: "u", Normal: "u", Shift: "U"},
			{Key: "i", Normal: "i", Shift: "I"},
			{Key: "o", Normal: "o", Shift: "O"},
			{Key: "p", Normal: "p", Shift: "P"},
			types.EmptyKey("["),
			types.EmptyKey("]"),
			types.EmptyKey(`\`),
		},
		// Home row: a-l + ñ remapped; ' not present in Spanish standard
		{
			{Key: "a", Normal: "a", Shift: "A"},
			{Key: "s", Normal: "s", Shift: "S"},
			{Key: "d", Normal: "d", Shift: "D"},
			{Key: "f", Normal: "f", Shift: "F"},
			{Key: "g", Normal: "g", Shift: "G"},
			{Key: "h", Normal: "h", Shift: "H"},
			{Key: "j", Normal: "j", Shift: "J"},
			{Key: "k", Normal: "k", Shift: "K"},
			{Key: "l", Normal: "l", Shift: "L"},
			{Key: ";", Normal: "ñ", Shift: "Ñ"},
			types.EmptyKey("'"),
		},
		// Bottom row: z-m + ç/punct; . and / derived from DirectMap/ShiftDirectMap
		{
			{Key: "z", Normal: "z", Shift: "Z"},
			{Key: "x", Normal: "x", Shift: "X"},
			{Key: "c", Normal: "c", Shift: "C"},
			{Key: "v", Normal: "v", Shift: "V"},
			{Key: "b", Normal: "b", Shift: "B"},
			{Key: "n", Normal: "n", Shift: "N"},
			{Key: "m", Normal: "m", Shift: "M"},
			{Key: ",", Normal: "ç", Shift: "Ç"},
			{Key: ".", Normal: ".", Shift: ":"},
			{Key: "/", Normal: "-", Shift: "_"},
		},
	},
}
