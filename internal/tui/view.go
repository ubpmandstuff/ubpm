package tui

import (
	"dura5ka/ubpm/internal/vault"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type viewKeymap struct {
	SeePass key.Binding
	Back    key.Binding
	// Help    key.Binding
	Quit key.Binding
}

func (k viewKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.SeePass, k.Back, k.Quit}
}

func (k viewKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SeePass, k.Back, k.Quit},
	}
}

var viewKeys = viewKeymap{
	SeePass: key.NewBinding(
		key.WithKeys(" ", "f2"),
		key.WithHelp("space/f2", "peek passwd"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to entry list"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type viewState struct {
	keys          viewKeymap
	passwdVisible bool
	entry         vault.Entry
}

func initViewState(e vault.Entry) viewState {
	return viewState{
		keys:          viewKeys,
		passwdVisible: false,
		entry:         e,
	}
}

func (m model) viewUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.state.view.keys.SeePass):
			m.state.view.passwdVisible = !m.state.view.passwdVisible
		case key.Matches(msg, m.state.view.keys.Back):
			return m, m.switchList()
		case key.Matches(msg, m.state.view.keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) viewView() string {
	var b1 strings.Builder

	e := m.state.view.entry
	fmt.Fprintf(&b1, "viewing entry %s (%s)\n\n", e.Title, e.ID[:8])

	var pass, notes string
	if m.state.view.passwdVisible {
		pass = e.Password
	} else {
		pass = "********"
	}
	if e.Notes != "" {
		notes = "\n" + e.Notes
	} else {
		notes = "[empty]"
	}

	fmt.Fprintf(&b1,
		"title:    %s\nusername: %s\npassword: %s\nnotes:    %s\n\n---\ncreated at: %s\nmodified at: %s\nid: %s\n\n",
		e.Title,
		e.Username,
		pass,
		notes,
		e.CreatedAt.Format("2006-01-02 15:04 MST"),
		e.ModifiedAt.Format("2006-01-02 15:04 MST"),
		e.ID)

	b1.WriteString(m.help.View(m.state.view.keys))

	return b1.String()
}
