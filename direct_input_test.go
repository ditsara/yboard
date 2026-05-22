package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ditsara/yboard/internal/types"
)

func makeKey(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func TestHandleDirectInput_ThaiLookup(t *testing.T) {
	m := testModel()
	m.inputMode = types.DirectMode
	m.activeIndex = 0 // Thai
	m = handleDirectInput(m, makeKey("d"))
	// Thai DirectMap "d" → "ก"
	if string(m.wordBuffer) != "ก" {
		t.Errorf("DirectMap 'd': want ก, got %q", string(m.wordBuffer))
	}
}

func TestHandleDirectInput_ThaiShiftLookup(t *testing.T) {
	m := testModel()
	m.inputMode = types.DirectMode
	m.activeIndex = 0 // Thai
	m = handleDirectInput(m, makeKey("D"))
	// Thai ShiftDirectMap "D" → "ฏ"
	if string(m.wordBuffer) != "ฏ" {
		t.Errorf("ShiftDirectMap 'D': want ฏ, got %q", string(m.wordBuffer))
	}
}

func TestHandleDirectInput_SpanishNTilde(t *testing.T) {
	m := testModel()
	m.inputMode = types.DirectMode
	m.activeIndex = 1 // Spanish
	m = handleDirectInput(m, makeKey(";"))
	// Spanish DirectMap ";" → "ñ"
	if string(m.wordBuffer) != "ñ" {
		t.Errorf("DirectMap ';': want ñ, got %q", string(m.wordBuffer))
	}
}

func TestHandleDirectInput_UnknownKey(t *testing.T) {
	m := testModel()
	m.inputMode = types.DirectMode
	m.activeIndex = 0
	m = handleDirectInput(m, makeKey("0")) // "0" not in Thai DirectMap
	if !strings.Contains(m.statusMessage, "Unknown key") {
		t.Errorf("unknown key should set warning status, got: %q", m.statusMessage)
	}
	if string(m.wordBuffer) != "0" {
		t.Errorf("unknown key should pass through as literal, got: %q", string(m.wordBuffer))
	}
}

func TestHandleDirectInput_NoLanguages(t *testing.T) {
	m := testModel()
	for i := range m.languages {
		m.languages[i].Enabled = false
	}
	m = handleDirectInput(m, makeKey("a"))
	if !strings.Contains(m.statusMessage, "No languages enabled") {
		t.Errorf("should show no-languages warning, got: %q", m.statusMessage)
	}
}

func TestUpdateTyping_Backspace_RemovesRune(t *testing.T) {
	m := testModel()
	m.wordBuffer = []rune("สวัสดี")
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyBackspace})
	rm := result.(model)
	if string(rm.wordBuffer) != "สวัสด" {
		t.Errorf("backspace should remove last rune, got: %q", string(rm.wordBuffer))
	}
}

func TestUpdateTyping_Backspace_EmptyBuffer(t *testing.T) {
	m := testModel()
	m.wordBuffer = nil
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyBackspace})
	rm := result.(model)
	if len(rm.wordBuffer) != 0 {
		t.Error("backspace on empty buffer should not crash")
	}
}

func TestUpdateTyping_Space_AppendsSpace(t *testing.T) {
	m := testModel()
	m.wordBuffer = []rune("ก")
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeySpace})
	rm := result.(model)
	if string(rm.wordBuffer) != "ก " {
		t.Errorf("space should append ' ', got: %q", string(rm.wordBuffer))
	}
}

func TestUpdateTyping_Tab_TogglesMode(t *testing.T) {
	m := testModel()
	m.inputMode = types.DirectMode
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyTab})
	rm := result.(model)
	if rm.inputMode != types.SearchMode {
		t.Error("Tab should switch DirectMode → SearchMode")
	}
	result2, _ := updateTyping(rm, tea.KeyMsg{Type: tea.KeyTab})
	rm2 := result2.(model)
	if rm2.inputMode != types.DirectMode {
		t.Error("Tab should switch SearchMode → DirectMode")
	}
}

func TestUpdateTyping_Tab_ClearsSearchQuery(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.searchQuery = "s"
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyTab})
	rm := result.(model)
	if rm.searchQuery != "" {
		t.Errorf("Tab from SearchMode should clear searchQuery, got: %q", rm.searchQuery)
	}
}

func TestUpdateTyping_DirectInput_AppendsToBuffer(t *testing.T) {
	m := testModel()
	m.inputMode = types.DirectMode
	m.activeIndex = 0 // Thai
	result, _ := updateTyping(m, makeKey("a"))
	rm := result.(model)
	// Thai "a" → "ฟ"
	if string(rm.wordBuffer) != "ฟ" {
		t.Errorf("'a' in DirectMode: want ฟ, got %q", string(rm.wordBuffer))
	}
}
