package pages

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	"gorm.io/gorm"

	"github.com/chand1012/arail/pkg/ai"
	"github.com/chand1012/arail/pkg/db"
	"github.com/chand1012/arail/pkg/db/models"
)

func ProcessURL(url, query string, i int, content chan<- string, wg *sync.WaitGroup, database *db.Database) {
	defer wg.Done()

	s, err := database.GetSummaryByURL(url)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error(err)
			content <- ""
			return
		}
	}

	if s.Summary != "" {
		log.Info("Summary for " + url + " already exists in the database, skipping web request and summarize...")
		content <- s.Summary
		return
	}

	existingTexts, err := database.GetTextByURL(url)
	if err != nil {
		log.Error("Error fetching existing chunks from database:", err)
		content <- ""
		return
	}

	var chunks []string
	if len(existingTexts) > 0 {
		log.Info("Chunks for " + url + " already exist in the database, skipping web request...")
		for _, text := range existingTexts {
			chunks = append(chunks, text.Text)
		}
	} else {
		log.Info("Getting page data for " + url + " on thread " + fmt.Sprint(i) + "...")
		resp, err := ExtractPageData(url)
		if err != nil {
			log.Error(err)
			content <- ""
			return
		}
		log.Info("Summarizing page data for " + url + " on thread " + fmt.Sprint(i) + "...")
		chunks = ai.ChunkSite(resp)
		for index, chunk := range chunks {
			siteChunk := models.SiteChunk{
				Text:      chunk,
				URL:       url,
				TextIndex: index,
			}
			err := database.PostSite(siteChunk)
			if err != nil {
				log.Error("Error saving chunk to database:", err)
			}
		}
	}

	resp, err := ai.SummarizeSite(chunks, query)
	if err != nil {
		log.Error(err)
		content <- ""
		return
	}
	s = models.Summary{
		URL:     url,
		Summary: resp,
	}
	err = database.PostSummary(s)
	if err != nil {
		log.Error("Error saving summary to database:", err)
	}
	content <- resp
	log.Info("Finished with " + url + "...")
}
