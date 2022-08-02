package cmdping

import (
	"google.golang.org/api/chat/v1"
	"harago/gservice/gchat"
)

type CmdPing struct {
	name string
}

func NewDeployCommand() *CmdPing {
	return &CmdPing{name: "/ping"}
}

func (c *CmdPing) GetName() string {
	return c.name
}

func (c *CmdPing) Run(_ *gchat.ChatEvent) *chat.Message {
	return &chat.Message{Text: "pong"}
}

func (c *CmdPing) Help() *chat.Message {
	return &chat.Message{Text: "HELP!"}
}
