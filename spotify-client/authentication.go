package spotify_client

import (
	"context"
	"fmt"
	"github.com/jacobdrury/spotify-duplicates/utils"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
)

func (c *SpotifyClient) Authenticate() *spotify.Client {
	c.auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(c.redirectUri),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
		))

	// first start an HTTP server
	http.HandleFunc("/callback", c.completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := c.auth.AuthURL(c.state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// TODO: Use webscrapping to automatically login with creads in .env
	utils.OpenBrowser(url)

	// Wait for auth to complete
	client := <-c.clientChannel

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	c.currentUser = user

	fmt.Println("You are logged in as:", c.currentUser.ID)

	return client
}

func (c *SpotifyClient) completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := c.auth.Token(r.Context(), c.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != c.state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, c.state)
	}

	// use the token to get an authenticated client
	client := spotify.New(c.auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	c.clientChannel <- client
}
