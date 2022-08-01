package gchat

import "google.golang.org/api/chat/v1"

type Handler interface {
	ProcessMessage(event *ChatEvent) *chat.Message
}
