package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm id [-i path] [--noconfirm]",
	Short: "remove an entry from a ubpm vault",
	RunE:  runRm,
}

func init() {
	rootCmd.AddCommand(rmCmd)
	// default path flag
	rmCmd.Flags().StringP("path", "i", ".ubpm/vault.ubpm.json", "where the vault is located (defaults to .ubpm/vault.ubpm.json)")
	// command-specific flags
	rmCmd.Flags().Bool("noconfirm", false, "skip confirmation")
}

func runRm(cmd *cobra.Command, args []string) error {
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

	// find and remove entry
	e, _, err := v.FindEntry(id)
	if err != nil {
		return err
	}

	value, err := cmd.Flags().GetBool("noconfirm")
	if err == nil && value == true {
		err := v.RemoveEntry(id)
		if err != nil {
			return err
		}
	} else if err == nil {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("are you sure you'd like to delete entry %s? [y/N]: ", e.ID)
		ans, err := reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		fmt.Println(ans)
		if err != nil {
			return err
		}
		if strings.ToLower(ans) == "y" {
			err := v.RemoveEntry(id)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("operation cancelled")
		}
	}

	return nil
}
