package tui

import (
	"dura5ka/ubpm/internal/vault"

	"github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor int
	vault  vault.Vault
}
