package main

import (
	"errors"

	grpccommand "github.com/bcurnow/go-plugin-command-line/grpc/command"
	grpcservice "github.com/bcurnow/go-plugin-command-line/grpc/service"

	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	goplugin "github.com/hashicorp/go-plugin"
)

var (
	services = make(map[string]service.Service)
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

func (c *CommandLoggerVersion) SetServices(reattaches map[string]plugin.Reattach) error {
	serviceMap, err := grpcservice.Services(reattaches, &grpcservice.Plugin{})
	if err != nil {
		return err
	}
	services = serviceMap
	return nil
}

// This isn't ideal in terms of service lookup but is sufficient for an example
func logger() service.Service {
	if nil == services {
		// This command was not setup properly by the main program
		panic(errors.New("services not initialized properly for command loggerversion"))
	}

	if nil == services["logger"] {
		panic(errors.New("service 'logger' could not be found in command loggerversion"))
	}

	return services["logger"]
}

// Starts the RCP server
func main() {
	plugin.Start(&grpccommand.Plugin{Impl: &CommandLoggerVersion{}}, "loggerversion", command.HandshakeConfig, goplugin.DefaultGRPCServer)
}
