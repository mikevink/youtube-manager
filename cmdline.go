package main

import "flag"

type Args struct {
	AuthOnly        bool
	AddChannel      bool
	InspectChannels bool
	SampleConfig    bool
	Verbose         bool
}

func parseArgs() Args {
	args := Args{
		AuthOnly:        false,
		AddChannel:      false,
		InspectChannels: false,
		SampleConfig:    false,
		Verbose:         false,
	}

	flag.BoolVar(
		&args.AuthOnly, "auth-only", args.AuthOnly,
		"Just authenticate, nothing more",
	)
	flag.BoolVar(
		&args.AddChannel, "add-channel", args.AddChannel,
		"Add channel",
	)
	flag.BoolVar(
		&args.InspectChannels, "inspect-channels", args.InspectChannels,
		"Given the channels in the config, list their playlists",
	)
	flag.BoolVar(
		&args.SampleConfig, "sample-config", args.SampleConfig,
		"Generate a sample config, if one doesn't exist",
	)
	flag.BoolVar(
		&args.Verbose, "verbose", args.Verbose,
		"Verbose logging",
	)

	flag.Parse()

	return args
}
