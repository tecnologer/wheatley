package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/adeithe/go-twitch/api"
	"github.com/tecnologer/wheatley/pkg/utils/log"
)

var ErrNotFound = fmt.Errorf("stream not found")

type API interface {
	StreamByName(ctx context.Context, username string) (*api.Stream, error)
	Token(ctx context.Context) (*Token, error)
	RenewToken(request *Request) (*Token, error)
}

type Streamer interface {
	List() *api.StreamsListCall
}

type Config struct {
	ClientID     string
	ClientSecret string
	IsMock       bool
}

type Twitch struct {
	*Config
	Streams Streamer
	token   *Token
}

func New(config *Config) API {
	if config.IsMock {
		return &MockAPI{}
	}

	return &Twitch{
		Streams: api.New(config.ClientID).Streams,
		Config:  config,
	}
}

func (t *Twitch) StreamByName(ctx context.Context, username string) (*api.Stream, error) {
	token, err := t.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	streams, err := t.Streams.List().Username([]string{username}).Do(ctx, api.WithBearerToken(token.AccessToken))
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

	token, err := t.RenewToken(&Request{
		HTTPClient:   t.httpClient(),
		Context:      ctx,
		ClientID:     t.ClientID,
		ClientSecret: t.ClientSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("get new token: %w", err)
	}

	t.token = token

	return t.token, nil
}

// MockRoundTripper implements http.RoundTripper to mock HTTP responses.
type MockRoundTripper struct {
	Response string
	Status   int
	Err      error
}

// RoundTrip mocks an HTTP request.
func (m *MockRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, fmt.Errorf("do request: %w", m.Err)
	}

	return &http.Response{
		StatusCode: m.Status,
		Body:       io.NopCloser(bytes.NewBufferString(m.Response)),
		Header:     make(http.Header),
	}, nil
}

func (t *Twitch) httpClient() *http.Client {
	if t.IsMock {
		return &http.Client{
			Transport: &MockRoundTripper{},
		}
	}

	return http.DefaultClient
}

type Request struct {
	HTTPClient   *http.Client
	Context      context.Context //nolint: containedctx // This is a context.Context that will be used to make requests to Twitch API
	ClientID     string
	ClientSecret string
}

func (t *Twitch) RenewToken(request *Request) (*Token, error) {
	request.ClientID = strings.TrimSpace(request.ClientID)
	request.ClientSecret = strings.TrimSpace(request.ClientSecret)

	req, err := http.NewRequestWithContext(request.Context, http.MethodPost, authTokenURL, formURLEncoded(request.ClientID, request.ClientSecret))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	res, err := request.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer func() {
		err := res.Body.Close()
		log.Warnf("error closing response body: %v", err)
	}()

	var token Token

	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if token.HasError() {
		return nil, fmt.Errorf("error response: Code: %d - %s", token.Status, token.Message)
	}

	token.lastUpdate = time.Now()

	return &token, nil
}
