package twitch

import (
	"context"
	"fmt"

	"github.com/adeithe/go-twitch/api"
)

var ErrNotFound = fmt.Errorf("stream not found")

type API interface {
	StreamByName(ctx context.Context, username string) (*api.Stream, error)
	Token(ctx context.Context) (*Token, error)
}

type Config struct {
	ClientID     string
	ClientSecret string
}

type Twitch struct {
	*Config
	client *api.Client
	token  *Token
}

func New(config *Config) *Twitch {
	return &Twitch{
		client: api.New(config.ClientID),
		Config: config,
	}
}

func (t *Twitch) StreamByName(ctx context.Context, username string) (*api.Stream, error) {
	token, err := t.Token(ctx)
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

func (t *Twitch) Token(ctx context.Context) (*Token, error) {
	if t.token.IsValid() {
		return t.token, nil
	}

	token, err := t.token.Renew(ctx, t.ClientID, t.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("get new token: %w", err)
	}

	t.token = token

	return t.token, nil
}
