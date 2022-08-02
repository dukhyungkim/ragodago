package cmddown

import (
	"fmt"
	"harago/gservice/gchat"
	"harago/stream"
	"harago/util"
	"strings"

	pbAct "github.com/dukhyungkim/libharago/gen/go/proto/action"
	"github.com/jessevdk/go-flags"
	"google.golang.org/api/chat/v1"
)

type CmdDown struct {
	name         string
	streamClient *stream.Client
}

func NewDownCommand(streamClient *stream.Client) *CmdDown {
	return &CmdDown{
		name:         "/down",
		streamClient: streamClient,
	}
}

func (c *CmdDown) GetName() string {
	return c.name
}

type Opts struct {
	Company string `long:"company" short:"c"`
}

func (c *CmdDown) Run(event *gchat.ChatEvent) *chat.Message {
	fields := strings.Fields(event.Message.Text)
	if fields == nil {
		return c.Help()
	}

	var opts Opts
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)

	args, err := parser.ParseArgs(fields[1:])
	if err != nil {
		return &chat.Message{Text: err.Error()}
	}

	if len(args) == 0 {
		return &chat.Message{Text: "invalid ResourceURL"}
	}
	resourceURL := args[0]

	subject := stream.SharedSubject
	if opts.Company != "" {
		subject = fmt.Sprintf(stream.SpecificCompanySubject, opts.Company)
	}

	pbAction := &pbAct.ActionRequest{
		Type:  pbAct.ActionType_DOWN,
		Space: event.Space.Name,
		Request_OneOf: &pbAct.ActionRequest_ReqDeploy{
			ReqDeploy: &pbAct.ActionRequest_DeployRequest{
				Name:        util.ParseName(resourceURL),
				ResourceUrl: resourceURL,
			},
		},
	}
	if err = c.streamClient.PublishAction(subject, pbAction); err != nil {
		return &chat.Message{Text: err.Error()}
	}

	if subject == stream.SharedSubject {
		return &chat.Message{Text: fmt.Sprintf("publish to %s, ResourceURL: %s", subject, resourceURL)}
	}
	return &chat.Message{Text: fmt.Sprintf("publish to %s, Company: %s, ResourceURL: %s", subject, opts.Company, resourceURL)}
}

func (c *CmdDown) Help() *chat.Message {
	return &chat.Message{Text: "HELP!"}
}
