package converser

import (
	aiTypes "github.com/CasualCodersProjects/gopenai/types"
	"gorm.io/gorm"

	"github.com/chand1012/arail/pkg/db"
	"github.com/chand1012/arail/pkg/db/models"
	"github.com/chand1012/arail/pkg/utils"
)

// A basic converser
// adds as much context as possible to the conversation

func Converse(systemPrompt, userPrompt string, database *db.Database) ([]aiTypes.ChatMessage, error) {
	var chats []models.Chat
	var total int

	// add the first chat, the system prompt
	chats = append(chats, models.Chat{
		Role:    "system",
		Content: systemPrompt,
		Tokens:  utils.EstimateTokens(systemPrompt),
	})

	// start with the number of tokens in the system prompt and user prompt
	total = chats[0].Tokens + utils.EstimateTokens(userPrompt)
	offset := 0
	for total < 2000 {
		chat, err := database.GetChat(offset)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				break
			}
			return nil, err
		}
		if chat.Tokens == 0 {
			break
		}
		chats = append(chats, chat)
		total += chat.Tokens
		offset++
	}

	finalChat := models.Chat{
		Role:    "user",
		Content: userPrompt,
		Tokens:  utils.EstimateTokens(userPrompt),
	}

	// add the last chat, the user prompt
	chats = append(chats, finalChat)

	// add the last chat, the user prompt to the DB
	err := database.PostChat(finalChat)

	return models.ChatsToMessages(chats), err
}
