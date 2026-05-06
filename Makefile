#!/usr/bin/make

SHELL := /bin/bash
currentDir := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

set-rpc:
	$(eval SERVICE_TYPE := rpc)
	@echo "SERVICE_TYPE is set to $(SERVICE_TYPE)"

set-grpc:
	$(eval SERVICE_TYPE := grpc)
	@echo "SERVICE_TYPE is set to $(SERVICE_TYPE)"

build-example-command-version:
	go build -o ${currentDir}/commands/version ./examples/$(SERVICE_TYPE)/commands/version/version.go

build-example-command-version2:
	go build -o ${currentDir}/commands/version2 ./examples/$(SERVICE_TYPE)/commands/version2/version2.go

build-example-command-loggerversion:
	go build -o ${currentDir}/commands/loggerversion ./examples/$(SERVICE_TYPE)/commands/loggerversion/loggerversion.go

build-example-service-logger:
	go build -o ${currentDir}/services/logger ./examples/$(SERVICE_TYPE)/services/logger/logger.go

build-rpc-example-commands: set-rpc clean-example-commands build-example-command-version build-example-command-version2 build-example-command-loggerversion

build-rpc-example-services: set-rpc clean-example-services build-example-service-logger

build-rpc-examples: build-rpc-example-commands build-rpc-example-services

build-grpc-example-commands: set-grpc clean-example-commands build-example-command-version build-example-command-version2 build-example-command-loggerversion

build-grpc-example-services: set-grpc clean-example-services build-example-service-logger

build-grpc-examples: build-grpc-example-commands build-grpc-example-services


clean-example-commands:
	rm -f ${currentDir}/commands/*

clean-example-services:
	rm -f ${currentDir}/services/*

clean-examples: clean-example-commands clean-example-services

run-rpc:
	go run ${currentDir}/rpc/main.go

run-grpc:
	go run ${currentDir}/grpc/main.go