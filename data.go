package main

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
	"time"
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

type Playlist struct {
	Id      string
	Title   string
	Channel Channel `yaml:",flow"`
}

func (s Playlist) String() string {
	return fmt.Sprintf("%s :: %s", s.Id, s.Title)
}

func (s Playlist) FromPlaylistSnippet(playlist *youtube.Playlist) Playlist {
	s.Id = playlist.Id
	s.Title = playlist.Snippet.Title
	s.Channel = Channel{
		Id:    playlist.Snippet.ChannelId,
		Title: playlist.Snippet.ChannelTitle,
	}
	return s
}

type Video struct {
	Id       string
	Title    string
	Playlist Playlist
}

func (v Video) String() string {
	return fmt.Sprintf("%s :: %s :: %s", v.Title, v.Playlist.Title, v.Playlist.Channel.Title)
}

type MergedPlaylist struct {
	Id      string `yaml:",omitempty"`
	Title   string
	After   string `yaml:",omitempty"`
	Sources []Playlist
}

func (m MergedPlaylist) String() string {
	return fmt.Sprintf("%s :: %s", m.Id, m.Title)
}

func (m MergedPlaylist) WithDetails(id string, title string, after string) MergedPlaylist {
	m.Id = id
	m.Title = title
	m.After = after
	return m
}

func (m MergedPlaylist) WithSources(sources []Playlist) MergedPlaylist {
	m.Sources = sources
	return m
}

func (m MergedPlaylist) PublishedAfter() time.Time {
	if 0 == len(m.After) {
		return time.UnixMilli(0)
	}
	tm, err := time.Parse("2006-01-02", m.After)
	onError(err, fmt.Sprintf("Could not parse time %s for playlist %v", m.After, m))
	return tm
}
