package main

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
)

func listVideos(service *youtube.Service, playlists ...Playlist) map[string]Video {
	videos := make(map[string]Video)
	for _, playlist := range playlists {
		request := service.PlaylistItems.List([]string{"snippet", "status"}).PlaylistId(playlist.Id).MaxResults(50)
		fetching := true
		for fetching {
			result, err := request.Do()
			onError(err, fmt.Sprintf("Could not get videos for %s", playlist))
			for _, item := range result.Items {
				status := item.Status.PrivacyStatus
				if "public" == status || "unlisted" == status {
					videoId := item.Snippet.ResourceId.VideoId
					videos[videoId] = Video{
						Id:       videoId,
						Title:    item.Snippet.Title,
						Playlist: playlist,
					}
				}
			}
			if 0 != len(result.NextPageToken) {
				request = request.PageToken(result.NextPageToken)
			} else {
				fetching = false
			}
		}
	}
	return videos
}

func determineVideosToAdd(service *youtube.Service, playlist MergedPlaylist) []Video {
	mine := listVideos(service, Playlist{
		Id:      playlist.Id,
		Title:   playlist.Title,
		Channel: Channel{"mine", "mine"},
	})
	sources := listVideos(service, playlist.Sources...)
	removed := false
	fmt.Println("Videos removed:")
	for id := range mine {
		if _, pres := sources[id]; !pres {
			fmt.Printf("\t - %s\n", mine[id])
			removed = true
		}
	}
	if !removed {
		fmt.Println("\t - None")
	}
	var added []Video
	for id := range sources {
		if _, pres := mine[id]; !pres {
			added = append(added, sources[id])
		}
	}
	return added
}
