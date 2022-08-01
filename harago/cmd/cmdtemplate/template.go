package cmdtemplate

import (
	"bytes"
	"harago/gservice/gchat"
	"harago/repository"
	"html/template"
	"strings"

	"github.com/jessevdk/go-flags"
	"google.golang.org/api/chat/v1"
)

type CmdTemplate struct {
	name       string
	etcdClient *repository.Etcd
}

func NewTemplateCommand(etcdClient *repository.Etcd) *CmdTemplate {
	return &CmdTemplate{
		name:       "/template",
		etcdClient: etcdClient,
	}
}

func (c *CmdTemplate) GetName() string {
	return c.name
}

type ShowOpts struct {
	Base        string `long:"base" default:"{{.Base}}"`
	Company     string `long:"company" default:"{{.Company}}"`
	Name        string `long:"name" default:"{{.Name}}"`
	ResourceURL string `long:"resource_url" default:"{{.ResourceURL}}"`
	Host        string `long:"host" default:"{{.Host}}"`
}

type Opts struct {
	List struct{} `command:"list" alias:"ls"`
	Show ShowOpts `command:"show"`
}

const (
	subCmdList = "list"
	subCmdShow = "show"
)

func (c *CmdTemplate) Run(event *gchat.ChatEvent) *chat.Message {
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

	switch parser.Active.Name {
	case subCmdList:
		return c.handleList()
	case subCmdShow:
		if len(args) == 0 {
			return &chat.Message{Text: "need name"}
		}
		return c.handleShow(args[0], &opts.Show)
	default:
		return c.Help()
	}
}

func (c *CmdTemplate) Help() *chat.Message {
	return &chat.Message{Text: "HELP!"}
}

func (c *CmdTemplate) handleList() *chat.Message {
	templates, err := c.etcdClient.ListTemplates()
	if err != nil {
		return &chat.Message{Text: err.Error()}
	}

	return &chat.Message{Text: strings.Join(templates, "\n")}
}

func (c *CmdTemplate) handleShow(templateName string, opts *ShowOpts) *chat.Message {
	templateStr, err := c.etcdClient.GetTemplate(templateName)
	if err != nil {
		return &chat.Message{Text: err.Error()}
	}

	tmpl, err := template.New(templateName).Parse(templateStr)
	if err != nil {
		return &chat.Message{Text: err.Error()}
	}

	var tmplBuffer bytes.Buffer
	err = tmpl.Execute(&tmplBuffer, opts)
	if err != nil {
		return &chat.Message{Text: err.Error()}
	}

	return &chat.Message{Text: "```" + tmplBuffer.String() + "```"}
}
