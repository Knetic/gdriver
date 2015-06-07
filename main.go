package main

import (
	"fmt"
	"net/http"
	"os"
	"code.google.com/p/google-api-go-client/drive/v2"
	"code.google.com/p/goauth2/oauth"
	"errors"
)

var config = &oauth.Config{
	ClientId: "179778203598-f4ntihkomqs6c4jbeehadpil35sfv8ea.apps.googleusercontent.com",
	ClientSecret: "KOknojIaqFBMG1EDI8ht-ozR",
	Scope:"https://www.googleapis.com/auth/drive",
	RedirectURL:"urn:ietf:wg:oauth:2.0:oob",
	AuthURL:"https://accounts.google.com/o/oauth2/auth",
	TokenURL: "https://accounts.google.com/o/oauth2/token",

	AccessType: "offline",
}

func main() {

	var service *drive.Service
	var storedFiles []*drive.File
	var err error

	service, err = createServiceClient()

	if(err != nil) {
		msg := fmt.Sprintf("Unable to authenticate with Drive: %s\n", err)
		fmt.Fprintf(os.Stderr, msg)
		return
	}

	storedFiles, err = retrieveBackupFileList(service, "backups")

	if(err != nil) {
		msg := fmt.Sprintf("Unable to get list of files: %s\n", err)
		fmt.Fprintf(os.Stderr, msg)
		return
	}

	for _, file := range(storedFiles) {

		fmt.Printf("%d: %s\n", file.FileSize, file.Title)
	}
}

func createServiceClient() (*drive.Service, error) {

	var transport *oauth.Transport
	var service *drive.Service
	var err error

	transport = &oauth.Transport{
		Config: config,
		Transport: http.DefaultTransport,
	}

	authenticateTransport(transport)

	// Create a new authorized Drive client.
	service, err = drive.New(transport.Client())
	if err != nil {
		msg := fmt.Sprintf("Unable to create drive client: %s\n", err)
		return nil, errors.New(msg)
	}

	return service, nil
}

func authenticateTransport(transport *oauth.Transport) (error) {

	var tokenCache oauth.CacheFile
	var token *oauth.Token
	var verificationCode string
	var err error

	tokenCache = "token.json"

	// try to read cached token
	if _, err := os.Stat("token.json"); !os.IsNotExist(err) {

		token, err = tokenCache.Token()

		if(err != nil) {
			msg := fmt.Sprintf("Unable to read token: %s\n", err)
			return errors.New(msg)
		}

		transport.Token = token
		return nil
	}

	// not cached, prompt user.
	authUrl := config.AuthCodeURL("state")

	fmt.Printf("Go to the following link in your browser: \n%v\n\n", authUrl)
	fmt.Printf("Enter verification code: ")

	fmt.Scanln(&verificationCode)

	token, err = transport.Exchange(verificationCode)

	if err != nil {
		msg := fmt.Sprintf("An error occurred exchanging the code: %v\n", err)
		return errors.New(msg)
	}

	tokenCache.PutToken(token)
	return nil
}

func retrieveBackupFileList(service *drive.Service, path string) ([]*drive.File, error) {

	var ret []*drive.File
	var listQuery *drive.FilesListCall
	var files *drive.FileList
	var pageToken string
	var err error

	pageToken = ""

	for {
		listQuery = service.Files.List()

		// if we're on a new page, use it in the query
		if(pageToken != "") {
			listQuery = listQuery.PageToken(pageToken)
		}

		files, err = listQuery.Do()
		if(err != nil) {
			return nil, err
		}

		ret = append(ret, files.Items...)
		pageToken = files.NextPageToken

		if(pageToken == "") {
			break
		}
	}

	return ret, nil
}
