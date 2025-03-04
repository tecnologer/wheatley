package twitch

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const authTokenURL = "https://id.twitch.tv/oauth2/token" //nolint:gosec

type Token struct {
	*ErrResponse
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	lastUpdate  time.Time
}

func (t *Token) IsValid() bool {
	return t != nil && time.Since(t.lastUpdate) < time.Duration(t.ExpiresIn)*time.Second
}

func formURLEncoded(clientID, clientSecret string) io.Reader {
	return strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", clientID, clientSecret))
}

func (t *Token) HasError() bool {
	return t != nil && t.ErrResponse != nil
}
