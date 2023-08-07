package spotifyapi

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"log"
)

type ItemsPaginator struct {
	playlistId spotify.ID
}

type PlaylistItem struct {
	Track    spotify.SimpleTrack
	Position int
}

func NewItemsPaginator(playlistId spotify.ID) *ItemsPaginator {
	return &ItemsPaginator{playlistId: playlistId}
}

func (p *ItemsPaginator) RequestData(client *spotify.Client, options *PageOptions, ch chan PlaylistItem) bool {
	fullPlaylist, err := client.GetPlaylistItems(context.Background(), p.playlistId, spotify.Limit(options.Limit), spotify.Offset(options.Offset))
	if err != nil {
		log.Fatal(err)
	}

	for i, item := range fullPlaylist.Items {
		ch <- PlaylistItem{
			Track:    item.Track.Track.SimpleTrack,
			Position: i + options.Offset,
		}
	}

	return fullPlaylist.Next != ""
}
