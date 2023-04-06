package research

import (
	"context"
	"sync"

	"github.com/charmbracelet/log"
	googlesearch "github.com/rocketlaunchr/google-search"

	"github.com/chand1012/arail/pkg/db"
	"github.com/chand1012/arail/pkg/pages"
	"github.com/chand1012/arail/pkg/utils"
)

func Research(query string, database *db.Database) ([]string, error) {
	results, err := googlesearch.Search(context.TODO(), query)
	if err != nil {
		return nil, err
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

	// Use a separate goroutine to read from the content channel concurrently
	var texts []string
	go func() {
		for t := range content {
			texts = append(texts, t)
		}
	}()

	for i, url := range urls {
		wg.Add(1)
		go pages.ProcessURL(url, query, i, content, &wg, database)
	}

	log.Info("Waiting for site summaries to finish...")
	wg.Wait()
	// log.Info("Closing channel...")
	close(content)

	return texts, nil
}
