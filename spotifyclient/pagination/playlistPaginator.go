package pagination

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"log"
)

type PlaylistPaginator struct {
}

func NewPlayListIterator() *PlaylistPaginator {
	return &PlaylistPaginator{}
}

func (_ *PlaylistPaginator) RequestData(client *spotify.Client, options *PageOptions, ch chan spotify.SimplePlaylist) bool {
	playlistPage, err := client.CurrentUsersPlaylists(context.Background(), spotify.Limit(options.Limit), spotify.Offset(options.Offset))

	if err != nil {
		log.Fatal(err)
	}

	for _, playlist := range playlistPage.Playlists {
		ch <- playlist
	}

	return playlistPage.Next != ""
}
