.PHONY: docs test build
.DEFAULT_GOAL := help

GOPATH ?= `go env GOPATH`

help:
	@perl -nle'print $& if m{^[a-zA-Z_-]+:.*?## .*$$}' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

deps: ## install development dependencies
	./tools/deps.sh install

install: ## build and install project in OS
	./build/install.sh ./cmd/anticycle

uninstall: clean ## uninstall project from OS
	rm -f $(GOPATH)/bin/anticycle

clean: clean-build ## remove artifacts
	go clean
	./tools/deps.sh uninstall

build: clean ## create artifacts
	./build/artifacts.sh ./dist ./cmd/anticycle
	chmod a+x -R ./dist

clean-build: ## remove linker artifacts
	rm -rf ./dist/*

test: ## run tests
	go test -v ./pkg/...
	go test -v ./internal/...

test-sanity: ## run sanity tests on builded binary
	go test -v ./test/sanity_test.go

test-all: test test-sanity ## run all tests

golden-update: ## update golden files
	go test ./test/sanity_test.go -update

coverage: clean-coverage ## make test coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean-coverage: ## clean coverage report
	rm -f *.out
	rm -f coverage.html

benchmark: ## run test benchmark
	go test -run=xxx -bench=. ./... > new.bench
	benchcmp old.bench new.bench

benchmark-save: ## save new benchmark as old
	mv new.bench old.bench

format: ## reformat sourcecode
	go fmt ./...

lint: ## static checks
	./tools/linter.sh
