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
		database, err := db.New()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

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

		texts, err := research.Research(query, database)

		if err != nil {
			log.Error(err)
			os.Exit(1)
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
