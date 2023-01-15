#!/usr/bin/env bash
cd "$( cd "$( dirname "$0" )" >/dev/null 2>&1 && pwd )" || exit
go mod tidy
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o release/m3u8Find_windows.exe
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o release/m3u8Find_darwin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o release/m3u8Find_linux
chmod +x release/m3u8Find_linux
release/m3u8Find_linux -oa -om -pt girls
release/m3u8Find_linux -oa -om -pt men
release/m3u8Find_linux -oa -om -pt couples
release/m3u8Find_linux -oa -om -pt trans

