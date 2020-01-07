#!/bin/sh

set -ex

COMMID=$(git rev-parse --short HEAD)
NAME='cocoapods-cache-proxy-server'
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main._version_=v1 -X main._commit_=$COMMID" -o ./bin/$NAME-linux-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main._version_=v1 -X main._commit_=$COMMID" -o ./bin/$NAME-darwin-amd64