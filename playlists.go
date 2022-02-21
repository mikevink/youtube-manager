package main

import (
	"bufio"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"os"
)

func maybeCreatePlaylist(service *youtube.Service, id string, title string) (string, string) {
	if 0 != len(id) {
		listed := listPlaylists(service.Playlists.List([]string{"snippet"}).Id(id))
		if 0 != len(listed) {
			return id, listed[id].Title
		}
	}
	result, err := service.Playlists.Insert(
		[]string{"snippet", "status"},
		&youtube.Playlist{
			Snippet: &youtube.PlaylistSnippet{Title: title},
			Status:  &youtube.PlaylistStatus{PrivacyStatus: "private"},
		},
	).Do()
	onError(err, fmt.Sprintf("Could not create playlist %s", title))
	fmt.Printf("Created playlist %s: https://www.youtube.com/playlist?list=%s\n", result.Snippet.Title, result.Id)
	return result.Id, result.Snippet.Title
}

func listPlaylists(request *youtube.PlaylistsListCall) map[string]Playlist {
	request = request.MaxResults(50)
	playlists := make(map[string]Playlist)
	fetching := true
	for fetching {
		result, err := request.Do()
		onError(err, fmt.Sprintf("Could not get playlists"))
		for _, item := range result.Items {
			playlist := Playlist{}.FromPlaylistSnippet(item)
			playlists[playlist.Id] = playlist
		}
		if 0 != len(result.NextPageToken) {
			request = request.PageToken(result.NextPageToken)
		} else {
			fetching = false
		}
	}
	return playlists
}

func resolveSourcePlaylists(service *youtube.Service, playlists []Playlist) []Playlist {
	var ids []string
	for _, playlist := range playlists {
		if 0 == len(playlist.Title) || playlist.Channel.Unresolved() {
			ids = append(ids, playlist.Id)
		}
	}
	if 0 == len(ids) {
		return playlists
	}
	listed := listPlaylists(service.Playlists.List([]string{"snippet"}).Id(ids...))
	resolved := make([]Playlist, 0, len(playlists))
	for _, playlist := range playlists {
		list, prs := listed[playlist.Id]
		if prs {
			resolved = append(resolved, list)
		} else {
			resolved = append(resolved, playlist)
		}
	}
	return resolved
}

func zipPlaylist(service *youtube.Service, playlist MergedPlaylist, verbose bool) {
	reader := bufio.NewReader(os.Stdin)
	videos := determineVideosToAdd(service, playlist, verbose)
	fmt.Println("Adding videos:")
	if 0 == len(videos) {
		fmt.Println("\t - None")
		return
	}
	total := len(videos)
	for i, video := range videos {
		fmt.Printf("\t - %s\n", video)
		_, err := service.PlaylistItems.Insert(
			[]string{"snippet"},
			&youtube.PlaylistItem{
				Snippet: &youtube.PlaylistItemSnippet{
					PlaylistId: playlist.Id,
					ResourceId: &youtube.ResourceId{
						Kind:    "youtube#video",
						VideoId: video.Id,
					},
				},
			},
		).Do()
		onError(err, fmt.Sprintf("Could not add video %s", video))
		fmt.Printf("\t   Done %d/%d\n", i, total)
		if verbose {
			_ = quittingInput(reader, "Continue? [y]")
		}
	}
}

func zipPlaylists(service *youtube.Service, playlists []MergedPlaylist, verbose bool) []MergedPlaylist {
	resolved := make([]MergedPlaylist, 0, len(playlists))
	for _, playlist := range playlists {
		merged := MergedPlaylist{}.WithDetails(
			maybeCreatePlaylist(service, playlist.Id, playlist.Title),
		).WithSources(
			resolveSourcePlaylists(service, playlist.Sources),
		)
		zipPlaylist(service, merged, verbose)
		resolved = append(resolved, merged)
	}
	return resolved
}

func resolvePlaylists(service *youtube.Service, playlists []MergedPlaylist) []MergedPlaylist {
	resolved := make([]MergedPlaylist, 0, len(playlists))
	for _, playlist := range playlists {
		merged := MergedPlaylist{}.WithDetails(
			maybeCreatePlaylist(service, playlist.Id, playlist.Title),
		).WithSources(
			resolveSourcePlaylists(service, playlist.Sources),
		)
		resolved = append(resolved, merged)
	}
	return resolved
}
