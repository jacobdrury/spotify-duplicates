package main

import (
	spotify_client "github.com/jacobdrury/spotify-duplicates/spotify-client"
	"github.com/jacobdrury/spotify-duplicates/utils"
)

func main() {
	utils.LoadEnvVariables()

	spotifyClient := spotify_client.NewSpotifyClient()
	spotifyClient.RemoveDuplicatesFromPlaylists()
}
