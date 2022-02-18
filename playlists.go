package main

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
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

func listPlaylists(request *youtube.PlaylistsListCall) map[string]SourcePlaylist {
	request = request.MaxResults(50)
	playlists := make(map[string]SourcePlaylist)
	fetching := true
	for fetching {
		result, err := request.Do()
		onError(err, fmt.Sprintf("Could not get playlists"))
		for _, item := range result.Items {
			playlist := SourcePlaylist{}.FromPlaylistSnippet(item)
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

func resolveSourcePlaylists(service *youtube.Service, playlists []SourcePlaylist) []SourcePlaylist {
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
	resolved := make([]SourcePlaylist, 0, len(playlists))
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

func zipPlaylists(service *youtube.Service, playlists []MergedPlaylist) []MergedPlaylist {
	resolved := make([]MergedPlaylist, len(playlists))
	for i, playlist := range playlists {
		resolved[i] = MergedPlaylist{}.WithDetails(
			maybeCreatePlaylist(service, playlist.Id, playlist.Title),
		).WithSources(
			resolveSourcePlaylists(service, playlist.Sources),
		)
	}
	return resolved
}
