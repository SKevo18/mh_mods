#!/bin/sh

export GOOS=$(uname -s | tr '[:upper:]' '[:lower:]')
export GOARCH=$(uname -m | sed 's/aarch64/arm64/' | sed 's/x86_64/amd64/')

echo "Building mhmods CLI for $GOOS/$GOARCH"
cd "$(dirname "$0")/cli/"
go build -o ../build/

echo "Building mhmods GUI for $GOOS/$GOARCH"
cd "$(dirname "$0")/gui/"
go build -o ../build/
