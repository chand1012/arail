package memory

import (
	"regexp"
	"strings"
)

func SimpleTokenize(sentence string) []string {
	// Define a lookup table of common words that are not part of the subject
	nonSubjectWords := map[string]bool{
		"a":      true,
		"an":     true,
		"the":    true,
		"is":     true,
		"are":    true,
		"was":    true,
		"were":   true,
		"am":     true,
		"be":     true,
		"been":   true,
		"being":  true,
		"has":    true,
		"have":   true,
		"had":    true,
		"do":     true,
		"does":   true,
		"did":    true,
		"may":    true,
		"might":  true,
		"must":   true,
		"shall":  true,
		"should": true,
		"will":   true,
		"would":  true,
	}

	// Define a regular expression that matches the subject of a sentence
	re := regexp.MustCompile(`^[A-Z][^\.?!]*`)

	// Find the subject of the sentence using the regular expression
	match := re.FindString(sentence)

	// If the regular expression didn't match, return an empty list
	if match == "" {
		return []string{}
	}

	// Split the subject into individual words
	words := strings.Fields(match)

	// Remove any filler words from the list of words
	var filteredWords []string
	for _, word := range words {
		if !nonSubjectWords[strings.ToLower(word)] {
			filteredWords = append(filteredWords, word)
		}
	}

	// Return the list of filtered words
	return filteredWords
}
