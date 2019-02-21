#!/usr/bin/env bash

# unzip4win build script for mac

cd $(dirname $0)

echo 'Cleaning the output directory.'
mkdir -p out/unzip4win
rm -rf out/unzip4win/**

echo 'Creating assets'
touch config.toml
go generate -v
echo 'Building binary for mac'
go build -v -o out/unzip4win/unzip4win main.go
echo 'Building binary for windows'
GOOS=windows GOARCH=amd64 go build -v -o out/unzip4win/unzip4win.exe main.go
echo 'Finish'
