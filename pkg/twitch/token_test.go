package twitch_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tecnologer/wheatley/pkg/twitch"
)

func TestRenew(t *testing.T) { //nolint: funlen
	t.Parallel()

	tests := []struct {
		name        string
		mockResp    twitch.Token
		mockStatus  int
		mockErr     error
		wantErr     bool
		expectedErr string
	}{
		{
			name:       "successful_token_renewal",
			mockResp:   twitch.Token{AccessToken: "new-token", ExpiresIn: 3600, TokenType: "bearer"},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:        "error_response",
			mockResp:    twitch.Token{ErrResponse: &twitch.ErrResponse{Status: 400, Message: "invalid grant type"}},
			mockStatus:  http.StatusBadRequest,
			wantErr:     true,
			expectedErr: "error response: Code: 400 - invalid credentials",
		},
		{
			name:        "network_failure",
			mockErr:     errors.New("network failure"),
			wantErr:     true,
			expectedErr: "do request: network failure",
		},
		{
			name:        "invalid_json",
			mockResp:    twitch.Token{},
			mockStatus:  http.StatusOK,
			mockErr:     errors.New("invalid json"),
			wantErr:     true,
			expectedErr: "decode response:",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Serialize mock response
			mockRespData, _ := json.Marshal(test.mockResp)

			// Create a mock HTTP client
			client := &http.Client{
				Transport: &twitch.MockRoundTripper{
					Response: string(mockRespData),
					Status:   test.mockStatus,
					Err:      test.mockErr,
				},
			}

			// Call the function
			service := twitch.New(&twitch.Config{
				ClientID:     "client_id",
				ClientSecret: "client_secret",
			})
			got, err := service.RenewToken(&twitch.Request{
				HTTPClient:   client,
				Context:      context.Background(),
				ClientID:     "client_id",
				ClientSecret: "client_secret",
			})

			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, test.mockResp.AccessToken, got.AccessToken)
		})
	}
}
