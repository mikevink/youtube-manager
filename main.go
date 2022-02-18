package main

import (
	"log"
)

func main() {

	args := parseArgs()

	if args.SampleConfig {
		maybeWriteSampleConfig()
		return
	}

	service := getService()

	if args.AuthOnly {
		log.Println("OAuth successful")
		return
	}

	config := loadConfig()

	if 0 == len(config.Channels) && 0 == len(config.Playlists) {
		args.AddChannel = true
	}

	if args.AddChannel {
		config.Channels = addChannels(service, config.Channels)
	}

	if args.InspectChannels {
		inspectChannels(service, config.Channels)
	}

	saveConfig(config)
}
