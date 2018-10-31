#!/usr/bin/env bash
# USAGE
# ./deps.sh <command>
#
# COMMANDS:
# install       Install all dependencies
# uninstall     Uninstall dependencies and its binaries

set -o errexit
set -o xtrace

GOPATH=$(go env GOPATH)

function install() {
    # golint
    go get -u golang.org/x/lint/golint
}

function uninstall() {
    # golint
    rm -rf ${GOPATH}/src/golang.org/x/lint
	rm -f ${GOPATH}/bin/golint
}

while :; do
	case "$1" in
		'install') install; shift;;
		'uninstall') uninstall; shift;;
		*) break;;
	esac
done
