#!/usr/bin/env bash
# Copyright 2018 The Anticycle Authors. All rights reserved.
# Use of this source code is governed by a GPL-style
# license that can be found in the LICENSE file.
#
# Build binaries for each OS and architecture
declare -a OSARCHS=("linux/amd64" "linux/arm" "darwin/amd64" "windows/amd64")
version=($("$(dirname "$0")/version.sh"))
version_tag=${version[0]}
version_hash=${version[1]}

ld=(
    "-X main.version=$version_tag"
    "-X main.build=$version_hash"
    "-s"
    "-w"
)
out=$1
in=$2
mkdir -p "./dist/"

for osarch in "${OSARCHS[@]}"
do
    echo "Build artifacts: ${osarch}"

    oa=(${osarch//// })  # replace slash to space and split to array
    os_name=${oa[0]}
    os_arch=${oa[1]}

    filename="anticycle-$version_tag.$os_name-$os_arch"
    if [[ ${os_name} == "windows" ]]; then
        filename="${filename}.exe"
    fi

    env GOOS=${os_name} GOARCH=${os_arch} go build -ldflags="${ld[*]}" \
                                                   -o ${out}/${filename} ${in}
done

chmod a+x ${out}/*
