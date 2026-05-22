package main

import (
	"strings"
	"testing"

	"github.com/ditsara/yboard/internal/types"
	"github.com/ditsara/yboard/modules"
)

func testModel() model {
	return model{
		state:     types.StateTyping,
		inputMode: types.DirectMode,
		languages: []types.LanguageModule{modules.ThaiModule, modules.SpanishModule},
		termWidth: 100,
		termHeight: 40,
	}
}

func TestViewTyping_HasFooter(t *testing.T) {
	m := testModel()
	out := viewTyping(m)
	if !strings.Contains(out, "F10/^Q:Quit") {
		t.Error("footer not found in typing view")
	}
	if !strings.Contains(out, "F2:Setup") {
		t.Error("F2:Setup not found in footer")
	}
}

func TestViewTyping_PlaceholderWhenEmpty(t *testing.T) {
	m := testModel()
	out := viewTyping(m)
	if !strings.Contains(out, "Start typing") {
		t.Error("placeholder text not shown when buffer is empty")
	}
}

func TestViewTyping_ShowsBufferContent(t *testing.T) {
	m := testModel()
	m.wordBuffer = []rune("สวัสดี")
	out := viewTyping(m)
	if !strings.Contains(out, "สวัสดี") {
		t.Error("word buffer content not shown")
	}
}

func TestViewTyping_NoLanguages(t *testing.T) {
	m := testModel()
	for i := range m.languages {
		m.languages[i].Enabled = false
	}
	out := viewTyping(m)
	if !strings.Contains(out, "No languages enabled") {
		t.Error("no-languages warning not shown")
	}
	if strings.Contains(out, "DIRECT") || strings.Contains(out, "SEARCH") {
		t.Error("mode badge should be hidden when no languages enabled")
	}
}

func TestViewTyping_DirectModeBadge(t *testing.T) {
	m := testModel()
	m.inputMode = types.DirectMode
	out := viewTyping(m)
	if !strings.Contains(out, "DIRECT") {
		t.Error("DIRECT badge not shown in DirectMode")
	}
}

func TestViewTyping_SearchModeBadge(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.searchQuery = "s"
	out := viewTyping(m)
	if !strings.Contains(out, "SEARCH") {
		t.Error("SEARCH badge not shown in SearchMode")
	}
	if !strings.Contains(out, "Search:") {
		t.Error("search query line not shown in SearchMode")
	}
}

func TestViewTyping_SearchCandidates(t *testing.T) {
	m := testModel()
	m.inputMode = types.SearchMode
	m.searchQuery = "s" // Thai PhoneticMap has "s": {"ส","ษ","ศ","ซ"}
	out := viewTyping(m)
	if !strings.Contains(out, "ส") {
		t.Error("phonetic candidate ส not shown")
	}
}

func TestViewTyping_StatusMessage(t *testing.T) {
	m := testModel()
	m.statusMessage = "📋 Copied!"
	out := viewTyping(m)
	if !strings.Contains(out, "📋 Copied!") {
		t.Error("status message not shown")
	}
}

func TestViewTyping_DefaultWidth(t *testing.T) {
	m := testModel()
	m.termWidth = 0 // should fall back to defaultTermWidth
	out := viewTyping(m)
	if out == "" {
		t.Error("view should not be empty even with zero termWidth")
	}
}
