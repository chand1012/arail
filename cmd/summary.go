/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/chand1012/arail/pkg/ai"
	"github.com/chand1012/arail/pkg/db"
	"github.com/chand1012/arail/pkg/db/models"
	"github.com/chand1012/arail/pkg/pages"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Summarize a given webpage",
	Long: `Summarize a given webpage. Takes a URL as a positional argument.
Example: arail summary https://en.wikipedia.org/wiki/Go_(programming_language)	
`,
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		database, err := db.New()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		s, err := database.GetSummaryByURL(url)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Error(err)
				os.Exit(1)
			}
		}

		if s.Summary != "" {
			fmt.Println("Summary of " + url + ":")
			fmt.Println(s.Summary)
			os.Exit(0)
		}

		log.Info("Getting page data...")

		p, err := database.GetTextByURL(url)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Error(err)
				os.Exit(1)
			}
		}

		chunks := []string{}

		if len(p) == 0 {
			page, err := pages.ExtractPageData(url)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			chunks = ai.ChunkSite(page)
		} else {
			for _, chunk := range p {
				chunks = append(chunks, chunk.Text)
			}
		}
		// post the site to the db
		for i, chunk := range chunks {
			c := models.SiteChunk{
				Text:      chunk,
				TextIndex: i,
				URL:       url,
			}
			err = database.PostSite(c)
			if err != nil {
				log.Error(err)
			}
		}

		log.Info("Summarizing...")
		summary, err := ai.SummarizeSite(chunks, "What is this page about?")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		fmt.Println("Summary of " + url + ":")
		fmt.Println(summary)

		// post the summary to the db
		s = models.Summary{
			URL:     url,
			Summary: summary,
		}

		err = database.PostSummary(s)
		if err != nil {
			log.Error(err)
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(summaryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// summaryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// summaryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
