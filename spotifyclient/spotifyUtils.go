package spotifyclient

import (
	"github.com/zmb3/spotify/v2"
	"strings"
)

func StringIdToSpotifyId(id string) spotify.ID {
	return spotify.ID(strings.TrimSpace(id))
}

func StringIdsToSpotifyIds(ids []string) []spotify.ID {
	spotifyIds := make([]spotify.ID, 0)
	for _, id := range ids {
		spotifyIds = append(spotifyIds, StringIdToSpotifyId(id))
	}

	return spotifyIds
}
