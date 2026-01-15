package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// addKeymap is a struct defining the keys used in the add view
type addKeymap struct {
	Help key.Binding
	Quit key.Binding
}

// ShortHelp returns the keys to show in compact help view
func (k addKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns the keys to show in complete help view
func (k addKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
	}
}

var addKeys = addKeymap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "show keybindings"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type addState struct {
	keys addKeymap
	form *huh.Form
}

func initAddState() addState {
	return addState{
		keys: addKeys,
		form: huh.NewForm(),
	}
}

func (m model) addUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.state.add.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.state.add.keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) addView() string {
	out := ""
	return out
}
