#!/bin/sh
# Builds the project. Requires Go 1.22

cd "$(dirname "$0")"

export GOOS=$(uname -s | tr '[:upper:]' '[:lower:]')
export GOARCH=$(uname -m | sed 's/aarch64/arm64/' | sed 's/x86_64/amd64/')

go build -o build/
