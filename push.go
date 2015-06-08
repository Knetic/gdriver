package main

import (
	"flag"
	"errors"
)

var filePath = flag.String("f", "", "The local path to a file to be uploaded")
var parentFolderName = flag.String("p", "", "The name of a parent folder")

func PushDriveFile() (error) {

	// flags are already parsed by this time
	if(filePath == nil || *filePath == "") {
		return errors.New("File path not specified (-f)")
	}

	return UploadFile(*filePath, *parentFolderName)
}
