package service

import (
	"net/rpc"

	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	goplugin "github.com/hashicorp/go-plugin"
)

// The goplugin.Plugin implementation for a Service which returns the RCP Client or Server
type Plugin struct{ Impl service.Service }

func (p *Plugin) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (Plugin) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// The RPC client implementation of a Service
type RPCClient struct{ client *rpc.Client }

func (c *RPCClient) Name() string {
	defer plugin.HandlePanic()
	var resp string
	err := c.client.Call("Plugin.Name", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func (c *RPCClient) Log(val string) {
	defer plugin.HandlePanic()
	var resp interface{}

	err := c.client.Call("Plugin.Log", map[string]interface{}{
		"val": val,
	}, &resp)

	if err != nil {
		service.Logger.Error("Error during Log", "Error", err)
		panic(err)
	}
}

// The RPC server implementation of a Service
// NOTE: While this struct will have implementations of the Service methods, they will have different signatures
// required by the RPC package.
type RPCServer struct{ Impl service.Service }

// The first argument, args interface{}, is RCP speak for no parameters
// resp is the return value and the type should match the Service method (e.g. string)
func (s *RPCServer) Name(args interface{}, resp *string) error {
	defer plugin.HandlePanic()
	*resp = s.Impl.Name()
	return nil
}

// The first agument is the set of args which were provided to the method, these will match
// with the RPCClient definition above
// The second argument is unused
func (s *RPCServer) Log(args map[string]interface{}, resp *interface{}) error {
	defer plugin.HandlePanic()
	s.Impl.Log(args["val"].(string))
	return nil
}
