package main

import (
	"errors"
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
