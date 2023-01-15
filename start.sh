#!/usr/bin/env bash
cd "$( cd "$( dirname "$0" )" >/dev/null 2>&1 && pwd )" || exit
go mod tidy
go build
./m3u8Find -oa -om -pt girls
./m3u8Find -oa -om -pt men
./m3u8Find -oa -om -pt couples
./m3u8Find -oa -om -pt trans

