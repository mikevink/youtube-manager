package main

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
)

func listVideos(service *youtube.Service, verbose bool, playlists ...Playlist) map[string]Video {
	videos := make(map[string]Video)
	if verbose {
		fmt.Println("Finding videos")
	}
	for _, playlist := range playlists {
		request := service.PlaylistItems.List([]string{"snippet", "status"}).PlaylistId(playlist.Id).MaxResults(50)
		fetching := true
		for fetching {
			result, err := request.Do()
			onError(err, fmt.Sprintf("Could not get videos for %s", playlist))
			for _, item := range result.Items {
				status := item.Status.PrivacyStatus
				videoId := item.Snippet.ResourceId.VideoId
				video := Video{
					Id:       videoId,
					Title:    item.Snippet.Title,
					Playlist: playlist,
				}
				if "public" == status || "unlisted" == status {
					if _, prs := videos[videoId]; !prs {
						videos[videoId] = video
					} else {
						if verbose {
							fmt.Printf("\t - Skipping [status=duplicate] %s\n", video)
						}
					}
				} else {
					if verbose {
						fmt.Printf("\t - Skipping [status=%s] %s\n", status, video)
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
	if verbose {
		fmt.Println("\t - Done")
	}
	return videos
}

func determineVideosToAdd(service *youtube.Service, playlist MergedPlaylist, verbose bool) []Video {
	mine := listVideos(service, verbose, Playlist{
		Id:      playlist.Id,
		Title:   playlist.Title,
		Channel: Channel{"mine", "mine"},
	})
	sources := listVideos(service, verbose, playlist.Sources...)
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
