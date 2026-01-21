package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// addKeymap is a struct defining the keys used in the add view
type addKeymap struct {
	Help key.Binding
	Back key.Binding
	Quit key.Binding
}

type cAddKeymap struct {
	k addKeymap // k means keys
	f *huh.Form // f means form
}

// ShortHelp returns the keys to show in compact help view
func (c cAddKeymap) ShortHelp() []key.Binding {
	binds := c.f.KeyBinds()
	return append(binds, c.k.Help, c.k.Back)
}

// FullHelp returns the keys to show in complete help view
func (c cAddKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		c.f.KeyBinds(),
		{c.k.Help, c.k.Back, c.k.Quit},
	}
}

var addKeys = addKeymap{
	Help: key.NewBinding(
		key.WithKeys("f1"),
		key.WithHelp("f1", "show keybindings"),
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

type addState struct {
	keys addKeymap
	form *huh.Form
	val  *addValues
}

type addValues struct {
	title    string
	username string
	password string
	notes    string
}

func initAddState() addState {
	s := addState{
		keys: addKeys,
		val:  &addValues{},
	}

	s.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("title").
				Value(&s.val.title),
			huh.NewInput().
				Title("username").
				Value(&s.val.username),
			huh.NewInput().
				Title("password").
				EchoMode(huh.EchoModePassword).
				Value(&s.val.password),
			huh.NewText().
				Title("notes").
				Lines(3).
				Value(&s.val.notes),
		),
	).WithShowHelp(false)

	return s
}

func (m model) addUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.state.add.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.state.add.form = f
	}

	if m.state.add.form.State == huh.StateCompleted {
		v := m.state.add.val
		m.vault.AddEntry(v.title, v.username, v.password, v.notes)
		return m, m.switchList()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.state.add.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.state.add.keys.Back):
			return m, m.switchList()
		case key.Matches(msg, m.state.add.keys.Quit):
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) addView() string {
	b1Style := lipgloss.NewStyle().Padding(1, 2)

	var b1 strings.Builder

	b1.WriteString("adding entry\n\n")
	b1.WriteString(m.state.add.form.View() + "\n\n")

	combined := cAddKeymap{
		k: m.state.add.keys,
		f: m.state.add.form,
	}
	b1.WriteString(m.help.View(combined))

	out := b1Style.Render(b1.String())

	return out
}
