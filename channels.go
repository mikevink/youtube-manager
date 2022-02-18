package main

import (
	"bufio"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
	"strconv"
)

//func lookupChannelId(service *youtube.Service, title string) string {
//	search, err := service.Search.List([]string{"snippet"}).Type("channel").Q(title).Do()
//	onError(err, fmt.Sprintf("Failed to search for channel '%s'", title))
//	for _, item := range search.Items {
//		if "youtube#channel" == item.Id.Kind && title == item.Snippet.Title {
//			return item.Snippet.ChannelId
//		}
//	}
//	log.Fatalf("Could not find channel titled '%s' is first 25 results", title)
//	return ""
//}

func getOptions(title string, request *youtube.SearchListCall) (string, []Channel, *youtube.SearchListResponse) {
	result, err := request.Do()
	onError(err, fmt.Sprintf("Could not search for channel '%s'", title))
	prompt := "Please choose an option in []:"
	options := make([]Channel, len(result.Items))
	for i, item := range result.Items {
		options[i] = Channel{}.FromSearchSnippet(item.Snippet)
		prompt = fmt.Sprintf(
			"%s\n\t[%d] %s - Description: %s",
			prompt, i, options[i], item.Snippet.Description,
		)
	}
	return prompt, options, result
}

func lookupChannel(service *youtube.Service, reader *bufio.Reader) Channel {
	title := quittingInput(reader, "Enter channel title")
	request := service.Search.List([]string{"snippet"}).Type("channel").MaxResults(10).Q(title)
	for true {
		prompt, options, result := getOptions(title, request)
		if 0 != len(result.PrevPageToken) {
			prompt = fmt.Sprintf("%s\n\t[p] Prev", prompt)
		}
		if 0 != len(result.NextPageToken) {
			prompt = fmt.Sprintf("%s\n\t[n] Next", prompt)
		}
		log.Println(prompt)
		option := quittingInput(reader, "Option")
		if "n" == option {
			if 0 != len(result.NextPageToken) {
				request = request.PageToken(result.NextPageToken)
			} else {
				log.Fatalln("No next page available, bailing")
			}
		} else if "p" == option {
			if 0 != len(result.PrevPageToken) {
				request = request.PageToken(result.PrevPageToken)
			} else {
				log.Fatalln("No prev page available, bailing")
			}
		} else {
			inx, err := strconv.Atoi(option)
			onError(err, "Atoi failed")
			return options[inx]
		}
	}
	return Channel{}
}

func addChannels(service *youtube.Service, channels []Channel) []Channel {
	reader := bufio.NewReader(os.Stdin)
	option := "y"
	for "y" == option {
		result := lookupChannel(service, reader)
		isNew := true
		for _, channel := range channels {
			if result.Id == channel.Id {
				isNew = false
				break
			}
		}
		if isNew {
			channels = append(channels, result)
		} else {
			fmt.Printf("Channel '%s' is already in your list\n", result)
		}
		option = quittingInput(reader, "Do you want to add more channels? [y]/n")
		if 0 == len(option) {
			option = "y"
		}
	}
	return channels
}

//func getPlaylists(service *youtube.Service, channelId string) []*youtube.Playlist {
//	request := service.Playlists.List([]string{"snippet"}).MaxResults(50).ChannelId(channelId)
//	var playlists []*youtube.Playlist
//	fetching := true
//	for fetching {
//		result, err := request.Do()
//		onError(err, fmt.Sprintf("Could not get playlists for channel 'https://www.youtube.com/channel/%s'", channelId))
//		playlists = append(playlists, result.Items...)
//		if 0 != len(result.NextPageToken) {
//			request = request.PageToken(result.NextPageToken)
//		} else {
//			fetching = false
//		}
//	}
//	return playlists
//}

//func inspectChannel(service *youtube.Service, channelId string) {
//	playlists := getPlaylists(service, channelId)
//	if 0 == len(playlists) {
//		log.Printf("No playlists found for channel '%s'\n", channelId)
//		return
//	}
//	channelTitle := playlists[0].Snippet.ChannelTitle
//	msg := fmt.Sprintf("Channel %s(%s):\n", channelTitle, channelId)
//	for _, playlist := range playlists {
//		msg = fmt.Sprintf("%s\t%s -- %s\n", msg, playlist.Snippet.Title, playlist.Id)
//	}
//	log.Println(msg)
//}

//func inspectChannels(service *youtube.Service, channels []Channel) []Channel {
//	channelIds := make([]string, len(channels))
//	for i, channel := range channels {
//		if strings.HasPrefix(channel, "title:") {
//			channelIds[i] = lookupChannelId(service, strings.Replace(channel, "title:", "", 1))
//		} else {
//			channelIds[i] = channel
//		}
//	}
//	for _, channelId := range channelIds {
//		inspectChannel(service, channelId)
//	}
//	return channelIds
//}
