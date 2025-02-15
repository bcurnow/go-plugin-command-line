package main

import (
	rpccommand "github.com/bcurnow/go-plugin-command-line/rpc/command"
	rpcservice "github.com/bcurnow/go-plugin-command-line/rpc/service"
	"github.com/bcurnow/go-plugin-command-line/shared/cmd"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type RPCCommandRegister struct {
	command.CommandRegister
}

func (r *RPCCommandRegister) RegisterCommands(pluginDir string, cobraCmd *cobra.Command, serviceInfo map[string]service.ReconnectInfo) error {
	return cmd.RegisterCommands(pluginDir, cobraCmd, &rpccommand.CommandPlugin{}, serviceInfo)
}

type RPCServiceRegister struct {
	service.ServiceRegister
}

func (r *RPCServiceRegister) RegisterServices(serviceDir string) (map[string]service.ReconnectInfo, error) {
	return service.RegisterServices(serviceDir, &rpcservice.ServicePlugin{}, func(pluginClient *plugin.Client, pluginName string, service service.Service) service.ReconnectInfo {
		return &rpcservice.RPCClientInfo{
			ReattachConfig: rpcservice.ReattachConfig{
				Protocol:        pluginClient.ReattachConfig().Protocol,
				ProtocolVersion: pluginClient.ReattachConfig().ProtocolVersion,
				Addr: rpcservice.Addr{
					Net:  pluginClient.ReattachConfig().Addr.Network(),
					Name: pluginClient.ReattachConfig().Addr.String(),
				},
				Pid:          pluginClient.ReattachConfig().Pid,
				ReattachFunc: pluginClient.ReattachConfig().ReattachFunc,
				Test:         pluginClient.ReattachConfig().Test,
			},
			PluginName: pluginName,
			PluginType: service.Type(),
		}
	})
}

// Simply call the Execute function in cmd (the Cobra package)
func main() {
	cmd.Execute(&RPCServiceRegister{}, &RPCCommandRegister{})
}
