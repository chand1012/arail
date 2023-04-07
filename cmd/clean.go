/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the database",
	Long:  `Clean the database by deleting the database file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// ask for a confirmation
		// if yes, delete the database
		if !f {
			fmt.Println("Are you sure you want to delete the database?")
			fmt.Print("This action cannot be undone. [y/N]: ")
			var response string
			fmt.Scanln(&response)
			if response != "y" {
				fmt.Println("Aborting...")
				os.Exit(0)
			}
		}
		currentUser, err := user.Current()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		homeDir := currentUser.HomeDir
		dbPath := filepath.Join(homeDir, ".arail", "arail.db")

		err = os.Remove(dbPath)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		fmt.Println("Database deleted.")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.Flags().BoolVarP(&f, "force", "f", false, "Force clean the database with no prompt.")
}
