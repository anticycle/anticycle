# Anticycle

[![Godoc](https://godoc.org/github.com/anticycle/anticycle?status.svg)](https://godoc.org/github.com/anticycle/anticycle)
[![CircleCI](https://circleci.com/gh/anticycle/anticycle.svg?style=shield)](https://circleci.com/gh/anticycle/anticycle)
[![Goreportcard](https://goreportcard.com/badge/github.com/anticycle/anticycle)](https://goreportcard.com/report/github.com/anticycle/anticycle)
[![Release](https://img.shields.io/github/release/anticycle/anticycle.svg)](https://github.com/anticycle/anticycle/releases/latest)
[![License](https://img.shields.io/github/license/anticycle/anticycle.svg)](https://github.com/anticycle/anticycle/blob/master/LICENSE)
[![Platforms](https://img.shields.io/badge/platforms-linux%20%7C%20osx%20%7C%20windows-red.svg)](https://github.com/anticycle/anticycle/releases/latest)

Anticycle is a tool for static code analysis which search for
dependency cycles. It scans recursively all source files and
parses theirs dependencies. Anticycle does not compile the code,
so it is ideal for searching for complex, difficult to debug cycles.

Anticycle is published under GNU GENERAL PUBLIC LICENSE Version 3, from 29 June 2007

## Installation

**You may need sudo privileges to unpack into system specific directories.**

1. Download binary from [latest release](https://github.com/anticycle/anticycle/releases/latest).
2. Unpack binary to your PATH for example: `tar xvf anticycle.linux-amd64.tar.gz -C /usr/local/bin`
3. Verify installation: `anticycle -version`

## Usage

```
anticycle [options] [directory]
```

### Options

```
-all                 Output all packages, with and without cycles.

-format="text"       Output format. Available: text, json.

-exclude=""          A space-separated list of directories that should 
                     not be scanned. The list will be added to the 
                     default list of directories.
-excludeDefault=""   A space-separated list of directories that should 
                     not be scanned. The list will override the default.
-showExcluded        Shows list of excluded directories.

-help                Shows this help text.
-version             Shows version and build hash.
```

### Directory

An optional path to the analyzed project. If the directory is not
defined, the current working directory will be used.

### Example

Analyze recursively from current working directory but skip `internal/` anywhere in dir tree.

```bash
$ anticycle -exclude="internal"
```

Analyze recursively given directory with JSON output format

```bash
$ anticycle -all -format=json $GOPATH/src/github.com/anticycle/anticycle
```

### Example output

**Real case scenario:**

Command without flags and directories.

```
$ cd $GOPATH/src/github.com/Juniper/contrail
$ anticycle
Found 4 cycles

db -> models -> db
ipam -> models -> db -> models
models -> db -> models
services -> models -> db -> models

Details

[db -> models] "github.com/Juniper/contrail/pkg/models"
   pkg/db/address_manager.go
   pkg/db/address_manager_test.go
   pkg/db/db.go
   pkg/db/db_test.go
   pkg/db/useragent_kv.go

[models -> db] "github.com/Juniper/contrail/pkg/db"
   pkg/models/validation.go

[services -> models] "github.com/Juniper/contrail/pkg/models"
   pkg/services/contrail_service_test.go
   pkg/services/event_test.go
   pkg/services/list_response_test.go

[ipam -> models] "github.com/Juniper/contrail/pkg/models"
   pkg/types/ipam/address_manager.go
```

**How to read:**

```
[package -> wants] "fully/qualified/import/name"
   path/to/affected/file.go
   path/to/another/file.go
```

This gives us a hint that few files in `db` package wants to import `models`, but
`validation.go` file in `models` want to import `db` package.

The cycle looks like: `db -> models -> db`.

## Development

**Require GO v1.11.x**

Make sure you have GO in version 1.11. If not, [follow official instructions](https://golang.org/doc/install).

### Makefile

Use `Make` to run tests, benchmarks, build process and more.

```bash
sudo apt install make
$ make help
```

### Download project

```bash
$ go get github.com/anticycle/anticycle
$ make devdeps install test-all
```

After each change use `make install` to update dev binary. Then run sanity tests.
**Sanity tests are called on built binary, not the source code.**

### Run tests

Make sure you had installed anticycle in your OS before run sanity or acceptance tests.

```bash
$ make test-unit test-sanity test-acceptance
```

or just

```bash
$ make test
```

### Build artifacts

Artifacts will be moved to /bin directory after building.

```bash
$ make build tarball
```

## Contribution

Your help is very desirable. We are open to all kinds of contribution:

- Reporting found bugs
- Requesting new features
- Typos or bug fixing
- Filling gaps in documentation
- Creating new features
- Building tools for other contributors

If you don't know how to GO, you can report a bug or request a new feature by
creating new issue at https://github.com/anticycle/anticycle/issues

If you want to be involved in development:

1. Fork this repository.
2. Create new issue with description of what you are plan to do.
3. Create new feature or bug branch in your forked repository.
4. Make changes, and create pull request to our master branch.
5. Attach pull request to issue and wait for community response.

### Remember to

- always follow conventions found in this repository
- be kind to other contributors and never attack them for any reason
- if you are reviewing pull requests, focus on source code, never on developer who wrote it
- always write a unittests
- if it is required create new sanity or end-to-end test
