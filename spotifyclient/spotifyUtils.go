package spotifyclient

import "github.com/zmb3/spotify/v2"

func StringIdToSpotifyId(id string) spotify.ID {
	return spotify.ID(id)
}

func StringIdsToSpotifyIds(ids []string) []spotify.ID {
	spotifyIds := make([]spotify.ID, 0)
	for _, id := range ids {
		spotifyIds = append(spotifyIds, StringIdToSpotifyId(id))
	}

	return spotifyIds
}
