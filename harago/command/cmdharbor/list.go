package cmdharbor

import (
	"fmt"
	"harago/common"
	"harago/util"

	"github.com/dukhyungkim/harbor-client"
	harborModel "github.com/dukhyungkim/harbor-client/model"
	"google.golang.org/api/chat/v1"
)

func (c *CmdHarbor) handleList(opts *SubCmdOpts) *chat.Message {
	if opts.RepoName != "" {
		return listArtifacts(c.harborClient, opts)
	}

	if opts.ProjectName != "" {
		return listRepositories(c.harborClient, opts)
	}

	return listProjects(c.harborClient, opts)
}

func listProjects(client *harbor.Client, opts *SubCmdOpts) *chat.Message {
	projectsParams := harborModel.NewListProjectsParams()
	if opts.Page != 0 {
		projectsParams.Page = opts.Page
	}
	projectsParams.PageSize = 15
	if opts.Size != 0 {
		projectsParams.PageSize = opts.Size
	}

	projects, err := client.ListProjects(projectsParams)
	if err != nil {
		return &chat.Message{Text: common.ErrHarborResponse(err).Error()}
	}

	cards := make([]*chat.Card, len(projects))
	for i := range projects {
		cards[i] = makeProjectCard(projects[i])
	}
	return &chat.Message{Text: "list of projects", Cards: cards}
}

func listRepositories(client *harbor.Client, opts *SubCmdOpts) *chat.Message {
	repositoriesParams := harborModel.NewListRepositoriesParams()
	if opts.Page != 0 {
		repositoriesParams.Page = opts.Page
	}
	repositoriesParams.PageSize = 15
	if opts.Size != 0 {
		repositoriesParams.PageSize = opts.Size
	}

	repositories, err := client.ListRepositories(opts.ProjectName, repositoriesParams)
	if err != nil {
		return &chat.Message{Text: common.ErrHarborResponse(err).Error()}
	}

	cards := make([]*chat.Card, len(repositories))
	for i := range repositories {
		cards[i] = makeRepositoryCard(repositories[i])
	}
	return &chat.Message{Text: fmt.Sprintf("list of repositories in %s", opts.ProjectName), Cards: cards}
}

func listArtifacts(client *harbor.Client, opts *SubCmdOpts) *chat.Message {
	artifactsParams := harborModel.NewListArtifactsParams()
	if opts.Page != 0 {
		artifactsParams.Page = opts.Page
	}
	artifactsParams.PageSize = 15
	if opts.Size != 0 {
		artifactsParams.PageSize = opts.Size
	}

	artifacts, err := client.ListArtifacts(opts.ProjectName, opts.RepoName, artifactsParams)
	if err != nil {
		return &chat.Message{Text: common.ErrHarborResponse(err).Error()}
	}

	cards := make([]*chat.Card, len(artifacts))
	for i := range artifacts {
		tags, err := client.ListTags(opts.ProjectName, opts.RepoName, artifacts[i].Digest, nil)
		if err != nil {
			return &chat.Message{Text: common.ErrHarborResponse(err).Error()}
		}
		cards[i] = makeArtifactCard(artifacts[i], tags)
	}

	return &chat.Message{Text: fmt.Sprintf("list of artifacts in %s/%s", opts.ProjectName, opts.RepoName), Cards: cards}
}

func makeProjectCard(project *harborModel.Project) *chat.Card {
	return &chat.Card{
		Header: &chat.CardHeader{
			Title: project.Name,
		},
		Sections: []*chat.Section{
			{
				Widgets: []*chat.WidgetMarkup{
					{
						KeyValue: &chat.KeyValue{
							TopLabel:         "RepoCount",
							Content:          fmt.Sprint(project.RepoCount),
							ContentMultiline: true,
						},
					},
					{
						KeyValue: &chat.KeyValue{
							TopLabel:         "OwnerName",
							Content:          project.OwnerName,
							ContentMultiline: true,
						},
					},
					{
						KeyValue: &chat.KeyValue{
							TopLabel:         "UpdateTime",
							Content:          project.UpdateTime.Local().String(),
							ContentMultiline: true,
						},
					},
				},
			},
		},
	}
}

func makeRepositoryCard(repository *harborModel.Repository) *chat.Card {
	return &chat.Card{
		Header: &chat.CardHeader{
			Title: repository.Name,
		},
		Sections: []*chat.Section{
			{
				Widgets: []*chat.WidgetMarkup{
					{
						KeyValue: &chat.KeyValue{
							TopLabel:         "ArtifactCount",
							Content:          fmt.Sprint(repository.ArtifactCount),
							ContentMultiline: true,
						},
					},
					{
						KeyValue: &chat.KeyValue{
							TopLabel:         "PullCount",
							Content:          fmt.Sprint(repository.PullCount),
							ContentMultiline: true,
						},
					},
					{
						KeyValue: &chat.KeyValue{
							TopLabel:         "UpdateTime",
							Content:          repository.UpdateTime.Local().String(),
							ContentMultiline: true,
						},
					},
				},
			},
		},
	}
}

func makeArtifactCard(artifact *harborModel.Artifact, tags []*harborModel.Tag) *chat.Card {
	return &chat.Card{
		Header: &chat.CardHeader{
			Title: artifact.Digest[:15],
		},
		Sections: []*chat.Section{
			{
				Widgets: []*chat.WidgetMarkup{
					{
						KeyValue: &chat.KeyValue{
							TopLabel:         "Tags",
							Content:          tags[0].Name,
							ContentMultiline: true,
						},
					},
					{
						KeyValue: &chat.KeyValue{
							TopLabel:         "Size",
							Content:          util.ByteCountIEC(int64(artifact.Size)),
							ContentMultiline: true,
						},
					},
					{
						KeyValue: &chat.KeyValue{
							TopLabel:         "PushTime",
							Content:          artifact.PushTime.Local().String(),
							ContentMultiline: true,
						},
					},
				},
			},
		},
	}
}
