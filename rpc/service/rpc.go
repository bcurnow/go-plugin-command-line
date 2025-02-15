package service

import (
	"net"
	"net/rpc"
	"os"

	"github.com/bcurnow/go-plugin-command-line/shared/service"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-plugin/runner"
)

type Addr struct {
	Net  string
	Name string
}

// This struct mirrors plugin.ReattachConfig but only includes types that are registered with gob to avoid issues
type ReattachConfig struct {
	Protocol        plugin.Protocol
	ProtocolVersion int
	Addr            Addr
	Pid             int
	ReattachFunc    runner.ReattachFunc
	Test            bool
}

func (rc *ReattachConfig) ToReattachConfig() *plugin.ReattachConfig {
	return &plugin.ReattachConfig{
		Protocol:        rc.Protocol,
		ProtocolVersion: rc.ProtocolVersion,
		Addr:            &net.UnixAddr{Name: rc.Addr.Name, Net: rc.Addr.Net},
		Pid:             rc.Pid,
		ReattachFunc:    rc.ReattachFunc,
		Test:            rc.Test,
	}
}

type RPCClientInfo struct {
	ReattachConfig ReattachConfig
	PluginName     string
	PluginType     string
}

func (i *RPCClientInfo) Name() string {
	return i.PluginName
}

func (i *RPCClientInfo) Type() string {
	return i.PluginType
}

// The plugin.Plugin implementation for a Service which returns the RCP Client or Server
type ServicePlugin struct{ Impl service.Service }

func (p *ServicePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ServiceRPCServer{Impl: p.Impl}, nil
}

func (ServicePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ServiceRPCClient{client: c}, nil
}

// The RPC client implementation of a Service
type ServiceRPCClient struct{ client *rpc.Client }

func (c *ServiceRPCClient) Name() string {
	var resp string
	err := c.client.Call("Plugin.Name", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func (c *ServiceRPCClient) Type() string {
	var resp string
	err := c.client.Call("Plugin.Type", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func (c *ServiceRPCClient) Execute(arg string) error {
	var resp interface{}
	return c.client.Call("Plugin.Execute", map[string]interface{}{
		"arg": arg,
	}, &resp)
}

// The RPC server implementation of a Service
// NOTE: While this struct will have implementations of the Service methods, they will have different signatures
// required by the RPC package.
type ServiceRPCServer struct{ Impl service.Service }

// The first argument, args interface{}, is RCP speak for no parameters
// resp is the return value and the type should match the Service method (e.g. string)
func (s *ServiceRPCServer) Name(args interface{}, resp *string) error {
	*resp = s.Impl.Name()
	return nil
}

// The first argument, args interface{}, is RCP speak for no parameters
// resp is the return value and the type should match the Service method (e.g. string)
func (s *ServiceRPCServer) Type(args interface{}, resp *string) error {
	*resp = s.Impl.Type()
	return nil
}

func (s *ServiceRPCServer) Execute(args map[string]interface{}, resp *interface{}) error {
	return s.Impl.Execute(args["arg"].(string))
}

func ToServices(serviceInfo map[string]service.ReconnectInfo, theService plugin.Plugin) (map[string]service.Service, error) {
	services := make(map[string]service.Service)

	for name, serviceInfo := range serviceInfo {
		rpcClientInfo := serviceInfo.(*RPCClientInfo)
		// Create a new PluginDef using the ReattachConfig instead of a Cmd
		pluginDef := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: service.HandshakeConfig,
			Plugins: map[string]plugin.Plugin{
				name: theService,
			},
			Reattach:   rpcClientInfo.ReattachConfig.ToReattachConfig(),
			Logger:     service.Logger,
			Managed:    true,      // Allow the plugin runtime to manage this plugin
			SyncStdout: os.Stdout, // Print any extra output to Stdout from the plugin to the host processes Stdout
			// AutoMTLS:    true,      // Ensure that we're using MTLS for communication between the plugin and the host
			SkipHostEnv: true, // Don't pass the host environment to the plugin to avoid security issues
		})

		service, err := service.ToService(pluginDef, name)
		if err != nil {
			return nil, err
		}

		services[name] = service
	}
	return services, nil
}
