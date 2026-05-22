package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ditsara/yboard/internal/types"
)

func updateSetup(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.setupCursor > 0 {
			m.setupCursor--
		}
	case tea.KeyDown:
		if m.setupCursor < len(m.languages)-1 {
			m.setupCursor++
		}
	case tea.KeySpace:
		if m.setupCursor < len(m.languages) {
			m.languages[m.setupCursor].Enabled = !m.languages[m.setupCursor].Enabled
		}
	case tea.KeyF2, tea.KeyEnter, tea.KeyEscape:
		m.state = types.StateTyping
		m = snapToEnabledLanguage(m)
	}

	// vim-style navigation
	switch msg.String() {
	case "k":
		if m.setupCursor > 0 {
			m.setupCursor--
		}
	case "j":
		if m.setupCursor < len(m.languages)-1 {
			m.setupCursor++
		}
	}

	return m, nil
}
