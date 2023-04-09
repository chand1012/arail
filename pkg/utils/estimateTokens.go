package utils

import (
	"math"
	"regexp"
)

func removeWhitespace(input string) string {
	re := regexp.MustCompile(`[\t\n\f\r\v]+`)
	return re.ReplaceAllString(input, " ")
}

func EstimateTokens(text string) int {
	// remove all non space whitespace
	text = removeWhitespace(text)
	// get the length as a float and divide by 4
	tokens := float64(len(text)) / 4

	tokens = math.Ceil(tokens)

	// return as an int
	return int(tokens)
}
