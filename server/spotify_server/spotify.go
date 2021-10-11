package spotify_server

import (
	"context"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type Spotify_server struct {
	ch     chan *spotify.Client
	Client *spotify.Client
}

func New_Spotify_Server(redirectURI string) *Spotify_server {
	s := &Spotify_server{
		ch: make(chan *spotify.Client),
	}
	return s
}

func (spotify_server *Spotify_server) Authenticate(token *oauth2.Token) {
	token = &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    "",
		RefreshToken: "",
		Expiry:       time.Time{},
	}
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	spotifyClient := spotify.NewClient(httpClient)
	spotify_server.Client = &spotifyClient
}
