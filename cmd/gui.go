/*
Copyright © 2025 dura5ka
*/
package cmd

import (
	gui "dura5ka/ubpm/internal/GUI"

	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui [-i path]",
	Short: "initialize a ubpm vault",
	RunE:  runGui,
}

func runGui(cmd *cobra.Command, args []string) error {
	// init path variable
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}

	// prompt user for password and decrypt vault
	_, err = loadVault(path)
	if err != nil {
		return err
	}

	window := gui.MakeWindow()
	window.ShowAndRun()
	return nil
}
