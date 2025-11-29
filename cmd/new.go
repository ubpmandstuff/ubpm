package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var newCmd = &cobra.Command{
	Use:     "new [-i path]",
	Short:   "add an entry to your vault",
	Aliases: []string{"add"},
	RunE:    runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
	// default path flag
	newCmd.Flags().StringP("path", "i", ".ubpm/vault.ubpm.json", "where the vault is located (defaults to .ubpm/vault.ubpm.json)")
}

func runNew(cmd *cobra.Command, args []string) error {
	// define path
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}

	// load the vault
	v, err := loadVault(path)
	if err != nil {
		return err
	}

	// announce to user
	fmt.Println("creating new entry")

	// ask user for data
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stderr, "enter entry title: ")
	title, _ := reader.ReadString('\n')
	if len(title) == 0 {
		return fmt.Errorf("title cannot be empty")
	}
	title = strings.TrimSpace(title)

	fmt.Fprint(os.Stderr, "enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Fprint(os.Stderr, "enter password (will not echo): ")
	passB, _ := term.ReadPassword(int(os.Stdin.Fd()))
	pass := string(passB)
	fmt.Println()

	fmt.Fprint(os.Stderr, "notes (optional): ")
	notes, _ := reader.ReadString('\n')
	notes = strings.TrimSpace(notes)

	// add entry to vault
	if err := v.AddEntry(title, username, pass, notes); err != nil {
		return err
	}

	return nil
}
