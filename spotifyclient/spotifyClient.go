package spotifyclient

import (
	"context"
	"fmt"
	"github.com/jacobdrury/spotify-duplicates/authentication"
	"github.com/jacobdrury/spotify-duplicates/spotifyclient/pagination"
	"github.com/jacobdrury/spotify-duplicates/utils"
	"github.com/zmb3/spotify/v2"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/exp/maps"
	"golang.org/x/oauth2"
	"log"
	"sync"
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

func (c *SpotifyClient) RemoveDuplicatesFromPlaylists() {
	c.client = authentication.NewHandler().Login(c)

	usersPlayLists := make([]spotify.SimplePlaylist, 0)
	wg := sync.WaitGroup{}

	playlistCh := pagination.ConsumePaginatedEndpoint[spotify.SimplePlaylist](
		c.client,
		pagination.NewPlayListIterator(),
		pagination.DefaultPageOptions())

	for playlist := range playlistCh {
		// Skip collaborative and liked playlists
		if playlist.Owner.ID != c.currentUser.ID {
			continue
		}

		// Restrict to just my testing playlists so i don't mess up my personal ones while developing this :)
		if playlist.ID != "3QgOKxAuYfpdqBPdYp6wBl" && playlist.ID != "4GaM4VDRN25luGjxvjrsIx" {
			continue
		}

		wg.Add(1)
		playlist := playlist
		usersPlayLists = append(usersPlayLists, playlist)
		go func() {
			defer wg.Done()
			c.removeDuplicates(&playlist)
		}()
	}

	wg.Wait()
	fmt.Printf("Successfully removed all duplicates from %d playlists", len(usersPlayLists))
}

func (c *SpotifyClient) removeDuplicates(playlist *spotify.SimplePlaylist) {
	duplicateTracks := c.findDuplicates(playlist)

	if len(duplicateTracks) > 0 {
		c.removeTracks(playlist, duplicateTracks)
	}
}

func (c *SpotifyClient) findDuplicates(playlist *spotify.SimplePlaylist) map[spotify.ID]spotify.TrackToRemove {
	hashMap := make(map[spotify.ID]spotify.ID)
	duplicateTracks := make(map[spotify.ID]spotify.TrackToRemove)

	// check for duplicates
	itemCh := pagination.ConsumePaginatedEndpoint[pagination.PlaylistItem](
		c.client,
		pagination.NewItemsPaginator(playlist.ID),
		pagination.DefaultPageOptions())

	for item := range itemCh {
		// TODO: Update to check for same name and artist instead of ID
		if hashMap[item.Track.ID] != "" {
			addDuplicateTrack(duplicateTracks, &item)
			continue
		}

		hashMap[item.Track.ID] = item.Track.ID
	}

	return duplicateTracks
}

func (c *SpotifyClient) removeTracks(playlist *spotify.SimplePlaylist, duplicateTracks map[spotify.ID]spotify.TrackToRemove) {
	fmt.Printf("Removing %d duplicate tracks in %s\b\n", len(duplicateTracks), playlist.Name)

	// Spotify limits delete to 100 values
	chunkSize := 100
	chunks := utils.ChunkSlice(maps.Values(duplicateTracks), chunkSize)

	_, err := c.client.RemoveTracksFromPlaylistOpt(context.Background(), playlist.ID, chunks[0], "")

	if err != nil {
		log.Fatal(err)
	}

	// If we have more than 100 tracks to delete, we need to re-index the playlist in order to send more delete requests
	// If we send multiple delete requests without re-indexing and updating the position values on the tracksToDelete,
	// they will be off since the position is just the index of the track in the playlist
	if len(chunks) > 1 {
		fmt.Println("Number of tracks exceeds the max allowed to delete at once, re-indexing playlist and trying again...")
		c.removeDuplicates(playlist)
	}
}

func addDuplicateTrack(duplicateTracks map[spotify.ID]spotify.TrackToRemove, trackPosition *pagination.PlaylistItem) {
	existingDuplicateTrack, found := duplicateTracks[trackPosition.Track.ID]
	if found {
		// The track is already in duplicateTracks. Add the position to the Positions slice.
		existingDuplicateTrack.Positions = append(existingDuplicateTrack.Positions, trackPosition.Position)
		duplicateTracks[trackPosition.Track.ID] = existingDuplicateTrack
		return
	}

	// The track is not yet in duplicateTracks. Add it with the current position.
	duplicateTracks[trackPosition.Track.ID] = spotify.TrackToRemove{
		URI:       string(trackPosition.Track.URI),
		Positions: []int{trackPosition.Position},
	}
}
