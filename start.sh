#!/usr/bin/env bash
cd "$( cd "$( dirname "$0" )" >/dev/null 2>&1 && pwd )" || exit
go mod tidy
go build
./m3u8Find -oa -oaa -om -pt girls
./m3u8Find -oa -oaa -om -pt couples
./m3u8Find -oa -oaa -om -pt new
./m3u8Find -oa -oaa -om -pt autoTagHd

