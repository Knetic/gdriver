#!/bin/bash

export GOPATH=$(pwd)

go get code.google.com/p/google-api-go-client/drive/v2
go get code.google.com/p/goauth2/oauth

mkdir -p output

# If you're building from windows, don't forget to install the bootstrappers:
#go tool dist install -v pkg/runtime
#go install -v -a std

export GOARCH=amd64

export GOOS=windows
go build -o output/gdriver.exe .

export GOOS=darwin
go build -o output/gdriver.bin .

export GOOS=linux
go build -o output/gdriver .
