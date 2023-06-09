package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
        log.Fatalf("Error loading environment variables file")
    }

	var spotify = SpotifyClient{}
	spotify.Init()
	
	// ctx := context.Background()
	// config := &clientcredentials.Config{
	// 	ClientID:     os.Getenv("SPOTIFY_ID"),
	// 	ClientSecret: os.Getenv("SPOTIFY_SECRET"),
	// 	TokenURL:     spotifyauth.TokenURL,
	// }
	// token, err := config.Token(ctx)
	// if err != nil {
	// 	log.Fatalf("couldn't get token: %v", err)
	// }

	// httpClient := spotifyauth.New().Client(ctx, token)
	// client := spotify.New(httpClient)

	// // Public playlist owned by noah.stride:
	// // "Long playlist for testing pagination"
	// playlistID := "3LsZV4IGAzA0yi59XEEPr3?si=a7c167d4381e4162"
	// if id := os.Getenv("SPOTIFY_PLAYLIST"); id != "" {
	// 	playlistID = id
	// }

	// tracks, err := client.GetPlaylistItems(
	// 	ctx,
	// 	spotify.ID(playlistID),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("Playlist has %d total tracks", tracks.Total)
	// for page := 1; ; page++ {
	// 	log.Printf("  Page %d has %d tracks", page, len(tracks.Items))
	// 	err = client.NextPage(ctx, tracks)
	// 	if err == spotify.ErrNoMorePages {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
}