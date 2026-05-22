package modules

import "github.com/ditsara/yboard/internal/types"

// ThaiModule is the Thai Kedmanee keyboard layout.
// Row 0 has 12 remapped keys (backtick through -); the = key is not used in Kedmanee,
// so it is represented as an EmptyKey at position 12 to satisfy StandardRowLengths.
var ThaiModule = types.LanguageModule{
	ID: "thai", Name: "Thai (Kedmanee)", Enabled: true,
	DirectMap: map[string]string{
		"q": "ๆ", "w": "ไ", "e": "ำ", "r": "พ", "t": "ะ", "y": "ั", "u": "ี", "i": "ร", "o": "น", "p": "ย",
		"a": "ฟ", "s": "ห", "d": "ก", "f": "ด", "g": "เ", "h": "้", "j": "่", "k": "า", "l": "ส", ";": "ว", "'": "ง",
		"z": "ผ", "x": "ป", "c": "แ", "v": "อ", "b": "ิ", "n": "ื", "m": "ท",
	},
	ShiftDirectMap: map[string]string{
		"Q": "๐", "W": "\"", "E": "ฎ", "R": "ฑ", "T": "ธ", "Y": "ํ", "U": "๊", "I": "ณ", "O": "ฯ", "P": "ญ",
		"A": "ฤ", "S": "ฆ", "D": "ฏ", "F": "โ", "G": "ฌ", "H": "็", "J": "๋", "K": "ษ", "L": "ศ", ":": "ซ", "\"": ".",
		"Z": "(", "X": ")", "C": "ฉ", "V": "ฮ", "B": "ฺ", "N": "์", "M": "?",
	},
	PhoneticMap: map[string][]string{
		"s": {"ส", "ษ", "ศ", "ซ"}, "t": {"ต", "ถ", "ท", "ธ", "ฑ", "ฒ"},
		"k": {"ก", "ข", "ค", "ฆ"}, "p": {"ป", "ผ", "พ", "ภ", "ฝ", "ฟ"},
		"a": {"ะ", "า", "ำ", "แ", "โ", "ใ", "ไ"}, "i": {"ิ", "ี", "ึ", "ื"}, "u": {"ุ", "ู"},
	},
	KeyboardRows: [][]types.VisualKey{
		// Number row: backtick through - remapped; = not used in Kedmanee
		{
			{Key: "`", Normal: "ๅ", Shift: "+"},
			{Key: "1", Normal: "/", Shift: "๑"},
			{Key: "2", Normal: "_", Shift: "๒"},
			{Key: "3", Normal: "ภ", Shift: "๓"},
			{Key: "4", Normal: "ถ", Shift: "๔"},
			{Key: "5", Normal: "ุ", Shift: "ู"},
			{Key: "6", Normal: "ึ", Shift: "฿"},
			{Key: "7", Normal: "ค", Shift: "๕"},
			{Key: "8", Normal: "ต", Shift: "๖"},
			{Key: "9", Normal: "จ", Shift: "๗"},
			{Key: "0", Normal: "ข", Shift: "๘"},
			{Key: "-", Normal: "ช", Shift: "๙"},
			types.EmptyKey("="),
		},
		// QWERTY row: all 13 positions remapped
		{
			{Key: "q", Normal: "ๆ", Shift: "๐"},
			{Key: "w", Normal: "ไ", Shift: "\""},
			{Key: "e", Normal: "ำ", Shift: "ฎ"},
			{Key: "r", Normal: "พ", Shift: "ฑ"},
			{Key: "t", Normal: "ะ", Shift: "ธ"},
			{Key: "y", Normal: "ั", Shift: "ํ"},
			{Key: "u", Normal: "ี", Shift: "๊"},
			{Key: "i", Normal: "ร", Shift: "ณ"},
			{Key: "o", Normal: "น", Shift: "ฯ"},
			{Key: "p", Normal: "ย", Shift: "ญ"},
			{Key: "[", Normal: "บ", Shift: "ฐ"},
			{Key: "]", Normal: "ล", Shift: ","},
			{Key: `\`, Normal: "ฃ", Shift: "ฅ"},
		},
		// Home row: all 11 positions remapped
		{
			{Key: "a", Normal: "ฟ", Shift: "ฤ"},
			{Key: "s", Normal: "ห", Shift: "ฆ"},
			{Key: "d", Normal: "ก", Shift: "ฏ"},
			{Key: "f", Normal: "ด", Shift: "โ"},
			{Key: "g", Normal: "เ", Shift: "ฌ"},
			{Key: "h", Normal: "้", Shift: "็"},
			{Key: "j", Normal: "่", Shift: "๋"},
			{Key: "k", Normal: "า", Shift: "ษ"},
			{Key: "l", Normal: "ส", Shift: "ศ"},
			{Key: ";", Normal: "ว", Shift: "ซ"},
			{Key: "'", Normal: "ง", Shift: "."},
		},
		// Bottom row: all 10 positions remapped
		{
			{Key: "z", Normal: "ผ", Shift: "("},
			{Key: "x", Normal: "ป", Shift: ")"},
			{Key: "c", Normal: "แ", Shift: "ฉ"},
			{Key: "v", Normal: "อ", Shift: "ฮ"},
			{Key: "b", Normal: "ิ", Shift: "ฺ"},
			{Key: "n", Normal: "ื", Shift: "์"},
			{Key: "m", Normal: "ท", Shift: "?"},
			{Key: ",", Normal: "ม", Shift: "ฒ"},
			{Key: ".", Normal: "ใ", Shift: "ฬ"},
			{Key: "/", Normal: "ฝ", Shift: "฾"},
		},
	},
}
