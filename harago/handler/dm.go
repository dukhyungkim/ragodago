package handler

import (
	"harago/command"
	"harago/entity"
	"harago/gservice/gchat"
	"harago/repository"

	"google.golang.org/api/chat/v1"
)

type DMHandler struct {
	cmdExecutor *command.Executor
	repo        *repository.DB
}

func NewDMHandler(cmdExecutor *command.Executor, repo *repository.DB) gchat.Handler {
	return &DMHandler{cmdExecutor: cmdExecutor, repo: repo}
}

func (h *DMHandler) ProcessMessage(event *gchat.ChatEvent) *chat.Message {
	switch event.Type {
	case gchat.AddedToSpace:
		userSpace := &entity.UserSpace{
			Name:  event.User.DisplayName,
			Email: event.User.Email,
			Space: event.Space.Name,
		}
		if err := h.repo.SaveSpace(userSpace); err != nil {
			return &chat.Message{Text: err.Error()}
		}
		return &chat.Message{Text: "Save Space"}

	case gchat.Message:
		return h.cmdExecutor.Run(event)

	case gchat.RemovedFromSpace:
		h.repo.DeleteSpace(event.User.Email)
	}

	return &chat.Message{}
}
