package main

import (
	"code.google.com/p/google-api-go-client/drive/v2"
	"fmt"
	"errors"
)
/*
	Prints a list of all files present in the drive on stdout.
*/
func ListDriveFiles() (error) {

	var service *drive.Service
	var storedFiles []*drive.File
	var file *drive.File
	var err error

	service, err = createServiceClient()

	if err != nil {
		msg := fmt.Sprintf("Unable to authenticate with Drive: %s\n", err)
		return errors.New(msg)
	}

	storedFiles, err = RetrieveFileList(service)
	if(err != nil) {
		msg := fmt.Sprintf("Unable to pull list of files: %v\n", err)
		return errors.New(msg)
	}

	for _, file = range(storedFiles) {
		fmt.Printf("%s\n", file.Title)
	}

	return nil
}
