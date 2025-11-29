/*
Copyright © 2025 dura5ka
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// listCmd represents
var listCmd = &cobra.Command{
	Use:   "list [-i path]",
	Short: "list all items in vault",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	// default path flag
	listCmd.Flags().StringP("path", "i", ".ubpm/vault.ubpm.json", "where the vault is located (defaults to .ubpm/vault.ubpm.json)")
}

func runList(cmd *cobra.Command, args []string) error {
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

	if len(v.Data.Entries) == 0 {
		fmt.Println("vault is empty")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tUSERNAME\tCREATED")
	for _, e := range v.Data.Entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.ID[:8],
			e.Title,
			e.Username,
			e.CreatedAt.Format("2006-01-02"),
		)
	}
	w.Flush()

	return nil
}
