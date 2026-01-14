#!/usr/bin/env bash
set -e

export GOOS=windows
export GOARCH=386
export CGO_ENABLED=1
export CC=i686-w64-mingw32-gcc
go build dmc.go
