# Anticycle

[![Godoc](https://godoc.org/github.com/anticycle/anticycle?status.svg)](https://godoc.org/github.com/anticycle/anticycle)
[![CircleCI](https://circleci.com/gh/anticycle/anticycle.svg?style=svg)](https://circleci.com/gh/anticycle/anticycle)

Anticycle is a tool for static code analysis which search for
dependency cycles. It scans recursively all source files and
parses theirs dependencies. Anticycle does not compile the code,
so it is ideal for searching for complex, difficult to debug cycles.

Anticycle is published under GNU GENERAL PUBLIC LICENSE Version 3, from 29 June 2007

## Usage

```
anticycle [options] [directory]
```

### Options

```
-all               Output all packages, with and without cycles.
-related           Output cycles and packages which lead to,
                   but are not involved in cycle.
                   Ignored if -all is present.

-format            Output format. Available: text,json. Default: text.

-exclude=""        A comma separated list of directories that should
                   not be scanned. The list will be added to the
                   default list of directories.
-excludeOnly=""    A comma separated list of directories that should
                   not be scanned. The list will override the default.
-showExclude       Shows default list of excluded directories.

-help              Shows this help text.
-version           Shows version tag.
```

### Directory

An optional path to the analyzed project. If the directory is not
defined, the current working directory will be used.

### Example

Analyze recursively from current working directory but skip `internal/` anywhere in dir tree.

```console
$ anticycle -exclude="internal"
```

Analyze recursively given directory

```console
$ anticycle $GOPATH/src/github.com/anticycle/anticycle -all -format=json
```

### Example output

**How to read:**

```console
[package -> wants] "fully/qualified/import/name"
   path/to/affected/file.go
   path/to/another/file.go
```

**Real case output:**

```console
[db -> models] "github.com/Juniper/contrail/pkg/models"
   pkg/db/address_manager.go
   pkg/db/address_manager_test.go
   pkg/db/db.go
   pkg/db/db_test.go
   pkg/db/useragent_kv.go

[models -> db] "github.com/Juniper/contrail/pkg/db"
   pkg/models/validation.go
```

## Development

**Require GO v1.11.x**

Make sure you have GO in version 1.11. If not, [follow official instructions](https://golang.org/doc/install).

### Download project

```console
$ go get github.com/anticycle/anticycle
$ make devdeps install
```

After each change use `make install` to update dev binary. Then run sanity tests.
**Sanity tests are called on built binary, not the source code.**

### Run tests

```console
$ make test-all
```

### Build artifacts

Artifacts will be moved to /bin directory after building.

```console
$ make build
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
