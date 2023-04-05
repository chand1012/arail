/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	googlesearch "github.com/rocketlaunchr/google-search"
	"github.com/spf13/cobra"

	"github.com/chand1012/arail/pkg/ai"
	"github.com/chand1012/arail/pkg/pages"
	"github.com/chand1012/arail/pkg/utils"
)

var (
	outFile   string
	notesMode bool
	f         bool
)

// researchCmd represents the research command
var researchCmd = &cobra.Command{
	Use:   "research",
	Short: "Search, summarize, and report a topic",
	Long: `Uses Google and OpenAI to create notes or a report on a topic. Requires a topic as a positional argument.
Example: arail research "How to make a website"

Search queries are recommended to be in quotes to avoid errors.
`,
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		if outFile != "" {
			// check if file exists.
			if _, err := os.Stat(outFile); err == nil {
				if f {
					// remove the file
					log.Warn("Removing file '" + outFile + "'...")
					err := os.Remove(outFile)
					if err != nil {
						log.Error(err)
						os.Exit(1)
					}
				} else {
					// If it exists then exit
					log.Error("File '" + outFile + "' already exists. Please rename or remove it.")
					os.Exit(1)
				}
			}
		}

		results, err := googlesearch.Search(context.TODO(), query)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// get all urls. Make sure there are no duplicates
		var urls []string
		for _, result := range results {
			if !utils.Contains(urls, result.URL) {
				urls = append(urls, result.URL)
			}
		}

		var wg sync.WaitGroup

		content := make(chan string, len(urls))

		for i, url := range urls {
			wg.Add(1)
			go func(url string, i int) {
				defer wg.Done()
				log.Info("Getting page data for " + url + "on thread " + fmt.Sprint(i) + "...")
				resp, err := pages.ExtractPageData(url)
				if err != nil {
					log.Error(err)
					content <- ""
					return
				}
				log.Info("Summarizing page data for " + url + "on thread " + fmt.Sprint(i) + "...")
				resp, err = ai.SummarizeSite(resp, query)
				if err != nil {
					log.Error(err)
					content <- ""
					return
				}
				content <- resp
				log.Info("Finished with " + url + "...")
			}(url, i)
		}

		log.Info("Waiting for site summaries to finish...")
		wg.Wait()
		// log.Info("Closing channel...")
		close(content)

		// log.Info("Extracting content from channel...")
		var texts []string
		for t := range content {
			texts = append(texts, t)
		}

		log.Info("Summarizing final notes...")
		summary, err := ai.SummarizeFinal(texts, query)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if notesMode {
			if outFile == "" {
				fmt.Println(summary)
			} else {
				log.Info("Writing notes to " + outFile + "...")
				err = os.WriteFile(outFile, []byte(summary), 0644)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
				log.Info("Done!")
			}
		} else {
			log.Info("Generating report...")
			report, err := ai.Report(summary, query)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			if outFile == "" {
				fmt.Println(report)
			} else {
				log.Info("Writing report to " + outFile + "...")
				err = os.WriteFile(outFile, []byte(summary), 0644)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
				log.Info("Done!")
			}
		}
	},
	// requires 1 positional argument
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(researchCmd)

	researchCmd.Flags().StringVarP(&outFile, "out", "o", "", "Output file")
	researchCmd.Flags().BoolVarP(&notesMode, "notes", "n", false, "Do not generate a report. Only generate notes. Required if GPT-4 is unavailable")
	researchCmd.Flags().BoolVarP(&f, "force", "f", false, "Force overwrite of output file")
}
