package main

import (
	"fmt"
	"os"
	"strconv"

	rpcservice "github.com/bcurnow/go-plugin-command-line/rpc/service"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	"github.com/bcurnow/go-plugin-command-line/shared/util"
)

// This is the instance of the logger service
type LoggerService struct {
	service.Service
}

func (c *LoggerService) Name() string {
	return fmt.Sprintf("LoggerService pid: %s", strconv.Itoa(os.Getpid()))
}

func (c *LoggerService) Type() string {
	return "LoggerService"
}

func (c *LoggerService) Execute(arg string) error {
	fmt.Println("3.0.0")
	return nil
}

// Starts the RCP server
func main() {
	util.StartPlugin(&rpcservice.ServicePlugin{Impl: &LoggerService{}}, "logger", service.HandshakeConfig)
}
