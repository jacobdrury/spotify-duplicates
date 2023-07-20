package spotify_client

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/exp/maps"
	"log"
	"os"
)

type SpotifyClient struct {
	auth          *spotifyAuth.Authenticator
	clientChannel chan *spotify.Client
	client        *spotify.Client
	currentUser   *spotify.PrivateUser
	state         string

	// URI spotify will redirect user to after successful auth
	redirectUri string
}

func NewSpotifyClient() *SpotifyClient {
	redirectUri := os.Getenv("REDIRECT_URI")
	if redirectUri == "" {
		log.Fatal("REDIRECT_URI environment variable not found")
	}

	return &SpotifyClient{
		auth: spotifyAuth.New(
			spotifyAuth.WithRedirectURL(redirectUri),
			spotifyAuth.WithScopes(
				spotifyAuth.ScopeUserReadPrivate,
				spotifyAuth.ScopePlaylistModifyPrivate,
				spotifyAuth.ScopePlaylistModifyPublic,
				spotifyAuth.ScopePlaylistReadPrivate,
				spotifyAuth.ScopePlaylistReadCollaborative,
			)),
		clientChannel: make(chan *spotify.Client),
		state:         "abc123",
		redirectUri:   redirectUri,
	}
}

func (c *SpotifyClient) RemoveDuplicatesFromPlaylists() {
	c.client = c.Authenticate()
	c.processPlaylists()
}

func (c *SpotifyClient) processPlaylists() {
	usersPlayLists := make([]spotify.SimplePlaylist, 0)

	pageLimit := 50
	pageOffSet := 0
	hasNext := true
	for hasNext {
		playlistPage, err := c.client.CurrentUsersPlaylists(context.Background(), spotify.Limit(pageLimit), spotify.Offset(pageOffSet))

		if err != nil {
			log.Fatal(err)
		}

		//fmt.Printf("Found %d of playlists\n", len(playlistPage.Playlists))
		for _, playlist := range playlistPage.Playlists {
			// && playlist.ID == "3QgOKxAuYfpdqBPdYp6wBl"
			if playlist.Owner.ID == c.currentUser.ID {
				usersPlayLists = append(usersPlayLists, playlist)
			}
		}

		hasNext = playlistPage.Next != ""
		pageOffSet += (pageOffSet + 1) * pageLimit
	}

	fmt.Printf("Found a total of %d playlists owned by %s\n", len(usersPlayLists), c.currentUser.ID)

	for _, playlist := range usersPlayLists {
		duplicateTracks := make(map[spotify.ID]spotify.TrackToRemove)
		hashMap := make(map[string]spotify.ID)
		itemsLimit := 50
		itemsPageOffset := 0
		itemsHasNext := true

		for itemsHasNext {
			fullPlaylist, err := c.client.GetPlaylistItems(context.Background(), playlist.ID, spotify.Limit(itemsLimit), spotify.Offset(itemsPageOffset))
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
			}

			itemsHasNext = fullPlaylist.Next != ""
			itemsPageOffset += (itemsPageOffset + 1) * itemsLimit
		}

		fmt.Printf("Removing %d duplicate tracks in %s\b\n", len(duplicateTracks), playlist.Name)
		_, err := c.client.RemoveTracksFromPlaylistOpt(context.Background(), playlist.ID, maps.Values(duplicateTracks), "")

		if err != nil {
			log.Fatal(err)
		}
	}
}
