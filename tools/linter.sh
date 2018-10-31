#!/usr/bin/env bash
set -o xtrace
set -o history
set -o histexpand

# Exit on any errors so that errors don't compound
trap err_trap ERR
function err_trap {
    local r=$?
    set +o xtrace
    echo "command: '$BASH_COMMAND' failed with status code $r"
    exit ${r}
}

# Vet examines Go source code and reports suspicious constructs.
# https://golang.org/cmd/vet/
function check_vet() {
    for dir in "${@:2}"
    do
        go vet -source "$1$dir/..."
    done
}

# Prepare directories to check
current_dir=$(dirname "$0")
if [[ ${current_dir} == "." ]]; then
    root_dir="../"
else
    root_dir="$current_dir/../"
fi

source_dirs=(
    cmd
    pkg
    internal
    test
)

# Run linters
check_vet ${root_dir} ${source_dirs[@]}
