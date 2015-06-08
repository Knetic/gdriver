package main

import (
	"fmt"
	"os"
	"errors"
	"path/filepath"
	"os/user"
)

func main() {

	var globals GlobalFlags
	var err error

	globals = ParseGlobalFlags()

	if(globals.Verb == VERB_UNKNOWN) {
		fmt.Fprintf(os.Stderr, "No action specified\n")
		return
	}

	if(globals.OAuthConfigPath == "") {
		err = loadConfig()
	} else {
		err = LoadOAuthConfig(globals.OAuthConfigPath)
	}

	if(err != nil) {

		msg := fmt.Sprintf("Unable to load configuration: %s\n", err)
		fmt.Fprintf(os.Stderr, msg)
		return
	}

	switch(globals.Verb) {
		case VERB_LIST: err = ListDriveFiles()
		case VERB_PUSH: err = PushDriveFile()
		default: err = errors.New("Verb not yet supported")
	}


	if err != nil {
		msg := fmt.Sprintf("%s\n", err)
		fmt.Fprintf(os.Stderr, msg)
		return
	}
}

/*
	Attempts to load the OAuth config from a few places, returning an error if none are valid.
	First, checks the current directory for "oauth.json",
	then checks "$HOME/.oauth/oauth.json".
*/
func loadConfig() (error) {

	var currentUser *user.User
	var homePath string
	var err error

	err = LoadOAuthConfig("oauth.json")
	if(err == nil) {
		return nil
	}

	currentUser, err = user.Current()
	if(err != nil) {
		fmt.Printf("Error getting user\n")
		return err
	}

	homePath = fmt.Sprintf("%s/.oauth/oauth.json", currentUser.HomeDir)
	homePath = filepath.FromSlash(homePath)

	err = LoadOAuthConfig(homePath)
	if(err == nil) {
		return nil
	}

	return errors.New("Unable to find valid 'oauth.json' in current directory nor ~/.oauth/oauth.json")
}
