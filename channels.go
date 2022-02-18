package main

import (
	"bufio"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
	"strings"
)

func lookupChannel(service *youtube.Service, reader *bufio.Reader) Channel {
	title := quittingInput(reader, "Enter channel title")
	request := service.Search.List([]string{"snippet"}).Type("channel").MaxResults(10).Q(title)
	for true {
		result, err := request.Do()
		onError(err, fmt.Sprintf("Could not search for channel '%s'", title))
		fmt.Println("Please choose an option in []:")
		options := make([]Channel, len(result.Items))
		for i, item := range result.Items {
			options[i] = Channel{}.FromSearchSnippet(item.Snippet)
			description := ""
			if 0 != len(item.Snippet.Description) {
				description = fmt.Sprintf("\t     %s\n", item.Snippet.Description)
			}
			fmt.Printf(
				"\t[%2d] %s\n%s\t     https://www.youtube.com/channel/%s\n",
				i, options[i], description, options[i].Id,
			)
		}
		if 0 != len(result.PrevPageToken) {
			fmt.Println("\t[ p] Prev")
		}
		if 0 != len(result.NextPageToken) {
			fmt.Println("\t[ n] Next")
		}
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
			return options[atoi(option)]
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

func getChannelPlaylists(service *youtube.Service, channel Channel) []Playlist {
	listed := listPlaylists(service.Playlists.List([]string{"snippet"}).ChannelId(channel.Id))
	playlists := make([]Playlist, 0, len(listed))
	for _, playlist := range listed {
		playlists = append(playlists, playlist)
	}
	return playlists
}

func inspectChannel(service *youtube.Service, channel Channel) {
	fmt.Printf("Channel %s:\n", channel)
	playlists := getChannelPlaylists(service, channel)
	if 0 == len(playlists) {
		fmt.Println("\t- No playlists")
		return
	}
	for _, playlist := range playlists {
		fmt.Printf("\t- %s\n\t  https://www.youtube.com/playlist?list=%s\n", playlist, playlist.Id)
	}
}

func inspectChannels(service *youtube.Service, channels []Channel) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Which channel(s) do you want to inspect?")
	for i, channel := range channels {
		fmt.Printf("\t[%d] %s\n", i, channel)
	}
	fmt.Println("\t[*] All\n\tc,s,v accepted")
	choice := quittingInput(reader, "Option [*]")
	if 0 == len(choice) {
		choice = "*"
	}
	var choices []Channel
	if strings.Contains(choice, "*") {
		choices = channels
	} else {
		inxs := strings.Split(choice, ",")
		choices = make([]Channel, len(inxs))
		for i, inx := range inxs {
			choices[i] = channels[atoi(inx)]
		}
	}
	for _, channel := range choices {
		inspectChannel(service, channel)
	}
}
