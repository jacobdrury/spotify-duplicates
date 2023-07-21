package spotify_client

import (
	"fmt"
	"github.com/zmb3/spotify/v2"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
	"log"
	"os"
	"sync"
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

	wg := sync.WaitGroup{}
	for playlist := range PlayListIterator(c.client, c.currentUser) {
		wg.Add(1)
		playlist := playlist
		usersPlayLists = append(usersPlayLists, playlist)
		go c.removeDuplicates(&wg, &playlist)
	}

	wg.Wait()
	fmt.Printf("Successfully removed all duplicates from %d playlists", len(usersPlayLists))
}

func (c *SpotifyClient) removeDuplicates(wg *sync.WaitGroup, playlist *spotify.SimplePlaylist) {
	hashMap := make(map[string]spotify.ID)
	duplicateTracks := make(map[spotify.ID]spotify.TrackToRemove)

	// check for duplicates
	for trackPosition := range ItemsIterator(c.client, playlist.ID) {
		trackId := trackPosition.Track.ID

		if hashMap[trackPosition.Track.Name] != "" {
			addDuplicateTrack(duplicateTracks, &trackPosition)
			continue
		}

		hashMap[trackPosition.Track.Name] = trackId
	}

	if playlist.ID == "3QgOKxAuYfpdqBPdYp6wBl" || playlist.ID == "4GaM4VDRN25luGjxvjrsIx" {
		fmt.Printf("Removing %d duplicate tracks in %s\b\n", len(duplicateTracks), playlist.Name)
		//_, err := c.client.RemoveTracksFromPlaylistOpt(context.Background(), playlist.ID, maps.Values(duplicateTracks), "")
		//
		//if err != nil {
		//	log.Fatal(err)
		//}
	}

	defer wg.Done()
}

func addDuplicateTrack(duplicateTracks map[spotify.ID]spotify.TrackToRemove, trackPosition *TrackPosition) {
	existingTrackToRemove, found := duplicateTracks[trackPosition.Track.ID]
	if found {
		// The track is already in duplicateTracks. Add the position to the Positions slice.
		existingTrackToRemove.Positions = append(existingTrackToRemove.Positions, trackPosition.Position)
		duplicateTracks[trackPosition.Track.ID] = existingTrackToRemove
		return
	}

	// The track is not yet in duplicateTracks. Add it with the current position.
	duplicateTracks[trackPosition.Track.ID] = spotify.TrackToRemove{
		URI:       string(trackPosition.Track.URI),
		Positions: []int{trackPosition.Position},
	}
}
