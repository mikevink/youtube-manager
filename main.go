package main

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
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

func configFile() string {
	return filepath.Join(configDir(), "config.yaml")
}

func credentialFile() string {
	return filepath.Join(configDir(), ".credentials")
}

func tokenFile() string {
	return filepath.Join(configDir(), ".token")
}

func getCachedToken() (*oauth2.Token, error) {
	file, err := os.Open(tokenFile())
	if nil != err {
		return nil, err
	}
	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	return token, err
}

func saveToken(token *oauth2.Token) {
	fileName := tokenFile()
	log.Printf("Saving token to: %s\n", fileName)

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	onError(err, "Unable to cache oauth token")
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	onError(json.NewEncoder(file).Encode(token), "Unable to encode oauth token")
}

func getNewToken(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n\t%v\n", authURL)

	var code string
	_, err := fmt.Scan(&code)
	onError(err, "Unable to read authorization code")

	token, err := config.Exchange(context.TODO(), code)
	onError(err, "Unable to retrieve token from web")

	saveToken(token)

	return token
}

func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	token, err := getCachedToken()
	if nil != err {
		token = getNewToken(config)
	}
	return config.Client(ctx, token)
}

func listPlaylists(service *youtube.Service) {
	call := service.Playlists.List([]string{"snippet", "contentDetails"}).ChannelId("-").MaxResults(50)
	response, err := call.Do()
	onError(err, "")
	for _, list := range response.Items {
		log.Printf("Playlist - ID: %s Name: %s\n", list.Id, list.Snippet.Title)
	}
}

func main() {
	ctx := context.Background()

	bytes, err := ioutil.ReadFile(credentialFile())
	onError(err, "Unable to read client credentials file")

	config, err := google.ConfigFromJSON(bytes, youtube.YoutubeScope)
	onError(err, "Unable to parse client credentials file to config")

	client := getClient(ctx, config)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	onError(err, "Error creating YouTube client")

	listPlaylists(service)
}
