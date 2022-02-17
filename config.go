package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type MergedPlaylist struct {
	Id      string
	Sources []string `yaml:",omitempty"`
}

type Config struct {
	Channels  []string         `yaml:",omitempty"`
	Playlists []MergedPlaylist `yaml:",omitempty"`
}

func configFile() string {
	return filepath.Join(configDir(), "config.yaml")
}

func maybeWriteSampleConfig() {
	path := configFile() + ".sample"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		sample := Config{
			Channels: []string{"Usernames", "Of", "Channels", "To", "Watch"},
			Playlists: []MergedPlaylist{
				{
					Id:      "name:<name-of-new-playlist> # will be replaced",
					Sources: []string{"Youtube", "IDs", "Of", "Source", "Playlists"},
				},
				{
					Id:      "id-of-existing-playlist",
					Sources: []string{"Youtube", "IDs", "Of", "Source", "Playlists"},
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
	if 0 == len(config.Playlists) && 0 == len(config.Channels) {
		return errors.New("no playlists or channels configured")
	}
	// omitempty likely saves us from this
	for _, channel := range config.Channels {
		if 0 == len(channel) {
			return errors.New("empty channel found")
		}
	}
	for _, playlist := range config.Playlists {
		if 0 == len(playlist.Id) {
			return errors.New("playlist with empty id found")
		}
		if 0 == len(playlist.Sources) {
			return errors.New("playlist with no sources found")
		}
		for _, source := range playlist.Sources {
			// omitempty likely saves us from this
			if 0 == len(source) {
				return errors.New(fmt.Sprintf("playlist %s has empty sources", playlist.Id))
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
