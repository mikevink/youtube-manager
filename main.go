package main

import (
	"log"
	"path/filepath"

	"google.golang.org/api/youtube/v3"
)

func configFile() string {
	return filepath.Join(configDir(), "config.yaml")
}

func listPlaylists(service *youtube.Service) {
	call := service.Playlists.List([]string{"snippet", "contentDetails"}).ChannelId("-").MaxResults(50)
	response, err := call.Do()
	onError(err, "")
	for _, list := range response.Items {
		log.Printf("Playlist - ID: %s Name: %s\n", list.Id, list.Snippet.Title)
	}
}

func main() {

	listPlaylists(getService())
}
