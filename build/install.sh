#!/usr/bin/env bash
go install -ldflags="-X main.version=$(./build/version.sh)" $1
