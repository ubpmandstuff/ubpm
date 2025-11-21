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
	Use:   "init path",
	Short: "initialize a ubpm vault",
	Args:  cobra.MaximumNArgs(1),
	RunE:  InitRun,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func InitRun(cmd *cobra.Command, args []string) error {
	// init path variable
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path must be a directory")
	}

	log.Info("initializing vault", "path", path)

	fmt.Print("enter master password: ")
	pass1, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if len(pass1) == 0 {
		return fmt.Errorf("password cannot be empty")
	}
	if err != nil {
		return err
	}

	fmt.Print("confirm password: ")
	pass2, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return err
	}

	if string(pass1) != string(pass2) {
		return fmt.Errorf("passwords do not match")
	}

	vaultPath, err := vault.InitVault(path, pass1)
	if err != nil {
		return err
	}

	fmt.Printf("vault created at: %s\n", vaultPath)

	return nil
}
