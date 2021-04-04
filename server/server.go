package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"spotifyMusicVideo/server/spotify_server"
	"spotifyMusicVideo/server/youtube_server"
)

type Article struct {
	Id      string `json:"Id"`
	Title string `json:"Title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}

type ID struct {
	Id string `json:"id"`
}

type Song struct {
	Id string `json:"id"`
}

var Track spotify.SimpleTrack
var Playlists spotify.SimplePlaylistPage

type Server struct {
	Router *mux.Router
	Address string
	SpotifyServer *spotify_server.Spotify_server
}

func NewServer() Server{
	spotifyServer := spotify_server.New_Spotify_Server(
		spotify.NewAuthenticator(
			"http://localhost:8080/callback", spotify.ScopeUserReadPrivate),
		"abc123",
		"http://localhost:8080/callback",
	)
	spotifyServer.Configure()

	return Server{
		Router: mux.NewRouter(),
		Address: ":8181",
		SpotifyServer: spotifyServer,
	}
}
// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var Articles []Article
var Id []ID
var SongID []Song

func (server *Server)InitialiseRoutes() {
	server.Router.HandleFunc("/home", homePage)
	server.Router.HandleFunc("/all", returnAllArticles)
	server.Router.HandleFunc("/id", server.returnID)
	server.Router.HandleFunc("/article/{id}", returnSingleArticle)
	server.Router.HandleFunc("/track", server.returnTrack)
	server.Router.HandleFunc("/playlists", server.returnCurrentUserPlaylists)
	server.Router.HandleFunc("/video", server.returnVideo)
}

func InitialiseData()  {
	Articles = []Article{
		Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	Id = []ID{
		ID{Id: "bob"},
		ID{Id: "jon"},
	}
}

func (server *Server)Run()  {
	log.Fatal(http.ListenAndServe(server.Address, server.Router))
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func (server *Server)returnID(w http.ResponseWriter, r *http.Request)  {
	client := server.SpotifyServer.Client
	currentUser, _ := client.CurrentUser()
	Id = append(Id, ID{
		Id: currentUser.ID,
	})
	json.NewEncoder(w).Encode(Id)
}

func (server *Server)returnVideo(w http.ResponseWriter, r *http.Request)  {
	client := server.SpotifyServer.Client
	track, _ := client.GetTrack("3n3Ppam7vgaVa1iaRUc9Lp")
	Track = track.SimpleTrack
	youtube_client := youtube_server.NewYoutubeClient()
	song, _ := youtube_client.GetSong(Track.Name)
	SongID = append(SongID, Song{Id: "https://www.youtube.com/watch?v=" + song})
	json.NewEncoder(w).Encode(SongID)
}

func (server *Server)returnTrack(w http.ResponseWriter, r *http.Request)  {
	client := server.SpotifyServer.Client
	track, _ := client.GetTrack("3n3Ppam7vgaVa1iaRUc9Lp")
	Track = track.SimpleTrack
	json.NewEncoder(w).Encode(Track)
}

func (server *Server)returnCurrentUserPlaylists(w http.ResponseWriter, r *http.Request)  {
	client := server.SpotifyServer.Client
	Playlists, _ := client.CurrentUsersPlaylists()
	json.NewEncoder(w).Encode(Playlists)
}
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
