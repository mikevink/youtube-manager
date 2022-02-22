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

	if !args.AddChannel && !args.InspectChannels {
		if args.ResolvePlaylists {
			config.Playlists = resolvePlaylists(service, config.Playlists)
		} else {
			config.Playlists = zipPlaylists(service, config.Playlists, config.compileExcludes(), args.Verbose, args.Noop)
		}
	}

	saveConfig(config)
}
