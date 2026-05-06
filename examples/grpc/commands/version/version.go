package main

import (
	"fmt"

	grpccommand "github.com/bcurnow/go-plugin-command-line/grpc/command"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	goplugin "github.com/hashicorp/go-plugin"
)

// This is the instance of the version command
type CommandVersion struct {
	command.Command
}

func (c *CommandVersion) Help() string {
	return "Returns the version"
}

func (c *CommandVersion) Execute(args []string) error {
	fmt.Println("1.0.0")
	return nil
}

func (c *CommandVersion) SetServices(serviceInfo map[string]plugin.Reattach) error {
	return nil
}

// Starts the RCP server
func main() {
	plugin.Start(&grpccommand.Plugin{Impl: &CommandVersion{}}, "version", command.HandshakeConfig, goplugin.DefaultGRPCServer)
}
