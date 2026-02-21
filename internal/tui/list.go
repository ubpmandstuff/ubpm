package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// listKeymap is a struct defining the keys used in the list view
type listKeymap struct {
	Up      key.Binding
	Down    key.Binding
	Add     key.Binding
	Edit    key.Binding
	Rm      key.Binding
	SeePass key.Binding
	Help    key.Binding
	Quit    key.Binding
}

// ShortHelp returns the keys to show in compact help view
func (k listKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns the keys to show in complete help view
func (k listKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.SeePass},
		{k.Add, k.Edit, k.Rm},
		{k.Help, k.Quit},
	}
}

var listKeys = listKeymap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add entry"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit entry"),
	),
	Rm: key.NewBinding(
		key.WithKeys("d", "r"),
		key.WithHelp("d/r", "delete entry"),
	),
	SeePass: key.NewBinding(
		key.WithKeys(" ", "v"),
		key.WithHelp("v/space", "peek passwd"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "show keybindings"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type listState struct {
	cursor        int
	keys          listKeymap
	passwdVisible bool
}

func initListState() listState {
	return listState{
		cursor:        0,
		keys:          listKeys,
		passwdVisible: false,
	}
}

func (m model) listUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.state.list.keys.Up):
			m.state.list.passwdVisible = false
			if m.state.list.cursor > 0 {
				m.state.list.cursor--
			}
		case key.Matches(msg, m.state.list.keys.Down):
			m.state.list.passwdVisible = false
			if m.state.list.cursor < len(m.vault.Data.Entries)-1 {
				m.state.list.cursor++
			}
		case key.Matches(msg, m.state.list.keys.Add):
			m.state.list.passwdVisible = false
			return m, m.switchAdd()
		case key.Matches(msg, m.state.list.keys.Edit):
			m.state.list.passwdVisible = false
			return m, m.switchEdit(m.vault.Data.Entries[m.state.list.cursor])
		case key.Matches(msg, m.state.list.keys.Rm):
			m.state.list.passwdVisible = false
			return m, m.switchRm(m.vault.Data.Entries[m.state.list.cursor])
		case key.Matches(msg, m.state.list.keys.SeePass):
			m.state.list.passwdVisible = !m.state.list.passwdVisible
		case key.Matches(msg, m.state.list.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.state.list.keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) listView() string {
	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7e98e8"))
	b1Style := lipgloss.NewStyle().Padding(1, 2).MarginRight(2)
	b2Style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7e98e8")).
		Padding(0, 1)
	errBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#d8647e")).
		Padding(0, 1)

	var b1 strings.Builder

	b1.WriteString("ubpm, a usb-based password manager\n\n")

	if len(m.vault.Data.Entries) > 0 {
		for i, e := range m.vault.Data.Entries {
			if m.state.list.cursor == i {
				fmt.Fprintf(&b1, "%s %s\n", cursorStyle.Render("> "), e.Title)
			} else {
				fmt.Fprintf(&b1, "%s %s\n", "  ", e.Title)
			}
		}
	} else {
		b1.WriteString("no entries\n")
	}

	if m.errMsg != nil {
		b1.WriteString("\n" + errBox.Render(m.errMsg.Error()) + "\n")
	}

	b1.WriteString("\n")
	b1.WriteString(m.help.View(m.state.list.keys))

	if len(m.vault.Data.Entries) > 0 {
		var b2 strings.Builder
		e := m.vault.Data.Entries[m.state.list.cursor]
		var pass string
		if m.state.list.passwdVisible {
			pass = e.Password
		} else {
			pass = "********"
		}

		fmt.Fprintf(&b2,
			"title: %s\nusername: %s\npassword: %s\nnotes: %s\n\ncreated at: %s\nmodified at: %s\nid: %s",
			e.Title,
			e.Username,
			pass,
			e.Notes,
			e.CreatedAt.Format("2006-01-02 15:04 MST"),
			e.ModifiedAt.Format("2006-01-02 15:04 MST"),
			e.ID)

		// boxStyle.Render(b2.String())
		return lipgloss.JoinHorizontal(lipgloss.Left, b1Style.Render(b1.String()), b2Style.Render(b2.String()))
	} else {
		return b1Style.Render(b1.String())
	}
}
