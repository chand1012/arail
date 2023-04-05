package ai

import (
	"errors"
	"os"

	openai "github.com/CasualCodersProjects/gopenai"
	aiTypes "github.com/CasualCodersProjects/gopenai/types"

	"github.com/chand1012/arail/pkg/config"
)

func initOpenAI() (openai.OpenAI, string) {

	var key string
	conf, err := config.Load()
	if err != nil {
		if os.IsNotExist(err) {
			key = os.Getenv("OPENAI_API_KEY")
		} else {
			panic(err)
		}
	} else {
		key = conf.APIKey
	}

	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
	}

	if key == "" {
		panic("OPENAI_API_KEY not set")
	}

	o := openai.NewOpenAI(&openai.OpenAIOpts{
		APIKey: key,
	})

	return o, conf.Model
}

func req(sysPrompt, userPrompt string) (string, error) {
	o, model := initOpenAI()
	messages := []aiTypes.ChatMessage{
		{Role: "system", Content: sysPrompt},
		{Role: "user", Content: userPrompt},
	}
	req := aiTypes.NewDefaultChatRequest("")
	req.Messages = messages

	if model != "" {
		req.Model = model
	}

	resp, err := o.CreateChat(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) < 1 {
		return "", errors.New("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

func reqModelOverride(sysPrompt, userPrompt, model string) (string, error) {
	o, _ := initOpenAI()
	messages := []aiTypes.ChatMessage{
		{Role: "system", Content: sysPrompt},
		{Role: "user", Content: userPrompt},
	}
	req := aiTypes.NewDefaultChatRequest("")
	req.Messages = messages
	req.Model = model

	resp, err := o.CreateChat(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) < 1 {
		return "", errors.New("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}
