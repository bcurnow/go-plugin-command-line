package main

import (
	"fmt"

	grpccommand "github.com/bcurnow/go-plugin-command-line/rpc/command"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	goplugin "github.com/hashicorp/go-plugin"
)

// This is the instance of the version2 command
type CommandVersion2 struct {
	command.Command
}

func (c *CommandVersion2) Help() string {
	return "Returns the version"
}

func (c *CommandVersion2) Execute(args []string) error {
	fmt.Println("2.0.0")
	return nil
}

func (c *CommandVersion2) SetServices(serviceInfo map[string]plugin.Reattach) error {
	return nil
}

// Starts the RCP server
func main() {
	plugin.Start(&grpccommand.Plugin{Impl: &CommandVersion2{}}, "version2", command.HandshakeConfig, goplugin.DefaultGRPCServer)
}
