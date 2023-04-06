package ai

import (
	"errors"
	"strings"
)

var (
	ParseQueryPrompt = "What is the query about? Ignore prepositions and articles and just return the main subject."
)

func ParseQuery(query string) ([]string, error) {
	o, _ := initOpenAI()
	resp, err := o.CreateChatSimple(ParseQueryPrompt+"\n"+query, 256)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) < 1 {
		return nil, errors.New("no response from OpenAI")
	}

	queries := []string{resp.Choices[0].Message.Content}
	// split the string on spaces and append
	// each word to the queries slice
	words := strings.Split(resp.Choices[0].Message.Content, " ")
	queries = append(queries, words...)

	return queries, nil
}
