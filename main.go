package main

import (
	"flag"
	"fmt"
	spotify "github.com/jacobdrury/spotify-duplicates/spotifyclient"
	"github.com/joho/godotenv"
	"log"
	"strings"
)

func main() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	// Define flags for playlistIds and all
	var playlistIds string
	var all bool

	flag.StringVar(&playlistIds, "playlistIds", "", "Comma-separated list of playlist IDs")
	flag.BoolVar(&all, "all", false, "Remove duplicates from all playlists")

	// Parse the command-line arguments
	flag.Parse()
	if !all && playlistIds == "" {
		fmt.Println("Please provide either -playlistIds or -all as command-line arguments.")
		return
	}

	spotifyClient := spotify.NewClient()

	if all {
		spotifyClient.RemoveDuplicatesFromAllPlaylists()
	} else if playlistIds != "" {
		// Split playlistIds by comma to get individual IDs
		ids := strings.Split(playlistIds, ",")
		fmt.Println("Remove duplicates from playlists with IDs:", ids)

		spotifyIds := spotify.StringIdsToSpotifyIds(ids)
		spotifyClient.RemoveDuplicatesFromPlaylistsById(spotifyIds)
	}
}
