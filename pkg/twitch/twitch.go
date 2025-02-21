package twitch

import (
	"context"
	"fmt"

	"github.com/adeithe/go-twitch/api"
)

var ErrNotFound = fmt.Errorf("stream not found")

type Twitch struct {
	client       *api.Client
	clientSecret string
	clientID     string
	token        *Token
}

func New(clientID, clientSecret string) *Twitch {
	return &Twitch{
		client:       api.New(clientID),
		clientSecret: clientSecret,
		clientID:     clientID,
	}
}

func (t *Twitch) StreamByName(ctx context.Context, username string) (*api.Stream, error) {
	token, err := t.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	streams, err := t.client.Streams.List().Username([]string{username}).Do(ctx, api.WithBearerToken(token.AccessToken))
	if err != nil {
		return nil, fmt.Errorf("get streams: %w", err)
	}

	if len(streams.Data) == 0 {
		return nil, ErrNotFound
	}

	return &streams.Data[0], nil
}

func (t *Twitch) getToken(ctx context.Context) (*Token, error) {
	if t.token.IsValid() {
		return t.token, nil
	}

	token, err := t.token.Renew(ctx, t.clientID, t.clientSecret)
	if err != nil {
		return nil, fmt.Errorf("get new token: %w", err)
	}

	t.token = token

	return t.token, nil
}
