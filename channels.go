package main

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
	"log"
	"strings"
)

func lookupChannelId(service *youtube.Service, title string) string {
	search, err := service.Search.List([]string{"snippet"}).Type("channel").Q(title).Do()
	onError(err, fmt.Sprintf("Failed to search for channel '%s'", title))
	for _, item := range search.Items {
		if "youtube#channel" == item.Id.Kind && title == item.Snippet.Title {
			return item.Snippet.ChannelId
		}
	}
	log.Fatalf("Could not find channel titled '%s' is first 25 results", title)
	return ""
}

func getPlaylists(service *youtube.Service, channelId string) []*youtube.Playlist {
	request := service.Playlists.List([]string{"snippet"}).MaxResults(50).ChannelId(channelId)
	var playlists []*youtube.Playlist
	fetching := true
	for fetching {
		result, err := request.Do()
		onError(err, fmt.Sprintf("Could not get playlists for channel 'https://www.youtube.com/channel/%s'", channelId))
		playlists = append(playlists, result.Items...)
		if 0 != len(result.NextPageToken) {
			request = request.PageToken(result.NextPageToken)
		} else {
			fetching = false
		}
	}
	return playlists
}

func inspectChannel(service *youtube.Service, channelId string) {
	playlists := getPlaylists(service, channelId)
	if 0 == len(playlists) {
		log.Printf("No playlists found for channel '%s'\n", channelId)
		return
	}
	channelTitle := playlists[0].Snippet.ChannelTitle
	msg := fmt.Sprintf("Channel %s(%s):\n", channelTitle, channelId)
	for _, playlist := range playlists {
		msg = fmt.Sprintf("%s\t%s -- %s\n", msg, playlist.Snippet.Title, playlist.Id)
	}
	log.Println(msg)
}

func inspectChannels(service *youtube.Service, channels []string) []string {
	channelIds := make([]string, len(channels))
	for i, channel := range channels {
		if strings.HasPrefix(channel, "title:") {
			channelIds[i] = lookupChannelId(service, strings.Replace(channel, "title:", "", 1))
		} else {
			channelIds[i] = channel
		}
	}
	for _, channelId := range channelIds {
		inspectChannel(service, channelId)
	}
	return channelIds
}
