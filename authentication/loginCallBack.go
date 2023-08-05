package authentication

import (
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type LoginCallBack interface {
	OnLoginSuccess(chan *spotify.Client, *spotifyauth.Authenticator, *oauth2.Token)
}
