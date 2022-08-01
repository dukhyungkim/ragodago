package cmdharbor

import (
	"harago/common"

	"github.com/dukhyungkim/harbor-client"
	"google.golang.org/api/chat/v1"
)

func (c *CmdHarbor) handleInfo(opts *SubCmdOpts) *chat.Message {
	if opts.ProjectName != "" && opts.RepoName != "" && opts.ArtifactName != "" {
		return infoArtifact(c.harborClient, opts)
	}

	if opts.ProjectName != "" && opts.RepoName != "" {
		return infoRepository(c.harborClient, opts)
	}

	if opts.ProjectName != "" {
		return infoProject(c.harborClient, opts)
	}

	return c.Help()
}

func infoProject(client *harbor.Client, opts *SubCmdOpts) *chat.Message {
	project, err := client.GetProject(opts.ProjectName)
	if err != nil {
		return &chat.Message{Text: common.ErrHarborResponse(err).Error()}
	}

	return &chat.Message{Text: "project info", Cards: []*chat.Card{makeProjectCard(project)}}
}

func infoRepository(client *harbor.Client, opts *SubCmdOpts) *chat.Message {
	repository, err := client.GetRepository(opts.ProjectName, opts.RepoName)
	if err != nil {
		return &chat.Message{Text: common.ErrHarborResponse(err).Error()}
	}

	return &chat.Message{Text: "repository info", Cards: []*chat.Card{makeRepositoryCard(repository)}}
}

func infoArtifact(client *harbor.Client, opts *SubCmdOpts) *chat.Message {
	artifact, err := client.GetArtifact(opts.ProjectName, opts.RepoName, opts.ArtifactName)
	if err != nil {
		return &chat.Message{Text: common.ErrHarborResponse(err).Error()}
	}

	tags, err := client.ListTags(opts.ProjectName, opts.RepoName, artifact.Digest, nil)
	if err != nil {
		return &chat.Message{Text: common.ErrHarborResponse(err).Error()}
	}

	return &chat.Message{Text: "artifact info", Cards: []*chat.Card{makeArtifactCard(artifact, tags)}}
}
