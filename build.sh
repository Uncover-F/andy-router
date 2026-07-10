#!/bin/bash

APP=andy-router

mkdir -p dist

GOOS=linux GOARCH=amd64 go build -o dist/${APP}-linux-amd64 ./cmd
GOOS=linux GOARCH=arm64 go build -o dist/${APP}-linux-arm64 ./cmd

GOOS=darwin GOARCH=amd64 go build -o dist/${APP}-darwin-amd64 ./cmd
GOOS=darwin GOARCH=arm64 go build -o dist/${APP}-darwin-arm64 ./cmd

GOOS=windows GOARCH=amd64 go build -o dist/${APP}-windows-amd64.exe ./cmd
GOOS=windows GOARCH=arm64 go build -o dist/${APP}-windows-arm64.exe ./cmd

echo "Build complete"