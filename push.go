package main

import (
	"flag"
	"errors"
)

func PushDriveFile() (error) {

	var filePath = flag.String("f", "", "The local path to a file to be uploaded")
	var parentFolderName = flag.String("p", "", "The name of a parent folder")

	flag.Parse()

	if(filePath == nil || *filePath == "") {
		return errors.New("File path not specified (-f)")
	}

	return UploadFile(*filePath, *parentFolderName)
}
