package main

import (
	spotify "github.com/jacobdrury/spotify-duplicates/spotifyclient"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// TODO: Parse cmd line args for specific playlists

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	spotify.
		NewClient().
		RemoveDuplicatesFromPlaylists()
}
