package authentication

import (
	"context"
	"fmt"
	"github.com/jacobdrury/spotify-duplicates/utils"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
	"os"
	"sync"
)

type Handler struct {
	state       string
	baseUrl     string
	redirectUrl string
	auth        *spotifyauth.Authenticator
	ch          chan *spotify.Client
}

func NewHandler() *Handler {
	baseUri := os.Getenv("BASE_URI")
	if baseUri == "" {
		log.Fatal("BASE_URI environment variable not found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("REDIRECT_PORT environment variable not found")
	}

	state := os.Getenv("SPOTIFY_STATE")
	if state == "" {
		log.Fatal("SPOTIFY_STATE environment variable not found")
	}

	baseUrl := fmt.Sprintf("%s:%s", baseUri, port)
	redirectUrl := fmt.Sprintf("http://%s/callback", baseUrl)

	return &Handler{
		state:       state,
		baseUrl:     baseUrl,
		redirectUrl: redirectUrl,
	}
}

func (h *Handler) Login(cb LoginCallBack) *spotify.Client {
	h.ch = make(chan *spotify.Client)
	h.auth = newSpotifyOAuthClient(h.redirectUrl)

	// first start an HTTP server
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := h.auth.Token(r.Context(), h.state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}

		if st := r.FormValue("state"); st != h.state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, h.state)
		}

		_, _ = fmt.Fprintf(w, "Login Completed!")
		cb.OnLoginSuccess(h.ch, h.auth, token)
	})

	// Start Web Server
	httpServer := &http.Server{Addr: h.baseUrl}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Make sure we shut down the web server once the function finishes
	defer func() {
		if err := httpServer.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
		wg.Wait()
	}()

	url := h.auth.AuthURL(h.state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	utils.OpenBrowser(url)

	// Wait for auth to complete
	client := <-h.ch

	return client
}

func newSpotifyOAuthClient(redirectUrl string) *spotifyauth.Authenticator {
	return spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectUrl),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
		))
}
