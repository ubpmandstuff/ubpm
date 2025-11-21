/*
Copyright © 2025 dura5ka
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"dura5ka/ubpm/internal/vault"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// listCmd represents
var listCmd = &cobra.Command{
	Use:   "list path",
	Short: "list all items in vault",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	// define path
	path := ".ubpm/vault.ubpm.json"
	if len(args) > 0 {
		path = args[0]
	}

	// prompt user for password and decrypt vault
	fmt.Print("enter password: ")
	pass, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return err
	}

	v, err := vault.Open(path, pass)
	if err != nil {
		return err
	}

	if len(v.Data.Entries) == 0 {
		fmt.Println("vault is empty")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tWEBSITE\tUSERNAME\tCREATED")
	for _, e := range v.Data.Entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.ID[:8],
			e.Website,
			e.Username,
			e.CreatedAt.Format("2006-01-02"),
		)
	}
	w.Flush()

	return nil
}
