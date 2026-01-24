package tui

import (
	"fmt"
	"strings"

	"dura5ka/ubpm/internal/vault"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type rmKeymap struct {
	Help key.Binding
	Back key.Binding
	Quit key.Binding
}
type cRmKeymap struct {
	k rmKeymap  // k means keys
	f *huh.Form // f means form
}

// ShortHelp returns the keys to show in compact help view
func (c cRmKeymap) ShortHelp() []key.Binding {
	binds := c.f.KeyBinds()
	return append(binds, c.k.Help, c.k.Back)
}

// FullHelp returns the keys to show in complete help view
func (c cRmKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		c.f.KeyBinds(),
		{c.k.Help, c.k.Back, c.k.Quit},
	}
}

var rmKeys = rmKeymap{
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

type rmState struct {
	keys    rmKeymap
	form    *huh.Form
	confirm bool
	id      string
}

func initRmState(e vault.Entry) *rmState {
	s := &rmState{
		keys:    rmKeys,
		confirm: false,
		id:      e.ID,
	}

	title := fmt.Sprintf("are you sure you want to delete entry %s?", e.ID[:8])

	s.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Value(&s.confirm),
		),
	).WithShowHelp(false)

	return s
}

func (m model) rmUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.state.rm.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.state.rm.form = f
	}

	if m.state.rm.form.State == huh.StateCompleted {
		if m.state.rm.confirm {
			m.vault.RemoveEntry(m.state.rm.id)
			return m, m.switchList()
		} else {
			m.errMsg = fmt.Errorf("operation aborted: delete entry %s", m.state.rm.id[:8])
			return m, tea.Batch(m.switchList(), clearErrDelayed())
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.state.rm.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.state.rm.keys.Back):
			return m, m.switchList()
		case key.Matches(msg, m.state.rm.keys.Quit):
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) rmView() string {
	b1Style := lipgloss.NewStyle().Padding(1, 2)

	var b1 strings.Builder

	b1.WriteString(m.state.rm.form.View() + "\n\n")

	combined := cRmKeymap{
		k: m.state.rm.keys,
		f: m.state.rm.form,
	}
	b1.WriteString(m.help.View(combined))

	out := b1Style.Render(b1.String())

	return out
}
