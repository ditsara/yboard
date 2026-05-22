package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ditsara/yboard/internal/types"
)

var (
	keyShiftStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	keyNormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	keyLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")) // gold — the physical key to press

	mutedLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555")) // signals "not remapped by this language"

	emptySlotStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333"))

	keyBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)
)

// renderKey renders a single VisualKey as a 5-wide × 4-tall rounded box.
// The English key label is injected into the center of the top border.
// Empty keys (no Normal/Shift) are rendered dim to signal "not remapped".
func renderKey(vk types.VisualKey) string {
	isEmpty := vk.Normal == "" && vk.Shift == ""

	labelSty := keyLabelStyle
	shiftSty := keyShiftStyle
	normSty := keyNormalStyle
	if isEmpty {
		labelSty = mutedLabelStyle
		shiftSty = emptySlotStyle
		normSty = emptySlotStyle
	}

	shiftChar := vk.Shift
	if shiftChar == "" {
		shiftChar = "·"
		if !isEmpty {
			// Partially remapped: Normal exists but Shift does not
			shiftSty = lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
		}
	}
	normalChar := vk.Normal
	if normalChar == "" {
		normalChar = "·"
	}

	inner := lipgloss.JoinVertical(lipgloss.Center,
		shiftSty.Render(shiftChar),
		normSty.Render(normalChar),
	)
	box := keyBoxStyle.Render(inner)

	// Inject the English key label into the center dash of the top border.
	// Top border of a 5-wide rounded box is "╭───╮"; we replace rune at index 2.
	lines := strings.Split(box, "\n")
	label := labelSty.Render(vk.Key)
	runes := []rune(lines[0])
	if len(runes) >= 5 {
		lines[0] = string(runes[:2]) + label + string(runes[3:])
	}
	return strings.Join(lines, "\n")
}

// renderKeyboard renders the full keyboard grid for the active language module.
// Rows are centered relative to the widest row using lipgloss.PlaceHorizontal.
// Returns empty string when no language is active or the active language is disabled.
func renderKeyboard(m model) string {
	if m.activeIndex >= len(m.languages) {
		return ""
	}
	mod := m.languages[m.activeIndex]
	if !mod.Enabled || len(mod.KeyboardRows) == 0 {
		return ""
	}

	maxWidth := 0
	renderedRows := make([]string, len(mod.KeyboardRows))
	for i, row := range mod.KeyboardRows {
		keys := make([]string, len(row))
		for j, vk := range row {
			keys[j] = renderKey(vk)
		}
		renderedRows[i] = lipgloss.JoinHorizontal(lipgloss.Top, keys...)
		if w := lipgloss.Width(renderedRows[i]); w > maxWidth {
			maxWidth = w
		}
	}

	// Center shorter rows relative to the widest row
	for i, row := range renderedRows {
		renderedRows[i] = lipgloss.PlaceHorizontal(maxWidth, lipgloss.Center, row)
	}

	kb := lipgloss.JoinVertical(lipgloss.Left, renderedRows...)

	// Indent the entire keyboard block by 2 spaces
	lines := strings.Split(kb, "\n")
	for i, line := range lines {
		lines[i] = "  " + line
	}
	return strings.Join(lines, "\n")
}

