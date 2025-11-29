package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show id [-i path] [-j] [flags]",
	Short: "show a specific entry",
	RunE:  runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
	// default path flag
	showCmd.Flags().StringP("path", "i", ".ubpm/vault.ubpm.json", "where the vault is located (defaults to .ubpm/vault.ubpm.json)")
	// individual field flags
	showCmd.Flags().Bool("id", false, "show only id")
	showCmd.Flags().Bool("title", false, "show only title")
	showCmd.Flags().Bool("username", false, "show only username")
	showCmd.Flags().Bool("password", false, "show only password")
	showCmd.Flags().Bool("notes", false, "show only notes")
	showCmd.Flags().Bool("created", false, "show only creation date")
	showCmd.Flags().Bool("modified", false, "show only modification date")
	// json output flag
	showCmd.Flags().BoolP("json", "j", false, "output as json")
}

func runShow(cmd *cobra.Command, args []string) error {
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

	// find an entry in vault
	e, _, err := v.FindEntry(id)
	if err != nil {
		return err
	}

	// check which flags are set
	showID, _ := cmd.Flags().GetBool("id")
	showTitle, _ := cmd.Flags().GetBool("title")
	showUsername, _ := cmd.Flags().GetBool("username")
	showPassword, _ := cmd.Flags().GetBool("password")
	showNotes, _ := cmd.Flags().GetBool("notes")
	showCreated, _ := cmd.Flags().GetBool("created")
	showModified, _ := cmd.Flags().GetBool("modified")
	asJSON, _ := cmd.Flags().GetBool("json")

	// if no specific flags are set, show everything
	showAll := !showID && !showTitle && !showUsername && !showPassword &&
		!showNotes && !showCreated && !showModified

	if asJSON {
		// build json output based on flags
		output := make(map[string]interface{})

		if showAll || showID {
			output["id"] = e.ID
		}
		if showAll || showTitle {
			output["title"] = e.Title
		}
		if showAll || showUsername {
			output["username"] = e.Username
		}
		if showAll || showPassword {
			output["password"] = e.Password
		}
		if showAll || showNotes {
			output["notes"] = e.Notes
		}
		if showAll || showCreated {
			output["created_at"] = e.CreatedAt
		}
		if showAll || showModified {
			output["modified_at"] = e.ModifiedAt
		}

		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	}

	// regular text output
	if showAll {
		fmt.Printf("showing entry %s\n", e.ID)
		fmt.Printf("title: %s\n", e.Title)
		fmt.Printf("username: %s\n", e.Username)
		fmt.Printf("password: %s\n", e.Password)
		fmt.Println("---------")
		fmt.Printf("notes:\n%s\n", e.Notes)
		fmt.Println("---------")
		fmt.Printf("created %s\n", e.CreatedAt.Format("on Mon, Jan 2, 2006 at 15:04:05 MST"))
		fmt.Printf("last modified %s\n", e.ModifiedAt.Format("on Mon, Jan 2, 2006 at 15:04:05 MST"))
	} else {
		// show only requested fields
		if showID {
			fmt.Println(e.ID)
		}
		if showTitle {
			fmt.Println(e.Title)
		}
		if showUsername {
			fmt.Println(e.Username)
		}
		if showPassword {
			fmt.Println(e.Password)
		}
		if showNotes {
			fmt.Println(e.Notes)
		}
		if showCreated {
			fmt.Println(e.CreatedAt.Format(time.RFC3339))
		}
		if showModified {
			fmt.Println(e.ModifiedAt.Format(time.RFC3339))
		}
	}

	return nil
}
