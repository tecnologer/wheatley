package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/tecnologer/wheatley/pkg/twitch"
	"log"
)

var (
	clientID     string
	clientSecret string
)

func init() {
	flag.StringVar(&clientID, "client-id", "", "Twitch Client ID")
	flag.StringVar(&clientSecret, "client-secret", "", "Twitch client secret")
	flag.Parse()
}

func main() {
	stream, err := twitch.New(clientID, clientSecret).StreamByName(context.Background(), "texarcane")
	if err != nil && !errors.Is(err, twitch.ErrNotFound) {
		log.Fatalln(err)
	}

	if errors.Is(err, twitch.ErrNotFound) {
		fmt.Printf("texarcane is not live\n")
		return
	}

	fmt.Printf("%s is live! - Streaming %s to %d viewers\n", stream.UserDisplayName, stream.GameName, stream.ViewerCount)
}
