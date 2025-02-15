package command

import (
	"net/rpc"

	_ "github.com/bcurnow/go-plugin-command-line/rpc/util"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	"github.com/hashicorp/go-plugin"
)

// The plugin.Plugin implementation for a Command which returns the RCP Client or Server
type CommandPlugin struct {
	plugin.Plugin
	Impl command.Command
}

func (p *CommandPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CommandRPCServer{Impl: p.Impl}, nil
}

func (p *CommandPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CommandRPCClient{client: c}, nil
}

// The RPC client implementation of a Command
type CommandRPCClient struct {
	command.Command
	client *rpc.Client
}

func (c *CommandRPCClient) Help() string {
	var resp string
	err := c.client.Call("Plugin.Help", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func (c *CommandRPCClient) Execute(args []string) error {
	var resp interface{}
	return c.client.Call("Plugin.Execute", map[string]interface{}{
		"args": args,
	}, &resp)
}

func (c *CommandRPCClient) SetServices(serviceInfo map[string]service.ReconnectInfo) error {
	var resp interface{}
	return c.client.Call("Plugin.SetServices", map[string]interface{}{
		"serviceInfo": serviceInfo,
	}, &resp)
}

// The RPC server implementation of a Command
// NOTE: While this struct will have implementations of the Command methods, they will have different signatures
// required by the RPC package.
type CommandRPCServer struct{ Impl command.Command }

// The first argument, args interface{}, is RCP speak for no parameters
// resp is the return value and the type should match the Command method (e.g. string)
func (s *CommandRPCServer) Help(args interface{}, resp *string) error {
	*resp = s.Impl.Help()
	return nil
}

// The first argument, map[string]interface{}, is RCP speak for parameters, each parameter will be
// mapped to a key in the map.
// resp is a required argument but is not used in this case
func (s *CommandRPCServer) Execute(args map[string]interface{}, resp *interface{}) error {
	return s.Impl.Execute(args["args"].([]string))
}

func (s *CommandRPCServer) SetServices(args map[string]interface{}, resp *interface{}) error {
	return s.Impl.SetServices(args["serviceInfo"].(map[string]service.ReconnectInfo))
}
