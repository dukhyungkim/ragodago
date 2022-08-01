package handler

import (
	"google.golang.org/api/chat/v1"
	"harago/cmd"
	"harago/gservice/gchat"
	"harago/repository"
)

type RoomHandler struct {
	cmdExecutor *cmd.Executor
	repo        *repository.DB
}

func NewRoomHandler(cmdExecutor *cmd.Executor, repo *repository.DB) gchat.Handler {
	return &RoomHandler{cmdExecutor: cmdExecutor, repo: repo}
}

func (h *RoomHandler) ProcessMessage(event *gchat.ChatEvent) *chat.Message {
	var chatMessage *chat.Message

	switch event.Type {
	case gchat.AddedToSpace:

	case gchat.Message:
		h.cmdExecutor.Run(event)

	case gchat.RemovedFromSpace:
		chatMessage = &chat.Message{}
	}

	return chatMessage
}
