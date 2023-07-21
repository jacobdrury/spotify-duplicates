package spotify_client

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"log"
)

func PlayListIterator(client *spotify.Client, currentUser *spotify.PrivateUser) <-chan spotify.SimplePlaylist {
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
				// && playlist.ID == "3QgOKxAuYfpdqBPdYp6wBl"
				// filter out collaborative and liked playlists
				if playlist.Owner.ID == currentUser.ID /*&& (playlist.ID == "3QgOKxAuYfpdqBPdYp6wBl" || playlist.ID == "4GaM4VDRN25luGjxvjrsIx")*/ {
					ch <- playlist
				}
			}

			offset += limit

			if playlistPage.Next == "" {
				log.Println("Finished fetching all playlists")
				close(ch)
				return
			}
		}
	}()

	return ch
}

type TrackPosition struct {
	Track    spotify.SimpleTrack
	Position int
}

func ItemsIterator(client *spotify.Client, playlistId spotify.ID) <-chan TrackPosition {
	ch := make(chan TrackPosition)
	offset := 0
	limit := 50

	go func() {
		for {
			fullPlaylist, err := client.GetPlaylistItems(context.Background(), playlistId, spotify.Limit(limit), spotify.Offset(offset))
			if err != nil {
				log.Fatal(err)
			}

			for i, item := range fullPlaylist.Items {
				ch <- TrackPosition{
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
