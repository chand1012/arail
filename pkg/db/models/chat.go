package models

import (
	"github.com/CasualCodersProjects/gopenai/types"
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Role    string
	Content string
}

func ChatFromChatMessage(m types.ChatMessage) *Chat {
	return &Chat{
		Role:    m.Role,
		Content: m.Content,
	}
}

func ChatsFromChatMessages(ms []types.ChatMessage) []*Chat {
	var chats []*Chat
	for _, m := range ms {
		chats = append(chats, ChatFromChatMessage(m))
	}
	return chats
}

func (c *Chat) ToMessage() types.ChatMessage {
	return types.ChatMessage{
		Role:    c.Role,
		Content: c.Content,
	}
}

func ChatsToMessages(cs []*Chat) []types.ChatMessage {
	var ms []types.ChatMessage
	for _, c := range cs {
		ms = append(ms, c.ToMessage())
	}
	return ms
}
