package cmdharbor

import (
	"harago/gservice/gchat"
	"strings"

	"github.com/dukhyungkim/harbor-client"
	"github.com/jessevdk/go-flags"
	"google.golang.org/api/chat/v1"
)

type CmdHarbor struct {
	name         string
	harborClient *harbor.Client
}

func NewHarborCommand(harborClient *harbor.Client) *CmdHarbor {
	return &CmdHarbor{
		name:         "/harbor",
		harborClient: harborClient,
	}
}

func (c *CmdHarbor) GetName() string {
	return c.name
}

type Opts struct {
	Info SubCmdOpts `command:"info"`
	List SubCmdOpts `command:"list" alias:"ls"`
}

type SubCmdOpts struct {
	ProjectName  string `long:"project" alias:"proj"`
	RepoName     string `long:"repository" alias:"repo"`
	ArtifactName string `long:"artifact"`
	Page         int64  `long:"page"`
	Size         int64  `long:"size"`
}

const (
	subCmdInfo = "info"
	subCmdList = "list"
)

func (c *CmdHarbor) Run(event *gchat.ChatEvent) *chat.Message {
	fields := strings.Fields(event.Message.Text)
	if fields == nil {
		return c.Help()
	}

	var opts Opts
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)

	_, err := parser.ParseArgs(fields[1:])
	if err != nil {
		return &chat.Message{Text: err.Error()}
	}

	switch parser.Active.Name {
	case subCmdList:
		return c.handleList(&opts.List)
	case subCmdInfo:
		return c.handleInfo(&opts.Info)
	default:
		return c.Help()
	}
}

func (c *CmdHarbor) Help() *chat.Message {
	return &chat.Message{Text: "HELP!"}
}
