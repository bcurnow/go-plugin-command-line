#!/usr/bin/make

SHELL := /bin/bash
currentDir := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

build-example-command-version:
	go build -o ${currentDir}/commands/version ./examples/commands/version/version.go

build-example-command-version2:
	go build -o ${currentDir}/commands/version2 ./examples/commands/version2/version2.go

build-example-command-loggerversion:
	go build -o ${currentDir}/commands/loggerversion ./examples/commands/loggerversion/loggerversion.go

build-example-service-logger:
	go build -o ${currentDir}/services/logger ./examples/services/logger/logger.go

build-example-commands: clean-example-commands build-example-command-version build-example-command-version2 build-example-command-loggerversion

build-example-services: clean-example-services build-example-service-logger

build-examples: build-example-commands build-example-services

clean-example-commands:
	rm -f ${currentDir}/commands/*

clean-example-services:
	rm -f ${currentDir}/services/*

clean-examples: clean-example-commands clean-example-services