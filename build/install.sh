#!/usr/bin/env bash
# Copyright 2018 The Anticycle Authors. All rights reserved.
# Use of this source code is governed by a GPL-style
# license that can be found in the LICENSE file.
#
# Install project in $GOPATH/bin
version=($("$(dirname "$0")/version.sh"))
ld=(
    "-X main.version=${version[0]}"
    "-X main.build=${version[1]}"
)
go install -ldflags="${ld[*]}" $1
