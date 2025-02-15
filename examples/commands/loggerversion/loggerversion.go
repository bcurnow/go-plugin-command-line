package main

import (
	rpccommand "github.com/bcurnow/go-plugin-command-line/rpc/command"
	rpcservice "github.com/bcurnow/go-plugin-command-line/rpc/service"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	"github.com/bcurnow/go-plugin-command-line/shared/util"
	"github.com/hashicorp/go-plugin"
)

var services = make(map[string]service.Service)

// This is the instance of the loggerversion command
type CommandLoggerVersion struct {
	command.Command
}

func (c *CommandLoggerVersion) Help() string {
	return "Returns the version using the logger service"
}

func (c *CommandLoggerVersion) Execute(args []string) error {
	services["logger"].Execute("3.0.0")
	return nil
}

func (c *CommandLoggerVersion) SetServices(serviceInfo map[string]service.ReconnectInfo) error {
	serviceMap, err := rpcservice.ToServices(serviceInfo, &rpcservice.ServicePlugin{})
	if err != nil {
		return err
	}
	services = serviceMap
	return nil
}

// Starts the RCP server
func main() {
	defer plugin.CleanupClients()
	util.StartPlugin(&rpccommand.CommandPlugin{Impl: &CommandLoggerVersion{}}, "loggerversion", command.HandshakeConfig)
}
