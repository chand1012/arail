/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
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

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask a question about a topic and get an answer",
	Long: `Ask a question about a topic, and arail will search for relevant information and return an answer.
Example: arail ask "What is the capital of France?"`,
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
