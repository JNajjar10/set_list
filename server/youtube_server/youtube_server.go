package youtube_server

import (
	"fmt"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
)

type YoutubeClient struct {
	youtubeService *youtube.Service
}

func NewYoutubeClient() YoutubeClient  {
	client := &http.Client{
		Transport: &transport.APIKey{Key: "AIzaSyAqQce_UVA3HC11-7yrlsCB7kJI_b6ALV4"},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	return YoutubeClient{youtubeService: service}
}

func (youtubeClient *YoutubeClient)GetListOfSongs(search string) (map[string]string, error) {
	// Make the API call to YouTube.
	var part []string
	part = append(part, "id,snippet")
	call := youtubeClient.youtubeService.Search.List(part).
		Q(search).
		MaxResults(25).Order("relevance")
	response, _ := call.Do()
	videos := make(map[string]string)
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		}
	}
	return videos, nil
}

func (youtubeClient *YoutubeClient)GetSong(search string) (string, error) {
	// Make the API call to YouTube.
	var part []string
	part = append(part, "id,snippet")
	call := youtubeClient.youtubeService.Search.List(part).
		Q(search).
		MaxResults(1).Order("relevance")
	response, err := call.Do()
	if err != nil {
		fmt.Print(err)
	}
	var url string
	url = response.Items[0].Id.VideoId
	return url, nil
}
