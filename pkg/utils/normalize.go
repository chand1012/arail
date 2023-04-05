package utils

import (
	"bytes"
	"regexp"
	"strings"
)

func NormalizeNewlines(s string) string {
	// Replace all sequences of newlines and carriage returns with a single newline
	var buf bytes.Buffer
	re := regexp.MustCompile(`[\r\n]+`)
	buf.WriteString(re.ReplaceAllString(s, "\n"))
	return buf.String()
}

func NormalizeSpaces(s string) string {
	// Replace all sequences of spaces with a single space
	var buf bytes.Buffer
	re := regexp.MustCompile(`[ ]+`)
	buf.WriteString(re.ReplaceAllString(s, " "))
	return buf.String()
}

func NormalizeParagraphs(s string) string {
	// Replace all sequences of newlines and spaces with two newlines
	var buf bytes.Buffer
	re := regexp.MustCompile(`[\r\n]+`)
	paragraphs := strings.Split(re.ReplaceAllString(s, "\n"), "\n")
	for i, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p != "" {
			if i > 0 {
				buf.WriteString("\n\n")
			}
			buf.WriteString(p)
		}
	}
	return buf.String()
}

func NormalizeString(s string) string {
	return NormalizeSpaces(NormalizeNewlines(s))
}
