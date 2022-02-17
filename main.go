package main

import (
	"google.golang.org/api/youtube/v3"
	"log"
)

func listPlaylists(service *youtube.Service) {
	call := service.Playlists.List([]string{"snippet", "contentDetails"}).ChannelId("-").MaxResults(50)
	response, err := call.Do()
	onError(err, "")
	for _, list := range response.Items {
		log.Printf("Playlist - ID: %s Name: %s\n", list.Id, list.Snippet.Title)
	}
}

func main() {

	args := parseArgs()

	if args.SampleConfig {
		maybeWriteSampleConfig()
		return
	}

	log.Println("Nothing to do")

}
