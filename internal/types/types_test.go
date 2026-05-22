package types

import "testing"

func TestEmptyKey(t *testing.T) {
vk := EmptyKey("q")
if vk.Key != "q" {
t.Errorf("Key: want %q, got %q", "q", vk.Key)
}
if vk.Normal != "" || vk.Shift != "" {
t.Errorf("EmptyKey should have empty Normal and Shift, got Normal=%q Shift=%q", vk.Normal, vk.Shift)
}
}

func TestValidateModule_Valid(t *testing.T) {
m := LanguageModule{
ID: "test",
KeyboardRows: [][]VisualKey{
make([]VisualKey, 13),
make([]VisualKey, 13),
make([]VisualKey, 11),
make([]VisualKey, 10),
},
}
if err := ValidateModule(m); err != nil {
t.Errorf("expected nil, got: %v", err)
}
}

func TestValidateModule_WrongRowCount(t *testing.T) {
m := LanguageModule{ID: "bad", KeyboardRows: [][]VisualKey{make([]VisualKey, 13)}}
if err := ValidateModule(m); err == nil {
t.Error("expected error for wrong row count, got nil")
}
}

func TestValidateModule_WrongRowLength(t *testing.T) {
m := LanguageModule{
ID: "bad",
KeyboardRows: [][]VisualKey{
make([]VisualKey, 10), // wrong: should be 13
make([]VisualKey, 13),
make([]VisualKey, 11),
make([]VisualKey, 10),
},
}
if err := ValidateModule(m); err == nil {
t.Error("expected error for wrong row 0 length, got nil")
}
}

func TestBuildDirectMaps_AlphaAndSymbol(t *testing.T) {
rows := [][]VisualKey{
{{Key: "1", Normal: "!", Shift: "@"}, EmptyKey("2")},
{{Key: "a", Normal: "x", Shift: "Y"}},
}
direct, shift := BuildDirectMaps(rows)

if direct["1"] != "!" {
t.Errorf("direct[1]: want !, got %q", direct["1"])
}
if direct["2"] != "" {
t.Errorf("EmptyKey should not appear in direct map, got %q", direct["2"])
}
if shift["!"] != "@" { // USKeyShifted["1"] == "!"
t.Errorf("shift[!]: want @, got %q", shift["!"])
}
if direct["a"] != "x" {
t.Errorf("direct[a]: want x, got %q", direct["a"])
}
if shift["A"] != "Y" { // strings.ToUpper("a") == "A"
t.Errorf("shift[A]: want Y, got %q", shift["A"])
}
}

func TestBuildDirectMaps_EmptyShift(t *testing.T) {
rows := [][]VisualKey{{{Key: "q", Normal: "ๆ", Shift: ""}}}
direct, shift := BuildDirectMaps(rows)
if direct["q"] != "ๆ" {
t.Errorf("direct[q]: want ๆ, got %q", direct["q"])
}
if _, ok := shift["Q"]; ok {
t.Error("shift[Q] should not exist when Shift is empty")
}
}
