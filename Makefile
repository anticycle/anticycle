.DEFAULT_GOAL := help
SHELL:=/bin/bash
GOPATH ?= `go env GOPATH`

.PHONY: docs test build

##@ Development

devdeps: ## install development dependencies
	./tools/devdeps.sh install

install: uninstall ## build and install project in $GOPATH/bin/
	./build/install.sh ./cmd/anticycle

##@ Testing

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

benchmark: ## run test benchmark
	go test -run=xxx -bench=. ./pkg/... ./internal/... > new.bench
	benchcmp old.bench new.bench

benchmark-save: ## save new benchmark as old
	mv new.bench old.bench

##@ Building

build: clean ## create artifacts
	./build/artifacts.sh ./dist ./cmd/anticycle

tarball: ## create tar.gz files
	rm -rf ./dist/release
	./build/tarball.sh ./dist

##@ Cleanup

clean: clean-build clean-coverage clean-devdeps uninstall ## clean cache, binaries, dev deps and coverage
	go clean

clean-build: ## celan built binaries
	rm -rf ./dist/

clean-devdeps: ## uninstall development dependencies
	./tools/devdeps.sh uninstall

clean-coverage: ## clean coverage report
	rm -f *.out
	rm -f coverage.html

uninstall: ## uninstall project from $GOPATH/bin/
	rm -f $(GOPATH)/bin/anticycle

##@ Helpers

format: ## reformat sourcecode
	go fmt ./...

lint: ## static checks
	./tools/linter.sh

help: ## display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
