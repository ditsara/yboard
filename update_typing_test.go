package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ditsara/yboard/internal/types"
)

func TestNextEnabledLanguage_WrapsAround(t *testing.T) {
	m := testModel()
	m.activeIndex = 1 // Spanish (last)
	m = nextEnabledLanguage(m)
	if m.activeIndex != 0 {
		t.Errorf("should wrap to Thai (0), got %d", m.activeIndex)
	}
}

func TestNextEnabledLanguage_SkipsDisabled(t *testing.T) {
	m := testModel()
	m.languages[1].Enabled = false // disable Spanish
	m.activeIndex = 0              // Thai
	m = nextEnabledLanguage(m)
	if m.activeIndex != 0 {
		t.Errorf("should stay on Thai (0) when Spanish is disabled, got %d", m.activeIndex)
	}
}

func TestPrevEnabledLanguage_WrapsAround(t *testing.T) {
	m := testModel()
	m.activeIndex = 0 // Thai (first)
	m = prevEnabledLanguage(m)
	if m.activeIndex != 1 {
		t.Errorf("should wrap to Spanish (1), got %d", m.activeIndex)
	}
}

func TestNextEnabledLanguage_SetsStatusMessage(t *testing.T) {
	m := testModel()
	m.activeIndex = 0
	m = nextEnabledLanguage(m)
	if !strings.Contains(m.statusMessage, "Switched to") {
		t.Errorf("should set status message on switch, got: %q", m.statusMessage)
	}
}

func TestNextEnabledLanguage_ClearsSearchQuery(t *testing.T) {
	m := testModel()
	m.searchQuery = "s"
	m.activeIndex = 0
	m = nextEnabledLanguage(m)
	if m.searchQuery != "" {
		t.Errorf("searchQuery should be cleared on switch, got: %q", m.searchQuery)
	}
}

func TestNextEnabledLanguage_NoOp_WhenAllDisabled(t *testing.T) {
	m := testModel()
	for i := range m.languages {
		m.languages[i].Enabled = false
	}
	original := m.activeIndex
	m = nextEnabledLanguage(m)
	if m.activeIndex != original {
		t.Errorf("should not change activeIndex when all disabled")
	}
}

func TestSnapToEnabledLanguage_FixesDisabledIndex(t *testing.T) {
	m := testModel()
	m.languages[0].Enabled = false // disable Thai
	m.activeIndex = 0              // pointing at disabled Thai
	m = snapToEnabledLanguage(m)
	if m.activeIndex != 1 {
		t.Errorf("should snap to Spanish (1), got %d", m.activeIndex)
	}
}

func TestSnapToEnabledLanguage_NoChangeWhenValid(t *testing.T) {
	m := testModel()
	m.activeIndex = 0
	m = snapToEnabledLanguage(m)
	if m.activeIndex != 0 {
		t.Errorf("should stay at 0 when Thai is enabled, got %d", m.activeIndex)
	}
}

func TestUpdateTyping_F3(t *testing.T) {
	m := testModel()
	m.activeIndex = 1
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyF3})
	rm := result.(model)
	if rm.activeIndex != 0 {
		t.Errorf("F3 should go to previous language, got activeIndex=%d", rm.activeIndex)
	}
}

func TestUpdateTyping_F4(t *testing.T) {
	m := testModel()
	m.activeIndex = 0
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyF4})
	rm := result.(model)
	if rm.activeIndex != 1 {
		t.Errorf("F4 should go to next language, got activeIndex=%d", rm.activeIndex)
	}
}

func TestUpdateTyping_F2_GoesToSetup(t *testing.T) {
	m := testModel()
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyF2})
	rm := result.(model)
	if rm.state != types.StateSetup {
		t.Error("F2 should switch to StateSetup")
	}
}

func TestUpdateTyping_ClearsStatusOnKeypress(t *testing.T) {
	m := testModel()
	m.statusMessage = "old message"
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyF4})
	rm := result.(model)
	// After F4, status should be the switch message (not the old one)
	if strings.Contains(rm.statusMessage, "old message") {
		t.Error("statusMessage should be cleared on keypress")
	}
}
