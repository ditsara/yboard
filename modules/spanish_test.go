package modules

import (
	"testing"

	"github.com/ditsara/yboard/internal/types"
)

func TestSpanishModule_ValidateRows(t *testing.T) {
	if err := types.ValidateModule(SpanishModule); err != nil {
		t.Fatalf("SpanishModule failed validation: %v", err)
	}
}

func TestSpanishModule_KeyLabels(t *testing.T) {
	for r, row := range SpanishModule.KeyboardRows {
		for c, key := range row {
			if key.Key == "" {
				t.Errorf("row %d col %d: missing Key label (Normal=%q Shift=%q)", r, c, key.Normal, key.Shift)
			}
		}
	}
}

func TestSpanishModule_EmptyKeys(t *testing.T) {
	empties := map[[2]int]string{
		{0, 0}: "`", {0, 11}: "-", {0, 12}: "=",
		{1, 10}: "[", {1, 11}: "]", {1, 12}: `\`,
		{2, 10}: "'",
	}
	for pos, wantKey := range empties {
		r, c := pos[0], pos[1]
		key := SpanishModule.KeyboardRows[r][c]
		if key.Key != wantKey {
			t.Errorf("row %d col %d: want Key=%q, got %q", r, c, wantKey, key.Key)
		}
		if key.Normal != "" || key.Shift != "" {
			t.Errorf("row %d col %d (%q): expected empty key, got Normal=%q Shift=%q", r, c, wantKey, key.Normal, key.Shift)
		}
	}
}

func TestSpanishModule_NKey(t *testing.T) {
	// ñ should be on the ; key (row 2, position 9)
	nKey := SpanishModule.KeyboardRows[2][9]
	if nKey.Key != ";" {
		t.Errorf("ñ key: want Key=\";\", got %q", nKey.Key)
	}
	if nKey.Normal != "ñ" || nKey.Shift != "Ñ" {
		t.Errorf("ñ key: want Normal=ñ Shift=Ñ, got Normal=%q Shift=%q", nKey.Normal, nKey.Shift)
	}
}
