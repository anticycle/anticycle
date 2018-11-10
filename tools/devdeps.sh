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

    # golang tools
    go get -u golang.org/x/tools/...
}

function uninstall() {
    # golint
    rm -rf ${GOPATH}/src/golang.org/x/lint
	rm -f ${GOPATH}/bin/golint

    # golang tools
    rm -rf ${GOPATH}/src/golang.org/x/tools
    cd ${GOPATH}/bin
    rm -f benchcmp bundle callgraph compilebench cover digraph eg fiximports getgo \
          go-contrib-init godex godoc goimports golsp gomvpkg gorename gotype goyacc \
          guru heapview html2article present ssadump stress stringer tip toolstash vet vet-lite
}

while :; do
	case "$1" in
		'install') install; shift;;
		'uninstall') uninstall; shift;;
		*) break;;
	esac
done
