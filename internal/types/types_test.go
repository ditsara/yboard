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
