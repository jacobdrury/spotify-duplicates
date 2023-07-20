package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jacobdrury/spotify-duplicates/utils"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	auth  *spotifyauth.Authenticator
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
		))

	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	utils.OpenBrowser(url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	usersPlayLists := make([]spotify.SimplePlaylist, 0)

	pageLimit := 50
	pageOffSet := 0
	hasNext := true
	for hasNext {
		playlistPage, err := client.CurrentUsersPlaylists(context.Background(), spotify.Limit(pageLimit), spotify.Offset(pageOffSet))

		if err != nil {
			log.Fatal(err)
		}

		//fmt.Printf("Found %d of playlists\n", len(playlistPage.Playlists))

		for _, playlist := range playlistPage.Playlists {
			if playlist.Owner.ID == user.ID && playlist.ID == "52RPoUq4YfBUXbZO7WKix0" {
				usersPlayLists = append(usersPlayLists, playlist)
			}
		}

		hasNext = playlistPage.Next != ""
		pageOffSet += (pageOffSet + 1) * pageLimit
	}

	fmt.Printf("Found a total of %d playlists owned by %s\n", len(usersPlayLists), user.ID)

	for _, playlist := range usersPlayLists {
		duplicateTracks := make(map[spotify.ID]spotify.TrackToRemove, 0)
		hashMap := make(map[string]spotify.ID)
		itemsLimit := 50
		itemsPageOffset := 0
		itemsHasNext := true

		for itemsHasNext {
			fullPlaylist, err := client.GetPlaylistItems(context.Background(), playlist.ID, spotify.Limit(itemsLimit), spotify.Offset(itemsPageOffset))
			if err != nil {
				log.Fatal(err)
			}

			for i, item := range fullPlaylist.Items {
				if hashMap[item.Track.Track.SimpleTrack.Name] != "" {
					if existingTrackToRemove, found := duplicateTracks[item.Track.Track.SimpleTrack.ID]; found {
						// The track is already in duplicateTracks. Add the position to the Positions slice.
						existingTrackToRemove.Positions = append(existingTrackToRemove.Positions, i+itemsPageOffset)
						duplicateTracks[item.Track.Track.SimpleTrack.ID] = existingTrackToRemove
					} else {
						// The track is not yet in duplicateTracks. Add it with the current position.
						duplicateTracks[item.Track.Track.SimpleTrack.ID] = spotify.TrackToRemove{
							URI:       string(item.Track.Track.SimpleTrack.URI),
							Positions: []int{i + itemsPageOffset},
						}
					}
				}

				hashMap[item.Track.Track.Name] = item.Track.Track.SimpleTrack.ID
				//item.AddedAt
			}

			itemsHasNext = fullPlaylist.Next != ""
			itemsPageOffset += (itemsPageOffset + 1) * itemsLimit
		}

		prettyPrintJson(duplicateTracks)
		fmt.Printf("Removing %d duplicate tracks in %s\b\n", len(duplicateTracks), playlist.Name)
		//_, err := client.RemoveTracksFromPlaylistOpt(context.Background(), playlist.ID, maps.Values(duplicateTracks), "")
		//
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

func prettyPrintJson(j any) {
	marshaled, err := json.MarshalIndent(j, "", "   ")
	if err != nil {
		log.Fatalf("marshaling error: %s", err)
	}
	fmt.Println(string(marshaled))
}
