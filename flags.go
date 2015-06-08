package main

import (
	"fmt"
	"flag"
	"errors"
	"os/user"
)

/*
	Describes the operations that can be performed with this tool.
*/
type OperationVerb int
const (

	VERB_UNKNOWN OperationVerb = iota
	VERB_LIST
	VERB_PUSH
	VERB_PULL
)

/*
	The flags that are used globally by the tool.
*/
type GlobalFlags struct {
	OAuthConfigPath string
	TokenPath string
}

func ParseOperationVerb(verb string) (OperationVerb, error) {

	switch(verb) {
		case "list": return VERB_LIST, nil
		case "push": return VERB_PUSH, nil
		case "pull": return VERB_PULL, nil
	}

	return VERB_UNKNOWN, errors.New(fmt.Sprintf("Unrecognized verb: '%s'", verb))
}

func ParseGlobalFlags() (GlobalFlags) {

	var oauthPath = flag.String("c", "", "The path to the OAuth config json file to use")
	var tokenPath = flag.String("t", getUserTokenPath(), "The path to the token cache file to use")

	flag.Parse()

	return GlobalFlags {
		OAuthConfigPath: *oauthPath,
		TokenPath: *tokenPath,
	}
}

func getUserTokenPath() (string) {

	var currentUser *user.User
	var err error

	currentUser, err = user.Current()
	if(err != nil) {
		return ""
	}

	return fmt.Sprintf("%s/.oauth/token.json", currentUser.HomeDir)
}
