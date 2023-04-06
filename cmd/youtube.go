/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/chand1012/yt_transcript"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/chand1012/arail/pkg/ai"
	"github.com/chand1012/arail/pkg/db"
	"github.com/chand1012/arail/pkg/db/models"
	"github.com/chand1012/arail/pkg/utils"
)

// youtubeCmd represents the youtube command
var youtubeCmd = &cobra.Command{
	Use:   "youtube",
	Short: "Summarize a YouTube video",
	Long: `Uses YouTube transcripts and OpenAI to summarize a YouTube video.
Takes a YouTube URL as a parameter.`,
	Aliases: []string{"yt"},
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		database, err := db.New()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		videoId, err := yt_transcript.GetVideoID(url)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		url, err = utils.ShortToFullYouTubeURL(url)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// check if the video summary is in the DB
		s, err := database.GetSummaryByURL(url)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Error(err)
				os.Exit(1)
			}
		}

		if s.Summary != "" {
			fmt.Println(s.Summary)
			os.Exit(0)
		}

		// check if the video transcript is in the DB
		v, err := database.GetVideo(videoId)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Error(err)
				os.Exit(1)
			}
		}
		var text string
		if v.Transcript == "" {
			// check and see if the video transcript is in the DB
			log.Info("Getting transcript...")
			transcript, title, err := yt_transcript.FetchTranscript(videoId, "en", "US")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			log.Info("Summarizing transcript for '" + title + "'...")
			text = ""
			for _, t := range transcript {
				text += t.Text + " "
			}

			err = database.PostVideo(models.Video{
				Title:      title,
				VideoID:    videoId,
				Transcript: text,
			})

			if err != nil {
				log.Error(err)
			}
		} else {
			log.Info("Summarizing transcript for '" + v.Title + "'...")
			text = v.Transcript
		}

		summary, err := ai.Summarize(text)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		s = models.Summary{
			URL:     url,
			Summary: summary,
		}

		err = database.PostSummary(s)

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		fmt.Println(summary)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(youtubeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// youtubeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// youtubeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
