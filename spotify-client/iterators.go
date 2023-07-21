package spotify_client

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"log"
)

func PlayListIterator(client *spotify.Client) <-chan spotify.SimplePlaylist {
	ch := make(chan spotify.SimplePlaylist)
	offset := 0
	limit := 50

	go func() {
		for {

			playlistPage, err := client.CurrentUsersPlaylists(context.Background(), spotify.Limit(limit), spotify.Offset(offset))

			if err != nil {
				log.Fatal(err)
			}

			for _, playlist := range playlistPage.Playlists {
				ch <- playlist
			}

			offset += limit

			if playlistPage.Next == "" {
				close(ch)
				return
			}
		}
	}()

	return ch
}

type Item struct {
	Track    spotify.SimpleTrack
	Position int
}

func ItemsIterator(client *spotify.Client, playlistId spotify.ID) <-chan Item {
	ch := make(chan Item)
	offset := 0
	limit := 50

	go func() {
		for {
			fullPlaylist, err := client.GetPlaylistItems(context.Background(), playlistId, spotify.Limit(limit), spotify.Offset(offset))
			if err != nil {
				log.Fatal(err)
			}

			for i, item := range fullPlaylist.Items {
				ch <- Item{
					Track:    item.Track.Track.SimpleTrack,
					Position: i + offset,
				}
			}

			offset += limit

			if fullPlaylist.Next == "" {
				close(ch)
				return
			}
		}
	}()

	return ch
}
