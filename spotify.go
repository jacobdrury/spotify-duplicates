package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type SpotifyClient struct {
	state string
	auth *spotifyauth.Authenticator

}

func (s *SpotifyClient) Init() {
	redirectURL := "http://localhost:5543/callback"

	http.HandleFunc("/callback", s.redirectHandler)
    log.Fatal(http.ListenAndServe(":10000", nil))

	
	// the redirect URL must be an exact match of a URL you've registered for your application
	// scopes determine which permissions the user is prompted to authorize
	s.auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURL), spotifyauth.WithScopes(
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopePlaylistModifyPrivate,
		spotifyauth.ScopePlaylistModifyPublic,
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopePlaylistReadCollaborative,
	))

	// get the user to this URL - how you do that is up to you
	// you should specify a unique state string to identify the session

	s.state = fmt.Sprintf("%s", uuid.New())
	url := s.auth.AuthURL(s.state)

	openBrowser(url)
}

func (s *SpotifyClient) redirectHandler(w http.ResponseWriter, r *http.Request) {
	// use the same state string here that you used to generate the URL
	token, err := s.auth.Token(r.Context(), s.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	// create a client using the specified token
	client := spotify.New(s.auth.Client(r.Context(), token))

	// Public playlist owned by noah.stride:
	// "Long playlist for testing pagination"
	playlistID := "3LsZV4IGAzA0yi59XEEPr3?si=a7c167d4381e4162"
	if id := os.Getenv("SPOTIFY_PLAYLIST"); id != "" {
		playlistID = id
	}

	tracks, err := client.GetPlaylistItems(
		ctx,
		spotify.ID(playlistID),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Playlist has %d total tracks", tracks.Total)
	for page := 1; ; page++ {
		log.Printf("  Page %d has %d tracks", page, len(tracks.Items))
		err = client.NextPage(ctx, tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}