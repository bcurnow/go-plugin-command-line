# go-plugin-command-line

<!-- TOC depthfrom:2 bulletcharacter:* -->

* [Usage](#usage)
  * [Example plugins](#example-plugins)
  * [Example Services](#example-services)
  * [Running](#running)
    * [RPC](#rpc)
    * [gRPC](#grpc)
      * [Protocol Buffers](#protocol-buffers)

<!-- /TOC -->

This is an example of how to use (https://github.com/hashicorp/go-plugin/).

There are versions in RPC and gRPC, both with bi-directional communication using the service plugin type.

NOTE: In the case of RPC, this isn't really bi-directional with the main process as there's not a good way to run an RPC plugin in-process. (There is a way but I couldn't find a good way to stop that RPC plugin when running in-process)

## Usage

In order to use, you'll first need to compile the plugin(s) and service(s) and put them in the correct directories: `make build-examples`

### Example plugins

There are three plugins included in the examples/commands directory:

* version
* version2
* loggerversion

These are very simple examples which print a hard-coded version number to stdout.

To build them: `make build-example-plugins`

### Example Services

There is one service included in the examples/services directory: logger

This is a service which allows a plugin to log using a plugin maintained by the parent process

To build it: `make build-example-services`

### Running

#### RPC

To run the program: `go run rpc/main.go <command>`

`<command>` should be replaced with one of the example commands:

* version
* version2
* loggerversion

To access the built-in Cobra help: `go run rpc/main.go help`

#### gRPC

##### Protocol Buffers

The gRPC version uses Google Protocol buffers which need to be generated:
```
cd grpc
buf generate
```

This is a work in progress...
