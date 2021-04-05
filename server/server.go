package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
	"net/http"
	"spotifyMusicVideo/server/spotify_server"
	"spotifyMusicVideo/server/youtube_server"
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


var Playlists spotify.SimplePlaylistPage
var Track spotify.SimpleTrack

type Server struct {
	Router *gin.Engine
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
		Router: gin.Default(),
		Address: ":3000",
		SpotifyServer: spotifyServer,
	}
}

func (server *Server)InitialiseRoutes() {
	api := server.Router.Group("/api")
	{
		api.GET("/track", server.returnTrack)
		api.GET("/video", server.returnVideo)
		api.GET("/user/current", server.returnCurrentUser)
		api.POST("/video/:song", server.returnVideoForSong)
		api.GET("/user/playlists", server.returnCurrentUserPlaylists)
	}
}

func (server *Server)returnCurrentUser(c *gin.Context) {
	client := server.SpotifyServer.Client
	currentUser, _ := client.CurrentUser()
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, currentUser)
}

func (server *Server)returnVideo(c *gin.Context)  {
	client := server.SpotifyServer.Client
	track, _ := client.GetTrack("3n3Ppam7vgaVa1iaRUc9Lp")
	Track = track.SimpleTrack
	youtube_client := youtube_server.NewYoutubeClient()
	song, _ := youtube_client.GetSong(Track.Name)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, Song{Id: "https://www.youtube.com/watch?v=" + song})
}

func (server *Server)returnTrack(c *gin.Context)  {
	client := server.SpotifyServer.Client
	track, _ := client.GetTrack("3n3Ppam7vgaVa1iaRUc9Lp")
	Track = track.SimpleTrack

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, Track)
}

func (server *Server)returnVideoForSong(c *gin.Context) {
	song := c.Param("song")
	client := server.SpotifyServer.Client
	track, _ := client.Search(song, spotify.SearchTypeTrack)
	trackName := track.Tracks.Tracks[0].Name
	youtube_client := youtube_server.NewYoutubeClient()
	youtubeSearchResult, _ := youtube_client.GetSong(trackName)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, Song{Id:"https://www.youtube.com/watch?v=" + youtubeSearchResult})

}

func (server *Server)returnCurrentUserPlaylists(c *gin.Context)  {
	client := server.SpotifyServer.Client
	playlistPage, _ := client.CurrentUsersPlaylists()
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK,  playlistPage.Playlists)
}

