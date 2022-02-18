package main

import "flag"

type Args struct {
	AuthOnly        bool
	InspectChannels bool
	SampleConfig    bool
}

func parseArgs() Args {
	args := Args{
		AuthOnly:        false,
		InspectChannels: false,
		SampleConfig:    false,
	}

	flag.BoolVar(
		&args.AuthOnly, "auth-only", args.AuthOnly,
		"Just authenticate, nothing more",
	)
	flag.BoolVar(
		&args.InspectChannels, "inspect-channels", args.InspectChannels,
		"Given the channels in the config, list their playlists",
	)
	flag.BoolVar(
		&args.SampleConfig, "sample-config", args.SampleConfig,
		"Generate a sample config, if one doesn't exist",
	)

	flag.Parse()

	return args
}
