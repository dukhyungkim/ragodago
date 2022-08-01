package gservice

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"harago/common"
	"net/http"
	"os"
)

type GService struct {
	client *http.Client
}

func NewGService(credential string) (*GService, error) {
	b, err := os.ReadFile(credential)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	const ChatScope = "https://www.googleapis.com/auth/chat.bot"
	var scope = []string{ChatScope}
	configFromJSON, err := google.JWTConfigFromJSON(b, scope...)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to configFromJSON: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultTimeout)
	defer cancel()

	httpClient := configFromJSON.Client(ctx)
	return &GService{client: httpClient}, nil
}

func (svc *GService) GetClient() *http.Client {
	return svc.client
}
