package cmd

import (
	"dura5ka/ubpm/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui [-i path]",
	Short: "open the app's tui",
	RunE:  runTui,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	// default path flag
	tuiCmd.Flags().StringP("path", "i", ".ubpm/vault.ubpm.json", "where the vault is located")
}

func runTui(cmd *cobra.Command, args []string) error {
	// init path variable
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}

	// prompt user for password and decrypt vault
	v, err := loadVault(path)
	if err != nil {
		return err
	}

	// create a new program and run it
	p := tea.NewProgram(tui.InitialModel(v), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
