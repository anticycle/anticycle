.PHONY: docs test build
.DEFAULT_GOAL := help

GOPATH ?= `go env GOPATH`

help:
	@perl -nle'print $& if m{^[a-zA-Z_-]+:.*?## .*$$}' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

devdeps: ## install development dependencies
	./tools/devdeps.sh install

clean-devdeps: ## uninstall development dependencies
	./tools/devdeps.sh uninstall

install: uninstall ## build and install project in $GOPATH/bin/
	./build/install.sh ./cmd/anticycle

uninstall: ## uninstall project from $GOPATH/bin/
	rm -f $(GOPATH)/bin/anticycle

clean: clean-build ## remove artifacts
	go clean
	./tools/devdeps.sh uninstall

build: clean ## create artifacts
	./build/artifacts.sh ./dist ./cmd/anticycle

clean-build: ## remove linker artifacts
	rm -rf ./dist/

tarball: ## create tar.gz files
	rm -rf ./dist/release
	./build/tarball.sh ./dist

test: ## run tests
	go test -race -covermode=atomic ./pkg/... ./internal/...

test-sanity: ## run sanity tests on builded binary
	go test ./test/sanity_test.go

test-all: test test-sanity ## run all tests

golden-update: ## update golden files
	go test ./test/sanity_test.go -update

coverage: clean-coverage ## make test coverage report
	go test -race -covermode=atomic -coverprofile=coverage.out ./pkg/... ./internal/...
	cover -html=coverage.out -o coverage.html

clean-coverage: ## clean coverage report
	rm -f *.out
	rm -f coverage.html

benchmark: ## run test benchmark
	go test -run=xxx -bench=. ./pkg/... ./internal/... > new.bench
	benchcmp old.bench new.bench

benchmark-save: ## save new benchmark as old
	mv new.bench old.bench

format: ## reformat sourcecode
	go fmt ./...

lint: ## static checks
	./tools/linter.sh
