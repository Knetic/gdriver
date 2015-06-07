package main

import (
	"fmt"
	"net/http"
	"os"
	"code.google.com/p/google-api-go-client/drive/v2"
	"code.google.com/p/goauth2/oauth"
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

	service = createServiceClient()
	if(service == nil) {
		return
	}

	files, _ := AllFiles(service)

	for file, _ := range(files) {

		fmt.Printf("File: %v\n", file)
	}
}

func createServiceClient() (*drive.Service) {

	t := &oauth.Transport{
		Config: config,
		Transport: http.DefaultTransport,
	}

	authenticateTransport(t)

	// Create a new authorized Drive client.
	svc, err := drive.New(t.Client())
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred creating Drive client: %v\n", err)
		return nil
	}

	return svc
}

func authenticateTransport(transport *oauth.Transport) {

	// cached?
	var tokenCache oauth.CacheFile

	tokenCache = "token.json"

	if _, err := os.Stat("token.json"); !os.IsNotExist(err) {

		token, err := tokenCache.Token()
		if(err != nil) {
			fmt.Fprintf(os.Stderr, "Unable to read token: %s\n", err)
			return
		}
		if(token != nil) {
			fmt.Printf("Using token: '%v'\n\n", token)
			transport.Token = token
			return
		}
	}

	// not cached, prompt user.
	authUrl := config.AuthCodeURL("state")
	fmt.Printf("Go to the following link in your browser: %v\n", authUrl)

	fmt.Printf("\nEnter verification code: ")
	var code string
	fmt.Scanln(&code)

	freshToken, err := transport.Exchange(code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred exchanging the code: %v\n", err)
		return
	}

	tokenCache.PutToken(freshToken)
}

// AllFiles fetches and displays all files
func AllFiles(d *drive.Service) ([]*drive.File, error) {
  var fs []*drive.File
  pageToken := ""
  for {
    q := d.Files.List()
    // If we have a pageToken set, apply it to the query
    if pageToken != "" {
      q = q.PageToken(pageToken)
    }
    r, err := q.Do()
    if err != nil {
      fmt.Printf("An error occurred: %v\n", err)
      return fs, err
    }
    fs = append(fs, r.Items...)
    pageToken = r.NextPageToken
    if pageToken == "" {
      break
    }
  }
  return fs, nil
}
