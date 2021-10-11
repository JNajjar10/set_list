package server

import (
	"fmt"
	"log"
	"net/http"
	"spotifyMusicVideo/server/spotify_server"
	"spotifyMusicVideo/server/youtube_server"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type ID struct {
	Id string `json:"id"`
}

type Song struct {
	Id string `json:"id"`
}

type PlaylistName struct {
	Name string `json:"name"`
}

type LoginLink struct {
	Link string `json:url`
}

var Playlists spotify.SimplePlaylistPage
var Track spotify.SimpleTrack

type Server struct {
	Router        *gin.Engine
	Address       string
	SpotifyServer *spotify_server.Spotify_server
}

func NewServer() Server {
	spotifyServer := spotify_server.New_Spotify_Server(
		"http://localhost:8080/callback",
	)

	server := Server{
		Router:        gin.Default(),
		Address:       ":8000",
		SpotifyServer: spotifyServer,
	}
	server.Router.Use(CORSMiddleware())

	return server
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func (server *Server) InitialiseRoutes() {
	api := server.Router.Group("/api")
	{
		api.GET("/track", server.returnTrack)
		api.GET("/video", server.returnVideo)
		api.GET("/user/current", server.returnCurrentUser)
		api.POST("/video/:song", server.returnVideoForSong)
		api.GET("/user/playlists", server.returnCurrentUserPlaylists)
		api.GET("/playlists/trending", server.returnTrendingPlaylists)
		api.GET("/playlists/festivals", server.returnFestivalPlaylists)
		api.GET("/playlist/:id/tracks", server.returnPlaylistTracks)
		api.GET("/authenticate", server.authenticateSpotify)
	}
}

func (server *Server) returnTrendingPlaylists(c *gin.Context) {
	client := server.SpotifyServer.Client
	message, playlistPage, _ := client.FeaturedPlaylists()
	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	log.Print(message)
	c.JSON(http.StatusOK, playlistPage.Playlists)
}

func (server *Server) returnFestivalPlaylists(c *gin.Context) {
	//playlistIDs := [5]string{
	//	"4BSj7IwrHLRpZeVNodMh5Z",
	//	"3VXReCeetN58c1clj9u8ZK",
	//	"1hjFALWWkfBrUlTlj1pleA",
	//	"3fgH9GrbLUOMaqcEczyolQ",
	//	"4RvrsTD46yvJYzS3HBGkRB",}
	//var festivalPlaylsits []spotify.SimplePlaylist
	//client := server.SpotifyServer.Client
	//for i, v := range playlistIDs {
	//	client.Search()
	//	playlistPage, _ := client.GetPlaylist(v)
	//
	//}
	//c.Header("Content-Type", "application/json")
	//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//log.Print(message)
	//c.JSON(http.StatusOK,  playlistPage.Playlists)
}

func (server *Server) returnCurrentUser(c *gin.Context) {
	client := server.SpotifyServer.Client
	currentUser, _ := client.CurrentUser()
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, currentUser)
}

func (server *Server) returnVideo(c *gin.Context) {
	client := server.SpotifyServer.Client
	track, _ := client.GetTrack("3n3Ppam7vgaVa1iaRUc9Lp")
	Track = track.SimpleTrack
	youtube_client := youtube_server.NewYoutubeClient()
	song, _ := youtube_client.GetSong(Track.Name)

	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.JSON(http.StatusOK, Song{Id: song})
}

func (server *Server) returnTrack(c *gin.Context) {
	client := server.SpotifyServer.Client
	track, _ := client.GetTrack("3n3Ppam7vgaVa1iaRUc9Lp")
	Track = track.SimpleTrack

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, Track)
}

func (server *Server) returnVideoForSong(c *gin.Context) {
	song := c.Param("song")
	song = strings.ReplaceAll(song, "%20", "+")
	song = strings.ReplaceAll(song, " ", "+")
	//client := server.SpotifyServer.Client
	//track, _ := client.Search(song, spotify.SearchTypeTrack)
	//trackName := track.Tracks.Tracks[0].Name
	youtube_client := youtube_server.NewYoutubeClient()
	youtubeSearchResult, err := youtube_client.GetSong(song)
	if err != nil {
		fmt.Print(err)
	}
	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, Song{Id: youtubeSearchResult})

}

func (server *Server) returnCurrentUserPlaylists(c *gin.Context) {
	client := server.SpotifyServer.Client
	playlistPage, _ := client.CurrentUsersPlaylists()
	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.JSON(http.StatusOK, playlistPage.Playlists)
}

func (server *Server) returnPlaylistTracks(c *gin.Context) {
	playlistID := spotify.ID(c.Param("id"))
	client := server.SpotifyServer.Client
	playlistTracksPage, _ := client.GetPlaylistTracks(playlistID)
	playlistTracks := playlistTracksPage.Tracks
	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, playlistTracks)
}

//Authenticates Spotify user's token with the server
func (server *Server) authenticateSpotify(c *gin.Context) {
	accessToken := c.Request.Header.Get("Authorization")
	server.SpotifyServer.Authenticate(&oauth2.Token{AccessToken: accessToken})
	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusNoContent, "")
}
