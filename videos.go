package main

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
	"regexp"
	"strings"
	"time"
)

func maybeAddVideo(
	status string, video Video, videos map[string]Video, videoId string, exclude []*regexp.Regexp, verbose bool,
) map[string]Video {
	if "public" == status || "unlisted" == status {
		filterTitle := strings.ToLower(video.Title)
		if _, prs := videos[videoId]; !prs {
			good := true
			if 0 != len(exclude) {
				for _, exc := range exclude {
					if exc.MatchString(filterTitle) {
						good = false
						break
					}
				}
			}
			if good {
				videos[videoId] = video
			} else {
				if verbose {
					fmt.Printf("\t - Skipping [status=excluded] %s\n", video)
				}
			}
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
	return videos
}

func listPlaylistVideos(
	service *youtube.Service, verbose bool, exclude []*regexp.Regexp, publishedAfter time.Time, playlist Playlist,
	videos map[string]Video,
) map[string]Video {
	request := service.PlaylistItems.List([]string{"snippet", "status"}).PlaylistId(playlist.Id).MaxResults(50)
	fetching := true
	for fetching {
		result, err := request.Do()
		onError(err, fmt.Sprintf("Could not get videos for %s", playlist))
		for _, item := range result.Items {
			status := item.Status.PrivacyStatus
			videoId := item.Snippet.ResourceId.VideoId
			tm, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
			onError(err, fmt.Sprintf("Could not parse time for %s, from %v", item.Snippet.PublishedAt, item))
			video := Video{
				Id:       videoId,
				Title:    item.Snippet.Title,
				Playlist: playlist,
			}
			if tm.Before(publishedAfter) {
				if verbose {
					fmt.Printf("\t - Skipping [status=publishedEarly] %s\n", video)
				}
			}
			videos = maybeAddVideo(status, video, videos, videoId, exclude, verbose)
		}
		if 0 != len(result.NextPageToken) {
			request = request.PageToken(result.NextPageToken)
		} else {
			fetching = false
		}
	}
	return videos
}

func listChannelVideos(
	service *youtube.Service, verbose bool, exclude []*regexp.Regexp, publishedAfter time.Time, playlist Playlist,
	videos map[string]Video,
) map[string]Video {
	request := service.Search.List([]string{"snippet"}).ChannelId(playlist.Channel.Id).Type("video").PublishedAfter(
		publishedAfter.Format(time.RFC3339),
	).MaxResults(50)
	fetching := true
	for fetching {
		result, err := request.Do()
		onError(err, fmt.Sprintf("Could not get videos for %s", playlist.Channel))
		for _, item := range result.Items {
			videoId := item.Id.VideoId
			video := Video{
				Id:       videoId,
				Title:    item.Snippet.Title,
				Playlist: playlist,
			}
			videos = maybeAddVideo("public", video, videos, videoId, exclude, verbose)
		}
		if 0 != len(result.NextPageToken) {
			request = request.PageToken(result.NextPageToken)
		} else {
			fetching = false
		}
	}
	return videos
}

func listVideos(
	service *youtube.Service, verbose bool, exclude []*regexp.Regexp, publishedAfter time.Time, playlists ...Playlist,
) map[string]Video {
	videos := make(map[string]Video)
	if verbose {
		fmt.Println("Finding videos")
	}
	for _, playlist := range playlists {
		if "fetch-from-channel" == playlist.Id {
			videos = listChannelVideos(service, verbose, exclude, publishedAfter, playlist, videos)
		} else {
			videos = listPlaylistVideos(service, verbose, exclude, publishedAfter, playlist, videos)
		}
	}
	if verbose {
		fmt.Println("\t - Done")
	}
	return videos
}

func determineVideosToAdd(service *youtube.Service, playlist MergedPlaylist, exclude []*regexp.Regexp, verbose bool) []Video {
	mine := listVideos(service, verbose, []*regexp.Regexp{}, time.UnixMilli(0), Playlist{
		Id:      playlist.Id,
		Title:   playlist.Title,
		Channel: Channel{"mine", "mine"},
	})
	sources := listVideos(service, verbose, exclude, playlist.PublishedAfter(), playlist.Sources...)
	removed := false
	if 0 == len(playlist.After) {
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
	}
	var added []Video
	for id := range sources {
		if _, pres := mine[id]; !pres {
			added = append(added, sources[id])
		}
	}
	return added
}
