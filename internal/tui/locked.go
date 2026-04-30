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

type lockedKeymap struct {
	SeePass key.Binding
	Quit    key.Binding
}

type cLockedKeymap struct {
	k lockedKeymap
	f *huh.Form
}

func (c cLockedKeymap) ShortHelp() []key.Binding {
	binds := c.f.KeyBinds()
	return append(binds, c.k.SeePass, c.k.Quit)
}

func (c cLockedKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		c.f.KeyBinds(),
		{c.k.SeePass, c.k.Quit},
	}
}

var lockedKeys = lockedKeymap{
	SeePass: key.NewBinding(
		key.WithKeys("f2"),
		key.WithHelp("f2", "peek passwd"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type lockedState struct {
	keys     lockedKeymap
	form     *huh.Form
	inputPtr *huh.Input
	passwd   string
	path     string
	peekPw   bool
}

func (s *lockedState) rebuildForm() tea.Cmd {
	s.passwd = ""
	s.peekPw = false

	in := huh.NewInput().
		Placeholder("enter password").
		Value(&s.passwd).
		EchoMode(huh.EchoModePassword)

	s.form = huh.NewForm(huh.NewGroup(in)).WithShowHelp(false)
	s.inputPtr = in

	return s.form.Init()
}

func initLockedState(path string) *lockedState {
	s := &lockedState{
		keys: lockedKeys,
		path: path,
	}

	s.rebuildForm()

	return s
}

func (m model) lockedUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.state.locked.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.state.locked.form = f
	}

	if m.state.locked.form.State == huh.StateCompleted {
		v, err := vault.Open(m.state.locked.path, []byte(m.state.locked.passwd))
		if err != nil {
			m.errMsg = err
			m.state.locked.rebuildForm()
			return m, tea.Batch(cmd, clearErrDelayed())
		}

		m.vault = v
		m.errMsg = nil
		return m, m.switchList()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.state.locked.keys.SeePass):
			if m.state.locked.peekPw {
				m.state.locked.inputPtr.EchoMode(huh.EchoModePassword)
			} else {
				m.state.locked.inputPtr.EchoMode(huh.EchoModeNormal)
			}
			m.state.locked.peekPw = !m.state.locked.peekPw
		
		// we do not add toggle help function here as short and full help
		// report the same keys
		case key.Matches(msg, m.state.locked.keys.Quit):
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) lockedView() string {
	b1Style := lipgloss.NewStyle().Padding(1, 2)
	errBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#d8647e")).
		Padding(0, 1)

	var b1 strings.Builder

	fmt.Fprintf(&b1, "vault locked\n\n")
	b1.WriteString(m.state.locked.form.View() + "\n\n")

	if m.errMsg != nil {
		b1.WriteString("\n" + errBox.Render(m.errMsg.Error()) + "\n")
	}

	combined := cLockedKeymap{
		k: m.state.locked.keys,
		f: m.state.locked.form,
	}
	b1.WriteString("\n" + m.help.View(combined))

	out := b1Style.Render(b1.String())

	return out
}

