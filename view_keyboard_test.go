package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/ditsara/yboard/internal/types"
	"github.com/ditsara/yboard/modules"
)

func TestRenderKey_VisualWidth(t *testing.T) {
	vk := types.VisualKey{Key: "q", Normal: "ๆ", Shift: "๐"}
	rendered := renderKey(vk)
	lines := strings.Split(rendered, "\n")
	// Should be 4 lines: top border, shift line, normal line, bottom border
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d:\n%s", len(lines), rendered)
	}
	// Visual width of top border should be 5 (╭─q─╮)
	topWidth := lipgloss.Width(lines[0])
	if topWidth != 5 {
		t.Errorf("top border visual width: want 5, got %d (line: %q)", topWidth, lines[0])
	}
}

func TestRenderKey_EmptyKey(t *testing.T) {
	vk := types.EmptyKey("=")
	rendered := renderKey(vk)
	// Empty key should still render as a 5-wide box
	lines := strings.Split(rendered, "\n")
	if len(lines) != 4 {
		t.Errorf("empty key: expected 4 lines, got %d", len(lines))
	}
	// Should contain dot placeholders
	if !strings.Contains(rendered, "·") {
		t.Error("empty key should contain · placeholder")
	}
}

func TestRenderKey_LabelInTopBorder(t *testing.T) {
	vk := types.VisualKey{Key: "a", Normal: "ฟ", Shift: "ฤ"}
	rendered := renderKey(vk)
	lines := strings.Split(rendered, "\n")
	// Top border should contain the key label "a"
	if !strings.Contains(lines[0], "a") {
		t.Errorf("top border should contain key label 'a': %q", lines[0])
	}
}

func TestRenderKeyboard_HasContent(t *testing.T) {
	m := testModel()
	m.activeIndex = 0 // Thai
	kb := renderKeyboard(m)
	if kb == "" {
		t.Error("renderKeyboard should return non-empty string for Thai module")
	}
	// Should contain some Thai characters
	if !strings.Contains(kb, "ก") {
		t.Error("Thai keyboard should contain ก")
	}
}

func TestRenderKeyboard_EmptyWhenNoLanguages(t *testing.T) {
	m := testModel()
	for i := range m.languages {
		m.languages[i].Enabled = false
	}
	kb := renderKeyboard(m)
	if kb != "" {
		t.Error("renderKeyboard should return empty string when no languages enabled")
	}
}

func TestRenderKeyboard_RowsIndented(t *testing.T) {
	m := testModel()
	kb := renderKeyboard(m)
	lines := strings.Split(kb, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "  ") {
			t.Errorf("line %d should start with 2-space indent: %q", i, line)
		}
	}
}

func TestRenderKeyboard_Spanish(t *testing.T) {
	m := model{
		state:     types.StateTyping,
		inputMode: types.DirectMode,
		languages: []types.LanguageModule{modules.SpanishModule},
		termWidth: 100,
	}
	kb := renderKeyboard(m)
	if kb == "" {
		t.Error("Spanish keyboard should not be empty")
	}
	if !strings.Contains(kb, "ñ") {
		t.Error("Spanish keyboard should contain ñ")
	}
}

func TestRenderKeyboard_RowCount(t *testing.T) {
	m := testModel()
	kb := renderKeyboard(m)
	// Each row is 4 lines tall; 4 rows → 16 lines (plus potential extra newlines)
	lines := strings.Split(kb, "\n")
	// Top borders contain ╭, check we have 4 rows worth
	topBorders := 0
	for _, line := range lines {
		if strings.Contains(line, "╭") {
			topBorders++
			break // just check first line of each row
		}
	}
	if topBorders == 0 {
		t.Error("keyboard should contain rounded box top borders")
	}
	// Rough check: should have at least 16 lines for 4 rows × 4 lines each
	if len(lines) < 16 {
		t.Errorf("keyboard should have at least 16 lines for 4 rows, got %d", len(lines))
	}
}
