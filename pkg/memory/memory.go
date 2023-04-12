package memory

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve"
)

type Memory struct {
	Index bleve.Index
}

type MemoryFragment struct {
	Content string
	ID      string
	Score   float64
}

func New() (*Memory, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	homeDir := currentUser.HomeDir
	configDir := filepath.Join(homeDir, ".arail")
	indexPath := filepath.Join(configDir, "index.bleve")
	// attempt to open the index
	index, err := bleve.Open(indexPath)
	if err != nil {
		// create a mapping
		mapping := bleve.NewIndexMapping()
		// create the index
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, err
		}
	}

	m := &Memory{
		Index: index,
	}

	return m, nil
}

func (m *Memory) Close() error {
	return m.Index.Close()
}

func (m *Memory) Add(id string, content string) error {
	return m.Index.Index(id, strings.ToLower(content))
}

func (m *Memory) Search(query string) ([]MemoryFragment, error) {
	q := bleve.NewQueryStringQuery(query)
	search := bleve.NewSearchRequest(q)
	searchResults, err := m.Index.Search(search)
	if err != nil {
		return nil, err
	}

	// tokenize and search again
	tokens := SimpleTokenize(query)
	for _, token := range tokens {
		q = bleve.NewQueryStringQuery(token)
		search = bleve.NewSearchRequest(q)
		r, err := m.Index.Search(search)
		if err != nil {
			return nil, err
		}
		searchResults.Hits = append(searchResults.Hits, r.Hits...)
	}

	var results []MemoryFragment
	for _, hit := range searchResults.Hits {
		var result MemoryFragment
		result.Content = hit.Fields["Content"].(string)
		result.ID = hit.ID
		result.Score = hit.Score
		results = append(results, result)
	}
	return results, nil
}

func TopFragment(fragments []MemoryFragment) MemoryFragment {
	var top MemoryFragment
	for _, fragment := range fragments {
		if fragment.Score > top.Score {
			top = fragment
		}
	}
	return top
}
