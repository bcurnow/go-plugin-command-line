package command

import (
	"net/rpc"

	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	goplugin "github.com/hashicorp/go-plugin"
)

// The plugin.Plugin implementation for a Command which returns the RCP Client or Server
type Plugin struct {
	Impl command.Command
}

func (p *Plugin) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (p *Plugin) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// The RPC client implementation of a Command
type RPCClient struct {
	client *rpc.Client
}

func (c *RPCClient) Help() string {
	defer plugin.HandlePanic()
	var resp string
	err := c.client.Call("Plugin.Help", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func (c *RPCClient) Execute(args []string) error {
	defer plugin.HandlePanic()
	var resp interface{}
	return c.client.Call("Plugin.Execute", map[string]interface{}{
		"args": args,
	}, &resp)
}

func (c *RPCClient) SetServices(serviceInfo map[string]plugin.Reattach) error {
	defer plugin.HandlePanic()
	var resp interface{}
	return c.client.Call("Plugin.SetServices", map[string]interface{}{
		"serviceInfo": serviceInfo,
	}, &resp)
}

// The RPC server implementation of a Command
// NOTE: While this struct will have implementations of the Command methods, they will have different signatures
// required by the RPC package.
type RPCServer struct{ Impl command.Command }

// The first argument, args interface{}, is RCP speak for no parameters
// resp is the return value and the type should match the Command method (e.g. string)
func (s *RPCServer) Help(args interface{}, resp *string) error {
	defer plugin.HandlePanic()
	*resp = s.Impl.Help()
	return nil
}

// The first argument, map[string]interface{}, is RCP speak for parameters, each parameter will be
// mapped to a key in the map.
// resp is a required argument but is not used in this case
func (s *RPCServer) Execute(args map[string]interface{}, resp *interface{}) error {
	defer plugin.HandlePanic()
	return s.Impl.Execute(args["args"].([]string))
}

func (s *RPCServer) SetServices(args map[string]interface{}, resp *interface{}) error {
	defer plugin.HandlePanic()
	return s.Impl.SetServices(args["serviceInfo"].(map[string]plugin.Reattach))
}
