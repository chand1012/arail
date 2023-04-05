package utils

import "strings"

func SplitStringOnSpace(text string, maxTokens int) []string {
	words := strings.Fields(text)
	var sentences []string

	currentPrompt := ""
	currentTokens := 0

	for _, word := range words {
		wordTokens := len(word) + 1
		if currentTokens+wordTokens > maxTokens {
			sentences = append(sentences, currentPrompt)
			currentPrompt = ""
			currentTokens = 0
		}

		currentPrompt += word + " "
		currentTokens += wordTokens
	}

	if currentPrompt != "" {
		sentences = append(sentences, currentPrompt)
	}

	return sentences
}
