package main

import (
	"fmt"

	rpccommand "github.com/bcurnow/go-plugin-command-line/rpc/command"
	// This is an ugly hack but we need to ensure that the gob registration occurs
	_ "github.com/bcurnow/go-plugin-command-line/rpc/util"

	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
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

func (c *CommandVersion) SetServices(serviceInfo map[string]service.ReconnectInfo) error {
	return nil
}

// Starts the RCP server
func main() {
	plugin.Start(&rpccommand.Plugin{Impl: &CommandVersion{}}, "version", command.HandshakeConfig)
}
