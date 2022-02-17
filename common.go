package main

import (
	"log"
	"os"
	"path/filepath"
)

const ApiError = "Error making API call"

func onError(err error, message string) {
	if nil != err {
		if "" == message {
			message = ApiError
		}
		log.Fatalf(message+": %v", err.Error())
	}
}

func configDir() string {
	userConfigDir, err := os.UserConfigDir()
	onError(err, "Could not get user config dir")
	configDir := filepath.Join(userConfigDir, "youtube_manager")
	onError(os.MkdirAll(configDir, 0700), "Could not create config dir at "+configDir)
	return configDir
}
