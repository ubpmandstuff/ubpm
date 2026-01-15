package tui

import (
	// "fmt"

	"dura5ka/ubpm/internal/vault"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
)

// ::::: types :::::

// model is the central model of the bubble tea app
type model struct {
	cursor int
	view   string // "list"|"add"|"edit"|"rm"
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
	list listState
	add  addState
	// edit editState
	// rm   rmState
}

// ::::: view switch funcs :::::

func (m *model) switchList() {
	m.view = "list"
	m.help.ShowAll = false
	m.state.list = initListState()
}

func (m *model) switchAdd() {
	m.view = "add"
	m.help.ShowAll = false
	m.state.add = initAddState()
}

// ::::: utils :::::

func (m *model) isViewportGood() bool {
	return m.vp.width >= 40 && m.vp.height >= 10
}

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
	default:
		return m, nil
	}
}

func (m model) View() string {
	if !m.isViewportGood() {
		return fmt.Sprintf("given viewport is not large enough for comfortable use.\n\nviewport is %dx%d.\n\npress q to quit", m.vp.width, m.vp.height)
	}
	switch m.view {
	case "list":
		return m.listView()
	default:
		return m.listView()
	}
}
