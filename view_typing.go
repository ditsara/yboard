package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ditsara/yboard/internal/types"
)

const defaultTermWidth = 80

const footerText = "F2:Setup  F3:← Lang  F4:→ Lang  F9:Copy  Enter:Copy+Clear  Tab:Mode  F10:Quit"

var (
	bufferBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	directBadgeStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				Background(lipgloss.Color("#3C3C3C")).
				Foreground(lipgloss.Color("#00FF88")).
				Bold(true).
				Padding(0, 1)

	searchBadgeStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				Background(lipgloss.Color("#3C3C3C")).
				Foreground(lipgloss.Color("#FF8800")).
				Bold(true).
				Padding(0, 1)

	langNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	statusMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555555"))

	candidateNumStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888"))

	candidateCharStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)

	searchQueryStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
)

// enabledLanguages returns the subset of languages with Enabled == true.
func enabledLanguages(m model) []types.LanguageModule {
	var out []types.LanguageModule
	for _, lang := range m.languages {
		if lang.Enabled {
			out = append(out, lang)
		}
	}
	return out
}

func viewTyping(m model) string {
	w := m.termWidth
	if w == 0 {
		w = defaultTermWidth
	}
	// Inner width: subtract 2 for the rounded border (left + right)
	innerW := w - 4 // 2 border chars + 2 padding chars

	enabled := enabledLanguages(m)
	noLangs := len(enabled) == 0

	var parts []string

	// --- Word buffer box ---
	var bufContent string
	if noLangs {
		bufContent = warningStyle.Render("⚠ No languages enabled — press F2 to configure")
	} else if len(m.wordBuffer) == 0 {
		bufContent = placeholderStyle.Render("Start typing…")
	} else {
		bufContent = string(m.wordBuffer)
	}
	parts = append(parts, bufferBoxStyle.Width(innerW).Render(bufContent))
	parts = append(parts, "")

	if noLangs {
		// Mode badge and keyboard are hidden when no languages enabled
		parts = append(parts, statusMsgStyle.Render(safeStatusLine(m)))
		parts = append(parts, "")
	} else {
		// --- Mode badge + language name ---
		// activeIndex is an index into m.languages (all, not just enabled)
		var badge string
		var langName string
		if m.activeIndex < len(m.languages) {
			langName = m.languages[m.activeIndex].Name
		}
		if m.inputMode == types.SearchMode {
			badge = searchBadgeStyle.Render("SEARCH")
		} else {
			badge = directBadgeStyle.Render("DIRECT")
		}
		badgeRow := lipgloss.JoinHorizontal(lipgloss.Center, badge, "  "+langNameStyle.Render(langName))
		parts = append(parts, lipgloss.NewStyle().PaddingLeft(1).Render(badgeRow))
		parts = append(parts, "")

		// --- Search query + candidates (SearchMode only) ---
		if m.inputMode == types.SearchMode {
			searchLine := "Search: " + searchQueryStyle.Render(m.searchQuery) + "▌"
			parts = append(parts, " "+searchLine)
			parts = append(parts, " "+buildCandidateLine(m))
			parts = append(parts, "")
		}

		// --- Status message (always reserves 1 line to prevent layout jump) ---
		parts = append(parts, " "+statusMsgStyle.Render(safeStatusLine(m)))
		parts = append(parts, "")

		// --- Visual keyboard grid ---
		kb := renderKeyboard(m)
		if kb != "" {
			parts = append(parts, kb)
			parts = append(parts, "")
		}
	}

	// --- Hotkeys footer (always visible) ---
	parts = append(parts, " "+footerStyle.Render(footerText))

	return strings.Join(parts, "\n")
}

// safeStatusLine returns the status message or a single space to hold the line height.
func safeStatusLine(m model) string {
	if m.statusMessage != "" {
		return m.statusMessage
	}
	return " "
}

// buildCandidateLine formats the phonetic search candidates as "1:ส  2:ษ  3:ศ  4:ซ".
// At most 10 candidates are shown; the 10th uses key label "0".
// Returns an empty string when there are no matches.
func buildCandidateLine(m model) string {
	if m.activeIndex >= len(m.languages) || m.searchQuery == "" {
		return ""
	}
	candidates, ok := m.languages[m.activeIndex].PhoneticMap[m.searchQuery]
	if !ok || len(candidates) == 0 {
		return ""
	}
	if len(candidates) > 10 {
		candidates = candidates[:10]
	}
	var sb strings.Builder
	for i, c := range candidates {
		if i > 0 {
			sb.WriteString("  ")
		}
		var numLabel string
		if i == 9 {
			numLabel = "0"
		} else {
			numLabel = fmt.Sprintf("%d", i+1)
		}
		sb.WriteString(candidateNumStyle.Render(numLabel + ":"))
		sb.WriteString(candidateCharStyle.Render(c))
	}
	return sb.String()
}
