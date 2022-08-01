package handler

import (
	"harago/repository"
	"harago/stream"
	"log"

	harborModel "github.com/dukhyungkim/harbor-client/model"
	pbAct "github.com/dukhyungkim/libharago/gen/go/proto/action"
)

type subjectMapper struct {
	f       func(string) bool
	subject string
}

type HarborEventHandler struct {
	streamClient *stream.Client
	etcdClient   *repository.Etcd
	mappers      []subjectMapper
}

func NewHarborEventHandler(streamClient *stream.Client, etcdClient *repository.Etcd) *HarborEventHandler {
	fss := []subjectMapper{
		{
			f:       etcdClient.IsShared,
			subject: stream.SharedSubject,
		},
		{
			f:       etcdClient.IsCompany,
			subject: stream.CompanySubject,
		},
		{
			f:       etcdClient.IsInternal,
			subject: stream.InternalSubject,
		},
		{
			f:       etcdClient.IsExternal,
			subject: stream.ExternalSubject,
		},
	}

	return &HarborEventHandler{
		streamClient: streamClient,
		etcdClient:   etcdClient,
		mappers:      fss,
	}
}

func (h *HarborEventHandler) HandleHarborEvent(event *harborModel.WebhookEvent) {
	name := event.EventData.Repository.Name
	request := &pbAct.ActionRequest{
		Type: pbAct.ActionType_UP,
		Request_OneOf: &pbAct.ActionRequest_ReqDeploy{
			ReqDeploy: &pbAct.ActionRequest_DeployRequest{
				Name:        name,
				ResourceUrl: event.EventData.Resources[0].ResourceURL,
			},
		},
	}
	log.Println("pbAction:", request.String())

	if h.etcdClient.IsIgnore(name) {
		log.Printf("%s is in ignoredList\n", name)
		return
	}

	for _, mapper := range h.mappers {
		if !mapper.f(name) {
			continue
		}

		if err := h.streamClient.PublishAction(mapper.subject, request); err != nil {
			log.Println(err)
		}
		log.Printf("sent action to subject: %s\n", mapper.subject)
	}
}
