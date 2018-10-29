#!/usr/bin/env bash
# Install project in $GOPATH/bin
version=($("$(dirname "$0")/version.sh"))
ld=(
    "-X main.version=${version[0]}"
    "-X main.build=${version[1]}"
)
go install -ldflags="${ld[*]}" $1
