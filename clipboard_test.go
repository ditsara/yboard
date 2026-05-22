package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestExecuteClipboardCopy_EmptyBuffer(t *testing.T) {
	result := executeClipboardCopy("")
	if result != "Buffer empty" {
		t.Errorf("empty buffer: want %q, got %q", "Buffer empty", result)
	}
}

func TestExecuteClipboardCopy_NonEmpty(t *testing.T) {
	// This test verifies the function executes without panic.
	// Whether it succeeds depends on available clipboard tools in the test environment.
	result := executeClipboardCopy("สวัสดี")
	if result == "" {
		t.Error("should return non-empty status message")
	}
	if result != "📋 Copied text block to clipboard!" && result != "❌ Clipboard process pipeline error" {
		t.Errorf("unexpected return value: %q", result)
	}
}

func TestUpdateTyping_F9_SetsStatus(t *testing.T) {
	m := testModel()
	m.wordBuffer = []rune("สวัสดี")
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyF9})
	rm := result.(model)
	if rm.statusMessage == "" {
		t.Error("F9 should set a status message")
	}
	// Buffer should NOT be cleared on F9
	if string(rm.wordBuffer) != "สวัสดี" {
		t.Errorf("F9 should not clear buffer, got: %q", string(rm.wordBuffer))
	}
}

func TestUpdateTyping_Enter_ClearsBuffer(t *testing.T) {
	m := testModel()
	m.wordBuffer = []rune("test")
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyEnter})
	rm := result.(model)
	if len(rm.wordBuffer) != 0 {
		t.Errorf("Enter should clear buffer, got: %q", string(rm.wordBuffer))
	}
	if rm.statusMessage == "" {
		t.Error("Enter should set a status message")
	}
}

func TestUpdateTyping_Enter_EmptyBuffer(t *testing.T) {
	m := testModel()
	m.wordBuffer = nil
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyEnter})
	rm := result.(model)
	if !strings.Contains(rm.statusMessage, "Buffer empty") {
		t.Errorf("Enter with empty buffer should show 'Buffer empty', got: %q", rm.statusMessage)
	}
}
