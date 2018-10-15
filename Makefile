.PHONY: docs test build
.DEFAULT_GOAL := help

help:
	@perl -nle'print $& if m{^[a-zA-Z_-]+:.*?## .*$$}' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

install: uninstall ## build and install project in OS
	go install cmd/anticycle

uninstall: ## uninstall project from OS
	rm -f $GOPATH/bin/anticycle

clean: clean-build clean-coverage ## remove artifacts

build: clean-build ## create artifacts
	./build/artifacts.sh

clean-build: ## remove linker artifacts
	rm -rf ./dist/*

test: ## run tests
	go test ./...

coverage: clean-coverage ## make test coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean-coverage: ## clean coverage report
	rm -f *.out
	rm -f coverage.html

format: ## reformat sourcecode
	go fmt ./...
