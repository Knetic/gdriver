package main

import (
	"fmt"
	"flag"
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
	Verb OperationVerb
}

func ParseGlobalFlags() (GlobalFlags) {

	var help = flag.Bool("h", false, "Shows the possible flags")
	var oauthPath = flag.String("c", "", "The path to the OAuth config json file to use")
	var tokenPath = flag.String("t", getUserTokenPath(), "The path to the token cache file to use")
	var verbList = flag.Bool("l", false, "Specifies that drive files should be listed")
	var verbPush = flag.Bool("u", false, "Specifies that a file should be uploaded")
	var verb OperationVerb

	flag.Parse()

	if(*help == true) {
		flag.PrintDefaults()
	}

	verb = VERB_UNKNOWN

	if(*verbList) {
		verb = VERB_LIST
	}
	if(*verbPush) {
		verb = VERB_PUSH
	}

	return GlobalFlags {
		OAuthConfigPath: *oauthPath,
		TokenPath: *tokenPath,
		Verb: verb,
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
