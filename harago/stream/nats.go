package stream

import (
	"harago/config"
	"log"
	"strings"
	"time"

	pbAct "github.com/dukhyungkim/libharago/gen/go/proto/action"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const (
	responseSubject = "handago.response"

	SharedSubject          = "harago.shared.action"
	CompanySubject         = "harago.company.action"
	InternalSubject        = "harago.internal.action"
	ExternalSubject        = "harago.external.action"
	SpecificCompanySubject = "harago.%s.action"
)

type Client struct {
	nc      *nats.Conn
	timeout time.Duration
}

func NewClient(cfg *config.Nats) (*Client, error) {
	nc, err := nats.Connect(strings.Join(cfg.Servers, ","),
		nats.UserInfo(cfg.Username, cfg.Password))
	if err != nil {
		return nil, err
	}

	return &Client{nc: nc, timeout: cfg.Timeout}, nil
}

func (s *Client) Close() {
	s.nc.Close()
}

func (s *Client) PublishAction(subject string, request *pbAct.ActionRequest) error {
	msg, err := proto.Marshal(request)
	if err != nil {
		return err
	}

	err = s.nc.Publish(subject, msg)
	if err != nil {
		return err
	}
	return nil
}

type ResponseHandler func(message *pbAct.ActionResponse)

func (s *Client) ClamResponse(handler ResponseHandler) error {
	if _, err := s.nc.QueueSubscribe(responseSubject, "harago", func(msg *nats.Msg) {
		var message pbAct.ActionResponse
		if err := proto.Unmarshal(msg.Data, &message); err != nil {
			log.Println(err)
			return
		}
		handler(&message)
	}); err != nil {
		return err
	}
	return nil
}
