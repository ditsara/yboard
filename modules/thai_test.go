package modules

import (
	"testing"

	"github.com/ditsara/yboard/internal/types"
)

func TestThaiModule_ValidateRows(t *testing.T) {
	if err := types.ValidateModule(ThaiModule); err != nil {
		t.Fatalf("ThaiModule failed validation: %v", err)
	}
}

func TestThaiModule_KeyLabels(t *testing.T) {
	// Spot check: every key in every row must have a non-empty Key field
	for r, row := range ThaiModule.KeyboardRows {
		for c, key := range row {
			if key.Key == "" {
				t.Errorf("row %d col %d: missing Key label (Normal=%q Shift=%q)", r, c, key.Normal, key.Shift)
			}
		}
	}
}

func TestThaiModule_Row0Has13Keys(t *testing.T) {
	if got := len(ThaiModule.KeyboardRows[0]); got != 13 {
		t.Errorf("row 0: want 13 keys, got %d", got)
	}
}

func TestThaiModule_EqualKeyEmpty(t *testing.T) {
	equalKey := ThaiModule.KeyboardRows[0][12]
	if equalKey.Key != "=" {
		t.Errorf("row 0 position 12: want Key=%q, got %q", "=", equalKey.Key)
	}
	if equalKey.Normal != "" || equalKey.Shift != "" {
		t.Errorf("= key should be empty (not used in Kedmanee), got Normal=%q Shift=%q", equalKey.Normal, equalKey.Shift)
	}
}
