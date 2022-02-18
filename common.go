package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func fclose(file *os.File) {
	_ = file.Close()
}

func input(reader *bufio.Reader, prompt string) string {
	fmt.Printf("%s: ", prompt)
	r, err := reader.ReadString('\n')
	onError(err, "Could not process user input")
	return strings.TrimSpace(r)
}

func quittingInput(reader *bufio.Reader, prompt string) string {
	r := input(reader, fmt.Sprintf("%s ([iq]uit to exit)", prompt))
	if "iquit" == r || "iq" == r {
		os.Exit(0)
	}
	return r
}
