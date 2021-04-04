package main

import (
	"fmt"
	my_server "spotifyMusicVideo/server"
)

const developerKey = "AIzaSyAqQce_UVA3HC11-7yrlsCB7kJI_b6ALV4"

func main(){
	//Initialize the server
	server := my_server.NewServer()

	//Setup data
	my_server.InitialiseData()
	name, _ := server.SpotifyServer.Client.CurrentUser()
	my_server.Id = append(my_server.Id, my_server.ID{name.ID})

	//Start Server
	go func() {
		server.InitialiseRoutes()
		server.Run()
	}()

	//Youtube stuff

	//youtube_client := youtube_server.NewYoutubeClient()
	//song, _ := youtube_client.GetSong("Laroi")
	//fmt.Println(song)

	select {
	}
}

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}