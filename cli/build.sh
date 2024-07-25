#!/bin/bash

# export VERSION=$(git describe --tags --abbrev=0)

# Build the application
# go build -buildvcs -ldflags "-X main.Version=${VERSION}" .
go build -buildvcs .