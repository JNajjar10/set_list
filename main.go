package main

import (
	my_server "spotifyMusicVideo/server"
)

const developerKey = "AIzaSyAqQce_UVA3HC11-7yrlsCB7kJI_b6ALV4"

func main(){
	//Initialize the server
	server := my_server.NewServer()
	server.InitialiseRoutes()

	// Start and run the server
	server.Router.Run(server.Address)
	select {
	}
}
