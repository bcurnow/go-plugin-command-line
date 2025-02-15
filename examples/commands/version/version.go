package main

import (
	"fmt"

	rpccommand "github.com/bcurnow/go-plugin-command-line/rpc/command"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	"github.com/bcurnow/go-plugin-command-line/shared/util"
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
	util.StartPlugin(&rpccommand.CommandPlugin{Impl: &CommandVersion{}}, "version", command.HandshakeConfig)
}
