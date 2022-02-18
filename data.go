package main

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
)

type Channel struct {
	Id    string
	Title string
}

func (c Channel) String() string {
	return fmt.Sprintf("%s :: %s", c.Id, c.Title)
}

func (c Channel) Unresolved() bool {
	return 0 == len(c.Id) || 0 == len(c.Title)
}

func (c Channel) FromSearchSnippet(snippet *youtube.SearchResultSnippet) Channel {
	c.Id = snippet.ChannelId
	c.Title = snippet.ChannelTitle
	return c
}

type SourcePlaylist struct {
	Id      string
	Title   string
	Channel Channel
}

func (s SourcePlaylist) String() string {
	return fmt.Sprintf("%s :: %s", s.Id, s.Title)
}

func (s SourcePlaylist) FromPlaylistSnippet(playlist *youtube.Playlist) SourcePlaylist {
	s.Id = playlist.Id
	s.Title = playlist.Snippet.Title
	s.Channel = Channel{
		Id:    playlist.Snippet.ChannelId,
		Title: playlist.Snippet.ChannelTitle,
	}
	return s
}

type MergedPlaylist struct {
	Id      string `yaml:",omitempty"`
	Title   string
	Sources []SourcePlaylist
}

func (m MergedPlaylist) String() string {
	return fmt.Sprintf("%s :: %s", m.Id, m.Title)
}

func (m MergedPlaylist) WithDetails(id string, title string) MergedPlaylist {
	m.Id = id
	m.Title = title
	return m
}

func (m MergedPlaylist) WithSources(sources []SourcePlaylist) MergedPlaylist {
	m.Sources = sources
	return m
}
