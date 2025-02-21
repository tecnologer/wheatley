package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tecnologer/wheatley/pkg/twitch"
)

var (
	clientID     string
	clientSecret string
)

func init() {
	flag.StringVar(&clientID, "client-id", os.Getenv("TWITCH_CLIENT_ID"), "Twitch Client ID")
	flag.StringVar(&clientSecret, "client-secret", os.Getenv("TWITCH_CLIENT_SECRET"), "Twitch Client Secret")
	flag.Parse()
}

func main() {
	streamerName := "tokejay2"

	stream, err := twitch.New(clientID, clientSecret).StreamByName(context.Background(), streamerName)
	if err != nil && !errors.Is(err, twitch.ErrNotFound) {
		log.Fatalln(err)
	}

	if errors.Is(err, twitch.ErrNotFound) {
		fmt.Printf(streamerName + " is not live\n")
		return
	}

	fmt.Printf("%s is live! - Streaming %s to %d viewers\n", stream.UserDisplayName, stream.GameName, stream.ViewerCount)
}
