#!/bin/bash

export GOPATH=$(pwd)

go get code.google.com/p/google-api-go-client/drive/v2
go get code.google.com/p/goauth2/oauth

mkdir -p output
go build -o output/griveBackup.exe .
