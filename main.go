package main

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/drive/v2"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
)

var config = &oauth.Config{
	ClientId:     "179778203598-f4ntihkomqs6c4jbeehadpil35sfv8ea.apps.googleusercontent.com",
	ClientSecret: "KOknojIaqFBMG1EDI8ht-ozR",
	Scope:        "https://www.googleapis.com/auth/drive",
	RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
	AuthURL:      "https://accounts.google.com/o/oauth2/auth",
	TokenURL:     "https://accounts.google.com/o/oauth2/token",

	AccessType: "offline",
}

func main() {

	var err error

	err = uploadFile("2015_06_07_codeBackup.tar.xz", "stagger")
	if err != nil {
		msg := fmt.Sprintf("%s", err)
		fmt.Fprintf(os.Stderr, msg)
	}
}

func uploadFile(sourceFilePath string, parentFolderName string) error {

	var storedFiles []*drive.File
	var service *drive.Service
	var parentFolderId string
	var err error

	service, err = createServiceClient()

	if err != nil {
		msg := fmt.Sprintf("Unable to authenticate with Drive: %s\n", err)
		return errors.New(msg)
	}

	// if the user wants a parent folder, find its id.
	if parentFolderName != "" {

		storedFiles, err = retrieveFileList(service)
		if err != nil {
			msg := fmt.Sprintf("Unable to get list of files from Drive: %s\n", err)
			return errors.New(msg)
		}

		parentFolderId = findFileId(storedFiles, parentFolderName)

		if parentFolderId == "" {

			msg := fmt.Sprintf("Unable to find parent folder named '%s'\n", parentFolderName)
			return errors.New(msg)
		}
	} else {
		parentFolderId = ""
	}

	return uploadLocalFile(service, sourceFilePath, parentFolderId)
}

func createServiceClient() (*drive.Service, error) {

	var transport *oauth.Transport
	var service *drive.Service
	var err error

	transport = &oauth.Transport{
		Config:    config,
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

func authenticateTransport(transport *oauth.Transport) error {

	var tokenCache oauth.CacheFile
	var token *oauth.Token
	var verificationCode string
	var err error

	tokenCache = "token.json"

	// try to read cached token
	if _, err := os.Stat("token.json"); !os.IsNotExist(err) {

		token, err = tokenCache.Token()

		if err != nil {
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

func retrieveFileList(service *drive.Service) ([]*drive.File, error) {

	var ret []*drive.File
	var listQuery *drive.FilesListCall
	var files *drive.FileList
	var pageToken string
	var err error

	pageToken = ""

	for {
		listQuery = service.Files.List()

		// if we're on a new page, use it in the query
		if pageToken != "" {
			listQuery = listQuery.PageToken(pageToken)
		}

		files, err = listQuery.Do()
		if err != nil {
			return nil, err
		}

		ret = append(ret, files.Items...)
		pageToken = files.NextPageToken

		if pageToken == "" {
			break
		}
	}

	return ret, nil
}

func findFileId(storedFiles []*drive.File, fileName string) string {

	var file *drive.File

	for _, file = range storedFiles {

		if file.Title == fileName {
			return file.Id
		}
	}

	return ""
}

func uploadLocalFile(service *drive.Service, sourceFilePath string, parentFolderId string) error {

	// upload.
	var fileName string
	var mimeType string

	m, err := os.Open(sourceFilePath)

	if err != nil {
		return err
	}

	fileName = path.Base(sourceFilePath)
	mimeType = determineMimeType(fileName)

	f := &drive.File{
		Title:    fileName,
		MimeType: mimeType,
	}

	if parentFolderId != "" {
		p := &drive.ParentReference{
			Id: parentFolderId,
		}
		f.Parents = []*drive.ParentReference{p}
	}

	_, err = service.Files.Insert(f).Media(m).Do()
	if err != nil {
		return err
	}
	return nil
}

/*
	Figures out and returns the mime type of the given [filePath].
	This generally defers to mime.TypeByExtension, but if this is a "*.tar.*"
	archive, it will return the mime type for "compressed".
*/
func determineMimeType(filePath string) string {

	if strings.Contains(filePath, ".tar.") {
		return "application/x-gzip"
	}

	return mime.TypeByExtension(path.Ext(filePath))
}
