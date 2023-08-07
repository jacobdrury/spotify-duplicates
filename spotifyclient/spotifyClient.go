package spotifyclient

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"log"
)

type SpotifyClient struct {
	client      *spotify.Client
	currentUser *spotify.PrivateUser
}

func NewClient() *SpotifyClient {
	return &SpotifyClient{}
}

func (c *SpotifyClient) OnLoginSuccess(ch chan *spotify.Client, auth *spotifyAuth.Authenticator, tok *oauth2.Token) {
	// Use the token to get an authenticated client
	client := spotify.New(auth.Client(context.TODO(), tok))

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	c.currentUser = user

	fmt.Println("You are logged in as:", c.currentUser.ID)
	ch <- client
}
