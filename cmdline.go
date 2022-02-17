package main

import "flag"

type Args struct {
	AuthOnly     bool
	SampleConfig bool
}

func parseArgs() Args {
	args := Args{
		AuthOnly:     false,
		SampleConfig: false,
	}

	flag.BoolVar(
		&args.AuthOnly, "auth-only", args.AuthOnly,
		"Just authenticate, nothing more",
	)
	flag.BoolVar(
		&args.SampleConfig, "sample-config", args.SampleConfig,
		"Generate a sample config, if one doesn't exist",
	)

	flag.Parse()

	return args
}
