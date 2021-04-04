package spotify_server

import (
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

type Spotify_server struct {
	auth spotify.Authenticator
	ch chan *spotify.Client
	state string
	redirectURI string
	Client *spotify.Client
}

func New_Spotify_Server(auth spotify.Authenticator, state string, redirectURI string) *Spotify_server {
	s := &Spotify_server{
		auth:  auth,
		ch:    make(chan *spotify.Client),
		state: state,
		redirectURI: redirectURI,
	}
	return s
}


func (spotify_server *Spotify_server)Configure() {
	//set auth secret and ID
	spotify_server.auth.SetAuthInfo("bd2af774145b483eba8eada67529fa37", "e5d0f05a55a4486a9f865b5547d09cb0")
	// first start an HTTP server
	http.HandleFunc("/callback", spotify_server.completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := spotify_server.auth.AuthURL(spotify_server.state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	// wait for auth to complete
	client := <-spotify_server.ch
	spotify_server.Client = client
}

func (Spotify_server *Spotify_server)completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := Spotify_server.auth.Token(Spotify_server.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != Spotify_server.state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, Spotify_server.state)
	}
	// use the token to get an authenticated client
	client := Spotify_server.auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	Spotify_server.ch <- &client
}
