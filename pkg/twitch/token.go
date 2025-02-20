package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const authTokenURL = "https://id.twitch.tv/oauth2/token"

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	lastUpdate  time.Time
}

func (t *Token) IsValid() bool {
	return t != nil && time.Since(t.lastUpdate) < time.Duration(t.ExpiresIn)*time.Second
}

func (t *Token) Renew(clientID, clientSecret string) (*Token, error) {
	req, err := http.NewRequest("POST", authTokenURL, formURLEncoded(clientID, clientSecret))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer res.Body.Close()

	var token Token
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	token.lastUpdate = time.Now()

	return &token, nil
}

func formURLEncoded(clientID, clientSecret string) io.Reader {
	return strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", clientID, clientSecret))
}
