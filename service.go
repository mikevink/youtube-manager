package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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
	defer fclose(file)
	return token, err
}

func saveToken(token *oauth2.Token) {
	path := tokenFile()
	log.Printf("Saving token to: %s\n", path)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	onError(err, "Unable to cache oauth token")
	defer fclose(file)

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

func getService() *youtube.Service {
	ctx := context.Background()

	bytes, err := ioutil.ReadFile(credentialFile())
	onError(err, "Unable to read client credentials file")

	config, err := google.ConfigFromJSON(bytes, youtube.YoutubeScope)
	onError(err, "Unable to parse client credentials file to config")

	client := getClient(ctx, config)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	onError(err, "Error creating YouTube client")

	return service
}
