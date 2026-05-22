package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ditsara/yboard/internal/types"
	"github.com/ditsara/yboard/modules"
)

type model struct {
	state         types.AppState
	inputMode     types.InputMode
	languages     []types.LanguageModule
	activeIndex   int
	wordBuffer    []rune
	searchQuery   string
	statusMessage string
	termWidth     int
	termHeight    int
	setupCursor   int
}

func initialModel() model {
	langs := []types.LanguageModule{
		modules.ThaiModule,
		modules.SpanishModule,
	}
	for i := range langs {
		langs[i].DirectMap, langs[i].ShiftDirectMap = types.BuildDirectMaps(langs[i].KeyboardRows)
	}
	return model{
		state:     types.StateTyping,
		inputMode: types.DirectMode,
		languages: langs,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyF10, tea.KeyCtrlQ:
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlL:
			return m, tea.ClearScreen
		}
		switch m.state {
		case types.StateTyping:
			return updateTyping(m, msg)
		case types.StateSetup:
			return updateSetup(m, msg)
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case types.StateSetup:
		return viewSetup(m)
	default:
		return viewTyping(m)
	}
}


func main() {
	m := initialModel()
	// Validate language modules at startup
	for _, lang := range m.languages {
		if err := types.ValidateModule(lang); err != nil {
			fmt.Fprintf(os.Stderr, "module error: %v\n", err)
			os.Exit(1)
		}
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
