package main

import (
	"errors"

	rpccommand "github.com/bcurnow/go-plugin-command-line/rpc/command"
	rpcservice "github.com/bcurnow/go-plugin-command-line/rpc/service"

	// This is an ugly hack but we need to ensure that the gob registration occurs
	_ "github.com/bcurnow/go-plugin-command-line/rpc/util"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
)

var (
	services = make(map[plugin.Type]service.Service)
)

// This is the instance of the loggerversion command
type CommandLoggerVersion struct {
	command.Command
}

func (c *CommandLoggerVersion) Help() string {
	return "Returns the version using the logger service"
}

func (c *CommandLoggerVersion) Execute(args []string) error {
	logger().Log("3.0.0")
	return nil
}

func (c *CommandLoggerVersion) SetServices(serviceInfos map[string]service.ReconnectInfo) error {
	serviceMap, err := rpcservice.Services(serviceInfos, &rpcservice.Plugin{})
	if err != nil {
		return err
	}
	services = serviceMap
	return nil
}

// This is a hacky way to do this but works for an example
// Probably need either a more complex schema-based mechanism (e.g. Terraform plugins)
// Or individual typed Set methods for the different types of services
func logger() service.LoggerService {
	if nil == services {
		// This command was not setup properly by the main program
		panic(errors.New("services not initialized properly for command loggerversion"))
	}

	if nil == services[plugin.Logger] {
		panic(errors.New("service of type Logger could not be found in command loggerversion"))
	}

	return services[plugin.Logger].(service.LoggerService)
}

// Starts the RCP server
func main() {
	plugin.Start(&rpccommand.Plugin{Impl: &CommandLoggerVersion{}}, "loggerversion", command.HandshakeConfig)
}
