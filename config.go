package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	Channels  []Channel        `yaml:",omitempty"`
	Playlists []MergedPlaylist `yaml:",omitempty"`
	Exclude   []string         `yaml:",omitempty"`
}

func (c Config) compileExcludes() []*regexp.Regexp {
	var expressions []*regexp.Regexp
	for _, exp := range c.Exclude {
		rexp, err := regexp.Compile(exp)
		onError(err, fmt.Sprintf("Could not parse expression '%s'", exp))
		expressions = append(expressions, rexp)
	}
	return expressions
}

func configFile() string {
	return filepath.Join(configDir(), "config.yaml")
}

func maybeWriteSampleConfig() {
	path := configFile() + ".sample"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		sample := Config{
			Channels: []Channel{
				{"ChannelId", "Channel Title"},
			},
			Playlists: []MergedPlaylist{
				{
					Id:    "merged playlist id. can be empty, as will fill / create",
					Title: "title of playlist",
					Sources: []Playlist{
						{"PlaylistId", "PlaylistTitle", Channel{"ChannelId", "ChannelTitle"}},
					},
				},
			},
		}
		data, err := yaml.Marshal(sample)
		onError(err, "Could not marshal sample config")
		onError(ioutil.WriteFile(path, data, 0600), "Could not write sample config")
		log.Printf("Sample config file written at: %s\n", path)
	} else if nil == err {
		log.Printf("Sample config already exists at: %s\n", path)
	} else {
		onError(err, "Writing sample config failed miserably :(")
	}
}

func validate(config Config) error {
	for _, channel := range config.Channels {
		if 0 == len(channel.Id) {
			return errors.New("channel found with empty id")
		}
		if 0 == len(channel.Title) {
			return errors.New("channel found with empty title")
		}
	}
	for _, playlist := range config.Playlists {
		if 0 == len(playlist.Sources) {
			return errors.New("playlist with no sources found")
		}
		for _, source := range playlist.Sources {
			if 0 == len(source.Id) {
				return errors.New(fmt.Sprintf("source playlist for %s with empty id found", playlist.Title))
			}
		}
	}
	return nil
}

func loadConfig() Config {
	yml, err := ioutil.ReadFile(configFile())
	onError(err, "Could not read config file")
	config := Config{}
	onError(yaml.Unmarshal(yml, &config), "Could not parse config file")
	onError(validate(config), "Config not valid")
	return config
}

func saveConfig(config Config) {
	data, err := yaml.Marshal(config)
	onError(err, "Could not marshal config")
	onError(ioutil.WriteFile(configFile(), data, 0600), "Could not write config")
}
