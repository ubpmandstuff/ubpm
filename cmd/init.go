/*
Copyright © 2025 dura5ka
*/
package cmd

import (
	"fmt"
	"os"

	"dura5ka/ubpm/internal/vault"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [-i path]",
	Short: "initialize a ubpm vault",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	// default path flag
	initCmd.Flags().StringP("path", "i", ".", "where to create vault")
}

func runInit(cmd *cobra.Command, args []string) error {
	// init path variable
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path must be a directory")
	}

	log.Info("initializing vault", "path", path)

	fmt.Fprint(os.Stderr, "enter master password: ")
	pass1, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if len(pass1) == 0 {
		return fmt.Errorf("password cannot be empty")
	}
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stderr, "confirm password: ")
	pass2, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return err
	}

	if string(pass1) != string(pass2) {
		return fmt.Errorf("passwords do not match")
	}

	vaultPath, err := vault.Init(path, pass1)
	if err != nil {
		return err
	}

	fmt.Printf("vault created at: %s\n", vaultPath)

	return nil
}
