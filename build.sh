#!/bin/sh

rm -r ./pkg
mkdir ./pkg

GOOS=windows GOARCH=amd64 go build -o pkg/shipmenttimes-amd64.exe
GOOS=windows GOARCH=386 go build -o pkg/shipmenttimes-386.exe
GOOS=darwin GOARCH=amd64 go build -o pkg/shipmenttimes-darwin-amd64
GOOS=darwin GOARCH=386 go build -o pkg/shipmenttimes-darwin-386
