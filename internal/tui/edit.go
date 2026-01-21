package tui

import (
	"strings"

	"dura5ka/ubpm/internal/vault"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	// "github.com/charmbracelet/lipgloss"
)

// addKeymap is a struct defining the keys used in the add view
type editKeymap struct {
	Help key.Binding
	Back key.Binding
	Quit key.Binding
}

type cEditKeymap struct {
	k editKeymap // k means keys
	f *huh.Form  // f means form
}

// ShortHelp returns the keys to show in compact help view
func (c cEditKeymap) ShortHelp() []key.Binding {
	binds := c.f.KeyBinds()
	return append(binds, c.k.Help, c.k.Back)
}

// FullHelp returns the keys to show in complete help view
func (c cEditKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		c.f.KeyBinds(),
		{c.k.Help, c.k.Back, c.k.Quit},
	}
}

var editKeys = editKeymap{
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

type editState struct {
	keys    editKeymap
	form    *huh.Form
	val     *editValues
	origval editValues
	id      string
}

type editValues struct {
	title    string
	username string
	password string
	notes    string
}

func initEditState(e vault.Entry) editState {
	s := editState{
		keys: editKeys,
		val: &editValues{
			title:    e.Title,
			username: e.Username,
			password: e.Password,
			notes:    e.Notes,
		},
		origval: editValues{
			title:    e.Title,
			username: e.Username,
			password: e.Password,
			notes:    e.Notes,
		},
		id: e.ID,
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

func (m model) editUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.state.edit.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.state.edit.form = f
	}

	if m.state.edit.form.State == huh.StateCompleted {
		opts := make([]vault.EntryOption, 0)

		v := m.state.edit.val
		ov := m.state.edit.origval

		if v.title != ov.title {
			opts = append(opts, vault.WithTitle(v.title))
		}
		if v.username != ov.username {
			opts = append(opts, vault.WithUsername(v.username))
		}
		if v.password != ov.password {
			opts = append(opts, vault.WithPassword(v.password))
		}
		if v.notes != ov.notes {
			opts = append(opts, vault.WithNotes(v.notes))
		}

		if len(opts) > 0 {
			m.vault.EditEntry(m.state.edit.id, opts...)
		}

		return m, m.switchList()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.state.edit.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.state.edit.keys.Back):
			return m, m.switchList()
		case key.Matches(msg, m.state.edit.keys.Quit):
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) editView() string {
	// coolBorder := lipgloss.Border{
	// 	Top:         "/",
	// 	Bottom:      "/",
	// 	Left:        "",
	// 	Right:       "",
	// 	TopLeft:     "/",
	// 	TopRight:    "/",
	// 	BottomLeft:  "/",
	// 	BottomRight: "/",
	// }
	// b1Style := lipgloss.NewStyle().
	// 	Padding(1, 2).BorderStyle(coolBorder).
	// 	BorderForeground(lipgloss.Color("#7e98e8")).
	// 	Margin(1, 1)

	var b1 strings.Builder

	b1.WriteString("adding entry\n\n")
	b1.WriteString(m.state.edit.form.View() + "\n\n")

	combined := cEditKeymap{
		k: m.state.edit.keys,
		f: m.state.edit.form,
	}
	b1.WriteString(m.help.View(combined))

	out := b1.String()

	return out
}
