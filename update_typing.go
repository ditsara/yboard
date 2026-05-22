package main

import (
"fmt"
"strings"

tea "github.com/charmbracelet/bubbletea"
"github.com/ditsara/yboard/internal/types"
)

// handleSearchInput processes a key press in SearchMode.
// Alphabetic keys (normalized to lowercase) extend the searchQuery.
// Number keys 1-9 select candidate at index n-1; 0 selects the 10th (index 9).
func handleSearchInput(m model, msg tea.KeyMsg) model {
	if m.activeIndex >= len(m.languages) || !m.languages[m.activeIndex].Enabled {
		if len(enabledLanguages(m)) == 0 {
			m.statusMessage = "⚠ No languages enabled — press F2 to configure"
			return m
		}
	}

	keyStr := msg.String()

	// Number key: select a candidate and clear the query
	if len(keyStr) == 1 && keyStr[0] >= '0' && keyStr[0] <= '9' {
		var idx int
		if keyStr[0] == '0' {
			idx = 9
		} else {
			idx = int(keyStr[0]-'1') // '1'→0 … '9'→8
		}
		candidates := m.languages[m.activeIndex].PhoneticMap[m.searchQuery]
		if idx < len(candidates) {
			m.wordBuffer = append(m.wordBuffer, []rune(candidates[idx])...)
		}
		m.searchQuery = ""
		return m
	}

	// Alphabetic key: extend searchQuery (Shift ignored — normalize to lowercase)
	lower := strings.ToLower(keyStr)
	runes := []rune(lower)
	if len(runes) == 1 && runes[0] >= 'a' && runes[0] <= 'z' {
		m.searchQuery += lower
	}
	return m
}

// handleDirectInput processes a printable key press in DirectMode.
// Looks up the key in DirectMap then ShiftDirectMap; falls through as literal on miss.
func handleDirectInput(m model, msg tea.KeyMsg) model {
if m.activeIndex >= len(m.languages) || !m.languages[m.activeIndex].Enabled {
if len(enabledLanguages(m)) == 0 {
m.statusMessage = "⚠ No languages enabled — press F2 to configure"
return m
}
}

lang := m.languages[m.activeIndex]
keyStr := msg.String()

if result, ok := lang.DirectMap[keyStr]; ok {
m.wordBuffer = append(m.wordBuffer, []rune(result)...)
} else if result, ok := lang.ShiftDirectMap[keyStr]; ok {
m.wordBuffer = append(m.wordBuffer, []rune(result)...)
} else {
m.wordBuffer = append(m.wordBuffer, []rune(keyStr)...)
m.statusMessage = fmt.Sprintf("⚠ Unknown key: '%s' — passed through", keyStr)
}
return m
}

// nextEnabledLanguage advances activeIndex to the next enabled language (with wrap-around).
// If no language is enabled, the model is returned unchanged.
func nextEnabledLanguage(m model) model {
n := len(m.languages)
for i := 1; i <= n; i++ {
idx := (m.activeIndex + i) % n
if m.languages[idx].Enabled {
m.activeIndex = idx
m.searchQuery = ""
m.statusMessage = fmt.Sprintf("🌐 Switched to %s", m.languages[idx].Name)
return m
}
}
return m
}

// prevEnabledLanguage retreats activeIndex to the previous enabled language (with wrap-around).
// If no language is enabled, the model is returned unchanged.
func prevEnabledLanguage(m model) model {
n := len(m.languages)
for i := 1; i <= n; i++ {
idx := (m.activeIndex - i + n) % n
if m.languages[idx].Enabled {
m.activeIndex = idx
m.searchQuery = ""
m.statusMessage = fmt.Sprintf("🌐 Switched to %s", m.languages[idx].Name)
return m
}
}
return m
}

// snapToEnabledLanguage ensures activeIndex points to an enabled language.
// Called on re-entry to StateTyping after potential module toggling in Setup.
func snapToEnabledLanguage(m model) model {
n := len(m.languages)
if n == 0 {
return m
}
if m.activeIndex < n && m.languages[m.activeIndex].Enabled {
return m
}
for i := 0; i < n; i++ {
if m.languages[i].Enabled {
m.activeIndex = i
return m
}
}
return m
}

// updateTyping handles key events in StateTyping.
// Search input logic (YB-8) will be added to the SearchMode branch.
func updateTyping(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
// Clear any transient status on next keypress
m.statusMessage = ""

switch msg.Type {
case tea.KeyF2:
m.state = types.StateSetup
return m, nil
case tea.KeyF3:
m = prevEnabledLanguage(m)
return m, nil
case tea.KeyF4:
m = nextEnabledLanguage(m)
return m, nil
case tea.KeyF9:
m.statusMessage = executeClipboardCopy(string(m.wordBuffer))
return m, nil
case tea.KeyEnter:
m.statusMessage = executeClipboardCopy(string(m.wordBuffer))
m.wordBuffer = nil
return m, nil
case tea.KeyTab:
if m.inputMode == types.DirectMode {
m.inputMode = types.SearchMode
} else {
m.inputMode = types.DirectMode
m.searchQuery = ""
}
return m, nil
case tea.KeySpace:
if m.inputMode == types.DirectMode {
m.wordBuffer = append(m.wordBuffer, ' ')
}
return m, nil
case tea.KeyBackspace:
if m.inputMode == types.SearchMode && m.searchQuery != "" {
runes := []rune(m.searchQuery)
m.searchQuery = string(runes[:len(runes)-1])
} else if len(m.wordBuffer) > 0 {
m.wordBuffer = m.wordBuffer[:len(m.wordBuffer)-1]
}
return m, nil
case tea.KeyRunes:
switch m.inputMode {
case types.DirectMode:
m = handleDirectInput(m, msg)
case types.SearchMode:
m = handleSearchInput(m, msg)
}
}
return m, nil
}
