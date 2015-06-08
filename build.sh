#!/bin/bash

export GOPATH=$(pwd)

echo "Getting dependencies"
go get code.google.com/p/google-api-go-client/drive/v2
go get code.google.com/p/goauth2/oauth

echo "Compiling"
mkdir -p output

# If you're building from windows, don't forget to install the bootstrappers:
#go tool dist install -v pkg/runtime
#go install -v -a std

export GOARCH=amd64

echo "Building windows"
export GOOS=windows
go build -o output/gdriver.exe .

echo "Building darwin"
export GOOS=darwin
go build -o output/gdriver.bin .

echo "Building linux"
export GOOS=linux
go build -o output/gdriver .

echo "Packaging"
pushd output > /dev/null

rm ./*.rpm
rm ./*.deb

fpm -n gdriver \
-s dir \
-t rpm \
-v 1.0 \
--log error \
./gdriver=/usr/local/bin/gdriver

fpm -n gdriver \
-s dir \
-t deb \
-v 1.0 \
--log error \
./gdriver=/usr/local/bin/gdriver

popd > /dev/null
