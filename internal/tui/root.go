/* Package tui provides a terminal-based user interface for ubpm */
package tui

import (
	"fmt"
	"strings"
	"time"

	"dura5ka/ubpm/internal/vault"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ::::: types :::::

// minWidth is minimum width for app to function properly
const minWidth int = 40

// minHeight is minimum height for app to function properly
const minHeight int = 10

// model is the central model of the bubble tea app
type model struct {
	cursor int
	view   string // "list"|"add"|"edit"|"rm"|"locked"
	vault  *vault.Vault
	help   help.Model
	vp     viewport
	state  state
	errMsg error
}

// viewport is a struct storing the width and height of the user's viewport,
// that usually being the terminal
type viewport struct {
	height int
	width  int
}

type state struct {
	list   listState
	view   viewState
	add    addState
	edit   editState
	rm     *rmState
	locked *lockedState
}

// ::::: view switch funcs :::::

func (m *model) switchList() tea.Cmd {
	m.view = "list"
	m.help.ShowAll = false
	m.state.list = initListState()
	return nil
}

func (m *model) switchView(e vault.Entry) tea.Cmd {
	m.view = "view"
	m.help.ShowAll = false
	m.state.view = initViewState(e)
	return nil
}

func (m *model) switchAdd() tea.Cmd {
	m.view = "add"
	m.help.ShowAll = false
	m.state.add = initAddState()
	return m.state.add.form.Init()
}

func (m *model) switchEdit(e vault.Entry) tea.Cmd {
	m.view = "edit"
	m.help.ShowAll = false
	m.state.edit = initEditState(e)
	return m.state.edit.form.Init()
}

func (m *model) switchRm(e vault.Entry) tea.Cmd {
	m.view = "rm"
	m.help.ShowAll = false
	m.state.rm = initRmState(e)
	return m.state.rm.form.Init()
}

func (m *model) switchLocked(path string) tea.Cmd {
	m.view = "locked"
	m.help.ShowAll = false
	m.state.locked = initLockedState(path)
	m.vault.Wipe()
	m.vault = nil
	return tea.Sequence(clearErr, m.state.locked.form.Init())
}

// ::::: utils :::::

// isViewportGood checks whether the viewport is larger or equal to
// minimum width and height
func (m *model) isViewportGood() bool {
	return m.vp.width >= minWidth && m.vp.height >= minHeight
}

// clearErrMsg is a bubbletea message clearing the error
type clearErrMsg struct{}

func clearErrDelayed() tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return clearErrMsg{}
	})
}

func clearErr() tea.Msg {
	return clearErrMsg{}
}

const headerBanner = "██ UBPM ████▓▒░"

// brandStyle defines the default brand style
var brandStyle lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7e98e8"))

// errorStyle defines the default error text style
var errorStyle lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#d8647e"))

// boxStyle is the default outer padding applied to every screen in View
var boxStyle lipgloss.Style = lipgloss.NewStyle().Padding(1, 2).MarginRight(2)

// noSuchView returns the default view shown when one switches to a
// non-existent view (intentionally or due to an error)
func noSuchView() string {
	var b1 strings.Builder
	b1.WriteString("hey! you're trying to access a non-existent view.\n")
	b1.WriteString("this must have happened due to some kind of error.\n")
	b1.WriteString("\npress esc to load back into the list view.")
	return b1.String()
}

// noSuchViewUpd provides keybinding handling for the noSuchView view
func (m model) noSuchViewUpd(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, m.switchList()
		}
	}
	return m, nil
}

// ::::: the app itself :::::

// InitialModel initializes the model for launching the app
func InitialModel(v *vault.Vault) model {
	h := help.New()
	return model{
		view:   "list",
		vault:  v,
		help:   h,
		errMsg: nil,
		state: state{
			list: initListState(),
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.vp.width = msg.Width
		m.vp.height = msg.Height
	case clearErrMsg:
		m.errMsg = nil
	}
	if !m.isViewportGood() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
		return m, nil
	}
	switch m.view {
	case "list":
		return m.listUpdate(msg)
	case "view":
		return m.viewUpdate(msg)
	case "add":
		return m.addUpdate(msg)
	case "edit":
		return m.editUpdate(msg)
	case "rm":
		return m.rmUpdate(msg)
	case "locked":
		return m.lockedUpdate(msg)
	default:
		return m.noSuchViewUpd(msg)
	}
}

func (m model) View() string {
	if !m.isViewportGood() {
		return fmt.Sprintf(
			"viewport is too small for the application\n\nviewport is %dx%d\n\nrecommended is %dx%d\n\npress q to quit",
			m.vp.width, m.vp.height, minWidth, minHeight,
		)
	}

	var out string
	banner := brandStyle
	switch m.view {
	case "list":
		out = m.listView()
	case "view":
		out = m.viewView()
	case "add":
		out = m.addView()
	case "edit":
		out = m.editView()
	case "rm":
		out = m.rmView()
	case "locked":
		out = m.lockedView()
	default:
		out = noSuchView()
		banner = errorStyle
	}
	if m.errMsg != nil {
		banner = errorStyle
	}

	return boxStyle.Render(banner.Render(headerBanner) + "\n\n" + out)
}
