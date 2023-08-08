package spotifyapi

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"log"
)

type AllUserPlaylistPaginator struct {
}

func NewAllUserPlaylistPaginator() *AllUserPlaylistPaginator {
	return &AllUserPlaylistPaginator{}
}

func (_ *AllUserPlaylistPaginator) RequestData(client *spotify.Client, options *PageOptions, ch chan []spotify.SimplePlaylist) bool {
	playlistPage, err := client.CurrentUsersPlaylists(context.Background(), spotify.Limit(options.Limit), spotify.Offset(options.Offset))

	if err != nil {
		log.Fatal(err)
	}

	ch <- playlistPage.Playlists

	return playlistPage.Next != ""
}
