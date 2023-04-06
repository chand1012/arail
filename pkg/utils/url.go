package utils

import "github.com/chand1012/yt_transcript"

func ShortToFullYouTubeURL(url string) (string, error) {
	id, err := yt_transcript.GetVideoID(url)
	if err != nil {
		return "", err
	}

	return "https://www.youtube.com/watch?v=" + id, nil
}
