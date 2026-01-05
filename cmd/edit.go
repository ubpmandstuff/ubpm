package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"dura5ka/ubpm/internal/vault"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var editCmd = &cobra.Command{
	Use:   "edit id [-i path] [-tupn]",
	Short: "edit an entry",
	RunE:  runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
	// default path flag
	editCmd.Flags().StringP("path", "i", ".ubpm/vault.ubpm.json", "where the vault is located")
	// command-specific flags
	editCmd.Flags().BoolP("title", "t", false, "change only title")
	editCmd.Flags().BoolP("username", "u", false, "change only username")
	editCmd.Flags().BoolP("password", "p", false, "change only password")
	editCmd.Flags().BoolP("notes", "n", false, "change only notes")
}

func runEdit(cmd *cobra.Command, args []string) error {
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

	// get id (assume only arg is id)
	id := ""
	if len(args) > 0 {
		id = args[0]
	} else {
		return fmt.Errorf("no id provided")
	}

	// check if entry exists
	e, _, err := v.FindEntry(id)
	if err != nil {
		return err
	}

	// announce action
	fmt.Fprintf(os.Stderr, "editing entry %s\n", e.ID[8:])

	// check what needs to be edited
	changeTitle, _ := cmd.Flags().GetBool("title")
	changeUsername, _ := cmd.Flags().GetBool("username")
	changePassword, _ := cmd.Flags().GetBool("password")
	changeNotes, _ := cmd.Flags().GetBool("notes")
	var opts []vault.EntryOption

	// prompt user
	reader := bufio.NewReader(os.Stdin)
	if changeTitle {
		fmt.Fprint(os.Stderr, "enter new title: ")
		title, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		title = strings.TrimSpace(title)
		opts = append(opts, vault.WithTitle(title))
	}
	if changeUsername {
		fmt.Fprint(os.Stderr, "enter new username: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		username = strings.TrimSpace(username)
		opts = append(opts, vault.WithUsername(username))
	}
	if changePassword {
		fmt.Fprint(os.Stderr, "enter new password (will not echo): ")
		passB, err := term.ReadPassword(int(os.Stdin.Fd()))
		pass := string(passB)
		if err != nil {
			return err
		}
		pass = strings.TrimSpace(pass)
		opts = append(opts, vault.WithPassword(pass))
	}
	if changeNotes {
		fmt.Fprint(os.Stderr, "enter new notes: ")
		notes, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		notes = strings.TrimSpace(notes)
		opts = append(opts, vault.WithNotes(notes))
	}

	// check if there is anything to edit
	if len(opts) == 0 {
		return fmt.Errorf("nothing to edit")
	}

	// edit entry and save it
	if err := v.EditEntry(id, opts...); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "entry edited successfully")

	return nil
}
