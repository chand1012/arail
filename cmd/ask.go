/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/chand1012/arail/pkg/ai"
	"github.com/chand1012/arail/pkg/db"
	"github.com/chand1012/arail/pkg/research"
)

// add code to create and search text from sites
// all text from all sites gets added and then searched
// this way we can ask questions about the topic and get answers

// this is broken. See pkg/db/db.go

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		q := args[0]
		database, err := db.New()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Info("Parsing query...")
		queries, err := ai.ParseQuery(q)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		_, err = research.Research(queries[0], database)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Info("Searching database...")
		texts, err := database.SearchSummarySlice(queries)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		prompt := ""
		for _, text := range texts {
			if len(prompt) >= 10000 {
				break
			}
			prompt += text.Summary + "\n"
		}

		prompt += "\nQ: " + q + "\nA:"

		log.Info("Asking...")
		answer, err := ai.Ask(prompt)

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		fmt.Println("Answer: ")
		fmt.Println(answer)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(askCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// askCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// askCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
