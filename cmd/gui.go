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

func init() {
	rootCmd.AddCommand(guiCmd)
	// default path flag
	guiCmd.Flags().StringP("path", "i", ".ubpm/vault.ubpm.json", "where the vault is located (defaults to .ubpm/vault.ubpm.json)")
}

func runGui(cmd *cobra.Command, args []string) error {
	// init path variable
	_, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}

	// prompt user for password and decrypt vault
	// commented out, to be implemented in the gui
	// _, err = loadVault(path)
	// if err != nil {
	// 	return err
	// }

	gui.MakeWindow()
	return nil
}
