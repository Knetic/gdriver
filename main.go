package main

import (
	"fmt"
	"os"
	"errors"
	"os/user"
)

func main() {

	var globals GlobalFlags;
	var verb OperationVerb
	var err error

	if(len(os.Args) <= 1) {
		fmt.Fprintf(os.Stderr, "No verb specified\n")
		return
	}

	globals = ParseGlobalFlags()

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

	verb, err = ParseOperationVerb(os.Args[1])
	if(err != nil) {
		msg := fmt.Sprintf("%s\n", err)
		fmt.Fprintf(os.Stderr, msg)
		return
	}

	switch(verb) {
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
		return err
	}

	homePath = fmt.Sprintf("%s/.oauth/oauth.json", currentUser.HomeDir)

	err = LoadOAuthConfig(homePath)
	if(err == nil) {
		return nil
	}

	return errors.New("Unable to find valid 'oauth.json' in current directory nor ~/.oauth/oauth.json")
}
