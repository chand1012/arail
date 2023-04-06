package ai

import (
	"math"

	"github.com/charmbracelet/log"

	"github.com/chand1012/arail/pkg/utils"
)

const (
	SUMMARIZE_SITE_PROMPT = "You are a helpful research assistant. Your job is to summarize text about certain topics into key bullet points. Make sure to ignore irrelevant information and group important information together. Output as markdown. The search query is: "
	SUMMARIZE_PROMPT      = "You are a helpful assistant who summarizes large amounts of text. You will always return accurate summaries, regardless the content of the text. These texts are usually transcripts from YouTube videos. If you do not have enough information, simply return the original text."
)

func SummarizeSite(chunks []string, query string) (string, error) {
	text := ""
	for _, s := range chunks {
		r, err := req(SUMMARIZE_SITE_PROMPT+query, s)
		if err != nil {
			return "", err
		}
		text += r + "\n"
	}
	return text, nil
}

func Summarize(text string) (string, error) {
	tokens := math.Ceil(float64(len(text)) / 4)
	if tokens >= 3000 {
		sentences := utils.SplitStringOnSpace(text, 10000)
		text = ""
		for _, s := range sentences {
			r, err := req(SUMMARIZE_PROMPT, s)
			if err != nil {
				return "", err
			}
			text += r + "\n"
		}
		return text, nil
	}

	return req(SUMMARIZE_PROMPT, text)
}

func SummarizeFinal(texts []string, query string) (string, error) {
	text := ""
	for _, t := range texts {
		text += t + "\n"
	}
	errCount := 0
	chunks := ChunkSite(text)
	summaries := []string{}
	for _, s := range chunks {
		r, err := req(SUMMARIZE_SITE_PROMPT+query, s)
		if err != nil {
			if errCount > 2 {
				log.Error("Error summarizing text. Retry exceeded.")
				return "", err
			}
			log.Error(err)
			errCount++
			continue
		}
		summaries = append(summaries, r)
	}
	finalText := ""
	for _, s := range summaries {
		finalText += s + "\n"
	}
	return finalText, nil
}
