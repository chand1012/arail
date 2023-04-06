package ai

import "github.com/chand1012/arail/pkg/utils"

var (
	MAX_CHARS   = 12000
	SPLIT_CHARS = 10000
)

func ChunkSite(text string) []string {
	if len(text) <= MAX_CHARS {
		return []string{text}
	}
	sentences := utils.SplitStringOnSpace(text, SPLIT_CHARS)
	return sentences
}
