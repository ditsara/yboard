package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	setupCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF88")).
				Bold(true)

	setupCheckedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF88"))

	setupUncheckedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555555"))

	setupSelectedNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)

	setupNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	setupWarningStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF4444"))
)

const setupFooterText = "Space: toggle   ↑/k: up   ↓/j: down   F2/Enter/Esc: return to typing"

func viewSetup(m model) string {
	w := m.termWidth
	if w == 0 {
		w = defaultTermWidth
	}
	innerW := w - 4 // 2 border chars + 2 padding

	var parts []string

	// --- Title box ---
	titleBox := bufferBoxStyle.Width(innerW).Render("Language Modules")
	parts = append(parts, titleBox)
	parts = append(parts, "")

	// --- Module list ---
	for i, lang := range m.languages {
		isSelected := i == m.setupCursor

		// Cursor indicator column (5 chars wide)
		var cursor string
		if isSelected {
			cursor = "  " + setupCursorStyle.Render("▶") + "  "
		} else {
			cursor = "     "
		}

		// Checkbox
		var checkbox string
		if lang.Enabled {
			checkbox = setupCheckedStyle.Render("[✓]")
		} else {
			checkbox = setupUncheckedStyle.Render("[ ]")
		}

		// Name
		var name string
		if isSelected {
			name = setupSelectedNameStyle.Render(lang.Name)
		} else {
			name = setupNameStyle.Render(lang.Name)
		}

		parts = append(parts, cursor+checkbox+"  "+name)
	}

	// --- Warning if all modules disabled ---
	allDisabled := true
	for _, lang := range m.languages {
		if lang.Enabled {
			allDisabled = false
			break
		}
	}
	if allDisabled {
		parts = append(parts, "")
		parts = append(parts, "  "+setupWarningStyle.Render("⚠  No languages enabled. Return to typing and configure at least one."))
	}

	parts = append(parts, "")
	parts = append(parts, "")

	// --- Footer ---
	parts = append(parts, " "+footerStyle.Render(setupFooterText))

	return strings.Join(parts, "\n")
}
