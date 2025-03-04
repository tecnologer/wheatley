package twitch

import (
	"context"
	"fmt"

	"github.com/adeithe/go-twitch/api"
	"github.com/stretchr/testify/mock"
)

type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) StreamByName(ctx context.Context, username string) (*api.Stream, error) {
	args := m.Called(ctx, username)

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("getting stream: %w", err)
	}

	stream, isStream := args.Get(0).(*api.Stream)
	if !isStream {
		return nil, fmt.Errorf("invalid return type, expected *api.Stream, got %T", args.Get(0))
	}

	return stream, nil
}

func (m *MockAPI) Token(ctx context.Context) (*Token, error) {
	args := m.Called(ctx)

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("getting token: %w", err)
	}

	token, isToken := args.Get(0).(*Token)
	if !isToken {
		return nil, fmt.Errorf("invalid return type, expected *Token, got %T", args.Get(0))
	}

	return token, nil
}

func (m *MockAPI) RenewToken(request *Request) (*Token, error) {
	args := m.Called(request)

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("renewing token: %w", err)
	}

	token, isToken := args.Get(0).(*Token)
	if !isToken {
		return nil, fmt.Errorf("invalid return type, expected *Token, got %T", args.Get(0))
	}

	return token, nil
}
