package repository

import (
	"context"
	"fmt"
	"harago/common"
	"harago/config"
	"log"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	sharedListKey   = "/config/shared"
	companyListKey  = "/config/company"
	internalListKey = "/config/internal"
	externalListKey = "/config/external"
	ignoreListKey   = "/config/.ignore"
)

type Etcd struct {
	client       *clientv3.Client
	sharedList   map[string]struct{}
	companyList  map[string]struct{}
	internalList map[string]struct{}
	ignoreList   map[string]struct{}
	externalList map[string]struct{}
}

func NewEtcd(cfg *config.Etcd) (*Etcd, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: common.DefaultTimeout,
		Username:    cfg.Username,
		Password:    cfg.Password,
	})
	if err != nil {
		return nil, common.ErrConnEtcd(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultTimeout)
	defer cancel()

	sharedList, err := fetchKeyAndParse(ctx, etcdClient, sharedListKey)
	if err != nil {
		return nil, err
	}

	companyList, err := fetchKeyAndParse(ctx, etcdClient, companyListKey)
	if err != nil {
		return nil, err
	}

	internalList, err := fetchKeyAndParse(ctx, etcdClient, internalListKey)
	if err != nil {
		return nil, err
	}

	externalList, err := fetchKeyAndParse(ctx, etcdClient, externalListKey)
	if err != nil {
		return nil, err
	}

	ignoreList, err := fetchKeyAndParse(ctx, etcdClient, ignoreListKey)
	if err != nil {
		return nil, err
	}

	return &Etcd{
		client:       etcdClient,
		sharedList:   sharedList,
		companyList:  companyList,
		internalList: internalList,
		externalList: externalList,
		ignoreList:   ignoreList,
	}, nil
}

func fetchKeyAndParse(ctx context.Context, etcdClient *clientv3.Client, key string) (map[string]struct{}, error) {
	resp, err := etcdClient.Get(ctx, key)
	if err != nil {
		log.Println(fmt.Errorf("failed to get kv; %w", err))
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		log.Println(fmt.Errorf("failed to find value from key: %s", key))
	}

	return parseListToMap(string(resp.Kvs[0].Value)), nil
}

func (e *Etcd) Close() {
	if err := e.client.Close(); err != nil {
		log.Println(err)
	}
}

func (e *Etcd) IsShared(name string) bool {
	_, has := e.sharedList[name]
	return has
}

func (e *Etcd) IsCompany(name string) bool {
	_, has := e.companyList[name]
	return has
}

func (e *Etcd) IsInternal(name string) bool {
	_, has := e.internalList[name]
	return has
}

func (e *Etcd) IsExternal(name string) bool {
	_, has := e.externalList[name]
	return has
}

func (e *Etcd) IsIgnore(name string) bool {
	_, has := e.ignoreList[name]
	return has
}

func (e *Etcd) WatchConfigList() {
	ctx := context.Background()
	sharedListChan := e.client.Watch(ctx, sharedListKey)
	companyListChan := e.client.Watch(ctx, companyListKey)
	internalListChan := e.client.Watch(ctx, internalListKey)
	externalListChan := e.client.Watch(ctx, externalListKey)
	ignoreListChan := e.client.Watch(ctx, ignoreListKey)

	for {
		select {
		case watchResp := <-sharedListChan:
			if len(watchResp.Events) == 0 {
				continue
			}
			e.sharedList = parseListToMap(string(watchResp.Events[0].Kv.Value))
			log.Println("sharedList:", e.sharedList)

		case watchResp := <-companyListChan:
			if len(watchResp.Events) == 0 {
				continue
			}
			e.companyList = parseListToMap(string(watchResp.Events[0].Kv.Value))
			log.Println("companyList", e.companyList)

		case watchResp := <-internalListChan:
			if len(watchResp.Events) == 0 {
				continue
			}
			e.internalList = parseListToMap(string(watchResp.Events[0].Kv.Value))
			log.Println("internalList", e.internalList)

		case watchResp := <-externalListChan:
			if len(watchResp.Events) == 0 {
				continue
			}
			e.externalList = parseListToMap(string(watchResp.Events[0].Kv.Value))
			log.Println("external", e.externalList)

		case watchResp := <-ignoreListChan:
			if len(watchResp.Events) == 0 {
				continue
			}
			e.ignoreList = parseListToMap(string(watchResp.Events[0].Kv.Value))
			log.Println("ignoreList", e.ignoreList)
		}
	}
}

const templatePrefix = "/templates/"

func (e *Etcd) ListTemplates() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultTimeout)
	defer cancel()

	resp, err := e.client.Get(ctx, templatePrefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("failed to find value from key: %s", templatePrefix)
	}

	templates := make([]string, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		templates[i] = strings.TrimPrefix(string(kv.Key), templatePrefix)
	}

	return templates, nil
}

func (e *Etcd) GetTemplate(name string) (string, error) {
	key := templatePrefix + name

	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultTimeout)
	defer cancel()

	resp, err := e.client.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("failed to find value from key: %s", key)
	}

	return string(resp.Kvs[0].Value), nil
}

func parseListToMap(s string) map[string]struct{} {
	fields := strings.Fields(s)
	sharedList := make(map[string]struct{})
	for _, field := range fields {
		sharedList[field] = struct{}{}
	}
	return sharedList
}
