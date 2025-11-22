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
	v, err := loadVault(path)
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
