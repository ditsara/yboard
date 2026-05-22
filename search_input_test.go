package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ditsara/yboard/internal/types"
)

// --- handleSearchInput tests ---

func TestHandleSearchInput_AlphaAppendsToQuery(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.activeIndex = 0 // Thai
	m = handleSearchInput(m, makeKey("s"))
	if m.searchQuery != "s" {
		t.Errorf("want searchQuery=%q, got %q", "s", m.searchQuery)
	}
}

func TestHandleSearchInput_ShiftNormalizedToLower(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.activeIndex = 0
	// Pass "S" — should normalize to "s"
	m = handleSearchInput(m, makeKey("S"))
	if m.searchQuery != "s" {
		t.Errorf("shifted alpha should normalize to lowercase, got %q", m.searchQuery)
	}
}

func TestHandleSearchInput_SelectCandidate1(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.activeIndex = 0
	m.searchQuery = "s"
	// Thai PhoneticMap["s"] = {"ส", "ษ", "ศ", "ซ"}; key "1" → index 0 → "ส"
	m = handleSearchInput(m, makeKey("1"))
	if string(m.wordBuffer) != "ส" {
		t.Errorf("'1' should select first candidate ส, got %q", string(m.wordBuffer))
	}
	if m.searchQuery != "" {
		t.Errorf("selecting candidate should clear searchQuery, got %q", m.searchQuery)
	}
}

func TestHandleSearchInput_SelectCandidate4(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.activeIndex = 0
	m.searchQuery = "s"
	// Thai PhoneticMap["s"][3] = "ซ"
	m = handleSearchInput(m, makeKey("4"))
	if string(m.wordBuffer) != "ซ" {
		t.Errorf("'4' should select 4th candidate ซ, got %q", string(m.wordBuffer))
	}
}

func TestHandleSearchInput_SelectCandidate0Is10th(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.activeIndex = 0
	// Use a query with a known short list; "0" out of range → no append but still clears query
	m.searchQuery = "s"
	m = handleSearchInput(m, makeKey("0"))
	// Thai PhoneticMap["s"] has only 4 entries; index 9 is out of range
	if len(m.wordBuffer) != 0 {
		t.Errorf("'0' out of range should not append anything, got %q", string(m.wordBuffer))
	}
	if m.searchQuery != "" {
		t.Errorf("selecting (even out-of-range) should clear searchQuery")
	}
}

func TestHandleSearchInput_NoMatch_QueryStillBuilds(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.activeIndex = 0
	m.searchQuery = "s"
	m = handleSearchInput(m, makeKey("s"))
	if m.searchQuery != "ss" {
		t.Errorf("second alpha should extend query to 'ss', got %q", m.searchQuery)
	}
}

func TestHandleSearchInput_NoLanguages(t *testing.T) {
	m := testModel()
	for i := range m.languages {
		m.languages[i].Enabled = false
	}
	m.inputMode = types.SearchMode
	m = handleSearchInput(m, makeKey("s"))
	if !strings.Contains(m.statusMessage, "No languages enabled") {
		t.Errorf("should show no-languages warning, got: %q", m.statusMessage)
	}
}

// --- SearchMode integration via updateTyping ---

func TestUpdateTyping_SearchMode_AlphaBuildsQuery(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.activeIndex = 0
	result, _ := updateTyping(m, makeKey("t"))
	rm := result.(model)
	if rm.searchQuery != "t" {
		t.Errorf("SearchMode alpha should extend searchQuery, got %q", rm.searchQuery)
	}
}

func TestUpdateTyping_SearchMode_BackspacePoopsQuery(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.searchQuery = "st"
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyBackspace})
	rm := result.(model)
	if rm.searchQuery != "s" {
		t.Errorf("backspace in SearchMode should pop query char, got %q", rm.searchQuery)
	}
}

func TestUpdateTyping_SearchMode_BackspaceEmptyQueryPopsBuffer(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.searchQuery = ""
	m.wordBuffer = []rune("ส")
	result, _ := updateTyping(m, tea.KeyMsg{Type: tea.KeyBackspace})
	rm := result.(model)
	if len(rm.wordBuffer) != 0 {
		t.Errorf("backspace on empty query should pop wordBuffer, got %q", string(rm.wordBuffer))
	}
}

func TestUpdateTyping_SearchMode_NumberSelectsCandidate(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.activeIndex = 0
	m.searchQuery = "s"
	result, _ := updateTyping(m, makeKey("2"))
	rm := result.(model)
	// Thai PhoneticMap["s"][1] = "ษ"
	if string(rm.wordBuffer) != "ษ" {
		t.Errorf("'2' in SearchMode should select 2nd candidate ษ, got %q", string(rm.wordBuffer))
	}
	if rm.searchQuery != "" {
		t.Error("searchQuery should be cleared after selection")
	}
}

// --- buildCandidateLine rendering ---

func TestBuildCandidateLine_EmptyQuery(t *testing.T) {
	m := testModel()
	m.activeIndex = 0
	m.searchQuery = ""
	line := buildCandidateLine(m)
	if line != "" {
		t.Errorf("empty query should produce empty candidate line, got %q", line)
	}
}

func TestBuildCandidateLine_NoMatch(t *testing.T) {
	m := testModel()
	m.activeIndex = 0
	m.searchQuery = "zzz"
	line := buildCandidateLine(m)
	if line != "" {
		t.Errorf("unmatched query should produce empty candidate line, got %q", line)
	}
}

func TestBuildCandidateLine_ShowsCandidates(t *testing.T) {
	m := testModel()
	m.activeIndex = 0
	m.searchQuery = "s"
	line := buildCandidateLine(m)
	// Should contain "1:" and "ส"
	if !strings.Contains(line, "1:") {
		t.Errorf("candidate line should contain '1:', got %q", line)
	}
	if !strings.Contains(line, "ส") {
		t.Errorf("candidate line should contain Thai candidate ส, got %q", line)
	}
}

func TestBuildCandidateLine_CapsAt10(t *testing.T) {
	m := testModel()
	m.activeIndex = 0
	// Inject a 12-entry phonetic map to test capping
	m.languages[0].PhoneticMap["x"] = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
	m.searchQuery = "x"
	line := buildCandidateLine(m)
	// Should not contain "11:" (index beyond 10)
	if strings.Contains(line, "11:") {
		t.Errorf("candidate line should cap at 10, but found '11:' in: %q", line)
	}
	// Should contain "0:" for the 10th candidate
	if !strings.Contains(line, "0:") {
		t.Errorf("10th candidate should use '0:' label, got: %q", line)
	}
}
