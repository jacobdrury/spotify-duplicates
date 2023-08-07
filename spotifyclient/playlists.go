package spotifyclient

import (
	"context"
	"fmt"
	"github.com/jacobdrury/spotify-duplicates/authentication"
	"github.com/jacobdrury/spotify-duplicates/spotifyclient/spotifyapi"
	"github.com/jacobdrury/spotify-duplicates/utils"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/exp/maps"
	"log"
	"sync"
)

func (c *SpotifyClient) RemoveDuplicatesFromAllPlaylists() {
	c.client = authentication.NewHandler().Login(c)
	wg := sync.WaitGroup{}

	playlistCh := spotifyapi.ConsumePaginatedEndpoint[[]spotify.SimplePlaylist](
		c.client,
		spotifyapi.AllUserPlaylistIterator(),
		spotifyapi.DefaultPageOptions())

	// Process each page as we get them from the paginated endpoint
	for playlists := range playlistCh {
		wg.Add(1)
		go func(playlists []spotify.SimplePlaylist) {
			defer wg.Done()
			c.removeDuplicatesFromPlaylists(playlists)
		}(playlists)
	}

	wg.Wait()
}

func (c *SpotifyClient) RemoveDuplicatesFromPlaylistsById(ids []spotify.ID) {
	c.client = authentication.NewHandler().Login(c)

	playlists := make([]spotify.SimplePlaylist, len(ids))
	for _, playlistId := range ids {
		playlist, err := c.client.GetPlaylist(context.Background(), playlistId)
		if err != nil {
			fmt.Printf("Error occurred when fetching playlistId: %s, Err: %s\n", playlistId, err)
			continue
		}

		playlists = append(playlists, playlist.SimplePlaylist)
	}

	c.removeDuplicatesFromPlaylists(playlists)
}

func (c *SpotifyClient) removeDuplicatesFromPlaylists(playlists []spotify.SimplePlaylist) {
	usersPlayLists := make([]spotify.SimplePlaylist, 0)
	wg := sync.WaitGroup{}

	for _, playlist := range playlists {
		// Skip collaborative and liked playlists
		if playlist.Owner.ID != c.currentUser.ID {
			continue
		}

		// Restrict to just my testing playlists so i don't mess up my personal ones while developing this :)
		if playlist.ID != "3QgOKxAuYfpdqBPdYp6wBl" && playlist.ID != "4GaM4VDRN25luGjxvjrsIx" {
			continue
		}

		wg.Add(1)
		usersPlayLists = append(usersPlayLists, playlist)
		go func(playlist spotify.SimplePlaylist) {
			defer wg.Done()
			c.removeDuplicateTracksFromPlaylist(&playlist)
		}(playlist)
	}

	wg.Wait()
	fmt.Printf("Successfully removed all duplicates from %d playlists\n", len(usersPlayLists))
}

func (c *SpotifyClient) removeDuplicateTracksFromPlaylist(playlist *spotify.SimplePlaylist) {
	duplicateTracks := c.findDuplicates(playlist)

	if len(duplicateTracks) > 0 {
		c.removeTracks(playlist, duplicateTracks)
	}
}

func (c *SpotifyClient) findDuplicates(playlist *spotify.SimplePlaylist) map[spotify.ID]spotify.TrackToRemove {
	hashMap := make(map[spotify.ID]spotify.ID)
	duplicateTracks := make(map[spotify.ID]spotify.TrackToRemove)

	// check for duplicates
	itemCh := spotifyapi.ConsumePaginatedEndpoint[spotifyapi.PlaylistItem](
		c.client,
		spotifyapi.NewItemsPaginator(playlist.ID),
		spotifyapi.DefaultPageOptions())

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
		c.removeDuplicateTracksFromPlaylist(playlist)
	}
}

func addDuplicateTrack(duplicateTracks map[spotify.ID]spotify.TrackToRemove, trackPosition *spotifyapi.PlaylistItem) {
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
