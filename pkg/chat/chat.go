package chat

import (
	"github.com/CasualCodersProjects/gopenai/types"
	"github.com/chand1012/arail/pkg/ai"
	"github.com/chand1012/arail/pkg/converser"
	"github.com/chand1012/arail/pkg/db"
	"github.com/chand1012/arail/pkg/db/models"
)

var (
	CHAT_SYSTEM_PROMPT = "You are a helpful assistant."
)

func Chat(prompt string, database *db.Database) (string, error) {

	chats, err := converser.Converse(CHAT_SYSTEM_PROMPT, prompt, database)
	if err != nil {
		return "", err
	}

	o, _ := ai.InitOpenAI()

	// create a default chat object
	req := types.NewDefaultChatRequest("")
	req.Messages = chats
	req.MaxTokens = 4000

	resp, err := o.CreateChat(req)

	if err != nil {
		return "", err
	}

	// get the response from the message as a ChatMessage
	c := resp.Choices[0].Message
	message := models.ChatFromChatMessage(c)

	// add the new chat to the database
	err = database.PostChat(message)

	return message.Content, err
}
