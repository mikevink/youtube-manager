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
	return fmt.Sprintf("%s #%s", c.Title, c.Id)
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
	return fmt.Sprintf("%s #%s [%s]", s.Title, s.Id, s.Channel)
}

type MergedPlaylist struct {
	Id      string `yaml:",omitempty"`
	Title   string
	Sources []SourcePlaylist
}

func (m MergedPlaylist) String() string {
	return fmt.Sprintf("%s #%s", m.Title, m.Id)
}
