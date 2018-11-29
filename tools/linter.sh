#!/usr/bin/env bash
# Copyright 2018 The Anticycle Authors. All rights reserved.
# Use of this source code is governed by a GPL-style
# license that can be found in the LICENSE file.
#
# Vet examines Go source code and reports suspicious constructs.
# https://golang.org/cmd/vet/
function check_vet() {
    echo "Check: go vet"
    for dir in "${@:2}"
    do
        go vet -source "$1$dir/..."
        if [[ $? > 0 ]]; then
            FAILED=1
        fi
    done
}

# Golint is a linter for Go source code.
# https://github.com/golang/lint
function check_golint() {
    echo "Check: golint"
    golint -set_exit_status "$1/doc.go"
    for dir in  "${@:2}"
    do
        golint -set_exit_status "$1$dir/..."
        if [[ $? > 0 ]]; then
            FAILED=1
        fi
    done
}

function check_gofmt() {
    echo "Check: gofmt"
    make -C $1 format
    for dir in  "${@:2}"
    do
        git diff --exit-code -- "$1$dir"
        if [[ $? > 0 ]]; then
            FAILED=1
        fi
    done
}

# Failure flag
FAILED=0

# Prepare relative path to match project root directory
current_dir=$(pwd)
if [[ -z "${current_dir##*tools*}" ]]; then
    # if "tools" in pwd then we want to go one directory up
    root_dir="../"
else
    # else stay where you are
    root_dir="./"
fi

source_dirs=(
    cmd
    pkg
    internal
    test
)

# Run linters
check_vet ${root_dir} ${source_dirs[@]} >&2
check_golint ${root_dir} ${source_dirs[@]} >&2
check_gofmt ${root_dir} ${source_dirs[@]} >&2

if [[ ${FAILED} > 0 ]]; then
    exit 1
fi
