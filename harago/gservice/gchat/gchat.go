package gchat

import (
	"context"
	"harago/common"
	"harago/gservice"
	"log"
	"time"

	"google.golang.org/api/chat/v1"
	"google.golang.org/api/option"
)

type GChat struct {
	service     *chat.Service
	dmHandler   Handler
	roomHandler Handler
}

func NewGChat(gService *gservice.GService, dmHandler, roomHandler Handler) (*GChat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultTimeout)
	defer cancel()

	service, err := chat.NewService(ctx, option.WithHTTPClient(gService.GetClient()))
	if err != nil {
		return nil, err
	}

	return &GChat{service: service, dmHandler: dmHandler, roomHandler: roomHandler}, nil
}

func (c *GChat) HandleMessage(event *ChatEvent) *chat.Message {
	log.Println(event)
	start := time.Now()

	var chatMessage *chat.Message
	if event.Space.Type == DM {
		chatMessage = c.dmHandler.ProcessMessage(event)
	} else {
		chatMessage = c.roomHandler.ProcessMessage(event)
	}

	log.Printf("elapsed: %v", time.Since(start))
	return chatMessage
}

func (c *GChat) SendMessage(space string, message *chat.Message) {
	if _, err := c.service.Spaces.Messages.Create(space, message).Do(); err != nil {
		log.Printf("failed send message: %v\n", err)
	}
}
