/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/chand1012/arail/pkg/chat"
	"github.com/chand1012/arail/pkg/db"
)

// Does not work.

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Have a chat with Arail",
	Long:  `Have a chat with Arail.`,
	Run: func(cmd *cobra.Command, args []string) {
		tempDB, err := db.NewTemp()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		var prompt string

		fmt.Println("Arail> Hello! I'm Arail, your friendly neighborhood AI. What would you like to ask me?")
		for {
			fmt.Print("You> ")
			fmt.Scanln(&prompt)

			if prompt == "/exit" {
				break
			}

			response, err := chat.Chat(prompt, tempDB)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			fmt.Println("Arail>", response)
		}
		fmt.Println("Arail> Thanks for chatting with me! Bye!")
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
