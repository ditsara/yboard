package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ditsara/yboard/internal/types"
)

func setupTestModel() model {
	m := testModel()
	m.state = types.StateSetup
	return m
}

func TestViewSetup_HasTitle(t *testing.T) {
	m := setupTestModel()
	out := viewSetup(m)
	if !strings.Contains(out, "Language Modules") {
		t.Error("setup screen should contain 'Language Modules' title")
	}
}

func TestViewSetup_ShowsModuleNames(t *testing.T) {
	m := setupTestModel()
	out := viewSetup(m)
	if !strings.Contains(out, "Thai (Kedmanee)") {
		t.Error("setup screen should show Thai module name")
	}
	if !strings.Contains(out, "Spanish Standard") {
		t.Error("setup screen should show Spanish module name")
	}
}

func TestViewSetup_CursorOnFirstRow(t *testing.T) {
	m := setupTestModel()
	m.setupCursor = 0
	out := viewSetup(m)
	if !strings.Contains(out, "▶") {
		t.Error("cursor indicator ▶ should appear")
	}
}

func TestViewSetup_CheckboxEnabled(t *testing.T) {
	m := setupTestModel()
	out := viewSetup(m)
	if !strings.Contains(out, "[✓]") {
		t.Error("enabled modules should show [✓]")
	}
}

func TestViewSetup_CheckboxDisabled(t *testing.T) {
	m := setupTestModel()
	m.languages[0].Enabled = false
	out := viewSetup(m)
	if !strings.Contains(out, "[ ]") {
		t.Error("disabled module should show [ ]")
	}
}

func TestViewSetup_AllDisabledWarning(t *testing.T) {
	m := setupTestModel()
	for i := range m.languages {
		m.languages[i].Enabled = false
	}
	out := viewSetup(m)
	if !strings.Contains(out, "No languages enabled") {
		t.Error("should show warning when all languages disabled")
	}
}

func TestViewSetup_HasFooter(t *testing.T) {
	m := setupTestModel()
	out := viewSetup(m)
	if !strings.Contains(out, "Space: toggle") {
		t.Error("setup footer should contain 'Space: toggle'")
	}
}

func TestUpdateSetup_NavigateDown(t *testing.T) {
	m := setupTestModel()
	m.setupCursor = 0
	result, _ := updateSetup(m, tea.KeyMsg{Type: tea.KeyDown})
	rm := result.(model)
	if rm.setupCursor != 1 {
		t.Errorf("down arrow: want cursor=1, got %d", rm.setupCursor)
	}
}

func TestUpdateSetup_NavigateUp(t *testing.T) {
	m := setupTestModel()
	m.setupCursor = 1
	result, _ := updateSetup(m, tea.KeyMsg{Type: tea.KeyUp})
	rm := result.(model)
	if rm.setupCursor != 0 {
		t.Errorf("up arrow: want cursor=0, got %d", rm.setupCursor)
	}
}

func TestUpdateSetup_ClampAtTop(t *testing.T) {
	m := setupTestModel()
	m.setupCursor = 0
	result, _ := updateSetup(m, tea.KeyMsg{Type: tea.KeyUp})
	rm := result.(model)
	if rm.setupCursor != 0 {
		t.Errorf("cursor should clamp at 0, got %d", rm.setupCursor)
	}
}

func TestUpdateSetup_ClampAtBottom(t *testing.T) {
	m := setupTestModel()
	m.setupCursor = len(m.languages) - 1
	result, _ := updateSetup(m, tea.KeyMsg{Type: tea.KeyDown})
	rm := result.(model)
	if rm.setupCursor != len(m.languages)-1 {
		t.Errorf("cursor should clamp at bottom, got %d", rm.setupCursor)
	}
}

func TestUpdateSetup_ToggleWithSpace(t *testing.T) {
	m := setupTestModel()
	m.setupCursor = 0
	initialEnabled := m.languages[0].Enabled
	result, _ := updateSetup(m, tea.KeyMsg{Type: tea.KeySpace})
	rm := result.(model)
	if rm.languages[0].Enabled == initialEnabled {
		t.Error("space should toggle the enabled state")
	}
}

func TestUpdateSetup_ReturnWithEsc(t *testing.T) {
	m := setupTestModel()
	m.state = types.StateSetup
	result, _ := updateSetup(m, tea.KeyMsg{Type: tea.KeyEscape})
	rm := result.(model)
	if rm.state != types.StateTyping {
		t.Error("Esc should return to StateTyping")
	}
}

func TestUpdateSetup_ReturnWithF2(t *testing.T) {
	m := setupTestModel()
	result, _ := updateSetup(m, tea.KeyMsg{Type: tea.KeyF2})
	rm := result.(model)
	if rm.state != types.StateTyping {
		t.Error("F2 should return to StateTyping")
	}
}

func TestUpdateSetup_VimNavigation(t *testing.T) {
	m := setupTestModel()
	m.setupCursor = 0
	result, _ := updateSetup(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	rm := result.(model)
	if rm.setupCursor != 1 {
		t.Errorf("j should move cursor down: want 1, got %d", rm.setupCursor)
	}
	result2, _ := updateSetup(rm, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	rm2 := result2.(model)
	if rm2.setupCursor != 0 {
		t.Errorf("k should move cursor up: want 0, got %d", rm2.setupCursor)
	}
}
