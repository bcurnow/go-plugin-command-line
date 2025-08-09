package main

import (
	"fmt"

	rpcservice "github.com/bcurnow/go-plugin-command-line/rpc/service"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
)

// This is the instance of the logger service
type LoggerService struct {
	service.Service
}

func (c *LoggerService) Name() string {
	return "logger"
}

func (c *LoggerService) Log(val string) {
	fmt.Println(val)
}

// Starts the RCP server
func main() {
	plugin.Start(&rpcservice.Plugin{Impl: &LoggerService{}}, "logger", service.HandshakeConfig, nil)
}
